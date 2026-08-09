package modulemgmt

import (
	"fmt"
	"testing"

	sqlite "github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database and auto-migrates
// PlatformModule for repository-level tests.
func setupTestDB() (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}
	if err := db.AutoMigrate(&PlatformModule{}); err != nil {
		panic(fmt.Sprintf("failed to migrate test db: %v", err))
	}
	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
	return db, cleanup
}

func createTestModule(db *gorm.DB, name, slug, moduleType string) *PlatformModule {
	m := &PlatformModule{
		Name:       name,
		Slug:       slug,
		Version:    "1.0.0",
		ModuleType: moduleType,
	}
	if err := db.Create(m).Error; err != nil {
		panic(fmt.Sprintf("failed to create test module: %v", err))
	}
	return m
}

// TestRepo_FindAll_Pagination guards the server-side pagination contract the
// FE relies on: only `perPage` rows come back per call, and `total` reflects
// the full unpaginated count (not just what's on the current page) — the FE
// Modules.vue page uses these to drive PrimeVue's lazy DataTable paginator.
func TestRepo_FindAll_Pagination(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(db)

	for i := 0; i < 5; i++ {
		createTestModule(db, fmt.Sprintf("Module %d", i), fmt.Sprintf("module-%d", i), "tenant")
	}

	modules, total, err := repo.FindAll(1, 2, "", "")
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(modules) != 2 {
		t.Errorf("expected 2 modules on page 1 (per_page=2), got %d", len(modules))
	}

	page2, total2, err := repo.FindAll(2, 2, "", "")
	if err != nil {
		t.Fatalf("FindAll page 2 failed: %v", err)
	}
	if total2 != 5 {
		t.Errorf("expected total 5 on page 2, got %d", total2)
	}
	if len(page2) != 2 {
		t.Errorf("expected 2 modules on page 2, got %d", len(page2))
	}
}

// TestRepo_FindAll_SearchFilter guards the new `search` param added
// alongside server-side pagination — without it, moving the FE's search box
// to filter against only the currently-loaded page (instead of the full
// dataset) would silently break search once paging became server-driven.
func TestRepo_FindAll_SearchFilter(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(db)

	createTestModule(db, "Attendance Management", "attendance", "tenant")
	createTestModule(db, "Payroll Management", "payroll", "tenant")
	createTestModule(db, "Leave Management", "leave", "tenant")

	modules, total, err := repo.FindAll(1, 20, "", "attend")
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected 1 match for 'attend', got %d", total)
	}
	if len(modules) != 1 || modules[0].Slug != "attendance" {
		t.Errorf("expected the attendance module, got %+v", modules)
	}
}

// TestRepo_FindAll_ModuleTypeAndSearchCombined ensures the search filter
// composes with the existing module_type filter rather than overriding it.
func TestRepo_FindAll_ModuleTypeAndSearchCombined(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(db)

	createTestModule(db, "Approval Engine", "approval", "platform")
	createTestModule(db, "Approval Requests", "approval-tenant", "tenant")

	modules, total, err := repo.FindAll(1, 20, "tenant", "approval")
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected 1 match (tenant + 'approval'), got %d", total)
	}
	if len(modules) != 1 || modules[0].Slug != "approval-tenant" {
		t.Errorf("expected only the tenant-type module, got %+v", modules)
	}
}
