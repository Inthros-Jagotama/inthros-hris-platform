package attendance

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

// ApprovalEngine abstracts the central approval module so overtime requests
// can be routed through it. Implemented via an adapter wrapping
// approval.Service in main.go (same narrow-interface-plus-adapter pattern
// payroll/leave/reimbursement/employeemovement already use).
type ApprovalEngine interface {
	CreateApprovalInstance(ctx context.Context, module, documentID, flowID string) (string, error)
	GetApprovalInstanceStatus(ctx context.Context, instanceID string) (string, error)
}

type Service struct {
	repo           *Repository
	logger         *zap.Logger
	approvalEngine ApprovalEngine
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// SetApprovalEngine wires the central approval module into this service.
func (s *Service) SetApprovalEngine(ae ApprovalEngine) {
	s.approvalEngine = ae
}

// =========================================================================
// Company Settings
// =========================================================================

func (s *Service) UpsertCompanySetting(ctx context.Context, req CreateCompanySettingRequest) (*CompanySettingResponse, error) {
	setting := &AttendanceCompanySetting{}
	if req.IsLocationRequired != nil {
		setting.IsLocationRequired = *req.IsLocationRequired
	}
	if req.IsFaceRequired != nil {
		setting.IsFaceRequired = *req.IsFaceRequired
	}
	if req.IsOvertimeEnabled != nil {
		setting.IsOvertimeEnabled = *req.IsOvertimeEnabled
	}
	setting.Latitude = req.Latitude
	setting.Longitude = req.Longitude
	setting.MaxDistanceMeter = req.MaxDistanceMeter
	setting.LateToleranceMinutes = req.LateToleranceMinutes
	setting.OvertimeMinMinutes = req.OvertimeMinMinutes
	setting.CreatedBy = authctx.GetUserID(ctx)
	setting.UpdatedBy = setting.CreatedBy

	if err := s.repo.UpsertCompanySetting(ctx, setting); err != nil {
		return nil, err
	}

	s.logger.Info("Company attendance setting updated", zap.String("id", setting.ID.String()))
	return settingToResponse(setting), nil
}

func (s *Service) GetCompanySetting(ctx context.Context) (*CompanySettingResponse, error) {
	setting, err := s.repo.FindCompanySetting(ctx)
	if err != nil {
		return nil, err
	}
	return settingToResponse(setting), nil
}

// =========================================================================
// Company Shifts
// =========================================================================

func (s *Service) CreateShift(ctx context.Context, req CreateCompanyShiftRequest) (*CompanyShiftResponse, error) {
	shift := &AttendanceCompanyShift{
		ShiftName:    req.ShiftName,
		CheckInTime:  req.CheckInTime,
		CheckOutTime: req.CheckOutTime,
	}
	if req.IsCrossMidnight != nil {
		shift.IsCrossMidnight = *req.IsCrossMidnight
	}
	shift.CreatedBy = authctx.GetUserID(ctx)
	shift.UpdatedBy = shift.CreatedBy

	if err := s.repo.CreateShift(ctx, shift); err != nil {
		return nil, err
	}
	return shiftToResponse(shift), nil
}

func (s *Service) GetShiftByID(ctx context.Context, id string) (*CompanyShiftResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	shift, err := s.repo.FindShiftByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return shiftToResponse(shift), nil
}

func (s *Service) ListShifts(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	shifts, total, err := s.repo.ListShifts(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]CompanyShiftResponse, 0, len(shifts))
	for _, shift := range shifts {
		responses = append(responses, *shiftToResponse(&shift))
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

func (s *Service) UpdateShift(ctx context.Context, id string, req UpdateCompanyShiftRequest) (*CompanyShiftResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	shift, err := s.repo.FindShiftByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	shift.UpdatedBy = authctx.GetUserID(ctx)
	if req.ShiftName != nil {
		shift.ShiftName = *req.ShiftName
	}
	if req.CheckInTime != nil {
		shift.CheckInTime = *req.CheckInTime
	}
	if req.CheckOutTime != nil {
		shift.CheckOutTime = *req.CheckOutTime
	}
	if req.IsCrossMidnight != nil {
		shift.IsCrossMidnight = *req.IsCrossMidnight
	}
	if err := s.repo.UpdateShift(ctx, shift); err != nil {
		return nil, err
	}
	return shiftToResponse(shift), nil
}

func (s *Service) DeleteShift(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteShift(ctx, uid)
}

// =========================================================================
// Employee Shifts
// =========================================================================

func (s *Service) CreateEmployeeShift(ctx context.Context, req CreateEmployeeShiftRequest) (*EmployeeShiftResponse, error) {
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	shiftID, err := uuid.Parse(req.AttendanceShiftID)
	if err != nil {
		return nil, fmt.Errorf("invalid attendance_shift_id: %w", err)
	}
	if req.EffectiveDateTo != nil && *req.EffectiveDateTo < req.EffectiveDateFrom {
		return nil, fmt.Errorf("effective_date_to must not be before effective_date_from")
	}
	overlapCount, err := s.repo.CountOverlappingEmployeeShifts(ctx, empID, req.EffectiveDateFrom, req.EffectiveDateTo, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check overlapping shift assignments: %w", err)
	}
	if overlapCount > 0 {
		return nil, fmt.Errorf("employee already has a shift assignment overlapping this date range")
	}

	es := &AttendanceEmployeeShift{
		EmployeeID:        empID,
		AttendanceShiftID: shiftID,
		EffectiveDateFrom: req.EffectiveDateFrom,
		EffectiveDateTo:   req.EffectiveDateTo,
		DaysOfWeekMask:    req.DaysOfWeekMask,
		IsDayOff:          req.IsDayOff,
	}
	es.CreatedBy = authctx.GetUserID(ctx)
	es.UpdatedBy = es.CreatedBy
	if err := s.repo.CreateEmployeeShift(ctx, es); err != nil {
		return nil, err
	}
	return employeeShiftToResponse(es), nil
}

func (s *Service) GetEmployeeShiftByID(ctx context.Context, id string) (*EmployeeShiftResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	es, err := s.repo.FindEmployeeShiftByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return employeeShiftToResponse(es), nil
}

func (s *Service) ListEmployeeShifts(ctx context.Context, employeeID *string, page, perPage int) (*PaginatedResponse, error) {
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
	shifts, total, err := s.repo.ListEmployeeShifts(ctx, empUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]EmployeeShiftResponse, 0, len(shifts))
	for _, es := range shifts {
		responses = append(responses, *employeeShiftToResponse(&es))
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

func (s *Service) UpdateEmployeeShift(ctx context.Context, id string, req UpdateEmployeeShiftRequest) (*EmployeeShiftResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	es, err := s.repo.FindEmployeeShiftByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	es.UpdatedBy = authctx.GetUserID(ctx)
	if req.AttendanceShiftID != nil {
		sid, err := uuid.Parse(*req.AttendanceShiftID)
		if err != nil {
			return nil, fmt.Errorf("invalid attendance_shift_id: %w", err)
		}
		es.AttendanceShiftID = sid
	}
	if req.EffectiveDateFrom != nil {
		es.EffectiveDateFrom = *req.EffectiveDateFrom
	}
	if req.EffectiveDateTo != nil {
		es.EffectiveDateTo = req.EffectiveDateTo
	}
	if req.DaysOfWeekMask != nil {
		es.DaysOfWeekMask = req.DaysOfWeekMask
	}
	if req.IsDayOff != nil {
		es.IsDayOff = req.IsDayOff
	}
	if es.EffectiveDateTo != nil && *es.EffectiveDateTo < es.EffectiveDateFrom {
		return nil, fmt.Errorf("effective_date_to must not be before effective_date_from")
	}
	overlapCount, err := s.repo.CountOverlappingEmployeeShifts(ctx, es.EmployeeID, es.EffectiveDateFrom, es.EffectiveDateTo, &es.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check overlapping shift assignments: %w", err)
	}
	if overlapCount > 0 {
		return nil, fmt.Errorf("employee already has a shift assignment overlapping this date range")
	}
	if err := s.repo.UpdateEmployeeShift(ctx, es); err != nil {
		return nil, err
	}
	return employeeShiftToResponse(es), nil
}

func (s *Service) DeleteEmployeeShift(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteEmployeeShift(ctx, uid)
}

// =========================================================================
// Locations
// =========================================================================

func (s *Service) CreateLocation(ctx context.Context, req CreateLocationRequest) (*LocationResponse, error) {
	loc := &AttendanceLocation{
		Name:      req.Name,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		RadiusM:   50,
	}
	if req.RadiusM != nil && *req.RadiusM > 0 {
		loc.RadiusM = *req.RadiusM
	}
	loc.CreatedBy = authctx.GetUserID(ctx)
	loc.UpdatedBy = loc.CreatedBy
	if err := s.repo.CreateLocation(ctx, loc); err != nil {
		return nil, err
	}
	return locationToResponse(loc), nil
}

func (s *Service) GetLocationByID(ctx context.Context, id string) (*LocationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	loc, err := s.repo.FindLocationByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return locationToResponse(loc), nil
}

func (s *Service) ListLocations(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	locs, total, err := s.repo.ListLocations(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]LocationResponse, 0, len(locs))
	for _, loc := range locs {
		responses = append(responses, *locationToResponse(&loc))
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

func (s *Service) UpdateLocation(ctx context.Context, id string, req UpdateLocationRequest) (*LocationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	loc, err := s.repo.FindLocationByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	loc.UpdatedBy = authctx.GetUserID(ctx)
	if req.Name != nil {
		loc.Name = *req.Name
	}
	if req.Latitude != nil {
		loc.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		loc.Longitude = *req.Longitude
	}
	if req.RadiusM != nil && *req.RadiusM > 0 {
		loc.RadiusM = *req.RadiusM
	}
	if err := s.repo.UpdateLocation(ctx, loc); err != nil {
		return nil, err
	}
	return locationToResponse(loc), nil
}

func (s *Service) DeleteLocation(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteLocation(ctx, uid)
}

// =========================================================================
// Events (Check-in / Check-out)
// =========================================================================

func (s *Service) CreateEvent(ctx context.Context, req CreateEventRequest) (*EventResponse, error) {
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	eventType := EventType(req.EventType)
	lastEvent, err := s.checkEventSequence(ctx, empID, eventType)
	if err != nil {
		return nil, err
	}

	utcTime, err := time.Parse(time.RFC3339, req.EventTimeUTC)
	if err != nil {
		return nil, fmt.Errorf("invalid event_time_utc: %w", err)
	}
	localTime, err := time.Parse(time.RFC3339, req.EventTimeLocal)
	if err != nil {
		return nil, fmt.Errorf("invalid event_time_local: %w", err)
	}

	event := &AttendanceEvent{
		EmployeeID:       empID,
		EventType:        eventType,
		EventTimeUTC:     utcTime,
		EventTimeLocal:   localTime,
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
		AccuracyM:        req.AccuracyM,
		LocationProvider: req.LocationProvider,
		ValidationStatus: ValidationPending,
	}
	if req.DeviceID != nil && *req.DeviceID != "" {
		did, err := uuid.Parse(*req.DeviceID)
		if err == nil {
			event.DeviceID = &did
		}
	}

	s.applyEventValidation(ctx, event)

	if err := s.repo.CreateEvent(ctx, event); err != nil {
		return nil, err
	}

	// Work date is the CHECKIN's local calendar date (§24) - for a CHECKOUT
	// that closes a cross-midnight shift, that's lastEvent's date, not this
	// event's own (which lands on the following day).
	workDate := event.EventTimeLocal.Format(workDateLayout)
	if eventType == EventTypeCheckOut && lastEvent != nil {
		workDate = lastEvent.EventTimeLocal.Format(workDateLayout)
	}
	if err := s.recalculateSession(ctx, empID, workDate); err != nil {
		s.logger.Warn("Failed to recalculate attendance session",
			zap.String("employee_id", empID.String()),
			zap.String("work_date", workDate),
			zap.Error(err),
		)
	}

	return eventToResponse(event), nil
}

// checkEventSequence rejects out-of-sequence check-in/check-out submissions
// (§5/§18's "Duplicate Event Detection"): an employee can't check in twice
// without checking out in between, and can't check out without an open
// check-in. This is a lightweight sequence check on raw events - it doesn't
// require resolving the employee's shift/work-date (that's session
// calculation, Phase 6), just the last event's type.
func (s *Service) checkEventSequence(ctx context.Context, employeeID uuid.UUID, eventType EventType) (*AttendanceEvent, error) {
	last, err := s.repo.FindLastEventForEmployee(ctx, employeeID)
	if err != nil {
		// No prior event - CHECKIN is fine, CHECKOUT has nothing to close.
		if eventType == EventTypeCheckOut {
			return nil, fmt.Errorf("cannot check out without an open check-in")
		}
		return nil, nil
	}
	switch eventType {
	case EventTypeCheckIn:
		if last.EventType == EventTypeCheckIn {
			return nil, fmt.Errorf("already checked in; check out before checking in again")
		}
	case EventTypeCheckOut:
		if last.EventType == EventTypeCheckOut {
			return nil, fmt.Errorf("cannot check out without an open check-in")
		}
	}
	return last, nil
}

// applyEventValidation runs the checks §17-18 describe against a not-yet-
// persisted event, setting its ValidationStatus/DistanceM/IsInGeofence
// fields in place. Only location validation is implemented here - face and
// device validation stay PENDING because there is no face-matching provider
// or employee-device mapping in this codebase yet (see docs/module-attendance-plan.md
// §17, Phase 2/4), so this can't silently decide those checks passed.
func (s *Service) applyEventValidation(ctx context.Context, event *AttendanceEvent) {
	setting, err := s.repo.FindCompanySetting(ctx)
	if err != nil {
		// No settings configured yet - nothing to validate against.
		event.ValidationStatus = ValidationPending
		return
	}

	locationOK := true
	if setting.IsLocationRequired {
		locations, err := s.repo.FindAllLocations(ctx)
		if err != nil {
			locations = nil
		}
		locationOK = validateEventLocation(event, locations, setting)
	}

	switch {
	case setting.IsLocationRequired && !locationOK:
		event.ValidationStatus = ValidationInvalid
		note := "outside allowed check-in location"
		event.ValidationNote = &note
	case setting.IsFaceRequired:
		// Face verification has no implementation to consult - leave PENDING
		// rather than mark VALID for a check that was never actually run.
		event.ValidationStatus = ValidationPending
	default:
		event.ValidationStatus = ValidationValid
	}
}

func (s *Service) GetEventByID(ctx context.Context, id string) (*EventResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	event, err := s.repo.FindEventByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return eventToResponse(event), nil
}

func (s *Service) ListEvents(ctx context.Context, employeeID *string, page, perPage int) (*PaginatedResponse, error) {
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
	events, total, err := s.repo.ListEvents(ctx, empUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]EventResponse, 0, len(events))
	for _, event := range events {
		responses = append(responses, *eventToResponse(&event))
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
// Sessions
// =========================================================================

func (s *Service) GetSession(ctx context.Context, employeeID, workDate string) (*SessionResponse, error) {
	empUUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	session, err := s.repo.FindSessionByEmployeeAndDate(ctx, empUUID, workDate)
	if err != nil {
		return nil, err
	}
	return sessionToResponse(session), nil
}

func (s *Service) ListSessions(ctx context.Context, employeeID *string, page, perPage int) (*PaginatedResponse, error) {
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
	sessions, total, err := s.repo.ListSessions(ctx, empUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]SessionResponse, 0, len(sessions))
	for _, session := range sessions {
		responses = append(responses, *sessionToResponse(&session))
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

// GetEmployeeCalendar returns one employee's sessions across a date range
// (§36 Employee/§37 Attendance Calendar), Phase 10. Manager/HR calendars and
// dashboards (§38-39) that aggregate across an organization are NOT
// implemented here - they'd need a cross-module read into the employee
// module to resolve "which employees report to this manager" /
// "which employees belong to this organization", and no such interface
// exists anywhere in this codebase yet (same category of gap already
// documented for Exempt status in Phase 6).
func (s *Service) GetEmployeeCalendar(ctx context.Context, employeeID, fromDate, toDate string) ([]SessionResponse, error) {
	empUUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	sessions, err := s.repo.FindSessionsForEmployeeInRange(ctx, empUUID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	responses := make([]SessionResponse, 0, len(sessions))
	for _, session := range sessions {
		responses = append(responses, *sessionToResponse(&session))
	}
	return responses, nil
}

// GetAttendanceReport returns every employee's sessions in a date range,
// tenant-wide (§40, Phase 11: Daily Attendance when from==to, Monthly
// Attendance for a month-long range). Late/Missing Attendance/Attendance
// Anomaly "reports" are the same underlying data filtered by
// status/lateness_minutes client-side - no separate endpoint per report
// name, since they're all views over the same session rows. Export
// (Excel/CSV/PDF) is a presentation-layer concern left to whatever consumes
// this endpoint, not implemented here.
func (s *Service) GetAttendanceReport(ctx context.Context, fromDate, toDate string) ([]SessionResponse, error) {
	sessions, err := s.repo.FindSessionsInRange(ctx, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	responses := make([]SessionResponse, 0, len(sessions))
	for _, session := range sessions {
		responses = append(responses, *sessionToResponse(&session))
	}
	return responses, nil
}

// GetEmployeeSummary aggregates one employee's sessions over a date range
// into the counts §37's Employee Dashboard shows. See SummaryResponse for
// why Absent isn't included.
func (s *Service) GetEmployeeSummary(ctx context.Context, employeeID, fromDate, toDate string) (*SummaryResponse, error) {
	empUUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	sessions, err := s.repo.FindSessionsForEmployeeInRange(ctx, empUUID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	summary := &SummaryResponse{
		EmployeeID: employeeID,
		FromDate:   fromDate,
		ToDate:     toDate,
	}
	for _, session := range sessions {
		summary.TotalSessions++
		summary.TotalWorkMinutes += session.WorkMinutes
		summary.TotalOvertimeMinutes += session.OvertimeMinutes
		if session.LeaveFraction != nil {
			summary.LeaveDays += *session.LeaveFraction
		}
		switch session.Status {
		case SessionStatusClosed:
			if session.LatenessMinutes > 0 {
				summary.LateDays++
			} else {
				summary.PresentDays++
			}
		case SessionStatusMissingCheckIn:
			summary.MissingCheckinDays++
		case SessionStatusMissingCheckOut:
			summary.MissingCheckoutDays++
		case SessionStatusDayOff:
			summary.DayOffDays++
		}
	}
	return summary, nil
}

// =========================================================================
// Overtime Requests
// =========================================================================

func (s *Service) CreateOvertimeRequest(ctx context.Context, req CreateOvertimeRequest) (*OvertimeResponse, error) {
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTimeLocal)
	if err != nil {
		return nil, fmt.Errorf("invalid start_time_local: %w", err)
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTimeLocal)
	if err != nil {
		return nil, fmt.Errorf("invalid end_time_local: %w", err)
	}

	overtime := &AttendanceOvertimeRequest{
		EmployeeID:       empID,
		WorkDate:         req.WorkDate,
		StartTimeLocal:   startTime,
		EndTimeLocal:     endTime,
		RequestedMinutes: req.RequestedMinutes,
		Reason:           strPtr(req.Reason),
		Status:           OvertimeSubmitted,
	}
	if err := s.repo.CreateOvertimeRequest(ctx, overtime); err != nil {
		return nil, err
	}

	// Route through the central approval module when a flow is selected.
	if s.approvalEngine != nil && req.FlowID != nil && *req.FlowID != "" {
		instanceID, err := s.approvalEngine.CreateApprovalInstance(ctx, "attendance", overtime.ID.String(), *req.FlowID)
		if err != nil {
			s.logger.Warn("Failed to create approval instance for overtime request, continuing without approval",
				zap.String("overtime_request_id", overtime.ID.String()),
				zap.Error(err),
			)
		} else {
			if parsedInstanceID, parseErr := uuid.Parse(instanceID); parseErr == nil {
				overtime.ApprovalInstanceID = &parsedInstanceID
			}
			overtime.Status = OvertimePendingApproval
			if err := s.repo.UpdateOvertimeRequest(ctx, overtime); err != nil {
				return nil, err
			}
		}
	}

	return overtimeToResponse(overtime), nil
}

// HandleApprovalStatusChange is invoked by the approval module's push-based
// status callback when an overtime or correction request's approval
// instance reaches a final state. Both request types share the "attendance"
// module slug, so documentID could belong to either table - overtime is
// tried first (the original type), then correction requests.
func (s *Service) HandleApprovalStatusChange(ctx context.Context, documentID uuid.UUID, status string, note string) error {
	if o, err := s.repo.FindOvertimeRequestByID(ctx, documentID); err == nil {
		return s.handleOvertimeApprovalStatusChange(ctx, o, status, note)
	}
	if c, err := s.repo.FindCorrectionRequestByID(ctx, documentID); err == nil {
		return s.handleCorrectionApprovalStatusChange(ctx, c, status, note)
	}
	return fmt.Errorf("no overtime or correction request found for document %s", documentID)
}

func (s *Service) handleOvertimeApprovalStatusChange(ctx context.Context, o *AttendanceOvertimeRequest, status string, note string) error {
	if o.Status != OvertimePendingApproval {
		return nil
	}

	now := time.Now()
	switch status {
	case "APPROVED":
		o.Status = OvertimeApproved
		o.ApprovedAt = &now
		s.applyOvertimeCalculation(ctx, o)
	case "REJECTED":
		o.Status = OvertimeRejected
		if note != "" {
			o.ApprovalNote = &note
		}
	case "CANCELLED":
		o.Status = OvertimeRejected
	default:
		return nil
	}

	s.logger.Info("Overtime request status updated via approval status handler",
		zap.String("overtime_request_id", o.ID.String()),
		zap.String("approval_status", status),
	)
	return s.repo.UpdateOvertimeRequest(ctx, o)
}

// applyOvertimeCalculation fills in ActualMinutes/CalculatedMinutes (§31-32)
// once an overtime request is approved: actual overtime is derived from that
// day's session (actual checkout minus planned checkout, per §31's formula),
// then capped by the approved/requested minutes. If the session doesn't
// exist yet or has no checkout recorded (e.g. approval happened before the
// employee actually checked out), both fields are left nil - there's
// nothing to calculate from yet, and re-approving isn't a workflow this
// endpoint supports, so this is a best-effort snapshot at approval time.
func (s *Service) applyOvertimeCalculation(ctx context.Context, o *AttendanceOvertimeRequest) {
	session, err := s.repo.FindSessionByEmployeeAndDate(ctx, o.EmployeeID, normalizeWorkDate(o.WorkDate))
	if err != nil || session.PlannedEndLocal == nil || session.CheckoutEventID == nil {
		return
	}
	checkoutEvent, err := s.repo.FindEventByID(ctx, *session.CheckoutEventID)
	if err != nil {
		return
	}
	actual := int(checkoutEvent.EventTimeLocal.Sub(*session.PlannedEndLocal).Minutes())
	if actual < 0 {
		actual = 0
	}
	calculated := actual
	if calculated > o.RequestedMinutes {
		calculated = o.RequestedMinutes
	}
	o.ActualMinutes = &actual
	o.CalculatedMinutes = &calculated

	session.IsOvertimeDay = true
	session.OvertimeRequestID = &o.ID
	session.ApprovedOvertimeStartLocal = session.PlannedEndLocal
	overtimeEnd := session.PlannedEndLocal.Add(time.Duration(calculated) * time.Minute)
	session.ApprovedOvertimeEndLocal = &overtimeEnd
	session.OvertimeMinutes = calculated
	if err := s.repo.UpsertSession(ctx, session); err != nil {
		s.logger.Warn("Failed to update session with overtime calculation",
			zap.String("employee_id", o.EmployeeID.String()),
			zap.String("work_date", o.WorkDate),
			zap.Error(err),
		)
	}
}

func (s *Service) GetOvertimeRequestByID(ctx context.Context, id string) (*OvertimeResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	o, err := s.repo.FindOvertimeRequestByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return overtimeToResponse(o), nil
}

func (s *Service) ListOvertimeRequests(ctx context.Context, employeeID *string, page, perPage int) (*PaginatedResponse, error) {
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
	requests, total, err := s.repo.ListOvertimeRequests(ctx, empUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]OvertimeResponse, 0, len(requests))
	for _, o := range requests {
		responses = append(responses, *overtimeToResponse(&o))
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
// Correction Requests (§16/§33-34) - Phase 8
// =========================================================================

func (s *Service) CreateCorrectionRequest(ctx context.Context, req CreateCorrectionRequest) (*CorrectionResponse, error) {
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	sessionID, err := uuid.Parse(req.AttendanceSessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid attendance_session_id: %w", err)
	}

	c := &AttendanceCorrectionRequest{
		EmployeeID:          empID,
		AttendanceSessionID: sessionID,
		CorrectionType:      CorrectionType(req.CorrectionType),
		Reason:              req.Reason,
		Status:              CorrectionSubmitted,
	}
	c.CreatedBy = authctx.GetUserID(ctx)
	if req.RequestedCheckin != nil {
		t, err := time.Parse(time.RFC3339, *req.RequestedCheckin)
		if err != nil {
			return nil, fmt.Errorf("invalid requested_checkin: %w", err)
		}
		c.RequestedCheckin = &t
	}
	if req.RequestedCheckout != nil {
		t, err := time.Parse(time.RFC3339, *req.RequestedCheckout)
		if err != nil {
			return nil, fmt.Errorf("invalid requested_checkout: %w", err)
		}
		c.RequestedCheckout = &t
	}
	if err := s.repo.CreateCorrectionRequest(ctx, c); err != nil {
		return nil, err
	}

	// Route through the central approval module when a flow is selected,
	// same "attendance" module slug the overtime flow uses.
	if s.approvalEngine != nil && req.FlowID != nil && *req.FlowID != "" {
		instanceID, err := s.approvalEngine.CreateApprovalInstance(ctx, "attendance", c.ID.String(), *req.FlowID)
		if err != nil {
			s.logger.Warn("Failed to create approval instance for correction request, continuing without approval",
				zap.String("correction_request_id", c.ID.String()),
				zap.Error(err),
			)
		} else {
			if parsedInstanceID, parseErr := uuid.Parse(instanceID); parseErr == nil {
				c.ApprovalInstanceID = &parsedInstanceID
			}
			c.Status = CorrectionPendingApproval
			if err := s.repo.UpdateCorrectionRequest(ctx, c); err != nil {
				return nil, err
			}
		}
	}

	return correctionToResponse(c), nil
}

func (s *Service) handleCorrectionApprovalStatusChange(ctx context.Context, c *AttendanceCorrectionRequest, status string, note string) error {
	if c.Status != CorrectionPendingApproval {
		return nil
	}

	now := time.Now()
	switch status {
	case "APPROVED":
		c.Status = CorrectionApproved
		c.ApprovedAt = &now
		s.applyCorrectionToSession(ctx, c)
	case "REJECTED", "CANCELLED":
		c.Status = CorrectionRejected
	default:
		return nil
	}

	s.logger.Info("Correction request status updated via approval status handler",
		zap.String("correction_request_id", c.ID.String()),
		zap.String("approval_status", status),
	)
	return s.repo.UpdateCorrectionRequest(ctx, c)
}

// applyCorrectionToSession applies an approved correction by inserting a new
// OVERRIDDEN attendance_event with the requested time - never mutating the
// original raw event(s), per §15's immutability principle - then
// recalculating that day's session so the correction actually takes effect.
//
// Only MISSING_CHECKIN/MISSING_CHECKOUT are applied automatically: there's
// no existing event to reconcile against, so inserting the corrected one is
// unambiguous. WRONG_CHECKIN/WRONG_CHECKOUT are intentionally NOT
// auto-applied - recalculateSession's selectCheckinCheckout always takes the
// first CHECKIN/first-following-CHECKOUT of the day, so simply inserting a
// second event wouldn't reliably override the wrong one without either
// excluding the original from selection (no tracking for that exists) or
// mutating it (forbidden). The request is still recorded and approved for
// audit purposes; a human applies the session-level fix outside this flow
// until that selection logic is extended to support it.
func (s *Service) applyCorrectionToSession(ctx context.Context, c *AttendanceCorrectionRequest) {
	session, err := s.repo.FindSessionByID(ctx, c.AttendanceSessionID)
	if err != nil {
		s.logger.Warn("Failed to load session for correction",
			zap.String("correction_request_id", c.ID.String()),
			zap.Error(err),
		)
		return
	}

	var eventType EventType
	var eventTime *time.Time
	switch c.CorrectionType {
	case CorrectionTypeMissingCheckin:
		eventType, eventTime = EventTypeCheckIn, c.RequestedCheckin
	case CorrectionTypeMissingCheckout:
		eventType, eventTime = EventTypeCheckOut, c.RequestedCheckout
	default:
		return // WRONG_CHECKIN/WRONG_CHECKOUT - not auto-applied, see doc comment above.
	}
	if eventTime == nil {
		return
	}

	event := &AttendanceEvent{
		EmployeeID:       c.EmployeeID,
		EventType:        eventType,
		EventTimeUTC:     eventTime.UTC(),
		EventTimeLocal:   *eventTime,
		ValidationStatus: ValidationOverridden,
	}
	if err := s.repo.CreateEvent(ctx, event); err != nil {
		s.logger.Warn("Failed to create corrected event",
			zap.String("correction_request_id", c.ID.String()),
			zap.Error(err),
		)
		return
	}

	if err := s.recalculateSession(ctx, c.EmployeeID, normalizeWorkDate(session.WorkDate)); err != nil {
		s.logger.Warn("Failed to recalculate session after correction",
			zap.String("correction_request_id", c.ID.String()),
			zap.Error(err),
		)
	}
}

func (s *Service) GetCorrectionRequestByID(ctx context.Context, id string) (*CorrectionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	c, err := s.repo.FindCorrectionRequestByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return correctionToResponse(c), nil
}

func (s *Service) ListCorrectionRequests(ctx context.Context, employeeID *string, page, perPage int) (*PaginatedResponse, error) {
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
	requests, total, err := s.repo.ListCorrectionRequests(ctx, empUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]CorrectionResponse, 0, len(requests))
	for _, c := range requests {
		responses = append(responses, *correctionToResponse(&c))
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

func correctionToResponse(c *AttendanceCorrectionRequest) *CorrectionResponse {
	resp := &CorrectionResponse{
		ID:                  c.ID.String(),
		EmployeeID:          c.EmployeeID.String(),
		AttendanceSessionID: c.AttendanceSessionID.String(),
		CorrectionType:      string(c.CorrectionType),
		RequestedCheckin:    c.RequestedCheckin,
		RequestedCheckout:   c.RequestedCheckout,
		Reason:              c.Reason,
		Status:              string(c.Status),
		ApprovedAt:          c.ApprovedAt,
		CreatedAt:           c.CreatedAt,
	}
	if c.ApprovalInstanceID != nil {
		aiID := c.ApprovalInstanceID.String()
		resp.ApprovalInstanceID = &aiID
	}
	return resp
}

// =========================================================================
// Exempt Positions
// =========================================================================

func (s *Service) CreateExemptPosition(ctx context.Context, req CreateExemptPositionRequest) (*ExemptPositionResponse, error) {
	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_id: %w", err)
	}
	p := &AttendanceExemptPosition{
		OrganizationID: orgID,
		IsExempt:       true,
		Note:           strPtr(req.Note),
	}
	if req.IsExempt != nil {
		p.IsExempt = *req.IsExempt
	}
	p.CreatedBy = authctx.GetUserID(ctx)
	p.UpdatedBy = p.CreatedBy
	if err := s.repo.CreateExemptPosition(ctx, p); err != nil {
		return nil, err
	}
	return exemptPositionToResponse(p), nil
}

func (s *Service) GetExemptPositionByID(ctx context.Context, id string) (*ExemptPositionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindExemptPositionByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return exemptPositionToResponse(p), nil
}

func (s *Service) ListExemptPositions(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	positions, total, err := s.repo.ListExemptPositions(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]ExemptPositionResponse, 0, len(positions))
	for _, p := range positions {
		responses = append(responses, *exemptPositionToResponse(&p))
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

func (s *Service) UpdateExemptPosition(ctx context.Context, id string, req UpdateExemptPositionRequest) (*ExemptPositionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindExemptPositionByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	p.UpdatedBy = authctx.GetUserID(ctx)
	if req.IsExempt != nil {
		p.IsExempt = *req.IsExempt
	}
	if req.Note != nil {
		p.Note = req.Note
	}
	if err := s.repo.UpdateExemptPosition(ctx, p); err != nil {
		return nil, err
	}
	return exemptPositionToResponse(p), nil
}

func (s *Service) DeleteExemptPosition(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteExemptPosition(ctx, uid)
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

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func settingToResponse(s *AttendanceCompanySetting) *CompanySettingResponse {
	return &CompanySettingResponse{
		ID:                   s.ID.String(),
		Latitude:             s.Latitude,
		Longitude:            s.Longitude,
		IsLocationRequired:   s.IsLocationRequired,
		IsFaceRequired:       s.IsFaceRequired,
		IsOvertimeEnabled:    s.IsOvertimeEnabled,
		MaxDistanceMeter:     s.MaxDistanceMeter,
		LateToleranceMinutes: s.LateToleranceMinutes,
		OvertimeMinMinutes:   s.OvertimeMinMinutes,
		CreatedAt:            s.CreatedAt,
		UpdatedAt:            s.UpdatedAt,
	}
}

func shiftToResponse(s *AttendanceCompanyShift) *CompanyShiftResponse {
	return &CompanyShiftResponse{
		ID:              s.ID.String(),
		ShiftName:       s.ShiftName,
		CheckInTime:     s.CheckInTime,
		CheckOutTime:    s.CheckOutTime,
		IsCrossMidnight: s.IsCrossMidnight,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
}

func employeeShiftToResponse(es *AttendanceEmployeeShift) *EmployeeShiftResponse {
	resp := &EmployeeShiftResponse{
		ID:                es.ID.String(),
		EmployeeID:        es.EmployeeID.String(),
		AttendanceShiftID: es.AttendanceShiftID.String(),
		EffectiveDateFrom: es.EffectiveDateFrom,
		EffectiveDateTo:   es.EffectiveDateTo,
		DaysOfWeekMask:    es.DaysOfWeekMask,
		IsDayOff:          es.IsDayOff,
		CreatedAt:         es.CreatedAt,
		UpdatedAt:         es.UpdatedAt,
	}
	return resp
}

func locationToResponse(l *AttendanceLocation) *LocationResponse {
	return &LocationResponse{
		ID:        l.ID.String(),
		Name:      l.Name,
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
		RadiusM:   l.RadiusM,
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
	}
}

func eventToResponse(e *AttendanceEvent) *EventResponse {
	resp := &EventResponse{
		ID:               e.ID.String(),
		EmployeeID:       e.EmployeeID.String(),
		EventType:        string(e.EventType),
		EventTimeUTC:     e.EventTimeUTC,
		EventTimeLocal:   e.EventTimeLocal,
		Latitude:         e.Latitude,
		Longitude:        e.Longitude,
		AccuracyM:        e.AccuracyM,
		LocationProvider: e.LocationProvider,
		IsInGeofence:     e.IsInGeofence,
		ValidationStatus: string(e.ValidationStatus),
		ValidationNote:   e.ValidationNote,
		CreatedAt:        e.CreatedAt,
	}
	if e.DeviceID != nil {
		did := e.DeviceID.String()
		resp.DeviceID = &did
	}
	if e.ValidatedLocationID != nil {
		vlid := e.ValidatedLocationID.String()
		resp.ValidatedLocationID = &vlid
	}
	if e.DistanceM != nil {
		resp.DistanceM = e.DistanceM
	}
	return resp
}

func sessionToResponse(s *AttendanceSession) *SessionResponse {
	resp := &SessionResponse{
		ID:                s.ID.String(),
		EmployeeID:        s.EmployeeID.String(),
		WorkDate:          s.WorkDate,
		IsOvertimeDay:     s.IsOvertimeDay,
		Status:            string(s.Status),
		LatenessMinutes:   s.LatenessMinutes,
		EarlyLeaveMinutes: s.EarlyLeaveMinutes,
		WorkMinutes:       s.WorkMinutes,
		BreakMinutes:      s.BreakMinutes,
		OvertimeMinutes:   s.OvertimeMinutes,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
	if s.ShiftID != nil {
		sid := s.ShiftID.String()
		resp.ShiftID = &sid
	}
	if s.CheckinEventID != nil {
		ceid := s.CheckinEventID.String()
		resp.CheckinEventID = &ceid
	}
	if s.CheckoutEventID != nil {
		ceid := s.CheckoutEventID.String()
		resp.CheckoutEventID = &ceid
	}
	if s.LeaveRequestID != nil {
		lrID := s.LeaveRequestID.String()
		resp.LeaveRequestID = &lrID
	}
	if s.PlannedStartLocal != nil {
		resp.PlannedStartLocal = s.PlannedStartLocal
	}
	if s.PlannedEndLocal != nil {
		resp.PlannedEndLocal = s.PlannedEndLocal
	}
	return resp
}

func overtimeToResponse(o *AttendanceOvertimeRequest) *OvertimeResponse {
	resp := &OvertimeResponse{
		ID:                o.ID.String(),
		EmployeeID:        o.EmployeeID.String(),
		WorkDate:          o.WorkDate,
		StartTimeLocal:    o.StartTimeLocal,
		EndTimeLocal:      o.EndTimeLocal,
		RequestedMinutes:  o.RequestedMinutes,
		ActualMinutes:     o.ActualMinutes,
		CalculatedMinutes: o.CalculatedMinutes,
		Reason:            o.Reason,
		Status:            string(o.Status),
		ApprovedAt:        o.ApprovedAt,
		ApprovalNote:      o.ApprovalNote,
		CreatedAt:         o.CreatedAt,
	}
	if o.ApprovedBy != nil {
		ab := o.ApprovedBy.String()
		resp.ApprovedBy = &ab
	}
	if o.ApprovalInstanceID != nil {
		aiID := o.ApprovalInstanceID.String()
		resp.ApprovalInstanceID = &aiID
	}
	return resp
}

func exemptPositionToResponse(p *AttendanceExemptPosition) *ExemptPositionResponse {
	return &ExemptPositionResponse{
		ID:             p.ID.String(),
		OrganizationID: p.OrganizationID.String(),
		IsExempt:       p.IsExempt,
		Note:           p.Note,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}
