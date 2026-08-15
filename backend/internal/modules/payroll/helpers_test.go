package payroll

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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

	// AutoMigrate all payroll models + read models dari modul lain
	// (employees/employments/positions/gradings) yang dipakai kalkulasi run.
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
		&Ptkp{},
		&Pph21TaxBracket{},
		&TerRate{},
		&PayrollRun{},
		&PayrollRunEmployee{},
		&PayrollRunItem{},
		&PayrollPayslip{},
		&PayrollPayment{},
		&Pph21CalculationLog{},
		&PayrollProfileChangeLog{},
		&EmployeeRead{},
		&EmploymentRead{},
		&PositionRead{},
		&GradingRead{},
		&AttendanceSessionRead{},
		&LeaveRequestRead{},
		&LeaveRequestDetailRead{},
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

// =============================================================================
// Helpers data master untuk test kalkulasi payroll run
// =============================================================================

// createTestGrading inserts a grading (read model) and returns it.
func createTestGrading(ctx context.Context, repo *Repository, code, name string) *GradingRead {
	g := &GradingRead{ID: uuid.New(), Code: code, Name: name}
	db, _ := repo.getDB(ctx)
	if err := db.Create(g).Error; err != nil {
		panic(fmt.Sprintf("failed to create test grading: %v", err))
	}
	return g
}

// createTestPosition inserts a position (read model) linked to a grading.
func createTestPosition(ctx context.Context, repo *Repository, title string, gradingID *uuid.UUID) *PositionRead {
	p := &PositionRead{ID: uuid.New(), Title: title, GradingID: gradingID}
	db, _ := repo.getDB(ctx)
	if err := db.Create(p).Error; err != nil {
		panic(fmt.Sprintf("failed to create test position: %v", err))
	}
	return p
}

// createTestEmployee inserts an employee (read model) and returns it.
func createTestEmployee(ctx context.Context, repo *Repository, code, name string) *EmployeeRead {
	e := &EmployeeRead{ID: uuid.New(), EmployeeID: code, Name: name, Status: "active"}
	db, _ := repo.getDB(ctx)
	if err := db.Create(e).Error; err != nil {
		panic(fmt.Sprintf("failed to create test employee: %v", err))
	}
	return e
}

// createTestEmployment inserts an employment (read model) for an employee.
func createTestEmployment(ctx context.Context, repo *Repository, employeeID, positionID *uuid.UUID, effectiveDate string) *EmploymentRead {
	e := &EmploymentRead{
		ID:            uuid.New(),
		EmployeeID:    employeeID,
		PositionID:    positionID,
		EffectiveDate: effectiveDate,
	}
	db, _ := repo.getDB(ctx)
	if err := db.Create(e).Error; err != nil {
		panic(fmt.Sprintf("failed to create test employment: %v", err))
	}
	return e
}

// createTestPayrollProfile inserts an employee payroll profile covering a period.
func createTestPayrollProfile(ctx context.Context, repo *Repository, employeeID uuid.UUID, startDate string) *EmployeePayrollProfile {
	p := &EmployeePayrollProfile{
		EmployeeID:         employeeID,
		PayrollGroupCode:   "MONTHLY",
		PayrollFrequency:   "MONTHLY",
		PaymentMethod:      "BANK_TRANSFER",
		SalaryCurrency:     "IDR",
		IsPayrollActive:    true,
		EffectiveStartDate: startDate,
		Status:             "ACTIVE",
	}
	if err := repo.CreateEmployeePayrollProfile(ctx, p); err != nil {
		panic(fmt.Sprintf("failed to create test payroll profile: %v", err))
	}
	return p
}

// createTestGradeComponent inserts a grade-level default component amount.
func createTestGradeComponent(ctx context.Context, repo *Repository, gradingID, componentID uuid.UUID, amount float64) *SalaryGradeComponent {
	gc := &SalaryGradeComponent{
		GradingID:          &gradingID,
		SalaryComponentID:  componentID,
		Amount:             amount,
		EffectiveStartDate: "2020-01-01",
		IsMandatory:        true,
		IsDefault:          true,
		Status:             "ACTIVE",
	}
	if err := repo.CreateSalaryGradeComponent(ctx, gc); err != nil {
		panic(fmt.Sprintf("failed to create test grade component: %v", err))
	}
	return gc
}

// createTestEmployeeComponent inserts an employee-level component amount.
func createTestEmployeeComponent(ctx context.Context, repo *Repository, employeeID, componentID uuid.UUID, amount float64) *SalaryEmployeeComponent {
	ec := &SalaryEmployeeComponent{
		EmployeeID:         employeeID,
		SalaryComponentID:  componentID,
		Amount:             amount,
		EffectiveStartDate: "2020-01-01",
		Status:             "ACTIVE",
	}
	if err := repo.CreateSalaryEmployeeComponent(ctx, ec); err != nil {
		panic(fmt.Sprintf("failed to create test employee component: %v", err))
	}
	return ec
}

// createTestAdjustment inserts a one-time adjustment for a period.
func createTestAdjustment(ctx context.Context, repo *Repository, employeeID, componentID uuid.UUID, year, month int, amount float64) *SalaryEmployeeAdjustment {
	a := &SalaryEmployeeAdjustment{
		EmployeeID:        employeeID,
		SalaryComponentID: componentID,
		PeriodYear:        year,
		PeriodMonth:       month,
		Amount:            amount,
		Status:            "APPROVED",
	}
	if err := repo.CreateSalaryEmployeeAdjustment(ctx, a); err != nil {
		panic(fmt.Sprintf("failed to create test adjustment: %v", err))
	}
	return a
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

// createTestBpjsRateComponentLinked inserts a BPJS rate component yang terhubung
// ke salary component (dibutuhkan kalkulator BPJS agar item punya code/name).
// rateCode harus unik (kolom rate_code ber-index unique).
func createTestBpjsRateComponentLinked(ctx context.Context, repo *Repository, setting *BpjsSetting, comp *SalaryComponent, program, paidBy, rateCode string, ratePercent float64) *BpjsRateComponent {
	br := &BpjsRateComponent{
		BpjsSettingID:         setting.ID,
		RateCode:              rateCode,
		RateName:              "BPJS " + program + " (" + paidBy + ")",
		BpjsProgram:           program,
		PaidBy:                paidBy,
		SalaryComponentID:     &comp.ID,
		RatePercent:           ratePercent,
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

// createTestEmployeeBpjsProfile inserts profil BPJS aktif untuk employee.
func createTestEmployeeBpjsProfile(ctx context.Context, repo *Repository, employeeID, payrollProfileID uuid.UUID) *EmployeeBpjsProfile {
	p := &EmployeeBpjsProfile{
		EmployeeID:               employeeID,
		EmployeePayrollProfileID: payrollProfileID,
		BpjsHealthActive:         true,
		BpjsTkActive:             true,
		JkkRiskClass:             "LOW",
		PensionActive:            true,
		EffectiveStartDate:       "2026-01-01",
		Status:                   "ACTIVE",
	}
	if err := repo.CreateEmployeeBpjsProfile(ctx, p); err != nil {
		panic(fmt.Sprintf("failed to create test BPJS profile: %v", err))
	}
	return p
}

// createTestBankProfile inserts bank profile utama aktif untuk employee.
func createTestBankProfile(ctx context.Context, repo *Repository, employeeID, payrollProfileID uuid.UUID, accountNumber, holderName string) *EmployeeBankProfile {
	b := &EmployeeBankProfile{
		EmployeeID:             employeeID,
		EmployeePayrollProfileID: payrollProfileID,
		BankName:               "Bank Maju",
		BankAccountNumber:      accountNumber,
		BankAccountHolderName:  holderName,
		IsPrimary:              true,
		EffectiveStartDate:     "2020-01-01",
		Status:                 "ACTIVE",
	}
	if err := repo.CreateEmployeeBankProfile(ctx, b); err != nil {
		panic(fmt.Sprintf("failed to create test bank profile: %v", err))
	}
	return b
}

// createTestAttendanceSession inserts session kehadiran (read model) untuk
// summary workforce. status kosong → CLOSED (hadir).
func createTestAttendanceSession(ctx context.Context, repo *Repository, employeeID uuid.UUID, workDate, status string, overtimeMinutes int) *AttendanceSessionRead {
	sess := &AttendanceSessionRead{
		ID:              uuid.New(),
		EmployeeID:      employeeID,
		WorkDate:        workDate,
		Status:          status,
		OvertimeMinutes: overtimeMinutes,
	}
	if sess.Status == "" {
		sess.Status = "CLOSED"
	}
	db, _ := repo.getDB(ctx)
	if err := db.Create(sess).Error; err != nil {
		panic(fmt.Sprintf("failed to create test attendance session: %v", err))
	}
	return sess
}

// createTestLeaveRequest inserts leave request (read model) berstatus tertentu.
func createTestLeaveRequest(ctx context.Context, repo *Repository, employeeID uuid.UUID, status string) *LeaveRequestRead {
	lr := &LeaveRequestRead{ID: uuid.New(), EmployeeID: employeeID, Status: status}
	db, _ := repo.getDB(ctx)
	if err := db.Create(lr).Error; err != nil {
		panic(fmt.Sprintf("failed to create test leave request: %v", err))
	}
	return lr
}

// createTestLeaveRequestDetail inserts detail cuti (hari, is_paid, fraksi).
func createTestLeaveRequestDetail(ctx context.Context, repo *Repository, leaveRequestID, employeeID uuid.UUID, leaveDate string, isPaid bool, dayFraction float64) *LeaveRequestDetailRead {
	d := &LeaveRequestDetailRead{
		LeaveRequestID: leaveRequestID,
		EmployeeID:     employeeID,
		LeaveDate:      leaveDate,
		IsPaid:         isPaid,
		DayFraction:    dayFraction,
	}
	db, _ := repo.getDB(ctx)
	if err := db.Create(d).Error; err != nil {
		panic(fmt.Sprintf("failed to create test leave detail: %v", err))
	}
	return d
}

// createTestEmployeeTaxProfile inserts profil pajak aktif untuk employee.
func createTestEmployeeTaxProfile(ctx context.Context, repo *Repository, employeeID, payrollProfileID uuid.UUID, ptkpStatus string, hasNpwp bool) *EmployeeTaxProfile {
	p := &EmployeeTaxProfile{
		EmployeeID:               employeeID,
		EmployeePayrollProfileID: payrollProfileID,
		PtkpStatus:               &ptkpStatus,
		TaxMethod:                "GROSS",
		IsTaxable:                true,
		HasNpwp:                  hasNpwp,
		EffectiveStartDate:       "2026-01-01",
		Status:                   "ACTIVE",
	}
	if err := repo.CreateEmployeeTaxProfile(ctx, p); err != nil {
		panic(fmt.Sprintf("failed to create test tax profile: %v", err))
	}
	return p
}

// createTestPph21SettingCustom inserts PPh21 setting dengan komponen potongan
// pajak yang ditandai flag IsPph21Component (sumber kebenaran di komponen gaji).
func createTestPph21SettingCustom(ctx context.Context, repo *Repository, pph21Comp *SalaryComponent) *Pph21Setting {
	return createTestPph21SettingCustomMethod(ctx, repo, pph21Comp, "REGULAR_GROSS_ANNUALIZED")
}

// createTestPph21SettingCustomMethod sama seperti createTestPph21SettingCustom
// tapi metode kalkulasi bisa dipilih (TER / REGULAR_GROSS_ANNUALIZED).
func createTestPph21SettingCustomMethod(ctx context.Context, repo *Repository, pph21Comp *SalaryComponent, method string) *Pph21Setting {
	pph21Comp.IsPph21Component = true
	if err := repo.UpdateSalaryComponent(ctx, pph21Comp); err != nil {
		panic(fmt.Sprintf("failed to mark test PPh21 component: %v", err))
	}
	ps := &Pph21Setting{
		SettingCode:                    "PPH21-DEFAULT",
		SettingName:                    "Default PPh21 Setting",
		CalculationMethod:              method,
		DefaultTaxMethod:               "GROSS",
		OccupationalExpenseRatePercent: 5.0,
		OccupationalExpenseMaxMonthly:  500000,
		OccupationalExpenseMaxYearly:   6000000,
		DeductBpjsHealthEmployee:       false,
		DeductBpjsJhtEmployee:          true,
		DeductBpjsJpEmployee:           true,
		AnnualizationMonths:            12,
		PkpRoundingUnit:                1000,
		NonNpwpMultiplierPercent:       100,
		RoundingMode:                   "ROUND",
		EffectiveStartDate:             "2026-01-01",
		Status:                         "ACTIVE",
	}
	if err := repo.CreatePph21Setting(ctx, ps); err != nil {
		panic(fmt.Sprintf("failed to create test PPh21 setting: %v", err))
	}
	return ps
}

// createTestTerRate inserts satu baris tarif TER (tabel ters).
func createTestTerRate(ctx context.Context, repo *Repository, group string, brutoMin, brutoMax *int64, rate float64) {
	db, err := repo.getDB(ctx)
	if err != nil {
		panic(err)
	}
	tr := &TerRate{ID: uuid.New(), Group: group, BrutoMin: brutoMin, BrutoMax: brutoMax, Rate: rate}
	if err := db.Create(tr).Error; err != nil {
		panic(fmt.Sprintf("failed to create test TER rate: %v", err))
	}
}

// createTestPayrollPeriodCustom membuat periode payroll untuk bulan & tahun tertentu.
func createTestPayrollPeriodCustom(ctx context.Context, repo *Repository, year, month int) *PayrollPeriod {
	p := &PayrollPeriod{
		PeriodCode:  fmt.Sprintf("%04d%02d", year, month),
		PeriodYear:  year,
		PeriodMonth: month,
		StartDate:   fmt.Sprintf("%04d-%02d-01", year, month),
		EndDate:     fmt.Sprintf("%04d-%02d-28", year, month),
		AsOfDate:    fmt.Sprintf("%04d-%02d-28", year, month),
		Status:      "OPEN",
	}
	if err := repo.CreatePayrollPeriod(ctx, p); err != nil {
		panic(fmt.Sprintf("failed to create test payroll period: %v", err))
	}
	return p
}

// createTestPtkp inserts baris PTKP ke tabel ptkps (satu sumber kebenaran).
func createTestPtkp(ctx context.Context, repo *Repository, code, name string, annual int64, group string) {
	db, err := repo.getDB(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to get db: %v", err))
	}
	if err := db.Table("ptkps").Create(map[string]interface{}{
		"id": uuid.NewString(), "code": code, "name": name, "ptkp": annual, "group": group,
	}).Error; err != nil {
		panic(fmt.Sprintf("failed to create test PTKP: %v", err))
	}
}

// createTestPph21TaxBracket inserts tax bracket progresif aktif. upper nil =
// bracket terbuka (sampai tak hingga).
func createTestPph21TaxBracket(ctx context.Context, repo *Repository, order int, lower float64, upper *float64, ratePercent float64) *Pph21TaxBracket {
	tb := &Pph21TaxBracket{
		BracketOrder:       order,
		LowerBound:         lower,
		UpperBound:         upper,
		RatePercent:        ratePercent,
		EffectiveStartDate: "2020-01-01",
		Status:             "ACTIVE",
	}
	if err := repo.CreatePph21TaxBracket(ctx, tb); err != nil {
		panic(fmt.Sprintf("failed to create test tax bracket: %v", err))
	}
	return tb
}

// createTestPph21Setting inserts a test PPh21 setting and returns it.
func createTestPph21Setting(ctx context.Context, repo *Repository) *Pph21Setting {
	ps := &Pph21Setting{
		SettingCode:                    "PPH21-DEFAULT",
		SettingName:                    "Default PPh21 Setting",
		CalculationMethod:              "REGULAR_GROSS_ANNUALIZED",
		DefaultTaxMethod:               "GROSS",
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
