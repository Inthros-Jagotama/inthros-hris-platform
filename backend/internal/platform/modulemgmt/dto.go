package modulemgmt

import "time"

// CreateModuleRequest untuk mendaftarkan modul baru.
type CreateModuleRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=255"`
	Slug        string `json:"slug" binding:"required,min=2,max=100"`
	Version     string `json:"version" binding:"required,max=20"`
	Description string `json:"description,omitempty"`
	ModuleType  string `json:"module_type,omitempty" binding:"omitempty,oneof=platform tenant"`
	IsCore      bool   `json:"is_core,omitempty"`
	DependsOn   string `json:"depends_on,omitempty"`
}

// UpdateModuleRequest untuk update modul.
type UpdateModuleRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=3,max=255"`
	Version     *string `json:"version,omitempty" binding:"omitempty,max=20"`
	Description *string `json:"description,omitempty"`
	ModuleType  *string `json:"module_type,omitempty" binding:"omitempty,oneof=platform tenant"`
	IsCore      *bool   `json:"is_core,omitempty"`
	DependsOn   *string `json:"depends_on,omitempty"`
}

// ModuleResponse untuk response data modul.
type ModuleResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Version     string    `json:"version"`
	Description string    `json:"description,omitempty"`
	ModuleType  string    `json:"module_type"`
	IsCore      bool      `json:"is_core"`
	DependsOn   string    `json:"depends_on,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToggleModuleRequest untuk activate/deactivate modul untuk company.
type ToggleModuleRequest struct {
	CompanyID string `json:"company_id" binding:"required"`
}

// CompanyModuleResponse untuk response company-module association.
type CompanyModuleResponse struct {
	CompanyID   string     `json:"company_id"`
	ModuleID    string     `json:"module_id"`
	ModuleName  string     `json:"module_name"`
	ModuleSlug  string     `json:"module_slug"`
	Enabled     bool       `json:"enabled"`
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
}

// ToResponse mengonversi PlatformModule ke ModuleResponse.
func (m *PlatformModule) ToResponse() ModuleResponse {
	desc := ""
	dep := ""
	if m.Description != nil {
		desc = *m.Description
	}
	if m.DependsOn != nil {
		dep = *m.DependsOn
	}
	return ModuleResponse{
		ID:          m.ID,
		Name:        m.Name,
		Slug:        m.Slug,
		Version:     m.Version,
		Description: desc,
		ModuleType:  m.ModuleType,
		IsCore:      m.IsCore,
		DependsOn:   dep,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
