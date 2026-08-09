package approval

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// capturingNotifier records every invocation of Notifier.Notify.
type capturingNotifier struct {
	calls []struct {
		recipientUserID uuid.UUID
		notifType       string
		params          []string
		referenceType   string
		referenceID     uuid.UUID
	}
}

func (c *capturingNotifier) Notify(ctx context.Context, recipientUserID uuid.UUID, notifType string, params []string, referenceType string, referenceID uuid.UUID) error {
	c.calls = append(c.calls, struct {
		recipientUserID uuid.UUID
		notifType       string
		params          []string
		referenceType   string
		referenceID     uuid.UUID
	}{recipientUserID, notifType, params, referenceType, referenceID})
	return nil
}

func (c *capturingNotifier) hasRecipient(uid uuid.UUID) bool {
	for _, call := range c.calls {
		if call.recipientUserID == uid {
			return true
		}
	}
	return false
}

// TestService_CreateInstance_NotifiesFirstStepApprover guards the root cause
// of the reported gap: assignees never got any "you have a pending approval"
// notification at all — the approval module never called Notify anywhere.
func TestService_CreateInstance_NotifiesFirstStepApprover(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	notifier := &capturingNotifier{}
	svc.SetNotifier(notifier)

	flow := createTestFlow(repo, "leave")
	approverID := uuidStr()
	_, err := svc.CreateStep(ctx(), flow.ID.String(), CreateStepRequest{
		StepName:       "Manager Approval",
		ApproverType:   "USER",
		ApproverUserID: &approverID,
	})
	if err != nil {
		t.Fatalf("CreateStep failed: %v", err)
	}

	docID := uuidStr()
	_, err = svc.CreateInstance(ctx(), CreateInstanceRequest{
		Module:     "leave",
		DocumentID: docID,
		FlowID:     flow.ID.String(),
	})
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(notifier.calls))
	}
	call := notifier.calls[0]
	if call.recipientUserID.String() != approverID {
		t.Errorf("expected recipient %s, got %s", approverID, call.recipientUserID)
	}
	if call.notifType != "APPROVAL_TASK_ASSIGNED" {
		t.Errorf("expected notif_type APPROVAL_TASK_ASSIGNED, got %q", call.notifType)
	}
	if call.referenceType != "leave" {
		t.Errorf("expected reference_type 'leave' (the owning module), got %q", call.referenceType)
	}
	if call.referenceID.String() != docID {
		t.Errorf("expected reference_id to be the leave document id %s, got %s", docID, call.referenceID)
	}
}

// TestService_CreateInstance_NotifiesWatcherStepDifferently guards that a
// WATCHER-step assignee still gets notified (they're a real recipient of a
// PENDING task now, per the earlier visibility fix) but with a distinct
// notif_type, since they have nothing to action.
func TestService_CreateInstance_NotifiesWatcherStepDifferently(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	notifier := &capturingNotifier{}
	svc.SetNotifier(notifier)

	flow := createTestFlow(repo, "leave")

	watcherID := uuid.New()
	watcherStep := &ApprovalFlowStep{
		FlowID:            flow.ID,
		StepOrder:         1,
		StepName:          "HR Informed",
		ApproverType:      ApproverTypeUser,
		ApproverUserID:    &watcherID,
		ApprovalMode:      ApprovalModeAnyOne,
		ParticipationType: ParticipationTypeWatcher,
		AllowReject:       true,
	}
	if err := repo.CreateStep(ctx(), watcherStep); err != nil {
		t.Fatalf("failed to create watcher step: %v", err)
	}

	approverID := uuid.New()
	approverStep := &ApprovalFlowStep{
		FlowID:            flow.ID,
		StepOrder:         2,
		StepName:          "Manager Approval",
		ApproverType:      ApproverTypeUser,
		ApproverUserID:    &approverID,
		ApprovalMode:      ApprovalModeAnyOne,
		ParticipationType: ParticipationTypeApprover,
		AllowReject:       true,
	}
	if err := repo.CreateStep(ctx(), approverStep); err != nil {
		t.Fatalf("failed to create approver step: %v", err)
	}

	_, err := svc.CreateInstance(ctx(), CreateInstanceRequest{
		Module:     "leave",
		DocumentID: uuidStr(),
		FlowID:     flow.ID.String(),
	})
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}

	if len(notifier.calls) != 2 {
		t.Fatalf("expected 2 notifications (watcher + approver), got %d", len(notifier.calls))
	}
	for _, call := range notifier.calls {
		switch call.recipientUserID {
		case watcherID:
			if call.notifType != "APPROVAL_WATCHER_ASSIGNED" {
				t.Errorf("expected watcher notif_type APPROVAL_WATCHER_ASSIGNED, got %q", call.notifType)
			}
		case approverID:
			if call.notifType != "APPROVAL_TASK_ASSIGNED" {
				t.Errorf("expected approver notif_type APPROVAL_TASK_ASSIGNED, got %q", call.notifType)
			}
		default:
			t.Errorf("unexpected notification recipient %s", call.recipientUserID)
		}
	}
}

// TestService_CreateInstance_NotifiesAllRoleMembers guards the other half of
// the reported gap: a ROLE-assigned step's task is stored against the role's
// own ID (one row, not one per member — see resolveStepAssignees), so naively
// treating assignee_id as a user id would either notify nobody or notify the
// wrong "user". Every actual member of the role must be notified.
func TestService_CreateInstance_NotifiesAllRoleMembers(t *testing.T) {
	svc, repo, db, cleanup := newTestServiceWithDB()
	defer cleanup()

	notifier := &capturingNotifier{}
	svc.SetNotifier(notifier)

	flow := createTestFlow(repo, "leave")
	roleID := uuid.New()
	step := &ApprovalFlowStep{
		FlowID:       flow.ID,
		StepOrder:    1,
		StepName:     "HR Role Approval",
		ApproverType: ApproverTypeRole,
		RoleID:       &roleID,
		ApprovalMode: ApprovalModeAnyOne,
		AllowReject:  true,
	}
	if err := repo.CreateStep(ctx(), step); err != nil {
		t.Fatalf("failed to create step: %v", err)
	}

	member1 := uuid.New()
	member2 := uuid.New()
	seedUserRole(db, member1, roleID)
	seedUserRole(db, member2, roleID)

	_, err := svc.CreateInstance(ctx(), CreateInstanceRequest{
		Module:     "leave",
		DocumentID: uuidStr(),
		FlowID:     flow.ID.String(),
	})
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}

	if len(notifier.calls) != 2 {
		t.Fatalf("expected 2 notifications (one per role member), got %d", len(notifier.calls))
	}
	if !notifier.hasRecipient(member1) {
		t.Errorf("expected member1 (%s) to be notified", member1)
	}
	if !notifier.hasRecipient(member2) {
		t.Errorf("expected member2 (%s) to be notified", member2)
	}
}

// TestService_SubmitAction_Approve_NotifiesNextStepAssignee guards that
// advancing to a subsequent step (not just the very first, at instance
// creation) also notifies the new step's assignee(s).
func TestService_SubmitAction_Approve_NotifiesNextStepAssignee(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	notifier := &capturingNotifier{}
	svc.SetNotifier(notifier)

	flow := createTestFlow(repo, "leave")
	approver1ID := uuidStr()
	approver2ID := uuidStr()

	if _, err := svc.CreateStep(ctx(), flow.ID.String(), CreateStepRequest{
		StepName:       "Step 1",
		ApproverType:   "USER",
		ApproverUserID: &approver1ID,
	}); err != nil {
		t.Fatalf("CreateStep (step1) failed: %v", err)
	}
	if _, err := svc.CreateStep(ctx(), flow.ID.String(), CreateStepRequest{
		StepName:       "Step 2",
		ApproverType:   "USER",
		ApproverUserID: &approver2ID,
	}); err != nil {
		t.Fatalf("CreateStep (step2) failed: %v", err)
	}

	inst, err := svc.CreateInstance(ctx(), CreateInstanceRequest{
		Module:     "leave",
		DocumentID: uuidStr(),
		FlowID:     flow.ID.String(),
	})
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}

	// Instance creation already notified approver1 — reset to isolate the
	// notification triggered by the SubmitAction advance below.
	notifier.calls = nil

	if _, err := svc.SubmitAction(ctx(), inst.ID, approver1ID, SubmitActionRequest{Action: "APPROVE"}); err != nil {
		t.Fatalf("SubmitAction failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notification for step2's assignee, got %d", len(notifier.calls))
	}
	if notifier.calls[0].recipientUserID.String() != approver2ID {
		t.Errorf("expected recipient %s, got %s", approver2ID, notifier.calls[0].recipientUserID)
	}
}

// TestService_CreateInstance_NoNotifier_NoPanic ensures the notifier
// remains fully optional — best-effort, matching every other Notifier
// integration in this codebase.
func TestService_CreateInstance_NoNotifier_NoPanic(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)

	if _, err := svc.CreateInstance(ctx(), CreateInstanceRequest{
		Module:     "leave",
		DocumentID: uuidStr(),
		FlowID:     flow.ID.String(),
	}); err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
}
