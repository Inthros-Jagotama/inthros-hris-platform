package payroll

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// =============================================================================
// DTO laporan (read-only, dari snapshot run — tidak ada tulis-balik)
// =============================================================================

// PayrollSummaryReport — Payroll Summary sebuah run.
type PayrollSummaryReport struct {
	RunID                string  `json:"run_id"`
	RunCode              string  `json:"run_code"`
	PeriodCode           string  `json:"period_code"`
	Status               string  `json:"status"`
	TotalEmployees       int     `json:"total_employees"`
	GrossSalary          float64 `json:"gross_salary"`
	EmployeeDeduction    float64 `json:"employee_deduction"`
	EmployerContribution float64 `json:"employer_contribution"`
	NetSalary            float64 `json:"net_salary"`
	TotalCompanyCost     float64 `json:"total_company_cost"`
}

// PayrollDetailReportRow — satu baris Payroll Detail (per employee per komponen).
type PayrollDetailReportRow struct {
	EmployeeID    string   `json:"employee_id"`
	EmployeeCode  string   `json:"employee_code"`
	EmployeeName  string   `json:"employee_name"`
	ComponentCode string   `json:"component_code"`
	ComponentName string   `json:"component_name"`
	ItemCategory  string   `json:"item_category"`
	ComponentType string   `json:"component_type"`
	PaidBy        string   `json:"paid_by"`
	SourceGroup   string   `json:"source_group"`
	BaseAmount    float64  `json:"base_amount"`
	Rate          *float64 `json:"rate,omitempty"`
	Amount        float64  `json:"amount"`
}

// BpjsReportRow — satu baris BPJS Report per employee.
type BpjsReportRow struct {
	EmployeeCode           string  `json:"employee_code"`
	EmployeeName           string  `json:"employee_name"`
	BpjsNumber             string  `json:"bpjs_number,omitempty"`
	WageBasis              float64 `json:"wage_basis"`
	EmployeeContribution   float64 `json:"employee_contribution"`
	EmployerContribution   float64 `json:"employer_contribution"`
	TotalContribution      float64 `json:"total_contribution"`
}

// TaxReportRow — satu baris Tax Report per employee.
type TaxReportRow struct {
	EmployeeCode   string  `json:"employee_code"`
	EmployeeName   string  `json:"employee_name"`
	TaxableIncome  float64 `json:"taxable_income"`
	Pph21          float64 `json:"pph21"`
}

// BankTransferReportRow — satu baris Bank Transfer report per payment.
type BankTransferReportRow struct {
	EmployeeCode       string  `json:"employee_code"`
	EmployeeName       string  `json:"employee_name"`
	BankName           string  `json:"bank_name,omitempty"`
	AccountNumber      string  `json:"account_number"`
	AccountHolderName  string  `json:"account_holder_name"`
	Amount             float64 `json:"amount"`
	Status             string  `json:"status"`
}

// =============================================================================
// Service
// =============================================================================

// GetPayrollSummaryReport menghasilkan ringkasan run (Payroll Summary).
func (s *Service) GetPayrollSummaryReport(ctx context.Context, runID string) (*PayrollSummaryReport, error) {
	uid, err := uuid.Parse(runID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	run, err := s.repo.FindPayrollRunByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	period, err := s.repo.FindPayrollPeriodByID(ctx, run.PayrollPeriodID)
	if err != nil {
		return nil, err
	}
	return &PayrollSummaryReport{
		RunID:                run.ID.String(),
		RunCode:              run.RunCode,
		PeriodCode:           period.PeriodCode,
		Status:               run.Status,
		TotalEmployees:       run.TotalEmployees,
		GrossSalary:          run.TotalEarning,
		EmployeeDeduction:    run.TotalDeduction,
		EmployerContribution: run.TotalEmployerContribution,
		NetSalary:            run.TotalNet,
		TotalCompanyCost:     run.TotalCompanyCost,
	}, nil
}

// GetPayrollDashboard menghasilkan agregat dashboard run (sama dengan summary,
// plus status — data sudah tersedia di payroll_runs).
func (s *Service) GetPayrollDashboard(ctx context.Context, runID string) (*PayrollSummaryReport, error) {
	return s.GetPayrollSummaryReport(ctx, runID)
}

// GetPayrollDetailReport menghasilkan Payroll Detail: per employee per komponen
// (base, rate, amount, kategori) dari snapshot payroll_run_items.
func (s *Service) GetPayrollDetailReport(ctx context.Context, runID string) ([]PayrollDetailReportRow, error) {
	uid, err := uuid.Parse(runID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	if _, err := s.repo.FindPayrollRunByID(ctx, uid); err != nil {
		return nil, err
	}
	items, err := s.repo.FindPayrollRunItemsByRunID(ctx, uid)
	if err != nil {
		return nil, err
	}
	runEmps, err := s.repo.FindPayrollRunEmployeesByRunID(ctx, uid)
	if err != nil {
		return nil, err
	}
	nameByEmp := map[uuid.UUID]PayrollRunEmployee{}
	for _, e := range runEmps {
		nameByEmp[e.ID] = e
	}
	rows := make([]PayrollDetailReportRow, 0, len(items))
	for _, it := range items {
		code, name := "", ""
		if emp, ok := nameByEmp[it.PayrollRunEmployeeID]; ok {
			code, name = emp.EmployeeCode, emp.EmployeeName
		}
		rows = append(rows, PayrollDetailReportRow{
			EmployeeID:    it.EmployeeID.String(),
			EmployeeCode:  code,
			EmployeeName:  name,
			ComponentCode: it.ComponentCode,
			ComponentName: it.ComponentName,
			ItemCategory:  it.ItemCategory,
			ComponentType: it.ComponentType,
			PaidBy:        it.PaidBy,
			SourceGroup:   it.SourceGroup,
			BaseAmount:    it.BaseAmount,
			Rate:          it.Rate,
			Amount:        it.Amount,
		})
	}
	return rows, nil
}

// GetBpjsReport menghasilkan BPJS Report per employee dari item statutori
// (source_group=STATUTORY, program BPJS) + nomor kepesertaan dari profil BPJS.
func (s *Service) GetBpjsReport(ctx context.Context, runID string) ([]BpjsReportRow, error) {
	uid, err := uuid.Parse(runID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	runEmps, err := s.repo.FindPayrollRunEmployeesByRunID(ctx, uid)
	if err != nil {
		return nil, err
	}
	items, err := s.repo.FindPayrollRunItemsByRunID(ctx, uid)
	if err != nil {
		return nil, err
	}

	// Nomor kepesertaan dari profil BPJS (satu query untuk semua employee).
	empIDs := make([]uuid.UUID, 0, len(runEmps))
	for _, e := range runEmps {
		empIDs = append(empIDs, e.EmployeeID)
	}
	profiles, err := s.repo.FindEmployeeBpjsProfilesByEmployeeIDs(ctx, empIDs)
	if err != nil {
		return nil, err
	}
	numberByEmp := map[uuid.UUID]string{}
	for _, p := range profiles {
		num := ""
		if p.BpjsHealthNo != nil {
			num = *p.BpjsHealthNo
		} else if p.BpjsTkNo != nil {
			num = *p.BpjsTkNo
		}
		if _, seen := numberByEmp[p.EmployeeID]; !seen {
			numberByEmp[p.EmployeeID] = num
		}
	}

	// Agregat per employee dari item BPJS (dibedakan employee vs employer).
	type agg struct {
		wageBasis, empContrib, erContrib float64
	}
	byEmp := map[uuid.UUID]*agg{}
	for _, it := range items {
		if it.SourceGroup != SourceGroupStatutory || it.SourceType == nil {
			continue
		}
		if !isBpjsProgram(*it.SourceType) {
			continue
		}
		a, ok := byEmp[it.EmployeeID]
		if !ok {
			a = &agg{}
			byEmp[it.EmployeeID] = a
		}
		if it.BaseAmount > a.wageBasis {
			a.wageBasis = it.BaseAmount
		}
		switch it.ItemCategory {
		case ItemCategoryEmployeeDeduction:
			a.empContrib += it.Amount
		case ItemCategoryEmployerContribution:
			a.erContrib += it.Amount
		}
	}

	rows := make([]BpjsReportRow, 0, len(runEmps))
	for _, e := range runEmps {
		a, ok := byEmp[e.EmployeeID]
		if !ok {
			continue // employee tanpa BPJS tidak masuk laporan
		}
		rows = append(rows, BpjsReportRow{
			EmployeeCode:         e.EmployeeCode,
			EmployeeName:         e.EmployeeName,
			BpjsNumber:           numberByEmp[e.EmployeeID],
			WageBasis:            a.wageBasis,
			EmployeeContribution: a.empContrib,
			EmployerContribution: a.erContrib,
			TotalContribution:    a.empContrib + a.erContrib,
		})
	}
	return rows, nil
}

// GetTaxReport menghasilkan Tax Report per employee dari pph21_calculation_logs
// (penghasilan kena pajak bulanan + PPh21 bulanan).
func (s *Service) GetTaxReport(ctx context.Context, runID string) ([]TaxReportRow, error) {
	uid, err := uuid.Parse(runID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	logs, err := s.repo.FindPph21CalculationLogsByRunID(ctx, uid)
	if err != nil {
		return nil, err
	}
	runEmps, err := s.repo.FindPayrollRunEmployeesByRunID(ctx, uid)
	if err != nil {
		return nil, err
	}
	nameByEmp := map[uuid.UUID]PayrollRunEmployee{}
	for _, e := range runEmps {
		nameByEmp[e.EmployeeID] = e
	}
	rows := make([]TaxReportRow, 0, len(logs))
	for _, l := range logs {
		code, name := "", ""
		if emp, ok := nameByEmp[l.EmployeeID]; ok {
			code, name = emp.EmployeeCode, emp.EmployeeName
		}
		rows = append(rows, TaxReportRow{
			EmployeeCode:  code,
			EmployeeName:  name,
			TaxableIncome: l.GrossMonthly,
			Pph21:         l.Pph21Monthly,
		})
	}
	return rows, nil
}

// GetBankTransferReport menghasilkan Bank Transfer report dari payroll_payments.
func (s *Service) GetBankTransferReport(ctx context.Context, runID string) ([]BankTransferReportRow, error) {
	uid, err := uuid.Parse(runID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	payments, err := s.repo.FindPayrollPaymentsByRunID(ctx, uid)
	if err != nil {
		return nil, err
	}
	rows := make([]BankTransferReportRow, 0, len(payments))
	for _, p := range payments {
		bankName := ""
		if p.BankName != nil {
			bankName = *p.BankName
		}
		rows = append(rows, BankTransferReportRow{
			EmployeeCode:      p.EmployeeCode,
			EmployeeName:      p.EmployeeName,
			BankName:          bankName,
			AccountNumber:     p.BankAccountNumber,
			AccountHolderName: p.BankAccountHolderName,
			Amount:            p.Amount,
			Status:            p.Status,
		})
	}
	return rows, nil
}

// isBpjsProgram true jika source type adalah program BPJS.
func isBpjsProgram(sourceType string) bool {
	switch sourceType {
	case BpjsProgramHealth, BpjsProgramJHT, BpjsProgramJP, BpjsProgramJKK, BpjsProgramJKM, BpjsProgramJKP:
		return true
	default:
		return false
	}
}
