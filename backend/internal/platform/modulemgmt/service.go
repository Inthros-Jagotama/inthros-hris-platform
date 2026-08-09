package modulemgmt

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/cache"
	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/middleware"
)

// Service untuk business logic Module Management.
type Service struct {
	repo         *Repository
	dbManager    *database.Manager
	logger       *zap.Logger
	cacheManager *cache.Cache
}

// NewService membuat Service baru.
func NewService(repo *Repository, dbManager *database.Manager, logger *zap.Logger) *Service {
	return &Service{
		repo:      repo,
		dbManager: dbManager,
		logger:    logger,
	}
}

// SetCacheManager menginjeksi cache manager untuk invalidasi license cache
// saat modul diaktifkan/dinonaktifkan. Dipanggil dari main.go setelah cache siap.
func (s *Service) SetCacheManager(cm *cache.Cache) {
	s.cacheManager = cm
}

// invalidateLicenseCache menghapus license cache PlatformLicenseMiddleware
// untuk company agar perubahan modul langsung berlaku (tanpa menunggu TTL).
func (s *Service) invalidateLicenseCache(companyID string) {
	if s.cacheManager == nil {
		return
	}
	if err := s.cacheManager.Invalidate(context.Background(), middleware.LicenseCacheKey(companyID)); err != nil {
		s.logger.Warn("Failed to invalidate license cache",
			zap.String("company_id", companyID),
			zap.Error(err),
		)
	}
}

// CreateModule mendaftarkan modul baru di platform.
func (s *Service) CreateModule(req CreateModuleRequest) (*ModuleResponse, error) {
	// Cek duplikasi slug
	if existing, _ := s.repo.FindBySlug(req.Slug); existing != nil {
		return nil, fmt.Errorf("module with slug '%s' already exists", req.Slug)
	}

	desc := req.Description
	dep := req.DependsOn
	module := &PlatformModule{
		Name:    req.Name,
		Slug:    req.Slug,
		Version: req.Version,
		IsCore:  req.IsCore,
	}
	if desc != "" {
		module.Description = &desc
	}
	if dep != "" {
		module.DependsOn = &dep
	}

	if err := s.repo.Create(module); err != nil {
		return nil, err
	}

	s.logger.Info("Module registered",
		zap.String("module_id", module.ID),
		zap.String("name", module.Name),
		zap.String("slug", module.Slug),
	)

	response := module.ToResponse()
	return &response, nil
}

// GetModule mengembalikan modul berdasarkan ID.
func (s *Service) GetModule(id string) (*ModuleResponse, error) {
	module, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := module.ToResponse()
	return &response, nil
}

// ListModules mengembalikan daftar semua modul dengan pagination dan optional filter module_type.
// moduleType bisa "platform", "tenant", atau "" untuk semua.
func (s *Service) ListModules(page, perPage int, moduleType, search string) (*PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	modules, total, err := s.repo.FindAll(page, perPage, moduleType, search)
	if err != nil {
		return nil, err
	}

	var responses []ModuleResponse
	for _, m := range modules {
		responses = append(responses, m.ToResponse())
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

// UpdateModule mengupdate modul.
func (s *Service) UpdateModule(id string, req UpdateModuleRequest) (*ModuleResponse, error) {
	module, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		module.Name = *req.Name
	}
	if req.Version != nil {
		module.Version = *req.Version
	}
	if req.Description != nil {
		module.Description = req.Description
	}
	if req.IsCore != nil {
		module.IsCore = *req.IsCore
	}
	if req.DependsOn != nil {
		module.DependsOn = req.DependsOn
	}

	if err := s.repo.Update(module); err != nil {
		return nil, err
	}

	response := module.ToResponse()
	return &response, nil
}

// ActivateModule mengaktifkan modul untuk company tertentu.
func (s *Service) ActivateModule(moduleID, companyID string) (*CompanyModuleResponse, error) {
	cid, err := uuid.Parse(companyID)
	if err != nil {
		return nil, fmt.Errorf("invalid company id: %w", err)
	}
	mid, err := uuid.Parse(moduleID)
	if err != nil {
		return nil, fmt.Errorf("invalid module id: %w", err)
	}

	// Cek apakah modul sudah terdaftar
	module, err := s.repo.FindByID(moduleID)
	if err != nil {
		return nil, fmt.Errorf("module not found: %w", err)
	}

	// Upsert company-module relation
	cm, err := s.repo.UpsertCompanyModule(cid, mid, true)
	if err != nil {
		return nil, err
	}

	s.invalidateLicenseCache(companyID)

	s.logger.Info("Module activated for company",
		zap.String("module_id", moduleID),
		zap.String("module_name", module.Name),
		zap.String("company_id", companyID),
	)

	return &CompanyModuleResponse{
		CompanyID:   cm.CompanyID.String(),
		ModuleID:    cm.ModuleID.String(),
		ModuleName:  module.Name,
		Enabled:     cm.Enabled,
		ActivatedAt: cm.ActivatedAt,
	}, nil
}

// DeactivateModule menonaktifkan modul untuk company tertentu.
func (s *Service) DeactivateModule(moduleID, companyID string) (*CompanyModuleResponse, error) {
	cid, err := uuid.Parse(companyID)
	if err != nil {
		return nil, fmt.Errorf("invalid company id: %w", err)
	}
	mid, err := uuid.Parse(moduleID)
	if err != nil {
		return nil, fmt.Errorf("invalid module id: %w", err)
	}

	module, err := s.repo.FindByID(moduleID)
	if err != nil {
		return nil, fmt.Errorf("module not found: %w", err)
	}

	cm, err := s.repo.UpsertCompanyModule(cid, mid, false)
	if err != nil {
		return nil, err
	}

	s.invalidateLicenseCache(companyID)

	s.logger.Info("Module deactivated for company",
		zap.String("module_id", moduleID),
		zap.String("module_name", module.Name),
		zap.String("company_id", companyID),
	)

	return &CompanyModuleResponse{
		CompanyID:   cm.CompanyID.String(),
		ModuleID:    cm.ModuleID.String(),
		ModuleName:  module.Name,
		Enabled:     cm.Enabled,
		ActivatedAt: cm.ActivatedAt,
	}, nil
}

// ListCompanyModules mengembalikan daftar modul yang terdaftar untuk company.
func (s *Service) ListCompanyModules(companyID string) ([]CompanyModuleResponse, error) {
	cid, err := uuid.Parse(companyID)
	if err != nil {
		return nil, fmt.Errorf("invalid company id: %w", err)
	}

	modules, err := s.repo.FindCompanyModules(cid)
	if err != nil {
		return nil, err
	}

	return modules, nil
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
