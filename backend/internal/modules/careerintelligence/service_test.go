package careerintelligence

import (
	"context"
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

func TestService_CreateCareerPathLadder_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	posA, posB := seedCareerPathPositions(t, repo)
	minMonths := 12
	desc := "Jenjang career IT"

	req := CreateCareerPathLadderRequest{
		Name:        "Staff to Supervisor",
		Description: &desc,
		Steps: []CreateCareerPathStepRequest{
			{PositionID: posA.String(), Sequence: 1, MinimumServiceMonths: &minMonths},
			{PositionID: posB.String(), Sequence: 2},
		},
	}

	resp, err := svc.CreateCareerPathLadder(ctx(), req)
	if err != nil {
		t.Fatalf("CreateCareerPathLadder failed: %v", err)
	}

	if resp.Name != "Staff to Supervisor" {
		t.Errorf("expected name 'Staff to Supervisor', got '%s'", resp.Name)
	}
	if len(resp.Steps) != 2 {
		t.Errorf("expected 2 steps, got %d", len(resp.Steps))
	}
	if resp.Steps[0].PositionName == "" || resp.Steps[1].PositionName == "" {
		t.Error("expected steps to be enriched with position names")
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

	resp, err := svc.ListCareerPaths(ctx(), 0, 0, "")
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
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	orgID := uuid.New()
	seedGapAnalysisCompetencyData(t, repo, employeeID, orgID)

	req := GapAnalysisRequest{
		EmployeeID:    employeeID.String(),
		TargetTitleID: orgID.String(),
	}

	resp, err := svc.GetGapAnalysis(ctx(), req)
	if err != nil {
		t.Fatalf("GetGapAnalysis failed: %v", err)
	}

	if resp.EmployeeID != req.EmployeeID {
		t.Errorf("expected employee ID '%s', got '%s'", req.EmployeeID, resp.EmployeeID)
	}
	// Fixture: 2 syarat (Leadership lvl4, Communication lvl3); employee
	// memenuhi Leadership (lvl4) tapi tidak Communication (lvl2) -> 1/2 gap.
	if resp.TotalRequired != 2 {
		t.Errorf("expected 2 total required, got %d", resp.TotalRequired)
	}
	if resp.MatchedSkills != 1 {
		t.Errorf("expected 1 matched skill, got %d", resp.MatchedSkills)
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

func TestService_GetGapAnalysis_NoRequirements(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := GapAnalysisRequest{
		EmployeeID:    uuidStr(),
		TargetTitleID: uuidStr(),
	}

	// competency_score_details table doesn't even exist in this test's DB
	// (never seeded) -- GetOrgCompetencyRequirements must surface this as a
	// clear error, not a misleading 0%-gap result.
	if _, err := svc.GetGapAnalysis(ctx(), req); err == nil {
		t.Error("expected error when target org has no competency assessment data")
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

// fakeNotifier implements Notifier for tests, recording every call.
type fakeNotifier struct {
	calls []fakeNotifyCall
}

type fakeNotifyCall struct {
	RecipientUserID uuid.UUID
	Type            string
	Params          []string
}

func (f *fakeNotifier) Notify(_ context.Context, recipientUserID uuid.UUID, notifType string, params []string, _ string, _ uuid.UUID) error {
	f.calls = append(f.calls, fakeNotifyCall{RecipientUserID: recipientUserID, Type: notifType, Params: params})
	return nil
}

// TestService_CreateSuccessionPlan_NotifiesSuccessor verifies the module
// notification integration (module-career-intelligence-plan.md §9 #7):
// naming a successor sends a SUCCESSION_PLAN_NAMED notification to their
// linked user account.
func TestService_CreateSuccessionPlan_NotifiesSuccessor(t *testing.T) {
	svc, _, db, cleanup := newTestService()
	defer cleanup()
	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	if err := db.Exec("CREATE TABLE IF NOT EXISTS employee_accounts (employee_id CHAR(36), user_id CHAR(36))").Error; err != nil {
		t.Fatalf("failed to create employee_accounts table: %v", err)
	}
	successorID := uuid.New()
	userID := uuid.New()
	if err := db.Exec("INSERT INTO employee_accounts (employee_id, user_id) VALUES (?, ?)", successorID.String(), userID.String()).Error; err != nil {
		t.Fatalf("failed to seed employee account: %v", err)
	}

	req := CreateSuccessionPlanRequest{
		PositionID:     uuidStr(),
		SuccessorID:    successorID.String(),
		ReadinessLevel: "READY_NOW",
	}
	if _, err := svc.CreateSuccessionPlan(ctx(), req); err != nil {
		t.Fatalf("CreateSuccessionPlan failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(notifier.calls))
	}
	call := notifier.calls[0]
	if call.Type != "SUCCESSION_PLAN_NAMED" {
		t.Errorf("expected notification type SUCCESSION_PLAN_NAMED, got %s", call.Type)
	}
	if call.RecipientUserID != userID {
		t.Errorf("expected recipient %s, got %s", userID, call.RecipientUserID)
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

// =========================================================================
// Internal Candidate Eligibility (S-4)
// =========================================================================

func TestService_GetEligibleEmployeesForPath_Success(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()
	createEligibilityTables(t, db)
	ctx := ctx()

	// Path: Staff (1) → Supervisor (2) → Manager (3 = target)
	srcStaff := uuid.New()
	srcSupervisor := uuid.New()
	targetManager := uuid.New()
	cp := &CareerPath{Name: "PROMOTION: Staff → Manager", IsActive: true}
	steps := []CareerPathStep{
		{PositionID: srcStaff, Sequence: 1},
		{PositionID: srcSupervisor, Sequence: 2},
		{PositionID: targetManager, Sequence: 3, PathType: "PROMOTION"},
	}
	if err := repo.CreateCareerPathTx(ctx, cp, steps); err != nil {
		t.Fatalf("CreateCareerPathTx failed: %v", err)
	}

	// Employee A di posisi Staff (source), Employee B di posisi Supervisor (source)
	past := "2020-01-01"
	for _, emp := range []struct {
		id   uuid.UUID
		name string
		pos  uuid.UUID
	}{{uuid.New(), "Andi", srcStaff}, {uuid.New(), "Budi", srcSupervisor}} {
		db.Exec("INSERT INTO employees (id, name) VALUES (?, ?)", emp.id.String(), emp.name)
		db.Exec("INSERT INTO employments (id, employee_id, position_id, effective_date, effective_end_date) VALUES (?, ?, ?, ?, NULL)", uuid.New().String(), emp.id.String(), emp.pos.String(), past)
	}

	result, err := svc.GetEligibleEmployeesForPath(ctx, cp.ID.String())
	if err != nil {
		t.Fatalf("GetEligibleEmployeesForPath failed: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 eligible employees, got %d", len(result))
	}
	for _, e := range result {
		if e.TargetPositionID != targetManager.String() {
			t.Errorf("expected target %s, got %s", targetManager, e.TargetPositionID)
		}
		if e.PathID != cp.ID.String() {
			t.Errorf("expected path %s, got %s", cp.ID, e.PathID)
		}
	}
}

func TestService_GetEligibleEmployeesForPath_RequiresTwoSteps(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	cp := &CareerPath{Name: "SINGLE STEP", IsActive: true}
	steps := []CareerPathStep{
		{PositionID: uuid.New(), Sequence: 1},
	}
	if err := repo.CreateCareerPathTx(ctx(), cp, steps); err != nil {
		t.Fatalf("CreateCareerPathTx failed: %v", err)
	}

	_, err := svc.GetEligibleEmployeesForPath(ctx(), cp.ID.String())
	if err == nil {
		t.Fatal("expected error for path with < 2 steps")
	}
}

func TestService_GetEligibleEmployeesForPath_NotFound(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()
	_, err := svc.GetEligibleEmployeesForPath(ctx(), uuid.New().String())
	if err == nil {
		t.Fatal("expected error for non-existent path")
	}
}

func TestService_GetEligibleEmployeesByPosition_Success(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()
	createEligibilityTables(t, db)
	ctx := ctx()

	srcStaff := uuid.New()
	targetManager := uuid.New()
	cp := &CareerPath{Name: "PROMOTION: Staff → Manager", IsActive: true}
	steps := []CareerPathStep{
		{PositionID: srcStaff, Sequence: 1},
		{PositionID: targetManager, Sequence: 2, PathType: "PROMOTION"},
	}
	if err := repo.CreateCareerPathTx(ctx, cp, steps); err != nil {
		t.Fatalf("CreateCareerPathTx failed: %v", err)
	}

	emp := uuid.New()
	db.Exec("INSERT INTO employees (id, name) VALUES (?, 'Andi')", emp.String())
	db.Exec("INSERT INTO employments (id, employee_id, position_id, effective_date, effective_end_date) VALUES (?, ?, ?, '2020-01-01', NULL)", uuid.New().String(), emp.String(), srcStaff.String())

	result, err := svc.GetEligibleEmployeesByPosition(ctx, targetManager.String())
	if err != nil {
		t.Fatalf("GetEligibleEmployeesByPosition failed: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 eligible employee for target, got %d", len(result))
	}
	if result[0].EmployeeID != emp.String() {
		t.Errorf("expected employee %s, got %s", emp, result[0].EmployeeID)
	}
}

func TestService_GetEligibleEmployeesByPosition_NoPathReturnsEmpty(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()
	result, err := svc.GetEligibleEmployeesByPosition(ctx(), uuid.New().String())
	if err != nil {
		t.Fatalf("GetEligibleEmployeesByPosition failed: %v", err)
	}
	if result == nil || len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

// uuidStr generates a random UUID string for test use.
func uuidStr() string {
	return uuid.New().String()
}

// =========================================================================
// Succession Gap (S-5) Service Tests
// =========================================================================

func TestService_GetSuccessionGaps_FlagsExternalRecruitment(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()
	ctx := ctx()

	if err := db.Exec(`CREATE TABLE organizations (
		id CHAR(36) PRIMARY KEY,
		nomenclature VARCHAR(200),
		full_code VARCHAR(50)
	)`).Error; err != nil {
		t.Fatalf("failed to create organizations table: %v", err)
	}
	posGap := uuid.New()
	posReady := uuid.New()
	db.Exec("INSERT INTO organizations (id, nomenclature) VALUES (?, 'CTO')", posGap.String())
	db.Exec("INSERT INTO organizations (id, nomenclature) VALUES (?, 'CFO')", posReady.String())

	// posGap: successor belum siap → requires_external_recruitment = true
	repo.CreateSuccessionPlan(ctx, &CareerSuccessionPlan{
		PositionID: posGap, SuccessorID: uuid.New(), ReadinessLevel: "READY_1YR", Status: "ACTIVE",
	})
	// posReady: ada successor READY_NOW → tidak butuh fallback
	repo.CreateSuccessionPlan(ctx, &CareerSuccessionPlan{
		PositionID: posReady, SuccessorID: uuid.New(), ReadinessLevel: "READY_NOW", Status: "ACTIVE",
	})

	gaps, err := svc.GetSuccessionGaps(ctx)
	if err != nil {
		t.Fatalf("GetSuccessionGaps failed: %v", err)
	}
	if len(gaps) != 2 {
		t.Fatalf("expected 2 key positions, got %d", len(gaps))
	}
	byID := map[string]SuccessionGapResponse{}
	for _, g := range gaps {
		byID[g.PositionID] = g
	}
	if g, ok := byID[posGap.String()]; !ok || !g.RequiresExternalRecruitment {
		t.Error("expected gap position to require external recruitment")
	}
	if g, ok := byID[posGap.String()]; ok && g.HasReadySuccessor {
		t.Error("expected gap position to have no ready successor")
	}
	if g, ok := byID[posReady.String()]; !ok || g.RequiresExternalRecruitment {
		t.Error("expected ready position to NOT require external recruitment")
	}
}

func TestService_GetSuccessionGaps_Empty(t *testing.T) {
	svc, _, db, cleanup := newTestService()
	defer cleanup()
	if err := db.Exec(`CREATE TABLE organizations (
		id CHAR(36) PRIMARY KEY,
		nomenclature VARCHAR(200),
		full_code VARCHAR(50)
	)`).Error; err != nil {
		t.Fatalf("failed to create organizations table: %v", err)
	}

	gaps, err := svc.GetSuccessionGaps(ctx())
	if err != nil {
		t.Fatalf("GetSuccessionGaps failed: %v", err)
	}
	if len(gaps) != 0 {
		t.Errorf("expected empty list, got %d", len(gaps))
	}
}

func TestService_CheckSuccessionGapByPosition(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()
	ctx := ctx()

	posGap := uuid.New()
	posReady := uuid.New()
	repo.CreateSuccessionPlan(ctx, &CareerSuccessionPlan{
		PositionID: posGap, SuccessorID: uuid.New(), ReadinessLevel: "POTENTIAL", Status: "ACTIVE",
	})
	repo.CreateSuccessionPlan(ctx, &CareerSuccessionPlan{
		PositionID: posReady, SuccessorID: uuid.New(), ReadinessLevel: "READY_NOW", Status: "ACTIVE",
	})

	gap, err := svc.CheckSuccessionGapByPosition(ctx, posGap.String())
	if err != nil {
		t.Fatalf("CheckSuccessionGapByPosition failed: %v", err)
	}
	if !gap {
		t.Error("expected gap=true for position without ready successor")
	}

	ready, err := svc.CheckSuccessionGapByPosition(ctx, posReady.String())
	if err != nil {
		t.Fatalf("CheckSuccessionGapByPosition failed: %v", err)
	}
	if ready {
		t.Error("expected gap=false for position with READY_NOW successor")
	}
}

func TestService_CheckSuccessionGapByPosition_InvalidID(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	if _, err := svc.CheckSuccessionGapByPosition(ctx(), "not-a-uuid"); err == nil {
		t.Error("expected error for invalid position_id")
	}
}
