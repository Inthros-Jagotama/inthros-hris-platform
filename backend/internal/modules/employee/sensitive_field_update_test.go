package employee

import (
	"context"
	"testing"

	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/crypto"
	"github.com/inthros/hris-platform/internal/pkg/mask"
)

// =============================================================================
// Regression: masked value echoed back on Update must NOT overwrite real PII.
//
// Skenario nyata: caller punya "employee.update" tapi TIDAK punya
// "employee.view_nik". GET mengembalikan NIK ter-mask; form edit mengirim
// balik seluruh objek (termasuk mask itu) saat user hanya mengubah nama.
// Sebelum perbaikan, NIK asli tertimpa string mask secara permanen.
// =============================================================================

const testKey = "00000000000000000000000000000000000000000000000000000000000000aa"

// ctxWithPerms membuat context dengan permissions claim seperti yang di-set
// middleware AuthJWT/TenantRequired.
func ctxWithPerms(perms ...string) context.Context {
	return context.WithValue(context.Background(), "permissions", perms) //nolint:staticcheck // pola key string sama dengan middleware
}

func newMaskTestService(t *testing.T) (*Service, *Repository, func(dest interface{}, id string)) {
	t.Helper()
	db, repo := setupEncryptionTestDB(t)
	logger, _ := zap.NewDevelopment()
	return NewService(repo, logger), repo, func(dest interface{}, id string) {
		db.First(dest, "id = ?", id)
	}
}

func TestUpdate_MaskedNIKEchoDoesNotOverwriteStoredValue(t *testing.T) {
	t.Setenv(crypto.EnvEncryptionKey, testKey)
	svc, _, load := newMaskTestService(t)

	// Admin (punya semua view_*) membuat employee dengan NIK asli, enkripsi on.
	admin := ctxWithPerms("*")
	if err := svc.SetSensitiveFieldEnabled(admin, "employee.nik", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}
	const originalNIK = "3201010101985678"
	nik := originalNIK
	created, err := svc.Create(admin, CreateEmployeeRequest{
		EmployeeID: "MASK-001", Name: "Mask Guard", NIK: &nik,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Caller tanpa view_nik melakukan GET → menerima NIK ter-mask.
	limited := ctxWithPerms("employee.view", "employee.update")
	got, err := svc.GetByID(limited, created.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	wantMask := mask.PartialMask(originalNIK)
	if got.NIK != wantMask {
		t.Fatalf("GetByID NIK = %q, want masked %q", got.NIK, wantMask)
	}

	// ...lalu PUT balik seluruh objek, hanya nama yang benar-benar berubah.
	echoed := got.NIK
	newName := "Mask Guard Renamed"
	if _, err := svc.Update(limited, created.ID, UpdateEmployeeRequest{
		Name: &newName, NIK: &echoed,
	}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	var stored Employee
	load(&stored, created.ID)
	if stored.NIK == nil {
		t.Fatal("stored NIK is nil after Update")
	}
	if got := crypto.DecryptIfLooksEncrypted(*stored.NIK); got != originalNIK {
		t.Errorf("stored NIK decrypts to %q, want original %q (masked echo overwrote real PII)", got, originalNIK)
	}
	if stored.Name != newName {
		t.Errorf("Name = %q, want %q (non-sensitive edit must still apply)", stored.Name, newName)
	}
}

func TestUpdate_GenuineNewNIKIsWrittenAndEncrypted(t *testing.T) {
	t.Setenv(crypto.EnvEncryptionKey, testKey)
	svc, _, load := newMaskTestService(t)
	admin := ctxWithPerms("*")

	if err := svc.SetSensitiveFieldEnabled(admin, "employee.nik", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}
	nik := "3201010101985678"
	created, err := svc.Create(admin, CreateEmployeeRequest{
		EmployeeID: "MASK-002", Name: "Genuine Edit", NIK: &nik,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	const newNIK = "3201010101991234"
	updated := newNIK
	if _, err := svc.Update(admin, created.ID, UpdateEmployeeRequest{NIK: &updated}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	var stored Employee
	load(&stored, created.ID)
	if stored.NIK == nil || !crypto.LooksEncrypted(*stored.NIK) {
		t.Fatalf("expected new NIK stored encrypted, got %v", stored.NIK)
	}
	if got := crypto.DecryptIfLooksEncrypted(*stored.NIK); got != newNIK {
		t.Errorf("stored NIK decrypts to %q, want %q", got, newNIK)
	}
}

// Passport/PhoneNumber/Email pada model employee memakai jalur yang sama.
func TestUpdate_MaskedEchoGuardCoversAllEmployeeSensitiveFields(t *testing.T) {
	t.Setenv(crypto.EnvEncryptionKey, testKey)
	svc, _, load := newMaskTestService(t)
	admin := ctxWithPerms("*")

	passport, phone, email := "A1234567890", "081234567890", "real.person@example.com"
	created, err := svc.Create(admin, CreateEmployeeRequest{
		EmployeeID: "MASK-003", Name: "All Fields",
		Passport: &passport, PhoneNumber: &phone, Email: &email,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	limited := ctxWithPerms("employee.view", "employee.update")
	got, err := svc.GetByID(limited, created.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if _, err := svc.Update(limited, created.ID, UpdateEmployeeRequest{
		Passport: &got.Passport, PhoneNumber: &got.PhoneNumber, Email: &got.Email,
	}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	var stored Employee
	load(&stored, created.ID)
	fields := []struct {
		name string
		want string
		have *string
	}{
		{"Passport", passport, stored.Passport},
		{"PhoneNumber", phone, stored.PhoneNumber},
		{"Email", email, stored.Email},
	}
	for _, tc := range fields {
		if tc.have == nil {
			t.Errorf("%s = nil, want %q", tc.name, tc.want)
			continue
		}
		if *tc.have != tc.want {
			t.Errorf("%s = %q, want %q (masked echo overwrote real PII)", tc.name, *tc.have, tc.want)
		}
	}
}

func TestUpdateFamily_MaskedNIKEchoDoesNotOverwriteStoredValue(t *testing.T) {
	t.Setenv(crypto.EnvEncryptionKey, testKey)
	svc, repo, load := newMaskTestService(t)
	admin := ctxWithPerms("*")
	emp := createTestEmployee(admin, repo)

	if err := svc.SetSensitiveFieldEnabled(admin, "employee_family.nik", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}
	const originalNIK = "3271010101990001"
	nik := originalNIK
	fam, err := svc.CreateFamily(admin, emp.ID.String(), CreateFamilyRequest{Name: "Spouse", NIK: &nik})
	if err != nil {
		t.Fatalf("CreateFamily() error = %v", err)
	}

	limited := ctxWithPerms("employee.view", "employee.update")
	masked := mask.PartialMask(originalNIK)
	newName := "Spouse Renamed"
	if _, err := svc.UpdateFamily(limited, emp.ID.String(), fam.ID, UpdateFamilyRequest{
		Name: &newName, NIK: &masked,
	}); err != nil {
		t.Fatalf("UpdateFamily() error = %v", err)
	}

	var stored EmployeeFamily
	load(&stored, fam.ID)
	if stored.NIK == nil {
		t.Fatal("family NIK is nil after UpdateFamily")
	}
	if got := crypto.DecryptIfLooksEncrypted(*stored.NIK); got != originalNIK {
		t.Errorf("family NIK = %q, want %q", got, originalNIK)
	}
}

func TestUpdateBank_MaskedEchoDoesNotOverwriteStoredValues(t *testing.T) {
	t.Setenv(crypto.EnvEncryptionKey, testKey)
	svc, repo, load := newMaskTestService(t)
	admin := ctxWithPerms("*")
	emp := createTestEmployee(admin, repo)

	if err := svc.SetSensitiveFieldEnabled(admin, "employee_bank_account.account_number", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}
	const accNo, accName = "1234567890123", "REAL ACCOUNT HOLDER"
	bank, err := svc.CreateBank(admin, emp.ID.String(), CreateBankRequest{
		AccountNumber: accNo, AccountName: accName,
	})
	if err != nil {
		t.Fatalf("CreateBank() error = %v", err)
	}

	limited := ctxWithPerms("employee.view", "employee.update")
	maskedNo := mask.PartialMask(accNo)
	maskedName := mask.PartialMask(accName)
	if _, err := svc.UpdateBank(limited, emp.ID.String(), bank.ID, UpdateBankRequest{
		AccountNumber: &maskedNo, AccountName: &maskedName,
	}); err != nil {
		t.Fatalf("UpdateBank() error = %v", err)
	}

	var stored EmployeeBankAccount
	load(&stored, bank.ID)
	if got := crypto.DecryptIfLooksEncrypted(stored.AccountNumber); got != accNo {
		t.Errorf("account_number = %q, want %q", got, accNo)
	}
	if stored.AccountName != accName {
		t.Errorf("account_name = %q, want %q", stored.AccountName, accName)
	}

	// Edit sungguhan tetap tersimpan.
	newName := "NEW ACCOUNT HOLDER"
	if _, err := svc.UpdateBank(admin, emp.ID.String(), bank.ID, UpdateBankRequest{AccountName: &newName}); err != nil {
		t.Fatalf("UpdateBank() genuine edit error = %v", err)
	}
	load(&stored, bank.ID)
	if stored.AccountName != newName {
		t.Errorf("account_name = %q, want genuine edit %q", stored.AccountName, newName)
	}
}

func TestUpdateEmergencyContact_MaskedPhoneEchoDoesNotOverwriteStoredValue(t *testing.T) {
	t.Setenv(crypto.EnvEncryptionKey, testKey)
	svc, repo, load := newMaskTestService(t)
	admin := ctxWithPerms("*")
	emp := createTestEmployee(admin, repo)

	if err := svc.SetSensitiveFieldEnabled(admin, "emergency_contact.phone_number", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}
	const originalPhone = "081298765432"
	contact, err := svc.CreateEmergencyContact(admin, emp.ID.String(), CreateEmergencyContactRequest{
		Name: "Contact", PhoneNumber: originalPhone,
	})
	if err != nil {
		t.Fatalf("CreateEmergencyContact() error = %v", err)
	}

	limited := ctxWithPerms("employee.view", "employee.update")
	masked := mask.PartialMask(originalPhone)
	if _, err := svc.UpdateEmergencyContact(limited, emp.ID.String(), contact.ID, UpdateEmergencyContactRequest{
		PhoneNumber: &masked,
	}); err != nil {
		t.Fatalf("UpdateEmergencyContact() error = %v", err)
	}

	var stored EmergencyContact
	load(&stored, contact.ID)
	if got := crypto.DecryptIfLooksEncrypted(stored.PhoneNumber); got != originalPhone {
		t.Errorf("phone_number = %q, want %q", got, originalPhone)
	}
}

// isMaskedEcho tidak boleh menganggap nilai plaintext yang kebetulan mirip
// sebagai mask, dan harus tahan terhadap nilai kosong.
func TestIsMaskedEcho(t *testing.T) {
	t.Setenv(crypto.EnvEncryptionKey, testKey)
	cases := []struct {
		name, stored, incoming string
		want                   bool
	}{
		{"masked echo of plaintext", "3201010101985678", "************5678", true},
		{"genuine new value", "3201010101985678", "3201010101991234", false},
		{"identical value", "3201010101985678", "3201010101985678", false},
		{"empty stored", "", "************5678", false},
		{"empty incoming", "3201010101985678", "", false},
		{"user typed literal stars but different length", "3201010101985678", "****5678", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isMaskedEcho(tc.stored, tc.incoming); got != tc.want {
				t.Errorf("isMaskedEcho(%q, %q) = %v, want %v", tc.stored, tc.incoming, got, tc.want)
			}
		})
	}
}

func TestIsMaskedEcho_StoredCiphertext(t *testing.T) {
	t.Setenv(crypto.EnvEncryptionKey, testKey)
	const plain = "3201010101985678"
	ct, err := crypto.EncryptString(plain)
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}
	if !isMaskedEcho(ct, mask.PartialMask(plain)) {
		t.Error("expected masked echo to be detected against encrypted stored value")
	}
	if isMaskedEcho(ct, "3201010101991234") {
		t.Error("genuine new value must not be treated as masked echo")
	}
}
