package recruitment

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =========================================================================
// CandidateStatus
// =========================================================================

type CandidateStatus string

const (
	CandStatusNew         CandidateStatus = "NEW"
	CandStatusScreened    CandidateStatus = "SCREENED"
	CandStatusShortlisted CandidateStatus = "SHORTLISTED"
	CandStatusInterviewed CandidateStatus = "INTERVIEWED"
	CandStatusOffered     CandidateStatus = "OFFERED"
	CandStatusAccepted    CandidateStatus = "ACCEPTED"
	CandStatusRejected    CandidateStatus = "REJECTED"
	CandStatusWithdrawn   CandidateStatus = "WITHDRAWN"
)

// =========================================================================
// RequisitionStatus
// =========================================================================

type RequisitionStatus string

const (
	ReqStatusDraft     RequisitionStatus = "DRAFT"
	ReqStatusOpen      RequisitionStatus = "OPEN"
	ReqStatusInProgress RequisitionStatus = "IN_PROGRESS"
	ReqStatusFilled    RequisitionStatus = "FILLED"
	ReqStatusCancelled RequisitionStatus = "CANCELLED"
)

// =========================================================================
// InterviewStatus
// =========================================================================

type InterviewStatus string

const (
	IntStatusScheduled InterviewStatus = "SCHEDULED"
	IntStatusCompleted InterviewStatus = "COMPLETED"
	IntStatusCancelled InterviewStatus = "CANCELLED"
	IntStatusRescheduled InterviewStatus = "RESCHEDULED"
)

// =========================================================================
// JobRequisition (Lowongan Pekerjaan)
// =========================================================================

type JobRequisition struct {
	ID                uuid.UUID         `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID    uuid.UUID         `gorm:"type:char(36);not null;index:idx_req_org" json:"organization_id"`
	Title             string            `gorm:"type:varchar(255);not null" json:"title"`
	Department        string            `gorm:"type:varchar(150)" json:"department"`
	EmploymentType    string            `gorm:"type:varchar(50)" json:"employment_type"`
	Location          string            `gorm:"type:varchar(255)" json:"location"`
	MinSalary         float64           `gorm:"type:decimal(15,2);default:0" json:"min_salary"`
	MaxSalary         float64           `gorm:"type:decimal(15,2);default:0" json:"max_salary"`
	Description       string            `gorm:"type:text" json:"description"`
	Requirements      string            `gorm:"type:text" json:"requirements"`
	Responsibilities  string            `gorm:"type:text" json:"responsibilities"`
	SlotsAvailable    int               `gorm:"type:int;default:1" json:"slots_available"`
	SlotsFilled       int               `gorm:"type:int;default:0" json:"slots_filled"`
	Status            RequisitionStatus `gorm:"type:varchar(20);default:DRAFT" json:"status"`
	RequestedBy       *uuid.UUID        `gorm:"type:char(36)" json:"requested_by,omitempty"`
	ApprovedBy        *uuid.UUID        `gorm:"type:char(36)" json:"approved_by,omitempty"`
	TargetStartDate   *string           `gorm:"type:date" json:"target_start_date,omitempty"`
	ClosedAt          *int64            `gorm:"type:bigint;default:0" json:"-"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

func (JobRequisition) TableName() string {
	return "job_requisitions"
}

func (r *JobRequisition) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Candidate (Kandidat Pelamar)
// =========================================================================

type Candidate struct {
	ID            uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	FirstName     string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName      string         `gorm:"type:varchar(100);not null" json:"last_name"`
	Email         string         `gorm:"type:varchar(255);not null;uniqueIndex:idx_cand_email" json:"email"`
	Phone         string         `gorm:"type:varchar(50)" json:"phone"`
	Address       string         `gorm:"type:text" json:"address"`
	CurrentCompany *string       `gorm:"type:varchar(255)" json:"current_company,omitempty"`
	CurrentTitle  *string        `gorm:"type:varchar(255)" json:"current_title,omitempty"`
	ResumeURL     *string        `gorm:"type:text" json:"resume_url,omitempty"`
	PortfolioURL  *string        `gorm:"type:text" json:"portfolio_url,omitempty"`
	LinkedInURL   *string        `gorm:"type:text" json:"linkedin_url,omitempty"`
	Source        string         `gorm:"type:varchar(50);default:direct" json:"source"`
	Notes         string         `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (Candidate) TableName() string {
	return "candidates"
}

func (c *Candidate) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// JobApplication (Lamaran Pekerjaan)
// =========================================================================

type JobApplication struct {
	ID              uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	RequisitionID   uuid.UUID      `gorm:"type:char(36);not null;index:idx_app_req" json:"requisition_id"`
	CandidateID     uuid.UUID      `gorm:"type:char(36);not null;index:idx_app_cand" json:"candidate_id"`
	Status          CandidateStatus `gorm:"type:varchar(50);default:NEW;index:idx_app_status" json:"status"`
	AppliedAt       int64          `gorm:"type:bigint;default:0" json:"-"`
	ScreenedAt      *int64         `gorm:"type:bigint;default:0" json:"-"`
	ShortlistedAt   *int64         `gorm:"type:bigint;default:0" json:"-"`
	OfferedAt       *int64         `gorm:"type:bigint;default:0" json:"-"`
	AcceptedAt      *int64         `gorm:"type:bigint;default:0" json:"-"`
	RejectedAt      *int64         `gorm:"type:bigint;default:0" json:"-"`
	WithdrawnAt     *int64         `gorm:"type:bigint;default:0" json:"-"`
	RejectionReason string         `gorm:"type:text" json:"rejection_reason"`
	Notes           string         `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (JobApplication) TableName() string {
	return "job_applications"
}

func (a *JobApplication) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Interview (Wawancara)
// =========================================================================

type Interview struct {
	ID              uuid.UUID       `gorm:"type:char(36);primaryKey" json:"id"`
	ApplicationID   uuid.UUID       `gorm:"type:char(36);not null;index:idx_int_app" json:"application_id"`
	InterviewerID   uuid.UUID       `gorm:"type:char(36);not null" json:"interviewer_id"`
	Stage           string          `gorm:"type:varchar(50);not null" json:"stage"`
	ScheduledAt     int64           `gorm:"type:bigint;not null;default:0" json:"-"`
	DurationMinutes int             `gorm:"type:int;default:60" json:"duration_minutes"`
	Location        string          `gorm:"type:varchar(255)" json:"location"`
	MeetingLink     *string         `gorm:"type:text" json:"meeting_link,omitempty"`
	Status          InterviewStatus `gorm:"type:varchar(20);default:SCHEDULED" json:"status"`
	Score           *float64        `gorm:"type:decimal(5,2)" json:"score,omitempty"`
	Feedback        string          `gorm:"type:text" json:"feedback"`
	CompletedAt     *int64          `gorm:"type:bigint;default:0" json:"-"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

func (Interview) TableName() string {
	return "interviews"
}

func (i *Interview) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// OnboardingTaskTemplate (Template Tugas Onboarding)
// =========================================================================

type OnboardingTaskTemplate struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"type:varchar(50)" json:"category"`
	DayOffset   int       `gorm:"type:int;default:0" json:"day_offset"`
	AssignedRole string   `gorm:"type:varchar(50)" json:"assigned_role"`
	IsMandatory bool      `gorm:"not null;default:1" json:"is_mandatory"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (OnboardingTaskTemplate) TableName() string {
	return "onboarding_task_templates"
}

func (t *OnboardingTaskTemplate) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// EmployeeOnboarding (Onboarding Karyawan Baru)
// =========================================================================

type EmployeeOnboarding struct {
	ID              uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID      uuid.UUID  `gorm:"type:char(36);not null;index:idx_onb_emp" json:"employee_id"`
	ApplicationID   uuid.UUID  `gorm:"type:char(36);not null;index:idx_onb_app" json:"application_id"`
	StartDate       string     `gorm:"type:date;not null" json:"start_date"`
	Status          string     `gorm:"type:varchar(20);default:PENDING" json:"status"`
	BuddyID         *uuid.UUID `gorm:"type:char(36)" json:"buddy_id,omitempty"`
	CompletedAt     *int64     `gorm:"type:bigint;default:0" json:"-"`
	Notes           string     `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (EmployeeOnboarding) TableName() string {
	return "employee_onboardings"
}

func (o *EmployeeOnboarding) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// OnboardingTaskItem (Tugas spesifik dalam onboarding)
// =========================================================================

type OnboardingTaskItem struct {
	ID                uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeOnboardingID uuid.UUID `gorm:"type:char(36);not null;index:idx_onb_task_item" json:"employee_onboarding_id"`
	TemplateID        *uuid.UUID `gorm:"type:char(36)" json:"template_id,omitempty"`
	Name              string    `gorm:"type:varchar(255);not null" json:"name"`
	Description       string    `gorm:"type:text" json:"description"`
	AssignedTo        *uuid.UUID `gorm:"type:char(36)" json:"assigned_to,omitempty"`
	DueDate           *int64    `gorm:"type:bigint;default:0" json:"-"`
	IsCompleted       bool      `gorm:"not null;default:0" json:"is_completed"`
	CompletedAt       *int64    `gorm:"type:bigint;default:0" json:"-"`
	Notes             string    `gorm:"type:text" json:"notes"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (OnboardingTaskItem) TableName() string {
	return "onboarding_task_items"
}

func (t *OnboardingTaskItem) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
