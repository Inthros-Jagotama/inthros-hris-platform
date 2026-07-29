package payroll

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

// ApprovalEngine defines the interface for approval operations needed by Payroll.
// This avoids circular dependencies between payroll and approval modules.
type ApprovalEngine interface {
	// CreateApprovalInstance creates an approval instance for a document.
	// Returns the instance ID if successful.
	CreateApprovalInstance(ctx context.Context, module, documentID, flowID string) (string, error)

	// GetApprovalInstanceStatus returns the status of an approval instance.
	GetApprovalInstanceStatus(ctx context.Context, instanceID string) (string, error)
}

type Service struct {
	repo            *Repository
	logger          *zap.Logger
	approvalEngine  ApprovalEngine
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// SetApprovalEngine injects the approval engine after service creation.
func (s *Service) SetApprovalEngine(ae ApprovalEngine) {
	s.approvalEngine = ae
}

// =============================================================================
// Salary Components
// =============================================================================

func (s *Service) CreateSalaryComponent(ctx context.Context, req CreateSalaryComponentRequest) (*SalaryComponentResponse, error) {
	sc := &SalaryComponent{
		Code:                   req.Code,
		Name:                   req.Name,
		ComponentType:          req.ComponentType,
		CalculationType:        "FIXED",
		IsTaxable:              true,
		IsBpjsBase:             false,
		IsRecurring:            true,
		IsProratable:           true,
		PrintOnSalaryStructure: true,
		DisplayOrder:           100,
		Status:                 "ACTIVE",
	}
	if req.Description != nil {
		sc.Description = req.Description
	}
	if req.CalculationType != "" {
		sc.CalculationType = req.CalculationType
	}
	if req.IsTaxable != nil {
		sc.IsTaxable = *req.IsTaxable
	}
	if req.IsBpjsBase != nil {
		sc.IsBpjsBase = *req.IsBpjsBase
	}
	if req.IsRecurring != nil {
		sc.IsRecurring = *req.IsRecurring
	}
	if req.IsProratable != nil {
		sc.IsProratable = *req.IsProratable
	}
	if req.PrintOnSalaryStructure != nil {
		sc.PrintOnSalaryStructure = *req.PrintOnSalaryStructure
	}
	if req.DisplayOrder != nil {
		sc.DisplayOrder = *req.DisplayOrder
	}
	if req.Status != nil {
		sc.Status = *req.Status
	}
	sc.CreatedBy = authctx.GetUserID(ctx)
	sc.UpdatedBy = sc.CreatedBy

	if err := s.repo.CreateSalaryComponent(ctx, sc); err != nil {
		return nil, err
	}

	s.logger.Info("Salary component created", zap.String("code", sc.Code), zap.String("name", sc.Name))
	response := toSalaryComponentResponse(sc)
	return &response, nil
}

func (s *Service) GetSalaryComponentByID(ctx context.Context, id string) (*SalaryComponentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	sc, err := s.repo.FindSalaryComponentByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := toSalaryComponentResponse(sc)
	return &response, nil
}

func (s *Service) ListSalaryComponents(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	items, total, err := s.repo.FindAllSalaryComponents(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	var responses []SalaryComponentResponse
	for _, item := range items {
		responses = append(responses, toSalaryComponentResponse(&item))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: totalPages,
	}, nil
}

func (s *Service) UpdateSalaryComponent(ctx context.Context, id string, req UpdateSalaryComponentRequest) (*SalaryComponentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	sc, err := s.repo.FindSalaryComponentByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	sc.UpdatedBy = authctx.GetUserID(ctx)
	if req.Name != nil {
		sc.Name = *req.Name
	}
	if req.Description != nil {
		sc.Description = req.Description
	}
	if req.ComponentType != nil {
		sc.ComponentType = *req.ComponentType
	}
	if req.CalculationType != nil {
		sc.CalculationType = *req.CalculationType
	}
	if req.IsTaxable != nil {
		sc.IsTaxable = *req.IsTaxable
	}
	if req.IsBpjsBase != nil {
		sc.IsBpjsBase = *req.IsBpjsBase
	}
	if req.IsRecurring != nil {
		sc.IsRecurring = *req.IsRecurring
	}
	if req.IsProratable != nil {
		sc.IsProratable = *req.IsProratable
	}
	if req.PrintOnSalaryStructure != nil {
		sc.PrintOnSalaryStructure = *req.PrintOnSalaryStructure
	}
	if req.DisplayOrder != nil {
		sc.DisplayOrder = *req.DisplayOrder
	}
	if req.Status != nil {
		sc.Status = *req.Status
	}
	if err := s.repo.UpdateSalaryComponent(ctx, sc); err != nil {
		return nil, err
	}
	response := toSalaryComponentResponse(sc)
	return &response, nil
}

func (s *Service) DeleteSalaryComponent(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteSalaryComponent(ctx, uid)
}

// =============================================================================
// Payroll Periods
// =============================================================================

func (s *Service) CreatePayrollPeriod(ctx context.Context, req CreatePayrollPeriodRequest) (*PayrollPeriodResponse, error) {
	periodCode := fmt.Sprintf("%d%02d", req.PeriodYear, req.PeriodMonth)
	status := req.Status
	if status == "" {
		status = "OPEN"
	}
	p := &PayrollPeriod{
		PeriodCode:  periodCode,
		PeriodYear:  req.PeriodYear,
		PeriodMonth: req.PeriodMonth,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		AsOfDate:    req.AsOfDate,
		Status:      status,
	}
	p.CreatedBy = authctx.GetUserID(ctx)
	p.UpdatedBy = p.CreatedBy
	if err := s.repo.CreatePayrollPeriod(ctx, p); err != nil {
		return nil, err
	}
	s.logger.Info("Payroll period created", zap.String("period", p.PeriodCode))
	response := toPayrollPeriodResponse(p)
	return &response, nil
}

func (s *Service) ListPayrollPeriods(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	items, total, err := s.repo.FindAllPayrollPeriods(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	var responses []PayrollPeriodResponse
	for _, item := range items {
		responses = append(responses, toPayrollPeriodResponse(&item))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: totalPages,
	}, nil
}

func (s *Service) UpdatePayrollPeriod(ctx context.Context, id string, req UpdatePayrollPeriodRequest) (*PayrollPeriodResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindPayrollPeriodByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	p.UpdatedBy = authctx.GetUserID(ctx)
	if req.StartDate != nil {
		p.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		p.EndDate = *req.EndDate
	}
	if req.AsOfDate != nil {
		p.AsOfDate = *req.AsOfDate
	}
	if req.Status != nil {
		p.Status = *req.Status
	}
	if err := s.repo.UpdatePayrollPeriod(ctx, p); err != nil {
		return nil, err
	}
	response := toPayrollPeriodResponse(p)
	return &response, nil
}

// =============================================================================
// Employee Payroll Profiles
// =============================================================================

func (s *Service) CreateEmployeePayrollProfile(ctx context.Context, req CreateEmployeePayrollProfileRequest) (*EmployeePayrollProfileResponse, error) {
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	p := &EmployeePayrollProfile{
		EmployeeID:         empID,
		PayrollGroupCode:   req.PayrollGroupCode,
		PayrollFrequency:   "MONTHLY",
		PaymentMethod:      "BANK_TRANSFER",
		SalaryCurrency:     "IDR",
		IsPayrollActive:    true,
		EffectiveStartDate: req.EffectiveStartDate,
		Status:             "ACTIVE",
	}
	if req.EmploymentID != nil && *req.EmploymentID != "" {
		id, _ := uuid.Parse(*req.EmploymentID)
		p.EmploymentID = &id
	}
	if req.PayrollFrequency != "" {
		p.PayrollFrequency = req.PayrollFrequency
	}
	if req.PaymentMethod != "" {
		p.PaymentMethod = req.PaymentMethod
	}
	if req.SalaryCurrency != "" {
		p.SalaryCurrency = req.SalaryCurrency
	}
	if req.IsPayrollActive != nil {
		p.IsPayrollActive = *req.IsPayrollActive
	}
	if req.EffectiveEndDate != nil {
		p.EffectiveEndDate = req.EffectiveEndDate
	}
	if req.Status != "" {
		p.Status = req.Status
	}
	if req.Notes != nil {
		p.Notes = req.Notes
	}
	p.CreatedBy = authctx.GetUserID(ctx)
	p.UpdatedBy = p.CreatedBy
	if err := s.repo.CreateEmployeePayrollProfile(ctx, p); err != nil {
		return nil, err
	}
	response := toEmployeePayrollProfileResponse(p)
	return &response, nil
}

func (s *Service) GetEmployeePayrollProfileByID(ctx context.Context, id string) (*EmployeePayrollProfileResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindEmployeePayrollProfileByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := toEmployeePayrollProfileResponse(p)
	return &response, nil
}

func (s *Service) ListEmployeePayrollProfiles(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	items, total, err := s.repo.FindAllEmployeePayrollProfiles(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	var responses []EmployeePayrollProfileResponse
	for _, item := range items {
		responses = append(responses, toEmployeePayrollProfileResponse(&item))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: totalPages,
	}, nil
}

func (s *Service) DeleteEmployeePayrollProfile(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteEmployeePayrollProfile(ctx, uid)
}

// =============================================================================
// Employee Bank Profiles
// =============================================================================

func (s *Service) GetEmployeeBankProfileByID(ctx context.Context, id string) (*EmployeeBankProfileResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	b, err := s.repo.FindEmployeeBankProfileByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := toEmployeeBankProfileResponse(b)
	return &response, nil
}

func (s *Service) UpdateEmployeeBankProfile(ctx context.Context, id string, req UpdateEmployeeBankProfileRequest) (*EmployeeBankProfileResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	b, err := s.repo.FindEmployeeBankProfileByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	b.UpdatedBy = authctx.GetUserID(ctx)
	if req.BankCode != nil {
		b.BankCode = req.BankCode
	}
	if req.BankName != nil {
		b.BankName = *req.BankName
	}
	if req.BankBranch != nil {
		b.BankBranch = req.BankBranch
	}
	if req.BankAccountNumber != nil {
		b.BankAccountNumber = *req.BankAccountNumber
	}
	if req.BankAccountHolderName != nil {
		b.BankAccountHolderName = *req.BankAccountHolderName
	}
	if req.IsPrimary != nil {
		b.IsPrimary = *req.IsPrimary
	}
	if req.EffectiveStartDate != nil {
		b.EffectiveStartDate = *req.EffectiveStartDate
	}
	if req.EffectiveEndDate != nil {
		b.EffectiveEndDate = req.EffectiveEndDate
	}
	if req.Status != nil {
		b.Status = *req.Status
	}
	if err := s.repo.UpdateEmployeeBankProfile(ctx, b); err != nil {
		return nil, err
	}
	response := toEmployeeBankProfileResponse(b)
	return &response, nil
}

func (s *Service) DeleteEmployeeBankProfile(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteEmployeeBankProfile(ctx, uid)
}

func (s *Service) CreateEmployeeBankProfile(ctx context.Context, req CreateEmployeeBankProfileRequest) (*EmployeeBankProfileResponse, error) {
	empID, _ := uuid.Parse(req.EmployeeID)
	profileID, _ := uuid.Parse(req.EmployeePayrollProfileID)
	b := &EmployeeBankProfile{
		EmployeeID:             empID,
		EmployeePayrollProfileID: profileID,
		BankName:               req.BankName,
		BankAccountNumber:      req.BankAccountNumber,
		BankAccountHolderName:  req.BankAccountHolderName,
		IsPrimary:              true,
		EffectiveStartDate:     req.EffectiveStartDate,
		Status:                 "ACTIVE",
	}
	if req.BankCode != nil {
		b.BankCode = req.BankCode
	}
	if req.BankBranch != nil {
		b.BankBranch = req.BankBranch
	}
	if req.IsPrimary != nil {
		b.IsPrimary = *req.IsPrimary
	}
	if req.EffectiveEndDate != nil {
		b.EffectiveEndDate = req.EffectiveEndDate
	}
	if req.Status != "" {
		b.Status = req.Status
	}
	b.CreatedBy = authctx.GetUserID(ctx)
	b.UpdatedBy = b.CreatedBy
	if err := s.repo.CreateEmployeeBankProfile(ctx, b); err != nil {
		return nil, err
	}
	response := toEmployeeBankProfileResponse(b)
	return &response, nil
}

func (s *Service) GetEmployeeBpjsProfileByID(ctx context.Context, id string) (*EmployeeBpjsProfileResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	b, err := s.repo.FindEmployeeBpjsProfileByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := toEmployeeBpjsProfileResponse(b)
	return &response, nil
}

func (s *Service) UpdateEmployeeBpjsProfile(ctx context.Context, id string, req UpdateEmployeeBpjsProfileRequest) (*EmployeeBpjsProfileResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	b, err := s.repo.FindEmployeeBpjsProfileByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	b.UpdatedBy = authctx.GetUserID(ctx)
	if req.BpjsHealthActive != nil {
		b.BpjsHealthActive = *req.BpjsHealthActive
	}
	if req.BpjsHealthNo != nil {
		b.BpjsHealthNo = req.BpjsHealthNo
	}
	if req.BpjsHealthRegisteredName != nil {
		b.BpjsHealthRegisteredName = req.BpjsHealthRegisteredName
	}
	if req.BpjsTkActive != nil {
		b.BpjsTkActive = *req.BpjsTkActive
	}
	if req.BpjsTkNo != nil {
		b.BpjsTkNo = req.BpjsTkNo
	}
	if req.BpjsTkRegisteredName != nil {
		b.BpjsTkRegisteredName = req.BpjsTkRegisteredName
	}
	if req.JkkRiskClass != nil {
		b.JkkRiskClass = *req.JkkRiskClass
	}
	if req.PensionActive != nil {
		b.PensionActive = *req.PensionActive
	}
	if req.EffectiveStartDate != nil {
		b.EffectiveStartDate = *req.EffectiveStartDate
	}
	if req.EffectiveEndDate != nil {
		b.EffectiveEndDate = req.EffectiveEndDate
	}
	if req.Status != nil {
		b.Status = *req.Status
	}
	if err := s.repo.UpdateEmployeeBpjsProfile(ctx, b); err != nil {
		return nil, err
	}
	response := toEmployeeBpjsProfileResponse(b)
	return &response, nil
}

func (s *Service) DeleteEmployeeBpjsProfile(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteEmployeeBpjsProfile(ctx, uid)
}

func (s *Service) CreateEmployeeBpjsProfile(ctx context.Context, req CreateEmployeeBpjsProfileRequest) (*EmployeeBpjsProfileResponse, error) {
	empID, _ := uuid.Parse(req.EmployeeID)
	profileID, _ := uuid.Parse(req.EmployeePayrollProfileID)
	b := &EmployeeBpjsProfile{
		EmployeeID:               empID,
		EmployeePayrollProfileID: profileID,
		JkkRiskClass:             "LOW",
		PensionActive:            true,
		EffectiveStartDate:       req.EffectiveStartDate,
		Status:                   "ACTIVE",
	}
	if req.BpjsHealthActive != nil {
		b.BpjsHealthActive = *req.BpjsHealthActive
	}
	if req.BpjsHealthNo != nil {
		b.BpjsHealthNo = req.BpjsHealthNo
	}
	if req.BpjsHealthRegisteredName != nil {
		b.BpjsHealthRegisteredName = req.BpjsHealthRegisteredName
	}
	if req.BpjsTkActive != nil {
		b.BpjsTkActive = *req.BpjsTkActive
	}
	if req.BpjsTkNo != nil {
		b.BpjsTkNo = req.BpjsTkNo
	}
	if req.BpjsTkRegisteredName != nil {
		b.BpjsTkRegisteredName = req.BpjsTkRegisteredName
	}
	if req.JkkRiskClass != nil {
		b.JkkRiskClass = *req.JkkRiskClass
	}
	if req.PensionActive != nil {
		b.PensionActive = *req.PensionActive
	}
	if req.EffectiveEndDate != nil {
		b.EffectiveEndDate = req.EffectiveEndDate
	}
	if req.Status != "" {
		b.Status = req.Status
	}
	if req.Notes != nil {
		b.Notes = req.Notes
	}
	b.CreatedBy = authctx.GetUserID(ctx)
	b.UpdatedBy = b.CreatedBy
	if err := s.repo.CreateEmployeeBpjsProfile(ctx, b); err != nil {
		return nil, err
	}
	response := toEmployeeBpjsProfileResponse(b)
	return &response, nil
}

func (s *Service) GetEmployeeTaxProfileByID(ctx context.Context, id string) (*EmployeeTaxProfileResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	t, err := s.repo.FindEmployeeTaxProfileByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := toEmployeeTaxProfileResponse(t)
	return &response, nil
}

func (s *Service) UpdateEmployeeTaxProfile(ctx context.Context, id string, req UpdateEmployeeTaxProfileRequest) (*EmployeeTaxProfileResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	t, err := s.repo.FindEmployeeTaxProfileByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	t.UpdatedBy = authctx.GetUserID(ctx)
	if req.Npwp != nil {
		t.Npwp = req.Npwp
	}
	if req.NpwpRegisteredName != nil {
		t.NpwpRegisteredName = req.NpwpRegisteredName
	}
	if req.PtkpStatus != nil {
		t.PtkpStatus = req.PtkpStatus
	}
	if req.TaxMethod != nil {
		t.TaxMethod = *req.TaxMethod
	}
	if req.IsTaxable != nil {
		t.IsTaxable = *req.IsTaxable
	}
	if req.HasNpwp != nil {
		t.HasNpwp = *req.HasNpwp
	}
	if req.EffectiveStartDate != nil {
		t.EffectiveStartDate = *req.EffectiveStartDate
	}
	if req.EffectiveEndDate != nil {
		t.EffectiveEndDate = req.EffectiveEndDate
	}
	if req.Status != nil {
		t.Status = *req.Status
	}
	if err := s.repo.UpdateEmployeeTaxProfile(ctx, t); err != nil {
		return nil, err
	}
	response := toEmployeeTaxProfileResponse(t)
	return &response, nil
}

func (s *Service) DeleteEmployeeTaxProfile(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteEmployeeTaxProfile(ctx, uid)
}

func (s *Service) CreateEmployeeTaxProfile(ctx context.Context, req CreateEmployeeTaxProfileRequest) (*EmployeeTaxProfileResponse, error) {
	empID, _ := uuid.Parse(req.EmployeeID)
	profileID, _ := uuid.Parse(req.EmployeePayrollProfileID)
	t := &EmployeeTaxProfile{
		EmployeeID:               empID,
		EmployeePayrollProfileID: profileID,
		TaxMethod:                "GROSS",
		IsTaxable:                true,
		HasNpwp:                  false,
		EffectiveStartDate:       req.EffectiveStartDate,
		Status:                   "ACTIVE",
	}
	if req.Npwp != nil {
		t.Npwp = req.Npwp
	}
	if req.NpwpRegisteredName != nil {
		t.NpwpRegisteredName = req.NpwpRegisteredName
	}
	if req.PtkpStatus != nil {
		t.PtkpStatus = req.PtkpStatus
	}
	if req.TaxMethod != "" {
		t.TaxMethod = req.TaxMethod
	}
	if req.IsTaxable != nil {
		t.IsTaxable = *req.IsTaxable
	}
	if req.HasNpwp != nil {
		t.HasNpwp = *req.HasNpwp
	}
	if req.EffectiveEndDate != nil {
		t.EffectiveEndDate = req.EffectiveEndDate
	}
	if req.Status != "" {
		t.Status = req.Status
	}
	if req.Notes != nil {
		t.Notes = req.Notes
	}
	t.CreatedBy = authctx.GetUserID(ctx)
	t.UpdatedBy = t.CreatedBy
	if err := s.repo.CreateEmployeeTaxProfile(ctx, t); err != nil {
		return nil, err
	}
	response := toEmployeeTaxProfileResponse(t)
	return &response, nil
}

// =============================================================================
// BPJS Settings
// =============================================================================

func (s *Service) CreateBpjsSetting(ctx context.Context, req CreateBpjsSettingRequest) (*BpjsSettingResponse, error) {
	bs := &BpjsSetting{
		SettingCode:         req.SettingCode,
		SettingName:         req.SettingName,
		BaseSource:          "BPJS_BASE_COMPONENTS",
		DefaultJkkRiskClass: "LOW",
		RoundingMode:        "ROUND",
		EffectiveStartDate:  req.EffectiveStartDate,
		Status:              "ACTIVE",
	}
	if req.BaseSource != "" {
		bs.BaseSource = req.BaseSource
	}
	if req.HealthMaxBaseAmount != nil {
		bs.HealthMaxBaseAmount = req.HealthMaxBaseAmount
	}
	if req.PensionMaxBaseAmount != nil {
		bs.PensionMaxBaseAmount = req.PensionMaxBaseAmount
	}
	if req.DefaultJkkRiskClass != "" {
		bs.DefaultJkkRiskClass = req.DefaultJkkRiskClass
	}
	if req.RoundingMode != "" {
		bs.RoundingMode = req.RoundingMode
	}
	if req.EffectiveEndDate != nil {
		bs.EffectiveEndDate = req.EffectiveEndDate
	}
	if req.Status != "" {
		bs.Status = req.Status
	}
	if req.Notes != nil {
		bs.Notes = req.Notes
	}
	bs.CreatedBy = authctx.GetUserID(ctx)
	bs.UpdatedBy = bs.CreatedBy
	if err := s.repo.CreateBpjsSetting(ctx, bs); err != nil {
		return nil, err
	}
	response := toBpjsSettingResponse(bs)
	return &response, nil
}

// =============================================================================
// PPh21 Settings
// =============================================================================

func (s *Service) CreatePph21Setting(ctx context.Context, req CreatePph21SettingRequest) (*Pph21SettingResponse, error) {
	compID, _ := uuid.Parse(req.Pph21ComponentID)
	ps := &Pph21Setting{
		SettingCode:                    req.SettingCode,
		SettingName:                    req.SettingName,
		CalculationMethod:              "REGULAR_GROSS_ANNUALIZED",
		DefaultTaxMethod:               "GROSS",
		Pph21ComponentID:               compID,
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
		EffectiveStartDate:             req.EffectiveStartDate,
		Status:                         "ACTIVE",
	}
	if req.CalculationMethod != "" {
		ps.CalculationMethod = req.CalculationMethod
	}
	if req.DefaultTaxMethod != "" {
		ps.DefaultTaxMethod = req.DefaultTaxMethod
	}
	if req.OccupationalExpenseRatePercent != nil {
		ps.OccupationalExpenseRatePercent = *req.OccupationalExpenseRatePercent
	}
	if req.OccupationalExpenseMaxMonthly != nil {
		ps.OccupationalExpenseMaxMonthly = *req.OccupationalExpenseMaxMonthly
	}
	if req.OccupationalExpenseMaxYearly != nil {
		ps.OccupationalExpenseMaxYearly = *req.OccupationalExpenseMaxYearly
	}
	if req.DeductBpjsHealthEmployee != nil {
		ps.DeductBpjsHealthEmployee = *req.DeductBpjsHealthEmployee
	}
	if req.DeductBpjsJhtEmployee != nil {
		ps.DeductBpjsJhtEmployee = *req.DeductBpjsJhtEmployee
	}
	if req.DeductBpjsJpEmployee != nil {
		ps.DeductBpjsJpEmployee = *req.DeductBpjsJpEmployee
	}
	if req.AnnualizationMonths != nil {
		ps.AnnualizationMonths = *req.AnnualizationMonths
	}
	if req.PkpRoundingUnit != nil {
		ps.PkpRoundingUnit = *req.PkpRoundingUnit
	}
	if req.NonNpwpMultiplierPercent != nil {
		ps.NonNpwpMultiplierPercent = *req.NonNpwpMultiplierPercent
	}
	if req.RoundingMode != nil {
		ps.RoundingMode = *req.RoundingMode
	}
	if req.EffectiveEndDate != nil {
		ps.EffectiveEndDate = req.EffectiveEndDate
	}
	if req.Status != "" {
		ps.Status = req.Status
	}
	ps.CreatedBy = authctx.GetUserID(ctx)
	ps.UpdatedBy = ps.CreatedBy
	if err := s.repo.CreatePph21Setting(ctx, ps); err != nil {
		return nil, err
	}
	response := toPph21SettingResponse(ps)
	return &response, nil
}

// =============================================================================
// Payroll Runs
// =============================================================================

func (s *Service) CreatePayrollRun(ctx context.Context, req CreatePayrollRunRequest) (*PayrollRunResponse, error) {
	periodID, err := uuid.Parse(req.PayrollPeriodID)
	if err != nil {
		return nil, fmt.Errorf("invalid payroll_period_id: %w", err)
	}
	pr := &PayrollRun{
		PayrollPeriodID: periodID,
		RunCode:         req.RunCode,
		RunType:         "REGULAR",
		Status:          "DRAFT",
	}
	if req.RunType != "" {
		pr.RunType = req.RunType
	}
	pr.CreatedBy = authctx.GetUserID(ctx)
	pr.UpdatedBy = pr.CreatedBy
	if err := s.repo.CreatePayrollRun(ctx, pr); err != nil {
		return nil, err
	}
	s.logger.Info("Payroll run created", zap.String("code", pr.RunCode))
	response := toPayrollRunResponse(pr)
	return &response, nil
}

func (s *Service) ListPayrollRuns(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	items, total, err := s.repo.FindAllPayrollRuns(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	var responses []PayrollRunResponse
	for _, item := range items {
		responses = append(responses, toPayrollRunResponse(&item))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: totalPages,
	}, nil
}

func (s *Service) GetPayrollRunByID(ctx context.Context, id string) (*PayrollRunResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	pr, err := s.repo.FindPayrollRunByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := toPayrollRunResponse(pr)
	return &response, nil
}

func (s *Service) UpdatePayrollRunStatus(ctx context.Context, id string, req UpdatePayrollRunStatusRequest) (*PayrollRunResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	pr, err := s.repo.FindPayrollRunByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	// Auto-create approval instance when transitioning to CALCULATED
	if req.Status == "CALCULATED" && pr.Status == "DRAFT" {
		now := time.Now()
		pr.CalculatedAt = &now

		// If approval engine is available, create approval instance
		if s.approvalEngine != nil && req.FlowID != nil && *req.FlowID != "" {
			instanceID, err := s.approvalEngine.CreateApprovalInstance(ctx, "payroll", pr.ID.String(), *req.FlowID)
			if err != nil {
				s.logger.Warn("Failed to create approval instance for payroll run, continuing without approval",
					zap.String("run_id", pr.ID.String()),
					zap.Error(err),
				)
			} else {
				s.logger.Info("Approval instance created for payroll run",
					zap.String("run_id", pr.ID.String()),
					zap.String("instance_id", instanceID),
				)
				pr.Status = "CALCULATED"
				if err := s.repo.UpdatePayrollRun(ctx, pr); err != nil {
					return nil, err
				}
				response := toPayrollRunResponse(pr)
				return &response, nil
			}
		}

		// Without approval engine, move directly to REVIEWED
		pr.Status = "REVIEWED"
	} else if req.Status == "APPROVED" && pr.Status == "CALCULATED" {
		// Coming from approval completion — check if approval was granted
		now := time.Now()
		pr.ApprovedAt = &now
		pr.Status = "APPROVED"
	} else if req.Status == "REVIEWED" && pr.Status == "CALCULATED" {
		// Manual review (no approval engine configured)
		now := time.Now()
		pr.ReviewedAt = &now
		pr.Status = "REVIEWED"
	} else if req.Status == "LOCKED" && pr.Status == "APPROVED" {
		now := time.Now()
		pr.LockedAt = &now
		pr.Status = "LOCKED"
	} else {
		// Direct status update for other transitions
		pr.Status = req.Status
	}

	if err := s.repo.UpdatePayrollRun(ctx, pr); err != nil {
		return nil, err
	}

	response := toPayrollRunResponse(pr)
	return &response, nil
}

// CheckPayrollRunApproval checks the approval status of a payroll run.
// If the approval instance is fully approved, transitions to APPROVED.
// If rejected, transitions back to DRAFT.
func (s *Service) CheckPayrollRunApproval(ctx context.Context, id string, instanceID string) (*PayrollRunResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	pr, err := s.repo.FindPayrollRunByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if pr.Status != "CALCULATED" {
		return nil, fmt.Errorf("payroll run is not in CALCULATED status (current: %s)", pr.Status)
	}

	if s.approvalEngine == nil {
		return nil, fmt.Errorf("approval engine not configured")
	}

	approvalStatus, err := s.approvalEngine.GetApprovalInstanceStatus(ctx, instanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get approval status: %w", err)
	}

	switch approvalStatus {
	case "APPROVED":
		now := time.Now()
		pr.ReviewedAt = &now
		pr.ApprovedAt = &now
		pr.Status = "APPROVED"
		s.logger.Info("Payroll run approved via approval engine",
			zap.String("run_id", pr.ID.String()),
		)
	case "REJECTED":
		pr.Status = "DRAFT"
		s.logger.Info("Payroll run rejected, returning to DRAFT",
			zap.String("run_id", pr.ID.String()),
		)
	default:
		// Still PENDING — no status change
		response := toPayrollRunResponse(pr)
		return &response, nil
	}

	if err := s.repo.UpdatePayrollRun(ctx, pr); err != nil {
		return nil, err
	}

	response := toPayrollRunResponse(pr)
	return &response, nil
}

// =============================================================================
// Simplified CRUD proxies (for entities that are mostly create/list/get/delete)
// =============================================================================

// BpjsSetting CRUD
func (s *Service) UpdateBpjsSetting(ctx context.Context, id string, req UpdateBpjsSettingRequest) (*BpjsSettingResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	bs, err := s.repo.FindBpjsSettingByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	bs.UpdatedBy = authctx.GetUserID(ctx)
	if req.SettingCode != nil {
		bs.SettingCode = *req.SettingCode
	}
	if req.SettingName != nil {
		bs.SettingName = *req.SettingName
	}
	if req.BaseSource != nil {
		bs.BaseSource = *req.BaseSource
	}
	if req.HealthMaxBaseAmount != nil {
		bs.HealthMaxBaseAmount = req.HealthMaxBaseAmount
	}
	if req.PensionMaxBaseAmount != nil {
		bs.PensionMaxBaseAmount = req.PensionMaxBaseAmount
	}
	if req.DefaultJkkRiskClass != nil {
		bs.DefaultJkkRiskClass = *req.DefaultJkkRiskClass
	}
	if req.RoundingMode != nil {
		bs.RoundingMode = *req.RoundingMode
	}
	if req.EffectiveStartDate != nil {
		bs.EffectiveStartDate = *req.EffectiveStartDate
	}
	if req.EffectiveEndDate != nil {
		bs.EffectiveEndDate = req.EffectiveEndDate
	}
	if req.Status != nil {
		bs.Status = *req.Status
	}
	if req.Notes != nil {
		bs.Notes = req.Notes
	}
	if err := s.repo.UpdateBpjsSetting(ctx, bs); err != nil {
		return nil, err
	}
	response := toBpjsSettingResponse(bs)
	return &response, nil
}

func (s *Service) GetBpjsSettingByID(ctx context.Context, id string) (*BpjsSettingResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	bs, err := s.repo.FindBpjsSettingByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := toBpjsSettingResponse(bs)
	return &response, nil
}

func (s *Service) ListBpjsSettings(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	items, total, err := s.repo.FindAllBpjsSettings(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	var responses []BpjsSettingResponse
	for _, item := range items {
		responses = append(responses, toBpjsSettingResponse(&item))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: totalPages,
	}, nil
}

func (s *Service) DeleteBpjsSetting(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteBpjsSetting(ctx, uid)
}

// BpjsRateComponent CRUD
func (s *Service) GetBpjsRateComponentByID(ctx context.Context, id string) (*BpjsRateComponentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	br, err := s.repo.FindBpjsRateComponentByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := toBpjsRateComponentResponse(br)
	return &response, nil
}

func (s *Service) UpdateBpjsRateComponent(ctx context.Context, id string, req UpdateBpjsRateComponentRequest) (*BpjsRateComponentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	br, err := s.repo.FindBpjsRateComponentByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	br.UpdatedBy = authctx.GetUserID(ctx)
	if req.BpjsSettingID != nil && *req.BpjsSettingID != "" {
		sid, _ := uuid.Parse(*req.BpjsSettingID)
		br.BpjsSettingID = sid
	}
	if req.RateCode != nil {
		br.RateCode = *req.RateCode
	}
	if req.RateName != nil {
		br.RateName = *req.RateName
	}
	if req.BpjsProgram != nil {
		br.BpjsProgram = *req.BpjsProgram
	}
	if req.PaidBy != nil {
		br.PaidBy = *req.PaidBy
	}
	if req.SalaryComponentID != nil {
		if *req.SalaryComponentID != "" {
			id, _ := uuid.Parse(*req.SalaryComponentID)
			br.SalaryComponentID = &id
		} else {
			br.SalaryComponentID = nil
		}
	}
	if req.RatePercent != nil {
		br.RatePercent = *req.RatePercent
	}
	if req.FixedAmount != nil {
		br.FixedAmount = req.FixedAmount
	}
	if req.MinBaseAmount != nil {
		br.MinBaseAmount = req.MinBaseAmount
	}
	if req.MaxBaseAmount != nil {
		br.MaxBaseAmount = req.MaxBaseAmount
	}
	if req.JkkRiskClass != nil {
		br.JkkRiskClass = req.JkkRiskClass
	}
	if req.IsEmployeeDeduction != nil {
		br.IsEmployeeDeduction = *req.IsEmployeeDeduction
	}
	if req.IsEmployerContribution != nil {
		br.IsEmployerContribution = *req.IsEmployerContribution
	}
	if req.GenerateToPayrollItem != nil {
		br.GenerateToPayrollItem = *req.GenerateToPayrollItem
	}
	if req.PrintOnPayslip != nil {
		br.PrintOnPayslip = *req.PrintOnPayslip
	}
	if req.DisplayOrder != nil {
		br.DisplayOrder = *req.DisplayOrder
	}
	if req.EffectiveStartDate != nil {
		br.EffectiveStartDate = *req.EffectiveStartDate
	}
	if req.EffectiveEndDate != nil {
		br.EffectiveEndDate = req.EffectiveEndDate
	}
	if req.Status != nil {
		br.Status = *req.Status
	}
	if err := s.repo.UpdateBpjsRateComponent(ctx, br); err != nil {
		return nil, err
	}
	response := toBpjsRateComponentResponse(br)
	return &response, nil
}

func (s *Service) DeleteBpjsRateComponent(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteBpjsRateComponent(ctx, uid)
}

func (s *Service) CreateBpjsRateComponent(ctx context.Context, req CreateBpjsRateComponentRequest) (*BpjsRateComponentResponse, error) {
	settingID, _ := uuid.Parse(req.BpjsSettingID)
	br := &BpjsRateComponent{
		BpjsSettingID:      settingID,
		RateCode:           req.RateCode,
		RateName:           req.RateName,
		BpjsProgram:        req.BpjsProgram,
		PaidBy:             req.PaidBy,
		RatePercent:        req.RatePercent,
		IsEmployeeDeduction:     false,
		IsEmployerContribution:  false,
		GenerateToPayrollItem:   true,
		PrintOnPayslip:          true,
		DisplayOrder:            0,
		EffectiveStartDate:      req.EffectiveStartDate,
		Status:                  "ACTIVE",
	}
	if req.SalaryComponentID != nil && *req.SalaryComponentID != "" {
		id, _ := uuid.Parse(*req.SalaryComponentID)
		br.SalaryComponentID = &id
	}
	if req.FixedAmount != nil {
		br.FixedAmount = req.FixedAmount
	}
	if req.MinBaseAmount != nil {
		br.MinBaseAmount = req.MinBaseAmount
	}
	if req.MaxBaseAmount != nil {
		br.MaxBaseAmount = req.MaxBaseAmount
	}
	if req.JkkRiskClass != nil {
		br.JkkRiskClass = req.JkkRiskClass
	}
	if req.IsEmployeeDeduction != nil {
		br.IsEmployeeDeduction = *req.IsEmployeeDeduction
	}
	if req.IsEmployerContribution != nil {
		br.IsEmployerContribution = *req.IsEmployerContribution
	}
	if req.GenerateToPayrollItem != nil {
		br.GenerateToPayrollItem = *req.GenerateToPayrollItem
	}
	if req.PrintOnPayslip != nil {
		br.PrintOnPayslip = *req.PrintOnPayslip
	}
	if req.DisplayOrder != nil {
		br.DisplayOrder = *req.DisplayOrder
	}
	if req.EffectiveEndDate != nil {
		br.EffectiveEndDate = req.EffectiveEndDate
	}
	if req.Status != "" {
		br.Status = req.Status
	}
	br.CreatedBy = authctx.GetUserID(ctx)
	br.UpdatedBy = br.CreatedBy
	if err := s.repo.CreateBpjsRateComponent(ctx, br); err != nil {
		return nil, err
	}
	response := toBpjsRateComponentResponse(br)
	return &response, nil
}

// Pph21Setting CRUD
func (s *Service) GetPph21SettingByID(ctx context.Context, id string) (*Pph21SettingResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	ps, err := s.repo.FindPph21SettingByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := toPph21SettingResponse(ps)
	return &response, nil
}

func (s *Service) ListPph21Settings(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	items, total, err := s.repo.FindAllPph21Settings(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	var responses []Pph21SettingResponse
	for _, item := range items {
		responses = append(responses, toPph21SettingResponse(&item))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: totalPages,
	}, nil
}

func (s *Service) UpdatePph21Setting(ctx context.Context, id string, req UpdatePph21SettingRequest) (*Pph21SettingResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	ps, err := s.repo.FindPph21SettingByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	ps.UpdatedBy = authctx.GetUserID(ctx)
	if req.SettingCode != nil {
		ps.SettingCode = *req.SettingCode
	}
	if req.SettingName != nil {
		ps.SettingName = *req.SettingName
	}
	if req.CalculationMethod != nil {
		ps.CalculationMethod = *req.CalculationMethod
	}
	if req.DefaultTaxMethod != nil {
		ps.DefaultTaxMethod = *req.DefaultTaxMethod
	}
	if req.Pph21ComponentID != nil && *req.Pph21ComponentID != "" {
		cid, _ := uuid.Parse(*req.Pph21ComponentID)
		ps.Pph21ComponentID = cid
	}
	if req.OccupationalExpenseRatePercent != nil {
		ps.OccupationalExpenseRatePercent = *req.OccupationalExpenseRatePercent
	}
	if req.OccupationalExpenseMaxMonthly != nil {
		ps.OccupationalExpenseMaxMonthly = *req.OccupationalExpenseMaxMonthly
	}
	if req.OccupationalExpenseMaxYearly != nil {
		ps.OccupationalExpenseMaxYearly = *req.OccupationalExpenseMaxYearly
	}
	if req.DeductBpjsHealthEmployee != nil {
		ps.DeductBpjsHealthEmployee = *req.DeductBpjsHealthEmployee
	}
	if req.DeductBpjsJhtEmployee != nil {
		ps.DeductBpjsJhtEmployee = *req.DeductBpjsJhtEmployee
	}
	if req.DeductBpjsJpEmployee != nil {
		ps.DeductBpjsJpEmployee = *req.DeductBpjsJpEmployee
	}
	if req.AnnualizationMonths != nil {
		ps.AnnualizationMonths = *req.AnnualizationMonths
	}
	if req.PkpRoundingUnit != nil {
		ps.PkpRoundingUnit = *req.PkpRoundingUnit
	}
	if req.NonNpwpMultiplierPercent != nil {
		ps.NonNpwpMultiplierPercent = *req.NonNpwpMultiplierPercent
	}
	if req.RoundingMode != nil {
		ps.RoundingMode = *req.RoundingMode
	}
	if req.EffectiveStartDate != nil {
		ps.EffectiveStartDate = *req.EffectiveStartDate
	}
	if req.EffectiveEndDate != nil {
		ps.EffectiveEndDate = req.EffectiveEndDate
	}
	if req.Status != nil {
		ps.Status = *req.Status
	}
	if err := s.repo.UpdatePph21Setting(ctx, ps); err != nil {
		return nil, err
	}
	response := toPph21SettingResponse(ps)
	return &response, nil
}

func (s *Service) DeletePph21Setting(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeletePph21Setting(ctx, uid)
}

// Pph21PtkpRate CRUD
func (s *Service) CreatePph21PtkpRate(ctx context.Context, req CreatePph21PtkpRateRequest) (*Pph21PtkpRateResponse, error) {
	pr := &Pph21PtkpRate{
		PtkpStatus:         req.PtkpStatus,
		AnnualAmount:       req.AnnualAmount,
		EffectiveStartDate: req.EffectiveStartDate,
		Status:             "ACTIVE",
	}
	if req.Description != nil {
		pr.Description = req.Description
	}
	if req.EffectiveEndDate != nil {
		pr.EffectiveEndDate = req.EffectiveEndDate
	}
	if req.Status != "" {
		pr.Status = req.Status
	}
	pr.CreatedBy = authctx.GetUserID(ctx)
	if err := s.repo.CreatePph21PtkpRate(ctx, pr); err != nil {
		return nil, err
	}
	response := toPph21PtkpRateResponse(pr)
	return &response, nil
}

func (s *Service) ListPph21PtkpRates(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	items, total, err := s.repo.FindAllPph21PtkpRates(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	var responses []Pph21PtkpRateResponse
	for _, item := range items {
		responses = append(responses, toPph21PtkpRateResponse(&item))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: totalPages,
	}, nil
}

// Pph21TaxBracket CRUD
func (s *Service) CreatePph21TaxBracket(ctx context.Context, req CreatePph21TaxBracketRequest) (*Pph21TaxBracketResponse, error) {
	tb := &Pph21TaxBracket{
		BracketOrder:       req.BracketOrder,
		LowerBound:         req.LowerBound,
		RatePercent:        req.RatePercent,
		EffectiveStartDate: req.EffectiveStartDate,
		Status:             "ACTIVE",
	}
	if req.UpperBound != nil {
		tb.UpperBound = req.UpperBound
	}
	if req.EffectiveEndDate != nil {
		tb.EffectiveEndDate = req.EffectiveEndDate
	}
	if req.Status != "" {
		tb.Status = req.Status
	}
	tb.CreatedBy = authctx.GetUserID(ctx)
	if err := s.repo.CreatePph21TaxBracket(ctx, tb); err != nil {
		return nil, err
	}
	response := toPph21TaxBracketResponse(tb)
	return &response, nil
}

func (s *Service) ListPph21TaxBrackets(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	items, total, err := s.repo.FindAllPph21TaxBrackets(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	var responses []Pph21TaxBracketResponse
	for _, item := range items {
		responses = append(responses, toPph21TaxBracketResponse(&item))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: totalPages,
	}, nil
}
