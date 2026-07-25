package workforceintelligence

import (
	"context"
	"fmt"

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
	// Read from employees table in tenant DB
	if err := db.WithContext(ctx).Table("employees").Where("deleted_at IS NULL").Count(&count).Error; err != nil {
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
		Where("e.deleted_at IS NULL").
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
		Where("deleted_at IS NULL").
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
		Where("deleted_at IS NULL").
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
		Where("deleted_at IS NULL").
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
	query := db.WithContext(ctx).Table("employee_movements").
		Select("movement_type, COUNT(id) as count").
		Where("deleted_at IS NULL AND status = 'executed'")
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
