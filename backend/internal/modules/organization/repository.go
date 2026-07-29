package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	dbResolver func(ctx context.Context) (*gorm.DB, error)
}

func NewRepository(dbResolver func(ctx context.Context) (*gorm.DB, error)) *Repository {
	return &Repository{dbResolver: dbResolver}
}

func (r *Repository) getDB(ctx context.Context) (*gorm.DB, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required for tenant database resolution")
	}
	return r.dbResolver(ctx)
}

func (r *Repository) Create(ctx context.Context, org *Organization) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(org).Error
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*Organization, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var org Organization
	if err := db.Preload("Parent").First(&org, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}
	return &org, nil
}

// FindTree returns root organizations (parent_id IS NULL), optionally filtered by summary_id.
func (r *Repository) FindTree(ctx context.Context, summaryID string) ([]Organization, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var roots []Organization
	query := db.Where("parent_id IS NULL")
	if summaryID != "" {
		query = query.Where("organization_summary_id = ?", summaryID)
	}
	if err := query.
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Children.Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Children.Children.Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Order("sort_order ASC").
		Find(&roots).Error; err != nil {
		return nil, fmt.Errorf("failed to load organization tree: %w", err)
	}
	return roots, nil
}

// FindAll returns paginated organizations, optionally filtered by summary_id.
func (r *Repository) FindAll(ctx context.Context, page, perPage int, summaryID string) ([]Organization, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var orgs []Organization
	var total int64

	query := db.Model(&Organization{})
	if summaryID != "" {
		query = query.Where("organization_summary_id = ?", summaryID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("full_code ASC").Find(&orgs).Error; err != nil {
		return nil, 0, err
	}

	return orgs, total, nil
}

func (r *Repository) Update(ctx context.Context, org *Organization) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(org).Error
}

func (r *Repository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&Organization{}).Error
}

// =========================================================================
// History Repository Methods
// =========================================================================

func (r *Repository) CreateHistory(ctx context.Context, history *OrganizationHistory) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(history).Error
}

func (r *Repository) FindHistoryByOrgID(ctx context.Context, orgID uuid.UUID, page, perPage int) ([]OrganizationHistory, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var histories []OrganizationHistory
	var total int64

	query := db.Model(&OrganizationHistory{}).Where("organization_id = ?", orgID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (r *Repository) FindAllHistory(ctx context.Context, page, perPage int) ([]OrganizationHistory, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var histories []OrganizationHistory
	var total int64

	query := db.Model(&OrganizationHistory{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// =========================================================================
// Version Repository Methods
// =========================================================================

func (r *Repository) CreateVersion(ctx context.Context, version *OrganizationVersion) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(version).Error
}

func (r *Repository) FindVersionByID(ctx context.Context, id uuid.UUID) (*OrganizationVersion, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var version OrganizationVersion
	if err := db.First(&version, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}
	return &version, nil
}

func (r *Repository) FindAllVersions(ctx context.Context, page, perPage int) ([]OrganizationVersion, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var versions []OrganizationVersion
	var total int64

	query := db.Model(&OrganizationVersion{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&versions).Error; err != nil {
		return nil, 0, err
	}

	return versions, total, nil
}

func (r *Repository) UpdateVersion(ctx context.Context, version *OrganizationVersion) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(version).Error
}

func (r *Repository) FindAllOrganizationsFlat(ctx context.Context) ([]Organization, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var orgs []Organization
	if err := db.Order("full_code ASC").Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

// RestoreAllFromSnapshot replaces the entire organization tree atomically.
// Performs hard delete + bulk create in a single transaction to prevent data loss.
func (r *Repository) RestoreAllFromSnapshot(ctx context.Context, newOrgs []Organization) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. Remove all existing organizations
		if err := tx.Unscoped().Where("1 = 1").Delete(&Organization{}).Error; err != nil {
			return err
		}
		// 2. Create all organizations from snapshot
		for i := range newOrgs {
			if err := tx.Create(&newOrgs[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// =========================================================================
// Unique Validation Repository Methods
// =========================================================================

func (r *Repository) FindByFullCodeAndSummary(ctx context.Context, fullCode string, summaryID uuid.UUID) (*Organization, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var org Organization
	if err := db.Where("full_code = ? AND organization_summary_id = ?", fullCode, summaryID).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *Repository) FindByFullCodeAndSummaryExcludeSelf(ctx context.Context, fullCode string, summaryID uuid.UUID, excludeID uuid.UUID) (*Organization, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var org Organization
	if err := db.Where("full_code = ? AND organization_summary_id = ? AND id != ?", fullCode, summaryID, excludeID).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

// =========================================================================
// Clone Helper Repository Methods
// =========================================================================

func (r *Repository) FindAllBySummaryID(ctx context.Context, summaryID uuid.UUID) ([]Organization, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var orgs []Organization
	if err := db.Where("organization_summary_id = ?", summaryID).
		Order("full_code ASC").
		Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

func (r *Repository) BulkCreateOrganizations(ctx context.Context, orgs []Organization) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		for i := range orgs {
			if err := tx.Create(&orgs[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
