package leave

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =========================================================================
// 5.1 LeaveType (Jenis Cuti)
// =========================================================================

type QuotaPeriod string

const (
	QuotaPeriodYear  QuotaPeriod = "YEAR"
	QuotaPeriodMonth QuotaPeriod = "MONTH"
	QuotaPeriodNone  QuotaPeriod = "NONE"
)

type LeaveType struct {
	ID                 uuid.UUID   `gorm:"type:char(36);primaryKey" json:"id"`
	Code               string      `gorm:"type:varchar(50);not null;default:''" json:"code"`
	Name               string      `gorm:"type:varchar(50);not null;uniqueIndex:uk_leave_type_name" json:"name"`
	Description        string      `gorm:"type:varchar(200);not null" json:"description"`
	IsPaid             bool        `gorm:"not null;default:1" json:"is_paid"`
	RequiresAttachment bool        `gorm:"not null;default:0" json:"requires_attachment"`
	AllowHalfDay       bool        `gorm:"not null;default:1" json:"allow_half_day"`
	DefaultQuotaDays   *int        `json:"default_quota_days,omitempty"`
	QuotaPeriod        QuotaPeriod `gorm:"type:varchar(255);not null;default:YEAR" json:"quota_period"`
	CountsAgainstQuota bool        `gorm:"not null;default:1" json:"counts_against_quota"`
	AllowHourly        bool        `gorm:"not null;default:0" json:"allow_hourly"`
	IsActive           bool        `gorm:"not null;default:1;index:idx_leave_type_active" json:"is_active"`
	DeletedAt          gorm.DeletedAt `gorm:"index:idx_leave_type_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
}

func (LeaveType) TableName() string {
	return "leave_types"
}

func (t *LeaveType) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 5.2 LeaveAccrualPolicy (Kebijakan Akrual Cuti)
// =========================================================================

type LeaveAccrualPolicy struct {
	ID             uuid.UUID     `gorm:"type:char(36);primaryKey" json:"id"`
	LeaveTypeID    uuid.UUID     `gorm:"type:char(36);not null;uniqueIndex:uk_accrual_policy;index:idx_accrual_leave_type" json:"leave_type_id"`
	BaseQuotaDays  float64       `gorm:"type:decimal(6,2);not null" json:"base_quota_days"`
	ExtraEveryYears int          `gorm:"not null;default:2" json:"extra_every_years"`
	ExtraDays      float64       `gorm:"type:decimal(6,2);not null;default:1.00" json:"extra_days"`
	MaxExtraDays   *float64      `gorm:"type:decimal(6,2)" json:"max_extra_days,omitempty"`
	EffectiveFrom  string        `gorm:"type:date;not null;uniqueIndex:uk_accrual_policy" json:"effective_from"`
	EffectiveTo    *string       `gorm:"type:date" json:"effective_to,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

func (LeaveAccrualPolicy) TableName() string {
	return "leave_accrual_policies"
}

func (p *LeaveAccrualPolicy) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 5.3 LeaveReason (Alasan Cuti)
// =========================================================================

type LeaveReason struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	IsOther   bool           `gorm:"not null;default:0" json:"is_other"`
	SortOrder int            `gorm:"not null;default:0" json:"sort_order"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_leave_reason_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (LeaveReason) TableName() string {
	return "leave_reasons"
}

func (r *LeaveReason) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 5.4 LeaveRequest (Pengajuan Cuti)
// =========================================================================

type DurationMode string

const (
	DurationFullDay   DurationMode = "FULL_DAY"
	DurationHalfDayAM DurationMode = "HALF_DAY_AM"
	DurationHalfDayPM DurationMode = "HALF_DAY_PM"
	DurationHourly    DurationMode = "HOURLY"
)

type LeaveStatus string

const (
	LeaveStatusDraft           LeaveStatus = "DRAFT"
	LeaveStatusSubmitted       LeaveStatus = "SUBMITTED"
	LeaveStatusPendingApproval LeaveStatus = "PENDING_APPROVAL"
	LeaveStatusApprovedFinal   LeaveStatus = "APPROVED_FINAL"
	LeaveStatusRejectedFinal   LeaveStatus = "REJECTED_FINAL"
	LeaveStatusCancelled       LeaveStatus = "CANCELLED"
)

type LeaveRequest struct {
	ID                uuid.UUID    `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID        uuid.UUID    `gorm:"type:char(36);not null;index:idx_leave_req_employee_date,idx_leave_req_employee_created" json:"employee_id"`
	LeaveTypeID       uuid.UUID    `gorm:"type:char(36);not null;index:idx_leave_req_type" json:"leave_type_id"`
	RequestStartDate  string       `gorm:"type:date;not null" json:"request_start_date"`
	RequestEndDate    string       `gorm:"type:date;not null" json:"request_end_date"`
	DurationMode      DurationMode `gorm:"type:varchar(255);not null;default:FULL_DAY" json:"duration_mode"`
	RequestedDays     float64      `gorm:"type:decimal(6,2);not null;default:0" json:"requested_days"`
	LeaveReasonID     *uuid.UUID   `gorm:"type:char(36);index:idx_leave_req_reason" json:"leave_reason_id,omitempty"`
	LeaveReasonNote   *string      `gorm:"type:varchar(255)" json:"leave_reason_note,omitempty"`
	AttachmentURL     *string      `gorm:"type:text" json:"attachment_url,omitempty"`
	Status            LeaveStatus  `gorm:"type:varchar(255);not null;default:SUBMITTED;index:idx_leave_req_status" json:"status"`
	SupervisorID      *uuid.UUID   `gorm:"type:char(36)" json:"supervisor_id,omitempty"`
	SupervisorActionAt *time.Time  `gorm:"type:timestamp(6)" json:"supervisor_action_at,omitempty"`
	SupervisorNote    *string      `gorm:"type:varchar(255)" json:"supervisor_note,omitempty"`
	HrID              *uuid.UUID   `gorm:"type:char(36)" json:"hr_id,omitempty"`
	HrActionAt        *time.Time   `gorm:"type:timestamp(6)" json:"hr_action_at,omitempty"`
	HrNote            *string      `gorm:"type:varchar(255)" json:"hr_note,omitempty"`
	ApprovalInstanceID *uuid.UUID  `gorm:"type:char(36);index:idx_leave_req_approval_instance" json:"approval_instance_id,omitempty"`
	StartTime         *string      `gorm:"type:time" json:"start_time,omitempty"`
	EndTime           *string      `gorm:"type:time" json:"end_time,omitempty"`
	SubmittedAt       *time.Time   `gorm:"type:timestamp(6)" json:"submitted_at,omitempty"`
	ApprovedAt        *time.Time   `gorm:"type:timestamp(6)" json:"approved_at,omitempty"`
	RejectedAt        *time.Time   `gorm:"type:timestamp(6)" json:"rejected_at,omitempty"`
	CancelledAt       *time.Time   `gorm:"type:timestamp(6)" json:"cancelled_at,omitempty"`
	DeletedAt         gorm.DeletedAt `gorm:"index:idx_leave_req_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

func (LeaveRequest) TableName() string {
	return "leave_requests"
}

func (r *LeaveRequest) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 5.5 LeaveRequestDetail (Detail Hari Cuti)
// =========================================================================

type LeaveRequestDetail struct {
	ID                 uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	LeaveRequestID     uuid.UUID      `gorm:"type:char(36);not null;uniqueIndex:uk_leave_req_day;index:idx_leave_detail_request" json:"leave_request_id"`
	EmployeeID         uuid.UUID      `gorm:"type:char(36);not null;index:idx_leave_detail_emp_date" json:"employee_id"`
	LeaveDate          string         `gorm:"type:date;not null;uniqueIndex:uk_leave_req_day" json:"leave_date"`
	DayFraction        float64        `gorm:"type:decimal(4,2);not null;default:1.00" json:"day_fraction"`
	IsPaid             bool           `gorm:"not null;default:1" json:"is_paid"`
	ApprovalInstanceID *uuid.UUID     `gorm:"type:char(36)" json:"approval_instance_id,omitempty"`
	DeletedAt          gorm.DeletedAt `gorm:"index:idx_leave_detail_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

func (LeaveRequestDetail) TableName() string {
	return "leave_request_details"
}

func (d *LeaveRequestDetail) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 5.6 EmployeeLeaveBalance (Saldo Cuti Karyawan)
// =========================================================================

type EmployeeLeaveBalance struct {
	ID                  uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID          uuid.UUID      `gorm:"type:char(36);not null;uniqueIndex:uk_leave_balance" json:"employee_id"`
	LeaveTypeID         uuid.UUID      `gorm:"type:char(36);not null;uniqueIndex:uk_leave_balance;index:idx_leave_balance_type" json:"leave_type_id"`
	PeriodYear          int            `gorm:"not null;uniqueIndex:uk_leave_balance" json:"period_year"`
	QuotaDays           float64        `gorm:"type:decimal(6,2);not null;default:0" json:"quota_days"`
	UsedDays            float64        `gorm:"type:decimal(6,2);not null;default:0" json:"used_days"`
	RemainingDays       float64        `gorm:"type:decimal(6,2);not null;default:0" json:"remaining_days"`
	LastAdjustmentRef   *string        `gorm:"type:varchar(50)" json:"last_adjustment_ref,omitempty"`
	LastAdjustmentRefID *uuid.UUID     `gorm:"type:char(36)" json:"last_adjustment_ref_id,omitempty"`
	LastAdjustmentAt    *time.Time     `gorm:"type:timestamp(6)" json:"last_adjustment_at,omitempty"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
}

func (EmployeeLeaveBalance) TableName() string {
	return "employee_leave_balances"
}

func (b *EmployeeLeaveBalance) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
