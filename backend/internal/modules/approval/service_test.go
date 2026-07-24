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

	resp, err := svc.ListMyPendingTasks(ctx(), actorID.String(), 1, 10)
	if err != nil {
		t.Fatalf("ListMyPendingTasks failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2 tasks, got %d", resp.Total)
	}
}

func TestService_ListMyPendingTasks_DefaultPagination(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	resp, err := svc.ListMyPendingTasks(ctx(), uuidStr(), 0, 0)
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
