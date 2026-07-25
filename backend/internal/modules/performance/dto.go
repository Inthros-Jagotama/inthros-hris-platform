package performance

import "time"

// =========================================================================
// Performance Period DTOs
// =========================================================================

type CreatePerformancePeriodRequest struct {
	PeriodCode string `json:"period_code" binding:"required,max=10"`
	PeriodType string `json:"period_type" binding:"required,oneof=MONTHLY QUARTERLY SEMESTER ANNUAL"`
	Year       int    `json:"year" binding:"required"`
	StartDate  *string `json:"start_date"`
	EndDate    *string `json:"end_date"`
	Status     *string `json:"status" binding:"omitempty,oneof=draft active closed"`
}

type UpdatePerformancePeriodRequest struct {
	PeriodCode *string `json:"period_code" binding:"omitempty,max=10"`
	PeriodType *string `json:"period_type" binding:"omitempty,oneof=MONTHLY QUARTERLY SEMESTER ANNUAL"`
	Year       *int    `json:"year"`
	StartDate  *string `json:"start_date"`
	EndDate    *string `json:"end_date"`
	Status     *string `json:"status" binding:"omitempty,oneof=draft active closed"`
}

type PerformancePeriodResponse struct {
	ID         string `json:"id"`
	PeriodCode string `json:"period_code"`
	PeriodType string `json:"period_type"`
	Year       int    `json:"year"`
	StartDate  string `json:"start_date,omitempty"`
	EndDate    string `json:"end_date,omitempty"`
	Status     string `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// =========================================================================
// Performance Perspective DTOs
// =========================================================================

type CreatePerformancePerspectiveRequest struct {
	Name        string  `json:"name" binding:"required,max=100"`
	Description *string `json:"description"`
	SortOrder   *int    `json:"sort_order"`
}

type UpdatePerformancePerspectiveRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=100"`
	Description *string `json:"description"`
	SortOrder   *int    `json:"sort_order"`
}

type PerformancePerspectiveResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// =========================================================================
// Performance Template DTOs
// =========================================================================

type CreatePerformanceTemplateRequest struct {
	OrganizationID string  `json:"organization_id" binding:"required"`
	Name           string  `json:"name" binding:"required,max=200"`
	Description    *string `json:"description"`
}

type UpdatePerformanceTemplateRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=200"`
	Description *string `json:"description"`
	Status      *string `json:"status" binding:"omitempty,oneof=DRAFT PUBLISHED ARCHIVED"`
}

type PerformanceTemplateResponse struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description,omitempty"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// =========================================================================
// Performance Indicator DTOs
// =========================================================================

type CreatePerformanceIndicatorRequest struct {
	PerformanceTemplateID string   `json:"performance_template_id" binding:"required"`
	PerspectiveID        string   `json:"perspective_id" binding:"required"`
	IndicatorType        string   `json:"indicator_type" binding:"required,oneof=MAXIMIZATION MINIMIZATION"`
	Title                string   `json:"title" binding:"required,max=255"`
	Description          *string  `json:"description"`
	Weight               float64  `json:"weight"`
	TargetValue          float64  `json:"target_value"`
	UnitOfMeasurement    *string  `json:"unit_of_measurement" binding:"omitempty,max=50"`
	SortOrder            *int     `json:"sort_order"`
}

type UpdatePerformanceIndicatorRequest struct {
	IndicatorType     *string  `json:"indicator_type" binding:"omitempty,oneof=MAXIMIZATION MINIMIZATION"`
	Title             *string  `json:"title" binding:"omitempty,max=255"`
	Description       *string  `json:"description"`
	Weight            *float64 `json:"weight"`
	TargetValue       *float64 `json:"target_value"`
	UnitOfMeasurement *string  `json:"unit_of_measurement" binding:"omitempty,max=50"`
	SortOrder         *int     `json:"sort_order"`
}

type PerformanceIndicatorResponse struct {
	ID                    string    `json:"id"`
	PerformanceTemplateID string    `json:"performance_template_id"`
	PerspectiveID         string    `json:"perspective_id"`
	IndicatorType         string    `json:"indicator_type"`
	Title                 string    `json:"title"`
	Description           string    `json:"description,omitempty"`
	Weight                float64   `json:"weight"`
	TargetValue           float64   `json:"target_value"`
	UnitOfMeasurement     string    `json:"unit_of_measurement,omitempty"`
	SortOrder             int       `json:"sort_order"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// =========================================================================
// Performance Evaluation DTOs
// =========================================================================

type CreatePerformanceEvaluationRequest struct {
	EmployeeID     string `json:"employee_id" binding:"required"`
	OrganizationID string `json:"organization_id" binding:"required"`
	PeriodID       string `json:"period_id" binding:"required"`
	TemplateID     string `json:"template_id" binding:"required"`
	SupervisorID   *string `json:"supervisor_id"`
	Notes          *string `json:"notes"`
}

type UpdatePerformanceEvaluationRequest struct {
	SupervisorID *string `json:"supervisor_id"`
	Notes        *string `json:"notes"`
}

type UpdateEvaluationStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=DRAFT PLAN_SUBMITTED PLAN_APPROVED ACTUAL_SUBMITTED ACTUAL_APPROVED COMPLETED"`
	Notes  string `json:"notes"`
}

type PerformanceEvaluationResponse struct {
	ID              string    `json:"id"`
	EmployeeID      string    `json:"employee_id"`
	OrganizationID  string    `json:"organization_id"`
	PeriodID        string    `json:"period_id"`
	TemplateID      string    `json:"template_id"`
	SupervisorID    string    `json:"supervisor_id,omitempty"`
	FinalScore      float64   `json:"final_score"`
	Status          string    `json:"status"`
	Notes           string    `json:"notes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// =========================================================================
// Performance EvaluationDetail DTOs
// =========================================================================

type CreateEvaluationDetailRequest struct {
	PerformanceEvaluationID string  `json:"performance_evaluation_id" binding:"required"`
	PerspectiveID          string  `json:"perspective_id" binding:"required"`
	AchievementPercentage  float64 `json:"achievement_percentage"`
	Weight                 float64 `json:"weight"`
	Score                  float64 `json:"score"`
	Description            *string `json:"description"`
}

type UpdateEvaluationDetailRequest struct {
	AchievementPercentage *float64 `json:"achievement_percentage"`
	Weight                *float64 `json:"weight"`
	Score                 *float64 `json:"score"`
	Description           *string  `json:"description"`
}

type EvaluationDetailResponse struct {
	ID                      string    `json:"id"`
	PerformanceEvaluationID string    `json:"performance_evaluation_id"`
	PerspectiveID          string    `json:"perspective_id"`
	AchievementPercentage  float64   `json:"achievement_percentage"`
	Weight                 float64   `json:"weight"`
	Score                  float64   `json:"score"`
	Description            string    `json:"description,omitempty"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// =========================================================================
// Performance Target DTOs
// =========================================================================

type CreatePerformanceTargetRequest struct {
	PerformanceEvaluationID string  `json:"performance_evaluation_id" binding:"required"`
	IndicatorID            string  `json:"indicator_id" binding:"required"`
	TargetValue            float64 `json:"target_value"`
	ActualValue            *float64 `json:"actual_value"`
	UnitOfMeasurement      *string `json:"unit_of_measurement" binding:"omitempty,max=50"`
	Weight                 float64 `json:"weight"`
}

type UpdatePerformanceTargetRequest struct {
	TargetValue   *float64 `json:"target_value"`
	ActualValue   *float64 `json:"actual_value"`
	UnitOfMeasurement *string `json:"unit_of_measurement" binding:"omitempty,max=50"`
	Weight        *float64 `json:"weight"`
}

type PerformanceTargetResponse struct {
	ID                      string    `json:"id"`
	PerformanceEvaluationID string    `json:"performance_evaluation_id"`
	IndicatorID            string    `json:"indicator_id"`
	TargetValue            float64   `json:"target_value"`
	ActualValue            float64   `json:"actual_value,omitempty"`
	UnitOfMeasurement      string    `json:"unit_of_measurement,omitempty"`
	AchievementPercentage  float64   `json:"achievement_percentage"`
	Weight                 float64   `json:"weight"`
	Score                  float64   `json:"score"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
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
