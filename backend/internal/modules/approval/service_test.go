package approval

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func ctx() context.Context {
	return context.Background()
}

// =========================================================================
// Approval Flow Service Tests
// =========================================================================

func TestService_CreateFlow_Success(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	req := CreateFlowRequest{
		Module:  "leave",
		Name:    "Leave Approval Flow",
		Version: 1,
	}

	resp, err := svc.CreateFlow(ctx(), req)
	if err != nil {
		t.Fatalf("CreateFlow failed: %v", err)
	}

	if resp.Module != "leave" {
		t.Errorf("expected module 'leave', got '%s'", resp.Module)
	}
	if !resp.IsActive {
		t.Error("expected default is_active = true")
	}
}

func TestService_CreateFlow_DefaultVersion(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	req := CreateFlowRequest{
		Module: "overtime",
		Name:   "Overtime Approval Flow",
		// Version = 0 → should default to 1
	}

	resp, err := svc.CreateFlow(ctx(), req)
	if err != nil {
		t.Fatalf("CreateFlow failed: %v", err)
	}

	if resp.Version != 1 {
		t.Errorf("expected default version 1, got %d", resp.Version)
	}
}

func TestService_GetFlowByID_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	created := createTestFlow(repo, "leave")

	found, err := svc.GetFlowByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetFlowByID failed: %v", err)
	}

	if found.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), found.ID)
	}
}

func TestService_GetFlowByID_InvalidUUID(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetFlowByID(ctx(), "not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_GetFlowByID_NotFound(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetFlowByID(ctx(), uuidStr())
	if err == nil {
		t.Fatal("expected error for non-existent flow")
	}
}

func TestService_ListFlows_DefaultPagination(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	for i := 0; i < 3; i++ {
		createTestFlow(repo, "leave")
	}

	resp, err := svc.ListFlows(ctx(), 0, 0)
	if err != nil {
		t.Fatalf("ListFlows failed: %v", err)
	}

	if resp.Page != 1 {
		t.Errorf("expected page 1, got %d", resp.Page)
	}
	if resp.PerPage != 20 {
		t.Errorf("expected per_page 20 (default), got %d", resp.PerPage)
	}
	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}
}

func TestService_UpdateFlow_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	created := createTestFlow(repo, "leave")

	updated, err := svc.UpdateFlow(ctx(), created.ID.String(), UpdateFlowRequest{
		Name: strPtr("Updated Flow Name"),
	})
	if err != nil {
		t.Fatalf("UpdateFlow failed: %v", err)
	}

	if updated.Name != "Updated Flow Name" {
		t.Errorf("expected name 'Updated Flow Name', got '%s'", updated.Name)
	}
}

func TestService_DeleteFlow_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	created := createTestFlow(repo, "leave")

	if err := svc.DeleteFlow(ctx(), created.ID.String()); err != nil {
		t.Fatalf("DeleteFlow failed: %v", err)
	}

	// Should not be found after soft delete
	_, err := svc.GetFlowByID(ctx(), created.ID.String())
	if err == nil {
		t.Fatal("expected error after deleting flow")
	}
}

// =========================================================================
// Approval Flow Step Service Tests
// =========================================================================

func TestService_CreateStep_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	approverID := uuidStr()

	resp, err := svc.CreateStep(ctx(), flow.ID.String(), CreateStepRequest{
		StepName:       "Manager Approval",
		ApproverType:   "USER",
		ApproverUserID: &approverID,
	})
	if err != nil {
		t.Fatalf("CreateStep failed: %v", err)
	}

	if resp.StepOrder != 1 {
		t.Errorf("expected step_order 1, got %d", resp.StepOrder)
	}
	if resp.StepName != "Manager Approval" {
		t.Errorf("expected step_name 'Manager Approval', got '%s'", resp.StepName)
	}
	if resp.ApprovalMode != "ANY_ONE" {
		t.Errorf("expected default approval_mode 'ANY_ONE', got '%s'", resp.ApprovalMode)
	}
}

func TestService_ListStepsByFlow_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)
	createTestStep(repo, flow.ID, 2)

	steps, err := svc.ListStepsByFlow(ctx(), flow.ID.String())
	if err != nil {
		t.Fatalf("ListStepsByFlow failed: %v", err)
	}

	if len(steps) != 2 {
		t.Errorf("expected 2 steps, got %d", len(steps))
	}
}

func TestService_UpdateStep_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	step := createTestStep(repo, flow.ID, 1)

	updated, err := svc.UpdateStep(ctx(), flow.ID.String(), step.ID.String(), UpdateStepRequest{
		StepName: strPtr("Updated Step Name"),
	})
	if err != nil {
		t.Fatalf("UpdateStep failed: %v", err)
	}

	if updated.StepName != "Updated Step Name" {
		t.Errorf("expected step_name 'Updated Step Name', got '%s'", updated.StepName)
	}
}

func TestService_DeleteStep_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	step := createTestStep(repo, flow.ID, 1)

	if err := svc.DeleteStep(ctx(), flow.ID.String(), step.ID.String()); err != nil {
		t.Fatalf("DeleteStep failed: %v", err)
	}

	// Step should not be visible after soft delete
	steps, _ := svc.ListStepsByFlow(ctx(), flow.ID.String())
	if len(steps) != 0 {
		t.Errorf("expected 0 steps after delete, got %d", len(steps))
	}
}

// =========================================================================
// Approval Instance Service Tests
// =========================================================================

func TestService_CreateInstance_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	step := createTestStep(repo, flow.ID, 1)

	docID := uuidStr()
	resp, err := svc.CreateInstance(ctx(), CreateInstanceRequest{
		Module:     "leave",
		DocumentID: docID,
		FlowID:     flow.ID.String(),
	})
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}

	if resp.Module != "leave" {
		t.Errorf("expected module 'leave', got '%s'", resp.Module)
	}
	if resp.Status != "PENDING" {
		t.Errorf("expected status 'PENDING', got '%s'", resp.Status)
	}
	if resp.CurrentStep != 1 {
		t.Errorf("expected current_step 1, got %d", resp.CurrentStep)
	}
	// Step should have been resolved, so first step approver should be set
	if len(resp.Tasks) == 0 {
		t.Error("expected tasks to be created for the first step")
	}
	_ = step
}

func TestService_CreateInstance_DuplicateDocument_Error(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)

	docID := uuidStr()

	// First instance succeeds
	_, err := svc.CreateInstance(ctx(), CreateInstanceRequest{
		Module:     "leave",
		DocumentID: docID,
		FlowID:     flow.ID.String(),
	})
	if err != nil {
		t.Fatalf("First CreateInstance failed: %v", err)
	}

	// Second instance for same document should fail
	_, err = svc.CreateInstance(ctx(), CreateInstanceRequest{
		Module:     "leave",
		DocumentID: docID,
		FlowID:     flow.ID.String(),
	})
	if err == nil {
		t.Fatal("expected error for duplicate document")
	}
}

func TestService_CreateInstance_InactiveFlow_Error(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	// Deactivate flow
	flow.IsActive = false
	repo.UpdateFlow(ctx(), flow)

	_, err := svc.CreateInstance(ctx(), CreateInstanceRequest{
		Module:     "leave",
		DocumentID: uuidStr(),
		FlowID:     flow.ID.String(),
	})
	if err == nil {
		t.Fatal("expected error when creating instance with inactive flow")
	}
}

func TestService_CreateInstance_ZeroAssignees_Error(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	// ORGANIZATION step with no Organizations linked resolves to zero
	// assignees — must not silently create an instance nobody can act on.
	step := &ApprovalFlowStep{
		FlowID:       flow.ID,
		StepOrder:    1,
		StepName:     "Unreachable Step",
		ApproverType: ApproverTypeOrganization,
		ApprovalMode: ApprovalModeAnyOne,
		AllowReject:  true,
	}
	if err := repo.CreateStep(ctx(), step); err != nil {
		t.Fatalf("failed to create step: %v", err)
	}

	_, err := svc.CreateInstance(ctx(), CreateInstanceRequest{
		Module:     "leave",
		DocumentID: uuidStr(),
		FlowID:     flow.ID.String(),
	})
	if err == nil {
		t.Fatal("expected error when the landed step resolves to zero assignees")
	}
}

func TestService_GetInstanceByID_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)

	docID := uuid.New()
	inst := createTestInstance(repo, flow, docID)

	found, err := svc.GetInstanceByID(ctx(), inst.ID.String())
	if err != nil {
		t.Fatalf("GetInstanceByID failed: %v", err)
	}

	if found.ID != inst.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", inst.ID.String(), found.ID)
	}
}

// TestService_GetInstanceByID_EnrichesActionsWithActorInfo guards the
// approval-preview requirement: once a step has been approved, its action
// entry must carry the approver's name/employee code/organization so the FE
// can show "who approved this step" instead of just a raw actor_user_id.
func TestService_GetInstanceByID_EnrichesActionsWithActorInfo(t *testing.T) {
	svc, repo, db, cleanup := newTestServiceWithDB()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)

	inst := createTestInstance(repo, flow, uuid.New())

	approverID := uuid.New()
	createTestTask(repo, inst.ID, 1, approverID)

	orgID := uuid.New()
	seedOrganizationNamed(db, orgID, nil, "Finance Department")
	employeeID := uuid.New()
	seedEmployee(db, employeeID, "Budi Santoso", "EMP-042")
	seedEmployment(db, employeeID, orgID)
	seedEmployeeAccount(db, employeeID, approverID)

	if _, err := svc.SubmitAction(ctx(), inst.ID.String(), approverID.String(), SubmitActionRequest{
		Action: "APPROVE",
	}); err != nil {
		t.Fatalf("SubmitAction failed: %v", err)
	}

	found, err := svc.GetInstanceByID(ctx(), inst.ID.String())
	if err != nil {
		t.Fatalf("GetInstanceByID failed: %v", err)
	}
	if len(found.Actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(found.Actions))
	}
	action := found.Actions[0]
	if action.ActorName != "Budi Santoso" {
		t.Errorf("expected actor_name 'Budi Santoso', got %q", action.ActorName)
	}
	if action.ActorEmployeeCode != "EMP-042" {
		t.Errorf("expected actor_employee_code 'EMP-042', got %q", action.ActorEmployeeCode)
	}
	if action.ActorOrganizationName != "Finance Department" {
		t.Errorf("expected actor_organization_name 'Finance Department', got %q", action.ActorOrganizationName)
	}
}

func TestService_ListInstances_DefaultPagination(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	for i := 0; i < 3; i++ {
		createTestInstance(repo, flow, uuid.New())
	}

	resp, err := svc.ListInstances(ctx(), 0, 0, "", "")
	if err != nil {
		t.Fatalf("ListInstances failed: %v", err)
	}

	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}
}

func TestService_CancelInstance_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	inst := createTestInstance(repo, flow, uuid.New())

	if err := svc.CancelInstance(ctx(), inst.ID.String()); err != nil {
		t.Fatalf("CancelInstance failed: %v", err)
	}

	found, _ := svc.GetInstanceByID(ctx(), inst.ID.String())
	if found.Status != "CANCELLED" {
		t.Errorf("expected status 'CANCELLED', got '%s'", found.Status)
	}
}

// =========================================================================
// Approval Action Service Tests
// =========================================================================

func TestService_SubmitAction_Approve_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	// Setup: flow with 1 step
	flow := createTestFlow(repo, "leave")
	step := createTestStep(repo, flow.ID, 1)

	// Create instance
	docID := uuid.New()
	inst := createTestInstance(repo, flow, docID)

	// Manually create a task for the actor
	actorID := uuid.New()
	createTestTask(repo, inst.ID, 1, actorID)

	// Submit approve action
	resp, err := svc.SubmitAction(ctx(), inst.ID.String(), actorID.String(), SubmitActionRequest{
		Action: "APPROVE",
	})
	if err != nil {
		t.Fatalf("SubmitAction failed: %v", err)
	}

	if resp.Status != "APPROVED" {
		t.Errorf("expected status 'APPROVED', got '%s'", resp.Status)
	}
	_ = step
}

func TestService_SubmitAction_RoleAssignee_Success(t *testing.T) {
	svc, repo, db, cleanup := newTestServiceWithDB()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)

	// Task di-assign ke ROLE (bukan user langsung) — pola alur dengan
	// approver_type ROLE (mis. role SDM).
	roleID := uuid.New()
	inst := createTestInstance(repo, flow, uuid.New())
	createTestTaskROLE(repo, inst.ID, 1, roleID)

	// Actor punya role tersebut → berhak approve.
	actorID := uuid.New()
	seedUserRole(db, actorID, roleID)

	resp, err := svc.SubmitAction(ctx(), inst.ID.String(), actorID.String(), SubmitActionRequest{
		Action: "APPROVE",
	})
	if err != nil {
		t.Fatalf("SubmitAction for ROLE-assigned task failed: %v", err)
	}
	if resp.Status != "APPROVED" {
		t.Errorf("expected status 'APPROVED', got '%s'", resp.Status)
	}

	// Actor TANPA role tersebut → tetap ditolak.
	outsiderID := uuid.New()
	inst2 := createTestInstance(repo, flow, uuid.New())
	createTestTaskROLE(repo, inst2.ID, 1, roleID)
	if _, err := svc.SubmitAction(ctx(), inst2.ID.String(), outsiderID.String(), SubmitActionRequest{
		Action: "APPROVE",
	}); err == nil {
		t.Fatal("expected error for user without the role, got nil")
	}
}

func TestService_SubmitAction_Reject_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)

	inst := createTestInstance(repo, flow, uuid.New())

	actorID := uuid.New()
	createTestTask(repo, inst.ID, 1, actorID)

	note := "Not approved"
	resp, err := svc.SubmitAction(ctx(), inst.ID.String(), actorID.String(), SubmitActionRequest{
		Action: "REJECT",
		Note:   &note,
	})
	if err != nil {
		t.Fatalf("SubmitAction failed: %v", err)
	}

	if resp.Status != "REJECTED" {
		t.Errorf("expected status 'REJECTED', got '%s'", resp.Status)
	}
}

func TestService_SubmitAction_AllowRejectFalse_Error(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	approverID := uuidStr()

	// Create step with AllowReject=true (default), then update to false via service
	// Note: We do it this way because GORM skips false (zero value) for bool fields on create
	step, err := svc.CreateStep(ctx(), flow.ID.String(), CreateStepRequest{
		StepName:       "Manager Approval",
		ApproverType:   "USER",
		ApproverUserID: &approverID,
	})
	if err != nil {
		t.Fatalf("CreateStep failed: %v", err)
	}

	allowReject := false
	_, err = svc.UpdateStep(ctx(), flow.ID.String(), step.ID, UpdateStepRequest{
		AllowReject: &allowReject,
	})
	if err != nil {
		t.Fatalf("UpdateStep failed to set AllowReject=false: %v", err)
	}

	// Create instance via service so tasks are properly set up
	inst, err := svc.CreateInstance(ctx(), CreateInstanceRequest{
		Module:     "leave",
		DocumentID: uuidStr(),
		FlowID:     flow.ID.String(),
	})
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}

	// The actor should be the approver_user_id from the step
	_, err = svc.SubmitAction(ctx(), inst.ID, approverID, SubmitActionRequest{
		Action: "REJECT",
	})
	if err == nil {
		t.Fatal("expected error when AllowReject is false")
	}
}

func TestService_SubmitAction_NoTask_Error(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)

	inst := createTestInstance(repo, flow, uuid.New())

	// Actor without a task
	actorID := uuid.New()

	_, err := svc.SubmitAction(ctx(), inst.ID.String(), actorID.String(), SubmitActionRequest{
		Action: "APPROVE",
	})
	if err == nil {
		t.Fatal("expected error when user has no pending task")
	}
}

func TestService_SubmitAction_NonPendingInstance_Error(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)

	inst := createTestInstance(repo, flow, uuid.New())
	inst.Status = InstanceStatusApproved
	repo.UpdateInstance(ctx(), inst)

	actorID := uuid.New()

	_, err := svc.SubmitAction(ctx(), inst.ID.String(), actorID.String(), SubmitActionRequest{
		Action: "APPROVE",
	})
	if err == nil {
		t.Fatal("expected error when instance is not PENDING")
	}
}

// =========================================================================
// Approval Task Service Tests
// =========================================================================

func TestService_ListMyPendingTasks_Success(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)

	actorID := uuid.New()

	// Create instances with tasks for this actor
	inst1 := createTestInstance(repo, flow, uuid.New())
	createTestTask(repo, inst1.ID, 1, actorID)

	inst2 := createTestInstance(repo, flow, uuid.New())
	createTestTask(repo, inst2.ID, 1, actorID)

	// Create a task for another actor (should not be counted)
	inst3 := createTestInstance(repo, flow, uuid.New())
	createTestTask(repo, inst3.ID, 1, uuid.New())

	resp, err := svc.ListMyPendingTasks(ctx(), actorID.String(), 1, 10, "", nil)
	if err != nil {
		t.Fatalf("ListMyPendingTasks failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2 tasks, got %d", resp.Total)
	}
}

func TestService_ListMyPendingTasks_FiltersByStatusAndFlow(t *testing.T) {
	svc, repo, db, cleanup := newTestServiceWithDB()
	defer cleanup()

	flow1 := createTestFlow(repo, "leave")
	createTestStep(repo, flow1.ID, 1)
	flow2 := createTestFlow(repo, "payroll")
	createTestStep(repo, flow2.ID, 1)

	actorID := uuid.New()

	// Instance PENDING di flow1
	inst1 := createTestInstance(repo, flow1, uuid.New())
	createTestTask(repo, inst1.ID, 1, actorID)

	// Instance PENDING di flow2
	inst2 := createTestInstance(repo, flow2, uuid.New())
	createTestTask(repo, inst2.ID, 1, actorID)

	// Instance APPROVED di flow1 (task tetap pending — watcher/progress case)
	inst3 := createTestInstance(repo, flow1, uuid.New())
	inst3.Status = InstanceStatusApproved
	if err := db.Save(inst3).Error; err != nil {
		t.Fatalf("failed to set instance APPROVED: %v", err)
	}
	createTestTask(repo, inst3.ID, 1, actorID)

	// Filter by flow1 → 2 tasks (inst1 PENDING + inst3 APPROVED)
	resp, err := svc.ListMyPendingTasks(ctx(), actorID.String(), 1, 10, "", &flow1.ID)
	if err != nil {
		t.Fatalf("ListMyPendingTasks by flow failed: %v", err)
	}
	if resp.Total != 2 {
		t.Errorf("expected 2 tasks for flow1, got %d", resp.Total)
	}

	// Filter by status APPROVED → 1 task (inst3)
	resp, err = svc.ListMyPendingTasks(ctx(), actorID.String(), 1, 10, "APPROVED", nil)
	if err != nil {
		t.Fatalf("ListMyPendingTasks by status failed: %v", err)
	}
	if resp.Total != 1 {
		t.Errorf("expected 1 APPROVED task, got %d", resp.Total)
	}

	// Filter by status + flow (APPROVED + flow1) → 1 task
	resp, err = svc.ListMyPendingTasks(ctx(), actorID.String(), 1, 10, "APPROVED", &flow1.ID)
	if err != nil {
		t.Fatalf("ListMyPendingTasks by status+flow failed: %v", err)
	}
	if resp.Total != 1 {
		t.Errorf("expected 1 APPROVED task in flow1, got %d", resp.Total)
	}

	// Filter by status REJECTED → 0 tasks
	resp, err = svc.ListMyPendingTasks(ctx(), actorID.String(), 1, 10, "REJECTED", nil)
	if err != nil {
		t.Fatalf("ListMyPendingTasks by rejected status failed: %v", err)
	}
	if resp.Total != 0 {
		t.Errorf("expected 0 REJECTED tasks, got %d", resp.Total)
	}
}

func TestService_ListMyDoneTasks_OnlyProcessedTasks(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)

	actorID := uuid.New()

	// Task PENDING (belum diproses) — tidak masuk riwayat.
	instPending := createTestInstance(repo, flow, uuid.New())
	createTestTask(repo, instPending.ID, 1, actorID)

	// Task DONE (sudah diproses) — masuk riwayat.
	instDone := createTestInstance(repo, flow, uuid.New())
	createTestTask(repo, instDone.ID, 1, actorID)
	doneTasks, err := repo.FindTasksByInstanceID(ctx(), instDone.ID)
	if err != nil || len(doneTasks) != 1 {
		t.Fatalf("failed to load done task: %v", err)
	}
	doneTask := doneTasks[0]
	if err := svc.repo.UpdateTaskStatus(ctx(), doneTask.ID, TaskStatusDone); err != nil {
		t.Fatalf("failed to mark task done: %v", err)
	}

	resp, err := svc.ListMyDoneTasks(ctx(), actorID.String(), 1, 10, nil)
	if err != nil {
		t.Fatalf("ListMyDoneTasks failed: %v", err)
	}
	if resp.Total != 1 {
		t.Errorf("expected 1 done task, got %d", resp.Total)
	}
	tasks, ok := resp.Data.([]TaskResponse)
	if !ok || len(tasks) != 1 {
		t.Fatalf("expected 1 task response, got %+v", resp.Data)
	}
	if tasks[0].ID != doneTask.ID.String() {
		t.Errorf("expected done task %s, got %s", doneTask.ID, tasks[0].ID)
	}
}

func TestService_ListMyPendingTasks_EnrichesFlowNameAndSubmitter(t *testing.T) {
	svc, repo, db, cleanup := newTestServiceWithDB()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)

	orgID := uuid.New()
	seedOrganizationNamed(db, orgID, nil, "Finance Department")

	submitterEmployeeID := uuid.New()
	submitterUserID := uuid.New()
	seedEmployee(db, submitterEmployeeID, "Jane Doe", "EMP-001")
	seedEmployment(db, submitterEmployeeID, orgID)
	seedEmployeeAccount(db, submitterEmployeeID, submitterUserID)

	actorID := uuid.New()
	inst := createTestInstance(repo, flow, uuid.New())
	inst.CreatedBy = &submitterUserID
	if err := db.Save(inst).Error; err != nil {
		t.Fatalf("failed to set instance CreatedBy: %v", err)
	}
	createTestTask(repo, inst.ID, 1, actorID)

	resp, err := svc.ListMyPendingTasks(ctx(), actorID.String(), 1, 10, "", nil)
	if err != nil {
		t.Fatalf("ListMyPendingTasks failed: %v", err)
	}
	tasks, ok := resp.Data.([]TaskResponse)
	if !ok || len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %+v", resp.Data)
	}
	task := tasks[0]
	if task.FlowName != flow.Name {
		t.Errorf("expected flow_name %q, got %q", flow.Name, task.FlowName)
	}
	if task.SubmitterName != "Jane Doe" {
		t.Errorf("expected submitter_name 'Jane Doe', got %q", task.SubmitterName)
	}
	if task.SubmitterEmployeeCode != "EMP-001" {
		t.Errorf("expected submitter_employee_code 'EMP-001', got %q", task.SubmitterEmployeeCode)
	}
	if task.SubmitterOrganizationName != "Finance Department" {
		t.Errorf("expected submitter_organization_name 'Finance Department', got %q", task.SubmitterOrganizationName)
	}
}

// =========================================================================
// WATCHER Participation Type Tests
//
// Regression coverage for the bug where a step2=WATCHER task was created
// with status DONE at the moment step1 (APPROVER) was approved, making it
// invisible to ListMyPendingTasks/FindPendingTasksByAssignee (which filter
// strictly on status=PENDING) — so the watcher never saw anything land in
// their task list even though the instance had genuinely progressed past
// their step. Fixed in advanceThroughWatcherSteps: WATCHER-step tasks are
// now created PENDING like any other task, while still never gating
// progression (a WATCHER step is never "landed" on / never becomes the
// instance's current_step).
// =========================================================================

func TestService_SubmitAction_Approve_WatcherStepBecomesVisibleAfterStep1(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")

	approverID := uuid.New()
	watcherID := uuid.New()

	// Step 1: APPROVER (gates progression)
	approverIDStr := approverID.String()
	step1, err := svc.CreateStep(ctx(), flow.ID.String(), CreateStepRequest{
		StepName:       "Manager Approval",
		ApproverType:   "USER",
		ApproverUserID: &approverIDStr,
	})
	if err != nil {
		t.Fatalf("CreateStep (step1) failed: %v", err)
	}
	if step1.ParticipationType != string(ParticipationTypeApprover) {
		t.Fatalf("expected step1 participation_type APPROVER by default, got %q", step1.ParticipationType)
	}

	// Step 2: WATCHER (informational only, must not gate/land, but its task
	// must be visible to the watcher once created)
	watcherIDStr := watcherID.String()
	_, err = svc.CreateStep(ctx(), flow.ID.String(), CreateStepRequest{
		StepName:          "Notify HR",
		ApproverType:      "USER",
		ApproverUserID:    &watcherIDStr,
		ParticipationType: string(ParticipationTypeWatcher),
	})
	if err != nil {
		t.Fatalf("CreateStep (step2/watcher) failed: %v", err)
	}

	inst, err := svc.CreateInstance(ctx(), CreateInstanceRequest{
		Module:     "leave",
		DocumentID: uuidStr(),
		FlowID:     flow.ID.String(),
	})
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	if inst.CurrentStep != 1 {
		t.Fatalf("expected instance to land on step 1 (APPROVER), got current_step=%d", inst.CurrentStep)
	}

	// Approve step 1 — this is the exact scenario reported: after step 1 is
	// approved, the watcher's task must reach them.
	resp, err := svc.SubmitAction(ctx(), inst.ID, approverID.String(), SubmitActionRequest{
		Action: "APPROVE",
	})
	if err != nil {
		t.Fatalf("SubmitAction failed: %v", err)
	}
	if resp.Status != "APPROVED" {
		t.Errorf("expected instance status 'APPROVED' (no more gating steps left), got '%s'", resp.Status)
	}

	// The watcher's task must now be visible in their own pending task list.
	watcherTasks, err := svc.ListMyPendingTasks(ctx(), watcherID.String(), 1, 10, "", nil)
	if err != nil {
		t.Fatalf("ListMyPendingTasks (watcher) failed: %v", err)
	}
	if watcherTasks.Total != 1 {
		t.Fatalf("expected 1 visible task for the watcher after step1 approval, got %d", watcherTasks.Total)
	}

	// The row must be enriched with participation_type=WATCHER and the
	// instance's real status, so the FE can avoid showing an ambiguous
	// "PENDING" tag (task PENDING doesn't mean there's an approve/reject
	// action available to a watcher — the FE displays instance_status
	// instead for WATCHER rows).
	watcherTaskList, ok := watcherTasks.Data.([]TaskResponse)
	if !ok || len(watcherTaskList) != 1 {
		t.Fatalf("expected 1 TaskResponse for the watcher, got %+v", watcherTasks.Data)
	}
	watcherRow := watcherTaskList[0]
	if watcherRow.ParticipationType != string(ParticipationTypeWatcher) {
		t.Errorf("expected participation_type WATCHER, got %q", watcherRow.ParticipationType)
	}
	if watcherRow.InstanceStatus != string(InstanceStatusApproved) {
		t.Errorf("expected instance_status APPROVED, got %q", watcherRow.InstanceStatus)
	}
}

func TestService_SubmitAction_Approve_Step1MultipleApprovers_AllMode(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	step := &ApprovalFlowStep{
		FlowID:       flow.ID,
		StepOrder:    1,
		StepName:     "Multi-Approver Step",
		ApproverType: ApproverTypeUser,
		ApprovalMode: ApprovalModeAll,
		AllowReject:  true,
	}
	if err := repo.CreateStep(ctx(), step); err != nil {
		t.Fatalf("failed to create step: %v", err)
	}

	inst := createTestInstance(repo, flow, uuid.New())

	approver1 := uuid.New()
	approver2 := uuid.New()
	createTestTask(repo, inst.ID, 1, approver1)
	createTestTask(repo, inst.ID, 1, approver2)

	// First approver approves — ALL mode, one task still pending, must not proceed.
	resp, err := svc.SubmitAction(ctx(), inst.ID.String(), approver1.String(), SubmitActionRequest{
		Action: "APPROVE",
	})
	if err != nil {
		t.Fatalf("SubmitAction (approver1) failed: %v", err)
	}
	if resp.Status != "PENDING" {
		t.Errorf("expected instance to remain PENDING with one approver still outstanding, got '%s'", resp.Status)
	}
	if resp.CurrentStep != 1 {
		t.Errorf("expected current_step to remain 1, got %d", resp.CurrentStep)
	}

	// Second approver approves — all tasks for step 1 now done, should proceed.
	resp, err = svc.SubmitAction(ctx(), inst.ID.String(), approver2.String(), SubmitActionRequest{
		Action: "APPROVE",
	})
	if err != nil {
		t.Fatalf("SubmitAction (approver2) failed: %v", err)
	}
	if resp.Status != "APPROVED" {
		t.Errorf("expected instance status 'APPROVED' once all step1 approvers are done, got '%s'", resp.Status)
	}
}

func TestService_SubmitAction_Approve_Step1ThenStep2_BothApprover(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "leave")

	approver1ID := uuid.New()
	approver2ID := uuid.New()
	approver1IDStr := approver1ID.String()
	approver2IDStr := approver2ID.String()

	if _, err := svc.CreateStep(ctx(), flow.ID.String(), CreateStepRequest{
		StepName:       "Step 1 Approval",
		ApproverType:   "USER",
		ApproverUserID: &approver1IDStr,
	}); err != nil {
		t.Fatalf("CreateStep (step1) failed: %v", err)
	}
	if _, err := svc.CreateStep(ctx(), flow.ID.String(), CreateStepRequest{
		StepName:       "Step 2 Approval",
		ApproverType:   "USER",
		ApproverUserID: &approver2IDStr,
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
	if inst.CurrentStep != 1 {
		t.Fatalf("expected instance to land on step 1, got current_step=%d", inst.CurrentStep)
	}

	// Approver 2 must not be able to act while step 1 is still current.
	_, err = svc.SubmitAction(ctx(), inst.ID, approver2ID.String(), SubmitActionRequest{
		Action: "APPROVE",
	})
	if err == nil {
		t.Fatal("expected error when step2's approver tries to act before step1 is resolved")
	}

	// Approver 1 approves step 1 — should advance to step 2 (still PENDING overall).
	resp, err := svc.SubmitAction(ctx(), inst.ID, approver1ID.String(), SubmitActionRequest{
		Action: "APPROVE",
	})
	if err != nil {
		t.Fatalf("SubmitAction (approver1) failed: %v", err)
	}
	if resp.Status != "PENDING" {
		t.Errorf("expected instance status 'PENDING' after step1 approval (step2 still gating), got '%s'", resp.Status)
	}
	if resp.CurrentStep != 2 {
		t.Errorf("expected current_step to advance to 2, got %d", resp.CurrentStep)
	}

	// Approver 1 must not be able to act again on step 2's task.
	_, err = svc.SubmitAction(ctx(), inst.ID, approver1ID.String(), SubmitActionRequest{
		Action: "APPROVE",
	})
	if err == nil {
		t.Fatal("expected error when step1's approver tries to act on step2's task")
	}

	// Approver 2 approves step 2 — no more steps left, instance fully approved.
	resp, err = svc.SubmitAction(ctx(), inst.ID, approver2ID.String(), SubmitActionRequest{
		Action: "APPROVE",
	})
	if err != nil {
		t.Fatalf("SubmitAction (approver2) failed: %v", err)
	}
	if resp.Status != "APPROVED" {
		t.Errorf("expected instance status 'APPROVED' after step2 approval, got '%s'", resp.Status)
	}
}

// TestService_ListMyPendingTasks_RoleAssignedTask_VisibleToRoleMember guards
// a related defect found while fixing the WATCHER-visibility bug above: a
// task routed to a ROLE approver/watcher (assignee_type=ROLE, assignee_id=
// roleID) previously never appeared in ANY user's pending-tasks list, since
// ListMyPendingTasks only ever queried assignee_type=USER matching the
// caller's own user ID literally. Fixed by resolving the caller's RBAC role
// IDs (via model_has_roles) and matching ROLE-assigned tasks against them too.
func TestService_ListMyPendingTasks_RoleAssignedTask_VisibleToRoleMember(t *testing.T) {
	svc, repo, db, cleanup := newTestServiceWithDB()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	roleID := uuid.New()
	step := &ApprovalFlowStep{
		FlowID:       flow.ID,
		StepOrder:    1,
		StepName:     "HR Role Watcher",
		ApproverType: ApproverTypeRole,
		RoleID:       &roleID,
		ApprovalMode: ApprovalModeAnyOne,
		AllowReject:  true,
	}
	if err := repo.CreateStep(ctx(), step); err != nil {
		t.Fatalf("failed to create step: %v", err)
	}

	inst, err := svc.CreateInstance(ctx(), CreateInstanceRequest{
		Module:     "leave",
		DocumentID: uuidStr(),
		FlowID:     flow.ID.String(),
	})
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	if len(inst.Tasks) != 1 || inst.Tasks[0].AssigneeID != roleID.String() {
		t.Fatalf("expected a single ROLE-assigned task, got %+v", inst.Tasks)
	}

	memberUserID := uuid.New()
	seedUserRole(db, memberUserID, roleID)

	resp, err := svc.ListMyPendingTasks(ctx(), memberUserID.String(), 1, 10, "", nil)
	if err != nil {
		t.Fatalf("ListMyPendingTasks failed: %v", err)
	}
	if resp.Total != 1 {
		t.Fatalf("expected 1 visible task for a member of the assigned role, got %d", resp.Total)
	}

	// A user who does NOT hold the role must not see the task.
	outsiderUserID := uuid.New()
	resp, err = svc.ListMyPendingTasks(ctx(), outsiderUserID.String(), 1, 10, "", nil)
	if err != nil {
		t.Fatalf("ListMyPendingTasks (outsider) failed: %v", err)
	}
	if resp.Total != 0 {
		t.Errorf("expected 0 visible tasks for a non-member of the role, got %d", resp.Total)
	}
}

func TestService_ListMyPendingTasks_DefaultPagination(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	resp, err := svc.ListMyPendingTasks(ctx(), uuidStr(), 0, 0, "", nil)
	if err != nil {
		t.Fatalf("ListMyPendingTasks failed: %v", err)
	}

	if resp.Page != 1 {
		t.Errorf("expected page 1, got %d", resp.Page)
	}
	if resp.PerPage != 20 {
		t.Errorf("expected per_page 20, got %d", resp.PerPage)
	}
}
