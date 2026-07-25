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
	AttendStatusExcused AttendanceStatus = "EXCUSED"
)

type SessionStatus string

const (
	SessStatusScheduled  SessionStatus = "SCHEDULED"
	SessStatusInProgress SessionStatus = "IN_PROGRESS"
	SessStatusCompleted  SessionStatus = "COMPLETED"
	SessStatusCancelled  SessionStatus = "CANCELLED"
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
	ID            uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	CategoryID    uuid.UUID      `gorm:"type:char(36);not null;index:idx_trn_course_cat" json:"category_id"`
	Code          string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_trn_course_code" json:"code"`
	Name          string         `gorm:"type:varchar(200);not null" json:"name"`
	Description   *string        `gorm:"type:text" json:"description,omitempty"`
	DurationHour  *float64       `gorm:"type:decimal(8,2)" json:"duration_hour,omitempty"`
	MinScore      *float64       `gorm:"type:decimal(5,2)" json:"min_score,omitempty"`
	Cost          *float64       `gorm:"type:decimal(18,2)" json:"cost,omitempty"`
	IsCertified   bool           `gorm:"not null;default:0" json:"is_certified"`
	ExternalVendor *string       `gorm:"type:varchar(200)" json:"external_vendor,omitempty"`
	IsActive      bool           `gorm:"not null;default:1" json:"is_active"`
	DeletedAt     gorm.DeletedAt `gorm:"index:idx_trn_course_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
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
	Location    *string        `gorm:"type:varchar(255)" json:"location,omitempty"`
	StartDate   string         `gorm:"type:date;not null" json:"start_date"`
	EndDate     string         `gorm:"type:date;not null" json:"end_date"`
	MaxQuota    int            `gorm:"type:int;not null;default:30" json:"max_quota"`
	Status      SessionStatus  `gorm:"type:varchar(20);not null;default:SCHEDULED;index:idx_trn_sess_status" json:"status"`
	DeletedAt   gorm.DeletedAt `gorm:"index:idx_trn_sess_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

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
	AttendanceStatus AttendanceStatus `gorm:"type:varchar(20);not null;default:PRESENT" json:"attendance_status"`
	Score            float64        `gorm:"type:decimal(5,2);default:0" json:"score"`
	CompletedAt      *string        `gorm:"type:date" json:"completed_at,omitempty"`
	DeletedAt        gorm.DeletedAt `gorm:"index:idx_trn_part_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
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
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	SessionID uuid.UUID      `gorm:"type:char(36);not null;index:idx_trn_mat_sess" json:"session_id"`
	Title     string         `gorm:"type:varchar(200);not null" json:"title"`
	FileURL   *string        `gorm:"type:text" json:"file_url,omitempty"`
	FileType  *string        `gorm:"type:varchar(50)" json:"file_type,omitempty"`
	SortOrder int            `gorm:"type:smallint;default:0" json:"sort_order"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_trn_mat_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
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
