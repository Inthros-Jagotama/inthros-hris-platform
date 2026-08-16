package competency

import (
	"context"
	"fmt"
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
		Name: "Tpl",
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

// setupManagerOrgTree menyiapkan tabel minimal (employees, employee_accounts,
// organizations, employments) + relasi organisasi: manager bekerja di org root,
// bawahan bekerja di org anak (subtree). Mengembalikan managerID & userID login.
func setupManagerOrgTree(t *testing.T, repo *Repository, ctx context.Context, subordinates ...uuid.UUID) (managerID, userID uuid.UUID) {
	t.Helper()
	db, err := repo.getDB(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("CREATE TABLE IF NOT EXISTS employees (id CHAR(36) PRIMARY KEY, name VARCHAR(255))").Error; err != nil {
		t.Fatalf("create employees: %v", err)
	}
	if err := db.Exec("CREATE TABLE IF NOT EXISTS employee_accounts (id CHAR(36) PRIMARY KEY, user_id CHAR(36), employee_id CHAR(36))").Error; err != nil {
		t.Fatalf("create employee_accounts: %v", err)
	}
	if err := db.Exec("CREATE TABLE IF NOT EXISTS organizations (id CHAR(36) PRIMARY KEY, parent_id CHAR(36), deleted_at TIMESTAMP)").Error; err != nil {
		t.Fatalf("create organizations: %v", err)
	}
	if err := db.Exec("CREATE TABLE IF NOT EXISTS employments (id CHAR(36) PRIMARY KEY, employee_id CHAR(36), organization_id CHAR(36), effective_date DATE, effective_end_date DATE)").Error; err != nil {
		t.Fatalf("create employments: %v", err)
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	managerID = uuid.New()
	userID = uuid.New()
	managerOrgID := uuid.New()
	childOrgID := uuid.New()

	// Org tree: managerOrg (root) → childOrg.
	insertOrg := func(id uuid.UUID, parent interface{}) {
		if err := db.Exec("INSERT INTO organizations (id, parent_id, deleted_at) VALUES (?, ?, NULL)", id.String(), parent).Error; err != nil {
			t.Fatalf("seed organization: %v", err)
		}
	}
	insertOrg(managerOrgID, nil)
	insertOrg(childOrgID, managerOrgID.String())

	insertEmp := func(id uuid.UUID, name string) {
		if err := db.Exec("INSERT INTO employees (id, name) VALUES (?, ?)", id.String(), name).Error; err != nil {
			t.Fatalf("seed employee %s: %v", name, err)
		}
	}
	insertEmp(managerID, "Manager")
	for i, sub := range subordinates {
		insertEmp(sub, fmt.Sprintf("Bawahan %d", i+1))
	}

	// Employment saat ini: manager di org root, bawahan di org anak (subtree).
	insertEmployment := func(empID, orgID uuid.UUID) {
		if err := db.Exec("INSERT INTO employments (id, employee_id, organization_id, effective_date, effective_end_date) VALUES (?, ?, ?, ?, NULL)",
			uuid.New().String(), empID.String(), orgID.String(), now,
		).Error; err != nil {
			t.Fatalf("seed employment: %v", err)
		}
	}
	insertEmployment(managerID, managerOrgID)
	for _, sub := range subordinates {
		insertEmployment(sub, childOrgID)
	}

	if err := db.Exec("INSERT INTO employee_accounts (id, user_id, employee_id) VALUES (?, ?, ?)",
		uuid.New().String(), userID.String(), managerID.String(),
	).Error; err != nil {
		t.Fatalf("seed employee_accounts: %v", err)
	}
	return managerID, userID
}

// TestService_ManagerAssessments_SubordinatesFromOrg memverifikasi
// ManagerAssessments mengambil daftar bawahan dari struktur organisasi
// (seluruh employee di subtree org tempat manager bekerja — bukan kolom
// supervisor manual) dan mengisi status rater superior manager.
func TestService_ManagerAssessments_SubordinatesFromOrg(t *testing.T) {
	svc, targetID, _, cleanup := setup360Scenario(t)
	defer cleanup()
	ctx := context.Background()
	repo := svc.repo

	target, err := repo.FindCompetencyEventTargetByID(ctx, mustParseUUID(t, targetID))
	if err != nil {
		t.Fatalf("get target: %v", err)
	}
	subjectID := uuid.New()
	if target.EmployeeID != nil {
		subjectID = *target.EmployeeID
	}
	otherSubID := uuid.New()

	// Manager (user login) dengan 2 bawahan di subtree org; bawahan #1 adalah
	// subject target scenario, bawahan #2 belum punya target sama sekali.
	_, userID := setupManagerOrgTree(t, repo, ctx, subjectID, otherSubID)

	assessments, err := svc.ManagerAssessments(ctxWithUserID(userID), "")
	if err != nil {
		t.Fatalf("ManagerAssessments failed: %v", err)
	}
	// Hanya bawahan #1 yang punya target (Bawahan 2 di-skip — belum subject).
	if len(assessments) != 1 {
		t.Fatalf("expected 1 subordinate with target, got %d: %v", len(assessments), assessments)
	}
	if assessments[0].EmployeeID != subjectID.String() {
		t.Errorf("expected employee %s, got %s", subjectID, assessments[0].EmployeeID)
	}
	if assessments[0].EmployeeName != "Bawahan 1" {
		t.Errorf("expected name 'Bawahan 1', got %q", assessments[0].EmployeeName)
	}
	if assessments[0].TargetID == "" || assessments[0].CompetencyEventID == "" {
		t.Errorf("expected target & event populated, got target=%q event=%q", assessments[0].TargetID, assessments[0].CompetencyEventID)
	}
	// Rater superior manager belum di-assign pada target → RaterID kosong.
	if assessments[0].RaterID != "" {
		t.Errorf("expected empty rater_id (manager not assigned), got %q", assessments[0].RaterID)
	}
}

// TestService_ManagerAssessments_RaterStatusFilled memverifikasi status rater
// superior manager terisi bila manager sudah di-assign pada target bawahan.
func TestService_ManagerAssessments_RaterStatusFilled(t *testing.T) {
	svc, targetID, _, cleanup := setup360Scenario(t)
	defer cleanup()
	ctx := context.Background()
	repo := svc.repo

	target, err := repo.FindCompetencyEventTargetByID(ctx, mustParseUUID(t, targetID))
	if err != nil {
		t.Fatalf("get target: %v", err)
	}
	subjectID := uuid.New()
	if target.EmployeeID != nil {
		subjectID = *target.EmployeeID
	}
	managerID, userID := setupManagerOrgTree(t, repo, ctx, subjectID)

	// Assign manager sebagai superior rater pada target bawahan.
	assignedAt := time.Now()
	rat := &CompetencyAssessmentRater{
		CompetencyEventTargetID: target.ID,
		RaterEmployeeID:         managerID,
		RaterType:               string(RaterTypeSuperior),
		Weight:                  1,
		Status:                  string(RaterStatusAssigned),
		AssignedAt:              &assignedAt,
	}
	if err := repo.CreateRater(ctx, rat); err != nil {
		t.Fatalf("assign manager as superior: %v", err)
	}

	assessments, err := svc.ManagerAssessments(ctxWithUserID(userID), "")
	if err != nil {
		t.Fatalf("ManagerAssessments failed: %v", err)
	}
	if len(assessments) != 1 {
		t.Fatalf("expected 1 assessment, got %d: %v", len(assessments), assessments)
	}
	if assessments[0].RaterID != rat.ID.String() {
		t.Errorf("expected rater_id %s, got %q", rat.ID, assessments[0].RaterID)
	}
	if assessments[0].RaterStatus != "assigned" {
		t.Errorf("expected rater_status 'assigned', got %q", assessments[0].RaterStatus)
	}
}

// ctxWithUserID menaruh user_id (string) di context — pola persis
// authctx.GetUserID yang membaca ctx.Value("user_id") sebagai string.
func ctxWithUserID(userID uuid.UUID) context.Context {
	return context.WithValue(context.Background(), "user_id", userID.String())
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

// TestService_SuggestedRaters memverifikasi saran rater dari struktur
// organisasi: superior (employee di parent org subject) dan subordinate
// (employee di subtree org subject), dengan rater yang sudah di-assign
// dikecualikan.
func TestService_SuggestedRaters(t *testing.T) {
	svc, targetID, _, cleanup := setup360Scenario(t)
	defer cleanup()
	ctx := context.Background()
	repo := svc.repo

	target, err := repo.FindCompetencyEventTargetByID(ctx, mustParseUUID(t, targetID))
	if err != nil {
		t.Fatalf("get target: %v", err)
	}
	subjectID := uuid.New()
	if target.EmployeeID != nil {
		subjectID = *target.EmployeeID
	}

	// Org tree: root (parent) → child → grandchild.
	// Atasan (manager) bekerja di root, subject di child, bawahan di grandchild.
	db, err := repo.getDB(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, tbl := range []string{
		"CREATE TABLE IF NOT EXISTS employees (id CHAR(36) PRIMARY KEY, name VARCHAR(255))",
		"CREATE TABLE IF NOT EXISTS organizations (id CHAR(36) PRIMARY KEY, parent_id CHAR(36), deleted_at TIMESTAMP)",
		"CREATE TABLE IF NOT EXISTS employments (id CHAR(36) PRIMARY KEY, employee_id CHAR(36), organization_id CHAR(36), effective_date DATE, effective_end_date DATE)",
	} {
		if err := db.Exec(tbl).Error; err != nil {
			t.Fatalf("create table: %v", err)
		}
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	rootID := uuid.New()
	childID := uuid.New()
	grandchildID := uuid.New()
	managerID := uuid.New()
	subordinateID := uuid.New()

	insertOrg := func(id uuid.UUID, parent interface{}) {
		if err := db.Exec("INSERT INTO organizations (id, parent_id, deleted_at) VALUES (?, ?, NULL)", id.String(), parent).Error; err != nil {
			t.Fatalf("seed org: %v", err)
		}
	}
	insertOrg(rootID, nil)
	insertOrg(childID, rootID.String())
	insertOrg(grandchildID, childID.String())

	insertEmp := func(id uuid.UUID, name string) {
		if err := db.Exec("INSERT INTO employees (id, name) VALUES (?, ?)", id.String(), name).Error; err != nil {
			t.Fatalf("seed employee: %v", err)
		}
	}
	insertEmp(managerID, "Manager")
	insertEmp(subjectID, "Subject")
	insertEmp(subordinateID, "Bawahan")

	insertEmpRel := func(empID, orgID uuid.UUID) {
		if err := db.Exec("INSERT INTO employments (id, employee_id, organization_id, effective_date, effective_end_date) VALUES (?, ?, ?, ?, NULL)",
			uuid.New().String(), empID.String(), orgID.String(), now,
		).Error; err != nil {
			t.Fatalf("seed employment: %v", err)
		}
	}
	insertEmpRel(managerID, rootID)
	insertEmpRel(subjectID, childID)
	insertEmpRel(subordinateID, grandchildID)

	// Assign manager sebagai superior — harus dikecualikan dari saran.
	if err := repo.CreateRater(ctx, &CompetencyAssessmentRater{
		CompetencyEventTargetID: target.ID,
		RaterEmployeeID:         managerID,
		RaterType:               string(RaterTypeSuperior),
		Status:                  string(RaterStatusAssigned),
	}); err != nil {
		t.Fatalf("assign manager: %v", err)
	}

	sug, err := svc.SuggestedRaters(ctx, targetID)
	if err != nil {
		t.Fatalf("SuggestedRaters failed: %v", err)
	}
	// Manager sudah di-assign → superior kosong.
	if len(sug.Superior) != 0 {
		t.Errorf("expected no superior suggestion (already assigned), got %+v", sug.Superior)
	}
	// Bawahan di subtree org subject → tersedia.
	if len(sug.Subordinates) != 1 {
		t.Fatalf("expected 1 subordinate suggestion, got %d: %+v", len(sug.Subordinates), sug.Subordinates)
	}
	if sug.Subordinates[0].ID != subordinateID.String() {
		t.Errorf("expected subordinate %s, got %s", subordinateID, sug.Subordinates[0].ID)
	}
	if sug.Subordinates[0].Name != "Bawahan" {
		t.Errorf("expected name 'Bawahan', got %q", sug.Subordinates[0].Name)
	}
}
