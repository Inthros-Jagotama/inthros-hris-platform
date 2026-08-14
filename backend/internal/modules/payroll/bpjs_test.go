package payroll

import (
	"context"
	"math"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// setupBpjsEnv menyiapkan env kalkulasi + profil BPJS employee.
func setupBpjsEnv(t *testing.T, repo *Repository, ctx context.Context) (*GradingRead, *EmployeeRead, *EmployeeBpjsProfile) {
	t.Helper()
	grading := createTestGrading(ctx, repo, "G-1", "Staff")
	position := createTestPosition(ctx, repo, "Staff HR", &grading.ID)
	employee := createTestEmployee(ctx, repo, "EMP-BPJS", "Cicih")
	createTestEmployment(ctx, repo, &employee.ID, &position.ID, "2020-01-01")
	profile := createTestPayrollProfile(ctx, repo, employee.ID, "2020-01-01")
	bpjsProfile := createTestEmployeeBpjsProfile(ctx, repo, employee.ID, profile.ID)
	return grading, employee, bpjsProfile
}

// setupBpjsRates memasang rate HEALTH employee 1% + employer 4% yang terhubung
// ke salary component, plus komponen BASIC 10jt di struktur grade.
func setupBpjsRates(t *testing.T, repo *Repository, ctx context.Context, gradingID uuid.UUID) (*SalaryComponent, *BpjsRateComponent, *BpjsRateComponent) {
	t.Helper()
	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	basic.IsBpjsBase = true
	if err := repo.UpdateSalaryComponent(ctx, basic); err != nil {
		t.Fatalf("update basic: %v", err)
	}
	createTestGradeComponent(ctx, repo, gradingID, basic.ID, 10000000)

	setting := createTestBpjsSetting(ctx, repo)
	empComp := createTestComponent(ctx, repo, "BPJS_KES_EMP", "BPJS Kesehatan - Employee", "DEDUCTION", "PERCENTAGE", 40)
	empRate := createTestBpjsRateComponentLinked(ctx, repo, setting, empComp, "HEALTH", "EMPLOYEE", "BPJS-KES-EMP", 1.0)
	erComp := createTestComponent(ctx, repo, "BPJS_KES_EMP_ER", "BPJS Kesehatan - Employer", "EMPLOYER_CONTRIBUTION", "PERCENTAGE", 50)
	erRate := createTestBpjsRateComponentLinked(ctx, repo, setting, erComp, "HEALTH", "EMPLOYER", "BPJS-KES-EMP-ER", 4.0)
	return basic, empRate, erRate
}

// TestBpjsHealthContributions: HEALTH 1% employee + 4% employer atas dasar upah
// 10jt → potongan 100rb, kontribusi perusahaan 400rb.
func TestBpjsHealthContributions(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, _ := setupBpjsEnv(t, repo, ctx)
	setupBpjsRates(t, repo, ctx, grading.ID)

	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-BPJS-HEALTH",
	})

	resp, err := svc.CalculatePayrollRun(ctx, run.ID)
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	if resp.TotalEarning != 10000000 {
		t.Errorf("expected earning 10000000, got %v", resp.TotalEarning)
	}
	if resp.TotalDeduction != 100000 {
		t.Errorf("expected deduction 100000 (1%% of 10jt), got %v", resp.TotalDeduction)
	}
	if resp.TotalEmployerContribution != 400000 {
		t.Errorf("expected employer contribution 400000 (4%% of 10jt), got %v", resp.TotalEmployerContribution)
	}
	if resp.TotalNet != 9900000 {
		t.Errorf("expected net 9900000, got %v", resp.TotalNet)
	}
	if resp.TotalCompanyCost != 10400000 {
		t.Errorf("expected company cost 10400000, got %v", resp.TotalCompanyCost)
	}

	// Snapshot item: BASIC + 2 BPJS, semuanya source_group=STATUTORY untuk BPJS.
	items, err := svc.ListPayrollRunItems(ctx, run.ID)
	if err != nil {
		t.Fatalf("ListPayrollRunItems: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("expected 3 items (BASIC + 2 BPJS), got %d", len(items))
	}
	var empItem, erItem *PayrollRunItemResponse
	for i := range items {
		it := &items[i]
		switch it.ComponentCode {
		case "BPJS_KES_EMP":
			empItem = it
		case "BPJS_KES_EMP_ER":
			erItem = it
		}
	}
	if empItem == nil || erItem == nil {
		t.Fatal("BPJS items not found in snapshot")
	}
	if empItem.SourceGroup != "STATUTORY" || erItem.SourceGroup != "STATUTORY" {
		t.Errorf("expected source_group STATUTORY, got %s / %s", empItem.SourceGroup, erItem.SourceGroup)
	}
	if empItem.ItemCategory != "EMPLOYEE_DEDUCTION" {
		t.Errorf("expected EMPLOYEE_DEDUCTION, got %s", empItem.ItemCategory)
	}
	if erItem.ItemCategory != "EMPLOYER_CONTRIBUTION" {
		t.Errorf("expected EMPLOYER_CONTRIBUTION, got %s", erItem.ItemCategory)
	}
	if empItem.Amount != 100000 || erItem.Amount != 400000 {
		t.Errorf("expected amounts 100000/400000, got %v/%v", empItem.Amount, erItem.Amount)
	}
	if empItem.BaseAmount != 10000000 {
		t.Errorf("expected base 10000000, got %v", empItem.BaseAmount)
	}
	if empItem.Rate == nil || math.Abs(*empItem.Rate-0.01) > 1e-9 {
		t.Errorf("expected rate 0.01, got %v", empItem.Rate)
	}
}

// TestBpjsHealthCap: HealthMaxBaseAmount membatasi dasar upah (8jt → 1% = 80rb).
func TestBpjsHealthCap(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, _ := setupBpjsEnv(t, repo, ctx)
	_, _, _ = setupBpjsRates(t, repo, ctx, grading.ID)

	// Setting dengan cap kesehatan 8jt.
	db, _ := repo.getDB(ctx)
	var setting BpjsSetting
	if err := db.First(&setting, "setting_code = ?", "BPJS-DEFAULT").Error; err != nil {
		t.Fatalf("find setting: %v", err)
	}
	capVal := 8000000.0
	setting.HealthMaxBaseAmount = &capVal
	if err := repo.UpdateBpjsSetting(ctx, &setting); err != nil {
		t.Fatalf("update setting: %v", err)
	}

	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-BPJS-CAP",
	})

	resp, err := svc.CalculatePayrollRun(ctx, run.ID)
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	if resp.TotalDeduction != 80000 {
		t.Errorf("expected deduction 80000 (1%% of capped 8jt), got %v", resp.TotalDeduction)
	}
	if resp.TotalEmployerContribution != 320000 {
		t.Errorf("expected employer 320000 (4%% of capped 8jt), got %v", resp.TotalEmployerContribution)
	}
}

// TestBpjsJkkRiskClass: rate JKK khusus risk class HIGH dilewati employee LOW;
// rate JKK umum (tanpa risk class) tetap dihitung.
func TestBpjsJkkRiskClass(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, _ := setupBpjsEnv(t, repo, ctx)
	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	basic.IsBpjsBase = true
	if err := repo.UpdateSalaryComponent(ctx, basic); err != nil {
		t.Fatalf("update basic: %v", err)
	}
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 10000000)

	setting := createTestBpjsSetting(ctx, repo)
	highComp := createTestComponent(ctx, repo, "JKK_HIGH", "JKK Risiko Tinggi", "EMPLOYER_CONTRIBUTION", "PERCENTAGE", 30)
	highRate := createTestBpjsRateComponentLinked(ctx, repo, setting, highComp, "JKK", "EMPLOYER", "BPJS-JKK-HIGH", 0.54)
	highRisk := "HIGH"
	highRate.JkkRiskClass = &highRisk
	if err := repo.UpdateBpjsRateComponent(ctx, highRate); err != nil {
		t.Fatalf("update high risk rate: %v", err)
	}

	genericComp := createTestComponent(ctx, repo, "JKK", "JKK", "EMPLOYER_CONTRIBUTION", "PERCENTAGE", 31)
	genericRate := createTestBpjsRateComponentLinked(ctx, repo, setting, genericComp, "JKK", "EMPLOYER", "BPJS-JKK-GENERIC", 0.24)

	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-BPJS-JKK",
	})

	resp, err := svc.CalculatePayrollRun(ctx, run.ID)
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	// Hanya rate umum yang dipakai (0.24% dari 10jt = 24rb); rate HIGH dilewati.
	if resp.TotalEmployerContribution != 24000 {
		t.Errorf("expected employer contribution 24000, got %v (high-risk rate harus di-skip)", resp.TotalEmployerContribution)
	}

	items, err := svc.ListPayrollRunItems(ctx, run.ID)
	if err != nil {
		t.Fatalf("ListPayrollRunItems: %v", err)
	}
	foundHigh := false
	for _, it := range items {
		if it.ComponentCode == "JKK_HIGH" {
			foundHigh = true
		}
	}
	if foundHigh {
		t.Error("expected JKK_HIGH rate skipped for LOW risk employee")
	}
	_ = genericRate
}

// TestBpjsFixedAmount: rate dengan fixed_amount menghasilkan nominal tetap.
func TestBpjsFixedAmount(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, _ := setupBpjsEnv(t, repo, ctx)
	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	basic.IsBpjsBase = true
	if err := repo.UpdateSalaryComponent(ctx, basic); err != nil {
		t.Fatalf("update basic: %v", err)
	}
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 10000000)

	setting := createTestBpjsSetting(ctx, repo)
	comp := createTestComponent(ctx, repo, "JKM", "JKM", "EMPLOYER_CONTRIBUTION", "FIXED", 30)
	rate := createTestBpjsRateComponentLinked(ctx, repo, setting, comp, "JKM", "EMPLOYER", "BPJS-JKM", 0)
	fixed := 24000.0
	rate.FixedAmount = &fixed
	if err := repo.UpdateBpjsRateComponent(ctx, rate); err != nil {
		t.Fatalf("update fixed rate: %v", err)
	}

	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-BPJS-FIXED",
	})

	resp, err := svc.CalculatePayrollRun(ctx, run.ID)
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	if resp.TotalEmployerContribution != 24000 {
		t.Errorf("expected fixed employer contribution 24000, got %v", resp.TotalEmployerContribution)
	}

	items, err := svc.ListPayrollRunItems(ctx, run.ID)
	if err != nil {
		t.Fatalf("ListPayrollRunItems: %v", err)
	}
	for _, it := range items {
		if it.ComponentCode == "JKM" && it.CalculationType != "FIXED" {
			t.Errorf("expected calculation_type FIXED, got %s", it.CalculationType)
		}
	}
}

// TestBpjsNoProfile: employee tanpa profil BPJS → tidak ada item BPJS.
func TestBpjsNoProfile(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	// Env kalkulasi standar TANPA profil BPJS employee.
	grading := createTestGrading(ctx, repo, "G-1", "Staff")
	position := createTestPosition(ctx, repo, "Staff HR", &grading.ID)
	employee := createTestEmployee(ctx, repo, "EMP-NOBPJS", "Dede")
	createTestEmployment(ctx, repo, &employee.ID, &position.ID, "2020-01-01")
	createTestPayrollProfile(ctx, repo, employee.ID, "2020-01-01")
	setupBpjsRates(t, repo, ctx, grading.ID)

	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-BPJS-NOPROFILE",
	})

	resp, err := svc.CalculatePayrollRun(ctx, run.ID)
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	if resp.TotalEarning != 10000000 || resp.TotalDeduction != 0 {
		t.Errorf("expected earning 10jt tanpa deduction, got earning=%v deduction=%v", resp.TotalEarning, resp.TotalDeduction)
	}

	items, err := svc.ListPayrollRunItems(ctx, run.ID)
	if err != nil {
		t.Fatalf("ListPayrollRunItems: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected only BASIC item, got %d items", len(items))
	}
}
