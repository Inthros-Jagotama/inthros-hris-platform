package payroll

import (
	"context"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
)

// PayslipStatus — status payslip (sinkron dengan kolom status payroll_payslips).
const (
	PayslipStatusDRAFT     = "DRAFT"
	PayslipStatusPublished = "PUBLISHED"
	PayslipStatusCancelled = "CANCELLED"
)

// GeneratePayslips membuat satu payslip per run employee dari snapshot run
// (payroll_run_employees + payroll_run_items). Run harus sudah dihitung
// (CALCULATED/APPROVED/LOCKED). Bisa dipanggil ulang: payslip lama dihapus
// lalu diisi ulang (regenerasi). Mengembalikan daftar payslip yang dibuat.
func (s *Service) GeneratePayslips(ctx context.Context, runID string) ([]PayrollPayslipResponse, error) {
	uid, err := uuid.Parse(runID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	run, err := s.repo.FindPayrollRunByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if run.Status != "CALCULATED" && run.Status != "REVIEWED" && run.Status != "APPROVED" && run.Status != "LOCKED" {
		return nil, &ValidationError{Message: fmt.Sprintf("payslip hanya bisa dibuat dari run final (CALCULATED/REVIEWED/APPROVED/LOCKED), saat ini %s", run.Status)}
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
		return nil, &ValidationError{Message: "run tidak punya employee — tidak ada payslip yang bisa dibuat"}
	}
	items, err := s.repo.FindPayrollRunItemsByRunID(ctx, uid)
	if err != nil {
		return nil, err
	}
	itemsByEmp := map[uuid.UUID][]PayrollRunItem{}
	for _, it := range items {
		itemsByEmp[it.PayrollRunEmployeeID] = append(itemsByEmp[it.PayrollRunEmployeeID], it)
	}

	actor := authctx.GetUserID(ctx)
	now := time.Now()
	// Hapus payslip lama dulu (regenerasi aman).
	if err := s.repo.DeletePayrollPayslipsByRunID(ctx, uid); err != nil {
		return nil, err
	}

	payslips := make([]PayrollPayslip, 0, len(runEmps))
	for i, emp := range runEmps {
		if emp.Status == "EXCLUDED" {
			continue
		}
		ps := PayrollPayslip{
			PayrollRunID:              uid,
			PayrollRunEmployeeID:      emp.ID,
			EmployeeID:                emp.EmployeeID,
			PayslipNumber:             fmt.Sprintf("SLP-%s-%03d", period.PeriodCode, i+1),
			PeriodYear:                period.PeriodYear,
			PeriodMonth:               period.PeriodMonth,
			PeriodCode:                period.PeriodCode,
			EmployeeCode:              emp.EmployeeCode,
			EmployeeName:              emp.EmployeeName,
			PositionTitle:             emp.PositionTitle,
			GradingName:               emp.GradingName,
			TotalEarning:              emp.TotalEarning,
			TotalDeduction:            emp.TotalDeduction,
			TotalEmployerContribution: emp.TotalEmployerContribution,
			NetAmount:                 emp.NetAmount,
			Status:                    PayslipStatusDRAFT,
			GeneratedAt:               &now,
			CreatedBy:                 actor,
			UpdatedBy:                 actor,
		}
		payslips = append(payslips, ps)
	}
	if len(payslips) == 0 {
		return nil, &ValidationError{Message: "tidak ada employee yang layak mendapat payslip"}
	}

	// Simpan payslip satu-satu supaya ID hasil insert bisa dipakai response
	// (jumlah payslip per run umumnya kecil).
	responses := make([]PayrollPayslipResponse, 0, len(payslips))
	for i := range payslips {
		if err := s.repo.CreatePayrollPayslip(ctx, &payslips[i]); err != nil {
			return nil, err
		}
		responses = append(responses, toPayrollPayslipResponse(&payslips[i], itemsByEmp[payslips[i].PayrollRunEmployeeID]))
	}
	s.logger.Info("Payslips generated",
		zap.String("run_id", runID),
		zap.Int("count", len(responses)),
	)
	return responses, nil
}

// PublishPayslip memindahkan payslip DRAFT → PUBLISHED (bisa dilihat employee).
func (s *Service) PublishPayslip(ctx context.Context, id string) (*PayrollPayslipResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	ps, err := s.repo.FindPayrollPayslipByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if ps.Status != PayslipStatusDRAFT {
		return nil, &ValidationError{Message: fmt.Sprintf("hanya payslip DRAFT yang bisa dipublikasikan, saat ini %s", ps.Status)}
	}
	now := time.Now()
	ps.Status = PayslipStatusPublished
	ps.PublishedAt = &now
	ps.UpdatedBy = authctx.GetUserID(ctx)
	if err := s.repo.UpdatePayrollPayslip(ctx, ps); err != nil {
		return nil, err
	}
	items, _ := s.repo.FindPayrollRunItemsByEmployeeID(ctx, ps.PayrollRunID, ps.EmployeeID)
	response := toPayrollPayslipResponse(ps, items)
	return &response, nil
}

// CancelPayslip membatalkan payslip (DRAFT/PUBLISHED → CANCELLED).
func (s *Service) CancelPayslip(ctx context.Context, id string) (*PayrollPayslipResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	ps, err := s.repo.FindPayrollPayslipByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if ps.Status == PayslipStatusCancelled {
		return nil, &ValidationError{Message: "payslip sudah dibatalkan"}
	}
	now := time.Now()
	ps.Status = PayslipStatusCancelled
	ps.CancelledAt = &now
	ps.UpdatedBy = authctx.GetUserID(ctx)
	if err := s.repo.UpdatePayrollPayslip(ctx, ps); err != nil {
		return nil, err
	}
	items, _ := s.repo.FindPayrollRunItemsByEmployeeID(ctx, ps.PayrollRunID, ps.EmployeeID)
	response := toPayrollPayslipResponse(ps, items)
	return &response, nil
}

// ListPayslipsByRun mengembalikan seluruh payslip sebuah run.
func (s *Service) ListPayslipsByRun(ctx context.Context, runID string) ([]PayrollPayslipResponse, error) {
	uid, err := uuid.Parse(runID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	payslips, err := s.repo.FindPayrollPayslipsByRunID(ctx, uid)
	if err != nil {
		return nil, err
	}
	responses := make([]PayrollPayslipResponse, 0, len(payslips))
	for i := range payslips {
		responses = append(responses, toPayrollPayslipResponse(&payslips[i], nil))
	}
	return responses, nil
}

// GetPayslipByID mengambil detail payslip (termasuk item earnings/deductions).
func (s *Service) GetPayslipByID(ctx context.Context, id string) (*PayrollPayslipResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	ps, err := s.repo.FindPayrollPayslipByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	items, _ := s.repo.FindPayrollRunItemsByEmployeeID(ctx, ps.PayrollRunID, ps.EmployeeID)
	response := toPayrollPayslipResponse(ps, items)
	return &response, nil
}

// GetPayslipHTML merender payslip sebagai HTML (format server-side) — bisa
// diserve ke employee portal atau dikonversi PDF oleh lapisan lain.
func (s *Service) GetPayslipHTML(ctx context.Context, id string) (string, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return "", fmt.Errorf("invalid id: %w", err)
	}
	ps, err := s.repo.FindPayrollPayslipByID(ctx, uid)
	if err != nil {
		return "", err
	}
	items, err := s.repo.FindPayrollRunItemsByEmployeeID(ctx, ps.PayrollRunID, ps.EmployeeID)
	if err != nil {
		return "", err
	}
	return renderPayslipHTML(ps, items)
}

// renderPayslipHTML membangun HTML payslip dari snapshot run (item per
// kategori: earning, deduction, employer contribution).
func renderPayslipHTML(ps *PayrollPayslip, items []PayrollRunItem) (string, error) {
	type row struct {
		Name   string
		Amount float64
	}
	var earnings, deductions, contributions []row
	for _, it := range items {
		r := row{Name: it.ComponentName, Amount: it.Amount}
		switch it.ItemCategory {
		case ItemCategoryEmployeeEarning:
			earnings = append(earnings, r)
		case ItemCategoryEmployeeDeduction:
			deductions = append(deductions, r)
		case ItemCategoryEmployerContribution:
			contributions = append(contributions, r)
		}
	}

	position := ""
	if ps.PositionTitle != nil {
		position = *ps.PositionTitle
	}
	grading := ""
	if ps.GradingName != nil {
		grading = *ps.GradingName
	}

	tmpl, err := template.New("payslip").Parse(payslipHTMLTemplate)
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	data := map[string]any{
		"PayslipNumber":   ps.PayslipNumber,
		"PeriodCode":      ps.PeriodCode,
		"EmployeeCode":    ps.EmployeeCode,
		"EmployeeName":    ps.EmployeeName,
		"PositionTitle":   position,
		"GradingName":     grading,
		"Earnings":        earnings,
		"Deductions":      deductions,
		"Contributions":   contributions,
		"TotalEarning":    ps.TotalEarning,
		"TotalDeduction":  ps.TotalDeduction,
		"NetAmount":       ps.NetAmount,
		"CompanyCost":     ps.TotalEarning + ps.TotalEmployerContribution,
		"Status":          ps.Status,
		"GeneratedAt":     formatTime(ps.GeneratedAt),
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04")
}

const payslipHTMLTemplate = `<!DOCTYPE html>
<html lang="id">
<head><meta charset="utf-8"><title>Slip Gaji {{.PayslipNumber}}</title></head>
<body style="font-family: Arial, sans-serif; max-width: 720px; margin: 0 auto;">
  <h2 style="text-align:center; margin-bottom:0;">SLIP GAJI</h2>
  <p style="text-align:center; margin-top:2px;">{{.PeriodCode}} · {{.PayslipNumber}}</p>
  <table style="width:100%; border-collapse:collapse; margin: 12px 0;">
    <tr><td><b>Nama</b></td><td>{{.EmployeeName}} ({{.EmployeeCode}})</td></tr>
    <tr><td><b>Posisi</b></td><td>{{.PositionTitle}}{{if .GradingName}} · {{.GradingName}}{{end}}</td></tr>
    <tr><td><b>Status</b></td><td>{{.Status}} · dibuat {{.GeneratedAt}}</td></tr>
  </table>
  <h3>Penghasilan</h3>
  <table style="width:100%; border-collapse:collapse;">
    {{range .Earnings}}<tr><td>{{.Name}}</td><td style="text-align:right;">Rp {{printf "%.0f" .Amount}}</td></tr>{{end}}
    <tr style="font-weight:bold; border-top:1px solid #000;"><td>Total Penghasilan</td><td style="text-align:right;">Rp {{printf "%.0f" .TotalEarning}}</td></tr>
  </table>
  <h3>Potongan</h3>
  <table style="width:100%; border-collapse:collapse;">
    {{range .Deductions}}<tr><td>{{.Name}}</td><td style="text-align:right;">Rp {{printf "%.0f" .Amount}}</td></tr>{{end}}
    <tr style="font-weight:bold; border-top:1px solid #000;"><td>Total Potongan</td><td style="text-align:right;">Rp {{printf "%.0f" .TotalDeduction}}</td></tr>
  </table>
  <h3>Iuran Perusahaan</h3>
  <table style="width:100%; border-collapse:collapse;">
    {{range .Contributions}}<tr><td>{{.Name}}</td><td style="text-align:right;">Rp {{printf "%.0f" .Amount}}</td></tr>{{end}}
  </table>
  <table style="width:100%; border-collapse:collapse; margin-top:16px; font-size:1.1em;">
    <tr><td><b>GAJI BERSIH (NET)</b></td><td style="text-align:right;"><b>Rp {{printf "%.0f" .NetAmount}}</b></td></tr>
    <tr><td>Total Biaya Perusahaan</td><td style="text-align:right;">Rp {{printf "%.0f" .CompanyCost}}</td></tr>
  </table>
</body>
</html>`
