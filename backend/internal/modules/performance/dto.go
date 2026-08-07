package performance

import "time"

// =========================================================================
// Performance Period DTOs
// =========================================================================

type CreatePerformancePeriodRequest struct {
	PeriodCode string  `json:"period_code" binding:"required,max=10"`
	PeriodType string  `json:"period_type" binding:"required,oneof=MONTHLY QUARTERLY SEMESTER ANNUAL"`
	Year       int     `json:"year" binding:"required"`
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
	ID         string    `json:"id"`
	PeriodCode string    `json:"period_code"`
	PeriodType string    `json:"period_type"`
	Year       int       `json:"year"`
	StartDate  string    `json:"start_date,omitempty"`
	EndDate    string    `json:"end_date,omitempty"`
	Status     string    `json:"status"`
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
	PeriodID       *string `json:"period_id"`
	Name           string  `json:"name" binding:"required,max=200"`
	Description    *string `json:"description"`
	EffectiveDate  *string `json:"effective_date"`
	ExpiredDate    *string `json:"expired_date"`
}

type UpdatePerformanceTemplateRequest struct {
	OrganizationID *string `json:"organization_id" binding:"omitempty,uuid"`
	PeriodID       *string `json:"period_id"`
	Name           *string `json:"name" binding:"omitempty,max=200"`
	Description    *string `json:"description"`
	Status         *string `json:"status" binding:"omitempty,oneof=DRAFT PUBLISHED ARCHIVED"`
	EffectiveDate  *string `json:"effective_date"`
	ExpiredDate    *string `json:"expired_date"`
}

type PerformanceTemplateResponse struct {
	ID               string    `json:"id"`
	OrganizationID   string    `json:"organization_id"`
	OrganizationName string    `json:"organization_name,omitempty"`
	PeriodID         string    `json:"period_id,omitempty"`
	PeriodCode       string    `json:"period_code,omitempty"`
	Name             string    `json:"name"`
	Description      string    `json:"description,omitempty"`
	Status           string    `json:"status"`
	EffectiveDate    string    `json:"effective_date,omitempty"`
	ExpiredDate      string    `json:"expired_date,omitempty"`
	IndicatorCount   int       `json:"indicator_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// MyKPIContextResponse resolves the logged-in user to their current employee
// record and Organization (posisi jabatan terakhir), plus the PUBLISHED
// templates available for that Organization — used by the employee
// self-assessment page to know whether it can offer "fill evaluation" at all.
type MyKPIContextResponse struct {
	HasPosition      bool                          `json:"has_position"`
	EmployeeID       string                        `json:"employee_id,omitempty"`
	OrganizationID   string                        `json:"organization_id,omitempty"`
	OrganizationName string                        `json:"organization_name,omitempty"`
	Templates        []PerformanceTemplateResponse `json:"templates"`
}

// =========================================================================
// Performance Indicator DTOs
// =========================================================================

type CreatePerformanceIndicatorRequest struct {
	PerformanceTemplateID string   `json:"performance_template_id" binding:"required"`
	PerspectiveID         string   `json:"perspective_id" binding:"required"`
	Code                  *string  `json:"code" binding:"omitempty,max=50"`
	IndicatorType         string   `json:"indicator_type" binding:"required,oneof=MAXIMIZATION MINIMIZATION"`
	Title                 string   `json:"title" binding:"required,max=255"`
	Description           *string  `json:"description"`
	Weight                float64  `json:"weight"`
	TargetValue           float64  `json:"target_value"`
	UnitOfMeasurement     *string  `json:"unit_of_measurement" binding:"omitempty,max=50"`
	FormulaType           *string  `json:"formula_type" binding:"omitempty,oneof=MANUAL HIGHER_BETTER LOWER_BETTER RANGE"`
	MinimumScore          *float64 `json:"minimum_score"`
	MaximumScore          *float64 `json:"maximum_score"`
	TargetType            *string  `json:"target_type" binding:"omitempty,oneof=NUMBER CURRENCY PERCENTAGE DURATION BOOLEAN"`
	IsRequired            *bool    `json:"is_required"`
	SortOrder             *int     `json:"sort_order"`
}

type UpdatePerformanceIndicatorRequest struct {
	PerspectiveID     *string  `json:"perspective_id" binding:"omitempty,uuid"`
	Code              *string  `json:"code" binding:"omitempty,max=50"`
	IndicatorType     *string  `json:"indicator_type" binding:"omitempty,oneof=MAXIMIZATION MINIMIZATION"`
	Title             *string  `json:"title" binding:"omitempty,max=255"`
	Description       *string  `json:"description"`
	Weight            *float64 `json:"weight"`
	TargetValue       *float64 `json:"target_value"`
	UnitOfMeasurement *string  `json:"unit_of_measurement" binding:"omitempty,max=50"`
	FormulaType       *string  `json:"formula_type" binding:"omitempty,oneof=MANUAL HIGHER_BETTER LOWER_BETTER RANGE"`
	MinimumScore      *float64 `json:"minimum_score"`
	MaximumScore      *float64 `json:"maximum_score"`
	TargetType        *string  `json:"target_type" binding:"omitempty,oneof=NUMBER CURRENCY PERCENTAGE DURATION BOOLEAN"`
	IsRequired        *bool    `json:"is_required"`
	SortOrder         *int     `json:"sort_order"`
}

type PerformanceIndicatorResponse struct {
	ID                    string    `json:"id"`
	PerformanceTemplateID string    `json:"performance_template_id"`
	PerspectiveID         string    `json:"perspective_id"`
	Code                  string    `json:"code,omitempty"`
	IndicatorType         string    `json:"indicator_type"`
	Title                 string    `json:"title"`
	Description           string    `json:"description,omitempty"`
	Weight                float64   `json:"weight"`
	TargetValue           float64   `json:"target_value"`
	UnitOfMeasurement     string    `json:"unit_of_measurement,omitempty"`
	FormulaType           string    `json:"formula_type"`
	MinimumScore          float64   `json:"minimum_score"`
	MaximumScore          float64   `json:"maximum_score"`
	TargetType            string    `json:"target_type"`
	IsRequired            bool      `json:"is_required"`
	SortOrder             int       `json:"sort_order"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// =========================================================================
// Performance Evaluation DTOs
// =========================================================================

type CreatePerformanceEvaluationRequest struct {
	EmployeeID     string  `json:"employee_id" binding:"required"`
	OrganizationID string  `json:"organization_id" binding:"required"`
	PeriodID       string  `json:"period_id" binding:"required"`
	TemplateID     string  `json:"template_id" binding:"required"`
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
	ID             string    `json:"id"`
	EmployeeID     string    `json:"employee_id"`
	OrganizationID string    `json:"organization_id"`
	PeriodID       string    `json:"period_id"`
	TemplateID     string    `json:"template_id"`
	SupervisorID   string    `json:"supervisor_id,omitempty"`
	FinalScore     float64   `json:"final_score"`
	RatingID       string    `json:"rating_id,omitempty"`
	Status         string    `json:"status"`
	SubmittedAt    string    `json:"submitted_at,omitempty"`
	ApprovedAt     string    `json:"approved_at,omitempty"`
	Notes          string    `json:"notes,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// =========================================================================
// Performance EvaluationDetail DTOs
// =========================================================================

type CreateEvaluationDetailRequest struct {
	PerformanceEvaluationID string  `json:"performance_evaluation_id" binding:"required"`
	PerspectiveID           string  `json:"perspective_id" binding:"required"`
	IndicatorID             *string `json:"indicator_id"`
	IndicatorName           *string `json:"indicator_name"`
	AchievementPercentage   float64 `json:"achievement_percentage"`
	Weight                  float64 `json:"weight"`
	Target                  float64 `json:"target"`
	Actual                  float64 `json:"actual"`
	Achievement             float64 `json:"achievement"`
	Score                   float64 `json:"score"`
	Description             *string `json:"description"`
	Remarks                 *string `json:"remarks"`
}

type UpdateEvaluationDetailRequest struct {
	AchievementPercentage *float64 `json:"achievement_percentage"`
	Weight                *float64 `json:"weight"`
	Target                *float64 `json:"target"`
	Actual                *float64 `json:"actual"`
	Achievement           *float64 `json:"achievement"`
	Score                 *float64 `json:"score"`
	Description           *string  `json:"description"`
	Remarks               *string  `json:"remarks"`
}

type EvaluationDetailResponse struct {
	ID                      string    `json:"id"`
	PerformanceEvaluationID string    `json:"performance_evaluation_id"`
	PerspectiveID           string    `json:"perspective_id"`
	PerspectiveName         string    `json:"perspective_name,omitempty"`
	IndicatorID             string    `json:"indicator_id,omitempty"`
	IndicatorName           string    `json:"indicator_name,omitempty"`
	UnitOfMeasurement       string    `json:"unit_of_measurement,omitempty"`
	FormulaType             string    `json:"formula_type,omitempty"`
	AchievementPercentage   float64   `json:"achievement_percentage"`
	Weight                  float64   `json:"weight"`
	Target                  float64   `json:"target"`
	Actual                  float64   `json:"actual"`
	Achievement             float64   `json:"achievement"`
	Score                   float64   `json:"score"`
	Description             string    `json:"description,omitempty"`
	Remarks                 string    `json:"remarks,omitempty"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

// =========================================================================
// Performance Target DTOs
// =========================================================================

type CreatePerformanceTargetRequest struct {
	PerformanceEvaluationID string   `json:"performance_evaluation_id" binding:"required"`
	IndicatorID             string   `json:"indicator_id" binding:"required"`
	TargetValue             float64  `json:"target_value"`
	ActualValue             *float64 `json:"actual_value"`
	UnitOfMeasurement       *string  `json:"unit_of_measurement" binding:"omitempty,max=50"`
	Weight                  float64  `json:"weight"`
}

type UpdatePerformanceTargetRequest struct {
	TargetValue       *float64 `json:"target_value"`
	ActualValue       *float64 `json:"actual_value"`
	UnitOfMeasurement *string  `json:"unit_of_measurement" binding:"omitempty,max=50"`
	Weight            *float64 `json:"weight"`
}

type PerformanceTargetResponse struct {
	ID                      string    `json:"id"`
	PerformanceEvaluationID string    `json:"performance_evaluation_id"`
	IndicatorID             string    `json:"indicator_id"`
	TargetValue             float64   `json:"target_value"`
	ActualValue             float64   `json:"actual_value,omitempty"`
	UnitOfMeasurement       string    `json:"unit_of_measurement,omitempty"`
	AchievementPercentage   float64   `json:"achievement_percentage"`
	Weight                  float64   `json:"weight"`
	Score                   float64   `json:"score"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
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

// =========================================================================
// Performance Progress DTOs
// =========================================================================

type CreatePerformanceProgressRequest struct {
	EvaluationDetailID string  `json:"evaluation_detail_id" binding:"required"`
	ProgressDate       string  `json:"progress_date" binding:"required"`
	ActualValue        float64 `json:"actual_value"`
	Achievement        float64 `json:"achievement"`
	Notes              *string `json:"notes"`
	CreatedBy          string  `json:"created_by" binding:"required"`
}

type UpdatePerformanceProgressRequest struct {
	ProgressDate *string  `json:"progress_date"`
	ActualValue  *float64 `json:"actual_value"`
	Achievement  *float64 `json:"achievement"`
	Notes        *string  `json:"notes"`
}

type PerformanceProgressResponse struct {
	ID                 string    `json:"id"`
	EvaluationDetailID string    `json:"evaluation_detail_id"`
	ProgressDate       string    `json:"progress_date"`
	ActualValue        float64   `json:"actual_value"`
	Achievement        float64   `json:"achievement"`
	Notes              string    `json:"notes,omitempty"`
	CreatedBy          string    `json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// =========================================================================
// Performance Comment DTOs
// =========================================================================

type CreatePerformanceCommentRequest struct {
	EvaluationID string `json:"evaluation_id" binding:"required"`
	EmployeeID   string `json:"employee_id" binding:"required"`
	Comment      string `json:"comment" binding:"required"`
	CreatedBy    string `json:"created_by" binding:"required"`
}

type UpdatePerformanceCommentRequest struct {
	Comment *string `json:"comment"`
}

type PerformanceCommentResponse struct {
	ID           string    `json:"id"`
	EvaluationID string    `json:"evaluation_id"`
	EmployeeID   string    `json:"employee_id"`
	Comment      string    `json:"comment"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// =========================================================================
// Performance Attachment DTOs
// =========================================================================

type CreatePerformanceAttachmentRequest struct {
	EvaluationDetailID string  `json:"evaluation_detail_id" binding:"required"`
	FilePath           string  `json:"file_path" binding:"required,max=500"`
	FileName           string  `json:"file_name" binding:"required,max=255"`
	FileType           *string `json:"file_type" binding:"omitempty,max=100"`
	FileSize           *int64  `json:"file_size"`
	Description        *string `json:"description"`
	UploadedBy         string  `json:"uploaded_by" binding:"required"`
}

type UpdatePerformanceAttachmentRequest struct {
	Description *string `json:"description"`
}

type PerformanceAttachmentResponse struct {
	ID                 string    `json:"id"`
	EvaluationDetailID string    `json:"evaluation_detail_id"`
	FilePath           string    `json:"file_path"`
	FileName           string    `json:"file_name"`
	FileType           string    `json:"file_type,omitempty"`
	FileSize           int64     `json:"file_size,omitempty"`
	Description        string    `json:"description,omitempty"`
	UploadedBy         string    `json:"uploaded_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// =========================================================================
// Performance Rating DTOs
// =========================================================================

type CreatePerformanceRatingRequest struct {
	Code        string  `json:"code" binding:"required,max=20"`
	Name        string  `json:"name" binding:"required,max=100"`
	MinScore    float64 `json:"min_score" binding:"required"`
	MaxScore    float64 `json:"max_score" binding:"required"`
	Color       *string `json:"color" binding:"omitempty,max=20"`
	Description *string `json:"description"`
	SortOrder   *int    `json:"sort_order"`
}

type UpdatePerformanceRatingRequest struct {
	Code        *string  `json:"code" binding:"omitempty,max=20"`
	Name        *string  `json:"name" binding:"omitempty,max=100"`
	MinScore    *float64 `json:"min_score"`
	MaxScore    *float64 `json:"max_score"`
	Color       *string  `json:"color" binding:"omitempty,max=20"`
	Description *string  `json:"description"`
	SortOrder   *int     `json:"sort_order"`
}

type PerformanceRatingResponse struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	MinScore    float64   `json:"min_score"`
	MaxScore    float64   `json:"max_score"`
	Color       string    `json:"color,omitempty"`
	Description string    `json:"description,omitempty"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// =========================================================================
// Performance Indicator Formula DTOs
// =========================================================================

type CreatePerformanceIndicatorFormulaRequest struct {
	Code        string  `json:"code" binding:"required,max=50"`
	Name        string  `json:"name" binding:"required,max=100"`
	FormulaType string  `json:"formula_type" binding:"required,max=30"`
	Expression  *string `json:"expression"`
	Description *string `json:"description"`
	SortOrder   *int    `json:"sort_order"`
}

type UpdatePerformanceIndicatorFormulaRequest struct {
	Code        *string `json:"code" binding:"omitempty,max=50"`
	Name        *string `json:"name" binding:"omitempty,max=100"`
	FormulaType *string `json:"formula_type" binding:"omitempty,max=30"`
	Expression  *string `json:"expression"`
	Description *string `json:"description"`
	SortOrder   *int    `json:"sort_order"`
}

type PerformanceIndicatorFormulaResponse struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	FormulaType string    `json:"formula_type"`
	Expression  string    `json:"expression,omitempty"`
	Description string    `json:"description,omitempty"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// =========================================================================
// Performance Log DTOs (Read-only)
// =========================================================================

type PerformanceLogResponse struct {
	ID           string    `json:"id"`
	EvaluationID string    `json:"evaluation_id,omitempty"`
	EntityType   string    `json:"entity_type"`
	EntityID     string    `json:"entity_id"`
	Action       string    `json:"action"`
	OldValues    string    `json:"old_values,omitempty"`
	NewValues    string    `json:"new_values,omitempty"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

// =========================================================================
// Phase 3 - Business Process DTOs
// =========================================================================

// CreateEvaluationWithSnapshotRequest creates evaluation and snapshots KPIs from template
type CreateEvaluationWithSnapshotRequest struct {
	EmployeeID     string  `json:"employee_id" binding:"required"`
	OrganizationID string  `json:"organization_id" binding:"required"`
	PeriodID       string  `json:"period_id" binding:"required"`
	TemplateID     string  `json:"template_id" binding:"required"`
	SupervisorID   *string `json:"supervisor_id"`
	Notes          *string `json:"notes"`
}

// UpdateEvaluationActualRequest updates actual values for evaluation detail
type UpdateEvaluationActualRequest struct {
	Actual  float64 `json:"actual" binding:"required"`
	Remarks *string `json:"remarks"`
}

// BulkUpdateEvaluationActualRequest updates multiple evaluation details
type BulkUpdateEvaluationActualRequest struct {
	Details []BulkUpdateDetailItem `json:"details" binding:"required,dive"`
}

type BulkUpdateDetailItem struct {
	DetailID string  `json:"detail_id" binding:"required"`
	Actual   float64 `json:"actual"`
	Remarks  *string `json:"remarks"`
}

// EvaluationWithDetailsResponse returns evaluation with all details
type EvaluationWithDetailsResponse struct {
	ID               string                     `json:"id"`
	EmployeeID       string                     `json:"employee_id"`
	EmployeeName     string                     `json:"employee_name,omitempty"`
	OrganizationID   string                     `json:"organization_id"`
	OrganizationName string                     `json:"organization_name,omitempty"`
	PeriodID         string                     `json:"period_id"`
	PeriodCode       string                     `json:"period_code,omitempty"`
	TemplateID       string                     `json:"template_id"`
	SupervisorID     string                     `json:"supervisor_id,omitempty"`
	FinalScore       float64                    `json:"final_score"`
	RatingID         string                     `json:"rating_id,omitempty"`
	RatingName       string                     `json:"rating_name,omitempty"`
	RatingColor      string                     `json:"rating_color,omitempty"`
	Status           string                     `json:"status"`
	SubmittedAt      string                     `json:"submitted_at,omitempty"`
	ApprovedAt       string                     `json:"approved_at,omitempty"`
	Notes            string                     `json:"notes,omitempty"`
	CreatedAt        time.Time                  `json:"created_at"`
	UpdatedAt        time.Time                  `json:"updated_at"`
	Details          []EvaluationDetailResponse `json:"details"`
}

// ProgressSummaryResponse returns progress summary for an evaluation
type ProgressSummaryResponse struct {
	EvaluationID       string  `json:"evaluation_id"`
	TotalIndicators    int     `json:"total_indicators"`
	CompletedCount     int     `json:"completed_count"`
	InProgressCount    int     `json:"in_progress_count"`
	NotStartedCount    int     `json:"not_started_count"`
	OverallProgress    float64 `json:"overall_progress"`
	AverageAchievement float64 `json:"average_achievement"`
}

// CalculateScoreRequest for recalculating scores
type CalculateScoreRequest struct {
	EvaluationID string `json:"evaluation_id" binding:"required"`
}

// =========================================================================
// Phase 4 - Dashboard DTOs
// =========================================================================

// EmployeeDashboardResponse - Dashboard untuk employee melihat KPI sendiri
type EmployeeDashboardResponse struct {
	EmployeeID       string             `json:"employee_id"`
	EmployeeName     string             `json:"employee_name,omitempty"`
	OrganizationID   string             `json:"organization_id"`
	OrganizationName string             `json:"organization_name,omitempty"`
	CurrentPeriod    *PeriodSummary     `json:"current_period,omitempty"`
	Evaluation       *EvaluationSummary `json:"evaluation,omitempty"`
	KPIProgress      []KPIProgressItem  `json:"kpi_progress"`
	RecentActivities []ActivityItem     `json:"recent_activities"`
}

type PeriodSummary struct {
	ID         string `json:"id"`
	PeriodCode string `json:"period_code"`
	PeriodType string `json:"period_type"`
	Year       int    `json:"year"`
	StartDate  string `json:"start_date,omitempty"`
	EndDate    string `json:"end_date,omitempty"`
	Status     string `json:"status"`
}

type EvaluationSummary struct {
	ID              string  `json:"id"`
	Status          string  `json:"status"`
	FinalScore      float64 `json:"final_score"`
	RatingName      string  `json:"rating_name,omitempty"`
	RatingColor     string  `json:"rating_color,omitempty"`
	TotalIndicators int     `json:"total_indicators"`
	CompletedCount  int     `json:"completed_count"`
	OverallProgress float64 `json:"overall_progress"`
	SubmittedAt     string  `json:"submitted_at,omitempty"`
	ApprovedAt      string  `json:"approved_at,omitempty"`
}

type KPIProgressItem struct {
	DetailID      string  `json:"detail_id"`
	IndicatorName string  `json:"indicator_name"`
	PerspectiveID string  `json:"perspective_id,omitempty"`
	Weight        float64 `json:"weight"`
	Target        float64 `json:"target"`
	Actual        float64 `json:"actual"`
	Achievement   float64 `json:"achievement"`
	Score         float64 `json:"score"`
	Status        string  `json:"status"` // NOT_STARTED, IN_PROGRESS, COMPLETED
}

type ActivityItem struct {
	ID          string `json:"id"`
	Action      string `json:"action"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	CreatedBy   string `json:"created_by,omitempty"`
}

// ManagerDashboardResponse - Dashboard untuk manager melihat tim
type ManagerDashboardResponse struct {
	ManagerID        string              `json:"manager_id"`
	ManagerName      string              `json:"manager_name,omitempty"`
	OrganizationID   string              `json:"organization_id"`
	OrganizationName string              `json:"organization_name,omitempty"`
	CurrentPeriod    *PeriodSummary      `json:"current_period,omitempty"`
	TeamSummary      TeamKPISummary      `json:"team_summary"`
	TeamMembers      []TeamMemberKPI     `json:"team_members"`
	PendingReviews   []PendingReviewItem `json:"pending_reviews"`
	OverdueReviews   []OverdueReviewItem `json:"overdue_reviews"`
}

type TeamKPISummary struct {
	TotalMembers    int     `json:"total_members"`
	CompletedCount  int     `json:"completed_count"`
	InProgressCount int     `json:"in_progress_count"`
	PendingCount    int     `json:"pending_count"`
	AverageScore    float64 `json:"average_score"`
	OverallProgress float64 `json:"overall_progress"`
}

type TeamMemberKPI struct {
	EmployeeID   string  `json:"employee_id"`
	EmployeeName string  `json:"employee_name,omitempty"`
	EvaluationID string  `json:"evaluation_id,omitempty"`
	Status       string  `json:"status"`
	FinalScore   float64 `json:"final_score"`
	RatingName   string  `json:"rating_name,omitempty"`
	RatingColor  string  `json:"rating_color,omitempty"`
	Progress     float64 `json:"progress"`
}

type PendingReviewItem struct {
	EvaluationID string `json:"evaluation_id"`
	EmployeeID   string `json:"employee_id"`
	EmployeeName string `json:"employee_name,omitempty"`
	SubmittedAt  string `json:"submitted_at"`
	DaysPending  int    `json:"days_pending"`
}

type OverdueReviewItem struct {
	EvaluationID string `json:"evaluation_id"`
	EmployeeID   string `json:"employee_id"`
	EmployeeName string `json:"employee_name,omitempty"`
	Status       string `json:"status"`
	DaysOverdue  int    `json:"days_overdue"`
}

// HRDashboardResponse - Dashboard untuk HR melihat keseluruhan
type HRDashboardResponse struct {
	CurrentPeriod      *PeriodSummary           `json:"current_period,omitempty"`
	CompletionStats    CompletionStats          `json:"completion_stats"`
	RatingDistribution []RatingDistributionItem `json:"rating_distribution"`
	OrganizationStats  []OrganizationStatsItem  `json:"organization_stats"`
	TopPerformers      []PerformerItem          `json:"top_performers"`
	BottomPerformers   []PerformerItem          `json:"bottom_performers"`
	TrendData          []TrendItem              `json:"trend_data"`
}

type CompletionStats struct {
	TotalEmployees     int     `json:"total_employees"`
	CompletedCount     int     `json:"completed_count"`
	ApprovedCount      int     `json:"approved_count"`
	SubmittedCount     int     `json:"submitted_count"`
	DraftCount         int     `json:"draft_count"`
	NotStartedCount    int     `json:"not_started_count"`
	CompletionRate     float64 `json:"completion_rate"`
	AverageScore       float64 `json:"average_score"`
	AverageAchievement float64 `json:"average_achievement"`
}

type RatingDistributionItem struct {
	RatingID    string  `json:"rating_id"`
	RatingCode  string  `json:"rating_code"`
	RatingName  string  `json:"rating_name"`
	RatingColor string  `json:"rating_color,omitempty"`
	Count       int     `json:"count"`
	Percentage  float64 `json:"percentage"`
}

type OrganizationStatsItem struct {
	OrganizationID     string  `json:"organization_id"`
	OrganizationName   string  `json:"organization_name,omitempty"`
	TotalEmployees     int     `json:"total_employees"`
	CompletedCount     int     `json:"completed_count"`
	AverageScore       float64 `json:"average_score"`
	AverageAchievement float64 `json:"average_achievement"`
	CompletionRate     float64 `json:"completion_rate"`
}

type PerformerItem struct {
	Rank             int     `json:"rank"`
	EmployeeID       string  `json:"employee_id"`
	EmployeeName     string  `json:"employee_name,omitempty"`
	OrganizationID   string  `json:"organization_id"`
	OrganizationName string  `json:"organization_name,omitempty"`
	FinalScore       float64 `json:"final_score"`
	RatingName       string  `json:"rating_name,omitempty"`
	RatingColor      string  `json:"rating_color,omitempty"`
}

type TrendItem struct {
	PeriodID       string  `json:"period_id"`
	PeriodCode     string  `json:"period_code"`
	Year           int     `json:"year"`
	AverageScore   float64 `json:"average_score"`
	CompletionRate float64 `json:"completion_rate"`
	TotalEmployees int     `json:"total_employees"`
}

// =========================================================================
// Performance Component DTOs (Phase 5 - Scoring Configuration)
// =========================================================================

type CreatePerformanceComponentRequest struct {
	Code        string  `json:"code" binding:"required,max=50"`
	Name        string  `json:"name" binding:"required,max=100"`
	Description *string `json:"description"`
	SortOrder   *int    `json:"sort_order"`
	IsActive    *bool   `json:"is_active"`
}

type UpdatePerformanceComponentRequest struct {
	Code        *string `json:"code" binding:"omitempty,max=50"`
	Name        *string `json:"name" binding:"omitempty,max=100"`
	Description *string `json:"description"`
	SortOrder   *int    `json:"sort_order"`
	IsActive    *bool   `json:"is_active"`
}

type PerformanceComponentResponse struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// =========================================================================
// Performance Organization Component DTOs
// =========================================================================

type UpsertOrganizationComponentRequest struct {
	OrganizationID string  `json:"organization_id" binding:"required"`
	ComponentID    string  `json:"component_id" binding:"required"`
	Weight         float64 `json:"weight"`
	IsEnabled      *bool   `json:"is_enabled"`
	SortOrder      *int    `json:"sort_order"`
}

type OrganizationComponentResponse struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	ComponentID    string    `json:"component_id"`
	ComponentCode  string    `json:"component_code,omitempty"`
	ComponentName  string    `json:"component_name,omitempty"`
	Weight         float64   `json:"weight"`
	IsEnabled      bool      `json:"is_enabled"`
	SortOrder      int       `json:"sort_order"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// =========================================================================
// Performance Evaluation Component DTOs
// =========================================================================

type PerformanceEvaluationComponentResponse struct {
	ID            string     `json:"id"`
	EvaluationID  string     `json:"evaluation_id"`
	ComponentID   string     `json:"component_id"`
	ComponentName string     `json:"component_name"`
	Score         float64    `json:"score"`
	Weight        float64    `json:"weight"`
	FinalScore    float64    `json:"final_score"`
	CalculatedAt  *time.Time `json:"calculated_at,omitempty"`
}

// UpdateEvaluationComponentScoreRequest digunakan untuk mengisi skor komponen
// yang tidak bisa dihitung otomatis (mis. Work Program — tidak punya sumber
// data lain di sistem, wajib diisi manual oleh reviewer).
type UpdateEvaluationComponentScoreRequest struct {
	Score float64 `json:"score" binding:"required,min=0,max=100"`
}
