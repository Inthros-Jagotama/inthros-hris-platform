package performance

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func newTestService() (*Service, func()) {
	_, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	return svc, func() { cleanup(); logger.Sync() }
}

// =========================================================================
// Performance Period Service Tests
// =========================================================================

func TestService_CreatePerformancePeriod(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreatePerformancePeriodRequest{
		PeriodCode: "Q1",
		PeriodType: "QUARTERLY",
		Year:       2026,
		StartDate:  strPtr("2026-01-01"),
		EndDate:    strPtr("2026-03-31"),
	}

	resp, err := svc.CreatePerformancePeriod(ctx, req)
	if err != nil {
		t.Fatalf("CreatePerformancePeriod failed: %v", err)
	}
	if resp.PeriodCode != "Q1" {
		t.Errorf("expected period_code 'Q1', got '%s'", resp.PeriodCode)
	}
	if resp.PeriodType != "QUARTERLY" {
		t.Errorf("expected period_type 'QUARTERLY', got '%s'", resp.PeriodType)
	}
	if resp.Year != 2026 {
		t.Errorf("expected year 2026, got %d", resp.Year)
	}
}

func TestService_GetPerformancePeriodByID_Success(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreatePerformancePeriod(ctx, CreatePerformancePeriodRequest{
		PeriodCode: "FY",
		PeriodType: "ANNUAL",
		Year:       2026,
	})

	found, err := svc.GetPerformancePeriodByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetPerformancePeriodByID failed: %v", err)
	}
	if found.PeriodCode != "FY" {
		t.Errorf("expected period_code 'FY', got '%s'", found.PeriodCode)
	}
}

func TestService_GetPerformancePeriodByID_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetPerformancePeriodByID(ctx, "bad-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_GetPerformancePeriodByID_NotFound(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetPerformancePeriodByID(ctx, uuid.New().String())
	if err == nil {
		t.Fatal("expected error for non-existent period")
	}
}

func TestService_ListPerformancePeriods_DefaultPagination(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		svc.CreatePerformancePeriod(ctx, CreatePerformancePeriodRequest{
			PeriodCode: "P",
			PeriodType: "MONTHLY",
			Year:       2026,
		})
	}

	resp, err := svc.ListPerformancePeriods(ctx, 0, 0)
	if err != nil {
		t.Fatalf("ListPerformancePeriods failed: %v", err)
	}
	if resp.Page != 1 {
		t.Errorf("expected page 1, got %d", resp.Page)
	}
	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}
}

func TestService_UpdatePerformancePeriod(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreatePerformancePeriod(ctx, CreatePerformancePeriodRequest{
		PeriodCode: "OLD",
		PeriodType: "QUARTERLY",
		Year:       2025,
	})

	newCode := "NEW"
	updated, err := svc.UpdatePerformancePeriod(ctx, created.ID, UpdatePerformancePeriodRequest{
		PeriodCode: &newCode,
	})
	if err != nil {
		t.Fatalf("UpdatePerformancePeriod failed: %v", err)
	}
	if updated.PeriodCode != "NEW" {
		t.Errorf("expected period_code 'NEW', got '%s'", updated.PeriodCode)
	}
}

func TestService_DeletePerformancePeriod(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreatePerformancePeriod(ctx, CreatePerformancePeriodRequest{
		PeriodCode: "DEL",
		PeriodType: "MONTHLY",
		Year:       2026,
	})

	if err := svc.DeletePerformancePeriod(ctx, created.ID); err != nil {
		t.Fatalf("DeletePerformancePeriod failed: %v", err)
	}
	_, err := svc.GetPerformancePeriodByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting period")
	}
}

// =========================================================================
// Performance Perspective Service Tests
// =========================================================================

func TestService_CreatePerformancePerspective(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreatePerformancePerspectiveRequest{
		Name:        "Financial",
		Description: strPtr("Financial metrics"),
		SortOrder:   intPtr(1),
	}

	resp, err := svc.CreatePerformancePerspective(ctx, req)
	if err != nil {
		t.Fatalf("CreatePerformancePerspective failed: %v", err)
	}
	if resp.Name != "Financial" {
		t.Errorf("expected name 'Financial', got '%s'", resp.Name)
	}
	if resp.SortOrder != 1 {
		t.Errorf("expected sort_order 1, got %d", resp.SortOrder)
	}
}

func TestService_UpdatePerformancePerspective(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreatePerformancePerspective(ctx, CreatePerformancePerspectiveRequest{
		Name: "Before",
	})

	newName := "After"
	updated, err := svc.UpdatePerformancePerspective(ctx, created.ID, UpdatePerformancePerspectiveRequest{
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("UpdatePerformancePerspective failed: %v", err)
	}
	if updated.Name != "After" {
		t.Errorf("expected name 'After', got '%s'", updated.Name)
	}
}

func TestService_DeletePerformancePerspective(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreatePerformancePerspective(ctx, CreatePerformancePerspectiveRequest{
		Name: "To Delete",
	})
	if err := svc.DeletePerformancePerspective(ctx, created.ID); err != nil {
		t.Fatalf("DeletePerformancePerspective failed: %v", err)
	}
	_, err := svc.GetPerformancePerspectiveByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting perspective")
	}
}

// =========================================================================
// Performance Template Service Tests
// =========================================================================

func TestService_CreatePerformanceTemplate(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreatePerformanceTemplateRequest{
		OrganizationID: createTestOrgID(),
		Name:           "Manager BSC Template",
		Description:    strPtr("Standard template for managers"),
	}

	resp, err := svc.CreatePerformanceTemplate(ctx, req)
	if err != nil {
		t.Fatalf("CreatePerformanceTemplate failed: %v", err)
	}
	if resp.Name != "Manager BSC Template" {
		t.Errorf("expected name 'Manager BSC Template', got '%s'", resp.Name)
	}
	if resp.Status != "DRAFT" {
		t.Errorf("expected default status 'DRAFT', got '%s'", resp.Status)
	}
}

func TestService_UpdatePerformanceTemplate_Status(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreatePerformanceTemplate(ctx, CreatePerformanceTemplateRequest{
		OrganizationID: createTestOrgID(),
		Name:           "Test Template",
	})

	published := "PUBLISHED"
	updated, err := svc.UpdatePerformanceTemplate(ctx, created.ID, UpdatePerformanceTemplateRequest{
		Status: &published,
	})
	if err != nil {
		t.Fatalf("UpdatePerformanceTemplate failed: %v", err)
	}
	if updated.Status != "PUBLISHED" {
		t.Errorf("expected status 'PUBLISHED', got '%s'", updated.Status)
	}
}

func TestService_UpdatePerformanceTemplate_OrganizationID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreatePerformanceTemplate(ctx, CreatePerformanceTemplateRequest{
		OrganizationID: createTestOrgID(),
		Name:           "Test Template",
	})

	newOrgID := createTestOrgID()
	updated, err := svc.UpdatePerformanceTemplate(ctx, created.ID, UpdatePerformanceTemplateRequest{
		OrganizationID: &newOrgID,
	})
	if err != nil {
		t.Fatalf("UpdatePerformanceTemplate failed: %v", err)
	}
	if updated.OrganizationID != newOrgID {
		t.Errorf("expected organization_id %q, got %q", newOrgID, updated.OrganizationID)
	}

	// Verify persistence via a fresh repository read (avoids GetPerformanceTemplateByID's
	// enrichTemplateResponses, which needs the raw "organizations" table this shared
	// test setup doesn't provision).
	stored, err := svc.repo.FindPerformanceTemplateByID(ctx, uuid.MustParse(created.ID))
	if err != nil {
		t.Fatalf("FindPerformanceTemplateByID failed: %v", err)
	}
	if stored.OrganizationID.String() != newOrgID {
		t.Errorf("expected persisted organization_id %q, got %q", newOrgID, stored.OrganizationID.String())
	}
}

func TestService_DeletePerformanceTemplate(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreatePerformanceTemplate(ctx, CreatePerformanceTemplateRequest{
		OrganizationID: createTestOrgID(),
		Name:           "To Delete",
	})
	if err := svc.DeletePerformanceTemplate(ctx, created.ID); err != nil {
		t.Fatalf("DeletePerformanceTemplate failed: %v", err)
	}
	_, err := svc.GetPerformanceTemplateByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting template")
	}
}

// =========================================================================
// Performance Indicator Service Tests
// =========================================================================

func TestService_CreatePerformanceIndicator(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	tmpl, _ := svc.CreatePerformanceTemplate(ctx, CreatePerformanceTemplateRequest{
		OrganizationID: createTestOrgID(),
		Name:           "Template",
	})
	persp, _ := svc.CreatePerformancePerspective(ctx, CreatePerformancePerspectiveRequest{
		Name: "Financial",
	})

	req := CreatePerformanceIndicatorRequest{
		PerformanceTemplateID: tmpl.ID,
		PerspectiveID:        persp.ID,
		IndicatorType:        "MAXIMIZATION",
		Title:                "Revenue Growth",
		Weight:               30.0,
		TargetValue:          15.0,
		UnitOfMeasurement:    strPtr("%"),
	}

	resp, err := svc.CreatePerformanceIndicator(ctx, req)
	if err != nil {
		t.Fatalf("CreatePerformanceIndicator failed: %v", err)
	}
	if resp.Title != "Revenue Growth" {
		t.Errorf("expected title 'Revenue Growth', got '%s'", resp.Title)
	}
	if resp.IndicatorType != "MAXIMIZATION" {
		t.Errorf("expected type 'MAXIMIZATION', got '%s'", resp.IndicatorType)
	}
	if resp.Weight != 30.0 {
		t.Errorf("expected weight 30.0, got %f", resp.Weight)
	}
}

func TestService_GetPerformanceIndicatorByID_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetPerformanceIndicatorByID(ctx, "bad-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_UpdatePerformanceIndicator(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	tmpl, _ := svc.CreatePerformanceTemplate(ctx, CreatePerformanceTemplateRequest{
		OrganizationID: createTestOrgID(),
		Name:           "T",
	})
	persp, _ := svc.CreatePerformancePerspective(ctx, CreatePerformancePerspectiveRequest{
		Name: "Customer",
	})
	created, _ := svc.CreatePerformanceIndicator(ctx, CreatePerformanceIndicatorRequest{
		PerformanceTemplateID: tmpl.ID,
		PerspectiveID:        persp.ID,
		IndicatorType:        "MAXIMIZATION",
		Title:                "Before",
	})

	newTitle := "After"
	updated, err := svc.UpdatePerformanceIndicator(ctx, created.ID, UpdatePerformanceIndicatorRequest{
		Title: &newTitle,
	})
	if err != nil {
		t.Fatalf("UpdatePerformanceIndicator failed: %v", err)
	}
	if updated.Title != "After" {
		t.Errorf("expected title 'After', got '%s'", updated.Title)
	}
}

func TestService_UpdatePerformanceIndicator_PerspectiveID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	tmpl, _ := svc.CreatePerformanceTemplate(ctx, CreatePerformanceTemplateRequest{
		OrganizationID: createTestOrgID(),
		Name:           "T",
	})
	perspA, _ := svc.CreatePerformancePerspective(ctx, CreatePerformancePerspectiveRequest{Name: "Financial"})
	perspB, _ := svc.CreatePerformancePerspective(ctx, CreatePerformancePerspectiveRequest{Name: "Customer"})
	created, _ := svc.CreatePerformanceIndicator(ctx, CreatePerformanceIndicatorRequest{
		PerformanceTemplateID: tmpl.ID,
		PerspectiveID:         perspA.ID,
		IndicatorType:         "MAXIMIZATION",
		Title:                 "Some Indicator",
	})

	updated, err := svc.UpdatePerformanceIndicator(ctx, created.ID, UpdatePerformanceIndicatorRequest{
		PerspectiveID: &perspB.ID,
	})
	if err != nil {
		t.Fatalf("UpdatePerformanceIndicator failed: %v", err)
	}
	if updated.PerspectiveID != perspB.ID {
		t.Errorf("expected perspective_id %q, got %q", perspB.ID, updated.PerspectiveID)
	}

	refetched, err := svc.GetPerformanceIndicatorByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetPerformanceIndicatorByID failed: %v", err)
	}
	if refetched.PerspectiveID != perspB.ID {
		t.Errorf("expected persisted perspective_id %q, got %q", perspB.ID, refetched.PerspectiveID)
	}
}

func TestService_DeletePerformanceIndicator(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	tmpl, _ := svc.CreatePerformanceTemplate(ctx, CreatePerformanceTemplateRequest{
		OrganizationID: createTestOrgID(),
		Name:           "T",
	})
	persp, _ := svc.CreatePerformancePerspective(ctx, CreatePerformancePerspectiveRequest{
		Name: "Internal",
	})
	created, _ := svc.CreatePerformanceIndicator(ctx, CreatePerformanceIndicatorRequest{
		PerformanceTemplateID: tmpl.ID,
		PerspectiveID:        persp.ID,
		IndicatorType:        "MINIMIZATION",
		Title:                "To Delete",
	})
	if err := svc.DeletePerformanceIndicator(ctx, created.ID); err != nil {
		t.Fatalf("DeletePerformanceIndicator failed: %v", err)
	}
	_, err := svc.GetPerformanceIndicatorByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting indicator")
	}
}

// =========================================================================
// Performance Evaluation Service Tests
// =========================================================================

func TestService_CreatePerformanceEvaluation(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	period, _ := svc.CreatePerformancePeriod(ctx, CreatePerformancePeriodRequest{
		PeriodCode: "FY", PeriodType: "ANNUAL", Year: 2026,
	})
	tmpl, _ := svc.CreatePerformanceTemplate(ctx, CreatePerformanceTemplateRequest{
		OrganizationID: createTestOrgID(), Name: "Template",
	})

	req := CreatePerformanceEvaluationRequest{
		EmployeeID:     createTestUUID(),
		OrganizationID: createTestOrgID(),
		PeriodID:       period.ID,
		TemplateID:     tmpl.ID,
		Notes:          strPtr("Annual review"),
	}

	resp, err := svc.CreatePerformanceEvaluation(ctx, req)
	if err != nil {
		t.Fatalf("CreatePerformanceEvaluation failed: %v", err)
	}
	if resp.Status != "DRAFT" {
		t.Errorf("expected status 'DRAFT', got '%s'", resp.Status)
	}
}

func TestService_UpdateEvaluationStatus_InvalidTransition(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	period, _ := svc.CreatePerformancePeriod(ctx, CreatePerformancePeriodRequest{
		PeriodCode: "FY", PeriodType: "ANNUAL", Year: 2026,
	})
	tmpl, _ := svc.CreatePerformanceTemplate(ctx, CreatePerformanceTemplateRequest{
		OrganizationID: createTestOrgID(), Name: "T",
	})
	eval, _ := svc.CreatePerformanceEvaluation(ctx, CreatePerformanceEvaluationRequest{
		EmployeeID:     createTestUUID(),
		OrganizationID: createTestOrgID(),
		PeriodID:       period.ID,
		TemplateID:     tmpl.ID,
	})

	// Can't go from DRAFT directly to COMPLETED
	_, err := svc.UpdateEvaluationStatus(ctx, eval.ID, "COMPLETED", "")
	if err == nil {
		t.Fatal("expected error for invalid status transition")
	}
}

func TestService_UpdateEvaluationStatus_ValidTransition(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	period, _ := svc.CreatePerformancePeriod(ctx, CreatePerformancePeriodRequest{
		PeriodCode: "FY", PeriodType: "ANNUAL", Year: 2026,
	})
	tmpl, _ := svc.CreatePerformanceTemplate(ctx, CreatePerformanceTemplateRequest{
		OrganizationID: createTestOrgID(), Name: "T",
	})
	eval, _ := svc.CreatePerformanceEvaluation(ctx, CreatePerformanceEvaluationRequest{
		EmployeeID:     createTestUUID(),
		OrganizationID: createTestOrgID(),
		PeriodID:       period.ID,
		TemplateID:     tmpl.ID,
	})

	updated, err := svc.UpdateEvaluationStatus(ctx, eval.ID, "PLAN_SUBMITTED", "Ready for review")
	if err != nil {
		t.Fatalf("UpdateEvaluationStatus failed: %v", err)
	}
	if updated.Status != "PLAN_SUBMITTED" {
		t.Errorf("expected status 'PLAN_SUBMITTED', got '%s'", updated.Status)
	}
}

func TestService_DeletePerformanceEvaluation(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	period, _ := svc.CreatePerformancePeriod(ctx, CreatePerformancePeriodRequest{
		PeriodCode: "FY", PeriodType: "ANNUAL", Year: 2026,
	})
	tmpl, _ := svc.CreatePerformanceTemplate(ctx, CreatePerformanceTemplateRequest{
		OrganizationID: createTestOrgID(), Name: "T",
	})
	eval, _ := svc.CreatePerformanceEvaluation(ctx, CreatePerformanceEvaluationRequest{
		EmployeeID:     createTestUUID(),
		OrganizationID: createTestOrgID(),
		PeriodID:       period.ID,
		TemplateID:     tmpl.ID,
	})

	if err := svc.DeletePerformanceEvaluation(ctx, eval.ID); err != nil {
		t.Fatalf("DeletePerformanceEvaluation failed: %v", err)
	}
	_, err := svc.GetPerformanceEvaluationByID(ctx, eval.ID)
	if err == nil {
		t.Fatal("expected error after deleting evaluation")
	}
}

// =========================================================================
// Evaluation Detail & Target Service Tests
// =========================================================================

func TestService_CreateEvaluationDetail(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	period, _ := svc.CreatePerformancePeriod(ctx, CreatePerformancePeriodRequest{
		PeriodCode: "FY", PeriodType: "ANNUAL", Year: 2026,
	})
	tmpl, _ := svc.CreatePerformanceTemplate(ctx, CreatePerformanceTemplateRequest{
		OrganizationID: createTestOrgID(), Name: "T",
	})
	persp, _ := svc.CreatePerformancePerspective(ctx, CreatePerformancePerspectiveRequest{
		Name: "Financial",
	})
	eval, _ := svc.CreatePerformanceEvaluation(ctx, CreatePerformanceEvaluationRequest{
		EmployeeID:     createTestUUID(),
		OrganizationID: createTestOrgID(),
		PeriodID:       period.ID,
		TemplateID:     tmpl.ID,
	})

	req := CreateEvaluationDetailRequest{
		PerformanceEvaluationID: eval.ID,
		PerspectiveID:          persp.ID,
		AchievementPercentage:  85.0,
		Weight:                 40.0,
		Score:                  34.0,
	}

	resp, err := svc.CreateEvaluationDetail(ctx, req)
	if err != nil {
		t.Fatalf("CreateEvaluationDetail failed: %v", err)
	}
	if resp.AchievementPercentage != 85.0 {
		t.Errorf("expected achievement 85.0, got %f", resp.AchievementPercentage)
	}
	if resp.Score != 34.0 {
		t.Errorf("expected score 34.0, got %f", resp.Score)
	}
}

func TestService_CreatePerformanceTarget(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	period, _ := svc.CreatePerformancePeriod(ctx, CreatePerformancePeriodRequest{
		PeriodCode: "FY", PeriodType: "ANNUAL", Year: 2026,
	})
	tmpl, _ := svc.CreatePerformanceTemplate(ctx, CreatePerformanceTemplateRequest{
		OrganizationID: createTestOrgID(), Name: "T",
	})
	persp, _ := svc.CreatePerformancePerspective(ctx, CreatePerformancePerspectiveRequest{
		Name: "Financial",
	})
	eval, _ := svc.CreatePerformanceEvaluation(ctx, CreatePerformanceEvaluationRequest{
		EmployeeID:     createTestUUID(),
		OrganizationID: createTestOrgID(),
		PeriodID:       period.ID,
		TemplateID:     tmpl.ID,
	})
	ind, _ := svc.CreatePerformanceIndicator(ctx, CreatePerformanceIndicatorRequest{
		PerformanceTemplateID: tmpl.ID,
		PerspectiveID:        persp.ID,
		IndicatorType:        "MAXIMIZATION",
		Title:                "Revenue",
		Weight:               100.0,
	})

	req := CreatePerformanceTargetRequest{
		PerformanceEvaluationID: eval.ID,
		IndicatorID:            ind.ID,
		TargetValue:            1000000,
		UnitOfMeasurement:      strPtr("IDR"),
		Weight:                 50.0,
	}

	resp, err := svc.CreatePerformanceTarget(ctx, req)
	if err != nil {
		t.Fatalf("CreatePerformanceTarget failed: %v", err)
	}
	if resp.TargetValue != 1000000 {
		t.Errorf("expected target 1000000, got %f", resp.TargetValue)
	}
	if resp.Weight != 50.0 {
		t.Errorf("expected weight 50.0, got %f", resp.Weight)
	}
}

func TestService_ListPerformanceTargets(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	period, _ := svc.CreatePerformancePeriod(ctx, CreatePerformancePeriodRequest{
		PeriodCode: "FY", PeriodType: "ANNUAL", Year: 2026,
	})
	tmpl, _ := svc.CreatePerformanceTemplate(ctx, CreatePerformanceTemplateRequest{
		OrganizationID: createTestOrgID(), Name: "T",
	})
	persp, _ := svc.CreatePerformancePerspective(ctx, CreatePerformancePerspectiveRequest{
		Name: "Financial",
	})
	eval, _ := svc.CreatePerformanceEvaluation(ctx, CreatePerformanceEvaluationRequest{
		EmployeeID:     createTestUUID(),
		OrganizationID: createTestOrgID(),
		PeriodID:       period.ID,
		TemplateID:     tmpl.ID,
	})
	ind1, _ := svc.CreatePerformanceIndicator(ctx, CreatePerformanceIndicatorRequest{
		PerformanceTemplateID: tmpl.ID,
		PerspectiveID:        persp.ID,
		IndicatorType:        "MAXIMIZATION",
		Title:                "Revenue",
		Weight:               50.0,
	})
	ind2, _ := svc.CreatePerformanceIndicator(ctx, CreatePerformanceIndicatorRequest{
		PerformanceTemplateID: tmpl.ID,
		PerspectiveID:        persp.ID,
		IndicatorType:        "MINIMIZATION",
		Title:                "Cost",
		Weight:               50.0,
	})

	svc.CreatePerformanceTarget(ctx, CreatePerformanceTargetRequest{
		PerformanceEvaluationID: eval.ID,
		IndicatorID:            ind1.ID,
		TargetValue:            1000000,
		Weight:                 50.0,
	})
	svc.CreatePerformanceTarget(ctx, CreatePerformanceTargetRequest{
		PerformanceEvaluationID: eval.ID,
		IndicatorID:            ind2.ID,
		TargetValue:            500000,
		Weight:                 50.0,
	})

	targets, err := svc.ListPerformanceTargets(ctx, eval.ID)
	if err != nil {
		t.Fatalf("ListPerformanceTargets failed: %v", err)
	}
	if len(targets) != 2 {
		t.Errorf("expected 2 targets, got %d", len(targets))
	}
}
