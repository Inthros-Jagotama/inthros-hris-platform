package performance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/inthros/hris-platform/internal/modules/approval"
)

// fakeApprovalEngine is a test double for ApprovalEngine.
type fakeApprovalEngine struct {
	flowIDByModule map[string]string
	createErr      error
	createCalls    []struct {
		module     string
		documentID string
		flowID     string
	}
}

func (f *fakeApprovalEngine) CreateApprovalInstance(ctx context.Context, module, documentID, flowID string) (string, error) {
	f.createCalls = append(f.createCalls, struct {
		module     string
		documentID string
		flowID     string
	}{module, documentID, flowID})
	if f.createErr != nil {
		return "", f.createErr
	}
	return uuid.New().String(), nil
}

func (f *fakeApprovalEngine) GetActiveFlowIDForModule(ctx context.Context, module string) (string, error) {
	flowID, ok := f.flowIDByModule[module]
	if !ok {
		return "", fmt.Errorf("no active flow configured for module %q", module)
	}
	return flowID, nil
}

func TestService_SubmitTarget_RoutesThroughApprovalEngine(t *testing.T) {
	svc, _, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()
	ctx := context.Background()

	fake := &fakeApprovalEngine{flowIDByModule: map[string]string{
		ApprovalModuleKPITarget: uuid.New().String(),
	}}
	svc.SetApprovalEngine(fake)

	evalID := setupTwoPhaseEvaluation(t, svc)
	full, _ := svc.GetEvaluationWithDetails(ctx, evalID)
	if _, err := svc.UpdateEvaluationTarget(ctx, full.Details[0].ID, UpdateEvaluationTargetRequest{Target: 1000}); err != nil {
		t.Fatalf("UpdateEvaluationTarget failed: %v", err)
	}

	resp, err := svc.SubmitTarget(ctx, evalID)
	if err != nil {
		t.Fatalf("SubmitTarget failed: %v", err)
	}
	if resp.Status != "TARGET_SUBMITTED" {
		t.Fatalf("expected TARGET_SUBMITTED, got %s", resp.Status)
	}
	if resp.TargetApprovalInstanceID == "" {
		t.Error("expected target_approval_instance_id to be set")
	}
	if len(fake.createCalls) != 1 {
		t.Fatalf("expected 1 CreateApprovalInstance call, got %d", len(fake.createCalls))
	}
	if fake.createCalls[0].module != ApprovalModuleKPITarget || fake.createCalls[0].documentID != evalID {
		t.Errorf("unexpected call params: %+v", fake.createCalls[0])
	}
}

// TestService_SubmitTarget_FlowConfiguredButResolutionFails_HardFails
// validates that when a flow IS configured for the module but resolving it
// into an actual approval instance fails (e.g. the step's approver is
// vacant), the whole submission is rejected — status stays DRAFT — instead
// of silently degrading to the no-approval-instance manual fallback, which
// would let the submission through with nobody able to approve it.
func TestService_SubmitTarget_FlowConfiguredButResolutionFails_HardFails(t *testing.T) {
	svc, _, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()
	ctx := context.Background()

	fake := &fakeApprovalEngine{
		flowIDByModule: map[string]string{ApprovalModuleKPITarget: uuid.New().String()},
		createErr:      fmt.Errorf("approval step resolved to zero approvers"),
	}
	svc.SetApprovalEngine(fake)

	evalID := setupTwoPhaseEvaluation(t, svc)
	full, _ := svc.GetEvaluationWithDetails(ctx, evalID)
	if _, err := svc.UpdateEvaluationTarget(ctx, full.Details[0].ID, UpdateEvaluationTargetRequest{Target: 1000}); err != nil {
		t.Fatalf("UpdateEvaluationTarget failed: %v", err)
	}

	if _, err := svc.SubmitTarget(ctx, evalID); err == nil {
		t.Fatal("expected SubmitTarget to fail when the configured flow fails to resolve")
	}

	after, err := svc.GetEvaluationWithDetails(ctx, evalID)
	if err != nil {
		t.Fatalf("GetEvaluationWithDetails failed: %v", err)
	}
	if after.Status != "DRAFT" {
		t.Errorf("expected status to remain DRAFT after a failed submission, got %s", after.Status)
	}
}

func TestService_SubmitTarget_NoFlowConfigured_FallsBackToManual(t *testing.T) {
	svc, _, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// No flow registered for ApprovalModuleKPITarget in flowIDByModule.
	fake := &fakeApprovalEngine{flowIDByModule: map[string]string{}}
	svc.SetApprovalEngine(fake)

	evalID := setupTwoPhaseEvaluation(t, svc)
	full, _ := svc.GetEvaluationWithDetails(ctx, evalID)
	if _, err := svc.UpdateEvaluationTarget(ctx, full.Details[0].ID, UpdateEvaluationTargetRequest{Target: 1000}); err != nil {
		t.Fatalf("UpdateEvaluationTarget failed: %v", err)
	}

	resp, err := svc.SubmitTarget(ctx, evalID)
	if err != nil {
		t.Fatalf("SubmitTarget failed: %v", err)
	}
	if resp.Status != "TARGET_SUBMITTED" {
		t.Fatalf("expected TARGET_SUBMITTED, got %s", resp.Status)
	}
	if resp.TargetApprovalInstanceID != "" {
		t.Error("expected no target_approval_instance_id when no flow is configured")
	}
	if len(fake.createCalls) != 0 {
		t.Errorf("expected no CreateApprovalInstance calls, got %d", len(fake.createCalls))
	}

	// Manual ApproveTarget fallback must still work.
	approved, err := svc.ApproveTarget(ctx, evalID)
	if err != nil {
		t.Fatalf("ApproveTarget fallback failed: %v", err)
	}
	if approved.Status != "TARGET_APPROVED" {
		t.Fatalf("expected TARGET_APPROVED, got %s", approved.Status)
	}
}

func TestService_HandleTargetApprovalStatusChange(t *testing.T) {
	svc, _, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()
	ctx := context.Background()

	fake := &fakeApprovalEngine{flowIDByModule: map[string]string{
		ApprovalModuleKPITarget: uuid.New().String(),
	}}
	svc.SetApprovalEngine(fake)

	evalID := setupTwoPhaseEvaluation(t, svc)
	full, _ := svc.GetEvaluationWithDetails(ctx, evalID)
	if _, err := svc.UpdateEvaluationTarget(ctx, full.Details[0].ID, UpdateEvaluationTargetRequest{Target: 1000}); err != nil {
		t.Fatalf("UpdateEvaluationTarget failed: %v", err)
	}
	if _, err := svc.SubmitTarget(ctx, evalID); err != nil {
		t.Fatalf("SubmitTarget failed: %v", err)
	}

	evalUUID, _ := uuid.Parse(evalID)
	if err := svc.HandleTargetApprovalStatusChange(ctx, evalUUID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleTargetApprovalStatusChange failed: %v", err)
	}
	updated, err := svc.GetEvaluationWithDetails(ctx, evalID)
	if err != nil {
		t.Fatalf("GetEvaluationWithDetails failed: %v", err)
	}
	if updated.Status != "TARGET_APPROVED" {
		t.Fatalf("expected TARGET_APPROVED, got %s", updated.Status)
	}
}

func TestService_HandleTargetApprovalStatusChange_Rejected_RevertsToDraft(t *testing.T) {
	svc, _, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()
	ctx := context.Background()

	fake := &fakeApprovalEngine{flowIDByModule: map[string]string{
		ApprovalModuleKPITarget: uuid.New().String(),
	}}
	svc.SetApprovalEngine(fake)

	evalID := setupTwoPhaseEvaluation(t, svc)
	full, _ := svc.GetEvaluationWithDetails(ctx, evalID)
	if _, err := svc.UpdateEvaluationTarget(ctx, full.Details[0].ID, UpdateEvaluationTargetRequest{Target: 1000}); err != nil {
		t.Fatalf("UpdateEvaluationTarget failed: %v", err)
	}
	if _, err := svc.SubmitTarget(ctx, evalID); err != nil {
		t.Fatalf("SubmitTarget failed: %v", err)
	}

	evalUUID, _ := uuid.Parse(evalID)
	if err := svc.HandleTargetApprovalStatusChange(ctx, evalUUID, "REJECTED", "please revise"); err != nil {
		t.Fatalf("HandleTargetApprovalStatusChange failed: %v", err)
	}
	updated, err := svc.GetEvaluationWithDetails(ctx, evalID)
	if err != nil {
		t.Fatalf("GetEvaluationWithDetails failed: %v", err)
	}
	if updated.Status != "DRAFT" {
		t.Fatalf("expected DRAFT after target rejection, got %s", updated.Status)
	}
}

// TestEmitKPIStatusError_RoutingErrorBilingual verifies that approval routing
// failures (e.g. "no supervisor found: every organization ... is vacant") are
// emitted with a translated message following the request language instead of
// leaking the raw English service error.
func TestEmitKPIStatusError_RoutingErrorBilingual(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tenant/performance/kpi/evaluations/x/submit-target", nil)
	c.Request.Header.Set("Accept-Language", "id")

	emitKPIStatusError(c, &approval.RoutingError{
		Key:    "approval.no_supervisor_vacant",
		Params: []string{"Persetujuan Supervisor"},
	})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	errObj, ok := body["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing error object: %+v", body)
	}
	if errObj["code"] != "APPROVAL_ROUTING_FAILED" {
		t.Errorf("expected code APPROVAL_ROUTING_FAILED, got %v", errObj["code"])
	}
	msg, _ := errObj["message"].(string)
	if !strings.Contains(msg, "Supervisor tidak ditemukan untuk langkah \"Persetujuan Supervisor\"") {
		t.Errorf("expected Indonesian supervisor message, got: %s", msg)
	}
}

// TestEmitKPIStatusError_NonRoutingFallsBackToRaw ensures non-routing errors
// keep their existing raw-message behavior.
func TestEmitKPIStatusError_NonRoutingFallsBackToRaw(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tenant/performance/kpi/evaluations/x/submit-target", nil)

	emitKPIStatusError(c, fmt.Errorf("target can only be submitted from DRAFT status, current: SUBMITTED"))

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	errObj := body["error"].(map[string]interface{})
	if msg, _ := errObj["message"].(string); !strings.Contains(msg, "target can only be submitted from DRAFT status") {
		t.Errorf("expected raw fallback message, got: %s", msg)
	}
}

func TestService_HandleRealizationApprovalStatusChange(t *testing.T) {
	svc, _, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()
	ctx := context.Background()

	fake := &fakeApprovalEngine{flowIDByModule: map[string]string{
		ApprovalModuleKPITarget:      uuid.New().String(),
		ApprovalModuleKPIRealization: uuid.New().String(),
	}}
	svc.SetApprovalEngine(fake)

	evalID := setupTwoPhaseEvaluation(t, svc)
	full, _ := svc.GetEvaluationWithDetails(ctx, evalID)
	detailID := full.Details[0].ID
	if _, err := svc.UpdateEvaluationTarget(ctx, detailID, UpdateEvaluationTargetRequest{Target: 1000}); err != nil {
		t.Fatalf("UpdateEvaluationTarget failed: %v", err)
	}
	if _, err := svc.SubmitTarget(ctx, evalID); err != nil {
		t.Fatalf("SubmitTarget failed: %v", err)
	}
	evalUUID, _ := uuid.Parse(evalID)
	if err := svc.HandleTargetApprovalStatusChange(ctx, evalUUID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleTargetApprovalStatusChange failed: %v", err)
	}
	if _, err := svc.UpdateEvaluationActual(ctx, detailID, UpdateEvaluationActualRequest{Actual: 1000}); err != nil {
		t.Fatalf("UpdateEvaluationActual failed: %v", err)
	}

	resp, err := svc.SubmitEvaluation(ctx, evalID)
	if err != nil {
		t.Fatalf("SubmitEvaluation failed: %v", err)
	}
	if resp.RealizationApprovalInstanceID == "" {
		t.Error("expected realization_approval_instance_id to be set")
	}

	if err := svc.HandleRealizationApprovalStatusChange(ctx, evalUUID, "REJECTED", "not accurate"); err != nil {
		t.Fatalf("HandleRealizationApprovalStatusChange failed: %v", err)
	}
	updated, err := svc.GetEvaluationWithDetails(ctx, evalID)
	if err != nil {
		t.Fatalf("GetEvaluationWithDetails failed: %v", err)
	}
	if updated.Status != "TARGET_APPROVED" {
		t.Fatalf("expected TARGET_APPROVED after realization rejection, got %s", updated.Status)
	}
}
