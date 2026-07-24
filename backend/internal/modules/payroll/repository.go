package payroll

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	dbResolver func(ctx context.Context) (*gorm.DB, error)
}

func NewRepository(dbResolver func(ctx context.Context) (*gorm.DB, error)) *Repository {
	return &Repository{dbResolver: dbResolver}
}

func (r *Repository) getDB(ctx context.Context) (*gorm.DB, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required for tenant database resolution")
	}
	return r.dbResolver(ctx)
}

// =============================================================================
// Salary Components
// =============================================================================

func (r *Repository) CreateSalaryComponent(ctx context.Context, sc *SalaryComponent) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(sc).Error
}

func (r *Repository) FindSalaryComponentByID(ctx context.Context, id uuid.UUID) (*SalaryComponent, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var sc SalaryComponent
	if err := db.First(&sc, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("salary component not found: %w", err)
	}
	return &sc, nil
}

func (r *Repository) FindAllSalaryComponents(ctx context.Context, page, perPage int) ([]SalaryComponent, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var items []SalaryComponent
	var total int64
	query := db.Model(&SalaryComponent{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("display_order ASC, code ASC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) UpdateSalaryComponent(ctx context.Context, sc *SalaryComponent) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(sc).Error
}

func (r *Repository) DeleteSalaryComponent(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&SalaryComponent{}).Error
}

// =============================================================================
// Payroll Periods
// =============================================================================

func (r *Repository) CreatePayrollPeriod(ctx context.Context, p *PayrollPeriod) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(p).Error
}

func (r *Repository) FindPayrollPeriodByID(ctx context.Context, id uuid.UUID) (*PayrollPeriod, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var p PayrollPeriod
	if err := db.First(&p, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("payroll period not found: %w", err)
	}
	return &p, nil
}

func (r *Repository) FindAllPayrollPeriods(ctx context.Context, page, perPage int) ([]PayrollPeriod, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var items []PayrollPeriod
	var total int64
	query := db.Model(&PayrollPeriod{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("period_year DESC, period_month DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) UpdatePayrollPeriod(ctx context.Context, p *PayrollPeriod) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(p).Error
}

// =============================================================================
// Employee Payroll Profiles
// =============================================================================

func (r *Repository) CreateEmployeePayrollProfile(ctx context.Context, p *EmployeePayrollProfile) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(p).Error
}

func (r *Repository) FindEmployeePayrollProfileByID(ctx context.Context, id uuid.UUID) (*EmployeePayrollProfile, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var p EmployeePayrollProfile
	if err := db.First(&p, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("employee payroll profile not found: %w", err)
	}
	return &p, nil
}

func (r *Repository) FindEmployeePayrollProfilesByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]EmployeePayrollProfile, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []EmployeePayrollProfile
	if err := db.Where("employee_id = ?", employeeID).Order("effective_start_date DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) FindAllEmployeePayrollProfiles(ctx context.Context, page, perPage int) ([]EmployeePayrollProfile, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var items []EmployeePayrollProfile
	var total int64
	query := db.Model(&EmployeePayrollProfile{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) UpdateEmployeePayrollProfile(ctx context.Context, p *EmployeePayrollProfile) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(p).Error
}

func (r *Repository) DeleteEmployeePayrollProfile(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmployeePayrollProfile{}).Error
}

// =============================================================================
// Employee Bank Profiles
// =============================================================================

func (r *Repository) CreateEmployeeBankProfile(ctx context.Context, b *EmployeeBankProfile) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(b).Error
}

func (r *Repository) FindEmployeeBankProfileByID(ctx context.Context, id uuid.UUID) (*EmployeeBankProfile, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var b EmployeeBankProfile
	if err := db.First(&b, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("employee bank profile not found: %w", err)
	}
	return &b, nil
}

func (r *Repository) FindEmployeeBankProfilesByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]EmployeeBankProfile, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []EmployeeBankProfile
	if err := db.Where("employee_id = ?", employeeID).Order("is_primary DESC, effective_start_date DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) UpdateEmployeeBankProfile(ctx context.Context, b *EmployeeBankProfile) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(b).Error
}

func (r *Repository) DeleteEmployeeBankProfile(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmployeeBankProfile{}).Error
}

// =============================================================================
// Employee BPJS Profiles
// =============================================================================

func (r *Repository) CreateEmployeeBpjsProfile(ctx context.Context, b *EmployeeBpjsProfile) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(b).Error
}

func (r *Repository) FindEmployeeBpjsProfileByID(ctx context.Context, id uuid.UUID) (*EmployeeBpjsProfile, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var b EmployeeBpjsProfile
	if err := db.First(&b, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("employee BPJS profile not found: %w", err)
	}
	return &b, nil
}

func (r *Repository) FindEmployeeBpjsProfilesByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]EmployeeBpjsProfile, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []EmployeeBpjsProfile
	if err := db.Where("employee_id = ?", employeeID).Order("effective_start_date DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) UpdateEmployeeBpjsProfile(ctx context.Context, b *EmployeeBpjsProfile) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(b).Error
}

func (r *Repository) DeleteEmployeeBpjsProfile(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmployeeBpjsProfile{}).Error
}

// =============================================================================
// Employee Tax Profiles
// =============================================================================

func (r *Repository) CreateEmployeeTaxProfile(ctx context.Context, t *EmployeeTaxProfile) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(t).Error
}

func (r *Repository) FindEmployeeTaxProfileByID(ctx context.Context, id uuid.UUID) (*EmployeeTaxProfile, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var t EmployeeTaxProfile
	if err := db.First(&t, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("employee tax profile not found: %w", err)
	}
	return &t, nil
}

func (r *Repository) FindEmployeeTaxProfilesByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]EmployeeTaxProfile, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []EmployeeTaxProfile
	if err := db.Where("employee_id = ?", employeeID).Order("effective_start_date DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) UpdateEmployeeTaxProfile(ctx context.Context, t *EmployeeTaxProfile) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(t).Error
}

func (r *Repository) DeleteEmployeeTaxProfile(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmployeeTaxProfile{}).Error
}

// =============================================================================
// BPJS Settings
// =============================================================================

func (r *Repository) CreateBpjsSetting(ctx context.Context, bs *BpjsSetting) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(bs).Error
}

func (r *Repository) FindBpjsSettingByID(ctx context.Context, id uuid.UUID) (*BpjsSetting, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var bs BpjsSetting
	if err := db.First(&bs, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("BPJS setting not found: %w", err)
	}
	return &bs, nil
}

func (r *Repository) FindAllBpjsSettings(ctx context.Context, page, perPage int) ([]BpjsSetting, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var items []BpjsSetting
	var total int64
	query := db.Model(&BpjsSetting{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("setting_code ASC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) UpdateBpjsSetting(ctx context.Context, bs *BpjsSetting) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(bs).Error
}

func (r *Repository) DeleteBpjsSetting(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&BpjsSetting{}).Error
}

// =============================================================================
// BPJS Rate Components
// =============================================================================

func (r *Repository) CreateBpjsRateComponent(ctx context.Context, br *BpjsRateComponent) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(br).Error
}

func (r *Repository) FindBpjsRateComponentByID(ctx context.Context, id uuid.UUID) (*BpjsRateComponent, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var br BpjsRateComponent
	if err := db.First(&br, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("BPJS rate component not found: %w", err)
	}
	return &br, nil
}

func (r *Repository) FindBpjsRateComponentsBySettingID(ctx context.Context, settingID uuid.UUID) ([]BpjsRateComponent, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []BpjsRateComponent
	if err := db.Where("bpjs_setting_id = ?", settingID).Order("display_order ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) UpdateBpjsRateComponent(ctx context.Context, br *BpjsRateComponent) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(br).Error
}

func (r *Repository) DeleteBpjsRateComponent(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&BpjsRateComponent{}).Error
}

// =============================================================================
// PPh21 Settings
// =============================================================================

func (r *Repository) CreatePph21Setting(ctx context.Context, ps *Pph21Setting) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(ps).Error
}

func (r *Repository) FindPph21SettingByID(ctx context.Context, id uuid.UUID) (*Pph21Setting, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var ps Pph21Setting
	if err := db.First(&ps, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("PPh21 setting not found: %w", err)
	}
	return &ps, nil
}

func (r *Repository) FindAllPph21Settings(ctx context.Context, page, perPage int) ([]Pph21Setting, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var items []Pph21Setting
	var total int64
	query := db.Model(&Pph21Setting{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("setting_code ASC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) UpdatePph21Setting(ctx context.Context, ps *Pph21Setting) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(ps).Error
}

func (r *Repository) DeletePph21Setting(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&Pph21Setting{}).Error
}

// =============================================================================
// PPh21 PTKP Rates
// =============================================================================

func (r *Repository) CreatePph21PtkpRate(ctx context.Context, pr *Pph21PtkpRate) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(pr).Error
}

func (r *Repository) FindPph21PtkpRateByID(ctx context.Context, id uuid.UUID) (*Pph21PtkpRate, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var pr Pph21PtkpRate
	if err := db.First(&pr, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("PPh21 PTKP rate not found: %w", err)
	}
	return &pr, nil
}

func (r *Repository) FindAllPph21PtkpRates(ctx context.Context, page, perPage int) ([]Pph21PtkpRate, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var items []Pph21PtkpRate
	var total int64
	query := db.Model(&Pph21PtkpRate{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("ptkp_status ASC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) UpdatePph21PtkpRate(ctx context.Context, pr *Pph21PtkpRate) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(pr).Error
}

func (r *Repository) DeletePph21PtkpRate(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&Pph21PtkpRate{}).Error
}

// =============================================================================
// PPh21 Tax Brackets
// =============================================================================

func (r *Repository) CreatePph21TaxBracket(ctx context.Context, tb *Pph21TaxBracket) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(tb).Error
}

func (r *Repository) FindPph21TaxBracketByID(ctx context.Context, id uuid.UUID) (*Pph21TaxBracket, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var tb Pph21TaxBracket
	if err := db.First(&tb, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("PPh21 tax bracket not found: %w", err)
	}
	return &tb, nil
}

func (r *Repository) FindAllPph21TaxBrackets(ctx context.Context, page, perPage int) ([]Pph21TaxBracket, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var items []Pph21TaxBracket
	var total int64
	query := db.Model(&Pph21TaxBracket{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("bracket_order ASC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) UpdatePph21TaxBracket(ctx context.Context, tb *Pph21TaxBracket) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(tb).Error
}

func (r *Repository) DeletePph21TaxBracket(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&Pph21TaxBracket{}).Error
}

// =============================================================================
// Payroll Runs
// =============================================================================

func (r *Repository) CreatePayrollRun(ctx context.Context, pr *PayrollRun) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(pr).Error
}

func (r *Repository) FindPayrollRunByID(ctx context.Context, id uuid.UUID) (*PayrollRun, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var pr PayrollRun
	if err := db.First(&pr, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("payroll run not found: %w", err)
	}
	return &pr, nil
}

func (r *Repository) FindAllPayrollRuns(ctx context.Context, page, perPage int) ([]PayrollRun, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var items []PayrollRun
	var total int64
	query := db.Model(&PayrollRun{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) UpdatePayrollRun(ctx context.Context, pr *PayrollRun) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(pr).Error
}

// =============================================================================
// Payroll Run Employees
// =============================================================================

func (r *Repository) CreatePayrollRunEmployee(ctx context.Context, pre *PayrollRunEmployee) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(pre).Error
}

func (r *Repository) BulkCreatePayrollRunEmployees(ctx context.Context, employees []PayrollRunEmployee) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.CreateInBatches(employees, 50).Error
}

func (r *Repository) FindPayrollRunEmployeesByRunID(ctx context.Context, runID uuid.UUID) ([]PayrollRunEmployee, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []PayrollRunEmployee
	if err := db.Where("payroll_run_id = ?", runID).Order("employee_name ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// =============================================================================
// Payroll Run Items
// =============================================================================

func (r *Repository) BulkCreatePayrollRunItems(ctx context.Context, items []PayrollRunItem) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.CreateInBatches(items, 100).Error
}

func (r *Repository) FindPayrollRunItemsByRunID(ctx context.Context, runID uuid.UUID) ([]PayrollRunItem, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []PayrollRunItem
	if err := db.Where("payroll_run_id = ?", runID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) FindPayrollRunItemsByEmployeeID(ctx context.Context, runID, employeeID uuid.UUID) ([]PayrollRunItem, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []PayrollRunItem
	if err := db.Where("payroll_run_id = ? AND employee_id = ?", runID, employeeID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// =============================================================================
// Payroll Payslips
// =============================================================================

func (r *Repository) CreatePayrollPayslip(ctx context.Context, p *PayrollPayslip) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(p).Error
}

func (r *Repository) FindPayrollPayslipByID(ctx context.Context, id uuid.UUID) (*PayrollPayslip, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var p PayrollPayslip
	if err := db.First(&p, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("payslip not found: %w", err)
	}
	return &p, nil
}

func (r *Repository) FindPayrollPayslipsByRunID(ctx context.Context, runID uuid.UUID) ([]PayrollPayslip, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []PayrollPayslip
	if err := db.Where("payroll_run_id = ?", runID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) UpdatePayrollPayslip(ctx context.Context, p *PayrollPayslip) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(p).Error
}
