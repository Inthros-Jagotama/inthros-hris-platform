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

func TestService_GetEvaluationWithDetails_EnrichesDetailRows(t *testing.T) {
	svc, dbResolver, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()

	db, err := dbResolver(context.Background())
	if err != nil {
		t.Fatalf("failed to get db: %v", err)
	}

	orgID := uuid.New()
	if err := db.Exec("INSERT INTO organizations (id, nomenclature) VALUES (?, ?)", orgID.String(), "Ops").Error; err != nil {
		t.Fatalf("failed to seed organization: %v", err)
	}
	period := &PerformancePeriod{PeriodCode: "2026-Q1", PeriodType: "QUARTERLY", Year: 2026}
	if err := db.Create(period).Error; err != nil {
		t.Fatalf("failed to seed period: %v", err)
	}
	tmpl := &PerformanceTemplate{OrganizationID: orgID, Name: "Tmpl", Status: "PUBLISHED"}
	if err := db.Create(tmpl).Error; err != nil {
		t.Fatalf("failed to seed template: %v", err)
	}
	persp := &PerformancePerspective{Name: "Financial"}
	if err := db.Create(persp).Error; err != nil {
		t.Fatalf("failed to seed perspective: %v", err)
	}
	unit := "%"
	indicator := &PerformanceIndicator{
		PerformanceTemplateID: tmpl.ID,
		PerspectiveID:         persp.ID,
		IndicatorType:         "MAXIMIZATION",
		Title:                 "Revenue Growth",
		FormulaType:           "HIGHER_BETTER",
		UnitOfMeasurement:     &unit,
	}
	if err := db.Create(indicator).Error; err != nil {
		t.Fatalf("failed to seed indicator: %v", err)
	}
	eval := &PerformanceEvaluation{
		EmployeeID:     uuid.New(),
		OrganizationID: orgID,
		PeriodID:       period.ID,
		TemplateID:     tmpl.ID,
		Status:         "DRAFT",
	}
	if err := db.Create(eval).Error; err != nil {
		t.Fatalf("failed to seed evaluation: %v", err)
	}
	detail := &PerformanceEvaluationDetail{
		PerformanceEvaluationID: eval.ID,
		PerspectiveID:           persp.ID,
		IndicatorID:             &indicator.ID,
		Weight:                  30,
		Target:                  15,
	}
	if err := db.Create(detail).Error; err != nil {
		t.Fatalf("failed to seed detail: %v", err)
	}

	resp, err := svc.GetEvaluationWithDetails(context.Background(), eval.ID.String())
	if err != nil {
		t.Fatalf("GetEvaluationWithDetails failed: %v", err)
	}
	if len(resp.Details) != 1 {
		t.Fatalf("expected 1 detail, got %d", len(resp.Details))
	}
	d := resp.Details[0]
	if d.PerspectiveName != "Financial" {
		t.Errorf("expected perspective_name 'Financial', got %q", d.PerspectiveName)
	}
	if d.Target != 15 {
		t.Errorf("expected target 15, got %v", d.Target)
	}
	if d.UnitOfMeasurement != "%" {
		t.Errorf("expected unit_of_measurement '%%', got %q", d.UnitOfMeasurement)
	}
	if d.FormulaType != "HIGHER_BETTER" {
		t.Errorf("expected formula_type 'HIGHER_BETTER', got %q", d.FormulaType)
	}
}
