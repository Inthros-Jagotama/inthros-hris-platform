package leave

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

type fakeNotifier struct {
	calls []fakeNotifyCall
	err   error
}

type fakeNotifyCall struct {
	recipientUserID uuid.UUID
	notifType       string
	referenceType   string
	referenceID     uuid.UUID
}

func (f *fakeNotifier) Notify(ctx context.Context, recipientUserID uuid.UUID, notifType string, params []string, referenceType string, referenceID uuid.UUID) error {
	f.calls = append(f.calls, fakeNotifyCall{recipientUserID, notifType, referenceType, referenceID})
	return f.err
}

func TestService_HandleApprovalStatusChange_NotifiesEmployeeOnApproval(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	userID := uuid.New()
	createTestEmployeeAccount(db, empID, userID)

	lType := createTestLeaveType(repo)
	lr := createTestLeaveRequest(repo, empID, lType.ID)
	lr.Status = LeaveStatusPendingApproval
	if err := repo.UpdateLeaveRequest(context.Background(), lr); err != nil {
		t.Fatalf("failed to seed pending status: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(context.Background(), lr.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notify call, got %d", len(notifier.calls))
	}
	call := notifier.calls[0]
	if call.recipientUserID != userID {
		t.Errorf("expected recipient %s, got %s", userID, call.recipientUserID)
	}
	if call.notifType != "LEAVE_APPROVED" {
		t.Errorf("expected notif type LEAVE_APPROVED, got %s", call.notifType)
	}
	if call.referenceType != "leave" || call.referenceID != lr.ID {
		t.Errorf("expected reference leave/%s, got %s/%s", lr.ID, call.referenceType, call.referenceID)
	}
}

func TestService_HandleApprovalStatusChange_SkipsNotifyWhenNoLinkedUserAccount(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New() // no employee_accounts row created
	lType := createTestLeaveType(repo)
	lr := createTestLeaveRequest(repo, empID, lType.ID)
	lr.Status = LeaveStatusPendingApproval
	if err := repo.UpdateLeaveRequest(context.Background(), lr); err != nil {
		t.Fatalf("failed to seed pending status: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(context.Background(), lr.ID, "REJECTED", "not eligible"); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	if len(notifier.calls) != 0 {
		t.Fatalf("expected 0 notify calls when employee has no linked user account, got %d", len(notifier.calls))
	}
}

func TestService_HandleApprovalStatusChange_NotifyFailureDoesNotFailApproval(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()

	notifier := &fakeNotifier{err: context.DeadlineExceeded}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	createTestEmployeeAccount(db, empID, uuid.New())

	lType := createTestLeaveType(repo)
	lr := createTestLeaveRequest(repo, empID, lType.ID)
	lr.Status = LeaveStatusPendingApproval
	if err := repo.UpdateLeaveRequest(context.Background(), lr); err != nil {
		t.Fatalf("failed to seed pending status: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(context.Background(), lr.ID, "REJECTED", ""); err != nil {
		t.Fatalf("expected HandleApprovalStatusChange to succeed despite notify failure, got: %v", err)
	}
	if len(notifier.calls) != 1 {
		t.Fatalf("expected the notifier to still be called once despite it returning an error, got %d calls", len(notifier.calls))
	}
}

func TestService_HandleApprovalStatusChange_NotifiesEmployeeOnCancellation(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	userID := uuid.New()
	createTestEmployeeAccount(db, empID, userID)

	lType := createTestLeaveType(repo)
	lr := createTestLeaveRequest(repo, empID, lType.ID)
	lr.Status = LeaveStatusPendingApproval
	if err := repo.UpdateLeaveRequest(context.Background(), lr); err != nil {
		t.Fatalf("failed to seed pending status: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(context.Background(), lr.ID, "CANCELLED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notify call, got %d", len(notifier.calls))
	}
	if notifier.calls[0].notifType != "LEAVE_CANCELLED" {
		t.Errorf("expected notif type LEAVE_CANCELLED, got %s", notifier.calls[0].notifType)
	}
}

// UpdateLeaveRequestStatus is the manual generic status endpoint (§31) used
// e.g. by HR overriding a request directly, bypassing the approval
// push-callback path entirely. Notification coverage should be consistent
// regardless of which path drove the transition.
func TestService_UpdateLeaveRequestStatus_NotifiesEmployeeOnApproval(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	userID := uuid.New()
	createTestEmployeeAccount(db, empID, userID)

	lType := createTestLeaveType(repo)
	lr := createTestLeaveRequest(repo, empID, lType.ID)

	if _, err := svc.UpdateLeaveRequestStatus(context.Background(), lr.ID.String(), "APPROVED_FINAL", "approved by HR"); err != nil {
		t.Fatalf("UpdateLeaveRequestStatus failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notify call, got %d", len(notifier.calls))
	}
	call := notifier.calls[0]
	if call.recipientUserID != userID {
		t.Errorf("expected recipient %s, got %s", userID, call.recipientUserID)
	}
	if call.notifType != "LEAVE_APPROVED" {
		t.Errorf("expected notif type LEAVE_APPROVED, got %s", call.notifType)
	}
}

func TestService_UpdateLeaveRequestStatus_DoesNotNotifyWhenStatusUnchanged(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	createTestEmployeeAccount(db, empID, uuid.New())

	lType := createTestLeaveType(repo)
	lr := createTestLeaveRequest(repo, empID, lType.ID) // starts SUBMITTED

	if _, err := svc.UpdateLeaveRequestStatus(context.Background(), lr.ID.String(), "SUBMITTED", ""); err != nil {
		t.Fatalf("UpdateLeaveRequestStatus failed: %v", err)
	}

	if len(notifier.calls) != 0 {
		t.Fatalf("expected 0 notify calls for a no-op status transition, got %d", len(notifier.calls))
	}
}
