package approval

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

// setupTestDB creates an in-memory SQLite database and auto-migrates all models.
func setupTestDB() (*gorm.DB, func(ctx context.Context) (*gorm.DB, error), func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	if err := db.AutoMigrate(
		&ApprovalFlow{},
		&ApprovalFlowStep{},
		&ApprovalInstance{},
		&ApprovalAction{},
		&ApprovalTask{},
	); err != nil {
		panic(fmt.Sprintf("failed to migrate test db: %v", err))
	}

	dbResolver := func(ctx context.Context) (*gorm.DB, error) {
		return db, nil
	}

	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	return db, dbResolver, cleanup
}

// newTestService creates a Service with in-memory SQLite repository.
func newTestService() (*Service, *Repository, func()) {
	_, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	return svc, repo, func() {
		cleanup()
		_ = logger.Sync()
	}
}

// createTestFlow inserts a test approval flow.
func createTestFlow(repo *Repository, module string) *ApprovalFlow {
	ctx := context.Background()
	f := &ApprovalFlow{
		Module:   module,
		Name:     fmt.Sprintf("Flow %s", module),
		Version:  1,
		IsActive: true,
	}
	if err := repo.CreateFlow(ctx, f); err != nil {
		panic(fmt.Sprintf("failed to create test flow: %v", err))
	}
	return f
}

// createTestStep inserts a test approval flow step.
func createTestStep(repo *Repository, flowID uuid.UUID, stepOrder int) *ApprovalFlowStep {
	ctx := context.Background()
	s := &ApprovalFlowStep{
		FlowID:       flowID,
		StepOrder:    stepOrder,
		StepName:     fmt.Sprintf("Step %d", stepOrder),
		ApproverType: ApproverTypeUser,
		ApproverUserID: func() *uuid.UUID {
			uid := uuid.New()
			return &uid
		}(),
		ApprovalMode: ApprovalModeAnyOne,
		AllowReject:  true,
	}
	if err := repo.CreateStep(ctx, s); err != nil {
		panic(fmt.Sprintf("failed to create test step: %v", err))
	}
	return s
}

// createTestInstance inserts a test approval instance for the given flow and document.
func createTestInstance(repo *Repository, flow *ApprovalFlow, documentID uuid.UUID) *ApprovalInstance {
	ctx := context.Background()
	inst := &ApprovalInstance{
		Module:      flow.Module,
		DocumentID:  documentID,
		FlowID:      flow.ID,
		Status:      InstanceStatusPending,
		CurrentStep: 1,
	}
	if err := repo.CreateInstance(ctx, inst); err != nil {
		panic(fmt.Sprintf("failed to create test instance: %v", err))
	}
	return inst
}

// createTestAction inserts a test approval action.
func createTestAction(repo *Repository, instanceID uuid.UUID, stepOrder int, actorID uuid.UUID, action ApprovalActionType) *ApprovalAction {
	ctx := context.Background()
	a := &ApprovalAction{
		InstanceID:  instanceID,
		StepOrder:   stepOrder,
		ActorUserID: actorID,
		Action:      action,
	}
	if err := repo.CreateAction(ctx, a); err != nil {
		panic(fmt.Sprintf("failed to create test action: %v", err))
	}
	return a
}

// createTestTask inserts a test approval task.
func createTestTask(repo *Repository, instanceID uuid.UUID, stepOrder int, assigneeID uuid.UUID) *ApprovalTask {
	ctx := context.Background()
	t := &ApprovalTask{
		InstanceID:   instanceID,
		StepOrder:    stepOrder,
		AssigneeType: "USER",
		AssigneeID:   assigneeID,
		Status:       TaskStatusPending,
	}
	if err := repo.CreateTasks(ctx, []ApprovalTask{*t}); err != nil {
		panic(fmt.Sprintf("failed to create test task: %v", err))
	}
	return t
}

// uuidStr returns a UUID string for test use.
func uuidStr() string {
	return uuid.New().String()
}

// strPtr returns a pointer to the given string.
func strPtr(s string) *string {
	return &s
}

// intPtr returns a pointer to the given int.
func intPtr(i int) *int {
	return &i
}
