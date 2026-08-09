package approval

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// capturingHandler records every invocation of a StatusChangeHandler.
type capturingHandler struct {
	calls []struct {
		documentID uuid.UUID
		status     InstanceStatus
		note       string
	}
}

func (c *capturingHandler) handle(ctx context.Context, documentID uuid.UUID, status InstanceStatus, note string) error {
	c.calls = append(c.calls, struct {
		documentID uuid.UUID
		status     InstanceStatus
		note       string
	}{documentID, status, note})
	return nil
}

func TestService_SubmitAction_Approve_NotifiesStatusHandler(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	captured := &capturingHandler{}
	svc.RegisterStatusHandler("leave", captured.handle)

	flow := createTestFlow(repo, "leave")
	step := createTestStep(repo, flow.ID, 1)
	docID := uuid.New()
	inst := createTestInstance(repo, flow, docID)
	createTestTask(repo, inst.ID, 1, *step.ApproverUserID)

	_, err := svc.SubmitAction(ctx(), inst.ID.String(), step.ApproverUserID.String(), SubmitActionRequest{Action: "APPROVE"})
	if err != nil {
		t.Fatalf("SubmitAction failed: %v", err)
	}

	if len(captured.calls) != 1 {
		t.Fatalf("expected 1 status handler call, got %d", len(captured.calls))
	}
	if captured.calls[0].documentID != docID {
		t.Errorf("expected documentID %s, got %s", docID, captured.calls[0].documentID)
	}
	if captured.calls[0].status != InstanceStatusApproved {
		t.Errorf("expected status APPROVED, got %s", captured.calls[0].status)
	}
}

// TestService_SubmitAction_Approve_WithNote_ForwardsNoteToStatusHandler guards
// the fix for a note entered on the final APPROVE action being silently
// discarded — SubmitAction previously always notified with a hardcoded empty
// note on full approval, dropping whatever the approver actually wrote.
func TestService_SubmitAction_Approve_WithNote_ForwardsNoteToStatusHandler(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	captured := &capturingHandler{}
	svc.RegisterStatusHandler("leave", captured.handle)

	flow := createTestFlow(repo, "leave")
	step := createTestStep(repo, flow.ID, 1)
	docID := uuid.New()
	inst := createTestInstance(repo, flow, docID)
	createTestTask(repo, inst.ID, 1, *step.ApproverUserID)

	note := "approved, enjoy your leave"
	_, err := svc.SubmitAction(ctx(), inst.ID.String(), step.ApproverUserID.String(), SubmitActionRequest{Action: "APPROVE", Note: &note})
	if err != nil {
		t.Fatalf("SubmitAction failed: %v", err)
	}

	if len(captured.calls) != 1 {
		t.Fatalf("expected 1 status handler call, got %d", len(captured.calls))
	}
	if captured.calls[0].note != note {
		t.Errorf("expected note %q, got %q", note, captured.calls[0].note)
	}
}

// TestService_SubmitAction_Approve_IntermediateStepNote_NotifiesStatusHandlerAsPending
// guards a related gap: in a multi-step flow, a note entered when approving
// a NON-final step (the instance stays PENDING, more steps remain) used to
// be dropped entirely — notifyStatusChange was only ever called when the
// instance reached a terminal state. Now an intermediate approve with a note
// still notifies (status stays PENDING) so the note isn't lost while waiting
// for the rest of the flow to resolve.
func TestService_SubmitAction_Approve_IntermediateStepNote_NotifiesStatusHandlerAsPending(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	captured := &capturingHandler{}
	svc.RegisterStatusHandler("leave", captured.handle)

	flow := createTestFlow(repo, "leave")
	step1 := createTestStep(repo, flow.ID, 1)
	createTestStep(repo, flow.ID, 2)

	docID := uuid.New()
	inst := createTestInstance(repo, flow, docID)
	createTestTask(repo, inst.ID, 1, *step1.ApproverUserID)

	note := "looks fine so far, forwarding to next approver"
	_, err := svc.SubmitAction(ctx(), inst.ID.String(), step1.ApproverUserID.String(), SubmitActionRequest{Action: "APPROVE", Note: &note})
	if err != nil {
		t.Fatalf("SubmitAction failed: %v", err)
	}

	if len(captured.calls) != 1 {
		t.Fatalf("expected 1 status handler call for the intermediate approval, got %d", len(captured.calls))
	}
	if captured.calls[0].status != InstanceStatusPending {
		t.Errorf("expected status PENDING (flow not finalized), got %s", captured.calls[0].status)
	}
	if captured.calls[0].note != note {
		t.Errorf("expected note %q, got %q", note, captured.calls[0].note)
	}
}

// TestService_SubmitAction_Approve_IntermediateStepNoNote_DoesNotNotify
// ensures the intermediate-step notification only fires when there's
// actually a note to carry — an intermediate approve without a note has
// nothing new to tell consumer modules, so it stays silent as before.
func TestService_SubmitAction_Approve_IntermediateStepNoNote_DoesNotNotify(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	captured := &capturingHandler{}
	svc.RegisterStatusHandler("leave", captured.handle)

	flow := createTestFlow(repo, "leave")
	step1 := createTestStep(repo, flow.ID, 1)
	createTestStep(repo, flow.ID, 2)

	docID := uuid.New()
	inst := createTestInstance(repo, flow, docID)
	createTestTask(repo, inst.ID, 1, *step1.ApproverUserID)

	_, err := svc.SubmitAction(ctx(), inst.ID.String(), step1.ApproverUserID.String(), SubmitActionRequest{Action: "APPROVE"})
	if err != nil {
		t.Fatalf("SubmitAction failed: %v", err)
	}

	if len(captured.calls) != 0 {
		t.Errorf("expected no status handler calls for a note-less intermediate approval, got %d", len(captured.calls))
	}
}

func TestService_SubmitAction_Reject_NotifiesStatusHandler(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	captured := &capturingHandler{}
	svc.RegisterStatusHandler("leave", captured.handle)

	flow := createTestFlow(repo, "leave")
	step := createTestStep(repo, flow.ID, 1)
	docID := uuid.New()
	inst := createTestInstance(repo, flow, docID)
	createTestTask(repo, inst.ID, 1, *step.ApproverUserID)

	note := "insufficient balance"
	_, err := svc.SubmitAction(ctx(), inst.ID.String(), step.ApproverUserID.String(), SubmitActionRequest{Action: "REJECT", Note: &note})
	if err != nil {
		t.Fatalf("SubmitAction failed: %v", err)
	}

	if len(captured.calls) != 1 {
		t.Fatalf("expected 1 status handler call, got %d", len(captured.calls))
	}
	if captured.calls[0].status != InstanceStatusRejected {
		t.Errorf("expected status REJECTED, got %s", captured.calls[0].status)
	}
	if captured.calls[0].note != note {
		t.Errorf("expected note %q, got %q", note, captured.calls[0].note)
	}
}

func TestService_CancelInstance_NotifiesStatusHandler(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	captured := &capturingHandler{}
	svc.RegisterStatusHandler("reimbursement", captured.handle)

	flow := createTestFlow(repo, "reimbursement")
	docID := uuid.New()
	inst := createTestInstance(repo, flow, docID)

	if err := svc.CancelInstance(ctx(), inst.ID.String()); err != nil {
		t.Fatalf("CancelInstance failed: %v", err)
	}

	if len(captured.calls) != 1 {
		t.Fatalf("expected 1 status handler call, got %d", len(captured.calls))
	}
	if captured.calls[0].status != InstanceStatusCancelled {
		t.Errorf("expected status CANCELLED, got %s", captured.calls[0].status)
	}
}

// TestService_CreateInstance_WatcherOnlyFlow_NotifiesApproved validates that
// a flow made entirely of WATCHER steps auto-approves at creation time and
// still fires the status handler (the consumer module never gets a chance to
// call anything else — creation itself is the final state).
func TestService_CreateInstance_WatcherOnlyFlow_NotifiesApproved(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	captured := &capturingHandler{}
	svc.RegisterStatusHandler("leave", captured.handle)

	flow := createTestFlow(repo, "leave")
	watcherUserID := uuid.New()
	step := &ApprovalFlowStep{
		FlowID:            flow.ID,
		StepOrder:         1,
		StepName:          "FYI Only",
		ApproverType:      ApproverTypeUser,
		ApproverUserID:    &watcherUserID,
		ApprovalMode:      ApprovalModeAnyOne,
		ParticipationType: ParticipationTypeWatcher,
		AllowReject:       true,
	}
	if err := repo.CreateStep(context.Background(), step); err != nil {
		t.Fatalf("failed to create step: %v", err)
	}

	docID := uuid.New()
	resp, err := svc.CreateInstance(ctx(), CreateInstanceRequest{
		Module:     "leave",
		DocumentID: docID.String(),
		FlowID:     flow.ID.String(),
	})
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	if resp.Status != "APPROVED" {
		t.Fatalf("expected instance auto-approved, got status %s", resp.Status)
	}

	if len(captured.calls) != 1 {
		t.Fatalf("expected 1 status handler call, got %d", len(captured.calls))
	}
	if captured.calls[0].documentID != docID || captured.calls[0].status != InstanceStatusApproved {
		t.Errorf("unexpected call: %+v", captured.calls[0])
	}
}

func TestService_SubmitAction_NoHandlerRegistered_NoPanic(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	// Deliberately do NOT register a handler for this module.
	flow := createTestFlow(repo, "leave")
	step := createTestStep(repo, flow.ID, 1)
	inst := createTestInstance(repo, flow, uuid.New())
	createTestTask(repo, inst.ID, 1, *step.ApproverUserID)

	_, err := svc.SubmitAction(ctx(), inst.ID.String(), step.ApproverUserID.String(), SubmitActionRequest{Action: "APPROVE"})
	if err != nil {
		t.Fatalf("SubmitAction failed: %v", err)
	}
}
