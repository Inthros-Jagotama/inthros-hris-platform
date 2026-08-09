package employeemovement

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// fakeApprovalEngine is a test double for ApprovalEngine.
type fakeApprovalEngine struct {
	createCalls []struct {
		module     string
		documentID string
		flowID     string
	}
	instanceID      string
	createErr       error
	resolvedFlowID  string
	resolvedFlowErr error
	resolveCalls    []string // module per call
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
	if f.instanceID == "" {
		f.instanceID = uuid.New().String()
	}
	return f.instanceID, nil
}

func (f *fakeApprovalEngine) GetApprovalInstanceStatus(ctx context.Context, instanceID string) (string, error) {
	return "PENDING", nil
}

func (f *fakeApprovalEngine) GetActiveFlowIDForModule(ctx context.Context, module string) (string, error) {
	f.resolveCalls = append(f.resolveCalls, module)
	if f.resolvedFlowErr != nil {
		return "", f.resolvedFlowErr
	}
	return f.resolvedFlowID, nil
}

func TestService_SubmitMovement_WithApprovalEngine_CreatesInstance(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{}
	svc.SetApprovalEngine(fake)

	empID := uuid.New()
	movement := createTestMovement(repo, empID)
	flowID := uuid.New().String()

	resp, err := svc.SubmitMovement(ctx(), movement.ID.String(), SubmitMovementRequest{FlowID: &flowID})
	if err != nil {
		t.Fatalf("SubmitMovement failed: %v", err)
	}

	if resp.Status != "pending_approval" {
		t.Errorf("expected status pending_approval, got '%s'", resp.Status)
	}
	if resp.ApprovalInstanceID == nil || *resp.ApprovalInstanceID != fake.instanceID {
		t.Errorf("expected approval_instance_id %s, got %v", fake.instanceID, resp.ApprovalInstanceID)
	}
	if len(fake.createCalls) != 1 {
		t.Fatalf("expected 1 CreateApprovalInstance call, got %d", len(fake.createCalls))
	}
	if fake.createCalls[0].module != "employeemovement" || fake.createCalls[0].flowID != flowID {
		t.Errorf("unexpected call params: %+v", fake.createCalls[0])
	}
}

func TestService_SubmitMovement_NoFlowID_AutoResolves(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	resolvedFlow := uuid.New().String()
	fake := &fakeApprovalEngine{resolvedFlowID: resolvedFlow}
	svc.SetApprovalEngine(fake)

	empID := uuid.New()
	movement := createTestMovement(repo, empID)

	resp, err := svc.SubmitMovement(ctx(), movement.ID.String(), SubmitMovementRequest{})
	if err != nil {
		t.Fatalf("SubmitMovement failed: %v", err)
	}

	if resp.Status != "pending_approval" {
		t.Errorf("expected status pending_approval, got '%s'", resp.Status)
	}
	// Auto-resolve called for module employeemovement.
	if len(fake.resolveCalls) != 1 || fake.resolveCalls[0] != "employeemovement" {
		t.Errorf("expected GetActiveFlowIDForModule('employeemovement') call, got %v", fake.resolveCalls)
	}
	// Instance created with the resolved flow.
	if len(fake.createCalls) != 1 || fake.createCalls[0].flowID != resolvedFlow {
		t.Errorf("expected CreateApprovalInstance with resolved flow %s, got %+v", resolvedFlow, fake.createCalls)
	}
	if resp.ApprovalInstanceID == nil || *resp.ApprovalInstanceID != fake.instanceID {
		t.Errorf("expected approval_instance_id %s, got %v", fake.instanceID, resp.ApprovalInstanceID)
	}
}

func TestService_SubmitMovement_NoFlowResolved_ReturnsError(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	// No active flow configured — auto-resolve returns empty.
	fake := &fakeApprovalEngine{}
	svc.SetApprovalEngine(fake)

	empID := uuid.New()
	movement := createTestMovement(repo, empID)

	_, err := svc.SubmitMovement(ctx(), movement.ID.String(), SubmitMovementRequest{})
	if err == nil {
		t.Fatal("expected error when no flow can be resolved")
	}
	if len(fake.createCalls) != 0 {
		t.Errorf("expected no CreateApprovalInstance calls, got %d", len(fake.createCalls))
	}
}

func TestService_SubmitMovement_ExplicitFlowID_BeatsAutoResolve(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{resolvedFlowID: uuid.New().String()}
	svc.SetApprovalEngine(fake)

	empID := uuid.New()
	movement := createTestMovement(repo, empID)
	explicitFlow := uuid.New().String()

	resp, err := svc.SubmitMovement(ctx(), movement.ID.String(), SubmitMovementRequest{FlowID: &explicitFlow})
	if err != nil {
		t.Fatalf("SubmitMovement failed: %v", err)
	}

	// Explicit flow_id wins — no auto-resolve call made.
	if len(fake.resolveCalls) != 0 {
		t.Errorf("expected no auto-resolve calls when flow_id supplied, got %v", fake.resolveCalls)
	}
	if len(fake.createCalls) != 1 || fake.createCalls[0].flowID != explicitFlow {
		t.Errorf("expected CreateApprovalInstance with explicit flow %s, got %+v", explicitFlow, fake.createCalls)
	}
	if resp.Status != "pending_approval" {
		t.Errorf("expected status pending_approval, got '%s'", resp.Status)
	}
}

func TestService_SubmitMovement_NotDraft_ReturnsError(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{}
	svc.SetApprovalEngine(fake)

	empID := uuid.New()
	movement := createTestMovement(repo, empID)
	movement.Status = MovementStatusApproved
	if err := repo.UpdateMovement(ctx(), movement); err != nil {
		t.Fatalf("failed to seed movement: %v", err)
	}

	flowID := uuid.New().String()
	_, err := svc.SubmitMovement(ctx(), movement.ID.String(), SubmitMovementRequest{FlowID: &flowID})
	if err == nil {
		t.Fatal("expected error when movement is not draft")
	}
}

func TestService_HandleApprovalStatusChange_Approved(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	empID := uuid.New()
	movement := createTestMovement(repo, empID)
	movement.Status = MovementStatusPendingApproval
	if err := repo.UpdateMovement(ctx(), movement); err != nil {
		t.Fatalf("failed to seed movement: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), movement.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetMovementByID(ctx(), movement.ID.String())
	if err != nil {
		t.Fatalf("GetMovementByID failed: %v", err)
	}
	if updated.Status != "approved" {
		t.Errorf("expected status approved, got '%s'", updated.Status)
	}
	if updated.ApprovedAt == nil {
		t.Error("expected ApprovedAt to be set")
	}
}

func TestService_HandleApprovalStatusChange_Rejected(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	empID := uuid.New()
	movement := createTestMovement(repo, empID)
	movement.Status = MovementStatusPendingApproval
	if err := repo.UpdateMovement(ctx(), movement); err != nil {
		t.Fatalf("failed to seed movement: %v", err)
	}

	note := "position no longer available"
	if err := svc.HandleApprovalStatusChange(ctx(), movement.ID, "REJECTED", note); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetMovementByID(ctx(), movement.ID.String())
	if err != nil {
		t.Fatalf("GetMovementByID failed: %v", err)
	}
	if updated.Status != "cancelled" {
		t.Errorf("expected status cancelled (rejection reuses cancelled), got '%s'", updated.Status)
	}
	if updated.Notes == nil || *updated.Notes != note {
		t.Errorf("expected notes %q, got %v", note, updated.Notes)
	}
}

func TestService_HandleApprovalStatusChange_NotPendingApproval_NoOp(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	empID := uuid.New()
	movement := createTestMovement(repo, empID)
	// Status is left as draft (the default from createTestMovement).

	if err := svc.HandleApprovalStatusChange(ctx(), movement.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetMovementByID(ctx(), movement.ID.String())
	if err != nil {
		t.Fatalf("GetMovementByID failed: %v", err)
	}
	if updated.Status != "draft" {
		t.Errorf("expected status to remain draft, got '%s'", updated.Status)
	}
}
