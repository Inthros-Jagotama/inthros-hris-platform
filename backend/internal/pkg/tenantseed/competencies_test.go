package tenantseed

import (
	"testing"

	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

// newCompetencyTestDB membuat in-memory SQLite dengan tabel competencies
// (mengikuti pola setupTestDB di module lain — glebarez/sqlite).
func newCompetencyTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	// DSN unik per test — hindari shared-cache SQLite yang berbagi DB antar test.
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	// Pool 1 koneksi agar DB memory tidak hilang antar query.
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	if err := db.Exec(`CREATE TABLE competencies (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		field TEXT,
		cluster TEXT,
		definition TEXT,
		created_by TEXT,
		updated_by TEXT,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME
	)`).Error; err != nil {
		t.Fatalf("failed to create competencies table: %v", err)
	}
	return db
}

// TestSeedCompetencies_InsertsAllData memastikan seluruh kompetensi dari
// docs/seeder/CompetenciesTableSeeder.php (165 item) ter-seed.
func TestSeedCompetencies_InsertsAllData(t *testing.T) {
	db := newCompetencyTestDB(t)

	inserted, skipped, err := seedCompetencies(db)
	if err != nil {
		t.Fatalf("seedCompetencies failed: %v", err)
	}
	if inserted != 165 {
		t.Fatalf("expected 165 competencies inserted (match PHP seeder), got %d", inserted)
	}
	if skipped != 0 {
		t.Fatalf("expected 0 skipped on first run, got %d", skipped)
	}

	var total int64
	if err := db.Table("competencies").Count(&total).Error; err != nil {
		t.Fatalf("count failed: %v", err)
	}
	if total != 165 {
		t.Fatalf("expected 165 rows, got %d", total)
	}
}

// TestSeedCompetencies_SpotCheck memverifikasi beberapa record penting:
// field/cluster terbaca benar.
func TestSeedCompetencies_SpotCheck(t *testing.T) {
	db := newCompetencyTestDB(t)
	if _, _, err := seedCompetencies(db); err != nil {
		t.Fatalf("seedCompetencies failed: %v", err)
	}

	checks := []struct {
		name, field, cluster string
	}{
		{"Tenacity (TEN)", "Potential", "Potential"},
		{"Integrity (INT)", "Core", "Core"},
		{"Leadership (LEA)", "Manajerial", "Manajerial"},
		{"Data Center Management (DCM)", "Technical Competency", "Technology"},
		{"Financial Accounting (ACC)", "Technical Competency", "Finance & Accounting"},
	}

	for _, c := range checks {
		var gotField, gotCluster string
		if err := db.Table("competencies").
			Select("field, cluster").
			Where("name = ?", c.name).
			Row().Scan(&gotField, &gotCluster); err != nil {
			t.Fatalf("competency %q not found: %v", c.name, err)
		}
		if gotField != c.field {
			t.Errorf("competency %q: expected field %q, got %q", c.name, c.field, gotField)
		}
		if gotCluster != c.cluster {
			t.Errorf("competency %q: expected cluster %q, got %q", c.name, c.cluster, gotCluster)
		}
	}
}

// TestSeedCompetencies_DefinitionNotEmpty memastikan tidak ada definition kosong.
func TestSeedCompetencies_DefinitionNotEmpty(t *testing.T) {
	db := newCompetencyTestDB(t)
	if _, _, err := seedCompetencies(db); err != nil {
		t.Fatalf("seedCompetencies failed: %v", err)
	}

	var empty int64
	if err := db.Table("competencies").
		Where("definition IS NULL OR definition = ''").Count(&empty).Error; err != nil {
		t.Fatalf("count failed: %v", err)
	}
	if empty != 0 {
		t.Fatalf("expected 0 competencies with empty definition, got %d", empty)
	}
}

// TestSeedCompetencies_Idempotent memastikan seeder aman dijalankan ulang
// (UUID deterministik dari nama → semua record di-skip, bukan duplikat).
func TestSeedCompetencies_Idempotent(t *testing.T) {
	db := newCompetencyTestDB(t)

	inserted1, _, err := seedCompetencies(db)
	if err != nil {
		t.Fatalf("first seed failed: %v", err)
	}
	if inserted1 != 165 {
		t.Fatalf("expected 165 inserted on first run, got %d", inserted1)
	}

	inserted2, skipped2, err := seedCompetencies(db)
	if err != nil {
		t.Fatalf("second seed failed: %v", err)
	}
	if inserted2 != 0 {
		t.Fatalf("second run must not insert new records, got %d", inserted2)
	}
	if skipped2 != 165 {
		t.Fatalf("second run must skip all 165 existing records, got %d", skipped2)
	}

	var total int64
	if err := db.Table("competencies").Count(&total).Error; err != nil {
		t.Fatalf("count failed: %v", err)
	}
	if total != 165 {
		t.Fatalf("expected 165 rows after re-seed, got %d", total)
	}
}
