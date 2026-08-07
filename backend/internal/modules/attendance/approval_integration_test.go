package attendance

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

func TestService_CreateOvertimeRequest_WithApprovalEngine_CreatesInstance(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{}
	svc.SetApprovalEngine(fake)

	flowID := uuid.New().String()
	req := CreateOvertimeRequest{
		EmployeeID:       uuid.New().String(),
		WorkDate:         "2026-01-15",
		StartTimeLocal:   "2026-01-15T18:00:00+07:00",
		EndTimeLocal:     "2026-01-15T20:00:00+07:00",
		RequestedMinutes: 120,
		Reason:           "Deadline crunch",
		FlowID:           &flowID,
	}

	resp, err := svc.CreateOvertimeRequest(ctx(), req)
	if err != nil {
		t.Fatalf("CreateOvertimeRequest failed: %v", err)
	}

	if resp.Status != "PENDING_APPROVAL" {
		t.Errorf("expected status PENDING_APPROVAL, got '%s'", resp.Status)
	}
	if resp.ApprovalInstanceID == nil || *resp.ApprovalInstanceID != fake.instanceID {
		t.Errorf("expected approval_instance_id %s, got %v", fake.instanceID, resp.ApprovalInstanceID)
	}
	if len(fake.createCalls) != 1 {
		t.Fatalf("expected 1 CreateApprovalInstance call, got %d", len(fake.createCalls))
	}
	if fake.createCalls[0].module != "attendance" || fake.createCalls[0].flowID != flowID {
		t.Errorf("unexpected call params: %+v", fake.createCalls[0])
	}
}

func TestService_CreateOvertimeRequest_NoFlowID_SkipsApproval(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{}
	svc.SetApprovalEngine(fake)

	req := CreateOvertimeRequest{
		EmployeeID:       uuid.New().String(),
		WorkDate:         "2026-01-15",
		StartTimeLocal:   "2026-01-15T18:00:00+07:00",
		EndTimeLocal:     "2026-01-15T20:00:00+07:00",
		RequestedMinutes: 120,
		Reason:           "Deadline crunch",
	}

	resp, err := svc.CreateOvertimeRequest(ctx(), req)
	if err != nil {
		t.Fatalf("CreateOvertimeRequest failed: %v", err)
	}
	if resp.Status != "SUBMITTED" {
		t.Errorf("expected status SUBMITTED, got '%s'", resp.Status)
	}
	if len(fake.createCalls) != 0 {
		t.Errorf("expected no CreateApprovalInstance calls, got %d", len(fake.createCalls))
	}
}

func TestService_HandleApprovalStatusChange_Approved(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	o := createTestOvertimeRequest(repo, uuid.New())
	o.Status = OvertimePendingApproval
	if err := repo.UpdateOvertimeRequest(ctx(), o); err != nil {
		t.Fatalf("failed to seed overtime request: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), o.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetOvertimeRequestByID(ctx(), o.ID.String())
	if err != nil {
		t.Fatalf("GetOvertimeRequestByID failed: %v", err)
	}
	if updated.Status != "APPROVED" {
		t.Errorf("expected status APPROVED, got '%s'", updated.Status)
	}
	if updated.ApprovedAt == nil {
		t.Error("expected ApprovedAt to be set")
	}
}

func TestService_HandleApprovalStatusChange_Rejected(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	o := createTestOvertimeRequest(repo, uuid.New())
	o.Status = OvertimePendingApproval
	if err := repo.UpdateOvertimeRequest(ctx(), o); err != nil {
		t.Fatalf("failed to seed overtime request: %v", err)
	}

	note := "insufficient justification"
	if err := svc.HandleApprovalStatusChange(ctx(), o.ID, "REJECTED", note); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetOvertimeRequestByID(ctx(), o.ID.String())
	if err != nil {
		t.Fatalf("GetOvertimeRequestByID failed: %v", err)
	}
	if updated.Status != "REJECTED" {
		t.Errorf("expected status REJECTED, got '%s'", updated.Status)
	}
	if updated.ApprovalNote == nil || *updated.ApprovalNote != note {
		t.Errorf("expected approval note %q, got %v", note, updated.ApprovalNote)
	}
}

func TestService_HandleApprovalStatusChange_NotPendingApproval_NoOp(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	o := createTestOvertimeRequest(repo, uuid.New())
	// Status is left as SUBMITTED (the default from createTestOvertimeRequest).

	if err := svc.HandleApprovalStatusChange(ctx(), o.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetOvertimeRequestByID(ctx(), o.ID.String())
	if err != nil {
		t.Fatalf("GetOvertimeRequestByID failed: %v", err)
	}
	if updated.Status != "SUBMITTED" {
		t.Errorf("expected status to remain SUBMITTED, got '%s'", updated.Status)
	}
}
