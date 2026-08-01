package company

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository untuk operasi database Company.
type Repository struct {
	db *gorm.DB
}

// NewRepository membuat Repository baru.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create menyimpan company baru ke database.
func (r *Repository) Create(company *Company) error {
	if err := r.db.Create(company).Error; err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}
	return nil
}

// FindByID mencari company berdasarkan ID.
func (r *Repository) FindByID(id uuid.UUID) (*Company, error) {
	var company Company
	if err := r.db.Where("id = ?", id).First(&company).Error; err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}
	return &company, nil
}

// FindBySlug mencari company berdasarkan slug.
func (r *Repository) FindBySlug(slug string) (*Company, error) {
	var company Company
	if err := r.db.Where("slug = ?", slug).First(&company).Error; err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}
	return &company, nil
}

// FindBySubdomain mencari company berdasarkan subdomain (mis. pt-inthros-jago-utama).
func (r *Repository) FindBySubdomain(subdomain string) (*Company, error) {
	var company Company
	if err := r.db.Where("subdomain = ?", subdomain).First(&company).Error; err != nil {
		return nil, fmt.Errorf("company not found by subdomain: %w", err)
	}
	return &company, nil
}

// FindByDomain mencari company berdasarkan domain penuh (mis. hris.pt-inthros.com).
func (r *Repository) FindByDomain(domain string) (*Company, error) {
	var company Company
	if err := r.db.Where("domain = ?", domain).First(&company).Error; err != nil {
		return nil, fmt.Errorf("company not found by domain: %w", err)
	}
	return &company, nil
}

// FindAll mengembalikan semua company dengan pagination.
// Gunakan query chain terpisah untuk Count dan Find untuk
// menghindari GORM v2 issue di mana Count() (terminal operation)
// dapat mengkonsumsi state dari query chain.
func (r *Repository) FindAll(page, perPage int) ([]Company, int64, error) {
	var companies []Company
	var total int64

	// Count total — chain terpisah
	if err := r.db.Model(&Company{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count companies: %w", err)
	}

	// Get paginated data — chain terpisah
	offset := (page - 1) * perPage
	if err := r.db.Model(&Company{}).
		Offset(offset).
		Limit(perPage).
		Order("created_at DESC").
		Find(&companies).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list companies: %w", err)
	}

	return companies, total, nil
}

// Update mengupdate company.
func (r *Repository) Update(company *Company) error {
	if err := r.db.Save(company).Error; err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}
	return nil
}

// SoftDelete melakukan soft delete company.
func (r *Repository) SoftDelete(id uuid.UUID) error {
	if err := r.db.Where("id = ?", id).Delete(&Company{}).Error; err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}
	return nil
}

// FindTenantConnectionsByCompanyIDs mencari tenant connections untuk
// kumpulan company IDs. Mengembalikan map[companyID]TenantConnection.
func (r *Repository) FindTenantConnectionsByCompanyIDs(ids []uuid.UUID) (map[uuid.UUID]TenantConnection, error) {
	if len(ids) == 0 {
		return make(map[uuid.UUID]TenantConnection), nil
	}
	var conns []TenantConnection
	if err := r.db.Where("company_id IN ?", ids).Find(&conns).Error; err != nil {
		return nil, fmt.Errorf("failed to find tenant connections: %w", err)
	}
	result := make(map[uuid.UUID]TenantConnection, len(conns))
	for _, c := range conns {
		result[c.CompanyID] = c
	}
	return result, nil
}

// FindTenantConnectionByCompanyID mencari satu tenant connection berdasarkan company ID.
func (r *Repository) FindTenantConnectionByCompanyID(companyID uuid.UUID) (*TenantConnection, error) {
	var conn TenantConnection
	if err := r.db.Where("company_id = ?", companyID).First(&conn).Error; err != nil {
		return nil, fmt.Errorf("tenant connection not found: %w", err)
	}
	return &conn, nil
}
