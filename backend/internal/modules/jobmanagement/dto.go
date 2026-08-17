package jobmanagement

import "time"

// =========================================================================
// Request DTOs — Job Titles
// =========================================================================

type CreateJobTitleRequest struct {
	Name         string `json:"name" binding:"required,max=100"`
	Descriptions string `json:"descriptions"`
	Status       int8   `json:"status"`
}

type UpdateJobTitleRequest struct {
	Name         *string `json:"name" binding:"omitempty,max=100"`
	Descriptions *string `json:"descriptions"`
	Status       *int8   `json:"status"`
}

// =========================================================================
// Request DTOs — Job Title Subs
// =========================================================================

type CreateJobTitleSubRequest struct {
	JobManagementTitleID string `json:"job_management_title_id" binding:"required"`
	Name                 string `json:"name" binding:"required,max=100"`
	Descriptions         string `json:"descriptions"`
	Status               int8   `json:"status"`
}

type UpdateJobTitleSubRequest struct {
	Name         *string `json:"name" binding:"omitempty,max=100"`
	Descriptions *string `json:"descriptions"`
	Status       *int8   `json:"status"`
}

// =========================================================================
// Request DTOs — Job Values
// =========================================================================

type CreateJobValueRequest struct {
	JobManagementTitleSubID *string `json:"job_management_title_sub_id"`
	Type                    string  `json:"type" binding:"required"`
	TypeGroup               *string `json:"type_group"`
	Level                   *int    `json:"level"`
	Descriptions            string  `json:"descriptions"`
	DescriptionGroup        *string `json:"description_group"`
	Note                    string  `json:"note"`
	Sort                    *int    `json:"sort"`
	RefID                   *string `json:"ref_id"`
	RefType                 *string `json:"ref_type"`
}

type UpdateJobValueRequest struct {
	Type             *string `json:"type" binding:"omitempty"`
	TypeGroup        *string `json:"type_group"`
	Level            *int    `json:"level"`
	Descriptions     *string `json:"descriptions"`
	DescriptionGroup *string `json:"description_group"`
	Note             *string `json:"note"`
	Sort             *int    `json:"sort"`
	RefID            *string `json:"ref_id"`
	RefType          *string `json:"ref_type"`
}

// =========================================================================
// Request DTOs — Job Objectives
// =========================================================================

type CreateJobObjectiveRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	Nomenclature   string `json:"nomenclature" binding:"required,max=50"`
	FullCode       string `json:"full_code" binding:"required,max=20"`
	Objective      string `json:"objective"`
}

type UpdateJobObjectiveRequest struct {
	Nomenclature *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode     *string `json:"full_code" binding:"omitempty,max=20"`
	Objective    *string `json:"objective"`
}

// =========================================================================
// Request DTOs — Job Identifications
// =========================================================================

type CreateJobIdentificationRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	Nomenclature   string `json:"nomenclature" binding:"required,max=50"`
	FullCode       string `json:"full_code" binding:"required,max=20"`
	GradingID      string `json:"grading_id" binding:"required"`
}

type UpdateJobIdentificationRequest struct {
	Nomenclature *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode     *string `json:"full_code" binding:"omitempty,max=20"`
	GradingID    *string `json:"grading_id"`
}

// =========================================================================
// Request DTOs — Job Responsibilities
// =========================================================================

type CreateJobResponsibilityRequest struct {
	OrganizationID    string `json:"organization_id" binding:"required"`
	Nomenclature      string `json:"nomenclature" binding:"required,max=50"`
	FullCode          string `json:"full_code" binding:"required,max=20"`
	MainTask          string `json:"main_task"`
	Activities        string `json:"activities"`
	Outputs           string `json:"outputs"`
	SuccessIndicators string `json:"success_indicators"`
}

type UpdateJobResponsibilityRequest struct {
	Nomenclature      *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode          *string `json:"full_code" binding:"omitempty,max=20"`
	MainTask          *string `json:"main_task"`
	Activities        *string `json:"activities"`
	Outputs           *string `json:"outputs"`
	SuccessIndicators *string `json:"success_indicators"`
}

// =========================================================================
// Request DTOs — Job Education Experiences
// =========================================================================

type CreateJobEducationExperienceRequest struct {
	OrganizationID   string   `json:"organization_id" binding:"required"`
	Nomenclature     string   `json:"nomenclature" binding:"required,max=50"`
	FullCode         string   `json:"full_code" binding:"required,max=20"`
	EducationID      *string  `json:"education_id"`
	ExperienceID     *string  `json:"experience_id"`
	EducationMajorID []string `json:"education_major_id"`
	JobFamilyID      []string `json:"job_family_id"`
}

type UpdateJobEducationExperienceRequest struct {
	Nomenclature     *string  `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode         *string  `json:"full_code" binding:"omitempty,max=20"`
	EducationID      *string  `json:"education_id"`
	ExperienceID     *string  `json:"experience_id"`
	EducationMajorID []string `json:"education_major_id"`
	JobFamilyID      []string `json:"job_family_id"`
}

// =========================================================================
// Request DTOs — Job HR Authorities
// =========================================================================

type CreateJobHRAuthorityRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	Nomenclature   string `json:"nomenclature" binding:"required,max=50"`
	FullCode       string `json:"full_code" binding:"required,max=20"`
	Description    string `json:"description"`
}

type UpdateJobHRAuthorityRequest struct {
	Nomenclature *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode     *string `json:"full_code" binding:"omitempty,max=20"`
	Description  *string `json:"description"`
}

// =========================================================================
// Request DTOs — Job Operational Authorities
// =========================================================================

type CreateJobOperationalAuthorityRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	Nomenclature   string `json:"nomenclature" binding:"required,max=50"`
	FullCode       string `json:"full_code" binding:"required,max=20"`
	Description    string `json:"description"`
}

type UpdateJobOperationalAuthorityRequest struct {
	Nomenclature *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode     *string `json:"full_code" binding:"omitempty,max=20"`
	Description  *string `json:"description"`
}

// =========================================================================
// Request DTOs — Job Working Activities
// =========================================================================

type CreateJobWorkingActivityRequest struct {
	OrganizationID      string  `json:"organization_id" binding:"required"`
	Nomenclature        string  `json:"nomenclature" binding:"required,max=50"`
	FullCode            string  `json:"full_code" binding:"required,max=20"`
	JobManagementValueID *string `json:"job_management_value_id"`
}

type UpdateJobWorkingActivityRequest struct {
	Nomenclature        *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode            *string `json:"full_code" binding:"omitempty,max=20"`
	JobManagementValueID *string `json:"job_management_value_id"`
}

// =========================================================================
// Request DTOs — Job Working Risks
// =========================================================================

type CreateJobWorkingRiskRequest struct {
	OrganizationID                   string  `json:"organization_id" binding:"required"`
	Nomenclature                     string  `json:"nomenclature" binding:"required,max=50"`
	FullCode                         string  `json:"full_code" binding:"required,max=20"`
	JobManagementValueEnvironmentID  *string `json:"job_management_value_environment_id"`
	JobManagementValueHazardID       *string `json:"job_management_value_hazard_id"`
}

type UpdateJobWorkingRiskRequest struct {
	Nomenclature                     *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode                         *string `json:"full_code" binding:"omitempty,max=20"`
	JobManagementValueEnvironmentID  *string `json:"job_management_value_environment_id"`
	JobManagementValueHazardID       *string `json:"job_management_value_hazard_id"`
}

// =========================================================================
// Request DTOs — Job Relationships
// =========================================================================

type CreateJobRelationshipRequest struct {
	OrganizationID                  string  `json:"organization_id" binding:"required"`
	Nomenclature                    string  `json:"nomenclature" binding:"required,max=50"`
	FullCode                        string  `json:"full_code" binding:"required,max=20"`
	JobManagementValueRelationshipID *string `json:"job_management_value_relationship_id"`
	JobManagementValueFrequencyID   *string `json:"job_management_value_frequency_id"`
}

type UpdateJobRelationshipRequest struct {
	Nomenclature                    *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode                        *string `json:"full_code" binding:"omitempty,max=20"`
	JobManagementValueRelationshipID *string `json:"job_management_value_relationship_id"`
	JobManagementValueFrequencyID   *string `json:"job_management_value_frequency_id"`
}

// =========================================================================
// Request DTOs — Job Relationship Details
// =========================================================================

type CreateJobRelationshipDetailRequest struct {
	OrganizationID *string `json:"organization_id"`
	Activity       *string `json:"activity"`
}

type UpdateJobRelationshipDetailRequest struct {
	OrganizationID *string `json:"organization_id"`
	Activity       *string `json:"activity"`
}

// =========================================================================
// Request DTOs — Job Subordinate Controls
// =========================================================================

type CreateJobSubordinateControlRequest struct {
	OrganizationID      string  `json:"organization_id" binding:"required"`
	Nomenclature        string  `json:"nomenclature" binding:"required,max=50"`
	FullCode            string  `json:"full_code" binding:"required,max=20"`
	JobManagementValueID *string `json:"job_management_value_id"`
}

type UpdateJobSubordinateControlRequest struct {
	Nomenclature        *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode            *string `json:"full_code" binding:"omitempty,max=20"`
	JobManagementValueID *string `json:"job_management_value_id"`
}

// =========================================================================
// Request DTOs — Job Assets
// =========================================================================

type CreateJobAssetRequest struct {
	OrganizationID              string  `json:"organization_id" binding:"required"`
	Nomenclature                string  `json:"nomenclature" binding:"required,max=50"`
	FullCode                    string  `json:"full_code" binding:"required,max=20"`
	JobManagementValueAssetID   *string `json:"job_management_value_asset_id"`
	JobManagementValueAuthorityID *string `json:"job_management_value_authority_id"`
}

type UpdateJobAssetRequest struct {
	Nomenclature                *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode                    *string `json:"full_code" binding:"omitempty,max=20"`
	JobManagementValueAssetID   *string `json:"job_management_value_asset_id"`
	JobManagementValueAuthorityID *string `json:"job_management_value_authority_id"`
}

// =========================================================================
// Request DTOs — Job Financials
// =========================================================================

type CreateJobFinancialRequest struct {
	OrganizationID              string  `json:"organization_id" binding:"required"`
	Nomenclature                string  `json:"nomenclature" binding:"required,max=50"`
	FullCode                    string  `json:"full_code" binding:"required,max=20"`
	IsAuthorized                bool    `json:"is_authorized"`
	JobManagementValueCashID    *string `json:"job_management_value_cash_id"`
	JobManagementValueAuthorityID *string `json:"job_management_value_authority_id"`
	JobManagementValueImpactID  *string `json:"job_management_value_impact_id"`
}

type UpdateJobFinancialRequest struct {
	Nomenclature                *string `json:"nomenclature" binding:"omitempty,max=50"`
	FullCode                    *string `json:"full_code" binding:"omitempty,max=20"`
	IsAuthorized                *bool   `json:"is_authorized"`
	JobManagementValueCashID    *string `json:"job_management_value_cash_id"`
	JobManagementValueAuthorityID *string `json:"job_management_value_authority_id"`
	JobManagementValueImpactID  *string `json:"job_management_value_impact_id"`
}

// =========================================================================
// Request DTOs — Job Potency Competencies
// =========================================================================

type CreateJobPotencyCompetencyRequest struct {
	OrganizationID      string   `json:"organization_id" binding:"required"`
	JobManagementValueID *string  `json:"job_management_value_id"`
	CompetencyID        *string  `json:"competency_id"`
	Weight              *float64 `json:"weight"`
}

type UpdateJobPotencyCompetencyRequest struct {
	JobManagementValueID *string  `json:"job_management_value_id"`
	CompetencyID        *string  `json:"competency_id"`
	Weight              *float64 `json:"weight"`
}

// =========================================================================
// Request DTOs — Job Scores
// =========================================================================

type UpdateJobScoreRequest struct {
	JobValueWithFinancial    *uint64 `json:"job_value_with_financial"`
	JobValueWithoutFinancial *uint64 `json:"job_value_without_financial"`
	HasFinancialAuthority    *bool   `json:"has_financial_authority"`
	Components               *string `json:"components"`
	SubComponentPoints       *string `json:"sub_component_points"`
}

// =========================================================================
// Request DTOs — Job Competency Groups
// =========================================================================

type CreateJobCompetencyGroupRequest struct {
	OrganizationID string  `json:"organization_id" binding:"required"`
	Category       string  `json:"category" binding:"required,oneof=technical managerial"`
	Weight         float64 `json:"weight" binding:"required"`
}

type UpdateJobCompetencyGroupRequest struct {
	Category *string  `json:"category" binding:"omitempty,oneof=technical managerial"`
	Weight   *float64 `json:"weight"`
}

// =========================================================================
// Response DTOs
// =========================================================================

type JobTitleResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name,omitempty"`
	Descriptions string    `json:"descriptions,omitempty"`
	Status       int8      `json:"status,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Subs         []JobTitleSubResponse `json:"subs,omitempty"`
}

type JobTitleSubResponse struct {
	ID                      string    `json:"id"`
	JobManagementTitleID    string    `json:"job_management_title_id,omitempty"`
	JobManagementTitleName  string    `json:"job_management_title_name,omitempty"`
	Name                    string    `json:"name,omitempty"`
	Descriptions            string    `json:"descriptions,omitempty"`
	Status                  int8      `json:"status,omitempty"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type JobValueResponse struct {
	ID                         string    `json:"id"`
	JobManagementTitleSubID    string    `json:"job_management_title_sub_id,omitempty"`
	JobManagementTitleSubName  string    `json:"job_management_title_sub_name,omitempty"`
	Type                       string    `json:"type"`
	TypeGroup                  string    `json:"type_group,omitempty"`
	Level                      int       `json:"level,omitempty"`
	Descriptions               string    `json:"descriptions,omitempty"`
	DescriptionGroup           string    `json:"description_group,omitempty"`
	Note                       string    `json:"note,omitempty"`
	Sort                       int       `json:"sort,omitempty"`
	RefID                      string    `json:"ref_id,omitempty"`
	RefType                    string    `json:"ref_type,omitempty"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

// Request/Response DTOs — Job Value Clusters (mapping type ↔ cluster kompetensi)
type UpdateJobValueClustersRequest struct {
	// Tanpa binding:"required" — array kosong (menghapus seluruh mapping) adalah aksi valid.
	Clusters []string `json:"clusters"`
}

type JobValueClusterResponse struct {
	Type     string   `json:"type"`
	Clusters []string `json:"clusters"`
}

type JobObjectiveResponse struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id,omitempty"`
	Nomenclature   string    `json:"nomenclature"`
	FullCode       string    `json:"full_code"`
	Objective      string    `json:"objective,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type JobIdentificationResponse struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id,omitempty"`
	Nomenclature   string    `json:"nomenclature"`
	FullCode       string    `json:"full_code"`
	GradingID      string    `json:"grading_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type JobResponsibilityResponse struct {
	ID                string    `json:"id"`
	OrganizationID    string    `json:"organization_id,omitempty"`
	Nomenclature      string    `json:"nomenclature"`
	FullCode          string    `json:"full_code"`
	MainTask          string    `json:"main_task,omitempty"`
	Activities        string    `json:"activities,omitempty"`
	Outputs           string    `json:"outputs,omitempty"`
	SuccessIndicators string    `json:"success_indicators,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type JobEducationExperienceResponse struct {
	ID               string    `json:"id"`
	OrganizationID   string    `json:"organization_id,omitempty"`
	Nomenclature     string    `json:"nomenclature"`
	FullCode         string    `json:"full_code"`
	EducationID      string    `json:"education_id,omitempty"`
	EducationName    string    `json:"education_name,omitempty"`
	ExperienceID     string    `json:"experience_id,omitempty"`
	ExperienceName   string    `json:"experience_name,omitempty"`
	EducationMajorID []string  `json:"education_major_id,omitempty"`
	EducationMajorName []string `json:"education_major_name,omitempty"`
	JobFamilyID      []string  `json:"job_family_id,omitempty"`
	JobFamilyName    []string  `json:"job_family_name,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type JobHRAuthorityResponse struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id,omitempty"`
	Nomenclature   string    `json:"nomenclature"`
	FullCode       string    `json:"full_code"`
	Description    string    `json:"description,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type JobOperationalAuthorityResponse struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id,omitempty"`
	Nomenclature   string    `json:"nomenclature"`
	FullCode       string    `json:"full_code"`
	Description    string    `json:"description,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type JobWorkingActivityResponse struct {
	ID                  string    `json:"id"`
	OrganizationID      string    `json:"organization_id,omitempty"`
	Nomenclature        string    `json:"nomenclature"`
	FullCode            string    `json:"full_code"`
	JobManagementValueID string   `json:"job_management_value_id,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type JobWorkingRiskResponse struct {
	ID                              string    `json:"id"`
	OrganizationID                  string    `json:"organization_id,omitempty"`
	Nomenclature                    string    `json:"nomenclature"`
	FullCode                        string    `json:"full_code"`
	JobManagementValueEnvironmentID string    `json:"job_management_value_environment_id,omitempty"`
	JobManagementValueHazardID      string    `json:"job_management_value_hazard_id,omitempty"`
	CreatedAt                       time.Time `json:"created_at"`
	UpdatedAt                       time.Time `json:"updated_at"`
}

type JobRelationshipResponse struct {
	ID                              string    `json:"id"`
	OrganizationID                  string    `json:"organization_id,omitempty"`
	Nomenclature                    string    `json:"nomenclature"`
	FullCode                        string    `json:"full_code"`
	JobManagementValueRelationshipID string   `json:"job_management_value_relationship_id,omitempty"`
	JobManagementValueFrequencyID   string    `json:"job_management_value_frequency_id,omitempty"`
	CreatedAt                       time.Time `json:"created_at"`
	UpdatedAt                       time.Time `json:"updated_at"`
}

type JobSubordinateControlResponse struct {
	ID                  string    `json:"id"`
	OrganizationID      string    `json:"organization_id,omitempty"`
	Nomenclature        string    `json:"nomenclature"`
	FullCode            string    `json:"full_code"`
	JobManagementValueID string   `json:"job_management_value_id,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type JobRelationshipDetailResponse struct {
	ID                         string    `json:"id"`
	JobManagementRelationshipID string   `json:"job_management_relationship_id,omitempty"`
	OrganizationID             string    `json:"organization_id,omitempty"`
	OrganizationName           string    `json:"organization_name,omitempty"`
	OrganizationCode           string    `json:"organization_code,omitempty"`
	Activity                   string    `json:"activity,omitempty"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

type JobAssetResponse struct {
	ID                          string    `json:"id"`
	OrganizationID              string    `json:"organization_id,omitempty"`
	Nomenclature                string    `json:"nomenclature"`
	FullCode                    string    `json:"full_code"`
	JobManagementValueAssetID   string    `json:"job_management_value_asset_id,omitempty"`
	JobManagementValueAuthorityID string  `json:"job_management_value_authority_id,omitempty"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
}

type JobFinancialResponse struct {
	ID                          string    `json:"id"`
	OrganizationID              string    `json:"organization_id,omitempty"`
	Nomenclature                string    `json:"nomenclature"`
	FullCode                    string    `json:"full_code"`
	IsAuthorized                bool      `json:"is_authorized"`
	JobManagementValueCashID    string    `json:"job_management_value_cash_id,omitempty"`
	JobManagementValueAuthorityID string  `json:"job_management_value_authority_id,omitempty"`
	JobManagementValueImpactID  string    `json:"job_management_value_impact_id,omitempty"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
}

type JobPotencyCompetencyResponse struct {
	ID                  string    `json:"id"`
	OrganizationID      string    `json:"organization_id,omitempty"`
	JobManagementValueID string   `json:"job_management_value_id,omitempty"`
	CompetencyID        string    `json:"competency_id,omitempty"`
	CompetencyName      string    `json:"competency_name,omitempty"`
	Type                string    `json:"type,omitempty"`
	Level               int       `json:"level,omitempty"`
	LevelDescription    string    `json:"level_description,omitempty"`
	Weight              float64   `json:"weight,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type JobScoreResponse struct {
	ID                       string     `json:"id"`
	OrganizationID           string     `json:"organization_id"`
	JobValueWithFinancial    uint64     `json:"job_value_with_financial"`
	JobValueWithoutFinancial uint64     `json:"job_value_without_financial"`
	HasFinancialAuthority    bool       `json:"has_financial_authority"`
	Components               string     `json:"components,omitempty"`
	SubComponentPoints       string     `json:"sub_component_points,omitempty"`
	CalculatedAt             *time.Time `json:"calculated_at,omitempty"`
	IsComplete               bool       `json:"is_complete"`
	CompletedAt              *time.Time `json:"completed_at,omitempty"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}

type JobCompetencyGroupResponse struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Category       string    `json:"category"`
	Weight         float64   `json:"weight"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// =========================================================================
// Converter Helpers
// =========================================================================

func toJobTitleResponse(t *JobTitle) JobTitleResponse {
	r := JobTitleResponse{
		ID:        t.ID.String(),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
	if t.Name != nil {
		r.Name = *t.Name
	}
	if t.Descriptions != nil {
		r.Descriptions = *t.Descriptions
	}
	if t.Status != nil {
		r.Status = *t.Status
	}
	if len(t.Subs) > 0 {
		r.Subs = make([]JobTitleSubResponse, 0, len(t.Subs))
		for _, s := range t.Subs {
			r.Subs = append(r.Subs, toJobTitleSubResponse(&s))
		}
	}
	return r
}

func toJobTitleSubResponse(s *JobTitleSub) JobTitleSubResponse {
	r := JobTitleSubResponse{
		ID:        s.ID.String(),
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
	if s.JobManagementTitleID != nil {
		r.JobManagementTitleID = s.JobManagementTitleID.String()
	}
	if s.JobManagementTitleName != nil {
		r.JobManagementTitleName = *s.JobManagementTitleName
	}
	if s.Name != nil {
		r.Name = *s.Name
	}
	if s.Descriptions != nil {
		r.Descriptions = *s.Descriptions
	}
	if s.Status != nil {
		r.Status = *s.Status
	}
	return r
}

// =========================================================================
// Response DTOs — Job Values Tree (grouped by type_group)
// =========================================================================

// JobValueTreeOption — option level dalam satu tipe (level + deskripsi)
type JobValueTreeOption struct {
	ID           string `json:"id"`
	Level        int    `json:"level,omitempty"`
	Descriptions string `json:"descriptions,omitempty"`
}

// JobValueTreeType — satu tipe dalam grup (label = description_group)
type JobValueTreeType struct {
	Type             string               `json:"type"`
	DescriptionGroup string               `json:"description_group"`
	Options          []JobValueTreeOption `json:"options"`
}

// JobValueTreeGroup — satu group (type_group) berisi daftar tipe.
// DescriptionGroup bersifat per-tipe (mis. 'Kecerdasan', 'Innovation & Creativity');
// label group (DescriptionGroup) di sini hanyalah fallback dari tipe pertama.
type JobValueTreeGroup struct {
	TypeGroup        string             `json:"type_group"`
	DescriptionGroup string             `json:"description_group"`
	Types            []JobValueTreeType `json:"types"`
}

// JobValueTreeResponse — respons endpoint tree
// Group diurutkan secara konsisten (pendidikan → pengalaman → psychological →
// technical → managerial → communication → problem_solving → financial → asset →
// sisanya), tipe diurutkan berdasarkan description_group.
type JobValueTreeResponse struct {
	Success bool                `json:"success"`
	Data    []JobValueTreeGroup `json:"data"`
}

func toJobValueResponse(v *JobValue) JobValueResponse {
	r := JobValueResponse{
		ID:        v.ID.String(),
		Type:      v.Type,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
	if v.JobManagementTitleSubID != nil {
		r.JobManagementTitleSubID = v.JobManagementTitleSubID.String()
	}
	if v.JobManagementTitleSubName != nil {
		r.JobManagementTitleSubName = *v.JobManagementTitleSubName
	}
	if v.TypeGroup != nil {
		r.TypeGroup = *v.TypeGroup
	}
	if v.Level != nil {
		r.Level = *v.Level
	}
	if v.Descriptions != nil {
		r.Descriptions = *v.Descriptions
	}
	if v.DescriptionGroup != nil {
		r.DescriptionGroup = *v.DescriptionGroup
	}
	if v.Note != nil {
		r.Note = *v.Note
	}
	if v.Sort != nil {
		r.Sort = *v.Sort
	}
	if v.RefID != nil {
		r.RefID = v.RefID.String()
	}
	if v.RefType != nil {
		r.RefType = *v.RefType
	}
	return r
}

func toJobObjectiveResponse(o *JobObjective) JobObjectiveResponse {
	r := JobObjectiveResponse{
		ID:           o.ID.String(),
		Nomenclature: o.Nomenclature,
		FullCode:     o.FullCode,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}
	if o.OrganizationID != nil {
		r.OrganizationID = o.OrganizationID.String()
	}
	if o.Objective != nil {
		r.Objective = *o.Objective
	}
	return r
}

func toJobIdentificationResponse(i *JobIdentification) JobIdentificationResponse {
	r := JobIdentificationResponse{
		ID:           i.ID.String(),
		Nomenclature: i.Nomenclature,
		FullCode:     i.FullCode,
		GradingID:    i.GradingID.String(),
		CreatedAt:    i.CreatedAt,
		UpdatedAt:    i.UpdatedAt,
	}
	if i.OrganizationID != nil {
		r.OrganizationID = i.OrganizationID.String()
	}
	return r
}

func toJobResponsibilityResponse(r *JobResponsibility) JobResponsibilityResponse {
	resp := JobResponsibilityResponse{
		ID:           r.ID.String(),
		Nomenclature: r.Nomenclature,
		FullCode:     r.FullCode,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
	if r.OrganizationID != nil {
		resp.OrganizationID = r.OrganizationID.String()
	}
	if r.MainTask != nil {
		resp.MainTask = *r.MainTask
	}
	if r.Activities != nil {
		resp.Activities = *r.Activities
	}
	if r.Outputs != nil {
		resp.Outputs = *r.Outputs
	}
	if r.SuccessIndicators != nil {
		resp.SuccessIndicators = *r.SuccessIndicators
	}
	return resp
}

func toJobEducationExperienceResponse(e *JobEducationExperience) JobEducationExperienceResponse {
	r := JobEducationExperienceResponse{
		ID:           e.ID.String(),
		Nomenclature: e.Nomenclature,
		FullCode:     e.FullCode,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
	if e.OrganizationID != nil {
		r.OrganizationID = e.OrganizationID.String()
	}
	if e.EducationID != nil {
		r.EducationID = e.EducationID.String()
	}
	// Education → job_management_values (type=education); nama tampil dari Descriptions
	if e.Education != nil && e.Education.Descriptions != nil {
		r.EducationName = *e.Education.Descriptions
	}
	if e.ExperienceID != nil {
		r.ExperienceID = e.ExperienceID.String()
	}
	// Experience → job_management_values (type=experience); nama tampil dari Descriptions
	if e.Experience != nil && e.Experience.Descriptions != nil {
		r.ExperienceName = *e.Experience.Descriptions
	}
	// Jurusan (multiple) — via pivot job_management_majors
	if len(e.Majors) > 0 {
		ids := make([]string, 0, len(e.Majors))
		names := make([]string, 0, len(e.Majors))
		for _, m := range e.Majors {
			ids = append(ids, m.EducationMajorID.String())
			if m.EducationMajor != nil {
				names = append(names, m.EducationMajor.Name)
			}
		}
		r.EducationMajorID = ids
		r.EducationMajorName = names
	}
	// Bidang Pekerjaan (multiple) — via pivot job_management_job_family
	if len(e.JobFamilies) > 0 {
		ids := make([]string, 0, len(e.JobFamilies))
		names := make([]string, 0, len(e.JobFamilies))
		for _, jf := range e.JobFamilies {
			ids = append(ids, jf.JobFamilyID.String())
			if jf.JobFamily != nil {
				names = append(names, jf.JobFamily.Name)
			}
		}
		r.JobFamilyID = ids
		r.JobFamilyName = names
	}
	return r
}

func toJobHRAuthorityResponse(a *JobHRAuthority) JobHRAuthorityResponse {
	r := JobHRAuthorityResponse{
		ID:           a.ID.String(),
		Nomenclature: a.Nomenclature,
		FullCode:     a.FullCode,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
	if a.OrganizationID != nil {
		r.OrganizationID = a.OrganizationID.String()
	}
	if a.Description != nil {
		r.Description = *a.Description
	}
	return r
}

func toJobOperationalAuthorityResponse(a *JobOperationalAuthority) JobOperationalAuthorityResponse {
	r := JobOperationalAuthorityResponse{
		ID:           a.ID.String(),
		Nomenclature: a.Nomenclature,
		FullCode:     a.FullCode,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
	if a.OrganizationID != nil {
		r.OrganizationID = a.OrganizationID.String()
	}
	if a.Description != nil {
		r.Description = *a.Description
	}
	return r
}

func toJobWorkingActivityResponse(a *JobWorkingActivity) JobWorkingActivityResponse {
	r := JobWorkingActivityResponse{
		ID:           a.ID.String(),
		Nomenclature: a.Nomenclature,
		FullCode:     a.FullCode,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
	if a.OrganizationID != nil {
		r.OrganizationID = a.OrganizationID.String()
	}
	if a.JobManagementValueID != nil {
		r.JobManagementValueID = a.JobManagementValueID.String()
	}
	return r
}

func toJobWorkingRiskResponse(r *JobWorkingRisk) JobWorkingRiskResponse {
	resp := JobWorkingRiskResponse{
		ID:           r.ID.String(),
		Nomenclature: r.Nomenclature,
		FullCode:     r.FullCode,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
	if r.OrganizationID != nil {
		resp.OrganizationID = r.OrganizationID.String()
	}
	if r.JobManagementValueEnvironmentID != nil {
		resp.JobManagementValueEnvironmentID = r.JobManagementValueEnvironmentID.String()
	}
	if r.JobManagementValueHazardID != nil {
		resp.JobManagementValueHazardID = r.JobManagementValueHazardID.String()
	}
	return resp
}

func toJobRelationshipResponse(r *JobRelationship) JobRelationshipResponse {
	resp := JobRelationshipResponse{
		ID:           r.ID.String(),
		Nomenclature: r.Nomenclature,
		FullCode:     r.FullCode,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
	if r.OrganizationID != nil {
		resp.OrganizationID = r.OrganizationID.String()
	}
	if r.JobManagementValueRelationshipID != nil {
		resp.JobManagementValueRelationshipID = r.JobManagementValueRelationshipID.String()
	}
	if r.JobManagementValueFrequencyID != nil {
		resp.JobManagementValueFrequencyID = r.JobManagementValueFrequencyID.String()
	}
	return resp
}

func toJobRelationshipDetailResponse(d *JobManagementRelationshipDetail) JobRelationshipDetailResponse {
	r := JobRelationshipDetailResponse{
		ID:                         d.ID.String(),
		JobManagementRelationshipID: d.JobManagementRelationshipID.String(),
		CreatedAt:                  d.CreatedAt,
		UpdatedAt:                  d.UpdatedAt,
	}
	if d.OrganizationID != nil {
		r.OrganizationID = d.OrganizationID.String()
	}
	if d.Activity != nil {
		r.Activity = *d.Activity
	}
	if d.Organization != nil {
		r.OrganizationName = d.Organization.Nomenclature
		r.OrganizationCode = d.Organization.FullCode
	}
	return r
}

func toJobSubordinateControlResponse(c *JobSubordinateControl) JobSubordinateControlResponse {
	r := JobSubordinateControlResponse{
		ID:           c.ID.String(),
		Nomenclature: c.Nomenclature,
		FullCode:     c.FullCode,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
	if c.OrganizationID != nil {
		r.OrganizationID = c.OrganizationID.String()
	}
	if c.JobManagementValueID != nil {
		r.JobManagementValueID = c.JobManagementValueID.String()
	}
	return r
}

func toJobAssetResponse(a *JobAsset) JobAssetResponse {
	r := JobAssetResponse{
		ID:           a.ID.String(),
		Nomenclature: a.Nomenclature,
		FullCode:     a.FullCode,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
	if a.OrganizationID != nil {
		r.OrganizationID = a.OrganizationID.String()
	}
	if a.JobManagementValueAssetID != nil {
		r.JobManagementValueAssetID = a.JobManagementValueAssetID.String()
	}
	if a.JobManagementValueAuthorityID != nil {
		r.JobManagementValueAuthorityID = a.JobManagementValueAuthorityID.String()
	}
	return r
}

func toJobFinancialResponse(f *JobFinancial) JobFinancialResponse {
	r := JobFinancialResponse{
		ID:           f.ID.String(),
		Nomenclature: f.Nomenclature,
		FullCode:     f.FullCode,
		IsAuthorized: f.IsAuthorized,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    f.UpdatedAt,
	}
	if f.OrganizationID != nil {
		r.OrganizationID = f.OrganizationID.String()
	}
	if f.JobManagementValueCashID != nil {
		r.JobManagementValueCashID = f.JobManagementValueCashID.String()
	}
	if f.JobManagementValueAuthorityID != nil {
		r.JobManagementValueAuthorityID = f.JobManagementValueAuthorityID.String()
	}
	if f.JobManagementValueImpactID != nil {
		r.JobManagementValueImpactID = f.JobManagementValueImpactID.String()
	}
	return r
}

func toJobPotencyCompetencyResponse(c *JobPotencyCompetency) JobPotencyCompetencyResponse {
	r := JobPotencyCompetencyResponse{
		ID:        c.ID.String(),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
	if c.OrganizationID != nil {
		r.OrganizationID = c.OrganizationID.String()
	}
	if c.JobManagementValueID != nil {
		r.JobManagementValueID = c.JobManagementValueID.String()
	}
	if c.CompetencyID != nil {
		r.CompetencyID = c.CompetencyID.String()
	}
	if c.Competency != nil {
		r.CompetencyName = c.Competency.Name
	}
	if c.JobManagementValue != nil && c.JobManagementValue.Type != "" {
		r.Type = c.JobManagementValue.Type
	}
	if c.JobManagementValue != nil {
		if c.JobManagementValue.Level != nil {
			r.Level = *c.JobManagementValue.Level
		}
		if c.JobManagementValue.Descriptions != nil {
			r.LevelDescription = *c.JobManagementValue.Descriptions
		}
	}
	if c.Weight != nil {
		r.Weight = *c.Weight
	}
	return r
}

func toJobScoreResponse(s *JobScore) JobScoreResponse {
	r := JobScoreResponse{
		ID:                       s.ID.String(),
		JobValueWithFinancial:    s.JobValueWithFinancial,
		JobValueWithoutFinancial: s.JobValueWithoutFinancial,
		HasFinancialAuthority:    s.HasFinancialAuthority,
		CalculatedAt:             s.CalculatedAt,
		IsComplete:               s.IsComplete,
		CompletedAt:              s.CompletedAt,
		CreatedAt:                s.CreatedAt,
		UpdatedAt:                s.UpdatedAt,
	}
	if s.OrganizationID != nil {
		r.OrganizationID = s.OrganizationID.String()
	}
	if s.Components != nil {
		r.Components = *s.Components
	}
	if s.SubComponentPoints != nil {
		r.SubComponentPoints = *s.SubComponentPoints
	}
	return r
}

// =========================================================================
// Dashboard
// =========================================================================

// OrgSummaryInfo — ringkasan organisasi AKTIF (organization_summaries
// dengan status=active) yang menjadi acuan seluruh data dashboard Job
// Management.
type OrgSummaryInfo struct {
	ID         string `json:"id"`
	Code       string `json:"code"`
	DecreeNo   string `json:"decree_no"`
	DecreeDate string `json:"decree_date"`
}

// JobManagementDashboardResponse — ringkasan summary organisasi aktif untuk
// dashboard Job Management (GET /job-management/dashboard). Semua angka
// mengacu ke organisasi milik summary aktif tersebut.
type JobManagementDashboardResponse struct {
	// Summary aktif; null jika belum ada summary berstatus active.
	Summary            *OrgSummaryInfo `json:"summary,omitempty"`
	TotalOrganizations int             `json:"total_organizations"`
	// Organisasi yang sudah terisi karyawan (employment berjalan) vs belum.
	WithEmployees     int `json:"with_employees"`
	WithoutEmployees  int `json:"without_employees"`
	// Progres pengisian value (job_management_scores per organisasi).
	ValueNotStarted int `json:"value_not_started"`
	ValueOnProgress int `json:"value_on_progress"`
	ValueCompleted  int `json:"value_completed"`
	// Wewenang keuangan (job_management_financials.is_authorized).
	WithFinancialAuthority    int `json:"with_financial_authority"`
	WithoutFinancialAuthority int `json:"without_financial_authority"`
}

func toJobCompetencyGroupResponse(g *JobCompetencyGroup) JobCompetencyGroupResponse {
	r := JobCompetencyGroupResponse{
		ID:        g.ID.String(),
		Category:  g.Category,
		Weight:    g.Weight,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
	if g.OrganizationID != nil {
		r.OrganizationID = g.OrganizationID.String()
	}
	return r
}
