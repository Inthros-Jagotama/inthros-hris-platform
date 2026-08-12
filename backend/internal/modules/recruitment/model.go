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
	ReqStatusDraft       RequisitionStatus = "DRAFT"
	// ReqStatusSubmitted — requisition dirutekan ke Central Approval (G-1).
	// Transisi DRAFT → SUBMITTED terjadi via submit; hasil APPROVED/REJECTED
	// didorong balik oleh modul approval (push-callback).
	ReqStatusSubmitted   RequisitionStatus = "SUBMITTED"
	ReqStatusOpen        RequisitionStatus = "OPEN"
	ReqStatusInProgress  RequisitionStatus = "IN_PROGRESS"
	ReqStatusFilled      RequisitionStatus = "FILLED"
	ReqStatusRejected    RequisitionStatus = "REJECTED"
	ReqStatusCancelled   RequisitionStatus = "CANCELLED"
)

// =========================================================================
// RequisitionPriority (G-2 — job requisition enhancement)
// =========================================================================

type RequisitionPriority string

const (
	ReqPriorityLow    RequisitionPriority = "LOW"
	ReqPriorityMedium RequisitionPriority = "MEDIUM"
	ReqPriorityHigh   RequisitionPriority = "HIGH"
	ReqPriorityUrgent RequisitionPriority = "URGENT"
)

// =========================================================================
// RequisitionReason — alasan pembuatan requisition (strategic layer)
// =========================================================================
// WORKFORCE_GAP menautkan requisition ke hiring need dari Workforce
// Intelligence (workforce_gap_id / workforce_plan_id). Recruitment membaca
// hiring need dari WI via interface narrow — tidak menghitung gap sendiri.
// SUCCESSION_GAP (S-5) menautkan requisition ke posisi kunci tanpa successor
// siap dari Career Intelligence (succession_position_id) — recruitment menjadi
// fallback external recruitment; CI yang menandai gap.

type RequisitionReason string

const (
	ReqReasonNewPosition     RequisitionReason = "NEW_POSITION"
	ReqReasonReplacement     RequisitionReason = "REPLACEMENT"
	ReqReasonExpansion       RequisitionReason = "EXPANSION"
	ReqReasonWorkforceGap    RequisitionReason = "WORKFORCE_GAP"
	ReqReasonSuccessionGap   RequisitionReason = "SUCCESSION_GAP"
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
	// RequisitionNumber (G-2): nomor requisition auto-generated REQ-YYYYMM-XXXX
	// (pola nomor dokumen TRN-* training) — diisi di service bila client kosong.
	RequisitionNumber string            `gorm:"type:varchar(50)" json:"requisition_number,omitempty"`
	// Priority (G-2): LOW | MEDIUM | HIGH | URGENT (default MEDIUM).
	Priority          RequisitionPriority `gorm:"type:varchar(10);default:MEDIUM" json:"priority"`
	// PositionID (G-2): referensi master position (tabel positions). Tanpa FK
	// karena modul Organization tidak mengekspos CRUD position — referensi saja.
	PositionID        *uuid.UUID        `gorm:"type:char(36)" json:"position_id,omitempty"`
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
	// ApprovalInstanceID (G-1): id ApprovalInstance di Central Approval yang
	// memproses requisition ini. Approval module source of truth proses
	// persetujuan; kolom ini hanya referensi.
	ApprovalInstanceID *uuid.UUID       `gorm:"type:char(36)" json:"approval_instance_id,omitempty"`
	ReasonType        string            `gorm:"type:varchar(30)" json:"reason_type"`
	WorkforceGapID    *uuid.UUID        `gorm:"type:char(36)" json:"workforce_gap_id,omitempty"`
	WorkforcePlanID   *uuid.UUID        `gorm:"type:char(36)" json:"workforce_plan_id,omitempty"`
	SuccessionPositionID *uuid.UUID     `gorm:"type:char(36)" json:"succession_position_id,omitempty"`
	TargetStartDate   *string           `gorm:"type:date" json:"target_start_date,omitempty"`
	ClosedAt          *int64            `gorm:"type:bigint;default:0" json:"-"`
	// OpenedAt (G-2): unix nano saat requisition menjadi OPEN — diset otomatis
	// saat approval APPROVED (G-1) atau transisi status manual ke OPEN.
	OpenedAt          *int64            `gorm:"type:bigint;default:0" json:"-"`
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
// OfferStatus (G-3 — offer management)
// =========================================================================

type OfferStatus string

const (
	OfferStatusDraft           OfferStatus = "DRAFT"
	OfferStatusPendingApproval OfferStatus = "PENDING_APPROVAL"
	OfferStatusApproved        OfferStatus = "APPROVED"
	OfferStatusSent            OfferStatus = "SENT"
	OfferStatusAccepted        OfferStatus = "ACCEPTED"
	OfferStatusRejected        OfferStatus = "REJECTED"
	OfferStatusExpired         OfferStatus = "EXPIRED"
	OfferStatusWithdrawn       OfferStatus = "WITHDRAWN"
)

// =========================================================================
// JobOffer (G-3 — Penawaran Kerja)
// =========================================================================
// Offer dibuat dari application (kandidat + requisition), dirutekan melalui
// Central Approval (workflow OFFER_DRAFT → SUBMITTED → APPROVED → SENT),
// lalu kandidat menerima/menolak. Status ACCEPTED menautkan kembali ke
// application (ACCEPTED + slots_filled increment) — fondasi G-4 (employee).

type JobOffer struct {
	ID                uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	ApplicationID     uuid.UUID  `gorm:"type:char(36);not null;index:idx_offer_app" json:"application_id"`
	OfferNumber       string     `gorm:"type:varchar(50)" json:"offer_number,omitempty"`
	EmploymentType    string     `gorm:"type:varchar(50)" json:"employment_type"`
	Salary            float64    `gorm:"type:decimal(15,2);default:0" json:"salary"`
	Allowances        float64    `gorm:"type:decimal(15,2);default:0" json:"allowances"`
	Benefits          string     `gorm:"type:text" json:"benefits"`
	StartDate         string     `gorm:"type:varchar(10)" json:"start_date"`
	ExpiryDate        string     `gorm:"type:varchar(10)" json:"expiry_date"`
	Status            OfferStatus `gorm:"type:varchar(30);default:DRAFT;index:idx_offer_status" json:"status"`
	SentAt            *int64     `gorm:"type:bigint;default:0" json:"-"`
	AcceptedAt        *int64     `gorm:"type:bigint;default:0" json:"-"`
	RejectedAt        *int64     `gorm:"type:bigint;default:0" json:"-"`
	ApprovalInstanceID *uuid.UUID `gorm:"type:char(36)" json:"approval_instance_id,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func (JobOffer) TableName() string {
	return "job_offers"
}

func (o *JobOffer) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
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
	// G-4: jenis kandidat — EXTERNAL (default, dibuatkan employee baru saat
	// offer diterima) atau INTERNAL (menunjuk employee_id yang sudah ada;
	// hasil seleksi diteruskan ke Employee Movement, bukan employee baru).
	CandidateType string         `gorm:"type:varchar(20);default:EXTERNAL;index:idx_cand_type" json:"candidate_type"`
	EmployeeID    *uuid.UUID     `gorm:"type:char(36);index:idx_cand_employee" json:"employee_id,omitempty"`
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
// RecruitmentStage (G-5 — master stage, seeded dari CandidateStatus)
// =========================================================================
// Seeded 1:1 dari 8 CandidateStatus existing (bukan taxonomy baru) — lihat
// docs/superpowers/specs/2026-08-12-recruitment-stage-history-design.md.

type RecruitmentStage struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Code      string    `gorm:"type:varchar(20);not null;uniqueIndex:uq_recruitment_stages_code" json:"code"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	SortOrder int       `gorm:"type:int;default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (RecruitmentStage) TableName() string {
	return "recruitment_stages"
}

func (s *RecruitmentStage) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// ApplicationStageHistory (G-5 — audit trail transisi status aplikasi)
// =========================================================================

type ApplicationStageHistory struct {
	ID            uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	ApplicationID uuid.UUID  `gorm:"type:char(36);not null;index:idx_ash_app" json:"application_id"`
	FromStageID   *uuid.UUID `gorm:"type:char(36)" json:"from_stage_id,omitempty"`
	ToStageID     uuid.UUID  `gorm:"type:char(36);not null" json:"to_stage_id"`
	ChangedBy     *uuid.UUID `gorm:"type:char(36)" json:"changed_by,omitempty"`
	Notes         string     `gorm:"type:text" json:"notes"`
	ChangedAt     int64      `gorm:"type:bigint;not null;index:idx_ash_changed_at" json:"changed_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (ApplicationStageHistory) TableName() string {
	return "job_application_stage_histories"
}

func (h *ApplicationStageHistory) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
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
