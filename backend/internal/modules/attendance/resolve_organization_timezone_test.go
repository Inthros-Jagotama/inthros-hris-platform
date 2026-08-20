package attendance

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestResolveOrganizationTimezone_UsesZoneOverride(t *testing.T) {
	companyID := uuid.New()
	repo, db, ctx := newTestRepository(t, companyID, "Asia/Jakarta")

	zoneID := uuid.New()
	seedZone(db, zoneID, "Asia/Jayapura")
	orgID := uuid.New()
	seedOrgWithZone(db, orgID, &zoneID, "Test Org")

	loc, err := repo.ResolveOrganizationTimezone(ctx, orgID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.String() != "Asia/Jayapura" {
		t.Errorf("got %s, want Asia/Jayapura", loc.String())
	}
}

func TestResolveOrganizationTimezone_FallsBackToCompanyTimezone(t *testing.T) {
	companyID := uuid.New()
	repo, db, ctx := newTestRepository(t, companyID, "Asia/Makassar")

	orgID := uuid.New()
	// Organization with no zone_id set at all.
	seedOrgWithZone(db, orgID, nil, "Test Org No Zone")

	loc, err := repo.ResolveOrganizationTimezone(ctx, orgID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.String() != "Asia/Makassar" {
		t.Errorf("got %s, want Asia/Makassar", loc.String())
	}
}

// TestResolveOrganizationTimezone_SoftDeletedZoneFallsBackToCompany verifies
// a soft-deleted zone's timezone override is ignored — the join must filter
// deleted_at IS NULL on both organizations and zones, matching the
// convention used by GetChildOrganizationIDs elsewhere in this file.
func TestResolveOrganizationTimezone_SoftDeletedZoneFallsBackToCompany(t *testing.T) {
	companyID := uuid.New()
	repo, db, ctx := newTestRepository(t, companyID, "Asia/Jakarta")

	zoneID := uuid.New()
	seedZone(db, zoneID, "Asia/Jayapura")
	if err := db.Exec("UPDATE zones SET deleted_at = ? WHERE id = ?", time.Now(), zoneID.String()).Error; err != nil {
		t.Fatalf("failed to soft-delete zone: %v", err)
	}
	orgID := uuid.New()
	seedOrgWithZone(db, orgID, &zoneID, "Test Org Soft-Deleted Zone")

	loc, err := repo.ResolveOrganizationTimezone(ctx, orgID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.String() != "Asia/Jakarta" {
		t.Errorf("got %s, want Asia/Jakarta (company default, soft-deleted zone override must be ignored)", loc.String())
	}
}

func TestResolveOrganizationTimezone_ZoneWithNoTimezoneOverrideFallsBackToCompany(t *testing.T) {
	companyID := uuid.New()
	repo, db, ctx := newTestRepository(t, companyID, "Asia/Jakarta")

	zoneID := uuid.New()
	seedZone(db, zoneID, "") // zone exists but has no timezone override
	orgID := uuid.New()
	seedOrgWithZone(db, orgID, &zoneID, "Test Org Zone No TZ")

	loc, err := repo.ResolveOrganizationTimezone(ctx, orgID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.String() != "Asia/Jakarta" {
		t.Errorf("got %s, want Asia/Jakarta", loc.String())
	}
}
