package setting

import "time"

// ── Zone DTOs ──
type CreateZoneRequest struct {
	Code      string `json:"code" binding:"required,max=20"`
	Name      string `json:"name" binding:"required,max=255"`
	Region    string `json:"region,omitempty" binding:"max=100"`
	IsActive  *bool  `json:"is_active,omitempty"`
	SortOrder int    `json:"sort_order,omitempty"`
}
type UpdateZoneRequest struct {
	Code      *string `json:"code,omitempty" binding:"omitempty,max=20"`
	Name      *string `json:"name,omitempty" binding:"omitempty,max=255"`
	Region    *string `json:"region,omitempty" binding:"omitempty,max=100"`
	IsActive  *bool   `json:"is_active,omitempty"`
	SortOrder *int    `json:"sort_order,omitempty"`
}
type ZoneResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Region    string    `json:"region,omitempty"`
	IsActive  bool      `json:"is_active"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
func (z *Zone) ToResponse() ZoneResponse {
	return ZoneResponse{ID: z.ID.String(), Code: z.Code, Name: z.Name, Region: z.Region, IsActive: z.IsActive, SortOrder: z.SortOrder, CreatedAt: z.CreatedAt, UpdatedAt: z.UpdatedAt}
}

// ── Province DTOs ──
type CreateProvinceRequest struct {
	Code string `json:"code" binding:"required,max=10"`
	Name string `json:"name" binding:"required,max=255"`
}
type UpdateProvinceRequest struct {
	Code *string `json:"code,omitempty" binding:"omitempty,max=10"`
	Name *string `json:"name,omitempty" binding:"omitempty,max=255"`
}
type ProvinceResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
func (p *Province) ToResponse() ProvinceResponse {
	return ProvinceResponse{ID: p.ID.String(), Code: p.Code, Name: p.Name, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}
}

// ── Regency DTOs ──
type CreateRegencyRequest struct {
	Code       string `json:"code" binding:"required,max=10"`
	Name       string `json:"name" binding:"required,max=255"`
	ProvinceID string `json:"province_id" binding:"required"`
}
type UpdateRegencyRequest struct {
	Code       *string `json:"code,omitempty" binding:"omitempty,max=10"`
	Name       *string `json:"name,omitempty" binding:"omitempty,max=255"`
	ProvinceID *string `json:"province_id,omitempty"`
}
type RegencyResponse struct {
	ID         string    `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	ProvinceID string    `json:"province_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
func (r *Regency) ToResponse() RegencyResponse {
	return RegencyResponse{ID: r.ID.String(), Code: r.Code, Name: r.Name, ProvinceID: r.ProvinceID.String(), CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt}
}

// ── District DTOs ──
type CreateDistrictRequest struct {
	Code      string `json:"code" binding:"required,max=10"`
	Name      string `json:"name" binding:"required,max=255"`
	RegencyID string `json:"regency_id" binding:"required"`
}
type UpdateDistrictRequest struct {
	Code      *string `json:"code,omitempty" binding:"omitempty,max=10"`
	Name      *string `json:"name,omitempty" binding:"omitempty,max=255"`
	RegencyID *string `json:"regency_id,omitempty"`
}
type DistrictResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	RegencyID string    `json:"regency_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
func (d *District) ToResponse() DistrictResponse {
	return DistrictResponse{ID: d.ID.String(), Code: d.Code, Name: d.Name, RegencyID: d.RegencyID.String(), CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt}
}

// ── Village DTOs ──
type CreateVillageRequest struct {
	Code       string `json:"code" binding:"required,max=15"`
	Name       string `json:"name" binding:"required,max=255"`
	DistrictID string `json:"district_id" binding:"required"`
}
type UpdateVillageRequest struct {
	Code       *string `json:"code,omitempty" binding:"omitempty,max=15"`
	Name       *string `json:"name,omitempty" binding:"omitempty,max=255"`
	DistrictID *string `json:"district_id,omitempty"`
}
type VillageResponse struct {
	ID         string    `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	DistrictID string    `json:"district_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
func (v *Village) ToResponse() VillageResponse {
	return VillageResponse{ID: v.ID.String(), Code: v.Code, Name: v.Name, DistrictID: v.DistrictID.String(), CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt}
}

// ── Education DTOs ──
type CreateEducationRequest struct {
	Code      string `json:"code" binding:"required,max=20"`
	Name      string `json:"name" binding:"required,max=255"`
	SortOrder int    `json:"sort_order,omitempty"`
}
type UpdateEducationRequest struct {
	Code      *string `json:"code,omitempty" binding:"omitempty,max=20"`
	Name      *string `json:"name,omitempty" binding:"omitempty,max=255"`
	SortOrder *int    `json:"sort_order,omitempty"`
}
type EducationResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
func (e *Education) ToResponse() EducationResponse {
	return EducationResponse{ID: e.ID.String(), Code: e.Code, Name: e.Name, SortOrder: e.SortOrder, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt}
}

// ── Religion DTOs ──
type CreateReligionRequest struct {
	Code      string `json:"code" binding:"required,max=20"`
	Name      string `json:"name" binding:"required,max=255"`
	SortOrder int    `json:"sort_order,omitempty"`
}
type UpdateReligionRequest struct {
	Code      *string `json:"code,omitempty" binding:"omitempty,max=20"`
	Name      *string `json:"name,omitempty" binding:"omitempty,max=255"`
	SortOrder *int    `json:"sort_order,omitempty"`
}
type ReligionResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
func (r *Religion) ToResponse() ReligionResponse {
	return ReligionResponse{ID: r.ID.String(), Code: r.Code, Name: r.Name, SortOrder: r.SortOrder, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt}
}

// ── MaritalStatus DTOs ──
type CreateMaritalStatusRequest struct {
	Code      string `json:"code" binding:"required,max=20"`
	Name      string `json:"name" binding:"required,max=255"`
	SortOrder int    `json:"sort_order,omitempty"`
}
type UpdateMaritalStatusRequest struct {
	Code      *string `json:"code,omitempty" binding:"omitempty,max=20"`
	Name      *string `json:"name,omitempty" binding:"omitempty,max=255"`
	SortOrder *int    `json:"sort_order,omitempty"`
}
type MaritalStatusResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
func (m *MaritalStatus) ToResponse() MaritalStatusResponse {
	return MaritalStatusResponse{ID: m.ID.String(), Code: m.Code, Name: m.Name, SortOrder: m.SortOrder, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}

// ── Paginated Responses ──
type ZonePaginatedResponse struct {
	Success    bool           `json:"success"`
	Data       []ZoneResponse `json:"data"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	Total      int64          `json:"total"`
	TotalPages int            `json:"total_pages"`
}
type ProvincePaginatedResponse struct {
	Success    bool               `json:"success"`
	Data       []ProvinceResponse `json:"data"`
	Page       int                `json:"page"`
	PerPage    int                `json:"per_page"`
	Total      int64              `json:"total"`
	TotalPages int                `json:"total_pages"`
}
type RegencyPaginatedResponse struct {
	Success    bool              `json:"success"`
	Data       []RegencyResponse `json:"data"`
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
	Total      int64             `json:"total"`
	TotalPages int               `json:"total_pages"`
}
type DistrictPaginatedResponse struct {
	Success    bool                `json:"success"`
	Data       []DistrictResponse  `json:"data"`
	Page       int                 `json:"page"`
	PerPage    int                 `json:"per_page"`
	Total      int64               `json:"total"`
	TotalPages int                 `json:"total_pages"`
}
type VillagePaginatedResponse struct {
	Success    bool              `json:"success"`
	Data       []VillageResponse `json:"data"`
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
	Total      int64             `json:"total"`
	TotalPages int               `json:"total_pages"`
}
type EducationPaginatedResponse struct {
	Success    bool                `json:"success"`
	Data       []EducationResponse `json:"data"`
	Page       int                 `json:"page"`
	PerPage    int                 `json:"per_page"`
	Total      int64               `json:"total"`
	TotalPages int                 `json:"total_pages"`
}
type ReligionPaginatedResponse struct {
	Success    bool               `json:"success"`
	Data       []ReligionResponse `json:"data"`
	Page       int                `json:"page"`
	PerPage    int                `json:"per_page"`
	Total      int64              `json:"total"`
	TotalPages int                `json:"total_pages"`
}
type MaritalStatusPaginatedResponse struct {
	Success    bool                    `json:"success"`
	Data       []MaritalStatusResponse `json:"data"`
	Page       int                     `json:"page"`
	PerPage    int                     `json:"per_page"`
	Total      int64                   `json:"total"`
	TotalPages int                     `json:"total_pages"`
}

// ── RelationshipType DTOs ──
type CreateRelationshipTypeRequest struct {
	Code      string `json:"code" binding:"required,max=20"`
	Name      string `json:"name" binding:"required,max=255"`
	SortOrder int    `json:"sort_order,omitempty"`
}
type UpdateRelationshipTypeRequest struct {
	Code      *string `json:"code,omitempty" binding:"omitempty,max=20"`
	Name      *string `json:"name,omitempty" binding:"omitempty,max=255"`
	SortOrder *int    `json:"sort_order,omitempty"`
}
type RelationshipTypeResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
func (rt *RelationshipType) ToResponse() RelationshipTypeResponse {
	return RelationshipTypeResponse{ID: rt.ID.String(), Code: rt.Code, Name: rt.Name, SortOrder: rt.SortOrder, CreatedAt: rt.CreatedAt, UpdatedAt: rt.UpdatedAt}
}

// ── EmploymentStatus DTOs ──
type CreateEmploymentStatusRequest struct {
	Code      string `json:"code" binding:"required,max=20"`
	Name      string `json:"name" binding:"required,max=255"`
	SortOrder int    `json:"sort_order,omitempty"`
}
type UpdateEmploymentStatusRequest struct {
	Code      *string `json:"code,omitempty" binding:"omitempty,max=20"`
	Name      *string `json:"name,omitempty" binding:"omitempty,max=255"`
	SortOrder *int    `json:"sort_order,omitempty"`
}
type EmploymentStatusResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
func (es *EmploymentStatus) ToResponse() EmploymentStatusResponse {
	return EmploymentStatusResponse{ID: es.ID.String(), Code: es.Code, Name: es.Name, SortOrder: es.SortOrder, CreatedAt: es.CreatedAt, UpdatedAt: es.UpdatedAt}
}

// ── Bank DTOs ──
type CreateBankRequest struct {
	Code      string `json:"code" binding:"required,max=20"`
	Name      string `json:"name" binding:"required,max=255"`
	SortOrder int    `json:"sort_order,omitempty"`
}
type UpdateBankRequest struct {
	Code      *string `json:"code,omitempty" binding:"omitempty,max=20"`
	Name      *string `json:"name,omitempty" binding:"omitempty,max=255"`
	SortOrder *int    `json:"sort_order,omitempty"`
}
type BankResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
func (b *Bank) ToResponse() BankResponse {
	return BankResponse{ID: b.ID.String(), Code: b.Code, Name: b.Name, SortOrder: b.SortOrder, CreatedAt: b.CreatedAt, UpdatedAt: b.UpdatedAt}
}

// ── Nationality DTOs ──
type CreateNationalityRequest struct {
	Code      string `json:"code" binding:"required,max=20"`
	Name      string `json:"name" binding:"required,max=255"`
	SortOrder int    `json:"sort_order,omitempty"`
}
type UpdateNationalityRequest struct {
	Code      *string `json:"code,omitempty" binding:"omitempty,max=20"`
	Name      *string `json:"name,omitempty" binding:"omitempty,max=255"`
	SortOrder *int    `json:"sort_order,omitempty"`
}
type NationalityResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
func (n *Nationality) ToResponse() NationalityResponse {
	return NationalityResponse{ID: n.ID.String(), Code: n.Code, Name: n.Name, SortOrder: n.SortOrder, CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt}
}

// ── Paginated Responses ──
type RelationshipTypePaginatedResponse struct {
	Success    bool                       `json:"success"`
	Data       []RelationshipTypeResponse `json:"data"`
	Page       int                        `json:"page"`
	PerPage    int                        `json:"per_page"`
	Total      int64                      `json:"total"`
	TotalPages int                        `json:"total_pages"`
}
type BankPaginatedResponse struct {
	Success    bool           `json:"success"`
	Data       []BankResponse `json:"data"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	Total      int64          `json:"total"`
	TotalPages int            `json:"total_pages"`
}

type NationalityPaginatedResponse struct {
	Success    bool                 `json:"success"`
	Data       []NationalityResponse `json:"data"`
	Page       int                  `json:"page"`
	PerPage    int                  `json:"per_page"`
	Total      int64                `json:"total"`
	TotalPages int                  `json:"total_pages"`
}

type EmploymentStatusPaginatedResponse struct {
	Success    bool                        `json:"success"`
	Data       []EmploymentStatusResponse  `json:"data"`
	Page       int                         `json:"page"`
	PerPage    int                         `json:"per_page"`
	Total      int64                       `json:"total"`
	TotalPages int                         `json:"total_pages"`
}

// ── JobFamily DTOs ──
type CreateJobFamilyRequest struct {
	Code        string `json:"code" binding:"required,max=20"`
	Name        string `json:"name" binding:"required,max=255"`
	Description string `json:"description,omitempty"`
	SortOrder   int    `json:"sort_order,omitempty"`
}
type UpdateJobFamilyRequest struct {
	Code        *string `json:"code,omitempty" binding:"omitempty,max=20"`
	Name        *string `json:"name,omitempty" binding:"omitempty,max=255"`
	Description *string `json:"description,omitempty"`
	SortOrder   *int    `json:"sort_order,omitempty"`
}
type JobFamilyResponse struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
func (jf *JobFamily) ToResponse() JobFamilyResponse {
	return JobFamilyResponse{ID: jf.ID.String(), Code: jf.Code, Name: jf.Name, Description: jf.Description, SortOrder: jf.SortOrder, CreatedAt: jf.CreatedAt, UpdatedAt: jf.UpdatedAt}
}

// ── SalaryGrade DTOs ──
type CreateSalaryGradeRequest struct {
	Code        string  `json:"code" binding:"required,max=20"`
	Name        string  `json:"name" binding:"required,max=255"`
	Description string  `json:"description,omitempty"`
	MinAmount   float64 `json:"min_amount,omitempty"`
	MaxAmount   float64 `json:"max_amount,omitempty"`
	SortOrder   int     `json:"sort_order,omitempty"`
}
type UpdateSalaryGradeRequest struct {
	Code        *string  `json:"code,omitempty" binding:"omitempty,max=20"`
	Name        *string  `json:"name,omitempty" binding:"omitempty,max=255"`
	Description *string  `json:"description,omitempty"`
	MinAmount   *float64 `json:"min_amount,omitempty"`
	MaxAmount   *float64 `json:"max_amount,omitempty"`
	SortOrder   *int     `json:"sort_order,omitempty"`
}
type SalaryGradeResponse struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	MinAmount   float64   `json:"min_amount,omitempty"`
	MaxAmount   float64   `json:"max_amount,omitempty"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
func (sg *SalaryGrade) ToResponse() SalaryGradeResponse {
	return SalaryGradeResponse{ID: sg.ID.String(), Code: sg.Code, Name: sg.Name, Description: sg.Description, MinAmount: sg.MinAmount, MaxAmount: sg.MaxAmount, SortOrder: sg.SortOrder, CreatedAt: sg.CreatedAt, UpdatedAt: sg.UpdatedAt}
}

// ── Paginated Responses ──
type JobFamilyPaginatedResponse struct {
	Success    bool                `json:"success"`
	Data       []JobFamilyResponse `json:"data"`
	Page       int                 `json:"page"`
	PerPage    int                 `json:"per_page"`
	Total      int64               `json:"total"`
	TotalPages int                 `json:"total_pages"`
}
type SalaryGradePaginatedResponse struct {
	Success    bool                  `json:"success"`
	Data       []SalaryGradeResponse `json:"data"`
	Page       int                   `json:"page"`
	PerPage    int                   `json:"per_page"`
	Total      int64                 `json:"total"`
	TotalPages int                   `json:"total_pages"`
}
