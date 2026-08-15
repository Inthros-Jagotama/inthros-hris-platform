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
	createTestPtkp(ctx, repo, "TK0", "Tidak Kawin (TK/0)", 54000000, "A")
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
func int64Ptr(v int64) *int64 { return &v }

// setupTerRatesA memasang tarif TER kategori A (beberapa bracket di sekitar
// bruto 10jt — 9.650.000–10.050.000 → 2%) + bracket terbuka untuk sisanya.
func setupTerRatesA(t *testing.T, repo *Repository, ctx context.Context) {
	t.Helper()
	createTestTerRate(ctx, repo, "A", int64Ptr(0), int64Ptr(5400000), 0.00)
	createTestTerRate(ctx, repo, "A", int64Ptr(9650000), int64Ptr(10050000), 2.00)
	createTestTerRate(ctx, repo, "A", int64Ptr(10050000), nil, 2.25)
}

// TestPph21TerMonthly: metode TER Jan–Nov — PPh21 = bruto bulanan × tarif TER
// (kategori A untuk TK/0). Bruto 10jt → bracket 9.65–10.05jt = 2% → 200rb.
// Tanpa biaya jabatan/BPJS/pensiun/annualisasi — hanya tarif × bruto.
func TestPph21TerMonthly(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, _, _ := setupPph21Env(t, repo, ctx, true)
	basic := createTestSalaryComponent(ctx, repo)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 10000000)
	pph21Comp := createTestComponent(ctx, repo, "PPH21", "PPh21", "DEDUCTION", "FORMULA", 50)
	createTestPph21SettingCustomMethod(ctx, repo, pph21Comp, "TER")
	setupTerRatesA(t, repo, ctx)

	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo) // 2026-01 (Januari)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-PPH21-TER-JAN",
	})

	resp, err := svc.CalculatePayrollRun(ctx, run.ID)
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	// Hanya potongan PPh21 TER (200rb) — tanpa BPJS di setup ini.
	assertClose(t, "total deduction", resp.TotalDeduction, 200000)

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
		t.Fatal("PPh21 item not found")
	}
	assertClose(t, "pph21 ter amount", pph21Item.Amount, 200000)
	assertClose(t, "pph21 ter base", pph21Item.BaseAmount, 10000000)

	logs := findPph21LogsForTest(t, repo, ctx, run.ID)
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}
	if logs[0].CalculationMethod != "TER" {
		t.Errorf("expected calculation_method TER, got %s", logs[0].CalculationMethod)
	}
	if logs[0].GrossMonthly != 10000000 {
		t.Errorf("expected gross 10jt, got %v", logs[0].GrossMonthly)
	}
}

// TestPph21TerDecember: Desember — pajak setahun (metode normal) dikurangi
// potongan TER Jan–Nov. Run Januari (TER, 200rb) lalu Desember → pajak
// Desember = pajak tahunan − 200rb.
func TestPph21TerDecember(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, _, _ := setupPph21Env(t, repo, ctx, true)
	basic := createTestSalaryComponent(ctx, repo)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 10000000)
	pph21Comp := createTestComponent(ctx, repo, "PPH21", "PPh21", "DEDUCTION", "FORMULA", 50)
	createTestPph21SettingCustomMethod(ctx, repo, pph21Comp, "TER")
	setupTerRatesA(t, repo, ctx)
	createTestPtkp(ctx, repo, "TK0", "Tidak Kawin (TK/0)", 54000000, "A")
	createTestPph21TaxBracket(ctx, repo, 1, 0, floatPtr(60000000), 5.0)

	svc := NewService(repo, zap.NewNop())

	// Run Januari (TER monthly) — potongan 200rb.
	jan := createTestPayrollPeriodCustom(ctx, repo, 2026, 1)
	runJan, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: jan.ID.String(),
		RunCode:         "RUN-TER-JAN",
	})
	if _, err := svc.CalculatePayrollRun(ctx, runJan.ID); err != nil {
		t.Fatalf("calculate january: %v", err)
	}

	// Run Desember: pajak setahun 3jt (gross 10jt − occ 500rb = 9.5jt × 12 = 114jt
	// − PTKP 54jt = PKP 60jt × 5%) − YTD potongan TER Jan–Nov (200rb) = 2.8jt.
	dec := createTestPayrollPeriodCustom(ctx, repo, 2026, 12)
	runDec, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: dec.ID.String(),
		RunCode:         "RUN-TER-DEC",
	})
	resp, err := svc.CalculatePayrollRun(ctx, runDec.ID)
	if err != nil {
		t.Fatalf("calculate december: %v", err)
	}
	assertClose(t, "total deduction dec", resp.TotalDeduction, 2800000)

	items, err := svc.ListPayrollRunItems(ctx, runDec.ID)
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
		t.Fatal("PPh21 item not found in december")
	}
	assertClose(t, "pph21 december amount", pph21Item.Amount, 2800000)

	// Log Desember mencatat TER + gross.
	logs := findPph21LogsForTest(t, repo, ctx, runDec.ID)
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}
	if logs[0].CalculationMethod != "TER" {
		t.Errorf("expected calculation_method TER in december, got %s", logs[0].CalculationMethod)
	}
}
