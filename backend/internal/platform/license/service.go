package license

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/database"
)

// ── PackageModuleManager ──

// PackageModuleManager mendefinisikan interface untuk mengaktifkan/menonaktifkan
// modul dari suatu package saat lisensi dibuat/diupdate.
// Diimplementasikan oleh adapter di main.go untuk menghindari circular dependency.
type PackageModuleManager interface {
	// ActivatePackageModules mengaktifkan semua modul dalam package untuk company.
	// Mengembalikan daftar nama modul yang berhasil diaktifkan.
	ActivatePackageModules(packageID, companyID string) ([]string, error)

	// DeactivatePackageModules menonaktifkan semua modul dalam package untuk company.
	// Mengembalikan daftar nama modul yang berhasil dinonaktifkan.
	DeactivatePackageModules(packageID, companyID string) ([]string, error)
}

// ── Service ──

// Service untuk business logic License Management.
type Service struct {
	repo          *Repository
	dbManager     *database.Manager
	moduleManager PackageModuleManager
	logger        *zap.Logger
}

// NewService membuat Service baru.
// moduleManager opsional — jika nil, auto-activation module tidak terjadi.
func NewService(repo *Repository, dbManager *database.Manager, moduleManager PackageModuleManager, logger *zap.Logger) *Service {
	return &Service{
		repo:          repo,
		dbManager:     dbManager,
		moduleManager: moduleManager,
		logger:        logger,
	}
}

// CreateLicense membuat lisensi baru untuk company.
// Jika PackageID disertakan, modul-modul dalam package akan auto-diaktifkan.
func (s *Service) CreateLicense(req CreateLicenseRequest) (*LicenseResponse, error) {
	cid, err := uuid.Parse(req.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("invalid company_id: %w", err)
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format (use YYYY-MM-DD): %w", err)
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format (use YYYY-MM-DD): %w", err)
	}

	if endDate.Before(startDate) {
		return nil, fmt.Errorf("end_date must be after start_date")
	}

	maxEmployees := req.MaxEmployees
	if maxEmployees <= 0 {
		switch req.PlanType {
		case "trial":
			maxEmployees = 10
		case "subscription":
			maxEmployees = 0 // unlimited
		case "basic":
			maxEmployees = 50
		case "pro":
			maxEmployees = 200
		case "enterprise":
			maxEmployees = 0 // unlimited
		}
	}

	maxModules := req.MaxModules
	if maxModules <= 0 {
		switch req.PlanType {
		case "trial":
			maxModules = 3
		case "subscription":
			maxModules = 0 // unlimited
		case "basic":
			maxModules = 8
		case "pro":
			maxModules = 15
		case "enterprise":
			maxModules = 0 // unlimited
		}
	}

	license := &License{
		CompanyID:    cid,
		LicenseKey:   xid.New().String(),
		PlanType:     req.PlanType,
		MaxEmployees: maxEmployees,
		MaxModules:   maxModules,
		StartDate:    startDate,
		EndDate:      endDate,
		Status:       string(LicenseActive),
	}

	// Optional: associate with a package
	if req.PackageID != "" {
		pid, err := uuid.Parse(req.PackageID)
		if err != nil {
			return nil, fmt.Errorf("invalid package_id: %w", err)
		}

		pkgName, err := s.repo.FindPackageNameByID(pid)
		if err != nil {
			return nil, fmt.Errorf("package not found: %w", err)
		}

		license.PackageID = &pid
		license.PackageName = pkgName

		s.logger.Info("License associated with package",
			zap.String("package_id", req.PackageID),
			zap.String("package_name", pkgName),
		)
	}

	if err := s.repo.Create(license); err != nil {
		return nil, err
	}

	s.logger.Info("License created",
		zap.String("license_id", license.ID.String()),
		zap.String("company_id", req.CompanyID),
		zap.String("plan_type", req.PlanType),
	)

	// Auto-activate package modules after license created
	if req.PackageID != "" && s.moduleManager != nil {
		activated, err := s.moduleManager.ActivatePackageModules(req.PackageID, req.CompanyID)
		if err != nil {
			// Log warning tapi tidak gagalkan pembuatan license
			s.logger.Warn("Failed to auto-activate package modules",
				zap.String("package_id", req.PackageID),
				zap.String("company_id", req.CompanyID),
				zap.Error(err),
			)
		} else {
			s.logger.Info("Package modules auto-activated",
				zap.String("package_id", req.PackageID),
				zap.String("company_id", req.CompanyID),
				zap.Strings("modules", activated),
			)
		}
	}

	response := license.ToResponse()
	return &response, nil
}

// GetLicense mengembalikan lisensi berdasarkan ID.
func (s *Service) GetLicense(id string) (*LicenseResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid license id: %w", err)
	}

	license, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	response := license.ToResponse()
	return &response, nil
}

// ListLicenses mengembalikan daftar lisensi dengan pagination.
func (s *Service) ListLicenses(page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	licenses, total, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return nil, err
	}

	// Collect company IDs untuk batch lookup
	companyIDs := make([]uuid.UUID, len(licenses))
	for i, l := range licenses {
		companyIDs[i] = l.CompanyID
	}
	companyNames := s.repo.FindCompanyNames(companyIDs)

	var responses []LicenseResponse
	for _, l := range licenses {
		resp := l.ToResponse()
		if name, ok := companyNames[l.CompanyID.String()]; ok {
			resp.CompanyName = name
		}
		responses = append(responses, resp)
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

// UpdateLicense mengupdate lisensi.
// Jika package_id berubah, module modules akan disesuaikan:
// - Old package modules di-deactivate
// - New package modules di-activate
func (s *Service) UpdateLicense(id string, req UpdateLicenseRequest) (*LicenseResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid license id: %w", err)
	}

	license, err := s.repo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	// Simpan old package_id sebelum diupdate
	var oldPackageID string
	if license.PackageID != nil {
		oldPackageID = license.PackageID.String()
	}
	companyID := license.CompanyID.String()

	if req.PlanType != nil {
		license.PlanType = *req.PlanType
	}
	if req.MaxEmployees != nil {
		license.MaxEmployees = *req.MaxEmployees
	}
	if req.MaxModules != nil {
		license.MaxModules = *req.MaxModules
	}
	if req.Status != nil {
		license.Status = *req.Status
	}
	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err == nil {
			license.StartDate = startDate
		}
	}
	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err == nil {
			license.EndDate = endDate
		}
	}

	// Track new package_id untuk perubahan
	var newPackageID string

	// Optional: update package association
	if req.PackageID != nil {
		if *req.PackageID == "" {
			// Remove package association
			license.PackageID = nil
			license.PackageName = ""
		} else {
			pid, err := uuid.Parse(*req.PackageID)
			if err != nil {
				return nil, fmt.Errorf("invalid package_id: %w", err)
			}
			pkgName, err := s.repo.FindPackageNameByID(pid)
			if err != nil {
				return nil, fmt.Errorf("package not found: %w", err)
			}
			license.PackageID = &pid
			license.PackageName = pkgName
			newPackageID = *req.PackageID
		}
	}

	if err := s.repo.Update(license); err != nil {
		return nil, err
	}

	// Handle module changes based on package_id transition
	if s.moduleManager != nil {
		// Case 1: Package changed from A to B (or A to none)
		if oldPackageID != "" && (newPackageID == "" || newPackageID != oldPackageID) {
			deactivated, err := s.moduleManager.DeactivatePackageModules(oldPackageID, companyID)
			if err != nil {
				s.logger.Warn("Failed to deactivate old package modules during license update",
					zap.String("license_id", id),
					zap.String("old_package_id", oldPackageID),
					zap.Error(err),
				)
			} else if len(deactivated) > 0 {
				s.logger.Info("Old package modules deactivated during license update",
					zap.String("old_package_id", oldPackageID),
					zap.Strings("modules", deactivated),
				)
			}
		}

		// Case 2: Package changed from none to B, or from A to B
		if newPackageID != "" && (oldPackageID == "" || newPackageID != oldPackageID) {
			activated, err := s.moduleManager.ActivatePackageModules(newPackageID, companyID)
			if err != nil {
				s.logger.Warn("Failed to activate new package modules during license update",
					zap.String("license_id", id),
					zap.String("new_package_id", newPackageID),
					zap.Error(err),
				)
			} else if len(activated) > 0 {
				s.logger.Info("New package modules activated during license update",
					zap.String("new_package_id", newPackageID),
					zap.Strings("modules", activated),
				)
			}
		}
	}

	response := license.ToResponse()
	return &response, nil
}

// DeleteLicense melakukan soft-delete lisensi.
// Jika lisensi memiliki package, modules dari package tersebut akan di-deactivate.
func (s *Service) DeleteLicense(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid license id: %w", err)
	}

	license, err := s.repo.FindByID(uid)
	if err != nil {
		return err
	}

	companyID := license.CompanyID.String()

	// Deactivate package modules if the license had a package association
	var oldPackageID string
	if license.PackageID != nil {
		oldPackageID = license.PackageID.String()
	}

	if oldPackageID != "" && s.moduleManager != nil {
		deactivated, err := s.moduleManager.DeactivatePackageModules(oldPackageID, companyID)
		if err != nil {
			s.logger.Warn("Failed to deactivate package modules during license deletion",
				zap.String("license_id", id),
				zap.String("package_id", oldPackageID),
				zap.Error(err),
			)
		} else if len(deactivated) > 0 {
			s.logger.Info("Package modules deactivated during license deletion",
				zap.String("package_id", oldPackageID),
				zap.Strings("modules", deactivated),
			)
		}
	}

	// Soft delete the license
	if err := s.repo.SoftDelete(uid); err != nil {
		return err
	}

	s.logger.Info("License soft-deleted",
		zap.String("license_id", id),
	)
	return nil
}
