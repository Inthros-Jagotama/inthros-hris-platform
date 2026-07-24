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
	SupervisorActionAt int64               `gorm:"type:bigint;default:0" json:"-"`
	SupervisorNote    *string              `gorm:"type:varchar(500)" json:"supervisor_note,omitempty"`
	HrID              *uuid.UUID           `gorm:"type:char(36)" json:"hr_id,omitempty"`
	HrActionAt        int64                `gorm:"type:bigint;default:0" json:"-"`
	HrNote            *string              `gorm:"type:varchar(500)" json:"hr_note,omitempty"`
	PaidAt            int64                `gorm:"type:bigint;default:0" json:"-"`
	PaidAmount        *float64             `gorm:"type:decimal(18,2)" json:"paid_amount,omitempty"`
	SubmittedAt       int64                `gorm:"type:bigint;default:0" json:"-"`
	ApprovedAt        int64                `gorm:"type:bigint;default:0" json:"-"`
	RejectedAt        int64                `gorm:"type:bigint;default:0" json:"-"`
	CancelledAt       int64                `gorm:"type:bigint;default:0" json:"-"`
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
