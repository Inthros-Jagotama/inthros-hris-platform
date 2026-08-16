package reimbursement

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

	rType := createTestReimbursementType(repo)
	rr := createTestReimbursementRequest(repo, empID, rType.ID)
	rr.Status = ReimbStatusPendingApproval
	if err := repo.UpdateReimbursementRequest(ctx(), rr); err != nil {
		t.Fatalf("failed to seed pending status: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), rr.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notify call, got %d", len(notifier.calls))
	}
	call := notifier.calls[0]
	if call.recipientUserID != userID {
		t.Errorf("expected recipient %s, got %s", userID, call.recipientUserID)
	}
	if call.notifType != "REIMBURSEMENT_APPROVED" {
		t.Errorf("expected notif type REIMBURSEMENT_APPROVED, got %s", call.notifType)
	}
	if call.referenceType != "reimbursement" || call.referenceID != rr.ID {
		t.Errorf("expected reference reimbursement/%s, got %s/%s", rr.ID, call.referenceType, call.referenceID)
	}
}

func TestService_HandleApprovalStatusChange_NotifiesEmployeeOnRejection(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	userID := uuid.New()
	createTestEmployeeAccount(db, empID, userID)

	rType := createTestReimbursementType(repo)
	rr := createTestReimbursementRequest(repo, empID, rType.ID)
	rr.Status = ReimbStatusPendingApproval
	if err := repo.UpdateReimbursementRequest(ctx(), rr); err != nil {
		t.Fatalf("failed to seed pending status: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), rr.ID, "REJECTED", "receipt missing"); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notify call, got %d", len(notifier.calls))
	}
	if notifier.calls[0].notifType != "REIMBURSEMENT_REJECTED" {
		t.Errorf("expected notif type REIMBURSEMENT_REJECTED, got %s", notifier.calls[0].notifType)
	}
}

func TestService_HandleApprovalStatusChange_SkipsNotifyWhenNoLinkedUserAccount(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New() // no employee_accounts row created
	rType := createTestReimbursementType(repo)
	rr := createTestReimbursementRequest(repo, empID, rType.ID)
	rr.Status = ReimbStatusPendingApproval
	if err := repo.UpdateReimbursementRequest(ctx(), rr); err != nil {
		t.Fatalf("failed to seed pending status: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), rr.ID, "APPROVED", ""); err != nil {
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

	rType := createTestReimbursementType(repo)
	rr := createTestReimbursementRequest(repo, empID, rType.ID)
	rr.Status = ReimbStatusPendingApproval
	if err := repo.UpdateReimbursementRequest(ctx(), rr); err != nil {
		t.Fatalf("failed to seed pending status: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), rr.ID, "APPROVED", ""); err != nil {
		t.Fatalf("expected HandleApprovalStatusChange to succeed despite notify failure, got: %v", err)
	}
	if len(notifier.calls) != 1 {
		t.Fatalf("expected the notifier to still be called once despite it returning an error, got %d calls", len(notifier.calls))
	}
}

// UpdateReimbursementRequestStatus is the manual generic status endpoint used
// e.g. by HR/Finance marking a request PAID directly. Notification coverage
// should be consistent regardless of which path drove the transition.
func TestService_UpdateReimbursementRequestStatus_NotifiesEmployeeOnPayment(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	userID := uuid.New()
	createTestEmployeeAccount(db, empID, userID)

	rType := createTestReimbursementType(repo)
	rr := createTestReimbursementRequest(repo, empID, rType.ID)
	rr.Status = ReimbStatusApproved
	if err := repo.UpdateReimbursementRequest(ctx(), rr); err != nil {
		t.Fatalf("failed to seed approved status: %v", err)
	}

	if _, err := svc.UpdateReimbursementRequestStatus(ctx(), rr.ID.String(), "PAID", "", nil, nil); err != nil {
		t.Fatalf("UpdateReimbursementRequestStatus failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notify call, got %d", len(notifier.calls))
	}
	call := notifier.calls[0]
	if call.recipientUserID != userID {
		t.Errorf("expected recipient %s, got %s", userID, call.recipientUserID)
	}
	if call.notifType != "REIMBURSEMENT_PAID" {
		t.Errorf("expected notif type REIMBURSEMENT_PAID, got %s", call.notifType)
	}
}

func TestService_UpdateReimbursementRequestStatus_DoesNotNotifyOnNonFinalTransitions(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	createTestEmployeeAccount(db, empID, uuid.New())

	rType := createTestReimbursementType(repo)
	rr := createTestReimbursementRequest(repo, empID, rType.ID) // DRAFT

	// Submit without a flow and without an active flow → SUBMITTED, no notify.
	if _, err := svc.UpdateReimbursementRequestStatus(ctx(), rr.ID.String(), "SUBMITTED", "", nil, nil); err != nil {
		t.Fatalf("UpdateReimbursementRequestStatus failed: %v", err)
	}
	// Cancel from SUBMITTED → CANCELLED, no outcome notify (only approved/
	// rejected/paid are notified per plan §6).
	if _, err := svc.UpdateReimbursementRequestStatus(ctx(), rr.ID.String(), "CANCELLED", "", nil, nil); err != nil {
		t.Fatalf("UpdateReimbursementRequestStatus failed: %v", err)
	}

	if len(notifier.calls) != 0 {
		t.Fatalf("expected 0 notify calls for non-final transitions, got %d", len(notifier.calls))
	}
}
