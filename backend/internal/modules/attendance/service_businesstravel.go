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
	case "REJECTED":
		travel.Status = TravelStatusRejected
		travel.ApprovalStatus = string(TravelStatusRejected)
	default:
		travel.ApprovalStatus = status
	}
	return s.repo.UpdateBusinessTravel(ctx, travel)
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
