package reimbursement

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/modules/approval"
	"github.com/inthros/hris-platform/internal/pkg/authctx"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

// ApprovalEngine abstracts the central approval module so reimbursement
// requests can be routed through it instead of this module's own ad-hoc
// SupervisorID/HrID fields. Implemented via an adapter wrapping
// approval.Service in main.go (same narrow-interface-plus-adapter pattern
// payroll/leave already use).
type ApprovalEngine interface {
	CreateApprovalInstance(ctx context.Context, module, documentID, flowID string) (string, error)
	GetApprovalInstanceStatus(ctx context.Context, instanceID string) (string, error)
	// GetActiveFlowIDForModule lets a reimbursement submission auto-resolve
	// which flow to route through when the client doesn't supply one
	// explicitly (same pattern leave/attendance use) — without this, a
	// request submitted without a flow_id would silently stay SUBMITTED
	// forever and never reach the Approval module.
	GetActiveFlowIDForModule(ctx context.Context, module string) (string, error)
}

// Notifier abstracts the notification module so Reimbursement can notify an
// employee of their request's approval/payment outcome
// (docs/module-reimbursement-development-plan.md §6, Phase 4).
// notification.Service satisfies this structurally.
type Notifier interface {
	Notify(ctx context.Context, recipientUserID uuid.UUID, notifType string, params []string, referenceType string, referenceID uuid.UUID) error
}

type Service struct {
	repo           *Repository
	logger         *zap.Logger
	approvalEngine ApprovalEngine
	notifier       Notifier
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// SetApprovalEngine wires the central approval module into this service.
func (s *Service) SetApprovalEngine(ae ApprovalEngine) {
	s.approvalEngine = ae
}

// SetNotifier wires the notification module into this service so status
// changes (approval outcome, payment) can notify the requesting employee
// (docs/module-reimbursement-development-plan.md §6, Phase 4).
func (s *Service) SetNotifier(n Notifier) {
	s.notifier = n
}

// =========================================================================
// Reimbursement Types
// =========================================================================

func (s *Service) CreateReimbursementType(ctx context.Context, req CreateReimbursementTypeRequest) (*ReimbursementTypeResponse, error) {
	code := req.Code
	if code == "" {
		// Auto-generate a unique code from the display name (e.g. "Medical"
		// -> "MEDICAL") so the create/edit form never needs to ask for one —
		// same pattern as Business Travel expense categories/funding methods.
		existingCodes, err := s.repo.ListReimbursementTypeCodes(ctx)
		if err != nil {
			return nil, err
		}
		existing := make(map[string]bool, len(existingCodes))
		for _, c := range existingCodes {
			existing[c] = true
		}
		code = generateReimbursementTypeCode(req.Name, existing)
	}
	t := &ReimbursementType{
		Code:        code,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
	}
	if req.IsActive != nil {
		t.IsActive = *req.IsActive
	}
	if err := s.repo.CreateReimbursementType(ctx, t); err != nil {
		return nil, err
	}
	s.logger.Info("Reimbursement type created", zap.String("id", t.ID.String()), zap.String("name", t.Name))
	return reimbTypeToResponse(t), nil
}

func (s *Service) GetReimbursementTypeByID(ctx context.Context, id string) (*ReimbursementTypeResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	t, err := s.repo.FindReimbursementTypeByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return reimbTypeToResponse(t), nil
}

func (s *Service) ListReimbursementTypes(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	types, total, err := s.repo.ListReimbursementTypes(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]ReimbursementTypeResponse, 0, len(types))
	for _, t := range types {
		responses = append(responses, *reimbTypeToResponse(&t))
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

func (s *Service) UpdateReimbursementType(ctx context.Context, id string, req UpdateReimbursementTypeRequest) (*ReimbursementTypeResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	t, err := s.repo.FindReimbursementTypeByID(ctx, uid)
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
	if req.IsActive != nil {
		t.IsActive = *req.IsActive
	}
	if err := s.repo.UpdateReimbursementType(ctx, t); err != nil {
		return nil, err
	}
	return reimbTypeToResponse(t), nil
}

func (s *Service) DeleteReimbursementType(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteReimbursementType(ctx, uid)
}

// =========================================================================
// Reimbursement Requests
// =========================================================================

func (s *Service) CreateReimbursementRequest(ctx context.Context, req CreateReimbursementRequest) (*ReimbursementRequestResponse, error) {
	userID := authctx.GetUserID(ctx)
	if userID == nil {
		return nil, fmt.Errorf("user context not found")
	}
	// JWT user_id adalah UUID akun user, bukan UUID employee — resolve
	// employee milik user yang login lewat employee_accounts (pola sama
	// dengan /user-accounts/me). Tanpa ini, request tersimpan dengan
	// employee_id = user_id sehingga tidak pernah muncul di list
	// filter milik karyawan.
	empID, err := s.repo.FindEmployeeIDByUserID(ctx, *userID)
	if err != nil {
		return nil, err
	}
	if empID == nil {
		return nil, fmt.Errorf("user is not linked to an employee account")
	}
	rTypeID, err := uuid.Parse(req.RequestTypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid request_type_id: %w", err)
	}

	rr := &ReimbursementRequest{
		EmployeeID:    *empID,
		RequestTypeID: rTypeID,
		Title:         req.Title,
		Description:   req.Description,
		Currency:      "IDR",
		Status:        ReimbStatusDraft,
	}
	if req.Currency != "" {
		rr.Currency = req.Currency
	}
	if req.SupervisorID != nil {
		sID, err := uuid.Parse(*req.SupervisorID)
		if err == nil {
			rr.SupervisorID = &sID
		}
	}

	if err := s.repo.CreateReimbursementRequest(ctx, rr); err != nil {
		return nil, err
	}
	s.logger.Info("Reimbursement request created", zap.String("id", rr.ID.String()))
	return reimbRequestToResponse(rr), nil
}

func (s *Service) GetReimbursementRequestByID(ctx context.Context, id string) (*ReimbursementRequestResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	rr, err := s.repo.FindReimbursementRequestByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return reimbRequestToResponse(rr), nil
}

func (s *Service) ListReimbursementRequests(ctx context.Context, employeeID *string, status *string, page, perPage int) (*PaginatedResponse, error) {
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
	requests, total, err := s.repo.ListReimbursementRequests(ctx, empUUID, status, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]ReimbursementRequestResponse, 0, len(requests))
	for _, r := range requests {
		responses = append(responses, *reimbRequestToResponse(&r))
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

func (s *Service) UpdateReimbursementRequest(ctx context.Context, id string, req UpdateReimbursementRequest) (*ReimbursementRequestResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	rr, err := s.repo.FindReimbursementRequestByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if rr.Status != ReimbStatusDraft {
		return nil, fmt.Errorf("cannot update request in %s status", rr.Status)
	}
	if req.Title != nil {
		rr.Title = *req.Title
	}
	if req.Description != nil {
		rr.Description = *req.Description
	}
	if req.Currency != nil {
		rr.Currency = *req.Currency
	}
	if req.SupervisorID != nil {
		sID, err := uuid.Parse(*req.SupervisorID)
		if err == nil {
			rr.SupervisorID = &sID
		}
	}
	if err := s.repo.UpdateReimbursementRequest(ctx, rr); err != nil {
		return nil, err
	}
	return reimbRequestToResponse(rr), nil
}

func (s *Service) UpdateReimbursementRequestStatus(ctx context.Context, id, status, note string, amount *float64, flowID *string, paymentDetails ...*string) (*ReimbursementRequestResponse, error) {
	// paymentDetails is variadic to stay compatible with existing callers:
	// optional [payment_method, payment_reference, payment_note] recorded
	// directly in this module when status becomes PAID (no payroll linkage).
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	rr, err := s.repo.FindReimbursementRequestByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	switch ReimbursementStatus(status) {
	case ReimbStatusSubmitted:
		if rr.Status != ReimbStatusDraft {
			return nil, fmt.Errorf("only draft requests can be submitted")
		}
		// Auto-calculate total from items
		total, err := s.repo.SumReimbursementItems(ctx, uid)
		if err != nil {
			return nil, err
		}
		rr.TotalAmount = total
		rr.SubmittedAt = &now

		// Route through the central approval module when a flow is selected,
		// instead of relying on this module's own ad-hoc SupervisorID/HrID
		// fields. If no flow_id was supplied, auto-resolve the active flow
		// for this module (same pattern business travel uses) so a plain
		// submit still reaches the Approval inbox instead of silently
		// staying SUBMITTED forever.
		flowIDToUse := ""
		if flowID != nil && *flowID != "" {
			flowIDToUse = *flowID
		} else if s.approvalEngine != nil {
			if resolved, resolveErr := s.approvalEngine.GetActiveFlowIDForModule(ctx, "reimbursement"); resolveErr == nil {
				flowIDToUse = resolved
			}
		}
		if s.approvalEngine != nil && flowIDToUse != "" {
			instanceID, err := s.approvalEngine.CreateApprovalInstance(ctx, "reimbursement", rr.ID.String(), flowIDToUse)
			if err != nil {
				// Routing/assignee-resolution failures (approval.RoutingError)
				// mean the configured flow can't reach an approver — fail
				// loudly instead of silently continuing (same policy as KPI
				// target submission). Any other approval error stays
				// best-effort: log and continue.
				var re *approval.RoutingError
				if errors.As(err, &re) {
					return nil, err
				}
				s.logger.Warn("Failed to create approval instance for reimbursement request, continuing without approval",
					zap.String("reimbursement_request_id", rr.ID.String()),
					zap.Error(err),
				)
				rr.Status = ReimbStatusSubmitted
			} else {
				if parsedInstanceID, parseErr := uuid.Parse(instanceID); parseErr == nil {
					rr.ApprovalInstanceID = &parsedInstanceID
				}
				rr.Status = ReimbStatusPendingApproval
			}
		} else {
			rr.Status = ReimbStatusSubmitted
		}
	case ReimbStatusApproved:
		if rr.Status != ReimbStatusSubmitted && rr.Status != ReimbStatusPendingApproval {
			return nil, fmt.Errorf("only submitted requests can be approved")
		}
		rr.Status = ReimbStatusApproved
		rr.ApprovedAt = &now
	case ReimbStatusRejected:
		rr.Status = ReimbStatusRejected
		rr.RejectedAt = &now
	case ReimbStatusCancelled:
		rr.Status = ReimbStatusCancelled
		rr.CancelledAt = &now
	case ReimbStatusPaid:
		if rr.Status != ReimbStatusApproved {
			return nil, fmt.Errorf("only approved requests can be paid")
		}
		rr.Status = ReimbStatusPaid
		rr.PaidAt = &now
		if amount != nil {
			rr.PaidAmount = amount
		} else {
			rr.PaidAmount = &rr.TotalAmount
		}
		// Payment details recorded directly in this module (no payroll
		// linkage — product decision 2026-08-16).
		if len(paymentDetails) > 0 && paymentDetails[0] != nil {
			rr.PaymentMethod = paymentDetails[0]
		}
		if len(paymentDetails) > 1 && paymentDetails[1] != nil {
			rr.PaymentReference = paymentDetails[1]
		}
		if len(paymentDetails) > 2 && paymentDetails[2] != nil {
			rr.PaymentNote = paymentDetails[2]
		}
	default:
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	if note != "" {
		rr.HrNote = &note
	}

	if err := s.repo.UpdateReimbursementRequest(ctx, rr); err != nil {
		return nil, err
	}
	s.logger.Info("Reimbursement request status updated",
		zap.String("id", rr.ID.String()),
		zap.String("status", string(rr.Status)),
	)
	if rr.Status == ReimbStatusPaid {
		s.notifyReimbursementOutcome(ctx, rr, "REIMBURSEMENT_PAID")
	}
	return reimbRequestToResponse(rr), nil
}

// HandleApprovalStatusChange is invoked by the approval module's push-based
// status callback when a reimbursement request's approval instance reaches
// a final state, so the request's own status field updates itself.
func (s *Service) HandleApprovalStatusChange(ctx context.Context, documentID uuid.UUID, status string, note string) error {
	rr, err := s.repo.FindReimbursementRequestByID(ctx, documentID)
	if err != nil {
		return err
	}
	if rr.Status != ReimbStatusPendingApproval {
		return nil
	}

	now := time.Now()
	switch status {
	case "APPROVED":
		rr.Status = ReimbStatusApproved
		rr.ApprovedAt = &now
	case "REJECTED":
		rr.Status = ReimbStatusRejected
		rr.RejectedAt = &now
		if note != "" {
			rr.SupervisorNote = &note
		}
	case "CANCELLED":
		rr.Status = ReimbStatusCancelled
		rr.CancelledAt = &now
	default:
		return nil
	}

	s.logger.Info("Reimbursement request status updated via approval status handler",
		zap.String("reimbursement_request_id", rr.ID.String()),
		zap.String("approval_status", status),
	)
	if err := s.repo.UpdateReimbursementRequest(ctx, rr); err != nil {
		return err
	}
	switch rr.Status {
	case ReimbStatusApproved:
		s.notifyReimbursementOutcome(ctx, rr, "REIMBURSEMENT_APPROVED")
	case ReimbStatusRejected:
		s.notifyReimbursementOutcome(ctx, rr, "REIMBURSEMENT_REJECTED")
	}
	return nil
}

// notifyReimbursementOutcome notifies the requesting employee of their
// reimbursement request's final outcome (approved/rejected/paid).
// Best-effort: if the notifier isn't wired, the employee has no linked user
// account, or Notify fails, this logs and moves on rather than failing the
// status transition itself.
func (s *Service) notifyReimbursementOutcome(ctx context.Context, rr *ReimbursementRequest, notifType string) {
	if s.notifier == nil {
		return
	}
	userID, err := s.repo.FindUserIDByEmployeeID(ctx, rr.EmployeeID)
	if err != nil {
		s.logger.Warn("Failed to resolve employee user id for reimbursement notification",
			zap.String("reimbursement_request_id", rr.ID.String()),
			zap.Error(err),
		)
		return
	}
	if userID == nil {
		return
	}
	if err := s.notifier.Notify(ctx, *userID, notifType, nil, "reimbursement", rr.ID); err != nil {
		s.logger.Warn("Failed to send reimbursement notification",
			zap.String("reimbursement_request_id", rr.ID.String()),
			zap.String("notif_type", notifType),
			zap.Error(err),
		)
	}
}

func (s *Service) DeleteReimbursementRequest(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteReimbursementRequest(ctx, uid)
}

// =========================================================================
// Reimbursement Items
// =========================================================================

func (s *Service) CreateReimbursementItem(ctx context.Context, requestID string, req CreateReimbursementItemRequest) (*ReimbursementItemResponse, error) {
	rID, err := uuid.Parse(requestID)
	if err != nil {
		return nil, fmt.Errorf("invalid request_id: %w", err)
	}
	// Verify request exists & is in draft
	rr, err := s.repo.FindReimbursementRequestByID(ctx, rID)
	if err != nil {
		return nil, err
	}
	if rr.Status != ReimbStatusDraft {
		return nil, fmt.Errorf("can only add items to draft requests")
	}

	item := &ReimbursementItem{
		ReimbursementRequestID: rID,
		ExpenseDate:            req.ExpenseDate,
		ExpenseType:            req.ExpenseType,
		Description:            req.Description,
		Amount:                 req.Amount,
		ReceiptURL:             req.ReceiptURL,
	}
	if err := s.repo.CreateReimbursementItem(ctx, item); err != nil {
		return nil, err
	}
	return reimbItemToResponse(item), nil
}

func (s *Service) ListReimbursementItems(ctx context.Context, requestID string) ([]ReimbursementItemResponse, error) {
	rID, err := uuid.Parse(requestID)
	if err != nil {
		return nil, fmt.Errorf("invalid request_id: %w", err)
	}
	items, err := s.repo.ListReimbursementItems(ctx, rID)
	if err != nil {
		return nil, err
	}
	responses := make([]ReimbursementItemResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, *reimbItemToResponse(&item))
	}
	return responses, nil
}

func (s *Service) UpdateReimbursementItem(ctx context.Context, requestID, itemID string, req UpdateReimbursementItemRequest) (*ReimbursementItemResponse, error) {
	rID, err := uuid.Parse(requestID)
	if err != nil {
		return nil, fmt.Errorf("invalid request_id: %w", err)
	}
	iID, err := uuid.Parse(itemID)
	if err != nil {
		return nil, fmt.Errorf("invalid item_id: %w", err)
	}
	// Verify request is in draft
	rr, err := s.repo.FindReimbursementRequestByID(ctx, rID)
	if err != nil {
		return nil, err
	}
	if rr.Status != ReimbStatusDraft {
		return nil, fmt.Errorf("can only update items in draft requests")
	}

	item, err := s.repo.FindReimbursementItemByID(ctx, iID)
	if err != nil {
		return nil, err
	}
	if item.ReimbursementRequestID != rID {
		return nil, fmt.Errorf("item does not belong to this request")
	}

	if req.ExpenseDate != nil {
		item.ExpenseDate = *req.ExpenseDate
	}
	if req.ExpenseType != nil {
		item.ExpenseType = *req.ExpenseType
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if req.Amount != nil {
		item.Amount = *req.Amount
	}
	if req.ReceiptURL != nil {
		item.ReceiptURL = req.ReceiptURL
	}
	if err := s.repo.UpdateReimbursementItem(ctx, item); err != nil {
		return nil, err
	}
	return reimbItemToResponse(item), nil
}

func (s *Service) DeleteReimbursementItem(ctx context.Context, requestID, itemID string) error {
	rID, err := uuid.Parse(requestID)
	if err != nil {
		return fmt.Errorf("invalid request_id: %w", err)
	}
	iID, err := uuid.Parse(itemID)
	if err != nil {
		return fmt.Errorf("invalid item_id: %w", err)
	}
	// Verify request is in draft
	rr, err := s.repo.FindReimbursementRequestByID(ctx, rID)
	if err != nil {
		return err
	}
	if rr.Status != ReimbStatusDraft {
		return fmt.Errorf("can only delete items from draft requests")
	}
	item, err := s.repo.FindReimbursementItemByID(ctx, iID)
	if err != nil {
		return err
	}
	if item.ReimbursementRequestID != rID {
		return fmt.Errorf("item does not belong to this request")
	}
	return s.repo.DeleteReimbursementItem(ctx, iID)
}

// generateReimbursementTypeCode slugifies a display name into an
// UPPER_SNAKE_CASE code (e.g. "Transport & Travel" -> "TRANSPORT_TRAVEL") and
// appends a numeric suffix if the base code already exists, so the type form
// never needs to ask the user for a code.
func generateReimbursementTypeCode(name string, existing map[string]bool) string {
	var b strings.Builder
	lastUnderscore := false
	for _, r := range strings.ToUpper(strings.TrimSpace(name)) {
		switch {
		case r >= 'A' && r <= 'Z' || r >= '0' && r <= '9':
			b.WriteRune(r)
			lastUnderscore = false
		default:
			if !lastUnderscore && b.Len() > 0 {
				b.WriteRune('_')
				lastUnderscore = true
			}
		}
	}
	base := strings.TrimSuffix(b.String(), "_")
	if base == "" {
		base = "OTHER"
	}
	code := base
	for i := 2; existing[code]; i++ {
		code = fmt.Sprintf("%s_%d", base, i)
	}
	return code
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

// =========================================================================
// Helpers
// =========================================================================

// =========================================================================
// Response converters
// =========================================================================

func reimbTypeToResponse(t *ReimbursementType) *ReimbursementTypeResponse {
	return &ReimbursementTypeResponse{
		ID:          t.ID.String(),
		Code:        t.Code,
		Name:        t.Name,
		Description: t.Description,
		IsActive:    t.IsActive,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func reimbRequestToResponse(r *ReimbursementRequest) *ReimbursementRequestResponse {
	resp := &ReimbursementRequestResponse{
		ID:                 r.ID.String(),
		EmployeeID:         r.EmployeeID.String(),
		RequestTypeID:      r.RequestTypeID.String(),
		Title:              r.Title,
		Description:        r.Description,
		TotalAmount:        r.TotalAmount,
		Currency:           r.Currency,
		Status:             string(r.Status),
		SupervisorNote:     r.SupervisorNote,
		SupervisorActionAt: r.SupervisorActionAt,
		HrNote:             r.HrNote,
		HrActionAt:         r.HrActionAt,
		PaidAt:             r.PaidAt,
		PaidAmount:         r.PaidAmount,
		PaymentMethod:      r.PaymentMethod,
		PaymentReference:   r.PaymentReference,
		PaymentNote:        r.PaymentNote,
		SubmittedAt:        r.SubmittedAt,
		ApprovedAt:         r.ApprovedAt,
		RejectedAt:         r.RejectedAt,
		CancelledAt:        r.CancelledAt,
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
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

func reimbItemToResponse(i *ReimbursementItem) *ReimbursementItemResponse {
	return &ReimbursementItemResponse{
		ID:                      i.ID.String(),
		ReimbursementRequestID:  i.ReimbursementRequestID.String(),
		ExpenseDate:             i.ExpenseDate,
		ExpenseType:             i.ExpenseType,
		Description:             i.Description,
		Amount:                  i.Amount,
		ReceiptURL:              i.ReceiptURL,
		CreatedAt:               i.CreatedAt,
		UpdatedAt:               i.UpdatedAt,
	}
}
