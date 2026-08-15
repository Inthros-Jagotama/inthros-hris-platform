package payroll

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =============================================================================
// Payroll Structure Models (Migration 006)
// =============================================================================

// SalaryComponent represents master data for salary components.
type SalaryComponent struct {
	ID                     uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Code                   string    `gorm:"type:varchar(50);not null;uniqueIndex:uk_salary_comp_code" json:"code"`
	Name                   string    `gorm:"type:varchar(150);not null" json:"name"`
	Description            *string   `gorm:"type:text" json:"description,omitempty"`
	ComponentType          string    `gorm:"type:varchar(255);not null" json:"component_type"` // EARNING, DEDUCTION, EMPLOYER_CONTRIBUTION, INFORMATION
	CalculationType        string    `gorm:"type:varchar(255);not null;default:FIXED" json:"calculation_type"` // FIXED, PERCENTAGE, FORMULA, REFERENCE, MANUAL
	Formula                *string   `gorm:"type:text" json:"formula,omitempty"`                            // expression utk PERCENTAGE/FORMULA (di-parse oleh Formula Engine)
	ReferenceComponentID   *uuid.UUID `gorm:"type:char(36);index" json:"reference_component_id,omitempty"`   // self-FK utk REFERENCE
	IsTaxable              bool      `gorm:"not null" json:"is_taxable"`
	IsBpjsBase             bool      `gorm:"not null" json:"is_bpjs_base"`
	IsRecurring            bool      `gorm:"not null" json:"is_recurring"`
	IsProratable           bool      `gorm:"not null" json:"is_proratable"`
	PrintOnSalaryStructure bool      `gorm:"not null" json:"print_on_salary_structure"`
	// IsPph21Component: komponen "wadah" baris potongan hasil pajak PPh21
	// (sumber kebenaran di komponen gaji, bukan di setting PPh21).
	IsPph21Component bool `gorm:"not null;default:0" json:"is_pph21_component"`
	// IsPph21Deductible: komponen pengurang penghasilan bruto PPh21
	// (mis. iuran pensiun — bisa lebih dari satu komponen).
	IsPph21Deductible bool `gorm:"not null;default:0" json:"is_pph21_deductible"`
	DisplayOrder           int       `gorm:"not null;default:100" json:"display_order"`
	Status                 string    `gorm:"type:varchar(255);not null;default:ACTIVE" json:"status"`
	CreatedBy              *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy              *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

func (SalaryComponent) TableName() string { return "salary_components" }

func (s *SalaryComponent) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// SalaryGradeComponent links salary components to grading levels.
type SalaryGradeComponent struct {
	ID                uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	GradingID         *uuid.UUID `gorm:"type:char(36);index" json:"grading_id,omitempty"`
	SalaryComponentID uuid.UUID  `gorm:"type:char(36);not null" json:"salary_component_id"`
	Amount            float64    `gorm:"type:decimal(18,2);not null;default:0" json:"amount"`
	CurrencyCode      string     `gorm:"type:char(3);not null;default:IDR" json:"currency_code"`
	EffectiveStartDate string    `gorm:"type:date;not null" json:"effective_start_date"`
	EffectiveEndDate  *string    `gorm:"type:date" json:"effective_end_date,omitempty"`
	IsMandatory       bool       `gorm:"not null;default:1" json:"is_mandatory"`
	IsDefault         bool       `gorm:"not null;default:1" json:"is_default"`
	Status            string     `gorm:"type:varchar(255);not null;default:ACTIVE" json:"status"`
	Notes             *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedBy         *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy         *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func (SalaryGradeComponent) TableName() string { return "salary_grade_components" }

func (s *SalaryGradeComponent) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// SalaryEmployeeComponent stores per-employee salary component amounts.
type SalaryEmployeeComponent struct {
	ID                uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID        uuid.UUID  `gorm:"type:char(36);not null;index" json:"employee_id"`
	EmploymentID      *uuid.UUID `gorm:"type:char(36);index" json:"employment_id,omitempty"`
	PositionID        *uuid.UUID `gorm:"type:char(36);index" json:"position_id,omitempty"`
	GradingID         *uuid.UUID `gorm:"type:char(36);index" json:"grading_id,omitempty"`
	SalaryComponentID uuid.UUID  `gorm:"type:char(36);not null" json:"salary_component_id"`
	Amount            float64    `gorm:"type:decimal(18,2);not null;default:0" json:"amount"`
	CurrencyCode      string     `gorm:"type:char(3);not null;default:IDR" json:"currency_code"`
	SourceType        string     `gorm:"type:varchar(255);not null;default:MANUAL" json:"source_type"`
	SourceRefID       *uuid.UUID `gorm:"type:char(36)" json:"source_ref_id,omitempty"`
	EffectiveStartDate string    `gorm:"type:date;not null" json:"effective_start_date"`
	EffectiveEndDate  *string    `gorm:"type:date" json:"effective_end_date,omitempty"`
	Notes             *string    `gorm:"type:text" json:"notes,omitempty"`
	Status            string     `gorm:"type:varchar(255);not null;default:ACTIVE" json:"status"`
	CreatedBy         *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy         *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func (SalaryEmployeeComponent) TableName() string { return "salary_employee_components" }

func (s *SalaryEmployeeComponent) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// SalaryChangeLog tracks all salary changes for audit.
type SalaryChangeLog struct {
	ID                       uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID               uuid.UUID  `gorm:"type:char(36);not null;index:idx_salary_change_employee" json:"employee_id"`
	EmployeeSalaryComponentID *uuid.UUID `gorm:"type:char(36)" json:"employee_salary_component_id,omitempty"`
	SalaryComponentID        *uuid.UUID `gorm:"type:char(36);index" json:"salary_component_id,omitempty"`
	ActionType               string     `gorm:"type:varchar(255);not null;index" json:"action_type"` // CREATE, UPDATE, END, REACTIVATE, DELETE
	OldAmount                *float64   `gorm:"type:decimal(18,2)" json:"old_amount,omitempty"`
	NewAmount                *float64   `gorm:"type:decimal(18,2)" json:"new_amount,omitempty"`
	OldEffectiveStartDate    *string    `gorm:"type:date" json:"old_effective_start_date,omitempty"`
	NewEffectiveStartDate    *string    `gorm:"type:date" json:"new_effective_start_date,omitempty"`
	OldEffectiveEndDate      *string    `gorm:"type:date" json:"old_effective_end_date,omitempty"`
	NewEffectiveEndDate      *string    `gorm:"type:date" json:"new_effective_end_date,omitempty"`
	Reason                   *string    `gorm:"type:varchar(255)" json:"reason,omitempty"`
	Notes                    *string    `gorm:"type:text" json:"notes,omitempty"`
	BeforeJSON               *string    `gorm:"type:json" json:"before_json,omitempty"`
	AfterJSON                *string    `gorm:"type:json" json:"after_json,omitempty"`
	ChangedBy                *uuid.UUID `gorm:"type:char(36)" json:"changed_by,omitempty"`
	ChangedAt                *time.Time `gorm:"type:timestamp" json:"changed_at,omitempty"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}

func (SalaryChangeLog) TableName() string { return "salary_change_logs" }

func (s *SalaryChangeLog) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// SalaryEmployeeAdjustment represents one-time adjustments per period.
type SalaryEmployeeAdjustment struct {
	ID                uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID        uuid.UUID  `gorm:"type:char(36);not null;index:idx_adj_employee_period" json:"employee_id"`
	EmploymentID      *uuid.UUID `gorm:"type:char(36);index" json:"employment_id,omitempty"`
	PositionID        *uuid.UUID `gorm:"type:char(36);index" json:"position_id,omitempty"`
	SalaryComponentID uuid.UUID  `gorm:"type:char(36);not null;index" json:"salary_component_id"`
	PeriodYear        int        `gorm:"not null" json:"period_year"`
	PeriodMonth       int        `gorm:"type:tinyint;not null" json:"period_month"`
	Amount            float64    `gorm:"type:decimal(18,2);not null;default:0" json:"amount"`
	CurrencyCode      string     `gorm:"type:char(3);not null;default:IDR" json:"currency_code"`
	SourceType        string     `gorm:"type:varchar(255);not null;default:MANUAL" json:"source_type"`
	Reason            *string    `gorm:"type:varchar(255)" json:"reason,omitempty"`
	Notes             *string    `gorm:"type:text" json:"notes,omitempty"`
	Status            string     `gorm:"type:varchar(255);not null;default:DRAFT" json:"status"` // DRAFT, APPROVED, CANCELLED, APPLIED
	ApprovedBy        *uuid.UUID `gorm:"type:char(36)" json:"approved_by,omitempty"`
	ApprovedAt        *time.Time `gorm:"type:timestamp" json:"approved_at,omitempty"`
	CreatedBy         *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy         *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func (SalaryEmployeeAdjustment) TableName() string { return "salary_employee_adjustments" }

func (s *SalaryEmployeeAdjustment) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// PayrollPeriod represents a payroll period.
type PayrollPeriod struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	PeriodCode  string     `gorm:"type:varchar(50);not null;uniqueIndex:uk_period_code" json:"period_code"`
	PeriodYear  int        `gorm:"not null;uniqueIndex:uk_period_year_month" json:"period_year"`
	PeriodMonth int        `gorm:"type:tinyint;not null;uniqueIndex:uk_period_year_month" json:"period_month"`
	StartDate   string     `gorm:"type:date;not null" json:"start_date"`
	EndDate     string     `gorm:"type:date;not null" json:"end_date"`
	AsOfDate    string     `gorm:"type:date;not null" json:"as_of_date"`
	Status      string     `gorm:"type:varchar(255);not null;default:OPEN" json:"status"` // OPEN, CLOSED
	CreatedBy   *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy   *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (PayrollPeriod) TableName() string { return "payroll_periods" }

func (p *PayrollPeriod) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// EmployeePayrollProfile stores payroll configuration per employee.
type EmployeePayrollProfile struct {
	ID                 uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID         uuid.UUID  `gorm:"type:char(36);not null;uniqueIndex:uk_payroll_profile_start" json:"employee_id"`
	EmploymentID       *uuid.UUID `gorm:"type:char(36);index" json:"employment_id,omitempty"`
	PayrollGroupCode   string     `gorm:"type:varchar(50);not null;default:MONTHLY" json:"payroll_group_code"`
	PayrollFrequency   string     `gorm:"type:varchar(255);not null;default:MONTHLY" json:"payroll_frequency"` // MONTHLY, WEEKLY, DAILY
	PaymentMethod      string     `gorm:"type:varchar(255);not null;default:BANK_TRANSFER" json:"payment_method"`
	SalaryCurrency     string     `gorm:"type:char(3);not null;default:IDR" json:"salary_currency"`
	IsPayrollActive    bool       `gorm:"not null;default:1" json:"is_payroll_active"`
	EffectiveStartDate string     `gorm:"type:date;not null" json:"effective_start_date"`
	EffectiveEndDate   *string    `gorm:"type:date" json:"effective_end_date,omitempty"`
	Status             string     `gorm:"type:varchar(255);not null;default:ACTIVE" json:"status"`
	Notes              *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedBy          *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy          *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func (EmployeePayrollProfile) TableName() string { return "employee_payroll_profiles" }

func (e *EmployeePayrollProfile) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// EmployeeBankProfile stores bank account info for payroll.
type EmployeeBankProfile struct {
	ID                     uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID             uuid.UUID  `gorm:"type:char(36);not null;uniqueIndex:uk_bank_profile_start" json:"employee_id"`
	EmployeePayrollProfileID uuid.UUID `gorm:"type:char(36);not null" json:"employee_payroll_profile_id"`
	BankCode               *string    `gorm:"type:varchar(50)" json:"bank_code,omitempty"`
	BankName               string     `gorm:"type:varchar(150);not null" json:"bank_name"`
	BankBranch             *string    `gorm:"type:varchar(150)" json:"bank_branch,omitempty"`
	BankAccountNumber      string     `gorm:"type:varchar(100);not null" json:"bank_account_number"`
	BankAccountHolderName  string     `gorm:"type:varchar(255);not null" json:"bank_account_holder_name"`
	IsPrimary              bool       `gorm:"not null;default:1" json:"is_primary"`
	EffectiveStartDate     string     `gorm:"type:date;not null" json:"effective_start_date"`
	EffectiveEndDate       *string    `gorm:"type:date" json:"effective_end_date,omitempty"`
	Status                 string     `gorm:"type:varchar(255);not null;default:ACTIVE" json:"status"`
	Notes                  *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedBy              *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy              *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

func (EmployeeBankProfile) TableName() string { return "employee_bank_profiles" }

func (e *EmployeeBankProfile) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// EmployeeBpjsProfile stores BPJS configuration per employee.
type EmployeeBpjsProfile struct {
	ID                       uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID               uuid.UUID  `gorm:"type:char(36);not null;uniqueIndex:uk_bpjs_profile_start" json:"employee_id"`
	EmployeePayrollProfileID uuid.UUID  `gorm:"type:char(36);not null" json:"employee_payroll_profile_id"`
	BpjsHealthActive         bool       `gorm:"not null;default:0" json:"bpjs_health_active"`
	BpjsHealthNo             *string    `gorm:"type:varchar(50)" json:"bpjs_health_no,omitempty"`
	BpjsHealthRegisteredName *string    `gorm:"type:varchar(255)" json:"bpjs_health_registered_name,omitempty"`
	BpjsTkActive             bool       `gorm:"not null;default:0" json:"bpjs_tk_active"`
	BpjsTkNo                 *string    `gorm:"type:varchar(50)" json:"bpjs_tk_no,omitempty"`
	BpjsTkRegisteredName     *string    `gorm:"type:varchar(255)" json:"bpjs_tk_registered_name,omitempty"`
	JkkRiskClass             string     `gorm:"type:varchar(255);not null;default:LOW" json:"jkk_risk_class"`
	PensionActive            bool       `gorm:"not null;default:1" json:"pension_active"`
	EffectiveStartDate       string     `gorm:"type:date;not null" json:"effective_start_date"`
	EffectiveEndDate         *string    `gorm:"type:date" json:"effective_end_date,omitempty"`
	Status                   string     `gorm:"type:varchar(255);not null;default:ACTIVE" json:"status"`
	Notes                    *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedBy                *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy                *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}

func (EmployeeBpjsProfile) TableName() string { return "employee_bpjs_profiles" }

func (e *EmployeeBpjsProfile) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// EmployeeTaxProfile stores tax configuration per employee.
type EmployeeTaxProfile struct {
	ID                       uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID               uuid.UUID  `gorm:"type:char(36);not null;uniqueIndex:uk_tax_profile_start" json:"employee_id"`
	EmployeePayrollProfileID uuid.UUID  `gorm:"type:char(36);not null" json:"employee_payroll_profile_id"`
	Npwp                     *string    `gorm:"type:varchar(50)" json:"npwp,omitempty"`
	NpwpRegisteredName       *string    `gorm:"type:varchar(255)" json:"npwp_registered_name,omitempty"`
	PtkpStatus               *string    `gorm:"type:varchar(20)" json:"ptkp_status,omitempty"`
	TaxMethod                string     `gorm:"type:varchar(255);not null;default:GROSS" json:"tax_method"` // GROSS, GROSS_UP, NETT
	IsTaxable                bool       `gorm:"not null;default:1" json:"is_taxable"`
	HasNpwp                  bool       `gorm:"not null;default:0" json:"has_npwp"`
	EffectiveStartDate       string     `gorm:"type:date;not null" json:"effective_start_date"`
	EffectiveEndDate         *string    `gorm:"type:date" json:"effective_end_date,omitempty"`
	Status                   string     `gorm:"type:varchar(255);not null;default:ACTIVE" json:"status"`
	Notes                    *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedBy                *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy                *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}

func (EmployeeTaxProfile) TableName() string { return "employee_tax_profiles" }

func (e *EmployeeTaxProfile) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// BpjsSetting stores BPJS system settings.
type BpjsSetting struct {
	ID                     uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	SettingCode            string     `gorm:"type:varchar(50);not null;uniqueIndex:uk_bpjs_setting_code" json:"setting_code"`
	SettingName            string     `gorm:"type:varchar(150);not null" json:"setting_name"`
	BaseSource             string     `gorm:"type:varchar(255);not null;default:BPJS_BASE_COMPONENTS" json:"base_source"`
	HealthMaxBaseAmount    *float64   `gorm:"type:decimal(18,2)" json:"health_max_base_amount,omitempty"`
	PensionMaxBaseAmount   *float64   `gorm:"type:decimal(18,2)" json:"pension_max_base_amount,omitempty"`
	DefaultJkkRiskClass    string     `gorm:"type:varchar(255);not null;default:LOW" json:"default_jkk_risk_class"`
	RoundingMode           string     `gorm:"type:varchar(255);not null;default:ROUND" json:"rounding_mode"`
	EffectiveStartDate     string     `gorm:"type:date;not null" json:"effective_start_date"`
	EffectiveEndDate       *string    `gorm:"type:date" json:"effective_end_date,omitempty"`
	Status                 string     `gorm:"type:varchar(255);not null;default:ACTIVE" json:"status"`
	Notes                  *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedBy              *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy              *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

func (BpjsSetting) TableName() string { return "bpjs_settings" }

func (b *BpjsSetting) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// BpjsRateComponent stores BPJS rate configurations per program.
type BpjsRateComponent struct {
	ID                     uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	BpjsSettingID          uuid.UUID  `gorm:"type:char(36);not null;index" json:"bpjs_setting_id"`
	RateCode               string     `gorm:"type:varchar(80);not null;uniqueIndex:uk_bpjs_rate_code_start" json:"rate_code"`
	RateName               string     `gorm:"type:varchar(180);not null" json:"rate_name"`
	BpjsProgram            string     `gorm:"type:varchar(255);not null" json:"bpjs_program"` // HEALTH, JHT, JP, JKK, JKM, JKP
	PaidBy                 string     `gorm:"type:varchar(255);not null" json:"paid_by"`      // EMPLOYEE, EMPLOYER
	SalaryComponentID      *uuid.UUID `gorm:"type:char(36);index" json:"salary_component_id,omitempty"`
	RatePercent            float64    `gorm:"type:decimal(8,4);not null;default:0" json:"rate_percent"`
	FixedAmount            *float64   `gorm:"type:decimal(18,2)" json:"fixed_amount,omitempty"`
	MinBaseAmount          *float64   `gorm:"type:decimal(18,2)" json:"min_base_amount,omitempty"`
	MaxBaseAmount          *float64   `gorm:"type:decimal(18,2)" json:"max_base_amount,omitempty"`
	JkkRiskClass           *string    `gorm:"type:varchar(255)" json:"jkk_risk_class,omitempty"`
	IsEmployeeDeduction    bool       `gorm:"not null;default:0" json:"is_employee_deduction"`
	IsEmployerContribution bool       `gorm:"not null;default:0" json:"is_employer_contribution"`
	GenerateToPayrollItem  bool       `gorm:"not null;default:1" json:"generate_to_payroll_item"`
	PrintOnPayslip         bool       `gorm:"not null;default:1" json:"print_on_payslip"`
	DisplayOrder           int        `gorm:"not null;default:0" json:"display_order"`
	EffectiveStartDate     string     `gorm:"type:date;not null" json:"effective_start_date"`
	EffectiveEndDate       *string    `gorm:"type:date" json:"effective_end_date,omitempty"`
	Status                 string     `gorm:"type:varchar(255);not null;default:ACTIVE" json:"status"`
	Notes                  *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedBy              *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy              *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

func (BpjsRateComponent) TableName() string { return "bpjs_rate_components" }

func (b *BpjsRateComponent) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// Pph21Setting stores PPh21 system configuration.
type Pph21Setting struct {
	ID                           uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	SettingCode                  string     `gorm:"type:varchar(50);not null;uniqueIndex:uk_pph21_setting_code" json:"setting_code"`
	SettingName                  string     `gorm:"type:varchar(150);not null" json:"setting_name"`
	CalculationMethod            string     `gorm:"type:varchar(255);not null;default:REGULAR_GROSS_ANNUALIZED" json:"calculation_method"`
	DefaultTaxMethod             string     `gorm:"type:varchar(255);not null;default:GROSS" json:"default_tax_method"`
	OccupationalExpenseRatePercent float64  `gorm:"type:decimal(8,4);not null;default:5" json:"occupational_expense_rate_percent"`
	OccupationalExpenseMaxMonthly float64   `gorm:"type:decimal(18,2);not null;default:500000" json:"occupational_expense_max_monthly"`
	OccupationalExpenseMaxYearly float64    `gorm:"type:decimal(18,2);not null;default:6000000" json:"occupational_expense_max_yearly"`
	DeductBpjsHealthEmployee     bool       `gorm:"not null;default:0" json:"deduct_bpjs_health_employee"`
	DeductBpjsJhtEmployee        bool       `gorm:"not null;default:1" json:"deduct_bpjs_jht_employee"`
	DeductBpjsJpEmployee         bool       `gorm:"not null;default:1" json:"deduct_bpjs_jp_employee"`
	AnnualizationMonths          int        `gorm:"type:tinyint;not null;default:12" json:"annualization_months"`
	PkpRoundingUnit              float64    `gorm:"type:decimal(18,2);not null;default:1000" json:"pkp_rounding_unit"`
	NonNpwpMultiplierPercent     float64    `gorm:"type:decimal(8,4);not null;default:100" json:"non_npwp_multiplier_percent"`
	RoundingMode                 string     `gorm:"type:varchar(255);not null;default:ROUND" json:"rounding_mode"`
	EffectiveStartDate           string     `gorm:"type:date;not null" json:"effective_start_date"`
	EffectiveEndDate             *string    `gorm:"type:date" json:"effective_end_date,omitempty"`
	Status                       string     `gorm:"type:varchar(255);not null;default:ACTIVE" json:"status"`
	Notes                        *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedBy                    *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy                    *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt                    time.Time  `json:"created_at"`
	UpdatedAt                    time.Time  `json:"updated_at"`
}

func (Pph21Setting) TableName() string { return "pph21_settings" }

func (p *Pph21Setting) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// TerRate membaca tabel ters (Tarif Efektif Rata-rata) — master data di luar
// modul payroll (migration 001). Satu baris = rentang bruto bulanan + rate.
type TerRate struct {
	ID       uuid.UUID `gorm:"column:id;type:char(36);primaryKey" json:"id"`
	Group    string    `gorm:"column:group;type:char(1);not null;index:idx_ter_group" json:"group"`
	BrutoMin *int64    `gorm:"column:bruto_min" json:"bruto_min,omitempty"`
	BrutoMax *int64    `gorm:"column:bruto_max" json:"bruto_max,omitempty"`
	Rate     float64   `gorm:"column:rate;type:decimal(10,2);not null" json:"rate"`
}

func (TerRate) TableName() string { return "ters" }

// Ptkp adalah pandangan read-only ke tabel ptkps (setting module) — satu-satunya
// sumber kebenaran PTKP untuk engine PPh21 (migration 121 menggabungkan
// pph21_ptkp_rates ke ptkps). Lookup engine memakai kolom code (mis. "TK0").
type Ptkp struct {
	ID    uuid.UUID `gorm:"column:id;type:char(36);primaryKey" json:"id"`
	Name  string    `gorm:"column:name" json:"name"`
	Code  string    `gorm:"column:code" json:"code"`
	Ptkp  int64     `gorm:"column:ptkp" json:"ptkp"`
	Group string    `gorm:"column:group" json:"group"`
}

func (Ptkp) TableName() string { return "ptkps" }

// Pph21TaxBracket stores progressive tax rate brackets.
type Pph21TaxBracket struct {
	ID                 uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	BracketOrder       int        `gorm:"not null;uniqueIndex:uk_pph21_bracket_order_start" json:"bracket_order"`
	LowerBound         float64    `gorm:"type:decimal(18,2);not null;default:0" json:"lower_bound"`
	UpperBound         *float64   `gorm:"type:decimal(18,2)" json:"upper_bound,omitempty"`
	RatePercent        float64    `gorm:"type:decimal(8,4);not null;default:0" json:"rate_percent"`
	EffectiveStartDate string     `gorm:"type:date;not null" json:"effective_start_date"`
	EffectiveEndDate   *string    `gorm:"type:date" json:"effective_end_date,omitempty"`
	Status             string     `gorm:"type:varchar(255);not null;default:ACTIVE" json:"status"`
	CreatedBy          *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy          *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func (Pph21TaxBracket) TableName() string { return "pph21_tax_brackets" }

func (p *Pph21TaxBracket) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// =============================================================================
// Payroll Run Models (Migration 007)
// =============================================================================

// PayrollRun represents a payroll run.
type PayrollRun struct {
	ID                    uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	PayrollPeriodID       uuid.UUID  `gorm:"type:char(36);not null;index" json:"payroll_period_id"`
	RunCode               string     `gorm:"type:varchar(50);not null;uniqueIndex:uk_payroll_run_code" json:"run_code"`
	RunType               string     `gorm:"type:varchar(255);not null;default:REGULAR" json:"run_type"` // REGULAR, OFF_CYCLE, THR, BONUS
	ProrationMethod       string     `gorm:"type:varchar(255);not null;default:CALENDAR_DAYS" json:"proration_method"` // CALENDAR_DAYS, WORKING_DAYS, FIXED_30_DAYS, ATTENDANCE_DAYS (migration 117)
	Status                string     `gorm:"type:varchar(255);not null;default:DRAFT" json:"status"`     // DRAFT, CALCULATED, REVIEWED, APPROVED, LOCKED, CANCELLED
	TotalEmployees        int        `gorm:"not null;default:0" json:"total_employees"`
	TotalEarning          float64    `gorm:"type:decimal(18,2);not null;default:0" json:"total_earning"`
	TotalDeduction        float64    `gorm:"type:decimal(18,2);not null;default:0" json:"total_deduction"`
	TotalEmployerContribution float64 `gorm:"type:decimal(18,2);not null;default:0" json:"total_employer_contribution"`
	TotalNet              float64    `gorm:"type:decimal(18,2);not null;default:0" json:"total_net"`
	TotalCompanyCost      float64    `gorm:"type:decimal(18,2);not null;default:0" json:"total_company_cost"`
	CalculatedAt          *time.Time `gorm:"type:timestamp" json:"calculated_at,omitempty"`
	ReviewedAt            *time.Time `gorm:"type:timestamp" json:"reviewed_at,omitempty"`
	ApprovedAt            *time.Time `gorm:"type:timestamp" json:"approved_at,omitempty"`
	LockedAt              *time.Time `gorm:"type:timestamp" json:"locked_at,omitempty"`
	ApprovalInstanceID    *uuid.UUID `gorm:"type:char(36)" json:"approval_instance_id,omitempty"`
	CreatedBy             *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy             *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

func (PayrollRun) TableName() string { return "payroll_runs" }

func (p *PayrollRun) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// PayrollRunEmployee represents an employee within a payroll run.
type PayrollRunEmployee struct {
	ID                     uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	PayrollRunID           uuid.UUID `gorm:"type:char(36);not null;uniqueIndex:uk_payroll_run_emp" json:"payroll_run_id"`
	EmployeeID             uuid.UUID `gorm:"type:char(36);not null;uniqueIndex:uk_payroll_run_emp;index" json:"employee_id"`
	EmploymentID           *uuid.UUID `gorm:"type:char(36);index" json:"employment_id,omitempty"`
	PositionID             *uuid.UUID `gorm:"type:char(36);index" json:"position_id,omitempty"`
	GradingID              *uuid.UUID `gorm:"type:char(36)" json:"grading_id,omitempty"`
	EmployeeCode           string     `gorm:"type:varchar(50);not null" json:"employee_code"`
	EmployeeName           string     `gorm:"type:varchar(255);not null" json:"employee_name"`
	PositionTitle          *string    `gorm:"type:varchar(200)" json:"position_title,omitempty"`
	GradingName            *string    `gorm:"type:varchar(255)" json:"grading_name,omitempty"`
	TotalEarning           float64    `gorm:"type:decimal(18,2);not null;default:0" json:"total_earning"`
	TotalDeduction         float64    `gorm:"type:decimal(18,2);not null;default:0" json:"total_deduction"`
	TotalEmployerContribution float64 `gorm:"type:decimal(18,2);not null;default:0" json:"total_employer_contribution"`
	NetAmount              float64    `gorm:"type:decimal(18,2);not null;default:0" json:"net_amount"`
	TotalCompanyCost       float64    `gorm:"type:decimal(18,2);not null;default:0" json:"total_company_cost"`
	Status                 string     `gorm:"type:varchar(255);not null;default:DRAFT" json:"status"` // DRAFT, CALCULATED, EXCLUDED
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

func (PayrollRunEmployee) TableName() string { return "payroll_run_employees" }

func (p *PayrollRunEmployee) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// PayrollRunItem represents a line item in a payroll run.
type PayrollRunItem struct {
	ID                   uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	PayrollRunID         uuid.UUID  `gorm:"type:char(36);not null;index:idx_payroll_items_run_employee" json:"payroll_run_id"`
	PayrollRunEmployeeID uuid.UUID  `gorm:"type:char(36);not null;index" json:"payroll_run_employee_id"`
	EmployeeID           uuid.UUID  `gorm:"type:char(36);not null;index:idx_payroll_items_run_employee" json:"employee_id"`
	SalaryComponentID    uuid.UUID  `gorm:"type:char(36);not null;index" json:"salary_component_id"`
	ComponentCode        string     `gorm:"type:varchar(50);not null" json:"component_code"`
	ComponentName        string     `gorm:"type:varchar(150);not null" json:"component_name"`
	ComponentType        string     `gorm:"type:varchar(255);not null" json:"component_type"`
	CalculationType      string     `gorm:"type:varchar(255);not null;default:FIXED" json:"calculation_type"` // FIXED, PERCENTAGE, FORMULA, REFERENCE, MANUAL
	ItemCategory         string     `gorm:"type:varchar(255);not null;default:EMPLOYEE_EARNING" json:"item_category"`
	PaidBy               string     `gorm:"type:varchar(255);not null;default:EMPLOYER" json:"paid_by"`
	AffectsGrossPay      bool       `gorm:"not null;default:0" json:"affects_gross_pay"`
	AffectsNetPay        bool       `gorm:"not null;default:0" json:"affects_net_pay"`
	AffectsCompanyCost   bool       `gorm:"not null;default:0" json:"affects_company_cost"`
	PrintOnPayslip       bool       `gorm:"not null;default:1" json:"print_on_payslip"`
	Amount               float64    `gorm:"type:decimal(18,2);not null;default:0" json:"amount"`
	BaseAmount           float64    `gorm:"type:decimal(18,2);not null;default:0" json:"base_amount"`
	Rate                 *float64   `gorm:"type:decimal(8,4)" json:"rate,omitempty"`
	Formula              *string    `gorm:"type:text" json:"formula,omitempty"`
	FormulaResult        *float64   `gorm:"type:decimal(18,2)" json:"formula_result,omitempty"`
	CurrencyCode         string     `gorm:"type:char(3);not null;default:IDR" json:"currency_code"`
	SourceGroup          string     `gorm:"type:varchar(255);not null" json:"source_group"` // STRUCTURE, ADJUSTMENT, STATUTORY
	SourceTable          *string    `gorm:"type:varchar(100)" json:"source_table,omitempty"`
	SourceID             *uuid.UUID `gorm:"type:char(36)" json:"source_id,omitempty"`
	SourceType           *string    `gorm:"type:varchar(50)" json:"source_type,omitempty"`
	Notes                *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

func (PayrollRunItem) TableName() string { return "payroll_run_items" }

func (p *PayrollRunItem) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// PayrollPayslip represents a generated payslip for an employee.
type PayrollPayslip struct {
	ID                   uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	PayrollRunID         uuid.UUID  `gorm:"type:char(36);not null;uniqueIndex:uk_payslip_run_employee" json:"payroll_run_id"`
	PayrollRunEmployeeID uuid.UUID  `gorm:"type:char(36);not null;index" json:"payroll_run_employee_id"`
	EmployeeID           uuid.UUID  `gorm:"type:char(36);not null;index:idx_payslip_employee_period" json:"employee_id"`
	PayslipNumber        string     `gorm:"type:varchar(100);not null;uniqueIndex:uk_payslip_number" json:"payslip_number"`
	PeriodYear           int        `gorm:"not null" json:"period_year"`
	PeriodMonth          int        `gorm:"type:tinyint;not null" json:"period_month"`
	PeriodCode           string     `gorm:"type:varchar(50);not null" json:"period_code"`
	EmployeeCode         string     `gorm:"type:varchar(50);not null" json:"employee_code"`
	EmployeeName         string     `gorm:"type:varchar(255);not null" json:"employee_name"`
	PositionTitle        *string    `gorm:"type:varchar(200)" json:"position_title,omitempty"`
	GradingName          *string    `gorm:"type:varchar(255)" json:"grading_name,omitempty"`
	TotalEarning              float64    `gorm:"type:decimal(18,2);not null;default:0" json:"total_earning"`
	TotalDeduction            float64    `gorm:"type:decimal(18,2);not null;default:0" json:"total_deduction"`
	TotalEmployerContribution float64    `gorm:"type:decimal(18,2);not null;default:0" json:"total_employer_contribution"` // migration 118
	NetAmount                 float64    `gorm:"type:decimal(18,2);not null;default:0" json:"net_amount"`
	Status                    string     `gorm:"type:varchar(255);not null;default:DRAFT" json:"status"` // DRAFT, PUBLISHED, CANCELLED
	GeneratedAt               *time.Time `gorm:"type:timestamp" json:"generated_at,omitempty"`
	PublishedAt               *time.Time `gorm:"type:timestamp" json:"published_at,omitempty"`
	CancelledAt               *time.Time `gorm:"type:timestamp" json:"cancelled_at,omitempty"`
	CreatedBy                 *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy                 *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
}

func (PayrollPayslip) TableName() string { return "payroll_payslips" }

// PayrollPayment mewakili satu baris payment/disbursement dalam batch
// (docs/payroll/07 §27). Nomor rekening adalah SNAPSHOT dari employee bank
// profile saat batch dibuat — perubahan rekening setelah run final tidak
// mengubah batch sebelumnya.
type PayrollPayment struct {
	ID                     uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	PayrollRunID           uuid.UUID  `gorm:"type:char(36);not null;index" json:"payroll_run_id"`
	PayrollRunEmployeeID   uuid.UUID  `gorm:"type:char(36);not null;uniqueIndex:uk_payment_run_employee" json:"payroll_run_employee_id"`
	EmployeeID             uuid.UUID  `gorm:"type:char(36);not null;index:idx_payment_employee" json:"employee_id"`
	EmployeeCode           string     `gorm:"type:varchar(50);not null" json:"employee_code"`
	EmployeeName           string     `gorm:"type:varchar(255);not null" json:"employee_name"`
	Amount                 float64    `gorm:"type:decimal(18,2);not null;default:0" json:"amount"`
	CurrencyCode           string     `gorm:"type:char(3);not null;default:IDR" json:"currency_code"`
	PaymentDate            string     `gorm:"type:date;not null" json:"payment_date"`
	EmployeeBankProfileID  *uuid.UUID `gorm:"type:char(36)" json:"employee_bank_profile_id,omitempty"`
	BankCode               *string    `gorm:"type:varchar(50)" json:"bank_code,omitempty"`
	BankName               *string    `gorm:"type:varchar(150)" json:"bank_name,omitempty"`
	BankBranch             *string    `gorm:"type:varchar(150)" json:"bank_branch,omitempty"`
	BankAccountNumber      string     `gorm:"type:varchar(100);not null" json:"bank_account_number"`
	BankAccountHolderName  string     `gorm:"type:varchar(255);not null" json:"bank_account_holder_name"`
	Status                 string     `gorm:"type:varchar(255);not null;default:PENDING;index:idx_payment_status" json:"status"` // PENDING, PROCESSING, PAID, FAILED, REVERSED
	Reference              *string    `gorm:"type:varchar(100)" json:"reference,omitempty"`
	ProcessedAt            *time.Time `gorm:"type:timestamp" json:"processed_at,omitempty"`
	PaidAt                 *time.Time `gorm:"type:timestamp" json:"paid_at,omitempty"`
	FailedAt               *time.Time `gorm:"type:timestamp" json:"failed_at,omitempty"`
	FailedReason           *string    `gorm:"type:varchar(255)" json:"failed_reason,omitempty"`
	ReversedAt             *time.Time `gorm:"type:timestamp" json:"reversed_at,omitempty"`
	CreatedBy              *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy              *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

func (PayrollPayment) TableName() string { return "payroll_payments" }

func (p *PayrollPayment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (p *PayrollPayslip) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// Pph21CalculationLog stores detailed PPh21 tax calculation logs.
type Pph21CalculationLog struct {
	ID                   uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	PayrollRunID         uuid.UUID  `gorm:"type:char(36);not null;index" json:"payroll_run_id"`
	PayrollRunEmployeeID uuid.UUID  `gorm:"type:char(36);not null;uniqueIndex:uk_pph21_log_run_employee" json:"payroll_run_employee_id"`
	EmployeeID           uuid.UUID  `gorm:"type:char(36);not null;index" json:"employee_id"`
	Pph21SettingID       uuid.UUID  `gorm:"type:char(36);not null;index" json:"pph21_setting_id"`
	EmployeeTaxProfileID uuid.UUID  `gorm:"type:char(36);not null;index" json:"employee_tax_profile_id"`
	PayrollRunItemID     *uuid.UUID `gorm:"type:char(36)" json:"payroll_run_item_id,omitempty"`
	CalculationMethod    string     `gorm:"type:varchar(255);not null;default:REGULAR_GROSS_ANNUALIZED" json:"calculation_method"`
	TaxMethod            string     `gorm:"type:varchar(20);not null" json:"tax_method"`
	PtkpStatus           string     `gorm:"type:varchar(20);not null" json:"ptkp_status"`
	HasNpwp              bool       `gorm:"not null;default:1" json:"has_npwp"`
	GrossMonthly         float64    `gorm:"type:decimal(18,2);not null;default:0" json:"gross_monthly"`
	OccupationalExpenseMonthly float64 `gorm:"type:decimal(18,2);not null;default:0" json:"occupational_expense_monthly"`
	BpjsTaxDeductibleMonthly float64  `gorm:"type:decimal(18,2);not null;default:0" json:"bpjs_tax_deductible_monthly"`
	PensionDeductibleMonthly float64  `gorm:"type:decimal(18,2);not null;default:0" json:"pension_deductible_monthly"`
	NetMonthly           float64    `gorm:"type:decimal(18,2);not null;default:0" json:"net_monthly"`
	NetAnnualized        float64    `gorm:"type:decimal(18,2);not null;default:0" json:"net_annualized"`
	PtkpAnnual           float64    `gorm:"type:decimal(18,2);not null;default:0" json:"ptkp_annual"`
	PkpAnnual            float64    `gorm:"type:decimal(18,2);not null;default:0" json:"pkp_annual"`
	AnnualTaxBeforeNpwpMult float64 `gorm:"type:decimal(18,2);not null;default:0" json:"annual_tax_before_npwp_mult"`
	NonNpwpMultiplierPercent float64 `gorm:"type:decimal(8,4);not null;default:100" json:"non_npwp_multiplier_percent"`
	AnnualTaxAfterNpwpMult  float64  `gorm:"type:decimal(18,2);not null;default:0" json:"annual_tax_after_npwp_mult"`
	Pph21Monthly         float64    `gorm:"type:decimal(18,2);not null;default:0" json:"pph21_monthly"`
	FormulaJSON          *string    `gorm:"type:json" json:"formula_json,omitempty"`
	Status               string     `gorm:"type:varchar(255);not null;default:CALCULATED" json:"status"` // CALCULATED, SKIPPED, ERROR
	Notes                *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedBy            *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy            *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

func (Pph21CalculationLog) TableName() string { return "pph21_calculation_logs" }

func (p *Pph21CalculationLog) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// PayrollProfileChangeLog logs changes to payroll profiles.
type PayrollProfileChangeLog struct {
	ID              uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	EmployeeID      uuid.UUID  `gorm:"type:char(36);not null;index:idx_payroll_profile_log_employee" json:"employee_id"`
	ProfileTable    string     `gorm:"type:varchar(100);not null;index" json:"profile_table"`
	ProfileRecordID *uuid.UUID `gorm:"type:char(36)" json:"profile_record_id,omitempty"`
	ActionType      string     `gorm:"type:varchar(255);not null" json:"action_type"` // CREATE, UPDATE, END, REACTIVATE, DELETE
	Reason          *string    `gorm:"type:varchar(255)" json:"reason,omitempty"`
	Notes           *string    `gorm:"type:text" json:"notes,omitempty"`
	BeforeJSON      *string    `gorm:"type:json" json:"before_json,omitempty"`
	AfterJSON       *string    `gorm:"type:json" json:"after_json,omitempty"`
	ChangedBy       *uuid.UUID `gorm:"type:char(36)" json:"changed_by,omitempty"`
	ChangedAt       *time.Time `gorm:"type:timestamp" json:"changed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (PayrollProfileChangeLog) TableName() string { return "payroll_profile_change_logs" }

func (p *PayrollProfileChangeLog) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// =============================================================================
// Read-only models untuk data dari modul lain (dipakai saat kalkulasi payroll
// run: snapshot employee, employment, position, grading). Bukan entity payroll.
// =============================================================================

// EmployeeRead menyimpan data minimal employee untuk snapshot payroll run.
type EmployeeRead struct {
	ID         uuid.UUID `gorm:"column:id;type:char(36);primaryKey" json:"id"`
	EmployeeID string    `gorm:"column:employee_id;type:varchar(50)" json:"employee_id"` // kode/nomor pegawai
	Name       string    `gorm:"column:name;type:varchar(255)" json:"name"`
	Status     string    `gorm:"column:status;type:varchar(20)" json:"status"`
}

func (EmployeeRead) TableName() string { return "employees" }

// EmploymentRead menyimpan employment aktif employee untuk resolusi position/grading.
type EmploymentRead struct {
	ID              uuid.UUID  `gorm:"column:id;type:char(36);primaryKey" json:"id"`
	EmployeeID      *uuid.UUID `gorm:"column:employee_id;type:char(36)" json:"employee_id,omitempty"`
	PositionID      *uuid.UUID `gorm:"column:position_id;type:char(36)" json:"position_id,omitempty"`
	OrganizationID  *uuid.UUID `gorm:"column:organization_id;type:char(36)" json:"organization_id,omitempty"`
	EffectiveDate   string     `gorm:"column:effective_date;type:date" json:"effective_date"`
	EffectiveEndDate *string   `gorm:"column:effective_end_date;type:date" json:"effective_end_date,omitempty"`
}

func (EmploymentRead) TableName() string { return "employments" }

// PositionRead menyimpan data position untuk snapshot payroll run.
type PositionRead struct {
	ID        uuid.UUID  `gorm:"column:id;type:char(36);primaryKey" json:"id"`
	Title     string     `gorm:"column:title;type:varchar(200)" json:"title"`
	GradingID *uuid.UUID `gorm:"column:grading_id;type:char(36)" json:"grading_id,omitempty"`
}

func (PositionRead) TableName() string { return "positions" }

// GradingRead menyimpan data grading untuk snapshot payroll run.
type GradingRead struct {
	ID   uuid.UUID `gorm:"column:id;type:char(36);primaryKey" json:"id"`
	Code string    `gorm:"column:code;type:varchar(20)" json:"code"`
	Name string    `gorm:"column:name;type:varchar(255)" json:"name"`
}

func (GradingRead) TableName() string { return "gradings" }

// =============================================================================
// Read models dari modul Workforce (attendance/leave) — payroll hanya
// mengonsumsi HASIL FINAL, tidak menghitung ulang (docs/payroll/06 §20-22).
// =============================================================================

// AttendanceSessionRead membaca baris attendance_sessions untuk summary
// kehadiran per employee per periode.
type AttendanceSessionRead struct {
	ID              uuid.UUID  `gorm:"column:id;type:char(36);primaryKey" json:"id"`
	EmployeeID      uuid.UUID  `gorm:"column:employee_id;type:char(36)" json:"employee_id"`
	WorkDate        string     `gorm:"column:work_date;type:date" json:"work_date"`
	Status          string     `gorm:"column:status;type:varchar(255)" json:"status"`
	OvertimeMinutes int        `gorm:"column:overtime_minutes" json:"overtime_minutes"`
	DeletedAt       *time.Time `gorm:"column:deleted_at;type:timestamp" json:"-"`
}

func (AttendanceSessionRead) TableName() string { return "attendance_sessions" }

// LeaveRequestRead membaca status leave_requests (untuk filter APPROVED_FINAL
// pada detail cuti).
type LeaveRequestRead struct {
	ID         uuid.UUID  `gorm:"column:id;type:char(36);primaryKey" json:"id"`
	EmployeeID uuid.UUID  `gorm:"column:employee_id;type:char(36)" json:"employee_id"`
	Status     string     `gorm:"column:status;type:varchar(255)" json:"status"`
	DeletedAt  *time.Time `gorm:"column:deleted_at;type:timestamp" json:"-"`
}

func (LeaveRequestRead) TableName() string { return "leave_requests" }

// LeaveRequestDetailRead membaca baris leave_request_details (hari per cuti,
// dengan flag is_paid & day_fraction).
type LeaveRequestDetailRead struct {
	LeaveRequestID uuid.UUID  `gorm:"column:leave_request_id;type:char(36)" json:"leave_request_id"`
	EmployeeID     uuid.UUID  `gorm:"column:employee_id;type:char(36)" json:"employee_id"`
	LeaveDate      string     `gorm:"column:leave_date;type:date" json:"leave_date"`
	DayFraction    float64    `gorm:"column:day_fraction;type:decimal(4,2)" json:"day_fraction"`
	IsPaid         bool       `gorm:"column:is_paid" json:"is_paid"`
	DeletedAt      *time.Time `gorm:"column:deleted_at;type:timestamp" json:"-"`
}

func (LeaveRequestDetailRead) TableName() string { return "leave_request_details" }
