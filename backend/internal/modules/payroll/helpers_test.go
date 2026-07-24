package payroll

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

// setupTestDB creates an in-memory SQLite database and auto-migrates all payroll models.
// Returns the GORM DB, a dbResolver function, and a cleanup function.
func setupTestDB() (*gorm.DB, func(ctx context.Context) (*gorm.DB, error), func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	// AutoMigrate all payroll models
	if err := db.AutoMigrate(
		&SalaryComponent{},
		&SalaryGradeComponent{},
		&SalaryEmployeeComponent{},
		&SalaryChangeLog{},
		&SalaryEmployeeAdjustment{},
		&PayrollPeriod{},
		&EmployeePayrollProfile{},
		&EmployeeBankProfile{},
		&EmployeeBpjsProfile{},
		&EmployeeTaxProfile{},
		&BpjsSetting{},
		&BpjsRateComponent{},
		&Pph21Setting{},
		&Pph21PtkpRate{},
		&Pph21TaxBracket{},
		&PayrollRun{},
		&PayrollRunEmployee{},
		&PayrollRunItem{},
		&PayrollPayslip{},
		&Pph21CalculationLog{},
		&PayrollProfileChangeLog{},
	); err != nil {
		panic(fmt.Sprintf("failed to migrate test db: %v", err))
	}

	dbResolver := func(ctx context.Context) (*gorm.DB, error) {
		return db, nil
	}

	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	return db, dbResolver, cleanup
}

// createTestSalaryComponent inserts a test salary component and returns it.
func createTestSalaryComponent(ctx context.Context, repo *Repository) *SalaryComponent {
	sc := &SalaryComponent{
		Code:                   "BASIC",
		Name:                   "Basic Salary",
		ComponentType:          "EARNING",
		CalculationType:        "FIXED",
		IsTaxable:              true,
		IsBpjsBase:             true,
		IsRecurring:            true,
		IsProratable:           true,
		PrintOnSalaryStructure: true,
		DisplayOrder:           10,
		Status:                 "ACTIVE",
	}
	if err := repo.CreateSalaryComponent(ctx, sc); err != nil {
		panic(fmt.Sprintf("failed to create test salary component: %v", err))
	}
	return sc
}

// createTestPayrollPeriod inserts a test payroll period and returns it.
func createTestPayrollPeriod(ctx context.Context, repo *Repository) *PayrollPeriod {
	p := &PayrollPeriod{
		PeriodCode:  "202601",
		PeriodYear:  2026,
		PeriodMonth: 1,
		StartDate:   "2026-01-01",
		EndDate:     "2026-01-31",
		AsOfDate:    "2026-01-31",
		Status:      "OPEN",
	}
	if err := repo.CreatePayrollPeriod(ctx, p); err != nil {
		panic(fmt.Sprintf("failed to create test payroll period: %v", err))
	}
	return p
}

// createTestBpjsSetting inserts a test BPJS setting and returns it.
func createTestBpjsSetting(ctx context.Context, repo *Repository) *BpjsSetting {
	bs := &BpjsSetting{
		SettingCode:         "BPJS-DEFAULT",
		SettingName:         "Default BPJS Setting",
		BaseSource:          "BPJS_BASE_COMPONENTS",
		DefaultJkkRiskClass: "LOW",
		RoundingMode:        "ROUND",
		EffectiveStartDate:  "2026-01-01",
		Status:              "ACTIVE",
	}
	if err := repo.CreateBpjsSetting(ctx, bs); err != nil {
		panic(fmt.Sprintf("failed to create test BPJS setting: %v", err))
	}
	return bs
}

// createTestPayrollRun inserts a test payroll run and returns it.
func createTestPayrollRun(ctx context.Context, repo *Repository) *PayrollRun {
	period := createTestPayrollPeriod(ctx, repo)
	pr := &PayrollRun{
		PayrollPeriodID: period.ID,
		RunCode:         "RUN-2026-01",
		RunType:         "REGULAR",
		Status:          "DRAFT",
	}
	if err := repo.CreatePayrollRun(ctx, pr); err != nil {
		panic(fmt.Sprintf("failed to create test payroll run: %v", err))
	}
	return pr
}

// createTestBpjsRateComponent inserts a test BPJS rate component and returns it.
func createTestBpjsRateComponent(ctx context.Context, repo *Repository) *BpjsRateComponent {
	setting := createTestBpjsSetting(ctx, repo)
	br := &BpjsRateComponent{
		BpjsSettingID:      setting.ID,
		RateCode:           "BPJS-HEALTH-EMP",
		RateName:           "BPJS Kesehatan - Employee",
		BpjsProgram:        "HEALTH",
		PaidBy:             "EMPLOYEE",
		RatePercent:        1.0,
		GenerateToPayrollItem: true,
		PrintOnPayslip:        true,
		DisplayOrder:          1,
		EffectiveStartDate:    "2026-01-01",
		Status:                "ACTIVE",
	}
	if err := repo.CreateBpjsRateComponent(ctx, br); err != nil {
		panic(fmt.Sprintf("failed to create test BPJS rate component: %v", err))
	}
	return br
}

// createTestPph21Setting inserts a test PPh21 setting and returns it.
func createTestPph21Setting(ctx context.Context, repo *Repository) *Pph21Setting {
	comp := createTestSalaryComponent(ctx, repo)
	ps := &Pph21Setting{
		SettingCode:                    "PPH21-DEFAULT",
		SettingName:                    "Default PPh21 Setting",
		CalculationMethod:              "REGULAR_GROSS_ANNUALIZED",
		DefaultTaxMethod:               "GROSS",
		Pph21ComponentID:               comp.ID,
		OccupationalExpenseRatePercent: 5.0,
		OccupationalExpenseMaxMonthly:  500000,
		OccupationalExpenseMaxYearly:   6000000,
		AnnualizationMonths:            12,
		PkpRoundingUnit:                1000,
		RoundingMode:                   "ROUND",
		EffectiveStartDate:             "2026-01-01",
		Status:                         "ACTIVE",
	}
	if err := repo.CreatePph21Setting(ctx, ps); err != nil {
		panic(fmt.Sprintf("failed to create test PPh21 setting: %v", err))
	}
	return ps
}
