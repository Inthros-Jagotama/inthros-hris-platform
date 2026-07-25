package performance

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func newTestRepo() (*Repository, func()) {
	_, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	return repo, cleanup
}

// =========================================================================
// Performance Period Repository Tests
// =========================================================================

func TestRepo_CreateAndFindPerformancePeriod(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	p := &PerformancePeriod{
		PeriodCode: "H1",
		PeriodType: "SEMESTER",
		Year:       2026,
	}

	if err := repo.CreatePerformancePeriod(ctx, p); err != nil {
		t.Fatalf("CreatePerformancePeriod failed: %v", err)
	}
	if p.ID == uuid.Nil {
		t.Fatal("expected ID to be generated")
	}

	found, err := repo.FindPerformancePeriodByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("FindPerformancePeriodByID failed: %v", err)
	}
	if found.PeriodCode != "H1" {
		t.Errorf("expected 'H1', got '%s'", found.PeriodCode)
	}
}

func TestRepo_FindPerformancePeriodByID_NotFound(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	_, err := repo.FindPerformancePeriodByID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent period")
	}
}

func TestRepo_ListPerformancePeriods(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		repo.CreatePerformancePeriod(ctx, &PerformancePeriod{
			PeriodCode: "P",
			PeriodType: "MONTHLY",
			Year:       2026,
		})
	}

	list, total, err := repo.ListPerformancePeriods(ctx, 1, 10)
	if err != nil {
		t.Fatalf("ListPerformancePeriods failed: %v", err)
	}
	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(list) != 5 {
		t.Errorf("expected 5 items, got %d", len(list))
	}
}

func TestRepo_UpdatePerformancePeriod(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	p := &PerformancePeriod{PeriodCode: "OLD", PeriodType: "MONTHLY", Year: 2025}
	repo.CreatePerformancePeriod(ctx, p)

	p.PeriodCode = "NEW"
	if err := repo.UpdatePerformancePeriod(ctx, p); err != nil {
		t.Fatalf("UpdatePerformancePeriod failed: %v", err)
	}

	found, _ := repo.FindPerformancePeriodByID(ctx, p.ID)
	if found.PeriodCode != "NEW" {
		t.Errorf("expected 'NEW', got '%s'", found.PeriodCode)
	}
}

func TestRepo_DeletePerformancePeriod(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	p := &PerformancePeriod{PeriodCode: "DEL", PeriodType: "MONTHLY", Year: 2026}
	repo.CreatePerformancePeriod(ctx, p)

	if err := repo.DeletePerformancePeriod(ctx, p.ID); err != nil {
		t.Fatalf("DeletePerformancePeriod failed: %v", err)
	}
}

// =========================================================================
// Performance Perspective Repository Tests
// =========================================================================

func TestRepo_CreateAndFindPerspective(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	p := &PerformancePerspective{Name: "Customer", SortOrder: 2}
	if err := repo.CreatePerformancePerspective(ctx, p); err != nil {
		t.Fatalf("CreatePerformancePerspective failed: %v", err)
	}

	found, err := repo.FindPerformancePerspectiveByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("FindPerformancePerspectiveByID failed: %v", err)
	}
	if found.Name != "Customer" {
		t.Errorf("expected 'Customer', got '%s'", found.Name)
	}
}

func TestRepo_DeletePerformancePerspective(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	p := &PerformancePerspective{Name: "Delete Me"}
	repo.CreatePerformancePerspective(ctx, p)
	if err := repo.DeletePerformancePerspective(ctx, p.ID); err != nil {
		t.Fatalf("DeletePerformancePerspective failed: %v", err)
	}
}

// =========================================================================
// Performance Template Repository Tests
// =========================================================================

func TestRepo_CreateAndFindTemplate(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	tmpl := &PerformanceTemplate{
		OrganizationID: uuid.New(),
		Name:           "Executive Template",
	}
	if err := repo.CreatePerformanceTemplate(ctx, tmpl); err != nil {
		t.Fatalf("CreatePerformanceTemplate failed: %v", err)
	}

	found, err := repo.FindPerformanceTemplateByID(ctx, tmpl.ID)
	if err != nil {
		t.Fatalf("FindPerformanceTemplateByID failed: %v", err)
	}
	if found.Name != "Executive Template" {
		t.Errorf("expected 'Executive Template', got '%s'", found.Name)
	}
	if found.Status != "DRAFT" {
		t.Errorf("expected default 'DRAFT', got '%s'", found.Status)
	}
}

func TestRepo_ListTemplates_ByOrg(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	orgA := uuid.New()
	orgB := uuid.New()

	repo.CreatePerformanceTemplate(ctx, &PerformanceTemplate{OrganizationID: orgA, Name: "A-1"})
	repo.CreatePerformanceTemplate(ctx, &PerformanceTemplate{OrganizationID: orgA, Name: "A-2"})
	repo.CreatePerformanceTemplate(ctx, &PerformanceTemplate{OrganizationID: orgB, Name: "B-1"})

	list, total, err := repo.ListPerformanceTemplates(ctx, &orgA, 1, 10)
	if err != nil {
		t.Fatalf("ListPerformanceTemplates failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2 templates for org A, got %d", total)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}
}

// =========================================================================
// Performance Indicator Repository Tests
// =========================================================================

func TestRepo_CreateAndFindIndicator(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	tmpl := &PerformanceTemplate{OrganizationID: uuid.New(), Name: "T"}
	repo.CreatePerformanceTemplate(ctx, tmpl)

	ind := &PerformanceIndicator{
		PerformanceTemplateID: tmpl.ID,
		PerspectiveID:        uuid.New(),
		IndicatorType:        "MAXIMIZATION",
		Title:                "Customer Satisfaction",
		Weight:               25.0,
	}
	if err := repo.CreatePerformanceIndicator(ctx, ind); err != nil {
		t.Fatalf("CreatePerformanceIndicator failed: %v", err)
	}

	found, err := repo.FindPerformanceIndicatorByID(ctx, ind.ID)
	if err != nil {
		t.Fatalf("FindPerformanceIndicatorByID failed: %v", err)
	}
	if found.Title != "Customer Satisfaction" {
		t.Errorf("expected 'Customer Satisfaction', got '%s'", found.Title)
	}
}

func TestRepo_ListIndicators_ByTemplate(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	tmpl := &PerformanceTemplate{OrganizationID: uuid.New(), Name: "T"}
	repo.CreatePerformanceTemplate(ctx, tmpl)

	repo.CreatePerformanceIndicator(ctx, &PerformanceIndicator{
		PerformanceTemplateID: tmpl.ID, PerspectiveID: uuid.New(),
		IndicatorType: "MAXIMIZATION", Title: "KPI 1", Weight: 50,
	})
	repo.CreatePerformanceIndicator(ctx, &PerformanceIndicator{
		PerformanceTemplateID: tmpl.ID, PerspectiveID: uuid.New(),
		IndicatorType: "MINIMIZATION", Title: "KPI 2", Weight: 50,
	})

	list, total, err := repo.ListPerformanceIndicators(ctx, tmpl.ID, 1, 10)
	if err != nil {
		t.Fatalf("ListPerformanceIndicators failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2 indicators, got %d", total)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}
}

// =========================================================================
// Performance Evaluation Repository Tests
// =========================================================================

func TestRepo_CreateAndFindEvaluation(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	eval := &PerformanceEvaluation{
		EmployeeID:     uuid.New(),
		OrganizationID: uuid.New(),
		PeriodID:       uuid.New(),
		TemplateID:     uuid.New(),
		Status:         "DRAFT",
	}
	if err := repo.CreatePerformanceEvaluation(ctx, eval); err != nil {
		t.Fatalf("CreatePerformanceEvaluation failed: %v", err)
	}

	found, err := repo.FindPerformanceEvaluationByID(ctx, eval.ID)
	if err != nil {
		t.Fatalf("FindPerformanceEvaluationByID failed: %v", err)
	}
	if found.Status != "DRAFT" {
		t.Errorf("expected 'DRAFT', got '%s'", found.Status)
	}
}

func TestRepo_ListEvaluations_ByEmployee(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	emp := uuid.New()
	period := uuid.New()

	repo.CreatePerformanceEvaluation(ctx, &PerformanceEvaluation{
		EmployeeID: emp, OrganizationID: uuid.New(), PeriodID: period, TemplateID: uuid.New(), Status: "DRAFT",
	})
	repo.CreatePerformanceEvaluation(ctx, &PerformanceEvaluation{
		EmployeeID: emp, OrganizationID: uuid.New(), PeriodID: period, TemplateID: uuid.New(), Status: "PLAN_SUBMITTED",
	})
	repo.CreatePerformanceEvaluation(ctx, &PerformanceEvaluation{
		EmployeeID: uuid.New(), OrganizationID: uuid.New(), PeriodID: period, TemplateID: uuid.New(), Status: "DRAFT",
	})

	list, total, err := repo.ListPerformanceEvaluations(ctx, &emp, nil, nil, 1, 10)
	if err != nil {
		t.Fatalf("ListPerformanceEvaluations failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2 evaluations for employee, got %d", total)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}
}

func TestRepo_DeletePerformanceEvaluation(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	eval := &PerformanceEvaluation{
		EmployeeID: uuid.New(), OrganizationID: uuid.New(),
		PeriodID: uuid.New(), TemplateID: uuid.New(), Status: "DRAFT",
	}
	repo.CreatePerformanceEvaluation(ctx, eval)

	if err := repo.DeletePerformanceEvaluation(ctx, eval.ID); err != nil {
		t.Fatalf("DeletePerformanceEvaluation failed: %v", err)
	}
}

// =========================================================================
// Evaluation Detail Repository Tests
// =========================================================================

func TestRepo_CreateAndFindEvaluationDetail(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	eval := &PerformanceEvaluation{
		EmployeeID: uuid.New(), OrganizationID: uuid.New(),
		PeriodID: uuid.New(), TemplateID: uuid.New(), Status: "DRAFT",
	}
	repo.CreatePerformanceEvaluation(ctx, eval)

	d := &PerformanceEvaluationDetail{
		PerformanceEvaluationID: eval.ID,
		PerspectiveID:          uuid.New(),
		AchievementPercentage:  90.0,
		Weight:                 50.0,
		Score:                  45.0,
	}
	if err := repo.CreateEvaluationDetail(ctx, d); err != nil {
		t.Fatalf("CreateEvaluationDetail failed: %v", err)
	}

	found, err := repo.FindEvaluationDetailByID(ctx, d.ID)
	if err != nil {
		t.Fatalf("FindEvaluationDetailByID failed: %v", err)
	}
	if found.Score != 45.0 {
		t.Errorf("expected score 45.0, got %f", found.Score)
	}
}

func TestRepo_ListEvaluationDetails(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	eval := &PerformanceEvaluation{
		EmployeeID: uuid.New(), OrganizationID: uuid.New(),
		PeriodID: uuid.New(), TemplateID: uuid.New(), Status: "DRAFT",
	}
	repo.CreatePerformanceEvaluation(ctx, eval)

	repo.CreateEvaluationDetail(ctx, &PerformanceEvaluationDetail{
		PerformanceEvaluationID: eval.ID, PerspectiveID: uuid.New(), Weight: 40, Score: 35,
	})
	repo.CreateEvaluationDetail(ctx, &PerformanceEvaluationDetail{
		PerformanceEvaluationID: eval.ID, PerspectiveID: uuid.New(), Weight: 60, Score: 50,
	})

	list, err := repo.ListEvaluationDetails(ctx, eval.ID)
	if err != nil {
		t.Fatalf("ListEvaluationDetails failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 details, got %d", len(list))
	}
}

// =========================================================================
// Performance Target Repository Tests
// =========================================================================

func TestRepo_CreateAndFindTarget(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	eval := &PerformanceEvaluation{
		EmployeeID: uuid.New(), OrganizationID: uuid.New(),
		PeriodID: uuid.New(), TemplateID: uuid.New(), Status: "DRAFT",
	}
	repo.CreatePerformanceEvaluation(ctx, eval)

	target := &PerformanceTarget{
		PerformanceEvaluationID: eval.ID,
		IndicatorID:            uuid.New(),
		TargetValue:            5000000,
		Weight:                 100.0,
	}
	if err := repo.CreatePerformanceTarget(ctx, target); err != nil {
		t.Fatalf("CreatePerformanceTarget failed: %v", err)
	}

	found, err := repo.FindPerformanceTargetByID(ctx, target.ID)
	if err != nil {
		t.Fatalf("FindPerformanceTargetByID failed: %v", err)
	}
	if found.TargetValue != 5000000 {
		t.Errorf("expected target 5000000, got %f", found.TargetValue)
	}
}

func TestRepo_ListPerformanceTargets(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	eval := &PerformanceEvaluation{
		EmployeeID: uuid.New(), OrganizationID: uuid.New(),
		PeriodID: uuid.New(), TemplateID: uuid.New(), Status: "DRAFT",
	}
	repo.CreatePerformanceEvaluation(ctx, eval)

	repo.CreatePerformanceTarget(ctx, &PerformanceTarget{
		PerformanceEvaluationID: eval.ID, IndicatorID: uuid.New(), TargetValue: 100, Weight: 50,
	})
	repo.CreatePerformanceTarget(ctx, &PerformanceTarget{
		PerformanceEvaluationID: eval.ID, IndicatorID: uuid.New(), TargetValue: 200, Weight: 50,
	})

	list, total, err := repo.ListPerformanceTargets(ctx, eval.ID)
	if err != nil {
		t.Fatalf("ListPerformanceTargets failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2 targets, got %d", total)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}
}

func TestRepo_UpdateEvaluationFinalScore(t *testing.T) {
	repo, cleanup := newTestRepo()
	defer cleanup()
	ctx := context.Background()

	eval := &PerformanceEvaluation{
		EmployeeID: uuid.New(), OrganizationID: uuid.New(),
		PeriodID: uuid.New(), TemplateID: uuid.New(), Status: "DRAFT",
	}
	repo.CreatePerformanceEvaluation(ctx, eval)

	repo.CreateEvaluationDetail(ctx, &PerformanceEvaluationDetail{
		PerformanceEvaluationID: eval.ID, PerspectiveID: uuid.New(), Score: 40, Weight: 50,
	})
	repo.CreateEvaluationDetail(ctx, &PerformanceEvaluationDetail{
		PerformanceEvaluationID: eval.ID, PerspectiveID: uuid.New(), Score: 35, Weight: 50,
	})

	score, err := repo.UpdateEvaluationFinalScore(ctx, eval.ID)
	if err != nil {
		t.Fatalf("UpdateEvaluationFinalScore failed: %v", err)
	}
	if score != 75.0 {
		t.Errorf("expected total score 75.0, got %f", score)
	}
}
