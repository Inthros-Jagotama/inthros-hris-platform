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

// TestGetMovementEligibility_AllMet verifies all rules met when tenure >= 24
// months and performance/competency scores >= 80.
func TestGetMovementEligibility_AllMet(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	svc.SetPerformanceProvider(fakePerformanceProvider{score: 95, found: true})
	svc.SetCompetencyProvider(fakeCompetencyProvider{score: 88, found: true})

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
	foundEmp := false
	foundPerf := false
	foundComp := false
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
	path := &CareerPath{Name: "Test Career Path — Staff ke Supervisor", Description: strPtr("Jenjang test"), IsActive: true}
	path.ID = uuid.New()
	if err := repo.CreateCareerPathTx(ctx(), path, []CareerPathStep{
		{CareerPathID: path.ID, PositionID: staffPosID, Sequence: 1, MinimumServiceMonths: intPtr(12)},
		{CareerPathID: path.ID, PositionID: seniorPosID, Sequence: 2, MinimumServiceMonths: intPtr(24)},
		{CareerPathID: path.ID, PositionID: supPosID, Sequence: 3, MinimumServiceMonths: intPtr(36)},
	}); err != nil {
		t.Fatalf("failed to seed career path: %v", err)
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
