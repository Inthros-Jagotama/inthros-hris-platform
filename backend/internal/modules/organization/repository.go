package organization

import (
	"context"
	"fmt"
	"sort"
	"time"

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
	// Catatan: sengaja TIDAK Preload("Parent") — tidak ada pemanggil yang membaca
	// asosiasi Parent dari FindByID (hanya field skalar org itu sendiri yang dipakai).
	// Preload asosiasi BelongsTo yang basi bisa membuat GORM Save menimpa ParentID
	// dengan ID asosiasi lama (lihat service.Update).
	if err := db.First(&org, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}
	return &org, nil
}

// FindTree returns root organizations (parent_id IS NULL), optionally filtered by summary_id.
// Tree dibangun in-memory dari daftar datar organisasi (diurutkan full_code ASC agar
// parent selalu muncul sebelum anak) sehingga kedalaman tree TIDAK TERBATAS — tidak
// lagi tergantung kedalaman preload GORM (sebelumnya hanya 3 level).
// Setiap level anak diurutkan berdasarkan sort_order lalu full_code.
func (r *Repository) FindTree(ctx context.Context, summaryID string) ([]Organization, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}

	var orgs []Organization
	query := db.Model(&Organization{})
	if summaryID != "" {
		query = query.Where("organization_summary_id = ?", summaryID)
	}
	if err := query.Order("full_code ASC").Find(&orgs).Error; err != nil {
		return nil, fmt.Errorf("failed to load organizations: %w", err)
	}

	// Mapping id → node (pointer ke elemen slice orgs) untuk penempelan anak.
	// Penting: penempelan dilakukan IN-PLACE pada elemen slice orgs; root baru
	// disalin SETELAH itu agar Children-nya ikut terbawa (bukan snapshot kosong).
	byID := make(map[uuid.UUID]*Organization, len(orgs))
	for i := range orgs {
		byID[orgs[i].ID] = &orgs[i]
	}

	// Tempelkan anak ke parent — proses dari node terdalam (urutan full_code
	// terbalik) agar saat sebuah node disalin ke Children parent-nya, seluruh
	// Children-nya sudah terpasang (menghindari snapshot Children kosong pada
	// rantai bertingkat). Parent yang tidak ditemukan (data tidak konsisten)
	// membuat node diabaikan — konsisten dengan perilaku lama.
	for i := len(orgs) - 1; i >= 0; i-- {
		o := &orgs[i]
		if o.ParentID != nil {
			if parent, ok := byID[*o.ParentID]; ok {
				parent.Children = append(parent.Children, *o)
			}
		}
	}

	// Root = parent_id IS NULL (anak sudah terpasang in-place pada elemen orgs).
	var roots []Organization
	for i := range orgs {
		o := &orgs[i]
		if o.ParentID == nil {
			roots = append(roots, *o)
		}
	}

	// Urutkan setiap level (sibling) berdasarkan sort_order lalu full_code —
	// diperlukan karena reverse-append menghasilkan urutan sibling terbalik.
	sortOrgsBySortOrder(roots)
	for i := range roots {
		sortChildrenRecursive(&roots[i])
	}
	return roots, nil
}

// sortOrgsBySortOrder mengurutkan slice organisasi berdasarkan sort_order lalu full_code.
func sortOrgsBySortOrder(orgs []Organization) {
	sort.SliceStable(orgs, func(i, j int) bool {
		if orgs[i].SortOrder != orgs[j].SortOrder {
			return orgs[i].SortOrder < orgs[j].SortOrder
		}
		return orgs[i].FullCode < orgs[j].FullCode
	})
}

// sortChildrenRecursive mengurutkan children setiap node hingga semua kedalaman.
func sortChildrenRecursive(node *Organization) {
	if len(node.Children) == 0 {
		return
	}
	sortOrgsBySortOrder(node.Children)
	for i := range node.Children {
		sortChildrenRecursive(&node.Children[i])
	}
}

// FindAll returns paginated organizations, optionally filtered by summary_id and active_only.
func (r *Repository) FindAll(ctx context.Context, page, perPage int, summaryID string, activeOnly bool) ([]Organization, int64, error) {
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
	if activeOnly {
		// Hanya tampilkan organisasi yang memiliki summary dengan status 'active'
		query = query.Joins("JOIN organization_summaries ON organization_summaries.id = organizations.organization_summary_id").
			Where("organization_summaries.status = ?", "active")
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

// UpdateWithDescendants meng-update organisasi beserta seluruh descendants-nya dalam
// satu transaksi. Jika full_code berubah (mis. karena field code diubah atau parent
// dipindahkan), semua keturunan ikut di-update prefix-nya (termasuk yang soft-deleted
// via Unscoped) agar full_code tetap merupakan chain kode dari root (parent full_code
// + code), dan level-nya disesuaikan dengan selisih kedalaman (delta).
// Mengembalikan jumlah descendant yang di-update.
func (r *Repository) UpdateWithDescendants(ctx context.Context, org *Organization, oldFullCode string, oldLevel int) (int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return 0, err
	}

	needsPropagation := oldFullCode != "" && oldFullCode != org.FullCode && org.OrganizationSummaryID != nil

	// Pre-check: full_code terpanjang di antara descendants tidak boleh melebihi
	// varchar(50) setelah prefix diganti (jika tidak, error DB 500 yang tidak jelas).
	if needsPropagation {
		var maxLen int
		row := db.Model(&Organization{}).Unscoped().
			Select("COALESCE(MAX(LENGTH(full_code)), 0)").
			Where("organization_summary_id = ? AND id <> ? AND SUBSTR(full_code, 1, LENGTH(?)) = ?",
				org.OrganizationSummaryID, org.ID, oldFullCode, oldFullCode).
			Row()
		if err := row.Scan(&maxLen); err != nil {
			return 0, err
		}
		if newMax := maxLen - len(oldFullCode) + len(org.FullCode); newMax > 50 {
			return 0, fmt.Errorf("generated full_code for a descendant exceeds maximum length of 50 characters")
		}
	}

	// Pre-check: tidak ada full_code descendant yang setelah dipropagasi bentrok
	// dengan organisasi lain dalam summary yang sama (mencegah 500/duplikat).
	if needsPropagation {
		var descCodes []string
		if err := db.Model(&Organization{}).Unscoped().
			Where("organization_summary_id = ? AND id <> ? AND SUBSTR(full_code, 1, LENGTH(?)) = ?",
				org.OrganizationSummaryID, org.ID, oldFullCode, oldFullCode).
			Pluck("full_code", &descCodes).Error; err != nil {
			return 0, err
		}
		if len(descCodes) > 0 {
			newCodes := make([]string, len(descCodes))
			for i, c := range descCodes {
				newCodes[i] = org.FullCode + c[len(oldFullCode):]
			}
			var collision int64
			if err := db.Model(&Organization{}).Unscoped().
				Where("organization_summary_id = ? AND full_code IN ?", org.OrganizationSummaryID, newCodes).
				Count(&collision).Error; err != nil {
				return 0, err
			}
			if collision > 0 {
				return 0, fmt.Errorf("moving organization would create a duplicate full_code for one of its descendants")
			}
		}
	}

	// Selisih kedalaman: dipakai untuk menggeser level seluruh descendants
	// (mis. pindah parent: anak naik/turun level mengikuti org).
	levelDelta := org.Level - oldLevel

	var affected int64
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(org).Error; err != nil {
			return err
		}

		// Recalculate full_code untuk seluruh subtree bila prefix berubah.
		// full_code = chain kode dari root, jadi semua record dengan prefix
		// oldFullCode adalah descendant organisasi ini (dalam summary yang sama).
		if needsPropagation {
			// Ekspresi concat menyesuaikan dialect DB (MySQL: CONCAT, SQLite: ||).
			replaceExpr := "CONCAT(?, SUBSTRING(full_code, LENGTH(?) + 1))"
			if tx.Dialector.Name() == "sqlite" {
				replaceExpr = "? || SUBSTR(full_code, LENGTH(?) + 1)"
			}
			updates := map[string]interface{}{
				"full_code":  gorm.Expr(replaceExpr, org.FullCode, oldFullCode),
				"updated_at": time.Now(),
			}
			// Level hanya digeser bila kedalaman berubah (hindari write tak perlu).
			if levelDelta != 0 {
				updates["level"] = gorm.Expr("level + ?", levelDelta)
			}
			res := tx.Model(&Organization{}).Unscoped().
				Where("organization_summary_id = ? AND id <> ? AND SUBSTR(full_code, 1, LENGTH(?)) = ?",
					org.OrganizationSummaryID, org.ID, oldFullCode, oldFullCode).
				Updates(updates)
			if res.Error != nil {
				return res.Error
			}
			affected = res.RowsAffected
		}
		return nil
	})
	return affected, err
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
