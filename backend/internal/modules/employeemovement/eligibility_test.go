package employeemovement

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/inthros/hris-platform/internal/modules/employee"
)

// fakePerformanceProvider implements PerformanceProvider for tests.
type fakePerformanceProvider struct {
	score float64
	found bool
}

func (f fakePerformanceProvider) LatestFinalScore(_ context.Context, _ uuid.UUID) (float64, bool, error) {
	return f.score, f.found, nil
}

// fakeCompetencyProvider implements CompetencyProvider for tests.
type fakeCompetencyProvider struct {
	score float64
	found bool
}

func (f fakeCompetencyProvider) LatestScore(_ context.Context, _ uuid.UUID) (float64, bool, error) {
	return f.score, f.found, nil
}

// fakeOKRProvider implements OKRProvider for tests.
type fakeOKRProvider struct {
	score float64
	found bool
}

func (f fakeOKRProvider) LatestScore(_ context.Context, _ uuid.UUID) (float64, bool, error) {
	return f.score, f.found, nil
}

// TestGetMovementEligibility_AllMet verifies all rules met when tenure >= 24
// months and performance/competency scores >= 80.
func TestGetMovementEligibility_AllMet(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	svc.SetPerformanceProvider(fakePerformanceProvider{score: 95, found: true})
	svc.SetCompetencyProvider(fakeCompetencyProvider{score: 88, found: true})
	svc.SetOKRProvider(fakeOKRProvider{score: 91, found: true})

	employeeID := uuid.New()
	orgID, posID, statusID := seedCareerReferenceTables(t, repo, employeeID)

	// Employment dengan effective date 3 tahun lalu → tenure > 24 bulan.
	seedEmployment(t, repo, &employee.Employment{
		EmployeeID:           &employeeID,
		OrganizationID:       &orgID,
		PositionID:           &posID,
		EmploymentStatusID:   &statusID,
		DecisionLetterNumber: "SK-TL-001",
		DecisionLetterDate:   "2023-01-01",
		EffectiveDate:        "2023-01-01",
	})

	resp, err := svc.GetMovementEligibility(ctx(), employeeID.String())
	if err != nil {
		t.Fatalf("GetMovementEligibility failed: %v", err)
	}
	if !resp.Data.Eligible {
		t.Errorf("expected eligible=true for met-all rules, got false")
	}
	if resp.Data.TenureMonths < 24 {
		t.Errorf("expected tenure >= 24 months, got %d", resp.Data.TenureMonths)
	}
	if resp.Data.OKRScore == nil || *resp.Data.OKRScore != 91 {
		t.Errorf("expected okr_score 91, got %v", resp.Data.OKRScore)
	}
	foundEmp := false
	foundPerf := false
	foundComp := false
	foundOKR := false
	for _, r := range resp.Data.Rules {
		if r.Code == "minimum_service" && r.Met {
			foundEmp = true
		}
		if r.Code == "performance" && r.Met {
			foundPerf = true
		}
		if r.Code == "competency" && r.Met {
			foundComp = true
		}
		if r.Code == "okr" && r.Met {
			foundOKR = true
		}
	}
	if !foundEmp {
		t.Error("expected minimum_service rule met")
	}
	if !foundPerf {
		t.Error("expected performance rule met")
	}
	if !foundComp {
		t.Error("expected competency rule met")
	}
	if !foundOKR {
		t.Error("expected okr rule met")
	}
}

// TestGetMovementEligibility_NotMet verifies rules not met when tenure
// < 24 months and scores < 80.
func TestGetMovementEligibility_NotMet(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	svc.SetPerformanceProvider(fakePerformanceProvider{score: 50, found: true})
	svc.SetCompetencyProvider(fakeCompetencyProvider{score: 60, found: true})

	employeeID := uuid.New()
	orgID, posID, statusID := seedCareerReferenceTables(t, repo, employeeID)

	// Employment dengan effective date 3 bulan lalu → tenure < 6 bulan.
	recentDate := time.Now().AddDate(0, -3, 0).Format("2006-01-02")
	seedEmployment(t, repo, &employee.Employment{
		EmployeeID:           &employeeID,
		OrganizationID:       &orgID,
		PositionID:           &posID,
		EmploymentStatusID:   &statusID,
		DecisionLetterNumber: "SK-TL-001",
		DecisionLetterDate:   recentDate,
		EffectiveDate:        recentDate,
	})

	resp, err := svc.GetMovementEligibility(ctx(), employeeID.String())
	if err != nil {
		t.Fatalf("GetMovementEligibility failed: %v", err)
	}
	if resp.Data.Eligible {
		t.Error("expected eligible=false when rules not met")
	}
	if resp.Data.TenureMonths >= 6 {
		t.Errorf("expected tenure < 6 months, got %d", resp.Data.TenureMonths)
	}
}

// TestGetMovementEligibility_NoProviders verifies nil providers are handled
// gracefully (all data rules show met=false with detail, but service doesn't
// panic and returns OK).
func TestGetMovementEligibility_NoProviders(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	// Providers not set (nil).
	employeeID := uuid.New()
	orgID, posID, statusID := seedCareerReferenceTables(t, repo, employeeID)
	seedEmployment(t, repo, &employee.Employment{
		EmployeeID:           &employeeID,
		OrganizationID:       &orgID,
		PositionID:           &posID,
		EmploymentStatusID:   &statusID,
		DecisionLetterNumber: "SK-TL-001",
		DecisionLetterDate:   "2023-01-01",
		EffectiveDate:        "2023-01-01",
	})

	resp, err := svc.GetMovementEligibility(ctx(), employeeID.String())
	if err != nil {
		t.Fatalf("GetMovementEligibility failed: %v", err)
	}
	if resp.Data.PerformanceScore != nil {
		t.Errorf("expected performance_score nil when provider not set, got %v", *resp.Data.PerformanceScore)
	}
	if resp.Data.CompetencyScore != nil {
		t.Errorf("expected competency_score nil when provider not set, got %v", *resp.Data.CompetencyScore)
	}
}

// TestGetPromotionEligibility_WithCareerPath verifies promotion eligibility
// when employee's current position is in a career path with a next step.
func TestGetPromotionEligibility_WithCareerPath(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	svc.SetPerformanceProvider(fakePerformanceProvider{score: 90, found: true})
	svc.SetCompetencyProvider(fakeCompetencyProvider{score: 85, found: true})
	svc.SetOKRProvider(fakeOKRProvider{score: 91, found: true})

	employeeID := uuid.New()
	orgID, posID, statusID := seedCareerReferenceTables(t, repo, employeeID)

	// Seed employments: join 3 years ago → tenure > 24 bulan.
	seedEmployment(t, repo, &employee.Employment{
		EmployeeID:           &employeeID,
		OrganizationID:       &orgID,
		PositionID:           &posID,
		EmploymentStatusID:   &statusID,
		DecisionLetterNumber: "SK-TL-001",
		DecisionLetterDate:   "2023-01-01",
		EffectiveDate:        "2023-01-01",
	})

	// Seed career path dengan staff → senior staff → supervisor step.
	// Posisi karyawan = posID (Staff), next step = senior staff.
	staffPosID := posID
	seniorPosID := uuid.New()
	supPosID := uuid.New()

	// Seed positions untuk career path steps.
	db, err := repo.getDB(ctx())
	if err != nil {
		t.Fatalf("failed to get test db: %v", err)
	}
	for _, p := range []struct {
		id    uuid.UUID
		title string
	}{
		{seniorPosID, "Senior Staff"},
		{supPosID, "Supervisor"},
	} {
		if err := db.Exec("INSERT INTO positions (id, title) VALUES (?, ?)", p.id.String(), p.title).Error; err != nil {
			t.Fatalf("failed to seed position %s: %v", p.title, err)
		}
	}
	// Seed career path langsung via gorm (CreateCareerPathTx dipindah ke modul
	// Career Intelligence — modul ini hanya membaca untuk eligibility).
	path := &CareerPath{ID: uuid.New(), Name: "Test Career Path — Staff ke Supervisor", IsActive: true}
	if err := db.Create(path).Error; err != nil {
		t.Fatalf("failed to seed career path: %v", err)
	}
	steps := []CareerPathStep{
		{ID: uuid.New(), CareerPathID: path.ID, PositionID: staffPosID, Sequence: 1, MinimumServiceMonths: intPtr(12)},
		{ID: uuid.New(), CareerPathID: path.ID, PositionID: seniorPosID, Sequence: 2, MinimumServiceMonths: intPtr(24)},
		{ID: uuid.New(), CareerPathID: path.ID, PositionID: supPosID, Sequence: 3, MinimumServiceMonths: intPtr(36)},
	}
	for _, st := range steps {
		if err := db.Create(&st).Error; err != nil {
			t.Fatalf("failed to seed career path step: %v", err)
		}
	}

	resp, err := svc.GetPromotionEligibility(ctx(), employeeID.String())
	if err != nil {
		t.Fatalf("GetPromotionEligibility failed: %v", err)
	}
	if resp.Data.NextPositionID == nil || *resp.Data.NextPositionID != seniorPosID.String() {
		t.Errorf("expected next position %s, got %v", seniorPosID.String(), resp.Data.NextPositionID)
	}
	if resp.Data.NextPositionName == nil || *resp.Data.NextPositionName != "Senior Staff" {
		t.Errorf("expected next position name 'Senior Staff', got %v", resp.Data.NextPositionName)
	}
	// Minimum service should be the next step's min_service_months (24).
	if resp.Data.MinimumServiceMonths != 24 {
		t.Errorf("expected minimum_service_months 24 (next step), got %d", resp.Data.MinimumServiceMonths)
	}
	// Employee with 3+ years tenure and performance 90 should be eligible.
	if !resp.Data.Eligible {
		t.Errorf("expected eligible=true for met-all promotion rules")
	}
}

// TestGetMovementEligibility_OKRNotMet verifies the OKR rule marks eligible=false
// when the OKR score is below threshold (plan §12.11 — OKR sebagai input
// eligibility bersama KPI & competency).
func TestGetMovementEligibility_OKRNotMet(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	svc.SetPerformanceProvider(fakePerformanceProvider{score: 90, found: true})
	svc.SetCompetencyProvider(fakeCompetencyProvider{score: 90, found: true})
	svc.SetOKRProvider(fakeOKRProvider{score: 40, found: true})

	employeeID := uuid.New()
	orgID, posID, statusID := seedCareerReferenceTables(t, repo, employeeID)
	seedEmployment(t, repo, &employee.Employment{
		EmployeeID:           &employeeID,
		OrganizationID:       &orgID,
		PositionID:           &posID,
		EmploymentStatusID:   &statusID,
		DecisionLetterNumber: "SK-TL-001",
		DecisionLetterDate:   "2023-01-01",
		EffectiveDate:        "2023-01-01",
	})

	resp, err := svc.GetMovementEligibility(ctx(), employeeID.String())
	if err != nil {
		t.Fatalf("GetMovementEligibility failed: %v", err)
	}
	if resp.Data.Eligible {
		t.Error("expected eligible=false when okr score below threshold")
	}
	okrMet := false
	for _, r := range resp.Data.Rules {
		if r.Code == "okr" && r.Met {
			okrMet = true
		}
	}
	if okrMet {
		t.Error("expected okr rule not met")
	}
}

// TestGetMovementEligibility_MissingOKRData_DoesNotBlock verifies the
// pragmatic nil-data policy: OKR rule tanpa data (tenant belum menjalankan
// OKR) tidak memblokir eligible selama rule lain yang punya data semuanya met.
func TestGetMovementEligibility_MissingOKRData_DoesNotBlock(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	svc.SetPerformanceProvider(fakePerformanceProvider{score: 90, found: true})
	svc.SetCompetencyProvider(fakeCompetencyProvider{score: 88, found: true})
	// OKR provider TIDAK di-set (nil) — tenant belum pakai OKR.

	employeeID := uuid.New()
	orgID, posID, statusID := seedCareerReferenceTables(t, repo, employeeID)
	seedEmployment(t, repo, &employee.Employment{
		EmployeeID:           &employeeID,
		OrganizationID:       &orgID,
		PositionID:           &posID,
		EmploymentStatusID:   &statusID,
		DecisionLetterNumber: "SK-TL-001",
		DecisionLetterDate:   "2023-01-01",
		EffectiveDate:        "2023-01-01",
	})

	resp, err := svc.GetMovementEligibility(ctx(), employeeID.String())
	if err != nil {
		t.Fatalf("GetMovementEligibility failed: %v", err)
	}
	if !resp.Data.Eligible {
		t.Error("expected eligible=true — okr rule tanpa data tidak memblokir")
	}
	// Rule OKR tetap dilaporkan sebagai available=false.
	found := false
	for _, r := range resp.Data.Rules {
		if r.Code == "okr" && !r.Available {
			found = true
		}
	}
	if !found {
		t.Error("expected okr rule reported with available=false")
	}
}

// TestComputeTenureMonths verifies tenure calculation from employments.
func TestComputeTenureMonths(t *testing.T) {
	tests := []struct {
		name     string
		date     string
		expected int // approximate; we set the start date far enough back
	}{
		{"employments 3 years ago, 3+ years tenure → expected >= 24", "2023-01-01", 24},
		{"employments 6 months ago → expected < 24", "", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows := []careerEmploymentRow{}
			if tt.date != "" {
				rows = append(rows, careerEmploymentRow{EffectiveDate: tt.date})
			}
			months := computeTenureMonths(rows)
			if len(rows) > 0 && months < tt.expected {
				t.Errorf("expected tenure >= %d months, got %d", tt.expected, months)
			}
			if len(rows) == 0 && months != 0 {
				t.Errorf("expected 0 months for empty employments, got %d", months)
			}
		})
	}
}

// TestEvaluatePromotionRules_PerStepThresholdOverride verifies
// career_path_requirements (module-career-intelligence-plan.md §9 #6): a
// per-step threshold (e.g. from career_path_steps.min_performance_score)
// must be used instead of the global 80.0 default when supplied.
func TestEvaluatePromotionRules_PerStepThresholdOverride(t *testing.T) {
	svc := &Service{}
	perf, comp, okr := 70.0, 70.0, 70.0

	// With global default thresholds (80/80/80), a score of 70 fails all three.
	rulesDefault := svc.evaluatePromotionRules(24, 24, &perf, &comp, &okr,
		eligibilityMinPerformanceScore, eligibilityMinCompetencyScore, eligibilityMinOKRScore)
	for _, r := range rulesDefault {
		if r.Code != "minimum_service" && r.Met {
			t.Errorf("expected rule %s to fail against default threshold 80, but it passed", r.Code)
		}
	}

	// With a lower per-step override (60/60/60), the same score of 70 passes.
	rulesOverride := svc.evaluatePromotionRules(24, 24, &perf, &comp, &okr, 60.0, 60.0, 60.0)
	for _, r := range rulesOverride {
		if r.Code != "minimum_service" && !r.Met {
			t.Errorf("expected rule %s to pass against overridden threshold 60, but it failed", r.Code)
		}
	}
}
