package payroll

import (
	"context"
	"math"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// setupCalcEnv menyiapkan environment kalkulasi: 1 grading, 1 position,
// 1 employee + employment + payroll profile. Periode dibuat sendiri per test
// (supaya period_code unik per test).
func setupCalcEnv(t *testing.T, repo *Repository, ctx context.Context) (*GradingRead, *PositionRead, *EmployeeRead) {
	t.Helper()
	grading := createTestGrading(ctx, repo, "G-1", "Staff")
	position := createTestPosition(ctx, repo, "Staff HR", &grading.ID)
	employee := createTestEmployee(ctx, repo, "EMP001", "Asep")
	createTestEmployment(ctx, repo, &employee.ID, &position.ID, "2020-01-01")
	createTestPayrollProfile(ctx, repo, employee.ID, "2020-01-01")
	return grading, position, employee
}

// createTestComponent membuat salary component dengan tipe & nilai.
func createTestComponent(ctx context.Context, repo *Repository, code, name, compType, calcType string, displayOrder int) *SalaryComponent {
	sc := &SalaryComponent{
		Code:                   code,
		Name:                   name,
		ComponentType:          compType,
		CalculationType:        calcType,
		IsTaxable:              true,
		IsRecurring:            true,
		IsProratable:           true,
		PrintOnSalaryStructure: true,
		DisplayOrder:           displayOrder,
		Status:                 "ACTIVE",
	}
	if err := repo.CreateSalaryComponent(ctx, sc); err != nil {
		panic(err)
	}
	return sc
}

func TestCalculatePayrollRun_BasicStructure(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, employee := setupCalcEnv(t, repo, ctx)

	// Struktur grade: BASIC 10jt, TRANSPORT 500rb.
	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	transport := createTestComponent(ctx, repo, "TRANSPORT", "Transport", "EARNING", "FIXED", 20)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 10000000)
	createTestGradeComponent(ctx, repo, grading.ID, transport.ID, 500000)

	// Override employee: TRANSPORT 750rb.
	createTestEmployeeComponent(ctx, repo, employee.ID, transport.ID, 750000)

	// Run + calculate.
	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-2026-01",
	})

	resp, err := svc.CalculatePayrollRun(ctx, run.ID)
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	if resp.TotalEmployees != 1 {
		t.Errorf("expected 1 employee, got %d", resp.TotalEmployees)
	}
	if resp.TotalEarning != 10750000 {
		t.Errorf("expected total earning 10750000, got %v", resp.TotalEarning)
	}
	if resp.TotalNet != 10750000 {
		t.Errorf("expected total net 10750000, got %v", resp.TotalNet)
	}

	// Snapshot item: BASIC dari grade, TRANSPORT dari override.
	items, err := svc.ListPayrollRunItems(ctx, run.ID)
	if err != nil {
		t.Fatalf("ListPayrollRunItems: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	amountByCode := map[string]float64{}
	for _, it := range items {
		amountByCode[it.ComponentCode] = it.Amount
	}
	if amountByCode["BASIC"] != 10000000 {
		t.Errorf("BASIC expected 10000000, got %v", amountByCode["BASIC"])
	}
	if amountByCode["TRANSPORT"] != 750000 {
		t.Errorf("TRANSPORT expected 750000 (override), got %v", amountByCode["TRANSPORT"])
	}

	emps, err := svc.ListPayrollRunEmployees(ctx, run.ID)
	if err != nil {
		t.Fatalf("ListPayrollRunEmployees: %v", err)
	}
	if len(emps) != 1 {
		t.Fatalf("expected 1 run employee, got %d", len(emps))
	}
	if emps[0].NetAmount != 10750000 {
		t.Errorf("expected net 10750000, got %v", emps[0].NetAmount)
	}
}

func TestCalculatePayrollRun_FormulaComponent(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, employee := setupCalcEnv(t, repo, ctx)

	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 10000000)

	// Allowance FORMULA: BASIC * 10%.
	formula := "BASIC * 10%"
	allowance := createTestComponent(ctx, repo, "ALLOWANCE", "Allowance", "EARNING", "FORMULA", 20)
	allowance.Formula = &formula
	if err := repo.UpdateSalaryComponent(ctx, allowance); err != nil {
		t.Fatalf("update allowance: %v", err)
	}
	createTestGradeComponent(ctx, repo, grading.ID, allowance.ID, 0)

	// Deduction FORMULA: BASIC * 2% (JHT).
	jhtFormula := "BASIC * 2%"
	jht := createTestComponent(ctx, repo, "JHT_EMP", "JHT Employee", "DEDUCTION", "FORMULA", 30)
	jht.Formula = &jhtFormula
	if err := repo.UpdateSalaryComponent(ctx, jht); err != nil {
		t.Fatalf("update jht: %v", err)
	}
	createTestGradeComponent(ctx, repo, grading.ID, jht.ID, 0)

	_ = employee
	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-FORMULA",
	})

	resp, err := svc.CalculatePayrollRun(ctx, run.ID)
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	if resp.TotalEarning != 11000000 {
		t.Errorf("expected earning 11000000 (10jt + 1jt), got %v", resp.TotalEarning)
	}
	if resp.TotalDeduction != 200000 {
		t.Errorf("expected deduction 200000 (2 pct of 10jt), got %v", resp.TotalDeduction)
	}
	if resp.TotalNet != 10800000 {
		t.Errorf("expected net 10800000, got %v", resp.TotalNet)
	}

	// Pastikan formula tersimpan di snapshot item.
	items, err := svc.ListPayrollRunItems(ctx, run.ID)
	if err != nil {
		t.Fatalf("ListPayrollRunItems: %v", err)
	}
	for _, it := range items {
		if it.ComponentCode == "ALLOWANCE" && (it.Formula == nil || *it.Formula != formula) {
			t.Errorf("ALLOWANCE formula not persisted in snapshot")
		}
		if it.ComponentCode == "ALLOWANCE" && (it.FormulaResult == nil || *it.FormulaResult != 1000000) {
			t.Errorf("ALLOWANCE formula_result expected 1000000, got %v", it.FormulaResult)
		}
	}
}

func TestCalculatePayrollRun_ReferenceComponent(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, employee := setupCalcEnv(t, repo, ctx)

	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 5000000)

	// BONUS REFERENCE → BASIC.
	bonus := createTestComponent(ctx, repo, "BONUS", "Bonus", "EARNING", "REFERENCE", 20)
	bonus.ReferenceComponentID = &basic.ID
	if err := repo.UpdateSalaryComponent(ctx, bonus); err != nil {
		t.Fatalf("update bonus: %v", err)
	}
	createTestGradeComponent(ctx, repo, grading.ID, bonus.ID, 0)

	_ = employee
	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-REF",
	})

	resp, err := svc.CalculatePayrollRun(ctx, run.ID)
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	if resp.TotalEarning != 10000000 {
		t.Errorf("expected earning 10000000 (BASIC 5jt + BONUS 5jt), got %v", resp.TotalEarning)
	}
}

func TestCalculatePayrollRun_Adjustment(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, employee := setupCalcEnv(t, repo, ctx)

	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 10000000)

	// Bonus sekali-jalan: komponen THR_BONUS tidak ada di struktur.
	bonusComp := createTestComponent(ctx, repo, "THR_BONUS", "THR Bonus", "EARNING", "FIXED", 20)
	createTestAdjustment(ctx, repo, employee.ID, bonusComp.ID, 2026, 1, 2500000)

	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-ADJ",
	})

	resp, err := svc.CalculatePayrollRun(ctx, run.ID)
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	if resp.TotalEarning != 12500000 {
		t.Errorf("expected earning 12500000 (10jt + 2.5jt), got %v", resp.TotalEarning)
	}
}

func TestCalculatePayrollRun_ProrationJoinMidPeriod(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading := createTestGrading(ctx, repo, "G-1", "Staff")
	position := createTestPosition(ctx, repo, "Staff HR", &grading.ID)
	employee := createTestEmployee(ctx, repo, "EMP002", "Budi")
	// Join 16 Jan 2026 — periode 1-31 Jan → eligible 16 hari dari 31.
	createTestEmployment(ctx, repo, &employee.ID, &position.ID, "2026-01-16")
	createTestPayrollProfile(ctx, repo, employee.ID, "2026-01-16")

	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 3100000)

	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-PRORATE",
	})

	resp, err := svc.CalculatePayrollRun(ctx, run.ID)
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	// 3.100.000 * 16/31 = 1.600.000
	expected := 1600000.0
	if math.Abs(resp.TotalEarning-expected) > 0.01 {
		t.Errorf("expected prorated earning %v, got %v", expected, resp.TotalEarning)
	}
}

func TestCalculatePayrollRun_Recalculation(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, _ := setupCalcEnv(t, repo, ctx)

	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 10000000)

	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-RECALC",
	})

	if _, err := svc.CalculatePayrollRun(ctx, run.ID); err != nil {
		t.Fatalf("first calculate: %v", err)
	}

	// Ubah grade component BASIC → 12jt lalu hitung ulang.
	db, _ := repo.getDB(ctx)
	var gc SalaryGradeComponent
	if err := db.First(&gc, "salary_component_id = ?", basic.ID).Error; err != nil {
		t.Fatalf("find grade component: %v", err)
	}
	gc.Amount = 12000000
	if err := repo.UpdateSalaryGradeComponent(ctx, &gc); err != nil {
		t.Fatalf("update grade component: %v", err)
	}

	resp2, err := svc.CalculatePayrollRun(ctx, run.ID)
	if err != nil {
		t.Fatalf("recalculate: %v", err)
	}
	if resp2.TotalEarning != 12000000 {
		t.Errorf("expected updated earning 12000000, got %v", resp2.TotalEarning)
	}

	// Snapshot lama harus terganti (tidak dobel).
	items2, err := svc.ListPayrollRunItems(ctx, run.ID)
	if err != nil {
		t.Fatalf("ListPayrollRunItems: %v", err)
	}
	if len(items2) != 1 {
		t.Errorf("expected 1 item after recalculation, got %d", len(items2))
	}
}

func TestCalculatePayrollRun_RejectsLocked(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	setupCalcEnv(t, repo, ctx)
	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-LOCKED",
	})
	// Paksa status LOCKED.
	statusReq := UpdatePayrollRunStatusRequest{Status: "LOCKED"}
	updated, err := svc.UpdatePayrollRunStatus(ctx, run.ID, statusReq)
	if err != nil {
		t.Fatalf("update status to LOCKED: %v", err)
	}
	if updated.Status != "LOCKED" {
		t.Fatalf("expected LOCKED, got %s", updated.Status)
	}

	_, err = svc.CalculatePayrollRun(ctx, run.ID)
	if err == nil {
		t.Fatal("expected error calculating a LOCKED run")
	}
}

// TestCalculatePayrollRun_StatusTransitionCalculates: transisi DRAFT→CALCULATED
// lewat UpdatePayrollRunStatus harus benar-benar mengisi snapshot (bukan cuma
// ganti status) — gap paling kritis yang diperbaiki sub-plan ini.
func TestCalculatePayrollRun_StatusTransitionCalculates(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, _ := setupCalcEnv(t, repo, ctx)
	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 10000000)

	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-TRANSITION",
	})

	updated, err := svc.UpdatePayrollRunStatus(ctx, run.ID, UpdatePayrollRunStatusRequest{
		Status: "CALCULATED",
	})
	if err != nil {
		t.Fatalf("UpdatePayrollRunStatus CALCULATED: %v", err)
	}
	if updated.TotalEarning != 10000000 {
		t.Errorf("expected run calculated with earning 10000000, got %v (status: %s)", updated.TotalEarning, updated.Status)
	}

	items, err := svc.ListPayrollRunItems(ctx, run.ID)
	if err != nil {
		t.Fatalf("ListPayrollRunItems: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected snapshot items after CALCULATED transition, got %d", len(items))
	}
	if updated.Status == "DRAFT" {
		t.Errorf("expected run to leave DRAFT, still DRAFT")
	}
}

var _ = uuid.New
