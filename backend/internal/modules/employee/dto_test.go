package employee

import (
	"testing"

	"github.com/inthros/hris-platform/internal/pkg/crypto"
)

func TestToFamilyResponse_DecryptsEncryptedNIK(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "00000000000000000000000000000000000000000000000000000000000000aa")
	plain := "3201010101985678"
	encrypted, err := crypto.EncryptString(plain)
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}
	fam := &EmployeeFamily{NIK: &encrypted}

	resp := toFamilyResponse(fam)

	if resp.NIK != plain {
		t.Errorf("toFamilyResponse().NIK = %q, want %q", resp.NIK, plain)
	}
}

func TestToFamilyResponse_LeavesPlaintextNIKAsIs(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "00000000000000000000000000000000000000000000000000000000000000aa")
	plain := "3201010101985678"
	fam := &EmployeeFamily{NIK: &plain}

	resp := toFamilyResponse(fam)

	if resp.NIK != plain {
		t.Errorf("toFamilyResponse().NIK = %q, want %q (unchanged plaintext)", resp.NIK, plain)
	}
}

func TestToBankResponse_DecryptsEncryptedAccountNumber(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "00000000000000000000000000000000000000000000000000000000000000aa")
	plain := "1234567890"
	encrypted, err := crypto.EncryptString(plain)
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}
	bank := &EmployeeBankAccount{AccountNumber: encrypted, AccountName: "Budi"}

	resp := toBankResponse(bank)

	if resp.AccountNumber != plain {
		t.Errorf("toBankResponse().AccountNumber = %q, want %q", resp.AccountNumber, plain)
	}
}

func TestToBankResponse_DecryptsEncryptedAccountName(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "00000000000000000000000000000000000000000000000000000000000000aa")
	plain := "Budi Santoso"
	encrypted, err := crypto.EncryptString(plain)
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}
	bank := &EmployeeBankAccount{AccountNumber: "1234567890", AccountName: encrypted}

	resp := toBankResponse(bank)

	if resp.AccountName != plain {
		t.Errorf("toBankResponse().AccountName = %q, want %q", resp.AccountName, plain)
	}
}

func TestToEmergencyContactResponse_DecryptsEncryptedPhoneNumber(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "00000000000000000000000000000000000000000000000000000000000000aa")
	plain := "081234567890"
	encrypted, err := crypto.EncryptString(plain)
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}
	contact := &EmergencyContact{Name: "Ani", PhoneNumber: encrypted}

	resp := toEmergencyContactResponse(contact)

	if resp.PhoneNumber != plain {
		t.Errorf("toEmergencyContactResponse().PhoneNumber = %q, want %q", resp.PhoneNumber, plain)
	}
}

func TestEmployeeToResponse_DecryptsEncryptedNIK(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "00000000000000000000000000000000000000000000000000000000000000aa")
	plain := "3201010101985678"
	encrypted, err := crypto.EncryptString(plain)
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}
	emp := &Employee{NIK: &encrypted}

	resp := emp.ToResponse()

	if resp.NIK != plain {
		t.Errorf("ToResponse().NIK = %q, want %q", resp.NIK, plain)
	}
}

func TestEmployeeToResponse_DecryptsEncryptedPassportPhoneEmail(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "00000000000000000000000000000000000000000000000000000000000000aa")
	plainPassport := "A1234567"
	plainPhone := "081234567890"
	plainEmail := "budi@example.com"

	encPassport, err := crypto.EncryptString(plainPassport)
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}
	encPhone, err := crypto.EncryptString(plainPhone)
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}
	encEmail, err := crypto.EncryptString(plainEmail)
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}

	emp := &Employee{Passport: &encPassport, PhoneNumber: &encPhone, Email: &encEmail}

	resp := emp.ToResponse()

	if resp.Passport != plainPassport {
		t.Errorf("ToResponse().Passport = %q, want %q", resp.Passport, plainPassport)
	}
	if resp.PhoneNumber != plainPhone {
		t.Errorf("ToResponse().PhoneNumber = %q, want %q", resp.PhoneNumber, plainPhone)
	}
	if resp.Email != plainEmail {
		t.Errorf("ToResponse().Email = %q, want %q", resp.Email, plainEmail)
	}
}
