package attendance

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =========================================================================
// 4.1 AttendanceCompanySetting (Konfigurasi Absensi Perusahaan)
// =========================================================================

type AttendanceCompanySetting struct {
	ID                   uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Latitude             *float64       `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude            *float64       `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	IsLocationRequired   bool           `gorm:"not null;default:0" json:"is_location_required"`
	IsFaceRequired       bool           `gorm:"not null;default:0" json:"is_face_required"`
	IsOvertimeEnabled    bool           `gorm:"not null;default:0" json:"is_overtime_enabled"`
	AllowCheckinOnDayOff bool           `gorm:"not null" json:"allow_checkin_on_day_off"`
	MaxDistanceMeter     *int           `json:"max_distance_meter,omitempty"`
	LateToleranceMinutes *int           `json:"late_tolerance_minutes,omitempty"`
	OvertimeMinMinutes   *int           `json:"overtime_min_minutes,omitempty"`
	CreatedBy            *uuid.UUID     `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy            *uuid.UUID     `gorm:"type:char(36)" json:"updated_by,omitempty"`
	DeletedBy            *uuid.UUID     `gorm:"type:char(36)" json:"deleted_by,omitempty"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (AttendanceCompanySetting) TableName() string {
	return "attendance_company_settings"
}

func (s *AttendanceCompanySetting) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 4.2 AttendanceCompanyShift (Shift Kerja Perusahaan)
// =========================================================================

type AttendanceCompanyShift struct {
	ID              uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	ShiftName       string         `gorm:"type:varchar(255);not null" json:"shift_name"`
	CheckInTime     string         `gorm:"type:time;not null" json:"check_in_time"`
	CheckOutTime    string         `gorm:"type:time;not null" json:"check_out_time"`
	IsCrossMidnight bool           `gorm:"not null;default:0" json:"is_cross_midnight"`
	CreatedBy       *uuid.UUID     `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy       *uuid.UUID     `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (AttendanceCompanyShift) TableName() string {
	return "attendance_company_shifts"
}

func (s *AttendanceCompanyShift) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 4.3 AttendanceEmployeeShift (Penugasan Shift ke Karyawan)
// =========================================================================

type AttendanceEmployeeShift struct {
	ID                uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID        uuid.UUID      `gorm:"type:char(36);not null;index" json:"employee_id"`
	AttendanceShiftID uuid.UUID      `gorm:"type:char(36);not null;index" json:"attendance_shift_id"`
	EffectiveDateFrom string         `gorm:"type:date;not null" json:"effective_date_from"`
	EffectiveDateTo   *string        `gorm:"type:date" json:"effective_date_to,omitempty"`
	DaysOfWeekMask    *int           `json:"days_of_week_mask,omitempty"`
	IsDayOff          *bool          `json:"is_day_off,omitempty"`
	CreatedBy         *uuid.UUID     `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy         *uuid.UUID     `gorm:"type:char(36)" json:"updated_by,omitempty"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

func (AttendanceEmployeeShift) TableName() string {
	return "attendance_employee_shifts"
}

func (s *AttendanceEmployeeShift) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 4.4 AttendanceLocation (Lokasi Geofence)
// =========================================================================

type AttendanceLocation struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Latitude  float64        `gorm:"type:decimal(10,7);not null" json:"latitude"`
	Longitude float64        `gorm:"type:decimal(10,7);not null" json:"longitude"`
	RadiusM   int            `gorm:"not null;default:50" json:"radius_m"`
	CreatedBy *uuid.UUID     `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID     `gorm:"type:char(36)" json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (AttendanceLocation) TableName() string {
	return "attendance_locations"
}

func (l *AttendanceLocation) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 4.5 AttendanceDeviceCapture (Perangkat Absensi)
// =========================================================================

type AttendanceDeviceCapture struct {
	ID         uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	DeviceUUID string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"device_uuid"`
	DeviceType string         `gorm:"type:varchar(20);not null" json:"device_type"`
	OsVersion  *string        `gorm:"type:varchar(50)" json:"os_version,omitempty"`
	Model      *string        `gorm:"type:varchar(100)" json:"model,omitempty"`
	AppVersion *string        `gorm:"type:varchar(50)" json:"app_version,omitempty"`
	LastSeenAt *time.Time     `json:"last_seen_at,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (AttendanceDeviceCapture) TableName() string {
	return "attendance_device_captures"
}

func (d *AttendanceDeviceCapture) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 4.6 AttendanceFaceCapture (Wajah Karyawan untuk Verifikasi)
// =========================================================================

type AttendanceFaceCapture struct {
	ID            uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID    uuid.UUID  `gorm:"type:char(36);not null;index:idx_face_employee_time" json:"employee_id"`
	CapturedAt    time.Time  `gorm:"not null" json:"captured_at"`
	ImageURL      string     `gorm:"type:text;not null" json:"image_url"`
	ImageSHA256   string     `gorm:"type:char(64);not null" json:"image_sha256"`
	LivenessScore *float64   `gorm:"type:decimal(6,3)" json:"liveness_score,omitempty"`
	MatchScore    *float64   `gorm:"type:decimal(6,3)" json:"match_score,omitempty"`
	Verified      *bool      `json:"verified,omitempty"`
	Provider      *string    `gorm:"type:varchar(50)" json:"provider,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (AttendanceFaceCapture) TableName() string {
	return "attendance_face_captures"
}

func (f *AttendanceFaceCapture) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 4.7 AttendanceEvent (Event Check-in / Check-out)
// =========================================================================

type EventType string

const (
	EventTypeCheckIn  EventType = "CHECKIN"
	EventTypeCheckOut EventType = "CHECKOUT"
)

type ValidationStatus string

const (
	ValidationPending    ValidationStatus = "PENDING"
	ValidationValid      ValidationStatus = "VALID"
	ValidationInvalid    ValidationStatus = "INVALID"
	ValidationOverridden ValidationStatus = "OVERRIDDEN"
)

type AttendanceEvent struct {
	ID                  uuid.UUID        `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID          uuid.UUID        `gorm:"type:char(36);not null;index:idx_att_event_employee" json:"employee_id"`
	OvertimeRequestID   *uuid.UUID       `gorm:"type:char(36);index:idx_att_event_overtime" json:"overtime_request_id,omitempty"`
	EventType           EventType        `gorm:"type:varchar(255);not null" json:"event_type"`
	EventTimeUTC        time.Time        `gorm:"not null" json:"event_time_utc"`
	EventTimeLocal      time.Time        `gorm:"not null" json:"event_time_local"`
	DeviceID            *uuid.UUID       `gorm:"type:char(36);index:idx_att_event_device" json:"device_id,omitempty"`
	Latitude            float64          `gorm:"type:decimal(10,7);not null" json:"latitude"`
	Longitude           float64          `gorm:"type:decimal(10,7);not null" json:"longitude"`
	AccuracyM           *int             `json:"accuracy_m,omitempty"`
	LocationProvider    *string          `gorm:"type:varchar(255)" json:"location_provider,omitempty"`
	ValidatedLocationID *uuid.UUID       `gorm:"type:char(36);index:idx_att_event_location" json:"validated_location_id,omitempty"`
	DistanceM           *int             `json:"distance_m,omitempty"`
	IsInGeofence        bool             `gorm:"not null;default:0" json:"is_in_geofence"`
	FaceCaptureID       *uuid.UUID       `gorm:"type:char(36);index:idx_att_event_face" json:"face_capture_id,omitempty"`
	ValidationStatus    ValidationStatus `gorm:"type:varchar(255);default:PENDING" json:"validation_status"`
	ValidationNote      *string          `gorm:"type:varchar(255)" json:"validation_note,omitempty"`
	DeletedAt           gorm.DeletedAt   `gorm:"index:idx_att_event_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt           time.Time        `json:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at"`
}

func (AttendanceEvent) TableName() string {
	return "attendance_events"
}

func (e *AttendanceEvent) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 4.8 AttendanceSession (Sesi Kerja Harian)
// =========================================================================

type SessionStatus string

const (
	SessionStatusOpen            SessionStatus = "OPEN"
	SessionStatusClosed          SessionStatus = "CLOSED"
	SessionStatusMissingCheckIn  SessionStatus = "MISSING_CHECKIN"
	SessionStatusMissingCheckOut SessionStatus = "MISSING_CHECKOUT"
	SessionStatusAbsent          SessionStatus = "ABSENT"
	SessionStatusDayOff          SessionStatus = "DAY_OFF"
	SessionStatusExempt          SessionStatus = "EXEMPT"
	SessionStatusLeave           SessionStatus = "LEAVE"
)

type AttendanceSession struct {
	ID                         uuid.UUID       `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID                 uuid.UUID       `gorm:"type:char(36);not null;uniqueIndex:uk_session_emp_date" json:"employee_id"`
	WorkDate                   string          `gorm:"not null;uniqueIndex:uk_session_emp_date" json:"work_date"`
	ShiftID                    *uuid.UUID      `gorm:"type:char(36);index:idx_att_session_shift" json:"shift_id,omitempty"`
	IsOvertimeDay              bool            `json:"is_overtime_day"`
	OvertimeRequestID          *uuid.UUID      `gorm:"type:char(36);index:idx_att_session_overtime" json:"overtime_request_id,omitempty"`
	ApprovedOvertimeStartLocal *time.Time      `json:"approved_overtime_start_local,omitempty"`
	ApprovedOvertimeEndLocal   *time.Time      `json:"approved_overtime_end_local,omitempty"`
	LeaveRequestID             *uuid.UUID      `gorm:"type:char(36);index:idx_att_session_leave" json:"leave_request_id,omitempty"`
	LeaveFraction              *float64        `gorm:"type:decimal(4,2)" json:"leave_fraction,omitempty"`
	PlannedStartLocal          *time.Time      `json:"planned_start_local,omitempty"`
	PlannedEndLocal            *time.Time      `json:"planned_end_local,omitempty"`
	CheckinEventID             *uuid.UUID      `gorm:"type:char(36);index:idx_att_session_checkin" json:"checkin_event_id,omitempty"`
	CheckoutEventID            *uuid.UUID      `gorm:"type:char(36);index:idx_att_session_checkout" json:"checkout_event_id,omitempty"`
	Status                     SessionStatus   `gorm:"type:varchar(255);default:OPEN" json:"status"`
	LatenessMinutes            int             `gorm:"not null;default:0" json:"lateness_minutes"`
	EarlyLeaveMinutes          int             `gorm:"not null;default:0" json:"early_leave_minutes"`
	WorkMinutes                int             `gorm:"not null;default:0" json:"work_minutes"`
	BreakMinutes               int             `gorm:"not null;default:0" json:"break_minutes"`
	OvertimeMinutes            int             `gorm:"not null;default:0" json:"overtime_minutes"`
	DeletedAt                  gorm.DeletedAt  `gorm:"index:idx_att_session_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt                  time.Time       `json:"created_at"`
	UpdatedAt                  time.Time       `json:"updated_at"`
}

func (AttendanceSession) TableName() string {
	return "attendance_sessions"
}

func (s *AttendanceSession) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 4.9 AttendanceOvertimeRequest (Pengajuan Lembur)
// =========================================================================

type OvertimeStatus string

const (
	OvertimeSubmitted      OvertimeStatus = "SUBMITTED"
	OvertimePendingApproval OvertimeStatus = "PENDING_APPROVAL"
	OvertimeApproved       OvertimeStatus = "APPROVED"
	OvertimeRejected       OvertimeStatus = "REJECTED"
)

type AttendanceOvertimeRequest struct {
	ID                uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID        uuid.UUID      `gorm:"type:char(36);not null;index:idx_att_overtime_employee" json:"employee_id"`
	WorkDate          string         `gorm:"type:date;not null;index:idx_att_overtime_date" json:"work_date"`
	StartTimeLocal    time.Time      `gorm:"not null" json:"start_time_local"`
	EndTimeLocal      time.Time      `gorm:"not null" json:"end_time_local"`
	RequestedMinutes  int            `gorm:"not null" json:"requested_minutes"`
	ActualMinutes     *int           `json:"actual_minutes,omitempty"`
	CalculatedMinutes *int           `json:"calculated_minutes,omitempty"`
	Reason            *string        `gorm:"type:varchar(255)" json:"reason,omitempty"`
	Status            OvertimeStatus `gorm:"type:varchar(255);default:SUBMITTED;index:idx_att_overtime_status" json:"status"`
	ApprovedBy        *uuid.UUID     `gorm:"type:char(36)" json:"approved_by,omitempty"`
	ApprovedAt        *time.Time     `json:"approved_at,omitempty"`
	ApprovalNote      *string        `gorm:"type:varchar(255)" json:"approval_note,omitempty"`
	ApprovalInstanceID *uuid.UUID    `gorm:"type:char(36);index:idx_att_overtime_approval_instance" json:"approval_instance_id,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

func (AttendanceOvertimeRequest) TableName() string {
	return "attendance_overtime_requests"
}

func (r *AttendanceOvertimeRequest) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 4.10 AttendanceExemptPosition (Posisi yang Dibebaskan dari Absensi)
// =========================================================================

type AttendanceExemptPosition struct {
	ID             uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID uuid.UUID      `gorm:"type:char(36);not null;uniqueIndex:uk_exempt_org" json:"organization_id"`
	IsExempt       bool           `gorm:"not null" json:"is_exempt"`
	Note           *string        `gorm:"type:varchar(255)" json:"note,omitempty"`
	CreatedBy      *uuid.UUID     `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy      *uuid.UUID     `gorm:"type:char(36)" json:"updated_by,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (AttendanceExemptPosition) TableName() string {
	return "attendance_exempt_positions"
}

func (p *AttendanceExemptPosition) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// AttendanceCorrectionRequest (§16/§33-34) - Phase 8
// =========================================================================

type CorrectionType string

const (
	CorrectionTypeMissingCheckin  CorrectionType = "MISSING_CHECKIN"
	CorrectionTypeMissingCheckout CorrectionType = "MISSING_CHECKOUT"
	CorrectionTypeWrongCheckin    CorrectionType = "WRONG_CHECKIN"
	CorrectionTypeWrongCheckout   CorrectionType = "WRONG_CHECKOUT"
)

type CorrectionStatus string

const (
	CorrectionSubmitted      CorrectionStatus = "SUBMITTED"
	CorrectionPendingApproval CorrectionStatus = "PENDING_APPROVAL"
	CorrectionApproved       CorrectionStatus = "APPROVED"
	CorrectionRejected       CorrectionStatus = "REJECTED"
)

// AttendanceCorrectionRequest lets an employee/HR request a fix to a day's
// session (missing/wrong check-in or check-out) without ever mutating the
// raw attendance_events history (§15's immutability principle) - approval
// applies the requested times to the session and triggers a recalculation.
type AttendanceCorrectionRequest struct {
	ID                  uuid.UUID        `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID          uuid.UUID        `gorm:"type:char(36);not null;index:idx_att_correction_employee" json:"employee_id"`
	AttendanceSessionID uuid.UUID        `gorm:"type:char(36);not null;index:idx_att_correction_session" json:"attendance_session_id"`
	CorrectionType      CorrectionType   `gorm:"type:varchar(50);not null" json:"correction_type"`
	RequestedCheckin    *time.Time       `json:"requested_checkin,omitempty"`
	RequestedCheckout   *time.Time       `json:"requested_checkout,omitempty"`
	Reason              string           `gorm:"type:varchar(255);not null" json:"reason"`
	Status              CorrectionStatus `gorm:"type:varchar(255);default:SUBMITTED;index:idx_att_correction_status" json:"status"`
	ApprovalInstanceID  *uuid.UUID       `gorm:"type:char(36);index:idx_att_correction_approval_instance" json:"approval_instance_id,omitempty"`
	CreatedBy           *uuid.UUID       `gorm:"type:char(36)" json:"created_by,omitempty"`
	ApprovedAt          *time.Time       `json:"approved_at,omitempty"`
	CreatedAt           time.Time        `json:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at"`
}

func (AttendanceCorrectionRequest) TableName() string {
	return "attendance_correction_requests"
}

func (c *AttendanceCorrectionRequest) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
