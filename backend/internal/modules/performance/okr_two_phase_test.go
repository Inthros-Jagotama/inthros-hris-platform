package performance

import (
	"context"
	"testing"

	sqlite "github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// setupOKRTestDB creates an in-memory SQLite database, auto-migrates all OKR
// models, and returns a service wired with a dbResolver usable by the
// ctx-based push-callback methods (HandleKeyResultApprovalStatusChange,
// HandleAssessmentApprovalStatusChange).
func setupOKRTestDB(t *testing.T) (OKRService, *gorm.DB, func()) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(
		&OKRTemplate{},
		&OKRObjective{},
		&OKRKeyResult{},
		&OKREvaluation{},
		&OKREvaluationDetail{},
		&PerformancePeriod{},
	); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	// Raw tables owned by other modules (organization/employee), queried
	// directly via db.Table(...) for cascading-objective hierarchy walks —
	// see okr_objective_scope_test.go.
	rawTables := []string{
		`CREATE TABLE IF NOT EXISTS organizations (
			id CHAR(36) PRIMARY KEY,
			parent_id CHAR(36) NULL,
			nomenclature VARCHAR(255) NOT NULL DEFAULT '',
			deleted_at TIMESTAMP NULL
		)`,
		`CREATE TABLE IF NOT EXISTS employments (
			id CHAR(36) PRIMARY KEY,
			employee_id CHAR(36) NOT NULL,
			organization_id CHAR(36) NOT NULL,
			effective_date DATE NOT NULL,
			effective_end_date DATE NULL
		)`,
		`CREATE TABLE IF NOT EXISTS employee_accounts (
			id CHAR(36) PRIMARY KEY,
			employee_id CHAR(36) NOT NULL,
			user_id CHAR(36) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS employees (
			id CHAR(36) PRIMARY KEY,
			name VARCHAR(255) NOT NULL DEFAULT ''
		)`,
	}
	for _, stmt := range rawTables {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("failed to create raw test table: %v", err)
		}
	}

	dbResolver := func(ctx context.Context) (*gorm.DB, error) {
		return db, nil
	}

	repo := NewOKRRepository()
	svc := NewOKRService(repo, dbResolver)

	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	return svc, db, cleanup
}

// seedOKROrganization inserts a minimal organizations row.
func seedOKROrganization(t *testing.T, db *gorm.DB, id uuid.UUID, parentID *uuid.UUID, name string) {
	t.Helper()
	var parent interface{}
	if parentID != nil {
		parent = parentID.String()
	}
	if err := db.Exec("INSERT INTO organizations (id, parent_id, nomenclature) VALUES (?, ?, ?)", id.String(), parent, name).Error; err != nil {
		t.Fatalf("failed to seed organization: %v", err)
	}
}

// seedOKREmployment links an employee to an Organization as their current
// (open-ended) assignment — makes the Organization "occupied".
func seedOKREmployment(t *testing.T, db *gorm.DB, employeeID, orgID uuid.UUID) {
	t.Helper()
	if err := db.Exec(
		"INSERT INTO employments (id, employee_id, organization_id, effective_date, effective_end_date) VALUES (?, ?, ?, ?, NULL)",
		uuid.New().String(), employeeID.String(), orgID.String(), "2026-01-01",
	).Error; err != nil {
		t.Fatalf("failed to seed employment: %v", err)
	}
}

// seedOKREmployeeAccount links an employee to their platform user login.
func seedOKREmployeeAccount(t *testing.T, db *gorm.DB, employeeID, userID uuid.UUID) {
	t.Helper()
	if err := db.Exec(
		"INSERT INTO employee_accounts (id, employee_id, user_id) VALUES (?, ?, ?)",
		uuid.New().String(), employeeID.String(), userID.String(),
	).Error; err != nil {
		t.Fatalf("failed to seed employee account: %v", err)
	}
}

// setupOKRTwoPhaseEvaluation creates a two-level org hierarchy (a root
// creator Organization, occupied by creatorUserID, and one occupied
// subordinate Organization), a Template+Objective created by the creator
// for the subordinate, and a DRAFT evaluation snapshot of it (no Key
// Results yet — those are employee-proposed). Returns the evaluation ID and
// the Objective's ID.
func setupOKRTwoPhaseEvaluation(t *testing.T, svc OKRService, db *gorm.DB) (evalID string, objectiveID string) {
	t.Helper()

	creatorOrgID := uuid.New()
	orgUUID := uuid.New()
	orgID := orgUUID.String()
	seedOKROrganization(t, db, creatorOrgID, nil, "Root Org")
	seedOKROrganization(t, db, orgUUID, &creatorOrgID, "Subordinate Org")

	creatorEmployeeID := uuid.New()
	creatorUserID := uuid.New()
	seedOKREmployment(t, db, creatorEmployeeID, creatorOrgID)
	seedOKREmployeeAccount(t, db, creatorEmployeeID, creatorUserID)

	subEmployeeID := uuid.New()
	seedOKREmployment(t, db, subEmployeeID, orgUUID)

	tmpl, err := svc.CreateTemplate(db, creatorUserID, &CreateOKRTemplateRequest{
		OrganizationID: orgID,
		Name:           "Two Phase OKR Template",
	})
	if err != nil {
		t.Fatalf("failed to create template: %v", err)
	}

	obj, err := svc.CreateObjective(db, &CreateOKRObjectiveRequest{
		TemplateID: tmpl.ID,
		Title:      "Grow Revenue",
		Weight:     100,
	})
	if err != nil {
		t.Fatalf("failed to create objective: %v", err)
	}

	periodID := uuid.New().String()
	eval, err := svc.CreateEvaluationWithSnapshot(db, &CreateOKREvaluationRequest{
		EmployeeID:     uuid.New().String(),
		OrganizationID: orgID,
		PeriodID:       periodID,
		TemplateID:     tmpl.ID,
	})
	if err != nil {
		t.Fatalf("failed to create evaluation snapshot: %v", err)
	}
	if eval.Status != "DRAFT" {
		t.Fatalf("expected DRAFT, got %s", eval.Status)
	}
	if len(eval.Details) != 0 {
		t.Fatalf("expected no Key Results copied from template, got %d", len(eval.Details))
	}

	return eval.ID, obj.ID
}

// proposeKeyResult adds one employee-authored Key Result under the given
// Objective, with weight 100 (fills the whole Objective).
func proposeKeyResult(t *testing.T, svc OKRService, db *gorm.DB, evalID, objectiveID string) {
	t.Helper()
	if _, err := svc.CreateEvaluationKeyResult(db, &CreateOKREvaluationKeyResultRequest{
		EvaluationID:    evalID,
		ObjectiveID:     objectiveID,
		ObjectiveTitle:  "Grow Revenue",
		ObjectiveWeight: 100,
		Title:           "Close 10 new deals",
		TargetType:      "NUMBER",
		TargetValue:     10,
		FormulaType:     "HIGHER_BETTER",
		Weight:          100,
	}); err != nil {
		t.Fatalf("failed to propose key result: %v", err)
	}
}

func TestOKRService_TwoPhaseWorkflow_FullCycle(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()
	ctx := context.Background()

	evalID, objectiveID := setupOKRTwoPhaseEvaluation(t, svc, db)
	proposeKeyResult(t, svc, db, evalID, objectiveID)

	evalUUID, _ := uuid.Parse(evalID)
	userID := uuid.New()

	// DRAFT -> KR_SUBMITTED (no approval engine wired: manual fallback).
	submitted, err := svc.SubmitKeyResults(ctx, db, evalUUID, userID)
	if err != nil {
		t.Fatalf("SubmitKeyResults failed: %v", err)
	}
	if submitted.Status != "KR_SUBMITTED" {
		t.Fatalf("expected KR_SUBMITTED, got %s", submitted.Status)
	}

	// KR_SUBMITTED -> KR_APPROVED (manual fallback).
	approved, err := svc.ApproveKeyResults(db, evalUUID, userID)
	if err != nil {
		t.Fatalf("ApproveKeyResults failed: %v", err)
	}
	if approved.Status != "KR_APPROVED" {
		t.Fatalf("expected KR_APPROVED, got %s", approved.Status)
	}

	// Fill actual on the proposed Key Result.
	full, err := svc.GetEvaluationWithDetails(db, evalUUID)
	if err != nil {
		t.Fatalf("GetEvaluationWithDetails failed: %v", err)
	}
	if len(full.Details) != 1 {
		t.Fatalf("expected 1 key result detail, got %d", len(full.Details))
	}
	detailID, _ := uuid.Parse(full.Details[0].ID)
	if _, err := svc.UpdateEvaluationDetailActual(db, detailID, &UpdateOKREvaluationDetailRequest{ActualValue: 10}); err != nil {
		t.Fatalf("UpdateEvaluationDetailActual failed: %v", err)
	}

	// KR_APPROVED -> SUBMITTED (manual fallback).
	realized, err := svc.SubmitEvaluation(ctx, db, evalUUID, userID)
	if err != nil {
		t.Fatalf("SubmitEvaluation failed: %v", err)
	}
	if realized.Status != "SUBMITTED" {
		t.Fatalf("expected SUBMITTED, got %s", realized.Status)
	}

	// SUBMITTED -> COMPLETED directly (final approval auto-completes, no
	// separate manual "Complete" step needed).
	completed, err := svc.ApproveEvaluation(db, evalUUID, userID)
	if err != nil {
		t.Fatalf("ApproveEvaluation failed: %v", err)
	}
	if completed.Status != "COMPLETED" {
		t.Fatalf("expected COMPLETED, got %s", completed.Status)
	}
}

func TestOKRService_UpdateEvaluationDetailActual_GatedToKRApproved(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	evalID, objectiveID := setupOKRTwoPhaseEvaluation(t, svc, db)
	proposeKeyResult(t, svc, db, evalID, objectiveID)

	full, err := svc.GetEvaluationWithDetails(db, uuidMustParse(t, evalID))
	if err != nil {
		t.Fatalf("GetEvaluationWithDetails failed: %v", err)
	}
	detailID, _ := uuid.Parse(full.Details[0].ID)

	// Still DRAFT — actual editing must be rejected.
	if _, err := svc.UpdateEvaluationDetailActual(db, detailID, &UpdateOKREvaluationDetailRequest{ActualValue: 5}); err == nil {
		t.Fatal("expected UpdateEvaluationDetailActual to fail while evaluation is DRAFT")
	}
}

func TestOKRService_KeyResultWeight_CannotExceed100PerObjective(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	evalID, objectiveID := setupOKRTwoPhaseEvaluation(t, svc, db)
	proposeKeyResult(t, svc, db, evalID, objectiveID)

	// A second key result under the same Objective pushing total weight over 100%.
	_, err := svc.CreateEvaluationKeyResult(db, &CreateOKREvaluationKeyResultRequest{
		EvaluationID:    evalID,
		ObjectiveID:     objectiveID,
		ObjectiveTitle:  "Grow Revenue",
		ObjectiveWeight: 100,
		Title:           "Another key result",
		TargetValue:     5,
		Weight:          1,
	})
	if err == nil {
		t.Fatal("expected weight-exceeds-100 validation error")
	}
}

func uuidMustParse(t *testing.T, s string) uuid.UUID {
	t.Helper()
	id, err := uuid.Parse(s)
	if err != nil {
		t.Fatalf("invalid uuid %q: %v", s, err)
	}
	return id
}

