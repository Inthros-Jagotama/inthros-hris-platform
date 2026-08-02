package organization

import (
	"context"
	"testing"

	sqlite "github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// newTestRepo membuat repository dengan in-memory SQLite
// (mengikuti pola glebarez/sqlite yang dipakai modul lain).
func newTestRepo(t *testing.T) *Repository {
	t.Helper()
	// DSN unik per test — hindari shared-cache SQLite yang berbagi DB antar test.
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	// Pool 1 koneksi agar DB memory tidak hilang antar query.
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&Organization{}, &OrganizationHistory{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return NewRepository(func(ctx context.Context) (*gorm.DB, error) { return db, nil })
}

func strPtrHelper(s string) *string { return &s }
