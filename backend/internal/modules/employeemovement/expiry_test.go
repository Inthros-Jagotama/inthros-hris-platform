package employeemovement

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// seedContractExpiryReferenceTables membuat tabel minimal rbac + employee
// accounts yang dibutuhkan ProcessContractExpiration untuk resolve user HR
// (permissions / model_has_roles / role_has_permissions) dan employee →
// user. Skema dibuat sekecil mungkin, konsisten dengan migration 011_settings
// dan 023_user_accounts.
func seedContractExpiryReferenceTables(t *testing.T, svc *Service, employeeID, hrUserID uuid.UUID) {
	t.Helper()
	db, err := svc.repo.getDB(ctx())
	if err != nil {
		t.Fatalf("failed to get test db: %v", err)
	}
	for _, stmt := range []string{
		"CREATE TABLE IF NOT EXISTS permissions (id CHAR(36) PRIMARY KEY, name VARCHAR(255) NOT NULL, guard_name VARCHAR(255) NOT NULL DEFAULT 'web')",
		"CREATE TABLE IF NOT EXISTS roles (id CHAR(36) PRIMARY KEY, name VARCHAR(255) NOT NULL, guard_name VARCHAR(255) NOT NULL DEFAULT 'web')",
		"CREATE TABLE IF NOT EXISTS role_has_permissions (permission_id CHAR(36) NOT NULL, role_id CHAR(36) NOT NULL, PRIMARY KEY (permission_id, role_id))",
		"CREATE TABLE IF NOT EXISTS model_has_roles (role_id CHAR(36) NOT NULL, model_type VARCHAR(255) NOT NULL, model_id CHAR(36) NOT NULL, PRIMARY KEY (role_id, model_id, model_type))",
		"CREATE TABLE IF NOT EXISTS model_has_permissions (permission_id CHAR(36) NOT NULL, model_type VARCHAR(255) NOT NULL, model_id CHAR(36) NOT NULL, PRIMARY KEY (permission_id, model_id, model_type))",
	} {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("failed to create rbac reference table: %v", err)
		}
	}
	// AutoMigrate test stand-in employee_accounts (pola approval_integration_test).
	if err := db.AutoMigrate(&EmployeeAccount{}); err != nil {
		t.Fatalf("failed to migrate employee_accounts: %v", err)
	}

	permissionID := uuid.New()
	roleID := uuid.New()
	if err := db.Exec("INSERT INTO permissions (id, name, guard_name) VALUES (?, ?, 'web')", permissionID.String(), contractExpiryHRPermission).Error; err != nil {
		t.Fatalf("failed to seed permission: %v", err)
	}
	if err := db.Exec("INSERT INTO roles (id, name, guard_name) VALUES (?, ?, 'web')", roleID.String(), "HR").Error; err != nil {
		t.Fatalf("failed to seed role: %v", err)
	}
	if err := db.Exec("INSERT INTO role_has_permissions (permission_id, role_id) VALUES (?, ?)", permissionID.String(), roleID.String()).Error; err != nil {
		t.Fatalf("failed to seed role_has_permissions: %v", err)
	}
	// HR user mendapat role HR (model_type=user).
	if err := db.Exec("INSERT INTO model_has_roles (role_id, model_type, model_id) VALUES (?, 'user', ?)", roleID.String(), hrUserID.String()).Error; err != nil {
		t.Fatalf("failed to seed model_has_roles: %v", err)
	}
	// Employee account → user (untuk notifikasi ke employee pemilik kontrak).
	if err := db.Create(&EmployeeAccount{
		EmployeeID: employeeID,
		UserID:     hrUserID, // pakai user yang sama agar HR == employee di test ini
		Email:      "expiry@test.local",
	}).Error; err != nil {
		t.Fatalf("failed to seed employee account: %v", err)
	}
}

// TestService_ProcessContractExpiration_MarksExpired verifies plan §12.13:
// kontrak active yang sudah lewat end_date dipindah ke status expired dan
// employee + HR menerima notifikasi CONTRACT_EXPIRED.
func TestService_ProcessContractExpiration_MarksExpired(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	hrUserID := uuid.New()
	seedContractExpiryReferenceTables(t, svc, employeeID, hrUserID)

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	// Kontrak yang sudah lewat end_date (kemarin).
	past := createTestContract(repo, employeeID)
	past.EndDate = strPtr(addDays(time.Now().Format("2006-01-02"), -1))
	if err := repo.UpdateContract(ctx(), past); err != nil {
		t.Fatalf("failed to seed expired contract: %v", err)
	}
	// Kontrak yang masih aktif (belum lewat) — tidak boleh berubah.
	future := createTestContract(repo, employeeID)
	future.EndDate = strPtr(addDays(time.Now().Format("2006-01-02"), 60))
	if err := repo.UpdateContract(ctx(), future); err != nil {
		t.Fatalf("failed to seed future contract: %v", err)
	}

	if err := svc.ProcessContractExpiration(ctx()); err != nil {
		t.Fatalf("ProcessContractExpiration failed: %v", err)
	}

	// Kontrak masa lalu → expired.
	updatedPast, err := repo.FindContractByID(ctx(), past.ID)
	if err != nil {
		t.Fatalf("failed to reload past contract: %v", err)
	}
	if updatedPast.Status != ContractStatusExpired {
		t.Errorf("expected past contract status expired, got '%s'", updatedPast.Status)
	}
	// Kontrak masa depan → tetap active.
	updatedFuture, err := repo.FindContractByID(ctx(), future.ID)
	if err != nil {
		t.Fatalf("failed to reload future contract: %v", err)
	}
	if updatedFuture.Status != ContractStatusActive {
		t.Errorf("expected future contract status active, got '%s'", updatedFuture.Status)
	}

	// Notifikasi CONTRACT_EXPIRED dikirim ke employee (== HR di test ini).
	found := false
	for _, call := range notifier.calls {
		if call.notifType == "CONTRACT_EXPIRED" && call.referenceID == past.ID {
			found = true
		}
	}
	if !found {
		t.Errorf("expected CONTRACT_EXPIRED notification for expired contract, got %+v", notifier.calls)
	}
}

// TestService_ProcessContractExpiration_Reminder verifies reminder H-30/H-14/
// H-7/H-1 mengirim notifikasi CONTRACT_EXPIRING untuk kontrak yang berakhir
// tepat N hari lagi.
func TestService_ProcessContractExpiration_Reminder(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	hrUserID := uuid.New()
	seedContractExpiryReferenceTables(t, svc, employeeID, hrUserID)

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	today := time.Now().Format("2006-01-02")
	// Kontrak berakhir H-30 → harus dapat reminder 30.
	c30 := createTestContract(repo, employeeID)
	c30.EndDate = strPtr(addDays(today, 30))
	if err := repo.UpdateContract(ctx(), c30); err != nil {
		t.Fatalf("failed to seed H-30 contract: %v", err)
	}
	// Kontrak berakhir H-7 → harus dapat reminder 7.
	c7 := createTestContract(repo, employeeID)
	c7.EndDate = strPtr(addDays(today, 7))
	if err := repo.UpdateContract(ctx(), c7); err != nil {
		t.Fatalf("failed to seed H-7 contract: %v", err)
	}
	// Kontrak berakhir H-15 — bukan jadwal reminder → tidak boleh dapat notif.
	c15 := createTestContract(repo, employeeID)
	c15.EndDate = strPtr(addDays(today, 15))
	if err := repo.UpdateContract(ctx(), c15); err != nil {
		t.Fatalf("failed to seed H-15 contract: %v", err)
	}

	if err := svc.ProcessContractExpiration(ctx()); err != nil {
		t.Fatalf("ProcessContractExpiration failed: %v", err)
	}

	got30, got7, got15 := false, false, false
	for _, call := range notifier.calls {
		switch {
		case call.notifType == "CONTRACT_EXPIRING" && call.referenceID == c30.ID:
			got30 = true
		case call.notifType == "CONTRACT_EXPIRING" && call.referenceID == c7.ID:
			got7 = true
		case call.notifType == "CONTRACT_EXPIRING" && call.referenceID == c15.ID:
			got15 = true
		}
	}
	if !got30 {
		t.Error("expected CONTRACT_EXPIRING for H-30 contract")
	}
	if !got7 {
		t.Error("expected CONTRACT_EXPIRING for H-7 contract")
	}
	if got15 {
		t.Error("did not expect CONTRACT_EXPIRING for H-15 contract (not a reminder day)")
	}

	// Semua kontrak tetap active (belum lewat end_date).
	for _, c := range []*EmployeeContract{c30, c7, c15} {
		loaded, err := repo.FindContractByID(ctx(), c.ID)
		if err != nil {
			t.Fatalf("failed to reload contract: %v", err)
		}
		if loaded.Status != ContractStatusActive {
			t.Errorf("expected contract %s active, got '%s'", c.ID.String(), loaded.Status)
		}
	}
}

// TestService_ProcessContractExpiration_NoNotifier verifies ProcessContractExpiration
// tetap berjalan (mark expired) ketika notifier tidak dikonfigurasi.
func TestService_ProcessContractExpiration_NoNotifier(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	// Tanpa SetNotifier — notifikasi skip, proses mark expired tetap jalan.

	past := createTestContract(repo, employeeID)
	past.EndDate = strPtr(addDays(time.Now().Format("2006-01-02"), -1))
	if err := repo.UpdateContract(ctx(), past); err != nil {
		t.Fatalf("failed to seed expired contract: %v", err)
	}

	if err := svc.ProcessContractExpiration(ctx()); err != nil {
		t.Fatalf("ProcessContractExpiration failed without notifier: %v", err)
	}
	loaded, err := repo.FindContractByID(ctx(), past.ID)
	if err != nil {
		t.Fatalf("failed to reload contract: %v", err)
	}
	if loaded.Status != ContractStatusExpired {
		t.Errorf("expected contract status expired, got '%s'", loaded.Status)
	}
}
