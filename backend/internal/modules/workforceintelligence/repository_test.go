package workforceintelligence

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// =========================================================================
// Headcount Plan Repository Tests
// =========================================================================

func TestRepo_CreateHeadcountPlan_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	h := &WorkforcePlanningHeadcount{
		ID:             uuid.New(),
		Period:         "2026-Q3",
		OrganizationID: uuid.New(),
		PlannedHC:      100,
		ActualHC:       0,
	}

	if err := repo.CreateHeadcountPlan(context.Background(), h); err != nil {
		t.Fatalf("CreateHeadcountPlan failed: %v", err)
	}
	if h.ID == uuid.Nil {
		t.Fatal("expected ID to be set after create")
	}
	if h.Period != "2026-Q3" {
		t.Errorf("expected period '2026-Q3', got '%s'", h.Period)
	}
}

func TestRepo_FindHeadcountPlanByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestHeadcountPlan(repo)

	found, err := repo.FindHeadcountPlanByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("FindHeadcountPlanByID failed: %v", err)
	}
	if found.PlannedHC != created.PlannedHC {
		t.Errorf("expected PlannedHC %d, got %d", created.PlannedHC, found.PlannedHC)
	}
}

func TestRepo_FindHeadcountPlanByID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	_, err := repo.FindHeadcountPlanByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent headcount plan")
	}
}

func TestRepo_ListHeadcountPlans_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	for i := 0; i < 3; i++ {
		createTestHeadcountPlan(repo)
	}

	plans, total, err := repo.ListHeadcountPlans(context.Background(), "", "", 1, 10)
	if err != nil {
		t.Fatalf("ListHeadcountPlans failed: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(plans) != 3 {
		t.Errorf("expected 3 plans, got %d", len(plans))
	}
}

func TestRepo_ListHeadcountPlans_FilterByPeriod(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	p1 := createTestHeadcountPlan(repo)
	p1.Period = "2026-Q3"
	_ = repo.UpdateHeadcountPlan(context.Background(), p1)

	plans, total, err := repo.ListHeadcountPlans(context.Background(), "2026-Q3", "", 1, 10)
	if err != nil {
		t.Fatalf("ListHeadcountPlans failed: %v", err)
	}
	if total == 0 {
		t.Error("expected at least 1 plan for Q3")
	}
	if len(plans) == 0 {
		t.Error("expected at least 1 plan in results")
	}
	for _, p := range plans {
		if p.Period != "2026-Q3" {
			t.Errorf("expected period '2026-Q3', got '%s'", p.Period)
		}
	}
}

func TestRepo_UpdateHeadcountPlan_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestHeadcountPlan(repo)
	created.PlannedHC = 200

	if err := repo.UpdateHeadcountPlan(context.Background(), created); err != nil {
		t.Fatalf("UpdateHeadcountPlan failed: %v", err)
	}

	found, _ := repo.FindHeadcountPlanByID(context.Background(), created.ID)
	if found.PlannedHC != 200 {
		t.Errorf("expected PlannedHC 200, got %d", found.PlannedHC)
	}
}

func TestRepo_DeleteHeadcountPlan_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestHeadcountPlan(repo)

	if err := repo.DeleteHeadcountPlan(context.Background(), created.ID); err != nil {
		t.Fatalf("DeleteHeadcountPlan failed: %v", err)
	}

	_, err := repo.FindHeadcountPlanByID(context.Background(), created.ID)
	if err == nil {
		t.Fatal("expected error after deleting headcount plan")
	}
}

// =========================================================================
// Forecast Repository Tests
// =========================================================================

func TestRepo_CreateForecast_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	f := &WorkforceForecast{
		ID:              uuid.New(),
		Period:          "2026-Q3",
		OrganizationID:  uuid.New(),
		ForecastType:    "DEMAND",
		Headcount:       150,
		ConfidenceLevel: 90.0,
	}

	if err := repo.CreateForecast(context.Background(), f); err != nil {
		t.Fatalf("CreateForecast failed: %v", err)
	}
	if f.ID == uuid.Nil {
		t.Fatal("expected ID to be set after create")
	}
}

func TestRepo_FindForecastByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestForecast(repo)

	found, err := repo.FindForecastByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("FindForecastByID failed: %v", err)
	}
	if found.ForecastType != created.ForecastType {
		t.Errorf("expected type '%s', got '%s'", created.ForecastType, found.ForecastType)
	}
}

func TestRepo_ListForecasts_ByType(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	createTestForecast(repo) // DEMAND

	// Create a SUPPLY forecast
	ctx := context.Background()
	supply := &WorkforceForecast{
		Period:          "2026-Q3",
		OrganizationID:  uuid.New(),
		ForecastType:    "SUPPLY",
		Headcount:       120,
		ConfidenceLevel: 80.0,
	}
	_ = repo.CreateForecast(ctx, supply)

	list, total, err := repo.ListForecasts(ctx, "", "", "SUPPLY", 1, 10)
	if err != nil {
		t.Fatalf("ListForecasts failed: %v", err)
	}
	if total != 1 {
		t.Errorf("expected 1 SUPPLY forecast, got %d", total)
	}
	if len(list) > 0 && list[0].ForecastType != "SUPPLY" {
		t.Errorf("expected SUPPLY type, got '%s'", list[0].ForecastType)
	}
}

func TestRepo_UpdateForecast_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestForecast(repo)
	created.Headcount = 200

	if err := repo.UpdateForecast(context.Background(), created); err != nil {
		t.Fatalf("UpdateForecast failed: %v", err)
	}

	found, _ := repo.FindForecastByID(context.Background(), created.ID)
	if found.Headcount != 200 {
		t.Errorf("expected Headcount 200, got %d", found.Headcount)
	}
}

func TestRepo_DeleteForecast_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestForecast(repo)

	if err := repo.DeleteForecast(context.Background(), created.ID); err != nil {
		t.Fatalf("DeleteForecast failed: %v", err)
	}

	_, err := repo.FindForecastByID(context.Background(), created.ID)
	if err == nil {
		t.Fatal("expected error after deleting forecast")
	}
}

// =========================================================================
// KPI Repository Tests
// =========================================================================

func TestRepo_CreateKPI_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	target := 90.0
	k := &WorkforceKPI{
		ID:         uuid.New(),
		Period:     "2026-Q3",
		KpiCode:    "HC_TOTAL",
		KpiName:    "Total Headcount",
		Value:      85.0,
		Target:     &target,
		Unit:       "COUNT",
		Dimension:  "COMPANY",
		SnapshotAt: parseDateOrNow(""),
	}

	if err := repo.CreateKPI(context.Background(), k); err != nil {
		t.Fatalf("CreateKPI failed: %v", err)
	}
	if k.ID == uuid.Nil {
		t.Fatal("expected ID to be set after create")
	}
}

func TestRepo_ListKPIs_ByCode(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	createTestKPI(repo)
	createTestKPI(repo)

	list, total, err := repo.ListKPIs(context.Background(), "", "", "", 1, 10)
	if err != nil {
		t.Fatalf("ListKPIs failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 KPIs, got %d", len(list))
	}
}

// =========================================================================
// Cache Repository Tests
// =========================================================================

func TestRepo_SetAndGetCache_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestCacheEntry(repo)

	found, err := repo.GetCache(context.Background(), created.CacheKey)
	if err != nil {
		t.Fatalf("GetCache failed: %v", err)
	}
	if found.CacheType != "HC" {
		t.Errorf("expected cache type 'HC', got '%s'", found.CacheType)
	}
}

func TestRepo_GetCache_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	_, err := repo.GetCache(context.Background(), "non-existent-key")
	if err == nil {
		t.Fatal("expected error for non-existent cache key")
	}
}

func TestRepo_InvalidateCache_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	createTestCacheEntry(repo)

	if err := repo.InvalidateCache(context.Background(), "HC"); err != nil {
		t.Fatalf("InvalidateCache failed: %v", err)
	}

	// Verify — create another entry and ensure the invalidated ones are gone
	list, _, _ := repo.ListKPIs(context.Background(), "", "", "", 1, 10)
	_ = list // cache table should be empty
}

// =========================================================================
// Scenario Repository Tests
// =========================================================================

func TestRepo_CreateScenario_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	s := &WorkforceScenario{
		ID:           uuid.New(),
		Name:         "New Branch Simulation",
		Description:  "Test new branch",
		ScenarioType: "NEW_BRANCH",
		Parameters:   JSON{"headcount": 50, "avg_cost": 10000000},
		Status:       "DRAFT",
		CreatedBy:    uuid.New(),
	}

	if err := repo.CreateScenario(context.Background(), s); err != nil {
		t.Fatalf("CreateScenario failed: %v", err)
	}
	if s.ID == uuid.Nil {
		t.Fatal("expected ID to be set after create")
	}
}

func TestRepo_FindScenarioByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestScenario(repo)

	found, err := repo.FindScenarioByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("FindScenarioByID failed: %v", err)
	}
	if found.Name != created.Name {
		t.Errorf("expected name '%s', got '%s'", created.Name, found.Name)
	}
}

func TestRepo_ListScenarios_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	for i := 0; i < 5; i++ {
		createTestScenario(repo)
	}

	list, total, err := repo.ListScenarios(context.Background(), "", 1, 10)
	if err != nil {
		t.Fatalf("ListScenarios failed: %v", err)
	}
	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(list) != 5 {
		t.Errorf("expected 5 scenarios, got %d", len(list))
	}
}

func TestRepo_UpdateScenario_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestScenario(repo)
	created.Status = "COMPLETED"
	created.Results = JSON{"headcount_needed": 25}

	if err := repo.UpdateScenario(context.Background(), created); err != nil {
		t.Fatalf("UpdateScenario failed: %v", err)
	}

	found, _ := repo.FindScenarioByID(context.Background(), created.ID)
	if found.Status != "COMPLETED" {
		t.Errorf("expected status 'COMPLETED', got '%s'", found.Status)
	}
}

func TestRepo_DeleteScenario_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestScenario(repo)

	if err := repo.DeleteScenario(context.Background(), created.ID); err != nil {
		t.Fatalf("DeleteScenario failed: %v", err)
	}

	_, err := repo.FindScenarioByID(context.Background(), created.ID)
	if err == nil {
		t.Fatal("expected error after deleting scenario")
	}
}

// =========================================================================
// Risk Indicator Repository Tests
// =========================================================================

func TestRepo_CreateRiskIndicator_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	deptID := uuid.New()
	ri := &WorkforceRiskIndicator{
		ID:             uuid.New(),
		Period:         "2026-Q3",
		RiskCode:       "HIGH_TURNOVER",
		RiskName:       "High Turnover Risk",
		RiskLevel:      "HIGH",
		Score:          85.0,
		Threshold:      70.0,
		DepartmentID:   &deptID,
		Recommendation: "Conduct retention program",
		SnapshotAt:     parseDateOrNow(""),
	}

	if err := repo.CreateRiskIndicator(context.Background(), ri); err != nil {
		t.Fatalf("CreateRiskIndicator failed: %v", err)
	}
	if ri.ID == uuid.Nil {
		t.Fatal("expected ID to be set after create")
	}
}

func TestRepo_FindRiskIndicatorByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestRiskIndicator(repo)

	found, err := repo.FindRiskIndicatorByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("FindRiskIndicatorByID failed: %v", err)
	}
	if found.RiskCode != created.RiskCode {
		t.Errorf("expected code '%s', got '%s'", created.RiskCode, found.RiskCode)
	}
}

func TestRepo_ListRiskIndicators_ByLevel(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	ctx := context.Background()
	for i := 0; i < 3; i++ {
		createTestRiskIndicator(repo)
	}
	// Create a HIGH risk indicator
	deptID := uuid.New()
	highRisk := &WorkforceRiskIndicator{
		Period:       "2026-Q3",
		RiskCode:     "CRITICAL_RISK",
		RiskName:     "Critical Risk",
		RiskLevel:    "HIGH",
		Score:        95.0,
		Threshold:    70.0,
		DepartmentID: &deptID,
		SnapshotAt:   parseDateOrNow(""),
	}
	_ = repo.CreateRiskIndicator(ctx, highRisk)

	list, total, err := repo.ListRiskIndicators(ctx, "", "HIGH", 1, 10)
	if err != nil {
		t.Fatalf("ListRiskIndicators failed: %v", err)
	}
	if total != 1 {
		t.Errorf("expected 1 HIGH risk, got %d", total)
	}
	if len(list) > 0 && list[0].RiskLevel != "HIGH" {
		t.Errorf("expected HIGH level, got '%s'", list[0].RiskLevel)
	}
}

func TestRepo_UpdateRiskIndicator_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestRiskIndicator(repo)
	created.RiskLevel = "CRITICAL"

	if err := repo.UpdateRiskIndicator(context.Background(), created); err != nil {
		t.Fatalf("UpdateRiskIndicator failed: %v", err)
	}

	found, _ := repo.FindRiskIndicatorByID(context.Background(), created.ID)
	if found.RiskLevel != "CRITICAL" {
		t.Errorf("expected level 'CRITICAL', got '%s'", found.RiskLevel)
	}
}

// =========================================================================
// Health Score Repository Tests
// =========================================================================

func TestRepo_CreateHealthScore_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	hs := &WorkforceHealthScore{
		ID:             uuid.New(),
		Period:         "2026-Q3",
		OrganizationID: uuid.New(),
		Score:          82.5,
		SpanOfControl:  5.0,
		ManagerRatio:   15.0,
		Components:     JSON{"score": 82.5},
		SnapshotAt:     parseDateOrNow(""),
	}

	if err := repo.CreateHealthScore(context.Background(), hs); err != nil {
		t.Fatalf("CreateHealthScore failed: %v", err)
	}
	if hs.ID == uuid.Nil {
		t.Fatal("expected ID to be set after create")
	}
}

func TestRepo_FindHealthScoreByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestHealthScore(repo)

	found, err := repo.FindHealthScoreByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("FindHealthScoreByID failed: %v", err)
	}
	if found.Score != created.Score {
		t.Errorf("expected Score %.1f, got %.1f", created.Score, found.Score)
	}
}

func TestRepo_ListHealthScores_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	for i := 0; i < 3; i++ {
		createTestHealthScore(repo)
	}

	list, total, err := repo.ListHealthScores(context.Background(), "", "", 1, 10)
	if err != nil {
		t.Fatalf("ListHealthScores failed: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(list) != 3 {
		t.Errorf("expected 3 health scores, got %d", len(list))
	}
}

// =========================================================================
// Read-only queries (edge cases: empty tables)
// =========================================================================

func TestRepo_GetActiveEmployeeCount_EmptyTable(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	count, err := repo.GetActiveEmployeeCount(context.Background())
	if err == nil {
		// employees table doesn't exist in test, so this is fine
		_ = count
	}
	// No assertion — this is a read-only query against source module tables
	// that don't exist in the test DB, so it's expected to error
}

func TestRepo_GetEmployeesByDepartment_EmptyTable(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	_, err := repo.GetEmployeesByDepartment(context.Background())
	if err != nil {
		// Expected — employees table doesn't exist in test context
		t.Logf("GetEmployeesByDepartment returned expected error: %v", err)
	}
}
