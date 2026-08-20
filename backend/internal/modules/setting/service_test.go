package setting

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/platform/company"
)

// =========================================================================
// Zone — duplicate code validation (Create + Update)
// =========================================================================

func TestService_CreateZone_Success(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	ctx := context.Background()
	resp, err := svc.CreateZone(ctx, CreateZoneRequest{
		Code: "Z001", Name: "Zone One", Region: "Region A",
	})
	if err != nil {
		t.Fatalf("CreateZone failed: %v", err)
	}
	if resp.Code != "Z001" {
		t.Errorf("expected code 'Z001', got %q", resp.Code)
	}
	if resp.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestService_CreateZone_DuplicateCode(t *testing.T) {
	svc, db, cleanup := newTestService()
	defer cleanup()

	createTestZone(db, "Z001", "Zone One")

	ctx := context.Background()
	_, err := svc.CreateZone(ctx, CreateZoneRequest{
		Code: "Z001", Name: "Zone Duplicate",
	})
	if err == nil {
		t.Fatal("expected error for duplicate zone code, got nil")
	}
}

func TestService_CreateZone_DifferentCode_OK(t *testing.T) {
	svc, db, cleanup := newTestService()
	defer cleanup()

	createTestZone(db, "Z001", "Zone One")

	ctx := context.Background()
	resp, err := svc.CreateZone(ctx, CreateZoneRequest{
		Code: "Z002", Name: "Zone Two",
	})
	if err != nil {
		t.Fatalf("CreateZone failed for different code: %v", err)
	}
	if resp.Code != "Z002" {
		t.Errorf("expected code 'Z002', got %q", resp.Code)
	}
}

func TestService_UpdateZone_ChangeCode_Success(t *testing.T) {
	svc, db, cleanup := newTestService()
	defer cleanup()

	zone := createTestZone(db, "Z010", "Zone Ten")

	ctx := context.Background()
	newCode := "Z010-NEW"
	resp, err := svc.UpdateZone(ctx, zone.ID.String(), UpdateZoneRequest{
		Code: &newCode,
	})
	if err != nil {
		t.Fatalf("UpdateZone failed: %v", err)
	}
	if resp.Code != "Z010-NEW" {
		t.Errorf("expected code 'Z010-NEW', got %q", resp.Code)
	}
}

func TestService_UpdateZone_DuplicateCode_ExcludeSelf(t *testing.T) {
	svc, db, cleanup := newTestService()
	defer cleanup()

	z1 := createTestZone(db, "Z020", "Zone Twenty")
	createTestZone(db, "Z021", "Zone Twenty One")

	// Updating z1 to use z1's own code should succeed (exclude self)
	ctx := context.Background()
	code := "Z020"
	_, err := svc.UpdateZone(ctx, z1.ID.String(), UpdateZoneRequest{
		Code: &code,
	})
	if err != nil {
		t.Fatalf("UpdateZone with same code (exclude self) should succeed: %v", err)
	}
}

func TestService_UpdateZone_DuplicateCode_Conflict(t *testing.T) {
	svc, db, cleanup := newTestService()
	defer cleanup()

	z1 := createTestZone(db, "Z030", "Zone Thirty")
	createTestZone(db, "Z031", "Zone Thirty One")

	// Updating z1 to use z031's code should fail (duplicate)
	ctx := context.Background()
	code := "Z031"
	_, err := svc.UpdateZone(ctx, z1.ID.String(), UpdateZoneRequest{
		Code: &code,
	})
	if err == nil {
		t.Fatal("expected error when updating to an existing code")
	}
}

func TestCreateZone_RejectsInvalidTimezone(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	tz := "Asia/Singapore"
	_, err := svc.CreateZone(context.Background(), CreateZoneRequest{
		Code: "TZ1", Name: "Zone TZ1", Timezone: &tz,
	})
	if err != ErrInvalidZoneTimezone {
		t.Errorf("got %v, want ErrInvalidZoneTimezone", err)
	}
}

func TestCreateZone_AllowsNilTimezone(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	resp, err := svc.CreateZone(context.Background(), CreateZoneRequest{Code: "TZ2", Name: "Zone TZ2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Timezone != nil {
		t.Errorf("got %v, want nil", resp.Timezone)
	}
}

func TestCreateZone_AllowsValidTimezone(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	tz := "Asia/Jakarta"
	resp, err := svc.CreateZone(context.Background(), CreateZoneRequest{
		Code: "TZ3", Name: "Zone TZ3", Timezone: &tz,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Timezone == nil || *resp.Timezone != "Asia/Jakarta" {
		t.Errorf("got %v, want Asia/Jakarta", resp.Timezone)
	}
}

func TestUpdateZone_RejectsInvalidTimezone(t *testing.T) {
	svc, db, cleanup := newTestService()
	defer cleanup()

	zone := createTestZone(db, "TZ4", "Zone TZ4")

	tz := "Asia/Singapore"
	_, err := svc.UpdateZone(context.Background(), zone.ID.String(), UpdateZoneRequest{Timezone: &tz})
	if err != ErrInvalidZoneTimezone {
		t.Errorf("got %v, want ErrInvalidZoneTimezone", err)
	}
}

// =========================================================================
// Education — duplicate code validation
// =========================================================================

func TestService_CreateEducation_Success(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	ctx := context.Background()
	resp, err := svc.CreateEducation(ctx, CreateEducationRequest{
		Code: "ED01", Name: "SMA",
	})
	if err != nil {
		t.Fatalf("CreateEducation failed: %v", err)
	}
	if resp.Code != "ED01" {
		t.Errorf("expected code 'ED01', got %q", resp.Code)
	}
}

func TestService_CreateEducation_DuplicateCode(t *testing.T) {
	svc, db, cleanup := newTestService()
	defer cleanup()

	createTestEducation(db, "ED01", "SMA")

	ctx := context.Background()
	_, err := svc.CreateEducation(ctx, CreateEducationRequest{
		Code: "ED01", Name: "SMA Duplicate",
	})
	if err == nil {
		t.Fatal("expected error for duplicate education code")
	}
}

func TestService_UpdateEducation_DuplicateCode_ExcludeSelf(t *testing.T) {
	svc, db, cleanup := newTestService()
	defer cleanup()

	e := createTestEducation(db, "ED10", "S1")

	ctx := context.Background()
	code := "ED10"
	_, err := svc.UpdateEducation(ctx, e.ID.String(), UpdateEducationRequest{
		Code: &code,
	})
	if err != nil {
		t.Fatalf("UpdateEducation with own code should succeed: %v", err)
	}
}

func TestService_UpdateEducation_DuplicateCode_Conflict(t *testing.T) {
	svc, db, cleanup := newTestService()
	defer cleanup()

	createTestEducation(db, "ED20", "SMP")
	e2 := createTestEducation(db, "ED21", "SMA")

	ctx := context.Background()
	code := "ED20"
	_, err := svc.UpdateEducation(ctx, e2.ID.String(), UpdateEducationRequest{
		Code: &code,
	})
	if err == nil {
		t.Fatal("expected error when updating education to existing code")
	}
}

// =========================================================================
// Religion — duplicate code validation
// =========================================================================

func TestService_CreateReligion_DuplicateCode(t *testing.T) {
	svc, db, cleanup := newTestService()
	defer cleanup()

	createTestReligion(db, "ISLAM", "Islam")

	ctx := context.Background()
	_, err := svc.CreateReligion(ctx, CreateReligionRequest{
		Code: "ISLAM", Name: "Islam (dup)",
	})
	if err == nil {
		t.Fatal("expected error for duplicate religion code")
	}
}

func TestService_UpdateReligion_DuplicateCode_Conflict(t *testing.T) {
	svc, db, cleanup := newTestService()
	defer cleanup()

	createTestReligion(db, "KATOLIK", "Katolik")
	r2 := createTestReligion(db, "PROTESTAN", "Protestan")

	ctx := context.Background()
	code := "KATOLIK"
	_, err := svc.UpdateReligion(ctx, r2.ID.String(), UpdateReligionRequest{
		Code: &code,
	})
	if err == nil {
		t.Fatal("expected error when updating religion to existing code")
	}
}

// =========================================================================
// MaritalStatus — duplicate code validation
// =========================================================================

func TestService_CreateMaritalStatus_DuplicateCode(t *testing.T) {
	svc, db, cleanup := newTestService()
	defer cleanup()

	createTestMaritalStatus(db, "KAWIN", "Kawin")

	ctx := context.Background()
	_, err := svc.CreateMaritalStatus(ctx, CreateMaritalStatusRequest{
		Code: "KAWIN", Name: "Kawin (dup)",
	})
	if err == nil {
		t.Fatal("expected error for duplicate marital status code")
	}
}

func TestService_UpdateMaritalStatus_DuplicateCode_ExcludeSelf(t *testing.T) {
	svc, db, cleanup := newTestService()
	defer cleanup()

	m := createTestMaritalStatus(db, "BELUM-KAWIN", "Belum Kawin")

	ctx := context.Background()
	code := "BELUM-KAWIN"
	_, err := svc.UpdateMaritalStatus(ctx, m.ID.String(), UpdateMaritalStatusRequest{
		Code: &code,
	})
	if err != nil {
		t.Fatalf("UpdateMaritalStatus with own code should succeed: %v", err)
	}
}

// =========================================================================
// Bank — duplicate code validation
// =========================================================================

func TestService_CreateBank_DuplicateCode(t *testing.T) {
	svc, db, cleanup := newTestService()
	defer cleanup()

	createTestBank(db, "BNI", "Bank Negara Indonesia")

	ctx := context.Background()
	_, err := svc.CreateBank(ctx, CreateBankRequest{
		Code: "BNI", Name: "BNI (dup)",
	})
	if err == nil {
		t.Fatal("expected error for duplicate bank code")
	}
}

func TestService_UpdateBank_DuplicateCode_Conflict(t *testing.T) {
	svc, db, cleanup := newTestService()
	defer cleanup()

	createTestBank(db, "BRI", "Bank Rakyat Indonesia")
	b2 := createTestBank(db, "MANDIRI", "Bank Mandiri")

	ctx := context.Background()
	code := "BRI"
	_, err := svc.UpdateBank(ctx, b2.ID.String(), UpdateBankRequest{
		Code: &code,
	})
	if err == nil {
		t.Fatal("expected error when updating bank to existing code")
	}
}

// =========================================================================
// Nationality — duplicate code validation
// =========================================================================

func TestService_CreateNationality_DuplicateCode(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	// Insert directly via repo to test service validation
	ctx := context.Background()
	svc.CreateNationality(ctx, CreateNationalityRequest{Code: "ID", Name: "Indonesia"})

	_, err := svc.CreateNationality(ctx, CreateNationalityRequest{Code: "ID", Name: "Indonesia (dup)"})
	if err == nil {
		t.Fatal("expected error for duplicate nationality code")
	}
}

// =========================================================================
// JobFamily — duplicate code validation
// =========================================================================

func TestService_CreateJobFamily_DuplicateCode(t *testing.T) {
	svc, db, cleanup := newTestService()
	defer cleanup()

	createTestJobFamily(db, "IT", "Information Technology")

	ctx := context.Background()
	_, err := svc.CreateJobFamily(ctx, CreateJobFamilyRequest{
		Code: "IT", Name: "IT (dup)",
	})
	if err == nil {
		t.Fatal("expected error for duplicate job family code")
	}
}

// =========================================================================
// Grading — duplicate code validation
// =========================================================================

func TestService_CreateGrading_Success(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	ctx := context.Background()
	resp, err := svc.CreateGrading(ctx, CreateGradingRequest{
		Code: "GRD-01", Name: "Grade 1", Description: "Entry Level",
	})
	if err != nil {
		t.Fatalf("CreateGrading failed: %v", err)
	}
	if resp.Code != "GRD-01" {
		t.Errorf("expected code 'GRD-01', got %q", resp.Code)
	}
}

func TestService_CreateGrading_DuplicateCode(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	ctx := context.Background()
	svc.CreateGrading(ctx, CreateGradingRequest{Code: "GRD-A", Name: "Grade A"})
	_, err := svc.CreateGrading(ctx, CreateGradingRequest{Code: "GRD-A", Name: "Grade A Dup"})
	if err == nil {
		t.Fatal("expected error for duplicate grading code")
	}
}

// =========================================================================
// SalaryGrade — duplicate code validation
// =========================================================================

func TestService_CreateSalaryGrade_Success(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	ctx := context.Background()
	resp, err := svc.CreateSalaryGrade(ctx, CreateSalaryGradeRequest{
		Code: "SG-01", Name: "Salary Grade 1", MinAmount: 1000000, MaxAmount: 5000000,
	})
	if err != nil {
		t.Fatalf("CreateSalaryGrade failed: %v", err)
	}
	if resp.Code != "SG-01" {
		t.Errorf("expected code 'SG-01', got %q", resp.Code)
	}
}

func TestService_CreateSalaryGrade_DuplicateCode(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	ctx := context.Background()
	svc.CreateSalaryGrade(ctx, CreateSalaryGradeRequest{Code: "SG-A", Name: "SG A"})
	_, err := svc.CreateSalaryGrade(ctx, CreateSalaryGradeRequest{Code: "SG-A", Name: "SG A Dup"})
	if err == nil {
		t.Fatal("expected error for duplicate salary grade code")
	}
}

// =========================================================================
// EmploymentStatus — duplicate code validation
// =========================================================================

func TestService_CreateEmploymentStatus_DuplicateCode(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	ctx := context.Background()
	svc.CreateEmploymentStatus(ctx, CreateEmploymentStatusRequest{Code: "PERM", Name: "Permanent"})
	_, err := svc.CreateEmploymentStatus(ctx, CreateEmploymentStatusRequest{Code: "PERM", Name: "Permanent (dup)"})
	if err == nil {
		t.Fatal("expected error for duplicate employment status code")
	}
}

// =========================================================================
// RelationshipType — duplicate code validation
// =========================================================================

func TestService_CreateRelationshipType_DuplicateCode(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	ctx := context.Background()
	svc.CreateRelationshipType(ctx, CreateRelationshipTypeRequest{Code: "SUAMI", Name: "Suami"})
	_, err := svc.CreateRelationshipType(ctx, CreateRelationshipTypeRequest{Code: "SUAMI", Name: "Suami (dup)"})
	if err == nil {
		t.Fatal("expected error for duplicate relationship type code")
	}
}

// =========================================================================
// ValidateUniqueCode edge cases
// =========================================================================

func TestService_ValidateUniqueCode_EmptyDB(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	ctx := context.Background()
	err := svc.validateUniqueCode(ctx, &Zone{}, "FIRST", "zones")
	if err != nil {
		t.Fatalf("validateUniqueCode on empty DB should succeed: %v", err)
	}
}

func TestService_ValidateUniqueCodeExcludeSelf_NonExistent(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	ctx := context.Background()
	err := svc.validateUniqueCodeExcludeSelf(ctx, &Zone{}, "NONEXIST", uuidStr(), "zones")
	if err != nil {
		t.Fatalf("validateUniqueCodeExcludeSelf for non-existent code should succeed: %v", err)
	}
}

// =========================================================================
// Company Timezone
// =========================================================================

func TestUpdateCompanyTimezone_RejectsInvalidValue(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	err := svc.UpdateCompanyTimezone(context.Background(), "Asia/Singapore")
	if err != ErrInvalidCompanyTimezone {
		t.Errorf("got %v, want ErrInvalidCompanyTimezone", err)
	}
}

func TestUpdateCompanyTimezone_PersistsAndGetReadsBack(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	platformDB, platformCleanup := setupTestPlatformDB()
	defer platformCleanup()

	comp := &company.Company{ID: uuid.New(), Name: "Acme", Slug: "acme", Timezone: "Asia/Jakarta"}
	if err := platformDB.Create(comp).Error; err != nil {
		t.Fatalf("failed to seed test company: %v", err)
	}

	repo := NewRepositoryWithPlatformDB(testDBResolver(db), platformDB)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)

	ctx := context.WithValue(context.Background(), "company_id", comp.ID.String())

	if err := svc.UpdateCompanyTimezone(ctx, "Asia/Makassar"); err != nil {
		t.Fatalf("UpdateCompanyTimezone should succeed: %v", err)
	}

	tz, err := svc.GetCompanyTimezone(ctx)
	if err != nil {
		t.Fatalf("GetCompanyTimezone should succeed: %v", err)
	}
	if tz != "Asia/Makassar" {
		t.Errorf("got timezone %q, want %q", tz, "Asia/Makassar")
	}
}

func TestGetCompanyTimezone_CompanyNotFound(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	ctx := context.WithValue(context.Background(), "company_id", uuid.New().String())
	_, err := svc.GetCompanyTimezone(ctx)
	if !errors.Is(err, ErrCompanyNotFound) {
		t.Errorf("got %v, want ErrCompanyNotFound", err)
	}
}

func TestUpdateCompanyTimezone_CompanyNotFound(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	ctx := context.WithValue(context.Background(), "company_id", uuid.New().String())
	err := svc.UpdateCompanyTimezone(ctx, "Asia/Makassar")
	if !errors.Is(err, ErrCompanyNotFound) {
		t.Errorf("got %v, want ErrCompanyNotFound", err)
	}
}
