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

// GenerateTalentMapRequest membuat talent map otomatis: performance diambil
// dari skor final Performance Management, potential dari skor Competency
// terbaru milik employee, lalu dibanding LOW/MEDIUM/HIGH memakai
// TalentMapSettings tenant. Tidak ada input performance/potential manual.
type GenerateTalentMapRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	Period     string `json:"period" binding:"required,max=7"`
	Notes      string `json:"notes"`
}

type TalentMapSettingsResponse struct {
	PerformanceLowMax  float64 `json:"performance_low_max"`
	PerformanceHighMin float64 `json:"performance_high_min"`
	PotentialLowMax    float64 `json:"potential_low_max"`
	PotentialHighMin   float64 `json:"potential_high_min"`
}

type UpdateTalentMapSettingsRequest struct {
	PerformanceLowMax  float64 `json:"performance_low_max" binding:"required,gte=0,lte=100"`
	PerformanceHighMin float64 `json:"performance_high_min" binding:"required,gte=0,lte=100"`
	PotentialLowMax    float64 `json:"potential_low_max" binding:"required,gte=0,lte=100"`
	PotentialHighMin   float64 `json:"potential_high_min" binding:"required,gte=0,lte=100"`
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

type CareerPathStepResponse struct {
	ID                   string  `json:"id"`
	PositionID           string  `json:"position_id"`
	PositionName         string  `json:"position_name,omitempty"`
	Sequence             int     `json:"sequence"`
	MinimumServiceMonths *int    `json:"minimum_service_months,omitempty"`
	Requirements         string  `json:"requirements,omitempty"`
	PathType             string  `json:"path_type,omitempty"`
	TypicalTenure        *int    `json:"typical_tenure,omitempty"`
	Competencies         string  `json:"competencies,omitempty"`
	Certifications       string  `json:"certifications,omitempty"`
}

type CareerPathResponse struct {
	ID            string                  `json:"id"`
	Name          string                  `json:"name"`
	Description   string                  `json:"description,omitempty"`
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
// Internal Candidate Eligibility DTOs (S-4 — CI → Recruitment)
// =========================================================================

// EligibleEmployeeResponse adalah employee internal yang eligible untuk posisi
// lowongan: saat ini memegang salah satu source step dari career path menuju
// target position (S-4 — Position Vacancy → Career Path → Eligible Employees).
// CI menentukan eligibility; Recruitment hanya mengeksekusi aplikasi internal.
type EligibleEmployeeResponse struct {
	EmployeeID         string `json:"employee_id"`
	Name               string `json:"name"`
	CurrentPositionID  string `json:"current_position_id"`
	CurrentPositionName string `json:"current_position_name,omitempty"`
	SourceStepSequence int    `json:"source_step_sequence"`
	TargetPositionID   string `json:"target_position_id"`
	TargetPositionName string `json:"target_position_name,omitempty"`
	PathID             string `json:"path_id"`
	PathName           string `json:"path_name,omitempty"`
}

// =========================================================================
// Career Path Ladder DTOs (enhancement plan §12.9 — strategical planning)
//
// Setelah unifikasi kepemilikan career paths ke module Career Intelligence,
// endpoint /career-intelligence/paths menerima bentuk ladder (nama jenjang +
// langkah berurutan) via CreateCareerPathLadderRequest. Menulis ke skema
// terpadu career_paths + career_path_steps.
// =========================================================================

// CreateCareerPathStepRequest satu langkah pada jenjang. Sequence unik per
// path (divalidasi service), position_id harus merujuk posisi yang ada.
type CreateCareerPathStepRequest struct {
	PositionID           string `json:"position_id" binding:"required,uuid"`
	Sequence             int    `json:"sequence" binding:"required,gte=1"`
	MinimumServiceMonths *int   `json:"minimum_service_months" binding:"omitempty,gte=0"`
	Requirements         string `json:"requirements"`
}

// CreateCareerPathLadderRequest membuat jenjang penuh: nama + daftar langkah
// berurutan. Inilah bentuk yang dipakai FE Career Paths (strategical).
type CreateCareerPathLadderRequest struct {
	Name        string                        `json:"name" binding:"required"`
	Description *string                       `json:"description"`
	IsActive    *bool                         `json:"is_active"`
	Steps       []CreateCareerPathStepRequest `json:"steps" binding:"required,min=1,dive"`
}

// UpdateCareerPathRequest semantik full-replace: klien mengirim daftar steps
// lengkap yang diinginkan; server mengganti seluruh steps path.
type UpdateCareerPathRequest struct {
	Name        *string                       `json:"name"`
	Description *string                       `json:"description"`
	IsActive    *bool                         `json:"is_active"`
	Steps       []CreateCareerPathStepRequest `json:"steps" binding:"required,min=1,dive"`
}

// PaginatedCareerPathResponse dipakai endpoint list ladder. Reuse
// PaginatedResponse agar shape konsisten dengan endpoint CI lain.
type PaginatedCareerPathResponse = PaginatedResponse

// =========================================================================
// Gap Analysis DTOs (Career Path)
// =========================================================================

type GapAnalysisRequest struct {
	EmployeeID string `json:"employee_id" form:"employee_id" binding:"required"`
	// TargetTitleID is actually an organizations.id (Organization = Position
	// in this app's data model) -- resolved via Repository.GetPositionTitle.
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

// SuccessionGapResponse — posisi kunci beserta status succession gap-nya
// (S-5). HasReadySuccessor=false + RequiresExternalRecruitment=true berarti
// posisi tsb tidak punya successor siap (READY_NOW) sehingga fallback external
// recruitment diperlukan.
type SuccessionGapResponse struct {
	PositionID                 string `json:"position_id"`
	PositionTitle              string `json:"position_title"`
	PositionCode               string `json:"position_code,omitempty"`
	OrganizationID             string `json:"organization_id,omitempty"`
	SuccessorCount             int    `json:"successor_count"`
	ReadyNowCount              int    `json:"ready_now_count"`
	HasReadySuccessor          bool   `json:"has_ready_successor"`
	RequiresExternalRecruitment bool  `json:"requires_external_recruitment"`
}
