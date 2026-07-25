package training

import "time"

// =========================================================================
// Training Category DTOs
// =========================================================================

type CreateTrainingCategoryRequest struct {
	Code        string  `json:"code" binding:"required,max=20"`
	Name        string  `json:"name" binding:"required,max=150"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateTrainingCategoryRequest struct {
	Code        *string `json:"code" binding:"omitempty,max=20"`
	Name        *string `json:"name" binding:"omitempty,max=150"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type TrainingCategoryResponse struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// =========================================================================
// Training Course DTOs
// =========================================================================

type CreateTrainingCourseRequest struct {
	CategoryID     string   `json:"category_id" binding:"required"`
	Code           string   `json:"code" binding:"required,max=20"`
	Name           string   `json:"name" binding:"required,max=200"`
	Description    *string  `json:"description"`
	DurationHour   *float64 `json:"duration_hour"`
	MinScore       *float64 `json:"min_score"`
	Cost           *float64 `json:"cost"`
	IsCertified    *bool    `json:"is_certified"`
	ExternalVendor *string  `json:"external_vendor" binding:"omitempty,max=200"`
}

type UpdateTrainingCourseRequest struct {
	CategoryID     *string  `json:"category_id"`
	Code           *string  `json:"code" binding:"omitempty,max=20"`
	Name           *string  `json:"name" binding:"omitempty,max=200"`
	Description    *string  `json:"description"`
	DurationHour   *float64 `json:"duration_hour"`
	MinScore       *float64 `json:"min_score"`
	Cost           *float64 `json:"cost"`
	IsCertified    *bool    `json:"is_certified"`
	ExternalVendor *string  `json:"external_vendor" binding:"omitempty,max=200"`
	IsActive       *bool    `json:"is_active"`
}

type TrainingCourseResponse struct {
	ID             string    `json:"id"`
	CategoryID     string    `json:"category_id"`
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	Description    string    `json:"description,omitempty"`
	DurationHour   float64   `json:"duration_hour,omitempty"`
	MinScore       float64   `json:"min_score,omitempty"`
	Cost           float64   `json:"cost,omitempty"`
	IsCertified    bool      `json:"is_certified"`
	ExternalVendor string    `json:"external_vendor,omitempty"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// =========================================================================
// Training Session DTOs
// =========================================================================

type CreateTrainingSessionRequest struct {
	CourseID    string  `json:"course_id" binding:"required"`
	SessionCode string  `json:"session_code" binding:"required,max=20"`
	TrainerName string  `json:"trainer_name" binding:"required,max=200"`
	Location    *string `json:"location" binding:"omitempty,max=255"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     string  `json:"end_date" binding:"required"`
	MaxQuota    int     `json:"max_quota" binding:"omitempty,min=1"`
}

type UpdateTrainingSessionRequest struct {
	SessionCode *string `json:"session_code" binding:"omitempty,max=20"`
	TrainerName *string `json:"trainer_name" binding:"omitempty,max=200"`
	Location    *string `json:"location" binding:"omitempty,max=255"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
	MaxQuota    *int    `json:"max_quota" binding:"omitempty,min=1"`
}

type UpdateSessionStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=SCHEDULED IN_PROGRESS COMPLETED CANCELLED"`
}

type TrainingSessionResponse struct {
	ID          string    `json:"id"`
	CourseID    string    `json:"course_id"`
	SessionCode string    `json:"session_code"`
	TrainerName string    `json:"trainer_name"`
	Location    string    `json:"location,omitempty"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
	MaxQuota    int       `json:"max_quota"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// =========================================================================
// Training Participant DTOs
// =========================================================================

type CreateTrainingParticipantRequest struct {
	SessionID  string `json:"session_id" binding:"required"`
	EmployeeID string `json:"employee_id" binding:"required"`
}

type UpdateTrainingParticipantRequest struct {
	AttendanceStatus *string `json:"attendance_status" binding:"omitempty,oneof=PRESENT ABSENT EXCUSED"`
	Score            *float64 `json:"score"`
}

type TrainingParticipantResponse struct {
	ID               string    `json:"id"`
	SessionID        string    `json:"session_id"`
	EmployeeID       string    `json:"employee_id"`
	AttendanceStatus string    `json:"attendance_status"`
	Score            float64   `json:"score"`
	CompletedAt      string    `json:"completed_at,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// =========================================================================
// Training Material DTOs
// =========================================================================

type CreateTrainingMaterialRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	Title     string `json:"title" binding:"required,max=200"`
	FileURL   *string `json:"file_url"`
	FileType  *string `json:"file_type" binding:"omitempty,max=50"`
	SortOrder *int    `json:"sort_order"`
}

type UpdateTrainingMaterialRequest struct {
	Title     *string `json:"title" binding:"omitempty,max=200"`
	FileURL   *string `json:"file_url"`
	FileType  *string `json:"file_type" binding:"omitempty,max=50"`
	SortOrder *int    `json:"sort_order"`
}

type TrainingMaterialResponse struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Title     string    `json:"title"`
	FileURL   string    `json:"file_url,omitempty"`
	FileType  string    `json:"file_type,omitempty"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// =========================================================================
// Training Evaluation DTOs
// =========================================================================

type CreateTrainingEvaluationRequest struct {
	SessionID  string `json:"session_id" binding:"required"`
	EmployeeID string `json:"employee_id" binding:"required"`
	Rating     int    `json:"rating" binding:"required,min=1,max=5"`
	Feedback   *string `json:"feedback"`
}

type UpdateTrainingEvaluationRequest struct {
	Rating   *int    `json:"rating" binding:"omitempty,min=1,max=5"`
	Feedback *string `json:"feedback"`
}

type TrainingEvaluationResponse struct {
	ID         string    `json:"id"`
	SessionID  string    `json:"session_id"`
	EmployeeID string    `json:"employee_id"`
	Rating     int       `json:"rating"`
	Feedback   string    `json:"feedback,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// =========================================================================
// Training Certificate DTOs
// =========================================================================

type CreateTrainingCertificateRequest struct {
	ParticipantID string `json:"participant_id" binding:"required"`
	CertificateNo string `json:"certificate_no" binding:"required,max=50"`
	IssuedDate    string `json:"issued_date" binding:"required"`
	ExpiryDate    *string `json:"expiry_date"`
}

type UpdateTrainingCertificateRequest struct {
	CertificateNo *string `json:"certificate_no" binding:"omitempty,max=50"`
	IssuedDate    *string `json:"issued_date"`
	ExpiryDate    *string `json:"expiry_date"`
}

type TrainingCertificateResponse struct {
	ID            string    `json:"id"`
	ParticipantID string    `json:"participant_id"`
	CertificateNo string    `json:"certificate_no"`
	IssuedDate    string    `json:"issued_date"`
	ExpiryDate    string    `json:"expiry_date,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
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
