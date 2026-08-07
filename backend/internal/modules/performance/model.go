package performance

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =========================================================================
// PerformanceIndicatorType (MAXIMIZATION / MINIMIZATION)
// =========================================================================

type IndicatorType string

const (
	IndicatorMaximization IndicatorType = "MAXIMIZATION"
	IndicatorMinimization IndicatorType = "MINIMIZATION"
)

// =========================================================================
// FormulaType (Jenis formula kalkulasi skor KPI)
// =========================================================================

type FormulaType string

const (
	FormulaTypeManual       FormulaType = "MANUAL"
	FormulaTypeHigherBetter FormulaType = "HIGHER_BETTER"
	FormulaTypeLowerBetter  FormulaType = "LOWER_BETTER"
	FormulaTypeRange        FormulaType = "RANGE"
)

// =========================================================================
// TargetType (Jenis target KPI)
// =========================================================================

type TargetType string

const (
	TargetTypeNumber     TargetType = "NUMBER"
	TargetTypeCurrency   TargetType = "CURRENCY"
	TargetTypePercentage TargetType = "PERCENTAGE"
	TargetTypeDuration   TargetType = "DURATION"
	TargetTypeBoolean    TargetType = "BOOLEAN"
)

// =========================================================================
// PerformancePeriod (Periode evaluasi: tahunan, kuartal, bulanan)
// =========================================================================

type PerformancePeriod struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	PeriodCode string    `gorm:"type:varchar(10);not null;index:idx_perf_period_code" json:"period_code"`
	PeriodType string    `gorm:"type:varchar(20);not null" json:"period_type"`
	Year       int       `gorm:"type:smallint;not null;index:idx_perf_period_year" json:"year"`
	StartDate  *string   `gorm:"type:date" json:"start_date,omitempty"`
	EndDate    *string   `gorm:"type:date" json:"end_date,omitempty"`
	Status     string    `gorm:"type:varchar(20);default:active" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (PerformancePeriod) TableName() string {
	return "performance_periods"
}

func (p *PerformancePeriod) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// PerformancePerspective (Perspektif BSC: Financial, Customer, Internal, L&G)
// =========================================================================

type PerformancePerspective struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	SortOrder   int       `gorm:"type:smallint;default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (PerformancePerspective) TableName() string {
	return "performance_perspectives"
}

func (p *PerformancePerspective) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// PerformanceTemplate (Template BSC per organisasi/posisi)
// =========================================================================

type PerformanceTemplate struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID uuid.UUID  `gorm:"type:char(36);not null;index:idx_perf_tmpl_org" json:"organization_id"`
	PeriodID       *uuid.UUID `gorm:"type:char(36);index:idx_perf_tmpl_period" json:"period_id,omitempty"`
	Name           string     `gorm:"type:varchar(200);not null" json:"name"`
	Description    *string    `gorm:"type:text" json:"description,omitempty"`
	Status         string     `gorm:"type:varchar(20);default:DRAFT" json:"status"`
	EffectiveDate  *string    `gorm:"type:date" json:"effective_date,omitempty"`
	ExpiredDate    *string    `gorm:"type:date" json:"expired_date,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (PerformanceTemplate) TableName() string {
	return "performance_templates"
}

func (t *PerformanceTemplate) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// PerformanceIndicator (Indikator KPI — linked to template)
// =========================================================================

type PerformanceIndicator struct {
	ID                    uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	PerformanceTemplateID uuid.UUID `gorm:"type:char(36);not null;index:idx_perf_ind_tmpl" json:"performance_template_id"`
	PerspectiveID         uuid.UUID `gorm:"type:char(36);not null;index:idx_perf_ind_persp" json:"perspective_id"`
	Code                  *string   `gorm:"type:varchar(50);index:idx_perf_ind_code" json:"code,omitempty"`
	IndicatorType         string    `gorm:"type:varchar(20);not null" json:"indicator_type"`
	Title                 string    `gorm:"type:varchar(255);not null" json:"title"`
	Description           *string   `gorm:"type:text" json:"description,omitempty"`
	Weight                float64   `gorm:"type:decimal(5,2);not null;default:0" json:"weight"`
	TargetValue           float64   `gorm:"type:decimal(12,2);default:0" json:"target_value"`
	UnitOfMeasurement     *string   `gorm:"type:varchar(50)" json:"unit_of_measurement,omitempty"`
	FormulaType           string    `gorm:"type:varchar(30);not null;default:MANUAL" json:"formula_type"`
	MinimumScore          float64   `gorm:"type:decimal(5,2);not null;default:0" json:"minimum_score"`
	MaximumScore          float64   `gorm:"type:decimal(5,2);not null;default:100" json:"maximum_score"`
	TargetType            string    `gorm:"type:varchar(30);not null;default:NUMBER" json:"target_type"`
	IsRequired            bool      `gorm:"type:smallint;not null;default:1" json:"is_required"`
	SortOrder             int       `gorm:"type:smallint;default:0" json:"sort_order"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func (PerformanceIndicator) TableName() string {
	return "performance_indicators"
}

func (i *PerformanceIndicator) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// PerformanceEvaluation (Evaluasi kinerja karyawan)
// =========================================================================

type PerformanceEvaluation struct {
	ID                uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID        uuid.UUID  `gorm:"type:char(36);not null;index:idx_perf_eval_emp" json:"employee_id"`
	OrganizationID    uuid.UUID  `gorm:"type:char(36);not null;index:idx_perf_eval_org" json:"organization_id"`
	PeriodID          uuid.UUID  `gorm:"type:char(36);not null;index:idx_perf_eval_period" json:"period_id"`
	TemplateID        uuid.UUID  `gorm:"type:char(36);not null;index:idx_perf_eval_tmpl" json:"template_id"`
	SupervisorID      *uuid.UUID `gorm:"type:char(36);index:idx_perf_eval_sup" json:"supervisor_id,omitempty"`
	FinalScore        float64    `gorm:"type:decimal(5,2);default:0" json:"final_score"`
	RatingID          *uuid.UUID `gorm:"type:char(36);index:idx_perf_eval_rating" json:"rating_id,omitempty"`
	Status            string     `gorm:"type:varchar(20);default:DRAFT" json:"status"`
	PlanSubmittedAt   *int64     `gorm:"type:bigint;default:0" json:"-"`
	PlanApprovedAt    *int64     `gorm:"type:bigint;default:0" json:"-"`
	ActualSubmittedAt *int64     `gorm:"type:bigint;default:0" json:"-"`
	ActualApprovedAt  *int64     `gorm:"type:bigint;default:0" json:"-"`
	SubmittedAt       *time.Time `gorm:"type:timestamp" json:"submitted_at,omitempty"`
	ApprovedAt        *time.Time `gorm:"type:timestamp" json:"approved_at,omitempty"`
	TargetSubmittedAt *time.Time `gorm:"type:timestamp" json:"target_submitted_at,omitempty"`
	TargetApprovedAt  *time.Time `gorm:"type:timestamp" json:"target_approved_at,omitempty"`
	// TargetApprovalInstanceID/RealizationApprovalInstanceID reference the
	// central approval module's ApprovalInstance for each phase, when a
	// flow is configured (module_kpi_target / performance_kpi_realization).
	TargetApprovalInstanceID      *uuid.UUID `gorm:"type:char(36);index:idx_perf_eval_target_approval" json:"target_approval_instance_id,omitempty"`
	RealizationApprovalInstanceID *uuid.UUID `gorm:"type:char(36);index:idx_perf_eval_realization_approval" json:"realization_approval_instance_id,omitempty"`
	Notes                         *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Relasi
	Details []PerformanceEvaluationDetail `gorm:"foreignKey:PerformanceEvaluationID" json:"details,omitempty"`
}

func (PerformanceEvaluation) TableName() string {
	return "performance_evaluations"
}

func (e *PerformanceEvaluation) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// PerformanceEvaluationDetail (Detail grup BSC dalam evaluasi)
// =========================================================================

type PerformanceEvaluationDetail struct {
	ID                      uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	PerformanceEvaluationID uuid.UUID  `gorm:"type:char(36);not null;index:idx_perf_detail_eval" json:"performance_evaluation_id"`
	PerspectiveID           uuid.UUID  `gorm:"type:char(36);not null" json:"perspective_id"`
	IndicatorID             *uuid.UUID `gorm:"type:char(36);index:idx_perf_detail_indicator" json:"indicator_id,omitempty"`
	IndicatorName           *string    `gorm:"type:varchar(255)" json:"indicator_name,omitempty"`
	AchievementPercentage   float64    `gorm:"type:decimal(5,2);default:0" json:"achievement_percentage"`
	Weight                  float64    `gorm:"type:decimal(5,2);not null;default:0" json:"weight"`
	UnitOfMeasurement       *string    `gorm:"type:varchar(50)" json:"unit_of_measurement,omitempty"`
	Target                  float64    `gorm:"type:decimal(18,2);not null;default:0" json:"target"`
	Actual                  float64    `gorm:"type:decimal(18,2);not null;default:0" json:"actual"`
	Achievement             float64    `gorm:"type:decimal(5,2);not null;default:0" json:"achievement"`
	Score                   float64    `gorm:"type:decimal(5,2);default:0" json:"score"`
	Description             *string    `gorm:"type:varchar(255)" json:"description,omitempty"`
	Remarks                 *string    `gorm:"type:text" json:"remarks,omitempty"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

func (PerformanceEvaluationDetail) TableName() string {
	return "performance_evaluation_details"
}

func (d *PerformanceEvaluationDetail) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// PerformanceEvaluationProgramItem — Program items the employee proposes
// themselves per evaluation (no HR-authored template, unlike indicators).
// Only relevant when the org's PROGRAM component is enabled.
// =========================================================================

type PerformanceEvaluationProgramItem struct {
	ID                      uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	PerformanceEvaluationID uuid.UUID `gorm:"type:char(36);not null;index:idx_perf_prog_eval" json:"performance_evaluation_id"`
	Title                   string    `gorm:"type:varchar(255);not null" json:"title"`
	FormulaType             string    `gorm:"type:varchar(30);not null;default:MANUAL" json:"formula_type"`
	Weight                  float64   `gorm:"type:decimal(5,2);not null;default:0" json:"weight"`
	UnitOfMeasurement       *string   `gorm:"type:varchar(50)" json:"unit_of_measurement,omitempty"`
	Target                  float64   `gorm:"type:decimal(18,2);not null;default:0" json:"target"`
	Actual                  float64   `gorm:"type:decimal(18,2);not null;default:0" json:"actual"`
	Achievement             float64   `gorm:"type:decimal(5,2);not null;default:0" json:"achievement"`
	Score                   float64   `gorm:"type:decimal(5,2);default:0" json:"score"`
	SortOrder               int       `gorm:"type:smallint;default:0" json:"sort_order"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

func (PerformanceEvaluationProgramItem) TableName() string {
	return "performance_evaluation_program_items"
}

func (p *PerformanceEvaluationProgramItem) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// PerformanceTarget (Target KPI individual — nilai target vs aktual)
// =========================================================================

type PerformanceTarget struct {
	ID                      uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	PerformanceEvaluationID uuid.UUID `gorm:"type:char(36);not null;index:idx_perf_tgt_eval" json:"performance_evaluation_id"`
	IndicatorID             uuid.UUID `gorm:"type:char(36);not null;index:idx_perf_tgt_ind" json:"indicator_id"`
	TargetValue             float64   `gorm:"type:decimal(12,2);default:0" json:"target_value"`
	ActualValue             *float64  `gorm:"type:decimal(12,2)" json:"actual_value,omitempty"`
	UnitOfMeasurement       *string   `gorm:"type:varchar(50)" json:"unit_of_measurement,omitempty"`
	AchievementPercent      float64   `gorm:"type:decimal(5,2);default:0" json:"achievement_percentage"`
	Weight                  float64   `gorm:"type:decimal(5,2);not null;default:0" json:"weight"`
	Score                   float64   `gorm:"type:decimal(5,2);default:0" json:"score"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

func (PerformanceTarget) TableName() string {
	return "performance_targets"
}

func (t *PerformanceTarget) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// PerformanceProgress (Progress monitoring KPI)
// =========================================================================

type PerformanceProgress struct {
	ID                 uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	EvaluationDetailID uuid.UUID `gorm:"type:char(36);not null;index:idx_perf_prog_detail" json:"evaluation_detail_id"`
	ProgressDate       string    `gorm:"type:date;not null" json:"progress_date"`
	ActualValue        float64   `gorm:"type:decimal(18,2);not null;default:0" json:"actual_value"`
	Achievement        float64   `gorm:"type:decimal(5,2);not null;default:0" json:"achievement"`
	Notes              *string   `gorm:"type:text" json:"notes,omitempty"`
	CreatedBy          uuid.UUID `gorm:"type:char(36);not null" json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (PerformanceProgress) TableName() string {
	return "performance_progress"
}

func (p *PerformanceProgress) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// PerformanceComment (Komentar evaluasi)
// =========================================================================

type PerformanceComment struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	EvaluationID uuid.UUID `gorm:"type:char(36);not null;index:idx_perf_comment_eval" json:"evaluation_id"`
	EmployeeID   uuid.UUID `gorm:"type:char(36);not null;index:idx_perf_comment_emp" json:"employee_id"`
	Comment      string    `gorm:"type:text;not null" json:"comment"`
	CreatedBy    uuid.UUID `gorm:"type:char(36);not null" json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (PerformanceComment) TableName() string {
	return "performance_comments"
}

func (c *PerformanceComment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// PerformanceAttachment (Evidence KPI)
// =========================================================================

type PerformanceAttachment struct {
	ID                 uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	EvaluationDetailID uuid.UUID `gorm:"type:char(36);not null;index:idx_perf_attach_detail" json:"evaluation_detail_id"`
	FilePath           string    `gorm:"type:varchar(500);not null" json:"file_path"`
	FileName           string    `gorm:"type:varchar(255);not null" json:"file_name"`
	FileType           *string   `gorm:"type:varchar(100)" json:"file_type,omitempty"`
	FileSize           *int64    `gorm:"type:bigint" json:"file_size,omitempty"`
	Description        *string   `gorm:"type:text" json:"description,omitempty"`
	UploadedBy         uuid.UUID `gorm:"type:char(36);not null" json:"uploaded_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (PerformanceAttachment) TableName() string {
	return "performance_attachments"
}

func (a *PerformanceAttachment) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// PerformanceRating (Master rating)
// =========================================================================

type PerformanceRating struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Code        string    `gorm:"type:varchar(20);not null;uniqueIndex:idx_perf_rating_code" json:"code"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	MinScore    float64   `gorm:"type:decimal(5,2);not null" json:"min_score"`
	MaxScore    float64   `gorm:"type:decimal(5,2);not null" json:"max_score"`
	Color       *string   `gorm:"type:varchar(20)" json:"color,omitempty"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	SortOrder   int       `gorm:"type:smallint;not null;default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (PerformanceRating) TableName() string {
	return "performance_ratings"
}

func (r *PerformanceRating) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// PerformanceIndicatorFormula (Master formula)
// =========================================================================

type PerformanceIndicatorFormula struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Code        string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_perf_formula_code" json:"code"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	FormulaType string    `gorm:"type:varchar(30);not null" json:"formula_type"`
	Expression  *string   `gorm:"type:text" json:"expression,omitempty"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	SortOrder   int       `gorm:"type:smallint;not null;default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (PerformanceIndicatorFormula) TableName() string {
	return "performance_indicator_formulas"
}

func (f *PerformanceIndicatorFormula) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// PerformanceLog (Audit trail)
// =========================================================================

type PerformanceLog struct {
	ID           uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EvaluationID *uuid.UUID `gorm:"type:char(36);index:idx_perf_log_eval" json:"evaluation_id,omitempty"`
	EntityType   string     `gorm:"type:varchar(50);not null;index:idx_perf_log_entity" json:"entity_type"`
	EntityID     uuid.UUID  `gorm:"type:char(36);not null" json:"entity_id"`
	Action       string     `gorm:"type:varchar(100);not null;index:idx_perf_log_action" json:"action"`
	OldValues    *string    `gorm:"type:text" json:"old_values,omitempty"`
	NewValues    *string    `gorm:"type:text" json:"new_values,omitempty"`
	CreatedBy    uuid.UUID  `gorm:"type:char(36);not null;index:idx_perf_log_user" json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (PerformanceLog) TableName() string {
	return "performance_logs"
}

func (l *PerformanceLog) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// PerformanceComponent (Master komponen penilaian: KPI, Work Program, dst)
// =========================================================================

type PerformanceComponent struct {
	ID          uuid.UUID       `gorm:"type:char(36);primaryKey" json:"id"`
	Code        string          `gorm:"type:varchar(50);not null;uniqueIndex:idx_perf_comp_code" json:"code"`
	Name        string          `gorm:"type:varchar(100);not null" json:"name"`
	Description *string         `gorm:"type:text" json:"description,omitempty"`
	SortOrder   int             `gorm:"type:int;not null;default:0;index:idx_perf_comp_sort" json:"sort_order"`
	IsActive    bool            `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (PerformanceComponent) TableName() string {
	return "performance_components"
}

func (c *PerformanceComponent) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// PerformanceOrganizationComponent (Konfigurasi bobot komponen per Organization)
// =========================================================================

type PerformanceOrganizationComponent struct {
	ID             uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID uuid.UUID `gorm:"type:char(36);not null;index:idx_perf_orgcomp_org" json:"organization_id"`
	ComponentID    uuid.UUID `gorm:"type:char(36);not null;index:idx_perf_orgcomp_comp" json:"component_id"`
	Weight         float64   `gorm:"type:decimal(5,2);not null;default:0" json:"weight"`
	IsEnabled      bool      `gorm:"type:boolean;not null;default:true" json:"is_enabled"`
	SortOrder      int       `gorm:"type:int;not null;default:0" json:"sort_order"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (PerformanceOrganizationComponent) TableName() string {
	return "performance_organization_components"
}

func (o *PerformanceOrganizationComponent) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// PerformanceEvaluationComponent (Snapshot hasil perhitungan komponen saat evaluasi)
// =========================================================================

type PerformanceEvaluationComponent struct {
	ID            uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EvaluationID  uuid.UUID  `gorm:"type:char(36);not null;index:idx_perf_evalcomp_eval" json:"evaluation_id"`
	ComponentID   uuid.UUID  `gorm:"type:char(36);not null;index:idx_perf_evalcomp_comp" json:"component_id"`
	ComponentName string     `gorm:"type:varchar(100);not null" json:"component_name"`
	Score         float64    `gorm:"type:decimal(5,2);not null;default:0" json:"score"`
	Weight        float64    `gorm:"type:decimal(5,2);not null;default:0" json:"weight"`
	FinalScore    float64    `gorm:"type:decimal(5,2);not null;default:0" json:"final_score"`
	CalculatedAt  *time.Time `json:"calculated_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (PerformanceEvaluationComponent) TableName() string {
	return "performance_evaluation_components"
}

func (e *PerformanceEvaluationComponent) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
