package payroll

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/modules/payroll/calculator"
)

// BpjsProgram — program BPJS (sinkron dengan kolom enum bpjs_rate_components).
const (
	BpjsProgramHealth = "HEALTH"
	BpjsProgramJHT    = "JHT"
	BpjsProgramJP     = "JP"
	BpjsProgramJKK    = "JKK"
	BpjsProgramJKM    = "JKM"
	BpjsProgramJKP    = "JKP"
)

// calculateBpjsContributions menghitung kontribusi BPJS (employee deduction &
// employer contribution) untuk satu employee berdasarkan bpjs_settings +
// bpjs_rate_components yang ACTIVE dan berlaku pada tanggal periode, serta
// profil BPJS employee. Menghasilkan PayrollRunItem dengan source_group =
// STATUTORY. Tidak ada setting/rate/profil → tidak ada item (bukan error).
func (s *Service) calculateBpjsContributions(ctx context.Context, period *PayrollPeriod, ec runEmployeeContext, components []SalaryComponent, values map[string]float64) ([]PayrollRunItem, error) {
	asOf := period.AsOfDate

	setting, err := s.repo.FindActiveBpjsSettingByDate(ctx, asOf)
	if err != nil {
		return nil, err
	}
	if setting == nil {
		s.logger.Debug("BPJS: no active setting, skipping",
			zap.String("employee_id", ec.employee.ID.String()),
			zap.String("as_of", asOf),
		)
		return nil, nil
	}

	rates, err := s.repo.FindActiveBpjsRateComponentsBySettingID(ctx, setting.ID, asOf)
	if err != nil {
		return nil, err
	}
	if len(rates) == 0 {
		return nil, nil
	}

	profile, err := s.repo.FindActiveEmployeeBpjsProfileByEmployeeID(ctx, ec.employee.ID, asOf)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		// Employee tanpa profil BPJS aktif → tidak ada iuran yang dihitung.
		return nil, nil
	}

	// Dasar upah BPJS = total nilai komponen yang bertanda is_bpjs_base.
	wageBase := 0.0
	for _, c := range components {
		if !c.IsBpjsBase {
			continue
		}
		if v, ok := values[c.Code]; ok {
			wageBase += v
		}
	}

	// Lookup komponen master (id → komponen) untuk code/name/type item.
	compByID := map[uuid.UUID]SalaryComponent{}
	for _, c := range components {
		compByID[c.ID] = c
	}

	rounding := calculator.NormalizeRoundingMode(setting.RoundingMode)

	var items []PayrollRunItem
	for _, rate := range rates {
		if !rate.GenerateToPayrollItem {
			continue
		}
		// Rate JKK khusus risk class: hanya untuk employee dengan class yang sama.
		if rate.JkkRiskClass != nil && *rate.JkkRiskClass != profile.JkkRiskClass {
			continue
		}
		if rate.SalaryComponentID == nil {
			s.logger.Warn("BPJS: rate tanpa salary component di-skip",
				zap.String("rate_code", rate.RateCode),
				zap.String("employee_id", ec.employee.ID.String()),
			)
			continue
		}
		comp, ok := compByID[*rate.SalaryComponentID]
		if !ok {
			s.logger.Warn("BPJS: rate merujuk komponen non-aktif di-skip",
				zap.String("rate_code", rate.RateCode),
				zap.String("employee_id", ec.employee.ID.String()),
			)
			continue
		}

		base := wageBase
		// Cap program: kesehatan & pensiun punya batas maksimum di setting.
		switch rate.BpjsProgram {
		case BpjsProgramHealth:
			if setting.HealthMaxBaseAmount != nil && base > *setting.HealthMaxBaseAmount {
				base = *setting.HealthMaxBaseAmount
			}
		case BpjsProgramJP, BpjsProgramJHT:
			if setting.PensionMaxBaseAmount != nil && base > *setting.PensionMaxBaseAmount {
				base = *setting.PensionMaxBaseAmount
			}
		}
		// Cap per rate (min/max base amount).
		if rate.MinBaseAmount != nil && base < *rate.MinBaseAmount {
			base = *rate.MinBaseAmount
		}
		if rate.MaxBaseAmount != nil && base > *rate.MaxBaseAmount {
			base = *rate.MaxBaseAmount
		}

		var amount float64
		var calcType string
		var rateFraction *float64
		if rate.FixedAmount != nil {
			amount = *rate.FixedAmount
			calcType = "FIXED"
		} else {
			amount = base * rate.RatePercent / 100
			calcType = "PERCENTAGE"
			f := rate.RatePercent / 100
			rateFraction = &f
		}
		amount = calculator.Round(amount, rounding)

		items = append(items, s.buildBpjsRunItem(ec, rate, comp, amount, base, calcType, rateFraction))
	}
	return items, nil
}

// buildBpjsRunItem membangun snapshot item payroll untuk satu rate BPJS.
// paid_by EMPLOYER → EMPLOYER_CONTRIBUTION (biaya perusahaan); EMPLOYEE →
// EMPLOYEE_DEDUCTION (potongan gaji).
func (s *Service) buildBpjsRunItem(ec runEmployeeContext, rate BpjsRateComponent, comp SalaryComponent, amount, base float64, calcType string, rateFraction *float64) PayrollRunItem {
	sourceTable := "bpjs_rate_components"
	sourceType := rate.BpjsProgram
	notes := rate.RateName
	item := PayrollRunItem{
		EmployeeID:        ec.employee.ID,
		SalaryComponentID: comp.ID,
		ComponentCode:     comp.Code,
		ComponentName:     comp.Name,
		ComponentType:     comp.ComponentType,
		CalculationType:   calcType,
		Amount:            amount,
		BaseAmount:        base,
		Rate:              rateFraction,
		CurrencyCode:      "IDR",
		SourceGroup:       SourceGroupStatutory,
		SourceTable:       &sourceTable,
		SourceID:          &rate.ID,
		SourceType:        &sourceType,
		PrintOnPayslip:    rate.PrintOnPayslip,
		Notes:             &notes,
	}

	switch rate.PaidBy {
	case "EMPLOYER":
		item.ItemCategory = ItemCategoryEmployerContribution
		item.PaidBy = "EMPLOYER"
		item.AffectsCompanyCost = true
	default:
		item.ItemCategory = ItemCategoryEmployeeDeduction
		item.PaidBy = "EMPLOYEE"
		item.AffectsNetPay = true
	}
	return item
}
