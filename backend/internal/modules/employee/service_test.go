package employee

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/crypto"
)

// =========================================================================
// Service Tests (using real SQLite repository for true integration)
// =========================================================================

func newTestService() (*Service, func()) {
	_, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	return svc, func() { cleanup(); logger.Sync() }
}

// =========================================================================
// Quota Tests (max_employees on-premise enforcement)
// =========================================================================

// fakeQuotaChecker adalah implementasi test dari EmployeeQuotaChecker.
type fakeQuotaChecker struct {
	max int
}

func (f fakeQuotaChecker) MaxEmployees() int { return f.max }

// TestService_Create_QuotaUnlimitedSaaS — tanpa checker (mode saas) → selalu lolos.
func TestService_Create_QuotaUnlimitedSaaS(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	// 25 employee — jauh di atas batas umum, tapi saas tanpa quota → semua berhasil
	for i := 0; i < 25; i++ {
		_, err := svc.Create(ctx, CreateEmployeeRequest{
			EmployeeID: fmt.Sprintf("QUOTA-SAAS-%03d", i),
			Name:       fmt.Sprintf("Quota SaaS %d", i),
		})
		if err != nil {
			t.Fatalf("create #%d failed: %v", i, err)
		}
	}
}

// TestService_Create_QuotaMaxZero — max <= 0 berarti unlimited.
func TestService_Create_QuotaMaxZero(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetQuotaChecker(fakeQuotaChecker{max: 0})
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		_, err := svc.Create(ctx, CreateEmployeeRequest{
			EmployeeID: fmt.Sprintf("QUOTA-ZERO-%03d", i),
			Name:       fmt.Sprintf("Quota Zero %d", i),
		})
		if err != nil {
			t.Fatalf("create #%d failed: %v", i, err)
		}
	}
}

// TestService_Create_QuotaBelowLimit — jumlah < max → berhasil.
func TestService_Create_QuotaBelowLimit(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetQuotaChecker(fakeQuotaChecker{max: 5})
	ctx := context.Background()

	for i := 0; i < 4; i++ {
		_, err := svc.Create(ctx, CreateEmployeeRequest{
			EmployeeID: fmt.Sprintf("QUOTA-BELOW-%03d", i),
			Name:       fmt.Sprintf("Quota Below %d", i),
		})
		if err != nil {
			t.Fatalf("create #%d failed: %v", i, err)
		}
	}
}

// TestService_Create_QuotaAtLimit — saat count == max, pembuatan berikutnya ditolak.
func TestService_Create_QuotaAtLimit(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetQuotaChecker(fakeQuotaChecker{max: 3})
	ctx := context.Background()

	// Isi sampai batas (3)
	for i := 0; i < 3; i++ {
		_, err := svc.Create(ctx, CreateEmployeeRequest{
			EmployeeID: fmt.Sprintf("QUOTA-AT-%03d", i),
			Name:       fmt.Sprintf("Quota At %d", i),
		})
		if err != nil {
			t.Fatalf("create #%d failed: %v", i, err)
		}
	}

	// Ke-4 ditolak
	_, err := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "QUOTA-AT-004",
		Name:       "Quota At 4",
	})
	if !errors.Is(err, ErrQuotaExceeded) {
		t.Fatalf("expected ErrQuotaExceeded, got: %v", err)
	}
}

// TestService_Create_QuotaExceeded — count sudah melebihi max → langsung ditolak
// bahkan untuk data employee yang berbeda (count diperhitungkan global).
func TestService_Create_QuotaExceeded(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetQuotaChecker(fakeQuotaChecker{max: 2})
	ctx := context.Background()

	svc.Create(ctx, CreateEmployeeRequest{EmployeeID: "QUOTA-EX-001", Name: "Ex 1"})
	svc.Create(ctx, CreateEmployeeRequest{EmployeeID: "QUOTA-EX-002", Name: "Ex 2"})

	_, err := svc.Create(ctx, CreateEmployeeRequest{EmployeeID: "QUOTA-EX-003", Name: "Ex 3"})
	if !errors.Is(err, ErrQuotaExceeded) {
		t.Fatalf("expected ErrQuotaExceeded, got: %v", err)
	}
}

func TestService_CreateEmployee_DefaultStatus(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateEmployeeRequest{
		EmployeeID: "SVC001",
		Name:       "Service Test",
	}

	resp, err := svc.Create(ctx, req)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if resp.Status != "active" {
		t.Errorf("expected status 'active', got '%s'", resp.Status)
	}
	if resp.Name != "Service Test" {
		t.Errorf("expected name 'Service Test', got '%s'", resp.Name)
	}
}

func TestService_CreateEmployee_WithOptionalFields(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	gender := "M"
	req := CreateEmployeeRequest{
		EmployeeID:  "SVC002",
		Name:        "Full Employee",
		NIK:         strPtr("1234567890123456"),
		Gender:      &gender,
		PhoneNumber: strPtr("08123456789"),
		Email:       strPtr("full@example.com"),
	}

	resp, err := svc.Create(ctx, req)
	if err != nil {
		t.Fatalf("Create with optional fields failed: %v", err)
	}

	if resp.NIK != "1234567890123456" {
		t.Errorf("expected NIK '1234567890123456', got '%s'", resp.NIK)
	}
	if resp.Gender != "M" {
		t.Errorf("expected Gender 'M', got '%s'", resp.Gender)
	}
	if resp.Email != "full@example.com" {
		t.Errorf("expected email 'full@example.com', got '%s'", resp.Email)
	}
}

func TestService_GetByID_Success(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateEmployeeRequest{
		EmployeeID: "SVC003",
		Name:       "Get By ID Test",
	}
	created, _ := svc.Create(ctx, req)

	found, err := svc.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if found.Name != "Get By ID Test" {
		t.Errorf("expected name 'Get By ID Test', got '%s'", found.Name)
	}
}

func TestService_GetByID_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetByID(ctx, "not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_GetByID_NotFound(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetByID(ctx, uuid.New().String())
	if err == nil {
		t.Fatal("expected error for non-existent employee")
	}
}

func TestService_List_DefaultPagination(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	// Create employees
	for i := 0; i < 5; i++ {
		svc.Create(ctx, CreateEmployeeRequest{
			EmployeeID: fmt.Sprintf("LST%03d", i+1),
			Name:       fmt.Sprintf("List Employee %d", i+1),
		})
	}

	// Test with invalid params (should use defaults)
	resp, err := svc.List(ctx, 0, 0, "", "", "")
	if err != nil {
		t.Fatalf("List failed: %v", err)
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
}

func TestService_UpdateEmployee(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC004",
		Name:       "Before Update",
	})

	newName := "After Update"
	updated, err := svc.Update(ctx, created.ID, UpdateEmployeeRequest{
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if updated.Name != "After Update" {
		t.Errorf("expected name 'After Update', got '%s'", updated.Name)
	}
}

func TestService_DeleteEmployee(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC005",
		Name:       "To Delete",
	})

	if err := svc.Delete(ctx, created.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deleted
	_, err := svc.GetByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting employee")
	}
}

// =========================================================================
// Sub-module Service Tests
// =========================================================================

func TestService_CreateAddress(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	emp, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC010",
		Name:       "Address Owner",
	})

	addr, err := svc.CreateAddress(ctx, emp.ID, CreateAddressRequest{
		Type:    "MAIN",
		Address: "Jl. Service Test No. 1",
	})
	if err != nil {
		t.Fatalf("CreateAddress failed: %v", err)
	}

	if addr.Type != "MAIN" {
		t.Errorf("expected type 'MAIN', got '%s'", addr.Type)
	}
}

func TestService_CreateEmergencyContact(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	emp, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC011",
		Name:       "Contact Owner",
	})

	contact, err := svc.CreateEmergencyContact(ctx, emp.ID, CreateEmergencyContactRequest{
		Name:        "Emergency Contact",
		PhoneNumber: "08111111111",
	})
	if err != nil {
		t.Fatalf("CreateEmergencyContact failed: %v", err)
	}

	if contact.Name != "Emergency Contact" {
		t.Errorf("expected name 'Emergency Contact', got '%s'", contact.Name)
	}
}

func TestService_CreateFamily(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	emp, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC012",
		Name:       "Family Owner",
	})

	fam, err := svc.CreateFamily(ctx, emp.ID, CreateFamilyRequest{
		Name: "Family Member",
	})
	if err != nil {
		t.Fatalf("CreateFamily failed: %v", err)
	}

	if fam.Name != "Family Member" {
		t.Errorf("expected name 'Family Member', got '%s'", fam.Name)
	}
}

func TestService_CreateEducation(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	emp, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC013",
		Name:       "Education Owner",
	})

	edu, err := svc.CreateEducation(ctx, emp.ID, CreateEducationRequest{
		Name: "S1 Informatics",
	})
	if err != nil {
		t.Fatalf("CreateEducation failed: %v", err)
	}

	if edu.Name != "S1 Informatics" {
		t.Errorf("expected name 'S1 Informatics', got '%s'", edu.Name)
	}
}

func TestService_CreateExperience(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	emp, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC014",
		Name:       "Experience Owner",
	})

	exp, err := svc.CreateExperience(ctx, emp.ID, CreateExperienceRequest{
		Company: "PT Test Corp",
	})
	if err != nil {
		t.Fatalf("CreateExperience failed: %v", err)
	}

	if exp.Company != "PT Test Corp" {
		t.Errorf("expected company 'PT Test Corp', got '%s'", exp.Company)
	}
}

func TestService_CreateDocument(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	emp, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC015",
		Name:       "Document Owner",
	})

	doc, err := svc.CreateDocument(ctx, emp.ID, CreateDocumentRequest{
		Name: "CV",
		File: "cv.pdf",
	})
	if err != nil {
		t.Fatalf("CreateDocument failed: %v", err)
	}

	if doc.Name != "CV" {
		t.Errorf("expected name 'CV', got '%s'", doc.Name)
	}
}

func TestService_CreateInsurance(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	emp, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC016",
		Name:       "Insurance Owner",
	})

	ins, err := svc.CreateInsurance(ctx, emp.ID, CreateInsuranceRequest{
		Number: "BPJS001",
	})
	if err != nil {
		t.Fatalf("CreateInsurance failed: %v", err)
	}

	if ins.Number != "BPJS001" {
		t.Errorf("expected number 'BPJS001', got '%s'", ins.Number)
	}
}

func TestService_CreateEmployment(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	emp, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC017",
		Name:       "Employment Owner",
	})

	empl, err := svc.CreateEmployment(ctx, emp.ID, CreateEmploymentRequest{
		DecisionLetterNumber: "SK-001",
		DecisionLetterDate:   "2024-01-01",
		EffectiveDate:        "2024-01-15",
	})
	if err != nil {
		t.Fatalf("CreateEmployment failed: %v", err)
	}

	if empl.DecisionLetterNumber != "SK-001" {
		t.Errorf("expected SK 'SK-001', got '%s'", empl.DecisionLetterNumber)
	}
}

func TestService_GetEmployeeWithSubModules(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	emp, _ := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "SVC020",
		Name:       "Full Data Employee",
	})

	// Add address
	svc.CreateAddress(ctx, emp.ID, CreateAddressRequest{
		Type:    "MAIN",
		Address: "Jl. Test",
	})

	// Add family
	svc.CreateFamily(ctx, emp.ID, CreateFamilyRequest{
		Name: "Spouse",
	})

	// Add document
	svc.CreateDocument(ctx, emp.ID, CreateDocumentRequest{
		Name: "ID Card",
		File: "ktp.pdf",
	})

	// Fetch employee with all sub-modules
	fullEmp, err := svc.GetByID(ctx, emp.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if len(fullEmp.Addresses) != 1 {
		t.Errorf("expected 1 address, got %d", len(fullEmp.Addresses))
	}
	if len(fullEmp.Families) != 1 {
		t.Errorf("expected 1 family, got %d", len(fullEmp.Families))
	}
	if len(fullEmp.Documents) != 1 {
		t.Errorf("expected 1 document, got %d", len(fullEmp.Documents))
	}
}

// G-4: employee dibuat dari offer recruitment eksternal yang diterima —
// referensi recruited_from_application_id tersimpan & terekspos di response
// (Employee → Application → Requisition → Position traceability).
func TestService_Create_WithRecruitedFromApplicationID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	appID := uuid.New().String()
	resp, err := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID:                 "HIRE-001",
		Name:                       "Hired From Offer",
		Email:                      strPtr("hired.offer@test.com"),
		RecruitedFromApplicationID: &appID,
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if resp.RecruitedFromApplicationID != appID {
		t.Errorf("expected recruited_from_application_id %s, got %s", appID, resp.RecruitedFromApplicationID)
	}

	// Persisted juga (baca ulang).
	persisted, err := svc.GetByID(ctx, resp.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if persisted.RecruitedFromApplicationID != appID {
		t.Errorf("expected persisted recruited_from_application_id %s, got %s", appID, persisted.RecruitedFromApplicationID)
	}
}

// =========================================================================
// Encrypt-on-write Tests (NIK, Passport, PhoneNumber, Email)
// =========================================================================

// setupEncryptionTestDB returns the shared employee test DB (setupTestDB
// already migrates and seeds SensitiveFieldSetting, all disabled by
// default) along with a Repository, so encryptIfEnabled has real rows to
// check against.
func setupEncryptionTestDB(t *testing.T) (*gorm.DB, *Repository) {
	t.Helper()
	db, dbResolver, cleanup := setupTestDB()
	t.Cleanup(cleanup)
	repo := NewRepository(dbResolver)
	return db, repo
}

func TestCreate_EncryptsNIKWhenEnabled(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "00000000000000000000000000000000000000000000000000000000000000aa")
	db, repo := setupEncryptionTestDB(t)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	ctx := context.Background()

	if err := svc.SetSensitiveFieldEnabled(ctx, "employee.nik", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}

	const originalNIK = "3201010101985678"
	nik := originalNIK
	resp, err := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "ENC-001",
		Name:       "Encrypt Test",
		NIK:        &nik,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	var stored Employee
	db.First(&stored, "id = ?", resp.ID)
	if stored.NIK == nil || *stored.NIK == originalNIK {
		t.Fatal("expected NIK to be stored encrypted, got plaintext or nil")
	}
	if !crypto.LooksEncrypted(*stored.NIK) {
		t.Errorf("stored NIK %q does not look encrypted", *stored.NIK)
	}
}

func TestCreate_StoresPlaintextWhenDisabled(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "00000000000000000000000000000000000000000000000000000000000000aa")
	db, repo := setupEncryptionTestDB(t)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	ctx := context.Background()
	// employee.nik defaults to disabled — no toggle call.

	nik := "3201010101985678"
	resp, err := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "ENC-002",
		Name:       "Plaintext Test",
		NIK:        &nik,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	var stored Employee
	db.First(&stored, "id = ?", resp.ID)
	if stored.NIK == nil || *stored.NIK != nik {
		t.Fatalf("expected NIK stored as plaintext %q, got %v", nik, stored.NIK)
	}
}

func TestUpdate_EncryptsEmailWhenEnabled(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "00000000000000000000000000000000000000000000000000000000000000aa")
	db, repo := setupEncryptionTestDB(t)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	ctx := context.Background()

	created, err := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "ENC-003",
		Name:       "Update Encrypt Test",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := svc.SetSensitiveFieldEnabled(ctx, "employee.email", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}

	const originalEmail = "encrypted.update@test.com"
	email := originalEmail
	if _, err := svc.Update(ctx, created.ID, UpdateEmployeeRequest{Email: &email}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	var stored Employee
	db.First(&stored, "id = ?", created.ID)
	if stored.Email == nil || *stored.Email == originalEmail {
		t.Fatal("expected Email to be stored encrypted, got plaintext or nil")
	}
	if !crypto.LooksEncrypted(*stored.Email) {
		t.Errorf("stored Email %q does not look encrypted", *stored.Email)
	}
}

// TestUpdate_DoesNotDoubleEncryptAlreadyEncryptedField guards against the bug
// where encryptIfEnabled was called unconditionally on emp.Email (etc.) after
// the FindEmployeeByID load, even when the incoming request did not touch
// that field. A second Update call (touching a different field) would
// re-encrypt the already-ciphertext value, corrupting it beyond recovery
// with a single Decrypt call.
func TestUpdate_DoesNotDoubleEncryptAlreadyEncryptedField(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "00000000000000000000000000000000000000000000000000000000000000aa")
	db, repo := setupEncryptionTestDB(t)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	ctx := context.Background()

	created, err := svc.Create(ctx, CreateEmployeeRequest{
		EmployeeID: "ENC-004",
		Name:       "Double Encrypt Guard Test",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := svc.SetSensitiveFieldEnabled(ctx, "employee.email", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}

	const originalEmail = "double.encrypt.guard@test.com"
	email := originalEmail
	if _, err := svc.Update(ctx, created.ID, UpdateEmployeeRequest{Email: &email}); err != nil {
		t.Fatalf("first Update() error = %v", err)
	}

	var afterFirst Employee
	db.First(&afterFirst, "id = ?", created.ID)
	if afterFirst.Email == nil || !crypto.LooksEncrypted(*afterFirst.Email) {
		t.Fatalf("expected Email to be encrypted after first Update, got %v", afterFirst.Email)
	}
	firstCiphertext := *afterFirst.Email

	// Second Update touches a different field entirely; it must not re-encrypt
	// the already-encrypted Email.
	name := "Double Encrypt Guard Test Updated"
	if _, err := svc.Update(ctx, created.ID, UpdateEmployeeRequest{Name: &name}); err != nil {
		t.Fatalf("second Update() error = %v", err)
	}

	var afterSecond Employee
	db.First(&afterSecond, "id = ?", created.ID)
	if afterSecond.Email == nil {
		t.Fatal("expected Email to remain set after second Update")
	}
	if *afterSecond.Email != firstCiphertext {
		t.Fatalf("Email ciphertext changed after unrelated Update: before=%q after=%q (indicates double-encryption)", firstCiphertext, *afterSecond.Email)
	}

	decrypted, err := crypto.DecryptString(*afterSecond.Email)
	if err != nil {
		t.Fatalf("DecryptString() error = %v (value may be double-encrypted)", err)
	}
	if decrypted != originalEmail {
		t.Fatalf("decrypted Email = %q, want %q", decrypted, originalEmail)
	}
}

// =========================================================================
// Task 10: Encrypt-on-write for Family, Bank Account, Emergency Contact
// =========================================================================

func TestCreateBank_EncryptsAccountNumberWhenEnabled(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "00000000000000000000000000000000000000000000000000000000000000aa")
	db, repo := setupEncryptionTestDB(t)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	ctx := context.Background()

	if err := svc.SetSensitiveFieldEnabled(ctx, "employee_bank_account.account_number", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}

	emp, err := svc.Create(ctx, CreateEmployeeRequest{EmployeeID: "ENC-BANK-001", Name: "Bank Encrypt Test"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	resp, err := svc.CreateBank(ctx, emp.ID, CreateBankRequest{AccountNumber: "1234567890", AccountName: "Budi"})
	if err != nil {
		t.Fatalf("CreateBank() error = %v", err)
	}

	var stored EmployeeBankAccount
	db.First(&stored, "id = ?", resp.ID)
	if stored.AccountNumber == "1234567890" {
		t.Fatal("expected account_number to be stored encrypted")
	}
	if !crypto.LooksEncrypted(stored.AccountNumber) {
		t.Errorf("stored account_number %q does not look encrypted", stored.AccountNumber)
	}
}

func TestCreateFamily_EncryptsNIKWhenEnabled(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "00000000000000000000000000000000000000000000000000000000000000aa")
	db, repo := setupEncryptionTestDB(t)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	ctx := context.Background()

	if err := svc.SetSensitiveFieldEnabled(ctx, "employee_family.nik", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}

	emp, err := svc.Create(ctx, CreateEmployeeRequest{EmployeeID: "ENC-FAM-001", Name: "Family Encrypt Test"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	const originalNIK = "3201010101985678"
	nik := originalNIK
	resp, err := svc.CreateFamily(ctx, emp.ID, CreateFamilyRequest{NIK: &nik, Name: "Anak Pertama"})
	if err != nil {
		t.Fatalf("CreateFamily() error = %v", err)
	}

	var stored EmployeeFamily
	db.First(&stored, "id = ?", resp.ID)
	if stored.NIK == nil || *stored.NIK == originalNIK {
		t.Fatal("expected family NIK to be stored encrypted")
	}
	if !crypto.LooksEncrypted(*stored.NIK) {
		t.Errorf("stored family NIK %q does not look encrypted", *stored.NIK)
	}
}

func TestCreateEmergencyContact_EncryptsPhoneWhenEnabled(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "00000000000000000000000000000000000000000000000000000000000000aa")
	db, repo := setupEncryptionTestDB(t)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	ctx := context.Background()

	if err := svc.SetSensitiveFieldEnabled(ctx, "emergency_contact.phone_number", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}

	emp, err := svc.Create(ctx, CreateEmployeeRequest{EmployeeID: "ENC-EC-001", Name: "Emergency Contact Encrypt Test"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	resp, err := svc.CreateEmergencyContact(ctx, emp.ID, CreateEmergencyContactRequest{Name: "Ibu", PhoneNumber: "081234567890"})
	if err != nil {
		t.Fatalf("CreateEmergencyContact() error = %v", err)
	}

	var stored EmergencyContact
	db.First(&stored, "id = ?", resp.ID)
	if stored.PhoneNumber == "081234567890" {
		t.Fatal("expected emergency contact phone to be stored encrypted")
	}
	if !crypto.LooksEncrypted(stored.PhoneNumber) {
		t.Errorf("stored emergency contact phone %q does not look encrypted", stored.PhoneNumber)
	}
}

// TestUpdateBank_DoesNotDoubleEncryptAlreadyEncryptedField guards against the
// same double-encryption bug fixed in Task 9, applied to bank account fields.
func TestUpdateBank_DoesNotDoubleEncryptAlreadyEncryptedField(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "00000000000000000000000000000000000000000000000000000000000000aa")
	db, repo := setupEncryptionTestDB(t)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	ctx := context.Background()

	if err := svc.SetSensitiveFieldEnabled(ctx, "employee_bank_account.account_number", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}

	emp, err := svc.Create(ctx, CreateEmployeeRequest{EmployeeID: "ENC-BANK-002", Name: "Bank Double Encrypt Guard"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	const originalAccountNumber = "9988776655"
	created, err := svc.CreateBank(ctx, emp.ID, CreateBankRequest{AccountNumber: originalAccountNumber, AccountName: "Budi"})
	if err != nil {
		t.Fatalf("CreateBank() error = %v", err)
	}

	var afterCreate EmployeeBankAccount
	db.First(&afterCreate, "id = ?", created.ID)
	if !crypto.LooksEncrypted(afterCreate.AccountNumber) {
		t.Fatalf("expected account_number to be encrypted after Create, got %v", afterCreate.AccountNumber)
	}
	firstCiphertext := afterCreate.AccountNumber

	// Update touches an unrelated field (AccountName); account_number must not
	// be re-encrypted.
	newName := "Budi Updated"
	if _, err := svc.UpdateBank(ctx, emp.ID, created.ID, UpdateBankRequest{AccountName: &newName}); err != nil {
		t.Fatalf("UpdateBank() error = %v", err)
	}

	var afterUpdate EmployeeBankAccount
	db.First(&afterUpdate, "id = ?", created.ID)
	if afterUpdate.AccountNumber != firstCiphertext {
		t.Fatalf("account_number ciphertext changed after unrelated Update: before=%q after=%q (indicates double-encryption)", firstCiphertext, afterUpdate.AccountNumber)
	}

	decrypted, err := crypto.DecryptString(afterUpdate.AccountNumber)
	if err != nil {
		t.Fatalf("DecryptString() error = %v (value may be double-encrypted)", err)
	}
	if decrypted != originalAccountNumber {
		t.Fatalf("decrypted account_number = %q, want %q", decrypted, originalAccountNumber)
	}
}
