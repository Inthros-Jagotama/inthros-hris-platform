package careerintelligence

import (
	"testing"

	"github.com/google/uuid"
)

// =========================================================================
// Talent Map Service Tests
// =========================================================================

func TestService_CreateTalentMap_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateTalentMapRequest{
		EmployeeID:  uuidStr(),
		Period:      "2026-07",
		Performance: "HIGH",
		Potential:   "HIGH",
		Notes:       "Star performer",
	}

	resp, err := svc.CreateTalentMap(ctx(), req)
	if err != nil {
		t.Fatalf("CreateTalentMap failed: %v", err)
	}

	if resp.Performance != "HIGH" {
		t.Errorf("expected performance 'HIGH', got '%s'", resp.Performance)
	}
	if resp.GridPosition != "9-BOX-1" {
		t.Errorf("expected grid position '9-BOX-1', got '%s'", resp.GridPosition)
	}
	if resp.ID == "" {
		t.Error("expected ID to be set")
	}
}

func TestService_CreateTalentMap_GridPositionComputation(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	tests := []struct {
		perf, pot, expected string
	}{
		{"HIGH", "HIGH", "9-BOX-1"},
		{"HIGH", "MEDIUM", "9-BOX-2"},
		{"HIGH", "LOW", "9-BOX-3"},
		{"MEDIUM", "HIGH", "9-BOX-4"},
		{"MEDIUM", "MEDIUM", "9-BOX-5"},
		{"MEDIUM", "LOW", "9-BOX-6"},
		{"LOW", "HIGH", "9-BOX-7"},
		{"LOW", "MEDIUM", "9-BOX-8"},
		{"LOW", "LOW", "9-BOX-9"},
	}

	for _, tt := range tests {
		req := CreateTalentMapRequest{
			EmployeeID:  uuidStr(),
			Period:      "2026-07",
			Performance: tt.perf,
			Potential:   tt.pot,
		}

		resp, err := svc.CreateTalentMap(ctx(), req)
		if err != nil {
			t.Fatalf("CreateTalentMap failed for %s/%s: %v", tt.perf, tt.pot, err)
		}

		if resp.GridPosition != tt.expected {
			t.Errorf("expected %s for %s/%s, got %s", tt.expected, tt.perf, tt.pot, resp.GridPosition)
		}
	}
}

func TestService_GetTalentMapByID_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestTalentMap(repo)

	found, err := svc.GetTalentMapByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetTalentMapByID failed: %v", err)
	}

	if found.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), found.ID)
	}
}

func TestService_GetTalentMapByID_NotFound(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetTalentMapByID(ctx(), uuidStr())
	if err == nil {
		t.Fatal("expected error for non-existent talent map")
	}
}

func TestService_ListTalentMaps_DefaultPagination(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	for i := 0; i < 3; i++ {
		createTestTalentMap(repo)
	}

	resp, err := svc.ListTalentMaps(ctx(), "", "", 0, 0)
	if err != nil {
		t.Fatalf("ListTalentMaps failed: %v", err)
	}

	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}
	if resp.Page != defaultPage {
		t.Errorf("expected default page %d, got %d", defaultPage, resp.Page)
	}
}

func TestService_UpdateTalentMap_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestTalentMap(repo)
	newPerf := "LOW"
	req := UpdateTalentMapRequest{Performance: &newPerf}

	updated, err := svc.UpdateTalentMap(ctx(), created.ID.String(), req)
	if err != nil {
		t.Fatalf("UpdateTalentMap failed: %v", err)
	}

	if updated.Performance != "LOW" {
		t.Errorf("expected performance 'LOW', got '%s'", updated.Performance)
	}
}

func TestService_DeleteTalentMap_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestTalentMap(repo)

	if err := svc.DeleteTalentMap(ctx(), created.ID.String()); err != nil {
		t.Fatalf("DeleteTalentMap failed: %v", err)
	}

	_, err := svc.GetTalentMapByID(ctx(), created.ID.String())
	if err == nil {
		t.Fatal("expected error after deleting talent map")
	}
}

func TestService_GetTalentGrid_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	for i := 0; i < 5; i++ {
		createTestTalentMap(repo)
	}

	grid, err := svc.GetTalentGrid(ctx(), "")
	if err != nil {
		t.Fatalf("GetTalentGrid failed: %v", err)
	}

	if len(grid.Quadrants) != 9 {
		t.Errorf("expected 9 grid quadrants, got %d", len(grid.Quadrants))
	}
	if grid.Total != 5 {
		t.Errorf("expected total 5, got %d", grid.Total)
	}
}

// =========================================================================
// Career Interest Service Tests
// =========================================================================

func TestService_CreateCareerInterest_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateCareerInterestRequest{
		EmployeeID:     uuidStr(),
		InterestType:   "LEADERSHIP",
		TargetPosition: "Senior Manager",
		ReadinessLevel: "1_YEAR",
	}

	resp, err := svc.CreateCareerInterest(ctx(), req)
	if err != nil {
		t.Fatalf("CreateCareerInterest failed: %v", err)
	}

	if resp.InterestType != "LEADERSHIP" {
		t.Errorf("expected type 'LEADERSHIP', got '%s'", resp.InterestType)
	}
	if resp.ID == "" {
		t.Error("expected ID to be set")
	}
}

func TestService_ListCareerInterests_DefaultPagination(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	for i := 0; i < 3; i++ {
		createTestCareerInterest(repo)
	}

	resp, err := svc.ListCareerInterests(ctx(), "", 0, 0)
	if err != nil {
		t.Fatalf("ListCareerInterests failed: %v", err)
	}

	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}
}

func TestService_ListCareerInterests_ByEmployee(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	ci := createTestCareerInterest(repo)

	resp, err := svc.ListCareerInterests(ctx(), ci.EmployeeID.String(), 1, 10)
	if err != nil {
		t.Fatalf("ListCareerInterests by employee failed: %v", err)
	}

	if resp.Total != 1 {
		t.Errorf("expected 1 interest for employee, got %d", resp.Total)
	}
}

// =========================================================================
// Career Path Service Tests
// =========================================================================

func TestService_CreateCareerPath_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateCareerPathRequest{
		SourceTitleID: uuidStr(),
		TargetTitleID: uuidStr(),
		PathType:      "PROMOTION",
		TypicalTenure: 24,
	}

	resp, err := svc.CreateCareerPath(ctx(), req)
	if err != nil {
		t.Fatalf("CreateCareerPath failed: %v", err)
	}

	if resp.PathType != "PROMOTION" {
		t.Errorf("expected path type 'PROMOTION', got '%s'", resp.PathType)
	}
	if resp.ID == "" {
		t.Error("expected ID to be set")
	}
}

func TestService_ListCareerPaths_DefaultPagination(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	for i := 0; i < 3; i++ {
		createTestCareerPath(repo)
	}

	resp, err := svc.ListCareerPaths(ctx(), 0, 0)
	if err != nil {
		t.Fatalf("ListCareerPaths failed: %v", err)
	}

	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}
}

func TestService_DeleteCareerPath_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestCareerPath(repo)

	if err := svc.DeleteCareerPath(ctx(), created.ID.String()); err != nil {
		t.Fatalf("DeleteCareerPath failed: %v", err)
	}
}

func TestService_GetGapAnalysis_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := GapAnalysisRequest{
		EmployeeID:    uuidStr(),
		TargetTitleID: uuidStr(),
	}

	resp, err := svc.GetGapAnalysis(ctx(), req)
	if err != nil {
		t.Fatalf("GetGapAnalysis failed: %v", err)
	}

	if resp.EmployeeID != req.EmployeeID {
		t.Errorf("expected employee ID '%s', got '%s'", req.EmployeeID, resp.EmployeeID)
	}
	if resp.GapPercentage <= 0 {
		t.Errorf("expected positive gap percentage, got %.2f", resp.GapPercentage)
	}
	if len(resp.Recommendations) == 0 {
		t.Error("expected at least one recommendation")
	}
	if resp.EstimatedTimeline == "" {
		t.Error("expected estimated timeline to be set")
	}
}

// =========================================================================
// Succession Plan Service Tests
// =========================================================================

func TestService_CreateSuccessionPlan_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateSuccessionPlanRequest{
		PositionID:     uuidStr(),
		SuccessorID:    uuidStr(),
		ReadinessLevel: "READY_NOW",
	}

	resp, err := svc.CreateSuccessionPlan(ctx(), req)
	if err != nil {
		t.Fatalf("CreateSuccessionPlan failed: %v", err)
	}

	if resp.ReadinessLevel != "READY_NOW" {
		t.Errorf("expected readiness 'READY_NOW', got '%s'", resp.ReadinessLevel)
	}
	if resp.Status != "ACTIVE" {
		t.Errorf("expected status 'ACTIVE', got '%s'", resp.Status)
	}
	if resp.ID == "" {
		t.Error("expected ID to be set")
	}
}

func TestService_ListSuccessionPlans_DefaultPagination(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	for i := 0; i < 3; i++ {
		createTestSuccessionPlan(repo)
	}

	resp, err := svc.ListSuccessionPlans(ctx(), 0, 0)
	if err != nil {
		t.Fatalf("ListSuccessionPlans failed: %v", err)
	}

	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}
}

func TestService_GetSuccessionPlanByID_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestSuccessionPlan(repo)

	found, err := svc.GetSuccessionPlanByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetSuccessionPlanByID failed: %v", err)
	}

	if found.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), found.ID)
	}
}

func TestService_UpdateSuccessionPlan_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestSuccessionPlan(repo)
	newReadiness := "READY_2YR"
	req := UpdateSuccessionPlanRequest{ReadinessLevel: &newReadiness}

	updated, err := svc.UpdateSuccessionPlan(ctx(), created.ID.String(), req)
	if err != nil {
		t.Fatalf("UpdateSuccessionPlan failed: %v", err)
	}

	if updated.ReadinessLevel != "READY_2YR" {
		t.Errorf("expected readiness 'READY_2YR', got '%s'", updated.ReadinessLevel)
	}
}

func TestService_DeleteSuccessionPlan_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestSuccessionPlan(repo)

	if err := svc.DeleteSuccessionPlan(ctx(), created.ID.String()); err != nil {
		t.Fatalf("DeleteSuccessionPlan failed: %v", err)
	}

	_, err := svc.GetSuccessionPlanByID(ctx(), created.ID.String())
	if err == nil {
		t.Fatal("expected error after deleting succession plan")
	}
}

// uuidStr generates a random UUID string for test use.
func uuidStr() string {
	return uuid.New().String()
}
