package workforceintelligence

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	dbFunc func(ctx context.Context) (*gorm.DB, error)
}

func NewRepository(dbFunc func(ctx context.Context) (*gorm.DB, error)) *Repository {
	return &Repository{dbFunc: dbFunc}
}

func (r *Repository) db(ctx context.Context) (*gorm.DB, error) {
	return r.dbFunc(ctx)
}

// =========================================================================
// Workforce Planning Headcount
// =========================================================================

func (r *Repository) CreateHeadcountPlan(ctx context.Context, h *WorkforcePlanningHeadcount) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(h).Error
}

func (r *Repository) FindHeadcountPlanByID(ctx context.Context, id uuid.UUID) (*WorkforcePlanningHeadcount, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var h WorkforcePlanningHeadcount
	if err := db.WithContext(ctx).First(&h, id).Error; err != nil {
		return nil, fmt.Errorf("headcount plan not found: %w", err)
	}
	return &h, nil
}

func (r *Repository) ListHeadcountPlans(ctx context.Context, period, orgID string, page, perPage int) ([]WorkforcePlanningHeadcount, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []WorkforcePlanningHeadcount
	var total int64
	query := db.WithContext(ctx).Model(&WorkforcePlanningHeadcount{})
	if period != "" {
		query = query.Where("period = ?", period)
	}
	if orgID != "" {
		query = query.Where("organization_id = ?", orgID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("period DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateHeadcountPlan(ctx context.Context, h *WorkforcePlanningHeadcount) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(h).Error
}

func (r *Repository) DeleteHeadcountPlan(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&WorkforcePlanningHeadcount{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("headcount plan not found")
	}
	return result.Error
}

// =========================================================================
// Workforce Forecast
// =========================================================================

func (r *Repository) CreateForecast(ctx context.Context, f *WorkforceForecast) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(f).Error
}

func (r *Repository) FindForecastByID(ctx context.Context, id uuid.UUID) (*WorkforceForecast, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var f WorkforceForecast
	if err := db.WithContext(ctx).First(&f, id).Error; err != nil {
		return nil, fmt.Errorf("forecast not found: %w", err)
	}
	return &f, nil
}

func (r *Repository) ListForecasts(ctx context.Context, period, orgID, forecastType string, page, perPage int) ([]WorkforceForecast, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []WorkforceForecast
	var total int64
	query := db.WithContext(ctx).Model(&WorkforceForecast{})
	if period != "" {
		query = query.Where("period = ?", period)
	}
	if orgID != "" {
		query = query.Where("organization_id = ?", orgID)
	}
	if forecastType != "" {
		query = query.Where("forecast_type = ?", forecastType)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("period DESC, confidence_level DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateForecast(ctx context.Context, f *WorkforceForecast) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(f).Error
}

func (r *Repository) DeleteForecast(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&WorkforceForecast{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("forecast not found")
	}
	return result.Error
}

// =========================================================================
// Workforce KPI
// =========================================================================

func (r *Repository) CreateKPI(ctx context.Context, k *WorkforceKPI) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(k).Error
}

func (r *Repository) ListKPIs(ctx context.Context, period, dimension, kpiCode string, page, perPage int) ([]WorkforceKPI, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []WorkforceKPI
	var total int64
	query := db.WithContext(ctx).Model(&WorkforceKPI{})
	if period != "" {
		query = query.Where("period = ?", period)
	}
	if dimension != "" {
		query = query.Where("dimension = ?", dimension)
	}
	if kpiCode != "" {
		query = query.Where("kpi_code = ?", kpiCode)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("period DESC, kpi_code ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// =========================================================================
// Analytics Cache
// =========================================================================

func (r *Repository) GetCache(ctx context.Context, cacheKey string) (*WorkforceAnalyticsCache, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var c WorkforceAnalyticsCache
	if err := db.WithContext(ctx).Where("cache_key = ?", cacheKey).First(&c).Error; err != nil {
		return nil, fmt.Errorf("cache not found: %w", err)
	}
	return &c, nil
}

func (r *Repository) SetCache(ctx context.Context, c *WorkforceAnalyticsCache) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	// Upsert: delete existing then create
	db.WithContext(ctx).Where("cache_key = ?", c.CacheKey).Delete(&WorkforceAnalyticsCache{})
	return db.WithContext(ctx).Create(c).Error
}

func (r *Repository) InvalidateCache(ctx context.Context, cacheType string) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Where("cache_type = ?", cacheType).Delete(&WorkforceAnalyticsCache{}).Error
}

// =========================================================================
// Scenario
// =========================================================================

func (r *Repository) CreateScenario(ctx context.Context, s *WorkforceScenario) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(s).Error
}

func (r *Repository) FindScenarioByID(ctx context.Context, id uuid.UUID) (*WorkforceScenario, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var s WorkforceScenario
	if err := db.WithContext(ctx).First(&s, id).Error; err != nil {
		return nil, fmt.Errorf("scenario not found: %w", err)
	}
	return &s, nil
}

func (r *Repository) ListScenarios(ctx context.Context, scenarioType string, page, perPage int) ([]WorkforceScenario, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []WorkforceScenario
	var total int64
	query := db.WithContext(ctx).Model(&WorkforceScenario{})
	if scenarioType != "" {
		query = query.Where("scenario_type = ?", scenarioType)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateScenario(ctx context.Context, s *WorkforceScenario) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(s).Error
}

func (r *Repository) DeleteScenario(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&WorkforceScenario{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("scenario not found")
	}
	return result.Error
}

// =========================================================================
// Risk Indicator
// =========================================================================

func (r *Repository) CreateRiskIndicator(ctx context.Context, ri *WorkforceRiskIndicator) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(ri).Error
}

func (r *Repository) FindRiskIndicatorByID(ctx context.Context, id uuid.UUID) (*WorkforceRiskIndicator, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var ri WorkforceRiskIndicator
	if err := db.WithContext(ctx).First(&ri, id).Error; err != nil {
		return nil, fmt.Errorf("risk indicator not found: %w", err)
	}
	return &ri, nil
}

func (r *Repository) ListRiskIndicators(ctx context.Context, period, riskLevel string, page, perPage int) ([]WorkforceRiskIndicator, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []WorkforceRiskIndicator
	var total int64
	query := db.WithContext(ctx).Model(&WorkforceRiskIndicator{})
	if period != "" {
		query = query.Where("period = ?", period)
	}
	if riskLevel != "" {
		query = query.Where("risk_level = ?", riskLevel)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("risk_level DESC, score DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateRiskIndicator(ctx context.Context, ri *WorkforceRiskIndicator) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(ri).Error
}

// =========================================================================
// Health Score
// =========================================================================

func (r *Repository) CreateHealthScore(ctx context.Context, hs *WorkforceHealthScore) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(hs).Error
}

func (r *Repository) FindHealthScoreByID(ctx context.Context, id uuid.UUID) (*WorkforceHealthScore, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var hs WorkforceHealthScore
	if err := db.WithContext(ctx).First(&hs, id).Error; err != nil {
		return nil, fmt.Errorf("health score not found: %w", err)
	}
	return &hs, nil
}

func (r *Repository) ListHealthScores(ctx context.Context, period, orgID string, page, perPage int) ([]WorkforceHealthScore, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []WorkforceHealthScore
	var total int64
	query := db.WithContext(ctx).Model(&WorkforceHealthScore{})
	if period != "" {
		query = query.Where("period = ?", period)
	}
	if orgID != "" {
		query = query.Where("organization_id = ?", orgID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("period DESC, score DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// =========================================================================
// Read-only queries to source modules
// =========================================================================

func (r *Repository) GetActiveEmployeeCount(ctx context.Context) (int, error) {
	db, err := r.db(ctx)
	if err != nil {
		return 0, err
	}
	var count int64
	// Read from employees table in tenant DB — employees tidak memakai soft delete
// (lihat employee/model.go); kolom status = 'active' adalah penanda aktif.
	if err := db.WithContext(ctx).Table("employees").Where("status = ?", "active").Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *Repository) GetDepartmentCount(ctx context.Context) (int, error) {
	db, err := r.db(ctx)
	if err != nil {
		return 0, err
	}
	var count int64
	if err := db.WithContext(ctx).Table("organizations").Where("deleted_at IS NULL AND parent_id IS NOT NULL").Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *Repository) GetEmployeesByDepartment(ctx context.Context) ([]DataPoint, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	type result struct {
		Name  string
		Count int
	}
	var rows []result
	if err := db.WithContext(ctx).Table("employees e").
		Select("o.name, COUNT(e.id) as count").
		Joins("LEFT JOIN organizations o ON o.id = e.organization_id").
		Where("e.status = ?", "active").
		Group("o.name").
		Order("count DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	var points []DataPoint
	for _, r := range rows {
		points = append(points, DataPoint{Label: r.Name, Value: float64(r.Count)})
	}
	return points, nil
}

func (r *Repository) GetEmployeeCountByType(ctx context.Context) ([]DataPoint, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	type result struct {
		EmploymentType string
		Count          int
	}
	var rows []result
	if err := db.WithContext(ctx).Table("employees").
		Select("employment_type, COUNT(id) as count").
		Where("status = ?", "active").
		Group("employment_type").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	var points []DataPoint
	for _, r := range rows {
		points = append(points, DataPoint{Label: r.EmploymentType, Value: float64(r.Count)})
	}
	return points, nil
}

func (r *Repository) GetMonthlyHeadcountTrend(ctx context.Context, months int) ([]DataPoint, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	type result struct {
		YearMonth string
		Count     int
	}
	var rows []result
	// Use application-level date extraction for cross-DB compatibility
	if err := db.WithContext(ctx).Table("employees").
		Select("SUBSTR(created_at, 1, 7) as year_month, COUNT(id) as count").
		Where("status = ?", "active").
		Group("year_month").
		Order("year_month DESC").
		Limit(months).
		Scan(&rows).Error; err != nil {
		// Fallback: return empty
		return []DataPoint{}, nil
	}
	var points []DataPoint
	for _, r := range rows {
		points = append(points, DataPoint{Label: r.YearMonth, Value: float64(r.Count)})
	}
	return points, nil
}

func (r *Repository) GetEmployeeGenderDistribution(ctx context.Context) ([]DataPoint, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	type result struct {
		Gender string
		Count  int
	}
	var rows []result
	if err := db.WithContext(ctx).Table("employees").
		Select("gender, COUNT(id) as count").
		Where("status = ?", "active").
		Group("gender").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	var points []DataPoint
	for _, r := range rows {
		points = append(points, DataPoint{Label: r.Gender, Value: float64(r.Count)})
	}
	return points, nil
}

func (r *Repository) GetMovementCountByType(ctx context.Context, period string) ([]DataPoint, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	type result struct {
		MovementType string
		Count        int
	}
	var rows []result
	// employee_movements tidak punya kolom deleted_at — cukup status executed.
	query := db.WithContext(ctx).Table("employee_movements").
		Select("movement_type, COUNT(id) as count").
		Where("status = 'executed'")
	if period != "" && len(period) >= 4 {
		query = query.Where("SUBSTR(created_at, 1, 4) = ?", period[:4])
	}
	if err := query.Group("movement_type").Scan(&rows).Error; err != nil {
		return nil, err
	}
	var points []DataPoint
	for _, r := range rows {
		points = append(points, DataPoint{Label: r.MovementType, Value: float64(r.Count)})
	}
	return points, nil
}

func (r *Repository) GetDepartments(ctx context.Context) ([]DepartmentGap, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	type result struct {
		ID   string
		Name string
	}
	var depts []result
	if err := db.WithContext(ctx).Table("organizations").
		Select("id, name").
		Where("deleted_at IS NULL AND parent_id IS NOT NULL").
		Find(&depts).Error; err != nil {
		return nil, err
	}
	var gaps []DepartmentGap
	for _, d := range depts {
		gaps = append(gaps, DepartmentGap{
			OrganizationID:   d.ID,
			OrganizationName: d.Name,
		})
	}
	return gaps, nil
}

// =========================================================================
// Recruitment pipeline reads (S-2 — expected hires → remaining gap)
// =========================================================================
// WI mengonsumsi data operasional Recruitment (job_requisitions +
// job_applications) untuk menghitung komponen hiring pipeline yang
// memperkaya gap analysis: open positions, accepted offers (expected hires),
// filled positions, dan pipeline per status.

type RecruitmentPipelineRow struct {
	Status string
	Count  int
}

// GetRecruitmentOpenPositions mengembalikan total slot terbuka dari requisition
// yang masih aktif (OPEN / IN_PROGRESS): sum(slots_available - slots_filled).
func (r *Repository) GetRecruitmentOpenPositions(ctx context.Context) (int, error) {
	db, err := r.db(ctx)
	if err != nil {
		return 0, err
	}
	var total int64
	if err := db.WithContext(ctx).Table("job_requisitions").
		Select("COALESCE(SUM(slots_available - slots_filled), 0)").
		Where("status IN ?", []string{"OPEN", "IN_PROGRESS"}).
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return int(total), nil
}

// GetRecruitmentAcceptedOffers mengembalikan jumlah aplikasi yang offer-nya
// diterima (status ACCEPTED) — ini adalah expected hires yang akan mengurangi
// remaining workforce gap.
func (r *Repository) GetRecruitmentAcceptedOffers(ctx context.Context) (int, error) {
	db, err := r.db(ctx)
	if err != nil {
		return 0, err
	}
	var count int64
	if err := db.WithContext(ctx).Table("job_applications").
		Where("status = ?", "ACCEPTED").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// GetRecruitmentFilledPositions mengembalikan total posisi yang sudah terisi
// (sum slots_filled pada requisition berstatus FILLED).
func (r *Repository) GetRecruitmentFilledPositions(ctx context.Context) (int, error) {
	db, err := r.db(ctx)
	if err != nil {
		return 0, err
	}
	var total int64
	if err := db.WithContext(ctx).Table("job_requisitions").
		Select("COALESCE(SUM(slots_filled), 0)").
		Where("status = ?", "FILLED").
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return int(total), nil
}

// GetRecruitmentPipeline mengembalikan jumlah aplikasi per status
// (NEW, SCREENED, SHORTLISTED, INTERVIEWED, OFFERED, ACCEPTED, REJECTED,
// WITHDRAWN) untuk funnel recruitment.
func (r *Repository) GetRecruitmentPipeline(ctx context.Context) ([]RecruitmentPipelineRow, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var rows []RecruitmentPipelineRow
	if err := db.WithContext(ctx).Table("job_applications").
		Select("status, COUNT(*) AS count").
		Group("status").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// =========================================================================
// Recruitment Analytics reads (S-3)
// =========================================================================

// GetRecruitmentTimeToHire mengembalikan rata-rata Time to Hire dalam hari:
// selisih accepted_at - applied_at (ms) untuk aplikasi berstatus ACCEPTED.
// Tidak ada aplikasi ter-hire → 0 (fail-safe, tidak error).
func (r *Repository) GetRecruitmentTimeToHire(ctx context.Context) (float64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return 0, err
	}
	var avgMS float64
	if err := db.WithContext(ctx).Table("job_applications").
		Select("COALESCE(AVG(accepted_at - applied_at), 0)").
		Where("status = ? AND applied_at > 0 AND accepted_at > 0", "ACCEPTED").
		Scan(&avgMS).Error; err != nil {
		return 0, err
	}
	return round1(avgMS / 86400000.0), nil
}

// FilledRequisitionDuration adalah durasi pengisian satu requisition FILLED
// (closed_at ms vs created_at timestamp). Dihitung di aplikasi agar kompatibel
// lintas dialek (BIGINT vs TIMESTAMP).
type FilledRequisitionDuration struct {
	ClosedAt  int64
	CreatedAt time.Time
}

// GetRecruitmentFilledRequisitionDurations mengembalikan pasangan
// closed_at/created_at requisition FILLED untuk menghitung Time to Fill.
func (r *Repository) GetRecruitmentFilledRequisitionDurations(ctx context.Context) ([]FilledRequisitionDuration, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var rows []FilledRequisitionDuration
	if err := db.WithContext(ctx).Table("job_requisitions").
		Select("closed_at, created_at").
		Where("status = ? AND closed_at IS NOT NULL AND closed_at > 0", "FILLED").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// RecruitmentOfferAcceptance adalah pasangan jumlah offer yang dibuat vs
// jumlah offer yang diterima (ACCEPTED).
type RecruitmentOfferAcceptance struct {
	Accepted int64
	Offered  int64
}

// GetRecruitmentOfferAcceptance menghitung Offer Acceptance Rate:
// accepted (status ACCEPTED) vs offered (offered_at terisi atau status
// OFFERED/ACCEPTED).
func (r *Repository) GetRecruitmentOfferAcceptance(ctx context.Context) (RecruitmentOfferAcceptance, error) {
	db, err := r.db(ctx)
	if err != nil {
		return RecruitmentOfferAcceptance{}, err
	}
	var row RecruitmentOfferAcceptance
	if err := db.WithContext(ctx).Table("job_applications").
		Select("COUNT(CASE WHEN status = ? THEN 1 END) AS accepted, COUNT(CASE WHEN offered_at > 0 OR status IN (?, ?) THEN 1 END) AS offered",
			"ACCEPTED", "OFFERED", "ACCEPTED").
		Scan(&row).Error; err != nil {
		return RecruitmentOfferAcceptance{}, err
	}
	return row, nil
}

// RecruitmentSourceConversion adalah satu baris konversi per source kandidat:
// jumlah kandidat (distinct) dan jumlah hire (aplikasi ACCEPTED).
type RecruitmentSourceConversion struct {
	Source     string
	Candidates int64
	Hires      int64
}

// GetRecruitmentSourceConversion menghitung Source Conversion per channel
// rekrutmen (candidates.source) → kandidat + hire, untuk conversion rate
// dan distribusi by source. Catatan: candidates di-DISTINCT, hires tidak
// (satu kandidat dengan >1 aplikasi ACCEPTED menghitung >1 hire — konsisten
// dengan semantik "accepted offer").
func (r *Repository) GetRecruitmentSourceConversion(ctx context.Context) ([]RecruitmentSourceConversion, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var rows []RecruitmentSourceConversion
	if err := db.WithContext(ctx).Table("candidates c").
		Select("c.source, COUNT(DISTINCT c.id) AS candidates, COUNT(CASE WHEN a.status = ? THEN 1 END) AS hires", "ACCEPTED").
		Joins("LEFT JOIN job_applications a ON a.candidate_id = c.id").
		Group("c.source").
		Order("candidates DESC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func round1(v float64) float64 {
	return math.Round(v*10) / 10
}

// =========================================================================
// Quality of Hire reads (S-6)
// =========================================================================

// QualityOfHireRow adalah satu hire (aplikasi ACCEPTED) beserta data kualitas
// yang tersedia lintas modul: interview score, status onboarding (proxy
// probation), performance final score, dan retensi employment. Komponen match
// score & assessment tidak ada di sini (G-9 & assessment belum dikumpulkan) —
// dihitung 0 di service (placeholder, pola sama S-3).
type QualityOfHireRow struct {
	ApplicationID    string
	Source           string
	RequisitionID    string
	RequisitionTitle string
	OrganizationID   string
	OrganizationName string
	EmployeeID       string
	OnboardingStatus string
	InterviewScore   float64
	PerformanceScore float64
	RetainedCount    int64
}

// GetQualityOfHireHires mengembalikan baris hire (job_applications ACCEPTED)
// dengan skor interview (AVG interviews.score per aplikasi), status onboarding
// (employee_onboardings via application_id — proxy hasil probation), skor
// performa (performance_evaluations final_score evaluasi selesai), dan jumlah
// employment aktif (effective_end_date IS NULL) untuk retensi. Komponen yang
// belum ada datanya otomatis 0/NULL-safe — tidak menggagalkan query.
func (r *Repository) GetQualityOfHireHires(ctx context.Context) ([]QualityOfHireRow, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var rows []QualityOfHireRow
	query := db.WithContext(ctx).Table("job_applications a").
		Select("a.id AS application_id, COALESCE(c.source, '') AS source, a.requisition_id AS requisition_id, "+
			"COALESCE(r.title, '') AS requisition_title, "+
			"COALESCE(r.organization_id, '') AS organization_id, COALESCE(o.nomenclature, '') AS organization_name, "+
			"COALESCE(eo.employee_id, '') AS employee_id, COALESCE(eo.status, '') AS onboarding_status, "+
			"(SELECT COALESCE(AVG(i.score), 0) FROM interviews i WHERE i.application_id = a.id AND i.score IS NOT NULL) AS interview_score, "+
			"(SELECT COALESCE(pe.final_score, 0) FROM performance_evaluations pe WHERE pe.employee_id = eo.employee_id AND pe.status IN (?, ?) ORDER BY pe.updated_at DESC, pe.id DESC LIMIT 1) AS performance_score, "+
			"(SELECT COUNT(*) FROM employments em WHERE em.employee_id = eo.employee_id AND em.effective_end_date IS NULL) AS retained_count",
			"ACTUAL_APPROVED", "COMPLETED").
		Joins("JOIN candidates c ON c.id = a.candidate_id").
		Joins("LEFT JOIN job_requisitions r ON r.id = a.requisition_id").
		Joins("LEFT JOIN organizations o ON o.id = r.organization_id").
		Joins("LEFT JOIN employee_onboardings eo ON eo.application_id = a.id").
		Where("a.status = ?", "ACCEPTED")
	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// =========================================================================
// Candidate Search
// =========================================================================

type candidateSearchOrgRow struct {
	OrganizationID   string
	OrganizationCode string
	OrganizationName string
	SummaryID        string
	SummaryCode      string
	SummaryDecreeNo  string
}

type candidateSearchPositionRow struct {
	OrganizationID string
	ID             string
	Title          string
	IsActive       bool
}

type candidateSearchCandidateRow struct {
	OrganizationID     string
	ID                 string
	FirstName          string
	LastName           string
	Email              string
	Phone              string
	CurrentTitle       *string
	CurrentCompany     *string
	Source             string
	ApplicationStatus  string
	RequisitionTitle   string
}

// CandidateSearchVacantOrgs mengembalikan organisasi yang LOWONG: tidak ada
// employment aktif (effective_date <= hari ini dan effective_end_date NULL/>= hari
// ini) dan berada di bawah Organization Summary berstatus active. Bisa difilter
// dengan kata kunci (kode/nama org atau kode/decree summary) dan/atau filter
// posisi (nama/kode organisasi — posisi yang sedang dicari kandidatnya, S-3).
func (r *Repository) CandidateSearchVacantOrgs(ctx context.Context, search, position *string, page, perPage int) ([]candidateSearchOrgRow, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	today := time.Now().Format("2006-01-02")

	query := db.WithContext(ctx).Table("organizations o").
		Select("o.id AS organization_id, o.code AS organization_code, o.nomenclature AS organization_name, os.id AS summary_id, os.code AS summary_code, os.decree_no AS summary_decree_no").
		Joins("JOIN organization_summaries os ON os.id = o.organization_summary_id").
		Where("o.deleted_at IS NULL AND os.deleted_at IS NULL AND os.status = ?", "active").
		Where("NOT EXISTS (SELECT 1 FROM employments e WHERE e.organization_id = o.id AND e.effective_date <= ? AND (e.effective_end_date IS NULL OR e.effective_end_date >= ?))", today, today)

	if search != nil && *search != "" {
		s := "%" + *search + "%"
		query = query.Where("(o.code LIKE ? OR o.full_code LIKE ? OR o.nomenclature LIKE ? OR os.code LIKE ? OR os.decree_no LIKE ?)", s, s, s, s, s)
	}
	if position != nil && *position != "" {
		p := "%" + *position + "%"
		query = query.Where("(o.nomenclature LIKE ? OR o.code LIKE ?)", p, p)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count vacant organizations: %w", err)
	}

	offset := (page - 1) * perPage
	var rows []candidateSearchOrgRow
	if err := query.Offset(offset).Limit(perPage).Order("o.code ASC").Scan(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list vacant organizations: %w", err)
	}
	return rows, total, nil
}

// CandidateSearchPositionsByOrgIDs mengembalikan daftar position (katalog
// positions milik modul organization) yang aktif pada organisasi lowong —
// dipakai untuk mencocokkan target career path (S-3 integrasi internal
// candidate eligible, data sumber Career Intelligence).
func (r *Repository) CandidateSearchPositionsByOrgIDs(ctx context.Context, orgIDs []uuid.UUID) ([]candidateSearchPositionRow, error) {
	if len(orgIDs) == 0 {
		return nil, nil
	}
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var rows []candidateSearchPositionRow
	if err := db.WithContext(ctx).Table("positions").
		Select("organization_id, id, title, is_active").
		Where("organization_id IN ? AND is_active = 1", orgIDs).
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("failed to load positions for vacant organizations: %w", err)
	}
	return rows, nil
}

// CandidateSearchCandidatesByOrgIDs mengembalikan kandidat recruitment yang
// melamar ke requisition AKTIF (OPEN/IN_PROGRESS) pada organisasi lowong tsb
// (status aplikasi bukan REJECTED/WITHDRAWN).
func (r *Repository) CandidateSearchCandidatesByOrgIDs(ctx context.Context, orgIDs []uuid.UUID) ([]candidateSearchCandidateRow, error) {
	if len(orgIDs) == 0 {
		return nil, nil
	}
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var rows []candidateSearchCandidateRow
	if err := db.WithContext(ctx).Table("candidates c").
		Select("r.organization_id, c.id, c.first_name, c.last_name, c.email, c.phone, c.current_title, c.current_company, c.source, a.status AS application_status, r.title AS requisition_title").
		Joins("JOIN job_applications a ON a.candidate_id = c.id").
		Joins("JOIN job_requisitions r ON r.id = a.requisition_id").
		Where("r.organization_id IN ?", orgIDs).
		Where("r.status IN ?", []string{"OPEN", "IN_PROGRESS"}).
		Where("a.status NOT IN ?", []string{"REJECTED", "WITHDRAWN"}).
		Order("a.created_at DESC").
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("failed to load candidates for vacant organizations: %w", err)
	}
	return rows, nil
}
