package approval

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =========================================================================
// Type Definitions
// =========================================================================

// ApproverType enum for step approver type.
type ApproverType string

const (
	ApproverTypeSupervisor ApproverType = "SUPERVISOR"
	ApproverTypeRole       ApproverType = "ROLE"
	ApproverTypeUser       ApproverType = "USER"
)

// ApprovalMode enum for step approval mode.
type ApprovalMode string

const (
	ApprovalModeAnyOne ApprovalMode = "ANY_ONE"
	ApprovalModeAll    ApprovalMode = "ALL"
	ApprovalModeNOfM   ApprovalMode = "N_OF_M"
)

// ApprovalActionType enum for action type.
type ApprovalActionType string

const (
	ActionApprove ApprovalActionType = "APPROVE"
	ActionReject  ApprovalActionType = "REJECT"
	ActionCancel  ApprovalActionType = "CANCEL"
)

// InstanceStatus enum for instance status.
type InstanceStatus string

const (
	InstanceStatusPending   InstanceStatus = "PENDING"
	InstanceStatusApproved  InstanceStatus = "APPROVED"
	InstanceStatusRejected  InstanceStatus = "REJECTED"
	InstanceStatusCancelled InstanceStatus = "CANCELLED"
)

// TaskStatus enum for task status.
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "PENDING"
	TaskStatusDone      TaskStatus = "DONE"
	TaskStatusCancelled TaskStatus = "CANCELLED"
)

// =========================================================================
// 10.1 ApprovalFlow (Master alur persetujuan)
// =========================================================================

type ApprovalFlow struct {
	ID        uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	Module    string     `gorm:"type:varchar(255);not null;index:idx_approval_flow_module" json:"module"`
	Name      string     `gorm:"type:varchar(255);not null" json:"name"`
	Version   int        `gorm:"type:int;default:1" json:"version"`
	IsActive  bool       `gorm:"type:tinyint(1);default:1;index:idx_approval_flow_active" json:"is_active"`
	DeletedAt *time.Time `gorm:"index:idx_approval_flow_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Relations (not stored)
	Steps []ApprovalFlowStep `gorm:"foreignKey:FlowID" json:"steps,omitempty"`
}

func (ApprovalFlow) TableName() string {
	return "approval_flows"
}

func (f *ApprovalFlow) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 10.2 ApprovalFlowStep (Langkah-langkah dalam alur)
// =========================================================================

type ApprovalFlowStep struct {
	ID                uuid.UUID     `gorm:"type:char(36);primaryKey" json:"id"`
	FlowID            uuid.UUID     `gorm:"type:char(36);not null;index:idx_approval_step_flow" json:"flow_id"`
	StepOrder         int           `gorm:"type:int;not null" json:"step_order"`
	StepName          string        `gorm:"type:varchar(100);not null" json:"step_name"`
	ApproverType      ApproverType  `gorm:"type:varchar(255);not null" json:"approver_type"`
	RoleID            *uuid.UUID    `gorm:"type:char(36);index:idx_approval_step_role" json:"role_id,omitempty"`
	ApproverUserID    *uuid.UUID    `gorm:"type:char(36);index:idx_approval_step_user" json:"approver_user_id,omitempty"`
	ApprovalMode      ApprovalMode  `gorm:"type:varchar(255);not null;default:ANY_ONE" json:"approval_mode"`
	RequiredApprovals *int          `gorm:"type:int" json:"required_approvals,omitempty"`
	AllowReject       bool          `gorm:"type:tinyint(1);default:1" json:"allow_reject"`
	ConditionsJSON    *string       `gorm:"type:json" json:"conditions_json,omitempty"`
	SLAHours          *int          `gorm:"type:int" json:"sla_hours,omitempty"`
	DeletedAt         *time.Time    `gorm:"index:idx_approval_step_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

func (ApprovalFlowStep) TableName() string {
	return "approval_flow_steps"
}

func (s *ApprovalFlowStep) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 10.3 ApprovalInstance (Instance persetujuan per dokumen)
// =========================================================================

type ApprovalInstance struct {
	ID          uuid.UUID       `gorm:"type:char(36);primaryKey" json:"id"`
	Module      string          `gorm:"type:varchar(50);not null" json:"module"`
	DocumentID  uuid.UUID       `gorm:"type:char(36);not null" json:"document_id"`
	FlowID      uuid.UUID       `gorm:"type:char(36);not null;index:idx_approval_instance_flow" json:"flow_id"`
	Status      InstanceStatus  `gorm:"type:varchar(255);not null;default:PENDING" json:"status"`
	CurrentStep int             `gorm:"type:int;default:1" json:"current_step"`
	CreatedBy   *uuid.UUID      `gorm:"type:char(36);index:idx_approval_instance_creator" json:"created_by,omitempty"`
	DeletedAt   *time.Time      `gorm:"index:idx_approval_instance_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`

	// Relations (not stored)
	Flow   *ApprovalFlow   `gorm:"foreignKey:FlowID" json:"flow,omitempty"`
	Steps  []ApprovalFlowStep `gorm:"-" json:"steps,omitempty"`
	Actions []ApprovalAction `gorm:"foreignKey:InstanceID" json:"actions,omitempty"`
	Tasks  []ApprovalTask  `gorm:"foreignKey:InstanceID" json:"tasks,omitempty"`
}

func (ApprovalInstance) TableName() string {
	return "approval_instances"
}

func (i *ApprovalInstance) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 10.4 ApprovalAction (Aksi yang dilakukan approver)
// =========================================================================

type ApprovalAction struct {
	ID          uuid.UUID           `gorm:"type:char(36);primaryKey" json:"id"`
	InstanceID  uuid.UUID           `gorm:"type:char(36);not null;index:idx_approval_action_instance" json:"instance_id"`
	StepOrder   int                 `gorm:"type:int;not null" json:"step_order"`
	ActorUserID uuid.UUID           `gorm:"type:char(36);not null;index:idx_approval_action_actor" json:"actor_user_id"`
	Action      ApprovalActionType  `gorm:"type:varchar(255);not null" json:"action"`
	Note        *string             `gorm:"type:varchar(255)" json:"note,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
}

func (ApprovalAction) TableName() string {
	return "approval_actions"
}

func (a *ApprovalAction) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 10.5 ApprovalTask (Task approval per approver)
// =========================================================================

type ApprovalTask struct {
	ID           uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	InstanceID   uuid.UUID  `gorm:"type:char(36);not null;index:idx_approval_task_instance" json:"instance_id"`
	StepOrder    int        `gorm:"type:int;not null" json:"step_order"`
	AssigneeType string     `gorm:"type:varchar(255);not null" json:"assignee_type"`
	AssigneeID   uuid.UUID  `gorm:"type:char(36);not null;index:idx_approval_task_assignee" json:"assignee_id"`
	Status       TaskStatus `gorm:"type:varchar(255);not null;default:PENDING" json:"status"`
	DeletedAt    *time.Time `gorm:"index:idx_approval_task_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (ApprovalTask) TableName() string {
	return "approval_tasks"
}

func (t *ApprovalTask) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
