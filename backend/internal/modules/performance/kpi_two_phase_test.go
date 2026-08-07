package performance

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// setupTwoPhaseEvaluation creates a period, template, perspective, indicator,
// and a DRAFT evaluation snapshot via CreateEvaluationWithSnapshot, returning
// the evaluation ID.
func setupTwoPhaseEvaluation(t *testing.T, svc *Service) string {
	t.Helper()
	ctx := context.Background()

	period, err := svc.CreatePerformancePeriod(ctx, CreatePerformancePeriodRequest{
		PeriodCode: "2026-Q1", PeriodType: "QUARTERLY", Year: 2026,
	})
	if err != nil {
		t.Fatalf("failed to create period: %v", err)
	}
	orgID := createTestOrgID()
	tmpl, err := svc.CreatePerformanceTemplate(ctx, CreatePerformanceTemplateRequest{
		OrganizationID: orgID, Name: "Two Phase Template",
	})
	if err != nil {
		t.Fatalf("failed to create template: %v", err)
	}
	persp, err := svc.CreatePerformancePerspective(ctx, CreatePerformancePerspectiveRequest{Name: "Financial"})
	if err != nil {
		t.Fatalf("failed to create perspective: %v", err)
	}
	if _, err := svc.CreatePerformanceIndicator(ctx, CreatePerformanceIndicatorRequest{
		PerformanceTemplateID: tmpl.ID,
		PerspectiveID:         persp.ID,
		IndicatorType:         "MAXIMIZATION",
		Title:                 "Revenue Growth",
		Weight:                100,
	}); err != nil {
		t.Fatalf("failed to create indicator: %v", err)
	}

	eval, err := svc.CreateEvaluationWithSnapshot(ctx, CreateEvaluationWithSnapshotRequest{
		EmployeeID:     uuid.New().String(),
		OrganizationID: orgID,
		PeriodID:       period.ID,
		TemplateID:     tmpl.ID,
	})
	if err != nil {
		t.Fatalf("failed to create evaluation snapshot: %v", err)
	}
	if eval.Status != "DRAFT" {
		t.Fatalf("expected new evaluation to start DRAFT, got %s", eval.Status)
	}
	if len(eval.Details) != 1 {
		t.Fatalf("expected 1 detail, got %d", len(eval.Details))
	}
	if eval.Details[0].Target != 0 {
		t.Errorf("expected snapshotted target to be 0 (employee-authored), got %v", eval.Details[0].Target)
	}
	return eval.ID
}

func TestService_SubmitTarget_RequiresAllTargetsFilled(t *testing.T) {
	svc, _, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()
	ctx := context.Background()

	evalID := setupTwoPhaseEvaluation(t, svc)

	if _, err := svc.SubmitTarget(ctx, evalID); err == nil {
		t.Fatal("expected error submitting target with unfilled indicator target")
	}
}

func TestService_TwoPhaseWorkflow_FullCycle(t *testing.T) {
	svc, _, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()
	ctx := context.Background()

	evalID := setupTwoPhaseEvaluation(t, svc)

	full, err := svc.GetEvaluationWithDetails(ctx, evalID)
	if err != nil {
		t.Fatalf("GetEvaluationWithDetails failed: %v", err)
	}
	detailID := full.Details[0].ID

	// Fill target (and its unit) while DRAFT.
	unit := "IDR"
	targetResp, err := svc.UpdateEvaluationTarget(ctx, detailID, UpdateEvaluationTargetRequest{Target: 1000, UnitOfMeasurement: &unit})
	if err != nil {
		t.Fatalf("UpdateEvaluationTarget failed: %v", err)
	}
	if targetResp.UnitOfMeasurement != "IDR" {
		t.Errorf("expected unit_of_measurement 'IDR', got %q", targetResp.UnitOfMeasurement)
	}

	// Actual cannot be filled before target is approved.
	if _, err := svc.UpdateEvaluationActual(ctx, detailID, UpdateEvaluationActualRequest{Actual: 500}); err == nil {
		t.Fatal("expected error filling actual before target approved")
	}

	// Submit target: DRAFT -> TARGET_SUBMITTED.
	submitted, err := svc.SubmitTarget(ctx, evalID)
	if err != nil {
		t.Fatalf("SubmitTarget failed: %v", err)
	}
	if submitted.Status != "TARGET_SUBMITTED" {
		t.Fatalf("expected TARGET_SUBMITTED, got %s", submitted.Status)
	}

	// Target can no longer be edited once submitted.
	if _, err := svc.UpdateEvaluationTarget(ctx, detailID, UpdateEvaluationTargetRequest{Target: 2000}); err == nil {
		t.Fatal("expected error editing target after submission")
	}

	// Reject target -> back to DRAFT, then resubmit.
	rejected, err := svc.RejectTarget(ctx, evalID, nil)
	if err != nil {
		t.Fatalf("RejectTarget failed: %v", err)
	}
	if rejected.Status != "DRAFT" {
		t.Fatalf("expected DRAFT after target rejection, got %s", rejected.Status)
	}
	if _, err := svc.SubmitTarget(ctx, evalID); err != nil {
		t.Fatalf("re-SubmitTarget failed: %v", err)
	}

	// Approve target: TARGET_SUBMITTED -> TARGET_APPROVED.
	approved, err := svc.ApproveTarget(ctx, evalID)
	if err != nil {
		t.Fatalf("ApproveTarget failed: %v", err)
	}
	if approved.Status != "TARGET_APPROVED" {
		t.Fatalf("expected TARGET_APPROVED, got %s", approved.Status)
	}

	// Now actual can be filled.
	if _, err := svc.UpdateEvaluationActual(ctx, detailID, UpdateEvaluationActualRequest{Actual: 1000}); err != nil {
		t.Fatalf("UpdateEvaluationActual failed: %v", err)
	}

	// Submit realization: TARGET_APPROVED -> SUBMITTED.
	realized, err := svc.SubmitEvaluation(ctx, evalID)
	if err != nil {
		t.Fatalf("SubmitEvaluation failed: %v", err)
	}
	if realized.Status != "SUBMITTED" {
		t.Fatalf("expected SUBMITTED, got %s", realized.Status)
	}

	// Reject realization -> TARGET_APPROVED (not all the way to DRAFT).
	rejectedRealization, err := svc.RejectEvaluation(ctx, evalID, nil)
	if err != nil {
		t.Fatalf("RejectEvaluation failed: %v", err)
	}
	if rejectedRealization.Status != "TARGET_APPROVED" {
		t.Fatalf("expected TARGET_APPROVED after realization rejection, got %s", rejectedRealization.Status)
	}

	// Resubmit and approve to COMPLETED.
	if _, err := svc.SubmitEvaluation(ctx, evalID); err != nil {
		t.Fatalf("re-SubmitEvaluation failed: %v", err)
	}
	finalApproved, err := svc.ApproveEvaluation(ctx, evalID)
	if err != nil {
		t.Fatalf("ApproveEvaluation failed: %v", err)
	}
	if finalApproved.Status != "APPROVED" {
		t.Fatalf("expected APPROVED, got %s", finalApproved.Status)
	}
	completed, err := svc.CompleteEvaluation(ctx, evalID)
	if err != nil {
		t.Fatalf("CompleteEvaluation failed: %v", err)
	}
	if completed.Status != "COMPLETED" {
		t.Fatalf("expected COMPLETED, got %s", completed.Status)
	}
}

func TestService_ProgramItems_CRUDAndStatusGating(t *testing.T) {
	svc, _, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()
	ctx := context.Background()

	evalID := setupTwoPhaseEvaluation(t, svc)

	// Create while DRAFT.
	unit := "unit"
	item, err := svc.CreateProgramItem(ctx, CreateProgramItemRequest{
		PerformanceEvaluationID: evalID,
		Title:                   "Improve onboarding flow",
		Weight:                  40,
		UnitOfMeasurement:       &unit,
		Target:                  10,
	})
	if err != nil {
		t.Fatalf("CreateProgramItem failed: %v", err)
	}
	if item.FormulaType != "MANUAL" {
		t.Errorf("expected default formula_type MANUAL, got %s", item.FormulaType)
	}
	if item.Weight != 40 {
		t.Errorf("expected weight 40, got %v", item.Weight)
	}
	if item.UnitOfMeasurement != "unit" {
		t.Errorf("expected unit_of_measurement 'unit', got %q", item.UnitOfMeasurement)
	}

	// Edit target while DRAFT.
	updated, err := svc.UpdateProgramItemTarget(ctx, item.ID, UpdateProgramItemTargetRequest{Target: float64Ptr(20)})
	if err != nil {
		t.Fatalf("UpdateProgramItemTarget failed: %v", err)
	}
	if updated.Target != 20 {
		t.Errorf("expected target 20, got %v", updated.Target)
	}

	// Fill the indicator target too so SubmitTarget's overall guard passes.
	full, _ := svc.GetEvaluationWithDetails(ctx, evalID)
	if _, err := svc.UpdateEvaluationTarget(ctx, full.Details[0].ID, UpdateEvaluationTargetRequest{Target: 1000}); err != nil {
		t.Fatalf("UpdateEvaluationTarget failed: %v", err)
	}

	if _, err := svc.SubmitTarget(ctx, evalID); err != nil {
		t.Fatalf("SubmitTarget failed: %v", err)
	}

	// Program item can't be edited/deleted once target is submitted.
	if _, err := svc.UpdateProgramItemTarget(ctx, item.ID, UpdateProgramItemTargetRequest{Target: float64Ptr(30)}); err == nil {
		t.Fatal("expected error editing program item target after target submission")
	}
	if err := svc.DeleteProgramItem(ctx, item.ID); err == nil {
		t.Fatal("expected error deleting program item after target submission")
	}

	if _, err := svc.ApproveTarget(ctx, evalID); err != nil {
		t.Fatalf("ApproveTarget failed: %v", err)
	}

	// Actual can now be filled.
	afterActual, err := svc.UpdateProgramItemActual(ctx, item.ID, UpdateProgramItemActualRequest{Actual: 10})
	if err != nil {
		t.Fatalf("UpdateProgramItemActual failed: %v", err)
	}
	if afterActual.Achievement != 50 {
		t.Errorf("expected achievement 50%% (10/20), got %v", afterActual.Achievement)
	}
	// Score = weight * achievement / 100 = 40 * 50 / 100 = 20.
	if afterActual.Score != 20 {
		t.Errorf("expected weighted score 20 (40%% weight * 50%% achievement), got %v", afterActual.Score)
	}
}

func TestService_ProgramItems_WeightCannotExceed100(t *testing.T) {
	svc, _, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()
	ctx := context.Background()

	evalID := setupTwoPhaseEvaluation(t, svc)

	first, err := svc.CreateProgramItem(ctx, CreateProgramItemRequest{
		PerformanceEvaluationID: evalID,
		Title:                   "First program",
		Weight:                  70,
		Target:                  10,
	})
	if err != nil {
		t.Fatalf("CreateProgramItem failed: %v", err)
	}

	// A second item pushing the total past 100% must be rejected.
	if _, err := svc.CreateProgramItem(ctx, CreateProgramItemRequest{
		PerformanceEvaluationID: evalID,
		Title:                   "Second program",
		Weight:                  40,
		Target:                  10,
	}); err == nil {
		t.Fatal("expected error creating program item that pushes total weight over 100%")
	}

	// Exactly reaching 100% is fine.
	second, err := svc.CreateProgramItem(ctx, CreateProgramItemRequest{
		PerformanceEvaluationID: evalID,
		Title:                   "Second program",
		Weight:                  30,
		Target:                  10,
	})
	if err != nil {
		t.Fatalf("expected CreateProgramItem to succeed at exactly 100%% total: %v", err)
	}

	// Editing an existing item's weight past the cap is also rejected.
	if _, err := svc.UpdateProgramItemTarget(ctx, second.ID, UpdateProgramItemTargetRequest{Weight: float64Ptr(31)}); err == nil {
		t.Fatal("expected error editing weight past the 100% cap")
	}

	// But re-saving the same weight (no-op) or reducing it is fine.
	if _, err := svc.UpdateProgramItemTarget(ctx, second.ID, UpdateProgramItemTargetRequest{Weight: float64Ptr(30)}); err != nil {
		t.Fatalf("expected no-op weight update to succeed: %v", err)
	}
	if _, err := svc.UpdateProgramItemTarget(ctx, first.ID, UpdateProgramItemTargetRequest{Weight: float64Ptr(60)}); err != nil {
		t.Fatalf("expected reducing weight to succeed: %v", err)
	}
}

func TestService_CreateProgramItem_RequiresDraftStatus(t *testing.T) {
	svc, _, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()
	ctx := context.Background()

	evalID := setupTwoPhaseEvaluation(t, svc)
	full, _ := svc.GetEvaluationWithDetails(ctx, evalID)
	if _, err := svc.UpdateEvaluationTarget(ctx, full.Details[0].ID, UpdateEvaluationTargetRequest{Target: 1000}); err != nil {
		t.Fatalf("UpdateEvaluationTarget failed: %v", err)
	}
	if _, err := svc.SubmitTarget(ctx, evalID); err != nil {
		t.Fatalf("SubmitTarget failed: %v", err)
	}

	if _, err := svc.CreateProgramItem(ctx, CreateProgramItemRequest{
		PerformanceEvaluationID: evalID,
		Title:                   "Too late",
		Target:                  5,
	}); err == nil {
		t.Fatal("expected error creating program item outside DRAFT status")
	}
}
