package approval

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func ctxAsUser(userID uuid.UUID) context.Context {
	return context.WithValue(context.Background(), "user_id", userID.String())
}

// TestService_CreateInstance_SupervisorHierarchy_DirectSupervisor validates the
// "atasan langsung" rule: a SUPERVISOR step with hierarchy_level=1 resolves to
// the platform user occupying the submitter's parent Organization.
func TestService_CreateInstance_SupervisorHierarchy_DirectSupervisor(t *testing.T) {
	svc, repo, db, cleanup := newTestServiceWithDB()
	defer cleanup()

	// Org tree: parentOrg -> childOrg (submitter's org)
	parentOrgID := uuid.New()
	childOrgID := uuid.New()
	seedOrganization(db, parentOrgID, nil)
	seedOrganization(db, childOrgID, &parentOrgID)

	submitterEmployeeID := uuid.New()
	submitterUserID := uuid.New()
	seedEmployment(db, submitterEmployeeID, childOrgID)
	seedEmployeeAccount(db, submitterEmployeeID, submitterUserID)

	supervisorEmployeeID := uuid.New()
	supervisorUserID := uuid.New()
	seedEmployment(db, supervisorEmployeeID, parentOrgID)
	seedEmployeeAccount(db, supervisorEmployeeID, supervisorUserID)

	flow := createTestFlow(repo, "leave")
	step := &ApprovalFlowStep{
		FlowID:            flow.ID,
		StepOrder:         1,
		StepName:          "Direct Supervisor",
		ApproverType:      ApproverTypeSupervisor,
		HierarchyLevel:    intPtr(1),
		ApprovalMode:      ApprovalModeAnyOne,
		ParticipationType: ParticipationTypeApprover,
		AllowReject:       true,
	}
	if err := repo.CreateStep(context.Background(), step); err != nil {
		t.Fatalf("failed to create step: %v", err)
	}

	resp, err := svc.CreateInstance(ctxAsUser(submitterUserID), CreateInstanceRequest{
		Module:     "leave",
		DocumentID: uuid.New().String(),
		FlowID:     flow.ID.String(),
	})
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}

	if len(resp.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(resp.Tasks))
	}
	if resp.Tasks[0].AssigneeID != supervisorUserID.String() {
		t.Errorf("expected assignee %s (direct supervisor), got %s", supervisorUserID, resp.Tasks[0].AssigneeID)
	}
}

// TestService_CreateInstance_SupervisorHierarchy_MultipleLevel validates
// "multiple level": hierarchy_level=2 skips the direct supervisor and resolves
// to the grandparent Organization's occupant.
func TestService_CreateInstance_SupervisorHierarchy_MultipleLevel(t *testing.T) {
	svc, repo, db, cleanup := newTestServiceWithDB()
	defer cleanup()

	grandparentOrgID := uuid.New()
	parentOrgID := uuid.New()
	childOrgID := uuid.New()
	seedOrganization(db, grandparentOrgID, nil)
	seedOrganization(db, parentOrgID, &grandparentOrgID)
	seedOrganization(db, childOrgID, &parentOrgID)

	submitterEmployeeID := uuid.New()
	submitterUserID := uuid.New()
	seedEmployment(db, submitterEmployeeID, childOrgID)
	seedEmployeeAccount(db, submitterEmployeeID, submitterUserID)

	level2EmployeeID := uuid.New()
	level2UserID := uuid.New()
	seedEmployment(db, level2EmployeeID, grandparentOrgID)
	seedEmployeeAccount(db, level2EmployeeID, level2UserID)

	flow := createTestFlow(repo, "leave")
	step := &ApprovalFlowStep{
		FlowID:            flow.ID,
		StepOrder:         1,
		StepName:          "Two Levels Up",
		ApproverType:      ApproverTypeSupervisor,
		HierarchyLevel:    intPtr(2),
		ApprovalMode:      ApprovalModeAnyOne,
		ParticipationType: ParticipationTypeApprover,
		AllowReject:       true,
	}
	if err := repo.CreateStep(context.Background(), step); err != nil {
		t.Fatalf("failed to create step: %v", err)
	}

	resp, err := svc.CreateInstance(ctxAsUser(submitterUserID), CreateInstanceRequest{
		Module:     "leave",
		DocumentID: uuid.New().String(),
		FlowID:     flow.ID.String(),
	})
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}

	if len(resp.Tasks) != 1 || resp.Tasks[0].AssigneeID != level2UserID.String() {
		t.Fatalf("expected task assigned to level-2 user %s, got %+v", level2UserID, resp.Tasks)
	}
}

// TestService_CreateInstance_OrganizationApprover_MultipleApprovers validates
// "bisa pilih organization" + "di satu level bisa lebih dari satu yang approve":
// an ORGANIZATION step targeting two explicit Organizations resolves to two tasks.
func TestService_CreateInstance_OrganizationApprover_MultipleApprovers(t *testing.T) {
	svc, repo, db, cleanup := newTestServiceWithDB()
	defer cleanup()

	orgAID := uuid.New()
	orgBID := uuid.New()
	submitterOrgID := uuid.New()
	seedOrganization(db, orgAID, nil)
	seedOrganization(db, orgBID, nil)
	seedOrganization(db, submitterOrgID, nil)

	empA, userA := uuid.New(), uuid.New()
	empB, userB := uuid.New(), uuid.New()
	seedEmployment(db, empA, orgAID)
	seedEmployeeAccount(db, empA, userA)
	seedEmployment(db, empB, orgBID)
	seedEmployeeAccount(db, empB, userB)

	submitterEmployeeID, submitterUserID := uuid.New(), uuid.New()
	seedEmployment(db, submitterEmployeeID, submitterOrgID)
	seedEmployeeAccount(db, submitterEmployeeID, submitterUserID)

	flow := createTestFlow(repo, "reimbursement")
	step := &ApprovalFlowStep{
		FlowID:            flow.ID,
		StepOrder:         1,
		StepName:          "Two Organizations",
		ApproverType:      ApproverTypeOrganization,
		ApprovalMode:      ApprovalModeAll,
		ParticipationType: ParticipationTypeApprover,
		AllowReject:       true,
	}
	if err := repo.CreateStep(context.Background(), step); err != nil {
		t.Fatalf("failed to create step: %v", err)
	}
	if err := repo.ReplaceStepOrganizations(context.Background(), step.ID, []uuid.UUID{orgAID, orgBID}); err != nil {
		t.Fatalf("failed to set step organizations: %v", err)
	}

	resp, err := svc.CreateInstance(ctxAsUser(submitterUserID), CreateInstanceRequest{
		Module:     "reimbursement",
		DocumentID: uuid.New().String(),
		FlowID:     flow.ID.String(),
	})
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}

	if len(resp.Tasks) != 2 {
		t.Fatalf("expected 2 tasks (one per organization), got %d: %+v", len(resp.Tasks), resp.Tasks)
	}
	got := map[string]bool{resp.Tasks[0].AssigneeID: true, resp.Tasks[1].AssigneeID: true}
	if !got[userA.String()] || !got[userB.String()] {
		t.Errorf("expected assignees %s and %s, got %+v", userA, userB, resp.Tasks)
	}
}

// TestService_CreateInstance_BothApprover_UnionAndDedup validates "bisa pilih
// keduanya": a BOTH step unions the hierarchy-resolved supervisor with the
// explicit Organization's occupant, deduplicating if they're the same person.
func TestService_CreateInstance_BothApprover_UnionAndDedup(t *testing.T) {
	svc, repo, db, cleanup := newTestServiceWithDB()
	defer cleanup()

	parentOrgID := uuid.New()
	childOrgID := uuid.New()
	seedOrganization(db, parentOrgID, nil)
	seedOrganization(db, childOrgID, &parentOrgID)

	submitterEmployeeID, submitterUserID := uuid.New(), uuid.New()
	seedEmployment(db, submitterEmployeeID, childOrgID)
	seedEmployeeAccount(db, submitterEmployeeID, submitterUserID)

	supervisorEmployeeID, supervisorUserID := uuid.New(), uuid.New()
	seedEmployment(db, supervisorEmployeeID, parentOrgID)
	seedEmployeeAccount(db, supervisorEmployeeID, supervisorUserID)

	flow := createTestFlow(repo, "leave")
	step := &ApprovalFlowStep{
		FlowID:            flow.ID,
		StepOrder:         1,
		StepName:          "Supervisor + Same Organization",
		ApproverType:      ApproverTypeBoth,
		HierarchyLevel:    intPtr(1),
		ApprovalMode:      ApprovalModeAll,
		ParticipationType: ParticipationTypeApprover,
		AllowReject:       true,
	}
	if err := repo.CreateStep(context.Background(), step); err != nil {
		t.Fatalf("failed to create step: %v", err)
	}
	// Explicit Organization is the SAME as the hierarchy-resolved parent org —
	// should dedup to a single task, not two, for the same supervisorUserID.
	if err := repo.ReplaceStepOrganizations(context.Background(), step.ID, []uuid.UUID{parentOrgID}); err != nil {
		t.Fatalf("failed to set step organizations: %v", err)
	}

	resp, err := svc.CreateInstance(ctxAsUser(submitterUserID), CreateInstanceRequest{
		Module:     "leave",
		DocumentID: uuid.New().String(),
		FlowID:     flow.ID.String(),
	})
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}

	if len(resp.Tasks) != 1 {
		t.Fatalf("expected 1 deduped task, got %d: %+v", len(resp.Tasks), resp.Tasks)
	}
	if resp.Tasks[0].AssigneeID != supervisorUserID.String() {
		t.Errorf("expected assignee %s, got %s", supervisorUserID, resp.Tasks[0].AssigneeID)
	}
}

// TestService_CreateInstance_WatcherStep_DoesNotBlock validates "hanya
// mengetahui": a WATCHER step's task is auto-completed and the instance lands
// directly on the next real (APPROVER) step instead of waiting on the watcher.
func TestService_CreateInstance_WatcherStep_DoesNotBlock(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")

	watcherUserID := uuid.New()
	watcherStep := &ApprovalFlowStep{
		FlowID:            flow.ID,
		StepOrder:         1,
		StepName:          "HR Informed",
		ApproverType:      ApproverTypeUser,
		ApproverUserID:    &watcherUserID,
		ApprovalMode:      ApprovalModeAnyOne,
		ParticipationType: ParticipationTypeWatcher,
		AllowReject:       true,
	}
	if err := repo.CreateStep(context.Background(), watcherStep); err != nil {
		t.Fatalf("failed to create watcher step: %v", err)
	}

	approverUserID := uuid.New()
	approverStep := &ApprovalFlowStep{
		FlowID:            flow.ID,
		StepOrder:         2,
		StepName:          "Manager Approval",
		ApproverType:      ApproverTypeUser,
		ApproverUserID:    &approverUserID,
		ApprovalMode:      ApprovalModeAnyOne,
		ParticipationType: ParticipationTypeApprover,
		AllowReject:       true,
	}
	if err := repo.CreateStep(context.Background(), approverStep); err != nil {
		t.Fatalf("failed to create approver step: %v", err)
	}

	resp, err := svc.CreateInstance(ctx(), CreateInstanceRequest{
		Module:     "leave",
		DocumentID: uuid.New().String(),
		FlowID:     flow.ID.String(),
	})
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}

	if resp.CurrentStep != 2 {
		t.Errorf("expected instance to land on step 2 (skipping the watcher step), got %d", resp.CurrentStep)
	}

	var watcherTask, approverTask *TaskResponse
	for i := range resp.Tasks {
		switch resp.Tasks[i].AssigneeID {
		case watcherUserID.String():
			watcherTask = &resp.Tasks[i]
		case approverUserID.String():
			approverTask = &resp.Tasks[i]
		}
	}
	if watcherTask == nil || watcherTask.Status != "DONE" {
		t.Errorf("expected watcher task auto-completed (DONE), got %+v", watcherTask)
	}
	if approverTask == nil || approverTask.Status != "PENDING" {
		t.Errorf("expected approver task PENDING, got %+v", approverTask)
	}

	// Watcher must not be able to actually approve/reject — they have no
	// pending task, so SubmitAction should reject the attempt.
	_, err = svc.SubmitAction(ctx(), resp.ID, watcherUserID.String(), SubmitActionRequest{Action: "APPROVE"})
	if err == nil {
		t.Error("expected error when a WATCHER attempts to submit an action, got nil")
	}
}
