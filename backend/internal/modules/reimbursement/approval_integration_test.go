package reimbursement

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/inthros/hris-platform/internal/modules/approval"
)

// fakeApprovalEngine is a test double for ApprovalEngine.
type fakeApprovalEngine struct {
	createCalls []struct {
		module     string
		documentID string
		flowID     string
	}
	instanceID   string
	createErr    error
	activeFlowID string
	resolveErr   error
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
	if f.resolveErr != nil {
		return "", f.resolveErr
	}
	return f.activeFlowID, nil
}

func TestService_UpdateReimbursementRequestStatus_SubmitWithApprovalEngine_CreatesInstance(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{}
	svc.SetApprovalEngine(fake)

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	flowID := uuidStr()
	updated, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "SUBMITTED", "", nil, &flowID)
	if err != nil {
		t.Fatalf("UpdateReimbursementRequestStatus failed: %v", err)
	}

	if updated.Status != "PENDING_APPROVAL" {
		t.Errorf("expected status PENDING_APPROVAL, got '%s'", updated.Status)
	}
	if updated.ApprovalInstanceID == nil || *updated.ApprovalInstanceID != fake.instanceID {
		t.Errorf("expected approval_instance_id %s, got %v", fake.instanceID, updated.ApprovalInstanceID)
	}
	if len(fake.createCalls) != 1 {
		t.Fatalf("expected 1 CreateApprovalInstance call, got %d", len(fake.createCalls))
	}
	if fake.createCalls[0].module != "reimbursement" || fake.createCalls[0].flowID != flowID {
		t.Errorf("unexpected call params: %+v", fake.createCalls[0])
	}
}

func TestService_UpdateReimbursementRequestStatus_SubmitNoFlowID_SkipsApproval(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{}
	svc.SetApprovalEngine(fake)

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	updated, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "SUBMITTED", "", nil, nil)
	if err != nil {
		t.Fatalf("UpdateReimbursementRequestStatus failed: %v", err)
	}
	if updated.Status != "SUBMITTED" {
		t.Errorf("expected status SUBMITTED, got '%s'", updated.Status)
	}
	if len(fake.createCalls) != 0 {
		t.Errorf("expected no CreateApprovalInstance calls, got %d", len(fake.createCalls))
	}
}

// TestService_UpdateReimbursementRequestStatus_SubmitAutoResolvesActiveFlow
// guards the auto-resolve behavior: submitting without an explicit flow_id
// resolves the module's active flow (same pattern business travel/leave use)
// so the request still reaches the Approval inbox.
func TestService_UpdateReimbursementRequestStatus_SubmitAutoResolvesActiveFlow(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	activeFlow := uuid.New().String()
	fake := &fakeApprovalEngine{activeFlowID: activeFlow}
	svc.SetApprovalEngine(fake)

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	updated, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "SUBMITTED", "", nil, nil)
	if err != nil {
		t.Fatalf("UpdateReimbursementRequestStatus failed: %v", err)
	}
	if updated.Status != "PENDING_APPROVAL" {
		t.Errorf("expected status PENDING_APPROVAL, got '%s'", updated.Status)
	}
	if len(fake.createCalls) != 1 {
		t.Fatalf("expected 1 CreateApprovalInstance call, got %d", len(fake.createCalls))
	}
	if fake.createCalls[0].module != "reimbursement" || fake.createCalls[0].flowID != activeFlow {
		t.Errorf("expected CreateApprovalInstance with auto-resolved flow %s, got %+v", activeFlow, fake.createCalls[0])
	}
}

// TestService_UpdateReimbursementRequestStatus_SubmitFlowResolveErrorSkipsApproval
// keeps the fallback contract: if the active-flow lookup itself fails, the
// request degrades to plain SUBMITTED instead of hard-failing.
func TestService_UpdateReimbursementRequestStatus_SubmitFlowResolveErrorSkipsApproval(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{resolveErr: errors.New("flow lookup failed")}
	svc.SetApprovalEngine(fake)

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	updated, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "SUBMITTED", "", nil, nil)
	if err != nil {
		t.Fatalf("UpdateReimbursementRequestStatus failed: %v", err)
	}
	if updated.Status != "SUBMITTED" {
		t.Errorf("expected status SUBMITTED, got '%s'", updated.Status)
	}
	if len(fake.createCalls) != 0 {
		t.Errorf("expected no CreateApprovalInstance calls, got %d", len(fake.createCalls))
	}
}

func TestService_HandleApprovalStatusChange_Approved(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	rr := createTestReimbursementRequest(repo, uuid.New(), rType.ID)
	rr.Status = ReimbStatusPendingApproval
	if err := repo.UpdateReimbursementRequest(ctx(), rr); err != nil {
		t.Fatalf("failed to seed reimbursement request: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), rr.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetReimbursementRequestByID(ctx(), rr.ID.String())
	if err != nil {
		t.Fatalf("GetReimbursementRequestByID failed: %v", err)
	}
	if updated.Status != "APPROVED" {
		t.Errorf("expected status APPROVED, got '%s'", updated.Status)
	}
}

func TestService_HandleApprovalStatusChange_Rejected(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	rr := createTestReimbursementRequest(repo, uuid.New(), rType.ID)
	rr.Status = ReimbStatusPendingApproval
	if err := repo.UpdateReimbursementRequest(ctx(), rr); err != nil {
		t.Fatalf("failed to seed reimbursement request: %v", err)
	}

	note := "receipt missing"
	if err := svc.HandleApprovalStatusChange(ctx(), rr.ID, "REJECTED", note); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetReimbursementRequestByID(ctx(), rr.ID.String())
	if err != nil {
		t.Fatalf("GetReimbursementRequestByID failed: %v", err)
	}
	if updated.Status != "REJECTED" {
		t.Errorf("expected status REJECTED, got '%s'", updated.Status)
	}
	if updated.SupervisorNote == nil || *updated.SupervisorNote != note {
		t.Errorf("expected supervisor note %q, got %v", note, updated.SupervisorNote)
	}
}

func TestService_HandleApprovalStatusChange_NotPendingApproval_NoOp(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	rr := createTestReimbursementRequest(repo, uuid.New(), rType.ID)
	// Status is left as DRAFT (the default from createTestReimbursementRequest).

	if err := svc.HandleApprovalStatusChange(ctx(), rr.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetReimbursementRequestByID(ctx(), rr.ID.String())
	if err != nil {
		t.Fatalf("GetReimbursementRequestByID failed: %v", err)
	}
	if updated.Status != "DRAFT" {
		t.Errorf("expected status to remain DRAFT, got '%s'", updated.Status)
	}
}

// TestService_UpdateReimbursementRequestStatus_ApprovalRoutingErrorFailsLoudly
// guards the fail-loudly policy: when the approval engine rejects routing, the
// RoutingError is propagated (instead of silently continuing as SUBMITTED) so
// the handler can show a bilingual message.
func TestService_UpdateReimbursementRequestStatus_ApprovalRoutingErrorFailsLoudly(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{
		createErr: &approval.RoutingError{
			Key:    "approval.zero_approvers",
			Params: []string{"Persetujuan Supervisor", "1"},
		},
	}
	svc.SetApprovalEngine(fake)

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	flowID := uuidStr()
	_, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "SUBMITTED", "", nil, &flowID)
	if err == nil {
		t.Fatal("expected error when approval routing fails")
	}
	var re *approval.RoutingError
	if !errors.As(err, &re) {
		t.Fatalf("expected approval.RoutingError, got: %v", err)
	}
}

// TestService_UpdateReimbursementRequestStatus_NonRoutingApprovalErrorStillSwallowed
// keeps the best-effort contract for non-routing approval failures.
func TestService_UpdateReimbursementRequestStatus_NonRoutingApprovalErrorStillSwallowed(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{createErr: errors.New("database connection lost")}
	svc.SetApprovalEngine(fake)

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	flowID := uuidStr()
	updated, err := svc.UpdateReimbursementRequestStatus(ctx(), created.ID.String(), "SUBMITTED", "", nil, &flowID)
	if err != nil {
		t.Fatalf("should not fail on non-routing approval error: %v", err)
	}
	if updated.Status != "SUBMITTED" {
		t.Errorf("expected status SUBMITTED without approval, got '%s'", updated.Status)
	}
}
