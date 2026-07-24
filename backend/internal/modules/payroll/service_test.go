package payroll

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// =============================================================================
// Service Tests (using real SQLite repository)
// =============================================================================

func newTestService() (*Service, func()) {
	_, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	return svc, func() { cleanup(); logger.Sync() }
}

// =============================================================================
// Salary Component Service Tests
// =============================================================================

func TestService_CreateSalaryComponent_Defaults(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateSalaryComponentRequest{
		Code:          "BONUS",
		Name:          "THR Bonus",
		ComponentType: "EARNING",
	}

	resp, err := svc.CreateSalaryComponent(ctx, req)
	if err != nil {
		t.Fatalf("CreateSalaryComponent failed: %v", err)
	}

	if resp.Status != "ACTIVE" {
		t.Errorf("expected status 'ACTIVE', got '%s'", resp.Status)
	}
	if resp.IsTaxable != true {
		t.Errorf("expected IsTaxable true, got false")
	}
	if resp.DisplayOrder != 100 {
		t.Errorf("expected DisplayOrder 100, got %d", resp.DisplayOrder)
	}
}

func TestService_CreateSalaryComponent_WithOptionalFields(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	isTaxable := false
	displayOrder := 50
	status := "INACTIVE"
	req := CreateSalaryComponentRequest{
		Code:          "DEDUCT-TEST",
		Name:          "Test Deduction",
		ComponentType: "DEDUCTION",
		IsTaxable:     &isTaxable,
		DisplayOrder:  &displayOrder,
		Status:        &status,
	}

	resp, err := svc.CreateSalaryComponent(ctx, req)
	if err != nil {
		t.Fatalf("CreateSalaryComponent failed: %v", err)
	}

	if resp.IsTaxable != false {
		t.Errorf("expected IsTaxable false, got true")
	}
	if resp.DisplayOrder != 50 {
		t.Errorf("expected DisplayOrder 50, got %d", resp.DisplayOrder)
	}
	if resp.Status != "INACTIVE" {
		t.Errorf("expected status 'INACTIVE', got '%s'", resp.Status)
	}
}

func TestService_GetSalaryComponentByID_Success(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
		Code:          "BASIC-SVC",
		Name:          "Basic Salary Service",
		ComponentType: "EARNING",
	})

	found, err := svc.GetSalaryComponentByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetSalaryComponentByID failed: %v", err)
	}
	if found.Name != "Basic Salary Service" {
		t.Errorf("expected name 'Basic Salary Service', got '%s'", found.Name)
	}
}

func TestService_GetSalaryComponentByID_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetSalaryComponentByID(ctx, "not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_GetSalaryComponentByID_NotFound(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetSalaryComponentByID(ctx, uuid.New().String())
	if err == nil {
		t.Fatal("expected error for non-existent salary component")
	}
}

func TestService_ListSalaryComponents_Pagination(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
			Code:          fmt.Sprintf("COMP-LST-%03d", i+1),
			Name:          fmt.Sprintf("Component %d", i+1),
			ComponentType: "EARNING",
		})
	}

	resp, err := svc.ListSalaryComponents(ctx, 0, 0)
	if err != nil {
		t.Fatalf("ListSalaryComponents failed: %v", err)
	}

	if resp.Page != 1 {
		t.Errorf("expected page 1, got %d", resp.Page)
	}
	if resp.PerPage != 20 {
		t.Errorf("expected per_page 20 (default), got %d", resp.PerPage)
	}
	if resp.Total != 5 {
		t.Errorf("expected total 5, got %d", resp.Total)
	}
	if !resp.Success {
		t.Error("expected Success true")
	}
}

func TestService_UpdateSalaryComponent(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
		Code:          "UPD-TEST",
		Name:          "Before Update",
		ComponentType: "EARNING",
	})

	newName := "After Update"
	updated, err := svc.UpdateSalaryComponent(ctx, created.ID, UpdateSalaryComponentRequest{
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("UpdateSalaryComponent failed: %v", err)
	}

	if updated.Name != "After Update" {
		t.Errorf("expected name 'After Update', got '%s'", updated.Name)
	}
}

func TestService_DeleteSalaryComponent(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
		Code:          "DEL-TEST",
		Name:          "To Delete",
		ComponentType: "EARNING",
	})

	if err := svc.DeleteSalaryComponent(ctx, created.ID); err != nil {
		t.Fatalf("DeleteSalaryComponent failed: %v", err)
	}

	_, err := svc.GetSalaryComponentByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting salary component")
	}
}

// =============================================================================
// Payroll Period Service Tests
// =============================================================================

func TestService_CreatePayrollPeriod_GeneratesCode(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	resp, err := svc.CreatePayrollPeriod(ctx, CreatePayrollPeriodRequest{
		PeriodYear:  2026,
		PeriodMonth: 2,
		StartDate:   "2026-02-01",
		EndDate:     "2026-02-28",
		AsOfDate:    "2026-02-28",
	})
	if err != nil {
		t.Fatalf("CreatePayrollPeriod failed: %v", err)
	}

	if resp.PeriodCode != "202602" {
		t.Errorf("expected period_code '202602', got '%s'", resp.PeriodCode)
	}
	if resp.Status != "OPEN" {
		t.Errorf("expected default status 'OPEN', got '%s'", resp.Status)
	}
}

// =============================================================================
// BPJS Setting Service Tests
// =============================================================================

func TestService_CreateBpjsSetting_Defaults(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	resp, err := svc.CreateBpjsSetting(ctx, CreateBpjsSettingRequest{
		SettingCode:        "BPJS-TEST",
		SettingName:        "Test BPJS",
		EffectiveStartDate: "2026-01-01",
	})
	if err != nil {
		t.Fatalf("CreateBpjsSetting failed: %v", err)
	}

	if resp.Status != "ACTIVE" {
		t.Errorf("expected status 'ACTIVE', got '%s'", resp.Status)
	}
	if resp.BaseSource != "BPJS_BASE_COMPONENTS" {
		t.Errorf("expected base_source 'BPJS_BASE_COMPONENTS', got '%s'", resp.BaseSource)
	}
	if resp.RoundingMode != "ROUND" {
		t.Errorf("expected rounding_mode 'ROUND', got '%s'", resp.RoundingMode)
	}
}

func TestService_UpdateBpjsSetting(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateBpjsSetting(ctx, CreateBpjsSettingRequest{
		SettingCode:        "BPJS-UPD",
		SettingName:        "Before Update",
		EffectiveStartDate: "2026-01-01",
	})

	newName := "After Update"
	updated, err := svc.UpdateBpjsSetting(ctx, created.ID, UpdateBpjsSettingRequest{
		SettingName: &newName,
	})
	if err != nil {
		t.Fatalf("UpdateBpjsSetting failed: %v", err)
	}

	if updated.SettingName != "After Update" {
		t.Errorf("expected 'After Update', got '%s'", updated.SettingName)
	}
}

func TestService_ListBpjsSettings(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	svc.CreateBpjsSetting(ctx, CreateBpjsSettingRequest{
		SettingCode: "BPJS-L1", SettingName: "BPJS 1", EffectiveStartDate: "2026-01-01",
	})
	svc.CreateBpjsSetting(ctx, CreateBpjsSettingRequest{
		SettingCode: "BPJS-L2", SettingName: "BPJS 2", EffectiveStartDate: "2026-01-01",
	})

	resp, err := svc.ListBpjsSettings(ctx, 1, 10)
	if err != nil {
		t.Fatalf("ListBpjsSettings failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
}

func TestService_DeleteBpjsSetting(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateBpjsSetting(ctx, CreateBpjsSettingRequest{
		SettingCode: "BPJS-DEL", SettingName: "To Delete", EffectiveStartDate: "2026-01-01",
	})

	if err := svc.DeleteBpjsSetting(ctx, created.ID); err != nil {
		t.Fatalf("DeleteBpjsSetting failed: %v", err)
	}
}

// =============================================================================
// BPJS Rate Component Service Tests
// =============================================================================

func TestService_CreateBpjsRateComponent(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	setting, _ := svc.CreateBpjsSetting(ctx, CreateBpjsSettingRequest{
		SettingCode: "BPJS-RC", SettingName: "Rate Test", EffectiveStartDate: "2026-01-01",
	})

	resp, err := svc.CreateBpjsRateComponent(ctx, CreateBpjsRateComponentRequest{
		BpjsSettingID:      setting.ID,
		RateCode:           "JHT-EMP",
		RateName:           "JHT Employee",
		BpjsProgram:        "JHT",
		PaidBy:             "EMPLOYEE",
		RatePercent:        2.0,
		EffectiveStartDate: "2026-01-01",
	})
	if err != nil {
		t.Fatalf("CreateBpjsRateComponent failed: %v", err)
	}

	if resp.RatePercent != 2.0 {
		t.Errorf("expected rate 2.0, got %f", resp.RatePercent)
	}
}

func TestService_GetBpjsRateComponentByID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	setting, _ := svc.CreateBpjsSetting(ctx, CreateBpjsSettingRequest{
		SettingCode: "BPJS-GET", SettingName: "Get Test", EffectiveStartDate: "2026-01-01",
	})

	created, _ := svc.CreateBpjsRateComponent(ctx, CreateBpjsRateComponentRequest{
		BpjsSettingID:      setting.ID,
		RateCode:           "GET-RATE",
		RateName:           "Get Rate",
		BpjsProgram:        "JKK",
		PaidBy:             "EMPLOYER",
		RatePercent:        0.24,
		EffectiveStartDate: "2026-01-01",
	})

	found, err := svc.GetBpjsRateComponentByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetBpjsRateComponentByID failed: %v", err)
	}
	if found.RateCode != "GET-RATE" {
		t.Errorf("expected 'GET-RATE', got '%s'", found.RateCode)
	}
}

// =============================================================================
// PPh21 Settings Service Tests
// =============================================================================

func TestService_CreatePph21Setting(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	comp, _ := svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
		Code: "PPH-COMP", Name: "PPh21 Component", ComponentType: "DEDUCTION",
	})

	resp, err := svc.CreatePph21Setting(ctx, CreatePph21SettingRequest{
		SettingCode:        "PPH-TEST",
		SettingName:        "Test PPh21",
		Pph21ComponentID:   comp.ID,
		EffectiveStartDate: "2026-01-01",
	})
	if err != nil {
		t.Fatalf("CreatePph21Setting failed: %v", err)
	}

	if resp.Status != "ACTIVE" {
		t.Errorf("expected status 'ACTIVE', got '%s'", resp.Status)
	}
	if resp.AnnualizationMonths != 12 {
		t.Errorf("expected annualization_months 12, got %d", resp.AnnualizationMonths)
	}
}

// =============================================================================
// Payroll Run Service Tests
// =============================================================================

func TestService_CreatePayrollRun(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	period, _ := svc.CreatePayrollPeriod(ctx, CreatePayrollPeriodRequest{
		PeriodYear: 2026, PeriodMonth: 3,
		StartDate: "2026-03-01", EndDate: "2026-03-31", AsOfDate: "2026-03-31",
	})

	resp, err := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID,
		RunCode:         "RUN-MAR-2026",
		RunType:         "REGULAR",
	})
	if err != nil {
		t.Fatalf("CreatePayrollRun failed: %v", err)
	}

	if resp.Status != "DRAFT" {
		t.Errorf("expected default status 'DRAFT', got '%s'", resp.Status)
	}
	if resp.RunCode != "RUN-MAR-2026" {
		t.Errorf("expected run_code 'RUN-MAR-2026', got '%s'", resp.RunCode)
	}
}

func TestService_UpdatePayrollRunStatus(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	period, _ := svc.CreatePayrollPeriod(ctx, CreatePayrollPeriodRequest{
		PeriodYear: 2026, PeriodMonth: 4,
		StartDate: "2026-04-01", EndDate: "2026-04-30", AsOfDate: "2026-04-30",
	})

	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID,
		RunCode:         "RUN-APR-2026",
	})

	updated, err := svc.UpdatePayrollRunStatus(ctx, run.ID, UpdatePayrollRunStatusRequest{
		Status: "CALCULATED",
	})
	if err != nil {
		t.Fatalf("UpdatePayrollRunStatus failed: %v", err)
	}

	if updated.Status != "CALCULATED" {
		t.Errorf("expected status 'CALCULATED', got '%s'", updated.Status)
	}
}

func TestService_ListPayrollRuns(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	period, _ := svc.CreatePayrollPeriod(ctx, CreatePayrollPeriodRequest{
		PeriodYear: 2026, PeriodMonth: 5,
		StartDate: "2026-05-01", EndDate: "2026-05-31", AsOfDate: "2026-05-31",
	})

	svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID, RunCode: "RUN-MAY-01",
	})
	svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID, RunCode: "RUN-MAY-02",
	})

	resp, err := svc.ListPayrollRuns(ctx, 1, 10)
	if err != nil {
		t.Fatalf("ListPayrollRuns failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
}

// =============================================================================
// PPh21 PTKP Rate & Tax Bracket Service Tests
// =============================================================================

func TestService_CreatePph21PtkpRate(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	resp, err := svc.CreatePph21PtkpRate(ctx, CreatePph21PtkpRateRequest{
		PtkpStatus:         "K/3",
		AnnualAmount:       72000000,
		EffectiveStartDate: "2026-01-01",
	})
	if err != nil {
		t.Fatalf("CreatePph21PtkpRate failed: %v", err)
	}

	if resp.PtkpStatus != "K/3" {
		t.Errorf("expected 'K/3', got '%s'", resp.PtkpStatus)
	}
}

func TestService_ListPph21PtkpRates(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	svc.CreatePph21PtkpRate(ctx, CreatePph21PtkpRateRequest{
		PtkpStatus: "TK/0", AnnualAmount: 54000000, EffectiveStartDate: "2026-01-01",
	})
	svc.CreatePph21PtkpRate(ctx, CreatePph21PtkpRateRequest{
		PtkpStatus: "K/0", AnnualAmount: 58500000, EffectiveStartDate: "2026-01-01",
	})

	resp, err := svc.ListPph21PtkpRates(ctx, 1, 10)
	if err != nil {
		t.Fatalf("ListPph21PtkpRates failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
}

func TestService_CreatePph21TaxBracket(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	resp, err := svc.CreatePph21TaxBracket(ctx, CreatePph21TaxBracketRequest{
		BracketOrder:       1,
		LowerBound:         0,
		UpperBound:         float64Ptr(60000000),
		RatePercent:        5,
		EffectiveStartDate: "2026-01-01",
	})
	if err != nil {
		t.Fatalf("CreatePph21TaxBracket failed: %v", err)
	}

	if resp.RatePercent != 5 {
		t.Errorf("expected rate 5, got %f", resp.RatePercent)
	}
}

func TestService_ListPph21TaxBrackets(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	svc.CreatePph21TaxBracket(ctx, CreatePph21TaxBracketRequest{
		BracketOrder: 1, LowerBound: 0, UpperBound: float64Ptr(60000000),
		RatePercent: 5, EffectiveStartDate: "2026-01-01",
	})
	svc.CreatePph21TaxBracket(ctx, CreatePph21TaxBracketRequest{
		BracketOrder: 2, LowerBound: 60000000, UpperBound: float64Ptr(250000000),
		RatePercent: 15, EffectiveStartDate: "2026-01-01",
	})

	resp, err := svc.ListPph21TaxBrackets(ctx, 1, 10)
	if err != nil {
		t.Fatalf("ListPph21TaxBrackets failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
}
