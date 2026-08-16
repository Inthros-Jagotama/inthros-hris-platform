package competency

import "time"

// =========================================================================
// Request DTOs — Rating Scale
// =========================================================================

type CreateRatingScaleRequest struct {
	// Code tidak diterima dari client — dibuat otomatis dari Name di service
	// (Service.generateUniqueCode), jadi form pembuatan tidak perlu memintanya.
	Name        string                    `json:"name" binding:"required,max=255"`
	Description *string                   `json:"description"`
	Status      string                    `json:"status" binding:"omitempty,oneof=active inactive"`
	Items       []RatingScaleItemRequest  `json:"items"`
}

type RatingScaleItemRequest struct {
	Value       int     `json:"value" binding:"required"`
	Label       string  `json:"label" binding:"required,max=255"`
	Description *string `json:"description"`
	Weight      float64 `json:"weight"`
	SortOrder   int     `json:"sort_order"`
}

type UpdateRatingScaleRequest struct {
	Name        *string                   `json:"name" binding:"omitempty,max=255"`
	Code        *string                   `json:"code" binding:"omitempty,max=50"`
	Description *string                   `json:"description"`
	Status      *string                   `json:"status" binding:"omitempty,oneof=active inactive"`
	Items       []RatingScaleItemRequest  `json:"items"`
}

// =========================================================================
// Request DTOs — Assessment Template
// =========================================================================

type TemplateCompetencyRequest struct {
	CompetencyID  string  `json:"competency_id" binding:"required"`
	RequiredLevel *int    `json:"required_level"`
	Weight        float64 `json:"weight"`
	SortOrder     int     `json:"sort_order"`
}

type TemplateRaterTypeRequest struct {
	RaterType string  `json:"rater_type" binding:"required,oneof=self superior peer subordinate other"`
	Weight    float64 `json:"weight"`
	MinRater  int     `json:"min_rater"`
	MaxRater  *int    `json:"max_rater"`
	Required  bool    `json:"required"`
	Anonymous bool    `json:"anonymous"`
}

type CreateAssessmentTemplateRequest struct {
	// Code tidak diterima dari client — dibuat otomatis dari Name di service
	// (Service.generateUniqueCode), jadi form pembuatan tidak perlu memintanya.
	Name         string                      `json:"name" binding:"required,max=255"`
	Description  *string                     `json:"description"`
	Status       string                      `json:"status" binding:"omitempty,oneof=active inactive"`
	ScaleID      *string                     `json:"scale_id"`
	Competencies []TemplateCompetencyRequest `json:"competencies"`
	RaterTypes   []TemplateRaterTypeRequest  `json:"rater_types"`
}

type UpdateAssessmentTemplateRequest struct {
	Name          *string                     `json:"name" binding:"omitempty,max=255"`
	Code          *string                     `json:"code" binding:"omitempty,max=50"`
	Description   *string                     `json:"description"`
	Status        *string                     `json:"status" binding:"omitempty,oneof=active inactive"`
	ScaleID       *string                     `json:"scale_id"`
	Competencies  []TemplateCompetencyRequest `json:"competencies"`
	RaterTypes    []TemplateRaterTypeRequest  `json:"rater_types"`
}

// =========================================================================
// Request DTOs — Indicator
// =========================================================================

type CreateIndicatorRequest struct {
	CompetencyID string  `json:"competency_id" binding:"required"`
	Code         *string `json:"code" binding:"omitempty,max=50"`
	Statement    string  `json:"statement" binding:"required,max=1000"`
	Description  *string `json:"description"`
	Status       string  `json:"status" binding:"omitempty,oneof=active inactive"`
	SortOrder    int     `json:"sort_order"`
}

type UpdateIndicatorRequest struct {
	CompetencyID *string `json:"competency_id"`
	Code         *string `json:"code" binding:"omitempty,max=50"`
	Statement    *string `json:"statement" binding:"omitempty,max=1000"`
	Description  *string `json:"description"`
	Status       *string `json:"status" binding:"omitempty,oneof=active inactive"`
	SortOrder    *int    `json:"sort_order"`
}

type TemplateIndicatorRequest struct {
	IndicatorID string  `json:"indicator_id" binding:"required"`
	Weight      float64 `json:"weight"`
	SortOrder   int     `json:"sort_order"`
}

// =========================================================================
// Response DTOs
// =========================================================================

type RatingScaleItemResponse struct {
	ID          string    `json:"id"`
	ScaleID     string    `json:"scale_id"`
	Value       int       `json:"value"`
	Label       string    `json:"label"`
	Description string    `json:"description,omitempty"`
	Weight      float64   `json:"weight"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RatingScaleResponse struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Code        string                   `json:"code"`
	Description string                   `json:"description,omitempty"`
	Status      string                   `json:"status"`
	Items       []RatingScaleItemResponse `json:"items,omitempty"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

type TemplateCompetencyResponse struct {
	ID            string    `json:"id"`
	TemplateID    string    `json:"template_id"`
	CompetencyID  string    `json:"competency_id"`
	CompetencyName string   `json:"competency_name,omitempty"`
	RequiredLevel int       `json:"required_level,omitempty"`
	Weight        float64   `json:"weight"`
	SortOrder     int       `json:"sort_order"`
}

type TemplateRaterTypeResponse struct {
	ID         string    `json:"id"`
	TemplateID string    `json:"template_id"`
	RaterType  string    `json:"rater_type"`
	Weight     float64   `json:"weight"`
	MinRater   int       `json:"min_rater"`
	MaxRater   int       `json:"max_rater,omitempty"`
	Required   bool      `json:"required"`
	Anonymous  bool      `json:"anonymous"`
}

type AssessmentTemplateResponse struct {
	ID           string                        `json:"id"`
	Name         string                        `json:"name"`
	Code         string                        `json:"code"`
	Description  string                        `json:"description,omitempty"`
	Status       string                        `json:"status"`
	ScaleID      string                        `json:"scale_id,omitempty"`
	Competencies []TemplateCompetencyResponse  `json:"competencies,omitempty"`
	RaterTypes   []TemplateRaterTypeResponse   `json:"rater_types,omitempty"`
	CreatedAt    time.Time                     `json:"created_at"`
	UpdatedAt    time.Time                     `json:"updated_at"`
}

type IndicatorResponse struct {
	ID            string    `json:"id"`
	CompetencyID  string    `json:"competency_id"`
	CompetencyName string   `json:"competency_name,omitempty"`
	Code          string    `json:"code,omitempty"`
	Statement     string    `json:"statement"`
	Description   string    `json:"description,omitempty"`
	Status        string    `json:"status"`
	SortOrder     int       `json:"sort_order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type EmployeeBriefDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SuggestedRatersDTO — saran rater untuk satu target: self (subject sendiri,
// bila template mewajibkan & belum di-assign), superior (atasan — parent org),
// dan subordinates (bawahan — subtree org).
type SuggestedRatersDTO struct {
	Self         *EmployeeBriefDTO  `json:"self,omitempty"`
	Superior     []EmployeeBriefDTO `json:"superior"`
	Subordinates []EmployeeBriefDTO `json:"subordinates"`
}

type TemplateIndicatorResponse struct {
	ID           string    `json:"id"`
	TemplateID   string    `json:"template_id"`
	IndicatorID  string    `json:"indicator_id"`
	Statement    string    `json:"statement,omitempty"`
	Weight       float64   `json:"weight"`
	SortOrder    int       `json:"sort_order"`
}

// =========================================================================
// Converter helpers
// =========================================================================

func (s *CompetencyRatingScale) ToResponse() RatingScaleResponse {
	r := RatingScaleResponse{
		ID:        s.ID.String(),
		Name:      s.Name,
		Code:      s.Code,
		Status:    s.Status,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
	if s.Description != nil {
		r.Description = *s.Description
	}
	if s.Items != nil {
		r.Items = make([]RatingScaleItemResponse, 0, len(s.Items))
		for _, item := range s.Items {
			r.Items = append(r.Items, item.ToResponse())
		}
	}
	return r
}

func (i *CompetencyRatingScaleItem) ToResponse() RatingScaleItemResponse {
	r := RatingScaleItemResponse{
		ID:        i.ID.String(),
		ScaleID:   i.ScaleID.String(),
		Value:     i.Value,
		Label:     i.Label,
		Weight:    i.Weight,
		SortOrder: i.SortOrder,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
	if i.Description != nil {
		r.Description = *i.Description
	}
	return r
}

func (t *CompetencyAssessmentTemplate) ToResponse() AssessmentTemplateResponse {
	r := AssessmentTemplateResponse{
		ID:        t.ID.String(),
		Name:      t.Name,
		Code:      t.Code,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
	if t.Description != nil {
		r.Description = *t.Description
	}
	if t.ScaleID != nil {
		r.ScaleID = t.ScaleID.String()
	}
	if t.Competencies != nil {
		r.Competencies = make([]TemplateCompetencyResponse, 0, len(t.Competencies))
		for _, c := range t.Competencies {
			resp := TemplateCompetencyResponse{
				ID:           c.ID.String(),
				TemplateID:   c.TemplateID.String(),
				CompetencyID: c.CompetencyID.String(),
				Weight:       c.Weight,
				SortOrder:    c.SortOrder,
			}
			if c.RequiredLevel != nil {
				resp.RequiredLevel = *c.RequiredLevel
			}
			if c.Competency != nil {
				resp.CompetencyName = c.Competency.Name
			}
			r.Competencies = append(r.Competencies, resp)
		}
	}
	if t.RaterTypes != nil {
		r.RaterTypes = make([]TemplateRaterTypeResponse, 0, len(t.RaterTypes))
		for _, rt := range t.RaterTypes {
			resp := TemplateRaterTypeResponse{
				ID:         rt.ID.String(),
				TemplateID: rt.TemplateID.String(),
				RaterType:  rt.RaterType,
				Weight:     rt.Weight,
				MinRater:   rt.MinRater,
				Required:   rt.Required,
				Anonymous:  rt.Anonymous,
			}
			if rt.MaxRater != nil {
				resp.MaxRater = *rt.MaxRater
			}
			r.RaterTypes = append(r.RaterTypes, resp)
		}
	}
	return r
}

func (i *CompetencyIndicator) ToResponse() IndicatorResponse {
	r := IndicatorResponse{
		ID:           i.ID.String(),
		CompetencyID: i.CompetencyID.String(),
		Statement:    i.Statement,
		Status:       i.Status,
		SortOrder:    i.SortOrder,
		CreatedAt:    i.CreatedAt,
		UpdatedAt:    i.UpdatedAt,
	}
	if i.Code != nil {
		r.Code = *i.Code
	}
	if i.Description != nil {
		r.Description = *i.Description
	}
	if i.Competency != nil {
		r.CompetencyName = i.Competency.Name
	}
	return r
}

func (ti *CompetencyAssessmentTemplateIndicator) ToResponse() TemplateIndicatorResponse {
	r := TemplateIndicatorResponse{
		ID:          ti.ID.String(),
		TemplateID:  ti.TemplateID.String(),
		IndicatorID: ti.IndicatorID.String(),
		Weight:      ti.Weight,
		SortOrder:   ti.SortOrder,
	}
	if ti.Indicator != nil {
		r.Statement = ti.Indicator.Statement
	}
	return r
}

// =========================================================================
// Request DTOs — Rater Assignment
// =========================================================================

type RaterAssignmentRequest struct {
	RaterEmployeeID string `json:"rater_employee_id" binding:"required"`
	RaterType       string `json:"rater_type" binding:"required,oneof=self superior peer subordinate other"`
	Weight          *float64 `json:"weight"`
}

type AssignRatersRequest struct {
	Raters []RaterAssignmentRequest `json:"raters" binding:"required,min=1"`
}

// =========================================================================
// Request DTOs — Assessment Response
// =========================================================================

type SaveResponseRequest struct {
	IndicatorID string  `json:"indicator_id" binding:"required"`
	RatingValue int     `json:"rating_value" binding:"required"`
	Comment     *string `json:"comment"`
}

type SaveResponsesRequest struct {
	Responses []SaveResponseRequest `json:"responses" binding:"required,min=1"`
}

// =========================================================================
// Response DTOs — Rater & Assessment
// =========================================================================

type RaterResponse struct {
	ID                      string    `json:"id"`
	CompetencyEventTargetID string    `json:"competency_event_target_id"`
	CompetencyEventID       string    `json:"competency_event_id,omitempty"`
	RaterEmployeeID         string    `json:"rater_employee_id"`
	RaterEmployeeName       string    `json:"rater_employee_name,omitempty"`
	SubjectEmployeeID       string    `json:"subject_employee_id,omitempty"`
	SubjectEmployeeName     string    `json:"subject_employee_name,omitempty"`
	RaterType               string    `json:"rater_type"`
	Weight                  float64   `json:"weight"`
	Status                  string    `json:"status"`
	AssignedAt              *time.Time `json:"assigned_at,omitempty"`
	SubmittedAt             *time.Time `json:"submitted_at,omitempty"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type AssessmentResponseDTO struct {
	ID           string    `json:"id"`
	RaterID      string    `json:"rater_id"`
	IndicatorID  string    `json:"indicator_id"`
	Statement    string    `json:"statement,omitempty"`
	RatingValue  int       `json:"rating_value"`
	Comment      string    `json:"comment,omitempty"`
	SubmittedAt  *time.Time `json:"submitted_at,omitempty"`
}

type AssessmentDetailDTO struct {
	Rater      RaterResponse                  `json:"rater"`
	Target     *CompetencyEventTargetResponse `json:"target,omitempty"`
	Indicators []TemplateIndicatorResponse    `json:"indicators,omitempty"`
	Responses  []AssessmentResponseDTO        `json:"responses,omitempty"`
	Scale      *RatingScaleResponse           `json:"scale,omitempty"`
}

// =========================================================================
// Converters — Rater & Response
// =========================================================================

func (r *CompetencyAssessmentRater) ToResponse() RaterResponse {
	resp := RaterResponse{
		ID:                      r.ID.String(),
		CompetencyEventTargetID: r.CompetencyEventTargetID.String(),
		RaterEmployeeID:         r.RaterEmployeeID.String(),
		RaterType:               r.RaterType,
		Weight:                  r.Weight,
		Status:                  r.Status,
		AssignedAt:              r.AssignedAt,
		SubmittedAt:             r.SubmittedAt,
		CreatedAt:               r.CreatedAt,
		UpdatedAt:               r.UpdatedAt,
	}
	if r.Target != nil {
		if r.Target.EmployeeID != nil {
			resp.SubjectEmployeeID = r.Target.EmployeeID.String()
		}
		resp.CompetencyEventID = r.Target.CompetencyEventID.String()
	}
	return resp
}

func (r *CompetencyAssessmentResponse) ToDTO() AssessmentResponseDTO {
	d := AssessmentResponseDTO{
		ID:          r.ID.String(),
		RaterID:     r.RaterID.String(),
		IndicatorID: r.IndicatorID.String(),
		RatingValue: r.RatingValue,
		SubmittedAt: r.SubmittedAt,
	}
	if r.Comment != nil {
		d.Comment = *r.Comment
	}
	if r.Indicator != nil {
		d.Statement = r.Indicator.Statement
	}
	return d
}

// =========================================================================
// Response DTOs — Manager Assessment (berbasis daftar bawahan)
// =========================================================================

type ManagerAssessmentItem struct {
	EmployeeID        string     `json:"employee_id"`
	EmployeeName      string     `json:"employee_name,omitempty"`
	TargetID          string     `json:"target_id"`
	CompetencyEventID string     `json:"competency_event_id"`
	// RaterID kosong bila manager belum di-assign sebagai superior rater
	// untuk bawahan tsb — frontend menawarkan tombol assign sendiri.
	RaterID     string     `json:"rater_id,omitempty"`
	RaterStatus string     `json:"rater_status,omitempty"`
	AssignedAt  *time.Time `json:"assigned_at,omitempty"`
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
}

// =========================================================================
// Response DTOs — Gap Analysis (§17/§18)
// =========================================================================

type GapItem struct {
	CompetencyID   string  `json:"competency_id"`
	CompetencyName string  `json:"competency_name,omitempty"`
	RequiredLevel  float64 `json:"required_level"`
	Score          float64 `json:"score"`
	Gap            float64 `json:"gap"`
	WeightedGap    float64 `json:"weighted_gap"`
}

type GapAnalysisResponse struct {
	TargetID         string    `json:"target_id"`
	EmployeeID       string    `json:"employee_id"`
	OverallScore     float64   `json:"overall_score"`
	TotalGap         float64   `json:"total_gap"`
	SelfScore        float64   `json:"self_score"`
	OthersScore      float64   `json:"others_score"`
	PerceptionGap    float64   `json:"perception_gap"`
	Strengths        []GapItem `json:"strengths"`
	DevelopmentAreas []GapItem `json:"development_areas"`
}
