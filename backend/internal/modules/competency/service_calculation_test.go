package competency

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// =========================================================================
// Calculation Engine Tests (§14-§18)
// =========================================================================

// setup360Scenario membangun skenario assessment 360 lengkap di SQLite:
// competency master, indicator, template (competencies + rater types +
// indicators), event, target, rater, dan response. Mengembalikan service
// beserta ID-ID yang dibutuhkan untuk CalculateTarget.
func setup360Scenario(t *testing.T) (*Service, string, *AssessmentResult, func()) {
	t.Helper()
	_, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	ctx := context.Background()

	// Competency master + indicator
	comp, err := svc.CreateCompetency(ctx, CreateCompetencyRequest{Name: "Leadership"})
	if err != nil {
		t.Fatalf("create competency: %v", err)
	}
	compUID := mustParseUUID(t, comp.ID)
	ind, err := svc.CreateIndicator(ctx, CreateIndicatorRequest{
		CompetencyID: comp.ID,
		Statement:    "Memimpin tim dengan jelas",
	})
	if err != nil {
		t.Fatalf("create indicator: %v", err)
	}
	indUID := mustParseUUID(t, ind.ID)

	// Template: 1 competency (required level 3, weight 1) + rater types
	// (self weight 0.2, superior weight 0.5, peer weight 0.3) + indicator.
	reqLevel := 3
	tpl, err := svc.CreateAssessmentTemplate(ctx, CreateAssessmentTemplateRequest{
		Name: "Template 360",
		Code: "TPL-360-001",
		Competencies: []TemplateCompetencyRequest{
			{CompetencyID: comp.ID, RequiredLevel: &reqLevel, Weight: 1, SortOrder: 0},
		},
		RaterTypes: []TemplateRaterTypeRequest{
			{RaterType: "self", Weight: 0.2, MinRater: 1},
			{RaterType: "superior", Weight: 0.5, MinRater: 1},
			{RaterType: "peer", Weight: 0.3, MinRater: 1},
		},
	})
	if err != nil {
		t.Fatalf("create template: %v", err)
	}
	if _, err := svc.SetTemplateIndicators(ctx, tpl.ID, []TemplateIndicatorRequest{
		{IndicatorID: ind.ID, Weight: 1, SortOrder: 0},
	}); err != nil {
		t.Fatalf("set template indicators: %v", err)
	}

	// Event + target (subject employee)
	event, err := svc.CreateCompetencyEvent(ctx, CreateCompetencyEventRequest{
		Type:       "manual",
		PeriodType: "annual",
		PeriodYear: 2026,
		TemplateID: &tpl.ID,
	})
	if err != nil {
		t.Fatalf("create event: %v", err)
	}
	subjectID := uuid.New()
	target, err := svc.CreateCompetencyEventTarget(ctx, CreateCompetencyEventTargetRequest{
		CompetencyEventID: event.ID,
		OrganizationID:    createTestOrgID(),
		EmployeeID:        uuidStrPtr(subjectID),
	})
	if err != nil {
		t.Fatalf("create target: %v", err)
	}
	targetUID := mustParseUUID(t, target.ID)

	// Rater: self (employee = subject), superior, peer — semua sudah submit.
	now := time.Now()
	raters := []struct {
		rtype  string
		rating int
	}{
		{"self", 4},
		{"superior", 5},
		{"peer", 3},
	}
	for _, r := range raters {
		rat := &CompetencyAssessmentRater{
			CompetencyEventTargetID: targetUID,
			RaterEmployeeID:         uuid.New(),
			RaterType:               r.rtype,
			Weight:                  1,
			Status:                  "submitted",
			SubmittedAt:             &now,
		}
		if err := repo.CreateRater(ctx, rat); err != nil {
			t.Fatalf("create rater %s: %v", r.rtype, err)
		}
		if err := repo.SaveResponse(ctx, &CompetencyAssessmentResponse{
			RaterID:     rat.ID,
			IndicatorID: indUID,
			RatingValue: r.rating,
			Comment:     strPtr("ok"),
			SubmittedAt: &now,
		}); err != nil {
			t.Fatalf("create response for %s: %v", r.rtype, err)
		}
	}

	result, err := svc.CalculateTarget(ctx, target.ID)
	if err != nil {
		t.Fatalf("CalculateTarget failed: %v", err)
	}
	_ = compUID
	return svc, target.ID, result, func() { cleanup(); _ = logger.Sync() }
}

func subjectIDStr(id uuid.UUID) string { return id.String() }

func uuidStrPtr(id uuid.UUID) *string {
	s := id.String()
	return &s
}

func mustParseUUID(t *testing.T, s string) uuid.UUID {
	t.Helper()
	uid, err := uuid.Parse(s)
	if err != nil {
		t.Fatalf("parse uuid %q: %v", s, err)
	}
	return uid
}

// TestCalculateTarget_RaterAggregation memverifikasi agregasi response rater:
// skor per tipe rater = rata-rata rating rater, dan competency score adalah
// weighted average antar tipe rater.
func TestCalculateTarget_RaterAggregation(t *testing.T) {
	_, _, result, cleanup := setup360Scenario(t)
	defer cleanup()

	if len(result.Competencies) != 1 {
		t.Fatalf("expected 1 competency in result, got %d", len(result.Competencies))
	}
	c := result.Competencies[0]

	// Skor per tipe rater (masing-masing 1 rater, rating 4/5/3).
	if c.RaterScores.Self != 4 {
		t.Errorf("expected self score 4, got %v", c.RaterScores.Self)
	}
	if c.RaterScores.Superior != 5 {
		t.Errorf("expected superior score 5, got %v", c.RaterScores.Superior)
	}
	if c.RaterScores.Peer != 3 {
		t.Errorf("expected peer score 3, got %v", c.RaterScores.Peer)
	}

	// Competency score = (4*0.2 + 5*0.5 + 3*0.3) / (0.2+0.5+3.0?) → (0.8+2.5+0.9)/1.0 = 4.2
	if c.Score != 4.2 {
		t.Errorf("expected weighted competency score 4.2, got %v", c.Score)
	}
	// Gap = score - required level = 4.2 - 3 = 1.2
	if c.Gap != 1.2 {
		t.Errorf("expected gap 1.2, got %v", c.Gap)
	}
}

// TestCalculateTarget_GapAndOverall memverifikasi gap, weighted gap, dan
// overall score (weighted average seluruh competency).
func TestCalculateTarget_GapAndOverall(t *testing.T) {
	_, _, result, cleanup := setup360Scenario(t)
	defer cleanup()

	if result.OverallScore != 4.2 {
		t.Errorf("expected overall score 4.2, got %v", result.OverallScore)
	}
	if result.TotalGap != 1.2 {
		t.Errorf("expected total gap 1.2, got %v", result.TotalGap)
	}
	if result.SelfScore != 4 {
		t.Errorf("expected self score 4, got %v", result.SelfScore)
	}
	if result.OthersScore != 4 { // (5+3)/2 = 4
		t.Errorf("expected others score 4, got %v", result.OthersScore)
	}
	if result.PerceptionGap != 0 { // 4 - 4
		t.Errorf("expected perception gap 0, got %v", result.PerceptionGap)
	}
}

// TestCalculateTarget_UnsubmittedRatersExcluded memverifikasi rater yang belum
// submit tidak ikut dalam perhitungan.
func TestCalculateTarget_UnsubmittedRatersExcluded(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	comp, _ := svc.CreateCompetency(ctx, CreateCompetencyRequest{Name: "Komunikasi"})
	ind, _ := svc.CreateIndicator(ctx, CreateIndicatorRequest{CompetencyID: comp.ID, Statement: "Berkomunikasi efektif"})
	indUID := mustParseUUID(t, ind.ID)
	reqLevel := 3
	tpl, _ := svc.CreateAssessmentTemplate(ctx, CreateAssessmentTemplateRequest{
		Name: "Tpl", Code: "TPL-002",
		Competencies: []TemplateCompetencyRequest{
			{CompetencyID: comp.ID, RequiredLevel: &reqLevel, Weight: 1},
		},
		RaterTypes: []TemplateRaterTypeRequest{
			{RaterType: "self", Weight: 0.5, MinRater: 1},
			{RaterType: "superior", Weight: 0.5, MinRater: 1},
		},
	})
	svc.SetTemplateIndicators(ctx, tpl.ID, []TemplateIndicatorRequest{{IndicatorID: ind.ID}})

	event, _ := svc.CreateCompetencyEvent(ctx, CreateCompetencyEventRequest{
		Type: "manual", PeriodType: "annual", PeriodYear: 2026, TemplateID: &tpl.ID,
	})
	subjectID := uuid.New()
	target, _ := svc.CreateCompetencyEventTarget(ctx, CreateCompetencyEventTargetRequest{
		CompetencyEventID: event.ID, OrganizationID: createTestOrgID(), EmployeeID: uuidStrPtr(subjectID),
	})
	targetUID := mustParseUUID(t, target.ID)

	// Self submitted (5), superior assigned tapi belum submit → hanya self dihitung.
	now := time.Now()
	selfRat := &CompetencyAssessmentRater{
		CompetencyEventTargetID: targetUID,
		RaterEmployeeID:         uuid.New(),
		RaterType:               "self",
		Weight:                  1,
		Status:                  "submitted",
		SubmittedAt:             &now,
	}
	repo := svc.repo
	repo.CreateRater(ctx, selfRat)
	repo.SaveResponse(ctx, &CompetencyAssessmentResponse{RaterID: selfRat.ID, IndicatorID: indUID, RatingValue: 5})

	supRat := &CompetencyAssessmentRater{
		CompetencyEventTargetID: targetUID,
		RaterEmployeeID:         uuid.New(),
		RaterType:               "superior",
		Weight:                  1,
		Status:                  "assigned",
	}
	repo.CreateRater(ctx, supRat)

	result, err := svc.CalculateTarget(ctx, target.ID)
	if err != nil {
		t.Fatalf("CalculateTarget failed: %v", err)
	}
	if len(result.Competencies) != 1 {
		t.Fatalf("expected 1 competency, got %d", len(result.Competencies))
	}
	c := result.Competencies[0]
	// Hanya self yang dihitung → score = 5 (weighted atas satu tipe = 5).
	if c.Score != 5 {
		t.Errorf("expected score 5 (unsubmitted superior excluded), got %v", c.Score)
	}
	if c.RaterScores.Superior != 0 {
		t.Errorf("expected superior score 0 (not submitted), got %v", c.RaterScores.Superior)
	}
}

// TestCalculateTarget_NoTemplateError memverifikasi target tanpa template
// menolak perhitungan.
func TestCalculateTarget_NoTemplateError(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	event, _ := svc.CreateCompetencyEvent(ctx, CreateCompetencyEventRequest{
		Type: "manual", PeriodType: "annual", PeriodYear: 2026,
	})
	target, _ := svc.CreateCompetencyEventTarget(ctx, CreateCompetencyEventTargetRequest{
		CompetencyEventID: event.ID, OrganizationID: createTestOrgID(),
	})

	_, err := svc.CalculateTarget(ctx, target.ID)
	if err == nil {
		t.Fatal("expected error when event has no template")
	}
}

// TestCalculateTarget_InvalidTargetID memverifikasi penolakan target ID invalid.
func TestCalculateTarget_InvalidTargetID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	_, err := svc.CalculateTarget(context.Background(), "not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid target id")
	}
}

// TestFinalizeTarget_Snapshot memverifikasi FinalizeTarget menulis snapshot ke
// competency_scores + competency_score_details (per employee per event).
func TestFinalizeTarget_Snapshot(t *testing.T) {
	svc, targetID, result, cleanup := setup360Scenario(t)
	defer cleanup()
	ctx := context.Background()

	finalized, err := svc.FinalizeTarget(ctx, targetID)
	if err != nil {
		t.Fatalf("FinalizeTarget failed: %v", err)
	}
	if finalized.OverallScore != result.OverallScore {
		t.Errorf("expected finalized overall score %v, got %v", result.OverallScore, finalized.OverallScore)
	}

	// Snapshot tersimpan di competency_scores untuk (event, employee).
	eventUID, _ := uuid.Parse(finalized.EventID)
	empUID, _ := uuid.Parse(finalized.EmployeeID)
	score, err := svc.repo.FindScoreByEventAndEmployee(ctx, eventUID, empUID)
	if err != nil {
		t.Fatalf("expected persisted score, got %v", err)
	}
	if score.TotalGradePercentage != finalized.OverallScore {
		t.Errorf("expected stored total_grade %v, got %v", finalized.OverallScore, score.TotalGradePercentage)
	}
	if len(score.Details) != 1 {
		t.Errorf("expected 1 score detail, got %d", len(score.Details))
	}
}
