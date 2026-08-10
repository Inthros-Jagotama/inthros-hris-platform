package careerintelligence

import "time"

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
// Talent Map DTOs
// =========================================================================

type CreateTalentMapRequest struct {
	EmployeeID  string `json:"employee_id" binding:"required"`
	Period      string `json:"period" binding:"required,max=7"`
	Performance string `json:"performance" binding:"required,oneof=LOW MEDIUM HIGH"`
	Potential   string `json:"potential" binding:"required,oneof=LOW MEDIUM HIGH"`
	Notes       string `json:"notes"`
}

type UpdateTalentMapRequest struct {
	Performance *string `json:"performance" binding:"omitempty,oneof=LOW MEDIUM HIGH"`
	Potential   *string `json:"potential" binding:"omitempty,oneof=LOW MEDIUM HIGH"`
	Notes       *string `json:"notes"`
}

type TalentMapResponse struct {
	ID           string `json:"id"`
	EmployeeID   string `json:"employee_id"`
	Period       string `json:"period"`
	Performance  string `json:"performance"`
	Potential    string `json:"potential"`
	GridPosition string `json:"grid_position"`
	Notes        string `json:"notes"`
	AssessorID   string `json:"assessor_id,omitempty"`
	AssessedAt   string `json:"assessed_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TalentGridResponse struct {
	Period    string            `json:"period"`
	Quadrants []TalentQuadrant  `json:"quadrants"`
	Total     int               `json:"total"`
}

type TalentQuadrant struct {
	Label       string `json:"label"`       // e.g. "High Performer - High Potential"
	Position    string `json:"position"`    // e.g. "9-BOX-1"
	Count       int    `json:"count"`
	Description string `json:"description"`
}

type EmployeeTalentProfileResponse struct {
	EmployeeID   string                `json:"employee_id"`
	CurrentMap   *TalentMapResponse    `json:"current_map,omitempty"`
	History      []TalentMapResponse   `json:"history,omitempty"`
	Interests    []CareerInterestResponse `json:"interests,omitempty"`
	ReadyFor     []string              `json:"ready_for,omitempty"` // recommended next positions
}

// =========================================================================
// Career Interest DTOs
// =========================================================================

type CreateCareerInterestRequest struct {
	EmployeeID      string `json:"employee_id" binding:"required"`
	InterestType    string `json:"interest_type" binding:"required,oneof=LEADERSHIP SPECIALIST INTERNATIONAL ENTREPRENEUR"`
	TargetPosition  string `json:"target_position"`
	TargetDepartment string `json:"target_department"`
	Motivation      string `json:"motivation"`
	ReadinessLevel  string `json:"readiness_level" binding:"omitempty,oneof=NOW 1_YEAR 2_3_YEARS 3_PLUS"`
}

type CareerInterestResponse struct {
	ID               string    `json:"id"`
	EmployeeID       string    `json:"employee_id"`
	InterestType     string    `json:"interest_type"`
	TargetPosition   string    `json:"target_position,omitempty"`
	TargetDepartment string    `json:"target_department,omitempty"`
	Motivation       string    `json:"motivation,omitempty"`
	ReadinessLevel   string    `json:"readiness_level,omitempty"`
	IsActive         bool      `json:"is_active"`
	RecordedAt       string    `json:"recorded_at"`
	CreatedAt        time.Time `json:"created_at"`
}

// =========================================================================
// Career Path DTOs
// =========================================================================

// CreateCareerPathRequest — CI membuat EDGE karier (source → target) yang pada
// skema terpadu disimpan sebagai path 2-langkah. Name opsional: bila kosong,
// service meng-generate "<PATH_TYPE>: <source> → <target>" (menghindari tabrakan
// uk_career_paths_name dengan menambahkan akhiran unik bila perlu).
type CreateCareerPathRequest struct {
	Name          string `json:"name"`
	SourceTitleID string `json:"source_title_id" binding:"required"`
	TargetTitleID string `json:"target_title_id" binding:"required"`
	PathType      string `json:"path_type" binding:"required,oneof=PROMOTION LATERAL DEMOTION CROSSFUNCTIONAL"`
	TypicalTenure int    `json:"typical_tenure"`
	Requirements  string `json:"requirements"`
	Competencies  string `json:"competencies"`
}

type CareerPathStepResponse struct {
	ID             string `json:"id"`
	PositionID     string `json:"position_id"`
	PositionName   string `json:"position_name,omitempty"`
	Sequence       int    `json:"sequence"`
	PathType       string `json:"path_type,omitempty"`
	TypicalTenure  *int   `json:"typical_tenure,omitempty"`
	Competencies   string `json:"competencies,omitempty"`
	Certifications string `json:"certifications,omitempty"`
}

type CareerPathResponse struct {
	ID            string                  `json:"id"`
	Name          string                  `json:"name"`
	SourceTitleID string                  `json:"source_title_id"`
	TargetTitleID string                  `json:"target_title_id"`
	PathType      string                  `json:"path_type"`
	TypicalTenure int                     `json:"typical_tenure"`
	Requirements  string                  `json:"requirements,omitempty"`
	Competencies  string                  `json:"competencies,omitempty"`
	Certifications string                 `json:"certifications,omitempty"`
	IsActive      bool                    `json:"is_active"`
	Steps         []CareerPathStepResponse `json:"steps,omitempty"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}

// =========================================================================
// Gap Analysis DTOs (Career Path)
// =========================================================================

type GapAnalysisRequest struct {
	EmployeeID   string `json:"employee_id" form:"employee_id" binding:"required"`
	TargetTitleID string `json:"target_title_id" form:"target_title_id" binding:"required"`
}

type GapAnalysisResponse struct {
	EmployeeID    string              `json:"employee_id"`
	TargetTitle   string              `json:"target_title"`
	MatchedSkills int                 `json:"matched_skills"`
	TotalRequired int                 `json:"total_required"`
	GapPercentage float64             `json:"gap_percentage"`
	Recommendations []GapRecommendation `json:"recommendations,omitempty"`
	EstimatedTimeline string          `json:"estimated_timeline"` // e.g. "12-18 months"
}

type GapRecommendation struct {
	Category    string `json:"category"`    // TRAINING / EXPERIENCE / CERTIFICATION
	Description string `json:"description"`
	Priority    string `json:"priority"`     // HIGH / MEDIUM / LOW
}

// =========================================================================
// Succession Plan DTOs
// =========================================================================

type CreateSuccessionPlanRequest struct {
	PositionID      string `json:"position_id" binding:"required"`
	SuccessorID     string `json:"successor_id" binding:"required"`
	ReadinessLevel  string `json:"readiness_level" binding:"required,oneof=READY_NOW READY_1YR READY_2YR POTENTIAL"`
	PriorityOrder   *int   `json:"priority_order"`
	TargetDate      string `json:"target_date"`
	DevelopmentPlan string `json:"development_plan"`
	Notes           string `json:"notes"`
}

type UpdateSuccessionPlanRequest struct {
	ReadinessLevel  *string `json:"readiness_level" binding:"omitempty,oneof=READY_NOW READY_1YR READY_2YR POTENTIAL"`
	PriorityOrder   *int    `json:"priority_order"`
	TargetDate      *string `json:"target_date"`
	DevelopmentPlan *string `json:"development_plan"`
	Notes           *string `json:"notes"`
}

type SuccessionPlanResponse struct {
	ID              string    `json:"id"`
	PositionID      string    `json:"position_id"`
	SuccessorID     string    `json:"successor_id"`
	ReadinessLevel  string    `json:"readiness_level"`
	PriorityOrder   int       `json:"priority_order"`
	TargetDate      *string   `json:"target_date,omitempty"`
	DevelopmentPlan string    `json:"development_plan,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
