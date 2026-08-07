package performance

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestService_GetEvaluationWithDetails_EnrichesNames(t *testing.T) {
	svc, dbResolver, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()

	db, err := dbResolver(context.Background())
	if err != nil {
		t.Fatalf("failed to get db: %v", err)
	}

	employeeID := uuid.New()
	orgID := uuid.New()
	if err := db.Exec("INSERT INTO employees (id, name) VALUES (?, ?)", employeeID.String(), "Jane Doe").Error; err != nil {
		t.Fatalf("failed to seed employee: %v", err)
	}
	if err := db.Exec("INSERT INTO organizations (id, nomenclature) VALUES (?, ?)", orgID.String(), "Finance Manager").Error; err != nil {
		t.Fatalf("failed to seed organization: %v", err)
	}

	period := &PerformancePeriod{PeriodCode: "2026-Q1", PeriodType: "QUARTERLY", Year: 2026}
	if err := db.Create(period).Error; err != nil {
		t.Fatalf("failed to seed period: %v", err)
	}
	tmpl := &PerformanceTemplate{OrganizationID: orgID, Name: "FY2026 BSC", Status: "PUBLISHED"}
	if err := db.Create(tmpl).Error; err != nil {
		t.Fatalf("failed to seed template: %v", err)
	}
	eval := &PerformanceEvaluation{
		EmployeeID:     employeeID,
		OrganizationID: orgID,
		PeriodID:       period.ID,
		TemplateID:     tmpl.ID,
		Status:         "DRAFT",
	}
	if err := db.Create(eval).Error; err != nil {
		t.Fatalf("failed to seed evaluation: %v", err)
	}

	resp, err := svc.GetEvaluationWithDetails(context.Background(), eval.ID.String())
	if err != nil {
		t.Fatalf("GetEvaluationWithDetails failed: %v", err)
	}
	if resp.EmployeeName != "Jane Doe" {
		t.Errorf("expected employee_name 'Jane Doe', got %q", resp.EmployeeName)
	}
	if resp.OrganizationName != "Finance Manager" {
		t.Errorf("expected organization_name 'Finance Manager', got %q", resp.OrganizationName)
	}
	if resp.PeriodCode != "2026-Q1" {
		t.Errorf("expected period_code '2026-Q1', got %q", resp.PeriodCode)
	}
}
