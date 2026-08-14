package payroll

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/modules/payroll/calculator"
)

// pph21FormulaBreakdown adalah jejak rinci kalkulasi PPh21 yang disimpan ke
// kolom formula_json di pph21_calculation_logs (migration 007).
type pph21FormulaBreakdown struct {
	GrossMonthly                float64 `json:"gross_monthly"`
	OccupationalExpenseMonthly  float64 `json:"occupational_expense_monthly"`
	BpjsTaxDeductibleMonthly    float64 `json:"bpjs_tax_deductible_monthly"`
	NetMonthly                  float64 `json:"net_monthly"`
	AnnualizationMonths         int     `json:"annualization_months"`
	NetAnnualized               float64 `json:"net_annualized"`
	PtkpStatus                  string  `json:"ptkp_status"`
	PtkpAnnual                  float64 `json:"ptkp_annual"`
	PkpAnnual                   float64 `json:"pkp_annual"`
	PkpRoundingUnit             float64 `json:"pkp_rounding_unit"`
	AnnualTaxBeforeNpwpMult     float64 `json:"annual_tax_before_npwp_mult"`
	NonNpwpMultiplierPercent    float64 `json:"non_npwp_multiplier_percent"`
	AnnualTaxAfterNpwpMult      float64 `json:"annual_tax_after_npwp_mult"`
	Pph21Monthly                float64 `json:"pph21_monthly"`
	AppliedNonNpwpMultiplier    float64 `json:"applied_non_npwp_multiplier_percent"`
	Brackets                    []pph21BracketBreakdown `json:"brackets"`
}

type pph21BracketBreakdown struct {
	Order       int     `json:"order"`
	LowerBound  float64 `json:"lower_bound"`
	UpperBound  *float64 `json:"upper_bound,omitempty"`
	RatePercent float64 `json:"rate_percent"`
	TaxAmount   float64 `json:"tax_amount"`
}

// calculatePph21 menghitung PPh21 bulanan (metode REGULAR_GROSS_ANNUALIZED) untuk
// satu employee dan menghasilkan:
//   - PayrollRunItem potongan PPh21 (source_group=STATUTORY), atau nil jika
//     employee tidak taxable / tidak ada setting / tidak ada profil pajak.
//   - Pph21CalculationLog berisi jejak rinci kalkulasi (PayrollRunEmployeeID &
//     PayrollRunItemID diisi saat persistRunSnapshot).
//
// Dipanggil SETELAH struktur & BPJS dihitung karena butuh gross taxable (values)
// dan iuran BPJS employee yang boleh dikurangkan (bpjsItems).
func (s *Service) calculatePph21(ctx context.Context, run *PayrollRun, period *PayrollPeriod, ec runEmployeeContext, components []SalaryComponent, values map[string]float64, bpjsItems []PayrollRunItem) (*PayrollRunItem, *Pph21CalculationLog, error) {
	asOf := period.AsOfDate

	setting, err := s.repo.FindActivePph21SettingByDate(ctx, asOf)
	if err != nil {
		return nil, nil, err
	}
	if setting == nil {
		s.logger.Debug("PPh21: no active setting, skipping",
			zap.String("employee_id", ec.employee.ID.String()),
			zap.String("as_of", asOf),
		)
		return nil, nil, nil
	}

	profile, err := s.repo.FindActiveEmployeeTaxProfileByEmployeeID(ctx, ec.employee.ID, asOf)
	if err != nil {
		return nil, nil, err
	}
	if profile == nil || !profile.IsTaxable {
		s.logger.Debug("PPh21: employee not taxable or no tax profile, skipping",
			zap.String("employee_id", ec.employee.ID.String()),
		)
		return nil, nil, nil
	}

	// Komponen potongan PPh21 (dari setting) harus aktif.
	var pph21Comp *SalaryComponent
	for i := range components {
		if components[i].ID == setting.Pph21ComponentID {
			pph21Comp = &components[i]
			break
		}
	}
	if pph21Comp == nil {
		s.logger.Warn("PPh21: setting merujuk komponen non-aktif, di-skip",
			zap.String("setting_code", setting.SettingCode),
			zap.String("employee_id", ec.employee.ID.String()),
		)
		return nil, nil, nil
	}

	// 1. Gross bulanan = total nilai komponen bertanda is_taxable.
	gross := 0.0
	for _, c := range components {
		if !c.IsTaxable {
			continue
		}
		if v, ok := values[c.Code]; ok {
			gross += v
		}
	}

	// 2. Biaya jabatan: min(gross * rate%, max bulanan).
	occ := gross * setting.OccupationalExpenseRatePercent / 100
	if occ > setting.OccupationalExpenseMaxMonthly {
		occ = setting.OccupationalExpenseMaxMonthly
	}

	// 3. Iuran BPJS employee yang boleh dikurangkan sesuai flag setting.
	bpjsDeductible := 0.0
	for _, item := range bpjsItems {
		if item.ItemCategory != ItemCategoryEmployeeDeduction || item.PaidBy != "EMPLOYEE" {
			continue
		}
		program := ""
		if item.SourceType != nil {
			program = *item.SourceType
		}
		switch program {
		case BpjsProgramHealth:
			if setting.DeductBpjsHealthEmployee {
				bpjsDeductible += item.Amount
			}
		case BpjsProgramJHT:
			if setting.DeductBpjsJhtEmployee {
				bpjsDeductible += item.Amount
			}
		case BpjsProgramJP:
			if setting.DeductBpjsJpEmployee {
				bpjsDeductible += item.Amount
			}
		}
	}

	// 4. Net bulanan → annualisasi.
	netMonthly := gross - occ - bpjsDeductible
	if netMonthly < 0 {
		netMonthly = 0
	}
	months := setting.AnnualizationMonths
	if months <= 0 {
		months = 12
	}
	netAnnualized := netMonthly * float64(months)

	// 5. PTKP tahunan dari status profil employee.
	ptkpRates, err := s.repo.FindActivePph21PtkpRatesByDate(ctx, asOf)
	if err != nil {
		return nil, nil, err
	}
	ptkpByStatus := map[string]float64{}
	for _, r := range ptkpRates {
		ptkpByStatus[r.PtkpStatus] = r.AnnualAmount
	}
	ptkpStatus := ""
	if profile.PtkpStatus != nil {
		ptkpStatus = *profile.PtkpStatus
	}
	ptkpAnnual := 0.0
	if ptkpStatus != "" {
		if amt, ok := ptkpByStatus[ptkpStatus]; ok {
			ptkpAnnual = amt
		} else {
			s.logger.Warn("PPh21: PTKP rate tidak ditemukan untuk status, dipakai 0",
				zap.String("ptkp_status", ptkpStatus),
				zap.String("employee_id", ec.employee.ID.String()),
			)
		}
	}

	// 6. PKP tahunan, dibulatkan ke bawah ke unit pembulatan.
	rounding := calculator.NormalizeRoundingMode(setting.RoundingMode)
	unit := setting.PkpRoundingUnit
	if unit <= 0 {
		unit = 1000
	}
	pkp := netAnnualized - ptkpAnnual
	if pkp < 0 {
		pkp = 0
	}
	pkp = calculator.RoundToUnit(pkp, unit, calculator.RoundingFloor)

	// 7. Pajak progresif per bracket.
	brackets, err := s.repo.FindActivePph21TaxBracketsByDate(ctx, asOf)
	if err != nil {
		return nil, nil, err
	}
	annualTax := 0.0
	bracketBreakdowns := make([]pph21BracketBreakdown, 0, len(brackets))
	for _, b := range brackets {
		upper := math.Inf(1)
		if b.UpperBound != nil {
			upper = *b.UpperBound
		}
		portion := math.Min(pkp, upper) - b.LowerBound
		taxAmount := 0.0
		if portion > 0 {
			taxAmount = portion * b.RatePercent / 100
			annualTax += taxAmount
		}
		bracketBreakdowns = append(bracketBreakdowns, pph21BracketBreakdown{
			Order:       b.BracketOrder,
			LowerBound:  b.LowerBound,
			UpperBound:  b.UpperBound,
			RatePercent: b.RatePercent,
			TaxAmount:   taxAmount,
		})
	}

	// 8. Multiplier non-NPWP (employee tanpa NPWP dikenakan lebih tinggi sesuai
	// konfigurasi setting; default 100 = tanpa penambahan).
	mult := setting.NonNpwpMultiplierPercent
	appliedMult := mult
	if profile.HasNpwp {
		appliedMult = 100
	}
	annualTaxAfter := annualTax * appliedMult / 100
	pph21Monthly := calculator.Round(annualTaxAfter/float64(months), rounding)

	breakdown := pph21FormulaBreakdown{
		GrossMonthly:             gross,
		OccupationalExpenseMonthly: occ,
		BpjsTaxDeductibleMonthly: bpjsDeductible,
		NetMonthly:               netMonthly,
		AnnualizationMonths:      months,
		NetAnnualized:            netAnnualized,
		PtkpStatus:               ptkpStatus,
		PtkpAnnual:               ptkpAnnual,
		PkpAnnual:                pkp,
		PkpRoundingUnit:          unit,
		AnnualTaxBeforeNpwpMult:  annualTax,
		NonNpwpMultiplierPercent: mult,
		AnnualTaxAfterNpwpMult:   annualTaxAfter,
		Pph21Monthly:             pph21Monthly,
		AppliedNonNpwpMultiplier: appliedMult,
		Brackets:                 bracketBreakdowns,
	}
	formulaJSON, _ := json.Marshal(breakdown)
	formulaStr := string(formulaJSON)

	item := s.buildPph21RunItem(ec, *setting, *pph21Comp, pph21Monthly, gross, &formulaStr)
	log := &Pph21CalculationLog{
		PayrollRunID:              run.ID,
		EmployeeID:                ec.employee.ID,
		Pph21SettingID:            setting.ID,
		EmployeeTaxProfileID:      profile.ID,
		TaxMethod:                 profile.TaxMethod,
		PtkpStatus:                ptkpStatus,
		HasNpwp:                   profile.HasNpwp,
		GrossMonthly:              gross,
		OccupationalExpenseMonthly: occ,
		BpjsTaxDeductibleMonthly:  bpjsDeductible,
		NetMonthly:                netMonthly,
		NetAnnualized:             netAnnualized,
		PtkpAnnual:                ptkpAnnual,
		PkpAnnual:                 pkp,
		AnnualTaxBeforeNpwpMult:   annualTax,
		NonNpwpMultiplierPercent:  mult,
		AnnualTaxAfterNpwpMult:    annualTaxAfter,
		Pph21Monthly:              pph21Monthly,
		FormulaJSON:               &formulaStr,
		Status:                    "CALCULATED",
	}
	return &item, log, nil
}

// buildPph21RunItem membangun snapshot item potongan PPh21 (EMPLOYEE_DEDUCTION,
// source_group=STATUTORY).
func (s *Service) buildPph21RunItem(ec runEmployeeContext, setting Pph21Setting, comp SalaryComponent, amount, gross float64, formulaJSON *string) PayrollRunItem {
	sourceTable := "pph21_settings"
	sourceType := "PPH21"
	notes := fmt.Sprintf("PPh21 %s", setting.CalculationMethod)
	item := PayrollRunItem{
		EmployeeID:        ec.employee.ID,
		SalaryComponentID: comp.ID,
		ComponentCode:     comp.Code,
		ComponentName:     comp.Name,
		ComponentType:     comp.ComponentType,
		CalculationType:   "FORMULA",
		ItemCategory:      ItemCategoryEmployeeDeduction,
		PaidBy:            "EMPLOYEE",
		AffectsNetPay:     true,
		Amount:            amount,
		BaseAmount:        gross,
		FormulaResult:     &amount,
		Formula:           formulaJSON,
		CurrencyCode:      "IDR",
		SourceGroup:       SourceGroupStatutory,
		SourceTable:       &sourceTable,
		SourceID:          &setting.ID,
		SourceType:        &sourceType,
		PrintOnPayslip:    true,
		Notes:             &notes,
	}
	return item
}
