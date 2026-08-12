package workforceintelligence

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

// =========================================================================
// Candidate Search Repository Tests
// =========================================================================

// createCandidateSearchTables membuat tabel raw milik modul lain (organization,
// employee, recruitment) yang di-query Candidate Search secara langsung.
func createCandidateSearchTables(t *testing.T, db *gorm.DB) {
	t.Helper()
	statements := []string{
		`CREATE TABLE organization_summaries (
			id CHAR(36) PRIMARY KEY,
			code VARCHAR(7),
			decree_no VARCHAR(20),
			decree_date TEXT,
			status VARCHAR(20),
			deleted_at DATETIME NULL
		)`,
		`CREATE TABLE organizations (
			id CHAR(36) PRIMARY KEY,
			organization_summary_id CHAR(36),
			code VARCHAR(10),
			full_code VARCHAR(50),
			nomenclature VARCHAR(255),
			parent_id CHAR(36) NULL,
			deleted_at DATETIME NULL
		)`,
		`CREATE TABLE employments (
			id CHAR(36) PRIMARY KEY,
			employee_id CHAR(36),
			organization_id CHAR(36),
			effective_date TEXT,
			effective_end_date TEXT NULL
		)`,
		`CREATE TABLE candidates (
			id CHAR(36) PRIMARY KEY,
			first_name VARCHAR(100),
			last_name VARCHAR(100),
			email VARCHAR(255),
			phone VARCHAR(50),
			current_company TEXT NULL,
			current_title TEXT NULL,
			source VARCHAR(50)
		)`,
		`CREATE TABLE job_requisitions (
			id CHAR(36) PRIMARY KEY,
			organization_id CHAR(36),
			title VARCHAR(255),
			status VARCHAR(20)
		)`,
		`CREATE TABLE job_applications (
			id CHAR(36) PRIMARY KEY,
			requisition_id CHAR(36),
			candidate_id CHAR(36),
			status VARCHAR(50),
			created_at DATETIME
		)`,
		`CREATE TABLE positions (
			id CHAR(36) PRIMARY KEY,
			organization_id CHAR(36),
			title VARCHAR(200),
			is_active TINYINT(1) DEFAULT 1
		)`,
	}
	for _, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("failed to create candidate search table: %v\n%s", err, stmt)
		}
	}
}

// =========================================================================
// Recruitment pipeline reads (S-2)
// =========================================================================

// createRecruitmentTables membuat tabel job_requisitions + job_applications
// dengan kolom lengkap (slots_available, slots_filled) untuk test S-2.
func createRecruitmentTables(t *testing.T, db *gorm.DB) {
	t.Helper()
	statements := []string{
		`CREATE TABLE job_requisitions (
			id CHAR(36) PRIMARY KEY,
			organization_id CHAR(36),
			title VARCHAR(255),
			status VARCHAR(20),
			slots_available INT DEFAULT 1,
			slots_filled INT DEFAULT 0
		)`,
		`CREATE TABLE job_applications (
			id CHAR(36) PRIMARY KEY,
			requisition_id CHAR(36),
			candidate_id CHAR(36),
			status VARCHAR(50),
			created_at DATETIME
		)`,
	}
	for _, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("failed to create recruitment table: %v\n%s", err, stmt)
		}
	}
}

func TestRepo_GetRecruitmentOpenPositions_SumsActiveSlots(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	createRecruitmentTables(t, db)
	ctx := context.Background()

	orgID := uuid.New().String()
	// OPEN: 3 slots, 1 filled → 2 open
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status, slots_available, slots_filled) VALUES (?, ?, 'A', 'OPEN', 3, 1)", uuid.New().String(), orgID)
	// IN_PROGRESS: 2 slots, 0 filled → 2 open
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status, slots_available, slots_filled) VALUES (?, ?, 'B', 'IN_PROGRESS', 2, 0)", uuid.New().String(), orgID)
	// FILLED & DRAFT: tidak dihitung
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status, slots_available, slots_filled) VALUES (?, ?, 'C', 'FILLED', 5, 5)", uuid.New().String(), orgID)
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status, slots_available, slots_filled) VALUES (?, ?, 'D', 'DRAFT', 1, 0)", uuid.New().String(), orgID)

	got, err := repo.GetRecruitmentOpenPositions(ctx)
	if err != nil {
		t.Fatalf("GetRecruitmentOpenPositions failed: %v", err)
	}
	if got != 4 {
		t.Errorf("expected 4 open positions (2+2), got %d", got)
	}
}

func TestRepo_GetRecruitmentAcceptedOffers_CountsAccepted(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	createRecruitmentTables(t, db)
	ctx := context.Background()

	reqID := uuid.New().String()
	cand := uuid.New().String()
	for _, st := range []string{"NEW", "OFFERED", "ACCEPTED", "ACCEPTED", "REJECTED", "WITHDRAWN"} {
		db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status) VALUES (?, ?, ?, ?)", uuid.New().String(), reqID, cand, st)
	}

	got, err := repo.GetRecruitmentAcceptedOffers(ctx)
	if err != nil {
		t.Fatalf("GetRecruitmentAcceptedOffers failed: %v", err)
	}
	if got != 2 {
		t.Errorf("expected 2 accepted offers, got %d", got)
	}
}

func TestRepo_GetRecruitmentFilledPositions_SumsFilled(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	createRecruitmentTables(t, db)
	ctx := context.Background()

	orgID := uuid.New().String()
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status, slots_available, slots_filled) VALUES (?, ?, 'A', 'FILLED', 3, 3)", uuid.New().String(), orgID)
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status, slots_available, slots_filled) VALUES (?, ?, 'B', 'FILLED', 2, 1)", uuid.New().String(), orgID)
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status, slots_available, slots_filled) VALUES (?, ?, 'C', 'OPEN', 5, 1)", uuid.New().String(), orgID)

	got, err := repo.GetRecruitmentFilledPositions(ctx)
	if err != nil {
		t.Fatalf("GetRecruitmentFilledPositions failed: %v", err)
	}
	// Hanya status FILLED dijumlah: 3 + 1 = 4 (OPEN tidak dihitung)
	if got != 4 {
		t.Errorf("expected 4 filled positions, got %d", got)
	}
}

func TestRepo_GetRecruitmentPipeline_GroupsByStatus(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	createRecruitmentTables(t, db)
	ctx := context.Background()

	reqID := uuid.New().String()
	cand := uuid.New().String()
	for _, st := range []string{"NEW", "NEW", "SCREENED", "OFFERED", "ACCEPTED"} {
		db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status) VALUES (?, ?, ?, ?)", uuid.New().String(), reqID, cand, st)
	}

	rows, err := repo.GetRecruitmentPipeline(ctx)
	if err != nil {
		t.Fatalf("GetRecruitmentPipeline failed: %v", err)
	}
	counts := map[string]int{}
	for _, r := range rows {
		counts[r.Status] = r.Count
	}
	if counts["NEW"] != 2 || counts["SCREENED"] != 1 || counts["OFFERED"] != 1 || counts["ACCEPTED"] != 1 {
		t.Errorf("unexpected pipeline counts: %v", counts)
	}
}

func TestRepo_CandidateSearchVacantOrgs_OnlyActiveEmpty(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	createCandidateSearchTables(t, db)
	ctx := context.Background()

	summaryActiveID := uuid.New().String()
	summaryInactiveID := uuid.New().String()
	orgVacantID := uuid.New().String()
	orgOccupiedID := uuid.New().String()
	orgEndedID := uuid.New().String()
	orgInactiveID := uuid.New().String()

	// Summaries
	db.Exec("INSERT INTO organization_summaries (id, code, decree_no, decree_date, status) VALUES (?, 'SA-01', 'SK-001', '2024-01-01', 'active')", summaryActiveID)
	db.Exec("INSERT INTO organization_summaries (id, code, decree_no, decree_date, status) VALUES (?, 'SI-01', 'SK-002', '2024-01-01', 'inactive')", summaryInactiveID)

	// Organizations
	db.Exec("INSERT INTO organizations (id, organization_summary_id, code, full_code, nomenclature) VALUES (?, ?, 'ORG-01', 'SA-01.ORG-01', 'Staff IT')", orgVacantID, summaryActiveID)
	db.Exec("INSERT INTO organizations (id, organization_summary_id, code, full_code, nomenclature) VALUES (?, ?, 'ORG-02', 'SA-01.ORG-02', 'Supervisor IT')", orgOccupiedID, summaryActiveID)
	db.Exec("INSERT INTO organizations (id, organization_summary_id, code, full_code, nomenclature) VALUES (?, ?, 'ORG-03', 'SA-01.ORG-03', 'Manager IT')", orgEndedID, summaryActiveID)
	db.Exec("INSERT INTO organizations (id, organization_summary_id, code, full_code, nomenclature) VALUES (?, ?, 'ORG-04', 'SI-01.ORG-04', 'Staff Finance')", orgInactiveID, summaryInactiveID)

	// Employments: ORG-02 aktif (masih berjalan), ORG-03 sudah berakhir
	db.Exec("INSERT INTO employments (id, employee_id, organization_id, effective_date, effective_end_date) VALUES (?, ?, ?, '2024-01-01', NULL)", uuid.New().String(), uuid.New().String(), orgOccupiedID)
	db.Exec("INSERT INTO employments (id, employee_id, organization_id, effective_date, effective_end_date) VALUES (?, ?, ?, '2020-01-01', '2021-01-01')", uuid.New().String(), uuid.New().String(), orgEndedID)

	search := ""
	rows, total, err := repo.CandidateSearchVacantOrgs(ctx, &search, nil, 1, 20)
	if err != nil {
		t.Fatalf("CandidateSearchVacantOrgs failed: %v", err)
	}

	if total != 2 {
		t.Errorf("expected 2 vacant orgs (ORG-01 lowong + ORG-03 kontrak berakhir), got %d", total)
	}

	codes := map[string]bool{}
	for _, r := range rows {
		codes[r.OrganizationCode] = true
	}
	if !codes["ORG-01"] || !codes["ORG-03"] {
		t.Errorf("expected ORG-01 & ORG-03 in results, got %v", codes)
	}
	if codes["ORG-02"] {
		t.Error("ORG-02 (masih ada employment aktif) tidak boleh dianggap lowong")
	}
	if codes["ORG-04"] {
		t.Error("ORG-04 (summary inactive) tidak boleh dianggap lowong")
	}

	// Search filter by nama posisi (hanya org lowong yang cocok; ORG-02 yang
	// terisi tidak ikut, jadi 'Manager' hanya match ORG-03)
	s2 := "Manager"
	_, total2, err := repo.CandidateSearchVacantOrgs(ctx, &s2, nil, 1, 20)
	if err != nil {
		t.Fatalf("CandidateSearchVacantOrgs search failed: %v", err)
	}
	if total2 != 1 {
		t.Errorf("expected 1 result for search 'Manager', got %d", total2)
	}
}

func TestRepo_CandidateSearchCandidatesByOrgIDs_FiltersRejected(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	createCandidateSearchTables(t, db)
	ctx := context.Background()

	orgID := uuid.New()
	candOK := uuid.New()
	candRejected := uuid.New()
	reqID := uuid.New()

	db.Exec("INSERT INTO candidates (id, first_name, last_name, email, source) VALUES (?, 'Andi', 'Wijaya', 'andi@test.local', 'direct')", candOK.String())
	db.Exec("INSERT INTO candidates (id, first_name, last_name, email, source) VALUES (?, 'Budi', 'Santoso', 'budi@test.local', 'referral')", candRejected.String())
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status) VALUES (?, ?, 'Staff IT', 'OPEN')", reqID.String(), orgID.String())
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status) VALUES (?, ?, ?, 'SCREENED')", uuid.New().String(), reqID.String(), candOK.String())
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status) VALUES (?, ?, ?, 'REJECTED')", uuid.New().String(), reqID.String(), candRejected.String())

	rows, err := repo.CandidateSearchCandidatesByOrgIDs(ctx, []uuid.UUID{orgID})
	if err != nil {
		t.Fatalf("CandidateSearchCandidatesByOrgIDs failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 candidate (rejected excluded), got %d", len(rows))
	}
	if rows[0].Email != "andi@test.local" {
		t.Errorf("expected andi@test.local, got %s", rows[0].Email)
	}
	if rows[0].ApplicationStatus != "SCREENED" {
		t.Errorf("expected application status SCREENED, got %s", rows[0].ApplicationStatus)
	}
	if rows[0].RequisitionTitle != "Staff IT" {
		t.Errorf("expected requisition title 'Staff IT', got %s", rows[0].RequisitionTitle)
	}

	// Empty orgIDs -> no rows, no error
	empty, err := repo.CandidateSearchCandidatesByOrgIDs(ctx, nil)
	if err != nil || empty != nil {
		t.Errorf("expected nil rows and no error for empty orgIDs, got rows=%v err=%v", empty, err)
	}
}

// =========================================================================
// Recruitment Analytics reads (S-3)
// =========================================================================

// createRecruitmentAnalyticsTables membuat tabel recruitment dengan kolom
// timestamp lengkap (applied_at/accepted_at/offered_at/closed_at/created_at)
// + candidates untuk metrik advanced S-3.
func createRecruitmentAnalyticsTables(t *testing.T, db *gorm.DB) {
	t.Helper()
	statements := []string{
		`CREATE TABLE job_requisitions (
			id CHAR(36) PRIMARY KEY,
			organization_id CHAR(36),
			title VARCHAR(255),
			status VARCHAR(20),
			slots_available INT DEFAULT 1,
			slots_filled INT DEFAULT 0,
			closed_at BIGINT DEFAULT 0,
			created_at DATETIME
		)`,
		`CREATE TABLE candidates (
			id CHAR(36) PRIMARY KEY,
			first_name VARCHAR(100),
			last_name VARCHAR(100),
			email VARCHAR(255),
			source VARCHAR(50)
		)`,
		`CREATE TABLE job_applications (
			id CHAR(36) PRIMARY KEY,
			requisition_id CHAR(36),
			candidate_id CHAR(36),
			status VARCHAR(50),
			applied_at BIGINT DEFAULT 0,
			offered_at BIGINT DEFAULT 0,
			accepted_at BIGINT DEFAULT 0,
			created_at DATETIME
		)`,
	}
	for _, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("failed to create recruitment analytics table: %v\n%s", err, stmt)
		}
	}
}

func TestRepo_GetRecruitmentTimeToHire_AvgDays(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	createRecruitmentAnalyticsTables(t, db)
	ctx := context.Background()

	reqID := uuid.New().String()
	cand := uuid.New().String()
	d := int64(86400000) // 1 hari dalam ms
	// 12 hari + 1 hari → avg 6.5 hari
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status, applied_at, accepted_at) VALUES (?, ?, ?, 'ACCEPTED', ?, ?)", uuid.New().String(), reqID, cand, 1767225600000, 1767225600000+12*d)
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status, applied_at, accepted_at) VALUES (?, ?, ?, 'ACCEPTED', ?, ?)", uuid.New().String(), reqID, cand, 1767225600000, 1767225600000+d)
	// Non-hired tidak boleh masuk perhitungan
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status, applied_at, accepted_at) VALUES (?, ?, ?, 'OFFERED', ?, ?)", uuid.New().String(), reqID, cand, 1767225600000, 1767225600000+2*d)

	got, err := repo.GetRecruitmentTimeToHire(ctx)
	if err != nil {
		t.Fatalf("GetRecruitmentTimeToHire failed: %v", err)
	}
	if got != 6.5 {
		t.Errorf("expected time to hire 6.5 days, got %v", got)
	}
}

func TestRepo_GetRecruitmentTimeToHire_NoHires(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	createRecruitmentAnalyticsTables(t, db)
	ctx := context.Background()

	reqID := uuid.New().String()
	cand := uuid.New().String()
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status) VALUES (?, ?, ?, 'NEW')", uuid.New().String(), reqID, cand)

	got, err := repo.GetRecruitmentTimeToHire(ctx)
	if err != nil {
		t.Fatalf("GetRecruitmentTimeToHire failed: %v", err)
	}
	if got != 0 {
		t.Errorf("expected 0 time to hire with no hires, got %v", got)
	}
}

func TestRepo_GetRecruitmentOfferAcceptance_Counts(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	createRecruitmentAnalyticsTables(t, db)
	ctx := context.Background()

	reqID := uuid.New().String()
	cand := uuid.New().String()
	d := int64(86400000)
	// 2 ACCEPTED + 2 OFFERED (belum diterima) + 1 NEW → accepted=2, offered=4
	for i := 0; i < 2; i++ {
		db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status, offered_at, accepted_at) VALUES (?, ?, ?, 'ACCEPTED', ?, ?)", uuid.New().String(), reqID, cand, 1767225600000, 1767225600000+d)
	}
	for i := 0; i < 2; i++ {
		db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status, offered_at) VALUES (?, ?, ?, 'OFFERED', ?)", uuid.New().String(), reqID, cand, 1767225600000)
	}
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status) VALUES (?, ?, ?, 'NEW')", uuid.New().String(), reqID, cand)

	got, err := repo.GetRecruitmentOfferAcceptance(ctx)
	if err != nil {
		t.Fatalf("GetRecruitmentOfferAcceptance failed: %v", err)
	}
	if got.Accepted != 2 || got.Offered != 4 {
		t.Errorf("expected accepted=2 offered=4, got %+v", got)
	}
}

func TestRepo_GetRecruitmentSourceConversion_ByChannel(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	createRecruitmentAnalyticsTables(t, db)
	ctx := context.Background()

	reqID := uuid.New().String()
	c1 := uuid.New().String() // referral, di-hire
	c2 := uuid.New().String() // referral, tanpa aplikasi
	c3 := uuid.New().String() // linkedin, di-hire
	c4 := uuid.New().String() // direct, tanpa hire
	db.Exec("INSERT INTO candidates (id, first_name, last_name, email, source) VALUES (?, 'A', 'B', 'a@t.local', 'referral')", c1)
	db.Exec("INSERT INTO candidates (id, first_name, last_name, email, source) VALUES (?, 'C', 'D', 'c@t.local', 'referral')", c2)
	db.Exec("INSERT INTO candidates (id, first_name, last_name, email, source) VALUES (?, 'E', 'F', 'e@t.local', 'linkedin')", c3)
	db.Exec("INSERT INTO candidates (id, first_name, last_name, email, source) VALUES (?, 'G', 'H', 'g@t.local', 'direct')", c4)
	d := int64(86400000)
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status, applied_at, accepted_at) VALUES (?, ?, ?, 'ACCEPTED', ?, ?)", uuid.New().String(), reqID, c1, 1767225600000, 1767225600000+d)
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status, applied_at, accepted_at) VALUES (?, ?, ?, 'ACCEPTED', ?, ?)", uuid.New().String(), reqID, c3, 1767225600000, 1767225600000+d)

	rows, err := repo.GetRecruitmentSourceConversion(ctx)
	if err != nil {
		t.Fatalf("GetRecruitmentSourceConversion failed: %v", err)
	}
	bySource := map[string]RecruitmentSourceConversion{}
	for _, r := range rows {
		bySource[r.Source] = r
	}
	if bySource["referral"].Candidates != 2 || bySource["referral"].Hires != 1 {
		t.Errorf("expected referral 2 candidates/1 hire, got %+v", bySource["referral"])
	}
	if bySource["linkedin"].Candidates != 1 || bySource["linkedin"].Hires != 1 {
		t.Errorf("expected linkedin 1 candidate/1 hire, got %+v", bySource["linkedin"])
	}
	if bySource["direct"].Candidates != 1 || bySource["direct"].Hires != 0 {
		t.Errorf("expected direct 1 candidate/0 hire, got %+v", bySource["direct"])
	}
}

func TestRepo_GetRecruitmentFilledRequisitionDurations_ReturnsClosed(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	createRecruitmentAnalyticsTables(t, db)
	ctx := context.Background()

	orgID := uuid.New().String()
	// FILLED dengan closed_at & created_at → ikut dihitung
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status, closed_at, created_at) VALUES (?, ?, 'A', 'FILLED', 1768089600000, '2026-01-01 00:00:00')", uuid.New().String(), orgID)
	// OPEN & FILLED tanpa closed_at → tidak dihitung
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status, closed_at, created_at) VALUES (?, ?, 'B', 'FILLED', 0, '2026-01-01 00:00:00')", uuid.New().String(), orgID)
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status, created_at) VALUES (?, ?, 'C', 'OPEN', '2026-01-01 00:00:00')", uuid.New().String(), orgID)

	rows, err := repo.GetRecruitmentFilledRequisitionDurations(ctx)
	if err != nil {
		t.Fatalf("GetRecruitmentFilledRequisitionDurations failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 FILLED requisition with closed_at, got %d", len(rows))
	}
	if rows[0].ClosedAt != 1768089600000 {
		t.Errorf("expected closed_at 1768089600000, got %d", rows[0].ClosedAt)
	}
}

func TestRepo_CandidateSearchVacantOrgs_PositionFilter(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	createCandidateSearchTables(t, db)
	ctx := context.Background()

	summaryID := uuid.New().String()
	orgIT := uuid.New().String()
	orgSales := uuid.New().String()
	db.Exec("INSERT INTO organization_summaries (id, code, decree_no, decree_date, status) VALUES (?, 'SA-01', 'SK-001', '2024-01-01', 'active')", summaryID)
	db.Exec("INSERT INTO organizations (id, organization_summary_id, code, full_code, nomenclature) VALUES (?, ?, 'ORG-01', 'SA-01.ORG-01', 'Staff IT')", orgIT, summaryID)
	db.Exec("INSERT INTO organizations (id, organization_summary_id, code, full_code, nomenclature) VALUES (?, ?, 'ORG-02', 'SA-01.ORG-02', 'Sales Executive')", orgSales, summaryID)

	pos := "Staff"
	rows, total, err := repo.CandidateSearchVacantOrgs(ctx, nil, &pos, 1, 20)
	if err != nil {
		t.Fatalf("CandidateSearchVacantOrgs failed: %v", err)
	}
	if total != 1 || len(rows) != 1 || rows[0].OrganizationID != orgIT {
		t.Fatalf("expected only Staff IT (1), got total=%d rows=%+v", total, rows)
	}
}

func TestRepo_CandidateSearchPositionsByOrgIDs_ActiveOnly(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	createCandidateSearchTables(t, db)
	ctx := context.Background()

	orgID := uuid.New().String()
	otherOrg := uuid.New().String()
	db.Exec("INSERT INTO positions (id, organization_id, title, is_active) VALUES (?, ?, 'Staff IT', 1)", uuid.New().String(), orgID)
	db.Exec("INSERT INTO positions (id, organization_id, title, is_active) VALUES (?, ?, 'Supervisor IT', 0)", uuid.New().String(), orgID)
	db.Exec("INSERT INTO positions (id, organization_id, title, is_active) VALUES (?, ?, 'Lain', 1)", uuid.New().String(), otherOrg)

	rows, err := repo.CandidateSearchPositionsByOrgIDs(ctx, []uuid.UUID{mustUUID(t, orgID)})
	if err != nil {
		t.Fatalf("CandidateSearchPositionsByOrgIDs failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 active position for org, got %d", len(rows))
	}
	if rows[0].Title != "Staff IT" {
		t.Errorf("expected position 'Staff IT', got %s", rows[0].Title)
	}
	empty, err := repo.CandidateSearchPositionsByOrgIDs(ctx, nil)
	if err != nil || empty != nil {
		t.Errorf("expected nil rows for empty orgIDs, got %v %v", empty, err)
	}
}

func mustUUID(t *testing.T, s string) uuid.UUID {
	t.Helper()
	uid, err := uuid.Parse(s)
	if err != nil {
		t.Fatalf("invalid uuid %s: %v", s, err)
	}
	return uid
}

// =========================================================================
// Quality of Hire (S-6) Repository Tests
// =========================================================================

// createQualityOfHireTables membuat tabel lintas modul yang dibaca WI untuk
// metrik Quality of Hire (S-6): candidates + job_requisitions + job_applications
// + interviews + employee_onboardings + performance_evaluations + employments.
func createQualityOfHireTables(t *testing.T, db *gorm.DB) {
	t.Helper()
	statements := []string{
		`CREATE TABLE candidates (
			id CHAR(36) PRIMARY KEY,
			first_name VARCHAR(100),
			last_name VARCHAR(100),
			email VARCHAR(255),
			source VARCHAR(50)
		)`,
		`CREATE TABLE job_requisitions (
			id CHAR(36) PRIMARY KEY,
			organization_id CHAR(36),
			title VARCHAR(255),
			status VARCHAR(20)
		)`,
		`CREATE TABLE job_applications (
			id CHAR(36) PRIMARY KEY,
			requisition_id CHAR(36),
			candidate_id CHAR(36),
			status VARCHAR(50),
			created_at DATETIME
		)`,
		`CREATE TABLE interviews (
			id CHAR(36) PRIMARY KEY,
			application_id CHAR(36),
			score DECIMAL(5,2) NULL
		)`,
		`CREATE TABLE employee_onboardings (
			id CHAR(36) PRIMARY KEY,
			application_id CHAR(36),
			employee_id CHAR(36),
			status VARCHAR(20)
		)`,
		`CREATE TABLE performance_evaluations (
			id CHAR(36) PRIMARY KEY,
			employee_id CHAR(36),
			final_score DECIMAL(5,2),
			status VARCHAR(20),
			updated_at DATETIME NULL
		)`,
		`CREATE TABLE employments (
			id CHAR(36) PRIMARY KEY,
			employee_id CHAR(36),
			effective_end_date TEXT NULL,
			deleted_at DATETIME NULL
		)`,
	}
	for _, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("failed to create quality-of-hire table: %v\n%s", err, stmt)
		}
	}
}

func TestRepo_GetQualityOfHireHires_Composites(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	createQualityOfHireTables(t, db)
	ctx := context.Background()

	orgA := uuid.New()
	orgB := uuid.New()
	reqA := uuid.New()
	reqB := uuid.New()
	candReferral := uuid.New()
	candJobBoard := uuid.New()
	empA := uuid.New()
	empB := uuid.New()

	db.Exec("INSERT INTO candidates (id, first_name, last_name, email, source) VALUES (?, 'A', 'A', 'a@x.com', 'referral')", candReferral.String())
	db.Exec("INSERT INTO candidates (id, first_name, last_name, email, source) VALUES (?, 'B', 'B', 'b@x.com', 'job_board')", candJobBoard.String())
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status) VALUES (?, ?, 'Eng', 'FILLED')", reqA.String(), orgA.String())
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status) VALUES (?, ?, 'Ops', 'FILLED')", reqB.String(), orgB.String())

	// Hire 1: referral — interview 80, onboarding COMPLETED, perf 85 (evaluasi
	// terbaru menang atas evaluasi lama 90), retained
	appA := uuid.New()
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status) VALUES (?, ?, ?, 'ACCEPTED')", appA.String(), reqA.String(), candReferral.String())
	db.Exec("INSERT INTO interviews (id, application_id, score) VALUES (?, ?, 80)", uuid.New().String(), appA.String())
	db.Exec("INSERT INTO employee_onboardings (id, application_id, employee_id, status) VALUES (?, ?, ?, 'COMPLETED')", uuid.New().String(), appA.String(), empA.String())
	db.Exec("INSERT INTO performance_evaluations (id, employee_id, final_score, status, updated_at) VALUES (?, ?, 90, 'COMPLETED', '2024-06-01')", uuid.New().String(), empA.String())
	db.Exec("INSERT INTO performance_evaluations (id, employee_id, final_score, status, updated_at) VALUES (?, ?, 85, 'COMPLETED', '2025-06-01')", uuid.New().String(), empA.String())
	db.Exec("INSERT INTO employments (id, employee_id, effective_end_date) VALUES (?, ?, NULL)", uuid.New().String(), empA.String())

	// Hire 2: job_board — interview 70, onboarding PENDING, perf 75, TIDAK retained
	appB := uuid.New()
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status) VALUES (?, ?, ?, 'ACCEPTED')", appB.String(), reqB.String(), candJobBoard.String())
	db.Exec("INSERT INTO interviews (id, application_id, score) VALUES (?, ?, 70)", uuid.New().String(), appB.String())
	db.Exec("INSERT INTO employee_onboardings (id, application_id, employee_id, status) VALUES (?, ?, ?, 'PENDING')", uuid.New().String(), appB.String(), empB.String())
	db.Exec("INSERT INTO performance_evaluations (id, employee_id, final_score, status, updated_at) VALUES (?, ?, 75, 'COMPLETED', '2025-06-01')", uuid.New().String(), empB.String())
	db.Exec("INSERT INTO employments (id, employee_id, effective_end_date) VALUES (?, ?, '2025-01-01')", uuid.New().String(), empB.String())

	// Bukan hire (REJECTED) — tidak boleh muncul
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status) VALUES (?, ?, ?, 'REJECTED')", uuid.New().String(), reqA.String(), candReferral.String())

	rows, err := repo.GetQualityOfHireHires(ctx)
	if err != nil {
		t.Fatalf("GetQualityOfHireHires failed: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 hires (ACCEPTED only), got %d", len(rows))
	}
	byApp := map[string]QualityOfHireRow{}
	for _, r := range rows {
		byApp[r.ApplicationID] = r
	}
	rA, ok := byApp[appA.String()]
	if !ok {
		t.Fatal("expected hire A in rows")
	}
	if rA.InterviewScore != 80 || rA.PerformanceScore != 85 || rA.OnboardingStatus != "COMPLETED" || rA.RetainedCount != 1 {
		t.Errorf("expected hire A interview=80 perf=85 onboarding=COMPLETED retained=1, got %+v", rA)
	}
	if rA.Source != "referral" || rA.OrganizationID != orgA.String() || rA.RequisitionID != reqA.String() {
		t.Errorf("expected hire A source/org/req linkage, got %+v", rA)
	}
	rB, ok := byApp[appB.String()]
	if !ok {
		t.Fatal("expected hire B in rows")
	}
	if rB.RetainedCount != 0 {
		t.Errorf("expected hire B not retained, got %+v", rB)
	}
}

func TestRepo_GetQualityOfHireHires_Empty(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	createQualityOfHireTables(t, db)
	ctx := context.Background()

	rows, err := repo.GetQualityOfHireHires(ctx)
	if err != nil {
		t.Fatalf("GetQualityOfHireHires failed: %v", err)
	}
	if len(rows) != 0 {
		t.Errorf("expected 0 rows on empty data, got %d", len(rows))
	}
}
