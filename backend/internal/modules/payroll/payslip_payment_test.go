package payroll

import (
	"context"
	"strings"
	"testing"

	"go.uber.org/zap"
)

// calcRunForPayslip menjalankan kalkulasi + transisi status ke CALCULATED
// (supaya run final dan layak dibuat payslip/payment).
func calcRunForPayslip(t *testing.T, svc *Service, ctx context.Context, repo *Repository, runCode string) *PayrollRunResponse {
	t.Helper()
	period := createTestPayrollPeriod(ctx, repo)
	run, err := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         runCode,
	})
	if err != nil {
		t.Fatalf("CreatePayrollRun: %v", err)
	}
	updated, err := svc.UpdatePayrollRunStatus(ctx, run.ID, UpdatePayrollRunStatusRequest{Status: "CALCULATED"})
	if err != nil {
		t.Fatalf("UpdatePayrollRunStatus: %v", err)
	}
	return updated
}

func TestGeneratePayslips(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, _ := setupCalcEnv(t, repo, ctx)
	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 10000000)

	svc := NewService(repo, zap.NewNop())
	run := calcRunForPayslip(t, svc, ctx, repo, "RUN-PAYSLIP")

	payslips, err := svc.GeneratePayslips(ctx, run.ID)
	if err != nil {
		t.Fatalf("GeneratePayslips: %v", err)
	}
	if len(payslips) != 1 {
		t.Fatalf("expected 1 payslip, got %d", len(payslips))
	}
	ps := payslips[0]
	if ps.PayslipNumber != "SLP-202601-001" {
		t.Errorf("expected number SLP-202601-001, got %s", ps.PayslipNumber)
	}
	if ps.TotalEarning != 10000000 || ps.NetAmount != 10000000 {
		t.Errorf("expected earning/net 10jt, got %v/%v", ps.TotalEarning, ps.NetAmount)
	}
	if ps.Status != "DRAFT" {
		t.Errorf("expected status DRAFT, got %s", ps.Status)
	}

	// Detail via GetPayslipByID — items ikut.
	detail, err := svc.GetPayslipByID(ctx, ps.ID)
	if err != nil {
		t.Fatalf("GetPayslipByID: %v", err)
	}
	if len(detail.Items) != 1 || detail.Items[0].ComponentCode != "BASIC" {
		t.Errorf("expected 1 item BASIC, got %+v", detail.Items)
	}

	// List by run.
	listed, err := svc.ListPayslipsByRun(ctx, run.ID)
	if err != nil {
		t.Fatalf("ListPayslipsByRun: %v", err)
	}
	if len(listed) != 1 {
		t.Errorf("expected 1 listed payslip, got %d", len(listed))
	}
}

func TestGeneratePayslipsRejectsDraftRun(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	setupCalcEnv(t, repo, ctx)
	svc := NewService(repo, zap.NewNop())
	period := createTestPayrollPeriod(ctx, repo)
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         "RUN-DRAFT-PS",
	})
	if _, err := svc.CalculatePayrollRun(ctx, run.ID); err != nil {
		t.Fatalf("calculate: %v", err)
	}
	_, err := svc.GeneratePayslips(ctx, run.ID)
	if err == nil {
		t.Fatal("expected error generating payslip from DRAFT run")
	}
}

func TestPublishCancelPayslip(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, _ := setupCalcEnv(t, repo, ctx)
	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 5000000)

	svc := NewService(repo, zap.NewNop())
	run := calcRunForPayslip(t, svc, ctx, repo, "RUN-PS-PUB")
	payslips, _ := svc.GeneratePayslips(ctx, run.ID)

	published, err := svc.PublishPayslip(ctx, payslips[0].ID)
	if err != nil {
		t.Fatalf("PublishPayslip: %v", err)
	}
	if published.Status != "PUBLISHED" || published.PublishedAt == nil {
		t.Errorf("expected PUBLISHED with timestamp, got %s", published.Status)
	}

	cancelled, err := svc.CancelPayslip(ctx, payslips[0].ID)
	if err != nil {
		t.Fatalf("CancelPayslip: %v", err)
	}
	if cancelled.Status != "CANCELLED" {
		t.Errorf("expected CANCELLED, got %s", cancelled.Status)
	}

	// Cancel lagi → error.
	if _, err := svc.CancelPayslip(ctx, payslips[0].ID); err == nil {
		t.Error("expected error cancelling already-cancelled payslip")
	}
}

func TestPayslipHTML(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, _ := setupCalcEnv(t, repo, ctx)
	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 10000000)

	svc := NewService(repo, zap.NewNop())
	run := calcRunForPayslip(t, svc, ctx, repo, "RUN-PS-HTML")
	payslips, _ := svc.GeneratePayslips(ctx, run.ID)

	html, err := svc.GetPayslipHTML(ctx, payslips[0].ID)
	if err != nil {
		t.Fatalf("GetPayslipHTML: %v", err)
	}
	if !strings.Contains(html, "SLIP GAJI") || !strings.Contains(html, "Asep") {
		t.Error("HTML payslip missing expected content")
	}
	if !strings.Contains(html, "10000000") {
		t.Errorf("HTML payslip missing amount: %s", html)
	}
}

func TestCreatePaymentBatch(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading := createTestGrading(ctx, repo, "G-1", "Staff")
	position := createTestPosition(ctx, repo, "Staff HR", &grading.ID)
	employee := createTestEmployee(ctx, repo, "EMP-PAY", "Gugun")
	createTestEmployment(ctx, repo, &employee.ID, &position.ID, "2020-01-01")
	profile := createTestPayrollProfile(ctx, repo, employee.ID, "2020-01-01")
	createTestBankProfile(ctx, repo, employee.ID, profile.ID, "1234567890", "Gugun Gumilar")

	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 10000000)

	svc := NewService(repo, zap.NewNop())
	run := calcRunForPayslip(t, svc, ctx, repo, "RUN-PAY")

	batch, err := svc.CreatePaymentBatch(ctx, run.ID)
	if err != nil {
		t.Fatalf("CreatePaymentBatch: %v", err)
	}
	if batch.Total != 1 || batch.Skipped != 0 {
		t.Errorf("expected total 1 skipped 0, got %d/%d", batch.Total, batch.Skipped)
	}
	if batch.TotalAmount != 10000000 {
		t.Errorf("expected total amount 10jt, got %v", batch.TotalAmount)
	}

	payments, err := svc.ListPaymentsByRun(ctx, run.ID)
	if err != nil {
		t.Fatalf("ListPaymentsByRun: %v", err)
	}
	if len(payments) != 1 {
		t.Fatalf("expected 1 payment, got %d", len(payments))
	}
	p := payments[0]
	if p.Status != "PENDING" || p.Amount != 10000000 {
		t.Errorf("expected PENDING 10jt, got %s %v", p.Status, p.Amount)
	}
	// Snapshot rekening.
	if p.BankAccountNumber != "1234567890" || p.BankAccountHolderName != "Gugun Gumilar" {
		t.Errorf("unexpected bank snapshot: %s / %s", p.BankAccountNumber, p.BankAccountHolderName)
	}
}

func TestCreatePaymentBatchSkipsNoBank(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, _ := setupCalcEnv(t, repo, ctx) // tanpa bank profile
	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 10000000)

	svc := NewService(repo, zap.NewNop())
	run := calcRunForPayslip(t, svc, ctx, repo, "RUN-PAY-NOBANK")

	_, err := svc.CreatePaymentBatch(ctx, run.ID)
	if err == nil {
		t.Fatal("expected error when no employee has bank profile")
	}
}

func TestUpdatePaymentStatus(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading := createTestGrading(ctx, repo, "G-1", "Staff")
	position := createTestPosition(ctx, repo, "Staff HR", &grading.ID)
	employee := createTestEmployee(ctx, repo, "EMP-PAY2", "Herman")
	createTestEmployment(ctx, repo, &employee.ID, &position.ID, "2020-01-01")
	profile := createTestPayrollProfile(ctx, repo, employee.ID, "2020-01-01")
	createTestBankProfile(ctx, repo, employee.ID, profile.ID, "555000111", "Herman")

	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 10000000)

	svc := NewService(repo, zap.NewNop())
	run := calcRunForPayslip(t, svc, ctx, repo, "RUN-PAY-STATUS")
	if _, err := svc.CreatePaymentBatch(ctx, run.ID); err != nil {
		t.Fatalf("create batch: %v", err)
	}
	payments, _ := svc.ListPaymentsByRun(ctx, run.ID)
	id := payments[0].ID

	// PENDING → PROCESSING → PAID (dengan reference).
	proc, err := svc.UpdatePaymentStatus(ctx, id, UpdatePaymentStatusRequest{Status: "PROCESSING"})
	if err != nil {
		t.Fatalf("to PROCESSING: %v", err)
	}
	if proc.Status != "PROCESSING" || proc.ProcessedAt == nil {
		t.Errorf("expected PROCESSING with timestamp, got %s", proc.Status)
	}
	paid, err := svc.UpdatePaymentStatus(ctx, id, UpdatePaymentStatusRequest{Status: "PAID", Reference: "TRX-001"})
	if err != nil {
		t.Fatalf("to PAID: %v", err)
	}
	if paid.Status != "PAID" || paid.PaidAt == nil || paid.Reference != "TRX-001" {
		t.Errorf("expected PAID with ref TRX-001, got %s / %s", paid.Status, paid.Reference)
	}

	// PAID → FAILED tidak diizinkan.
	if _, err := svc.UpdatePaymentStatus(ctx, id, UpdatePaymentStatusRequest{Status: "FAILED", Reason: "x"}); err == nil {
		t.Error("expected invalid transition PAID→FAILED")
	}

	// PAID → REVERSED valid.
	rev, err := svc.UpdatePaymentStatus(ctx, id, UpdatePaymentStatusRequest{Status: "REVERSED"})
	if err != nil {
		t.Fatalf("to REVERSED: %v", err)
	}
	if rev.Status != "REVERSED" || rev.ReversedAt == nil {
		t.Errorf("expected REVERSED, got %s", rev.Status)
	}
}

func TestExportPaymentsCSV(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading := createTestGrading(ctx, repo, "G-1", "Staff")
	position := createTestPosition(ctx, repo, "Staff HR", &grading.ID)
	employee := createTestEmployee(ctx, repo, "EMP-CSV", "Iwan")
	createTestEmployment(ctx, repo, &employee.ID, &position.ID, "2020-01-01")
	profile := createTestPayrollProfile(ctx, repo, employee.ID, "2020-01-01")
	createTestBankProfile(ctx, repo, employee.ID, profile.ID, "777888999", "Iwan")

	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 10000000)

	svc := NewService(repo, zap.NewNop())
	run := calcRunForPayslip(t, svc, ctx, repo, "RUN-PAY-CSV")
	if _, err := svc.CreatePaymentBatch(ctx, run.ID); err != nil {
		t.Fatalf("create batch: %v", err)
	}

	csvOut, err := svc.ExportPaymentsCSV(ctx, run.ID)
	if err != nil {
		t.Fatalf("ExportPaymentsCSV: %v", err)
	}
	if !strings.Contains(csvOut, "employee_code") || !strings.Contains(csvOut, "777888999") || !strings.Contains(csvOut, "10000000.00") {
		t.Errorf("CSV missing expected content:\n%s", csvOut)
	}
}
