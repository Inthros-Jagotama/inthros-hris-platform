package attendance

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

func TestService_HandleApprovalStatusChange_Overtime_NotifiesEmployeeOnApproval(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	userID := uuid.New()
	createTestEmployeeAccount(db, empID, userID)

	o := createTestOvertimeRequest(repo, empID)
	o.Status = OvertimePendingApproval
	if err := repo.UpdateOvertimeRequest(ctx(), o); err != nil {
		t.Fatalf("failed to seed overtime request: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), o.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notify call, got %d", len(notifier.calls))
	}
	call := notifier.calls[0]
	if call.recipientUserID != userID {
		t.Errorf("expected recipient %s, got %s", userID, call.recipientUserID)
	}
	if call.notifType != "OVERTIME_APPROVED" {
		t.Errorf("expected notif type OVERTIME_APPROVED, got %s", call.notifType)
	}
	if call.referenceType != "attendance_overtime" || call.referenceID != o.ID {
		t.Errorf("expected reference attendance_overtime/%s, got %s/%s", o.ID, call.referenceType, call.referenceID)
	}
}

func TestService_HandleApprovalStatusChange_Overtime_SkipsNotifyWhenNoLinkedUserAccount(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New() // no employee_accounts row created
	o := createTestOvertimeRequest(repo, empID)
	o.Status = OvertimePendingApproval
	if err := repo.UpdateOvertimeRequest(ctx(), o); err != nil {
		t.Fatalf("failed to seed overtime request: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), o.ID, "REJECTED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	if len(notifier.calls) != 0 {
		t.Fatalf("expected 0 notify calls when employee has no linked user account, got %d", len(notifier.calls))
	}
}

func TestService_HandleApprovalStatusChange_Overtime_NotifyFailureDoesNotFailApproval(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()

	notifier := &fakeNotifier{err: context.DeadlineExceeded}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	createTestEmployeeAccount(db, empID, uuid.New())

	o := createTestOvertimeRequest(repo, empID)
	o.Status = OvertimePendingApproval
	if err := repo.UpdateOvertimeRequest(ctx(), o); err != nil {
		t.Fatalf("failed to seed overtime request: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), o.ID, "REJECTED", ""); err != nil {
		t.Fatalf("expected HandleApprovalStatusChange to succeed despite notify failure, got: %v", err)
	}
	if len(notifier.calls) != 1 {
		t.Fatalf("expected the notifier to still be called once despite it returning an error, got %d calls", len(notifier.calls))
	}
}

func TestService_HandleApprovalStatusChange_Correction_NotifiesEmployeeOnApproval(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	shift := createTestShift(repo) // 08:00 - 17:00
	empID := uuid.New()
	userID := uuid.New()
	createTestEmployeeAccount(db, empID, userID)
	createTestEmployeeShift(repo, empID, shift.ID)

	checkin := CreateEventRequest{
		EmployeeID:     empID.String(),
		EventType:      "CHECKIN",
		EventTimeUTC:   "2026-01-15T01:00:00Z",
		EventTimeLocal: "2026-01-15T08:00:00+07:00",
		Latitude:       -6.2088,
		Longitude:      106.8456,
	}
	if _, err := svc.CreateEvent(ctx(), checkin); err != nil {
		t.Fatalf("checkin CreateEvent failed: %v", err)
	}
	session, err := repo.FindSessionByEmployeeAndDate(ctx(), empID, "2026-01-15")
	if err != nil {
		t.Fatalf("expected session: %v", err)
	}

	checkoutTime := "2026-01-15T17:00:00+07:00"
	c := &AttendanceCorrectionRequest{
		EmployeeID:          empID,
		AttendanceSessionID: session.ID,
		CorrectionType:      CorrectionTypeMissingCheckout,
		RequestedCheckout:   timePtr(parseTime(checkoutTime)),
		Reason:              "Forgot to check out",
		Status:              CorrectionPendingApproval,
	}
	if err := repo.CreateCorrectionRequest(ctx(), c); err != nil {
		t.Fatalf("failed to seed correction request: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), c.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notify call, got %d", len(notifier.calls))
	}
	call := notifier.calls[0]
	if call.recipientUserID != userID {
		t.Errorf("expected recipient %s, got %s", userID, call.recipientUserID)
	}
	if call.notifType != "CORRECTION_APPROVED" {
		t.Errorf("expected notif type CORRECTION_APPROVED, got %s", call.notifType)
	}
	if call.referenceType != "attendance_correction" || call.referenceID != c.ID {
		t.Errorf("expected reference attendance_correction/%s, got %s/%s", c.ID, call.referenceType, call.referenceID)
	}
}
