package payroll

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

// =============================================================================
// Salary Component Tests
// =============================================================================

func TestRepository_CreateSalaryComponent(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	sc := &SalaryComponent{
		Code:           "TEST-001",
		Name:           "Test Salary Component",
		ComponentType:  "EARNING",
		CalculationType: "FIXED",
		IsTaxable:       true,
		DisplayOrder:    100,
		Status:          "ACTIVE",
	}

	if err := repo.CreateSalaryComponent(ctx, sc); err != nil {
		t.Fatalf("CreateSalaryComponent failed: %v", err)
	}

	if sc.ID == uuid.Nil {
		t.Error("expected salary component ID to be generated")
	}
}

func TestRepository_FindSalaryComponentByID(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestSalaryComponent(ctx, repo)

	found, err := repo.FindSalaryComponentByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindSalaryComponentByID failed: %v", err)
	}

	if found.Name != "Basic Salary" {
		t.Errorf("expected name 'Basic Salary', got '%s'", found.Name)
	}
	if found.ComponentType != "EARNING" {
		t.Errorf("expected component_type 'EARNING', got '%s'", found.ComponentType)
	}
}

func TestRepository_FindSalaryComponentByID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	_, err := repo.FindSalaryComponentByID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent salary component")
	}
}

func TestRepository_FindAllSalaryComponents_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		sc := &SalaryComponent{
			Code:           fmt.Sprintf("COMP%03d", i+1),
			Name:           fmt.Sprintf("Component %d", i+1),
			ComponentType:  "EARNING",
			CalculationType: "FIXED",
			Status:         "ACTIVE",
			DisplayOrder:   100,
		}
		repo.CreateSalaryComponent(ctx, sc)
	}

	items, total, err := repo.FindAllSalaryComponents(ctx, 1, 3)
	if err != nil {
		t.Fatalf("FindAllSalaryComponents failed: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(items) != 3 {
		t.Errorf("expected 3 items on page 1, got %d", len(items))
	}
}

func TestRepository_UpdateSalaryComponent(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestSalaryComponent(ctx, repo)

	created.Name = "Updated Component"
	if err := repo.UpdateSalaryComponent(ctx, created); err != nil {
		t.Fatalf("UpdateSalaryComponent failed: %v", err)
	}

	found, _ := repo.FindSalaryComponentByID(ctx, created.ID)
	if found.Name != "Updated Component" {
		t.Errorf("expected name 'Updated Component', got '%s'", found.Name)
	}
}

func TestRepository_DeleteSalaryComponent(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestSalaryComponent(ctx, repo)

	if err := repo.DeleteSalaryComponent(ctx, created.ID); err != nil {
		t.Fatalf("DeleteSalaryComponent failed: %v", err)
	}

	_, err := repo.FindSalaryComponentByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting salary component")
	}
}

// =============================================================================
// Payroll Period Tests
// =============================================================================

func TestRepository_CreatePayrollPeriod(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	p := &PayrollPeriod{
		PeriodCode:  "PER-2026-01",
		PeriodYear:  2026,
		PeriodMonth: 1,
		StartDate:   "2026-01-01",
		EndDate:     "2026-01-31",
		AsOfDate:    "2026-01-31",
		Status:      "OPEN",
	}

	if err := repo.CreatePayrollPeriod(ctx, p); err != nil {
		t.Fatalf("CreatePayrollPeriod failed: %v", err)
	}

	if p.ID == uuid.Nil {
		t.Error("expected payroll period ID to be generated")
	}
}

func TestRepository_FindPayrollPeriodByID(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestPayrollPeriod(ctx, repo)

	found, err := repo.FindPayrollPeriodByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindPayrollPeriodByID failed: %v", err)
	}

	if found.PeriodCode != "202601" {
		t.Errorf("expected period_code '202601', got '%s'", found.PeriodCode)
	}
}

func TestRepository_FindAllPayrollPeriods_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	for i := 1; i <= 3; i++ {
		p := &PayrollPeriod{
			PeriodCode:  fmt.Sprintf("%d%02d", 2026, i),
			PeriodYear:  2026,
			PeriodMonth: i,
			StartDate:   fmt.Sprintf("2026-%02d-01", i),
			EndDate:     fmt.Sprintf("2026-%02d-28", i),
			AsOfDate:    fmt.Sprintf("2026-%02d-28", i),
			Status:      "OPEN",
		}
		repo.CreatePayrollPeriod(ctx, p)
	}

	items, total, err := repo.FindAllPayrollPeriods(ctx, 1, 10)
	if err != nil {
		t.Fatalf("FindAllPayrollPeriods failed: %v", err)
	}

	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(items) != 3 {
		t.Errorf("expected 3 periods, got %d", len(items))
	}
}

// =============================================================================
// BPJS Setting Tests
// =============================================================================

func TestRepository_BpjsSettingCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	// Create
	created := createTestBpjsSetting(ctx, repo)

	// Find
	found, err := repo.FindBpjsSettingByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindBpjsSettingByID failed: %v", err)
	}
	if found.SettingCode != "BPJS-DEFAULT" {
		t.Errorf("expected code 'BPJS-DEFAULT', got '%s'", found.SettingCode)
	}

	// Update
	found.SettingName = "Updated BPJS Setting"
	if err := repo.UpdateBpjsSetting(ctx, found); err != nil {
		t.Fatalf("UpdateBpjsSetting failed: %v", err)
	}
	updated, _ := repo.FindBpjsSettingByID(ctx, created.ID)
	if updated.SettingName != "Updated BPJS Setting" {
		t.Errorf("expected 'Updated BPJS Setting', got '%s'", updated.SettingName)
	}

	// List
	items, total, err := repo.FindAllBpjsSettings(ctx, 1, 10)
	if err != nil {
		t.Fatalf("FindAllBpjsSettings failed: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}

	// Delete
	if err := repo.DeleteBpjsSetting(ctx, created.ID); err != nil {
		t.Fatalf("DeleteBpjsSetting failed: %v", err)
	}
	_, err = repo.FindBpjsSettingByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting BPJS setting")
	}
}

// =============================================================================
// BPJS Rate Component Tests
// =============================================================================

func TestRepository_BpjsRateComponentCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestBpjsRateComponent(ctx, repo)

	// Find
	found, err := repo.FindBpjsRateComponentByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindBpjsRateComponentByID failed: %v", err)
	}
	if found.RateCode != "BPJS-HEALTH-EMP" {
		t.Errorf("expected code 'BPJS-HEALTH-EMP', got '%s'", found.RateCode)
	}

	// Find by Setting ID
	items, err := repo.FindBpjsRateComponentsBySettingID(ctx, created.BpjsSettingID)
	if err != nil {
		t.Fatalf("FindBpjsRateComponentsBySettingID failed: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 rate component, got %d", len(items))
	}

	// Update
	found.RateName = "Updated Rate"
	if err := repo.UpdateBpjsRateComponent(ctx, found); err != nil {
		t.Fatalf("UpdateBpjsRateComponent failed: %v", err)
	}

	// Delete
	if err := repo.DeleteBpjsRateComponent(ctx, created.ID); err != nil {
		t.Fatalf("DeleteBpjsRateComponent failed: %v", err)
	}
	_, err = repo.FindBpjsRateComponentByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting rate component")
	}
}

// =============================================================================
// PPh21 Setting Tests
// =============================================================================

func TestRepository_Pph21SettingCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestPph21Setting(ctx, repo)

	// Find
	found, err := repo.FindPph21SettingByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindPph21SettingByID failed: %v", err)
	}
	if found.SettingCode != "PPH21-DEFAULT" {
		t.Errorf("expected code 'PPH21-DEFAULT', got '%s'", found.SettingCode)
	}

	// Update
	found.SettingName = "Updated PPh21"
	if err := repo.UpdatePph21Setting(ctx, found); err != nil {
		t.Fatalf("UpdatePph21Setting failed: %v", err)
	}

	// Delete
	if err := repo.DeletePph21Setting(ctx, created.ID); err != nil {
		t.Fatalf("DeletePph21Setting failed: %v", err)
	}
	_, err = repo.FindPph21SettingByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting PPh21 setting")
	}
}

// =============================================================================
// PPh21 PTKP Rate Tests
// =============================================================================

func TestRepository_Pph21PtkpRateCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	pr := &Pph21PtkpRate{
		PtkpStatus:         "TK/0",
		AnnualAmount:       54000000,
		EffectiveStartDate: "2026-01-01",
		Status:             "ACTIVE",
	}
	if err := repo.CreatePph21PtkpRate(ctx, pr); err != nil {
		t.Fatalf("CreatePph21PtkpRate failed: %v", err)
	}

	found, err := repo.FindPph21PtkpRateByID(ctx, pr.ID)
	if err != nil {
		t.Fatalf("FindPph21PtkpRateByID failed: %v", err)
	}
	if found.AnnualAmount != 54000000 {
		t.Errorf("expected amount 54000000, got %f", found.AnnualAmount)
	}

	if err := repo.DeletePph21PtkpRate(ctx, pr.ID); err != nil {
		t.Fatalf("DeletePph21PtkpRate failed: %v", err)
	}
}

// =============================================================================
// PPh21 Tax Bracket Tests
// =============================================================================

func TestRepository_Pph21TaxBracketCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	tb := &Pph21TaxBracket{
		BracketOrder:       1,
		LowerBound:         0,
		UpperBound:         float64Ptr(60000000),
		RatePercent:        5,
		EffectiveStartDate: "2026-01-01",
		Status:             "ACTIVE",
	}
	if err := repo.CreatePph21TaxBracket(ctx, tb); err != nil {
		t.Fatalf("CreatePph21TaxBracket failed: %v", err)
	}

	found, err := repo.FindPph21TaxBracketByID(ctx, tb.ID)
	if err != nil {
		t.Fatalf("FindPph21TaxBracketByID failed: %v", err)
	}
	if found.RatePercent != 5 {
		t.Errorf("expected rate 5, got %f", found.RatePercent)
	}
}

// =============================================================================
// Payroll Run Tests
// =============================================================================

func TestRepository_PayrollRunCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestPayrollRun(ctx, repo)

	// Find
	found, err := repo.FindPayrollRunByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindPayrollRunByID failed: %v", err)
	}
	if found.RunCode != "RUN-2026-01" {
		t.Errorf("expected run_code 'RUN-2026-01', got '%s'", found.RunCode)
	}

	// Update
	found.Status = "CALCULATED"
	if err := repo.UpdatePayrollRun(ctx, found); err != nil {
		t.Fatalf("UpdatePayrollRun failed: %v", err)
	}

	// List
	items, total, err := repo.FindAllPayrollRuns(ctx, 1, 10)
	if err != nil {
		t.Fatalf("FindAllPayrollRuns failed: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 run, got %d", len(items))
	}
}

// =============================================================================
// Payroll Run Employee & Item Tests
// =============================================================================

func TestRepository_BulkCreateRunEmployees(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	run := createTestPayrollRun(ctx, repo)

	employees := []PayrollRunEmployee{
		{
			PayrollRunID: run.ID,
			EmployeeID:   uuid.New(),
			EmployeeCode: "EMP001",
			EmployeeName: "Employee One",
			Status:       "DRAFT",
		},
		{
			PayrollRunID: run.ID,
			EmployeeID:   uuid.New(),
			EmployeeCode: "EMP002",
			EmployeeName: "Employee Two",
			Status:       "DRAFT",
		},
	}

	if err := repo.BulkCreatePayrollRunEmployees(ctx, employees); err != nil {
		t.Fatalf("BulkCreatePayrollRunEmployees failed: %v", err)
	}

	items, err := repo.FindPayrollRunEmployeesByRunID(ctx, run.ID)
	if err != nil {
		t.Fatalf("FindPayrollRunEmployeesByRunID failed: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("expected 2 employees, got %d", len(items))
	}
}

func TestRepository_BulkCreateRunItems(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	run := createTestPayrollRun(ctx, repo)
	comp := createTestSalaryComponent(ctx, repo)

	runEmp := &PayrollRunEmployee{
		PayrollRunID: run.ID,
		EmployeeID:   uuid.New(),
		EmployeeCode: "EMP001",
		EmployeeName: "Test Employee",
		Status:       "DRAFT",
	}
	repo.CreatePayrollRunEmployee(ctx, runEmp)

	items := []PayrollRunItem{
		{
			PayrollRunID:         run.ID,
			PayrollRunEmployeeID: runEmp.ID,
			EmployeeID:           runEmp.EmployeeID,
			SalaryComponentID:    comp.ID,
			ComponentCode:        comp.Code,
			ComponentName:        comp.Name,
			ComponentType:        comp.ComponentType,
			ItemCategory:         "EMPLOYEE_EARNING",
			PaidBy:               "EMPLOYER",
			Amount:               5000000,
			CurrencyCode:         "IDR",
			SourceGroup:          "STRUCTURE",
		},
		{
			PayrollRunID:         run.ID,
			PayrollRunEmployeeID: runEmp.ID,
			EmployeeID:           runEmp.EmployeeID,
			SalaryComponentID:    comp.ID,
			ComponentCode:        "DEDUCT",
			ComponentName:        "Deduction",
			ComponentType:        "DEDUCTION",
			ItemCategory:         "EMPLOYEE_DEDUCTION",
			PaidBy:               "EMPLOYEE",
			Amount:               500000,
			CurrencyCode:         "IDR",
			SourceGroup:          "STRUCTURE",
		},
	}

	if err := repo.BulkCreatePayrollRunItems(ctx, items); err != nil {
		t.Fatalf("BulkCreatePayrollRunItems failed: %v", err)
	}

	runItems, err := repo.FindPayrollRunItemsByRunID(ctx, run.ID)
	if err != nil {
		t.Fatalf("FindPayrollRunItemsByRunID failed: %v", err)
	}

	if len(runItems) != 2 {
		t.Errorf("expected 2 items, got %d", len(runItems))
	}
}

// Helper for float64 pointer
func float64Ptr(f float64) *float64 {
	return &f
}


