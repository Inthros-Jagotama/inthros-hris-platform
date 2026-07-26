package pkgmgr

import "time"

// CreatePackageRequest untuk membuat paket baru.
type CreatePackageRequest struct {
	Name        string                `json:"name" binding:"required,min=3,max=255"`
	Slug        string                `json:"slug" binding:"required,min=2,max=100"`
	Description string                `json:"description,omitempty"`
	Price       float64               `json:"price" binding:"omitempty,min=0"`
	SortOrder   int                   `json:"sort_order,omitempty"`
	Modules     []PackageModuleInput  `json:"modules,omitempty"`
}

// UpdatePackageRequest untuk update paket.
type UpdatePackageRequest struct {
	Name        *string               `json:"name,omitempty" binding:"omitempty,min=3,max=255"`
	Slug        *string               `json:"slug,omitempty" binding:"omitempty,min=2,max=100"`
	Description *string               `json:"description,omitempty"`
	Price       *float64              `json:"price,omitempty" binding:"omitempty,min=0"`
	SortOrder   *int                  `json:"sort_order,omitempty"`
	Status      *string               `json:"status,omitempty" binding:"omitempty,oneof=draft published archived"`
	IsPublic    *bool                 `json:"is_public,omitempty"`
	Modules     []PackageModuleInput  `json:"modules,omitempty"`
}

// PackageModuleInput untuk input modul dalam paket.
type PackageModuleInput struct {
	ModuleID    string `json:"module_id" binding:"required"`
	IsMandatory bool   `json:"is_mandatory"`
	SortOrder   int    `json:"sort_order,omitempty"`
}

// ModuleDependency untuk return info dependensi.
type ModuleDependency struct {
	ModuleID   string `json:"module_id"`
	ModuleName string `json:"module_name"`
	DependsOn  string `json:"depends_on"`
	Resolved   bool   `json:"resolved"`
}

// PackageResponse untuk response data paket.
type PackageResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Slug        string                 `json:"slug"`
	Description string                 `json:"description,omitempty"`
	Price       float64                `json:"price"`
	Status      string                 `json:"status"`
	IsPublic    bool                   `json:"is_public"`
	SortOrder   int                    `json:"sort_order"`
	ModuleCount int                    `json:"module_count"`
	Modules     []PackageModuleResponse `json:"modules,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// PackageModuleResponse untuk response relasi modul dalam paket.
type PackageModuleResponse struct {
	ModuleID    string `json:"module_id"`
	ModuleName  string `json:"module_name"`
	ModuleSlug  string `json:"module_slug"`
	IsMandatory bool   `json:"is_mandatory"`
	SortOrder   int    `json:"sort_order"`
}

// PublicPackageResponse untuk halaman public (hanya published).
type PublicPackageResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Slug        string                 `json:"slug"`
	Description string                 `json:"description,omitempty"`
	Price       float64                `json:"price"`
	SortOrder   int                    `json:"sort_order"`
	ModuleCount int                    `json:"module_count"`
	Modules     []PackageModuleResponse `json:"modules,omitempty"`
}

// ToResponse mengonversi Package ke PackageResponse.
func (p *Package) ToResponse() PackageResponse {
	resp := PackageResponse{
		ID:          p.ID.String(),
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Price:       p.Price,
		Status:      p.Status,
		IsPublic:    p.IsPublic,
		SortOrder:   p.SortOrder,
		ModuleCount: len(p.Modules),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	if len(p.Modules) > 0 {
		var modules []PackageModuleResponse
		for _, m := range p.Modules {
			modules = append(modules, PackageModuleResponse{
				ModuleID:    m.ModuleID.String(),
				ModuleName:  m.ModuleName,
				ModuleSlug:  m.ModuleSlug,
				IsMandatory: m.IsMandatory,
				SortOrder:   m.SortOrder,
			})
		}
		resp.Modules = modules
	}

	return resp
}

// ToPublicResponse mengonversi Package ke PublicPackageResponse.
func (p *Package) ToPublicResponse() PublicPackageResponse {
	resp := PublicPackageResponse{
		ID:          p.ID.String(),
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Price:       p.Price,
		SortOrder:   p.SortOrder,
		ModuleCount: len(p.Modules),
	}

	if len(p.Modules) > 0 {
		var modules []PackageModuleResponse
		for _, m := range p.Modules {
			modules = append(modules, PackageModuleResponse{
				ModuleID:    m.ModuleID.String(),
				ModuleName:  m.ModuleName,
				ModuleSlug:  m.ModuleSlug,
				IsMandatory: m.IsMandatory,
				SortOrder:   m.SortOrder,
			})
		}
		resp.Modules = modules
	}

	return resp
}
