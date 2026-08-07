package performance

import (
	"context"
	"testing"

	sqlite "github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// setupMyKPIContextTestDB creates an in-memory SQLite database with the
// performance models plus the minimal raw tables (organizations/employments/
// employee_accounts) that GetCurrentEmployeeContextByUserID/
// GetOrganizationNamesByIDs query directly via db.Table(...) — performance
// must not import the organization/employee/useraccount packages.
func setupMyKPIContextTestDB(t *testing.T) (*Service, func(ctx context.Context) (*gorm.DB, error), func()) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(
		&PerformancePeriod{},
		&PerformancePerspective{},
		&PerformanceTemplate{},
		&PerformanceIndicator{},
		&PerformanceEvaluation{},
		&PerformanceEvaluationDetail{},
		&PerformanceEvaluationProgramItem{},
		&PerformanceComponent{},
		&PerformanceOrganizationComponent{},
		&PerformanceRating{},
	); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	rawTables := []string{
		`CREATE TABLE IF NOT EXISTS organizations (
			id CHAR(36) PRIMARY KEY,
			nomenclature VARCHAR(255) NOT NULL DEFAULT ''
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
			name VARCHAR(255) NOT NULL DEFAULT ''
		)`,
	}
	for _, stmt := range rawTables {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("failed to create raw test table: %v", err)
		}
	}

	dbResolver := func(ctx context.Context) (*gorm.DB, error) {
		return db, nil
	}
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)

	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		_ = logger.Sync()
	}
	return svc, dbResolver, cleanup
}

func ctxWithUser(userID uuid.UUID) context.Context {
	return context.WithValue(context.Background(), "user_id", userID.String())
}

func seedMyKPIContextEmployee(t *testing.T, dbResolver func(ctx context.Context) (*gorm.DB, error), userID, employeeID, orgID uuid.UUID, orgName string) {
	t.Helper()
	db, err := dbResolver(context.Background())
	if err != nil {
		t.Fatalf("failed to get db: %v", err)
	}
	if err := db.Exec("INSERT INTO organizations (id, nomenclature) VALUES (?, ?)", orgID.String(), orgName).Error; err != nil {
		t.Fatalf("failed to seed organization: %v", err)
	}
	if err := db.Exec(
		"INSERT INTO employments (id, employee_id, organization_id, effective_date, effective_end_date) VALUES (?, ?, ?, ?, NULL)",
		uuid.New().String(), employeeID.String(), orgID.String(), "2026-01-01",
	).Error; err != nil {
		t.Fatalf("failed to seed employment: %v", err)
	}
	if err := db.Exec(
		"INSERT INTO employee_accounts (id, employee_id, user_id) VALUES (?, ?, ?)",
		uuid.New().String(), employeeID.String(), userID.String(),
	).Error; err != nil {
		t.Fatalf("failed to seed employee account: %v", err)
	}
}

func TestService_GetMyKPIContext_NoPosition(t *testing.T) {
	svc, _, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()

	resp, err := svc.GetMyKPIContext(ctxWithUser(uuid.New()))
	if err != nil {
		t.Fatalf("GetMyKPIContext failed: %v", err)
	}
	if resp.HasPosition {
		t.Error("expected HasPosition false when no employment record exists")
	}
	if len(resp.Templates) != 0 {
		t.Errorf("expected no templates, got %d", len(resp.Templates))
	}
}

func TestService_GetMyKPIContext_Unauthenticated(t *testing.T) {
	svc, _, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()

	_, err := svc.GetMyKPIContext(context.Background())
	if err == nil {
		t.Fatal("expected error when no user_id in context")
	}
}

func TestService_GetMyKPIContext_WithPublishedTemplate(t *testing.T) {
	svc, dbResolver, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()

	userID := uuid.New()
	employeeID := uuid.New()
	orgID := uuid.New()
	seedMyKPIContextEmployee(t, dbResolver, userID, employeeID, orgID, "Finance Manager")

	db, _ := dbResolver(context.Background())
	published := &PerformanceTemplate{OrganizationID: orgID, Name: "FY2026 BSC", Status: "PUBLISHED"}
	if err := db.Create(published).Error; err != nil {
		t.Fatalf("failed to seed published template: %v", err)
	}
	draft := &PerformanceTemplate{OrganizationID: orgID, Name: "Draft template", Status: "DRAFT"}
	if err := db.Create(draft).Error; err != nil {
		t.Fatalf("failed to seed draft template: %v", err)
	}

	resp, err := svc.GetMyKPIContext(ctxWithUser(userID))
	if err != nil {
		t.Fatalf("GetMyKPIContext failed: %v", err)
	}
	if !resp.HasPosition {
		t.Fatal("expected HasPosition true")
	}
	if resp.EmployeeID != employeeID.String() {
		t.Errorf("expected employee_id %s, got %s", employeeID, resp.EmployeeID)
	}
	if resp.OrganizationID != orgID.String() {
		t.Errorf("expected organization_id %s, got %s", orgID, resp.OrganizationID)
	}
	if resp.OrganizationName != "Finance Manager" {
		t.Errorf("expected organization_name 'Finance Manager', got %q", resp.OrganizationName)
	}
	if len(resp.Templates) != 1 {
		t.Fatalf("expected exactly 1 PUBLISHED template (draft excluded), got %d", len(resp.Templates))
	}
	if resp.Templates[0].Name != "FY2026 BSC" {
		t.Errorf("expected published template 'FY2026 BSC', got %q", resp.Templates[0].Name)
	}
}
