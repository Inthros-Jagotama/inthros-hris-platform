package okr

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =========================================================================
// Reuse types from Performance module
// =========================================================================

type FormulaType string

const (
	FormulaTypeManual       FormulaType = "MANUAL"
	FormulaTypeHigherBetter FormulaType = "HIGHER_BETTER"
	FormulaTypeLowerBetter  FormulaType = "LOWER_BETTER"
	FormulaTypeRange        FormulaType = "RANGE"
	FormulaTypeBoolean      FormulaType = "BOOLEAN"
	FormulaTypePercentage   FormulaType = "PERCENTAGE"
)

type TargetType string

const (
	TargetTypeNumber     TargetType = "NUMBER"
	TargetTypeCurrency   TargetType = "CURRENCY"
	TargetTypePercentage TargetType = "PERCENTAGE"
	TargetTypeDuration   TargetType = "DURATION"
	TargetTypeBoolean    TargetType = "BOOLEAN"
)

type EvaluationStatus string

const (
	StatusDraft     EvaluationStatus = "DRAFT"
	StatusSubmitted EvaluationStatus = "SUBMITTED"
	StatusApproved  EvaluationStatus = "APPROVED"
	StatusCompleted EvaluationStatus = "COMPLETED"
	StatusRejected  EvaluationStatus = "REJECTED"
)

// =========================================================================
// OKRTemplate - Template OKR untuk setiap Organization
// =========================================================================

type OKRTemplate struct {
	ID             uuid.UUID       `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID uuid.UUID       `gorm:"type:char(36);not null;index:idx_okr_tpl_org" json:"organization_id"`
	PeriodID       *uuid.UUID      `gorm:"type:char(36);index:idx_okr_tpl_period" json:"period_id,omitempty"`
	Name           string          `gorm:"type:varchar(255);not null" json:"name"`
	Description    *string         `gorm:"type:text" json:"description,omitempty"`
	Status         int             `gorm:"type:smallint;not null;default:0;index:idx_okr_tpl_status" json:"status"`
	EffectiveDate  *string         `gorm:"type:date" json:"effective_date,omitempty"`
	ExpiredDate    *string         `gorm:"type:date" json:"expired_date,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Objectives []OKRObjective `gorm:"foreignKey:TemplateID" json:"objectives,omitempty"`
}

func (OKRTemplate) TableName() string {
	return "okr_templates"
}

func (t *OKRTemplate) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// OKRObjective - Objective yang ingin dicapai
// =========================================================================

type OKRObjective struct {
	ID          uuid.UUID       `gorm:"type:char(36);primaryKey" json:"id"`
	TemplateID  uuid.UUID       `gorm:"type:char(36);not null;index:idx_okr_obj_template" json:"template_id"`
	Code        *string         `gorm:"type:varchar(50)" json:"code,omitempty"`
	Title       string          `gorm:"type:varchar(255);not null" json:"title"`
	Description *string         `gorm:"type:text" json:"description,omitempty"`
	Weight      float64         `gorm:"type:decimal(5,2);not null;default:0" json:"weight"`
	SortOrder   int             `gorm:"type:int;not null;default:0;index:idx_okr_obj_sort" json:"sort_order"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	KeyResults []OKRKeyResult `gorm:"foreignKey:ObjectiveID" json:"key_results,omitempty"`
}

func (OKRObjective) TableName() string {
	return "okr_objectives"
}

func (o *OKRObjective) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// OKRKeyResult - Target terukur dari setiap Objective
// =========================================================================

type OKRKeyResult struct {
	ID           uuid.UUID       `gorm:"type:char(36);primaryKey" json:"id"`
	ObjectiveID  uuid.UUID       `gorm:"type:char(36);not null;index:idx_okr_kr_objective" json:"objective_id"`
	Code         *string         `gorm:"type:varchar(50)" json:"code,omitempty"`
	Title        string          `gorm:"type:varchar(255);not null" json:"title"`
	Description  *string         `gorm:"type:text" json:"description,omitempty"`
	TargetType   TargetType      `gorm:"type:varchar(30);not null;default:'NUMBER'" json:"target_type"`
	TargetValue  float64         `gorm:"type:decimal(18,2);not null;default:0" json:"target_value"`
	Unit         *string         `gorm:"type:varchar(50)" json:"unit,omitempty"`
	FormulaType  FormulaType     `gorm:"type:varchar(30);not null;default:'HIGHER_BETTER'" json:"formula_type"`
	Weight       float64         `gorm:"type:decimal(5,2);not null;default:0" json:"weight"`
	MinimumScore float64         `gorm:"type:decimal(5,2);not null;default:0" json:"minimum_score"`
	MaximumScore float64         `gorm:"type:decimal(5,2);not null;default:100" json:"maximum_score"`
	SortOrder    int             `gorm:"type:int;not null;default:0;index:idx_okr_kr_sort" json:"sort_order"`
	IsRequired   bool            `gorm:"type:boolean;not null;default:true" json:"is_required"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (OKRKeyResult) TableName() string {
	return "okr_key_results"
}

func (k *OKRKeyResult) BeforeCreate(tx *gorm.DB) error {
	if k.ID == uuid.Nil {
		k.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// OKREvaluation - Header penilaian OKR
// =========================================================================

type OKREvaluation struct {
	ID             uuid.UUID        `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID     uuid.UUID        `gorm:"type:char(36);not null;index:idx_okr_eval_employee" json:"employee_id"`
	OrganizationID uuid.UUID        `gorm:"type:char(36);not null;index:idx_okr_eval_org" json:"organization_id"`
	PeriodID       uuid.UUID        `gorm:"type:char(36);not null;index:idx_okr_eval_period" json:"period_id"`
	TemplateID     *uuid.UUID       `gorm:"type:char(36)" json:"template_id,omitempty"`
	Status         EvaluationStatus `gorm:"type:varchar(20);not null;default:'DRAFT';index:idx_okr_eval_status" json:"status"`
	SubmittedAt    *time.Time       `json:"submitted_at,omitempty"`
	SubmittedBy    *uuid.UUID       `gorm:"type:char(36)" json:"submitted_by,omitempty"`
	ApprovedAt     *time.Time       `json:"approved_at,omitempty"`
	ApprovedBy     *uuid.UUID       `gorm:"type:char(36)" json:"approved_by,omitempty"`
	FinalScore     float64          `gorm:"type:decimal(5,2);not null;default:0" json:"final_score"`
	RatingID       *uuid.UUID       `gorm:"type:char(36)" json:"rating_id,omitempty"`
	ReviewerNotes  *string          `gorm:"type:text" json:"reviewer_notes,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	DeletedAt      *gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Details []OKREvaluationDetail `gorm:"foreignKey:EvaluationID" json:"details,omitempty"`
}

func (OKREvaluation) TableName() string {
	return "okr_evaluations"
}

func (e *OKREvaluation) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// OKREvaluationDetail - Snapshot Objective dan Key Result saat evaluasi
// =========================================================================

type OKREvaluationDetail struct {
	ID              uuid.UUID   `gorm:"type:char(36);primaryKey" json:"id"`
	EvaluationID    uuid.UUID   `gorm:"type:char(36);not null;index:idx_okr_detail_eval" json:"evaluation_id"`
	ObjectiveID     *uuid.UUID  `gorm:"type:char(36);index:idx_okr_detail_obj" json:"objective_id,omitempty"`
	KeyResultID     *uuid.UUID  `gorm:"type:char(36);index:idx_okr_detail_kr" json:"key_result_id,omitempty"`
	ObjectiveTitle  string      `gorm:"type:varchar(255);not null" json:"objective_title"`
	KeyResultTitle  string      `gorm:"type:varchar(255);not null" json:"key_result_title"`
	ObjectiveWeight float64     `gorm:"type:decimal(5,2);not null;default:0" json:"objective_weight"`
	KeyResultWeight float64     `gorm:"type:decimal(5,2);not null;default:0" json:"key_result_weight"`
	TargetValue     float64     `gorm:"type:decimal(18,2);not null;default:0" json:"target_value"`
	TargetType      TargetType  `gorm:"type:varchar(30);not null;default:'NUMBER'" json:"target_type"`
	Unit            *string     `gorm:"type:varchar(50)" json:"unit,omitempty"`
	FormulaType     FormulaType `gorm:"type:varchar(30);not null;default:'HIGHER_BETTER'" json:"formula_type"`
	ActualValue     float64     `gorm:"type:decimal(18,2);not null;default:0" json:"actual_value"`
	Achievement     float64     `gorm:"type:decimal(5,2);not null;default:0" json:"achievement"`
	Score           float64     `gorm:"type:decimal(5,2);not null;default:0" json:"score"`
	Remarks         *string     `gorm:"type:text" json:"remarks,omitempty"`
	SortOrder       int         `gorm:"type:int;not null;default:0" json:"sort_order"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

func (OKREvaluationDetail) TableName() string {
	return "okr_evaluation_details"
}

func (d *OKREvaluationDetail) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// OKRProgress - Progress Check-in
// =========================================================================

type OKRProgress struct {
	ID                 uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	EvaluationDetailID uuid.UUID `gorm:"type:char(36);not null;index:idx_okr_prog_detail" json:"evaluation_detail_id"`
	ProgressDate       string    `gorm:"type:date;not null;index:idx_okr_prog_date" json:"progress_date"`
	ActualValue        float64   `gorm:"type:decimal(18,2);not null;default:0" json:"actual_value"`
	Achievement        float64   `gorm:"type:decimal(5,2);not null;default:0" json:"achievement"`
	Notes              *string   `gorm:"type:text" json:"notes,omitempty"`
	CreatedBy          uuid.UUID `gorm:"type:char(36);not null" json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (OKRProgress) TableName() string {
	return "okr_progress"
}

func (p *OKRProgress) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// OKRComment - Komentar antara Employee dan Reviewer
// =========================================================================

type OKRComment struct {
	ID           uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EvaluationID uuid.UUID  `gorm:"type:char(36);not null;index:idx_okr_comment_eval" json:"evaluation_id"`
	ParentID     *uuid.UUID `gorm:"type:char(36);index:idx_okr_comment_parent" json:"parent_id,omitempty"`
	Comment      string     `gorm:"type:text;not null" json:"comment"`
	CreatedBy    uuid.UUID  `gorm:"type:char(36);not null;index:idx_okr_comment_user" json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Relations
	Replies []OKRComment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

func (OKRComment) TableName() string {
	return "okr_comments"
}

func (c *OKRComment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// OKRAttachment - Evidence pencapaian Key Result
// =========================================================================

type OKRAttachment struct {
	ID                 uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	EvaluationDetailID uuid.UUID `gorm:"type:char(36);not null;index:idx_okr_attach_detail" json:"evaluation_detail_id"`
	FilePath           string    `gorm:"type:varchar(500);not null" json:"file_path"`
	FileName           string    `gorm:"type:varchar(255);not null" json:"file_name"`
	FileType           *string   `gorm:"type:varchar(100)" json:"file_type,omitempty"`
	FileSize           *int64    `gorm:"type:bigint" json:"file_size,omitempty"`
	Description        *string   `gorm:"type:text" json:"description,omitempty"`
	UploadedBy         uuid.UUID `gorm:"type:char(36);not null;index:idx_okr_attach_user" json:"uploaded_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (OKRAttachment) TableName() string {
	return "okr_attachments"
}

func (a *OKRAttachment) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
