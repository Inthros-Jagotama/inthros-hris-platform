package company

import (
	"fmt"
	"testing"
)

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
