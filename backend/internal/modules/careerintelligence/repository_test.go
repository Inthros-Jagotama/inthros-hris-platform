package careerintelligence

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// =========================================================================
// Talent Map Repository Tests
// =========================================================================

func TestRepo_CreateTalentMap_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	tm := &CareerTalentMap{
		EmployeeID:  uuid.New(),
		Period:      "2026-07",
		Performance: "HIGH",
		Potential:   "HIGH",
		GridPosition: "9-BOX-1",
	}

	if err := repo.CreateTalentMap(context.Background(), tm); err != nil {
		t.Fatalf("CreateTalentMap failed: %v", err)
	}

	if tm.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
	if tm.GridPosition != "9-BOX-1" {
		t.Errorf("expected GridPosition '9-BOX-1', got '%s'", tm.GridPosition)
	}
}

func TestRepo_FindTalentMapByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestTalentMap(repo)

	found, err := repo.FindTalentMapByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("FindTalentMapByID failed: %v", err)
	}

	if found.Performance != created.Performance {
		t.Errorf("expected performance '%s', got '%s'", created.Performance, found.Performance)
	}
}

func TestRepo_FindTalentMapByID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	_, err := repo.FindTalentMapByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent talent map")
	}
}

func TestRepo_ListTalentMaps_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	for i := 0; i < 3; i++ {
		createTestTalentMap(repo)
	}

	list, total, err := repo.ListTalentMaps(context.Background(), "", "", 1, 10)
	if err != nil {
		t.Fatalf("ListTalentMaps failed: %v", err)
	}

	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(list) != 3 {
		t.Errorf("expected 3 items, got %d", len(list))
	}
}

func TestRepo_ListTalentMaps_ByPeriod(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	tm1 := createTestTalentMap(repo)
	tm1.Period = "2026-07"
	repo.UpdateTalentMap(context.Background(), tm1)

	tm2 := createTestTalentMap(repo)
	tm2.Period = "2026-07"
	repo.UpdateTalentMap(context.Background(), tm2)

	createTestTalentMap(repo) // different period

	list, total, err := repo.ListTalentMaps(context.Background(), "2026-07", "", 1, 10)
	if err != nil {
		t.Fatalf("ListTalentMaps by period failed: %v", err)
	}

	if total != 2 {
		t.Errorf("expected 2 items for period 2026-07, got %d", total)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}
}

func TestRepo_UpdateTalentMap_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestTalentMap(repo)
	created.Performance = "LOW"
	created.Notes = "Updated notes"

	if err := repo.UpdateTalentMap(context.Background(), created); err != nil {
		t.Fatalf("UpdateTalentMap failed: %v", err)
	}

	found, _ := repo.FindTalentMapByID(context.Background(), created.ID)
	if found.Performance != "LOW" {
		t.Errorf("expected performance 'LOW', got '%s'", found.Performance)
	}
	if found.Notes != "Updated notes" {
		t.Errorf("expected notes 'Updated notes', got '%s'", found.Notes)
	}
}

func TestRepo_DeleteTalentMap_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestTalentMap(repo)

	if err := repo.DeleteTalentMap(context.Background(), created.ID); err != nil {
		t.Fatalf("DeleteTalentMap failed: %v", err)
	}

	_, err := repo.FindTalentMapByID(context.Background(), created.ID)
	if err == nil {
		t.Fatal("expected error after deleting talent map")
	}
}

func TestRepo_GetTalentGrid_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	createTestTalentMap(repo)
	createTestTalentMap(repo)
	createTestTalentMap(repo)

	grid, err := repo.GetTalentGrid(context.Background(), "")
	if err != nil {
		t.Fatalf("GetTalentGrid failed: %v", err)
	}

	if len(grid) != 3 {
		t.Errorf("expected 3 grid entries, got %d", len(grid))
	}
}

// =========================================================================
// Career Interest Repository Tests
// =========================================================================

func TestRepo_CreateCareerInterest_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	ci := &CareerInterest{
		EmployeeID:   uuid.New(),
		InterestType: "LEADERSHIP",
		IsActive:     true,
	}

	if err := repo.CreateCareerInterest(context.Background(), ci); err != nil {
		t.Fatalf("CreateCareerInterest failed: %v", err)
	}

	if ci.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindCareerInterestByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestCareerInterest(repo)

	found, err := repo.FindCareerInterestByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("FindCareerInterestByID failed: %v", err)
	}

	if found.InterestType != created.InterestType {
		t.Errorf("expected type '%s', got '%s'", created.InterestType, found.InterestType)
	}
}

func TestRepo_FindCareerInterestByID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	_, err := repo.FindCareerInterestByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent interest")
	}
}

func TestRepo_ListCareerInterests_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	for i := 0; i < 3; i++ {
		createTestCareerInterest(repo)
	}

	list, total, err := repo.ListCareerInterests(context.Background(), "", 1, 10)
	if err != nil {
		t.Fatalf("ListCareerInterests failed: %v", err)
	}

	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(list) != 3 {
		t.Errorf("expected 3 items, got %d", len(list))
	}
}

func TestRepo_ListCareerInterests_ByEmployee(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	empID := uuid.New()
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		repo.CreateCareerInterest(ctx, &CareerInterest{
			EmployeeID:   empID,
			InterestType: "LEADERSHIP",
			IsActive:     true,
		})
	}
	createTestCareerInterest(repo) // different employee

	list, total, err := repo.ListCareerInterests(context.Background(), empID.String(), 1, 10)
	if err != nil {
		t.Fatalf("ListCareerInterests by employee failed: %v", err)
	}

	if total != 2 {
		t.Errorf("expected 2 items for employee, got %d", total)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}
}

// =========================================================================
// Career Path Repository Tests
// =========================================================================

func TestRepo_CreateCareerPath_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	cp := &CareerPath{
		SourceTitleID:  uuid.New(),
		TargetTitleID:  uuid.New(),
		PathType:       "PROMOTION",
		TypicalTenure:  24,
		IsActive:       true,
	}

	if err := repo.CreateCareerPath(context.Background(), cp); err != nil {
		t.Fatalf("CreateCareerPath failed: %v", err)
	}

	if cp.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindCareerPathByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestCareerPath(repo)

	found, err := repo.FindCareerPathByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("FindCareerPathByID failed: %v", err)
	}

	if found.PathType != created.PathType {
		t.Errorf("expected path type '%s', got '%s'", created.PathType, found.PathType)
	}
}

func TestRepo_ListCareerPaths_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	for i := 0; i < 3; i++ {
		createTestCareerPath(repo)
	}

	list, total, err := repo.ListCareerPaths(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("ListCareerPaths failed: %v", err)
	}

	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(list) != 3 {
		t.Errorf("expected 3 items, got %d", len(list))
	}
}

func TestRepo_DeleteCareerPath_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestCareerPath(repo)

	if err := repo.DeleteCareerPath(context.Background(), created.ID); err != nil {
		t.Fatalf("DeleteCareerPath failed: %v", err)
	}

	_, err := repo.FindCareerPathByID(context.Background(), created.ID)
	if err == nil {
		t.Fatal("expected error after deleting career path")
	}
}

// =========================================================================
// Succession Plan Repository Tests
// =========================================================================

func TestRepo_CreateSuccessionPlan_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	sp := &CareerSuccessionPlan{
		PositionID:     uuid.New(),
		SuccessorID:    uuid.New(),
		ReadinessLevel: "READY_NOW",
		Status:         "ACTIVE",
	}

	if err := repo.CreateSuccessionPlan(context.Background(), sp); err != nil {
		t.Fatalf("CreateSuccessionPlan failed: %v", err)
	}

	if sp.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindSuccessionPlanByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestSuccessionPlan(repo)

	found, err := repo.FindSuccessionPlanByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("FindSuccessionPlanByID failed: %v", err)
	}

	if found.ReadinessLevel != created.ReadinessLevel {
		t.Errorf("expected readiness '%s', got '%s'", created.ReadinessLevel, found.ReadinessLevel)
	}
}

func TestRepo_FindSuccessionPlanByID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	_, err := repo.FindSuccessionPlanByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent succession plan")
	}
}

func TestRepo_ListSuccessionPlans_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	for i := 0; i < 3; i++ {
		createTestSuccessionPlan(repo)
	}

	list, total, err := repo.ListSuccessionPlans(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("ListSuccessionPlans failed: %v", err)
	}

	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(list) != 3 {
		t.Errorf("expected 3 items, got %d", len(list))
	}
}

func TestRepo_UpdateSuccessionPlan_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestSuccessionPlan(repo)
	created.ReadinessLevel = "READY_2YR"
	created.Notes = "Updated notes"

	if err := repo.UpdateSuccessionPlan(context.Background(), created); err != nil {
		t.Fatalf("UpdateSuccessionPlan failed: %v", err)
	}

	found, _ := repo.FindSuccessionPlanByID(context.Background(), created.ID)
	if found.ReadinessLevel != "READY_2YR" {
		t.Errorf("expected readiness 'READY_2YR', got '%s'", found.ReadinessLevel)
	}
	if found.Notes != "Updated notes" {
		t.Errorf("expected notes 'Updated notes', got '%s'", found.Notes)
	}
}

func TestRepo_DeleteSuccessionPlan_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestSuccessionPlan(repo)

	if err := repo.DeleteSuccessionPlan(context.Background(), created.ID); err != nil {
		t.Fatalf("DeleteSuccessionPlan failed: %v", err)
	}

	_, err := repo.FindSuccessionPlanByID(context.Background(), created.ID)
	if err == nil {
		t.Fatal("expected error after deleting succession plan")
	}
}
