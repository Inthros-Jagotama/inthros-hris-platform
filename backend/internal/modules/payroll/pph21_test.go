package payroll

import (
	"context"
	"math"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// setupPph21Env menyiapkan env kalkulasi lengkap: grading/position/employee/
// employment/payroll profile + profil BPJS + profil pajak.
func setupPph21Env(t *testing.T, repo *Repository, ctx context.Context, hasNpwp bool) (*GradingRead, *EmployeeRead, *EmployeeBpjsProfile, *EmployeeTaxProfile) {
	t.Helper()
	grading, employee, bpjsProfile := setupBpjsEnv(t, repo, ctx)
	taxProfile := createTestEmployeeTaxProfile(ctx, repo, employee.ID, bpjsProfile.EmployeePayrollProfileID, "TK/0", hasNpwp)
	return grading, employee, bpjsProfile, taxProfile
}

// setupPph21Rates memasang komponen BASIC + BPJS (health 1%, jht 2%) + komponen
// potongan PPh21 + setting + PTKP + bracket. Salary dikontrol via basicAmount.
func setupPph21Rates(t *testing.T, repo *Repository, ctx context.Context, gradingID uuid.UUID, basicAmount float64) (*SalaryComponent, *Pph21Setting) {
	t.Helper()
	basic := createTestSalaryComponent(ctx, repo)
	createTestGradeComponent(ctx, repo, gradingID, basic.ID, basicAmount)

	setting := createTestBpjsSetting(ctx, repo)
	healthComp := createTestComponent(ctx, repo, "BPJS_KES_EMP", "BPJS Kesehatan - Employee", "DEDUCTION", "PERCENTAGE", 40)
	createTestBpjsRateComponentLinked(ctx, repo, setting, healthComp, "HEALTH", "EMPLOYEE", "BPJS-KES-EMP", 1.0)
	jhtComp := createTestComponent(ctx, repo, "BPJS_JHT_EMP", "BPJS JHT - Employee", "DEDUCTION", "PERCENTAGE", 41)
	createTestBpjsRateComponentLinked(ctx, repo, setting, jhtComp, "JHT", "EMPLOYEE", "BPJS-JHT-EMP", 2.0)

	pph21Comp := createTestComponent(ctx, repo, "PPH21", "PPh21", "DEDUCTION", "FORMULA", 50)
	pph21Setting := createTestPph21SettingCustom(ctx, repo, pph21Comp)
	createTestPph21PtkpRate(ctx, repo, "TK/0", 54000000)
	return basic, pph21Setting
}

func assertClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 0.01 {
		t.Errorf("%s: expected %v, got %v", name, want, got)
	}
}

// TestPph21Basic: gross 10jt, biaya jabatan 500rb, BPJS deductible hanya JHT
// (2% = 200rb; health 1% TIDAK dikurangkan), PTKP TK/0 54jt, bracket 5% →
// PPh21 bulanan 240rb. Potongan total = health 100rb + jht 200rb + pajak 240rb.
func TestPph21Basic(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, _, _ := setupPph21Env(t, repo, ctx, true)
	_, _ = setupPph21Rates(t, repo, ctx, grading.ID, 10000000)
	createTestPph21TaxBracket(ctx, repo, 1, 0, floatPtr(60000000), 5.0)

	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-PPH21-BASIC",
	})

	resp, err := svc.CalculatePayrollRun(ctx, run.ID)
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	assertClose(t, "total earning", resp.TotalEarning, 10000000)
	assertClose(t, "total deduction", resp.TotalDeduction, 540000)
	assertClose(t, "total net", resp.TotalNet, 9460000)

	// Item PPh21 di snapshot.
	items, err := svc.ListPayrollRunItems(ctx, run.ID)
	if err != nil {
		t.Fatalf("ListPayrollRunItems: %v", err)
	}
	var pph21Item *PayrollRunItemResponse
	for i := range items {
		if items[i].ComponentCode == "PPH21" {
			pph21Item = &items[i]
		}
	}
	if pph21Item == nil {
		t.Fatal("PPh21 item not found in snapshot")
	}
	assertClose(t, "pph21 amount", pph21Item.Amount, 240000)
	if pph21Item.SourceGroup != "STATUTORY" {
		t.Errorf("expected source_group STATUTORY, got %s", pph21Item.SourceGroup)
	}
	if pph21Item.ItemCategory != "EMPLOYEE_DEDUCTION" {
		t.Errorf("expected EMPLOYEE_DEDUCTION, got %s", pph21Item.ItemCategory)
	}
	assertClose(t, "pph21 base (gross)", pph21Item.BaseAmount, 10000000)
}

// TestPph21ProgressiveAndNoNpwp: gross 100jt, bracket 5/15/25%, tanpa NPWP
// (multiplier 120%) → pajak tahunan 248jt × 1.2 = 297.6jt → bulanan 24.8jt.
func TestPph21ProgressiveAndNoNpwp(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, _, _ := setupPph21Env(t, repo, ctx, false)
	_, pph21Setting := setupPph21Rates(t, repo, ctx, grading.ID, 100000000)
	createTestPph21TaxBracket(ctx, repo, 1, 0, floatPtr(60000000), 5.0)
	createTestPph21TaxBracket(ctx, repo, 2, 60000000, floatPtr(250000000), 15.0)
	createTestPph21TaxBracket(ctx, repo, 3, 250000000, nil, 25.0)
	mult := 120.0
	pph21Setting.NonNpwpMultiplierPercent = mult
	if err := repo.UpdatePph21Setting(ctx, pph21Setting); err != nil {
		t.Fatalf("update pph21 setting: %v", err)
	}

	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-PPH21-PROG",
	})

	resp, err := svc.CalculatePayrollRun(ctx, run.ID)
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	// Potongan = health 1jt + JHT 2jt + PPh21 24.8jt.
	assertClose(t, "total deduction", resp.TotalDeduction, 27800000)
	assertClose(t, "total net", resp.TotalNet, 72200000)

}

// TestPph21NoTaxProfile: employee tanpa profil pajak → tidak ada item PPh21.
func TestPph21NoTaxProfile(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	// Env TANPA profil pajak (setupBpjsEnv tidak membuat tax profile).
	grading, _, _ := setupBpjsEnv(t, repo, ctx)
	_, _ = setupPph21Rates(t, repo, ctx, grading.ID, 10000000)
	createTestPph21TaxBracket(ctx, repo, 1, 0, floatPtr(60000000), 5.0)

	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-PPH21-NOPROFILE",
	})

	resp, err := svc.CalculatePayrollRun(ctx, run.ID)
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	// Hanya BPJS (health 100rb + jht 200rb) sebagai potongan.
	assertClose(t, "total deduction", resp.TotalDeduction, 300000)

	items, err := svc.ListPayrollRunItems(ctx, run.ID)
	if err != nil {
		t.Fatalf("ListPayrollRunItems: %v", err)
	}
	for _, it := range items {
		if it.ComponentCode == "PPH21" {
			t.Error("expected no PPh21 item without tax profile")
		}
	}
}

// TestPph21LogPersisted: log tertulis, ter-link ke run employee & item, dan
// aman di-recalculate (log lama dihapus, tidak dobel).
func TestPph21LogPersisted(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, _, _ := setupPph21Env(t, repo, ctx, true)
	_, _ = setupPph21Rates(t, repo, ctx, grading.ID, 10000000)
	createTestPph21TaxBracket(ctx, repo, 1, 0, floatPtr(60000000), 5.0)

	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-PPH21-LOG",
	})

	if _, err := svc.CalculatePayrollRun(ctx, run.ID); err != nil {
		t.Fatalf("first calculate: %v", err)
	}
	logs := findPph21LogsForTest(t, repo, ctx, run.ID)
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}
	log := logs[0]
	assertClose(t, "log gross", log.GrossMonthly, 10000000)
	assertClose(t, "log occ", log.OccupationalExpenseMonthly, 500000)
	assertClose(t, "log bpjs deductible", log.BpjsTaxDeductibleMonthly, 200000)
	assertClose(t, "log net monthly", log.NetMonthly, 9300000)
	assertClose(t, "log ptkp", log.PtkpAnnual, 54000000)
	assertClose(t, "log pkp", log.PkpAnnual, 57600000)
	assertClose(t, "log monthly tax", log.Pph21Monthly, 240000)
	if log.PayrollRunEmployeeID == uuid.Nil {
		t.Error("expected log linked to run employee")
	}
	if log.PayrollRunItemID == nil {
		t.Error("expected log linked to payroll run item")
	}
	if log.FormulaJSON == nil || *log.FormulaJSON == "" {
		t.Error("expected formula_json populated")
	}

	// Recalculate — harus tetap 1 log (log lama dihapus) tanpa error unique.
	if _, err := svc.CalculatePayrollRun(ctx, run.ID); err != nil {
		t.Fatalf("recalculate: %v", err)
	}
	logs2 := findPph21LogsForTest(t, repo, ctx, run.ID)
	if len(logs2) != 1 {
		t.Errorf("expected 1 log after recalc, got %d", len(logs2))
	}
}

// findPph21LogsForTest mengambil seluruh log PPh21 sebuah run (via repo).
func findPph21LogsForTest(t *testing.T, repo *Repository, ctx context.Context, runID string) []Pph21CalculationLog {
	t.Helper()
	uid := uuid.MustParse(runID)
	db, err := repo.getDB(ctx)
	if err != nil {
		t.Fatalf("getDB: %v", err)
	}
	var logs []Pph21CalculationLog
	if err := db.Where("payroll_run_id = ?", uid).Find(&logs).Error; err != nil {
		t.Fatalf("find logs: %v", err)
	}
	return logs
}

func floatPtr(v float64) *float64 { return &v }
