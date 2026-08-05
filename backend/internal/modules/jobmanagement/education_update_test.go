package jobmanagement

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// TestService_UpdateJobEducationExperience_PersistsFKChange — regression test.
//
// Bug lama: FindJobEducationExperienceByID me-Preload belongs-to Education/Experience,
// lalu repo Update memakai Omit("Majors","JobFamilies").Save(e) — GORM Save mengembalikan
// (revert) education_id/experience_id ke nilai association yang di-load, sehingga
// perubahan FK tidak tersimpan di DB (endpoint PUT sukses tapi data tidak terupdate).
func TestService_UpdateJobEducationExperience_PersistsFKChange(t *testing.T) {
	_, resolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(resolver)
	svc := NewService(repo, nil)
	ctx := context.Background()

	orgID := createTestOrgID()
	eduA := createTestJobValue(ctx, repo, "education")
	eduB := createTestJobValue(ctx, repo, "education")
	expA := createTestJobValue(ctx, repo, "experience")
	expB := createTestJobValue(ctx, repo, "experience")

	created, err := svc.CreateJobEducationExperience(ctx, CreateJobEducationExperienceRequest{
		OrganizationID:   orgID,
		Nomenclature:     "Org",
		FullCode:         "01",
		EducationID:      strPtr(eduA.ID.String()),
		ExperienceID:     strPtr(expA.ID.String()),
		EducationMajorID: []string{},
		JobFamilyID:      []string{},
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// Ganti pendidikan & pengalaman → harus tersimpan di DB
	if _, err := svc.UpdateJobEducationExperience(ctx, created.ID, UpdateJobEducationExperienceRequest{
		EducationID:      strPtr(eduB.ID.String()),
		ExperienceID:     strPtr(expB.ID.String()),
		EducationMajorID: []string{},
		JobFamilyID:      []string{},
	}); err != nil {
		t.Fatalf("update: %v", err)
	}

	got, err := repo.FindJobEducationExperienceByID(ctx, uuid.MustParse(created.ID))
	if err != nil {
		t.Fatalf("refetch: %v", err)
	}
	if got.EducationID == nil || got.EducationID.String() != eduB.ID.String() {
		t.Fatalf("education_id NOT updated: got %v want %s", got.EducationID, eduB.ID.String())
	}
	if got.ExperienceID == nil || got.ExperienceID.String() != expB.ID.String() {
		t.Fatalf("experience_id NOT updated: got %v want %s", got.ExperienceID, expB.ID.String())
	}
}
