package reimbursement

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

type Service struct {
	repo   *Repository
	logger *zap.Logger
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// =========================================================================
// Reimbursement Types
// =========================================================================

func (s *Service) CreateReimbursementType(ctx context.Context, req CreateReimbursementTypeRequest) (*ReimbursementTypeResponse, error) {
	t := &ReimbursementType{
		Code:        req.Code,
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
	empID, ok := ctx.Value("user_id").(string)
	if !ok || empID == "" {
		return nil, fmt.Errorf("user context not found")
	}
	employeeUUID, err := uuid.Parse(empID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id from context: %w", err)
	}
	rTypeID, err := uuid.Parse(req.RequestTypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid request_type_id: %w", err)
	}

	rr := &ReimbursementRequest{
		EmployeeID:    employeeUUID,
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

func (s *Service) UpdateReimbursementRequestStatus(ctx context.Context, id, status, note string, amount *float64) (*ReimbursementRequestResponse, error) {
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
		rr.Status = ReimbStatusSubmitted
		rr.SubmittedAt = now.UnixNano()
	case ReimbStatusApproved:
		if rr.Status != ReimbStatusSubmitted {
			return nil, fmt.Errorf("only submitted requests can be approved")
		}
		rr.Status = ReimbStatusApproved
		rr.ApprovedAt = now.UnixNano()
	case ReimbStatusRejected:
		rr.Status = ReimbStatusRejected
		rr.RejectedAt = now.UnixNano()
	case ReimbStatusCancelled:
		rr.Status = ReimbStatusCancelled
		rr.CancelledAt = now.UnixNano()
	case ReimbStatusPaid:
		if rr.Status != ReimbStatusApproved {
			return nil, fmt.Errorf("only approved requests can be paid")
		}
		rr.Status = ReimbStatusPaid
		rr.PaidAt = now.UnixNano()
		if amount != nil {
			rr.PaidAmount = amount
		} else {
			rr.PaidAmount = &rr.TotalAmount
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
	return reimbRequestToResponse(rr), nil
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

// unixNanoToTimePtr converts Unix nanoseconds to *time.Time.
// Returns nil if the value is 0.
func unixNanoToTimePtr(ns int64) *time.Time {
	if ns == 0 {
		return nil
	}
	t := time.Unix(0, ns)
	return &t
}

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
		ID:                r.ID.String(),
		EmployeeID:        r.EmployeeID.String(),
		RequestTypeID:     r.RequestTypeID.String(),
		Title:             r.Title,
		Description:       r.Description,
		TotalAmount:       r.TotalAmount,
		Currency:          r.Currency,
		Status:            string(r.Status),
		SupervisorNote:    r.SupervisorNote,
		SupervisorActionAt: unixNanoToTimePtr(r.SupervisorActionAt),
		HrNote:            r.HrNote,
		HrActionAt:        unixNanoToTimePtr(r.HrActionAt),
		PaidAt:            unixNanoToTimePtr(r.PaidAt),
		PaidAmount:        r.PaidAmount,
		SubmittedAt:       unixNanoToTimePtr(r.SubmittedAt),
		ApprovedAt:        unixNanoToTimePtr(r.ApprovedAt),
		RejectedAt:        unixNanoToTimePtr(r.RejectedAt),
		CancelledAt:       unixNanoToTimePtr(r.CancelledAt),
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
	}
	if r.SupervisorID != nil {
		v := r.SupervisorID.String()
		resp.SupervisorID = &v
	}
	if r.HrID != nil {
		v := r.HrID.String()
		resp.HrID = &v
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
