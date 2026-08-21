package careerintelligence

import (
	"context"
	"fmt"
	"time"

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
// Talent Map Settings (singleton, ambang batas banding skor)
// =========================================================================

// GetOrCreateTalentMapSettings returns the tenant's singleton row, creating
// one with defaults (50/80 for both performance & potential) on first read —
// same auto-create-default-on-read convention as employee_id_format_settings.
func (r *Repository) GetOrCreateTalentMapSettings(ctx context.Context) (*TalentMapSettings, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var s TalentMapSettings
	err = db.First(&s).Error
	if err == nil {
		return &s, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to load talent map settings: %w", err)
	}
	s = TalentMapSettings{
		ID:                 uuid.NewString(),
		PerformanceLowMax:  50,
		PerformanceHighMin: 80,
		PotentialLowMax:    50,
		PotentialHighMin:   80,
	}
	if err := db.Create(&s).Error; err != nil {
		return nil, fmt.Errorf("failed to create default talent map settings: %w", err)
	}
	return &s, nil
}

// UpdateTalentMapSettings persists changes to the singleton row.
func (r *Repository) UpdateTalentMapSettings(ctx context.Context, s *TalentMapSettings) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	if err := db.Save(s).Error; err != nil {
		return fmt.Errorf("failed to update talent map settings: %w", err)
	}
	return nil
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

// ListCareerPaths mengembalikan daftar path (is_active = true) terurut name
// ASC dengan pagination. Keyword opsional memfilter substring nama path.
func (r *Repository) ListCareerPaths(ctx context.Context, page, perPage int, keyword string) ([]CareerPath, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	query := db.Model(&CareerPath{}).Where("is_active = ?", true)
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	var total int64
	query.Count(&total)
	var list []CareerPath
	offset := (page - 1) * perPage
	if err := query.Order("name ASC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// UpdateCareerPathTx memperbarui header path dan mengganti SELURUH steps-nya
// dalam satu transaksi (semantik full-replace — pola sama EM). Header memakai
// map agar Description dapat dikosongkan eksplisit (NULL) dan IsActive
// di-toggle tanpa ambiguitas zero-value struct update.
func (r *Repository) UpdateCareerPathTx(ctx context.Context, cp *CareerPath, steps []CareerPathStep) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"name":      cp.Name,
			"is_active": cp.IsActive,
		}
		if cp.Description != "" {
			updates["description"] = cp.Description
		} else {
			updates["description"] = nil
		}
		if err := tx.Model(&CareerPath{}).Where("id = ?", cp.ID).Updates(updates).Error; err != nil {
			return err
		}
		if err := tx.Where("career_path_id = ?", cp.ID).Delete(&CareerPathStep{}).Error; err != nil {
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

// GetPositionNamesByIDs resolves position titles for a batch of ids — pola
// sama EM resolveNamesByIDs; dipakai enrichment position_name pada steps.
func (r *Repository) GetPositionNamesByIDs(ctx context.Context, ids []uuid.UUID) (map[string]string, error) {
	result := make(map[string]string, len(ids))
	if len(ids) == 0 {
		return result, nil
	}
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	type row struct {
		ID   string
		Name string
	}
	var rows []row
	// "positions" table tidak pernah dipakai aplikasi (Organization = Position,
	// lihat prinsip modul ini) -- position_id yang dikirim FE (dari dropdown
	// /organizations) sebenarnya adalah organizations.id, jadi resolve nama
	// dari organizations.nomenclature, bukan dari tabel positions yang mati.
	if err := db.Table("organizations").
		Where("id IN ?", ids).
		Select("id AS id, nomenclature AS name").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, rw := range rows {
		result[rw.ID] = rw.Name
	}
	return result, nil
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

// FindCareerPathsByTarget mengembalikan path aktif yang memiliki step dengan
// position_id = targetTitleID pada sequence = max (langkah terakhir) — dipakai
// S-4 untuk menemukan jenjang karier yang menuju ke posisi lowongan.
func (r *Repository) FindCareerPathsByTarget(ctx context.Context, targetTitleID uuid.UUID) ([]CareerPath, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []CareerPath
	// Step target = step dengan sequence maksimal pada path yang memiliki
	// position_id = targetTitleID.
	subQuery := db.Model(&CareerPathStep{}).
		Select("DISTINCT s1.career_path_id").
		Table("career_path_steps s1").
		Where("s1.position_id = ?", targetTitleID).
		Where("s1.sequence = (SELECT MAX(s2.sequence) FROM career_path_steps s2 WHERE s2.career_path_id = s1.career_path_id)")
	if err := db.Model(&CareerPath{}).
		Where("is_active = ? AND id IN (?)", true, subQuery).
		Order("name ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// EligibleEmployeeRow adalah baris hasil query employee yang saat ini (employment
// aktif) memegang posisi tertentu — dipakai S-4 internal candidate eligibility.
type EligibleEmployeeRow struct {
	EmployeeID   string
	Name         string
	PositionID   string
	PositionName string
}

// ListEligibleEmployeesByPositionIDs mengembalikan employee dengan employment
// aktif (effective_date <= hari ini, effective_end_date NULL/>= hari ini,
// bukan soft-delete) yang memegang salah satu posisi dalam daftar. Dipakai S-4:
// employee pada source step career path adalah kandidat internal menuju target.
func (r *Repository) ListEligibleEmployeesByPositionIDs(ctx context.Context, positionIDs []uuid.UUID) ([]EligibleEmployeeRow, error) {
	if len(positionIDs) == 0 {
		return []EligibleEmployeeRow{}, nil
	}
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	today := time.Now().Format("2006-01-02")
	var rows []EligibleEmployeeRow
	// employees/employments tidak punya kolom deleted_at (lihat employee/model.go
	// — tanpa soft delete); filter employee aktif pakai e.status.
	if err := db.Table("employees e").
		Select("e.id AS employee_id, e.name AS name, em.position_id AS position_id, o.nomenclature AS position_name").
		Joins("JOIN employments em ON em.employee_id = e.id").
		Joins("LEFT JOIN organizations o ON o.id = em.position_id").
		Where("e.status = ?", "active").
		Where("em.position_id IN ?", positionIDs).
		Where("em.effective_date <= ?", today).
		Where("(em.effective_end_date IS NULL OR em.effective_end_date >= ?)", today).
		Order("e.name ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
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

// SuccessionGapRow adalah agregasi per posisi kunci (positions) yang memiliki
// succession plan ACTIVE: jumlah successor terencana + berapa yang siap
// (READY_NOW). Dipakai S-5 untuk menandai posisi kunci tanpa successor siap.
type SuccessionGapRow struct {
	PositionID      string
	PositionTitle   string
	PositionCode    string
	OrganizationID  string
	SuccessorCount  int
	ReadyNowCount   int
}

// ListSuccessionGapPositions mengembalikan posisi kunci (positions yang memiliki
// ≥1 succession plan ACTIVE) beserta statistik successor-nya, TANPA memfilter
// readiness — service yang menentukan gap (tidak ada successor READY_NOW).
// Join organizations untuk nama (Organization = Position, lihat komentar
// GetPositionNamesByIDs) — organization_id di response sengaja sama dengan
// position_id itu sendiri karena tidak ada level "organization pemilik posisi"
// terpisah dalam model data ini.
func (r *Repository) ListSuccessionGapPositions(ctx context.Context) ([]SuccessionGapRow, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var rows []SuccessionGapRow
	if err := db.Table("career_succession_plans sp").
		Select("sp.position_id AS position_id, COALESCE(o.nomenclature, '') AS position_title, COALESCE(o.full_code, '') AS position_code, sp.position_id AS organization_id, COUNT(*) AS successor_count, COALESCE(SUM(CASE WHEN sp.readiness_level = ? THEN 1 ELSE 0 END), 0) AS ready_now_count", "READY_NOW").
		Joins("LEFT JOIN organizations o ON o.id = sp.position_id").
		Where("sp.status = ?", "ACTIVE").
		Where("sp.deleted_at IS NULL").
		Group("sp.position_id, o.nomenclature, o.full_code").
		Order("position_title ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// CheckSuccessionGapByPosition mengembalikan true bila posisi adalah posisi
// kunci dengan succession gap — yaitu memiliki ≥1 succession plan ACTIVE dan
// TIDAK ada satupun successor dengan readiness READY_NOW. Dipakai S-5 oleh
// Recruitment (via narrow provider) untuk fallback external recruitment.
func (r *Repository) CheckSuccessionGapByPosition(ctx context.Context, positionID uuid.UUID) (bool, error) {
	db, err := r.db(ctx)
	if err != nil {
		return false, err
	}
	var count int64
	if err := db.Model(&CareerSuccessionPlan{}).
		Where("position_id = ? AND status = ? AND deleted_at IS NULL", positionID, "ACTIVE").
		Count(&count).Error; err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	var readyNow int64
	if err := db.Model(&CareerSuccessionPlan{}).
		Where("position_id = ? AND status = ? AND deleted_at IS NULL AND readiness_level = ?", positionID, "ACTIVE", "READY_NOW").
		Count(&readyNow).Error; err != nil {
		return false, err
	}
	return readyNow == 0, nil
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
	// employments tidak punya kolom position — ambil nama via join organizations
	// (Organization = Position, lihat GetPositionNamesByIDs).
	if err := db.Table("employees").
		Select("COALESCE(o.nomenclature, '') as position_title").
		Joins("LEFT JOIN employments eo ON eo.employee_id = employees.id").
		Joins("LEFT JOIN organizations o ON o.id = eo.position_id").
		Where("employees.id = ?", employeeID).
		Scan(&res).Error; err != nil {
		return "", fmt.Errorf("employee not found: %w", err)
	}
	return res.PositionTitle, nil
}

// GetPositionTitle resolves the display name for a "target title" id used by
// Gap Analysis. "job_management_titles" is a dead table (never written by any
// FE flow -- Job Management is built entirely on organizations, see
// GetPositionNamesByIDs comment) -- resolve from organizations instead,
// consistent with how the rest of this module treats position/title ids.
func (r *Repository) GetPositionTitle(ctx context.Context, titleID uuid.UUID) (string, error) {
	db, err := r.db(ctx)
	if err != nil {
		return "", err
	}
	var name string
	if err := db.Table("organizations").
		Select("nomenclature").
		Where("id = ?", titleID).
		Scan(&name).Error; err != nil {
		return "", err
	}
	return name, nil
}

// CompetencyRequirement is one competency required by a target organization,
// derived from that organization's own latest finalized competency
// assessment (competency_scores.organization_id has a unique index -- one
// row per org -- so "latest" is simply the org's single existing row).
type CompetencyRequirement struct {
	CompetencyID   uuid.UUID
	CompetencyName string
	StandardLevel  int
}

// GetOrgCompetencyRequirements reads the required competency levels for a
// target organization, sourced from competency_score_details of that org's
// own finalized competency_scores row (real data -- no standalone "required
// competencies per position" master table exists in this codebase; see
// docs/module-career-intelligence-plan.md §8.5/§9). Returns an empty slice
// (not an error) if the target org has never been assessed.
func (r *Repository) GetOrgCompetencyRequirements(ctx context.Context, orgID uuid.UUID) ([]CompetencyRequirement, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var out []CompetencyRequirement
	err = db.Table("competency_score_details csd").
		Select("csd.competency_id AS competency_id, c.name AS competency_name, csd.standard_level AS standard_level").
		Joins("JOIN competencies c ON c.id = csd.competency_id").
		Where("csd.competency_score_id = (SELECT id FROM competency_scores WHERE organization_id = ? ORDER BY assessed_at DESC LIMIT 1)", orgID).
		Where("csd.standard_level IS NOT NULL").
		Scan(&out).Error
	if err != nil {
		return nil, fmt.Errorf("failed to load org competency requirements: %w", err)
	}
	return out, nil
}

// GetEmployeeCompetencyLevels reads an employee's own latest assessed
// competency levels, keyed by competency_id, from their most recent
// competency_scores row (across whichever organization they were last
// assessed in).
func (r *Repository) GetEmployeeCompetencyLevels(ctx context.Context, employeeID uuid.UUID) (map[uuid.UUID]int, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	type row struct {
		CompetencyID  uuid.UUID
		EmployeeLevel int
	}
	var rows []row
	err = db.Table("competency_score_details csd").
		Select("csd.competency_id AS competency_id, csd.employee_level AS employee_level").
		Where("csd.competency_score_id = (SELECT id FROM competency_scores WHERE employee_id = ? ORDER BY assessed_at DESC LIMIT 1)", employeeID).
		Where("csd.employee_level IS NOT NULL").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to load employee competency levels: %w", err)
	}
	levels := make(map[uuid.UUID]int, len(rows))
	for _, rw := range rows {
		levels[rw.CompetencyID] = rw.EmployeeLevel
	}
	return levels, nil
}
