package performance

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestOKRService_SubmitKeyResults_RoutesThroughApprovalEngine(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()
	ctx := context.Background()

	fake := &fakeApprovalEngine{flowIDByModule: map[string]string{
		ApprovalModuleOKRKeyResult: uuid.New().String(),
	}}
	svc.SetApprovalEngine(fake)

	evalID, objectiveID := setupOKRTwoPhaseEvaluation(t, svc, db)
	proposeKeyResult(t, svc, db, evalID, objectiveID)

	evalUUID, _ := uuid.Parse(evalID)
	resp, err := svc.SubmitKeyResults(ctx, db, evalUUID, uuid.New())
	if err != nil {
		t.Fatalf("SubmitKeyResults failed: %v", err)
	}
	if resp.Status != "KR_SUBMITTED" {
		t.Fatalf("expected KR_SUBMITTED, got %s", resp.Status)
	}
	if resp.KRApprovalInstanceID == "" {
		t.Error("expected kr_approval_instance_id to be set")
	}
	if len(fake.createCalls) != 1 {
		t.Fatalf("expected 1 CreateApprovalInstance call, got %d", len(fake.createCalls))
	}
	if fake.createCalls[0].module != ApprovalModuleOKRKeyResult || fake.createCalls[0].documentID != evalID {
		t.Errorf("unexpected call params: %+v", fake.createCalls[0])
	}
}

// TestOKRService_SubmitKeyResults_FlowConfiguredButResolutionFails_HardFails
// validates that when a flow IS configured but fails to resolve, the whole
// submission is rejected (status stays DRAFT) instead of silently degrading
// to the manual fallback — mirrors the equivalent KPI SubmitTarget test.
func TestOKRService_SubmitKeyResults_FlowConfiguredButResolutionFails_HardFails(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()
	ctx := context.Background()

	fake := &fakeApprovalEngine{
		flowIDByModule: map[string]string{ApprovalModuleOKRKeyResult: uuid.New().String()},
		createErr:      fmt.Errorf("approval step resolved to zero approvers"),
	}
	svc.SetApprovalEngine(fake)

	evalID, objectiveID := setupOKRTwoPhaseEvaluation(t, svc, db)
	proposeKeyResult(t, svc, db, evalID, objectiveID)

	evalUUID, _ := uuid.Parse(evalID)
	if _, err := svc.SubmitKeyResults(ctx, db, evalUUID, uuid.New()); err == nil {
		t.Fatal("expected SubmitKeyResults to fail when the configured flow fails to resolve")
	}

	after, err := svc.GetEvaluationWithDetails(db, evalUUID)
	if err != nil {
		t.Fatalf("GetEvaluationWithDetails failed: %v", err)
	}
	if after.Status != "DRAFT" {
		t.Errorf("expected status to remain DRAFT after a failed submission, got %s", after.Status)
	}
}

func TestOKRService_SubmitKeyResults_NoFlowConfigured_FallsBackToManual(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()
	ctx := context.Background()

	fake := &fakeApprovalEngine{flowIDByModule: map[string]string{}}
	svc.SetApprovalEngine(fake)

	evalID, objectiveID := setupOKRTwoPhaseEvaluation(t, svc, db)
	proposeKeyResult(t, svc, db, evalID, objectiveID)

	evalUUID, _ := uuid.Parse(evalID)
	resp, err := svc.SubmitKeyResults(ctx, db, evalUUID, uuid.New())
	if err != nil {
		t.Fatalf("SubmitKeyResults failed: %v", err)
	}
	if resp.Status != "KR_SUBMITTED" {
		t.Fatalf("expected KR_SUBMITTED, got %s", resp.Status)
	}
	if resp.KRApprovalInstanceID != "" {
		t.Error("expected no kr_approval_instance_id when no flow is configured")
	}
	if len(fake.createCalls) != 0 {
		t.Errorf("expected no CreateApprovalInstance calls, got %d", len(fake.createCalls))
	}

	approved, err := svc.ApproveKeyResults(db, evalUUID, uuid.New())
	if err != nil {
		t.Fatalf("ApproveKeyResults fallback failed: %v", err)
	}
	if approved.Status != "KR_APPROVED" {
		t.Fatalf("expected KR_APPROVED, got %s", approved.Status)
	}
}

func TestOKRService_HandleKeyResultApprovalStatusChange(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()
	ctx := context.Background()

	fake := &fakeApprovalEngine{flowIDByModule: map[string]string{
		ApprovalModuleOKRKeyResult: uuid.New().String(),
	}}
	svc.SetApprovalEngine(fake)

	evalID, objectiveID := setupOKRTwoPhaseEvaluation(t, svc, db)
	proposeKeyResult(t, svc, db, evalID, objectiveID)

	evalUUID, _ := uuid.Parse(evalID)
	if _, err := svc.SubmitKeyResults(ctx, db, evalUUID, uuid.New()); err != nil {
		t.Fatalf("SubmitKeyResults failed: %v", err)
	}

	if err := svc.HandleKeyResultApprovalStatusChange(ctx, evalUUID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleKeyResultApprovalStatusChange failed: %v", err)
	}
	updated, err := svc.GetEvaluationWithDetails(db, evalUUID)
	if err != nil {
		t.Fatalf("GetEvaluationWithDetails failed: %v", err)
	}
	if updated.Status != "KR_APPROVED" {
		t.Fatalf("expected KR_APPROVED, got %s", updated.Status)
	}
}

func TestOKRService_HandleKeyResultApprovalStatusChange_Rejected_RevertsToDraft(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()
	ctx := context.Background()

	fake := &fakeApprovalEngine{flowIDByModule: map[string]string{
		ApprovalModuleOKRKeyResult: uuid.New().String(),
	}}
	svc.SetApprovalEngine(fake)

	evalID, objectiveID := setupOKRTwoPhaseEvaluation(t, svc, db)
	proposeKeyResult(t, svc, db, evalID, objectiveID)

	evalUUID, _ := uuid.Parse(evalID)
	if _, err := svc.SubmitKeyResults(ctx, db, evalUUID, uuid.New()); err != nil {
		t.Fatalf("SubmitKeyResults failed: %v", err)
	}

	if err := svc.HandleKeyResultApprovalStatusChange(ctx, evalUUID, "REJECTED", "please revise"); err != nil {
		t.Fatalf("HandleKeyResultApprovalStatusChange failed: %v", err)
	}
	updated, err := svc.GetEvaluationWithDetails(db, evalUUID)
	if err != nil {
		t.Fatalf("GetEvaluationWithDetails failed: %v", err)
	}
	if updated.Status != "DRAFT" {
		t.Fatalf("expected DRAFT after key result rejection, got %s", updated.Status)
	}
}

// TestOKRService_HandleAssessmentApprovalStatusChange_ApprovedCompletesDirectly
// verifies that final approval on the assessment checkpoint completes the
// evaluation directly (SUBMITTED -> COMPLETED), with no separate manual
// "Complete" step — mirrors the KPI behavior fixed earlier this session.
func TestOKRService_HandleAssessmentApprovalStatusChange_ApprovedCompletesDirectly(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()
	ctx := context.Background()

	fake := &fakeApprovalEngine{flowIDByModule: map[string]string{
		ApprovalModuleOKRKeyResult:  uuid.New().String(),
		ApprovalModuleOKRAssessment: uuid.New().String(),
	}}
	svc.SetApprovalEngine(fake)

	evalID, objectiveID := setupOKRTwoPhaseEvaluation(t, svc, db)
	proposeKeyResult(t, svc, db, evalID, objectiveID)

	evalUUID, _ := uuid.Parse(evalID)
	if _, err := svc.SubmitKeyResults(ctx, db, evalUUID, uuid.New()); err != nil {
		t.Fatalf("SubmitKeyResults failed: %v", err)
	}
	if err := svc.HandleKeyResultApprovalStatusChange(ctx, evalUUID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleKeyResultApprovalStatusChange failed: %v", err)
	}

	full, err := svc.GetEvaluationWithDetails(db, evalUUID)
	if err != nil {
		t.Fatalf("GetEvaluationWithDetails failed: %v", err)
	}
	detailID, _ := uuid.Parse(full.Details[0].ID)
	if _, err := svc.UpdateEvaluationDetailActual(db, detailID, &UpdateOKREvaluationDetailRequest{ActualValue: 10}); err != nil {
		t.Fatalf("UpdateEvaluationDetailActual failed: %v", err)
	}

	if _, err := svc.SubmitEvaluation(ctx, db, evalUUID, uuid.New()); err != nil {
		t.Fatalf("SubmitEvaluation failed: %v", err)
	}

	if err := svc.HandleAssessmentApprovalStatusChange(ctx, evalUUID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleAssessmentApprovalStatusChange failed: %v", err)
	}
	updated, err := svc.GetEvaluationWithDetails(db, evalUUID)
	if err != nil {
		t.Fatalf("GetEvaluationWithDetails failed: %v", err)
	}
	if updated.Status != "COMPLETED" {
		t.Fatalf("expected COMPLETED directly after final assessment approval, got %s", updated.Status)
	}
}

func TestOKRService_HandleAssessmentApprovalStatusChange_Rejected_RevertsToKRApproved(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()
	ctx := context.Background()

	fake := &fakeApprovalEngine{flowIDByModule: map[string]string{
		ApprovalModuleOKRKeyResult:  uuid.New().String(),
		ApprovalModuleOKRAssessment: uuid.New().String(),
	}}
	svc.SetApprovalEngine(fake)

	evalID, objectiveID := setupOKRTwoPhaseEvaluation(t, svc, db)
	proposeKeyResult(t, svc, db, evalID, objectiveID)

	evalUUID, _ := uuid.Parse(evalID)
	if _, err := svc.SubmitKeyResults(ctx, db, evalUUID, uuid.New()); err != nil {
		t.Fatalf("SubmitKeyResults failed: %v", err)
	}
	if err := svc.HandleKeyResultApprovalStatusChange(ctx, evalUUID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleKeyResultApprovalStatusChange failed: %v", err)
	}

	full, err := svc.GetEvaluationWithDetails(db, evalUUID)
	if err != nil {
		t.Fatalf("GetEvaluationWithDetails failed: %v", err)
	}
	detailID, _ := uuid.Parse(full.Details[0].ID)
	if _, err := svc.UpdateEvaluationDetailActual(db, detailID, &UpdateOKREvaluationDetailRequest{ActualValue: 10}); err != nil {
		t.Fatalf("UpdateEvaluationDetailActual failed: %v", err)
	}
	if _, err := svc.SubmitEvaluation(ctx, db, evalUUID, uuid.New()); err != nil {
		t.Fatalf("SubmitEvaluation failed: %v", err)
	}

	if err := svc.HandleAssessmentApprovalStatusChange(ctx, evalUUID, "REJECTED", "not accurate"); err != nil {
		t.Fatalf("HandleAssessmentApprovalStatusChange failed: %v", err)
	}
	updated, err := svc.GetEvaluationWithDetails(db, evalUUID)
	if err != nil {
		t.Fatalf("GetEvaluationWithDetails failed: %v", err)
	}
	if updated.Status != "KR_APPROVED" {
		t.Fatalf("expected KR_APPROVED after assessment rejection, got %s", updated.Status)
	}
}
