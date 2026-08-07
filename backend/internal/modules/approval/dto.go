package approval

import (
	"time"

	"github.com/google/uuid"
)

// =========================================================================
// Approval Flow DTOs
// =========================================================================

type CreateFlowRequest struct {
	Module  string `json:"module" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Version int    `json:"version"`
}

type UpdateFlowRequest struct {
	Module   *string `json:"module"`
	Name     *string `json:"name"`
	Version  *int    `json:"version"`
	IsActive *bool   `json:"is_active"`
}

type FlowResponse struct {
	ID        string             `json:"id"`
	Module    string             `json:"module"`
	Name      string             `json:"name"`
	Version   int                `json:"version"`
	IsActive  bool               `json:"is_active"`
	Steps     []FlowStepResponse `json:"steps,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// =========================================================================
// Approval Flow Step DTOs
// =========================================================================

type CreateStepRequest struct {
	StepName          string   `json:"step_name" binding:"required"`
	ApproverType      string   `json:"approver_type" binding:"required,oneof=SUPERVISOR ROLE USER ORGANIZATION BOTH"`
	HierarchyLevel    *int     `json:"hierarchy_level"`
	OrganizationIDs   []string `json:"organization_ids" binding:"omitempty,dive,uuid"`
	RoleID            *string  `json:"role_id" binding:"omitempty,uuid"`
	ApproverUserID    *string  `json:"approver_user_id" binding:"omitempty,uuid"`
	ApprovalMode      string   `json:"approval_mode" binding:"omitempty,oneof=ANY_ONE ALL N_OF_M"`
	ParticipationType string   `json:"participation_type" binding:"omitempty,oneof=APPROVER WATCHER"`
	RequiredApprovals *int     `json:"required_approvals"`
	AllowReject       *bool    `json:"allow_reject"`
	ConditionsJSON    *string  `json:"conditions_json"`
	SLAHours          *int     `json:"sla_hours"`
}

type UpdateStepRequest struct {
	StepName          *string  `json:"step_name"`
	ApproverType      *string  `json:"approver_type" binding:"omitempty,oneof=SUPERVISOR ROLE USER ORGANIZATION BOTH"`
	HierarchyLevel    *int     `json:"hierarchy_level"`
	OrganizationIDs   []string `json:"organization_ids" binding:"omitempty,dive,uuid"`
	RoleID            *string  `json:"role_id" binding:"omitempty,uuid"`
	ApproverUserID    *string  `json:"approver_user_id" binding:"omitempty,uuid"`
	ApprovalMode      *string  `json:"approval_mode" binding:"omitempty,oneof=ANY_ONE ALL N_OF_M"`
	ParticipationType *string  `json:"participation_type" binding:"omitempty,oneof=APPROVER WATCHER"`
	RequiredApprovals *int     `json:"required_approvals"`
	AllowReject       *bool    `json:"allow_reject"`
	ConditionsJSON    *string  `json:"conditions_json"`
	SLAHours          *int     `json:"sla_hours"`
}

type FlowStepResponse struct {
	ID                string    `json:"id"`
	FlowID            string    `json:"flow_id"`
	StepOrder         int       `json:"step_order"`
	StepName          string    `json:"step_name"`
	ApproverType      string    `json:"approver_type"`
	HierarchyLevel    *int      `json:"hierarchy_level,omitempty"`
	OrganizationIDs   []string  `json:"organization_ids,omitempty"`
	RoleID            *string   `json:"role_id,omitempty"`
	ApproverUserID    *string   `json:"approver_user_id,omitempty"`
	ApprovalMode      string    `json:"approval_mode"`
	ParticipationType string    `json:"participation_type"`
	RequiredApprovals *int      `json:"required_approvals,omitempty"`
	AllowReject       bool      `json:"allow_reject"`
	ConditionsJSON    *string   `json:"conditions_json,omitempty"`
	SLAHours          *int      `json:"sla_hours,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// =========================================================================
// Approval Instance DTOs
// =========================================================================

type CreateInstanceRequest struct {
	Module     string `json:"module" binding:"required"`
	DocumentID string `json:"document_id" binding:"required,uuid"`
	FlowID     string `json:"flow_id" binding:"required,uuid"`
}

type InstanceResponse struct {
	ID          string             `json:"id"`
	Module      string             `json:"module"`
	DocumentID  string             `json:"document_id"`
	FlowID      string             `json:"flow_id"`
	FlowName    string             `json:"flow_name,omitempty"`
	Status      string             `json:"status"`
	CurrentStep int                `json:"current_step"`
	CreatedBy   *string            `json:"created_by,omitempty"`
	Steps       []FlowStepResponse `json:"steps,omitempty"`
	Actions     []ActionResponse   `json:"actions,omitempty"`
	Tasks       []TaskResponse     `json:"tasks,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// =========================================================================
// Approval Action DTOs
// =========================================================================

type SubmitActionRequest struct {
	Action string  `json:"action" binding:"required,oneof=APPROVE REJECT"`
	Note   *string `json:"note"`
}

type ActionResponse struct {
	ID          string    `json:"id"`
	InstanceID  string    `json:"instance_id"`
	StepOrder   int       `json:"step_order"`
	ActorUserID string    `json:"actor_user_id"`
	Action      string    `json:"action"`
	Note        *string   `json:"note,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// =========================================================================
// Approval Task DTOs
// =========================================================================

type TaskResponse struct {
	ID           string `json:"id"`
	InstanceID   string `json:"instance_id"`
	FlowName     string `json:"flow_name,omitempty"`
	StepOrder    int    `json:"step_order"`
	AssigneeType string `json:"assignee_type"`
	AssigneeID   string `json:"assignee_id"`
	Status       string `json:"status"`
	// Submitter* describe whoever created the instance this task belongs to.
	SubmitterName             string    `json:"submitter_name,omitempty"`
	SubmitterEmployeeCode     string    `json:"submitter_employee_code,omitempty"`
	SubmitterOrganizationName string    `json:"submitter_organization_name,omitempty"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

// =========================================================================
// Paginated Response
// =========================================================================

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// =========================================================================
// Converter Helpers
// =========================================================================

func (f *ApprovalFlow) ToResponse() FlowResponse {
	r := FlowResponse{
		ID:        f.ID.String(),
		Module:    f.Module,
		Name:      f.Name,
		Version:   f.Version,
		IsActive:  f.IsActive,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
	if len(f.Steps) > 0 {
		for _, s := range f.Steps {
			r.Steps = append(r.Steps, s.ToResponse())
		}
	}
	return r
}

func (s *ApprovalFlowStep) ToResponse() FlowStepResponse {
	r := FlowStepResponse{
		ID:                s.ID.String(),
		FlowID:            s.FlowID.String(),
		StepOrder:         s.StepOrder,
		StepName:          s.StepName,
		ApproverType:      string(s.ApproverType),
		HierarchyLevel:    s.HierarchyLevel,
		ApprovalMode:      string(s.ApprovalMode),
		ParticipationType: string(s.ParticipationType),
		AllowReject:       s.AllowReject,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
	if s.RoleID != nil {
		id := s.RoleID.String()
		r.RoleID = &id
	}
	if s.ApproverUserID != nil {
		id := s.ApproverUserID.String()
		r.ApproverUserID = &id
	}
	if s.RequiredApprovals != nil {
		r.RequiredApprovals = s.RequiredApprovals
	}
	if s.ConditionsJSON != nil {
		r.ConditionsJSON = s.ConditionsJSON
	}
	if s.SLAHours != nil {
		r.SLAHours = s.SLAHours
	}
	if len(s.Organizations) > 0 {
		for _, o := range s.Organizations {
			r.OrganizationIDs = append(r.OrganizationIDs, o.OrganizationID.String())
		}
	}
	return r
}

func (i *ApprovalInstance) ToResponse() InstanceResponse {
	r := InstanceResponse{
		ID:          i.ID.String(),
		Module:      i.Module,
		DocumentID:  i.DocumentID.String(),
		FlowID:      i.FlowID.String(),
		Status:      string(i.Status),
		CurrentStep: i.CurrentStep,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
	if i.CreatedBy != nil {
		uid := i.CreatedBy.String()
		r.CreatedBy = &uid
	}
	if i.Flow != nil {
		r.FlowName = i.Flow.Name
	}
	if len(i.Steps) > 0 {
		for _, s := range i.Steps {
			r.Steps = append(r.Steps, s.ToResponse())
		}
	}
	if len(i.Actions) > 0 {
		for _, a := range i.Actions {
			r.Actions = append(r.Actions, a.ToResponse())
		}
	}
	if len(i.Tasks) > 0 {
		for _, t := range i.Tasks {
			r.Tasks = append(r.Tasks, t.ToResponse())
		}
	}
	return r
}

func (a *ApprovalAction) ToResponse() ActionResponse {
	r := ActionResponse{
		ID:          a.ID.String(),
		InstanceID:  a.InstanceID.String(),
		StepOrder:   a.StepOrder,
		ActorUserID: a.ActorUserID.String(),
		Action:      string(a.Action),
		CreatedAt:   a.CreatedAt,
	}
	if a.Note != nil {
		r.Note = a.Note
	}
	return r
}

func (t *ApprovalTask) ToResponse() TaskResponse {
	r := TaskResponse{
		ID:           t.ID.String(),
		InstanceID:   t.InstanceID.String(),
		StepOrder:    t.StepOrder,
		AssigneeType: t.AssigneeType,
		AssigneeID:   t.AssigneeID.String(),
		Status:       string(t.Status),
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
	return r
}

// Helper untuk parse UUID pointer dari string pointer
func parseUUIDPtr(s *string) (*uuid.UUID, error) {
	if s == nil {
		return nil, nil
	}
	uid, err := uuid.Parse(*s)
	if err != nil {
		return nil, err
	}
	return &uid, nil
}
