package employeemovement

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/modules/employee"
)

// careerTestAdapter mirrors the employeeCareerAdapter wiring in cmd/server/main.go:
// it wraps a real employee.Repository so the ExecuteMovement integration test
// exercises the full G-1 flow against actual tables (not mocks).
type careerTestAdapter struct {
	repo *employee.Repository
}

func (a careerTestAdapter) FindCurrentEmployment(ctx context.Context, tx *gorm.DB, employeeID uuid.UUID) (*CareerEmployment, error) {
	e, err := a.repo.FindActiveEmploymentByEmployeeIDTx(ctx, tx, employeeID)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, nil
	}
	return &CareerEmployment{
		ID:                   e.ID,
		OrganizationID:       e.OrganizationID,
		PositionID:           e.PositionID,
		EmploymentStatusID:   e.EmploymentStatusID,
		DecisionLetterNumber: e.DecisionLetterNumber,
		DecisionLetterDate:   e.DecisionLetterDate,
		EffectiveDate:        e.EffectiveDate,
	}, nil
}

func (a careerTestAdapter) CloseEmployment(ctx context.Context, tx *gorm.DB, employmentID uuid.UUID, effectiveDate string) error {
	return a.repo.CloseEmploymentTx(ctx, tx, employmentID, effectiveDate)
}

func (a careerTestAdapter) CreateEmployment(ctx context.Context, tx *gorm.DB, employeeID uuid.UUID, data CareerEmployment) (uuid.UUID, error) {
	emp := &employee.Employment{
		EmployeeID:           &employeeID,
		OrganizationID:       data.OrganizationID,
		PositionID:           data.PositionID,
		EmploymentStatusID:   data.EmploymentStatusID,
		DecisionLetterNumber: data.DecisionLetterNumber,
		DecisionLetterDate:   data.DecisionLetterDate,
		EffectiveDate:        data.EffectiveDate,
	}
	if err := a.repo.CreateEmploymentTx(ctx, tx, emp); err != nil {
		return uuid.Nil, err
	}
	return emp.ID, nil
}

func (a careerTestAdapter) SetEmployeeInactive(ctx context.Context, tx *gorm.DB, employeeID uuid.UUID) error {
	return a.repo.SetEmployeeStatusTx(ctx, tx, employeeID, "inactive")
}

// =========================================================================
// Employee Movement Repository Tests
// =========================================================================

func TestRepo_CreateMovement_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	m := &EmployeeMovement{
		EmployeeID:           uuid.New(),
		MovementType:         MovementTypePromotion,
		DecisionLetterNumber: "SK-001",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-01",
	}

	if err := repo.CreateMovement(ctx, m); err != nil {
		t.Fatalf("CreateMovement failed: %v", err)
	}

	if m.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindMovementByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New())

	found, err := repo.FindMovementByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindMovementByID failed: %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("expected ID '%s', got '%s'", created.ID, found.ID)
	}
}

func TestRepo_FindMovementByID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	_, err := repo.FindMovementByID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent movement")
	}
}

func TestRepo_ListMovements_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	empID := uuid.New()
	for i := 0; i < 5; i++ {
		createTestMovement(repo, empID)
	}

	movements, total, err := repo.ListMovements(ctx, 1, 3, "", "", "")
	if err != nil {
		t.Fatalf("ListMovements failed: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(movements) != 3 {
		t.Errorf("expected 3 movements (page 1 of 3), got %d", len(movements))
	}
}

func TestRepo_FindMovementsByEmployeeID(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp1 := uuid.New()
	emp2 := uuid.New()

	createTestMovement(repo, emp1)
	createTestMovement(repo, emp1)
	createTestMovement(repo, emp2)

	movements, total, err := repo.FindMovementsByEmployeeID(ctx, emp1, 1, 10)
	if err != nil {
		t.Fatalf("FindMovementsByEmployeeID failed: %v", err)
	}

	if total != 2 {
		t.Errorf("expected total 2 for emp1, got %d", total)
	}
	if len(movements) != 2 {
		t.Errorf("expected 2 movements, got %d", len(movements))
	}
}

func TestRepo_UpdateMovement_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New())

	created.Reason = strPtr("Updated reason")
	if err := repo.UpdateMovement(ctx, created); err != nil {
		t.Fatalf("UpdateMovement failed: %v", err)
	}

	found, _ := repo.FindMovementByID(ctx, created.ID)
	if found.Reason == nil || *found.Reason != "Updated reason" {
		t.Errorf("expected reason 'Updated reason', got '%v'", found.Reason)
	}
}

func TestRepo_DeleteMovement_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New())

	if err := repo.DeleteMovement(ctx, created.ID); err != nil {
		t.Fatalf("DeleteMovement failed: %v", err)
	}

	_, err := repo.FindMovementByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting movement")
	}
}

func TestRepo_ExecuteMovement_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New())
	created.Status = MovementStatusApproved
	repo.UpdateMovement(ctx, created)

	executorID := uuid.New()
	if err := repo.ExecuteMovement(ctx, created.ID, executorID, nil); err != nil {
		t.Fatalf("ExecuteMovement failed: %v", err)
	}

	found, _ := repo.FindMovementByID(ctx, created.ID)
	if found.Status != MovementStatusExecuted {
		t.Errorf("expected status 'executed', got '%s'", found.Status)
	}
}

func TestRepo_ExecuteMovement_NonApproved_Error(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New())
	// Still draft — cannot execute without approval

	err := repo.ExecuteMovement(ctx, created.ID, uuid.New(), nil)
	if err == nil {
		t.Fatal("expected error when executing non-approved movement")
	}
}

func TestRepo_CancelMovement_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New())

	if err := repo.CancelMovement(ctx, created.ID); err != nil {
		t.Fatalf("CancelMovement failed: %v", err)
	}

	found, _ := repo.FindMovementByID(ctx, created.ID)
	if found.Status != MovementStatusCancelled {
		t.Errorf("expected status 'cancelled', got '%s'", found.Status)
	}
}

func TestRepo_CancelMovement_Executed_Error(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New())
	created.Status = MovementStatusExecuted
	repo.UpdateMovement(ctx, created)

	err := repo.CancelMovement(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error when cancelling executed movement")
	}
}

// TestRepo_CancelMovement_Approved_Error verifies plan §12.16: an approved
// movement cannot be cancelled directly at the repo level — it must go through
// the cancellation approval flow (SetCancellationRequested).
func TestRepo_CancelMovement_Approved_Error(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New())
	created.Status = MovementStatusApproved
	repo.UpdateMovement(ctx, created)

	err := repo.CancelMovement(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error when directly cancelling an approved movement")
	}

	// Status untouched.
	found, _ := repo.FindMovementByID(ctx, created.ID)
	if found.Status != MovementStatusApproved {
		t.Errorf("expected status to remain 'approved', got '%s'", found.Status)
	}
}

// TestRepo_SetCancellationRequested_Success verifies the cancellation request
// persistence (plan §12.16): approved → cancellation_pending + instance id.
func TestRepo_SetCancellationRequested_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New())
	created.Status = MovementStatusApproved
	repo.UpdateMovement(ctx, created)

	instanceID := uuid.New()
	reason := "Dibatalkan atas permintaan manajemen"
	if err := repo.SetCancellationRequested(ctx, created.ID, instanceID, &reason); err != nil {
		t.Fatalf("SetCancellationRequested failed: %v", err)
	}

	found, _ := repo.FindMovementByID(ctx, created.ID)
	if found.Status != MovementStatusCancellationPending {
		t.Errorf("expected status 'cancellation_pending', got '%s'", found.Status)
	}
	if found.CancellationApprovalInstanceID == nil || *found.CancellationApprovalInstanceID != instanceID {
		t.Errorf("expected cancellation_approval_instance_id %s, got %v", instanceID, found.CancellationApprovalInstanceID)
	}
	if found.Notes == nil || *found.Notes != reason {
		t.Errorf("expected notes '%s', got %v", reason, found.Notes)
	}
}

// TestRepo_SetCancellationRequested_NonApproved_Error verifies only approved
// movements can enter the cancellation approval flow.
func TestRepo_SetCancellationRequested_NonApproved_Error(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestMovement(repo, uuid.New()) // draft

	err := repo.SetCancellationRequested(ctx, created.ID, uuid.New(), nil)
	if err == nil {
		t.Fatal("expected error when requesting cancellation for a non-approved movement")
	}

	found, _ := repo.FindMovementByID(ctx, created.ID)
	if found.Status != MovementStatusDraft {
		t.Errorf("expected status to remain 'draft', got '%s'", found.Status)
	}
}

// =========================================================================
// Employee Contract Repository Tests
// =========================================================================

func TestRepo_CreateContract_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	c := &EmployeeContract{
		EmployeeID:     uuid.New(),
		ContractNumber: "CTR-001",
		ContractType:   ContractTypePKWT,
		StartDate:      "2026-01-01",
	}

	if err := repo.CreateContract(ctx, c); err != nil {
		t.Fatalf("CreateContract failed: %v", err)
	}

	if c.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindContractByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestContract(repo, uuid.New())

	found, err := repo.FindContractByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindContractByID failed: %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("expected ID '%s', got '%s'", created.ID, found.ID)
	}
}

func TestRepo_FindContractByID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	_, err := repo.FindContractByID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent contract")
	}
}

func TestRepo_UpdateContract_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestContract(repo, uuid.New())

	created.Status = ContractStatusExpired
	if err := repo.UpdateContract(ctx, created); err != nil {
		t.Fatalf("UpdateContract failed: %v", err)
	}

	found, _ := repo.FindContractByID(ctx, created.ID)
	if found.Status != ContractStatusExpired {
		t.Errorf("expected status 'expired', got '%s'", found.Status)
	}
}

func TestRepo_DeleteContract_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestContract(repo, uuid.New())

	if err := repo.DeleteContract(ctx, created.ID); err != nil {
		t.Fatalf("DeleteContract failed: %v", err)
	}

	_, err := repo.FindContractByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting contract")
	}
}

func TestRepo_ListContracts_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	empID := uuid.New()
	for i := 0; i < 5; i++ {
		createTestContract(repo, empID)
	}

	contracts, total, err := repo.ListContracts(ctx, 1, 2, "", "")
	if err != nil {
		t.Fatalf("ListContracts failed: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(contracts) != 2 {
		t.Errorf("expected 2 contracts (page 1 of 3), got %d", len(contracts))
	}
}

func TestRepo_ListContracts_Filters(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	empID := uuid.New()
	// Seed minimal employees row agar JOIN search berfungsi.
	db, err := dbResolver(ctx)
	if err != nil {
		t.Fatalf("dbResolver failed: %v", err)
	}
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS employees (id CHAR(36) PRIMARY KEY, employee_id VARCHAR(50) NOT NULL, name VARCHAR(255) NOT NULL)`).Error; err != nil {
		t.Fatalf("create employees table failed: %v", err)
	}
	if err := db.Exec(`INSERT INTO employees (id, employee_id, name) VALUES (?, ?, ?)`, empID.String(), "EMP-CON-001", "Budi Santoso").Error; err != nil {
		t.Fatalf("seed employee failed: %v", err)
	}

	// 2 aktif, 1 terminated.
	c1 := createTestContract(repo, empID)
	createTestContract(repo, empID)
	c3 := createTestContract(repo, empID)
	// Set status terminated untuk c3 via direct update (helper selalu active).
	if err := db.Model(&EmployeeContract{}).Where("id = ?", c3.ID).Update("status", "terminated").Error; err != nil {
		t.Fatalf("update status failed: %v", err)
	}

	// Filter status.
	contracts, total, err := repo.ListContracts(ctx, 1, 10, "terminated", "")
	if err != nil {
		t.Fatalf("ListContracts(status) failed: %v", err)
	}
	if total != 1 || len(contracts) != 1 || contracts[0].ID != c3.ID {
		t.Errorf("status filter: expected 1 terminated contract, got total=%d len=%d", total, len(contracts))
	}

	// Filter search by contract number.
	contracts, total, err = repo.ListContracts(ctx, 1, 10, "", c1.ContractNumber)
	if err != nil {
		t.Fatalf("ListContracts(search) failed: %v", err)
	}
	if total != 1 || len(contracts) != 1 || contracts[0].ID != c1.ID {
		t.Errorf("search filter: expected 1 contract by number, got total=%d len=%d", total, len(contracts))
	}

	// Filter kombinasi tidak boleh cocok.
	contracts, total, err = repo.ListContracts(ctx, 1, 10, "terminated", c1.ContractNumber)
	if err != nil {
		t.Fatalf("ListContracts(combo) failed: %v", err)
	}
	if total != 0 || len(contracts) != 0 {
		t.Errorf("combo filter: expected 0, got total=%d len=%d", total, len(contracts))
	}
}

func TestRepo_FindContractsByEmployeeID(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp1 := uuid.New()
	emp2 := uuid.New()

	createTestContract(repo, emp1)
	createTestContract(repo, emp1)
	createTestContract(repo, emp1)
	createTestContract(repo, emp2)

	contracts, total, err := repo.FindContractsByEmployeeID(ctx, emp1, 1, 10)
	if err != nil {
		t.Fatalf("FindContractsByEmployeeID failed: %v", err)
	}

	if total != 3 {
		t.Errorf("expected total 3 for emp1, got %d", total)
	}
	if len(contracts) != 3 {
		t.Errorf("expected 3 contracts, got %d", len(contracts))
	}
}

// =========================================================================
// G-1 Integration Test (real employee repository + adapter, SQLite)
// =========================================================================

func TestService_ExecuteMovement_Integration_RealEmployment(t *testing.T) {
	// Shared in-memory DB for both modules' repositories.
	db, _, cleanup := setupTestDB()
	defer cleanup()

	// Auto-migrate employee tables the movement flow touches.
	if err := db.AutoMigrate(&employee.Employee{}, &employee.Employment{}); err != nil {
		t.Fatalf("failed to migrate employee tables: %v", err)
	}

	ctx := context.Background()
	employeeRepo := employee.NewRepository(func(context.Context) (*gorm.DB, error) { return db, nil })
	movementRepo := NewRepository(func(context.Context) (*gorm.DB, error) { return db, nil })

	// Seed an active employee + current employment.
	emp := &employee.Employee{
		EmployeeID: "EMP-G1-001",
		Name:       "G1 Test Employee",
		Status:     "active",
	}
	if err := employeeRepo.CreateEmployee(ctx, emp); err != nil {
		t.Fatalf("CreateEmployee failed: %v", err)
	}

	oldEmployment := &employee.Employment{
		EmployeeID:           &emp.ID,
		OrganizationID:       uuidPtr(uuid.New()),
		DecisionLetterNumber: "SK-OLD-001",
		DecisionLetterDate:   "2026-01-01",
		EffectiveDate:        "2026-01-01",
	}
	if err := employeeRepo.CreateEmployment(ctx, oldEmployment); err != nil {
		t.Fatalf("CreateEmployment failed: %v", err)
	}

	// Service wired with the real-adapter (mirrors main.go wiring).
	svc := NewService(movementRepo, testLogger())
	svc.SetCareerExecutor(careerTestAdapter{repo: employeeRepo})

	// Approved promotion movement.
	newOrg := uuid.New()
	newPos := uuid.New()
	movement := createTestMovement(movementRepo, emp.ID)
	movement.Status = MovementStatusApproved
	movement.MovementType = MovementTypePromotion
	movement.ToOrganizationID = &newOrg
	movement.ToPositionID = &newPos
	movement.EffectiveDate = "2026-08-01"
	if err := movementRepo.UpdateMovement(ctx, movement); err != nil {
		t.Fatalf("UpdateMovement failed: %v", err)
	}

	if err := svc.ExecuteMovement(ctx, movement.ID.String(), uuid.New().String()); err != nil {
		t.Fatalf("ExecuteMovement failed: %v", err)
	}

	// 1) Old employment closed with effective_end_date = 2026-07-31.
	closed, err := employeeRepo.FindEmploymentByID(ctx, oldEmployment.ID)
	if err != nil {
		t.Fatalf("FindEmploymentByID failed: %v", err)
	}
	if closed.EffectiveEndDate == nil {
		t.Fatal("expected old employment to have effective_end_date")
	}
	// SQLite returns DATE columns as RFC3339 timestamps; compare date part only.
	if !dateStartsWith(*closed.EffectiveEndDate, "2026-07-31") {
		t.Errorf("expected effective_end_date 2026-07-31, got '%s'", *closed.EffectiveEndDate)
	}

	// 2) New employment created with movement's to_* fields.
	active, err := employeeRepo.FindActiveEmploymentByEmployeeID(ctx, emp.ID)
	if err != nil {
		t.Fatalf("FindActiveEmploymentByEmployeeID failed: %v", err)
	}
	if active == nil {
		t.Fatal("expected a new active employment after execution")
	}
	if active.OrganizationID == nil || *active.OrganizationID != newOrg {
		t.Errorf("expected organization_id %s, got %v", newOrg, active.OrganizationID)
	}
	if active.PositionID == nil || *active.PositionID != newPos {
		t.Errorf("expected position_id %s, got %v", newPos, active.PositionID)
	}

	// 3) Movement executed with to_employment_id persisted.
	executed, _ := movementRepo.FindMovementByID(ctx, movement.ID)
	if executed.Status != MovementStatusExecuted {
		t.Errorf("expected status 'executed', got '%s'", executed.Status)
	}
	if executed.ToEmploymentID == nil || *executed.ToEmploymentID != active.ID {
		t.Errorf("expected to_employment_id %s, got %v", active.ID, executed.ToEmploymentID)
	}
}

func TestService_ExecuteMovement_Integration_Offboarding(t *testing.T) {
	db, _, cleanup := setupTestDB()
	defer cleanup()
	if err := db.AutoMigrate(&employee.Employee{}, &employee.Employment{}); err != nil {
		t.Fatalf("failed to migrate employee tables: %v", err)
	}

	ctx := context.Background()
	employeeRepo := employee.NewRepository(func(context.Context) (*gorm.DB, error) { return db, nil })
	movementRepo := NewRepository(func(context.Context) (*gorm.DB, error) { return db, nil })

	emp := &employee.Employee{
		EmployeeID: "EMP-G1-002",
		Name:       "G1 Offboarding Employee",
		Status:     "active",
	}
	if err := employeeRepo.CreateEmployee(ctx, emp); err != nil {
		t.Fatalf("CreateEmployee failed: %v", err)
	}

	oldEmployment := &employee.Employment{
		EmployeeID:           &emp.ID,
		OrganizationID:       uuidPtr(uuid.New()),
		DecisionLetterNumber: "SK-OLD-002",
		DecisionLetterDate:   "2026-01-01",
		EffectiveDate:        "2026-01-01",
	}
	if err := employeeRepo.CreateEmployment(ctx, oldEmployment); err != nil {
		t.Fatalf("CreateEmployment failed: %v", err)
	}

	svc := NewService(movementRepo, testLogger())
	svc.SetCareerExecutor(careerTestAdapter{repo: employeeRepo})

	movement := createTestMovement(movementRepo, emp.ID)
	movement.Status = MovementStatusApproved
	movement.MovementType = MovementTypeOffboarding
	movement.EffectiveDate = "2026-08-01"
	if err := movementRepo.UpdateMovement(ctx, movement); err != nil {
		t.Fatalf("UpdateMovement failed: %v", err)
	}

	if err := svc.ExecuteMovement(ctx, movement.ID.String(), uuid.New().String()); err != nil {
		t.Fatalf("ExecuteMovement failed: %v", err)
	}

	// Employment closed, no new employment.
	closed, _ := employeeRepo.FindEmploymentByID(ctx, oldEmployment.ID)
	if closed.EffectiveEndDate == nil {
		t.Fatal("expected old employment to have effective_end_date")
	}
	active, _ := employeeRepo.FindActiveEmploymentByEmployeeID(ctx, emp.ID)
	if active != nil {
		t.Error("offboarding should not leave an active employment")
	}

	// Employee marked inactive (query status column directly to avoid the
	// heavy Preload chain in FindEmployeeByID).
	var empStatus string
	if err := db.Model(&employee.Employee{}).Where("id = ?", emp.ID).Pluck("status", &empStatus).Error; err != nil {
		t.Fatalf("failed to query employee status: %v", err)
	}
	if empStatus != "inactive" {
		t.Errorf("expected employee status 'inactive', got '%s'", empStatus)
	}
}

// TestRepo_ExtendContract_ChainCount verifies G-6: extension_count is derived
// from the previous contract (previous + 1), so chained extensions accumulate
// instead of being hardcoded to 1.
func TestRepo_ExtendContract_ChainCount(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	empID := uuid.New()
	c1 := &EmployeeContract{
		EmployeeID:     empID,
		ContractNumber: "CTR-001",
		ContractType:   ContractTypePKWT,
		StartDate:      "2026-01-01",
		Status:         ContractStatusActive,
	}
	if err := repo.CreateContract(ctx, c1); err != nil {
		t.Fatalf("CreateContract failed: %v", err)
	}

	c2 := &EmployeeContract{
		EmployeeID:         empID,
		ContractNumber:     "CTR-002",
		ContractType:       ContractTypePKWT,
		StartDate:          "2026-06-01",
		Status:             ContractStatusActive,
		PreviousContractID: &c1.ID,
	}
	if err := repo.ExtendContract(ctx, c2, c1.ID); err != nil {
		t.Fatalf("ExtendContract (1st) failed: %v", err)
	}
	if c2.ExtensionCount != 1 {
		t.Errorf("expected 1st extension_count 1, got %d", c2.ExtensionCount)
	}

	c3 := &EmployeeContract{
		EmployeeID:         empID,
		ContractNumber:     "CTR-003",
		ContractType:       ContractTypePKWT,
		StartDate:          "2027-01-01",
		Status:             ContractStatusActive,
		PreviousContractID: &c2.ID,
	}
	if err := repo.ExtendContract(ctx, c3, c2.ID); err != nil {
		t.Fatalf("ExtendContract (2nd) failed: %v", err)
	}
	if c3.ExtensionCount != 2 {
		t.Errorf("expected 2nd extension_count 2, got %d", c3.ExtensionCount)
	}

	// Previous contracts must be marked extended.
	var status1, status2 string
	if err := db.Model(&EmployeeContract{}).Where("id = ?", c1.ID).Pluck("status", &status1).Error; err != nil {
		t.Fatalf("failed to query c1 status: %v", err)
	}
	if status1 != string(ContractStatusExtended) {
		t.Errorf("expected c1 status 'extended', got '%s'", status1)
	}
	if err := db.Model(&EmployeeContract{}).Where("id = ?", c2.ID).Pluck("status", &status2).Error; err != nil {
		t.Fatalf("failed to query c2 status: %v", err)
	}
	if status2 != string(ContractStatusExtended) {
		t.Errorf("expected c2 status 'extended', got '%s'", status2)
	}
}

// TestRepo_ExtendContract_MissingPrevious verifies extending a non-existent
// contract fails and rolls back (no orphan contract created).
func TestRepo_ExtendContract_MissingPrevious(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	newC := &EmployeeContract{
		EmployeeID:         uuid.New(),
		ContractNumber:     "CTR-X",
		ContractType:       ContractTypePKWT,
		StartDate:          "2026-06-01",
		Status:             ContractStatusActive,
		PreviousContractID: uuidPtr(uuid.New()),
	}
	err := repo.ExtendContract(ctx, newC, *newC.PreviousContractID)
	if err == nil {
		t.Fatal("expected error when previous contract does not exist")
	}

	var count int64
	if err := db.Model(&EmployeeContract{}).Count(&count).Error; err != nil {
		t.Fatalf("failed to count contracts: %v", err)
	}
	if count != 0 {
		t.Errorf("expected no contracts persisted after rollback, got %d", count)
	}
}

// TestRepo_ListMovements_Filters verifies movement_type/status/search filters
// on the list endpoint (halaman FE Movements — langkah 9 plan).
func TestRepo_ListMovements_Filters(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	empID := uuid.New()
	// Seed a minimal employees row so the search JOIN works.
	db.Exec(`CREATE TABLE IF NOT EXISTS employees (id CHAR(36) PRIMARY KEY, employee_id VARCHAR(50) NOT NULL, name VARCHAR(255) NOT NULL)`)
	db.Exec(`INSERT INTO employees (id, employee_id, name) VALUES (?, ?, ?)`, empID.String(), "EMP-001", "Budi Santoso")

	m1 := createTestMovement(repo, empID) // type other, draft
	m2 := createTestMovement(repo, empID)
	m2.MovementType = MovementTypePromotion
	m2.Status = MovementStatusApproved
	if err := repo.UpdateMovement(ctx, m2); err != nil {
		t.Fatalf("failed to update m2: %v", err)
	}

	// Filter by movement_type.
	promotions, total, err := repo.ListMovements(ctx, 1, 20, "promotion", "", "")
	if err != nil {
		t.Fatalf("ListMovements by type failed: %v", err)
	}
	if total != 1 || promotions[0].ID != m2.ID {
		t.Errorf("expected 1 promotion, got %d", total)
	}

	// Filter by status.
	approved, total, err := repo.ListMovements(ctx, 1, 20, "", "approved", "")
	if err != nil {
		t.Fatalf("ListMovements by status failed: %v", err)
	}
	if total != 1 || approved[0].ID != m2.ID {
		t.Errorf("expected 1 approved, got %d", total)
	}

	// Search by employee name.
	_, total, err = repo.ListMovements(ctx, 1, 20, "", "", "Budi")
	if err != nil {
		t.Fatalf("ListMovements search failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2 movements for search 'Budi', got %d", total)
	}

	// Combined + no match.
	none, total, err := repo.ListMovements(ctx, 1, 20, "mutation", "executed", "")
	if err != nil {
		t.Fatalf("ListMovements combined failed: %v", err)
	}
	if total != 0 || len(none) != 0 {
		t.Errorf("expected 0 movements, got %d", total)
	}

	_ = m1
}

// uuidPtr returns a pointer to the given UUID (test helper).
func uuidPtr(u uuid.UUID) *uuid.UUID {
	return &u
}

// dateStartsWith reports whether s begins with the YYYY-MM-DD prefix (SQLite
// returns DATE columns as RFC3339 timestamps, MySQL as plain dates).
func dateStartsWith(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}
