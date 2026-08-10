package careerintelligence

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	dbResolver func(ctx context.Context) (*gorm.DB, error)
}

func NewRepository(resolver func(ctx context.Context) (*gorm.DB, error)) *Repository {
	return &Repository{dbResolver: resolver}
}

func (r *Repository) db(ctx context.Context) (*gorm.DB, error) {
	return r.dbResolver(ctx)
}

// =========================================================================
// Talent Map
// =========================================================================

func (r *Repository) CreateTalentMap(ctx context.Context, tm *CareerTalentMap) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.Create(tm).Error
}

func (r *Repository) FindTalentMapByID(ctx context.Context, id uuid.UUID) (*CareerTalentMap, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var tm CareerTalentMap
	if err := db.First(&tm, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tm, nil
}

func (r *Repository) ListTalentMaps(ctx context.Context, period, employeeID string, page, perPage int) ([]CareerTalentMap, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	query := db.Model(&CareerTalentMap{})
	if period != "" {
		query = query.Where("period = ?", period)
	}
	if employeeID != "" {
		query = query.Where("employee_id = ?", employeeID)
	}
	var total int64
	query.Count(&total)
	var list []CareerTalentMap
	offset := (page - 1) * perPage
	if err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateTalentMap(ctx context.Context, tm *CareerTalentMap) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.Save(tm).Error
}

func (r *Repository) DeleteTalentMap(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.Delete(&CareerTalentMap{}, "id = ?", id).Error
}

func (r *Repository) GetTalentGrid(ctx context.Context, period string) ([]CareerTalentMap, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []CareerTalentMap
	query := db.Model(&CareerTalentMap{})
	if period != "" {
		query = query.Where("period = ?", period)
	}
	if err := query.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) GetTalentHistoryByEmployee(ctx context.Context, employeeID uuid.UUID) ([]CareerTalentMap, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []CareerTalentMap
	if err := db.Where("employee_id = ?", employeeID).Order("period DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// =========================================================================
// Career Interest
// =========================================================================

func (r *Repository) CreateCareerInterest(ctx context.Context, ci *CareerInterest) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.Create(ci).Error
}

func (r *Repository) FindCareerInterestByID(ctx context.Context, id uuid.UUID) (*CareerInterest, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var ci CareerInterest
	if err := db.First(&ci, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &ci, nil
}

func (r *Repository) ListCareerInterests(ctx context.Context, employeeID string, page, perPage int) ([]CareerInterest, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	query := db.Model(&CareerInterest{}).Where("is_active = ?", true)
	if employeeID != "" {
		query = query.Where("employee_id = ?", employeeID)
	}
	var total int64
	query.Count(&total)
	var list []CareerInterest
	offset := (page - 1) * perPage
	if err := query.Order("recorded_at DESC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) GetInterestsByEmployee(ctx context.Context, employeeID uuid.UUID) ([]CareerInterest, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []CareerInterest
	if err := db.Where("employee_id = ? AND is_active = ?", employeeID, true).
		Order("recorded_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// =========================================================================
// Career Path (SKEMA TERPADU — header + steps)
// =========================================================================

// CreateCareerPathTx menyimpan header path + steps dalam satu transaksi
// (pola sama Employee Movement CreateCareerPathTx).
func (r *Repository) CreateCareerPathTx(ctx context.Context, cp *CareerPath, steps []CareerPathStep) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(cp).Error; err != nil {
			return err
		}
		for i := range steps {
			steps[i].CareerPathID = cp.ID
			if steps[i].ID == uuid.Nil {
				steps[i].ID = uuid.New()
			}
			if err := tx.Create(&steps[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) FindCareerPathByID(ctx context.Context, id uuid.UUID) (*CareerPath, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var cp CareerPath
	if err := db.First(&cp, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &cp, nil
}

// FindCareerPathByName mencari path berdasarkan name (termasuk yang sudah
// FindCareerPathByName mengembalikan path dengan nama tertentu. Memakai
// Unscoped agar nama dari path yang sudah soft-deleted tetap terdeteksi dan
// tetap "dipesan" (menghormati uk_career_paths_name saat buildCareerPathName
// menentukan nama unik — mencegah unique constraint violation saat nama yang
// sama hendak dipakai ulang).
func (r *Repository) FindCareerPathByName(ctx context.Context, name string) (*CareerPath, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var cp CareerPath
	if err := db.Unscoped().First(&cp, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &cp, nil
}

// ListCareerPathStepsByPathID mengembalikan steps satu path terurut sequence
// ascending.
func (r *Repository) ListCareerPathStepsByPathID(ctx context.Context, pathID uuid.UUID) ([]CareerPathStep, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var steps []CareerPathStep
	if err := db.Where("career_path_id = ?", pathID).
		Order("sequence ASC").Find(&steps).Error; err != nil {
		return nil, err
	}
	return steps, nil
}

// ListCareerPathStepsByPathIDs mengembalikan steps untuk banyak path sekaligus
// (batch query untuk ListCareerPaths).
func (r *Repository) ListCareerPathStepsByPathIDs(ctx context.Context, pathIDs []uuid.UUID) (map[uuid.UUID][]CareerPathStep, error) {
	if len(pathIDs) == 0 {
		return map[uuid.UUID][]CareerPathStep{}, nil
	}
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var steps []CareerPathStep
	if err := db.Where("career_path_id IN ?", pathIDs).
		Order("sequence ASC").Find(&steps).Error; err != nil {
		return nil, err
	}
	result := make(map[uuid.UUID][]CareerPathStep, len(pathIDs))
	for _, s := range steps {
		result[s.CareerPathID] = append(result[s.CareerPathID], s)
	}
	return result, nil
}

func (r *Repository) ListCareerPaths(ctx context.Context, page, perPage int) ([]CareerPath, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	query := db.Model(&CareerPath{}).Where("is_active = ?", true)
	var total int64
	query.Count(&total)
	var list []CareerPath
	offset := (page - 1) * perPage
	if err := query.Order("name ASC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// DeleteCareerPath menghapus header path (soft delete) beserta steps-nya
// (hard delete — steps tidak memiliki soft delete).
func (r *Repository) DeleteCareerPath(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("career_path_id = ?", id).Delete(&CareerPathStep{}).Error; err != nil {
			return err
		}
		return tx.Delete(&CareerPath{}, "id = ?", id).Error
	})
}

// FindCareerPathsBySource mengembalikan path aktif yang memiliki step dengan
// position_id = sourceTitleID dan sequence = 1 (langkah pertama) — ekuivalen
// edge CI "paths yang mulai dari source".
func (r *Repository) FindCareerPathsBySource(ctx context.Context, sourceTitleID uuid.UUID) ([]CareerPath, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []CareerPath
	subQuery := db.Model(&CareerPathStep{}).
		Select("career_path_id").
		Where("position_id = ? AND sequence = 1", sourceTitleID)
	if err := db.Model(&CareerPath{}).
		Where("is_active = ? AND id IN (?)", true, subQuery).
		Order("name ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}



// =========================================================================
// Succession Plan
// =========================================================================

func (r *Repository) CreateSuccessionPlan(ctx context.Context, sp *CareerSuccessionPlan) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.Create(sp).Error
}

func (r *Repository) FindSuccessionPlanByID(ctx context.Context, id uuid.UUID) (*CareerSuccessionPlan, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var sp CareerSuccessionPlan
	if err := db.First(&sp, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &sp, nil
}

func (r *Repository) ListSuccessionPlans(ctx context.Context, page, perPage int) ([]CareerSuccessionPlan, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	query := db.Model(&CareerSuccessionPlan{}).Where("status = ?", "ACTIVE")
	var total int64
	query.Count(&total)
	var list []CareerSuccessionPlan
	offset := (page - 1) * perPage
	if err := query.Order("priority_order ASC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateSuccessionPlan(ctx context.Context, sp *CareerSuccessionPlan) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.Save(sp).Error
}

func (r *Repository) DeleteSuccessionPlan(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.Delete(&CareerSuccessionPlan{}, "id = ?", id).Error
}

// =========================================================================
// Employee queries (cross-module)
// =========================================================================

func (r *Repository) GetEmployeePosition(ctx context.Context, employeeID uuid.UUID) (string, error) {
	db, err := r.db(ctx)
	if err != nil {
		return "", err
	}
	type result struct {
		PositionTitle string
	}
	var res result
	if err := db.Table("employees").
		Select("COALESCE(eo.position, '') as position_title").
		Joins("LEFT JOIN employments eo ON eo.employee_id = employees.id AND eo.deleted_at IS NULL").
		Where("employees.id = ?", employeeID).
		Scan(&res).Error; err != nil {
		return "", fmt.Errorf("employee not found: %w", err)
	}
	return res.PositionTitle, nil
}

func (r *Repository) GetPositionTitle(ctx context.Context, titleID uuid.UUID) (string, error) {
	db, err := r.db(ctx)
	if err != nil {
		return "", err
	}
	var name string
	if err := db.Table("job_management_titles").
		Select("name").
		Where("id = ?", titleID).
		Scan(&name).Error; err != nil {
		return "", err
	}
	return name, nil
}
