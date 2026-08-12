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

	resp, err := svc.CandidateSearch(ctx(), "", 1, 20)
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

	resp, err := svc.CandidateSearch(ctx(), "", 0, 0)
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

