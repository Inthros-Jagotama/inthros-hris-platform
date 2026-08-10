package employeemovement

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/modules/employee"
)

func ctx() context.Context {
	return context.Background()
}

// fakeCareerExecutor captures the HR data changes ExecuteMovement pushes
// through the CareerExecutor interface (plan G-1). The tx parameter is
// ignored — the fake only records calls. createErr lets tests simulate a
// mid-transaction failure to verify atomic rollback (plan §12.2).
type fakeCareerExecutor struct {
	currentEmployment *CareerEmployment
	closedID          *uuid.UUID
	closedEndDate     string
	createdData       *CareerEmployment
	createdEmployeeID *uuid.UUID
	createdID         uuid.UUID
	inactiveEmployee  *uuid.UUID
	createErr         error
}

func (f *fakeCareerExecutor) FindCurrentEmployment(_ context.Context, _ *gorm.DB, employeeID uuid.UUID) (*CareerEmployment, error) {
	return f.currentEmployment, nil
}

func (f *fakeCareerExecutor) CloseEmployment(_ context.Context, _ *gorm.DB, employmentID uuid.UUID, effectiveDate string) error {
	f.closedID = &employmentID
	f.closedEndDate = effectiveDate
	return nil
}

func (f *fakeCareerExecutor) CreateEmployment(_ context.Context, _ *gorm.DB, employeeID uuid.UUID, data CareerEmployment) (uuid.UUID, error) {
	f.createdEmployeeID = &employeeID
	f.createdData = &data
	if f.createErr != nil {
		return uuid.Nil, f.createErr
	}
	f.createdID = uuid.New()
	return f.createdID, nil
}

func (f *fakeCareerExecutor) SetEmployeeInactive(_ context.Context, _ *gorm.DB, employeeID uuid.UUID) error {
	f.inactiveEmployee = &employeeID
	return nil
}

// seedEmployment inserts a raw employment row via the employee repository into
// the test DB shared by the movement repository (both resolvers point at the
// same *gorm.DB from setupTestDB).
func seedEmployment(t *testing.T, repo *Repository, emp *employee.Employment) {
	t.Helper()
	db, err := repo.getDB(ctx())
	if err != nil {
		t.Fatalf("failed to get test db: %v", err)
	}
	employeeRepo := employee.NewRepository(func(context.Context) (*gorm.DB, error) { return db, nil })
	if err := employeeRepo.CreateEmployment(ctx(), emp); err != nil {
		t.Fatalf("CreateEmployment failed: %v", err)
	}
}

// =========================================================================
// Employee Movement Service Tests
// =========================================================================

func TestService_CreateMovement_Success(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	employeeID := uuidStr()
	toPosition := uuidStr()
	req := CreateMovementRequest{
		EmployeeID:           employeeID,
		MovementType:         "promotion",
		DecisionLetterNumber: "SK-001",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-01",
		ToPositionID:         &toPosition,
		Reason:               strPtr("Kinerja baik"),
	}

	resp, err := svc.CreateMovement(ctx(), req)
	if err != nil {
		t.Fatalf("CreateMovement failed: %v", err)
	}

	if resp.MovementType != "promotion" {
		t.Errorf("expected movement_type 'promotion', got '%s'", resp.MovementType)
	}
	if resp.Status != "draft" {
		t.Errorf("expected default status 'draft', got '%s'", resp.Status)
	}
	if resp.EmployeeID != employeeID {
		t.Errorf("expected employee_id '%s', got '%s'", employeeID, resp.EmployeeID)
	}
}

// =========================================================================
// G-7: Business validation per movement type
// =========================================================================

// newValidPromotionReq returns a promotion request that passes G-7 validation.
func newValidPromotionReq() CreateMovementRequest {
	toPos := uuidStr()
	return CreateMovementRequest{
		EmployeeID:           uuidStr(),
		MovementType:         "promotion",
		ToPositionID:         &toPos,
		DecisionLetterNumber: "SK-G7",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-01",
	}
}

// TestService_CreateMovement_Validation_Mutation requires to_org or to_position.
func TestService_CreateMovement_Validation_Mutation(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	req := newValidPromotionReq()
	req.MovementType = "mutation"
	req.ToPositionID = nil
	_, err := svc.CreateMovement(ctx(), req)
	if err == nil {
		t.Fatal("expected error: mutation requires to_organization_id or to_position_id")
	}
	var ve *MovementValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected MovementValidationError, got %T: %v", err, err)
	}

	// Providing to_organization_id satisfies it.
	toOrg := uuidStr()
	req.ToOrganizationID = &toOrg
	if _, err := svc.CreateMovement(ctx(), req); err != nil {
		t.Fatalf("mutation with to_organization_id should pass: %v", err)
	}
}

// TestService_CreateMovement_Validation_Promotion requires to_position_id.
func TestService_CreateMovement_Validation_Promotion(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	req := newValidPromotionReq()
	req.ToPositionID = nil
	_, err := svc.CreateMovement(ctx(), req)
	if err == nil {
		t.Fatal("expected error: promotion requires to_position_id")
	}
}

// TestService_CreateMovement_Validation_StatusChange requires to_employment_status_id.
func TestService_CreateMovement_Validation_StatusChange(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	req := newValidPromotionReq()
	req.MovementType = "status_change"
	_, err := svc.CreateMovement(ctx(), req)
	if err == nil {
		t.Fatal("expected error: status_change requires to_employment_status_id")
	}

	toStatus := uuidStr()
	req.ToEmploymentStatusID = &toStatus
	if _, err := svc.CreateMovement(ctx(), req); err != nil {
		t.Fatalf("status_change with to_employment_status_id should pass: %v", err)
	}
}

// TestService_CreateMovement_Validation_ContractExtension requires an active
// contract for the employee (plan G-7).
func TestService_CreateMovement_Validation_ContractExtension(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	req := newValidPromotionReq()
	req.EmployeeID = employeeID.String()
	req.MovementType = "contract_extension"

	// No active contract yet → validation error.
	_, err := svc.CreateMovement(ctx(), req)
	if err == nil {
		t.Fatal("expected error: contract_extension requires an active contract")
	}

	// Seed an active contract → passes.
	createTestContract(repo, employeeID)
	if _, err := svc.CreateMovement(ctx(), req); err != nil {
		t.Fatalf("contract_extension with active contract should pass: %v", err)
	}
}

// TestService_CreateMovement_Validation_Offboarding passes without to_* fields.
func TestService_CreateMovement_Validation_Offboarding(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	req := newValidPromotionReq()
	req.MovementType = "offboarding"
	req.ToPositionID = nil
	if _, err := svc.CreateMovement(ctx(), req); err != nil {
		t.Fatalf("offboarding without to_* should pass: %v", err)
	}
}

// TestService_UpdateMovement_Validation re-validates after type change.
func TestService_UpdateMovement_Validation(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	created := createTestMovement(repo, uuid.New()) // type = other, no to_*

	// Change to promotion without to_position_id → rejected.
	promotion := "promotion"
	_, err := svc.UpdateMovement(ctx(), created.ID.String(), UpdateMovementRequest{
		MovementType: &promotion,
	})
	if err == nil {
		t.Fatal("expected error: promotion requires to_position_id on update")
	}

	// With to_position_id → passes.
	toPos := uuidStr()
	if _, err := svc.UpdateMovement(ctx(), created.ID.String(), UpdateMovementRequest{
		MovementType: &promotion,
		ToPositionID: &toPos,
	}); err != nil {
		t.Fatalf("promotion with to_position_id should pass on update: %v", err)
	}
}

func TestService_CreateMovement_InvalidUUID(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	req := CreateMovementRequest{
		EmployeeID:           "not-a-uuid",
		MovementType:         "promotion",
		DecisionLetterNumber: "SK-001",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-01",
	}

	_, err := svc.CreateMovement(ctx(), req)
	if err == nil {
		t.Fatal("expected error for invalid employee UUID")
	}
}

func TestService_GetMovementByID_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)

	found, err := svc.GetMovementByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetMovementByID failed: %v", err)
	}

	if found.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), found.ID)
	}
}

func TestService_GetMovementByID_InvalidUUID(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetMovementByID(ctx(), "not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_GetMovementByID_NotFound(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetMovementByID(ctx(), uuidStr())
	if err == nil {
		t.Fatal("expected error for non-existent movement")
	}
}

func TestService_ListMovements_DefaultPagination(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	for i := 0; i < 3; i++ {
		createTestMovement(repo, employeeID)
	}

	resp, err := svc.ListMovements(ctx(), 0, 0, "", "", "")
	if err != nil {
		t.Fatalf("ListMovements failed: %v", err)
	}

	if resp.Page != 1 {
		t.Errorf("expected page 1, got %d", resp.Page)
	}
	if resp.PerPage != 20 {
		t.Errorf("expected per_page 20 (default), got %d", resp.PerPage)
	}
	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}
}

func TestService_ListMovementsByEmployee(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	emp1 := uuid.New()
	emp2 := uuid.New()

	createTestMovement(repo, emp1)
	createTestMovement(repo, emp1)
	createTestMovement(repo, emp2)

	resp, err := svc.ListMovementsByEmployee(ctx(), emp1.String(), 1, 10)
	if err != nil {
		t.Fatalf("ListMovementsByEmployee failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2 for emp1, got %d", resp.Total)
	}
}

func TestService_UpdateMovement_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)

	newReason := "Updated reason"
	toPosition := uuidStr()
	updated, err := svc.UpdateMovement(ctx(), created.ID.String(), UpdateMovementRequest{
		Reason:       &newReason,
		ToPositionID: &toPosition,
	})
	if err != nil {
		t.Fatalf("UpdateMovement failed: %v", err)
	}

	if updated.Reason == nil || *updated.Reason != "Updated reason" {
		t.Errorf("expected reason 'Updated reason', got '%v'", updated.Reason)
	}
}

func TestService_UpdateMovement_NonDraft_Error(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusExecuted
	repo.UpdateMovement(ctx(), created)

	newReason := "Should fail"
	_, err := svc.UpdateMovement(ctx(), created.ID.String(), UpdateMovementRequest{
		Reason: &newReason,
	})
	if err == nil {
		t.Fatal("expected error when updating non-draft movement")
	}
}

func TestService_DeleteMovement_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)

	if err := svc.DeleteMovement(ctx(), created.ID.String()); err != nil {
		t.Fatalf("DeleteMovement failed: %v", err)
	}

	_, err := svc.GetMovementByID(ctx(), created.ID.String())
	if err == nil {
		t.Fatal("expected error after deleting movement")
	}
}

func TestService_DeleteMovement_NonDraft_Error(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusExecuted
	repo.UpdateMovement(ctx(), created)

	err := svc.DeleteMovement(ctx(), created.ID.String())
	if err == nil {
		t.Fatal("expected error when deleting non-draft movement")
	}
}

func TestService_ExecuteMovement_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	// Promotion movements need a career executor for the HR data change.
	svc.SetCareerExecutor(&fakeCareerExecutor{})

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusApproved
	repo.UpdateMovement(ctx(), created)

	executorID := uuidStr()
	if err := svc.ExecuteMovement(ctx(), created.ID.String(), executorID); err != nil {
		t.Fatalf("ExecuteMovement failed: %v", err)
	}

	m, _ := repo.FindMovementByID(ctx(), created.ID)
	if m.Status != MovementStatusExecuted {
		t.Errorf("expected status 'executed', got '%s'", m.Status)
	}
}

func TestService_ExecuteMovement_NoExecutor_Error(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusApproved
	repo.UpdateMovement(ctx(), created)

	err := svc.ExecuteMovement(ctx(), created.ID.String(), uuidStr())
	if err == nil {
		t.Fatal("expected error when career executor is not configured")
	}

	// Movement stays approved (no partial execution).
	m, _ := repo.FindMovementByID(ctx(), created.ID)
	if m.Status != MovementStatusApproved {
		t.Errorf("expected status to remain 'approved', got '%s'", m.Status)
	}
}

func TestService_ExecuteMovement_Promotion_CreatesEmployment(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	executor := &fakeCareerExecutor{
		currentEmployment: &CareerEmployment{
			ID:                   uuid.New(),
			OrganizationID:       ptrUUID(uuid.New()),
			DecisionLetterNumber: "SK-OLD",
			DecisionLetterDate:   "2026-01-01",
			EffectiveDate:        "2026-01-01",
		},
	}
	svc.SetCareerExecutor(executor)

	employeeID := uuid.New()
	toOrg := uuid.New()
	toPos := uuid.New()
	toStatus := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusApproved
	created.MovementType = MovementTypePromotion
	created.ToOrganizationID = &toOrg
	created.ToPositionID = &toPos
	created.ToEmploymentStatusID = &toStatus
	created.EffectiveDate = "2026-08-01"
	repo.UpdateMovement(ctx(), created)

	if err := svc.ExecuteMovement(ctx(), created.ID.String(), uuidStr()); err != nil {
		t.Fatalf("ExecuteMovement failed: %v", err)
	}

	// Previous employment closed one day before effective date.
	if executor.closedID == nil {
		t.Fatal("expected previous employment to be closed")
	}
	if executor.closedEndDate != "2026-07-31" {
		t.Errorf("expected closed end date 2026-07-31, got '%s'", executor.closedEndDate)
	}

	// New employment created with the movement's to_* fields.
	if executor.createdData == nil {
		t.Fatal("expected new employment to be created")
	}
	if executor.createdData.OrganizationID == nil || *executor.createdData.OrganizationID != toOrg {
		t.Errorf("expected organization_id %s, got %v", toOrg, executor.createdData.OrganizationID)
	}
	if executor.createdData.PositionID == nil || *executor.createdData.PositionID != toPos {
		t.Errorf("expected position_id %s, got %v", toPos, executor.createdData.PositionID)
	}
	// SQLite returns DATE columns as RFC3339 timestamps; normalize before compare.
	if len(executor.createdData.EffectiveDate) >= 10 && executor.createdData.EffectiveDate[:10] != "2026-08-01" {
		t.Errorf("expected effective_date 2026-08-01, got '%s'", executor.createdData.EffectiveDate)
	}
	if executor.createdData.DecisionLetterNumber != created.DecisionLetterNumber {
		t.Errorf("expected decision_letter_number from movement, got '%s'", executor.createdData.DecisionLetterNumber)
	}

	// to_employment_id persisted on the movement.
	m, _ := repo.FindMovementByID(ctx(), created.ID)
	if m.Status != MovementStatusExecuted {
		t.Errorf("expected status 'executed', got '%s'", m.Status)
	}
	if m.ToEmploymentID == nil || *m.ToEmploymentID != executor.createdID {
		t.Errorf("expected to_employment_id %s on movement, got %v", executor.createdID, m.ToEmploymentID)
	}

	// Employee NOT deactivated for a promotion.
	if executor.inactiveEmployee != nil {
		t.Errorf("promotion should not deactivate the employee, got %v", *executor.inactiveEmployee)
	}
}

func TestService_ExecuteMovement_Offboarding_DeactivatesEmployee(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	executor := &fakeCareerExecutor{
		currentEmployment: &CareerEmployment{
			ID:                   uuid.New(),
			DecisionLetterNumber: "SK-OLD",
			DecisionLetterDate:   "2026-01-01",
			EffectiveDate:        "2026-01-01",
		},
	}
	svc.SetCareerExecutor(executor)

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusApproved
	created.MovementType = MovementTypeOffboarding
	created.EffectiveDate = "2026-08-01"
	repo.UpdateMovement(ctx(), created)

	if err := svc.ExecuteMovement(ctx(), created.ID.String(), uuidStr()); err != nil {
		t.Fatalf("ExecuteMovement failed: %v", err)
	}

	// Previous employment closed, no new employment created.
	if executor.closedID == nil {
		t.Fatal("expected previous employment to be closed")
	}
	if executor.createdData != nil {
		t.Error("offboarding should NOT create a new employment")
	}

	// Employee marked inactive.
	if executor.inactiveEmployee == nil || *executor.inactiveEmployee != employeeID {
		t.Errorf("expected employee %s to be marked inactive, got %v", employeeID, executor.inactiveEmployee)
	}

	m, _ := repo.FindMovementByID(ctx(), created.ID)
	if m.Status != MovementStatusExecuted {
		t.Errorf("expected status 'executed', got '%s'", m.Status)
	}
}

func TestService_ExecuteMovement_ContractExtension_SkipsEmployment(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	svc.SetCareerExecutor(&fakeCareerExecutor{})

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusApproved
	created.MovementType = MovementTypeContractExtension
	repo.UpdateMovement(ctx(), created)

	if err := svc.ExecuteMovement(ctx(), created.ID.String(), uuidStr()); err != nil {
		t.Fatalf("ExecuteMovement failed: %v", err)
	}

	m, _ := repo.FindMovementByID(ctx(), created.ID)
	if m.Status != MovementStatusExecuted {
		t.Errorf("expected status 'executed', got '%s'", m.Status)
	}
}

// =========================================================================
// Enhancement plan §12.2-§12.4 — transactional execution & conflict detection
// =========================================================================

// TestService_CreateMovement_PositionConflict_Rejected verifies the draft-time
// position conflict check (plan §12.3): creating a movement into a position
// currently held by another employee's open employment is rejected.
func TestService_CreateMovement_PositionConflict_Rejected(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	positionID := uuid.New()
	otherEmp := uuid.New()
	seedEmployment(t, repo, &employee.Employment{
		EmployeeID:           &otherEmp,
		PositionID:           &positionID,
		DecisionLetterNumber: "SK-OTHER",
		DecisionLetterDate:   "2026-01-01",
		EffectiveDate:        "2026-01-01",
	})

	toPos := positionID.String()
	req := CreateMovementRequest{
		EmployeeID:           uuid.New().String(),
		MovementType:         "promotion",
		ToPositionID:         &toPos,
		DecisionLetterNumber: "SK-CONFLICT",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-01",
	}
	_, err := svc.CreateMovement(ctx(), req)
	if err == nil {
		t.Fatal("expected conflict error when target position is occupied")
	}
	var ce *MovementConflictError
	if !errors.As(err, &ce) {
		t.Fatalf("expected MovementConflictError, got %T: %v", err, err)
	}
}

// TestService_CreateMovement_PositionConflict_PassesAfterRelieved verifies the
// same create succeeds when the position is freed before the effective date
// (employment already closed earlier).
func TestService_CreateMovement_PositionConflict_PassesAfterRelieved(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	positionID := uuid.New()
	otherEmp := uuid.New()
	end := "2026-07-31"
	seedEmployment(t, repo, &employee.Employment{
		EmployeeID:           &otherEmp,
		PositionID:           &positionID,
		DecisionLetterNumber: "SK-OTHER",
		DecisionLetterDate:   "2026-01-01",
		EffectiveDate:        "2026-01-01",
		EffectiveEndDate:     &end,
	})

	toPos := positionID.String()
	req := CreateMovementRequest{
		EmployeeID:           uuid.New().String(),
		MovementType:         "promotion",
		ToPositionID:         &toPos,
		DecisionLetterNumber: "SK-CONFLICT",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-01",
	}
	if _, err := svc.CreateMovement(ctx(), req); err != nil {
		t.Fatalf("movement into position freed before effective date should pass: %v", err)
	}
}

// TestService_ExecuteMovement_PositionConflict_Rejected verifies the atomic
// execute-time position conflict check (plan §12.3): execution is rejected and
// the movement stays approved when another employee occupies the target
// position at the effective date.
func TestService_ExecuteMovement_PositionConflict_Rejected(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	positionID := uuid.New()
	otherEmp := uuid.New()
	seedEmployment(t, repo, &employee.Employment{
		EmployeeID:           &otherEmp,
		PositionID:           &positionID,
		DecisionLetterNumber: "SK-OTHER",
		DecisionLetterDate:   "2026-01-01",
		EffectiveDate:        "2026-01-01",
	})

	executor := &fakeCareerExecutor{currentEmployment: &CareerEmployment{ID: uuid.New()}}
	svc.SetCareerExecutor(executor)

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusApproved
	created.MovementType = MovementTypePromotion
	created.ToPositionID = &positionID
	created.EffectiveDate = "2026-08-01"
	repo.UpdateMovement(ctx(), created)

	err := svc.ExecuteMovement(ctx(), created.ID.String(), uuidStr())
	if err == nil {
		t.Fatal("expected conflict error when target position is occupied")
	}
	var ce *MovementConflictError
	if !errors.As(err, &ce) {
		t.Fatalf("expected MovementConflictError, got %T: %v", err, err)
	}

	// Atomic: no HR data change applied, movement still approved.
	if executor.closedID != nil || executor.createdData != nil {
		t.Error("conflict must prevent any employment change")
	}
	m, _ := repo.FindMovementByID(ctx(), created.ID)
	if m.Status != MovementStatusApproved {
		t.Errorf("expected movement to stay approved, got '%s'", m.Status)
	}
}

// TestService_ExecuteMovement_EffectiveDateConflict_Rejected verifies plan
// §12.4: executing a movement whose effective date overlaps a future-dated
// open employment of the same employee is rejected atomically.
func TestService_ExecuteMovement_EffectiveDateConflict_Rejected(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	// Employee already has an open (future-dated) employment starting after
	// the movement's effective date → the new period would overlap.
	seedEmployment(t, repo, &employee.Employment{
		EmployeeID:           &employeeID,
		DecisionLetterNumber: "SK-FUTURE",
		DecisionLetterDate:   "2026-06-01",
		EffectiveDate:        "2026-09-01",
	})

	executor := &fakeCareerExecutor{currentEmployment: &CareerEmployment{ID: uuid.New()}}
	svc.SetCareerExecutor(executor)

	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusApproved
	created.MovementType = MovementTypePromotion
	created.ToPositionID = ptrUUID(uuid.New())
	created.EffectiveDate = "2026-08-15"
	repo.UpdateMovement(ctx(), created)

	err := svc.ExecuteMovement(ctx(), created.ID.String(), uuidStr())
	if err == nil {
		t.Fatal("expected conflict error when effective date overlaps existing employment")
	}
	var ce *MovementConflictError
	if !errors.As(err, &ce) {
		t.Fatalf("expected MovementConflictError, got %T: %v", err, err)
	}
	if executor.closedID != nil || executor.createdData != nil {
		t.Error("conflict must prevent any employment change")
	}
	m, _ := repo.FindMovementByID(ctx(), created.ID)
	if m.Status != MovementStatusApproved {
		t.Errorf("expected movement to stay approved, got '%s'", m.Status)
	}
}

// TestService_ExecuteMovement_BackdatedEffectiveDate_Rejected verifies that a
// movement effective on/before the current employment's start is rejected
// (would create a backdated/overlapping period — plan §12.4).
func TestService_ExecuteMovement_BackdatedEffectiveDate_Rejected(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	seedEmployment(t, repo, &employee.Employment{
		EmployeeID:           &employeeID,
		DecisionLetterNumber: "SK-CURRENT",
		DecisionLetterDate:   "2026-01-01",
		EffectiveDate:        "2026-08-01",
	})

	executor := &fakeCareerExecutor{currentEmployment: &CareerEmployment{ID: uuid.New()}}
	svc.SetCareerExecutor(executor)

	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusApproved
	created.MovementType = MovementTypePromotion
	created.ToPositionID = ptrUUID(uuid.New())
	created.EffectiveDate = "2026-07-01"
	repo.UpdateMovement(ctx(), created)

	err := svc.ExecuteMovement(ctx(), created.ID.String(), uuidStr())
	if err == nil {
		t.Fatal("expected conflict error for backdated effective date")
	}
	var ce *MovementConflictError
	if !errors.As(err, &ce) {
		t.Fatalf("expected MovementConflictError, got %T: %v", err, err)
	}
	m, _ := repo.FindMovementByID(ctx(), created.ID)
	if m.Status != MovementStatusApproved {
		t.Errorf("expected movement to stay approved, got '%s'", m.Status)
	}
}

// TestService_ExecuteMovement_Offboarding_FutureEmployment_Rejected verifies
// the effective-date conflict also guards offboarding/retirement: closing a
// future-dated open employment at a date before its own start would create an
// invalid period, so execution is rejected atomically (plan §12.4).
func TestService_ExecuteMovement_Offboarding_FutureEmployment_Rejected(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	// Employee has a future-dated open employment (from an earlier executed
	// movement) starting after the offboarding's effective date.
	seedEmployment(t, repo, &employee.Employment{
		EmployeeID:           &employeeID,
		DecisionLetterNumber: "SK-FUTURE",
		DecisionLetterDate:   "2026-06-01",
		EffectiveDate:        "2026-09-01",
	})

	executor := &fakeCareerExecutor{currentEmployment: &CareerEmployment{ID: uuid.New()}}
	svc.SetCareerExecutor(executor)

	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusApproved
	created.MovementType = MovementTypeOffboarding
	created.EffectiveDate = "2026-08-15"
	repo.UpdateMovement(ctx(), created)

	err := svc.ExecuteMovement(ctx(), created.ID.String(), uuidStr())
	if err == nil {
		t.Fatal("expected conflict error when offboarding collides with a future-dated employment")
	}
	var ce *MovementConflictError
	if !errors.As(err, &ce) {
		t.Fatalf("expected MovementConflictError, got %T: %v", err, err)
	}
	if executor.closedID != nil || executor.inactiveEmployee != nil {
		t.Error("conflict must prevent any employment/employee change")
	}
	m, _ := repo.FindMovementByID(ctx(), created.ID)
	if m.Status != MovementStatusApproved {
		t.Errorf("expected movement to stay approved, got '%s'", m.Status)
	}
}

// TestService_ExecuteMovement_Rollback_OnCreateFailure verifies plan §12.2:
// when creating the new employment fails mid-execution, the transaction rolls
// back — the movement stays approved (retryable by HR) instead of being left
// in a partially-executed state.
func TestService_ExecuteMovement_Rollback_OnCreateFailure(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	executor := &fakeCareerExecutor{
		currentEmployment: &CareerEmployment{ID: uuid.New()},
		createErr:         errors.New("simulated employment insert failure"),
	}
	svc.SetCareerExecutor(executor)

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusApproved
	created.MovementType = MovementTypePromotion
	created.ToPositionID = ptrUUID(uuid.New())
	created.EffectiveDate = "2026-08-01"
	repo.UpdateMovement(ctx(), created)

	err := svc.ExecuteMovement(ctx(), created.ID.String(), uuidStr())
	if err == nil {
		t.Fatal("expected error when employment creation fails")
	}

	m, _ := repo.FindMovementByID(ctx(), created.ID)
	if m.Status != MovementStatusApproved {
		t.Errorf("expected movement to stay approved after rollback, got '%s'", m.Status)
	}
}

func TestDayBefore(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"2026-08-01", "2026-07-31"},
		{"2026-03-01", "2026-02-28"},
		{"2026-01-01", "2025-12-31"},
	}
	for _, c := range cases {
		got, err := dayBefore(c.in)
		if err != nil {
			t.Fatalf("dayBefore(%q) unexpected error: %v", c.in, err)
		}
		if got != c.want {
			t.Errorf("dayBefore(%q) = %q, want %q", c.in, got, c.want)
		}
	}

	if _, err := dayBefore("not-a-date"); err == nil {
		t.Error("expected error for invalid date")
	}
}

func TestService_ExecuteMovement_Draft_Error(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	// Cannot execute a draft movement (must approve first)
	_, repo, cleanup2 := newTestService()
	defer cleanup2()
	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	// Keep as draft

	err := svc.ExecuteMovement(ctx(), created.ID.String(), uuidStr())
	if err == nil {
		t.Fatal("expected error when executing draft movement")
	}
}

func TestService_CancelMovement_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)

	resp, err := svc.CancelMovement(ctx(), created.ID.String(), CancelMovementRequest{})
	if err != nil {
		t.Fatalf("CancelMovement failed: %v", err)
	}

	m, _ := repo.FindMovementByID(ctx(), created.ID)
	if m.Status != MovementStatusCancelled {
		t.Errorf("expected status 'cancelled', got '%s'", m.Status)
	}
	if resp.Status != string(MovementStatusCancelled) {
		t.Errorf("expected response status 'cancelled', got '%s'", resp.Status)
	}
}

func TestService_CancelMovement_Executed_Error(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusExecuted
	repo.UpdateMovement(ctx(), created)

	_, err := svc.CancelMovement(ctx(), created.ID.String(), CancelMovementRequest{})
	if err == nil {
		t.Fatal("expected error when cancelling executed movement")
	}
}

// TestService_CancelMovement_Approved_CreatesCancellationInstance verifies
// plan §12.16: cancelling an APPROVED movement routes through the Central
// Approval module — a new instance for module employeemovement_cancellation is
// created and the movement enters cancellation_pending (NOT cancelled yet).
func TestService_CancelMovement_Approved_CreatesCancellationInstance(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{resolvedFlowID: uuid.New().String()}
	svc.SetApprovalEngine(fake)

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusApproved
	repo.UpdateMovement(ctx(), created)

	resp, err := svc.CancelMovement(ctx(), created.ID.String(), CancelMovementRequest{})
	if err != nil {
		t.Fatalf("CancelMovement (approved) failed: %v", err)
	}

	if resp.Status != string(MovementStatusCancellationPending) {
		t.Errorf("expected status 'cancellation_pending', got '%s'", resp.Status)
	}
	if resp.CancellationApprovalInstanceID == nil || *resp.CancellationApprovalInstanceID != fake.instanceID {
		t.Errorf("expected cancellation_approval_instance_id %s, got %v", fake.instanceID, resp.CancellationApprovalInstanceID)
	}
	if len(fake.createCalls) != 1 {
		t.Fatalf("expected 1 CreateApprovalInstance call, got %d", len(fake.createCalls))
	}
	if fake.createCalls[0].module != "employeemovement_cancellation" {
		t.Errorf("expected module 'employeemovement_cancellation', got '%s'", fake.createCalls[0].module)
	}
	if fake.createCalls[0].documentID != created.ID.String() {
		t.Errorf("expected document_id %s, got '%s'", created.ID.String(), fake.createCalls[0].documentID)
	}

	// Persisted: status + instance id on the row.
	m, _ := repo.FindMovementByID(ctx(), created.ID)
	if m.Status != MovementStatusCancellationPending {
		t.Errorf("expected persisted status 'cancellation_pending', got '%s'", m.Status)
	}
	if m.CancellationApprovalInstanceID == nil || *m.CancellationApprovalInstanceID != uuid.MustParse(fake.instanceID) {
		t.Errorf("expected persisted cancellation_approval_instance_id %s, got %v", fake.instanceID, m.CancellationApprovalInstanceID)
	}
}

// TestService_CancelMovement_Approved_NoEngine_Error verifies the cancellation
// request requires the approval engine to be wired.
func TestService_CancelMovement_Approved_NoEngine_Error(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusApproved
	repo.UpdateMovement(ctx(), created)

	_, err := svc.CancelMovement(ctx(), created.ID.String(), CancelMovementRequest{})
	if err == nil {
		t.Fatal("expected error when approval engine is not configured")
	}
}

// TestService_CancelMovement_Approved_NoFlow_Error verifies the cancellation
// request fails when no flow can be resolved for employeemovement_cancellation.
func TestService_CancelMovement_Approved_NoFlow_Error(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	svc.SetApprovalEngine(&fakeApprovalEngine{}) // no resolvedFlowID

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusApproved
	repo.UpdateMovement(ctx(), created)

	_, err := svc.CancelMovement(ctx(), created.ID.String(), CancelMovementRequest{})
	if err == nil {
		t.Fatal("expected error when no cancellation flow is configured")
	}
}

// TestService_HandleCancellationStatusChange_Approved verifies that when the
// cancellation approval instance resolves APPROVED, the movement is cancelled
// (plan §12.16).
func TestService_HandleCancellationStatusChange_Approved(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusCancellationPending
	repo.UpdateMovement(ctx(), created)

	if err := svc.HandleCancellationStatusChange(ctx(), created.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleCancellationStatusChange failed: %v", err)
	}

	m, _ := repo.FindMovementByID(ctx(), created.ID)
	if m.Status != MovementStatusCancelled {
		t.Errorf("expected status 'cancelled' after cancellation approved, got '%s'", m.Status)
	}
}

// TestService_HandleCancellationStatusChange_Rejected verifies a REJECTED
// cancellation request returns the movement to approved (plan §12.16).
func TestService_HandleCancellationStatusChange_Rejected(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusCancellationPending
	repo.UpdateMovement(ctx(), created)

	if err := svc.HandleCancellationStatusChange(ctx(), created.ID, "REJECTED", "Pembatalan ditolak"); err != nil {
		t.Fatalf("HandleCancellationStatusChange failed: %v", err)
	}

	m, _ := repo.FindMovementByID(ctx(), created.ID)
	if m.Status != MovementStatusApproved {
		t.Errorf("expected status back to 'approved' after cancellation rejected, got '%s'", m.Status)
	}
	if m.Notes == nil || *m.Notes != "Pembatalan ditolak" {
		t.Errorf("expected note from rejection, got %v", m.Notes)
	}
}

// TestService_HandleCancellationStatusChange_NotPending_Ignored verifies the
// cancellation callback is a no-op for movements not in cancellation_pending.
func TestService_HandleCancellationStatusChange_NotPending_Ignored(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)
	created.Status = MovementStatusApproved
	repo.UpdateMovement(ctx(), created)

	if err := svc.HandleCancellationStatusChange(ctx(), created.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleCancellationStatusChange should be a no-op: %v", err)
	}

	m, _ := repo.FindMovementByID(ctx(), created.ID)
	if m.Status != MovementStatusApproved {
		t.Errorf("expected status to remain 'approved', got '%s'", m.Status)
	}
}

// =========================================================================
// Employee Contract Service Tests
// =========================================================================

func TestService_CreateContract_Success(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	employeeID := uuidStr()
	endDate := "2026-12-31"
	req := CreateContractRequest{
		EmployeeID:     employeeID,
		ContractNumber: "CTR-001",
		ContractType:   "pkwt",
		StartDate:      "2026-01-01",
		EndDate:        &endDate,
	}

	resp, err := svc.CreateContract(ctx(), req)
	if err != nil {
		t.Fatalf("CreateContract failed: %v", err)
	}

	if resp.ContractNumber != "CTR-001" {
		t.Errorf("expected contract_number 'CTR-001', got '%s'", resp.ContractNumber)
	}
	if resp.Status != "active" {
		t.Errorf("expected default status 'active', got '%s'", resp.Status)
	}
}

// TestService_CreateContract_WithPrevious_ChainCount verifies G-6 at the
// service level: creating a contract with previous_contract_id routes through
// ExtendContract and the response carries the derived extension_count.
func TestService_CreateContract_WithPrevious_ChainCount(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	first := createTestContract(repo, employeeID) // 1st extension → count 1
	prevID := first.ID.String()
	second := CreateContractRequest{
		EmployeeID:         employeeID.String(),
		ContractNumber:     "CTR-002",
		ContractType:       "pkwt",
		StartDate:          "2026-06-01",
		PreviousContractID: &prevID,
	}
	resp1, err := svc.CreateContract(ctx(), second)
	if err != nil {
		t.Fatalf("CreateContract (1st ext) failed: %v", err)
	}
	if resp1.ExtensionCount != 1 {
		t.Errorf("expected extension_count 1, got %d", resp1.ExtensionCount)
	}

	// 2nd extension → count 2 (chained, plan G-6)
	third := CreateContractRequest{
		EmployeeID:         employeeID.String(),
		ContractNumber:     "CTR-003",
		ContractType:       "pkwt",
		StartDate:          "2027-01-01",
		PreviousContractID: &resp1.ID,
	}
	resp2, err := svc.CreateContract(ctx(), third)
	if err != nil {
		t.Fatalf("CreateContract (2nd ext) failed: %v", err)
	}
	if resp2.ExtensionCount != 2 {
		t.Errorf("expected extension_count 2, got %d", resp2.ExtensionCount)
	}

	// Previous contract marked extended in DB.
	db, err := svc.repo.getDB(ctx())
	if err != nil {
		t.Fatalf("failed to get test db: %v", err)
	}
	var status string
	if err := db.Model(&EmployeeContract{}).Where("id = ?", resp1.ID).Pluck("status", &status).Error; err != nil {
		t.Fatalf("failed to query contract status: %v", err)
	}
	if status != string(ContractStatusExtended) {
		t.Errorf("expected previous contract status 'extended', got '%s'", status)
	}
}

func TestService_CreateContract_InvalidEmployeeID(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	req := CreateContractRequest{
		EmployeeID:     "not-a-uuid",
		ContractNumber: "CTR-001",
		ContractType:   "pkwt",
		StartDate:      "2026-01-01",
	}

	_, err := svc.CreateContract(ctx(), req)
	if err == nil {
		t.Fatal("expected error for invalid employee UUID")
	}
}

func TestService_GetContractByID_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestContract(repo, employeeID)

	found, err := svc.GetContractByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetContractByID failed: %v", err)
	}

	if found.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), found.ID)
	}
}

func TestService_GetContractByID_NotFound(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetContractByID(ctx(), uuidStr())
	if err == nil {
		t.Fatal("expected error for non-existent contract")
	}
}

func TestService_ListContracts_DefaultPagination(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	for i := 0; i < 3; i++ {
		createTestContract(repo, employeeID)
	}

	resp, err := svc.ListContracts(ctx(), 0, 0, "", "")
	if err != nil {
		t.Fatalf("ListContracts failed: %v", err)
	}

	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}
}

func TestService_ListContractsByEmployee(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	emp1 := uuid.New()
	emp2 := uuid.New()

	createTestContract(repo, emp1)
	createTestContract(repo, emp2)
	createTestContract(repo, emp2)

	resp, err := svc.ListContractsByEmployee(ctx(), emp2.String(), 1, 10)
	if err != nil {
		t.Fatalf("ListContractsByEmployee failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2 for emp2, got %d", resp.Total)
	}
}

func TestService_UpdateContract_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestContract(repo, employeeID)

	newStatus := "expired"
	updated, err := svc.UpdateContract(ctx(), created.ID.String(), UpdateContractRequest{
		Status: &newStatus,
	})
	if err != nil {
		t.Fatalf("UpdateContract failed: %v", err)
	}

	if updated.Status != "expired" {
		t.Errorf("expected status 'expired', got '%s'", updated.Status)
	}
}

func TestService_DeleteContract_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestContract(repo, employeeID)

	if err := svc.DeleteContract(ctx(), created.ID.String()); err != nil {
		t.Fatalf("DeleteContract failed: %v", err)
	}

	_, err := svc.GetContractByID(ctx(), created.ID.String())
	if err == nil {
		t.Fatal("expected error after deleting contract")
	}
}

func TestService_DeleteContract_NotFound(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	err := svc.DeleteContract(ctx(), uuidStr())
	if err == nil {
		t.Fatal("expected error for non-existent contract")
	}
}
