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
		Code:            "TEST-001",
		Name:            "Test Salary Component",
		ComponentType:   "EARNING",
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
			Code:            fmt.Sprintf("COMP%03d", i+1),
			Name:            fmt.Sprintf("Component %d", i+1),
			ComponentType:   "EARNING",
			CalculationType: "FIXED",
			Status:          "ACTIVE",
			DisplayOrder:    100,
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

// TestRepository_FindActiveBpjsByDate: query effective-dated memilih setting &
// rate component yang ACTIVE dan berlaku pada tanggal tertentu, serta profil
// BPJS employee yang aktif per tanggal.
func TestRepository_FindActiveBpjsByDate(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	setting := createTestBpjsSetting(ctx, repo)

	// Berlaku pada Jan 2026.
	active, err := repo.FindActiveBpjsSettingByDate(ctx, "2026-01-31")
	if err != nil {
		t.Fatalf("FindActiveBpjsSettingByDate: %v", err)
	}
	if active == nil || active.ID != setting.ID {
		t.Fatalf("expected active setting %s, got %v", setting.ID, active)
	}

	// Tidak berlaku sebelum effective_start_date.
	before, err := repo.FindActiveBpjsSettingByDate(ctx, "2025-12-31")
	if err != nil {
		t.Fatalf("FindActiveBpjsSettingByDate(before): %v", err)
	}
	if before != nil {
		t.Fatalf("expected nil before effective date, got %v", before.ID)
	}

	// Rate component effective-dated (pakai setting yang sudah ada).
	comp := createTestSalaryComponent(ctx, repo)
	rate := createTestBpjsRateComponentLinked(ctx, repo, setting, comp, "HEALTH", "EMPLOYEE", "BPJS-REPO-EMP", 1.0)
	rates, err := repo.FindActiveBpjsRateComponentsBySettingID(ctx, setting.ID, "2026-01-31")
	if err != nil {
		t.Fatalf("FindActiveBpjsRateComponentsBySettingID: %v", err)
	}
	if len(rates) != 1 || rates[0].ID != rate.ID {
		t.Fatalf("expected 1 active rate, got %d", len(rates))
	}

	// Profil BPJS employee effective-dated.
	employee := createTestEmployee(ctx, repo, "EMP-BPJS-REPO", "Euis")
	profile := createTestPayrollProfile(ctx, repo, employee.ID, "2026-01-01")
	bpjsProfile := createTestEmployeeBpjsProfile(ctx, repo, employee.ID, profile.ID)
	found, err := repo.FindActiveEmployeeBpjsProfileByEmployeeID(ctx, employee.ID, "2026-01-31")
	if err != nil {
		t.Fatalf("FindActiveEmployeeBpjsProfileByEmployeeID: %v", err)
	}
	if found == nil || found.ID != bpjsProfile.ID {
		t.Fatalf("expected active BPJS profile %s, got %v", bpjsProfile.ID, found)
	}
}

// TestRepository_FindActivePph21ByDate: query effective-dated memilih setting,
// PTKP rate, tax bracket, dan profil pajak employee yang berlaku pada tanggal.
func TestRepository_FindActivePph21ByDate(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	comp := createTestSalaryComponent(ctx, repo)
	setting := createTestPph21SettingCustom(ctx, repo, comp)

	// Setting aktif pada Jan 2026.
	active, err := repo.FindActivePph21SettingByDate(ctx, "2026-01-31")
	if err != nil {
		t.Fatalf("FindActivePph21SettingByDate: %v", err)
	}
	if active == nil || active.ID != setting.ID {
		t.Fatalf("expected active setting %s, got %v", setting.ID, active)
	}

	// PTKP dari tabel ptkps (satu sumber kebenaran).
	createTestPtkp(ctx, repo, "TK0", "Tidak Kawin (TK/0)", 54000000, "A")
	ptkp, err := repo.FindPtkpByCode(ctx, "TK0")
	if err != nil {
		t.Fatalf("FindPtkpByCode: %v", err)
	}
	if ptkp == nil || ptkp.Ptkp != 54000000 || ptkp.Group != "A" {
		t.Fatalf("expected PTKP TK0 54000000 group A, got %+v", ptkp)
	}

	createTestPph21TaxBracket(ctx, repo, 1, 0, floatPtr(60000000), 5.0)
	brackets, err := repo.FindActivePph21TaxBracketsByDate(ctx, "2026-01-31")
	if err != nil {
		t.Fatalf("FindActivePph21TaxBracketsByDate: %v", err)
	}
	if len(brackets) != 1 || brackets[0].RatePercent != 5.0 {
		t.Fatalf("expected 1 bracket 5%%, got %d", len(brackets))
	}

	// Profil pajak employee effective-dated.
	employee := createTestEmployee(ctx, repo, "EMP-TAX-REPO", "Neneng")
	payrollProfile := createTestPayrollProfile(ctx, repo, employee.ID, "2026-01-01")
	profile := createTestEmployeeTaxProfile(ctx, repo, employee.ID, payrollProfile.ID, "K/1", true)
	found, err := repo.FindActiveEmployeeTaxProfileByEmployeeID(ctx, employee.ID, "2026-01-31")
	if err != nil {
		t.Fatalf("FindActiveEmployeeTaxProfileByEmployeeID: %v", err)
	}
	if found == nil || found.ID != profile.ID {
		t.Fatalf("expected active tax profile %s, got %v", profile.ID, found)
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

// =============================================================================
// Read models & Salary Structure Queries (dipakai kalkulasi run)
// =============================================================================

func TestRepository_FindActiveEmploymentByEmployeeID(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	employee := createTestEmployee(ctx, repo, "EMP001", "John")
	grading := createTestGrading(ctx, repo, "G-1", "Staff")
	position := createTestPosition(ctx, repo, "Staff HR", &grading.ID)
	createTestEmployment(ctx, repo, &employee.ID, &position.ID, "2026-01-01")

	// Berlaku pada Jan 2026.
	found, err := repo.FindActiveEmploymentByEmployeeID(ctx, employee.ID, "2026-01-31")
	if err != nil {
		t.Fatalf("FindActiveEmploymentByEmployeeID: %v", err)
	}
	if found == nil {
		t.Fatal("expected employment to be found")
	}
	if found.PositionID == nil || *found.PositionID != position.ID {
		t.Errorf("expected position %s, got %v", position.ID, found.PositionID)
	}

	// Tidak berlaku sebelum effective_date.
	found2, err := repo.FindActiveEmploymentByEmployeeID(ctx, employee.ID, "2025-12-31")
	if err != nil {
		t.Fatalf("FindActiveEmploymentByEmployeeID: %v", err)
	}
	if found2 != nil {
		t.Error("expected no employment before effective date")
	}
}

// TestRepository_FindEmploymentByEmployeeIDForPeriod: employment yang overlap
// dengan periode ditemukan walau resign tengah bulan (end_date < period end),
// dan employment yang berakhir sebelum periode tidak ikut.
func TestRepository_FindEmploymentByEmployeeIDForPeriod(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	employee := createTestEmployee(ctx, repo, "EMP-RESIGN-REPO", "Euis")
	grading := createTestGrading(ctx, repo, "G-1", "Staff")
	position := createTestPosition(ctx, repo, "Staff HR", &grading.ID)
	emp := createTestEmployment(ctx, repo, &employee.ID, &position.ID, "2020-01-01")
	db, _ := repo.getDB(ctx)
	if err := db.Model(&EmploymentRead{}).Where("id = ?", emp.ID).Update("effective_end_date", "2026-01-16").Error; err != nil {
		t.Fatalf("set end date: %v", err)
	}

	// Resign tengah bulan Januari → tetap ditemukan untuk periode Jan 2026.
	found, err := repo.FindEmploymentByEmployeeIDForPeriod(ctx, employee.ID, "2026-01-01", "2026-01-31")
	if err != nil {
		t.Fatalf("FindEmploymentByEmployeeIDForPeriod: %v", err)
	}
	if found == nil {
		t.Fatal("expected resigned employment to overlap January period")
	}

	// Employment yang sudah selesai sebelum periode → tidak ditemukan.
	none, err := repo.FindEmploymentByEmployeeIDForPeriod(ctx, employee.ID, "2026-03-01", "2026-03-31")
	if err != nil {
		t.Fatalf("FindEmploymentByEmployeeIDForPeriod(March): %v", err)
	}
	if none != nil {
		t.Error("expected no employment overlapping March period")
	}
}

func TestRepository_FindPositionAndGrading(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading := createTestGrading(ctx, repo, "G-2", "Supervisor")
	position := createTestPosition(ctx, repo, "Supervisor Ops", &grading.ID)

	foundPos, err := repo.FindPositionByID(ctx, position.ID)
	if err != nil {
		t.Fatalf("FindPositionByID: %v", err)
	}
	if foundPos.Title != "Supervisor Ops" {
		t.Errorf("expected title 'Supervisor Ops', got %q", foundPos.Title)
	}

	foundGrading, err := repo.FindGradingByID(ctx, grading.ID)
	if err != nil {
		t.Fatalf("FindGradingByID: %v", err)
	}
	if foundGrading.Name != "Supervisor" {
		t.Errorf("expected grading name 'Supervisor', got %q", foundGrading.Name)
	}
}

func TestRepository_FindSalaryStructureByGradingAndEmployee(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading := createTestGrading(ctx, repo, "G-1", "Staff")
	employee := createTestEmployee(ctx, repo, "EMP001", "John")

	comp := createTestSalaryComponent(ctx, repo)
	createTestGradeComponent(ctx, repo, grading.ID, comp.ID, 5000000)

	gradeComps, err := repo.FindAllSalaryGradeComponentsByGradingID(ctx, grading.ID, "2026-01-31")
	if err != nil {
		t.Fatalf("FindAllSalaryGradeComponentsByGradingID: %v", err)
	}
	if len(gradeComps) != 1 {
		t.Fatalf("expected 1 grade component, got %d", len(gradeComps))
	}

	empComps, err := repo.FindAllSalaryEmployeeComponentsByEmployeeID(ctx, employee.ID, "2026-01-31")
	if err != nil {
		t.Fatalf("FindAllSalaryEmployeeComponentsByEmployeeID: %v", err)
	}
	if len(empComps) != 0 {
		t.Errorf("expected 0 employee components, got %d", len(empComps))
	}

	adjustments, err := repo.FindAllSalaryEmployeeAdjustmentsByPeriod(ctx, employee.ID, 2026, 1)
	if err != nil {
		t.Fatalf("FindAllSalaryEmployeeAdjustmentsByPeriod: %v", err)
	}
	if len(adjustments) != 0 {
		t.Errorf("expected 0 adjustments, got %d", len(adjustments))
	}
}

func TestRepository_DeleteRunEmployeesAndItems(t *testing.T) {
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
	item := &PayrollRunItem{
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
	}
	if err := repo.BulkCreatePayrollRunItems(ctx, []PayrollRunItem{*item}); err != nil {
		t.Fatalf("BulkCreatePayrollRunItems: %v", err)
	}

	// Hapus item dulu (FK), lalu employee.
	if err := repo.DeletePayrollRunItemsByRunID(ctx, run.ID); err != nil {
		t.Fatalf("DeletePayrollRunItemsByRunID: %v", err)
	}
	if err := repo.DeletePayrollRunEmployeesByRunID(ctx, run.ID); err != nil {
		t.Fatalf("DeletePayrollRunEmployeesByRunID: %v", err)
	}

	emps, _ := repo.FindPayrollRunEmployeesByRunID(ctx, run.ID)
	items, _ := repo.FindPayrollRunItemsByRunID(ctx, run.ID)
	if len(emps) != 0 || len(items) != 0 {
		t.Errorf("expected empty snapshot after delete, got employees=%d items=%d", len(emps), len(items))
	}
}

// Helper for float64 pointer
func float64Ptr(f float64) *float64 {
	return &f
}
