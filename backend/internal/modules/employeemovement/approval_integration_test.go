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
	instanceID string
	createErr  error
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

func TestService_SubmitMovement_NoFlowID_ReturnsError(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{}
	svc.SetApprovalEngine(fake)

	empID := uuid.New()
	movement := createTestMovement(repo, empID)

	_, err := svc.SubmitMovement(ctx(), movement.ID.String(), SubmitMovementRequest{})
	if err == nil {
		t.Fatal("expected error when flow_id is missing")
	}
	if len(fake.createCalls) != 0 {
		t.Errorf("expected no CreateApprovalInstance calls, got %d", len(fake.createCalls))
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
