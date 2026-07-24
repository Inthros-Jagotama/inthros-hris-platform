package approval

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// =========================================================================
// Approval Flow Repository Tests
// =========================================================================

func TestRepo_CreateFlow_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	f := &ApprovalFlow{
		Module:   "leave",
		Name:     "Leave Approval Flow",
		Version:  1,
		IsActive: true,
	}

	if err := repo.CreateFlow(ctx, f); err != nil {
		t.Fatalf("CreateFlow failed: %v", err)
	}

	if f.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindFlowByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestFlow(repo, "leave")

	found, err := repo.FindFlowByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindFlowByID failed: %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("expected ID '%s', got '%s'", created.ID, found.ID)
	}
	if found.Module != "leave" {
		t.Errorf("expected module 'leave', got '%s'", found.Module)
	}
}

func TestRepo_FindFlowByID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	_, err := repo.FindFlowByID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent flow")
	}
}

func TestRepo_FindFlowByIDWithSteps_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)
	createTestStep(repo, flow.ID, 2)

	found, err := repo.FindFlowByIDWithSteps(ctx, flow.ID)
	if err != nil {
		t.Fatalf("FindFlowByIDWithSteps failed: %v", err)
	}

	if len(found.Steps) != 2 {
		t.Errorf("expected 2 steps, got %d", len(found.Steps))
	}
	if found.Steps[0].StepOrder != 1 {
		t.Errorf("expected step 1 order 1, got %d", found.Steps[0].StepOrder)
	}
}

func TestRepo_ListFlows_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		createTestFlow(repo, "leave")
	}

	flows, total, err := repo.ListFlows(ctx, 1, 3)
	if err != nil {
		t.Fatalf("ListFlows failed: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(flows) != 3 {
		t.Errorf("expected 3 flows (page 1 of 2), got %d", len(flows))
	}
}

func TestRepo_UpdateFlow_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestFlow(repo, "leave")

	created.Name = "Updated Flow Name"
	if err := repo.UpdateFlow(ctx, created); err != nil {
		t.Fatalf("UpdateFlow failed: %v", err)
	}

	found, _ := repo.FindFlowByID(ctx, created.ID)
	if found.Name != "Updated Flow Name" {
		t.Errorf("expected name 'Updated Flow Name', got '%s'", found.Name)
	}
}

func TestRepo_SoftDeleteFlow_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestFlow(repo, "leave")

	if err := repo.SoftDeleteFlow(ctx, created.ID); err != nil {
		t.Fatalf("SoftDeleteFlow failed: %v", err)
	}

	// Should not be found in standard query (deleted_at IS NULL)
	_, err := repo.FindFlowByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after soft-deleting flow")
	}
}

func TestRepo_ListFlowsByModule(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	createTestFlow(repo, "leave")
	createTestFlow(repo, "leave")
	createTestFlow(repo, "overtime")

	flows, err := repo.ListFlowsByModule(ctx, "leave")
	if err != nil {
		t.Fatalf("ListFlowsByModule failed: %v", err)
	}

	if len(flows) != 2 {
		t.Errorf("expected 2 flows for 'leave', got %d", len(flows))
	}
}

// =========================================================================
// Approval Flow Steps Repository Tests
// =========================================================================

func TestRepo_CreateStep_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	flow := createTestFlow(repo, "leave")

	s := &ApprovalFlowStep{
		FlowID:       flow.ID,
		StepOrder:    1,
		StepName:     "Manager Approval",
		ApproverType: ApproverTypeUser,
		ApproverUserID: func() *uuid.UUID {
			uid := uuid.New()
			return &uid
		}(),
		ApprovalMode: ApprovalModeAnyOne,
		AllowReject:  true,
	}

	if err := repo.CreateStep(ctx, s); err != nil {
		t.Fatalf("CreateStep failed: %v", err)
	}

	if s.ID == uuid.Nil {
		t.Error("expected step ID to be auto-generated")
	}
}

func TestRepo_ListStepsByFlowID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)
	createTestStep(repo, flow.ID, 2)
	createTestStep(repo, flow.ID, 3)

	steps, err := repo.ListStepsByFlowID(ctx, flow.ID)
	if err != nil {
		t.Fatalf("ListStepsByFlowID failed: %v", err)
	}

	if len(steps) != 3 {
		t.Errorf("expected 3 steps, got %d", len(steps))
	}
	if steps[0].StepOrder != 1 {
		t.Errorf("expected first step order 1, got %d", steps[0].StepOrder)
	}
}

func TestRepo_GetMaxStepOrder(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)
	createTestStep(repo, flow.ID, 2)

	maxOrder, err := repo.GetMaxStepOrder(ctx, flow.ID)
	if err != nil {
		t.Fatalf("GetMaxStepOrder failed: %v", err)
	}

	if maxOrder != 2 {
		t.Errorf("expected max order 2, got %d", maxOrder)
	}
}

func TestRepo_SoftDeleteStep_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	flow := createTestFlow(repo, "leave")
	step := createTestStep(repo, flow.ID, 1)

	if err := repo.SoftDeleteStep(ctx, step.ID); err != nil {
		t.Fatalf("SoftDeleteStep failed: %v", err)
	}

	_, err := repo.FindStepByID(ctx, step.ID)
	if err == nil {
		t.Fatal("expected error after soft-deleting step")
	}
}

// =========================================================================
// Approval Instance Repository Tests
// =========================================================================

func TestRepo_CreateInstance_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	flow := createTestFlow(repo, "leave")

	inst := &ApprovalInstance{
		Module:      "leave",
		DocumentID:  uuid.New(),
		FlowID:      flow.ID,
		Status:      InstanceStatusPending,
		CurrentStep: 1,
	}

	if err := repo.CreateInstance(ctx, inst); err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}

	if inst.ID == uuid.Nil {
		t.Error("expected instance ID to be auto-generated")
	}
}

func TestRepo_FindInstanceByIDWithRelations_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)

	docID := uuid.New()
	inst := createTestInstance(repo, flow, docID)

	// Create action and task
	actorID := uuid.New()
	createTestAction(repo, inst.ID, 1, actorID, ActionApprove)
	createTestTask(repo, inst.ID, 1, actorID)

	found, err := repo.FindInstanceByIDWithRelations(ctx, inst.ID)
	if err != nil {
		t.Fatalf("FindInstanceByIDWithRelations failed: %v", err)
	}

	if found.Module != "leave" {
		t.Errorf("expected module 'leave', got '%s'", found.Module)
	}
	if found.Flow == nil || found.Flow.Name == "" {
		t.Error("expected flow relation to be loaded")
	}
	if len(found.Actions) != 1 {
		t.Errorf("expected 1 action, got %d", len(found.Actions))
	}
	if len(found.Tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(found.Tasks))
	}
}

func TestRepo_FindInstanceByDocument_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	flow := createTestFlow(repo, "leave")
	docID := uuid.New()
	createTestInstance(repo, flow, docID)

	found, err := repo.FindInstanceByDocument(ctx, "leave", docID)
	if err != nil {
		t.Fatalf("FindInstanceByDocument failed: %v", err)
	}

	if found.DocumentID != docID {
		t.Errorf("expected document_id '%s', got '%s'", docID, found.DocumentID)
	}
}

func TestRepo_ListInstances_WithFilters(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	flow1 := createTestFlow(repo, "leave")
	flow2 := createTestFlow(repo, "overtime")

	createTestInstance(repo, flow1, uuid.New())
	createTestInstance(repo, flow1, uuid.New())
	createTestInstance(repo, flow2, uuid.New())

	// List all
	_, totalAll, err := repo.ListInstances(ctx, 1, 10, "", "")
	if err != nil {
		t.Fatalf("ListInstances failed: %v", err)
	}
	if totalAll != 3 {
		t.Errorf("expected total 3, got %d", totalAll)
	}

	// Filter by module
	_, totalLeave, err := repo.ListInstances(ctx, 1, 10, "leave", "")
	if err != nil {
		t.Fatalf("ListInstances by module failed: %v", err)
	}
	if totalLeave != 2 {
		t.Errorf("expected 2 leave instances, got %d", totalLeave)
	}

	// Filter by status
	_, totalPending, err := repo.ListInstances(ctx, 1, 10, "", string(InstanceStatusPending))
	if err != nil {
		t.Fatalf("ListInstances by status failed: %v", err)
	}
	if totalPending != 3 {
		t.Errorf("expected 3 pending instances, got %d", totalPending)
	}
}

func TestRepo_CancelInstance_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	flow := createTestFlow(repo, "leave")
	inst := createTestInstance(repo, flow, uuid.New())

	// Create a pending task
	assigneeID := uuid.New()
	createTestTask(repo, inst.ID, 1, assigneeID)

	if err := repo.CancelInstance(ctx, inst.ID); err != nil {
		t.Fatalf("CancelInstance failed: %v", err)
	}

	cancelled, _ := repo.FindInstanceByID(ctx, inst.ID)
	if cancelled.Status != InstanceStatusCancelled {
		t.Errorf("expected status CANCELLED, got '%s'", cancelled.Status)
	}

	// Task should be cancelled too
	tasks, _ := repo.FindTasksByInstanceID(ctx, inst.ID)
	if len(tasks) > 0 && tasks[0].Status != TaskStatusCancelled {
		t.Errorf("expected task status CANCELLED, got '%s'", tasks[0].Status)
	}
}

func TestRepo_CancelInstance_NonPending_Error(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	flow := createTestFlow(repo, "leave")
	inst := createTestInstance(repo, flow, uuid.New())
	inst.Status = InstanceStatusApproved
	repo.UpdateInstance(ctx, inst)

	err := repo.CancelInstance(ctx, inst.ID)
	if err == nil {
		t.Fatal("expected error when cancelling non-PENDING instance")
	}
}

// =========================================================================
// Approval Action Repository Tests
// =========================================================================

func TestRepo_CreateAction_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	flow := createTestFlow(repo, "leave")
	inst := createTestInstance(repo, flow, uuid.New())
	actorID := uuid.New()

	a := &ApprovalAction{
		InstanceID:  inst.ID,
		StepOrder:   1,
		ActorUserID: actorID,
		Action:      ActionApprove,
	}

	if err := repo.CreateAction(ctx, a); err != nil {
		t.Fatalf("CreateAction failed: %v", err)
	}

	if a.ID == uuid.Nil {
		t.Error("expected action ID to be auto-generated")
	}
}

// =========================================================================
// Approval Task Repository Tests
// =========================================================================

func TestRepo_CreateTasks_Bulk(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	flow := createTestFlow(repo, "leave")
	inst := createTestInstance(repo, flow, uuid.New())

	tasks := []ApprovalTask{
		{InstanceID: inst.ID, StepOrder: 1, AssigneeType: "USER", AssigneeID: uuid.New(), Status: TaskStatusPending},
		{InstanceID: inst.ID, StepOrder: 1, AssigneeType: "USER", AssigneeID: uuid.New(), Status: TaskStatusPending},
	}

	if err := repo.CreateTasks(ctx, tasks); err != nil {
		t.Fatalf("CreateTasks failed: %v", err)
	}

	loaded, err := repo.FindTasksByInstanceID(ctx, inst.ID)
	if err != nil {
		t.Fatalf("FindTasksByInstanceID failed: %v", err)
	}

	if len(loaded) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(loaded))
	}
}

func TestRepo_FindPendingTasksByAssignee_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	flow := createTestFlow(repo, "leave")
	inst1 := createTestInstance(repo, flow, uuid.New())
	inst2 := createTestInstance(repo, flow, uuid.New())

	assigneeID := uuid.New()
	createTestTask(repo, inst1.ID, 1, assigneeID)
	createTestTask(repo, inst1.ID, 1, uuid.New()) // different assignee
	createTestTask(repo, inst2.ID, 1, assigneeID)

	tasks, total, err := repo.FindPendingTasksByAssignee(ctx, "USER", assigneeID, 1, 10)
	if err != nil {
		t.Fatalf("FindPendingTasksByAssignee failed: %v", err)
	}

	if total != 2 {
		t.Errorf("expected total 2 tasks for assignee, got %d", total)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

// =========================================================================
// Transactional Operations Tests
// =========================================================================

func TestRepo_CreateInstanceWithTasks_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	flow := createTestFlow(repo, "leave")

	inst := &ApprovalInstance{
		Module:      "leave",
		DocumentID:  uuid.New(),
		FlowID:      flow.ID,
		Status:      InstanceStatusPending,
		CurrentStep: 1,
	}

	tasks := []ApprovalTask{
		{StepOrder: 1, AssigneeType: "USER", AssigneeID: uuid.New(), Status: TaskStatusPending},
		{StepOrder: 1, AssigneeType: "USER", AssigneeID: uuid.New(), Status: TaskStatusPending},
	}

	if err := repo.CreateInstanceWithTasks(ctx, inst, tasks); err != nil {
		t.Fatalf("CreateInstanceWithTasks failed: %v", err)
	}

	if inst.ID == uuid.Nil {
		t.Error("expected instance ID to be auto-generated")
	}

	loadedTasks, _ := repo.FindTasksByInstanceID(ctx, inst.ID)
	if len(loadedTasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(loadedTasks))
	}
	for _, task := range loadedTasks {
		if task.InstanceID != inst.ID {
			t.Errorf("expected task instance_id '%s', got '%s'", inst.ID, task.InstanceID)
		}
	}
}

func TestRepo_ApproveStep_AdvanceToNextStep(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)
	createTestStep(repo, flow.ID, 2)

	inst := createTestInstance(repo, flow, uuid.New())

	// Create task for step 1 (already done)
	actorID := uuid.New()
	task := createTestTask(repo, inst.ID, 1, actorID)
	repo.UpdateTaskStatus(ctx, task.ID, TaskStatusDone)

	// Create tasks for next step
	nextTasks := []ApprovalTask{
		{StepOrder: 2, AssigneeType: "USER", AssigneeID: uuid.New(), Status: TaskStatusPending},
	}

	if err := repo.ApproveStep(ctx, inst.ID, 1, 2, nextTasks); err != nil {
		t.Fatalf("ApproveStep failed: %v", err)
	}

	updated, _ := repo.FindInstanceByID(ctx, inst.ID)
	if updated.CurrentStep != 2 {
		t.Errorf("expected current_step 2, got %d", updated.CurrentStep)
	}

	// Check that step 2 tasks were created
	allTasks, _ := repo.FindTasksByInstanceID(ctx, inst.ID)
	step2Tasks := 0
	for _, t := range allTasks {
		if t.StepOrder == 2 {
			step2Tasks++
		}
	}
	if step2Tasks != 1 {
		t.Errorf("expected 1 task for step 2, got %d", step2Tasks)
	}
}

func TestRepo_ApproveStep_CompleteInstance(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)

	inst := createTestInstance(repo, flow, uuid.New())

	// Create task for step 1 (already done)
	actorID := uuid.New()
	task := createTestTask(repo, inst.ID, 1, actorID)
	repo.UpdateTaskStatus(ctx, task.ID, TaskStatusDone)

	// Approve step 1 (no more steps → instance approved)
	if err := repo.ApproveStep(ctx, inst.ID, 1, 1, nil); err != nil {
		t.Fatalf("ApproveStep failed: %v", err)
	}

	updated, _ := repo.FindInstanceByID(ctx, inst.ID)
	if updated.Status != InstanceStatusApproved {
		t.Errorf("expected status APPROVED, got '%s'", updated.Status)
	}
}

func TestRepo_RejectInstance_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	flow := createTestFlow(repo, "leave")
	inst := createTestInstance(repo, flow, uuid.New())

	// Create pending task
	createTestTask(repo, inst.ID, 1, uuid.New())

	if err := repo.RejectInstance(ctx, inst.ID); err != nil {
		t.Fatalf("RejectInstance failed: %v", err)
	}

	rejected, _ := repo.FindInstanceByID(ctx, inst.ID)
	if rejected.Status != InstanceStatusRejected {
		t.Errorf("expected status REJECTED, got '%s'", rejected.Status)
	}

	// Tasks should be cancelled
	tasks, _ := repo.FindTasksByInstanceID(ctx, inst.ID)
	if len(tasks) > 0 && tasks[0].Status != TaskStatusCancelled {
		t.Errorf("expected task status CANCELLED, got '%s'", tasks[0].Status)
	}
}
