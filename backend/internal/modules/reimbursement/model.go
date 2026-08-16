package reimbursement

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =========================================================================
// ReimbursementType (Jenis Klaim / Penggantian)
// =========================================================================

type ReimbursementType struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code        string         `gorm:"type:varchar(50);not null;default:''" json:"code"`
	Name        string         `gorm:"type:varchar(150);not null" json:"name"`
	Description string         `gorm:"type:varchar(500);not null" json:"description"`
	IsActive    bool           `gorm:"not null;default:1;index:idx_reimb_type_active" json:"is_active"`
	DeletedAt   gorm.DeletedAt `gorm:"index:idx_reimb_type_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (ReimbursementType) TableName() string {
	return "reimbursement_types"
}

func (t *ReimbursementType) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// ReimbursementStatus (Status Pengajuan Klaim)
// =========================================================================

type ReimbursementStatus string

const (
	ReimbStatusDraft           ReimbursementStatus = "DRAFT"
	ReimbStatusSubmitted       ReimbursementStatus = "SUBMITTED"
	ReimbStatusPendingApproval ReimbursementStatus = "PENDING_APPROVAL"
	ReimbStatusApproved        ReimbursementStatus = "APPROVED"
	ReimbStatusRejected        ReimbursementStatus = "REJECTED"
	ReimbStatusPaid            ReimbursementStatus = "PAID"
	ReimbStatusCancelled       ReimbursementStatus = "CANCELLED"
)

// =========================================================================
// ReimbursementRequest (Pengajuan Klaim)
// =========================================================================

type ReimbursementRequest struct {
	ID                uuid.UUID            `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID        uuid.UUID            `gorm:"type:char(36);not null;index:idx_reimb_req_employee" json:"employee_id"`
	RequestTypeID     uuid.UUID            `gorm:"type:char(36);not null;index:idx_reimb_req_type" json:"request_type_id"`
	Title             string               `gorm:"type:varchar(200);not null" json:"title"`
	Description       string               `gorm:"type:text" json:"description"`
	TotalAmount       float64              `gorm:"type:decimal(18,2);not null;default:0" json:"total_amount"`
	Currency          string               `gorm:"type:varchar(3);not null;default:IDR" json:"currency"`
	Status            ReimbursementStatus  `gorm:"type:varchar(50);not null;default:DRAFT;index:idx_reimb_req_status" json:"status"`
	SupervisorID      *uuid.UUID           `gorm:"type:char(36)" json:"supervisor_id,omitempty"`
	// Action timestamps are *time.Time (nullable) to match the SQL schema
	// (TIMESTAMP NULL) — the previous int64 unix-nano mapping made GORM write
	// 0 into the TIMESTAMP column, which MySQL strict mode rejects as
	// '0000-00-00 00:00:00' (Error 1292). No gorm type tag: an explicit
	// 'timestamp(6)' tag breaks *time.Time scanning on the SQLite test driver,
	// and tenant DBs are SQL-migrated anyway (GORM AutoMigrate never runs for
	// tenants), so the default datetime mapping is safe.
	SupervisorActionAt *time.Time  `json:"-"`
	SupervisorNote    *string     `gorm:"type:varchar(500)" json:"supervisor_note,omitempty"`
	HrID              *uuid.UUID  `gorm:"type:char(36)" json:"hr_id,omitempty"`
	HrActionAt        *time.Time  `json:"-"`
	HrNote            *string     `gorm:"type:varchar(500)" json:"hr_note,omitempty"`
	PaidAt            *time.Time  `json:"-"`
	PaidAmount        *float64    `gorm:"type:decimal(18,2)" json:"paid_amount,omitempty"`
	// Payment details recorded directly in this module (no payroll linkage —
	// product decision 2026-08-16). Values follow the payroll convention:
	// BANK_TRANSFER / CASH / CHEQUE.
	PaymentMethod     *string     `gorm:"type:varchar(50)" json:"payment_method,omitempty"`
	PaymentReference  *string     `gorm:"type:varchar(200)" json:"payment_reference,omitempty"`
	PaymentNote       *string     `gorm:"type:varchar(500)" json:"payment_note,omitempty"`
	SubmittedAt       *time.Time  `json:"-"`
	ApprovedAt        *time.Time  `json:"-"`
	RejectedAt        *time.Time  `json:"-"`
	CancelledAt       *time.Time  `json:"-"`
	ApprovalInstanceID *uuid.UUID          `gorm:"type:char(36);index:idx_reimb_req_approval_instance" json:"approval_instance_id,omitempty"`
	Items             []ReimbursementItem  `gorm:"foreignKey:ReimbursementRequestID" json:"items,omitempty"`
	DeletedAt         gorm.DeletedAt       `gorm:"index:idx_reimb_req_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
}

func (ReimbursementRequest) TableName() string {
	return "reimbursement_requests"
}

func (r *ReimbursementRequest) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// ReimbursementItem (Item Klaim — per pengeluaran)
// =========================================================================

type ReimbursementItem struct {
	ID                      uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	ReimbursementRequestID  uuid.UUID      `gorm:"type:char(36);not null;index:idx_reimb_item_request" json:"reimbursement_request_id"`
	ExpenseDate             string         `gorm:"type:date;not null" json:"expense_date"`
	ExpenseType             string         `gorm:"type:varchar(100);not null" json:"expense_type"`
	Description             string         `gorm:"type:varchar(500)" json:"description"`
	Amount                  float64        `gorm:"type:decimal(18,2);not null" json:"amount"`
	ReceiptURL              *string        `gorm:"type:text" json:"receipt_url,omitempty"`
	DeletedAt               gorm.DeletedAt `gorm:"index:idx_reimb_item_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
}

func (ReimbursementItem) TableName() string {
	return "reimbursement_items"
}

func (i *ReimbursementItem) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}
