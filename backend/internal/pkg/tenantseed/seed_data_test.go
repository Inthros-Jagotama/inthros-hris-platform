package tenantseed

import (
	"strings"
	"testing"
)

// TestEmbeddedSeedDataFS memastikan file SQL bulk districts/villages benar-benar
// ada di embed tree (seeddata/) dan tidak kosong.
//
// Regresi: sebelumnya seedFromEmbeddedSQL membaca dari migrator.MigrationsFS,
// padahal file besar itu TIDAK pernah ada di embed tree migrator (hanya ada di
// filesystem backend/migrations/seeders/ — folder duplikat yang kini sudah
// dihapus) — menyebabkan provisioning gagal "file does not exist" saat server
// dijalankan dari CWD selain backend/.
// Test ini mencegah kelas bug yang sama terulang.
func TestEmbeddedSeedDataFS(t *testing.T) {
	expected := []struct {
		file string
		min  int // minimal ukuran file (bytes)
	}{
		{"seeddata/002_seed_districts.sql", 100000},
		{"seeddata/003_seed_villages.sql", 1000000},
	}

	for _, e := range expected {
		t.Run(e.file, func(t *testing.T) {
			data, err := seedDataFS.ReadFile(e.file)
			if err != nil {
				t.Fatalf("embedded file %s not found: %v (files MUST be embedded, tidak boleh CWD-dependent)", e.file, err)
			}
			if len(data) < e.min {
				t.Fatalf("embedded file %s terlalu kecil: %d bytes (min %d) — kemungkinan file kosong/tidak lengkap", e.file, len(data), e.min)
			}
			if !strings.Contains(string(data), "INSERT INTO") {
				t.Fatalf("embedded file %s tidak mengandung statement INSERT INTO — format salah", e.file)
			}
		})
	}
}
