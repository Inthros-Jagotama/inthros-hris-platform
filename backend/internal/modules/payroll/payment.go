package payroll

import (
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
)

// PaymentStatus — status payment batch (sinkron dengan kolom status
// payroll_payments).
const (
	PaymentStatusPending    = "PENDING"
	PaymentStatusProcessing = "PROCESSING"
	PaymentStatusPaid       = "PAID"
	PaymentStatusFailed     = "FAILED"
	PaymentStatusReversed   = "REVERSED"
)

// validPaymentTransitions — transisi status yang diizinkan.
var validPaymentTransitions = map[string][]string{
	PaymentStatusPending:    {PaymentStatusProcessing, PaymentStatusFailed},
	PaymentStatusProcessing: {PaymentStatusPaid, PaymentStatusFailed, PaymentStatusPending},
	PaymentStatusPaid:       {PaymentStatusReversed},
	PaymentStatusFailed:     {PaymentStatusPending},
	PaymentStatusReversed:   {},
}

// CreatePaymentBatch membuat batch pembayaran dari run yang sudah final
// (CALCULATED/APPROVED/LOCKED): satu payment per run employee, nominal =
// net amount, rekening = SNAPSHOT dari bank profile utama yang aktif pada
// tanggal periode. Employee tanpa bank profile aktif dilewati. Bisa dipanggil
// ulang (batch lama dihapus lalu dibuat ulang).
func (s *Service) CreatePaymentBatch(ctx context.Context, runID string) (*PaymentBatchResponse, error) {
	uid, err := uuid.Parse(runID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	run, err := s.repo.FindPayrollRunByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if run.Status != "CALCULATED" && run.Status != "REVIEWED" && run.Status != "APPROVED" && run.Status != "LOCKED" {
		return nil, &ValidationError{Message: fmt.Sprintf("payment batch hanya bisa dibuat dari run final (CALCULATED/REVIEWED/APPROVED/LOCKED), saat ini %s", run.Status)}
	}
	period, err := s.repo.FindPayrollPeriodByID(ctx, run.PayrollPeriodID)
	if err != nil {
		return nil, err
	}
	runEmps, err := s.repo.FindPayrollRunEmployeesByRunID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if len(runEmps) == 0 {
		return nil, &ValidationError{Message: "run tidak punya employee — tidak ada payment yang bisa dibuat"}
	}

	actor := authctx.GetUserID(ctx)
	payments := make([]PayrollPayment, 0, len(runEmps))
	skipped := 0
	for _, emp := range runEmps {
		if emp.Status == "EXCLUDED" || emp.NetAmount <= 0 {
			continue
		}
		bank, err := s.repo.FindActivePrimaryBankProfileByEmployeeID(ctx, emp.EmployeeID, period.AsOfDate)
		if err != nil {
			return nil, err
		}
		if bank == nil {
			s.logger.Warn("Payment: employee tanpa bank profile aktif dilewati",
				zap.String("employee_id", emp.EmployeeID.String()),
				zap.String("run_id", runID),
			)
			skipped++
			continue
		}
		p := PayrollPayment{
			PayrollRunID:          uid,
			PayrollRunEmployeeID:  emp.ID,
			EmployeeID:            emp.EmployeeID,
			EmployeeCode:          emp.EmployeeCode,
			EmployeeName:          emp.EmployeeName,
			Amount:                emp.NetAmount,
			CurrencyCode:          "IDR",
			PaymentDate:           period.EndDate,
			EmployeeBankProfileID: &bank.ID,
			BankCode:              bank.BankCode,
			BankName:              stringPtr(bank.BankName),
			BankBranch:            bank.BankBranch,
			BankAccountNumber:     bank.BankAccountNumber,
			BankAccountHolderName: bank.BankAccountHolderName,
			Status:                PaymentStatusPending,
			CreatedBy:             actor,
			UpdatedBy:             actor,
		}
		payments = append(payments, p)
	}
	if len(payments) == 0 {
		return nil, &ValidationError{Message: "tidak ada employee dengan bank profile aktif yang layak dibayar"}
	}

	// Hapus batch lama dulu (regenerasi aman), lalu tulis batch baru.
	if err := s.repo.DeletePayrollPaymentsByRunID(ctx, uid); err != nil {
		return nil, err
	}
	if err := s.repo.BulkCreatePayrollPayments(ctx, payments); err != nil {
		return nil, err
	}

	total := 0.0
	for _, p := range payments {
		total += p.Amount
	}
	s.logger.Info("Payment batch created",
		zap.String("run_id", runID),
		zap.Int("payments", len(payments)),
		zap.Int("skipped_no_bank", skipped),
		zap.Float64("total", total),
	)
	return &PaymentBatchResponse{
		RunID:      runID,
		Total:      len(payments),
		TotalAmount: total,
		Skipped:    skipped,
		Status:     PaymentStatusPending,
	}, nil
}

// UpdatePaymentStatus memindahkan status payment sesuai transisi yang diizinkan.
// Reason dipakai saat FAILED, reference saat PAID/PROCESSING.
func (s *Service) UpdatePaymentStatus(ctx context.Context, id string, req UpdatePaymentStatusRequest) (*PayrollPaymentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindPayrollPaymentByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	allowed, ok := validPaymentTransitions[p.Status]
	if !ok {
		return nil, &ValidationError{Message: fmt.Sprintf("status payment tidak dikenal: %s", p.Status)}
	}
	if !containsString(allowed, req.Status) {
		return nil, &ValidationError{Message: fmt.Sprintf("transisi status %s → %s tidak diizinkan", p.Status, req.Status)}
	}

	now := time.Now()
	p.Status = req.Status
	p.UpdatedBy = authctx.GetUserID(ctx)
	switch req.Status {
	case PaymentStatusProcessing:
		p.ProcessedAt = &now
	case PaymentStatusPaid:
		p.PaidAt = &now
		if req.Reference != "" {
			p.Reference = &req.Reference
		}
	case PaymentStatusFailed:
		p.FailedAt = &now
		if req.Reason != "" {
			p.FailedReason = &req.Reason
		}
	case PaymentStatusReversed:
		p.ReversedAt = &now
	}
	if err := s.repo.UpdatePayrollPayment(ctx, p); err != nil {
		return nil, err
	}
	response := toPayrollPaymentResponse(p)
	return &response, nil
}

// ListPaymentsByRun mengembalikan seluruh payment sebuah run.
func (s *Service) ListPaymentsByRun(ctx context.Context, runID string) ([]PayrollPaymentResponse, error) {
	uid, err := uuid.Parse(runID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	payments, err := s.repo.FindPayrollPaymentsByRunID(ctx, uid)
	if err != nil {
		return nil, err
	}
	responses := make([]PayrollPaymentResponse, 0, len(payments))
	for i := range payments {
		responses = append(responses, toPayrollPaymentResponse(&payments[i]))
	}
	return responses, nil
}

// GetPaymentByID mengambil detail payment.
func (s *Service) GetPaymentByID(ctx context.Context, id string) (*PayrollPaymentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindPayrollPaymentByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := toPayrollPaymentResponse(p)
	return &response, nil
}

// ExportPaymentsCSV menghasilkan file bank transfer (CSV) dari payment batch
// sebuah run — berisi snapshot rekening & nominal per employee.
func (s *Service) ExportPaymentsCSV(ctx context.Context, runID string) (string, error) {
	uid, err := uuid.Parse(runID)
	if err != nil {
		return "", fmt.Errorf("invalid id: %w", err)
	}
	payments, err := s.repo.FindPayrollPaymentsByRunID(ctx, uid)
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	w := csv.NewWriter(&buf)
	header := []string{"employee_code", "employee_name", "bank_name", "account_number", "account_name", "amount", "status", "reference"}
	if err := w.Write(header); err != nil {
		return "", err
	}
	for _, p := range payments {
		bankName := ""
		if p.BankName != nil {
			bankName = *p.BankName
		}
		reference := ""
		if p.Reference != nil {
			reference = *p.Reference
		}
		row := []string{
			p.EmployeeCode,
			p.EmployeeName,
			bankName,
			p.BankAccountNumber,
			p.BankAccountHolderName,
			strconv.FormatFloat(p.Amount, 'f', 2, 64),
			p.Status,
			reference,
		}
		if err := w.Write(row); err != nil {
			return "", err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func containsString(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
