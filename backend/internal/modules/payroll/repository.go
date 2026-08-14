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
// Salary Grade Components & Salary Employee Components
// =============================================================================

func (r *Repository) CreateSalaryGradeComponent(ctx context.Context, gc *SalaryGradeComponent) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(gc).Error
}

func (r *Repository) UpdateSalaryGradeComponent(ctx context.Context, gc *SalaryGradeComponent) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(gc).Error
}

func (r *Repository) CreateSalaryEmployeeComponent(ctx context.Context, ec *SalaryEmployeeComponent) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(ec).Error
}

func (r *Repository) CreateSalaryEmployeeAdjustment(ctx context.Context, a *SalaryEmployeeAdjustment) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(a).Error
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

// FindActiveBpjsSettingByDate mengambil setting BPJS ACTIVE yang berlaku pada
// tanggal tertentu (paling baru berdasarkan effective_start_date).
func (r *Repository) FindActiveBpjsSettingByDate(ctx context.Context, asOfDate string) (*BpjsSetting, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var bs BpjsSetting
	err = db.Where(
		"status = ? AND effective_start_date <= ? AND (effective_end_date IS NULL OR effective_end_date >= ?)",
		"ACTIVE", asOfDate, asOfDate,
	).Order("effective_start_date DESC").First(&bs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &bs, nil
}

// FindActiveBpjsRateComponentsBySettingID mengambil rate component ACTIVE untuk
// sebuah setting yang berlaku pada tanggal tertentu.
func (r *Repository) FindActiveBpjsRateComponentsBySettingID(ctx context.Context, settingID uuid.UUID, asOfDate string) ([]BpjsRateComponent, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []BpjsRateComponent
	if err := db.Where(
		"bpjs_setting_id = ? AND status = ? AND effective_start_date <= ? AND (effective_end_date IS NULL OR effective_end_date >= ?)",
		settingID, "ACTIVE", asOfDate, asOfDate,
	).Order("display_order ASC, rate_code ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindActiveEmployeeBpjsProfileByEmployeeID mengambil profile BPJS employee yang
// berlaku pada tanggal tertentu (paling baru).
func (r *Repository) FindActiveEmployeeBpjsProfileByEmployeeID(ctx context.Context, employeeID uuid.UUID, asOfDate string) (*EmployeeBpjsProfile, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var p EmployeeBpjsProfile
	err = db.Where(
		"employee_id = ? AND status = ? AND effective_start_date <= ? AND (effective_end_date IS NULL OR effective_end_date >= ?)",
		employeeID, "ACTIVE", asOfDate, asOfDate,
	).Order("effective_start_date DESC").First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
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

// FindActivePph21SettingByDate mengambil setting PPh21 ACTIVE yang berlaku pada
// tanggal tertentu (paling baru berdasarkan effective_start_date).
func (r *Repository) FindActivePph21SettingByDate(ctx context.Context, asOfDate string) (*Pph21Setting, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var ps Pph21Setting
	err = db.Where(
		"status = ? AND effective_start_date <= ? AND (effective_end_date IS NULL OR effective_end_date >= ?)",
		"ACTIVE", asOfDate, asOfDate,
	).Order("effective_start_date DESC").First(&ps).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &ps, nil
}

// FindActivePph21PtkpRatesByDate mengambil seluruh PTKP rate ACTIVE yang berlaku
// pada tanggal tertentu.
func (r *Repository) FindActivePph21PtkpRatesByDate(ctx context.Context, asOfDate string) ([]Pph21PtkpRate, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []Pph21PtkpRate
	if err := db.Where(
		"status = ? AND effective_start_date <= ? AND (effective_end_date IS NULL OR effective_end_date >= ?)",
		"ACTIVE", asOfDate, asOfDate,
	).Order("ptkp_status ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindActivePph21TaxBracketsByDate mengambil seluruh tax bracket ACTIVE yang
// berlaku pada tanggal tertentu, urut bracket_order.
func (r *Repository) FindActivePph21TaxBracketsByDate(ctx context.Context, asOfDate string) ([]Pph21TaxBracket, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []Pph21TaxBracket
	if err := db.Where(
		"status = ? AND effective_start_date <= ? AND (effective_end_date IS NULL OR effective_end_date >= ?)",
		"ACTIVE", asOfDate, asOfDate,
	).Order("bracket_order ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindActiveEmployeeTaxProfileByEmployeeID mengambil profil pajak employee yang
// berlaku pada tanggal tertentu (paling baru).
func (r *Repository) FindActiveEmployeeTaxProfileByEmployeeID(ctx context.Context, employeeID uuid.UUID, asOfDate string) (*EmployeeTaxProfile, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var p EmployeeTaxProfile
	err = db.Where(
		"employee_id = ? AND status = ? AND effective_start_date <= ? AND (effective_end_date IS NULL OR effective_end_date >= ?)",
		employeeID, "ACTIVE", asOfDate, asOfDate,
	).Order("effective_start_date DESC").First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
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
// Read models dari modul lain (dipakai kalkulasi payroll run)
// =============================================================================

// FindEmployeesByIDs mengambil data employee (read-only) untuk snapshot run.
func (r *Repository) FindEmployeesByIDs(ctx context.Context, ids []uuid.UUID) ([]EmployeeRead, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}
	var items []EmployeeRead
	if err := db.Where("id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindAllActiveEmployees mengambil seluruh employee berstatus active.
func (r *Repository) FindAllActiveEmployees(ctx context.Context) ([]EmployeeRead, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []EmployeeRead
	if err := db.Where("status = ?", "active").Order("name ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindActiveEmploymentByEmployeeID mengambil employment yang berlaku pada
// tanggal tertentu (paling baru yang effective_date <= asOfDate).
func (r *Repository) FindActiveEmploymentByEmployeeID(ctx context.Context, employeeID uuid.UUID, asOfDate string) (*EmploymentRead, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var emp EmploymentRead
	err = db.Where(
		"employee_id = ? AND effective_date <= ? AND (effective_end_date IS NULL OR effective_end_date >= ?)",
		employeeID, asOfDate, asOfDate,
	).Order("effective_date DESC").First(&emp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &emp, nil
}

// FindEmploymentByEmployeeIDForPeriod mengambil employment yang OVERLAP dengan
// rentang periode: effective_date <= periodEnd dan berakhir sejak periodStart
// (atau tanpa tanggal akhir). Mencakup employee yang join dan/atau resign
// tengah bulan — dipakai kalkulasi payroll run agar keduanya tetap dihitung
// (dengan prorasi).
func (r *Repository) FindEmploymentByEmployeeIDForPeriod(ctx context.Context, employeeID uuid.UUID, periodStart, periodEnd string) (*EmploymentRead, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var emp EmploymentRead
	err = db.Where(
		"employee_id = ? AND effective_date <= ? AND (effective_end_date IS NULL OR effective_end_date >= ?)",
		employeeID, periodEnd, periodStart,
	).Order("effective_date DESC").First(&emp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &emp, nil
}

// FindPositionByID mengambil position (read-only) untuk resolusi grading.
func (r *Repository) FindPositionByID(ctx context.Context, id uuid.UUID) (*PositionRead, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var pos PositionRead
	if err := db.First(&pos, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &pos, nil
}

// FindGradingByID mengambil grading (read-only) untuk snapshot nama grading.
func (r *Repository) FindGradingByID(ctx context.Context, id uuid.UUID) (*GradingRead, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var g GradingRead
	if err := db.First(&g, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &g, nil
}

// FindAllSalaryGradeComponentsByGradingID mengambil komponen default per grade
// yang berlaku pada asOfDate (status ACTIVE).
func (r *Repository) FindAllSalaryGradeComponentsByGradingID(ctx context.Context, gradingID uuid.UUID, asOfDate string) ([]SalaryGradeComponent, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []SalaryGradeComponent
	if err := db.Where(
		"grading_id = ? AND status = ? AND effective_start_date <= ? AND (effective_end_date IS NULL OR effective_end_date >= ?)",
		gradingID, "ACTIVE", asOfDate, asOfDate,
	).Order("created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindAllSalaryEmployeeComponentsByEmployeeID mengambil komponen override per
// employee yang berlaku pada asOfDate (status ACTIVE).
func (r *Repository) FindAllSalaryEmployeeComponentsByEmployeeID(ctx context.Context, employeeID uuid.UUID, asOfDate string) ([]SalaryEmployeeComponent, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []SalaryEmployeeComponent
	if err := db.Where(
		"employee_id = ? AND status = ? AND effective_start_date <= ? AND (effective_end_date IS NULL OR effective_end_date >= ?)",
		employeeID, "ACTIVE", asOfDate, asOfDate,
	).Order("created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindAllSalaryEmployeeAdjustmentsByPeriod mengambil penyesuaian sekali-jalan
// untuk employee pada periode tertentu (status APPROVED atau APPLIED).
func (r *Repository) FindAllSalaryEmployeeAdjustmentsByPeriod(ctx context.Context, employeeID uuid.UUID, year, month int) ([]SalaryEmployeeAdjustment, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []SalaryEmployeeAdjustment
	if err := db.Where(
		"employee_id = ? AND period_year = ? AND period_month = ? AND status IN (?)",
		employeeID, year, month, []string{"APPROVED", "APPLIED"},
	).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
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

// UpdatePayrollRunEmployee menyimpan ulang row payroll_run_employees (totals).
func (r *Repository) UpdatePayrollRunEmployee(ctx context.Context, pre *PayrollRunEmployee) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(pre).Error
}

// DeletePayrollRunEmployeesByRunID menghapus seluruh snapshot employee sebuah
// run (dipakai recalculation agar snapshot diganti bersih).
func (r *Repository) DeletePayrollRunEmployeesByRunID(ctx context.Context, runID uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("payroll_run_id = ?", runID).Delete(&PayrollRunEmployee{}).Error
}

// DeletePayrollRunItemsByRunID menghapus seluruh item snapshot sebuah run.
func (r *Repository) DeletePayrollRunItemsByRunID(ctx context.Context, runID uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("payroll_run_id = ?", runID).Delete(&PayrollRunItem{}).Error
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

// FindPph21CalculationLogsByRunID mengambil seluruh log kalkulasi PPh21 sebuah
// run (dipakai laporan pajak).
func (r *Repository) FindPph21CalculationLogsByRunID(ctx context.Context, runID uuid.UUID) ([]Pph21CalculationLog, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []Pph21CalculationLog
	if err := db.Where("payroll_run_id = ?", runID).Order("created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindEmployeeBpjsProfilesByEmployeeIDs mengambil profil BPJS employee untuk
// sekumpulan employee (dipakai laporan BPJS untuk nomor kepesertaan).
func (r *Repository) FindEmployeeBpjsProfilesByEmployeeIDs(ctx context.Context, employeeIDs []uuid.UUID) ([]EmployeeBpjsProfile, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	if len(employeeIDs) == 0 {
		return nil, nil
	}
	var items []EmployeeBpjsProfile
	if err := db.Where("employee_id IN ?", employeeIDs).Order("employee_id ASC, effective_start_date DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// CountPayrollPayslipsByRunID menghitung payslip sebuah run (untuk nomor urut).
func (r *Repository) CountPayrollPayslipsByRunID(ctx context.Context, runID uuid.UUID) (int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return 0, err
	}
	var count int64
	if err := db.Model(&PayrollPayslip{}).Where("payroll_run_id = ?", runID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// DeletePayrollPayslipsByRunID menghapus seluruh payslip sebuah run (dipakai
// regenerasi payslip).
func (r *Repository) DeletePayrollPayslipsByRunID(ctx context.Context, runID uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("payroll_run_id = ?", runID).Delete(&PayrollPayslip{}).Error
}

// =============================================================================
// Payroll Payments
// =============================================================================

func (r *Repository) BulkCreatePayrollPayments(ctx context.Context, payments []PayrollPayment) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.CreateInBatches(payments, 50).Error
}

func (r *Repository) FindPayrollPaymentByID(ctx context.Context, id uuid.UUID) (*PayrollPayment, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var p PayrollPayment
	if err := db.First(&p, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}
	return &p, nil
}

func (r *Repository) FindPayrollPaymentsByRunID(ctx context.Context, runID uuid.UUID) ([]PayrollPayment, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []PayrollPayment
	if err := db.Where("payroll_run_id = ?", runID).Order("employee_name ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) UpdatePayrollPayment(ctx context.Context, p *PayrollPayment) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(p).Error
}

// DeletePayrollPaymentsByRunID menghapus seluruh payment sebuah run (dipakai
// regenerasi batch).
func (r *Repository) DeletePayrollPaymentsByRunID(ctx context.Context, runID uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("payroll_run_id = ?", runID).Delete(&PayrollPayment{}).Error
}

// =============================================================================
// Read model bank profile (dipakai snapshot payment batch)
// =============================================================================

// FindActivePrimaryBankProfileByEmployeeID mengambil bank profile utama yang
// berlaku pada tanggal tertentu (is_primary + ACTIVE + rentang effective).
func (r *Repository) FindActivePrimaryBankProfileByEmployeeID(ctx context.Context, employeeID uuid.UUID, asOfDate string) (*EmployeeBankProfile, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var b EmployeeBankProfile
	err = db.Where(
		"employee_id = ? AND is_primary = ? AND status = ? AND effective_start_date <= ? AND (effective_end_date IS NULL OR effective_end_date >= ?)",
		employeeID, true, "ACTIVE", asOfDate, asOfDate,
	).Order("effective_start_date DESC").First(&b).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}
