package pkgmgr

import (
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Service untuk business logic Package Management.
type Service struct {
	repo   *Repository
	logger *zap.Logger
}

// NewService membuat Service baru.
func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// CreatePackage membuat paket baru dengan validasi dependensi modul.
func (s *Service) CreatePackage(req CreatePackageRequest) (*PackageResponse, error) {
	// Cek duplikasi slug
	if existing, _ := s.repo.FindBySlug(req.Slug); existing != nil {
		return nil, fmt.Errorf("package with slug '%s' already exists", req.Slug)
	}

	pkg := &Package{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Price:       req.Price,
		Status:      string(PackageDraft),
		SortOrder:   req.SortOrder,
	}

	// Process modules
	for _, m := range req.Modules {
		mid, err := uuid.Parse(m.ModuleID)
		if err != nil {
			return nil, fmt.Errorf("invalid module_id '%s': %w", m.ModuleID, err)
		}

		// Validate module exists
		modInfo, err := s.repo.FindModuleInfo(mid)
		if err != nil {
			return nil, fmt.Errorf("module '%s' not found: %w", m.ModuleID, err)
		}

		pkg.Modules = append(pkg.Modules, PackageModule{
			ModuleID:    mid,
			IsMandatory: m.IsMandatory,
			SortOrder:   m.SortOrder,
			ModuleName:  modInfo.Name,
			ModuleSlug:  modInfo.Slug,
		})
	}

	// Validate dependencies
	if err := s.validateModuleDependencies(pkg.Modules); err != nil {
		return nil, fmt.Errorf("module dependency validation failed: %w", err)
	}

	// Sort modules
	sort.Slice(pkg.Modules, func(i, j int) bool {
		return pkg.Modules[i].SortOrder < pkg.Modules[j].SortOrder
	})

	if err := s.repo.Create(pkg); err != nil {
		return nil, err
	}

	// Re-fetch with modules loaded
	created, err := s.repo.FindByID(pkg.ID)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Package created",
		zap.String("package_id", created.ID.String()),
		zap.String("name", created.Name),
		zap.Float64("price", created.Price),
		zap.Int("modules", len(created.Modules)),
	)

	response := created.ToResponse()
	return &response, nil
}

// GetPackage mengembalikan paket berdasarkan ID.
func (s *Service) GetPackage(id string) (*PackageResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid package id: %w", err)
	}

	pkg, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	response := pkg.ToResponse()
	return &response, nil
}

// GetPackageBySlug mengembalikan paket berdasarkan slug.
func (s *Service) GetPackageBySlug(slug string) (*PackageResponse, error) {
	pkg, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}

	response := pkg.ToResponse()
	return &response, nil
}

// ListPackages mengembalikan daftar paket dengan pagination.
// moduleType opsional: filter paket yang mengandung modul dengan tipe tertentu ("platform" atau "tenant").
func (s *Service) ListPackages(page, perPage int, moduleType string) (*PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	packages, total, err := s.repo.FindAll(page, perPage, moduleType)
	if err != nil {
		return nil, err
	}

	var responses []PackageResponse
	for _, p := range packages {
		responses = append(responses, p.ToResponse())
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// ListPublishedPackages mengembalikan paket published untuk halaman public.
// moduleType opsional: filter paket yang mengandung modul dengan tipe tertentu ("platform" atau "tenant").
func (s *Service) ListPublishedPackages(moduleType string) ([]PublicPackageResponse, error) {
	packages, err := s.repo.FindPublished(moduleType)
	if err != nil {
		return nil, err
	}

	var responses []PublicPackageResponse
	for _, p := range packages {
		responses = append(responses, p.ToPublicResponse())
	}

	return responses, nil
}

// UpdatePackage mengupdate paket.
func (s *Service) UpdatePackage(id string, req UpdatePackageRequest) (*PackageResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid package id: %w", err)
	}

	pkg, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		pkg.Name = *req.Name
	}
	if req.Slug != nil {
		// Cek duplikasi slug jika berubah
		if *req.Slug != pkg.Slug {
			if existing, _ := s.repo.FindBySlug(*req.Slug); existing != nil {
				return nil, fmt.Errorf("package with slug '%s' already exists", *req.Slug)
			}
		}
		pkg.Slug = *req.Slug
	}
	if req.Description != nil {
		pkg.Description = *req.Description
	}
	if req.Price != nil {
		pkg.Price = *req.Price
	}
	if req.SortOrder != nil {
		pkg.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		pkg.Status = *req.Status
	}
	if req.IsPublic != nil {
		pkg.IsPublic = *req.IsPublic
	}

	// Process modules if provided
	if req.Modules != nil {
		var newModules []PackageModule
		for _, m := range req.Modules {
			mid, err := uuid.Parse(m.ModuleID)
			if err != nil {
				return nil, fmt.Errorf("invalid module_id '%s': %w", m.ModuleID, err)
			}

			modInfo, err := s.repo.FindModuleInfo(mid)
			if err != nil {
				return nil, fmt.Errorf("module '%s' not found: %w", m.ModuleID, err)
			}

			newModules = append(newModules, PackageModule{
				ModuleID:    mid,
				IsMandatory: m.IsMandatory,
				SortOrder:   m.SortOrder,
				ModuleName:  modInfo.Name,
				ModuleSlug:  modInfo.Slug,
			})
		}

		// Validate dependencies
		if err := s.validateModuleDependencies(newModules); err != nil {
			return nil, fmt.Errorf("module dependency validation failed: %w", err)
		}

		sort.Slice(newModules, func(i, j int) bool {
			return newModules[i].SortOrder < newModules[j].SortOrder
		})

		pkg.Modules = newModules
	}

	if err := s.repo.Update(pkg); err != nil {
		return nil, err
	}

	// Re-fetch
	updated, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Package updated",
		zap.String("package_id", updated.ID.String()),
		zap.String("name", updated.Name),
	)

	response := updated.ToResponse()
	return &response, nil
}

// DeletePackage soft-deletes a package.
func (s *Service) DeletePackage(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid package id: %w", err)
	}

	if err := s.repo.Delete(uid); err != nil {
		return err
	}

	s.logger.Info("Package deleted", zap.String("package_id", id))
	return nil
}

// PublishPackage mengubah status paket menjadi published.
func (s *Service) PublishPackage(id string) (*PackageResponse, error) {
	return s.updateStatus(id, string(PackagePublished))
}

// UnpublishPackage mengubah status paket menjadi draft.
func (s *Service) UnpublishPackage(id string) (*PackageResponse, error) {
	return s.updateStatus(id, string(PackageDraft))
}

// ValidatePackageDependencies memeriksa dependensi modul dalam paket.
func (s *Service) ValidatePackageDependencies(id string) ([]ModuleDependency, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid package id: %w", err)
	}

	pkg, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	// Collect all module IDs in this package
	moduleIDs := make(map[string]bool)
	for _, m := range pkg.Modules {
		moduleIDs[m.ModuleID.String()] = true
	}

	// Map module slug -> package module data for quick lookup
	moduleSlugMap := make(map[string]PackageModule)
	for _, m := range pkg.Modules {
		moduleSlugMap[m.ModuleSlug] = m
	}

	// Batch fetch all module info with depends_on
	allModules, err := s.repo.FindModuleInfoMap()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch module info: %w", err)
	}

	// For each module in the package, check its dependencies
	var deps []ModuleDependency
	for _, pm := range pkg.Modules {
		modID := pm.ModuleID.String()
		modInfo, ok := allModules[modID]
		if !ok {
			deps = append(deps, ModuleDependency{
				ModuleID:   modID,
				ModuleName: pm.ModuleName,
				DependsOn:  "[module data not found]",
				Resolved:   false,
			})
			continue
		}

		dependedSlugs := parseDependsOn(*modInfo.DependsOn)
		if len(dependedSlugs) == 0 {
			deps = append(deps, ModuleDependency{
				ModuleID:   modID,
				ModuleName: pm.ModuleName,
				DependsOn:  "(none)",
				Resolved:   true,
			})
			continue
		}

		// Check each dependency
		resolved := true
		var missing []string
		for _, depSlug := range dependedSlugs {
			if _, exists := moduleSlugMap[depSlug]; !exists {
				missing = append(missing, depSlug)
				resolved = false
			}
		}

		depInfo := fmt.Sprintf("needs: %s", *modInfo.DependsOn)
		deps = append(deps, ModuleDependency{
			ModuleID:   modID,
			ModuleName: pm.ModuleName,
			DependsOn:  depInfo,
			Resolved:   resolved,
		})

		if !resolved && len(missing) > 0 {
			s.logger.Warn("Unresolved dependencies in package",
				zap.String("package_id", id),
				zap.String("module", pm.ModuleName),
				zap.Strings("missing", missing),
			)
		}
	}

	return deps, nil
}

// updateStatus helper untuk mengubah status paket.
func (s *Service) updateStatus(id, status string) (*PackageResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid package id: %w", err)
	}

	pkg, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	// Validate: can only publish if has modules
	if status == string(PackagePublished) {
		if len(pkg.Modules) == 0 {
			return nil, fmt.Errorf("cannot publish package without modules")
		}
		// Validasi dependensi modul sebelum publish
		if err := s.validateModuleDependencies(pkg.Modules); err != nil {
			return nil, fmt.Errorf("cannot publish package: %w", err)
		}
	}

	pkg.Status = status

	// When publishing, auto-set is_public
	if status == string(PackagePublished) {
		pkg.IsPublic = true
	}

	if err := s.repo.Update(pkg); err != nil {
		return nil, err
	}

	updated, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Package status updated",
		zap.String("package_id", id),
		zap.String("status", status),
	)

	response := updated.ToResponse()
	return &response, nil
}

// validateModuleDependencies memeriksa dependensi antar modul dalam satu paket.
// Memverifikasi bahwa semua modul yang menjadi dependensi (DependsOn)
// sudah disertakan dalam paket yang sama.
func (s *Service) validateModuleDependencies(modules []PackageModule) error {
	// Build map of module slugs in this package for dependency lookup
	slugMap := make(map[string]PackageModule)
	idMap := make(map[string]bool)
	for _, m := range modules {
		slugMap[m.ModuleSlug] = m
		id := m.ModuleID.String()
		if idMap[id] {
			return fmt.Errorf("duplicate module in package: %s (%s)", m.ModuleName, id)
		}
		idMap[id] = true
	}

	// Batch fetch module info with depends_on
	allModuleInfo, err := s.repo.FindModuleInfoMap()
	if err != nil {
		s.logger.Warn("Could not fetch module dependency info, skipping validation", zap.Error(err))
		return nil
	}

	// For each module in the package, check its dependencies
	var missingDeps []string
	for _, m := range modules {
		modID := m.ModuleID.String()
		info, ok := allModuleInfo[modID]
		if !ok {
			// Module info not found in platform DB — skip
			continue
		}

		depSlugs := parseDependsOn(*info.DependsOn)
		for _, dep := range depSlugs {
			if _, exists := slugMap[dep]; !exists {
				missingDeps = append(missingDeps,
					fmt.Sprintf("'%s' requires module '%s' which is not in this package",
						m.ModuleName, dep))
			}
		}
	}

	if len(missingDeps) > 0 {
		return fmt.Errorf("missing module dependencies:\n  - %s", strings.Join(missingDeps, "\n  - "))
	}

	return nil
}

// PaginatedResponse untuk response pagination.
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}
