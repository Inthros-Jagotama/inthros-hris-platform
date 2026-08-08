package attendance

import "time"

// =========================================================================
// Company Settings DTOs
// =========================================================================

type CreateCompanySettingRequest struct {
	Latitude             *float64 `json:"latitude"`
	Longitude            *float64 `json:"longitude"`
	IsLocationRequired   *bool    `json:"is_location_required"`
	IsFaceRequired       *bool    `json:"is_face_required"`
	IsOvertimeEnabled    *bool    `json:"is_overtime_enabled"`
	MaxDistanceMeter     *int     `json:"max_distance_meter"`
	LateToleranceMinutes *int     `json:"late_tolerance_minutes"`
	OvertimeMinMinutes   *int     `json:"overtime_min_minutes"`
}

type UpdateCompanySettingRequest struct {
	Latitude             *float64 `json:"latitude"`
	Longitude            *float64 `json:"longitude"`
	IsLocationRequired   *bool    `json:"is_location_required"`
	IsFaceRequired       *bool    `json:"is_face_required"`
	IsOvertimeEnabled    *bool    `json:"is_overtime_enabled"`
	MaxDistanceMeter     *int     `json:"max_distance_meter"`
	LateToleranceMinutes *int     `json:"late_tolerance_minutes"`
	OvertimeMinMinutes   *int     `json:"overtime_min_minutes"`
}

type CompanySettingResponse struct {
	ID                   string    `json:"id"`
	Latitude             *float64  `json:"latitude,omitempty"`
	Longitude            *float64  `json:"longitude,omitempty"`
	IsLocationRequired   bool      `json:"is_location_required"`
	IsFaceRequired       bool      `json:"is_face_required"`
	IsOvertimeEnabled    bool      `json:"is_overtime_enabled"`
	MaxDistanceMeter     *int      `json:"max_distance_meter,omitempty"`
	LateToleranceMinutes *int      `json:"late_tolerance_minutes,omitempty"`
	OvertimeMinMinutes   *int      `json:"overtime_min_minutes,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// =========================================================================
// Company Shift DTOs
// =========================================================================

type CreateCompanyShiftRequest struct {
	ShiftName       string `json:"shift_name" binding:"required"`
	CheckInTime     string `json:"check_in_time" binding:"required"`
	CheckOutTime    string `json:"check_out_time" binding:"required"`
	IsCrossMidnight *bool  `json:"is_cross_midnight"`
}

type UpdateCompanyShiftRequest struct {
	ShiftName       *string `json:"shift_name"`
	CheckInTime     *string `json:"check_in_time"`
	CheckOutTime    *string `json:"check_out_time"`
	IsCrossMidnight *bool   `json:"is_cross_midnight"`
}

type CompanyShiftResponse struct {
	ID              string    `json:"id"`
	ShiftName       string    `json:"shift_name"`
	CheckInTime     string    `json:"check_in_time"`
	CheckOutTime    string    `json:"check_out_time"`
	IsCrossMidnight bool      `json:"is_cross_midnight"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// =========================================================================
// Employee Shift DTOs
// =========================================================================

type CreateEmployeeShiftRequest struct {
	EmployeeID        string  `json:"employee_id" binding:"required"`
	AttendanceShiftID string  `json:"attendance_shift_id" binding:"required"`
	EffectiveDateFrom string  `json:"effective_date_from" binding:"required"`
	EffectiveDateTo   *string `json:"effective_date_to"`
	DaysOfWeekMask    *int    `json:"days_of_week_mask"`
	IsDayOff          *bool   `json:"is_day_off"`
}

type UpdateEmployeeShiftRequest struct {
	AttendanceShiftID *string `json:"attendance_shift_id"`
	EffectiveDateFrom *string `json:"effective_date_from"`
	EffectiveDateTo   *string `json:"effective_date_to"`
	DaysOfWeekMask    *int    `json:"days_of_week_mask"`
	IsDayOff          *bool   `json:"is_day_off"`
}

type EmployeeShiftResponse struct {
	ID                string    `json:"id"`
	EmployeeID        string    `json:"employee_id"`
	AttendanceShiftID string    `json:"attendance_shift_id"`
	EffectiveDateFrom string    `json:"effective_date_from"`
	EffectiveDateTo   *string   `json:"effective_date_to,omitempty"`
	DaysOfWeekMask    *int      `json:"days_of_week_mask,omitempty"`
	IsDayOff          *bool     `json:"is_day_off,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// =========================================================================
// Location DTOs
// =========================================================================

type CreateLocationRequest struct {
	Name      string  `json:"name" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	RadiusM   *int    `json:"radius_m"`
}

type UpdateLocationRequest struct {
	Name      *string  `json:"name"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	RadiusM   *int     `json:"radius_m"`
}

type LocationResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	RadiusM   int       `json:"radius_m"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// =========================================================================
// Event DTOs
// =========================================================================

type CreateEventRequest struct {
	EmployeeID       string  `json:"employee_id" binding:"required"`
	EventType        string  `json:"event_type" binding:"required"`
	EventTimeUTC     string  `json:"event_time_utc" binding:"required"`
	EventTimeLocal   string  `json:"event_time_local" binding:"required"`
	DeviceID         *string `json:"device_id"`
	Latitude         float64 `json:"latitude" binding:"required"`
	Longitude        float64 `json:"longitude" binding:"required"`
	AccuracyM        *int    `json:"accuracy_m"`
	LocationProvider *string `json:"location_provider"`
}

type EventResponse struct {
	ID                  string    `json:"id"`
	EmployeeID          string    `json:"employee_id"`
	EventType           string    `json:"event_type"`
	EventTimeUTC        time.Time `json:"event_time_utc"`
	EventTimeLocal      time.Time `json:"event_time_local"`
	DeviceID            *string   `json:"device_id,omitempty"`
	Latitude            float64   `json:"latitude"`
	Longitude           float64   `json:"longitude"`
	AccuracyM           *int      `json:"accuracy_m,omitempty"`
	LocationProvider    *string   `json:"location_provider,omitempty"`
	ValidatedLocationID *string   `json:"validated_location_id,omitempty"`
	DistanceM           *int      `json:"distance_m,omitempty"`
	IsInGeofence        bool      `json:"is_in_geofence"`
	ValidationStatus    string    `json:"validation_status"`
	ValidationNote      *string   `json:"validation_note,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
}

// =========================================================================
// Session DTOs
// =========================================================================

type SessionResponse struct {
	ID                string     `json:"id"`
	EmployeeID        string     `json:"employee_id"`
	WorkDate          string     `json:"work_date"`
	ShiftID           *string    `json:"shift_id,omitempty"`
	IsOvertimeDay     bool       `json:"is_overtime_day"`
	OvertimeRequestID *string    `json:"overtime_request_id,omitempty"`
	LeaveRequestID    *string    `json:"leave_request_id,omitempty"`
	PlannedStartLocal *time.Time `json:"planned_start_local,omitempty"`
	PlannedEndLocal   *time.Time `json:"planned_end_local,omitempty"`
	CheckinEventID    *string    `json:"checkin_event_id,omitempty"`
	CheckoutEventID   *string    `json:"checkout_event_id,omitempty"`
	Status            string     `json:"status"`
	LatenessMinutes   int        `json:"lateness_minutes"`
	EarlyLeaveMinutes int        `json:"early_leave_minutes"`
	WorkMinutes       int        `json:"work_minutes"`
	BreakMinutes      int        `json:"break_minutes"`
	OvertimeMinutes   int        `json:"overtime_minutes"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// =========================================================================
// Summary DTOs (Phase 10 - Employee Dashboard/Calendar, §36-37)
// =========================================================================

// SummaryResponse aggregates one employee's sessions over a date range.
// Absent isn't included: attendance_sessions never gets a status of
// SessionStatusAbsent anywhere in this codebase yet (no scheduled
// ProcessDailyAttendance/DetectMissingAttendance job exists to set it - see
// §44-45), so a count would always read zero and misleadingly imply
// absences are being tracked when they aren't.
type SummaryResponse struct {
	EmployeeID          string  `json:"employee_id"`
	FromDate            string  `json:"from_date"`
	ToDate              string  `json:"to_date"`
	TotalSessions       int     `json:"total_sessions"`
	PresentDays         int     `json:"present_days"`
	LateDays            int     `json:"late_days"`
	MissingCheckinDays  int     `json:"missing_checkin_days"`
	MissingCheckoutDays int     `json:"missing_checkout_days"`
	DayOffDays          int     `json:"day_off_days"`
	LeaveDays           float64 `json:"leave_days"`
	TotalWorkMinutes    int     `json:"total_work_minutes"`
	TotalOvertimeMinutes int    `json:"total_overtime_minutes"`
}

// =========================================================================
// Overtime Request DTOs
// =========================================================================

type CreateOvertimeRequest struct {
	EmployeeID       string  `json:"employee_id" binding:"required"`
	WorkDate         string  `json:"work_date" binding:"required"`
	StartTimeLocal   string  `json:"start_time_local" binding:"required"`
	EndTimeLocal     string  `json:"end_time_local" binding:"required"`
	RequestedMinutes int     `json:"requested_minutes" binding:"required"`
	Reason           string  `json:"reason"`
	FlowID           *string `json:"flow_id"`
}

type OvertimeResponse struct {
	ID                 string     `json:"id"`
	EmployeeID         string     `json:"employee_id"`
	WorkDate           string     `json:"work_date"`
	StartTimeLocal     time.Time  `json:"start_time_local"`
	EndTimeLocal       time.Time  `json:"end_time_local"`
	RequestedMinutes   int        `json:"requested_minutes"`
	ActualMinutes      *int       `json:"actual_minutes,omitempty"`
	CalculatedMinutes  *int       `json:"calculated_minutes,omitempty"`
	Reason             *string    `json:"reason,omitempty"`
	Status             string     `json:"status"`
	ApprovedBy         *string    `json:"approved_by,omitempty"`
	ApprovedAt         *time.Time `json:"approved_at,omitempty"`
	ApprovalNote       *string    `json:"approval_note,omitempty"`
	ApprovalInstanceID *string    `json:"approval_instance_id,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
}

// =========================================================================
// Exempt Position DTOs
// =========================================================================

type CreateExemptPositionRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	IsExempt       *bool  `json:"is_exempt"`
	Note           string `json:"note"`
}

type UpdateExemptPositionRequest struct {
	IsExempt *bool   `json:"is_exempt"`
	Note     *string `json:"note"`
}

type ExemptPositionResponse struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	IsExempt       bool      `json:"is_exempt"`
	Note           *string   `json:"note,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// =========================================================================
// Correction Request DTOs
// =========================================================================

type CreateCorrectionRequest struct {
	EmployeeID          string  `json:"employee_id" binding:"required"`
	AttendanceSessionID string  `json:"attendance_session_id" binding:"required"`
	CorrectionType      string  `json:"correction_type" binding:"required"`
	RequestedCheckin    *string `json:"requested_checkin"`
	RequestedCheckout   *string `json:"requested_checkout"`
	Reason              string  `json:"reason" binding:"required"`
	FlowID              *string `json:"flow_id"`
}

type CorrectionResponse struct {
	ID                  string     `json:"id"`
	EmployeeID          string     `json:"employee_id"`
	AttendanceSessionID string     `json:"attendance_session_id"`
	CorrectionType      string     `json:"correction_type"`
	RequestedCheckin    *time.Time `json:"requested_checkin,omitempty"`
	RequestedCheckout   *time.Time `json:"requested_checkout,omitempty"`
	Reason              string     `json:"reason"`
	Status              string     `json:"status"`
	ApprovalInstanceID  *string    `json:"approval_instance_id,omitempty"`
	ApprovedAt          *time.Time `json:"approved_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
}

// =========================================================================
// Pagination
// =========================================================================

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}
