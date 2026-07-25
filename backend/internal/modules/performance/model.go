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
// PerformancePeriod (Periode evaluasi: tahunan, kuartal, bulanan)
// =========================================================================

type PerformancePeriod struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	PeriodCode  string    `gorm:"type:varchar(10);not null;index:idx_perf_period_code" json:"period_code"`
	PeriodType  string    `gorm:"type:varchar(20);not null" json:"period_type"`
	Year        int       `gorm:"type:smallint;not null;index:idx_perf_period_year" json:"year"`
	StartDate   *string   `gorm:"type:date" json:"start_date,omitempty"`
	EndDate     *string   `gorm:"type:date" json:"end_date,omitempty"`
	Status      string    `gorm:"type:varchar(20);default:active" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
	ID             uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID uuid.UUID `gorm:"type:char(36);not null;index:idx_perf_tmpl_org" json:"organization_id"`
	Name           string    `gorm:"type:varchar(200);not null" json:"name"`
	Description    *string   `gorm:"type:text" json:"description,omitempty"`
	Status         string    `gorm:"type:varchar(20);default:DRAFT" json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
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
	ID                  uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	PerformanceTemplateID uuid.UUID `gorm:"type:char(36);not null;index:idx_perf_ind_tmpl" json:"performance_template_id"`
	PerspectiveID       uuid.UUID `gorm:"type:char(36);not null;index:idx_perf_ind_persp" json:"perspective_id"`
	IndicatorType       string    `gorm:"type:varchar(20);not null" json:"indicator_type"`
	Title               string    `gorm:"type:varchar(255);not null" json:"title"`
	Description         *string   `gorm:"type:text" json:"description,omitempty"`
	Weight              float64   `gorm:"type:decimal(5,2);not null;default:0" json:"weight"`
	TargetValue         float64   `gorm:"type:decimal(12,2);default:0" json:"target_value"`
	UnitOfMeasurement   *string   `gorm:"type:varchar(50)" json:"unit_of_measurement,omitempty"`
	SortOrder           int       `gorm:"type:smallint;default:0" json:"sort_order"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
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
	ID              uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID      uuid.UUID  `gorm:"type:char(36);not null;index:idx_perf_eval_emp" json:"employee_id"`
	OrganizationID  uuid.UUID  `gorm:"type:char(36);not null;index:idx_perf_eval_org" json:"organization_id"`
	PeriodID        uuid.UUID  `gorm:"type:char(36);not null;index:idx_perf_eval_period" json:"period_id"`
	TemplateID      uuid.UUID  `gorm:"type:char(36);not null;index:idx_perf_eval_tmpl" json:"template_id"`
	SupervisorID    *uuid.UUID `gorm:"type:char(36);index:idx_perf_eval_sup" json:"supervisor_id,omitempty"`
	FinalScore      float64    `gorm:"type:decimal(5,2);default:0" json:"final_score"`
	Status          string     `gorm:"type:varchar(20);default:DRAFT" json:"status"`
	PlanSubmittedAt *int64     `gorm:"type:bigint;default:0" json:"-"`
	PlanApprovedAt  *int64     `gorm:"type:bigint;default:0" json:"-"`
	ActualSubmittedAt *int64   `gorm:"type:bigint;default:0" json:"-"`
	ActualApprovedAt  *int64   `gorm:"type:bigint;default:0" json:"-"`
	Notes           *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

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
	ID                      uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	PerformanceEvaluationID uuid.UUID `gorm:"type:char(36);not null;index:idx_perf_detail_eval" json:"performance_evaluation_id"`
	PerspectiveID           uuid.UUID `gorm:"type:char(36);not null" json:"perspective_id"`
	AchievementPercentage   float64   `gorm:"type:decimal(5,2);default:0" json:"achievement_percentage"`
	Weight                  float64   `gorm:"type:decimal(5,2);not null;default:0" json:"weight"`
	Score                   float64   `gorm:"type:decimal(5,2);default:0" json:"score"`
	Description             *string   `gorm:"type:varchar(255)" json:"description,omitempty"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
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
// PerformanceTarget (Target KPI individual — nilai target vs aktual)
// =========================================================================

type PerformanceTarget struct {
	ID                  uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	PerformanceEvaluationID uuid.UUID `gorm:"type:char(36);not null;index:idx_perf_tgt_eval" json:"performance_evaluation_id"`
	IndicatorID        uuid.UUID `gorm:"type:char(36);not null;index:idx_perf_tgt_ind" json:"indicator_id"`
	TargetValue        float64   `gorm:"type:decimal(12,2);default:0" json:"target_value"`
	ActualValue        *float64  `gorm:"type:decimal(12,2)" json:"actual_value,omitempty"`
	UnitOfMeasurement  *string   `gorm:"type:varchar(50)" json:"unit_of_measurement,omitempty"`
	AchievementPercent float64   `gorm:"type:decimal(5,2);default:0" json:"achievement_percentage"`
	Weight             float64   `gorm:"type:decimal(5,2);not null;default:0" json:"weight"`
	Score              float64   `gorm:"type:decimal(5,2);default:0" json:"score"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
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
