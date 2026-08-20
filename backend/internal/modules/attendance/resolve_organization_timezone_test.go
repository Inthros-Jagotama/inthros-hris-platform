package attendance

import (
	"testing"

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
