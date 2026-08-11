package training

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =========================================================================
// TrainingStatus constants
// =========================================================================

type AttendanceStatus string

const (
	AttendStatusPresent AttendanceStatus = "PRESENT"
	AttendStatusAbsent  AttendanceStatus = "ABSENT"
	AttendStatusLate    AttendanceStatus = "LATE"
	AttendStatusExcused AttendanceStatus = "EXCUSED"
)

type SessionStatus string

const (
	SessStatusDraft           SessionStatus = "DRAFT"
	SessStatusScheduled       SessionStatus = "SCHEDULED"
	SessStatusRegistrationOpen SessionStatus = "REGISTRATION_OPEN"
	SessStatusFull            SessionStatus = "FULL"
	SessStatusInProgress      SessionStatus = "IN_PROGRESS"
	SessStatusCompleted       SessionStatus = "COMPLETED"
	SessStatusCancelled       SessionStatus = "CANCELLED"
)

// ProviderType — tipe penyelenggara training (in-house / external).
type ProviderType string

const (
	ProviderTypeInternal ProviderType = "INTERNAL"
	ProviderTypeExternal ProviderType = "EXTERNAL"
)

// DeliveryMode — mode penyampaian session.
type DeliveryMode string

const (
	DeliveryModeOnsite    DeliveryMode = "ONSITE"
	DeliveryModeOnline    DeliveryMode = "ONLINE"
	DeliveryModeHybrid    DeliveryMode = "HYBRID"
	DeliveryModeSelfPaced DeliveryMode = "SELF_PACED"
)

// CourseType — kategori jenis kursus.
type CourseType string

const (
	CourseTypeTechnical     CourseType = "TECHNICAL"
	CourseTypeSoftSkill     CourseType = "SOFT_SKILL"
	CourseTypeCompliance    CourseType = "COMPLIANCE"
	CourseTypeManagement    CourseType = "MANAGEMENT"
	CourseTypeCertification CourseType = "CERTIFICATION"
	CourseTypeOther         CourseType = "OTHER"
)

// DeliveryType — default/preferred penyelenggaraan pada master course.
type DeliveryType string

const (
	DeliveryTypeInHouse  DeliveryType = "IN_HOUSE"
	DeliveryTypeExternal DeliveryType = "EXTERNAL"
	DeliveryTypeBoth     DeliveryType = "BOTH"
)

// RegistrationStatus — status pendaftaran/enrollment peserta.
type RegistrationStatus string

const (
	RegStatusNominated  RegistrationStatus = "NOMINATED"
	RegStatusRequested  RegistrationStatus = "REQUESTED"
	RegStatusApproved   RegistrationStatus = "APPROVED"
	RegStatusRegistered RegistrationStatus = "REGISTERED"
	RegStatusWaitlisted RegistrationStatus = "WAITLISTED"
	RegStatusCancelled  RegistrationStatus = "CANCELLED"
)

// CompletionStatus — status penyelesaian peserta.
type CompletionStatus string

const (
	CompletionNotStarted CompletionStatus = "NOT_STARTED"
	CompletionInProgress CompletionStatus = "IN_PROGRESS"
	CompletionCompleted  CompletionStatus = "COMPLETED"
	CompletionFailed     CompletionStatus = "FAILED"
)

// TrainerType — tipe trainer.
type TrainerType string

const (
	TrainerTypeInternal TrainerType = "INTERNAL"
	TrainerTypeExternal TrainerType = "EXTERNAL"
)

// SessionTrainerRole — peran trainer dalam satu session.
type SessionTrainerRole string

const (
	SessionTrainerRoleMain      SessionTrainerRole = "MAIN"
	SessionTrainerRoleAssistant SessionTrainerRole = "ASSISTANT"
)

// AssessmentType — tipe assessment.
type AssessmentType string

const (
	AssessTypePreTest   AssessmentType = "PRE_TEST"
	AssessTypePostTest  AssessmentType = "POST_TEST"
	AssessTypeFinal     AssessmentType = "FINAL"
	AssessTypePractical AssessmentType = "PRACTICAL"
	AssessTypeOther     AssessmentType = "OTHER"
)

// =========================================================================
// TrainingCategory — Kategori Pelatihan
// =========================================================================

type TrainingCategory struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code        string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_trn_cat_code" json:"code"`
	Name        string         `gorm:"type:varchar(150);not null" json:"name"`
	Description *string        `gorm:"type:varchar(500)" json:"description,omitempty"`
	IsActive    bool           `gorm:"not null;default:1" json:"is_active"`
	DeletedAt   gorm.DeletedAt `gorm:"index:idx_trn_cat_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (TrainingCategory) TableName() string { return "training_categories" }

func (c *TrainingCategory) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// TrainingCourse — Master Data Kursus/Pelatihan
// =========================================================================

type TrainingCourse struct {
	ID             uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	CategoryID     uuid.UUID      `gorm:"type:char(36);not null;index:idx_trn_course_cat" json:"category_id"`
	Code           string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_trn_course_code" json:"code"`
	Name           string         `gorm:"type:varchar(200);not null" json:"name"`
	Description    *string        `gorm:"type:text" json:"description,omitempty"`
	DurationHour   *float64       `gorm:"type:decimal(8,2)" json:"duration_hour,omitempty"`
	MinScore       *float64       `gorm:"type:decimal(5,2)" json:"min_score,omitempty"`
	Cost           *float64       `gorm:"type:decimal(18,2)" json:"cost,omitempty"`
	IsCertified    bool           `gorm:"not null;default:0" json:"is_certified"`
	ExternalVendor *string        `gorm:"type:varchar(200)" json:"external_vendor,omitempty"`
	// P0-BE (plan §7): tipe kursus + default delivery + mandatory flag.
	CourseType   *CourseType  `gorm:"type:varchar(20)" json:"course_type,omitempty"`
	DeliveryType *DeliveryType `gorm:"type:varchar(20)" json:"delivery_type,omitempty"`
	IsMandatory  bool         `gorm:"not null;default:0" json:"is_mandatory"`
	IsActive     bool         `gorm:"not null;default:1" json:"is_active"`
	DeletedAt    gorm.DeletedAt `gorm:"index:idx_trn_course_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func (TrainingCourse) TableName() string { return "training_courses" }

func (c *TrainingCourse) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// TrainingSession — Sesi/Kelas Pelatihan
// =========================================================================

type TrainingSession struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	CourseID    uuid.UUID      `gorm:"type:char(36);not null;index:idx_trn_sess_course" json:"course_id"`
	SessionCode string         `gorm:"type:varchar(20);not null" json:"session_code"`
	TrainerName string         `gorm:"type:varchar(200);not null" json:"trainer_name"`
	// P0-BE (plan §14): penyelenggara + mode delivery + link provider + datetime presisi.
	ProviderType         *ProviderType `gorm:"type:varchar(20)" json:"provider_type,omitempty"`
	DeliveryMode         *DeliveryMode `gorm:"type:varchar(20)" json:"delivery_mode,omitempty"`
	ProviderID           *uuid.UUID    `gorm:"type:char(36);index:idx_trn_sess_provider" json:"provider_id,omitempty"`
	StartDatetime        *time.Time    `gorm:"type:timestamp" json:"start_datetime,omitempty"`
	EndDatetime          *time.Time    `gorm:"type:timestamp" json:"end_datetime,omitempty"`
	MeetingURL           *string       `gorm:"type:text" json:"meeting_url,omitempty"`
	RegistrationDeadline *time.Time    `gorm:"type:timestamp" json:"registration_deadline,omitempty"`
	Location             *string       `gorm:"type:varchar(255)" json:"location,omitempty"`
	StartDate            string        `gorm:"type:date;not null" json:"start_date"`
	EndDate              string        `gorm:"type:date;not null" json:"end_date"`
	MaxQuota             int           `gorm:"type:int;not null;default:30" json:"max_quota"`
	Status               SessionStatus `gorm:"type:varchar(20);not null;default:SCHEDULED;index:idx_trn_sess_status" json:"status"`
	DeletedAt            gorm.DeletedAt `gorm:"index:idx_trn_sess_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`

	// Relasi
	Participants []TrainingParticipant `gorm:"foreignKey:SessionID" json:"participants,omitempty"`
}

func (TrainingSession) TableName() string { return "training_sessions" }

func (s *TrainingSession) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// TrainingParticipant — Peserta Pelatihan
// =========================================================================

type TrainingParticipant struct {
	ID               uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	SessionID        uuid.UUID      `gorm:"type:char(36);not null;index:idx_trn_part_sess" json:"session_id"`
	EmployeeID       uuid.UUID      `gorm:"type:char(36);not null;index:idx_trn_part_emp" json:"employee_id"`
	// P0-BE (plan §18): enrollment — registration + completion.
	RegistrationStatus RegistrationStatus `gorm:"type:varchar(20);not null;default:REGISTERED" json:"registration_status"`
	RegisteredAt       *time.Time         `gorm:"type:timestamp" json:"registered_at,omitempty"`
	ApprovedAt         *time.Time         `gorm:"type:timestamp" json:"approved_at,omitempty"`
	AttendanceStatus   AttendanceStatus   `gorm:"type:varchar(20);not null;default:PRESENT" json:"attendance_status"`
	Score              float64            `gorm:"type:decimal(5,2);default:0" json:"score"`
	CompletionStatus   CompletionStatus   `gorm:"type:varchar(20);not null;default:NOT_STARTED" json:"completion_status"`
	CompletionDate     *string            `gorm:"type:date" json:"completion_date,omitempty"`
	FinalScore         *float64           `gorm:"type:decimal(5,2)" json:"final_score,omitempty"`
	Passed             *bool              `json:"passed,omitempty"`
	Remarks            *string            `gorm:"type:text" json:"remarks,omitempty"`
	CompletedAt        *string            `gorm:"type:date" json:"completed_at,omitempty"`
	DeletedAt          gorm.DeletedAt     `gorm:"index:idx_trn_part_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}

func (TrainingParticipant) TableName() string { return "training_participants" }

func (p *TrainingParticipant) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// TrainingMaterial — Materi Pelatihan
// =========================================================================

type TrainingMaterial struct {
	ID           uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	SessionID    uuid.UUID      `gorm:"type:char(36);not null;index:idx_trn_mat_sess" json:"session_id"`
	Title        string         `gorm:"type:varchar(200);not null" json:"title"`
	// P0-BE (plan §20): deskripsi, wajib, ketersediaan.
	Description   *string    `gorm:"type:text" json:"description,omitempty"`
	IsRequired    bool       `gorm:"not null;default:0" json:"is_required"`
	AvailableFrom *time.Time `gorm:"type:timestamp" json:"available_from,omitempty"`
	FileURL       *string    `gorm:"type:text" json:"file_url,omitempty"`
	FileType      *string    `gorm:"type:varchar(50)" json:"file_type,omitempty"`
	SortOrder     int        `gorm:"type:smallint;default:0" json:"sort_order"`
	DeletedAt     gorm.DeletedAt `gorm:"index:idx_trn_mat_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (TrainingMaterial) TableName() string { return "training_materials" }

func (m *TrainingMaterial) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// TrainingEvaluation — Evaluasi Pelatihan
// =========================================================================

type TrainingEvaluation struct {
	ID         uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	SessionID  uuid.UUID      `gorm:"type:char(36);not null;index:idx_trn_eval_sess" json:"session_id"`
	EmployeeID uuid.UUID      `gorm:"type:char(36);not null;index:idx_trn_eval_emp" json:"employee_id"`
	Rating     int            `gorm:"type:tinyint;not null" json:"rating"`
	Feedback   *string        `gorm:"type:text" json:"feedback,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index:idx_trn_eval_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (TrainingEvaluation) TableName() string { return "training_evaluations" }

func (e *TrainingEvaluation) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// TrainingCertificate — Sertifikat Pelatihan
// =========================================================================

type TrainingCertificate struct {
	ID            uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	ParticipantID uuid.UUID      `gorm:"type:char(36);not null;index:idx_trn_cert_part" json:"participant_id"`
	CertificateNo string         `gorm:"type:varchar(50);not null;uniqueIndex:idx_trn_cert_no" json:"certificate_no"`
	IssuedDate    string         `gorm:"type:date;not null" json:"issued_date"`
	ExpiryDate    *string        `gorm:"type:date" json:"expiry_date,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"index:idx_trn_cert_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (TrainingCertificate) TableName() string { return "training_certificates" }

func (c *TrainingCertificate) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// TrainingProvider — Master Penyelenggara (INTERNAL | EXTERNAL)
// =========================================================================

type TrainingProvider struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code        string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_trn_provider_code" json:"code"`
	Name        string         `gorm:"type:varchar(200);not null" json:"name"`
	Type        ProviderType   `gorm:"type:varchar(20);not null;default:EXTERNAL;index:idx_trn_provider_type" json:"type"`
	ContactName *string        `gorm:"type:varchar(150)" json:"contact_name,omitempty"`
	Email       *string        `gorm:"type:varchar(150)" json:"email,omitempty"`
	Phone       *string        `gorm:"type:varchar(50)" json:"phone,omitempty"`
	Address     *string        `gorm:"type:text" json:"address,omitempty"`
	Website     *string        `gorm:"type:varchar(200)" json:"website,omitempty"`
	IsActive    bool           `gorm:"not null;default:1" json:"is_active"`
	DeletedAt   gorm.DeletedAt `gorm:"index:idx_trn_provider_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (TrainingProvider) TableName() string { return "training_providers" }

func (p *TrainingProvider) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// TrainingTrainer — Trainer (INTERNAL: employee / EXTERNAL: provider)
// =========================================================================

type TrainingTrainer struct {
	ID         uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Type       TrainerType    `gorm:"type:varchar(20);not null" json:"type"`
	EmployeeID *uuid.UUID     `gorm:"type:char(36);index:idx_trn_trainer_emp" json:"employee_id,omitempty"`
	ProviderID *uuid.UUID     `gorm:"type:char(36);index:idx_trn_trainer_provider" json:"provider_id,omitempty"`
	Name       string         `gorm:"type:varchar(200);not null" json:"name"`
	Email      *string        `gorm:"type:varchar(150)" json:"email,omitempty"`
	Phone      *string        `gorm:"type:varchar(50)" json:"phone,omitempty"`
	Bio        *string        `gorm:"type:text" json:"bio,omitempty"`
	IsActive   bool           `gorm:"not null;default:1" json:"is_active"`
	DeletedAt  gorm.DeletedAt `gorm:"index:idx_trn_trainer_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (TrainingTrainer) TableName() string { return "training_trainers" }

func (t *TrainingTrainer) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// TrainingSessionTrainer — relasi session ↔ trainer (banyak trainer per session)
// =========================================================================

type TrainingSessionTrainer struct {
	ID        uuid.UUID         `gorm:"type:char(36);primaryKey" json:"id"`
	SessionID uuid.UUID         `gorm:"type:char(36);not null;index:idx_trn_sess_trn_session" json:"session_id"`
	TrainerID uuid.UUID         `gorm:"type:char(36);not null;index:idx_trn_sess_trn_trainer" json:"trainer_id"`
	Role      SessionTrainerRole `gorm:"type:varchar(20);not null;default:MAIN" json:"role"`
	DeletedAt gorm.DeletedAt     `gorm:"index:idx_trn_sess_trn_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

func (TrainingSessionTrainer) TableName() string { return "training_session_trainers" }

func (st *TrainingSessionTrainer) BeforeCreate(tx *gorm.DB) error {
	if st.ID == uuid.Nil {
		st.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// TrainingAttendance — Attendance detail per hari (multi-day session)
// =========================================================================

type TrainingAttendance struct {
	ID             uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	ParticipantID  uuid.UUID      `gorm:"type:char(36);not null;index:idx_trn_att_part" json:"participant_id"`
	AttendanceDate string         `gorm:"type:date;not null" json:"attendance_date"`
	CheckIn        *time.Time     `gorm:"type:timestamp" json:"check_in,omitempty"`
	CheckOut       *time.Time     `gorm:"type:timestamp" json:"check_out,omitempty"`
	Status         AttendanceStatus `gorm:"type:varchar(20);not null;default:PRESENT" json:"status"`
	Remarks        *string        `gorm:"type:text" json:"remarks,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"index:idx_trn_att_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (TrainingAttendance) TableName() string { return "training_attendances" }

func (a *TrainingAttendance) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// TrainingAssessment — Definisi assessment per session
// =========================================================================

type TrainingAssessment struct {
	ID            uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	SessionID     uuid.UUID      `gorm:"type:char(36);not null;index:idx_trn_assess_session" json:"session_id"`
	Name          string         `gorm:"type:varchar(200);not null" json:"name"`
	Type          AssessmentType `gorm:"type:varchar(20);not null;default:OTHER" json:"type"`
	MaxScore      float64        `gorm:"type:decimal(8,2);not null;default:100" json:"max_score"`
	PassingScore  float64        `gorm:"type:decimal(8,2);not null;default:60" json:"passing_score"`
	AttemptLimit  int            `gorm:"type:int;not null;default:1" json:"attempt_limit"`
	IsRequired    bool           `gorm:"not null;default:1" json:"is_required"`
	DeletedAt     gorm.DeletedAt `gorm:"index:idx_trn_assess_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (TrainingAssessment) TableName() string { return "training_assessments" }

func (a *TrainingAssessment) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// TrainingAssessmentResult — Nilai per peserta per attempt
// =========================================================================

type TrainingAssessmentResult struct {
	ID            uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	AssessmentID  uuid.UUID      `gorm:"type:char(36);not null;index:idx_trn_assess_res_assessment" json:"assessment_id"`
	ParticipantID uuid.UUID      `gorm:"type:char(36);not null;index:idx_trn_assess_res_part" json:"participant_id"`
	Score         float64        `gorm:"type:decimal(8,2);not null;default:0" json:"score"`
	Passed        bool           `gorm:"not null;default:0" json:"passed"`
	Attempt       int            `gorm:"type:int;not null;default:1" json:"attempt"`
	CompletedAt   *time.Time     `gorm:"type:timestamp" json:"completed_at,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"index:idx_trn_assess_res_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (TrainingAssessmentResult) TableName() string { return "training_assessment_results" }

func (r *TrainingAssessmentResult) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
