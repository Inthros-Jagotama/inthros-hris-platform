package payroll

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

// TestReports_SummaryDetailBpjsTaxBank: run dengan struktur + BPJS + PPh21 +
// bank profile → semua laporan dari snapshot menghasilkan angka yang tepat.
func TestReports_SummaryDetailBpjsTaxBank(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, employee, bpjsProfile, _ := setupPph21Env(t, repo, ctx, true)
	_, _ = setupPph21Rates(t, repo, ctx, grading.ID, 10000000)
	createTestPph21TaxBracket(ctx, repo, 1, 0, floatPtr(60000000), 5.0)

	// Tambah rate employer kesehatan 4% + nomor BPJS di profil.
	db, _ := repo.getDB(ctx)
	var bs BpjsSetting
	if err := db.First(&bs, "setting_code = ?", "BPJS-DEFAULT").Error; err != nil {
		t.Fatalf("find bpjs setting: %v", err)
	}
	erComp := createTestComponent(ctx, repo, "BPJS_KES_ER", "BPJS Kesehatan - Employer", "EMPLOYER_CONTRIBUTION", "PERCENTAGE", 55)
	createTestBpjsRateComponentLinked(ctx, repo, &bs, erComp, "HEALTH", "EMPLOYER", "BPJS-KES-ER", 4.0)

	bpjsNo := "1234567890"
	bpjsProfile.BpjsHealthNo = &bpjsNo
	if err := repo.UpdateEmployeeBpjsProfile(ctx, bpjsProfile); err != nil {
		t.Fatalf("update bpjs profile: %v", err)
	}
	createTestBankProfile(ctx, repo, employee.ID, bpjsProfile.EmployeePayrollProfileID, "9988776655", "John")

	svc := NewService(repo, zap.NewNop())
	run := calcRunForPayslip(t, svc, ctx, repo, "RUN-REPORT")

	// Summary & dashboard.
	summary, err := svc.GetPayrollSummaryReport(ctx, run.ID)
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if summary.TotalEmployees != 1 || summary.GrossSalary != 10000000 {
		t.Errorf("summary employees/gross: %d / %v", summary.TotalEmployees, summary.GrossSalary)
	}
	if summary.EmployeeDeduction != 540000 { // health 100rb + jht 200rb + pph21 240rb
		t.Errorf("summary deduction: %v", summary.EmployeeDeduction)
	}
	if summary.EmployerContribution != 400000 {
		t.Errorf("summary employer contribution: %v", summary.EmployerContribution)
	}
	if summary.NetSalary != 9460000 {
		t.Errorf("summary net: %v", summary.NetSalary)
	}
	if summary.TotalCompanyCost != 10400000 {
		t.Errorf("summary company cost: %v", summary.TotalCompanyCost)
	}
	dash, err := svc.GetPayrollDashboard(ctx, run.ID)
	if err != nil {
		t.Fatalf("dashboard: %v", err)
	}
	if dash.NetSalary != summary.NetSalary || dash.Status != "REVIEWED" {
		t.Errorf("dashboard mismatch: %+v", dash)
	}

	// Detail per komponen.
	detail, err := svc.GetPayrollDetailReport(ctx, run.ID)
	if err != nil {
		t.Fatalf("detail: %v", err)
	}
	if len(detail) != 5 { // BASIC + 3 BPJS + PPH21
		t.Fatalf("expected 5 detail rows, got %d", len(detail))
	}
	amountByCode := map[string]float64{}
	for _, r := range detail {
		amountByCode[r.ComponentCode] = r.Amount
	}
	if amountByCode["BASIC"] != 10000000 || amountByCode["PPH21"] != 240000 {
		t.Errorf("detail amounts wrong: %+v", amountByCode)
	}

	// BPJS report.
	bpjsRows, err := svc.GetBpjsReport(ctx, run.ID)
	if err != nil {
		t.Fatalf("bpjs report: %v", err)
	}
	if len(bpjsRows) != 1 {
		t.Fatalf("expected 1 bpjs row, got %d", len(bpjsRows))
	}
	b := bpjsRows[0]
	if b.WageBasis != 10000000 || b.EmployeeContribution != 300000 || b.EmployerContribution != 400000 || b.TotalContribution != 700000 {
		t.Errorf("bpjs row wrong: %+v", b)
	}
	if b.BpjsNumber != bpjsNo {
		t.Errorf("bpjs number: expected %s, got %s", bpjsNo, b.BpjsNumber)
	}

	// Tax report.
	taxRows, err := svc.GetTaxReport(ctx, run.ID)
	if err != nil {
		t.Fatalf("tax report: %v", err)
	}
	if len(taxRows) != 1 || taxRows[0].TaxableIncome != 10000000 || taxRows[0].Pph21 != 240000 {
		t.Errorf("tax row wrong: %+v", taxRows)
	}

	// Bank transfer report (butuh payment batch dulu).
	if _, err := svc.CreatePaymentBatch(ctx, run.ID); err != nil {
		t.Fatalf("create batch: %v", err)
	}
	bankRows, err := svc.GetBankTransferReport(ctx, run.ID)
	if err != nil {
		t.Fatalf("bank report: %v", err)
	}
	if len(bankRows) != 1 {
		t.Fatalf("expected 1 bank row, got %d", len(bankRows))
	}
	if bankRows[0].AccountNumber != "9988776655" || bankRows[0].Amount != 9460000 {
		t.Errorf("bank row wrong: %+v", bankRows[0])
	}
}

// TestGoldenDatasetRegression — golden dataset: perubahan engine apa pun harus
// tetap menghasilkan angka berikut untuk input tetap ini.
//
//	BASIC 10jt (taxable + bpjs base)
//	BPJS: kesehatan 1% employee / 4% employer, JHT 2% employee
//	PPh21: TK/0 (PTKP 54jt), biaya jabatan 5% (cap 500rb), bracket 5%,
//	      deduct BPJS: health ✗, JHT ✓
//
//	Gross               = 10.000.000
//	BPJS employee       = 300.000   (100rb kesehatan + 200rb JHT)
//	PPh21               = 240.000
//	Total deduction     = 540.000
//	Net                 = 9.460.000
//	BPJS employer       = 400.000
//	Company cost        = 10.400.000
func TestGoldenDatasetRegression(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, employee, _, _ := setupPph21Env(t, repo, ctx, true)
	_, _ = setupPph21Rates(t, repo, ctx, grading.ID, 10000000)
	createTestPph21TaxBracket(ctx, repo, 1, 0, floatPtr(60000000), 5.0)
	db, _ := repo.getDB(ctx)
	var bs BpjsSetting
	if err := db.First(&bs, "setting_code = ?", "BPJS-DEFAULT").Error; err != nil {
		t.Fatalf("find bpjs setting: %v", err)
	}
	erComp := createTestComponent(ctx, repo, "BPJS_KES_ER", "BPJS Kesehatan - Employer", "EMPLOYER_CONTRIBUTION", "PERCENTAGE", 55)
	createTestBpjsRateComponentLinked(ctx, repo, &bs, erComp, "HEALTH", "EMPLOYER", "BPJS-KES-ER", 4.0)

	svc := NewService(repo, zap.NewNop())
	run := calcRunForPayslip(t, svc, ctx, repo, "RUN-GOLDEN")

	emps, err := svc.ListPayrollRunEmployees(ctx, run.ID)
	if err != nil {
		t.Fatalf("list employees: %v", err)
	}
	if len(emps) != 1 {
		t.Fatalf("expected 1 employee, got %d", len(emps))
	}
	emp := emps[0]
	assertClose(t, "gross (earning)", emp.TotalEarning, 10000000)
	assertClose(t, "bpjs employee deduction", emp.TotalDeduction, 540000)
	assertClose(t, "employer contribution", emp.TotalEmployerContribution, 400000)
	assertClose(t, "net salary", emp.NetAmount, 9460000)
	assertClose(t, "company cost", emp.TotalCompanyCost, 10400000)

	// Pastikan jumlahnya stabil saat payslip di-generate (golden, bukan cuma run).
	payslips, err := svc.GeneratePayslips(ctx, run.ID)
	if err != nil {
		t.Fatalf("generate payslips: %v", err)
	}
	if len(payslips) != 1 {
		t.Fatalf("expected 1 payslip, got %d", len(payslips))
	}
	assertClose(t, "payslip net", payslips[0].NetAmount, 9460000)

	_ = employee
}
