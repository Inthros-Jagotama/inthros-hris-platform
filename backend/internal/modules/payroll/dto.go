package payroll

import "time"

// =============================================================================
// Request DTOs — Salary Components
// =============================================================================

type CreateSalaryComponentRequest struct {
	Code                   string  `json:"code" binding:"required,max=50"`
	Name                   string  `json:"name" binding:"required,max=150"`
	Description            *string `json:"description"`
	ComponentType          string  `json:"component_type" binding:"required,oneof=EARNING DEDUCTION EMPLOYER_CONTRIBUTION INFORMATION"`
	CalculationType        string  `json:"calculation_type" binding:"omitempty,oneof=FIXED PERCENTAGE FORMULA REFERENCE MANUAL"`
	Formula                *string `json:"formula"`
	ReferenceComponentID   *string `json:"reference_component_id" binding:"omitempty,uuid"`
	IsTaxable              *bool   `json:"is_taxable"`
	IsBpjsBase             *bool   `json:"is_bpjs_base"`
	IsRecurring            *bool   `json:"is_recurring"`
	IsProratable           *bool   `json:"is_proratable"`
	PrintOnSalaryStructure *bool   `json:"print_on_salary_structure"`
	DisplayOrder           *int    `json:"display_order"`
	Status                 *string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

type UpdateSalaryComponentRequest struct {
	Name                   *string `json:"name" binding:"omitempty,max=150"`
	Description            *string `json:"description"`
	ComponentType          *string `json:"component_type" binding:"omitempty,oneof=EARNING DEDUCTION EMPLOYER_CONTRIBUTION INFORMATION"`
	CalculationType        *string `json:"calculation_type" binding:"omitempty,oneof=FIXED PERCENTAGE FORMULA REFERENCE MANUAL"`
	Formula                *string `json:"formula"`
	ReferenceComponentID   *string `json:"reference_component_id" binding:"omitempty,uuid"`
	IsTaxable              *bool   `json:"is_taxable"`
	IsBpjsBase             *bool   `json:"is_bpjs_base"`
	IsRecurring            *bool   `json:"is_recurring"`
	IsProratable           *bool   `json:"is_proratable"`
	PrintOnSalaryStructure *bool   `json:"print_on_salary_structure"`
	DisplayOrder           *int    `json:"display_order"`
	Status                 *string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

// =============================================================================
// Request DTOs — Payroll Periods
// =============================================================================

type CreatePayrollPeriodRequest struct {
	PeriodYear  int    `json:"period_year" binding:"required"`
	PeriodMonth int    `json:"period_month" binding:"required,min=1,max=12"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"required"`
	AsOfDate    string `json:"as_of_date" binding:"required"`
	Status      string `json:"status" binding:"omitempty,oneof=OPEN CLOSED"`
}

type UpdatePayrollPeriodRequest struct {
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
	AsOfDate  *string `json:"as_of_date"`
	Status    *string `json:"status" binding:"omitempty,oneof=OPEN CLOSED"`
}

// =============================================================================
// Request DTOs — Employee Payroll Profiles
// =============================================================================

type CreateEmployeePayrollProfileRequest struct {
	EmployeeID         string  `json:"employee_id" binding:"required"`
	EmploymentID       *string `json:"employment_id"`
	PayrollGroupCode   string  `json:"payroll_group_code" binding:"required,max=50"`
	PayrollFrequency   string  `json:"payroll_frequency" binding:"omitempty,oneof=MONTHLY WEEKLY DAILY"`
	PaymentMethod      string  `json:"payment_method" binding:"omitempty,oneof=BANK_TRANSFER CASH CHEQUE"`
	SalaryCurrency     string  `json:"salary_currency" binding:"omitempty,len=3"`
	IsPayrollActive    *bool   `json:"is_payroll_active"`
	EffectiveStartDate string  `json:"effective_start_date" binding:"required"`
	EffectiveEndDate   *string `json:"effective_end_date"`
	Status             string  `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
	Notes              *string `json:"notes"`
}

type UpdateEmployeePayrollProfileRequest struct {
	PayrollGroupCode   *string `json:"payroll_group_code" binding:"omitempty,max=50"`
	PayrollFrequency   *string `json:"payroll_frequency" binding:"omitempty,oneof=MONTHLY WEEKLY DAILY"`
	PaymentMethod      *string `json:"payment_method" binding:"omitempty,oneof=BANK_TRANSFER CASH CHEQUE"`
	SalaryCurrency     *string `json:"salary_currency" binding:"omitempty,len=3"`
	IsPayrollActive    *bool   `json:"is_payroll_active"`
	EffectiveStartDate *string `json:"effective_start_date"`
	EffectiveEndDate   *string `json:"effective_end_date"`
	Status             *string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
	Notes              *string `json:"notes"`
}

// =============================================================================
// Request DTOs — Employee Bank Profiles
// =============================================================================

type CreateEmployeeBankProfileRequest struct {
	EmployeeID              string  `json:"employee_id" binding:"required"`
	EmployeePayrollProfileID string `json:"employee_payroll_profile_id" binding:"required"`
	BankCode                *string `json:"bank_code"`
	BankName                string  `json:"bank_name" binding:"required,max=150"`
	BankBranch              *string `json:"bank_branch"`
	BankAccountNumber       string  `json:"bank_account_number" binding:"required,max=100"`
	BankAccountHolderName   string  `json:"bank_account_holder_name" binding:"required,max=255"`
	IsPrimary               *bool   `json:"is_primary"`
	EffectiveStartDate      string  `json:"effective_start_date" binding:"required"`
	EffectiveEndDate        *string `json:"effective_end_date"`
	Status                  string  `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

type UpdateEmployeeBankProfileRequest struct {
	BankCode                *string `json:"bank_code"`
	BankName                *string `json:"bank_name" binding:"omitempty,max=150"`
	BankBranch              *string `json:"bank_branch"`
	BankAccountNumber       *string `json:"bank_account_number" binding:"omitempty,max=100"`
	BankAccountHolderName   *string `json:"bank_account_holder_name" binding:"omitempty,max=255"`
	IsPrimary               *bool   `json:"is_primary"`
	EffectiveStartDate      *string `json:"effective_start_date"`
	EffectiveEndDate        *string `json:"effective_end_date"`
	Status                  *string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

// =============================================================================
// Request DTOs — Employee BPJS Profiles
// =============================================================================

type CreateEmployeeBpjsProfileRequest struct {
	EmployeeID               string  `json:"employee_id" binding:"required"`
	EmployeePayrollProfileID string  `json:"employee_payroll_profile_id" binding:"required"`
	BpjsHealthActive         *bool   `json:"bpjs_health_active"`
	BpjsHealthNo             *string `json:"bpjs_health_no"`
	BpjsHealthRegisteredName *string `json:"bpjs_health_registered_name"`
	BpjsTkActive             *bool   `json:"bpjs_tk_active"`
	BpjsTkNo                 *string `json:"bpjs_tk_no"`
	BpjsTkRegisteredName     *string `json:"bpjs_tk_registered_name"`
	JkkRiskClass             *string `json:"jkk_risk_class" binding:"omitempty,oneof=VERY_LOW LOW MEDIUM HIGH VERY_HIGH"`
	PensionActive            *bool   `json:"pension_active"`
	EffectiveStartDate       string  `json:"effective_start_date" binding:"required"`
	EffectiveEndDate         *string `json:"effective_end_date"`
	Status                   string  `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
	Notes                    *string `json:"notes"`
}

type UpdateEmployeeBpjsProfileRequest struct {
	BpjsHealthActive         *bool   `json:"bpjs_health_active"`
	BpjsHealthNo             *string `json:"bpjs_health_no"`
	BpjsHealthRegisteredName *string `json:"bpjs_health_registered_name"`
	BpjsTkActive             *bool   `json:"bpjs_tk_active"`
	BpjsTkNo                 *string `json:"bpjs_tk_no"`
	BpjsTkRegisteredName     *string `json:"bpjs_tk_registered_name"`
	JkkRiskClass             *string `json:"jkk_risk_class" binding:"omitempty,oneof=VERY_LOW LOW MEDIUM HIGH VERY_HIGH"`
	PensionActive            *bool   `json:"pension_active"`
	EffectiveStartDate       *string `json:"effective_start_date"`
	EffectiveEndDate         *string `json:"effective_end_date"`
	Status                   *string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

// =============================================================================
// Request DTOs — Employee Tax Profiles
// =============================================================================

type CreateEmployeeTaxProfileRequest struct {
	EmployeeID               string  `json:"employee_id" binding:"required"`
	EmployeePayrollProfileID string  `json:"employee_payroll_profile_id" binding:"required"`
	Npwp                     *string `json:"npwp"`
	NpwpRegisteredName       *string `json:"npwp_registered_name"`
	PtkpStatus               *string `json:"ptkp_status"`
	TaxMethod                string  `json:"tax_method" binding:"omitempty,oneof=GROSS GROSS_UP NETT"`
	IsTaxable                *bool   `json:"is_taxable"`
	HasNpwp                  *bool   `json:"has_npwp"`
	EffectiveStartDate       string  `json:"effective_start_date" binding:"required"`
	EffectiveEndDate         *string `json:"effective_end_date"`
	Status                   string  `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
	Notes                    *string `json:"notes"`
}

type UpdateEmployeeTaxProfileRequest struct {
	Npwp               *string `json:"npwp"`
	NpwpRegisteredName *string `json:"npwp_registered_name"`
	PtkpStatus         *string `json:"ptkp_status"`
	TaxMethod          *string `json:"tax_method" binding:"omitempty,oneof=GROSS GROSS_UP NETT"`
	IsTaxable          *bool   `json:"is_taxable"`
	HasNpwp            *bool   `json:"has_npwp"`
	EffectiveStartDate *string `json:"effective_start_date"`
	EffectiveEndDate   *string `json:"effective_end_date"`
	Status             *string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

// =============================================================================
// Request DTOs — BPJS Settings & Rates
// =============================================================================

type CreateBpjsSettingRequest struct {
	SettingCode           string   `json:"setting_code" binding:"required,max=50"`
	SettingName           string   `json:"setting_name" binding:"required,max=150"`
	BaseSource            string   `json:"base_source" binding:"omitempty,oneof=BPJS_BASE_COMPONENTS BASIC_SALARY GROSS_EARNING"`
	HealthMaxBaseAmount   *float64 `json:"health_max_base_amount"`
	PensionMaxBaseAmount  *float64 `json:"pension_max_base_amount"`
	DefaultJkkRiskClass   string   `json:"default_jkk_risk_class" binding:"omitempty,oneof=VERY_LOW LOW MEDIUM HIGH VERY_HIGH"`
	RoundingMode          string   `json:"rounding_mode" binding:"omitempty,oneof=NONE ROUND CEIL FLOOR"`
	EffectiveStartDate    string   `json:"effective_start_date" binding:"required"`
	EffectiveEndDate      *string  `json:"effective_end_date"`
	Status                string   `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
	Notes                 *string  `json:"notes"`
}

type CreateBpjsRateComponentRequest struct {
	BpjsSettingID           string   `json:"bpjs_setting_id" binding:"required"`
	RateCode                string   `json:"rate_code" binding:"required,max=80"`
	RateName                string   `json:"rate_name" binding:"required,max=180"`
	BpjsProgram             string   `json:"bpjs_program" binding:"required,oneof=HEALTH JHT JP JKK JKM JKP"`
	PaidBy                  string   `json:"paid_by" binding:"required,oneof=EMPLOYEE EMPLOYER"`
	SalaryComponentID       *string  `json:"salary_component_id"`
	RatePercent             float64  `json:"rate_percent" binding:"min=0"`
	FixedAmount             *float64 `json:"fixed_amount"`
	MinBaseAmount           *float64 `json:"min_base_amount"`
	MaxBaseAmount           *float64 `json:"max_base_amount"`
	JkkRiskClass            *string  `json:"jkk_risk_class" binding:"omitempty,oneof=VERY_LOW LOW MEDIUM HIGH VERY_HIGH"`
	IsEmployeeDeduction     *bool    `json:"is_employee_deduction"`
	IsEmployerContribution  *bool    `json:"is_employer_contribution"`
	GenerateToPayrollItem   *bool    `json:"generate_to_payroll_item"`
	PrintOnPayslip          *bool    `json:"print_on_payslip"`
	DisplayOrder            *int     `json:"display_order"`
	EffectiveStartDate      string   `json:"effective_start_date" binding:"required"`
	EffectiveEndDate        *string  `json:"effective_end_date"`
	Status                  string   `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

// =============================================================================
// Request DTOs — PPh21 Settings & Rates
// =============================================================================

type UpdateBpjsSettingRequest struct {
	SettingCode           *string  `json:"setting_code" binding:"omitempty,max=50"`
	SettingName           *string  `json:"setting_name" binding:"omitempty,max=150"`
	BaseSource            *string  `json:"base_source" binding:"omitempty,oneof=BPJS_BASE_COMPONENTS BASIC_SALARY GROSS_EARNING"`
	HealthMaxBaseAmount   *float64 `json:"health_max_base_amount"`
	PensionMaxBaseAmount  *float64 `json:"pension_max_base_amount"`
	DefaultJkkRiskClass   *string  `json:"default_jkk_risk_class" binding:"omitempty,oneof=VERY_LOW LOW MEDIUM HIGH VERY_HIGH"`
	RoundingMode          *string  `json:"rounding_mode" binding:"omitempty,oneof=NONE ROUND CEIL FLOOR"`
	EffectiveStartDate    *string  `json:"effective_start_date"`
	EffectiveEndDate      *string  `json:"effective_end_date"`
	Status                *string  `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
	Notes                 *string  `json:"notes"`
}

type UpdateBpjsRateComponentRequest struct {
	BpjsSettingID           *string  `json:"bpjs_setting_id"`
	RateCode                *string  `json:"rate_code" binding:"omitempty,max=80"`
	RateName                *string  `json:"rate_name" binding:"omitempty,max=180"`
	BpjsProgram             *string  `json:"bpjs_program" binding:"omitempty,oneof=HEALTH JHT JP JKK JKM JKP"`
	PaidBy                  *string  `json:"paid_by" binding:"omitempty,oneof=EMPLOYEE EMPLOYER"`
	SalaryComponentID       *string  `json:"salary_component_id"`
	RatePercent             *float64 `json:"rate_percent"`
	FixedAmount             *float64 `json:"fixed_amount"`
	MinBaseAmount           *float64 `json:"min_base_amount"`
	MaxBaseAmount           *float64 `json:"max_base_amount"`
	JkkRiskClass            *string  `json:"jkk_risk_class" binding:"omitempty,oneof=VERY_LOW LOW MEDIUM HIGH VERY_HIGH"`
	IsEmployeeDeduction     *bool    `json:"is_employee_deduction"`
	IsEmployerContribution  *bool    `json:"is_employer_contribution"`
	GenerateToPayrollItem   *bool    `json:"generate_to_payroll_item"`
	PrintOnPayslip          *bool    `json:"print_on_payslip"`
	DisplayOrder            *int     `json:"display_order"`
	EffectiveStartDate      *string  `json:"effective_start_date"`
	EffectiveEndDate        *string  `json:"effective_end_date"`
	Status                  *string  `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

type UpdatePph21SettingRequest struct {
	SettingCode                    *string  `json:"setting_code" binding:"omitempty,max=50"`
	SettingName                    *string  `json:"setting_name" binding:"omitempty,max=150"`
	CalculationMethod              *string  `json:"calculation_method" binding:"omitempty,oneof=REGULAR_GROSS_ANNUALIZED"`
	DefaultTaxMethod               *string  `json:"default_tax_method" binding:"omitempty,oneof=GROSS GROSS_UP NETT"`
	Pph21ComponentID               *string  `json:"pph21_component_id"`
	OccupationalExpenseRatePercent *float64 `json:"occupational_expense_rate_percent"`
	OccupationalExpenseMaxMonthly  *float64 `json:"occupational_expense_max_monthly"`
	OccupationalExpenseMaxYearly   *float64 `json:"occupational_expense_max_yearly"`
	DeductBpjsHealthEmployee       *bool    `json:"deduct_bpjs_health_employee"`
	DeductBpjsJhtEmployee          *bool    `json:"deduct_bpjs_jht_employee"`
	DeductBpjsJpEmployee           *bool    `json:"deduct_bpjs_jp_employee"`
	AnnualizationMonths            *int     `json:"annualization_months"`
	PkpRoundingUnit                *float64 `json:"pkp_rounding_unit"`
	NonNpwpMultiplierPercent       *float64 `json:"non_npwp_multiplier_percent"`
	RoundingMode                   *string  `json:"rounding_mode" binding:"omitempty,oneof=NONE ROUND CEIL FLOOR"`
	EffectiveStartDate             *string  `json:"effective_start_date"`
	EffectiveEndDate               *string  `json:"effective_end_date"`
	Status                         *string  `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

type CreatePph21SettingRequest struct {
	SettingCode                    string   `json:"setting_code" binding:"required,max=50"`
	SettingName                    string   `json:"setting_name" binding:"required,max=150"`
	CalculationMethod              string   `json:"calculation_method" binding:"omitempty,oneof=REGULAR_GROSS_ANNUALIZED"`
	DefaultTaxMethod               string   `json:"default_tax_method" binding:"omitempty,oneof=GROSS GROSS_UP NETT"`
	Pph21ComponentID               string   `json:"pph21_component_id" binding:"required"`
	OccupationalExpenseRatePercent *float64 `json:"occupational_expense_rate_percent"`
	OccupationalExpenseMaxMonthly  *float64 `json:"occupational_expense_max_monthly"`
	OccupationalExpenseMaxYearly   *float64 `json:"occupational_expense_max_yearly"`
	DeductBpjsHealthEmployee       *bool    `json:"deduct_bpjs_health_employee"`
	DeductBpjsJhtEmployee          *bool    `json:"deduct_bpjs_jht_employee"`
	DeductBpjsJpEmployee           *bool    `json:"deduct_bpjs_jp_employee"`
	AnnualizationMonths            *int     `json:"annualization_months"`
	PkpRoundingUnit                *float64 `json:"pkp_rounding_unit"`
	NonNpwpMultiplierPercent       *float64 `json:"non_npwp_multiplier_percent"`
	RoundingMode                   *string  `json:"rounding_mode" binding:"omitempty,oneof=NONE ROUND CEIL FLOOR"`
	EffectiveStartDate             string   `json:"effective_start_date" binding:"required"`
	EffectiveEndDate               *string  `json:"effective_end_date"`
	Status                         string   `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

type CreatePph21PtkpRateRequest struct {
	PtkpStatus         string  `json:"ptkp_status" binding:"required,max=20"`
	Description        *string `json:"description"`
	AnnualAmount       float64 `json:"annual_amount" binding:"required"`
	EffectiveStartDate string  `json:"effective_start_date" binding:"required"`
	EffectiveEndDate   *string `json:"effective_end_date"`
	Status             string  `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

type CreatePph21TaxBracketRequest struct {
	BracketOrder       int      `json:"bracket_order" binding:"required"`
	LowerBound         float64  `json:"lower_bound" binding:"required"`
	UpperBound         *float64 `json:"upper_bound"`
	RatePercent        float64  `json:"rate_percent" binding:"required"`
	EffectiveStartDate string   `json:"effective_start_date" binding:"required"`
	EffectiveEndDate   *string  `json:"effective_end_date"`
	Status             string   `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

// =============================================================================
// Request DTOs — Payroll Runs
// =============================================================================

type CreatePayrollRunRequest struct {
	PayrollPeriodID string   `json:"payroll_period_id" binding:"required"`
	RunCode         string   `json:"run_code" binding:"required,max=50"`
	RunType         string   `json:"run_type" binding:"omitempty,oneof=REGULAR OFF_CYCLE THR BONUS"`
	// ProrationMethod opsional: CALENDAR_DAYS (default) | WORKING_DAYS |
	// FIXED_30_DAYS | ATTENDANCE_DAYS. Validasi nilai di service.
	ProrationMethod string   `json:"proration_method" binding:"omitempty"`
	// EmployeeIDs opsional: daftar employee yang diikutkan. Kosong → semua
	// employee dengan payroll profile aktif saat kalkulasi dijalankan.
	EmployeeIDs []string `json:"employee_ids" binding:"omitempty"`
}

type UpdatePayrollRunStatusRequest struct {
	Status string  `json:"status" binding:"required,oneof=CALCULATED REVIEWED APPROVED LOCKED CANCELLED"`
	FlowID *string `json:"flow_id" binding:"omitempty,uuid"`
}

// UpdatePaymentStatusRequest — transisi status payment batch.
type UpdatePaymentStatusRequest struct {
	Status    string `json:"status" binding:"required,oneof=PENDING PROCESSING PAID FAILED REVERSED"`
	Reason    string `json:"reason,omitempty"`    // dipakai saat FAILED
	Reference string `json:"reference,omitempty"` // dipakai saat PAID/PROCESSING
}

// =============================================================================
// Response DTOs
// =============================================================================

type SalaryComponentResponse struct {
	ID                     string    `json:"id"`
	Code                   string    `json:"code"`
	Name                   string    `json:"name"`
	Description            string    `json:"description,omitempty"`
	ComponentType          string    `json:"component_type"`
	CalculationType        string    `json:"calculation_type"`
	Formula                string    `json:"formula,omitempty"`
	ReferenceComponentID   string    `json:"reference_component_id,omitempty"`
	IsTaxable              bool      `json:"is_taxable"`
	IsBpjsBase             bool      `json:"is_bpjs_base"`
	IsRecurring            bool      `json:"is_recurring"`
	IsProratable           bool      `json:"is_proratable"`
	PrintOnSalaryStructure bool      `json:"print_on_salary_structure"`
	DisplayOrder           int       `json:"display_order"`
	Status                 string    `json:"status"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type PayrollPeriodResponse struct {
	ID          string    `json:"id"`
	PeriodCode  string    `json:"period_code"`
	PeriodYear  int       `json:"period_year"`
	PeriodMonth int       `json:"period_month"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
	AsOfDate    string    `json:"as_of_date"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EmployeePayrollProfileResponse struct {
	ID                 string    `json:"id"`
	EmployeeID         string    `json:"employee_id"`
	EmploymentID       string    `json:"employment_id,omitempty"`
	PayrollGroupCode   string    `json:"payroll_group_code"`
	PayrollFrequency   string    `json:"payroll_frequency"`
	PaymentMethod      string    `json:"payment_method"`
	SalaryCurrency     string    `json:"salary_currency"`
	IsPayrollActive    bool      `json:"is_payroll_active"`
	EffectiveStartDate string    `json:"effective_start_date"`
	EffectiveEndDate   string    `json:"effective_end_date,omitempty"`
	Status             string    `json:"status"`
	Notes              string    `json:"notes,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type EmployeeBankProfileResponse struct {
	ID                        string    `json:"id"`
	EmployeeID                string    `json:"employee_id"`
	EmployeePayrollProfileID  string    `json:"employee_payroll_profile_id"`
	BankCode                  string    `json:"bank_code,omitempty"`
	BankName                  string    `json:"bank_name"`
	BankBranch                string    `json:"bank_branch,omitempty"`
	BankAccountNumber         string    `json:"bank_account_number"`
	BankAccountHolderName     string    `json:"bank_account_holder_name"`
	IsPrimary                 bool      `json:"is_primary"`
	EffectiveStartDate        string    `json:"effective_start_date"`
	EffectiveEndDate          string    `json:"effective_end_date,omitempty"`
	Status                    string    `json:"status"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

type EmployeeBpjsProfileResponse struct {
	ID                       string    `json:"id"`
	EmployeeID               string    `json:"employee_id"`
	EmployeePayrollProfileID string    `json:"employee_payroll_profile_id"`
	BpjsHealthActive         bool      `json:"bpjs_health_active"`
	BpjsHealthNo             string    `json:"bpjs_health_no,omitempty"`
	BpjsHealthRegisteredName string    `json:"bpjs_health_registered_name,omitempty"`
	BpjsTkActive             bool      `json:"bpjs_tk_active"`
	BpjsTkNo                 string    `json:"bpjs_tk_no,omitempty"`
	BpjsTkRegisteredName     string    `json:"bpjs_tk_registered_name,omitempty"`
	JkkRiskClass             string    `json:"jkk_risk_class"`
	PensionActive            bool      `json:"pension_active"`
	EffectiveStartDate       string    `json:"effective_start_date"`
	EffectiveEndDate         string    `json:"effective_end_date,omitempty"`
	Status                   string    `json:"status"`
	Notes                    string    `json:"notes,omitempty"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type EmployeeTaxProfileResponse struct {
	ID                       string    `json:"id"`
	EmployeeID               string    `json:"employee_id"`
	EmployeePayrollProfileID string    `json:"employee_payroll_profile_id"`
	Npwp                     string    `json:"npwp,omitempty"`
	NpwpRegisteredName       string    `json:"npwp_registered_name,omitempty"`
	PtkpStatus               string    `json:"ptkp_status,omitempty"`
	TaxMethod                string    `json:"tax_method"`
	IsTaxable                bool      `json:"is_taxable"`
	HasNpwp                  bool      `json:"has_npwp"`
	EffectiveStartDate       string    `json:"effective_start_date"`
	EffectiveEndDate         string    `json:"effective_end_date,omitempty"`
	Status                   string    `json:"status"`
	Notes                    string    `json:"notes,omitempty"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type BpjsSettingResponse struct {
	ID                    string    `json:"id"`
	SettingCode           string    `json:"setting_code"`
	SettingName           string    `json:"setting_name"`
	BaseSource            string    `json:"base_source"`
	HealthMaxBaseAmount   float64   `json:"health_max_base_amount,omitempty"`
	PensionMaxBaseAmount  float64   `json:"pension_max_base_amount,omitempty"`
	DefaultJkkRiskClass   string    `json:"default_jkk_risk_class"`
	RoundingMode          string    `json:"rounding_mode"`
	EffectiveStartDate    string    `json:"effective_start_date"`
	EffectiveEndDate      string    `json:"effective_end_date,omitempty"`
	Status                string    `json:"status"`
	Notes                 string    `json:"notes,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type BpjsRateComponentResponse struct {
	ID                      string    `json:"id"`
	BpjsSettingID           string    `json:"bpjs_setting_id"`
	RateCode                string    `json:"rate_code"`
	RateName                string    `json:"rate_name"`
	BpjsProgram             string    `json:"bpjs_program"`
	PaidBy                  string    `json:"paid_by"`
	SalaryComponentID       string    `json:"salary_component_id,omitempty"`
	RatePercent             float64   `json:"rate_percent"`
	FixedAmount             float64   `json:"fixed_amount,omitempty"`
	MinBaseAmount           float64   `json:"min_base_amount,omitempty"`
	MaxBaseAmount           float64   `json:"max_base_amount,omitempty"`
	JkkRiskClass            string    `json:"jkk_risk_class,omitempty"`
	IsEmployeeDeduction     bool      `json:"is_employee_deduction"`
	IsEmployerContribution  bool      `json:"is_employer_contribution"`
	GenerateToPayrollItem   bool      `json:"generate_to_payroll_item"`
	PrintOnPayslip          bool      `json:"print_on_payslip"`
	DisplayOrder            int       `json:"display_order"`
	EffectiveStartDate      string    `json:"effective_start_date"`
	EffectiveEndDate        string    `json:"effective_end_date,omitempty"`
	Status                  string    `json:"status"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type Pph21SettingResponse struct {
	ID                           string    `json:"id"`
	SettingCode                  string    `json:"setting_code"`
	SettingName                  string    `json:"setting_name"`
	CalculationMethod            string    `json:"calculation_method"`
	DefaultTaxMethod             string    `json:"default_tax_method"`
	Pph21ComponentID             string    `json:"pph21_component_id"`
	OccupationalExpenseRatePercent float64 `json:"occupational_expense_rate_percent"`
	OccupationalExpenseMaxMonthly float64  `json:"occupational_expense_max_monthly"`
	OccupationalExpenseMaxYearly float64   `json:"occupational_expense_max_yearly"`
	DeductBpjsHealthEmployee     bool      `json:"deduct_bpjs_health_employee"`
	DeductBpjsJhtEmployee        bool      `json:"deduct_bpjs_jht_employee"`
	DeductBpjsJpEmployee         bool      `json:"deduct_bpjs_jp_employee"`
	AnnualizationMonths          int       `json:"annualization_months"`
	PkpRoundingUnit              float64   `json:"pkp_rounding_unit"`
	NonNpwpMultiplierPercent     float64   `json:"non_npwp_multiplier_percent"`
	RoundingMode                 string    `json:"rounding_mode"`
	EffectiveStartDate           string    `json:"effective_start_date"`
	EffectiveEndDate             string    `json:"effective_end_date,omitempty"`
	Status                       string    `json:"status"`
	CreatedAt                    time.Time `json:"created_at"`
	UpdatedAt                    time.Time `json:"updated_at"`
}

type Pph21PtkpRateResponse struct {
	ID                 string    `json:"id"`
	PtkpStatus         string    `json:"ptkp_status"`
	Description        string    `json:"description,omitempty"`
	AnnualAmount       float64   `json:"annual_amount"`
	EffectiveStartDate string    `json:"effective_start_date"`
	EffectiveEndDate   string    `json:"effective_end_date,omitempty"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type Pph21TaxBracketResponse struct {
	ID                 string    `json:"id"`
	BracketOrder       int       `json:"bracket_order"`
	LowerBound         float64   `json:"lower_bound"`
	UpperBound         float64   `json:"upper_bound,omitempty"`
	RatePercent        float64   `json:"rate_percent"`
	EffectiveStartDate string    `json:"effective_start_date"`
	EffectiveEndDate   string    `json:"effective_end_date,omitempty"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type PayrollRunResponse struct {
	ID                    string     `json:"id"`
	PayrollPeriodID       string     `json:"payroll_period_id"`
	RunCode               string     `json:"run_code"`
	RunType               string     `json:"run_type"`
	ProrationMethod       string     `json:"proration_method"`
	Status                string     `json:"status"`
	TotalEmployees        int        `json:"total_employees"`
	TotalEarning          float64    `json:"total_earning"`
	TotalDeduction        float64    `json:"total_deduction"`
	TotalEmployerContribution float64 `json:"total_employer_contribution"`
	TotalNet              float64    `json:"total_net"`
	TotalCompanyCost      float64    `json:"total_company_cost"`
	CalculatedAt          *time.Time `json:"calculated_at,omitempty"`
	ReviewedAt            *time.Time `json:"reviewed_at,omitempty"`
	ApprovedAt            *time.Time `json:"approved_at,omitempty"`
	LockedAt              *time.Time `json:"locked_at,omitempty"`
	ApprovalInstanceID    string     `json:"approval_instance_id,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// PayslipItemResponse — satu baris komponen pada payslip (dari snapshot run).
type PayslipItemResponse struct {
	ComponentCode string  `json:"component_code"`
	ComponentName string  `json:"component_name"`
	ItemCategory  string  `json:"item_category"`
	Amount        float64 `json:"amount"`
}

// PayrollPayslipResponse — payslip + rincian earnings/deductions/contributions.
type PayrollPayslipResponse struct {
	ID                       string                `json:"id"`
	PayrollRunID             string                `json:"payroll_run_id"`
	EmployeeID               string                `json:"employee_id"`
	PayslipNumber            string                `json:"payslip_number"`
	PeriodCode               string                `json:"period_code"`
	PeriodYear               int                   `json:"period_year"`
	PeriodMonth              int                   `json:"period_month"`
	EmployeeCode             string                `json:"employee_code"`
	EmployeeName             string                `json:"employee_name"`
	PositionTitle            string                `json:"position_title,omitempty"`
	GradingName              string                `json:"grading_name,omitempty"`
	TotalEarning             float64               `json:"total_earning"`
	TotalDeduction           float64               `json:"total_deduction"`
	TotalEmployerContribution float64              `json:"total_employer_contribution"`
	NetAmount                float64               `json:"net_amount"`
	Status                   string                `json:"status"`
	GeneratedAt              *time.Time            `json:"generated_at,omitempty"`
	PublishedAt              *time.Time            `json:"published_at,omitempty"`
	CancelledAt              *time.Time            `json:"cancelled_at,omitempty"`
	Items                    []PayslipItemResponse `json:"items,omitempty"`
}

// PaymentBatchResponse — ringkasan batch pembayaran yang baru dibuat.
type PaymentBatchResponse struct {
	RunID       string  `json:"run_id"`
	Total       int     `json:"total"`
	TotalAmount float64 `json:"total_amount"`
	Skipped     int     `json:"skipped_no_bank_profile"`
	Status      string  `json:"status"`
}

// PayrollPaymentResponse — satu baris payment batch.
type PayrollPaymentResponse struct {
	ID                    string     `json:"id"`
	PayrollRunID          string     `json:"payroll_run_id"`
	EmployeeID            string     `json:"employee_id"`
	EmployeeCode          string     `json:"employee_code"`
	EmployeeName          string     `json:"employee_name"`
	Amount                float64    `json:"amount"`
	CurrencyCode          string     `json:"currency_code"`
	PaymentDate           string     `json:"payment_date"`
	BankName              string     `json:"bank_name,omitempty"`
	BankAccountNumber     string     `json:"bank_account_number"`
	BankAccountHolderName string     `json:"bank_account_holder_name"`
	Status                string     `json:"status"`
	Reference             string     `json:"reference,omitempty"`
	FailedReason          string     `json:"failed_reason,omitempty"`
	ProcessedAt           *time.Time `json:"processed_at,omitempty"`
	PaidAt                *time.Time `json:"paid_at,omitempty"`
	FailedAt              *time.Time `json:"failed_at,omitempty"`
	ReversedAt            *time.Time `json:"reversed_at,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
}

func toPayrollPaymentResponse(p *PayrollPayment) PayrollPaymentResponse {
	r := PayrollPaymentResponse{
		ID:                    p.ID.String(),
		PayrollRunID:          p.PayrollRunID.String(),
		EmployeeID:            p.EmployeeID.String(),
		EmployeeCode:          p.EmployeeCode,
		EmployeeName:          p.EmployeeName,
		Amount:                p.Amount,
		CurrencyCode:          p.CurrencyCode,
		PaymentDate:           p.PaymentDate,
		BankAccountNumber:     p.BankAccountNumber,
		BankAccountHolderName: p.BankAccountHolderName,
		Status:                p.Status,
		ProcessedAt:           p.ProcessedAt,
		PaidAt:                p.PaidAt,
		FailedAt:              p.FailedAt,
		ReversedAt:            p.ReversedAt,
		CreatedAt:             p.CreatedAt,
	}
	if p.BankName != nil {
		r.BankName = *p.BankName
	}
	if p.Reference != nil {
		r.Reference = *p.Reference
	}
	if p.FailedReason != nil {
		r.FailedReason = *p.FailedReason
	}
	return r
}

func toPayrollPayslipResponse(p *PayrollPayslip, items []PayrollRunItem) PayrollPayslipResponse {
	r := PayrollPayslipResponse{
		ID:                       p.ID.String(),
		PayrollRunID:             p.PayrollRunID.String(),
		EmployeeID:               p.EmployeeID.String(),
		PayslipNumber:            p.PayslipNumber,
		PeriodCode:               p.PeriodCode,
		PeriodYear:               p.PeriodYear,
		PeriodMonth:              p.PeriodMonth,
		EmployeeCode:             p.EmployeeCode,
		EmployeeName:             p.EmployeeName,
		TotalEarning:             p.TotalEarning,
		TotalDeduction:           p.TotalDeduction,
		TotalEmployerContribution: p.TotalEmployerContribution,
		NetAmount:                p.NetAmount,
		Status:                   p.Status,
		GeneratedAt:              p.GeneratedAt,
		PublishedAt:              p.PublishedAt,
		CancelledAt:              p.CancelledAt,
	}
	if p.PositionTitle != nil {
		r.PositionTitle = *p.PositionTitle
	}
	if p.GradingName != nil {
		r.GradingName = *p.GradingName
	}
	for _, it := range items {
		r.Items = append(r.Items, PayslipItemResponse{
			ComponentCode: it.ComponentCode,
			ComponentName: it.ComponentName,
			ItemCategory:  it.ItemCategory,
			Amount:        it.Amount,
		})
	}
	return r
}

type PayrollRunEmployeeResponse struct {
	ID                     string  `json:"id"`
	PayrollRunID           string  `json:"payroll_run_id"`
	EmployeeID             string  `json:"employee_id"`
	EmploymentID           string  `json:"employment_id,omitempty"`
	PositionID             string  `json:"position_id,omitempty"`
	GradingID              string  `json:"grading_id,omitempty"`
	EmployeeCode           string  `json:"employee_code"`
	EmployeeName           string  `json:"employee_name"`
	PositionTitle          string  `json:"position_title,omitempty"`
	GradingName            string  `json:"grading_name,omitempty"`
	TotalEarning           float64 `json:"total_earning"`
	TotalDeduction         float64 `json:"total_deduction"`
	TotalEmployerContribution float64 `json:"total_employer_contribution"`
	NetAmount              float64 `json:"net_amount"`
	TotalCompanyCost       float64 `json:"total_company_cost"`
	Status                 string  `json:"status"`
}

func toPayrollRunEmployeeResponse(e *PayrollRunEmployee) PayrollRunEmployeeResponse {
	r := PayrollRunEmployeeResponse{
		ID:                       e.ID.String(),
		PayrollRunID:             e.PayrollRunID.String(),
		EmployeeID:               e.EmployeeID.String(),
		EmployeeCode:             e.EmployeeCode,
		EmployeeName:             e.EmployeeName,
		TotalEarning:             e.TotalEarning,
		TotalDeduction:           e.TotalDeduction,
		TotalEmployerContribution: e.TotalEmployerContribution,
		NetAmount:                e.NetAmount,
		TotalCompanyCost:         e.TotalCompanyCost,
		Status:                   e.Status,
	}
	if e.EmploymentID != nil {
		r.EmploymentID = e.EmploymentID.String()
	}
	if e.PositionID != nil {
		r.PositionID = e.PositionID.String()
	}
	if e.GradingID != nil {
		r.GradingID = e.GradingID.String()
	}
	if e.PositionTitle != nil {
		r.PositionTitle = *e.PositionTitle
	}
	if e.GradingName != nil {
		r.GradingName = *e.GradingName
	}
	return r
}

type PayrollRunItemResponse struct {
	ID                   string   `json:"id"`
	PayrollRunID         string   `json:"payroll_run_id"`
	PayrollRunEmployeeID string   `json:"payroll_run_employee_id"`
	EmployeeID           string   `json:"employee_id"`
	SalaryComponentID    string   `json:"salary_component_id"`
	ComponentCode        string   `json:"component_code"`
	ComponentName        string   `json:"component_name"`
	ComponentType        string   `json:"component_type"`
	CalculationType      string   `json:"calculation_type"`
	ItemCategory         string   `json:"item_category"`
	PaidBy               string   `json:"paid_by"`
	AffectsGrossPay      bool     `json:"affects_gross_pay"`
	AffectsNetPay        bool     `json:"affects_net_pay"`
	AffectsCompanyCost   bool     `json:"affects_company_cost"`
	PrintOnPayslip       bool     `json:"print_on_payslip"`
	Amount               float64  `json:"amount"`
	BaseAmount           float64  `json:"base_amount"`
	Rate                 *float64 `json:"rate,omitempty"`
	Formula              *string  `json:"formula,omitempty"`
	FormulaResult        *float64 `json:"formula_result,omitempty"`
	CurrencyCode         string   `json:"currency_code"`
	SourceGroup          string   `json:"source_group"`
	Notes                *string  `json:"notes,omitempty"`
}

func toPayrollRunItemResponse(i *PayrollRunItem) PayrollRunItemResponse {
	r := PayrollRunItemResponse{
		ID:                   i.ID.String(),
		PayrollRunID:         i.PayrollRunID.String(),
		PayrollRunEmployeeID: i.PayrollRunEmployeeID.String(),
		EmployeeID:           i.EmployeeID.String(),
		SalaryComponentID:    i.SalaryComponentID.String(),
		ComponentCode:        i.ComponentCode,
		ComponentName:        i.ComponentName,
		ComponentType:        i.ComponentType,
		CalculationType:      i.CalculationType,
		ItemCategory:         i.ItemCategory,
		PaidBy:               i.PaidBy,
		AffectsGrossPay:      i.AffectsGrossPay,
		AffectsNetPay:        i.AffectsNetPay,
		AffectsCompanyCost:   i.AffectsCompanyCost,
		PrintOnPayslip:       i.PrintOnPayslip,
		Amount:               i.Amount,
		BaseAmount:           i.BaseAmount,
		Rate:                 i.Rate,
		Formula:              i.Formula,
		FormulaResult:        i.FormulaResult,
		CurrencyCode:         i.CurrencyCode,
		SourceGroup:          i.SourceGroup,
		Notes:                i.Notes,
	}
	return r
}

// =============================================================================
// Converter Helpers
// =============================================================================

func toSalaryComponentResponse(s *SalaryComponent) SalaryComponentResponse {
	r := SalaryComponentResponse{
		ID:              s.ID.String(),
		Code:            s.Code,
		Name:            s.Name,
		ComponentType:   s.ComponentType,
		CalculationType: s.CalculationType,
		IsTaxable:       s.IsTaxable,
		IsBpjsBase:      s.IsBpjsBase,
		IsRecurring:     s.IsRecurring,
		IsProratable:    s.IsProratable,
		DisplayOrder:    s.DisplayOrder,
		Status:          s.Status,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
	if s.Description != nil {
		r.Description = *s.Description
	}
	if s.Formula != nil {
		r.Formula = *s.Formula
	}
	if s.ReferenceComponentID != nil {
		r.ReferenceComponentID = s.ReferenceComponentID.String()
	}
	return r
}

func toPayrollPeriodResponse(p *PayrollPeriod) PayrollPeriodResponse {
	r := PayrollPeriodResponse{
		ID:          p.ID.String(),
		PeriodCode:  p.PeriodCode,
		PeriodYear:  p.PeriodYear,
		PeriodMonth: p.PeriodMonth,
		StartDate:   p.StartDate,
		EndDate:     p.EndDate,
		AsOfDate:    p.AsOfDate,
		Status:      p.Status,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	return r
}

func toEmployeePayrollProfileResponse(p *EmployeePayrollProfile) EmployeePayrollProfileResponse {
	r := EmployeePayrollProfileResponse{
		ID:                p.ID.String(),
		EmployeeID:        p.EmployeeID.String(),
		PayrollGroupCode:  p.PayrollGroupCode,
		PayrollFrequency:  p.PayrollFrequency,
		PaymentMethod:     p.PaymentMethod,
		SalaryCurrency:    p.SalaryCurrency,
		IsPayrollActive:   p.IsPayrollActive,
		EffectiveStartDate: p.EffectiveStartDate,
		Status:            p.Status,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}
	if p.EmploymentID != nil {
		r.EmploymentID = p.EmploymentID.String()
	}
	if p.EffectiveEndDate != nil {
		r.EffectiveEndDate = *p.EffectiveEndDate
	}
	if p.Notes != nil {
		r.Notes = *p.Notes
	}
	return r
}

func toEmployeeBankProfileResponse(b *EmployeeBankProfile) EmployeeBankProfileResponse {
	r := EmployeeBankProfileResponse{
		ID:                        b.ID.String(),
		EmployeeID:                b.EmployeeID.String(),
		EmployeePayrollProfileID:  b.EmployeePayrollProfileID.String(),
		BankName:                  b.BankName,
		BankAccountNumber:         b.BankAccountNumber,
		BankAccountHolderName:     b.BankAccountHolderName,
		IsPrimary:                 b.IsPrimary,
		EffectiveStartDate:        b.EffectiveStartDate,
		Status:                    b.Status,
		CreatedAt:                 b.CreatedAt,
		UpdatedAt:                 b.UpdatedAt,
	}
	if b.BankCode != nil {
		r.BankCode = *b.BankCode
	}
	if b.BankBranch != nil {
		r.BankBranch = *b.BankBranch
	}
	if b.EffectiveEndDate != nil {
		r.EffectiveEndDate = *b.EffectiveEndDate
	}
	return r
}

func toEmployeeBpjsProfileResponse(b *EmployeeBpjsProfile) EmployeeBpjsProfileResponse {
	r := EmployeeBpjsProfileResponse{
		ID:               b.ID.String(),
		EmployeeID:       b.EmployeeID.String(),
		EmployeePayrollProfileID: b.EmployeePayrollProfileID.String(),
		BpjsHealthActive: b.BpjsHealthActive,
		BpjsTkActive:     b.BpjsTkActive,
		JkkRiskClass:     b.JkkRiskClass,
		PensionActive:    b.PensionActive,
		EffectiveStartDate: b.EffectiveStartDate,
		Status:           b.Status,
		CreatedAt:        b.CreatedAt,
		UpdatedAt:        b.UpdatedAt,
	}
	if b.BpjsHealthNo != nil {
		r.BpjsHealthNo = *b.BpjsHealthNo
	}
	if b.BpjsHealthRegisteredName != nil {
		r.BpjsHealthRegisteredName = *b.BpjsHealthRegisteredName
	}
	if b.BpjsTkNo != nil {
		r.BpjsTkNo = *b.BpjsTkNo
	}
	if b.BpjsTkRegisteredName != nil {
		r.BpjsTkRegisteredName = *b.BpjsTkRegisteredName
	}
	if b.EffectiveEndDate != nil {
		r.EffectiveEndDate = *b.EffectiveEndDate
	}
	if b.Notes != nil {
		r.Notes = *b.Notes
	}
	return r
}

func toEmployeeTaxProfileResponse(t *EmployeeTaxProfile) EmployeeTaxProfileResponse {
	r := EmployeeTaxProfileResponse{
		ID:         t.ID.String(),
		EmployeeID: t.EmployeeID.String(),
		EmployeePayrollProfileID: t.EmployeePayrollProfileID.String(),
		TaxMethod:  t.TaxMethod,
		IsTaxable:  t.IsTaxable,
		HasNpwp:    t.HasNpwp,
		EffectiveStartDate: t.EffectiveStartDate,
		Status:     t.Status,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
	if t.Npwp != nil {
		r.Npwp = *t.Npwp
	}
	if t.NpwpRegisteredName != nil {
		r.NpwpRegisteredName = *t.NpwpRegisteredName
	}
	if t.PtkpStatus != nil {
		r.PtkpStatus = *t.PtkpStatus
	}
	if t.EffectiveEndDate != nil {
		r.EffectiveEndDate = *t.EffectiveEndDate
	}
	if t.Notes != nil {
		r.Notes = *t.Notes
	}
	return r
}

func toBpjsSettingResponse(b *BpjsSetting) BpjsSettingResponse {
	r := BpjsSettingResponse{
		ID:          b.ID.String(),
		SettingCode: b.SettingCode,
		SettingName: b.SettingName,
		BaseSource:  b.BaseSource,
		DefaultJkkRiskClass: b.DefaultJkkRiskClass,
		RoundingMode: b.RoundingMode,
		EffectiveStartDate: b.EffectiveStartDate,
		Status:      b.Status,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
	if b.HealthMaxBaseAmount != nil {
		r.HealthMaxBaseAmount = *b.HealthMaxBaseAmount
	}
	if b.PensionMaxBaseAmount != nil {
		r.PensionMaxBaseAmount = *b.PensionMaxBaseAmount
	}
	if b.EffectiveEndDate != nil {
		r.EffectiveEndDate = *b.EffectiveEndDate
	}
	if b.Notes != nil {
		r.Notes = *b.Notes
	}
	return r
}

func toBpjsRateComponentResponse(b *BpjsRateComponent) BpjsRateComponentResponse {
	r := BpjsRateComponentResponse{
		ID:         b.ID.String(),
		BpjsSettingID: b.BpjsSettingID.String(),
		RateCode:   b.RateCode,
		RateName:   b.RateName,
		BpjsProgram: b.BpjsProgram,
		PaidBy:     b.PaidBy,
		RatePercent: b.RatePercent,
		IsEmployeeDeduction: b.IsEmployeeDeduction,
		IsEmployerContribution: b.IsEmployerContribution,
		GenerateToPayrollItem: b.GenerateToPayrollItem,
		PrintOnPayslip: b.PrintOnPayslip,
		DisplayOrder: b.DisplayOrder,
		EffectiveStartDate: b.EffectiveStartDate,
		Status:      b.Status,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
	if b.SalaryComponentID != nil {
		r.SalaryComponentID = b.SalaryComponentID.String()
	}
	if b.FixedAmount != nil {
		r.FixedAmount = *b.FixedAmount
	}
	if b.MinBaseAmount != nil {
		r.MinBaseAmount = *b.MinBaseAmount
	}
	if b.MaxBaseAmount != nil {
		r.MaxBaseAmount = *b.MaxBaseAmount
	}
	if b.JkkRiskClass != nil {
		r.JkkRiskClass = *b.JkkRiskClass
	}
	if b.EffectiveEndDate != nil {
		r.EffectiveEndDate = *b.EffectiveEndDate
	}
	return r
}

func toPph21SettingResponse(p *Pph21Setting) Pph21SettingResponse {
	r := Pph21SettingResponse{
		ID:              p.ID.String(),
		SettingCode:     p.SettingCode,
		SettingName:     p.SettingName,
		CalculationMethod: p.CalculationMethod,
		DefaultTaxMethod: p.DefaultTaxMethod,
		Pph21ComponentID: p.Pph21ComponentID.String(),
		OccupationalExpenseRatePercent: p.OccupationalExpenseRatePercent,
		OccupationalExpenseMaxMonthly:  p.OccupationalExpenseMaxMonthly,
		OccupationalExpenseMaxYearly:   p.OccupationalExpenseMaxYearly,
		DeductBpjsHealthEmployee: p.DeductBpjsHealthEmployee,
		DeductBpjsJhtEmployee:    p.DeductBpjsJhtEmployee,
		DeductBpjsJpEmployee:     p.DeductBpjsJpEmployee,
		AnnualizationMonths:      p.AnnualizationMonths,
		PkpRoundingUnit:          p.PkpRoundingUnit,
		NonNpwpMultiplierPercent: p.NonNpwpMultiplierPercent,
		RoundingMode:             p.RoundingMode,
		EffectiveStartDate:       p.EffectiveStartDate,
		Status:                   p.Status,
		CreatedAt:                p.CreatedAt,
		UpdatedAt:                p.UpdatedAt,
	}
	if p.EffectiveEndDate != nil {
		r.EffectiveEndDate = *p.EffectiveEndDate
	}
	return r
}

func toPph21PtkpRateResponse(p *Pph21PtkpRate) Pph21PtkpRateResponse {
	r := Pph21PtkpRateResponse{
		ID:         p.ID.String(),
		PtkpStatus: p.PtkpStatus,
		AnnualAmount: p.AnnualAmount,
		EffectiveStartDate: p.EffectiveStartDate,
		Status:      p.Status,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	if p.Description != nil {
		r.Description = *p.Description
	}
	if p.EffectiveEndDate != nil {
		r.EffectiveEndDate = *p.EffectiveEndDate
	}
	return r
}

func toPph21TaxBracketResponse(p *Pph21TaxBracket) Pph21TaxBracketResponse {
	r := Pph21TaxBracketResponse{
		ID:           p.ID.String(),
		BracketOrder: p.BracketOrder,
		LowerBound:   p.LowerBound,
		RatePercent:  p.RatePercent,
		EffectiveStartDate: p.EffectiveStartDate,
		Status:       p.Status,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
	if p.UpperBound != nil {
		r.UpperBound = *p.UpperBound
	}
	if p.EffectiveEndDate != nil {
		r.EffectiveEndDate = *p.EffectiveEndDate
	}
	return r
}

func toPayrollRunResponse(p *PayrollRun) PayrollRunResponse {
	r := PayrollRunResponse{
		ID:              p.ID.String(),
		PayrollPeriodID: p.PayrollPeriodID.String(),
		RunCode:         p.RunCode,
		RunType:         p.RunType,
		ProrationMethod: p.ProrationMethod,
		Status:          p.Status,
		TotalEmployees:  p.TotalEmployees,
		TotalEarning:    p.TotalEarning,
		TotalDeduction:  p.TotalDeduction,
		TotalEmployerContribution: p.TotalEmployerContribution,
		TotalNet:        p.TotalNet,
		TotalCompanyCost: p.TotalCompanyCost,
		CalculatedAt:    p.CalculatedAt,
		ReviewedAt:      p.ReviewedAt,
		ApprovedAt:      p.ApprovedAt,
		LockedAt:        p.LockedAt,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
	if p.ApprovalInstanceID != nil {
		r.ApprovalInstanceID = p.ApprovalInstanceID.String()
	}
	return r
}

// =============================================================================
// Generic Paginated Response
// =============================================================================

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// StringPtr is a helper to convert string to *string
func StringPtr(s string) *string {
	return &s
}

// BoolPtr is a helper to convert bool to *bool
func BoolPtr(b bool) *bool {
	return &b
}

// IntPtr is a helper to convert int to *int
func IntPtr(i int) *int {
	return &i
}

// Float64Ptr is a helper to convert float64 to *float64
func Float64Ptr(f float64) *float64 {
	return &f
}
