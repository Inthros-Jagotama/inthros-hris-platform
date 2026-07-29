package setting

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// =========================================================================
// Generic helpers: findByCode
// =========================================================================

func TestRepository_FindByCode_Found(t *testing.T) {
	repo, db, cleanup := newTestRepo()
	defer cleanup()

	createTestZone(db, "Z001", "Zone One")

	ctx := context.Background()
	found, err := repo.findByCode(ctx, &Zone{}, "Z001", "zones")
	if err != nil {
		t.Fatalf("findByCode failed: %v", err)
	}
	if !found {
		t.Error("expected findByCode to return true for existing code")
	}
}

func TestRepository_FindByCode_NotFound(t *testing.T) {
	repo, _, cleanup := newTestRepo()
	defer cleanup()

	ctx := context.Background()
	found, err := repo.findByCode(ctx, &Zone{}, "NONEXISTENT", "zones")
	if err != nil {
		t.Fatalf("findByCode failed: %v", err)
	}
	if found {
		t.Error("expected findByCode to return false for non-existent code")
	}
}

func TestRepository_FindByCode_SoftDeleted(t *testing.T) {
	repo, db, cleanup := newTestRepo()
	defer cleanup()

	zone := createTestZone(db, "Z002", "Zone Two")

	// Soft delete
	if err := repo.DeleteZone(context.Background(), zone.ID); err != nil {
		t.Fatalf("DeleteZone failed: %v", err)
	}

	// findByCode should NOT find the soft-deleted record
	ctx := context.Background()
	found, err := repo.findByCode(ctx, &Zone{}, "Z002", "zones")
	if err != nil {
		t.Fatalf("findByCode failed: %v", err)
	}
	if found {
		t.Error("expected findByCode to return false for soft-deleted record")
	}
}

func TestRepository_FindByCodeExcludeSelf_Excludes(t *testing.T) {
	repo, db, cleanup := newTestRepo()
	defer cleanup()

	z1 := createTestZone(db, "Z003", "Zone Three")
	createTestZone(db, "Z004", "Zone Four")

	ctx := context.Background()
	// Z003 should be found by code "Z003" even when excluding itself
	// Wait — findByCodeExcludeSelf excludes the given ID.
	// Z003 exists with code Z003. If we exclude z1.ID, we should NOT find it.
	found, err := repo.findByCodeExcludeSelf(ctx, &Zone{}, "Z003", z1.ID, "zones")
	if err != nil {
		t.Fatalf("findByCodeExcludeSelf failed: %v", err)
	}
	if found {
		t.Error("expected findByCodeExcludeSelf to return false when excluding own ID")
	}

	// Z003's code should be found when excluding a different ID
	found, err = repo.findByCodeExcludeSelf(ctx, &Zone{}, "Z003", uuid.New(), "zones")
	if err != nil {
		t.Fatalf("findByCodeExcludeSelf failed: %v", err)
	}
	if !found {
		t.Error("expected findByCodeExcludeSelf to return true when code exists for other ID")
	}
}

// =========================================================================
// Zone CRUD
// =========================================================================

func TestRepository_CreateZone_Success(t *testing.T) {
	repo, db, cleanup := newTestRepo()
	defer cleanup()

	ctx := context.Background()
	zone := &Zone{Code: "ZA01", Name: "Zone Alpha", Zone: "Zone Alpha"}
	if err := repo.CreateZone(ctx, zone); err != nil {
		t.Fatalf("CreateZone failed: %v", err)
	}

	// Verify in DB
	var saved Zone
	if err := db.First(&saved, "id = ?", zone.ID).Error; err != nil {
		t.Fatalf("failed to find created zone: %v", err)
	}
	if saved.Code != "ZA01" {
		t.Errorf("expected code 'ZA01', got %q", saved.Code)
	}
}

func TestRepository_FindZoneByID_Found(t *testing.T) {
	repo, db, cleanup := newTestRepo()
	defer cleanup()

	zone := createTestZone(db, "ZB01", "Zone Beta")

	ctx := context.Background()
	found, err := repo.FindZoneByID(ctx, zone.ID)
	if err != nil {
		t.Fatalf("FindZoneByID failed: %v", err)
	}
	if found.Code != "ZB01" {
		t.Errorf("expected code 'ZB01', got %q", found.Code)
	}
}

func TestRepository_FindZoneByID_NotFound(t *testing.T) {
	repo, _, cleanup := newTestRepo()
	defer cleanup()

	ctx := context.Background()
	_, err := repo.FindZoneByID(ctx, uuid.New())
	if err == nil {
		t.Error("expected error for non-existent zone")
	}
}

func TestRepository_FindAllZones(t *testing.T) {
	repo, db, cleanup := newTestRepo()
	defer cleanup()

	createTestZone(db, "ZC01", "Zone Charlie")
	createTestZone(db, "ZC02", "Zone Delta")

	ctx := context.Background()
	zones, total, err := repo.FindAllZones(ctx, 1, 20)
	if err != nil {
		t.Fatalf("FindAllZones failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(zones) != 2 {
		t.Errorf("expected 2 zones, got %d", len(zones))
	}
}

func TestRepository_UpdateZone(t *testing.T) {
	repo, db, cleanup := newTestRepo()
	defer cleanup()

	zone := createTestZone(db, "ZE01", "Zone Echo")

	ctx := context.Background()
	zone.Name = "Zone Echo Updated"
	if err := repo.UpdateZone(ctx, zone); err != nil {
		t.Fatalf("UpdateZone failed: %v", err)
	}

	var saved Zone
	db.First(&saved, "id = ?", zone.ID)
	if saved.Name != "Zone Echo Updated" {
		t.Errorf("expected name 'Zone Echo Updated', got %q", saved.Name)
	}
}

func TestRepository_DeleteZone_SoftDelete(t *testing.T) {
	repo, db, cleanup := newTestRepo()
	defer cleanup()

	zone := createTestZone(db, "ZF01", "Zone Foxtrot")

	ctx := context.Background()
	if err := repo.DeleteZone(ctx, zone.ID); err != nil {
		t.Fatalf("DeleteZone failed: %v", err)
	}

	// Should not be found via normal query
	_, err := repo.FindZoneByID(ctx, zone.ID)
	if err == nil {
		t.Error("expected error for soft-deleted zone")
	}

	// But should still exist in DB (soft delete)
	var count int64
	db.Unscoped().Model(&Zone{}).Where("id = ?", zone.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 record (soft-deleted) in DB, got %d", count)
	}
}

// =========================================================================
// Education CRUD (standalone, no BeforeCreate complexity)
// =========================================================================

func TestRepository_CreateEducation_Success(t *testing.T) {
	repo, db, cleanup := newTestRepo()
	defer cleanup()

	ctx := context.Background()
	e := &Education{Code: "ED01", Name: "SMA"}
	if err := repo.CreateEducation(ctx, e); err != nil {
		t.Fatalf("CreateEducation failed: %v", err)
	}

	var saved Education
	if err := db.First(&saved, "id = ?", e.ID).Error; err != nil {
		t.Fatalf("failed to find created education: %v", err)
	}
	if saved.Code != "ED01" {
		t.Errorf("expected code 'ED01', got %q", saved.Code)
	}
}

func TestRepository_FindEducationByID_Found(t *testing.T) {
	repo, db, cleanup := newTestRepo()
	defer cleanup()

	e := createTestEducation(db, "ED02", "S1")

	ctx := context.Background()
	found, err := repo.FindEducationByID(ctx, e.ID)
	if err != nil {
		t.Fatalf("FindEducationByID failed: %v", err)
	}
	if found.Code != "ED02" {
		t.Errorf("expected code 'ED02', got %q", found.Code)
	}
}

func TestRepository_FindEducationByID_NotFound(t *testing.T) {
	repo, _, cleanup := newTestRepo()
	defer cleanup()

	ctx := context.Background()
	_, err := repo.FindEducationByID(ctx, uuid.New())
	if err == nil {
		t.Error("expected error for non-existent education")
	}
}

func TestRepository_FindAllEducations(t *testing.T) {
	repo, db, cleanup := newTestRepo()
	defer cleanup()

	createTestEducation(db, "ED03", "SMP")
	createTestEducation(db, "ED04", "SD")

	ctx := context.Background()
	list, total, err := repo.FindAllEducations(ctx, 1, 20)
	if err != nil {
		t.Fatalf("FindAllEducations failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 educations, got %d", len(list))
	}
}

func TestRepository_UpdateEducation(t *testing.T) {
	repo, db, cleanup := newTestRepo()
	defer cleanup()

	e := createTestEducation(db, "ED05", "Diploma")

	ctx := context.Background()
	e.Name = "Diploma Updated"
	if err := repo.UpdateEducation(ctx, e); err != nil {
		t.Fatalf("UpdateEducation failed: %v", err)
	}

	var saved Education
	db.First(&saved, "id = ?", e.ID)
	if saved.Name != "Diploma Updated" {
		t.Errorf("expected name 'Diploma Updated', got %q", saved.Name)
	}
}

func TestRepository_DeleteEducation_SoftDelete(t *testing.T) {
	repo, db, cleanup := newTestRepo()
	defer cleanup()

	e := createTestEducation(db, "ED06", "Doktor")

	ctx := context.Background()
	if err := repo.DeleteEducation(ctx, e.ID); err != nil {
		t.Fatalf("DeleteEducation failed: %v", err)
	}

	_, err := repo.FindEducationByID(ctx, e.ID)
	if err == nil {
		t.Error("expected error for soft-deleted education")
	}
}

// =========================================================================
// Province CRUD (string PK, no UUID)
// =========================================================================

func TestRepository_CreateProvince_Success(t *testing.T) {
	repo, db, cleanup := newTestRepo()
	defer cleanup()

	ctx := context.Background()
	p := &Province{ID: "11", Code: "11", Name: "Aceh"}
	if err := repo.CreateProvince(ctx, p); err != nil {
		t.Fatalf("CreateProvince failed: %v", err)
	}

	var saved Province
	if err := db.First(&saved, "id = ?", "11").Error; err != nil {
		t.Fatalf("failed to find created province: %v", err)
	}
	if saved.Name != "Aceh" {
		t.Errorf("expected name 'Aceh', got %q", saved.Name)
	}
}

func TestRepository_FindProvinceByID_Found(t *testing.T) {
	repo, _, cleanup := newTestRepo()
	defer cleanup()

	ctx := context.Background()
	p := &Province{ID: "12", Code: "12", Name: "Sumatera Utara"}
	if err := repo.CreateProvince(ctx, p); err != nil {
		t.Fatalf("CreateProvince failed: %v", err)
	}

	found, err := repo.FindProvinceByID(ctx, "12")
	if err != nil {
		t.Fatalf("FindProvinceByID failed: %v", err)
	}
	if found.Name != "Sumatera Utara" {
		t.Errorf("expected name 'Sumatera Utara', got %q", found.Name)
	}
}
