package reimbursement

import "time"

// =========================================================================
// Reimbursement Type DTOs
// =========================================================================

type CreateReimbursementTypeRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

type UpdateReimbursementTypeRequest struct {
	Code        *string `json:"code"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type ReimbursementTypeResponse struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// =========================================================================
// Reimbursement Request DTOs
// =========================================================================

type CreateReimbursementRequest struct {
	RequestTypeID string  `json:"request_type_id" binding:"required"`
	Title         string  `json:"title" binding:"required"`
	Description   string  `json:"description"`
	Currency      string  `json:"currency"`
	SupervisorID  *string `json:"supervisor_id"`
}

type UpdateReimbursementRequest struct {
	Title        *string `json:"title"`
	Description  *string `json:"description"`
	Currency     *string `json:"currency"`
	SupervisorID *string `json:"supervisor_id"`
}

type UpdateReimbursementStatusRequest struct {
	Status           string   `json:"status" binding:"required,oneof=SUBMITTED APPROVED REJECTED PAID CANCELLED"`
	Note             string   `json:"note"`
	Amount           *float64 `json:"amount"`
	FlowID           *string  `json:"flow_id"`
	PaymentMethod    *string  `json:"payment_method" binding:"omitempty,oneof=BANK_TRANSFER CASH CHEQUE"`
	PaymentReference *string  `json:"payment_reference"`
	PaymentNote      *string  `json:"payment_note"`
}

type ReimbursementRequestResponse struct {
	ID                string     `json:"id"`
	EmployeeID        string     `json:"employee_id"`
	RequestTypeID     string     `json:"request_type_id"`
	Title             string     `json:"title"`
	Description       string     `json:"description,omitempty"`
	TotalAmount       float64    `json:"total_amount"`
	Currency          string     `json:"currency"`
	Status            string     `json:"status"`
	SupervisorID      *string    `json:"supervisor_id,omitempty"`
	SupervisorNote    *string    `json:"supervisor_note,omitempty"`
	SupervisorActionAt *time.Time `json:"supervisor_action_at,omitempty"`
	HrID              *string    `json:"hr_id,omitempty"`
	HrNote            *string    `json:"hr_note,omitempty"`
	HrActionAt        *time.Time `json:"hr_action_at,omitempty"`
	PaidAt            *time.Time `json:"paid_at,omitempty"`
	PaidAmount        *float64   `json:"paid_amount,omitempty"`
	PaymentMethod     *string    `json:"payment_method,omitempty"`
	PaymentReference  *string    `json:"payment_reference,omitempty"`
	PaymentNote       *string    `json:"payment_note,omitempty"`
	SubmittedAt       *time.Time `json:"submitted_at,omitempty"`
	ApprovedAt        *time.Time `json:"approved_at,omitempty"`
	RejectedAt        *time.Time `json:"rejected_at,omitempty"`
	CancelledAt       *time.Time `json:"cancelled_at,omitempty"`
	ApprovalInstanceID *string   `json:"approval_instance_id,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// =========================================================================
// Reimbursement Item DTOs
// =========================================================================

type CreateReimbursementItemRequest struct {
	ExpenseDate string  `json:"expense_date" binding:"required"`
	ExpenseType string  `json:"expense_type" binding:"required"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount" binding:"required"`
	ReceiptURL  *string `json:"receipt_url"`
}

type UpdateReimbursementItemRequest struct {
	ExpenseDate *string  `json:"expense_date"`
	ExpenseType *string  `json:"expense_type"`
	Description *string  `json:"description"`
	Amount      *float64 `json:"amount"`
	ReceiptURL  *string  `json:"receipt_url"`
}

type ReimbursementItemResponse struct {
	ID                      string    `json:"id"`
	ReimbursementRequestID  string    `json:"reimbursement_request_id"`
	ExpenseDate             string    `json:"expense_date"`
	ExpenseType             string    `json:"expense_type"`
	Description             string    `json:"description,omitempty"`
	Amount                  float64   `json:"amount"`
	ReceiptURL              *string   `json:"receipt_url,omitempty"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
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
