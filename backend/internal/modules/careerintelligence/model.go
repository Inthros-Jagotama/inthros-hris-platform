package careerintelligence

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =========================================================================
// CareerTalentMap — 9-box grid placement for talent mapping
// =========================================================================

type CareerTalentMap struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID   uuid.UUID `gorm:"type:char(36);not null;index:idx_ctm_employee" json:"employee_id"`
	Period       string    `gorm:"type:char(7);not null;index:idx_ctm_period" json:"period"`
	Performance  string    `gorm:"type:varchar(20);not null" json:"performance"` // LOW / MEDIUM / HIGH
	Potential    string    `gorm:"type:varchar(20);not null" json:"potential"`   // LOW / MEDIUM / HIGH
	GridPosition string    `gorm:"type:varchar(30);not null" json:"grid_position"` // e.g. "9-BOX-1" through "9-BOX-9"
	Notes        string    `gorm:"type:text" json:"notes"`
	AssessorID   uuid.UUID `gorm:"type:char(36)" json:"assessor_id"`
	AssessedAt   time.Time `gorm:"type:date" json:"assessed_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index:idx_ctm_deleted_at" json:"deleted_at,omitempty"`
}

func (CareerTalentMap) TableName() string { return "career_talent_maps" }

func (tm *CareerTalentMap) BeforeCreate(tx *gorm.DB) error {
	if tm.ID == uuid.Nil {
		tm.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// CareerInterest — Employee career interests and aspirations
// =========================================================================

type CareerInterest struct {
	ID              uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID      uuid.UUID `gorm:"type:char(36);not null;index:idx_ci_employee" json:"employee_id"`
	InterestType    string    `gorm:"type:varchar(50);not null" json:"interest_type"` // LEADERSHIP / SPECIALIST / INTERNATIONAL / ENTREPRENEUR
	TargetPosition  string    `gorm:"type:varchar(100)" json:"target_position"`
	TargetDepartment string   `gorm:"type:varchar(100)" json:"target_department"`
	Motivation      string    `gorm:"type:text" json:"motivation"`
	ReadinessLevel  string    `gorm:"type:varchar(20)" json:"readiness_level"` // NOW / 1_YEAR / 2_3_YEARS / 3_PLUS
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	RecordedAt      time.Time `gorm:"type:date" json:"recorded_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index:idx_ci_deleted_at" json:"deleted_at,omitempty"`
}

func (CareerInterest) TableName() string { return "career_interests" }

func (ci *CareerInterest) BeforeCreate(tx *gorm.DB) error {
	if ci.ID == uuid.Nil {
		ci.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// CareerPath — Career path template defining progression routes
// =========================================================================

type CareerPath struct {
	ID            uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	SourceTitleID uuid.UUID `gorm:"type:char(36);not null;index:idx_cp_source" json:"source_title_id"`
	TargetTitleID uuid.UUID `gorm:"type:char(36);not null;index:idx_cp_target" json:"target_title_id"`
	PathType      string    `gorm:"type:varchar(30);not null" json:"path_type"` // PROMOTION / LATERAL / DEMOTION / CROSSFUNCTIONAL
	TypicalTenure int       `gorm:"default:0" json:"typical_tenure"` // months
	Requirements  string    `gorm:"type:text" json:"requirements"`
	Competencies  string    `gorm:"type:text" json:"competencies"` // JSON list of competency IDs
	Certifications string   `gorm:"type:text" json:"certifications"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index:idx_cp_deleted_at" json:"deleted_at,omitempty"`
}

func (CareerPath) TableName() string { return "career_paths" }

func (cp *CareerPath) BeforeCreate(tx *gorm.DB) error {
	if cp.ID == uuid.Nil {
		cp.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// CareerSuccessionPlan — Succession planning for key positions
// =========================================================================

type CareerSuccessionPlan struct {
	ID             uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	PositionID     uuid.UUID `gorm:"type:char(36);not null;index:idx_csp_position" json:"position_id"`
	SuccessorID    uuid.UUID `gorm:"type:char(36);not null;index:idx_csp_successor" json:"successor_id"`
	ReadinessLevel string    `gorm:"type:varchar(20);not null" json:"readiness_level"` // READY_NOW / READY_1YR / READY_2YR / POTENTIAL
	PriorityOrder  int       `gorm:"default:1" json:"priority_order"`
	TargetDate     *time.Time `gorm:"type:date" json:"target_date,omitempty"`
	DevelopmentPlan string   `gorm:"type:text" json:"development_plan"`
	Notes          string    `gorm:"type:text" json:"notes"`
	Status         string    `gorm:"type:varchar(20);default:ACTIVE" json:"status"` // ACTIVE / COMPLETED / REMOVED
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index:idx_csp_deleted_at" json:"deleted_at,omitempty"`
}

func (CareerSuccessionPlan) TableName() string { return "career_succession_plans" }

func (sp *CareerSuccessionPlan) BeforeCreate(tx *gorm.DB) error {
	if sp.ID == uuid.Nil {
		sp.ID = uuid.New()
	}
	return nil
}
