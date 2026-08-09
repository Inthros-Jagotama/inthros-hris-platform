package leave

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// fakeApprovalEngine is a test double for ApprovalEngine.
type fakeApprovalEngine struct {
	createCalls []struct {
		module     string
		documentID string
		flowID     string
	}
	instanceID    string
	createErr     error
	activeFlowID  string
	activeFlowErr error
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
	if f.activeFlowErr != nil {
		return "", f.activeFlowErr
	}
	if f.activeFlowID == "" {
		return "", fmt.Errorf("no active flow configured")
	}
	return f.activeFlowID, nil
}

func TestService_CreateLeaveRequest_WithApprovalEngine_CreatesInstance(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{}
	svc.SetApprovalEngine(fake)

	lType := createTestLeaveType(repo)
	flowID := uuidStr()

	req := CreateLeaveRequest{
		EmployeeID:       uuidStr(),
		LeaveTypeID:      lType.ID.String(),
		RequestStartDate: "2026-01-15",
		RequestEndDate:   "2026-01-16",
		RequestedDays:    2,
		FlowID:           &flowID,
	}

	resp, err := svc.CreateLeaveRequest(ctx(), req)
	if err != nil {
		t.Fatalf("CreateLeaveRequest failed: %v", err)
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
	if fake.createCalls[0].module != "leave" || fake.createCalls[0].flowID != flowID {
		t.Errorf("unexpected call params: %+v", fake.createCalls[0])
	}
}

// TestService_CreateLeaveRequest_NoFlowID_AutoResolvesActiveFlow guards the
// fix for the bug where a leave request created without an explicit flow_id
// (the normal case — no FE sends one) silently stayed SUBMITTED forever and
// never reached the Approval module, because CreateLeaveRequest previously
// only routed through approval when the client supplied flow_id itself.
func TestService_CreateLeaveRequest_NoFlowID_AutoResolvesActiveFlow(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	activeFlowID := uuidStr()
	fake := &fakeApprovalEngine{activeFlowID: activeFlowID}
	svc.SetApprovalEngine(fake)

	lType := createTestLeaveType(repo)

	req := CreateLeaveRequest{
		EmployeeID:       uuidStr(),
		LeaveTypeID:      lType.ID.String(),
		RequestStartDate: "2026-01-15",
		RequestEndDate:   "2026-01-16",
		RequestedDays:    2,
	}

	resp, err := svc.CreateLeaveRequest(ctx(), req)
	if err != nil {
		t.Fatalf("CreateLeaveRequest failed: %v", err)
	}
	if resp.Status != "PENDING_APPROVAL" {
		t.Errorf("expected status PENDING_APPROVAL, got '%s'", resp.Status)
	}
	if resp.ApprovalInstanceID == nil {
		t.Fatal("expected approval_instance_id to be set")
	}
	if len(fake.createCalls) != 1 {
		t.Fatalf("expected 1 CreateApprovalInstance call, got %d", len(fake.createCalls))
	}
	if fake.createCalls[0].flowID != activeFlowID {
		t.Errorf("expected auto-resolved flow_id %s, got %s", activeFlowID, fake.createCalls[0].flowID)
	}
}

func TestService_CreateLeaveRequest_NoFlowIDAndNoActiveFlow_SkipsApproval(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{} // no active flow configured
	svc.SetApprovalEngine(fake)

	lType := createTestLeaveType(repo)

	req := CreateLeaveRequest{
		EmployeeID:       uuidStr(),
		LeaveTypeID:      lType.ID.String(),
		RequestStartDate: "2026-01-15",
		RequestEndDate:   "2026-01-16",
		RequestedDays:    2,
	}

	resp, err := svc.CreateLeaveRequest(ctx(), req)
	if err != nil {
		t.Fatalf("CreateLeaveRequest failed: %v", err)
	}
	if resp.Status != "SUBMITTED" {
		t.Errorf("expected status SUBMITTED, got '%s'", resp.Status)
	}
	if len(fake.createCalls) != 0 {
		t.Errorf("expected no CreateApprovalInstance calls, got %d", len(fake.createCalls))
	}
}

// leaveRequestStatusNote is a narrow projection used to assert persisted
// state without scanning the *time.Time columns: the glebarez/sqlite test
// driver cannot round-trip GORM's `type:timestamp(6)` columns once they hold
// a non-null value (pre-existing, unrelated to approval integration — see
// TestService_UpdateLeaveRequestStatus_Approve, which sidesteps it the same
// way by asserting off the in-memory struct instead of a re-SELECT).
type leaveRequestStatusNote struct {
	Status         string
	SupervisorNote *string
}

func fetchLeaveRequestStatusNote(t *testing.T, db *gorm.DB, id uuid.UUID) leaveRequestStatusNote {
	t.Helper()
	var row leaveRequestStatusNote
	if err := db.Table("leave_requests").Select("status", "supervisor_note").Where("id = ?", id.String()).Scan(&row).Error; err != nil {
		t.Fatalf("failed to fetch leave request status: %v", err)
	}
	return row
}

func TestService_HandleApprovalStatusChange_Approved(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()

	lType := createTestLeaveType(repo)
	lr := createTestLeaveRequest(repo, uuid.New(), lType.ID)
	lr.Status = LeaveStatusPendingApproval
	if err := repo.UpdateLeaveRequest(ctx(), lr); err != nil {
		t.Fatalf("failed to seed leave request: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), lr.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	row := fetchLeaveRequestStatusNote(t, db, lr.ID)
	if row.Status != "APPROVED_FINAL" {
		t.Errorf("expected status APPROVED_FINAL, got '%s'", row.Status)
	}
}

func TestService_HandleApprovalStatusChange_Rejected(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()

	lType := createTestLeaveType(repo)
	lr := createTestLeaveRequest(repo, uuid.New(), lType.ID)
	lr.Status = LeaveStatusPendingApproval
	if err := repo.UpdateLeaveRequest(ctx(), lr); err != nil {
		t.Fatalf("failed to seed leave request: %v", err)
	}

	note := "insufficient balance"
	if err := svc.HandleApprovalStatusChange(ctx(), lr.ID, "REJECTED", note); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	row := fetchLeaveRequestStatusNote(t, db, lr.ID)
	if row.Status != "REJECTED_FINAL" {
		t.Errorf("expected status REJECTED_FINAL, got '%s'", row.Status)
	}
	if row.SupervisorNote == nil || *row.SupervisorNote != note {
		t.Errorf("expected supervisor note %q, got %v", note, row.SupervisorNote)
	}
}

func TestService_HandleApprovalStatusChange_NotPendingApproval_NoOp(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	lType := createTestLeaveType(repo)
	lr := createTestLeaveRequest(repo, uuid.New(), lType.ID)
	// Status is left as SUBMITTED (the default from createTestLeaveRequest).

	if err := svc.HandleApprovalStatusChange(ctx(), lr.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetLeaveRequestByID(ctx(), lr.ID.String())
	if err != nil {
		t.Fatalf("GetLeaveRequestByID failed: %v", err)
	}
	if updated.Status != "SUBMITTED" {
		t.Errorf("expected status to remain SUBMITTED, got '%s'", updated.Status)
	}
}
