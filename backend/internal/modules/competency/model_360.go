package competency

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =========================================================================
// Competency 360 — Assessment Template (plan generik §5)
// =========================================================================

type CompetencyAssessmentTemplate struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Code        string     `gorm:"type:varchar(50);not null;uniqueIndex:uq_comp_tpl_code" json:"code"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	Status      string     `gorm:"type:varchar(20);default:active" json:"status"`
	ScaleID     *uuid.UUID `gorm:"type:char(36);index" json:"scale_id,omitempty"`
	CreatedBy   *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy   *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relasi
	Competencies []CompetencyAssessmentTemplateCompetency `gorm:"foreignKey:TemplateID" json:"competencies,omitempty"`
	RaterTypes   []CompetencyAssessmentTemplateRaterType   `gorm:"foreignKey:TemplateID" json:"rater_types,omitempty"`
}

func (CompetencyAssessmentTemplate) TableName() string {
	return "competency_assessment_templates"
}

func (t *CompetencyAssessmentTemplate) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// CompetencyAssessmentTemplateCompetency menghubungkan template dengan
// competency yang dinilai + required level & weight (plan generik §5.2).
type CompetencyAssessmentTemplateCompetency struct {
	ID            uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	TemplateID    uuid.UUID `gorm:"type:char(36);not null;index" json:"template_id"`
	CompetencyID  uuid.UUID `gorm:"type:char(36);not null;index" json:"competency_id"`
	RequiredLevel *int      `gorm:"type:smallint" json:"required_level,omitempty"`
	Weight        float64   `gorm:"type:decimal(6,2);default:1" json:"weight"`
	SortOrder     int       `gorm:"type:int;default:0" json:"sort_order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relasi
	Competency *Competency `gorm:"foreignKey:CompetencyID" json:"competency,omitempty"`
}

func (CompetencyAssessmentTemplateCompetency) TableName() string {
	return "competency_assessment_template_competencies"
}

func (t *CompetencyAssessmentTemplateCompetency) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// CompetencyAssessmentTemplateRaterType mengkonfigurasi weight & anonymity
// per rater type di level template (plan generik §10).
type CompetencyAssessmentTemplateRaterType struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	TemplateID uuid.UUID `gorm:"type:char(36);not null;index" json:"template_id"`
	RaterType  string    `gorm:"type:varchar(20);not null" json:"rater_type"`
	Weight     float64   `gorm:"type:decimal(6,2);default:0" json:"weight"`
	MinRater   int       `gorm:"type:int;default:1" json:"min_rater"`
	MaxRater   *int      `gorm:"type:int" json:"max_rater,omitempty"`
	Required   bool      `gorm:"type:boolean;default:false" json:"required"`
	Anonymous  bool      `gorm:"type:boolean;default:false" json:"anonymous"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (CompetencyAssessmentTemplateRaterType) TableName() string {
	return "competency_assessment_template_rater_types"
}

func (t *CompetencyAssessmentTemplateRaterType) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Competency 360 — Rating Scale (plan generik §7)
// =========================================================================

type CompetencyRatingScale struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Code        string     `gorm:"type:varchar(50);not null;uniqueIndex:uq_comp_scale_code" json:"code"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	Status      string     `gorm:"type:varchar(20);default:active" json:"status"`
	CreatedBy   *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy   *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relasi
	Items []CompetencyRatingScaleItem `gorm:"foreignKey:ScaleID" json:"items,omitempty"`
}

func (CompetencyRatingScale) TableName() string {
	return "competency_rating_scales"
}

func (s *CompetencyRatingScale) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type CompetencyRatingScaleItem struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	ScaleID     uuid.UUID `gorm:"type:char(36);not null;index" json:"scale_id"`
	Value       int       `gorm:"type:smallint;not null" json:"value"`
	Label       string    `gorm:"type:varchar(255);not null" json:"label"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	Weight      float64   `gorm:"type:decimal(6,2);default:1" json:"weight"`
	SortOrder   int       `gorm:"type:int;default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (CompetencyRatingScaleItem) TableName() string {
	return "competency_rating_scale_items"
}

func (i *CompetencyRatingScaleItem) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Competency 360 — Indicator (plan generik §6)
// =========================================================================

type CompetencyIndicator struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	CompetencyID uuid.UUID `gorm:"type:char(36);not null;index" json:"competency_id"`
	Code         *string   `gorm:"type:varchar(50)" json:"code,omitempty"`
	Statement    string    `gorm:"type:varchar(1000);not null" json:"statement"`
	Description  *string   `gorm:"type:text" json:"description,omitempty"`
	Status       string    `gorm:"type:varchar(20);default:active" json:"status"`
	SortOrder    int       `gorm:"type:int;default:0" json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relasi
	Competency *Competency `gorm:"foreignKey:CompetencyID" json:"competency,omitempty"`
}

func (CompetencyIndicator) TableName() string {
	return "competency_indicators"
}

func (i *CompetencyIndicator) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

type CompetencyAssessmentTemplateIndicator struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	TemplateID  uuid.UUID `gorm:"type:char(36);not null;index" json:"template_id"`
	IndicatorID uuid.UUID `gorm:"type:char(36);not null;index" json:"indicator_id"`
	Weight      float64   `gorm:"type:decimal(6,2);default:1" json:"weight"`
	SortOrder   int       `gorm:"type:int;default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relasi
	Indicator *CompetencyIndicator `gorm:"foreignKey:IndicatorID" json:"indicator,omitempty"`
}

func (CompetencyAssessmentTemplateIndicator) TableName() string {
	return "competency_assessment_template_indicators"
}

func (i *CompetencyAssessmentTemplateIndicator) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Competency 360 — Rater Assignment (plan generik §9)
// =========================================================================

// RaterType adalah jenis rater pada assessment 360.
type RaterType string

const (
	RaterTypeSelf         RaterType = "self"
	RaterTypeSuperior     RaterType = "superior"
	RaterTypePeer         RaterType = "peer"
	RaterTypeSubordinate  RaterType = "subordinate"
	RaterTypeOther        RaterType = "other"
)

// RaterStatus adalah status siklus assignment rater.
type RaterStatus string

const (
	RaterStatusAssigned  RaterStatus = "assigned"
	RaterStatusStarted   RaterStatus = "started"
	RaterStatusSubmitted RaterStatus = "submitted"
)

type CompetencyAssessmentRater struct {
	ID                       uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	CompetencyEventTargetID  uuid.UUID  `gorm:"type:char(36);not null;index" json:"competency_event_target_id"`
	RaterEmployeeID          uuid.UUID  `gorm:"type:char(36);not null;index" json:"rater_employee_id"`
	RaterType                string     `gorm:"type:varchar(20);not null" json:"rater_type"`
	Weight                   float64    `gorm:"type:decimal(6,2);default:0" json:"weight"`
	Status                   string     `gorm:"type:varchar(20);default:assigned" json:"status"`
	AssignedAt               *time.Time `gorm:"type:timestamp" json:"assigned_at,omitempty"`
	SubmittedAt              *time.Time `gorm:"type:timestamp" json:"submitted_at,omitempty"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`

	// Relasi
	Target    *CompetencyEventTarget `gorm:"foreignKey:CompetencyEventTargetID" json:"target,omitempty"`
	Responses []CompetencyAssessmentResponse `gorm:"foreignKey:RaterID" json:"responses,omitempty"`
}

func (CompetencyAssessmentRater) TableName() string {
	return "competency_assessment_raters"
}

func (r *CompetencyAssessmentRater) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Competency 360 — Assessment Response (plan generik §11)
// =========================================================================

type CompetencyAssessmentResponse struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	RaterID     uuid.UUID  `gorm:"type:char(36);not null;index" json:"rater_id"`
	IndicatorID uuid.UUID  `gorm:"type:char(36);not null;index" json:"indicator_id"`
	RatingValue int        `gorm:"type:smallint;not null" json:"rating_value"`
	Comment     *string    `gorm:"type:text" json:"comment,omitempty"`
	SubmittedAt *time.Time `gorm:"type:timestamp" json:"submitted_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relasi
	Indicator *CompetencyIndicator `gorm:"foreignKey:IndicatorID" json:"indicator,omitempty"`
}

func (CompetencyAssessmentResponse) TableName() string {
	return "competency_assessment_responses"
}

func (r *CompetencyAssessmentResponse) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
