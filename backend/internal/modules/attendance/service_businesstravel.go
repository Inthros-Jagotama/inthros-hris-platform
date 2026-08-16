package attendance

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/modules/approval"
	"github.com/inthros/hris-platform/internal/pkg/authctx"
)

// businessTravelApprovalModule is the module slug registered with the
// central Approval module for the Travel Approval flow. Kept distinct from
// "attendance" (used by overtime/correction) so both flows can be routed
// and configured independently, per
// docs/module-attendance-business-travel-development-plan.md §54.3.
const businessTravelApprovalModule = "business_travel"

// ErrBusinessTravelInvalidState: aksi tidak valid untuk status travel saat ini
// (mis. submit travel yang bukan DRAFT, atau update travel yang sudah SUBMITTED).
var ErrBusinessTravelInvalidState = errors.New("business travel is not in a valid state for this action")

func (s *Service) CreateBusinessTravel(ctx context.Context, req CreateBusinessTravelRequest) (*BusinessTravelResponse, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date: %w", err)
	}
	if endDate.Before(startDate) {
		return nil, fmt.Errorf("end_date must not be before start_date")
	}

	requesterID := authctx.GetUserID(ctx)
	if requesterID == nil {
		return nil, fmt.Errorf("unable to resolve requester from context")
	}

	travel := &BusinessTravel{
		RequestNumber:  generateBusinessTravelRequestNumber(),
		RequesterID:    *requesterID,
		Title:          req.Title,
		StartDate:      startDate,
		EndDate:        endDate,
		Status:         TravelStatusDraft,
		ApprovalStatus: string(TravelStatusDraft),
		CreatedBy:      requesterID,
	}
	if req.Purpose != "" {
		travel.Purpose = &req.Purpose
	}
	if req.Description != "" {
		travel.Description = &req.Description
	}
	if req.Origin != "" {
		travel.Origin = &req.Origin
	}

	if err := s.repo.CreateBusinessTravel(ctx, travel); err != nil {
		return nil, err
	}

	for _, p := range req.Participants {
		participant := &BusinessTravelParticipant{
			BusinessTravelID: travel.ID,
			ParticipantType:  ParticipantType(strings.ToUpper(p.ParticipantType)),
			Role:             ParticipantRole(strings.ToUpper(p.Role)),
		}
		if p.EmployeeID != "" {
			if empID, err := uuid.Parse(p.EmployeeID); err == nil {
				participant.EmployeeID = &empID
			}
		}
		if p.Name != "" {
			participant.Name = &p.Name
		}
		if p.Organization != "" {
			participant.Organization = &p.Organization
		}
		if p.Position != "" {
			participant.Position = &p.Position
		}
		if p.IdentityNumber != "" {
			participant.IdentityNumber = &p.IdentityNumber
		}
		if p.Email != "" {
			participant.Email = &p.Email
		}
		if p.Phone != "" {
			participant.Phone = &p.Phone
		}
		if p.Notes != "" {
			participant.Notes = &p.Notes
		}
		if participant.Role == "" {
			participant.Role = ParticipantRoleMember
		}
		if err := s.repo.CreateParticipant(ctx, participant); err != nil {
			return nil, err
		}
	}

	for _, d := range req.Destinations {
		destination := &BusinessTravelDestination{
			BusinessTravelID: travel.ID,
			Sequence:         d.Sequence,
		}
		if d.Country != "" {
			destination.Country = &d.Country
		}
		if d.Province != "" {
			destination.Province = &d.Province
		}
		if d.City != "" {
			destination.City = &d.City
		}
		if d.Location != "" {
			destination.Location = &d.Location
		}
		if d.Purpose != "" {
			destination.Purpose = &d.Purpose
		}
		if d.Notes != "" {
			destination.Notes = &d.Notes
		}
		if d.ArrivalDate != "" {
			if parsed, err := time.Parse("2006-01-02", d.ArrivalDate); err == nil {
				destination.ArrivalDate = &parsed
			}
		}
		if d.DepartureDate != "" {
			if parsed, err := time.Parse("2006-01-02", d.DepartureDate); err == nil {
				destination.DepartureDate = &parsed
			}
		}
		if err := s.repo.CreateDestination(ctx, destination); err != nil {
			return nil, err
		}
	}

	return s.GetBusinessTravelByID(ctx, travel.ID.String())
}

func (s *Service) GetBusinessTravelByID(ctx context.Context, id string) (*BusinessTravelResponse, error) {
	travelID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	travel, err := s.repo.FindBusinessTravelByID(ctx, travelID)
	if err != nil {
		return nil, err
	}
	return businessTravelToResponse(travel), nil
}

func (s *Service) ListBusinessTravels(ctx context.Context, requesterID string, status string, page, perPage int) (*PaginatedResponse, error) {
	var requesterUUID *uuid.UUID
	if requesterID != "" {
		if parsed, err := uuid.Parse(requesterID); err == nil {
			requesterUUID = &parsed
		}
	}
	travels, total, err := s.repo.ListBusinessTravels(ctx, requesterUUID, status, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]*BusinessTravelResponse, 0, len(travels))
	for i := range travels {
		responses = append(responses, businessTravelToResponse(&travels[i]))
	}
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) UpdateBusinessTravel(ctx context.Context, id string, req UpdateBusinessTravelRequest) (*BusinessTravelResponse, error) {
	travelID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	travel, err := s.repo.FindBusinessTravelByID(ctx, travelID)
	if err != nil {
		return nil, err
	}
	if travel.Status != TravelStatusDraft {
		return nil, ErrBusinessTravelInvalidState
	}
	if req.Title != nil {
		travel.Title = *req.Title
	}
	if req.Purpose != nil {
		travel.Purpose = req.Purpose
	}
	if req.Description != nil {
		travel.Description = req.Description
	}
	if req.Origin != nil {
		travel.Origin = req.Origin
	}
	if req.StartDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date: %w", err)
		}
		travel.StartDate = parsed
	}
	if req.EndDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date: %w", err)
		}
		travel.EndDate = parsed
	}
	if err := s.repo.UpdateBusinessTravel(ctx, travel); err != nil {
		return nil, err
	}
	return businessTravelToResponse(travel), nil
}

// SubmitBusinessTravel routes a DRAFT travel request through the central
// Approval module (Rule 6, §52: no bespoke approval engine). Mirrors
// Service.CreateOvertimeRequest's approval-instance wiring.
func (s *Service) SubmitBusinessTravel(ctx context.Context, id string, req SubmitBusinessTravelRequest) (*BusinessTravelResponse, error) {
	travelID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	travel, err := s.repo.FindBusinessTravelByID(ctx, travelID)
	if err != nil {
		return nil, err
	}
	if travel.Status != TravelStatusDraft {
		return nil, ErrBusinessTravelInvalidState
	}

	travel.Status = TravelStatusSubmitted
	travel.ApprovalStatus = string(TravelStatusSubmitted)

	if s.approvalEngine != nil {
		flowID := ""
		if req.FlowID != nil && *req.FlowID != "" {
			flowID = *req.FlowID
		} else if resolved, err := s.approvalEngine.GetActiveFlowIDForModule(ctx, businessTravelApprovalModule); err == nil {
			flowID = resolved
		}
		if flowID != "" {
			instanceID, err := s.approvalEngine.CreateApprovalInstance(ctx, businessTravelApprovalModule, travel.ID.String(), flowID)
			if err != nil {
				var re *approval.RoutingError
				if errors.As(err, &re) {
					return nil, err
				}
				s.logger.Warn("Failed to create approval instance for business travel, continuing without approval",
					zap.String("business_travel_id", travel.ID.String()),
					zap.Error(err),
				)
			} else if parsedInstanceID, parseErr := uuid.Parse(instanceID); parseErr == nil {
				travel.ApprovalInstanceID = &parsedInstanceID
			}
		}
	}

	if err := s.repo.UpdateBusinessTravel(ctx, travel); err != nil {
		return nil, err
	}
	return businessTravelToResponse(travel), nil
}

func (s *Service) CancelBusinessTravel(ctx context.Context, id string) (*BusinessTravelResponse, error) {
	travelID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	travel, err := s.repo.FindBusinessTravelByID(ctx, travelID)
	if err != nil {
		return nil, err
	}
	if travel.Status == TravelStatusClosed || travel.Status == TravelStatusCancelled {
		return nil, ErrBusinessTravelInvalidState
	}
	if s.approvalEngine != nil && travel.ApprovalInstanceID != nil {
		if cerr := s.approvalEngine.CancelApprovalInstance(ctx, travel.ApprovalInstanceID.String()); cerr != nil {
			s.logger.Warn("Failed to cancel approval instance for business travel",
				zap.String("business_travel_id", travel.ID.String()),
				zap.Error(cerr),
			)
		}
	}
	travel.Status = TravelStatusCancelled
	travel.ApprovalStatus = string(TravelStatusCancelled)
	if err := s.repo.UpdateBusinessTravel(ctx, travel); err != nil {
		return nil, err
	}
	return businessTravelToResponse(travel), nil
}

// AddBusinessTravelActivity menambahkan satu item agenda/kegiatan ke travel
// yang sudah ada. Tidak dibatasi status DRAFT — agenda bisa disesuaikan
// selama perjalanan berlangsung (§9 plan doc tidak menyebut pembatasan status).
func (s *Service) AddBusinessTravelActivity(ctx context.Context, travelIDStr string, req CreateBusinessTravelActivityRequest) (*BusinessTravelActivityResponse, error) {
	travelID, err := uuid.Parse(travelIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	if _, err := s.repo.FindBusinessTravelByIDForOwnership(ctx, travelID); err != nil {
		return nil, err
	}
	activityDate, err := time.Parse("2006-01-02", req.ActivityDate)
	if err != nil {
		return nil, fmt.Errorf("invalid activity_date: %w", err)
	}
	activity := &BusinessTravelActivity{
		BusinessTravelID: travelID,
		ActivityDate:     activityDate,
		Title:            req.Title,
	}
	if req.StartTime != "" {
		activity.StartTime = &req.StartTime
	}
	if req.EndTime != "" {
		activity.EndTime = &req.EndTime
	}
	if req.Description != "" {
		activity.Description = &req.Description
	}
	if req.Location != "" {
		activity.Location = &req.Location
	}
	if req.Organizer != "" {
		activity.Organizer = &req.Organizer
	}
	if req.Notes != "" {
		activity.Notes = &req.Notes
	}
	if err := s.repo.CreateActivity(ctx, activity); err != nil {
		return nil, err
	}
	resp := businessTravelActivityToResponse(activity)
	return &resp, nil
}

func (s *Service) ListBusinessTravelActivities(ctx context.Context, travelIDStr string) ([]BusinessTravelActivityResponse, error) {
	travelID, err := uuid.Parse(travelIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	activities, err := s.repo.ListActivitiesByTravel(ctx, travelID)
	if err != nil {
		return nil, err
	}
	responses := make([]BusinessTravelActivityResponse, 0, len(activities))
	for i := range activities {
		responses = append(responses, businessTravelActivityToResponse(&activities[i]))
	}
	return responses, nil
}

// AddBusinessTravelSchedule menambahkan satu jadwal/transportasi ke travel
// yang sudah ada.
func (s *Service) AddBusinessTravelSchedule(ctx context.Context, travelIDStr string, req CreateBusinessTravelScheduleRequest) (*BusinessTravelScheduleResponse, error) {
	travelID, err := uuid.Parse(travelIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	if _, err := s.repo.FindBusinessTravelByIDForOwnership(ctx, travelID); err != nil {
		return nil, err
	}
	schedule := &BusinessTravelSchedule{
		BusinessTravelID: travelID,
		ScheduleType:     ScheduleType(strings.ToUpper(req.ScheduleType)),
	}
	if req.TransportationType != "" {
		schedule.TransportationType = TransportationType(strings.ToUpper(req.TransportationType))
	} else {
		schedule.TransportationType = TransportationOther
	}
	if req.Origin != "" {
		schedule.Origin = &req.Origin
	}
	if req.Destination != "" {
		schedule.Destination = &req.Destination
	}
	if req.Provider != "" {
		schedule.Provider = &req.Provider
	}
	if req.BookingReference != "" {
		schedule.BookingReference = &req.BookingReference
	}
	if req.Notes != "" {
		schedule.Notes = &req.Notes
	}
	if req.DepartureDatetime != "" {
		if parsed, err := time.Parse(time.RFC3339, req.DepartureDatetime); err == nil {
			schedule.DepartureDatetime = &parsed
		}
	}
	if req.ArrivalDatetime != "" {
		if parsed, err := time.Parse(time.RFC3339, req.ArrivalDatetime); err == nil {
			schedule.ArrivalDatetime = &parsed
		}
	}
	if err := s.repo.CreateSchedule(ctx, schedule); err != nil {
		return nil, err
	}
	resp := businessTravelScheduleToResponse(schedule)
	return &resp, nil
}

func (s *Service) ListBusinessTravelSchedules(ctx context.Context, travelIDStr string) ([]BusinessTravelScheduleResponse, error) {
	travelID, err := uuid.Parse(travelIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	schedules, err := s.repo.ListSchedulesByTravel(ctx, travelID)
	if err != nil {
		return nil, err
	}
	responses := make([]BusinessTravelScheduleResponse, 0, len(schedules))
	for i := range schedules {
		responses = append(responses, businessTravelScheduleToResponse(&schedules[i]))
	}
	return responses, nil
}

// =========================================================================
// Funding (§14-17 plan doc)
// =========================================================================

// ErrBusinessTravelNotApproved: funding hanya boleh dibuat setelah travel
// APPROVED — request tidak menentukan pembiayaan (Rule 1, §52).
var ErrBusinessTravelNotApproved = errors.New("business travel must be APPROVED before funding can be created")

// ErrFundingInvalidState: aksi tidak valid untuk status funding saat ini
// (mis. confirm funding yang sudah FUNDED/CANCELLED).
var ErrFundingInvalidState = errors.New("funding is not in a valid state for this action")

func (s *Service) CreateFundingMethod(ctx context.Context, req CreateFundingMethodRequest) (*FundingMethodResponse, error) {
	method := &FundingMethod{
		Code:   strings.ToUpper(req.Code),
		Name:   req.Name,
		Active: true,
	}
	if req.Description != "" {
		method.Description = &req.Description
	}
	if err := s.repo.CreateFundingMethod(ctx, method); err != nil {
		return nil, err
	}
	return fundingMethodToResponse(method), nil
}

func (s *Service) ListFundingMethods(ctx context.Context, activeOnly bool) ([]FundingMethodResponse, error) {
	methods, err := s.repo.ListFundingMethods(ctx, activeOnly)
	if err != nil {
		return nil, err
	}
	responses := make([]FundingMethodResponse, 0, len(methods))
	for i := range methods {
		responses = append(responses, *fundingMethodToResponse(&methods[i]))
	}
	return responses, nil
}

// CreateFunding mencatat pembiayaan untuk travel yang sudah APPROVED. FundedBy
// diisi dari user yang login (bisa Finance/Admin, BUKAN otomatis requester —
// Rule 2, §52). Amount/funding method di sini independen dari estimasi biaya
// yang diisi saat request (§11: estimasi bukan funding method/komitmen).
func (s *Service) CreateFunding(ctx context.Context, travelIDStr string, req CreateFundingRequest) (*FundingResponse, error) {
	travelID, err := uuid.Parse(travelIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	travel, err := s.repo.FindBusinessTravelByIDForOwnership(ctx, travelID)
	if err != nil {
		return nil, err
	}
	if travel.Status != TravelStatusApproved && travel.Status != TravelStatusInProgress && travel.Status != TravelStatusCompleted {
		return nil, ErrBusinessTravelNotApproved
	}
	fundingMethodID, err := uuid.Parse(req.FundingMethodID)
	if err != nil {
		return nil, fmt.Errorf("invalid funding_method_id: %w", err)
	}
	if _, err := s.repo.FindFundingMethodByID(ctx, fundingMethodID); err != nil {
		return nil, err
	}

	funding := &Funding{
		BusinessTravelID: travelID,
		FundingMethodID:  fundingMethodID,
		Amount:           req.Amount,
		Status:           FundingStatusPending,
		FundedBy:         authctx.GetUserID(ctx),
	}
	if req.ParticipantID != "" {
		if participantID, err := uuid.Parse(req.ParticipantID); err == nil {
			funding.ParticipantID = &participantID
		}
	}
	if req.PaymentMethod != "" {
		funding.PaymentMethod = &req.PaymentMethod
	}
	if req.PaymentReference != "" {
		funding.PaymentReference = &req.PaymentReference
	}
	if req.Notes != "" {
		funding.Notes = &req.Notes
	}
	if req.FundingDate != "" {
		if parsed, err := time.Parse("2006-01-02", req.FundingDate); err == nil {
			funding.FundingDate = &parsed
		}
	}

	if err := s.repo.CreateFunding(ctx, funding); err != nil {
		return nil, err
	}
	return fundingToResponse(funding), nil
}

func (s *Service) ListFundingsByTravel(ctx context.Context, travelIDStr string) ([]FundingResponse, error) {
	travelID, err := uuid.Parse(travelIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	fundings, err := s.repo.ListFundingsByTravel(ctx, travelID)
	if err != nil {
		return nil, err
	}
	responses := make([]FundingResponse, 0, len(fundings))
	for i := range fundings {
		responses = append(responses, *fundingToResponse(&fundings[i]))
	}
	return responses, nil
}

func (s *Service) UpdateFunding(ctx context.Context, fundingIDStr string, req UpdateFundingRequest) (*FundingResponse, error) {
	fundingID, err := uuid.Parse(fundingIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid funding id: %w", err)
	}
	funding, err := s.repo.FindFundingByID(ctx, fundingID)
	if err != nil {
		return nil, err
	}
	if funding.Status == FundingStatusFunded || funding.Status == FundingStatusCancelled || funding.Status == FundingStatusReversed {
		return nil, ErrFundingInvalidState
	}
	if req.Amount != nil {
		funding.Amount = *req.Amount
	}
	if req.PaymentMethod != nil {
		funding.PaymentMethod = req.PaymentMethod
	}
	if req.PaymentReference != nil {
		funding.PaymentReference = req.PaymentReference
	}
	if req.Notes != nil {
		funding.Notes = req.Notes
	}
	if req.FundingDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.FundingDate)
		if err != nil {
			return nil, fmt.Errorf("invalid funding_date: %w", err)
		}
		funding.FundingDate = &parsed
	}
	if err := s.repo.UpdateFunding(ctx, funding); err != nil {
		return nil, err
	}
	return fundingToResponse(funding), nil
}

// ConfirmFunding menandai funding sebagai FUNDED (§17: dana benar-benar
// sudah ditransfer/diberikan). Aksi terpisah dari Create supaya funding bisa
// dicatat dulu (PENDING/PROCESSING) sebelum benar-benar cair.
func (s *Service) ConfirmFunding(ctx context.Context, fundingIDStr string) (*FundingResponse, error) {
	fundingID, err := uuid.Parse(fundingIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid funding id: %w", err)
	}
	funding, err := s.repo.FindFundingByID(ctx, fundingID)
	if err != nil {
		return nil, err
	}
	if funding.Status == FundingStatusFunded || funding.Status == FundingStatusCancelled || funding.Status == FundingStatusReversed {
		return nil, ErrFundingInvalidState
	}
	funding.Status = FundingStatusFunded
	if err := s.repo.UpdateFunding(ctx, funding); err != nil {
		return nil, err
	}
	return fundingToResponse(funding), nil
}

// AddFundingDocument attaches proof of transfer to a funding. The URL is
// obtained by the client from the generic upload endpoint
// (POST /api/v1/tenant/uploads) beforehand — this module never handles raw
// file bytes itself, per §54.4.
func (s *Service) AddFundingDocument(ctx context.Context, fundingIDStr string, req AddFundingDocumentRequest) (*FundingResponse, error) {
	fundingID, err := uuid.Parse(fundingIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid funding id: %w", err)
	}
	if _, err := s.repo.FindFundingByID(ctx, fundingID); err != nil {
		return nil, err
	}
	doc := &FundingDocument{
		BusinessTravelFundingID: fundingID,
		DocumentType:            FundingDocumentType(strings.ToUpper(req.DocumentType)),
		FileName:                req.FileName,
		FilePath:                req.FilePath,
		UploadedBy:              authctx.GetUserID(ctx),
	}
	if req.MimeType != "" {
		doc.MimeType = &req.MimeType
	}
	if req.FileSize > 0 {
		doc.FileSize = &req.FileSize
	}
	now := time.Now()
	doc.UploadedAt = &now
	if err := s.repo.CreateFundingDocument(ctx, doc); err != nil {
		return nil, err
	}
	funding, err := s.repo.FindFundingByID(ctx, fundingID)
	if err != nil {
		return nil, err
	}
	return fundingToResponse(funding), nil
}

// =========================================================================
// Expense Category (master) & Actual Expense (§12, §21 plan doc)
// =========================================================================

// ErrExpenseInvalidState: aksi tidak valid untuk status expense saat ini.
var ErrExpenseInvalidState = errors.New("expense is not in a valid state for this action")

func (s *Service) CreateExpenseCategory(ctx context.Context, req CreateExpenseCategoryRequest) (*ExpenseCategoryResponse, error) {
	category := &ExpenseCategory{
		Code:            strings.ToUpper(req.Code),
		Name:            req.Name,
		RequiresReceipt: true,
		Reimbursable:    true,
		Active:          true,
	}
	if req.RequiresReceipt != nil {
		category.RequiresReceipt = *req.RequiresReceipt
	}
	if req.Reimbursable != nil {
		category.Reimbursable = *req.Reimbursable
	}
	if req.Description != "" {
		category.Description = &req.Description
	}
	if req.PayrollTreatment != "" {
		category.PayrollTreatment = &req.PayrollTreatment
	}
	if req.AccountCode != "" {
		category.AccountCode = &req.AccountCode
	}
	if err := s.repo.CreateExpenseCategory(ctx, category); err != nil {
		return nil, err
	}
	return expenseCategoryToResponse(category), nil
}

func (s *Service) ListExpenseCategories(ctx context.Context, activeOnly bool) ([]ExpenseCategoryResponse, error) {
	categories, err := s.repo.ListExpenseCategories(ctx, activeOnly)
	if err != nil {
		return nil, err
	}
	responses := make([]ExpenseCategoryResponse, 0, len(categories))
	for i := range categories {
		responses = append(responses, *expenseCategoryToResponse(&categories[i]))
	}
	return responses, nil
}

// CreateExpense mencatat actual expense pasca perjalanan (§21). Funding
// method boleh berbeda dari funding awal travel — mixed funding per item
// (§33) ditangani lewat FundingMethodID opsional per expense, bukan
// diwariskan otomatis dari travel.
func (s *Service) CreateExpense(ctx context.Context, travelIDStr string, req CreateExpenseRequest) (*ExpenseResponse, error) {
	travelID, err := uuid.Parse(travelIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	travel, err := s.repo.FindBusinessTravelByIDForOwnership(ctx, travelID)
	if err != nil {
		return nil, err
	}
	if travel.Status != TravelStatusApproved && travel.Status != TravelStatusInProgress && travel.Status != TravelStatusCompleted {
		return nil, ErrBusinessTravelNotApproved
	}
	categoryID, err := uuid.Parse(req.ExpenseCategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid expense_category_id: %w", err)
	}
	if _, err := s.repo.FindExpenseCategoryByID(ctx, categoryID); err != nil {
		return nil, err
	}
	expenseDate, err := time.Parse("2006-01-02", req.ExpenseDate)
	if err != nil {
		return nil, fmt.Errorf("invalid expense_date: %w", err)
	}

	expense := &Expense{
		BusinessTravelID:  travelID,
		ExpenseCategoryID: categoryID,
		ExpenseDate:       expenseDate,
		Quantity:          1,
		Amount:            req.Amount,
		Status:            ExpenseStatusDraft,
	}
	if req.Quantity > 0 {
		expense.Quantity = req.Quantity
	}
	if req.ParticipantID != "" {
		if participantID, err := uuid.Parse(req.ParticipantID); err == nil {
			expense.ParticipantID = &participantID
		}
	}
	if req.FundingMethodID != "" {
		if fundingMethodID, err := uuid.Parse(req.FundingMethodID); err == nil {
			expense.FundingMethodID = &fundingMethodID
		}
	}
	if req.Description != "" {
		expense.Description = &req.Description
	}
	if req.Unit != "" {
		expense.Unit = &req.Unit
	}
	if req.Vendor != "" {
		expense.Vendor = &req.Vendor
	}
	if req.ReceiptNumber != "" {
		expense.ReceiptNumber = &req.ReceiptNumber
	}
	if req.Notes != "" {
		expense.Notes = &req.Notes
	}

	if err := s.repo.CreateExpense(ctx, expense); err != nil {
		return nil, err
	}
	return expenseToResponse(expense), nil
}

func (s *Service) ListExpensesByTravel(ctx context.Context, travelIDStr string) ([]ExpenseResponse, error) {
	travelID, err := uuid.Parse(travelIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	expenses, err := s.repo.ListExpensesByTravel(ctx, travelID)
	if err != nil {
		return nil, err
	}
	responses := make([]ExpenseResponse, 0, len(expenses))
	for i := range expenses {
		responses = append(responses, *expenseToResponse(&expenses[i]))
	}
	return responses, nil
}

func (s *Service) UpdateExpense(ctx context.Context, expenseIDStr string, req UpdateExpenseRequest) (*ExpenseResponse, error) {
	expenseID, err := uuid.Parse(expenseIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid expense id: %w", err)
	}
	expense, err := s.repo.FindExpenseByID(ctx, expenseID)
	if err != nil {
		return nil, err
	}
	if expense.Status == ExpenseStatusApproved {
		return nil, ErrExpenseInvalidState
	}
	if req.ExpenseCategoryID != nil {
		if categoryID, err := uuid.Parse(*req.ExpenseCategoryID); err == nil {
			expense.ExpenseCategoryID = categoryID
		}
	}
	if req.ExpenseDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.ExpenseDate)
		if err != nil {
			return nil, fmt.Errorf("invalid expense_date: %w", err)
		}
		expense.ExpenseDate = parsed
	}
	if req.Description != nil {
		expense.Description = req.Description
	}
	if req.Quantity != nil {
		expense.Quantity = *req.Quantity
	}
	if req.Unit != nil {
		expense.Unit = req.Unit
	}
	if req.Amount != nil {
		expense.Amount = *req.Amount
	}
	if req.FundingMethodID != nil {
		if *req.FundingMethodID == "" {
			expense.FundingMethodID = nil
		} else if fundingMethodID, err := uuid.Parse(*req.FundingMethodID); err == nil {
			expense.FundingMethodID = &fundingMethodID
		}
	}
	if req.Vendor != nil {
		expense.Vendor = req.Vendor
	}
	if req.ReceiptNumber != nil {
		expense.ReceiptNumber = req.ReceiptNumber
	}
	if req.Notes != nil {
		expense.Notes = req.Notes
	}
	if err := s.repo.UpdateExpense(ctx, expense); err != nil {
		return nil, err
	}
	return expenseToResponse(expense), nil
}

func (s *Service) DeleteExpense(ctx context.Context, expenseIDStr string) error {
	expenseID, err := uuid.Parse(expenseIDStr)
	if err != nil {
		return fmt.Errorf("invalid expense id: %w", err)
	}
	expense, err := s.repo.FindExpenseByID(ctx, expenseID)
	if err != nil {
		return err
	}
	if expense.Status == ExpenseStatusApproved {
		return ErrExpenseInvalidState
	}
	return s.repo.DeleteExpense(ctx, expenseID)
}

// AddExpenseDocument attaches proof of expense (receipt/invoice/ticket/dst,
// §22). Sama seperti AddFundingDocument: URL didapat dari endpoint upload
// generik, module ini tidak menangani file mentah (§54.4).
func (s *Service) AddExpenseDocument(ctx context.Context, expenseIDStr string, req AddExpenseDocumentRequest) (*ExpenseResponse, error) {
	expenseID, err := uuid.Parse(expenseIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid expense id: %w", err)
	}
	if _, err := s.repo.FindExpenseByID(ctx, expenseID); err != nil {
		return nil, err
	}
	doc := &ExpenseDocument{
		BusinessTravelExpenseID: expenseID,
		DocumentType:            ExpenseDocumentType(strings.ToUpper(req.DocumentType)),
		FileName:                req.FileName,
		FilePath:                req.FilePath,
		UploadedBy:              authctx.GetUserID(ctx),
	}
	if req.MimeType != "" {
		doc.MimeType = &req.MimeType
	}
	if req.FileSize > 0 {
		doc.FileSize = &req.FileSize
	}
	now := time.Now()
	doc.UploadedAt = &now
	if err := s.repo.CreateExpenseDocument(ctx, doc); err != nil {
		return nil, err
	}
	expense, err := s.repo.FindExpenseByID(ctx, expenseID)
	if err != nil {
		return nil, err
	}
	return expenseToResponse(expense), nil
}

// =========================================================================
// Settlement (§24-33 plan doc)
// =========================================================================

// businessTravelSettlementApprovalModule is a separate module slug from
// businessTravelApprovalModule so the Settlement Approval flow can be
// configured independently from Travel Approval (§54.3, Rule 6/§52: reuse
// the central Approval module, no bespoke engine).
const businessTravelSettlementApprovalModule = "business_travel_settlement"

// ErrBusinessTravelNotCompleted: settlement hanya boleh dibuat setelah
// travel COMPLETED (§24: IN_PROGRESS -> COMPLETED -> SETTLEMENT).
var ErrBusinessTravelNotCompleted = errors.New("business travel must be COMPLETED before settlement can be created")

// ErrSettlementInvalidState: aksi tidak valid untuk status settlement saat ini.
var ErrSettlementInvalidState = errors.New("settlement is not in a valid state for this action")

// CreateSettlement menghitung total advance/actual/company-paid dan
// menentukan hasil awal (balanced/refund/reimbursement) dari funding &
// expense yang sudah tercatat (§25-27, §35-36). Company-paid expense
// dikeluarkan dari rekonsiliasi advance karena bukan hutang ke employee
// (§34). Hasil baru final setelah SubmitSettlement disetujui — lihat
// HandleSettlementApprovalStatusChange.
func (s *Service) CreateSettlement(ctx context.Context, travelIDStr string, req CreateSettlementRequest) (*SettlementResponse, error) {
	travelID, err := uuid.Parse(travelIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	travel, err := s.repo.FindBusinessTravelByIDForOwnership(ctx, travelID)
	if err != nil {
		return nil, err
	}
	if travel.Status != TravelStatusCompleted {
		return nil, ErrBusinessTravelNotCompleted
	}

	var participantID *uuid.UUID
	if req.ParticipantID != "" {
		parsed, err := uuid.Parse(req.ParticipantID)
		if err != nil {
			return nil, fmt.Errorf("invalid participant_id: %w", err)
		}
		participantID = &parsed
	}

	fundings, err := s.repo.ListFundingsByTravel(ctx, travelID)
	if err != nil {
		return nil, err
	}
	expenses, err := s.repo.ListExpensesByTravel(ctx, travelID)
	if err != nil {
		return nil, err
	}
	methods, err := s.repo.ListFundingMethods(ctx, false)
	if err != nil {
		return nil, err
	}
	methodCodeByID := make(map[uuid.UUID]string, len(methods))
	for _, m := range methods {
		methodCodeByID[m.ID] = m.Code
	}

	matchesParticipant := func(pid *uuid.UUID) bool {
		if participantID == nil {
			return true
		}
		return pid != nil && *pid == *participantID
	}

	var totalAdvance, totalActual, totalCompanyPaid float64
	relevantFundings := make([]Funding, 0, len(fundings))
	for _, f := range fundings {
		if !matchesParticipant(f.ParticipantID) || f.Status != FundingStatusFunded {
			continue
		}
		if methodCodeByID[f.FundingMethodID] == FundingMethodDeposit {
			totalAdvance += f.Amount
			relevantFundings = append(relevantFundings, f)
		}
	}
	relevantExpenses := make([]Expense, 0, len(expenses))
	for _, e := range expenses {
		if !matchesParticipant(e.ParticipantID) {
			continue
		}
		totalActual += e.Amount
		if e.FundingMethodID != nil && methodCodeByID[*e.FundingMethodID] == FundingMethodCompanyPaid {
			totalCompanyPaid += e.Amount
		}
		relevantExpenses = append(relevantExpenses, e)
	}

	// diff > 0: actual (net of company-paid) melebihi advance -> additional
	// reimbursement (§31 Scenario 4). diff < 0: advance tersisa -> refund
	// (§30 Scenario 3). diff == 0: balanced (§29 Scenario 2).
	diff := (totalActual - totalCompanyPaid) - totalAdvance

	settlement := &Settlement{
		BusinessTravelID:   travelID,
		ParticipantID:      participantID,
		TotalAdvance:       totalAdvance,
		TotalActualExpense: totalActual,
		TotalCompanyPaid:   totalCompanyPaid,
		Balance:            diff,
		Status:             SettlementStatusPending,
	}
	if diff > 0 {
		settlement.TotalReimbursement = diff
	} else if diff < 0 {
		settlement.TotalRefund = -diff
	}
	if req.Notes != "" {
		settlement.Notes = &req.Notes
	}

	if err := s.repo.CreateSettlement(ctx, settlement); err != nil {
		return nil, err
	}

	for _, f := range relevantFundings {
		fundingMethodID := f.FundingMethodID
		item := &SettlementItem{
			BusinessTravelSettlementID: settlement.ID,
			FundingMethodID:            &fundingMethodID,
			ItemType:                   SettlementItemAdvance,
			Amount:                     f.Amount,
		}
		if err := s.repo.CreateSettlementItem(ctx, item); err != nil {
			return nil, err
		}
	}
	for _, e := range relevantExpenses {
		itemType := SettlementItemActual
		if e.FundingMethodID != nil && methodCodeByID[*e.FundingMethodID] == FundingMethodCompanyPaid {
			itemType = SettlementItemCompanyPaid
		}
		expenseID := e.ID
		item := &SettlementItem{
			BusinessTravelSettlementID: settlement.ID,
			ExpenseID:                  &expenseID,
			FundingMethodID:            e.FundingMethodID,
			ItemType:                   itemType,
			Amount:                     e.Amount,
		}
		if err := s.repo.CreateSettlementItem(ctx, item); err != nil {
			return nil, err
		}
	}
	if diff > 0 {
		item := &SettlementItem{BusinessTravelSettlementID: settlement.ID, ItemType: SettlementItemReimbursement, Amount: diff}
		if err := s.repo.CreateSettlementItem(ctx, item); err != nil {
			return nil, err
		}
	} else if diff < 0 {
		item := &SettlementItem{BusinessTravelSettlementID: settlement.ID, ItemType: SettlementItemRefund, Amount: -diff}
		if err := s.repo.CreateSettlementItem(ctx, item); err != nil {
			return nil, err
		}
	}

	return s.GetSettlementByID(ctx, settlement.ID.String())
}

func (s *Service) GetSettlementByID(ctx context.Context, id string) (*SettlementResponse, error) {
	settlementID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	settlement, err := s.repo.FindSettlementByID(ctx, settlementID)
	if err != nil {
		return nil, err
	}
	return settlementToResponse(settlement), nil
}

func (s *Service) ListSettlementsByTravel(ctx context.Context, travelIDStr string) ([]SettlementResponse, error) {
	travelID, err := uuid.Parse(travelIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	settlements, err := s.repo.ListSettlementsByTravel(ctx, travelID)
	if err != nil {
		return nil, err
	}
	responses := make([]SettlementResponse, 0, len(settlements))
	for i := range settlements {
		responses = append(responses, *settlementToResponse(&settlements[i]))
	}
	return responses, nil
}

// SubmitSettlement routes a PENDING settlement through the central Approval
// module under its own module slug (Rule 6, §52), mirroring
// Service.SubmitBusinessTravel.
func (s *Service) SubmitSettlement(ctx context.Context, settlementIDStr string, req SubmitSettlementRequest) (*SettlementResponse, error) {
	settlementID, err := uuid.Parse(settlementIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid settlement id: %w", err)
	}
	settlement, err := s.repo.FindSettlementByID(ctx, settlementID)
	if err != nil {
		return nil, err
	}
	if settlement.Status != SettlementStatusPending {
		return nil, ErrSettlementInvalidState
	}

	now := time.Now()
	settlement.Status = SettlementStatusSubmitted
	settlement.SubmittedAt = &now

	if s.approvalEngine != nil {
		flowID := ""
		if req.FlowID != nil && *req.FlowID != "" {
			flowID = *req.FlowID
		} else if resolved, err := s.approvalEngine.GetActiveFlowIDForModule(ctx, businessTravelSettlementApprovalModule); err == nil {
			flowID = resolved
		}
		if flowID != "" {
			instanceID, err := s.approvalEngine.CreateApprovalInstance(ctx, businessTravelSettlementApprovalModule, settlement.ID.String(), flowID)
			if err != nil {
				var re *approval.RoutingError
				if errors.As(err, &re) {
					return nil, err
				}
				s.logger.Warn("Failed to create approval instance for settlement, continuing without approval",
					zap.String("settlement_id", settlement.ID.String()),
					zap.Error(err),
				)
			} else if parsedInstanceID, parseErr := uuid.Parse(instanceID); parseErr == nil {
				settlement.ApprovalInstanceID = &parsedInstanceID
			}
		}
	}

	if err := s.repo.UpdateSettlement(ctx, settlement); err != nil {
		return nil, err
	}
	return settlementToResponse(settlement), nil
}

// HandleSettlementApprovalStatusChange is invoked by the approval module's
// push-based status callback for the "business_travel_settlement" module
// slug. On approval, finalizes the settlement outcome computed at
// CreateSettlement time (§26) and creates the corresponding Refund or
// TravelReimbursement record (§35-36) so Phase 7/8 processing has something
// to act on.
func (s *Service) HandleSettlementApprovalStatusChange(ctx context.Context, documentID uuid.UUID, status string, note string) error {
	settlement, err := s.repo.FindSettlementByID(ctx, documentID)
	if err != nil {
		return err
	}
	switch status {
	case "APPROVED":
		now := time.Now()
		settlement.ApprovedAt = &now
		switch {
		case settlement.TotalReimbursement > 0:
			settlement.Status = SettlementStatusReimbursementRequired
			reimbursement := &TravelReimbursement{
				BusinessTravelID: settlement.BusinessTravelID,
				ParticipantID:    settlement.ParticipantID,
				SettlementID:     &settlement.ID,
				Amount:           settlement.TotalReimbursement,
				Status:           TravelReimbStatusRequested,
				RequestedAt:      &now,
			}
			if err := s.repo.CreateTravelReimbursement(ctx, reimbursement); err != nil {
				return err
			}
		case settlement.TotalRefund > 0:
			settlement.Status = SettlementStatusRefundRequired
			refund := &Refund{
				BusinessTravelID: settlement.BusinessTravelID,
				SettlementID:     &settlement.ID,
				ParticipantID:    settlement.ParticipantID,
				RefundAmount:     settlement.TotalRefund,
				Status:           RefundStatusPending,
			}
			if err := s.repo.CreateRefund(ctx, refund); err != nil {
				return err
			}
		default:
			settlement.Status = SettlementStatusBalanced
		}
	case "REJECTED":
		settlement.Status = SettlementStatusRejected
	default:
		return nil
	}
	return s.repo.UpdateSettlement(ctx, settlement)
}

// =========================================================================
// Refund (§35 plan doc)
// =========================================================================

// ErrRefundInvalidState: aksi tidak valid untuk status refund saat ini.
var ErrRefundInvalidState = errors.New("refund is not in a valid state for this action")

func (s *Service) ListRefundsByTravel(ctx context.Context, travelIDStr string) ([]RefundResponse, error) {
	travelID, err := uuid.Parse(travelIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	refunds, err := s.repo.ListRefundsByTravel(ctx, travelID)
	if err != nil {
		return nil, err
	}
	responses := make([]RefundResponse, 0, len(refunds))
	for i := range refunds {
		responses = append(responses, *refundToResponse(&refunds[i]))
	}
	return responses, nil
}

// ConfirmRefund menandai refund sebagai CONFIRMED (uang sudah dikembalikan
// employee ke perusahaan, §35), lalu mencoba menutup settlement & travel
// terkait (§26: SETTLED, §Kesimpulan: CLOSED) via maybeSettleAndCloseTravel.
func (s *Service) ConfirmRefund(ctx context.Context, refundIDStr string, req ConfirmRefundRequest) (*RefundResponse, error) {
	refundID, err := uuid.Parse(refundIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid refund id: %w", err)
	}
	refund, err := s.repo.FindRefundByID(ctx, refundID)
	if err != nil {
		return nil, err
	}
	if refund.Status != RefundStatusPending {
		return nil, ErrRefundInvalidState
	}
	now := time.Now()
	refund.Status = RefundStatusConfirmed
	refund.RefundDate = &now
	refund.RefundedBy = authctx.GetUserID(ctx)
	if req.RefundReference != "" {
		refund.RefundReference = &req.RefundReference
	}
	if req.RefundDocument != "" {
		refund.RefundDocument = &req.RefundDocument
	}
	if err := s.repo.UpdateRefund(ctx, refund); err != nil {
		return nil, err
	}
	if refund.SettlementID != nil {
		s.maybeSettleAndCloseTravel(ctx, *refund.SettlementID, refund.BusinessTravelID)
	}
	return refundToResponse(refund), nil
}

// =========================================================================
// Reimbursement (§36 plan doc + §54.7: cek subscription module Reimbursement)
// =========================================================================

// ErrTravelReimbursementInvalidState: aksi tidak valid untuk status
// reimbursement saat ini.
var ErrTravelReimbursementInvalidState = errors.New("reimbursement is not in a valid state for this action")

const reimbursementModuleSlug = "reimbursement"

func (s *Service) ListTravelReimbursements(ctx context.Context, travelIDStr string) ([]TravelReimbursementResponse, error) {
	travelID, err := uuid.Parse(travelIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	reimbursements, err := s.repo.ListTravelReimbursementsByTravel(ctx, travelID)
	if err != nil {
		return nil, err
	}
	responses := make([]TravelReimbursementResponse, 0, len(reimbursements))
	for i := range reimbursements {
		responses = append(responses, *travelReimbursementToResponse(&reimbursements[i]))
	}
	return responses, nil
}

func (s *Service) ApproveTravelReimbursement(ctx context.Context, reimbursementIDStr string) (*TravelReimbursementResponse, error) {
	reimbursementID, err := uuid.Parse(reimbursementIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid reimbursement id: %w", err)
	}
	reimbursement, err := s.repo.FindTravelReimbursementByID(ctx, reimbursementID)
	if err != nil {
		return nil, err
	}
	if reimbursement.Status != TravelReimbStatusRequested {
		return nil, ErrTravelReimbursementInvalidState
	}
	now := time.Now()
	reimbursement.Status = TravelReimbStatusApproved
	reimbursement.ApprovedAt = &now
	if err := s.repo.UpdateTravelReimbursement(ctx, reimbursement); err != nil {
		return nil, err
	}
	return travelReimbursementToResponse(reimbursement), nil
}

// ProcessTravelReimbursement transitions APPROVED -> PROCESSING. Per §54.7:
// if the tenant subscribes to the standalone Reimbursement module, the
// claim should ideally be handed off to it instead of duplicating payout
// logic here; this is checked and logged as a hint for operators, but
// actual cross-module claim creation isn't implemented yet (the
// Reimbursement module has no public API for accepting externally-created
// claims) — both paths currently process the claim internally. Not
// subscribed -> always processed internally, per §54.7's ELSE branch.
func (s *Service) ProcessTravelReimbursement(ctx context.Context, reimbursementIDStr string) (*TravelReimbursementResponse, error) {
	reimbursementID, err := uuid.Parse(reimbursementIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid reimbursement id: %w", err)
	}
	reimbursement, err := s.repo.FindTravelReimbursementByID(ctx, reimbursementID)
	if err != nil {
		return nil, err
	}
	if reimbursement.Status != TravelReimbStatusApproved {
		return nil, ErrTravelReimbursementInvalidState
	}

	if s.moduleChecker != nil {
		if companyID := authctx.GetCompanyID(ctx); companyID != "" {
			if active, err := s.moduleChecker.IsModuleActive(companyID, reimbursementModuleSlug); err != nil {
				s.logger.Warn("Failed to check reimbursement module subscription, processing internally",
					zap.String("business_travel_reimbursement_id", reimbursement.ID.String()),
					zap.Error(err),
				)
			} else if active {
				s.logger.Info("Tenant subscribes to the Reimbursement module; business travel reimbursement claim is still processed internally pending cross-module integration (§54.7)",
					zap.String("business_travel_reimbursement_id", reimbursement.ID.String()),
				)
			}
		}
	}

	reimbursement.Status = TravelReimbStatusProcessing
	if err := s.repo.UpdateTravelReimbursement(ctx, reimbursement); err != nil {
		return nil, err
	}
	return travelReimbursementToResponse(reimbursement), nil
}

func (s *Service) PayTravelReimbursement(ctx context.Context, reimbursementIDStr string, req PayReimbursementRequest) (*TravelReimbursementResponse, error) {
	reimbursementID, err := uuid.Parse(reimbursementIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid reimbursement id: %w", err)
	}
	reimbursement, err := s.repo.FindTravelReimbursementByID(ctx, reimbursementID)
	if err != nil {
		return nil, err
	}
	if reimbursement.Status != TravelReimbStatusProcessing {
		return nil, ErrTravelReimbursementInvalidState
	}
	now := time.Now()
	reimbursement.Status = TravelReimbStatusPaid
	reimbursement.PaidAt = &now
	reimbursement.PaidBy = authctx.GetUserID(ctx)
	if req.PaymentReference != "" {
		reimbursement.PaymentReference = &req.PaymentReference
	}
	if err := s.repo.UpdateTravelReimbursement(ctx, reimbursement); err != nil {
		return nil, err
	}
	if reimbursement.SettlementID != nil {
		s.maybeSettleAndCloseTravel(ctx, *reimbursement.SettlementID, reimbursement.BusinessTravelID)
	}
	return travelReimbursementToResponse(reimbursement), nil
}

// maybeSettleAndCloseTravel marks a settlement SETTLED once its refund/
// reimbursement outcome is finalized (confirmed/paid), then closes the
// parent travel (§Kesimpulan: ... -> CLOSED) once every settlement under it
// is BALANCED or SETTLED. Best-effort: logs and continues on error instead
// of failing the caller's refund/payment action, since the settlement/travel
// status here is a derived convenience, not the source of truth for the
// refund/reimbursement record itself.
func (s *Service) maybeSettleAndCloseTravel(ctx context.Context, settlementID, travelID uuid.UUID) {
	settlement, err := s.repo.FindSettlementByID(ctx, settlementID)
	if err != nil {
		s.logger.Warn("Failed to load settlement while finalizing close", zap.String("settlement_id", settlementID.String()), zap.Error(err))
		return
	}
	if settlement.Status == SettlementStatusRefundRequired || settlement.Status == SettlementStatusReimbursementRequired {
		now := time.Now()
		settlement.Status = SettlementStatusSettled
		settlement.SettledAt = &now
		if err := s.repo.UpdateSettlement(ctx, settlement); err != nil {
			s.logger.Warn("Failed to mark settlement SETTLED", zap.String("settlement_id", settlementID.String()), zap.Error(err))
			return
		}
	}

	settlements, err := s.repo.ListSettlementsByTravel(ctx, travelID)
	if err != nil {
		s.logger.Warn("Failed to list settlements while checking travel close", zap.String("business_travel_id", travelID.String()), zap.Error(err))
		return
	}
	if len(settlements) == 0 {
		return
	}
	for _, st := range settlements {
		if st.Status != SettlementStatusBalanced && st.Status != SettlementStatusSettled {
			return
		}
	}
	travel, err := s.repo.FindBusinessTravelByIDForOwnership(ctx, travelID)
	if err != nil {
		s.logger.Warn("Failed to load travel while checking close", zap.String("business_travel_id", travelID.String()), zap.Error(err))
		return
	}
	if travel.Status == TravelStatusClosed {
		return
	}
	travel.Status = TravelStatusClosed
	if err := s.repo.UpdateBusinessTravel(ctx, travel); err != nil {
		s.logger.Warn("Failed to close travel", zap.String("business_travel_id", travelID.String()), zap.Error(err))
	}
}

// HandleBusinessTravelApprovalStatusChange is invoked by the approval
// module's push-based status callback for the "business_travel" module slug
// (registered separately from HandleApprovalStatusChange, which handles the
// "attendance" slug used by overtime/correction).
func (s *Service) HandleBusinessTravelApprovalStatusChange(ctx context.Context, documentID uuid.UUID, status string, note string) error {
	travel, err := s.repo.FindBusinessTravelByID(ctx, documentID)
	if err != nil {
		return err
	}
	switch status {
	case "APPROVED":
		travel.Status = TravelStatusApproved
		travel.ApprovalStatus = string(TravelStatusApproved)
		s.pushBusinessTravelAttendance(ctx, travel)
	case "REJECTED":
		travel.Status = TravelStatusRejected
		travel.ApprovalStatus = string(TravelStatusRejected)
	default:
		travel.ApprovalStatus = status
	}
	return s.repo.UpdateBusinessTravel(ctx, travel)
}

// pushBusinessTravelAttendance marks every EMPLOYEE participant's attendance
// session BUSINESS_TRAVEL for each day of the trip (§37 plan doc), once the
// travel is APPROVED. Best-effort per participant/day: logs and continues
// rather than failing the approval callback, matching notifyRequestOutcome's
// policy for auxiliary side effects.
func (s *Service) pushBusinessTravelAttendance(ctx context.Context, travel *BusinessTravel) {
	for _, p := range travel.Participants {
		if p.ParticipantType != ParticipantTypeEmployee || p.EmployeeID == nil {
			continue
		}
		for d := travel.StartDate; !d.After(travel.EndDate); d = d.AddDate(0, 0, 1) {
			if err := s.ApplyApprovedBusinessTravel(ctx, *p.EmployeeID, d.Format("2006-01-02"), travel.ID); err != nil {
				s.logger.Warn("Failed to apply business travel to attendance session",
					zap.String("business_travel_id", travel.ID.String()),
					zap.String("employee_id", p.EmployeeID.String()),
					zap.String("work_date", d.Format("2006-01-02")),
					zap.Error(err),
				)
			}
		}
	}
}

func generateBusinessTravelRequestNumber() string {
	return fmt.Sprintf("BT-%s-%s", time.Now().Format("200601"), strings.ToUpper(uuid.New().String()[:8]))
}

func businessTravelToResponse(t *BusinessTravel) *BusinessTravelResponse {
	resp := &BusinessTravelResponse{
		ID:             t.ID.String(),
		RequestNumber:  t.RequestNumber,
		RequesterID:    t.RequesterID.String(),
		Title:          t.Title,
		Purpose:        t.Purpose,
		Description:    t.Description,
		StartDate:      t.StartDate.Format("2006-01-02"),
		EndDate:        t.EndDate.Format("2006-01-02"),
		Origin:         t.Origin,
		Status:         string(t.Status),
		ApprovalStatus: t.ApprovalStatus,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}
	if t.ApprovalInstanceID != nil {
		instanceID := t.ApprovalInstanceID.String()
		resp.ApprovalInstanceID = &instanceID
	}
	for _, p := range t.Participants {
		resp.Participants = append(resp.Participants, businessTravelParticipantToResponse(&p))
	}
	for _, d := range t.Destinations {
		resp.Destinations = append(resp.Destinations, businessTravelDestinationToResponse(&d))
	}
	return resp
}

func businessTravelParticipantToResponse(p *BusinessTravelParticipant) BusinessTravelParticipantResponse {
	resp := BusinessTravelParticipantResponse{
		ID:              p.ID.String(),
		ParticipantType: string(p.ParticipantType),
		Name:            p.Name,
		Organization:    p.Organization,
		Position:        p.Position,
		IdentityNumber:  p.IdentityNumber,
		Email:           p.Email,
		Phone:           p.Phone,
		Role:            string(p.Role),
		Notes:           p.Notes,
	}
	if p.EmployeeID != nil {
		empID := p.EmployeeID.String()
		resp.EmployeeID = &empID
	}
	return resp
}

func businessTravelDestinationToResponse(d *BusinessTravelDestination) BusinessTravelDestinationResponse {
	resp := BusinessTravelDestinationResponse{
		ID:       d.ID.String(),
		Sequence: d.Sequence,
		Country:  d.Country,
		Province: d.Province,
		City:     d.City,
		Location: d.Location,
		Purpose:  d.Purpose,
		Notes:    d.Notes,
	}
	if d.ArrivalDate != nil {
		arrival := d.ArrivalDate.Format("2006-01-02")
		resp.ArrivalDate = &arrival
	}
	if d.DepartureDate != nil {
		departure := d.DepartureDate.Format("2006-01-02")
		resp.DepartureDate = &departure
	}
	return resp
}

func businessTravelActivityToResponse(a *BusinessTravelActivity) BusinessTravelActivityResponse {
	return BusinessTravelActivityResponse{
		ID:           a.ID.String(),
		ActivityDate: a.ActivityDate.Format("2006-01-02"),
		StartTime:    a.StartTime,
		EndTime:      a.EndTime,
		Title:        a.Title,
		Description:  a.Description,
		Location:     a.Location,
		Organizer:    a.Organizer,
		Notes:        a.Notes,
	}
}

func refundToResponse(rf *Refund) *RefundResponse {
	resp := &RefundResponse{
		ID:               rf.ID.String(),
		BusinessTravelID: rf.BusinessTravelID.String(),
		RefundAmount:     rf.RefundAmount,
		RefundReference:  rf.RefundReference,
		RefundDocument:   rf.RefundDocument,
		Status:           string(rf.Status),
		Notes:            rf.Notes,
		CreatedAt:        rf.CreatedAt,
	}
	if rf.SettlementID != nil {
		settlementID := rf.SettlementID.String()
		resp.SettlementID = &settlementID
	}
	if rf.ParticipantID != nil {
		participantID := rf.ParticipantID.String()
		resp.ParticipantID = &participantID
	}
	if rf.RefundedBy != nil {
		refundedBy := rf.RefundedBy.String()
		resp.RefundedBy = &refundedBy
	}
	if rf.RefundDate != nil {
		refundDate := rf.RefundDate.Format("2006-01-02")
		resp.RefundDate = &refundDate
	}
	return resp
}

func travelReimbursementToResponse(tr *TravelReimbursement) *TravelReimbursementResponse {
	resp := &TravelReimbursementResponse{
		ID:               tr.ID.String(),
		BusinessTravelID: tr.BusinessTravelID.String(),
		Amount:           tr.Amount,
		Status:           string(tr.Status),
		RequestedAt:      tr.RequestedAt,
		ApprovedAt:       tr.ApprovedAt,
		PaidAt:           tr.PaidAt,
		PaymentReference: tr.PaymentReference,
		Notes:            tr.Notes,
		CreatedAt:        tr.CreatedAt,
	}
	if tr.ParticipantID != nil {
		participantID := tr.ParticipantID.String()
		resp.ParticipantID = &participantID
	}
	if tr.SettlementID != nil {
		settlementID := tr.SettlementID.String()
		resp.SettlementID = &settlementID
	}
	if tr.PaidBy != nil {
		paidBy := tr.PaidBy.String()
		resp.PaidBy = &paidBy
	}
	return resp
}

func settlementToResponse(s *Settlement) *SettlementResponse {
	resp := &SettlementResponse{
		ID:                 s.ID.String(),
		BusinessTravelID:   s.BusinessTravelID.String(),
		TotalAdvance:       s.TotalAdvance,
		TotalActualExpense: s.TotalActualExpense,
		TotalCompanyPaid:   s.TotalCompanyPaid,
		TotalReimbursement: s.TotalReimbursement,
		TotalRefund:        s.TotalRefund,
		Balance:            s.Balance,
		Status:             string(s.Status),
		SubmittedAt:        s.SubmittedAt,
		ApprovedAt:         s.ApprovedAt,
		SettledAt:          s.SettledAt,
		Notes:              s.Notes,
		CreatedAt:          s.CreatedAt,
		UpdatedAt:          s.UpdatedAt,
	}
	if s.ParticipantID != nil {
		participantID := s.ParticipantID.String()
		resp.ParticipantID = &participantID
	}
	if s.ApprovalInstanceID != nil {
		instanceID := s.ApprovalInstanceID.String()
		resp.ApprovalInstanceID = &instanceID
	}
	for _, item := range s.Items {
		resp.Items = append(resp.Items, settlementItemToResponse(&item))
	}
	return resp
}

func settlementItemToResponse(i *SettlementItem) SettlementItemResponse {
	resp := SettlementItemResponse{
		ID:       i.ID.String(),
		ItemType: string(i.ItemType),
		Category: i.Category,
		Amount:   i.Amount,
		Notes:    i.Notes,
	}
	if i.ExpenseID != nil {
		expenseID := i.ExpenseID.String()
		resp.ExpenseID = &expenseID
	}
	if i.FundingMethodID != nil {
		fundingMethodID := i.FundingMethodID.String()
		resp.FundingMethodID = &fundingMethodID
	}
	return resp
}

func expenseCategoryToResponse(c *ExpenseCategory) *ExpenseCategoryResponse {
	return &ExpenseCategoryResponse{
		ID:               c.ID.String(),
		Code:             c.Code,
		Name:             c.Name,
		Description:      c.Description,
		RequiresReceipt:  c.RequiresReceipt,
		Reimbursable:     c.Reimbursable,
		PayrollTreatment: c.PayrollTreatment,
		AccountCode:      c.AccountCode,
		Active:           c.Active,
	}
}

func expenseToResponse(e *Expense) *ExpenseResponse {
	resp := &ExpenseResponse{
		ID:                e.ID.String(),
		BusinessTravelID:  e.BusinessTravelID.String(),
		ExpenseCategoryID: e.ExpenseCategoryID.String(),
		ExpenseDate:       e.ExpenseDate.Format("2006-01-02"),
		Description:       e.Description,
		Quantity:          e.Quantity,
		Unit:              e.Unit,
		Amount:            e.Amount,
		Vendor:            e.Vendor,
		ReceiptNumber:     e.ReceiptNumber,
		Status:            string(e.Status),
		Notes:             e.Notes,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
	if e.ParticipantID != nil {
		participantID := e.ParticipantID.String()
		resp.ParticipantID = &participantID
	}
	if e.FundingMethodID != nil {
		fundingMethodID := e.FundingMethodID.String()
		resp.FundingMethodID = &fundingMethodID
	}
	for _, d := range e.Documents {
		resp.Documents = append(resp.Documents, expenseDocumentToResponse(&d))
	}
	return resp
}

func expenseDocumentToResponse(d *ExpenseDocument) ExpenseDocumentResponse {
	return ExpenseDocumentResponse{
		ID:           d.ID.String(),
		DocumentType: string(d.DocumentType),
		FileName:     d.FileName,
		FilePath:     d.FilePath,
		MimeType:     d.MimeType,
		FileSize:     d.FileSize,
	}
}

func fundingMethodToResponse(m *FundingMethod) *FundingMethodResponse {
	return &FundingMethodResponse{
		ID:          m.ID.String(),
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		Active:      m.Active,
	}
}

func fundingToResponse(f *Funding) *FundingResponse {
	resp := &FundingResponse{
		ID:               f.ID.String(),
		BusinessTravelID: f.BusinessTravelID.String(),
		FundingMethodID:  f.FundingMethodID.String(),
		Amount:           f.Amount,
		PaymentMethod:    f.PaymentMethod,
		PaymentReference: f.PaymentReference,
		Status:           string(f.Status),
		Notes:            f.Notes,
		CreatedAt:        f.CreatedAt,
		UpdatedAt:        f.UpdatedAt,
	}
	if f.ParticipantID != nil {
		participantID := f.ParticipantID.String()
		resp.ParticipantID = &participantID
	}
	if f.FundedBy != nil {
		fundedBy := f.FundedBy.String()
		resp.FundedBy = &fundedBy
	}
	if f.FundingDate != nil {
		fundingDate := f.FundingDate.Format("2006-01-02")
		resp.FundingDate = &fundingDate
	}
	for _, d := range f.Documents {
		resp.Documents = append(resp.Documents, fundingDocumentToResponse(&d))
	}
	return resp
}

func fundingDocumentToResponse(d *FundingDocument) FundingDocumentResponse {
	return FundingDocumentResponse{
		ID:           d.ID.String(),
		DocumentType: string(d.DocumentType),
		FileName:     d.FileName,
		FilePath:     d.FilePath,
		MimeType:     d.MimeType,
		FileSize:     d.FileSize,
	}
}

func businessTravelScheduleToResponse(s *BusinessTravelSchedule) BusinessTravelScheduleResponse {
	return BusinessTravelScheduleResponse{
		ID:                 s.ID.String(),
		ScheduleType:       string(s.ScheduleType),
		DepartureDatetime:  s.DepartureDatetime,
		ArrivalDatetime:    s.ArrivalDatetime,
		Origin:             s.Origin,
		Destination:        s.Destination,
		TransportationType: string(s.TransportationType),
		Provider:           s.Provider,
		BookingReference:   s.BookingReference,
		Notes:              s.Notes,
	}
}
