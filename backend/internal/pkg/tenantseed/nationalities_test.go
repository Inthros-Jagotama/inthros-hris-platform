package tenantseed

import (
	"testing"

	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

// newNationalityTestDB membuat in-memory SQLite dengan tabel nationalities
// (mengikuti pola newCompetencyTestDB — glebarez/sqlite).
func newNationalityTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	// DSN unik per test — hindari shared-cache SQLite yang berbagi DB antar test.
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	// Pool 1 koneksi agar DB memory tidak hilang antar query.
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	if err := db.Exec(`CREATE TABLE nationalities (
		id TEXT PRIMARY KEY,
		code TEXT NOT NULL,
		name TEXT NOT NULL,
		sort_order INTEGER,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME
	)`).Error; err != nil {
		t.Fatalf("failed to create nationalities table: %v", err)
	}
	return db
}

// TestSeedNationalities_InsertsAllData memastikan seluruh nationality dari
// docs/seeder/CountriesTableSeeder.php (240 negara ISO alpha-2) + "Lainnya / Other"
// ter-seed (total 241).
func TestSeedNationalities_InsertsAllData(t *testing.T) {
	db := newNationalityTestDB(t)

	inserted, skipped, err := seedNationalities(db)
	if err != nil {
		t.Fatalf("seedNationalities failed: %v", err)
	}
	if inserted != 241 {
		t.Fatalf("expected 241 nationalities inserted (240 + LNY), got %d", inserted)
	}
	if skipped != 0 {
		t.Fatalf("expected 0 skipped on first run, got %d", skipped)
	}

	var total int64
	if err := db.Table("nationalities").Count(&total).Error; err != nil {
		t.Fatalf("count failed: %v", err)
	}
	if total != 241 {
		t.Fatalf("expected 241 nationalities in table, got %d", total)
	}
}

// TestSeedNationalities_SpotCheck memastikan record yang diketahui ter-seed benar.
func TestSeedNationalities_SpotCheck(t *testing.T) {
	db := newNationalityTestDB(t)

	if _, _, err := seedNationalities(db); err != nil {
		t.Fatalf("seedNationalities failed: %v", err)
	}

	var name string
	if err := db.Table("nationalities").Where("code = ?", "US").Pluck("name", &name).Error; err != nil {
		t.Fatalf("query US failed: %v", err)
	}
	if name != "United States" {
		t.Fatalf("expected US = United States, got %q", name)
	}

	// LNY "Lainnya / Other" harus ada di sort_order terakhir (999)
	var sortOrder int64
	if err := db.Table("nationalities").Where("code = ?", "LNY").Pluck("sort_order", &sortOrder).Error; err != nil {
		t.Fatalf("query LNY failed: %v", err)
	}
	if sortOrder != 999 {
		t.Fatalf("expected LNY sort_order 999, got %d", sortOrder)
	}
}

// TestSeedNationalities_Idempotent memastikan menjalankan ulang seeder
// tidak menambah/menduplikasi data (skip semua, count tetap 241).
func TestSeedNationalities_Idempotent(t *testing.T) {
	db := newNationalityTestDB(t)

	if _, _, err := seedNationalities(db); err != nil {
		t.Fatalf("first run failed: %v", err)
	}

	inserted, skipped, err := seedNationalities(db)
	if err != nil {
		t.Fatalf("second run failed: %v", err)
	}
	if inserted != 0 {
		t.Fatalf("expected 0 inserted on second run, got %d", inserted)
	}
	if skipped != 241 {
		t.Fatalf("expected 241 skipped on second run, got %d", skipped)
	}

	var total int64
	if err := db.Table("nationalities").Count(&total).Error; err != nil {
		t.Fatalf("count failed: %v", err)
	}
	if total != 241 {
		t.Fatalf("expected 241 nationalities after re-run, got %d", total)
	}
}
