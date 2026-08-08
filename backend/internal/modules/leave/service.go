package leave

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

// ApprovalEngine abstracts the central approval module so leave requests can
// be routed through it instead of leave's own ad-hoc SupervisorID/HrID
// fields. Implemented via an adapter wrapping approval.Service in main.go
// (same narrow-interface-plus-adapter pattern payroll already uses).
type ApprovalEngine interface {
	CreateApprovalInstance(ctx context.Context, module, documentID, flowID string) (string, error)
	GetApprovalInstanceStatus(ctx context.Context, instanceID string) (string, error)
}

// HolidayProvider abstracts the setting module's company holiday calendar so
// the leave calculation engine can exclude holidays from working-day counts
// without importing setting's repository/service directly (same
// narrow-interface-plus-adapter pattern as ApprovalEngine above).
type HolidayProvider interface {
	ListHolidayDatesInRange(ctx context.Context, fromDate, toDate string) ([]string, error)
}

// AttendanceSessionUpdater lets Leave push an approved leave onto Attendance's
// session for each covered date (docs/module-attendance-plan.md §26/§50), so
// a day fully covered by approved leave gets marked without needing a
// check-in event to trigger session generation. Implemented via an adapter
// wrapping attendance.Service in main.go (same narrow-interface-plus-adapter
// pattern as ApprovalEngine/HolidayProvider above).
type AttendanceSessionUpdater interface {
	ApplyApprovedLeave(ctx context.Context, employeeID uuid.UUID, workDate string, leaveRequestID uuid.UUID, dayFraction float64) error
}

// Notifier abstracts the notification module so Leave can notify an employee
// of their leave request's approval outcome (docs/module-notification-plan.md
// §5.2/§8). notification.Service satisfies this structurally.
type Notifier interface {
	Notify(ctx context.Context, recipientUserID uuid.UUID, notifType, title, body, referenceType string, referenceID uuid.UUID) error
}

type Service struct {
	repo              *Repository
	logger            *zap.Logger
	approvalEngine    ApprovalEngine
	holidayProvider   HolidayProvider
	attendanceUpdater AttendanceSessionUpdater
	notifier          Notifier
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// SetApprovalEngine wires the central approval module into this service.
func (s *Service) SetApprovalEngine(ae ApprovalEngine) {
	s.approvalEngine = ae
}

// SetHolidayProvider wires the company holiday calendar into this service.
func (s *Service) SetHolidayProvider(hp HolidayProvider) {
	s.holidayProvider = hp
}

// SetNotifier wires the notification module into this service so
// HandleApprovalStatusChange can notify the requesting employee of the
// outcome (docs/module-notification-plan.md §5/§8, Phase 4).
func (s *Service) SetNotifier(n Notifier) {
	s.notifier = n
}

// SetAttendanceSessionUpdater wires the attendance module into this service.
func (s *Service) SetAttendanceSessionUpdater(u AttendanceSessionUpdater) {
	s.attendanceUpdater = u
}

// holidayDatesInRange returns holiday dates in the range, degrading to "no
// holidays known" (rather than failing the request) if the provider isn't
// wired or errors — working-day calculation should not hard-fail a leave
// request just because the holiday lookup is unavailable.
func (s *Service) holidayDatesInRange(ctx context.Context, fromDate, toDate string) map[string]bool {
	holidayDates := map[string]bool{}
	if s.holidayProvider == nil {
		return holidayDates
	}
	dates, err := s.holidayProvider.ListHolidayDatesInRange(ctx, fromDate, toDate)
	if err != nil {
		s.logger.Warn("Failed to load company holidays for leave calculation, proceeding without them", zap.Error(err))
		return holidayDates
	}
	for _, d := range dates {
		holidayDates[d] = true
	}
	return holidayDates
}

// =========================================================================
// Leave Types
// =========================================================================

func (s *Service) CreateLeaveType(ctx context.Context, req CreateLeaveTypeRequest) (*LeaveTypeResponse, error) {
	t := &LeaveType{
		Code:               req.Code,
		Name:               req.Name,
		Description:        req.Description,
		IsPaid:             true,
		RequiresAttachment: false,
		AllowHalfDay:       true,
		QuotaPeriod:        QuotaPeriodYear,
		CountsAgainstQuota: true,
		AllowHourly:        false,
		IsActive:           true,
	}
	if req.IsPaid != nil {
		t.IsPaid = *req.IsPaid
	}
	if req.RequiresAttachment != nil {
		t.RequiresAttachment = *req.RequiresAttachment
	}
	if req.AllowHalfDay != nil {
		t.AllowHalfDay = *req.AllowHalfDay
	}
	if req.DefaultQuotaDays != nil {
		t.DefaultQuotaDays = req.DefaultQuotaDays
	}
	if req.QuotaPeriod != "" {
		t.QuotaPeriod = QuotaPeriod(req.QuotaPeriod)
	}
	if req.CountsAgainstQuota != nil {
		t.CountsAgainstQuota = *req.CountsAgainstQuota
	}
	if req.AllowHourly != nil {
		t.AllowHourly = *req.AllowHourly
	}
	if req.IsActive != nil {
		t.IsActive = *req.IsActive
	}

	if err := s.repo.CreateLeaveType(ctx, t); err != nil {
		return nil, err
	}
	s.logger.Info("Leave type created", zap.String("id", t.ID.String()), zap.String("name", t.Name))
	return leaveTypeToResponse(t), nil
}

func (s *Service) GetLeaveTypeByID(ctx context.Context, id string) (*LeaveTypeResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	t, err := s.repo.FindLeaveTypeByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return leaveTypeToResponse(t), nil
}

func (s *Service) ListLeaveTypes(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	types, total, err := s.repo.ListLeaveTypes(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]LeaveTypeResponse, 0, len(types))
	for _, t := range types {
		responses = append(responses, *leaveTypeToResponse(&t))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateLeaveType(ctx context.Context, id string, req UpdateLeaveTypeRequest) (*LeaveTypeResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	t, err := s.repo.FindLeaveTypeByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Code != nil {
		t.Code = *req.Code
	}
	if req.Name != nil {
		t.Name = *req.Name
	}
	if req.Description != nil {
		t.Description = *req.Description
	}
	if req.IsPaid != nil {
		t.IsPaid = *req.IsPaid
	}
	if req.RequiresAttachment != nil {
		t.RequiresAttachment = *req.RequiresAttachment
	}
	if req.AllowHalfDay != nil {
		t.AllowHalfDay = *req.AllowHalfDay
	}
	if req.DefaultQuotaDays != nil {
		t.DefaultQuotaDays = req.DefaultQuotaDays
	}
	if req.QuotaPeriod != nil {
		t.QuotaPeriod = QuotaPeriod(*req.QuotaPeriod)
	}
	if req.CountsAgainstQuota != nil {
		t.CountsAgainstQuota = *req.CountsAgainstQuota
	}
	if req.AllowHourly != nil {
		t.AllowHourly = *req.AllowHourly
	}
	if req.IsActive != nil {
		t.IsActive = *req.IsActive
	}
	if err := s.repo.UpdateLeaveType(ctx, t); err != nil {
		return nil, err
	}
	return leaveTypeToResponse(t), nil
}

func (s *Service) DeleteLeaveType(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteLeaveType(ctx, uid)
}

// =========================================================================
// Accrual Policies
// =========================================================================

func (s *Service) CreateAccrualPolicy(ctx context.Context, req CreateAccrualPolicyRequest) (*AccrualPolicyResponse, error) {
	lTypeID, err := uuid.Parse(req.LeaveTypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid leave_type_id: %w", err)
	}

	p := &LeaveAccrualPolicy{
		LeaveTypeID:     lTypeID,
		BaseQuotaDays:   req.BaseQuotaDays,
		ExtraEveryYears: 2,
		ExtraDays:       1.00,
		EffectiveFrom:   req.EffectiveFrom,
		EffectiveTo:     req.EffectiveTo,
	}
	if req.ExtraEveryYears != nil {
		p.ExtraEveryYears = *req.ExtraEveryYears
	}
	if req.ExtraDays != nil {
		p.ExtraDays = *req.ExtraDays
	}
	p.MaxExtraDays = req.MaxExtraDays

	if err := s.repo.CreateAccrualPolicy(ctx, p); err != nil {
		return nil, err
	}
	return accrualPolicyToResponse(p), nil
}

func (s *Service) GetAccrualPolicyByID(ctx context.Context, id string) (*AccrualPolicyResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindAccrualPolicyByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return accrualPolicyToResponse(p), nil
}

func (s *Service) ListAccrualPolicies(ctx context.Context, leaveTypeID *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var lTypeUUID *uuid.UUID
	if leaveTypeID != nil && *leaveTypeID != "" {
		uid, err := uuid.Parse(*leaveTypeID)
		if err != nil {
			return nil, fmt.Errorf("invalid leave_type_id: %w", err)
		}
		lTypeUUID = &uid
	}
	policies, total, err := s.repo.ListAccrualPolicies(ctx, lTypeUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]AccrualPolicyResponse, 0, len(policies))
	for _, p := range policies {
		responses = append(responses, *accrualPolicyToResponse(&p))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateAccrualPolicy(ctx context.Context, id string, req UpdateAccrualPolicyRequest) (*AccrualPolicyResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindAccrualPolicyByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.LeaveTypeID != nil {
		lTypeID, err := uuid.Parse(*req.LeaveTypeID)
		if err != nil {
			return nil, fmt.Errorf("invalid leave_type_id: %w", err)
		}
		p.LeaveTypeID = lTypeID
	}
	if req.BaseQuotaDays != nil {
		p.BaseQuotaDays = *req.BaseQuotaDays
	}
	if req.ExtraEveryYears != nil {
		p.ExtraEveryYears = *req.ExtraEveryYears
	}
	if req.ExtraDays != nil {
		p.ExtraDays = *req.ExtraDays
	}
	if req.MaxExtraDays != nil {
		p.MaxExtraDays = req.MaxExtraDays
	}
	if req.EffectiveFrom != nil {
		p.EffectiveFrom = *req.EffectiveFrom
	}
	if req.EffectiveTo != nil {
		p.EffectiveTo = req.EffectiveTo
	}
	if err := s.repo.UpdateAccrualPolicy(ctx, p); err != nil {
		return nil, err
	}
	return accrualPolicyToResponse(p), nil
}

func (s *Service) DeleteAccrualPolicy(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteAccrualPolicy(ctx, uid)
}

// =========================================================================
// Leave Reasons
// =========================================================================

func (s *Service) CreateLeaveReason(ctx context.Context, req CreateLeaveReasonRequest) (*LeaveReasonResponse, error) {
	r := &LeaveReason{
		Name:      req.Name,
		IsOther:   false,
		SortOrder: 0,
	}
	if req.IsOther != nil {
		r.IsOther = *req.IsOther
	}
	if req.SortOrder != nil {
		r.SortOrder = *req.SortOrder
	}
	if err := s.repo.CreateLeaveReason(ctx, r); err != nil {
		return nil, err
	}
	return leaveReasonToResponse(r), nil
}

func (s *Service) GetLeaveReasonByID(ctx context.Context, id string) (*LeaveReasonResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	r, err := s.repo.FindLeaveReasonByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return leaveReasonToResponse(r), nil
}

func (s *Service) ListLeaveReasons(ctx context.Context) ([]LeaveReasonResponse, error) {
	reasons, err := s.repo.ListLeaveReasons(ctx)
	if err != nil {
		return nil, err
	}
	responses := make([]LeaveReasonResponse, 0, len(reasons))
	for _, r := range reasons {
		responses = append(responses, *leaveReasonToResponse(&r))
	}
	return responses, nil
}

func (s *Service) UpdateLeaveReason(ctx context.Context, id string, req UpdateLeaveReasonRequest) (*LeaveReasonResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	r, err := s.repo.FindLeaveReasonByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		r.Name = *req.Name
	}
	if req.IsOther != nil {
		r.IsOther = *req.IsOther
	}
	if req.SortOrder != nil {
		r.SortOrder = *req.SortOrder
	}
	if err := s.repo.UpdateLeaveReason(ctx, r); err != nil {
		return nil, err
	}
	return leaveReasonToResponse(r), nil
}

func (s *Service) DeleteLeaveReason(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteLeaveReason(ctx, uid)
}

// =========================================================================
// Leave Requests
// =========================================================================

func (s *Service) CreateLeaveRequest(ctx context.Context, req CreateLeaveRequest) (*LeaveRequestResponse, error) {
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	lTypeID, err := uuid.Parse(req.LeaveTypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid leave_type_id: %w", err)
	}
	leaveType, err := s.repo.FindLeaveTypeByID(ctx, lTypeID)
	if err != nil {
		return nil, fmt.Errorf("leave type not found: %w", err)
	}
	if !leaveType.IsActive {
		return nil, fmt.Errorf("leave type %q is not active", leaveType.Name)
	}
	if leaveType.RequiresAttachment && (req.AttachmentURL == nil || *req.AttachmentURL == "") {
		return nil, fmt.Errorf("leave type %q requires an attachment", leaveType.Name)
	}

	overlapCount, err := s.repo.CountOverlappingLeaveRequests(ctx, empID, req.RequestStartDate, req.RequestEndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to check overlapping leave requests: %w", err)
	}
	if overlapCount > 0 {
		return nil, fmt.Errorf("requested dates overlap with an existing leave request")
	}

	mode := DurationFullDay
	if req.DurationMode != "" {
		mode = DurationMode(req.DurationMode)
	}
	if (mode == DurationHalfDayAM || mode == DurationHalfDayPM) && !leaveType.AllowHalfDay {
		return nil, fmt.Errorf("leave type %q does not allow half day requests", leaveType.Name)
	}
	if mode == DurationHourly && !leaveType.AllowHourly {
		return nil, fmt.Errorf("leave type %q does not allow hourly requests", leaveType.Name)
	}

	holidayDates := s.holidayDatesInRange(ctx, req.RequestStartDate, req.RequestEndDate)
	requestedDays, dayDetails, err := CalculateLeaveDuration(req.RequestStartDate, req.RequestEndDate, mode, req.StartTime, req.EndTime, holidayDates)
	if err != nil {
		return nil, err
	}

	lr := &LeaveRequest{
		EmployeeID:       empID,
		LeaveTypeID:      lTypeID,
		RequestStartDate: req.RequestStartDate,
		RequestEndDate:   req.RequestEndDate,
		DurationMode:     mode,
		RequestedDays:    requestedDays,
		Status:           LeaveStatusSubmitted,
		SubmittedAt:      timeNowPtr(),
	}
	if req.LeaveReasonID != nil {
		rID, err := uuid.Parse(*req.LeaveReasonID)
		if err == nil {
			lr.LeaveReasonID = &rID
		}
	}
	lr.LeaveReasonNote = req.LeaveReasonNote
	lr.AttachmentURL = req.AttachmentURL
	if req.SupervisorID != nil {
		sID, err := uuid.Parse(*req.SupervisorID)
		if err == nil {
			lr.SupervisorID = &sID
		}
	}
	lr.StartTime = req.StartTime
	lr.EndTime = req.EndTime

	if err := s.repo.CreateLeaveRequest(ctx, lr); err != nil {
		return nil, err
	}

	for _, dd := range dayDetails {
		detail := &LeaveRequestDetail{
			LeaveRequestID: lr.ID,
			EmployeeID:     lr.EmployeeID,
			LeaveDate:      dd.Date,
			DayFraction:    dd.DayFraction,
			IsPaid:         leaveType.IsPaid,
		}
		if err := s.repo.CreateLeaveRequestDetail(ctx, detail); err != nil {
			return nil, fmt.Errorf("failed to create leave request detail: %w", err)
		}
	}

	// Route through the central approval module when a flow is selected,
	// instead of relying on leave's own ad-hoc SupervisorID/HrID fields.
	if s.approvalEngine != nil && req.FlowID != nil && *req.FlowID != "" {
		instanceID, err := s.approvalEngine.CreateApprovalInstance(ctx, "leave", lr.ID.String(), *req.FlowID)
		if err != nil {
			s.logger.Warn("Failed to create approval instance for leave request, continuing without approval",
				zap.String("leave_request_id", lr.ID.String()),
				zap.Error(err),
			)
		} else {
			if parsedInstanceID, parseErr := uuid.Parse(instanceID); parseErr == nil {
				lr.ApprovalInstanceID = &parsedInstanceID
			}
			lr.Status = LeaveStatusPendingApproval
			if err := s.repo.UpdateLeaveRequest(ctx, lr); err != nil {
				return nil, err
			}
		}
	}

	s.logger.Info("Leave request created", zap.String("id", lr.ID.String()), zap.String("employee_id", lr.EmployeeID.String()))
	return leaveRequestToResponse(lr), nil
}

func (s *Service) GetLeaveRequestByID(ctx context.Context, id string) (*LeaveRequestResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	lr, err := s.repo.FindLeaveRequestByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return leaveRequestToResponse(lr), nil
}

func (s *Service) ListLeaveRequests(ctx context.Context, employeeID *string, status *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var empUUID *uuid.UUID
	if employeeID != nil && *employeeID != "" {
		uid, err := uuid.Parse(*employeeID)
		if err != nil {
			return nil, fmt.Errorf("invalid employee_id: %w", err)
		}
		empUUID = &uid
	}
	requests, total, err := s.repo.ListLeaveRequests(ctx, empUUID, status, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]LeaveRequestResponse, 0, len(requests))
	for _, r := range requests {
		responses = append(responses, *leaveRequestToResponse(&r))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateLeaveRequestStatus(ctx context.Context, id, status, note string) (*LeaveRequestResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	lr, err := s.repo.FindLeaveRequestByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	previousStatus := lr.Status
	now := time.Now()
	switch LeaveStatus(status) {
	case LeaveStatusApprovedFinal:
		lr.Status = LeaveStatusApprovedFinal
		lr.ApprovedAt = &now
	case LeaveStatusRejectedFinal:
		lr.Status = LeaveStatusRejectedFinal
		lr.RejectedAt = &now
	case LeaveStatusCancelled:
		lr.Status = LeaveStatusCancelled
		lr.CancelledAt = &now
	default:
		lr.Status = LeaveStatus(status)
	}
	if note != "" {
		lr.HrNote = &note
	}

	if err := s.repo.UpdateLeaveRequest(ctx, lr); err != nil {
		return nil, err
	}
	if err := s.applyBalanceEffectOnStatusChange(ctx, lr, previousStatus); err != nil {
		return nil, err
	}
	return leaveRequestToResponse(lr), nil
}

// applyBalanceEffectOnStatusChange deducts or reverses the leave balance when
// a status transition crosses into/out of APPROVED_FINAL (§34/§18) — the
// leave balance ledger is a side effect of the status change, not of the
// specific endpoint/caller that triggered it.
func (s *Service) applyBalanceEffectOnStatusChange(ctx context.Context, lr *LeaveRequest, previousStatus LeaveStatus) error {
	if lr.Status == previousStatus {
		return nil
	}
	if lr.Status != LeaveStatusApprovedFinal && previousStatus != LeaveStatusApprovedFinal {
		return nil
	}
	leaveType, err := s.repo.FindLeaveTypeByID(ctx, lr.LeaveTypeID)
	if err != nil {
		return fmt.Errorf("leave type not found: %w", err)
	}
	if lr.Status == LeaveStatusApprovedFinal {
		return s.applyLeaveUsage(ctx, lr, leaveType)
	}
	// previousStatus was APPROVED_FINAL and it's moving away from it
	// (cancelled/rejected via manual override) — reverse the deduction.
	return s.reverseLeaveUsage(ctx, lr, leaveType)
}

// HandleApprovalStatusChange is invoked by the approval module's push-based
// status callback when a leave request's approval instance reaches a final
// state, so the leave request's own status field updates itself.
func (s *Service) HandleApprovalStatusChange(ctx context.Context, documentID uuid.UUID, status string, note string) error {
	lr, err := s.repo.FindLeaveRequestByID(ctx, documentID)
	if err != nil {
		return err
	}
	if lr.Status != LeaveStatusPendingApproval {
		return nil
	}

	now := time.Now()
	switch status {
	case "APPROVED":
		lr.Status = LeaveStatusApprovedFinal
		lr.ApprovedAt = &now
	case "REJECTED":
		lr.Status = LeaveStatusRejectedFinal
		lr.RejectedAt = &now
		if note != "" {
			lr.SupervisorNote = &note
		}
	case "CANCELLED":
		lr.Status = LeaveStatusCancelled
		lr.CancelledAt = &now
	default:
		return nil
	}

	s.logger.Info("Leave request status updated via approval status handler",
		zap.String("leave_request_id", lr.ID.String()),
		zap.String("approval_status", status),
	)
	if err := s.repo.UpdateLeaveRequest(ctx, lr); err != nil {
		return err
	}
	// Coming from PENDING_APPROVAL, only the APPROVED_FINAL transition can
	// require a balance deduction — REJECTED/CANCELLED from this state never
	// had the balance touched in the first place, so no reversal applies.
	if lr.Status == LeaveStatusApprovedFinal {
		leaveType, err := s.repo.FindLeaveTypeByID(ctx, lr.LeaveTypeID)
		if err != nil {
			return fmt.Errorf("leave type not found: %w", err)
		}
		if err := s.applyLeaveUsage(ctx, lr, leaveType); err != nil {
			return err
		}
		s.applyAttendanceIntegration(ctx, lr)
	}
	if lr.Status == LeaveStatusApprovedFinal || lr.Status == LeaveStatusRejectedFinal {
		s.notifyLeaveOutcome(ctx, lr)
	}
	return nil
}

// notifyLeaveOutcome notifies the requesting employee of their leave
// request's final approval outcome (docs/module-notification-plan.md §8).
// Best-effort: if the notifier isn't wired, the employee has no linked user
// account, or Notify fails, this logs and moves on rather than failing the
// approval itself.
func (s *Service) notifyLeaveOutcome(ctx context.Context, lr *LeaveRequest) {
	if s.notifier == nil {
		return
	}
	userID, err := s.repo.FindUserIDByEmployeeID(ctx, lr.EmployeeID)
	if err != nil {
		s.logger.Warn("Failed to resolve employee user id for leave notification",
			zap.String("leave_request_id", lr.ID.String()),
			zap.Error(err),
		)
		return
	}
	if userID == nil {
		return
	}

	var notifType, title, body string
	switch lr.Status {
	case LeaveStatusApprovedFinal:
		notifType = "LEAVE_APPROVED"
		title = "Leave Request Approved"
		body = "Your leave request has been approved."
	case LeaveStatusRejectedFinal:
		notifType = "LEAVE_REJECTED"
		title = "Leave Request Rejected"
		body = "Your leave request has been rejected."
	default:
		return
	}

	if err := s.notifier.Notify(ctx, *userID, notifType, title, body, "leave", lr.ID); err != nil {
		s.logger.Warn("Failed to send leave notification",
			zap.String("leave_request_id", lr.ID.String()),
			zap.String("notif_type", notifType),
			zap.Error(err),
		)
	}
}

// applyAttendanceIntegration pushes each of the leave request's per-date
// entries onto Attendance's session for that date (docs/module-attendance-plan.md
// §26/§50), so a day fully covered by approved leave is reflected in
// Attendance without needing a check-in event to trigger session
// generation. Best-effort: if the updater isn't wired, or a given date
// fails, this logs and moves on rather than failing the whole approval -
// the leave itself is still correctly approved either way.
func (s *Service) applyAttendanceIntegration(ctx context.Context, lr *LeaveRequest) {
	if s.attendanceUpdater == nil {
		return
	}
	details, err := s.repo.ListLeaveRequestDetails(ctx, lr.ID)
	if err != nil {
		s.logger.Warn("Failed to load leave request details for attendance integration",
			zap.String("leave_request_id", lr.ID.String()),
			zap.Error(err),
		)
		return
	}
	for _, d := range details {
		if err := s.attendanceUpdater.ApplyApprovedLeave(ctx, lr.EmployeeID, d.LeaveDate, lr.ID, d.DayFraction); err != nil {
			s.logger.Warn("Failed to apply approved leave to attendance session",
				zap.String("leave_request_id", lr.ID.String()),
				zap.String("leave_date", d.LeaveDate),
				zap.Error(err),
			)
		}
	}
}

func (s *Service) DeleteLeaveRequest(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteLeaveRequest(ctx, uid)
}

// =========================================================================
// Leave Request Details
// =========================================================================

func (s *Service) ListLeaveRequestDetails(ctx context.Context, leaveRequestID string) ([]LeaveRequestDetailResponse, error) {
	lrID, err := uuid.Parse(leaveRequestID)
	if err != nil {
		return nil, fmt.Errorf("invalid leave_request_id: %w", err)
	}
	details, err := s.repo.ListLeaveRequestDetails(ctx, lrID)
	if err != nil {
		return nil, err
	}
	responses := make([]LeaveRequestDetailResponse, 0, len(details))
	for _, d := range details {
		responses = append(responses, *leaveRequestDetailToResponse(&d))
	}
	return responses, nil
}

// =========================================================================
// Calendar (Phase 7 - Employee Calendar, §20)
// =========================================================================

// GetEmployeeCalendar returns the employee's leave dates in [fromDate,
// toDate], mirroring attendance.Service.GetEmployeeCalendar's shape/purpose.
// Team/Organization Calendar are intentionally not built here — they'd need
// a cross-module employee/organization read ("who reports to whom") that
// doesn't exist anywhere in this codebase yet, the same category of gap
// already deferred for Attendance's Manager/HR Dashboard (Phase 10).
func (s *Service) GetEmployeeCalendar(ctx context.Context, employeeID, fromDate, toDate string) ([]CalendarEntryResponse, error) {
	empUUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	entries, err := s.repo.FindCalendarEntriesForEmployeeInRange(ctx, empUUID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	responses := make([]CalendarEntryResponse, 0, len(entries))
	for _, e := range entries {
		responses = append(responses, CalendarEntryResponse{
			LeaveRequestID: e.LeaveRequestID.String(),
			LeaveDate:      normalizeLeaveDate(e.LeaveDate),
			DayFraction:    e.DayFraction,
			LeaveTypeID:    e.LeaveTypeID.String(),
			Status:         e.Status,
		})
	}
	return responses, nil
}

// =========================================================================
// Leave Balances
// =========================================================================

func (s *Service) GetLeaveBalance(ctx context.Context, employeeID, leaveTypeID string, year int) (*LeaveBalanceResponse, error) {
	empUUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	lTypeUUID, err := uuid.Parse(leaveTypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid leave_type_id: %w", err)
	}
	bal, err := s.repo.FindLeaveBalance(ctx, empUUID, lTypeUUID, year)
	if err != nil {
		return nil, err
	}
	return leaveBalanceToResponse(bal), nil
}

func (s *Service) ListLeaveBalances(ctx context.Context, employeeID *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var empUUID *uuid.UUID
	if employeeID != nil && *employeeID != "" {
		uid, err := uuid.Parse(*employeeID)
		if err != nil {
			return nil, fmt.Errorf("invalid employee_id: %w", err)
		}
		empUUID = &uid
	}
	balances, total, err := s.repo.ListLeaveBalances(ctx, empUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]LeaveBalanceResponse, 0, len(balances))
	for _, b := range balances {
		responses = append(responses, *leaveBalanceToResponse(&b))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

// =========================================================================
// Helpers
// =========================================================================

func calcTotalPages(total int64, perPage int) int {
	pages := int(math.Ceil(float64(total) / float64(perPage)))
	if pages < 1 {
		return 1
	}
	return pages
}

func timeNowPtr() *time.Time {
	now := time.Now()
	return &now
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// =========================================================================
// Response converters
// =========================================================================

func leaveTypeToResponse(t *LeaveType) *LeaveTypeResponse {
	return &LeaveTypeResponse{
		ID:                 t.ID.String(),
		Code:               t.Code,
		Name:               t.Name,
		Description:        t.Description,
		IsPaid:             t.IsPaid,
		RequiresAttachment: t.RequiresAttachment,
		AllowHalfDay:       t.AllowHalfDay,
		DefaultQuotaDays:   t.DefaultQuotaDays,
		QuotaPeriod:        string(t.QuotaPeriod),
		CountsAgainstQuota: t.CountsAgainstQuota,
		AllowHourly:        t.AllowHourly,
		IsActive:           t.IsActive,
		CreatedAt:          t.CreatedAt,
		UpdatedAt:          t.UpdatedAt,
	}
}

func accrualPolicyToResponse(p *LeaveAccrualPolicy) *AccrualPolicyResponse {
	return &AccrualPolicyResponse{
		ID:              p.ID.String(),
		LeaveTypeID:     p.LeaveTypeID.String(),
		BaseQuotaDays:   p.BaseQuotaDays,
		ExtraEveryYears: p.ExtraEveryYears,
		ExtraDays:       p.ExtraDays,
		MaxExtraDays:    p.MaxExtraDays,
		EffectiveFrom:   p.EffectiveFrom,
		EffectiveTo:     p.EffectiveTo,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

func leaveReasonToResponse(r *LeaveReason) *LeaveReasonResponse {
	return &LeaveReasonResponse{
		ID:        r.ID.String(),
		Name:      r.Name,
		IsOther:   r.IsOther,
		SortOrder: r.SortOrder,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func leaveRequestToResponse(r *LeaveRequest) *LeaveRequestResponse {
	resp := &LeaveRequestResponse{
		ID:               r.ID.String(),
		EmployeeID:       r.EmployeeID.String(),
		LeaveTypeID:      r.LeaveTypeID.String(),
		RequestStartDate: r.RequestStartDate,
		RequestEndDate:   r.RequestEndDate,
		DurationMode:     string(r.DurationMode),
		RequestedDays:    r.RequestedDays,
		LeaveReasonNote:  r.LeaveReasonNote,
		AttachmentURL:    r.AttachmentURL,
		Status:           string(r.Status),
		SupervisorNote:   r.SupervisorNote,
		HrNote:           r.HrNote,
		StartTime:        r.StartTime,
		EndTime:          r.EndTime,
		SubmittedAt:      r.SubmittedAt,
		ApprovedAt:       r.ApprovedAt,
		RejectedAt:       r.RejectedAt,
		CancelledAt:      r.CancelledAt,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}
	if r.LeaveReasonID != nil {
		v := r.LeaveReasonID.String()
		resp.LeaveReasonID = &v
	}
	if r.SupervisorID != nil {
		v := r.SupervisorID.String()
		resp.SupervisorID = &v
	}
	if r.HrID != nil {
		v := r.HrID.String()
		resp.HrID = &v
	}
	if r.ApprovalInstanceID != nil {
		v := r.ApprovalInstanceID.String()
		resp.ApprovalInstanceID = &v
	}
	return resp
}

func leaveRequestDetailToResponse(d *LeaveRequestDetail) *LeaveRequestDetailResponse {
	return &LeaveRequestDetailResponse{
		ID:             d.ID.String(),
		LeaveRequestID: d.LeaveRequestID.String(),
		EmployeeID:     d.EmployeeID.String(),
		LeaveDate:      d.LeaveDate,
		DayFraction:    d.DayFraction,
		IsPaid:         d.IsPaid,
		CreatedAt:      d.CreatedAt,
	}
}

func leaveBalanceToResponse(b *EmployeeLeaveBalance) *LeaveBalanceResponse {
	return &LeaveBalanceResponse{
		ID:                b.ID.String(),
		EmployeeID:        b.EmployeeID.String(),
		LeaveTypeID:       b.LeaveTypeID.String(),
		PeriodYear:        b.PeriodYear,
		QuotaDays:         b.QuotaDays,
		UsedDays:          b.UsedDays,
		RemainingDays:     b.RemainingDays,
		LastAdjustmentRef: b.LastAdjustmentRef,
		LastAdjustmentAt:  b.LastAdjustmentAt,
		CreatedAt:         b.CreatedAt,
		UpdatedAt:         b.UpdatedAt,
	}
}
