package employee

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"

	"github.com/inthros/hris-platform/internal/modules/organization"
	"github.com/inthros/hris-platform/internal/modules/setting"
)

// =========================================================================
// Employee CRUD Tests
// =========================================================================

func TestRepository_CreateEmployee(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := &Employee{
		EmployeeID:      "EMP001",
		Name:            "John Doe",
		Gender:          strPtr("M"),
		PhoneNumber:     strPtr("08123456789"),
		Email:           strPtr("john@example.com"),
		Status:          "active",
	}

	if err := repo.CreateEmployee(ctx, emp); err != nil {
		t.Fatalf("CreateEmployee failed: %v", err)
	}

	if emp.ID == uuid.Nil {
		t.Error("expected employee ID to be generated")
	}
}

func TestRepository_FindEmployeeByID(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestEmployee(ctx, repo)

	found, err := repo.FindEmployeeByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindEmployeeByID failed: %v", err)
	}

	if found.Name != "Test Employee" {
		t.Errorf("expected name 'Test Employee', got '%s'", found.Name)
	}
	if found.Status != "active" {
		t.Errorf("expected status 'active', got '%s'", found.Status)
	}
}

func TestRepository_FindEmployeeByID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	_, err := repo.FindEmployeeByID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent employee")
	}
}

func TestRepository_FindEmployeeByEmployeeID(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	createTestEmployee(ctx, repo)

	found, err := repo.FindEmployeeByEmployeeID(ctx, "EMP-TEST-001")
	if err != nil {
		t.Fatalf("FindEmployeeByEmployeeID failed: %v", err)
	}

	if found.Name != "Test Employee" {
		t.Errorf("expected name 'Test Employee', got '%s'", found.Name)
	}
}

func TestRepository_FindAllEmployees_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	// Create 3 employees
	for i := 0; i < 3; i++ {
		emp := &Employee{
			EmployeeID: fmt.Sprintf("EMP%03d", i+1),
			Name:       fmt.Sprintf("Employee %d", i+1),
			Status:     "active",
		}
		if err := repo.CreateEmployee(ctx, emp); err != nil {
			t.Fatalf("failed to create employee %d: %v", i, err)
		}
	}

	// Test page 1 with per_page 2
	emps, total, err := repo.FindAllEmployees(ctx, 1, 2, "", "", "")
	if err != nil {
		t.Fatalf("FindAllEmployees failed: %v", err)
	}

	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(emps) != 2 {
		t.Errorf("expected 2 employees on page 1, got %d", len(emps))
	}

	// Test status filter
	empsActive, totalActive, err := repo.FindAllEmployees(ctx, 1, 20, "", "active", "")
	if err != nil {
		t.Fatalf("FindAllEmployees (status filter) failed: %v", err)
	}
	if totalActive != 3 || len(empsActive) != 3 {
		t.Errorf("expected 3 active employees, got total=%d len=%d", totalActive, len(empsActive))
	}
	empsInactive, totalInactive, err := repo.FindAllEmployees(ctx, 1, 20, "", "inactive", "")
	if err != nil {
		t.Fatalf("FindAllEmployees (inactive) failed: %v", err)
	}
	if totalInactive != 0 || len(empsInactive) != 0 {
		t.Errorf("expected 0 inactive employees, got total=%d len=%d", totalInactive, len(empsInactive))
	}
}

func TestRepository_UpdateEmployee(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestEmployee(ctx, repo)

	created.Name = "Updated Name"
	if err := repo.UpdateEmployee(ctx, created); err != nil {
		t.Fatalf("UpdateEmployee failed: %v", err)
	}

	found, _ := repo.FindEmployeeByID(ctx, created.ID)
	if found.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got '%s'", found.Name)
	}
}

func TestRepository_DeleteEmployee(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestEmployee(ctx, repo)

	if err := repo.DeleteEmployee(ctx, created.ID); err != nil {
		t.Fatalf("DeleteEmployee failed: %v", err)
	}

	// Verify it's gone
	_, err := repo.FindEmployeeByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting employee")
	}
}

// =========================================================================
// Address Sub-module Tests
// =========================================================================

func TestRepository_CreateAddress(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)

	addr := &EmployeeAddress{
		EmployeeID: &emp.ID,
		Type:       strPtr("MAIN"),
		Address:    strPtr("Jl. Test No. 1"),
	}
	if err := repo.CreateAddress(ctx, addr); err != nil {
		t.Fatalf("CreateAddress failed: %v", err)
	}

	if addr.ID == uuid.Nil {
		t.Error("expected address ID to be generated")
	}
}

func TestRepository_FindAddressByID(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)
	addr := &EmployeeAddress{
		EmployeeID: &emp.ID,
		Type:       strPtr("MAIN"),
		Address:    strPtr("Jl. Test No. 1"),
	}
	repo.CreateAddress(ctx, addr)

	found, err := repo.FindAddressByID(ctx, addr.ID)
	if err != nil {
		t.Fatalf("FindAddressByID failed: %v", err)
	}
	if *found.Type != "MAIN" {
		t.Errorf("expected type 'MAIN', got '%s'", *found.Type)
	}
}

func TestRepository_UpdateAddress(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)
	addr := &EmployeeAddress{
		EmployeeID: &emp.ID,
		Type:       strPtr("MAIN"),
		Address:    strPtr("Jl. Test No. 1"),
	}
	repo.CreateAddress(ctx, addr)

	addr.Type = strPtr("DOMICILE")
	repo.UpdateAddress(ctx, addr)

	found, _ := repo.FindAddressByID(ctx, addr.ID)
	if *found.Type != "DOMICILE" {
		t.Errorf("expected type 'DOMICILE', got '%s'", *found.Type)
	}
}

func TestRepository_DeleteAddress(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)
	addr := &EmployeeAddress{
		EmployeeID: &emp.ID,
		Type:       strPtr("MAIN"),
		Address:    strPtr("Jl. Test No. 1"),
	}
	repo.CreateAddress(ctx, addr)

	repo.DeleteAddress(ctx, addr.ID)

	_, err := repo.FindAddressByID(ctx, addr.ID)
	if err == nil {
		t.Fatal("expected error after deleting address")
	}
}

// =========================================================================
// Emergency Contact Sub-module Tests
// =========================================================================

func TestRepository_EmergencyContactCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)

	// Create
	contact := &EmergencyContact{
		EmployeeID:  &emp.ID,
		Name:        "Emergency Contact",
		PhoneNumber: "08111111111",
	}
	if err := repo.CreateEmergencyContact(ctx, contact); err != nil {
		t.Fatalf("CreateEmergencyContact failed: %v", err)
	}

	// Find
	found, err := repo.FindEmergencyContactByID(ctx, contact.ID)
	if err != nil {
		t.Fatalf("FindEmergencyContactByID failed: %v", err)
	}
	if found.Name != "Emergency Contact" {
		t.Errorf("expected name 'Emergency Contact', got '%s'", found.Name)
	}

	// Update
	found.Name = "Updated Contact"
	repo.UpdateEmergencyContact(ctx, found)
	updated, _ := repo.FindEmergencyContactByID(ctx, contact.ID)
	if updated.Name != "Updated Contact" {
		t.Errorf("expected name 'Updated Contact', got '%s'", updated.Name)
	}

	// Delete
	repo.DeleteEmergencyContact(ctx, contact.ID)
	_, err = repo.FindEmergencyContactByID(ctx, contact.ID)
	if err == nil {
		t.Fatal("expected error after deleting emergency contact")
	}
}

// =========================================================================
// Family Sub-module Tests
// =========================================================================

func TestRepository_FamilyCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)

	// Create
	fam := &EmployeeFamily{
		EmployeeID: &emp.ID,
		Name:       "Family Member",
	}
	if err := repo.CreateFamily(ctx, fam); err != nil {
		t.Fatalf("CreateFamily failed: %v", err)
	}

	// Find
	found, err := repo.FindFamilyByID(ctx, fam.ID)
	if err != nil {
		t.Fatalf("FindFamilyByID failed: %v", err)
	}
	if found.Name != "Family Member" {
		t.Errorf("expected name 'Family Member', got '%s'", found.Name)
	}

	// Update
	found.Name = "Updated Family"
	repo.UpdateFamily(ctx, found)
	updated, _ := repo.FindFamilyByID(ctx, fam.ID)
	if updated.Name != "Updated Family" {
		t.Errorf("expected name 'Updated Family', got '%s'", updated.Name)
	}

	// Delete
	repo.DeleteFamily(ctx, fam.ID)
	_, err = repo.FindFamilyByID(ctx, fam.ID)
	if err == nil {
		t.Fatal("expected error after deleting family")
	}
}

// =========================================================================
// Education Sub-module Tests
// =========================================================================

func TestRepository_EducationCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)

	edu := &EmployeeEducation{
		EmployeeID: &emp.ID,
		Name:       "S1 Computer Science",
		Major:      strPtr("Computer Science"),
		GradYear:   intPtr(2020),
	}
	if err := repo.CreateEducation(ctx, edu); err != nil {
		t.Fatalf("CreateEducation failed: %v", err)
	}

	// CRUD cycle
	found, _ := repo.FindEducationByID(ctx, edu.ID)
	if found.Name != "S1 Computer Science" {
		t.Errorf("expected name 'S1 Computer Science', got '%s'", found.Name)
	}

	found.Name = "S2 Data Science"
	repo.UpdateEducation(ctx, found)
	updated, _ := repo.FindEducationByID(ctx, edu.ID)
	if updated.Name != "S2 Data Science" {
		t.Errorf("expected name 'S2 Data Science', got '%s'", updated.Name)
	}

	repo.DeleteEducation(ctx, edu.ID)
	_, err := repo.FindEducationByID(ctx, edu.ID)
	if err == nil {
		t.Fatal("expected error after deleting education")
	}
}

// =========================================================================
// Experience Sub-module Tests
// =========================================================================

func TestRepository_ExperienceCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)

	exp := &EmployeeExperience{
		EmployeeID: &emp.ID,
		Company:    "PT Test",
		Position:   strPtr("Manager"),
		StartYear:  intPtr(2020),
		EndYear:    intPtr(2024),
	}
	if err := repo.CreateExperience(ctx, exp); err != nil {
		t.Fatalf("CreateExperience failed: %v", err)
	}

	found, _ := repo.FindExperienceByID(ctx, exp.ID)
	if found.Company != "PT Test" {
		t.Errorf("expected company 'PT Test', got '%s'", found.Company)
	}

	// Update
	found.Company = "PT Updated"
	repo.UpdateExperience(ctx, found)
	updated, _ := repo.FindExperienceByID(ctx, exp.ID)
	if updated.Company != "PT Updated" {
		t.Errorf("expected company 'PT Updated', got '%s'", updated.Company)
	}

	// Delete
	repo.DeleteExperience(ctx, exp.ID)
	_, err := repo.FindExperienceByID(ctx, exp.ID)
	if err == nil {
		t.Fatal("expected error after deleting experience")
	}
}

// =========================================================================
// Document Sub-module Tests
// =========================================================================

func TestRepository_DocumentCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)

	doc := &EmployeeDocument{
		EmployeeID: &emp.ID,
		Name:       "CV",
		File:       "cv.pdf",
		Note:       strPtr("Original CV"),
	}
	if err := repo.CreateDocument(ctx, doc); err != nil {
		t.Fatalf("CreateDocument failed: %v", err)
	}

	found, _ := repo.FindDocumentByID(ctx, doc.ID)
	if found.Name != "CV" {
		t.Errorf("expected name 'CV', got '%s'", found.Name)
	}

	found.Name = "Updated CV"
	repo.UpdateDocument(ctx, found)
	updated, _ := repo.FindDocumentByID(ctx, doc.ID)
	if updated.Name != "Updated CV" {
		t.Errorf("expected name 'Updated CV', got '%s'", updated.Name)
	}

	repo.DeleteDocument(ctx, doc.ID)
	_, err := repo.FindDocumentByID(ctx, doc.ID)
	if err == nil {
		t.Fatal("expected error after deleting document")
	}
}

// =========================================================================
// Insurance Sub-module Tests
// =========================================================================

func TestRepository_InsuranceCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)

	insID := uuid.New()
	ins := &EmployeeInsurance{
		EmployeeID:  &emp.ID,
		InsuranceID: &insID,
		Number:      "BPJS001",
	}
	if err := repo.CreateInsurance(ctx, ins); err != nil {
		t.Fatalf("CreateInsurance failed: %v", err)
	}

	found, _ := repo.FindInsuranceByID(ctx, ins.ID)
	if found.Number != "BPJS001" {
		t.Errorf("expected number 'BPJS001', got '%s'", found.Number)
	}

	found.Number = "BPJS002"
	repo.UpdateInsurance(ctx, found)
	updated, _ := repo.FindInsuranceByID(ctx, ins.ID)
	if updated.Number != "BPJS002" {
		t.Errorf("expected number 'BPJS002', got '%s'", updated.Number)
	}

	repo.DeleteInsurance(ctx, ins.ID)
	_, err := repo.FindInsuranceByID(ctx, ins.ID)
	if err == nil {
		t.Fatal("expected error after deleting insurance")
	}
}

// =========================================================================
// Employment Sub-module Tests
// =========================================================================

func TestRepository_EmploymentCRUD(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	emp := createTestEmployee(ctx, repo)

	empl := &Employment{
		EmployeeID:           &emp.ID,
		DecisionLetterNumber: "SK-001",
		DecisionLetterDate:   "2024-01-01",
		EffectiveDate:        "2024-01-15",
	}
	if err := repo.CreateEmployment(ctx, empl); err != nil {
		t.Fatalf("CreateEmployment failed: %v", err)
	}

	found, _ := repo.FindEmploymentByID(ctx, empl.ID)
	if found.DecisionLetterNumber != "SK-001" {
		t.Errorf("expected SK 'SK-001', got '%s'", found.DecisionLetterNumber)
	}

	found.DecisionLetterNumber = "SK-002"
	repo.UpdateEmployment(ctx, found)
	updated, _ := repo.FindEmploymentByID(ctx, empl.ID)
	if updated.DecisionLetterNumber != "SK-002" {
		t.Errorf("expected SK 'SK-002', got '%s'", updated.DecisionLetterNumber)
	}

	repo.DeleteEmployment(ctx, empl.ID)
	_, err := repo.FindEmploymentByID(ctx, empl.ID)
	if err == nil {
		t.Fatal("expected error after deleting employment")
	}
}

func TestRepository_CountByGender(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	// 2 laki-laki, 1 perempuan, 1 tanpa gender — semuanya aktif.
	for _, g := range []*string{strPtr("M"), strPtr("M"), strPtr("F"), nil} {
		emp := &Employee{EmployeeID: uuid.NewString(), Name: "X", Status: "active", Gender: g}
		if err := repo.CreateEmployee(ctx, emp); err != nil {
			t.Fatalf("failed to create employee: %v", err)
		}
	}
	// 1 karyawan non-aktif dengan gender M — TIDAK dihitung.
	inactive := &Employee{EmployeeID: uuid.NewString(), Name: "Y", Status: "inactive", Gender: strPtr("M")}
	if err := repo.CreateEmployee(ctx, inactive); err != nil {
		t.Fatalf("failed to create inactive employee: %v", err)
	}

	male, female, other, err := repo.CountByGender(ctx)
	if err != nil {
		t.Fatalf("CountByGender failed: %v", err)
	}
	if male != 2 {
		t.Errorf("expected 2 male, got %d", male)
	}
	if female != 1 {
		t.Errorf("expected 1 female, got %d", female)
	}
	if other != 1 {
		t.Errorf("expected 1 other/unknown, got %d", other)
	}
	_ = db
}

func TestRepository_CountByEmploymentStatus(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	// 2 status: Tetap (PKWTT) & Kontrak (PKWT).
	tetapID := uuid.New()
	kontrakID := uuid.New()
	if err := db.Create(&setting.EmploymentStatus{ID: tetapID, Code: "PKWTT", Name: "Tetap"}).Error; err != nil {
		t.Fatalf("failed to create status: %v", err)
	}
	if err := db.Create(&setting.EmploymentStatus{ID: kontrakID, Code: "PKWT", Name: "Kontrak"}).Error; err != nil {
		t.Fatalf("failed to create status: %v", err)
	}

	// 2 karyawan Tetap (employment berjalan), 1 Kontrak, 1 tanpa employment berjalan.
	emp1 := createTestEmployee(ctx, repo)
	emp2 := createTestEmployee(ctx, repo)
	emp3 := createTestEmployee(ctx, repo)
	emp4 := createTestEmployee(ctx, repo)
	// Karyawan non-aktif dengan employment berjalan — TIDAK dihitung.
	inactive := &Employee{EmployeeID: uuid.NewString(), Name: "Inactive", Status: "inactive"}
	if err := repo.CreateEmployee(ctx, inactive); err != nil {
		t.Fatalf("failed to create inactive employee: %v", err)
	}

	for _, empl := range []*Employment{
		{EmployeeID: &emp1.ID, OrganizationID: &uuid.Nil, EmploymentStatusID: &tetapID, DecisionLetterNumber: "SK-1", DecisionLetterDate: "2024-01-01", EffectiveDate: "2024-01-01"},
		{EmployeeID: &emp2.ID, OrganizationID: &uuid.Nil, EmploymentStatusID: &tetapID, DecisionLetterNumber: "SK-2", DecisionLetterDate: "2024-01-01", EffectiveDate: "2024-01-01"},
		{EmployeeID: &emp3.ID, OrganizationID: &uuid.Nil, EmploymentStatusID: &kontrakID, DecisionLetterNumber: "SK-3", DecisionLetterDate: "2024-01-01", EffectiveDate: "2024-01-01"},
		// emp4: employment sudah berakhir (bukan berjalan) → tidak terhitung di grup.
		{EmployeeID: &emp4.ID, OrganizationID: &uuid.Nil, EmploymentStatusID: &kontrakID, DecisionLetterNumber: "SK-4", DecisionLetterDate: "2024-01-01", EffectiveDate: "2024-01-01", EffectiveEndDate: strPtr("2024-06-30")},
		// inactive: employment berjalan tapi karyawan non-aktif → tidak dihitung.
		{EmployeeID: &inactive.ID, OrganizationID: &uuid.Nil, EmploymentStatusID: &tetapID, DecisionLetterNumber: "SK-5", DecisionLetterDate: "2024-01-01", EffectiveDate: "2024-01-01"},
	} {
		if err := repo.CreateEmployment(ctx, empl); err != nil {
			t.Fatalf("failed to create employment: %v", err)
		}
	}

	groups, unclassified, err := repo.CountByEmploymentStatus(ctx)
	if err != nil {
		t.Fatalf("CountByEmploymentStatus failed: %v", err)
	}
	byName := map[string]int64{}
	for _, g := range groups {
		byName[g.Name] = g.Count
	}
	if byName["Tetap"] != 2 {
		t.Errorf("expected 2 Tetap, got %v", byName)
	}
	if byName["Kontrak"] != 1 {
		t.Errorf("expected 1 Kontrak, got %v", byName)
	}
	if unclassified != 1 {
		t.Errorf("expected 1 unclassified employee (no open employment), got %d", unclassified)
	}
}

func TestRepository_ResolveEmployeeRefNames(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()

	// Tabel referensi tidak ikut AutoMigrate di setupTestDB — buat di sini.
	if err := db.AutoMigrate(&setting.Religion{}, &setting.MaritalStatus{}, &setting.Nationality{}); err != nil {
		t.Fatalf("failed to migrate setting ref tables: %v", err)
	}

	religionID := uuid.New()
	maritalID := uuid.New()
	if err := db.Create(&setting.Religion{ID: religionID, Code: "ISL", Name: "Islam"}).Error; err != nil {
		t.Fatalf("failed to create religion: %v", err)
	}
	if err := db.Create(&setting.MaritalStatus{ID: maritalID, Code: "KWN", Name: "Kawin"}).Error; err != nil {
		t.Fatalf("failed to create marital status: %v", err)
	}
	if err := db.Create(&setting.Nationality{ID: uuid.New(), Code: "ID", Name: "Indonesia"}).Error; err != nil {
		t.Fatalf("failed to create nationality: %v", err)
	}

	repo := NewRepository(dbResolver)
	ctx := context.Background()

	religionName, maritalName, nationalityName := repo.ResolveEmployeeRefNames(ctx, &religionID, &maritalID, strPtr("ID"))
	if religionName != "Islam" {
		t.Errorf("expected religion name 'Islam', got '%s'", religionName)
	}
	if maritalName != "Kawin" {
		t.Errorf("expected marital status name 'Kawin', got '%s'", maritalName)
	}
	if nationalityName != "Indonesia" {
		t.Errorf("expected nationality name 'Indonesia', got '%s'", nationalityName)
	}

	// Referensi tidak ditemukan → string kosong (bukan ID mentah).
	unknownID := uuid.New()
	if name, _, _ := repo.ResolveEmployeeRefNames(ctx, &unknownID, nil, nil); name != "" {
		t.Errorf("expected empty religion name for unknown id, got '%s'", name)
	}
}

func TestRepository_FindEmployeeByID_ResolvesSubRefNames(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	db, _ := dbResolver(context.Background())

	// Referensi settings + organization.
	relTypeID := uuid.New()
	eduID := uuid.New()
	insID := uuid.New()
	bankID := uuid.New()
	empStatusID := uuid.New()
	orgID := uuid.New()
	for _, rec := range []map[string]interface{}{
		{"model": &setting.RelationshipType{ID: relTypeID, Code: "IST", Name: "Istri"}},
		{"model": &setting.Education{ID: eduID, Code: "S1", Name: "Strata 1"}},
		{"model": &setting.Insurance{ID: insID, Code: "01", Name: "BPJS Kesehatan"}},
		{"model": &setting.Bank{ID: bankID, Code: "BCA", Name: "Bank BCA"}},
		{"model": &setting.EmploymentStatus{ID: empStatusID, Code: "PKWTT", Name: "Pegawai Tetap"}},
	} {
		if err := db.Create(rec["model"]).Error; err != nil {
			t.Fatalf("failed to create setting ref: %v", err)
		}
	}
	if err := db.Create(&organization.Organization{ID: orgID, Code: "HQ", FullCode: "HQ-01", Nomenclature: "Head Office"}).Error; err != nil {
		t.Fatalf("failed to create organization: %v", err)
	}

	repo := NewRepository(dbResolver)
	ctx := context.Background()
	emp := createTestEmployee(ctx, repo)

	if err := repo.CreateEmergencyContact(ctx, &EmergencyContact{
		EmployeeID:         &emp.ID,
		Name:               "Siti",
		PhoneNumber:        "0812",
		RelationshipTypeID: &relTypeID,
	}); err != nil {
		t.Fatalf("CreateEmergencyContact failed: %v", err)
	}
	if err := repo.CreateFamily(ctx, &EmployeeFamily{
		EmployeeID:         &emp.ID,
		Name:               "Ayah",
		RelationshipTypeID: &relTypeID,
		EducationID:        &eduID,
	}); err != nil {
		t.Fatalf("CreateFamily failed: %v", err)
	}
	if err := repo.CreateEducation(ctx, &EmployeeEducation{
		EmployeeID:  &emp.ID,
		EducationID: &eduID,
		Name:        "",
	}); err != nil {
		t.Fatalf("CreateEducation failed: %v", err)
	}
	if err := repo.CreateInsurance(ctx, &EmployeeInsurance{EmployeeID: &emp.ID, InsuranceID: &insID, Number: "0001"}); err != nil {
		t.Fatalf("CreateInsurance failed: %v", err)
	}
	if err := repo.CreateBank(ctx, &EmployeeBankAccount{EmployeeID: &emp.ID, BankID: &bankID, AccountNumber: "123", AccountName: "Siti"}); err != nil {
		t.Fatalf("CreateBank failed: %v", err)
	}
	if err := repo.CreateEmployment(ctx, &Employment{
		EmployeeID:         &emp.ID,
		OrganizationID:     &orgID,
		EmploymentStatusID: &empStatusID,
		DecisionLetterNumber: "SK-001",
		DecisionLetterDate:   "2024-01-01",
		EffectiveDate:        "2024-01-01",
	}); err != nil {
		t.Fatalf("CreateEmployment failed: %v", err)
	}

	found, err := repo.FindEmployeeByID(ctx, emp.ID)
	if err != nil {
		t.Fatalf("FindEmployeeByID failed: %v", err)
	}
	resp := found.ToResponse()

	if len(resp.EmergencyContacts) != 1 || resp.EmergencyContacts[0].RelationshipTypeName != "Istri" {
		t.Errorf("expected contact relationship type name 'Istri', got %+v", resp.EmergencyContacts)
	}
	if len(resp.Families) != 1 || resp.Families[0].RelationshipTypeName != "Istri" || resp.Families[0].EducationName != "Strata 1" {
		t.Errorf("expected family names 'Istri'/'Strata 1', got %+v", resp.Families)
	}
	if len(resp.Educations) != 1 || resp.Educations[0].EducationName != "Strata 1" {
		t.Errorf("expected education name 'Strata 1', got %+v", resp.Educations)
	}
	if len(resp.Insurances) != 1 || resp.Insurances[0].InsuranceName != "BPJS Kesehatan" {
		t.Errorf("expected insurance name 'BPJS Kesehatan', got %+v", resp.Insurances)
	}
	if len(resp.Banks) != 1 || resp.Banks[0].BankName != "Bank BCA" {
		t.Errorf("expected bank name 'Bank BCA', got %+v", resp.Banks)
	}
	if len(resp.Employments) != 1 || resp.Employments[0].OrganizationName != "Head Office" || resp.Employments[0].EmploymentStatusName != "Pegawai Tetap" {
		t.Errorf("expected employment names 'Head Office'/'Pegawai Tetap', got %+v", resp.Employments)
	}
}
