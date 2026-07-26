package pkgmgr

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository untuk operasi database Package & PackageModule.
type Repository struct {
	db *gorm.DB
}

// NewRepository membuat Repository baru.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// packagePreloadModules adalah helper untuk preload Modules dengan order by sort_order.
func packagePreloadModules() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}
}

// FindBySlug mencari paket berdasarkan slug.
func (r *Repository) FindBySlug(slug string) (*Package, error) {
	var p Package
	if err := r.db.Where("slug = ?", slug).Preload("Modules", packagePreloadModules()).First(&p).Error; err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}
	return &p, nil
}

// FindByID mencari paket berdasarkan ID.
func (r *Repository) FindByID(id uuid.UUID) (*Package, error) {
	var p Package
	if err := r.db.Where("id = ?", id).Preload("Modules", packagePreloadModules()).First(&p).Error; err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}
	return &p, nil
}

// FindAll mengembalikan semua paket dengan pagination.
// Jika moduleType tidak kosong, hanya mengembalikan paket yang mengandung modul dengan tipe tersebut.
func (r *Repository) FindAll(page, perPage int, moduleType string) ([]Package, int64, error) {
	var packages []Package
	var total int64

	// Count total
	countQuery := r.db.Model(&Package{})
	if moduleType != "" {
		subquery := r.db.Table("package_modules").
			Select("package_id").
			Joins("JOIN modules ON modules.id = package_modules.module_id").
			Where("modules.module_type = ?", moduleType)
		countQuery = countQuery.Where("packages.id IN (?)", subquery)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count packages: %w", err)
	}

	// Get paginated data
	offset := (page - 1) * perPage
	dataQuery := r.db.Model(&Package{})
	if moduleType != "" {
		subquery := r.db.Table("package_modules").
			Select("package_id").
			Joins("JOIN modules ON modules.id = package_modules.module_id").
			Where("modules.module_type = ?", moduleType)
		dataQuery = dataQuery.Where("packages.id IN (?)", subquery)
	}
	if err := dataQuery.
		Preload("Modules", packagePreloadModules()).
		Offset(offset).
		Limit(perPage).
		Order("sort_order ASC, created_at DESC").
		Find(&packages).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list packages: %w", err)
	}

	return packages, total, nil
}

// FindPublished mengembalikan paket published yang is_public=true.
// Jika moduleType tidak kosong, hanya mengembalikan paket yang mengandung modul dengan tipe tersebut.
func (r *Repository) FindPublished(moduleType string) ([]Package, error) {
	var packages []Package

	query := r.db.Model(&Package{}).
		Where("status = ? AND is_public = ?", PackagePublished, true)

	if moduleType != "" {
		subquery := r.db.Table("package_modules").
			Select("package_id").
			Joins("JOIN modules ON modules.id = package_modules.module_id").
			Where("modules.module_type = ?", moduleType)
		query = query.Where("packages.id IN (?)", subquery)
	}

	if err := query.
		Preload("Modules", packagePreloadModules()).
		Order("sort_order ASC, created_at DESC").
		Find(&packages).Error; err != nil {
		return nil, fmt.Errorf("failed to list published packages: %w", err)
	}
	return packages, nil
}

// Create menyimpan paket baru beserta modul-modulnya.
func (r *Repository) Create(pkg *Package) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Omit Modules agar GORM tidak auto-create module associations
		// (kita handle manual dengan set PackageID lalu create)
		if err := tx.Omit("Modules").Create(pkg).Error; err != nil {
			return fmt.Errorf("failed to create package: %w", err)
		}

		// Create module associations manually
		for i := range pkg.Modules {
			pkg.Modules[i].PackageID = pkg.ID
		}
		if len(pkg.Modules) > 0 {
			if err := tx.Create(&pkg.Modules).Error; err != nil {
				return fmt.Errorf("failed to create package modules: %w", err)
			}
		}

		return nil
	})
}

// Update mengupdate paket dan modul-modulnya (replace).
func (r *Repository) Update(pkg *Package) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Omit Modules agar GORM tidak auto-update module associations
		if err := tx.Omit("Modules").Save(pkg).Error; err != nil {
			return fmt.Errorf("failed to update package: %w", err)
		}

		// Replace module associations: delete old, insert new
		if err := tx.Where("package_id = ?", pkg.ID).Delete(&PackageModule{}).Error; err != nil {
			return fmt.Errorf("failed to delete old package modules: %w", err)
		}

		for i := range pkg.Modules {
			pkg.Modules[i].PackageID = pkg.ID
		}
		if len(pkg.Modules) > 0 {
			if err := tx.Create(&pkg.Modules).Error; err != nil {
				return fmt.Errorf("failed to create package modules: %w", err)
			}
		}

		return nil
	})
}

// Delete soft-deletes a package.
func (r *Repository) Delete(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Soft delete package modules
		if err := tx.Where("package_id = ?", id).Delete(&PackageModule{}).Error; err != nil {
			return fmt.Errorf("failed to delete package modules: %w", err)
		}
		// Soft delete package
		if err := tx.Where("id = ?", id).Delete(&Package{}).Error; err != nil {
			return fmt.Errorf("failed to delete package: %w", err)
		}
		return nil
	})
}

// FindModuleByID mencari modul platform berdasarkan ID (cross-module query).
type moduleInfo struct {
	ID        string
	Name      string
	Slug      string
	DependsOn *string
}

func (r *Repository) FindModuleInfo(moduleID uuid.UUID) (*moduleInfo, error) {
	var m moduleInfo
	row := r.db.Table("modules").
		Select("id, name, slug, depends_on").
		Where("id = ?", moduleID).
		Row()
	if err := row.Scan(&m.ID, &m.Name, &m.Slug, &m.DependsOn); err != nil {
		return nil, fmt.Errorf("module not found: %w", err)
	}
	// Normalize nil pointer to empty string for validation logic
	if m.DependsOn == nil {
		empty := ""
		m.DependsOn = &empty
	}
	return &m, nil
}

// FindModuleInfoMap mengembalikan map[id]moduleInfo untuk semua modul.
// Digunakan untuk batch lookup dependency pada ValidatePackageDependencies.
func (r *Repository) FindModuleInfoMap() (map[string]*moduleInfo, error) {
	rows, err := r.db.Table("modules").
		Select("id, name, slug, depends_on").
		Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to query modules: %w", err)
	}
	defer rows.Close()

	result := make(map[string]*moduleInfo)
	for rows.Next() {
		var m moduleInfo
		if err := r.db.ScanRows(rows, &m); err != nil {
			return nil, fmt.Errorf("failed to scan module row: %w", err)
		}
		// Normalize nil pointer to empty string for validation logic
		if m.DependsOn == nil {
			empty := ""
			m.DependsOn = &empty
		}
		result[m.ID] = &m
	}

	return result, nil
}

// parseDependsOn mengurai string depends_on (comma-separated slugs) menjadi slice.
func parseDependsOn(dependsOn string) []string {
	if dependsOn == "" {
		return nil
	}
	var result []string
	for _, s := range strings.Split(dependsOn, ",") {
		trimmed := strings.TrimSpace(s)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
