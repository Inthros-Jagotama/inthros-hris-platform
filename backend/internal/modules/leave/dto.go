package leave

import "time"

// =========================================================================
// Leave Type DTOs
// =========================================================================

type CreateLeaveTypeRequest struct {
	Code               string `json:"code"`
	Name               string `json:"name" binding:"required"`
	Description        string `json:"description"`
	IsPaid             *bool  `json:"is_paid"`
	RequiresAttachment *bool  `json:"requires_attachment"`
	AllowHalfDay       *bool  `json:"allow_half_day"`
	DefaultQuotaDays   *int   `json:"default_quota_days"`
	QuotaPeriod        string `json:"quota_period"`
	CountsAgainstQuota *bool  `json:"counts_against_quota"`
	AllowHourly        *bool  `json:"allow_hourly"`
	IsActive           *bool  `json:"is_active"`
}

type UpdateLeaveTypeRequest struct {
	Code               *string `json:"code"`
	Name               *string `json:"name"`
	Description        *string `json:"description"`
	IsPaid             *bool   `json:"is_paid"`
	RequiresAttachment *bool   `json:"requires_attachment"`
	AllowHalfDay       *bool   `json:"allow_half_day"`
	DefaultQuotaDays   *int    `json:"default_quota_days"`
	QuotaPeriod        *string `json:"quota_period"`
	CountsAgainstQuota *bool   `json:"counts_against_quota"`
	AllowHourly        *bool   `json:"allow_hourly"`
	IsActive           *bool   `json:"is_active"`
}

type LeaveTypeResponse struct {
	ID                 string    `json:"id"`
	Code               string    `json:"code"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	IsPaid             bool      `json:"is_paid"`
	RequiresAttachment bool      `json:"requires_attachment"`
	AllowHalfDay       bool      `json:"allow_half_day"`
	DefaultQuotaDays   *int      `json:"default_quota_days,omitempty"`
	QuotaPeriod        string    `json:"quota_period"`
	CountsAgainstQuota bool      `json:"counts_against_quota"`
	AllowHourly        bool      `json:"allow_hourly"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// =========================================================================
// Accrual Policy DTOs
// =========================================================================

type CreateAccrualPolicyRequest struct {
	LeaveTypeID     string  `json:"leave_type_id" binding:"required"`
	BaseQuotaDays   float64 `json:"base_quota_days" binding:"required"`
	ExtraEveryYears *int    `json:"extra_every_years"`
	ExtraDays       *float64 `json:"extra_days"`
	MaxExtraDays    *float64 `json:"max_extra_days"`
	EffectiveFrom   string  `json:"effective_from" binding:"required"`
	EffectiveTo     *string `json:"effective_to"`
}

type UpdateAccrualPolicyRequest struct {
	LeaveTypeID     *string  `json:"leave_type_id"`
	BaseQuotaDays   *float64 `json:"base_quota_days"`
	ExtraEveryYears *int     `json:"extra_every_years"`
	ExtraDays       *float64 `json:"extra_days"`
	MaxExtraDays    *float64 `json:"max_extra_days"`
	EffectiveFrom   *string  `json:"effective_from"`
	EffectiveTo     *string  `json:"effective_to"`
}

type AccrualPolicyResponse struct {
	ID              string    `json:"id"`
	LeaveTypeID     string    `json:"leave_type_id"`
	BaseQuotaDays   float64   `json:"base_quota_days"`
	ExtraEveryYears int       `json:"extra_every_years"`
	ExtraDays       float64   `json:"extra_days"`
	MaxExtraDays    *float64  `json:"max_extra_days,omitempty"`
	EffectiveFrom   string    `json:"effective_from"`
	EffectiveTo     *string   `json:"effective_to,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// =========================================================================
// Leave Reason DTOs
// =========================================================================

type CreateLeaveReasonRequest struct {
	Name      string `json:"name" binding:"required"`
	IsOther   *bool  `json:"is_other"`
	SortOrder *int   `json:"sort_order"`
}

type UpdateLeaveReasonRequest struct {
	Name      *string `json:"name"`
	IsOther   *bool   `json:"is_other"`
	SortOrder *int    `json:"sort_order"`
}

type LeaveReasonResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	IsOther   bool      `json:"is_other"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// =========================================================================
// Leave Request DTOs
// =========================================================================

type CreateLeaveRequest struct {
	EmployeeID      string  `json:"employee_id" binding:"required"`
	LeaveTypeID     string  `json:"leave_type_id" binding:"required"`
	RequestStartDate string `json:"request_start_date" binding:"required"`
	RequestEndDate  string  `json:"request_end_date" binding:"required"`
	DurationMode    string  `json:"duration_mode"`
	RequestedDays   float64 `json:"requested_days"`
	LeaveReasonID   *string `json:"leave_reason_id"`
	LeaveReasonNote *string `json:"leave_reason_note"`
	AttachmentURL   *string `json:"attachment_url"`
	SupervisorID    *string `json:"supervisor_id"`
	StartTime       *string `json:"start_time"`
	EndTime         *string `json:"end_time"`
	FlowID          *string `json:"flow_id"`
}

type UpdateLeaveRequestStatus struct {
	Status string `json:"status" binding:"required"`
	Note   string `json:"note"`
}

type LeaveRequestResponse struct {
	ID                string     `json:"id"`
	EmployeeID        string     `json:"employee_id"`
	LeaveTypeID       string     `json:"leave_type_id"`
	RequestStartDate  string     `json:"request_start_date"`
	RequestEndDate    string     `json:"request_end_date"`
	DurationMode      string     `json:"duration_mode"`
	RequestedDays     float64    `json:"requested_days"`
	LeaveReasonID     *string    `json:"leave_reason_id,omitempty"`
	LeaveReasonNote   *string    `json:"leave_reason_note,omitempty"`
	AttachmentURL     *string    `json:"attachment_url,omitempty"`
	Status            string     `json:"status"`
	SupervisorID      *string    `json:"supervisor_id,omitempty"`
	SupervisorNote    *string    `json:"supervisor_note,omitempty"`
	HrID              *string    `json:"hr_id,omitempty"`
	HrNote            *string    `json:"hr_note,omitempty"`
	ApprovalInstanceID *string   `json:"approval_instance_id,omitempty"`
	StartTime         *string    `json:"start_time,omitempty"`
	EndTime           *string    `json:"end_time,omitempty"`
	SubmittedAt       *time.Time `json:"submitted_at,omitempty"`
	ApprovedAt        *time.Time `json:"approved_at,omitempty"`
	RejectedAt        *time.Time `json:"rejected_at,omitempty"`
	CancelledAt       *time.Time `json:"cancelled_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// =========================================================================
// Leave Request Detail DTOs
// =========================================================================

type LeaveRequestDetailResponse struct {
	ID                string    `json:"id"`
	LeaveRequestID    string    `json:"leave_request_id"`
	EmployeeID        string    `json:"employee_id"`
	LeaveDate         string    `json:"leave_date"`
	DayFraction       float64   `json:"day_fraction"`
	IsPaid            bool      `json:"is_paid"`
	CreatedAt         time.Time `json:"created_at"`
}

// =========================================================================
// Employee Leave Balance DTOs
// =========================================================================

type LeaveBalanceResponse struct {
	ID                string    `json:"id"`
	EmployeeID        string    `json:"employee_id"`
	LeaveTypeID       string    `json:"leave_type_id"`
	PeriodYear        int       `json:"period_year"`
	QuotaDays         float64   `json:"quota_days"`
	UsedDays          float64   `json:"used_days"`
	RemainingDays     float64   `json:"remaining_days"`
	LastAdjustmentRef *string   `json:"last_adjustment_ref,omitempty"`
	LastAdjustmentAt  *time.Time `json:"last_adjustment_at,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// =========================================================================
// Pagination
// =========================================================================

type PaginatedResponse struct {
	Success    bool          `json:"success"`
	Data       interface{}   `json:"data"`
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
	Total      int64         `json:"total"`
	TotalPages int           `json:"total_pages"`
}
