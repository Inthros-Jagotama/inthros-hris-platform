package employeemovement

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MovementType enum untuk tipe pergerakan karyawan.
type MovementType string

const (
	MovementTypePromotion         MovementType = "promotion"
	MovementTypeDemotion          MovementType = "demotion"
	MovementTypeMutation          MovementType = "mutation"
	MovementTypeContractExtension MovementType = "contract_extension"
	MovementTypeStatusChange      MovementType = "status_change"
	MovementTypeRetirement        MovementType = "retirement"
	MovementTypeOffboarding       MovementType = "offboarding"
	MovementTypeOther             MovementType = "other"
)

// MovementStatus enum untuk status pergerakan.
type MovementStatus string

const (
	MovementStatusDraft           MovementStatus = "draft"
	MovementStatusPendingApproval MovementStatus = "pending_approval"
	MovementStatusApproved        MovementStatus = "approved"
	MovementStatusRejected        MovementStatus = "rejected" // keputusan plan §11.4: status terpisah utk approval ditolak
	MovementStatusExecuted        MovementStatus = "executed"
	MovementStatusCancelled       MovementStatus = "cancelled"
)

// ContractType enum untuk tipe kontrak.
type ContractType string

const (
	ContractTypePKWT  ContractType = "pkwt"
	ContractTypePKWTT ContractType = "pkwtt"
	ContractTypeDaily ContractType = "daily"
	ContractTypeOther ContractType = "other"
)

// ContractStatus enum untuk status kontrak.
type ContractStatus string

const (
	ContractStatusActive     ContractStatus = "active"
	ContractStatusExpired    ContractStatus = "expired"
	ContractStatusExtended   ContractStatus = "extended"
	ContractStatusTerminated ContractStatus = "terminated"
)

// =========================================================================
// 12.1 EmployeeMovement (Riwayat Pergerakan Karyawan)
// =========================================================================

type EmployeeMovement struct {
	ID                     uuid.UUID    `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID             uuid.UUID    `gorm:"type:char(36);not null;index:idx_emp_mvmt_employee" json:"employee_id"`
	MovementType           MovementType `gorm:"type:varchar(50);not null;index:idx_emp_mvmt_type" json:"movement_type"`
	FromEmploymentID       *uuid.UUID   `gorm:"type:char(36);index" json:"from_employment_id,omitempty"`
	ToEmploymentID         *uuid.UUID   `gorm:"type:char(36);index" json:"to_employment_id,omitempty"`
	FromOrganizationID     *uuid.UUID   `gorm:"type:char(36);index:idx_emp_mvmt_from_org" json:"from_organization_id,omitempty"`
	ToOrganizationID       *uuid.UUID   `gorm:"type:char(36);index:idx_emp_mvmt_to_org" json:"to_organization_id,omitempty"`
	FromPositionID         *uuid.UUID   `gorm:"type:char(36);index:idx_emp_mvmt_from_pos" json:"from_position_id,omitempty"`
	ToPositionID           *uuid.UUID   `gorm:"type:char(36);index:idx_emp_mvmt_to_pos" json:"to_position_id,omitempty"`
	FromEmploymentStatusID *uuid.UUID   `gorm:"type:char(36);index" json:"from_employment_status_id,omitempty"`
	ToEmploymentStatusID   *uuid.UUID   `gorm:"type:char(36);index" json:"to_employment_status_id,omitempty"`
	// Movement Snapshot (plan §12.5): nama master data saat movement dibuat.
	// Disimpan agar histori movement tidak berubah ketika nomenclature
	// organisasi / title posisi / nama employment status diubah di kemudian hari.
	FromOrganizationName     string         `gorm:"type:varchar(255)" json:"from_organization_name,omitempty"`
	FromPositionName         string         `gorm:"type:varchar(255)" json:"from_position_name,omitempty"`
	FromEmploymentStatusName string         `gorm:"type:varchar(255)" json:"from_employment_status_name,omitempty"`
	ToOrganizationName       string         `gorm:"type:varchar(255)" json:"to_organization_name,omitempty"`
	ToPositionName           string         `gorm:"type:varchar(255)" json:"to_position_name,omitempty"`
	ToEmploymentStatusName   string         `gorm:"type:varchar(255)" json:"to_employment_status_name,omitempty"`
	Reason                   *string        `gorm:"type:text" json:"reason,omitempty"`
	DecisionLetterNumber     string         `gorm:"type:varchar(50);not null" json:"decision_letter_number"`
	DecisionLetterDate       string         `gorm:"type:date;not null" json:"decision_letter_date"`
	EffectiveDate            string         `gorm:"type:date;not null;index:idx_emp_mvmt_effective" json:"effective_date"`
	Status                   MovementStatus `gorm:"type:varchar(20);default:draft;index:idx_emp_mvmt_status" json:"status"`
	Notes                    *string        `gorm:"type:text" json:"notes,omitempty"`
	ApprovedBy               *uuid.UUID     `gorm:"type:char(36)" json:"approved_by,omitempty"`
	ApprovedAt               *time.Time     `gorm:"type:timestamp" json:"approved_at,omitempty"`
	ExecutedBy               *uuid.UUID     `gorm:"type:char(36)" json:"executed_by,omitempty"`
	ExecutedAt               *time.Time     `gorm:"type:timestamp" json:"executed_at,omitempty"`
	ApprovalInstanceID       *uuid.UUID     `gorm:"type:char(36);index:idx_emp_mvmt_approval_instance" json:"approval_instance_id,omitempty"`
	CreatedBy                *uuid.UUID     `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy                *uuid.UUID     `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
}

func (EmployeeMovement) TableName() string {
	return "employee_movements"
}

func (m *EmployeeMovement) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 12.6 EmployeeMovementAudit (Movement Audit Trail)
// =========================================================================

// MovementAuditAction enum untuk aksi pada audit trail movement
// (plan §12.6: CREATED, UPDATED, SUBMITTED, APPROVED, REJECTED, CANCELLED,
// EXECUTED).
type MovementAuditAction string

const (
	MovementAuditActionCreated   MovementAuditAction = "CREATED"
	MovementAuditActionUpdated   MovementAuditAction = "UPDATED"
	MovementAuditActionSubmitted MovementAuditAction = "SUBMITTED"
	MovementAuditActionApproved  MovementAuditAction = "APPROVED"
	MovementAuditActionRejected  MovementAuditAction = "REJECTED"
	MovementAuditActionCancelled MovementAuditAction = "CANCELLED"
	MovementAuditActionExecuted  MovementAuditAction = "EXECUTED"
)

// EmployeeMovementAudit mencatat setiap perubahan lifecycle movement.
// OldData/NewData disimpan sebagai JSON string (pola OrganizationHistory /
// payroll before_json/after_json) berisi snapshot movement sebelum/sesudah
// aksi, sehingga histori tetap dapat direkonstruksi meskipun movement diubah
// atau data master berganti.
type EmployeeMovementAudit struct {
	ID         uuid.UUID           `gorm:"type:char(36);primaryKey" json:"id"`
	MovementID uuid.UUID           `gorm:"type:char(36);not null;index:idx_emp_mvmt_audit_movement" json:"movement_id"`
	Action     MovementAuditAction `gorm:"type:varchar(30);not null;index:idx_emp_mvmt_audit_action" json:"action"`
	OldStatus  *string             `gorm:"type:varchar(20)" json:"old_status,omitempty"`
	NewStatus  *string             `gorm:"type:varchar(20)" json:"new_status,omitempty"`
	OldData    *string             `gorm:"type:json" json:"old_data,omitempty"`
	NewData    *string             `gorm:"type:json" json:"new_data,omitempty"`
	Reason     *string             `gorm:"type:text" json:"reason,omitempty"`
	ActedBy    *uuid.UUID          `gorm:"type:char(36)" json:"acted_by,omitempty"`
	ActedAt    time.Time           `json:"acted_at"`
}

func (EmployeeMovementAudit) TableName() string {
	return "employee_movement_audits"
}

func (a *EmployeeMovementAudit) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.ActedAt.IsZero() {
		a.ActedAt = time.Now()
	}
	return nil
}

// =========================================================================
// 12.15 EmployeeMovementDocument (Movement Documents)
// =========================================================================

// MovementDocumentType enum untuk jenis dokumen movement
// (plan §12.15: PROMOTION_SK, MUTATION_SK, DEMOTION_SK, RETIREMENT_LETTER,
// OFFBOARDING_LETTER, OTHER).
type MovementDocumentType string

const (
	MovementDocumentTypePromotionSK       MovementDocumentType = "PROMOTION_SK"
	MovementDocumentTypeMutationSK        MovementDocumentType = "MUTATION_SK"
	MovementDocumentTypeDemotionSK        MovementDocumentType = "DEMOTION_SK"
	MovementDocumentTypeRetirementLetter  MovementDocumentType = "RETIREMENT_LETTER"
	MovementDocumentTypeOffboardingLetter MovementDocumentType = "OFFBOARDING_LETTER"
	MovementDocumentTypeOther             MovementDocumentType = "OTHER"
)

// EmployeeMovementDocument menyimpan metadata dokumen pergerakan karyawan
// (plan §12.15). File fisik di-upload via endpoint upload generik
// (POST /api/v1/tenant/uploads) yang mengembalikan file_url publik; tabel ini
// hanya menyimpan referensi URL + metadata agar satu movement dapat memiliki
// banyak dokumen (SK promosi, SK mutasi, surat pensiun, dll).
type EmployeeMovementDocument struct {
	ID           uuid.UUID            `gorm:"type:char(36);primaryKey" json:"id"`
	MovementID   uuid.UUID            `gorm:"type:char(36);not null;index:idx_emp_mvmt_doc_movement" json:"movement_id"`
	DocumentType MovementDocumentType `gorm:"type:varchar(30);not null" json:"document_type"`
	FileName     string               `gorm:"type:varchar(255);not null" json:"file_name"`
	FileURL      string               `gorm:"type:varchar(255);not null" json:"file_url"`
	UploadedBy   *uuid.UUID           `gorm:"type:char(36)" json:"uploaded_by,omitempty"`
	CreatedAt    time.Time            `json:"created_at"`
}

func (EmployeeMovementDocument) TableName() string {
	return "employee_movement_documents"
}

func (d *EmployeeMovementDocument) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 12.2 EmployeeContract (PKWT & Perjanjian Kerja)
// =========================================================================

type EmployeeContract struct {
	ID                   uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID           uuid.UUID      `gorm:"type:char(36);not null;index:idx_emp_ctrct_employee" json:"employee_id"`
	ContractNumber       string         `gorm:"type:varchar(50);not null" json:"contract_number"`
	ContractType         ContractType   `gorm:"type:varchar(20);not null;index:idx_emp_ctrct_type" json:"contract_type"`
	StartDate            string         `gorm:"type:date;not null" json:"start_date"`
	EndDate              *string        `gorm:"type:date" json:"end_date,omitempty"`
	ExtensionCount       int            `gorm:"type:int;default:0" json:"extension_count"`
	PreviousContractID   *uuid.UUID     `gorm:"type:char(36);index" json:"previous_contract_id,omitempty"`
	DecisionLetterNumber *string        `gorm:"type:varchar(50)" json:"decision_letter_number,omitempty"`
	Notes                *string        `gorm:"type:text" json:"notes,omitempty"`
	DocumentURL          *string        `gorm:"type:varchar(255)" json:"document_url,omitempty"`
	Status               ContractStatus `gorm:"type:varchar(20);default:active;index:idx_emp_ctrct_status" json:"status"`
	CreatedBy            *uuid.UUID     `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy            *uuid.UUID     `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
}

func (EmployeeContract) TableName() string {
	return "employee_contracts"
}

func (c *EmployeeContract) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
