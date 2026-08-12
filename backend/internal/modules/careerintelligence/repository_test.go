package careerintelligence

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func TestRepo_CreateCareerPathTx_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	tenure := 24
	cp := &CareerPath{Name: "PROMOTION: Staff → Supervisor", IsActive: true}
	steps := []CareerPathStep{
		{PositionID: uuid.New(), Sequence: 1},
		{PositionID: uuid.New(), Sequence: 2, PathType: "PROMOTION", TypicalTenure: &tenure, Requirements: "Bachelor degree"},
	}

	if err := repo.CreateCareerPathTx(context.Background(), cp, steps); err != nil {
		t.Fatalf("CreateCareerPathTx failed: %v", err)
	}

	if cp.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
	steps, err := repo.ListCareerPathStepsByPathID(context.Background(), cp.ID)
	if err != nil {
		t.Fatalf("ListCareerPathStepsByPathID failed: %v", err)
	}
	if len(steps) != 2 {
		t.Errorf("expected 2 steps, got %d", len(steps))
	}
	if steps[0].Sequence != 1 || steps[1].Sequence != 2 {
		t.Errorf("expected sequences 1,2; got %d,%d", steps[0].Sequence, steps[1].Sequence)
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

	if found.Name != created.Name {
		t.Errorf("expected name '%s', got '%s'", created.Name, found.Name)
	}
}

func TestRepo_ListCareerPaths_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	for i := 0; i < 3; i++ {
		createTestCareerPath(repo)
	}

	list, total, err := repo.ListCareerPaths(context.Background(), 1, 10, "")
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

func TestRepo_FindCareerPathsByTarget_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	target := uuid.New()
	cp := &CareerPath{Name: "PROMOTION: A → B → C", IsActive: true}
	steps := []CareerPathStep{
		{PositionID: uuid.New(), Sequence: 1},
		{PositionID: uuid.New(), Sequence: 2},
		{PositionID: target, Sequence: 3, PathType: "PROMOTION"},
	}
	if err := repo.CreateCareerPathTx(context.Background(), cp, steps); err != nil {
		t.Fatalf("CreateCareerPathTx failed: %v", err)
	}

	list, err := repo.FindCareerPathsByTarget(context.Background(), target)
	if err != nil {
		t.Fatalf("FindCareerPathsByTarget failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 path targeting position, got %d", len(list))
	}
	if list[0].ID != cp.ID {
		t.Errorf("expected path %s, got %s", cp.ID, list[0].ID)
	}
}

func TestRepo_FindCareerPathsByTarget_IgnoresMidStep(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	midStep := uuid.New()
	cp := &CareerPath{Name: "PROMOTION: A → B → C", IsActive: true}
	steps := []CareerPathStep{
		{PositionID: midStep, Sequence: 1},
		{PositionID: uuid.New(), Sequence: 2},
		{PositionID: uuid.New(), Sequence: 3},
	}
	if err := repo.CreateCareerPathTx(context.Background(), cp, steps); err != nil {
		t.Fatalf("CreateCareerPathTx failed: %v", err)
	}

	// midStep bukan step terakhir → path tidak boleh terpilih sebagai target.
	list, err := repo.FindCareerPathsByTarget(context.Background(), midStep)
	if err != nil {
		t.Fatalf("FindCareerPathsByTarget failed: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected 0 paths (mid-step bukan target), got %d", len(list))
	}
}

// createEligibilityTables membuat tabel employees + employments + positions
// untuk test S-4 eligible employees.
func createEligibilityTables(t *testing.T, db *gorm.DB) {
	t.Helper()
	statements := []string{
		`CREATE TABLE employees (
			id CHAR(36) PRIMARY KEY,
			name VARCHAR(255),
			deleted_at DATETIME NULL
		)`,
		`CREATE TABLE employments (
			id CHAR(36) PRIMARY KEY,
			employee_id CHAR(36),
			position_id CHAR(36),
			effective_date DATE,
			effective_end_date DATE NULL,
			deleted_at DATETIME NULL
		)`,
		`CREATE TABLE positions (
			id CHAR(36) PRIMARY KEY,
			title VARCHAR(200)
		)`,
	}
	for _, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("failed to create eligibility table: %v\n%s", err, stmt)
		}
	}
}

func TestRepo_ListEligibleEmployeesByPositionIDs_ActiveOnly(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	createEligibilityTables(t, db)
	ctx := context.Background()

	posSource := uuid.New()
	posOther := uuid.New()
	past := "2020-01-01"

	// Employee aktif di posisi source
	empActive := uuid.New()
	db.Exec("INSERT INTO employees (id, name) VALUES (?, 'Andi')", empActive.String())
	db.Exec("INSERT INTO employments (id, employee_id, position_id, effective_date, effective_end_date) VALUES (?, ?, ?, ?, NULL)", uuid.New().String(), empActive.String(), posSource.String(), past)
	// Employee dengan kontrak sudah berakhir → TIDAK eligible
	empExpired := uuid.New()
	db.Exec("INSERT INTO employees (id, name) VALUES (?, 'Budi')", empExpired.String())
	db.Exec("INSERT INTO employments (id, employee_id, position_id, effective_date, effective_end_date) VALUES (?, ?, ?, ?, ?)", uuid.New().String(), empExpired.String(), posSource.String(), past, "2021-01-01")
	// Employee di posisi lain → TIDAK eligible untuk posSource
	empOther := uuid.New()
	db.Exec("INSERT INTO employees (id, name) VALUES (?, 'Citra')", empOther.String())
	db.Exec("INSERT INTO employments (id, employee_id, position_id, effective_date, effective_end_date) VALUES (?, ?, ?, ?, NULL)", uuid.New().String(), empOther.String(), posOther.String(), past)

	rows, err := repo.ListEligibleEmployeesByPositionIDs(ctx, []uuid.UUID{posSource})
	if err != nil {
		t.Fatalf("ListEligibleEmployeesByPositionIDs failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 eligible employee (active only), got %d", len(rows))
	}
	if rows[0].EmployeeID != empActive.String() {
		t.Errorf("expected employee %s, got %s", empActive, rows[0].EmployeeID)
	}
	if rows[0].PositionID != posSource.String() {
		t.Errorf("expected position %s, got %s", posSource, rows[0].PositionID)
	}
}

func TestRepo_ListEligibleEmployeesByPositionIDs_EmptyIDs(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	rows, err := repo.ListEligibleEmployeesByPositionIDs(context.Background(), []uuid.UUID{})
	if err != nil {
		t.Fatalf("ListEligibleEmployeesByPositionIDs failed: %v", err)
	}
	if len(rows) != 0 {
		t.Errorf("expected 0 rows for empty ids, got %d", len(rows))
	}
}

func TestRepo_FindCareerPathsBySource_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	src := uuid.New()
	tenure := 12
	cp := &CareerPath{Name: "PROMOTION: A → B", IsActive: true}
	steps := []CareerPathStep{
		{PositionID: src, Sequence: 1},
		{PositionID: uuid.New(), Sequence: 2, PathType: "PROMOTION", TypicalTenure: &tenure},
	}
	if err := repo.CreateCareerPathTx(context.Background(), cp, steps); err != nil {
		t.Fatalf("CreateCareerPathTx failed: %v", err)
	}

	list, err := repo.FindCareerPathsBySource(context.Background(), src)
	if err != nil {
		t.Fatalf("FindCareerPathsBySource failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 path from source, got %d", len(list))
	}
	if list[0].ID != cp.ID {
		t.Errorf("expected path %s, got %s", cp.ID, list[0].ID)
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
	steps, err := repo.ListCareerPathStepsByPathID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("ListCareerPathStepsByPathID after delete failed: %v", err)
	}
	if len(steps) != 0 {
		t.Errorf("expected steps deleted, got %d", len(steps))
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
