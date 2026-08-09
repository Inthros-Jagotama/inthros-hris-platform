package performance

import "time"

// =========================================================================
// OKR Template DTOs
// =========================================================================

type CreateOKRTemplateRequest struct {
	OrganizationID string  `json:"organization_id" binding:"required"`
	PeriodID       *string `json:"period_id"`
	Name           string  `json:"name" binding:"required,max=255"`
	Description    *string `json:"description"`
	Status         *int    `json:"status"`
	EffectiveDate  *string `json:"effective_date"`
	ExpiredDate    *string `json:"expired_date"`
}

type UpdateOKRTemplateRequest struct {
	PeriodID      *string `json:"period_id"`
	Name          *string `json:"name" binding:"omitempty,max=255"`
	Description   *string `json:"description"`
	Status        *int    `json:"status"`
	EffectiveDate *string `json:"effective_date"`
	ExpiredDate   *string `json:"expired_date"`
}

type OKRTemplateResponse struct {
	ID               string                 `json:"id"`
	OrganizationID   string                 `json:"organization_id"`
	OrganizationName string                 `json:"organization_name,omitempty"`
	PeriodID         string                 `json:"period_id,omitempty"`
	PeriodCode       string                 `json:"period_code,omitempty"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description,omitempty"`
	Status           int                    `json:"status"`
	EffectiveDate    string                 `json:"effective_date,omitempty"`
	ExpiredDate      string                 `json:"expired_date,omitempty"`
	ObjectiveCount   int                    `json:"objective_count,omitempty"`
	CreatedBy        string                 `json:"created_by"`
	CreatedByOrgID   string                 `json:"created_by_org_id,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	Objectives       []OKRObjectiveResponse `json:"objectives,omitempty"`
}

// OKRObjectiveScopeResponse — cascading objective creation scope: whether
// the calling employee is currently allowed to create an Objective (an
// OKRTemplate) for a subordinate Organization, and which Organizations
// qualify as their effective subordinates.
type OKRObjectiveScopeResponse struct {
	OrganizationID           string                       `json:"organization_id,omitempty"`
	OrganizationName         string                       `json:"organization_name,omitempty"`
	Eligible                 bool                         `json:"eligible"`
	IneligibleReasonKey      string                       `json:"ineligible_reason_key,omitempty"`
	SubordinateOrganizations []OrganizationOptionResponse `json:"subordinate_organizations"`
}

// MyOKRContextResponse — self-assessment context: resolves the current
// employee's Organization and the PUBLISHED (status=1) OKR templates
// configured for it, mirroring MyKPIContextResponse for the KPI module.
type MyOKRContextResponse struct {
	HasPosition      bool                  `json:"has_position"`
	EmployeeID       string                `json:"employee_id,omitempty"`
	OrganizationID   string                `json:"organization_id,omitempty"`
	OrganizationName string                `json:"organization_name,omitempty"`
	Templates        []OKRTemplateResponse `json:"templates"`
}

// =========================================================================
// OKR Objective DTOs
// =========================================================================

type CreateOKRObjectiveRequest struct {
	TemplateID  string  `json:"template_id" binding:"required"`
	Code        *string `json:"code" binding:"omitempty,max=50"`
	Title       string  `json:"title" binding:"required,max=255"`
	Description *string `json:"description"`
	Weight      float64 `json:"weight"`
	SortOrder   *int    `json:"sort_order"`
}

type UpdateOKRObjectiveRequest struct {
	Code        *string  `json:"code" binding:"omitempty,max=50"`
	Title       *string  `json:"title" binding:"omitempty,max=255"`
	Description *string  `json:"description"`
	Weight      *float64 `json:"weight"`
	SortOrder   *int     `json:"sort_order"`
}

type OKRObjectiveResponse struct {
	ID          string                 `json:"id"`
	TemplateID  string                 `json:"template_id"`
	Code        string                 `json:"code,omitempty"`
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	Weight      float64                `json:"weight"`
	SortOrder   int                    `json:"sort_order"`
	KeyResults  []OKRKeyResultResponse `json:"key_results,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// =========================================================================
// OKR Key Result DTOs
// =========================================================================

type CreateOKRKeyResultRequest struct {
	ObjectiveID  string   `json:"objective_id" binding:"required"`
	Code         *string  `json:"code" binding:"omitempty,max=50"`
	Title        string   `json:"title" binding:"required,max=255"`
	Description  *string  `json:"description"`
	TargetType   *string  `json:"target_type" binding:"omitempty,oneof=NUMBER CURRENCY PERCENTAGE DURATION BOOLEAN"`
	TargetValue  float64  `json:"target_value"`
	Unit         *string  `json:"unit" binding:"omitempty,max=50"`
	FormulaType  *string  `json:"formula_type" binding:"omitempty,oneof=MANUAL HIGHER_BETTER LOWER_BETTER RANGE BOOLEAN PERCENTAGE"`
	Weight       float64  `json:"weight"`
	MinimumScore *float64 `json:"minimum_score"`
	MaximumScore *float64 `json:"maximum_score"`
	SortOrder    *int     `json:"sort_order"`
	IsRequired   *bool    `json:"is_required"`
}

type UpdateOKRKeyResultRequest struct {
	Code         *string  `json:"code" binding:"omitempty,max=50"`
	Title        *string  `json:"title" binding:"omitempty,max=255"`
	Description  *string  `json:"description"`
	TargetType   *string  `json:"target_type" binding:"omitempty,oneof=NUMBER CURRENCY PERCENTAGE DURATION BOOLEAN"`
	TargetValue  *float64 `json:"target_value"`
	Unit         *string  `json:"unit" binding:"omitempty,max=50"`
	FormulaType  *string  `json:"formula_type" binding:"omitempty,oneof=MANUAL HIGHER_BETTER LOWER_BETTER RANGE BOOLEAN PERCENTAGE"`
	Weight       *float64 `json:"weight"`
	MinimumScore *float64 `json:"minimum_score"`
	MaximumScore *float64 `json:"maximum_score"`
	SortOrder    *int     `json:"sort_order"`
	IsRequired   *bool    `json:"is_required"`
}

type OKRKeyResultResponse struct {
	ID           string    `json:"id"`
	ObjectiveID  string    `json:"objective_id"`
	Code         string    `json:"code,omitempty"`
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	TargetType   string    `json:"target_type"`
	TargetValue  float64   `json:"target_value"`
	Unit         string    `json:"unit,omitempty"`
	FormulaType  string    `json:"formula_type"`
	Weight       float64   `json:"weight"`
	MinimumScore float64   `json:"minimum_score"`
	MaximumScore float64   `json:"maximum_score"`
	SortOrder    int       `json:"sort_order"`
	IsRequired   bool      `json:"is_required"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// =========================================================================
// OKR Evaluation DTOs
// =========================================================================

type CreateOKREvaluationRequest struct {
	EmployeeID     string `json:"employee_id" binding:"required"`
	OrganizationID string `json:"organization_id" binding:"required"`
	PeriodID       string `json:"period_id" binding:"required"`
	TemplateID     string `json:"template_id" binding:"required"`
}

type UpdateOKREvaluationRequest struct {
	Status        *string `json:"status" binding:"omitempty,oneof=DRAFT SUBMITTED APPROVED COMPLETED REJECTED"`
	ReviewerNotes *string `json:"reviewer_notes"`
}

type OKREvaluationResponse struct {
	ID                           string                        `json:"id"`
	EmployeeID                   string                        `json:"employee_id"`
	EmployeeName                 string                        `json:"employee_name,omitempty"`
	OrganizationID               string                        `json:"organization_id"`
	OrganizationName             string                        `json:"organization_name,omitempty"`
	PeriodID                     string                        `json:"period_id"`
	PeriodCode                   string                        `json:"period_code,omitempty"`
	TemplateID                   string                        `json:"template_id,omitempty"`
	TemplateName                 string                        `json:"template_name,omitempty"`
	Status                       string                        `json:"status"`
	SubmittedAt                  *time.Time                    `json:"submitted_at,omitempty"`
	ApprovedAt                   *time.Time                    `json:"approved_at,omitempty"`
	KRApprovalInstanceID         string                        `json:"kr_approval_instance_id,omitempty"`
	AssessmentApprovalInstanceID string                        `json:"assessment_approval_instance_id,omitempty"`
	FinalScore                   float64                       `json:"final_score"`
	RatingID                     string                        `json:"rating_id,omitempty"`
	RatingName                   string                        `json:"rating_name,omitempty"`
	RatingColor                  string                        `json:"rating_color,omitempty"`
	ReviewerNotes                string                        `json:"reviewer_notes,omitempty"`
	Details                      []OKREvaluationDetailResponse `json:"details,omitempty"`
	CreatedAt                    time.Time                     `json:"created_at"`
	UpdatedAt                    time.Time                     `json:"updated_at"`
}

// =========================================================================
// OKR Evaluation Detail DTOs
// =========================================================================

type UpdateOKREvaluationDetailRequest struct {
	ActualValue float64 `json:"actual_value"`
	Remarks     *string `json:"remarks"`
}

// CreateOKREvaluationKeyResultRequest — employee-proposed Key Result under a
// snapshotted Objective, created while the evaluation is DRAFT. Mirrors
// CreateProgramItemRequest for KPI.
type CreateOKREvaluationKeyResultRequest struct {
	EvaluationID    string  `json:"evaluation_id" binding:"required"`
	ObjectiveID     string  `json:"objective_id" binding:"required"`
	ObjectiveTitle  string  `json:"objective_title" binding:"required,max=255"`
	ObjectiveWeight float64 `json:"objective_weight"`
	Title           string  `json:"title" binding:"required,max=255"`
	TargetType      string  `json:"target_type" binding:"omitempty,oneof=NUMBER PERCENTAGE CURRENCY BOOLEAN DURATION SCORE"`
	TargetValue     float64 `json:"target_value"`
	Unit            *string `json:"unit" binding:"omitempty,max=50"`
	FormulaType     string  `json:"formula_type" binding:"omitempty,oneof=MANUAL HIGHER_BETTER LOWER_BETTER RANGE BOOLEAN PERCENTAGE"`
	Weight          float64 `json:"weight"`
}

type UpdateOKREvaluationKeyResultTargetRequest struct {
	Title       *string  `json:"title" binding:"omitempty,max=255"`
	TargetType  *string  `json:"target_type" binding:"omitempty,oneof=NUMBER PERCENTAGE CURRENCY BOOLEAN DURATION SCORE"`
	TargetValue *float64 `json:"target_value"`
	Unit        *string  `json:"unit" binding:"omitempty,max=50"`
	FormulaType *string  `json:"formula_type" binding:"omitempty,oneof=MANUAL HIGHER_BETTER LOWER_BETTER RANGE BOOLEAN PERCENTAGE"`
	Weight      *float64 `json:"weight"`
}

type OKRBulkUpdateActualsRequest struct {
	Details []OKRBulkActualItem `json:"details" binding:"required,dive"`
}

type OKRBulkActualItem struct {
	ID          string  `json:"id" binding:"required"`
	ActualValue float64 `json:"actual_value"`
}

type OKREvaluationDetailResponse struct {
	ID              string    `json:"id"`
	EvaluationID    string    `json:"evaluation_id"`
	ObjectiveID     string    `json:"objective_id,omitempty"`
	KeyResultID     string    `json:"key_result_id,omitempty"`
	ObjectiveTitle  string    `json:"objective_title"`
	KeyResultTitle  string    `json:"key_result_title"`
	ObjectiveWeight float64   `json:"objective_weight"`
	KeyResultWeight float64   `json:"key_result_weight"`
	TargetValue     float64   `json:"target_value"`
	TargetType      string    `json:"target_type"`
	Unit            string    `json:"unit,omitempty"`
	FormulaType     string    `json:"formula_type"`
	ActualValue     float64   `json:"actual_value"`
	Achievement     float64   `json:"achievement"`
	Score           float64   `json:"score"`
	Remarks         string    `json:"remarks,omitempty"`
	SortOrder       int       `json:"sort_order"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// =========================================================================
// OKR Progress DTOs
// =========================================================================

type CreateOKRProgressRequest struct {
	EvaluationDetailID string  `json:"evaluation_detail_id" binding:"required"`
	ProgressDate       string  `json:"progress_date" binding:"required"`
	ActualValue        float64 `json:"actual_value"`
	Notes              *string `json:"notes"`
}

type UpdateOKRProgressRequest struct {
	ProgressDate *string  `json:"progress_date"`
	ActualValue  *float64 `json:"actual_value"`
	Notes        *string  `json:"notes"`
}

type OKRProgressResponse struct {
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
// OKR Comment DTOs
// =========================================================================

type CreateOKRCommentRequest struct {
	EvaluationID string  `json:"evaluation_id" binding:"required"`
	ParentID     *string `json:"parent_id"`
	Comment      string  `json:"comment" binding:"required"`
}

type UpdateOKRCommentRequest struct {
	Comment string `json:"comment" binding:"required"`
}

type OKRCommentResponse struct {
	ID            string               `json:"id"`
	EvaluationID  string               `json:"evaluation_id"`
	ParentID      string               `json:"parent_id,omitempty"`
	Comment       string               `json:"comment"`
	CreatedBy     string               `json:"created_by"`
	CreatedByName string               `json:"created_by_name,omitempty"`
	Replies       []OKRCommentResponse `json:"replies,omitempty"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

// =========================================================================
// OKR Attachment DTOs
// =========================================================================

type CreateOKRAttachmentRequest struct {
	EvaluationDetailID string  `json:"evaluation_detail_id" binding:"required"`
	FilePath           string  `json:"file_path" binding:"required"`
	FileName           string  `json:"file_name" binding:"required"`
	FileType           *string `json:"file_type"`
	FileSize           *int64  `json:"file_size"`
	Description        *string `json:"description"`
}

type OKRAttachmentResponse struct {
	ID                 string    `json:"id"`
	EvaluationDetailID string    `json:"evaluation_detail_id"`
	FilePath           string    `json:"file_path"`
	FileName           string    `json:"file_name"`
	FileType           string    `json:"file_type,omitempty"`
	FileSize           int64     `json:"file_size,omitempty"`
	Description        string    `json:"description,omitempty"`
	UploadedBy         string    `json:"uploaded_by"`
	UploadedByName     string    `json:"uploaded_by_name,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// =========================================================================
// OKR Dashboard DTOs
// =========================================================================

type OKRDashboardEmployeeResponse struct {
	TotalObjectives    int                       `json:"total_objectives"`
	TotalKeyResults    int                       `json:"total_key_results"`
	OverallAchievement float64                   `json:"overall_achievement"`
	FinalScore         float64                   `json:"final_score"`
	RatingName         string                    `json:"rating_name,omitempty"`
	Status             string                    `json:"status"`
	Objectives         []OKRObjectiveProgressDTO `json:"objectives,omitempty"`
}

type OKRObjectiveProgressDTO struct {
	ObjectiveTitle string                    `json:"objective_title"`
	Weight         float64                   `json:"weight"`
	Achievement    float64                   `json:"achievement"`
	Score          float64                   `json:"score"`
	KeyResults     []OKRKeyResultProgressDTO `json:"key_results,omitempty"`
}

type OKRKeyResultProgressDTO struct {
	KeyResultTitle string  `json:"key_result_title"`
	TargetValue    float64 `json:"target_value"`
	ActualValue    float64 `json:"actual_value"`
	Achievement    float64 `json:"achievement"`
	Score          float64 `json:"score"`
}

type OKRDashboardHRResponse struct {
	TotalEvaluations   int                 `json:"total_evaluations"`
	CompletedCount     int                 `json:"completed_count"`
	ApprovedCount      int                 `json:"approved_count"`
	SubmittedCount     int                 `json:"submitted_count"`
	DraftCount         int                 `json:"draft_count"`
	AverageScore       float64             `json:"average_score"`
	AverageAchievement float64             `json:"average_achievement"`
	RatingDistribution []OKRRatingCountDTO `json:"rating_distribution,omitempty"`
}

type OKRRatingCountDTO struct {
	RatingName string `json:"rating_name"`
	Count      int    `json:"count"`
	Color      string `json:"color,omitempty"`
}
