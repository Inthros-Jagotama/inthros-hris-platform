package employeemovement

import (
	"testing"

	"github.com/google/uuid"

	"github.com/inthros/hris-platform/internal/pkg/crypto"
)

// TestGetEmployeeProfile_DecryptsSensitiveFields adalah regression test untuk
// temuan review: jalur generate document membaca employees.nik/passport/
// phone_number/email lewat raw SQL tanpa dekripsi, sehingga begitu toggle
// encrypt-at-rest per-field aktif, dokumen kontrak/SK mencetak blob ciphertext
// alih-alih nilai asli.
//
// Test ini menyimpan nilai TERENKRIPSI di tabel employees (persis seperti yang
// dilakukan encrypt-on-write modul employee), lalu memastikan value yang masuk
// ke variable {{employee.*}} sudah berupa plaintext.
func TestGetEmployeeProfile_DecryptsSensitiveFields(t *testing.T) {
	t.Setenv(crypto.EnvEncryptionKey, "00000000000000000000000000000000000000000000000000000000000000aa")

	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	employeeID := uuid.New()
	seedCareerReferenceTables(t, repo, employeeID)

	const (
		plainNIK      = "3201010101985678"
		plainPassport = "A1234567890"
		plainPhone    = "081234567890"
		plainEmail    = "real.person@example.com"
	)
	encrypted := map[string]string{}
	for col, plain := range map[string]string{
		"nik":          plainNIK,
		"passport":     plainPassport,
		"phone_number": plainPhone,
		"email":        plainEmail,
	} {
		ct, err := crypto.EncryptString(plain)
		if err != nil {
			t.Fatalf("EncryptString(%s) error = %v", col, err)
		}
		encrypted[col] = ct
	}

	db, err := repo.getDB(ctx())
	if err != nil {
		t.Fatalf("getDB() error = %v", err)
	}
	if err := db.Exec(
		"UPDATE employees SET nik = ?, passport = ?, phone_number = ?, email = ? WHERE id = ?",
		encrypted["nik"], encrypted["passport"], encrypted["phone_number"], encrypted["email"], employeeID.String(),
	).Error; err != nil {
		t.Fatalf("seed encrypted employee columns: %v", err)
	}

	profile, err := repo.GetEmployeeProfile(ctx(), employeeID)
	if err != nil {
		t.Fatalf("GetEmployeeProfile() error = %v", err)
	}

	vars := employeeValues(profile)
	for _, tc := range []struct {
		key  string
		want string
	}{
		{"employee.nik", plainNIK},
		{"employee.passport", plainPassport},
		{"employee.phone_number", plainPhone},
		{"employee.email", plainEmail},
	} {
		if got := vars[tc.key]; got != tc.want {
			t.Errorf("{{%s}} = %q, want decrypted %q (ciphertext leaked into generated document)", tc.key, got, tc.want)
		}
	}
}

// Data lama yang masih plaintext (toggle enkripsi belum pernah aktif) harus
// lewat apa adanya — dekripsi bersifat opportunistic.
func TestGetEmployeeProfile_LeavesPlaintextUnchanged(t *testing.T) {
	t.Setenv(crypto.EnvEncryptionKey, "00000000000000000000000000000000000000000000000000000000000000aa")

	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	employeeID := uuid.New()
	seedCareerReferenceTables(t, repo, employeeID)

	db, err := repo.getDB(ctx())
	if err != nil {
		t.Fatalf("getDB() error = %v", err)
	}
	const plainNIK = "3201010101985678"
	if err := db.Exec("UPDATE employees SET nik = ? WHERE id = ?", plainNIK, employeeID.String()).Error; err != nil {
		t.Fatalf("seed plaintext nik: %v", err)
	}

	profile, err := repo.GetEmployeeProfile(ctx(), employeeID)
	if err != nil {
		t.Fatalf("GetEmployeeProfile() error = %v", err)
	}
	if profile.NIK != plainNIK {
		t.Errorf("NIK = %q, want unchanged plaintext %q", profile.NIK, plainNIK)
	}
}
