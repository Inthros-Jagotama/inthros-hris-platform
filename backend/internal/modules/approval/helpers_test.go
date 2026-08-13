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
		&ApprovalFlowStepOrganization{},
		&ApprovalInstance{},
		&ApprovalAction{},
		&ApprovalTask{},
	); err != nil {
		panic(fmt.Sprintf("failed to migrate test db: %v", err))
	}

	// Raw tables owned by other modules (organization/employee/useraccount) that
	// the org-hierarchy resolution queries directly via db.Table(...) — approval
	// must not import those packages (circular dependency risk). Minimal schema
	// covering only the columns the raw queries touch.
	rawTables := []string{
		`CREATE TABLE IF NOT EXISTS organizations (
			id CHAR(36) PRIMARY KEY,
			parent_id CHAR(36) NULL,
			nomenclature VARCHAR(255) NOT NULL DEFAULT '',
			deleted_at TIMESTAMP NULL
		)`,
		`CREATE TABLE IF NOT EXISTS employments (
			id CHAR(36) PRIMARY KEY,
			employee_id CHAR(36) NOT NULL,
			organization_id CHAR(36) NOT NULL,
			effective_date DATE NOT NULL,
			effective_end_date DATE NULL
		)`,
		`CREATE TABLE IF NOT EXISTS employee_accounts (
			id CHAR(36) PRIMARY KEY,
			employee_id CHAR(36) NOT NULL,
			user_id CHAR(36) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS employees (
			id CHAR(36) PRIMARY KEY,
			name VARCHAR(255) NOT NULL DEFAULT '',
			employee_id VARCHAR(50) NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS model_has_roles (
			role_id CHAR(36) NOT NULL,
			model_type VARCHAR(255) NOT NULL,
			model_id CHAR(36) NOT NULL,
			PRIMARY KEY (role_id, model_type, model_id)
		)`,
	}
	for _, stmt := range rawTables {
		if err := db.Exec(stmt).Error; err != nil {
			panic(fmt.Sprintf("failed to create raw test table: %v", err))
		}
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
	svc, repo, _, cleanup := newTestServiceWithDB()
	return svc, repo, cleanup
}

// newTestServiceWithDB is like newTestService but also exposes the underlying
// *gorm.DB, needed by tests that seed raw organization/employment/employee_account
// fixtures for org-hierarchy approver resolution.
func newTestServiceWithDB() (*Service, *Repository, *gorm.DB, func()) {
	db, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	return svc, repo, db, func() {
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

// createTestTaskROLE like createTestTask but assigns the task to a ROLE
// (assignee_type ROLE) instead of a direct user — regression coverage for
// SubmitAction approving role-assigned tasks.
func createTestTaskROLE(repo *Repository, instanceID uuid.UUID, stepOrder int, roleID uuid.UUID) *ApprovalTask {
	ctx := context.Background()
	t := &ApprovalTask{
		InstanceID:   instanceID,
		StepOrder:    stepOrder,
		AssigneeType: "ROLE",
		AssigneeID:   roleID,
		Status:       TaskStatusPending,
	}
	if err := repo.CreateTasks(ctx, []ApprovalTask{*t}); err != nil {
		panic(fmt.Sprintf("failed to create test role task: %v", err))
	}
	return t
}

// seedOrganization inserts a minimal organizations row.
func seedOrganization(db *gorm.DB, id uuid.UUID, parentID *uuid.UUID) {
	if err := db.Exec("INSERT INTO organizations (id, parent_id) VALUES (?, ?)", id.String(), uuidPtrStr(parentID)).Error; err != nil {
		panic(fmt.Sprintf("failed to seed organization: %v", err))
	}
}

// seedOrganizationNamed is like seedOrganization but also sets nomenclature,
// needed by tests asserting on the resolved organization *name* (not just ID).
func seedOrganizationNamed(db *gorm.DB, id uuid.UUID, parentID *uuid.UUID, name string) {
	if err := db.Exec("INSERT INTO organizations (id, parent_id, nomenclature) VALUES (?, ?, ?)", id.String(), uuidPtrStr(parentID), name).Error; err != nil {
		panic(fmt.Sprintf("failed to seed organization: %v", err))
	}
}

// seedEmployee inserts a minimal employees row (name + employee code).
func seedEmployee(db *gorm.DB, id uuid.UUID, name, employeeCode string) {
	if err := db.Exec("INSERT INTO employees (id, name, employee_id) VALUES (?, ?, ?)", id.String(), name, employeeCode).Error; err != nil {
		panic(fmt.Sprintf("failed to seed employee: %v", err))
	}
}

// seedEmployment links an employee to an Organization as their current (open-ended) assignment.
func seedEmployment(db *gorm.DB, employeeID, orgID uuid.UUID) {
	if err := db.Exec(
		"INSERT INTO employments (id, employee_id, organization_id, effective_date, effective_end_date) VALUES (?, ?, ?, ?, NULL)",
		uuid.New().String(), employeeID.String(), orgID.String(), "2026-01-01",
	).Error; err != nil {
		panic(fmt.Sprintf("failed to seed employment: %v", err))
	}
}

// seedEmployeeAccount links an employee to their platform user login.
func seedEmployeeAccount(db *gorm.DB, employeeID, userID uuid.UUID) {
	if err := db.Exec(
		"INSERT INTO employee_accounts (id, employee_id, user_id) VALUES (?, ?, ?)",
		uuid.New().String(), employeeID.String(), userID.String(),
	).Error; err != nil {
		panic(fmt.Sprintf("failed to seed employee account: %v", err))
	}
}

// seedUserRole links a platform user to an RBAC role via model_has_roles.
func seedUserRole(db *gorm.DB, userID, roleID uuid.UUID) {
	if err := db.Exec(
		"INSERT INTO model_has_roles (role_id, model_type, model_id) VALUES (?, ?, ?)",
		roleID.String(), "user", userID.String(),
	).Error; err != nil {
		panic(fmt.Sprintf("failed to seed user role: %v", err))
	}
}

func uuidPtrStr(id *uuid.UUID) interface{} {
	if id == nil {
		return nil
	}
	return id.String()
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

// boolPtr returns a pointer to the given bool.
func boolPtr(b bool) *bool {
	return &b
}
