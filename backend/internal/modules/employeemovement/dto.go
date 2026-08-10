package employeemovement

import "time"

// =========================================================================
// Employee Movement DTOs
// =========================================================================

type CreateMovementRequest struct {
	EmployeeID             string  `json:"employee_id" binding:"required,uuid"`
	MovementType           string  `json:"movement_type" binding:"required,oneof=promotion demotion mutation contract_extension status_change retirement offboarding other"`
	FromEmploymentID       *string `json:"from_employment_id" binding:"omitempty,uuid"`
	ToEmploymentID         *string `json:"to_employment_id" binding:"omitempty,uuid"`
	FromOrganizationID     *string `json:"from_organization_id" binding:"omitempty,uuid"`
	ToOrganizationID       *string `json:"to_organization_id" binding:"omitempty,uuid"`
	FromPositionID         *string `json:"from_position_id" binding:"omitempty,uuid"`
	ToPositionID           *string `json:"to_position_id" binding:"omitempty,uuid"`
	FromEmploymentStatusID *string `json:"from_employment_status_id" binding:"omitempty,uuid"`
	ToEmploymentStatusID   *string `json:"to_employment_status_id" binding:"omitempty,uuid"`
	Reason                 *string `json:"reason"`
	DecisionLetterNumber   string  `json:"decision_letter_number" binding:"required"`
	DecisionLetterDate     string  `json:"decision_letter_date" binding:"required"`
	EffectiveDate          string  `json:"effective_date" binding:"required"`
	Notes                  *string `json:"notes"`
}

type UpdateMovementRequest struct {
	MovementType         *string `json:"movement_type" binding:"omitempty,oneof=promotion demotion mutation contract_extension status_change retirement offboarding other"`
	ToOrganizationID     *string `json:"to_organization_id" binding:"omitempty,uuid"`
	ToPositionID         *string `json:"to_position_id" binding:"omitempty,uuid"`
	ToEmploymentStatusID *string `json:"to_employment_status_id" binding:"omitempty,uuid"`
	Reason               *string `json:"reason"`
	DecisionLetterNumber *string `json:"decision_letter_number"`
	DecisionLetterDate   *string `json:"decision_letter_date"`
	EffectiveDate        *string `json:"effective_date"`
	Status               *string `json:"status" binding:"omitempty,oneof=draft approved executed cancelled"`
	Notes                *string `json:"notes"`
}

type MovementResponse struct {
	ID                       string     `json:"id"`
	EmployeeID               string     `json:"employee_id"`
	EmployeeName             string     `json:"employee_name"`
	EmployeeCode             string     `json:"employee_code"`
	MovementType             string     `json:"movement_type"`
	FromEmploymentID         *string    `json:"from_employment_id,omitempty"`
	ToEmploymentID           *string    `json:"to_employment_id,omitempty"`
	FromOrganizationID       *string    `json:"from_organization_id,omitempty"`
	ToOrganizationID         *string    `json:"to_organization_id,omitempty"`
	FromOrganizationName     string     `json:"from_organization_name,omitempty"`
	ToOrganizationName       string     `json:"to_organization_name,omitempty"`
	FromPositionID           *string    `json:"from_position_id,omitempty"`
	ToPositionID             *string    `json:"to_position_id,omitempty"`
	FromPositionName         string     `json:"from_position_name,omitempty"`
	ToPositionName           string     `json:"to_position_name,omitempty"`
	FromEmploymentStatusID   *string    `json:"from_employment_status_id,omitempty"`
	ToEmploymentStatusID     *string    `json:"to_employment_status_id,omitempty"`
	FromEmploymentStatusName string     `json:"from_employment_status_name,omitempty"`
	ToEmploymentStatusName   string     `json:"to_employment_status_name,omitempty"`
	Reason                   *string    `json:"reason,omitempty"`
	DecisionLetterNumber     string     `json:"decision_letter_number"`
	DecisionLetterDate       string     `json:"decision_letter_date"`
	EffectiveDate            string     `json:"effective_date"`
	Status                   string     `json:"status"`
	Notes                    *string    `json:"notes,omitempty"`
	ApprovedBy               *string    `json:"approved_by,omitempty"`
	ApprovedAt               *time.Time `json:"approved_at,omitempty"`
	ExecutedBy               *string    `json:"executed_by,omitempty"`
	ExecutedAt               *time.Time `json:"executed_at,omitempty"`
	ApprovalInstanceID       *string    `json:"approval_instance_id,omitempty"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}

// SubmitMovementRequest routes a draft movement through the central
// approval module — the single approval path (manual approve dihapus, G-5).
type SubmitMovementRequest struct {
	FlowID *string `json:"flow_id"`
}

// MovementAuditResponse membawa satu baris audit trail movement (plan §12.6).
// OldData/NewData adalah JSON snapshot movement sebelum/sesudah aksi.
type MovementAuditResponse struct {
	ID         string    `json:"id"`
	MovementID string    `json:"movement_id"`
	Action     string    `json:"action"`
	OldStatus  *string   `json:"old_status,omitempty"`
	NewStatus  *string   `json:"new_status,omitempty"`
	OldData    *string   `json:"old_data,omitempty"`
	NewData    *string   `json:"new_data,omitempty"`
	Reason     *string   `json:"reason,omitempty"`
	ActedBy    *string   `json:"acted_by,omitempty"`
	ActedAt    time.Time `json:"acted_at"`
}

type PaginatedMovementAuditResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

type PaginatedMovementResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// =========================================================================
// Employee Movement Document DTOs (plan §12.15)
// =========================================================================

// CreateMovementDocumentRequest menyimpan metadata dokumen movement. File
// fisik di-upload terpisah via endpoint upload generik (POST /uploads); di
// sini hanya file_url hasil upload + metadata yang dikirim.
type CreateMovementDocumentRequest struct {
	DocumentType string `json:"document_type" binding:"required,oneof=PROMOTION_SK MUTATION_SK DEMOTION_SK RETIREMENT_LETTER OFFBOARDING_LETTER OTHER"`
	FileName     string `json:"file_name" binding:"required"`
	FileURL      string `json:"file_url" binding:"required,startswith=/"`
}

type MovementDocumentResponse struct {
	ID           string    `json:"id"`
	MovementID   string    `json:"movement_id"`
	DocumentType string    `json:"document_type"`
	FileName     string    `json:"file_name"`
	FileURL      string    `json:"file_url"`
	UploadedBy   *string   `json:"uploaded_by,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type PaginatedMovementDocumentResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// =========================================================================
// Employee Contract DTOs
// =========================================================================

type CreateContractRequest struct {
	EmployeeID           string  `json:"employee_id" binding:"required,uuid"`
	ContractNumber       string  `json:"contract_number" binding:"required"`
	ContractType         string  `json:"contract_type" binding:"required,oneof=pkwt pkwtt daily other"`
	StartDate            string  `json:"start_date" binding:"required"`
	EndDate              *string `json:"end_date" binding:"omitempty"`
	PreviousContractID   *string `json:"previous_contract_id" binding:"omitempty,uuid"`
	DecisionLetterNumber *string `json:"decision_letter_number"`
	Notes                *string `json:"notes"`
	DocumentURL          *string `json:"document_url"`
}

type UpdateContractRequest struct {
	ContractNumber       *string `json:"contract_number"`
	ContractType         *string `json:"contract_type" binding:"omitempty,oneof=pkwt pkwtt daily other"`
	EndDate              *string `json:"end_date"`
	DecisionLetterNumber *string `json:"decision_letter_number"`
	Notes                *string `json:"notes"`
	DocumentURL          *string `json:"document_url"`
	Status               *string `json:"status" binding:"omitempty,oneof=active expired extended terminated"`
}

type ContractResponse struct {
	ID                     string    `json:"id"`
	EmployeeID             string    `json:"employee_id"`
	EmployeeName           string    `json:"employee_name"`
	EmployeeCode           string    `json:"employee_code"`
	ContractNumber         string    `json:"contract_number"`
	ContractType           string    `json:"contract_type"`
	StartDate              string    `json:"start_date"`
	EndDate                *string   `json:"end_date,omitempty"`
	ExtensionCount         int       `json:"extension_count"`
	PreviousContractID     *string   `json:"previous_contract_id,omitempty"`
	PreviousContractNumber string    `json:"previous_contract_number,omitempty"`
	DecisionLetterNumber   *string   `json:"decision_letter_number,omitempty"`
	Notes                  *string   `json:"notes,omitempty"`
	DocumentURL            *string   `json:"document_url,omitempty"`
	Status                 string    `json:"status"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type PaginatedContractResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// =========================================================================
// Converter helpers
// =========================================================================

func (m *EmployeeMovement) ToResponse() MovementResponse {
	r := MovementResponse{
		ID:                   m.ID.String(),
		EmployeeID:           m.EmployeeID.String(),
		MovementType:         string(m.MovementType),
		DecisionLetterNumber: m.DecisionLetterNumber,
		DecisionLetterDate:   m.DecisionLetterDate,
		EffectiveDate:        m.EffectiveDate,
		Status:               string(m.Status),
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
	if m.FromEmploymentID != nil {
		s := m.FromEmploymentID.String()
		r.FromEmploymentID = &s
	}
	if m.ToEmploymentID != nil {
		s := m.ToEmploymentID.String()
		r.ToEmploymentID = &s
	}
	if m.FromOrganizationID != nil {
		s := m.FromOrganizationID.String()
		r.FromOrganizationID = &s
	}
	if m.ToOrganizationID != nil {
		s := m.ToOrganizationID.String()
		r.ToOrganizationID = &s
	}
	if m.FromPositionID != nil {
		s := m.FromPositionID.String()
		r.FromPositionID = &s
	}
	if m.ToPositionID != nil {
		s := m.ToPositionID.String()
		r.ToPositionID = &s
	}
	if m.FromEmploymentStatusID != nil {
		s := m.FromEmploymentStatusID.String()
		r.FromEmploymentStatusID = &s
	}
	if m.ToEmploymentStatusID != nil {
		s := m.ToEmploymentStatusID.String()
		r.ToEmploymentStatusID = &s
	}
	// Snapshot names (plan §12.5) — persisted at creation so history is
	// immune to later master-data renames.
	r.FromOrganizationName = m.FromOrganizationName
	r.FromPositionName = m.FromPositionName
	r.FromEmploymentStatusName = m.FromEmploymentStatusName
	r.ToOrganizationName = m.ToOrganizationName
	r.ToPositionName = m.ToPositionName
	r.ToEmploymentStatusName = m.ToEmploymentStatusName
	if m.Reason != nil {
		r.Reason = m.Reason
	}
	if m.Notes != nil {
		r.Notes = m.Notes
	}
	if m.ApprovedBy != nil {
		s := m.ApprovedBy.String()
		r.ApprovedBy = &s
	}
	if m.ApprovedAt != nil {
		r.ApprovedAt = m.ApprovedAt
	}
	if m.ExecutedBy != nil {
		s := m.ExecutedBy.String()
		r.ExecutedBy = &s
	}
	if m.ExecutedAt != nil {
		r.ExecutedAt = m.ExecutedAt
	}
	if m.ApprovalInstanceID != nil {
		s := m.ApprovalInstanceID.String()
		r.ApprovalInstanceID = &s
	}
	return r
}

func (a *EmployeeMovementAudit) ToResponse() MovementAuditResponse {
	r := MovementAuditResponse{
		ID:         a.ID.String(),
		MovementID: a.MovementID.String(),
		Action:     string(a.Action),
		ActedAt:    a.ActedAt,
	}
	if a.OldStatus != nil {
		s := *a.OldStatus
		r.OldStatus = &s
	}
	if a.NewStatus != nil {
		s := *a.NewStatus
		r.NewStatus = &s
	}
	if a.OldData != nil {
		s := *a.OldData
		r.OldData = &s
	}
	if a.NewData != nil {
		s := *a.NewData
		r.NewData = &s
	}
	if a.Reason != nil {
		s := *a.Reason
		r.Reason = &s
	}
	if a.ActedBy != nil {
		s := a.ActedBy.String()
		r.ActedBy = &s
	}
	return r
}

func (d *EmployeeMovementDocument) ToResponse() MovementDocumentResponse {
	r := MovementDocumentResponse{
		ID:           d.ID.String(),
		MovementID:   d.MovementID.String(),
		DocumentType: string(d.DocumentType),
		FileName:     d.FileName,
		FileURL:      d.FileURL,
		CreatedAt:    d.CreatedAt,
	}
	if d.UploadedBy != nil {
		s := d.UploadedBy.String()
		r.UploadedBy = &s
	}
	return r
}

func (c *EmployeeContract) ToResponse() ContractResponse {
	r := ContractResponse{
		ID:             c.ID.String(),
		EmployeeID:     c.EmployeeID.String(),
		ContractNumber: c.ContractNumber,
		ContractType:   string(c.ContractType),
		StartDate:      c.StartDate,
		ExtensionCount: c.ExtensionCount,
		Status:         string(c.Status),
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
	if c.EndDate != nil {
		r.EndDate = c.EndDate
	}
	if c.PreviousContractID != nil {
		s := c.PreviousContractID.String()
		r.PreviousContractID = &s
	}
	if c.DecisionLetterNumber != nil {
		r.DecisionLetterNumber = c.DecisionLetterNumber
	}
	if c.Notes != nil {
		r.Notes = c.Notes
	}
	if c.DocumentURL != nil {
		r.DocumentURL = c.DocumentURL
	}
	return r
}
