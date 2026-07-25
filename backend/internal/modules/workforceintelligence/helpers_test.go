package workforceintelligence

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

// setupTestDB creates an in-memory SQLite database and auto-migrates all models.
func setupTestDB() (*gorm.DB, func(ctx context.Context) (*gorm.DB, error), func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	if err := db.AutoMigrate(
		&WorkforcePlanningHeadcount{},
		&WorkforceForecast{},
		&WorkforceKPI{},
		&WorkforceAnalyticsCache{},
		&WorkforceScenario{},
		&WorkforceRiskIndicator{},
		&WorkforceHealthScore{},
	); err != nil {
		panic(fmt.Sprintf("failed to migrate test db: %v", err))
	}

	dbResolver := func(ctx context.Context) (*gorm.DB, error) {
		return db, nil
	}

	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	return db, dbResolver, cleanup
}

// newTestService creates a Service with in-memory SQLite repository.
func newTestService() (*Service, *Repository, *gorm.DB, func()) {
	db, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	return svc, repo, db, func() {
		cleanup()
		_ = logger.Sync()
	}
}

// =========================================================================
// Fixture helpers
// =========================================================================

var _fixtureCounter int

func createTestHeadcountPlan(repo *Repository) *WorkforcePlanningHeadcount {
	ctx := context.Background()
	_fixtureCounter++
	h := &WorkforcePlanningHeadcount{
		ID:             uuid.New(),
		Period:         fmt.Sprintf("2026-Q%d", (_fixtureCounter%4)+1),
		OrganizationID: uuid.New(),
		PlannedHC:      50 + _fixtureCounter,
		ActualHC:       45,
		SnapshotDate:   time.Now().Truncate(24 * time.Hour),
	}
	if err := repo.CreateHeadcountPlan(ctx, h); err != nil {
		panic(fmt.Sprintf("failed to create test headcount plan: %v", err))
	}
	return h
}

func createTestForecast(repo *Repository) *WorkforceForecast {
	ctx := context.Background()
	_fixtureCounter++
	f := &WorkforceForecast{
		ID:              uuid.New(),
		Period:          fmt.Sprintf("2026-Q%d", (_fixtureCounter%4)+1),
		OrganizationID:  uuid.New(),
		ForecastType:    "DEMAND",
		Headcount:       60 + _fixtureCounter,
		ConfidenceLevel: 85.0,
		Parameters:      JSON{"growth": 0.05},
	}
	if err := repo.CreateForecast(ctx, f); err != nil {
		panic(fmt.Sprintf("failed to create test forecast: %v", err))
	}
	return f
}

func createTestKPI(repo *Repository) *WorkforceKPI {
	ctx := context.Background()
	_fixtureCounter++
	target := 100.0
	k := &WorkforceKPI{
		ID:         uuid.New(),
		Period:     fmt.Sprintf("2026-Q%d", (_fixtureCounter%4)+1),
		KpiCode:    fmt.Sprintf("KPI_%03d", _fixtureCounter),
		KpiName:    fmt.Sprintf("Test KPI %d", _fixtureCounter),
		Value:      85.5,
		Target:     &target,
		Unit:       "PCT",
		Dimension:  "COMPANY",
		SnapshotAt: time.Now().Truncate(24 * time.Hour),
	}
	if err := repo.CreateKPI(ctx, k); err != nil {
		panic(fmt.Sprintf("failed to create test KPI: %v", err))
	}
	return k
}

func createTestScenario(repo *Repository) *WorkforceScenario {
	ctx := context.Background()
	_fixtureCounter++
	s := &WorkforceScenario{
		ID:           uuid.New(),
		Name:         fmt.Sprintf("Scenario %d", _fixtureCounter),
		Description:  "Test scenario",
		ScenarioType: "GROWTH",
		Parameters:   JSON{"growth_rate": 10.0},
		Status:       "DRAFT",
		CreatedBy:    uuid.New(),
	}
	if err := repo.CreateScenario(ctx, s); err != nil {
		panic(fmt.Sprintf("failed to create test scenario: %v", err))
	}
	return s
}

func createTestRiskIndicator(repo *Repository) *WorkforceRiskIndicator {
	ctx := context.Background()
	_fixtureCounter++
	deptID := uuid.New()
	ri := &WorkforceRiskIndicator{
		ID:             uuid.New(),
		Period:         fmt.Sprintf("2026-Q%d", (_fixtureCounter%4)+1),
		RiskCode:       fmt.Sprintf("RISK_%03d", _fixtureCounter),
		RiskName:       fmt.Sprintf("Test Risk %d", _fixtureCounter),
		RiskLevel:      "MEDIUM",
		Score:          75.0,
		Threshold:      80.0,
		DepartmentID:   &deptID,
		Recommendation: "Monitor closely",
		SnapshotAt:     time.Now().Truncate(24 * time.Hour),
	}
	if err := repo.CreateRiskIndicator(ctx, ri); err != nil {
		panic(fmt.Sprintf("failed to create test risk indicator: %v", err))
	}
	return ri
}

func createTestHealthScore(repo *Repository) *WorkforceHealthScore {
	ctx := context.Background()
	_fixtureCounter++
	hs := &WorkforceHealthScore{
		ID:                 uuid.New(),
		Period:             fmt.Sprintf("2026-Q%d", (_fixtureCounter%4)+1),
		OrganizationID:     uuid.New(),
		Score:              78.5,
		SpanOfControl:      5.2,
		ManagerRatio:       14.8,
		PromotionRate:      8.5,
		InternalHiringRate: 45.3,
		SuccessionCoverage: 68.0,
		StabilityRatio:     72.0,
		Components:         JSON{"score": 78.5},
		SnapshotAt:         time.Now().Truncate(24 * time.Hour),
	}
	if err := repo.CreateHealthScore(ctx, hs); err != nil {
		panic(fmt.Sprintf("failed to create test health score: %v", err))
	}
	return hs
}

func createTestCacheEntry(repo *Repository) *WorkforceAnalyticsCache {
	ctx := context.Background()
	_fixtureCounter++
	c := &WorkforceAnalyticsCache{
		CacheKey:  fmt.Sprintf("hc_trend_%d", _fixtureCounter),
		CacheType: "HC",
		Data:      JSON{"total": 100, "trend": []int{95, 98, 100}},
		Period:    fmt.Sprintf("2026-Q%d", (_fixtureCounter%4)+1),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	if err := repo.SetCache(ctx, c); err != nil {
		panic(fmt.Sprintf("failed to create test cache entry: %v", err))
	}
	return c
}
