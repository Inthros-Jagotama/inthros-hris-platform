package company

import (
	"fmt"
	"testing"

	"github.com/inthros/hris-platform/internal/pkg/database"
)

// =========================================================================
// subdomainFromHost Tests
// =========================================================================

func TestSubdomainFromHost(t *testing.T) {
	cases := []struct {
		host     string
		expected string
	}{
		{"pt-inthros-jago-utama.localhost", "pt-inthros-jago-utama"},
		{"pt-inthros-jago-utama.localhost:5174", "pt-inthros-jago-utama"},
		{"localhost", ""},
		{"localhost:5174", ""},
		{"127.0.0.1", ""},
		{"127.0.0.1:5174", ""},
		{"hris.pt-inthros.com", "hris"},
		{"www.pt-inthros.com", ""}, // www diabaikan
		{"pt-inthros.com", "pt-inthros"},
	}
	for _, tc := range cases {
		got := subdomainFromHost(tc.host)
		if got != tc.expected {
			t.Errorf("subdomainFromHost(%q) = %q, want %q", tc.host, got, tc.expected)
		}
	}
}

// =========================================================================
// ResolveByHost Service Tests
// =========================================================================

func TestService_ResolveByHost_ByDomain(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	company := createTestCompany(svc.repo.db, "Resolve Domain")
	domain := "hris.pt-inthros.com"
	company.Domain = &domain
	if err := svc.repo.Update(company); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	resp, err := svc.ResolveByHost(domain)
	if err != nil {
		t.Fatalf("ResolveByHost by domain failed: %v", err)
	}
	if resp.ID != company.ID.String() {
		t.Errorf("expected company ID %s, got %s", company.ID.String(), resp.ID)
	}
	if resp.Domain == nil || *resp.Domain != domain {
		t.Errorf("expected domain %q in response", domain)
	}
}

func TestService_ResolveByHost_BySubdomain(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	company := createTestCompany(svc.repo.db, "Resolve Subdomain")
	sub := "pt-inthros-jago-utama"
	company.Subdomain = &sub
	if err := svc.repo.Update(company); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Host penuh (subdomain + base) → cocok via subdomain.
	resp, err := svc.ResolveByHost(sub + ".localhost")
	if err != nil {
		t.Fatalf("ResolveByHost by subdomain failed: %v", err)
	}
	if resp.ID != company.ID.String() {
		t.Errorf("expected company ID %s, got %s", company.ID.String(), resp.ID)
	}
}

func TestService_ResolveByHost_NonActive_Rejected(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	company := createTestCompany(svc.repo.db, "Resolve Inactive")
	company.Status = CompanyStatusSuspended
	sub := "inactive-sub"
	company.Subdomain = &sub
	if err := svc.repo.Update(company); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if _, err := svc.ResolveByHost(sub + ".localhost"); err == nil {
		t.Fatal("expected error for non-active company, got nil")
	}
}

func TestService_ResolveByHost_NotFound(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	if _, err := svc.ResolveByHost("unknown.localhost"); err == nil {
		t.Fatal("expected error for unknown host, got nil")
	}
}

func TestService_ResolveByHost_EmptyHost(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	if _, err := svc.ResolveByHost(""); err == nil {
		t.Fatal("expected error for empty host, got nil")
	}
}

// =========================================================================
// emptyToNil Tests
// =========================================================================

func TestEmptyToNil(t *testing.T) {
	// nil → nil
	if got := emptyToNil(nil); got != nil {
		t.Errorf("emptyToNil(nil) = %q, want nil", *got)
	}

	// "" → nil (hindari bentrok unique index '')
	empty := ""
	if got := emptyToNil(&empty); got != nil {
		t.Errorf("emptyToNil(\"\") = %q, want nil", *got)
	}

	// "abc" → pointer "abc" (tidak berubah)
	val := "abc"
	if got := emptyToNil(&val); got == nil || *got != "abc" {
		t.Errorf("emptyToNil(\"abc\") = %v, want pointer to \"abc\"", got)
	}
}

func TestService_Create_EmptySubdomainStoredAsNil(t *testing.T) {
	svc, fakeTM, cleanup := newTestService()
	defer cleanup()

	// Skip provisioning — cukup verifikasi penyimpanan company.
	fakeTM.ProvisionTenantFunc = func(companyID, dbName, dbUser, dbPassword, driverType string) (*database.TenantConnection, error) {
		return nil, fmt.Errorf("skip provisioning for test")
	}

	empty := ""
	resp, err := svc.Create(CreateCompanyRequest{
		Name:          "PT Alpha",
		Subdomain:     &empty,
		AdminName:     "Admin A",
		AdminEmail:    "admin.a@alpha.com",
		AdminPassword: "secret123",
	})
	if err != nil {
		t.Fatalf("Create #1 failed: %v", err)
	}
	if resp.Subdomain != nil {
		t.Errorf("expected subdomain nil (NULL), got %q", *resp.Subdomain)
	}

	// Company kedua dengan subdomain kosong juga harus sukses —
	// bukti tidak ada bentrok unique index '' (regresi potensial tanpa emptyToNil).
	if _, err := svc.Create(CreateCompanyRequest{
		Name:          "PT Beta",
		Subdomain:     &empty,
		AdminName:     "Admin B",
		AdminEmail:    "admin.b@beta.com",
		AdminPassword: "secret123",
	}); err != nil {
		t.Fatalf("Create #2 with empty subdomain failed (unique '' conflict?): %v", err)
	}
}

func TestService_Update_EmptySubdomainClearsToNil(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	company := createTestCompany(svc.repo.db, "Update Clear Sub")
	sub := "old-sub"
	company.Subdomain = &sub
	if err := svc.repo.Update(company); err != nil {
		t.Fatalf("Update setup failed: %v", err)
	}

	empty := ""
	resp, err := svc.Update(company.ID.String(), UpdateCompanyRequest{Subdomain: &empty})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if resp.Subdomain != nil {
		t.Errorf("expected subdomain cleared to nil, got %q", *resp.Subdomain)
	}

	// Verifikasi di DB.
	updated, err := svc.repo.FindByID(company.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if updated.Subdomain != nil {
		t.Errorf("expected subdomain NULL in DB, got %q", *updated.Subdomain)
	}
}

func TestService_Create_DuplicateSubdomain_Rejected(t *testing.T) {
	svc, fakeTM, cleanup := newTestService()
	defer cleanup()
	fakeTM.ProvisionTenantFunc = func(companyID, dbName, dbUser, dbPassword, driverType string) (*database.TenantConnection, error) {
		return nil, fmt.Errorf("skip provisioning for test")
	}

	sub := "pt-satu"
	if _, err := svc.Create(CreateCompanyRequest{
		Name:          "PT Satu",
		Subdomain:     &sub,
		AdminName:     "Admin 1",
		AdminEmail:    "admin1@satu.com",
		AdminPassword: "secret123",
	}); err != nil {
		t.Fatalf("Create #1 failed: %v", err)
	}

	if _, err := svc.Create(CreateCompanyRequest{
		Name:          "PT Dua",
		Subdomain:     &sub,
		AdminName:     "Admin 2",
		AdminEmail:    "admin2@dua.com",
		AdminPassword: "secret123",
	}); err == nil {
		t.Fatal("expected duplicate subdomain rejected, got nil")
	}
}

// =========================================================================
// Rotate Credentials Service Tests
// =========================================================================

func TestService_RotateCredentials_WithProvidedPassword(t *testing.T) {
	svc, fakeTM, cleanup := newTestService()
	defer cleanup()

	company := createTestCompany(svc.repo.db, "Rotate Test")

	var rotatedCompanyID, rotatedPassword string
	fakeTM.RotateTenantCredFunc = func(companyID, newPassword string) error {
		rotatedCompanyID = companyID
		rotatedPassword = newPassword
		return nil
	}

	resp, err := svc.RotateCredentials(company.ID.String(), "s3cure-pass-123")
	if err != nil {
		t.Fatalf("RotateCredentials failed: %v", err)
	}

	if !resp.Rotated {
		t.Error("expected Rotated=true")
	}
	if resp.NewPassword != "" {
		t.Errorf("expected empty NewPassword (user-provided), got %q", resp.NewPassword)
	}
	if rotatedCompanyID != company.ID.String() {
		t.Errorf("expected company ID %s, got %s", company.ID.String(), rotatedCompanyID)
	}
	if rotatedPassword != "s3cure-pass-123" {
		t.Errorf("expected password 's3cure-pass-123', got %q", rotatedPassword)
	}
}

func TestService_RotateCredentials_AutoGeneratePassword(t *testing.T) {
	svc, fakeTM, cleanup := newTestService()
	defer cleanup()

	company := createTestCompany(svc.repo.db, "Rotate Auto")

	var rotatedPassword string
	fakeTM.RotateTenantCredFunc = func(companyID, newPassword string) error {
		rotatedPassword = newPassword
		return nil
	}

	resp, err := svc.RotateCredentials(company.ID.String(), "")
	if err != nil {
		t.Fatalf("RotateCredentials failed: %v", err)
	}

	if resp.NewPassword == "" {
		t.Error("expected auto-generated NewPassword in response")
	}
	if resp.NewPassword != rotatedPassword {
		t.Errorf("expected returned password %q to match rotated %q", resp.NewPassword, rotatedPassword)
	}
	if len(resp.NewPassword) < 12 {
		t.Errorf("expected generated password length >= 12, got %d", len(resp.NewPassword))
	}
}

func TestService_RotateCredentials_TerminatedCompany(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	company := createTestCompany(svc.repo.db, "Rotate Terminated")
	company.Status = CompanyStatusTerminated
	if err := svc.repo.Update(company); err != nil {
		t.Fatalf("failed to update company: %v", err)
	}

	_, err := svc.RotateCredentials(company.ID.String(), "pass123456")
	if err == nil {
		t.Fatal("expected error for terminated company")
	}
}

func TestService_RotateCredentials_InvalidUUID(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.RotateCredentials("not-a-uuid", "pass123456")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_RotateCredentials_DBManagerError(t *testing.T) {
	svc, fakeTM, cleanup := newTestService()
	defer cleanup()

	company := createTestCompany(svc.repo.db, "Rotate DBError")

	fakeTM.RotateTenantCredFunc = func(companyID, newPassword string) error {
		return fmt.Errorf("simulated ALTER USER failure")
	}

	_, err := svc.RotateCredentials(company.ID.String(), "pass123456")
	if err == nil {
		t.Fatal("expected error when DB manager rotation fails")
	}
}

// =========================================================================
// Terminate Service Tests
// =========================================================================

func TestService_Terminate_ActiveCompany_Success(t *testing.T) {
	svc, fakeTM, cleanup := newTestService()
	defer cleanup()

	// Create an active company
	company := createTestCompany(svc.repo.db, "Test Company")

	droppedCalled := false
	removedCalled := false

	fakeTM.DropTenantDBFunc = func(companyID string) error {
		droppedCalled = true
		if companyID != company.ID.String() {
			t.Errorf("expected company ID %s, got %s", company.ID.String(), companyID)
		}
		return nil
	}

	fakeTM.RemoveTenantConnFunc = func(companyID string) error {
		removedCalled = true
		if companyID != company.ID.String() {
			t.Errorf("expected company ID %s, got %s", company.ID.String(), companyID)
		}
		return nil
	}

	resp, err := svc.Terminate(company.ID.String())
	if err != nil {
		t.Fatalf("Terminate failed: %v", err)
	}

	if resp.Status != "terminated" {
		t.Errorf("expected status 'terminated', got '%s'", resp.Status)
	}

	if !droppedCalled {
		t.Error("DropTenantDB was not called")
	}
	if !removedCalled {
		t.Error("RemoveTenantConnection was not called")
	}

	// Verify company record is updated in database
	updated, err := svc.repo.FindByID(company.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if updated.Status != CompanyStatusTerminated {
		t.Errorf("expected company status 'terminated', got '%s'", updated.Status)
	}
}

func TestService_Terminate_SuspendedCompany_Success(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	company := createTestCompany(svc.repo.db, "Suspended Co")
	company.Status = CompanyStatusSuspended
	if err := svc.repo.Update(company); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	resp, err := svc.Terminate(company.ID.String())
	if err != nil {
		t.Fatalf("Terminate suspended company failed: %v", err)
	}

	if resp.Status != "terminated" {
		t.Errorf("expected status 'terminated', got '%s'", resp.Status)
	}
}

func TestService_Terminate_AlreadyTerminated_Error(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	company := createTestCompany(svc.repo.db, "Already Terminated")
	company.Status = CompanyStatusTerminated
	if err := svc.repo.Update(company); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	_, err := svc.Terminate(company.ID.String())
	if err == nil {
		t.Fatal("expected error for already terminated company, got nil")
	}
}

func TestService_Terminate_InvalidUUID_Error(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.Terminate("not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID, got nil")
	}
}

func TestService_Terminate_NotFound_Error(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.Terminate(uuidStr())
	if err == nil {
		t.Fatal("expected error for non-existent company, got nil")
	}
}

func TestService_Terminate_DropDBFails_StillRemovesConnectionAndUpdatesStatus(t *testing.T) {
	svc, fakeTM, cleanup := newTestService()
	defer cleanup()

	company := createTestCompany(svc.repo.db, "Drop DB Fails Co")

	removedCalled := false

	// DropTenantDB returns error (simulates DB already gone)
	fakeTM.DropTenantDBFunc = func(companyID string) error {
		return fmt.Errorf("database not found")
	}

	fakeTM.RemoveTenantConnFunc = func(companyID string) error {
		removedCalled = true
		return nil
	}

	resp, err := svc.Terminate(company.ID.String())
	if err != nil {
		t.Fatalf("Terminate should succeed even if DropTenantDB fails: %v", err)
	}

	if resp.Status != "terminated" {
		t.Errorf("expected status 'terminated', got '%s'", resp.Status)
	}
	if !removedCalled {
		t.Error("RemoveTenantConnection should still be called even if DropTenantDB fails")
	}
}

func TestService_Terminate_RemoveConnFails_StillUpdatesStatus(t *testing.T) {
	svc, fakeTM, cleanup := newTestService()
	defer cleanup()

	company := createTestCompany(svc.repo.db, "Remove Conn Fails Co")

	fakeTM.DropTenantDBFunc = func(companyID string) error {
		return nil
	}

	fakeTM.RemoveTenantConnFunc = func(companyID string) error {
		return fmt.Errorf("connection record not found")
	}

	resp, err := svc.Terminate(company.ID.String())
	if err != nil {
		t.Fatalf("Terminate should succeed even if RemoveTenantConnection fails: %v", err)
	}

	if resp.Status != "terminated" {
		t.Errorf("expected status 'terminated', got '%s'", resp.Status)
	}
}
