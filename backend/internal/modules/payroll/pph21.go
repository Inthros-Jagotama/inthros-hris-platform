package payroll

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/modules/payroll/calculator"
)

// pph21FormulaBreakdown adalah jejak rinci kalkulasi PPh21 yang disimpan ke
// kolom formula_json di pph21_calculation_logs (migration 007).
type pph21FormulaBreakdown struct {
	CalculationMethod           string  `json:"calculation_method"`
	TerCategory                 string  `json:"ter_category,omitempty"`
	TerRate                     float64 `json:"ter_rate,omitempty"`
	GrossMonthly                float64 `json:"gross_monthly"`
	OccupationalExpenseMonthly  float64 `json:"occupational_expense_monthly"`
	BpjsTaxDeductibleMonthly    float64 `json:"bpjs_tax_deductible_monthly"`
	PensionDeductibleMonthly    float64 `json:"pension_deductible_monthly"`
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
	YtdDeductedMonthly          float64 `json:"ytd_deducted_monthly,omitempty"`
	Brackets                    []pph21BracketBreakdown `json:"brackets"`
}

type pph21BracketBreakdown struct {
	Order       int     `json:"order"`
	LowerBound  float64 `json:"lower_bound"`
	UpperBound  *float64 `json:"upper_bound,omitempty"`
	RatePercent float64 `json:"rate_percent"`
	TaxAmount   float64 `json:"tax_amount"`
}

// normalizePtkpCode menormalkan status PTKP menjadi kode lookup tabel ptkps:
// hilangkan slash/spasi/kapitalisasi ("TK/0" ≡ "TK0").
func normalizePtkpCode(status string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(status, "/", ""), " ", ""))
}

// terCategoryForPtkpStatus memetakan status PTKP ke kategori TER sesuai aturan
// resmi PER-2/PJ/2024 (PP 58/2023): A = TK/0, TK/1, K/0 · B = TK/2, TK/3, K/1,
// K/2 · C = K/3 dan K/I/*. Status tidak dikenal → kategori kosong.
// Dipakai sebagai fallback bila baris ptkps tidak ditemukan (engine biasanya
// memakai kolom `group` di tabel ptkps).
func terCategoryForPtkpStatus(status string) string {
	norm := normalizePtkpCode(status)
	switch norm {
	case "TK0", "TK1", "K0":
		return "A"
	case "TK2", "TK3", "K1", "K2":
		return "B"
	case "K3", "KI0", "KI1", "KI2", "KI3":
		return "C"
	}
	return ""
}

// calculatePph21 menghitung PPh21 bulanan (metode REGULAR_GROSS_ANNUALIZED atau
// TER PP 58/2023) untuk satu employee dan menghasilkan:
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

	// Komponen potongan PPh21 (wadah hasil pajak): ditandai flag
	// is_pph21_component di komponen gaji (sumber kebenaran tunggal).
	var pph21Comp *SalaryComponent
	for i := range components {
		if components[i].IsPph21Component {
			pph21Comp = &components[i]
			break
		}
	}
	if pph21Comp == nil {
		s.logger.Warn("PPh21: tidak ada komponen ber-flag is_pph21_component, di-skip",
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

	// Status PTKP & kategori TER (dibutuhkan oleh kedua metode). PTKP diambil
	// dari tabel ptkps (satu sumber kebenaran, migration 121) via kode
	// ternormalisasi — berisi nominal tahunan (ptkp) dan kategori TER (group).
	ptkpStatus := ""
	if profile.PtkpStatus != nil {
		ptkpStatus = *profile.PtkpStatus
	}
	ptkpRow, err := s.repo.FindPtkpByCode(ctx, normalizePtkpCode(ptkpStatus))
	if err != nil {
		return nil, nil, err
	}
	ptkpAnnual := 0.0
	terCategory := terCategoryForPtkpStatus(ptkpStatus)
	if ptkpRow != nil {
		ptkpAnnual = float64(ptkpRow.Ptkp)
		if ptkpRow.Group != "" {
			terCategory = ptkpRow.Group
		}
	} else if ptkpStatus != "" {
		s.logger.Warn("PPh21: PTKP tidak ditemukan di tabel ptkps, nominal dipakai 0",
			zap.String("ptkp_status", ptkpStatus),
			zap.String("employee_id", ec.employee.ID.String()),
		)
	}

	rounding := calculator.NormalizeRoundingMode(setting.RoundingMode)
	months := setting.AnnualizationMonths
	if months <= 0 {
		months = 12
	}

	// ── Metode TER (PP 58/2023) — Januari s.d. November ──
	// PPh21 = bruto bulanan × tarif TER (tanpa biaya jabatan / BPJS / pensiun /
	// annualisasi / PTKP). Desember memakai metode normal (di bawah).
	if setting.CalculationMethod == "TER" && period.PeriodMonth < 12 {
		return s.calculatePph21TerMonthly(ctx, run, period, ec, *setting, *pph21Comp, gross, ptkpStatus, terCategory, rounding, profile)
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

	// 3b. Pengurang non-BPJS (mis. iuran pensiun): semua komponen bertanda
	// is_pph21_deductible dijumlahkan sebagai pengurang penghasilan bruto.
	pensionDeductible := 0.0
	for _, c := range components {
		if !c.IsPph21Deductible {
			continue
		}
		if v, ok := values[c.Code]; ok {
			pensionDeductible += v
		}
	}

	// 4. Net bulanan → annualisasi.
	netMonthly := gross - occ - bpjsDeductible - pensionDeductible
	if netMonthly < 0 {
		netMonthly = 0
	}
	netAnnualized := netMonthly * float64(months)

	// 5. PTKP tahunan sudah diambil dari tabel ptkps (lihat atas).

	// 6. PKP tahunan, dibulatkan ke bawah ke unit pembulatan.
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

	// ── TER bulan Desember: pajak setahun (metode normal) − potongan Jan–Nov ──
	method := setting.CalculationMethod
	if method == "" {
		method = "REGULAR_GROSS_ANNUALIZED"
	}
	var pph21Monthly float64
	ytdDeducted := 0.0
	if method == "TER" && period.PeriodMonth == 12 {
		ytd, err := s.repo.SumPph21YtdByEmployeeAndYear(ctx, ec.employee.ID, period.PeriodYear)
		if err != nil {
			return nil, nil, err
		}
		ytdDeducted = ytd
		pph21Monthly = annualTaxAfter - ytd
		if pph21Monthly < 0 {
			pph21Monthly = 0
		}
		pph21Monthly = calculator.Round(pph21Monthly, rounding)
	} else {
		pph21Monthly = calculator.Round(annualTaxAfter/float64(months), rounding)
	}

	breakdown := pph21FormulaBreakdown{
		CalculationMethod:          method,
		TerCategory:                terCategory,
		GrossMonthly:               gross,
		OccupationalExpenseMonthly: occ,
		BpjsTaxDeductibleMonthly:   bpjsDeductible,
		PensionDeductibleMonthly:   pensionDeductible,
		NetMonthly:                 netMonthly,
		AnnualizationMonths:        months,
		NetAnnualized:              netAnnualized,
		PtkpStatus:                 ptkpStatus,
		PtkpAnnual:                 ptkpAnnual,
		PkpAnnual:                  pkp,
		PkpRoundingUnit:            unit,
		AnnualTaxBeforeNpwpMult:    annualTax,
		NonNpwpMultiplierPercent:   mult,
		AnnualTaxAfterNpwpMult:     annualTaxAfter,
		Pph21Monthly:               pph21Monthly,
		AppliedNonNpwpMultiplier:   appliedMult,
		YtdDeductedMonthly:         ytdDeducted,
		Brackets:                   bracketBreakdowns,
	}
	formulaJSON, _ := json.Marshal(breakdown)
	formulaStr := string(formulaJSON)

	item := s.buildPph21RunItem(ec, *setting, *pph21Comp, pph21Monthly, gross, &formulaStr)
	log := &Pph21CalculationLog{
		PayrollRunID:              run.ID,
		EmployeeID:                ec.employee.ID,
		Pph21SettingID:            setting.ID,
		EmployeeTaxProfileID:      profile.ID,
		CalculationMethod:         method,
		TaxMethod:                 profile.TaxMethod,
		PtkpStatus:                ptkpStatus,
		HasNpwp:                   profile.HasNpwp,
		GrossMonthly:              gross,
		OccupationalExpenseMonthly: occ,
		BpjsTaxDeductibleMonthly:  bpjsDeductible,
		PensionDeductibleMonthly:  pensionDeductible,
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

// calculatePph21TerMonthly menghitung PPh21 bulanan metode TER untuk
// Januari–November: PPh21 = bruto bulanan × tarif TER (kategori dari status
// PTKP). Tanpa biaya jabatan, BPJS, iuran pensiun, annualisasi, maupun PTKP.
func (s *Service) calculatePph21TerMonthly(ctx context.Context, run *PayrollRun, period *PayrollPeriod, ec runEmployeeContext, setting Pph21Setting, pph21Comp SalaryComponent, gross float64, ptkpStatus, terCategory string, rounding calculator.RoundingMode, profile *EmployeeTaxProfile) (*PayrollRunItem, *Pph21CalculationLog, error) {
	if terCategory == "" {
		s.logger.Warn("PPh21 TER: status PTKP tidak dikenal untuk kategori TER, di-skip",
			zap.String("ptkp_status", ptkpStatus),
			zap.String("employee_id", ec.employee.ID.String()),
		)
		return nil, nil, nil
	}
	rate, err := s.repo.FindTerRateByGroupAndBruto(ctx, terCategory, int64(gross))
	if err != nil {
		return nil, nil, err
	}
	if rate == nil {
		s.logger.Warn("PPh21 TER: tarif tidak ditemukan untuk kategori/bruto, di-skip",
			zap.String("ter_category", terCategory),
			zap.Float64("gross", gross),
			zap.String("employee_id", ec.employee.ID.String()),
		)
		return nil, nil, nil
	}

	terRate := rate.Rate
	pph21Monthly := calculator.Round(gross*terRate/100, rounding)

	breakdown := pph21FormulaBreakdown{
		CalculationMethod:         "TER",
		TerCategory:               terCategory,
		TerRate:                   terRate,
		GrossMonthly:              gross,
		OccupationalExpenseMonthly: 0,
		BpjsTaxDeductibleMonthly:  0,
		PensionDeductibleMonthly:  0,
		NetMonthly:                gross,
		AnnualizationMonths:       0,
		NetAnnualized:             0,
		PtkpStatus:                ptkpStatus,
		PtkpAnnual:                0,
		PkpAnnual:                 0,
		PkpRoundingUnit:           0,
		AnnualTaxBeforeNpwpMult:   0,
		NonNpwpMultiplierPercent:  0,
		AnnualTaxAfterNpwpMult:    0,
		Pph21Monthly:              pph21Monthly,
		AppliedNonNpwpMultiplier:  0,
		Brackets:                  nil,
	}
	formulaJSON, _ := json.Marshal(breakdown)
	formulaStr := string(formulaJSON)

	item := s.buildPph21RunItem(ec, setting, pph21Comp, pph21Monthly, gross, &formulaStr)
	log := &Pph21CalculationLog{
		PayrollRunID:              run.ID,
		EmployeeID:                ec.employee.ID,
		Pph21SettingID:            setting.ID,
		EmployeeTaxProfileID:      profile.ID,
		CalculationMethod:         "TER",
		TaxMethod:                 profile.TaxMethod,
		PtkpStatus:                ptkpStatus,
		HasNpwp:                   profile.HasNpwp,
		GrossMonthly:              gross,
		OccupationalExpenseMonthly: 0,
		BpjsTaxDeductibleMonthly:  0,
		PensionDeductibleMonthly:  0,
		NetMonthly:                gross,
		NetAnnualized:             0,
		PtkpAnnual:                0,
		PkpAnnual:                 0,
		AnnualTaxBeforeNpwpMult:   0,
		NonNpwpMultiplierPercent:  0,
		AnnualTaxAfterNpwpMult:    0,
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
