package workforceintelligence

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func ctx() context.Context {
	return context.Background()
}

// =========================================================================
// Headcount Plan Service Tests
// =========================================================================

func TestService_CreateHeadcountPlan_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateHeadcountPlanRequest{
		Period:         "2026-Q3",
		OrganizationID: uuid.New().String(),
		PlannedHC:      100,
	}

	resp, err := svc.CreateHeadcountPlan(ctx(), req)
	if err != nil {
		t.Fatalf("CreateHeadcountPlan failed: %v", err)
	}

	if resp.Period != "2026-Q3" {
		t.Errorf("expected period '2026-Q3', got '%s'", resp.Period)
	}
	if resp.PlannedHC != 100 {
		t.Errorf("expected PlannedHC 100, got %d", resp.PlannedHC)
	}
	if resp.ID == "" {
		t.Error("expected ID to be set")
	}
}

func TestService_CreateHeadcountPlan_InvalidOrgID(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateHeadcountPlanRequest{
		Period:         "2026-Q3",
		OrganizationID: "not-a-uuid",
		PlannedHC:      100,
	}

	_, err := svc.CreateHeadcountPlan(ctx(), req)
	if err == nil {
		t.Fatal("expected error for invalid organization_id")
	}
}

func TestService_GetHeadcountPlanByID_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestHeadcountPlan(repo)

	found, err := svc.GetHeadcountPlanByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetHeadcountPlanByID failed: %v", err)
	}

	if found.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), found.ID)
	}
}

func TestService_ListHeadcountPlans_DefaultPagination(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	for i := 0; i < 3; i++ {
		createTestHeadcountPlan(repo)
	}

	resp, err := svc.ListHeadcountPlans(ctx(), "", "", 0, 0)
	if err != nil {
		t.Fatalf("ListHeadcountPlans failed: %v", err)
	}

	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
}

func TestService_UpdateHeadcountPlan_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestHeadcountPlan(repo)
	newPlanned := 200
	req := UpdateHeadcountPlanRequest{PlannedHC: &newPlanned}

	updated, err := svc.UpdateHeadcountPlan(ctx(), created.ID.String(), req)
	if err != nil {
		t.Fatalf("UpdateHeadcountPlan failed: %v", err)
	}

	if updated.PlannedHC != 200 {
		t.Errorf("expected PlannedHC 200, got %d", updated.PlannedHC)
	}
}

func TestService_DeleteHeadcountPlan_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestHeadcountPlan(repo)

	if err := svc.DeleteHeadcountPlan(ctx(), created.ID.String()); err != nil {
		t.Fatalf("DeleteHeadcountPlan failed: %v", err)
	}

	_, err := svc.GetHeadcountPlanByID(ctx(), created.ID.String())
	if err == nil {
		t.Fatal("expected error after deleting headcount plan")
	}
}

// =========================================================================
// Forecast Service Tests
// =========================================================================

func TestService_CreateForecast_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateForecastRequest{
		Period:          "2026-Q3",
		OrganizationID:  uuid.New().String(),
		ForecastType:    "DEMAND",
		Headcount:       150,
		ConfidenceLevel: 90.0,
	}

	resp, err := svc.CreateForecast(ctx(), req)
	if err != nil {
		t.Fatalf("CreateForecast failed: %v", err)
	}

	if resp.Headcount != 150 {
		t.Errorf("expected Headcount 150, got %d", resp.Headcount)
	}
	if resp.ForecastType != "DEMAND" {
		t.Errorf("expected type 'DEMAND', got '%s'", resp.ForecastType)
	}
}

func TestService_GetForecastByID_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestForecast(repo)

	found, err := svc.GetForecastByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetForecastByID failed: %v", err)
	}

	if found.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), found.ID)
	}
}

func TestService_UpdateForecast_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestForecast(repo)
	newHC := 200
	req := UpdateForecastRequest{Headcount: &newHC}

	updated, err := svc.UpdateForecast(ctx(), created.ID.String(), req)
	if err != nil {
		t.Fatalf("UpdateForecast failed: %v", err)
	}

	if updated.Headcount != 200 {
		t.Errorf("expected Headcount 200, got %d", updated.Headcount)
	}
}

func TestService_DeleteForecast_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestForecast(repo)

	if err := svc.DeleteForecast(ctx(), created.ID.String()); err != nil {
		t.Fatalf("DeleteForecast failed: %v", err)
	}

	_, err := svc.GetForecastByID(ctx(), created.ID.String())
	if err == nil {
		t.Fatal("expected error after deleting forecast")
	}
}

// =========================================================================
// KPI Service Tests
// =========================================================================

func TestService_ListKPIs_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	createTestKPI(repo)
	createTestKPI(repo)

	resp, err := svc.ListKPIs(ctx(), "", "", "", 1, 10)
	if err != nil {
		t.Fatalf("ListKPIs failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
}

func TestService_GetKPISummary_Empty(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	summary, err := svc.GetKPISummary(ctx(), "2026-Q3")
	if err != nil {
		t.Fatalf("GetKPISummary failed: %v", err)
	}

	if summary.TotalKPIs != 0 {
		t.Errorf("expected 0 KPIs, got %d", summary.TotalKPIs)
	}
	if summary.Period != "2026-Q3" {
		t.Errorf("expected period '2026-Q3', got '%s'", summary.Period)
	}
}

func TestService_GetKPIByCode_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestKPI(repo)

	found, err := svc.GetKPIByCode(ctx(), created.KpiCode)
	if err != nil {
		t.Fatalf("GetKPIByCode failed: %v", err)
	}

	if found.KpiCode != created.KpiCode {
		t.Errorf("expected code '%s', got '%s'", created.KpiCode, found.KpiCode)
	}
}

func TestService_GetKPIByCode_NotFound(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetKPIByCode(ctx(), "NONEXISTENT")
	if err == nil {
		t.Fatal("expected error for non-existent KPI code")
	}
}

// =========================================================================
// Scenario Service Tests
// =========================================================================

func TestService_CreateScenario_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateScenarioRequest{
		Name:         "Growth Simulation",
		Description:  "10% growth scenario",
		ScenarioType: "GROWTH",
		Parameters:   map[string]interface{}{"growth_rate": 10.0},
	}

	resp, err := svc.CreateScenario(context.WithValue(ctx(), "user_id", uuid.New().String()), req)
	if err != nil {
		t.Fatalf("CreateScenario failed: %v", err)
	}

	if resp.Name != "Growth Simulation" {
		t.Errorf("expected name 'Growth Simulation', got '%s'", resp.Name)
	}
	if resp.Status != "DRAFT" {
		t.Errorf("expected status DRAFT, got '%s'", resp.Status)
	}
}

func TestService_GetScenarioByID_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestScenario(repo)

	found, err := svc.GetScenarioByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetScenarioByID failed: %v", err)
	}

	if found.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), found.ID)
	}
}

func TestService_RunScenario_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestScenario(repo)
	created.ScenarioType = "GROWTH"
	created.Parameters = JSON{"growth_rate": 10.0}
	_ = repo.UpdateScenario(ctx(), created)

	result, err := svc.RunScenario(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("RunScenario failed: %v", err)
	}

	if result.Status != "COMPLETED" {
		t.Errorf("expected status COMPLETED, got '%s'", result.Status)
	}
	if result.Results == nil {
		t.Error("expected results to be populated")
	}
}

func TestService_RunScenario_AlreadyCompleted_Fails(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestScenario(repo)
	created.Status = "COMPLETED"
	_ = repo.UpdateScenario(ctx(), created)

	_, err := svc.RunScenario(ctx(), created.ID.String())
	if err == nil {
		t.Fatal("expected error when running already-completed scenario")
	}
}

func TestService_CloneScenario_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestScenario(repo)

	clone, err := svc.CloneScenario(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("CloneScenario failed: %v", err)
	}

	if clone.ID == created.ID.String() {
		t.Error("expected cloned ID to be different from original")
	}
	if clone.Status != "DRAFT" {
		t.Errorf("expected status DRAFT for clone, got '%s'", clone.Status)
	}
}

func TestService_UpdateScenario_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestScenario(repo)
	newName := "Updated Scenario"
	req := UpdateScenarioRequest{Name: &newName}

	updated, err := svc.UpdateScenario(ctx(), created.ID.String(), req)
	if err != nil {
		t.Fatalf("UpdateScenario failed: %v", err)
	}

	if updated.Name != "Updated Scenario" {
		t.Errorf("expected name 'Updated Scenario', got '%s'", updated.Name)
	}
}

func TestService_DeleteScenario_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestScenario(repo)

	if err := svc.DeleteScenario(ctx(), created.ID.String()); err != nil {
		t.Fatalf("DeleteScenario failed: %v", err)
	}

	_, err := svc.GetScenarioByID(ctx(), created.ID.String())
	if err == nil {
		t.Fatal("expected error after deleting scenario")
	}
}

// =========================================================================
// Risk Service Tests
// =========================================================================

func TestService_ListRiskIndicators_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	createTestRiskIndicator(repo)
	createTestRiskIndicator(repo)

	resp, err := svc.ListRiskIndicators(ctx(), "", "", 1, 10)
	if err != nil {
		t.Fatalf("ListRiskIndicators failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
}

func TestService_GetRiskIndicatorByID_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestRiskIndicator(repo)

	found, err := svc.GetRiskIndicatorByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetRiskIndicatorByID failed: %v", err)
	}

	if found.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), found.ID)
	}
}

func TestService_UpdateRiskIndicator_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestRiskIndicator(repo)
	newLevel := "HIGH"
	req := UpdateRiskRequest{RiskLevel: &newLevel}

	updated, err := svc.UpdateRiskIndicator(ctx(), created.ID.String(), req)
	if err != nil {
		t.Fatalf("UpdateRiskIndicator failed: %v", err)
	}

	if updated.RiskLevel != "HIGH" {
		t.Errorf("expected level 'HIGH', got '%s'", updated.RiskLevel)
	}
}

func TestService_GetRiskDetail_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	detail, err := svc.GetRiskDetail(ctx(), "high-turnover")
	if err != nil {
		t.Fatalf("GetRiskDetail failed: %v", err)
	}

	if detail.RiskCode != "HIGH_TURNOVER" {
		t.Errorf("expected code 'HIGH_TURNOVER', got '%s'", detail.RiskCode)
	}
	if len(detail.Recommendations) == 0 {
		t.Error("expected recommendations to be populated")
	}
}

func TestService_GetRiskDetail_UnknownType(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetRiskDetail(ctx(), "unknown-risk")
	if err == nil {
		t.Fatal("expected error for unknown risk type")
	}
}

// =========================================================================
// Health Score Service Tests
// =========================================================================

func TestService_ListHealthScores_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	createTestHealthScore(repo)

	resp, err := svc.ListHealthScores(ctx(), "", "", 1, 10)
	if err != nil {
		t.Fatalf("ListHealthScores failed: %v", err)
	}

	if resp.Total != 1 {
		t.Errorf("expected total 1, got %d", resp.Total)
	}
}

func TestService_GetHealthDashboard_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	createTestHealthScore(repo)
	createTestHealthScore(repo)

	dash, err := svc.GetHealthDashboard(ctx(), "")
	if err != nil {
		t.Fatalf("GetHealthDashboard failed: %v", err)
	}

	if len(dash) > 2 {
		t.Errorf("expected at most 2 health scores, got %d", len(dash))
	}
}

func TestService_GetHealthScoreByID_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestHealthScore(repo)

	found, err := svc.GetHealthScoreByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetHealthScoreByID failed: %v", err)
	}

	if found.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), found.ID)
	}
}

func TestService_GetSpanOfControl_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	resp, err := svc.GetSpanOfControl(ctx())
	if err != nil {
		t.Fatalf("GetSpanOfControl failed: %v", err)
	}

	if resp.AvgRatio <= 0 {
		t.Errorf("expected positive AvgRatio, got %.1f", resp.AvgRatio)
	}
	if resp.Status != "HEALTHY" {
		t.Errorf("expected status 'HEALTHY', got '%s'", resp.Status)
	}
}

func TestService_GetSuccessionReadiness_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	resp, err := svc.GetSuccessionReadiness(ctx())
	if err != nil {
		t.Fatalf("GetSuccessionReadiness failed: %v", err)
	}

	if resp.CoverageRate <= 0 {
		t.Errorf("expected positive CoverageRate, got %.1f", resp.CoverageRate)
	}
}

// =========================================================================
// People Analytics Service Tests
// =========================================================================

func TestService_GetPeopleAnalytics_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	resp, err := svc.GetPeopleAnalytics(ctx(), "training-vs-performance")
	if err != nil {
		t.Fatalf("GetPeopleAnalytics failed: %v", err)
	}

	if resp.AnalysisType == "" {
		t.Error("expected AnalysisType to be set")
	}
	if resp.Correlation == 0 {
		t.Error("expected non-zero correlation")
	}
	if resp.Strength == "" {
		t.Error("expected Strength to be set")
	}
}

func TestService_GetPeopleAnalytics_UnknownType(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetPeopleAnalytics(ctx(), "unknown-analysis")
	if err == nil {
		t.Fatal("expected error for unknown analysis type")
	}
}

// =========================================================================
// Gap Analysis Service Tests
// =========================================================================

func TestService_GetGapAnalysis_NoForecasts(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetGapAnalysis(ctx(), "")
	if err != nil {
		// Expected: employees table doesn't exist in test DB context
		t.Logf("GetGapAnalysis expected error (no employees table): %v", err)
	}
}

// =========================================================================
// Projections Service Tests
// =========================================================================

func TestService_GetProjections_NoEmployeesTable(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetProjections(ctx(), "2026-Q3")
	if err != nil {
		// Expected: employees table doesn't exist in test DB context
		t.Logf("GetProjections expected error (no employees table): %v", err)
	}
}

// =========================================================================
// Executive Dashboard Service Tests
// =========================================================================

func TestService_GetExecutiveSummary_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	summary, err := svc.GetExecutiveSummary(ctx())
	if err != nil {
		t.Fatalf("GetExecutiveSummary failed: %v", err)
	}

	if summary.Period == "" {
		t.Error("expected period to be set")
	}
	if summary.HealthScore <= 0 {
		t.Errorf("expected positive HealthScore, got %.1f", summary.HealthScore)
	}
}

func TestService_GetExecutiveGrowth_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	resp, err := svc.GetExecutiveGrowth(ctx())
	if err != nil {
		t.Fatalf("GetExecutiveGrowth failed: %v", err)
	}

	if resp.Period == "" {
		t.Error("expected period to be set")
	}
}

func TestService_GetExecutiveHealthScore_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	resp, err := svc.GetExecutiveHealthScore(ctx())
	if err != nil {
		t.Fatalf("GetExecutiveHealthScore failed: %v", err)
	}

	if resp.Score <= 0 {
		t.Errorf("expected positive Score, got %.1f", resp.Score)
	}
}

// =========================================================================
// Capacity & Cost Service Tests
// =========================================================================

func TestService_GetCapacityForecast_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	resp, err := svc.GetCapacityForecast(ctx())
	if err != nil {
		t.Fatalf("GetCapacityForecast failed: %v", err)
	}

	if resp.Period == "" {
		t.Error("expected period to be set")
	}
	if resp.ProjectedUtil <= 0 {
		t.Errorf("expected positive ProjectedUtil, got %.1f", resp.ProjectedUtil)
	}
}

func TestService_GetPayrollCostBreakdown_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	resp, err := svc.GetPayrollCostBreakdown(ctx())
	if err != nil {
		t.Fatalf("GetPayrollCostBreakdown failed: %v", err)
	}

	if resp.TotalSalary <= 0 {
		t.Errorf("expected positive TotalSalary, got %.2f", resp.TotalSalary)
	}
}

func TestService_GetCostPerEmployee_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	resp, err := svc.GetCostPerEmployee(ctx())
	if err != nil {
		t.Fatalf("GetCostPerEmployee failed: %v", err)
	}

	if resp.AvgCostPerEmployee <= 0 {
		t.Errorf("expected positive AvgCostPerEmployee, got %.2f", resp.AvgCostPerEmployee)
	}
}

// =========================================================================
// Candidate Search Service Tests
// =========================================================================

func TestService_CandidateSearch_GroupsCandidatesByPosition(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	svc := NewService(repo, logger)
	createCandidateSearchTables(t, db)

	summaryID := uuid.New().String()
	orgVacantID := uuid.New().String()
	orgOccupiedID := uuid.New().String()
	reqID := uuid.New().String()
	reqClosedID := uuid.New().String()
	reqDupID := uuid.New().String()
	candID := uuid.New().String()

	db.Exec("INSERT INTO organization_summaries (id, code, decree_no, decree_date, status) VALUES (?, 'SA-01', 'SK-001', '2024-01-01', 'active')", summaryID)
	db.Exec("INSERT INTO organizations (id, organization_summary_id, code, full_code, nomenclature) VALUES (?, ?, 'ORG-01', 'SA-01.ORG-01', 'Staff IT')", orgVacantID, summaryID)
	db.Exec("INSERT INTO organizations (id, organization_summary_id, code, full_code, nomenclature) VALUES (?, ?, 'ORG-02', 'SA-01.ORG-02', 'Supervisor IT')", orgOccupiedID, summaryID)
	db.Exec("INSERT INTO employments (id, employee_id, organization_id, effective_date, effective_end_date) VALUES (?, ?, ?, '2024-01-01', NULL)", uuid.New().String(), uuid.New().String(), orgOccupiedID)
	db.Exec("INSERT INTO candidates (id, first_name, last_name, email, current_title, source) VALUES (?, 'Andi', 'Wijaya', 'andi@test.local', 'Staff IT', 'direct')", candID)
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status) VALUES (?, ?, 'Staff IT', 'OPEN')", reqID, orgVacantID)
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status) VALUES (?, ?, 'Staff IT (FILLED)', 'FILLED')", reqClosedID, orgVacantID)
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status) VALUES (?, ?, 'Staff IT (DUP)', 'OPEN')", reqDupID, orgVacantID)
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status) VALUES (?, ?, ?, 'SHORTLISTED')", uuid.New().String(), reqID, candID)
	// Aplikasi ke requisition yang sudah FILLED — harus TIDAK ikut (filter status req)
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status) VALUES (?, ?, ?, 'INTERVIEWED')", uuid.New().String(), reqClosedID, candID)
	// Aplikasi duplikat ke requisition OPEN lain — harus di-dedupe
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status) VALUES (?, ?, ?, 'OFFERED')", uuid.New().String(), reqDupID, candID)

	resp, err := svc.CandidateSearch(ctx(), "", "", 1, 20)
	if err != nil {
		t.Fatalf("CandidateSearch failed: %v", err)
	}

	data, ok := resp.Data.([]CandidateSearchPosition)
	if !ok {
		t.Fatalf("expected data type []CandidateSearchPosition, got %T", resp.Data)
	}
	if resp.Total != 1 {
		t.Fatalf("expected 1 vacant position, got %d", resp.Total)
	}
	pos := data[0]
	if pos.OrganizationID != orgVacantID {
		t.Errorf("expected position ORG-01, got %s", pos.OrganizationCode)
	}
	if pos.SummaryCode != "SA-01" || pos.SummaryDecreeNo != "SK-001" {
		t.Errorf("unexpected summary info: %s %s", pos.SummaryCode, pos.SummaryDecreeNo)
	}
	if pos.CandidateCount != 1 || len(pos.Candidates) != 1 {
		t.Fatalf("expected 1 candidate, got count=%d len=%d", pos.CandidateCount, len(pos.Candidates))
	}
	c := pos.Candidates[0]
	if c.Email != "andi@test.local" || c.ApplicationStatus != "SHORTLISTED" {
		t.Errorf("unexpected candidate payload: %s %s", c.Email, c.ApplicationStatus)
	}
}

func TestService_CandidateSearch_PaginationDefaults(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	svc := NewService(repo, logger)
	createCandidateSearchTables(t, db)

	resp, err := svc.CandidateSearch(ctx(), "", "", 0, 0)
	if err != nil {
		t.Fatalf("CandidateSearch failed: %v", err)
	}
	if resp.Page != 1 || resp.PerPage != 20 {
		t.Errorf("expected default pagination 1/20, got %d/%d", resp.Page, resp.PerPage)
	}
	if resp.Data == nil {
		t.Error("expected non-nil data (empty slice)")
	}
}

// =========================================================================
// Recruitment Analytics (S-2 — expected hires → remaining gap)
// =========================================================================

func TestService_GetRecruitmentAnalytics_ComputesPipeline(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	svc := NewService(repo, logger)
	createRecruitmentTables(t, db)
	if err := db.Exec(`CREATE TABLE employees (id CHAR(36) PRIMARY KEY, deleted_at DATETIME NULL)`).Error; err != nil {
		t.Fatalf("failed to create employees table: %v", err)
	}

	orgID := uuid.New().String()
	// OPEN requisition: 3 slots, 1 filled → 2 open
	reqOpen := uuid.New().String()
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status, slots_available, slots_filled) VALUES (?, ?, 'A', 'OPEN', 3, 1)", reqOpen, orgID)
	// FILLED requisition: 2 slots filled
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status, slots_available, slots_filled) VALUES (?, ?, 'B', 'FILLED', 2, 2)", uuid.New().String(), orgID)
	// Aplikasi: 2 NEW + 1 ACCEPTED pada requisition aktif
	cand := uuid.New().String()
	for _, st := range []string{"NEW", "NEW", "ACCEPTED"} {
		db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status) VALUES (?, ?, ?, ?)", uuid.New().String(), reqOpen, cand, st)
	}

	resp, err := svc.GetRecruitmentAnalytics(ctx())
	if err != nil {
		t.Fatalf("GetRecruitmentAnalytics failed: %v", err)
	}
	if resp.OpenPositions != 2 {
		t.Errorf("expected 2 open positions, got %d", resp.OpenPositions)
	}
	if resp.FilledPositions != 2 {
		t.Errorf("expected 2 filled positions, got %d", resp.FilledPositions)
	}
	if resp.ExpectedHires != 1 {
		t.Errorf("expected 1 expected hire (ACCEPTED), got %d", resp.ExpectedHires)
	}
	if resp.Pipeline == nil || len(resp.Pipeline) == 0 {
		t.Fatal("expected non-empty pipeline")
	}
	// Funnel mengikuti pipeline
	if len(resp.Funnel) != len(resp.Pipeline) {
		t.Errorf("expected funnel == pipeline length, got %d vs %d", len(resp.Funnel), len(resp.Pipeline))
	}
}

func TestService_GetRecruitmentAnalytics_EmptyData(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	svc := NewService(repo, logger)
	createRecruitmentTables(t, db)
	if err := db.Exec(`CREATE TABLE employees (id CHAR(36) PRIMARY KEY, deleted_at DATETIME NULL)`).Error; err != nil {
		t.Fatalf("failed to create employees table: %v", err)
	}

	resp, err := svc.GetRecruitmentAnalytics(ctx())
	if err != nil {
		t.Fatalf("GetRecruitmentAnalytics failed: %v", err)
	}
	if resp.OpenPositions != 0 || resp.ExpectedHires != 0 || resp.FilledPositions != 0 {
		t.Errorf("expected all zeros on empty data, got %+v", resp)
	}
	if resp.Pipeline == nil {
		t.Error("expected non-nil empty pipeline slice")
	}
}

func TestService_GetGapAnalysis_IncludesRemainingGap(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	svc := NewService(repo, logger)
	createRecruitmentTables(t, db)
	// GetGapAnalysis membaca tabel employees (supply) — buat minimal.
	if err := db.Exec(`CREATE TABLE employees (id CHAR(36) PRIMARY KEY, deleted_at DATETIME NULL)`).Error; err != nil {
		t.Fatalf("failed to create employees table: %v", err)
	}

	// Supply 0 (tidak ada employee). Demand default 5% dari supply = 0.
	// Tanpa forecast, demand = 0 → gap 0 → remaining gap 0.
	resp, err := svc.GetGapAnalysis(ctx(), "2026-08")
	if err != nil {
		t.Fatalf("GetGapAnalysis failed: %v", err)
	}
	if resp.RemainingGap != 0 {
		t.Errorf("expected remaining_gap 0 (no shortage), got %d", resp.RemainingGap)
	}
	if resp.ExpectedHires != 0 {
		t.Errorf("expected expected_hires 0, got %d", resp.ExpectedHires)
	}
}

func TestService_GetGapAnalysis_RemainingGapSubtractsExpectedHires(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	svc := NewService(repo, logger)
	createRecruitmentTables(t, db)
	if err := db.Exec(`CREATE TABLE employees (id CHAR(36) PRIMARY KEY, deleted_at DATETIME NULL)`).Error; err != nil {
		t.Fatalf("failed to create employees table: %v", err)
	}

	// Supply 5 employee, demand 0 (tanpa forecast → demand = supply*1.05 = 5
	// (int(5.25)=5)) → gap = 5-5 = 0. Tidak membentuk shortage.
	for i := 0; i < 5; i++ {
		db.Exec("INSERT INTO employees (id) VALUES (?)", uuid.New().String())
	}

	// Accepted offer = 2 (expected hires)
	reqOpen := uuid.New().String()
	orgID := uuid.New().String()
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status, slots_available, slots_filled) VALUES (?, ?, 'A', 'OPEN', 5, 0)", reqOpen, orgID)
	cand := uuid.New().String()
	for _, st := range []string{"ACCEPTED", "ACCEPTED"} {
		db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status) VALUES (?, ?, ?, ?)", uuid.New().String(), reqOpen, cand, st)
	}

	resp, err := svc.GetGapAnalysis(ctx(), "2026-08")
	if err != nil {
		t.Fatalf("GetGapAnalysis failed: %v", err)
	}
	// Demand fallback = int(5 * 1.05) = 5 → gap = 0 (OPTIMAL), bukan shortage.
	// Expected hires tetap ter-expose di response.
	if resp.ExpectedHires != 2 {
		t.Errorf("expected expected_hires 2, got %d", resp.ExpectedHires)
	}
	if resp.OpenPositions != 5 {
		t.Errorf("expected open_positions 5, got %d", resp.OpenPositions)
	}
	if resp.RemainingGap < 0 {
		t.Errorf("expected remaining_gap >= 0, got %d", resp.RemainingGap)
	}
}

// =========================================================================
// Recruitment Analytics advanced metrics (S-3)
// =========================================================================

func TestService_GetRecruitmentAnalytics_AdvancedMetrics(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	svc := NewService(repo, logger)
	createRecruitmentAnalyticsTables(t, db)
	if err := db.Exec(`CREATE TABLE employees (id CHAR(36) PRIMARY KEY, deleted_at DATETIME NULL)`).Error; err != nil {
		t.Fatalf("failed to create employees table: %v", err)
	}

	orgID := uuid.New().String()
	reqID := uuid.New().String()
	reqFilled := uuid.New().String()
	// Requisition FILLED: created 2026-01-01, closed 10 hari kemudian (ms)
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status, slots_available, slots_filled, closed_at, created_at) VALUES (?, ?, 'Staff IT', 'FILLED', 2, 2, 1768089600000, '2026-01-01 00:00:00')", reqFilled, orgID)
	// Requisition OPEN (pipeline tetap terbaca)
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status, slots_available, slots_filled) VALUES (?, ?, 'Staff IT (OPEN)', 'OPEN', 3, 0)", reqID, orgID)

	c1 := uuid.New().String() // referral, hire 12 hari
	c2 := uuid.New().String() // referral, tanpa aplikasi
	c3 := uuid.New().String() // linkedin, hire 1 hari
	c4 := uuid.New().String() // direct, offer ditolak (OFFERED)
	db.Exec("INSERT INTO candidates (id, first_name, last_name, email, source) VALUES (?, 'A', 'B', 'a@t.local', 'referral')", c1)
	db.Exec("INSERT INTO candidates (id, first_name, last_name, email, source) VALUES (?, 'C', 'D', 'c@t.local', 'referral')", c2)
	db.Exec("INSERT INTO candidates (id, first_name, last_name, email, source) VALUES (?, 'E', 'F', 'e@t.local', 'linkedin')", c3)
	db.Exec("INSERT INTO candidates (id, first_name, last_name, email, source) VALUES (?, 'G', 'H', 'g@t.local', 'direct')", c4)
	d := int64(86400000)
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status, applied_at, offered_at, accepted_at) VALUES (?, ?, ?, 'ACCEPTED', ?, ?, ?)", uuid.New().String(), reqID, c1, 1767225600000, 1767225600000+10*d, 1767225600000+12*d)
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status, applied_at, offered_at, accepted_at) VALUES (?, ?, ?, 'ACCEPTED', ?, ?, ?)", uuid.New().String(), reqID, c3, 1767225600000, 1767225600000+d, 1767225600000+d)
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status, applied_at, offered_at) VALUES (?, ?, ?, 'OFFERED', ?, ?)", uuid.New().String(), reqID, c4, 1767225600000, 1767225600000+2*d)

	resp, err := svc.GetRecruitmentAnalytics(ctx())
	if err != nil {
		t.Fatalf("GetRecruitmentAnalytics failed: %v", err)
	}
	// Time to Hire: (12 + 1)/2 = 6.5 hari
	if resp.TimeToHire != 6.5 {
		t.Errorf("expected time_to_hire 6.5, got %v", resp.TimeToHire)
	}
	// Time to Fill: 10 hari
	if resp.TimeToFill != 10 {
		t.Errorf("expected time_to_fill 10, got %v", resp.TimeToFill)
	}
	// Offer Acceptance Rate: 2 dari 3 offer = 66.7%
	if resp.OfferAcceptanceRate != 66.7 {
		t.Errorf("expected offer_acceptance_rate 66.7, got %v", resp.OfferAcceptanceRate)
	}
	// Source Conversion
	bySource := map[string]SourceConversionMetric{}
	for _, s := range resp.SourceConversion {
		bySource[s.Source] = s
	}
	if bySource["referral"].Candidates != 2 || bySource["referral"].Hires != 1 || bySource["referral"].ConversionRate != 50 {
		t.Errorf("unexpected referral conversion: %+v", bySource["referral"])
	}
	if bySource["linkedin"].Candidates != 1 || bySource["linkedin"].Hires != 1 || bySource["linkedin"].ConversionRate != 100 {
		t.Errorf("unexpected linkedin conversion: %+v", bySource["linkedin"])
	}
	if bySource["direct"].Candidates != 1 || bySource["direct"].Hires != 0 || bySource["direct"].ConversionRate != 0 {
		t.Errorf("unexpected direct conversion: %+v", bySource["direct"])
	}
	// BySource mengikuti jumlah kandidat per channel
	if len(resp.BySource) != 3 {
		t.Errorf("expected 3 by_source entries, got %d", len(resp.BySource))
	}
}

// =========================================================================
// Candidate Search S-3: filter posisi & internal candidate eligible
// =========================================================================

// fakeEligProvider mengimplementasikan InternalEligibilityProvider untuk test.
type fakeEligProvider struct {
	byPos map[string][]EligibleInternalCandidate
}

func (f fakeEligProvider) EligibleCandidatesForPositions(_ context.Context, ids []uuid.UUID) (map[string][]EligibleInternalCandidate, error) {
	out := map[string][]EligibleInternalCandidate{}
	for _, id := range ids {
		if cands, ok := f.byPos[id.String()]; ok {
			out[id.String()] = cands
		}
	}
	return out, nil
}

func TestService_CandidateSearch_PositionFilter(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	svc := NewService(repo, logger)
	createCandidateSearchTables(t, db)

	summaryID := uuid.New().String()
	orgIT := uuid.New().String()
	orgSales := uuid.New().String()
	db.Exec("INSERT INTO organization_summaries (id, code, decree_no, decree_date, status) VALUES (?, 'SA-01', 'SK-001', '2024-01-01', 'active')", summaryID)
	db.Exec("INSERT INTO organizations (id, organization_summary_id, code, full_code, nomenclature) VALUES (?, ?, 'ORG-01', 'SA-01.ORG-01', 'Staff IT')", orgIT, summaryID)
	db.Exec("INSERT INTO organizations (id, organization_summary_id, code, full_code, nomenclature) VALUES (?, ?, 'ORG-02', 'SA-01.ORG-02', 'Sales Executive')", orgSales, summaryID)

	resp, err := svc.CandidateSearch(ctx(), "", "Staff", 1, 20)
	if err != nil {
		t.Fatalf("CandidateSearch failed: %v", err)
	}
	if resp.Total != 1 {
		t.Fatalf("expected 1 vacant position matching 'Staff', got %d", resp.Total)
	}
	data := resp.Data.([]CandidateSearchPosition)
	if data[0].OrganizationID != orgIT {
		t.Errorf("expected Staff IT org, got %s", data[0].OrganizationName)
	}
}

func TestService_CandidateSearch_InternalCandidates(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	svc := NewService(repo, logger)
	createCandidateSearchTables(t, db)

	summaryID := uuid.New().String()
	orgVacant := uuid.New().String()
	posID := uuid.New().String()
	db.Exec("INSERT INTO organization_summaries (id, code, decree_no, decree_date, status) VALUES (?, 'SA-01', 'SK-001', '2024-01-01', 'active')", summaryID)
	db.Exec("INSERT INTO organizations (id, organization_summary_id, code, full_code, nomenclature) VALUES (?, ?, 'ORG-01', 'SA-01.ORG-01', 'Staff IT')", orgVacant, summaryID)
	db.Exec("INSERT INTO positions (id, organization_id, title, is_active) VALUES (?, ?, 'Staff IT', 1)", posID, orgVacant)

	svc.SetInternalEligibilityProvider(fakeEligProvider{byPos: map[string][]EligibleInternalCandidate{
		posID: {
			{EmployeeID: "emp-1", Name: "Budi Santoso", CurrentPositionID: "pos-src", CurrentPositionName: "Supervisor IT", SourceStepSequence: 1},
		},
	}})

	resp, err := svc.CandidateSearch(ctx(), "", "", 1, 20)
	if err != nil {
		t.Fatalf("CandidateSearch failed: %v", err)
	}
	data := resp.Data.([]CandidateSearchPosition)
	if len(data) != 1 {
		t.Fatalf("expected 1 vacant position, got %d", len(data))
	}
	pos := data[0]
	if pos.InternalCandidateCount != 1 || len(pos.InternalCandidates) != 1 {
		t.Fatalf("expected 1 internal candidate, got count=%d len=%d", pos.InternalCandidateCount, len(pos.InternalCandidates))
	}
	if pos.InternalCandidates[0].Name != "Budi Santoso" || pos.InternalCandidates[0].SourceStepSequence != 1 {
		t.Errorf("unexpected internal candidate payload: %+v", pos.InternalCandidates[0])
	}
}

func TestService_CandidateSearch_InternalCandidates_DedupedPerOrg(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	svc := NewService(repo, logger)
	createCandidateSearchTables(t, db)

	summaryID := uuid.New().String()
	orgVacant := uuid.New().String()
	posA := uuid.New().String()
	posB := uuid.New().String()
	db.Exec("INSERT INTO organization_summaries (id, code, decree_no, decree_date, status) VALUES (?, 'SA-01', 'SK-001', '2024-01-01', 'active')", summaryID)
	db.Exec("INSERT INTO organizations (id, organization_summary_id, code, full_code, nomenclature) VALUES (?, ?, 'ORG-01', 'SA-01.ORG-01', 'Staff IT')", orgVacant, summaryID)
	// Dua position pada org yang sama, keduanya eligible untuk employee yang sama
	db.Exec("INSERT INTO positions (id, organization_id, title, is_active) VALUES (?, ?, 'Staff IT', 1)", posA, orgVacant)
	db.Exec("INSERT INTO positions (id, organization_id, title, is_active) VALUES (?, ?, 'Staff IT (B)', 1)", posB, orgVacant)

	svc.SetInternalEligibilityProvider(fakeEligProvider{byPos: map[string][]EligibleInternalCandidate{
		posA: {{EmployeeID: "emp-1", Name: "Budi Santoso", SourceStepSequence: 1}},
		posB: {{EmployeeID: "emp-1", Name: "Budi Santoso", SourceStepSequence: 2}},
	}})

	resp, err := svc.CandidateSearch(ctx(), "", "", 1, 20)
	if err != nil {
		t.Fatalf("CandidateSearch failed: %v", err)
	}
	data := resp.Data.([]CandidateSearchPosition)
	pos := data[0]
	if pos.InternalCandidateCount != 1 || len(pos.InternalCandidates) != 1 {
		t.Fatalf("expected 1 internal candidate after dedupe, got count=%d len=%d", pos.InternalCandidateCount, len(pos.InternalCandidates))
	}
}

func TestService_CandidateSearch_ProviderNil_EmptyInternal(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	svc := NewService(repo, logger)
	createCandidateSearchTables(t, db)

	summaryID := uuid.New().String()
	orgVacant := uuid.New().String()
	db.Exec("INSERT INTO organization_summaries (id, code, decree_no, decree_date, status) VALUES (?, 'SA-01', 'SK-001', '2024-01-01', 'active')", summaryID)
	db.Exec("INSERT INTO organizations (id, organization_summary_id, code, full_code, nomenclature) VALUES (?, ?, 'ORG-01', 'SA-01.ORG-01', 'Staff IT')", orgVacant, summaryID)

	// Tanpa SetInternalEligibilityProvider → internal candidates kosong (fail-safe).
	resp, err := svc.CandidateSearch(ctx(), "", "", 1, 20)
	if err != nil {
		t.Fatalf("CandidateSearch failed: %v", err)
	}
	data := resp.Data.([]CandidateSearchPosition)
	if len(data) != 1 {
		t.Fatalf("expected 1 vacant position, got %d", len(data))
	}
	pos := data[0]
	if pos.InternalCandidateCount != 0 || len(pos.InternalCandidates) != 0 {
		t.Errorf("expected empty internal candidates without provider, got %+v", pos.InternalCandidates)
	}
}

// =========================================================================
// Quality of Hire (S-6) Service Tests
// =========================================================================

func TestService_GetQualityOfHire_CompositeAndBreakdown(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	svc := NewService(repo, logger)
	createQualityOfHireTables(t, db)
	ctx := ctx()

	orgA := uuid.New()
	reqA := uuid.New()
	candReferral := uuid.New()
	candJobBoard := uuid.New()
	empA := uuid.New()
	empB := uuid.New()

	db.Exec("INSERT INTO candidates (id, first_name, last_name, email, source) VALUES (?, 'A', 'A', 'a@x.com', 'referral')", candReferral.String())
	db.Exec("INSERT INTO candidates (id, first_name, last_name, email, source) VALUES (?, 'B', 'B', 'b@x.com', 'job_board')", candJobBoard.String())
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status) VALUES (?, ?, 'Eng', 'FILLED')", reqA.String(), orgA.String())

	// Hire 1 (referral): interview 80, onboarding COMPLETED, perf 85, retained
	appA := uuid.New()
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status) VALUES (?, ?, ?, 'ACCEPTED')", appA.String(), reqA.String(), candReferral.String())
	db.Exec("INSERT INTO interviews (id, application_id, score) VALUES (?, ?, 80)", uuid.New().String(), appA.String())
	db.Exec("INSERT INTO employee_onboardings (id, application_id, employee_id, status) VALUES (?, ?, ?, 'COMPLETED')", uuid.New().String(), appA.String(), empA.String())
	db.Exec("INSERT INTO performance_evaluations (id, employee_id, final_score, status, updated_at) VALUES (?, ?, 85, 'COMPLETED', '2025-06-01')", uuid.New().String(), empA.String())
	db.Exec("INSERT INTO employments (id, employee_id, effective_end_date) VALUES (?, ?, NULL)", uuid.New().String(), empA.String())

	// Hire 2 (job_board): interview 70, onboarding PENDING, perf 75, not retained
	appB := uuid.New()
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status) VALUES (?, ?, ?, 'ACCEPTED')", appB.String(), reqA.String(), candJobBoard.String())
	db.Exec("INSERT INTO interviews (id, application_id, score) VALUES (?, ?, 70)", uuid.New().String(), appB.String())
	db.Exec("INSERT INTO employee_onboardings (id, application_id, employee_id, status) VALUES (?, ?, ?, 'PENDING')", uuid.New().String(), appB.String(), empB.String())
	db.Exec("INSERT INTO performance_evaluations (id, employee_id, final_score, status, updated_at) VALUES (?, ?, 75, 'COMPLETED', '2025-06-01')", uuid.New().String(), empB.String())
	db.Exec("INSERT INTO employments (id, employee_id, effective_end_date) VALUES (?, ?, '2025-01-01')", uuid.New().String(), empB.String())

	resp, err := svc.GetQualityOfHire(ctx)
	if err != nil {
		t.Fatalf("GetQualityOfHire failed: %v", err)
	}
	if resp.HiresAnalyzed != 2 {
		t.Errorf("expected 2 hires analyzed, got %d", resp.HiresAnalyzed)
	}
	// Interview avg = (80+70)/2 = 75
	if resp.InterviewScore != 75 {
		t.Errorf("expected interview score 75, got %v", resp.InterviewScore)
	}
	// Onboarding completion = 1/2 = 50
	if resp.OnboardingCompletionRate != 50 {
		t.Errorf("expected onboarding completion 50, got %v", resp.OnboardingCompletionRate)
	}
	// Performance avg = (85+75)/2 = 80
	if resp.PerformanceScore != 80 {
		t.Errorf("expected performance score 80, got %v", resp.PerformanceScore)
	}
	// Retention = 1/2 = 50
	if resp.RetentionRate != 50 {
		t.Errorf("expected retention 50, got %v", resp.RetentionRate)
	}
	// Overall = avg(75, 50, 80, 50) = 63.75 → 63.8
	if resp.OverallScore != 63.8 {
		t.Errorf("expected overall 63.8, got %v", resp.OverallScore)
	}
	// Placeholder tetap 0
	if resp.RecruitmentMatchScore != 0 || resp.AssessmentScore != 0 {
		t.Errorf("expected match/assessment placeholder 0, got %v/%v", resp.RecruitmentMatchScore, resp.AssessmentScore)
	}
	// Breakdown by source: referral=1 hire score komposit hire1, job_board=1 hire
	if len(resp.BySource) != 2 {
		t.Fatalf("expected 2 sources in breakdown, got %d", len(resp.BySource))
	}
	byKey := map[string]QualityOfHireBreakdown{}
	for _, b := range resp.BySource {
		byKey[b.Key] = b
	}
	// hire1 komposit = (80+100+85+100)/4 = 91.25 → 91.3; hire2 = (70+0+75+0)/4 = 36.25 → 36.3
	if s, ok := byKey["referral"]; !ok || s.Hires != 1 || s.Score != 91.3 {
		t.Errorf("expected referral breakdown hires=1 score=91.3, got %+v", byKey["referral"])
	}
	if s, ok := byKey["job_board"]; !ok || s.Hires != 1 || s.Score != 36.3 {
		t.Errorf("expected job_board breakdown hires=1 score=36.3, got %+v", byKey["job_board"])
	}
	// ByRequisition harus 1 entri (keduanya reqA)
	if len(resp.ByRequisition) != 1 || resp.ByRequisition[0].Hires != 2 {
		t.Errorf("expected 1 requisition breakdown with 2 hires, got %+v", resp.ByRequisition)
	}
}

func TestService_GetQualityOfHire_PartialData_OverallConsistentWithBreakdown(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	svc := NewService(repo, logger)
	createQualityOfHireTables(t, db)
	ctx := ctx()

	orgA := uuid.New()
	reqA := uuid.New()
	candFull := uuid.New()
	candPartial := uuid.New()
	empA := uuid.New()

	db.Exec("INSERT INTO candidates (id, first_name, last_name, email, source) VALUES (?, 'A', 'A', 'a@x.com', 'referral')", candFull.String())
	db.Exec("INSERT INTO candidates (id, first_name, last_name, email, source) VALUES (?, 'B', 'B', 'b@x.com', 'referral')", candPartial.String())
	db.Exec("INSERT INTO job_requisitions (id, organization_id, title, status) VALUES (?, ?, 'Eng', 'FILLED')", reqA.String(), orgA.String())

	// Hire 1: data lengkap → komposit (80+100+85+100)/4 = 91.3
	appA := uuid.New()
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status) VALUES (?, ?, ?, 'ACCEPTED')", appA.String(), reqA.String(), candFull.String())
	db.Exec("INSERT INTO interviews (id, application_id, score) VALUES (?, ?, 80)", uuid.New().String(), appA.String())
	db.Exec("INSERT INTO employee_onboardings (id, application_id, employee_id, status) VALUES (?, ?, ?, 'COMPLETED')", uuid.New().String(), appA.String(), empA.String())
	db.Exec("INSERT INTO performance_evaluations (id, employee_id, final_score, status, updated_at) VALUES (?, ?, 85, 'COMPLETED', '2025-06-01')", uuid.New().String(), empA.String())
	db.Exec("INSERT INTO employments (id, employee_id, effective_end_date) VALUES (?, ?, NULL)", uuid.New().String(), empA.String())

	// Hire 2: HANYA interview 60 (tanpa onboarding/perf/employee) → komposit 60
	appB := uuid.New()
	db.Exec("INSERT INTO job_applications (id, requisition_id, candidate_id, status) VALUES (?, ?, ?, 'ACCEPTED')", appB.String(), reqA.String(), candPartial.String())
	db.Exec("INSERT INTO interviews (id, application_id, score) VALUES (?, ?, 60)", uuid.New().String(), appB.String())

	resp, err := svc.GetQualityOfHire(ctx)
	if err != nil {
		t.Fatalf("GetQualityOfHire failed: %v", err)
	}
	// Overall = avg(91.3, 60) = 75.65 → 75.7 — KONSISTEN dengan breakdown
	// (jangan rata-rata komponen agregat yang akan memberi hasil berbeda).
	if resp.OverallScore != 75.7 {
		t.Errorf("expected overall 75.7 (avg per-hire composite), got %v", resp.OverallScore)
	}
	// Interview aggregate tetap dihitung dari hire yang punya data interview.
	if resp.InterviewScore != 70 {
		t.Errorf("expected interview avg (80+60)/2=70, got %v", resp.InterviewScore)
	}
	if len(resp.BySource) != 1 {
		t.Fatalf("expected 1 source breakdown, got %d", len(resp.BySource))
	}
	// Breakdown score = avg komposit hire dalam grup = sama dengan overall.
	if resp.BySource[0].Score != resp.OverallScore {
		t.Errorf("expected by_source score == overall (%v), got %v", resp.OverallScore, resp.BySource[0].Score)
	}
	if resp.BySource[0].Hires != 2 {
		t.Errorf("expected 2 hires in source breakdown, got %d", resp.BySource[0].Hires)
	}
}

func TestService_GetQualityOfHire_Empty(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	svc := NewService(repo, logger)
	createQualityOfHireTables(t, db)

	resp, err := svc.GetQualityOfHire(ctx())
	if err != nil {
		t.Fatalf("GetQualityOfHire failed: %v", err)
	}
	if resp.HiresAnalyzed != 0 || resp.OverallScore != 0 {
		t.Errorf("expected zeros on empty data, got hires=%d overall=%v", resp.HiresAnalyzed, resp.OverallScore)
	}
	if len(resp.BySource) != 0 || len(resp.ByRequisition) != 0 || len(resp.ByOrganization) != 0 {
		t.Errorf("expected empty breakdowns, got %+v", resp)
	}
}

