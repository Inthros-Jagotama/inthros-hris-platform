// Package tenantseed menyediakan seeder untuk tenant database HRIS.
//
// Berisi dua entry point utama:
//   - SeedTenantMasterData: seeding master data (religion, education, bank,
//     region BPS, salary grade, PTKP, TER, BPJS, dll) ke tenant database.
//   - SeedTenantRBAC: seeding default roles tenant (Admin, Employee) beserta
//     permissions-nya ke tabel RBAC Level 2 (permissions, roles, role_has_permissions).
//
// Keduanya idempotent (menggunakan UUID deterministik + cek keberadaan) dan
// dipanggil saat provisioning tenant baru (CLI handleProvision, command
// seed-data, dan jalur API company.Service.provisionTenant).
package tenantseed

import (
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// =============================================================================
// Helper: Deterministic UUID dari kode
// =============================================================================

// codeToUUID menghasilkan UUID deterministik berdasarkan namespace tabel dan kode.
// Ini memastikan bahwa kode yang sama selalu menghasilkan UUID yang sama,
// sehingga referensi FK antar tabel yang menggunakan UUID bisa konsisten.
func codeToUUID(table, code string) string {
	return uuid.NewSHA1(uuid.NameSpaceDNS, []byte("hris-"+table+"-"+code)).String()
}

// =============================================================================
// Seeder untuk Tenant Database
// =============================================================================

// SeedTenantMasterData melakukan seeding semua master data ke tenant database.
// Mengembalikan aggregate error (errors.Join) jika ada satu atau lebih tabel
// yang gagal di-seed — caller (API provisionTenant / CLI) harus memperlakukan
// error ini sebagai kegagalan provisioning, bukan membiarkan tabel kosong.
func SeedTenantMasterData(tenantDB *gorm.DB, l *zap.Logger) error {
	l.Info("Seeding tenant master data...")

	seeders := []struct {
		table string
		fn    func(*gorm.DB) (int, int, error)
	}{
		{"religions", seedReligions},
		{"educations", seedEducations},
		{"education_majors", seedEducationMajors},
		{"insurances", seedInsurances},
		{"competencies", seedCompetencies},
		{"marital_statuses", seedMaritalStatuses},
		{"relationship_types", seedRelationshipTypes},
		{"employment_statuses", seedEmploymentStatuses},
		{"banks", seedBanks},
		{"nationalities", seedNationalities},
		{"job_families", seedJobFamilies},
		{"provinces", seedProvinces},
		{"regencies", seedRegencies},
		{"districts", seedDistrictsFromSQL},
		{"villages", seedVillagesFromSQL},
		{"salary_grades", seedSalaryGrades},
		{"ptkps", seedPTKPs},
		{"pph21_tax_brackets", seedPPh21TaxBrackets},
		{"ters", seedTERs},
		{"bpjs_settings", seedBPJSSettings},
		{"bpjs_rate_components", seedBPJSRateComponents},
		// Performance Management (KPI) Master Data
		{"performance_perspectives", seedPerformancePerspectives},
		{"performance_ratings", seedPerformanceRatings},
		{"performance_indicator_formulas", seedPerformanceIndicatorFormulas},
	}

	// Kumpulkan error per-tabel. Seeder TIDAK boleh gagal senyap: jika satu tabel
	// saja gagal di-seed (mis. key kolom salah seperti bug religion vs name),
	// provisioning harus gagal "loud" agar caller (API provisionTenant → company
	// suspended, atau CLI → Fatal) segera tahu, bukan membiarkan tabel kosong.
	var errs []error
	for _, s := range seeders {
		start := time.Now()
		inserted, skipped, err := s.fn(tenantDB)
		duration := time.Since(start)
		if err != nil {
			l.Error("Failed to seed table",
				zap.String("table", s.table),
				zap.Error(err),
			)
			errs = append(errs, fmt.Errorf("seed %s failed: %w", s.table, err))
			continue
		}
		l.Info(fmt.Sprintf("  %-25s inserted=%-5d skipped=%-5d %v",
			s.table, inserted, skipped, duration))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// =============================================================================
// Batch Insert Helper
// =============================================================================

func batchInsert(db *gorm.DB, table string, records []map[string]interface{}, batchSize int) (int, int, error) {
	inserted := 0
	skipped := 0

	for i := 0; i < len(records); i += batchSize {
		end := i + batchSize
		if end > len(records) {
			end = len(records)
		}
		batch := records[i:end]

		for _, record := range batch {
			var count int64
			db.Table(table).Where("id = ?", record["id"]).Count(&count)
			if count > 0 {
				skipped++
				continue
			}

			if err := db.Table(table).Create(record).Error; err != nil {
				return inserted, skipped, fmt.Errorf("insert %s failed: %w", table, err)
			}
			inserted++
		}
	}

	return inserted, skipped, nil
}

// =============================================================================
// Data Definitions
// =============================================================================

// ── Religions ──
// Schema (migration 001_master_data): id, code, name, sort_order, created_by, updated_by, timestamps, deleted_at
func seedReligions(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		{"id": codeToUUID("religion", "ISL"), "code": "ISL", "name": "Islam", "sort_order": 1, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("religion", "KTL"), "code": "KTL", "name": "Kristen", "sort_order": 2, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("religion", "KTK"), "code": "KTK", "name": "Katolik", "sort_order": 3, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("religion", "HND"), "code": "HND", "name": "Hindu", "sort_order": 4, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("religion", "BDH"), "code": "BDH", "name": "Budha", "sort_order": 5, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("religion", "KHC"), "code": "KHC", "name": "Konghucu", "sort_order": 6, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("religion", "LNY"), "code": "LNY", "name": "Lainnya", "sort_order": 7, "created_at": time.Now(), "updated_at": time.Now()},
	}
	return batchInsert(db, "religions", data, 50)
}

// ── Educations ──
// Schema (migration 001_master_data): id, code, name, sort_order, created_by, updated_by, timestamps, deleted_at
func seedEducations(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		{"id": codeToUUID("education", "TK"), "code": "TK", "name": "Taman Kanak-Kanak", "sort_order": 1, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education", "SD"), "code": "SD", "name": "Sekolah Dasar", "sort_order": 2, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education", "SMP"), "code": "SMP", "name": "Sekolah Menengah Pertama", "sort_order": 3, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education", "SMA"), "code": "SMA", "name": "Sekolah Menengah Atas", "sort_order": 4, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education", "SMK"), "code": "SMK", "name": "Sekolah Menengah Kejuruan", "sort_order": 5, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education", "D1"), "code": "D1", "name": "Diploma 1", "sort_order": 6, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education", "D2"), "code": "D2", "name": "Diploma 2", "sort_order": 7, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education", "D3"), "code": "D3", "name": "Diploma 3", "sort_order": 8, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education", "D4"), "code": "D4", "name": "Diploma 4", "sort_order": 9, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education", "S1"), "code": "S1", "name": "Strata 1", "sort_order": 10, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education", "S2"), "code": "S2", "name": "Strata 2", "sort_order": 11, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education", "S3"), "code": "S3", "name": "Strata 3", "sort_order": 12, "created_at": time.Now(), "updated_at": time.Now()},
	}
	return batchInsert(db, "educations", data, 50)
}

// ── Education Majors (jurusan pendidikan) ──
// Schema (migration 024_education_majors): id, code, name, sort_order, timestamps, deleted_at
func seedEducationMajors(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		{"id": codeToUUID("education_major", "001"), "code": "001", "name": "Teknik Informatika", "sort_order": 1, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "002"), "code": "002", "name": "Sistem Informasi", "sort_order": 2, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "003"), "code": "003", "name": "Teknik Komputer", "sort_order": 3, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "004"), "code": "004", "name": "Akuntansi", "sort_order": 4, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "005"), "code": "005", "name": "Manajemen", "sort_order": 5, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "006"), "code": "006", "name": "Ekonomi Pembangunan", "sort_order": 6, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "007"), "code": "007", "name": "Teknik Mesin", "sort_order": 7, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "008"), "code": "008", "name": "Teknik Elektro", "sort_order": 8, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "009"), "code": "009", "name": "Teknik Sipil", "sort_order": 9, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "010"), "code": "010", "name": "Arsitektur", "sort_order": 10, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "011"), "code": "011", "name": "Ilmu Hukum", "sort_order": 11, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "012"), "code": "012", "name": "Psikologi", "sort_order": 12, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "013"), "code": "013", "name": "Ilmu Komunikasi", "sort_order": 13, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "014"), "code": "014", "name": "Kedokteran", "sort_order": 14, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "015"), "code": "015", "name": "Keperawatan", "sort_order": 15, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "016"), "code": "016", "name": "Farmasi", "sort_order": 16, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "017"), "code": "017", "name": "Ilmu Pendidikan", "sort_order": 17, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "018"), "code": "018", "name": "Sastra Inggris", "sort_order": 18, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "019"), "code": "019", "name": "Pariwisata", "sort_order": 19, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("education_major", "020"), "code": "020", "name": "Bisnis Digital", "sort_order": 20, "created_at": time.Now(), "updated_at": time.Now()},
	}
	return batchInsert(db, "education_majors", data, 50)
}

// ── Insurances (asuransi) ──
// Schema (migration 021_insurances): id, code, name, sort_order, timestamps, deleted_at
func seedInsurances(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		{"id": codeToUUID("insurance", "01"), "code": "01", "name": "BPJS Kesehatan", "sort_order": 1, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("insurance", "02"), "code": "02", "name": "BPJS Ketenagakerjaan", "sort_order": 2, "created_at": time.Now(), "updated_at": time.Now()},
	}
	return batchInsert(db, "insurances", data, 50)
}

// ── Marital Statuses ──
// Schema (migration 001_master_data): id, code, name, sort_order, created_by, updated_by, timestamps, deleted_at
func seedMaritalStatuses(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		{"id": codeToUUID("marital_status", "BLM"), "code": "BLM", "name": "Belum Kawin", "sort_order": 1, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("marital_status", "KWN"), "code": "KWN", "name": "Kawin", "sort_order": 2, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("marital_status", "CRH"), "code": "CRH", "name": "Cerai Hidup", "sort_order": 3, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("marital_status", "CRM"), "code": "CRM", "name": "Cerai Mati", "sort_order": 4, "created_at": time.Now(), "updated_at": time.Now()},
	}
	return batchInsert(db, "marital_statuses", data, 50)
}

// ── Relationship Types ──
// Schema (migration 001_master_data): id, code, name, sort_order, created_by, updated_by, timestamps, deleted_at
func seedRelationshipTypes(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		{"id": codeToUUID("relationship", "SMI"), "code": "SMI", "name": "Suami", "sort_order": 1, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("relationship", "IST"), "code": "IST", "name": "Istri", "sort_order": 2, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("relationship", "ANC"), "code": "ANC", "name": "Anak", "sort_order": 3, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("relationship", "AYK"), "code": "AYK", "name": "Anak Kandung", "sort_order": 4, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("relationship", "OTA"), "code": "OTA", "name": "Orang Tua", "sort_order": 5, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("relationship", "MRT"), "code": "MRT", "name": "Mertua", "sort_order": 6, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("relationship", "SDR"), "code": "SDR", "name": "Saudara", "sort_order": 7, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("relationship", "KPN"), "code": "KPN", "name": "Keponakan", "sort_order": 8, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("relationship", "CUC"), "code": "CUC", "name": "Cucu", "sort_order": 9, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("relationship", "LNY"), "code": "LNY", "name": "Lainnya", "sort_order": 10, "created_at": time.Now(), "updated_at": time.Now()},
	}
	return batchInsert(db, "relationship_types", data, 50)
}

// ── Employment Statuses ──
// Schema: id(char(36) PK), code, name, sort_order, has_duration, duration, duration_type, created_by, updated_by, timestamps
func seedEmploymentStatuses(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		{"id": codeToUUID("employment_status", "TTP"), "code": "TTP", "name": "Tetap", "sort_order": 1, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("employment_status", "KTR"), "code": "KTR", "name": "Kontrak", "sort_order": 2, "has_duration": 1, "duration": 1, "duration_type": "years", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("employment_status", "MGD"), "code": "MGD", "name": "Magang", "sort_order": 3, "has_duration": 1, "duration": 6, "duration_type": "months", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("employment_status", "HRL"), "code": "HRL", "name": "Harian Lepas", "sort_order": 4, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("employment_status", "OUT"), "code": "OUT", "name": "Outsourcing", "sort_order": 5, "has_duration": 1, "duration": 1, "duration_type": "years", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("employment_status", "PDJ"), "code": "PDJ", "name": "Pensiunan", "sort_order": 6, "created_at": time.Now(), "updated_at": time.Now()},
	}
	return batchInsert(db, "employment_statuses", data, 50)
}

// ── Banks ──
// Schema: id(char(36) PK), code, name, sort_order, timestamps
func seedBanks(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		{"id": codeToUUID("bank", "002"), "code": "002", "name": "Bank Mandiri", "sort_order": 1, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "008"), "code": "008", "name": "Bank Central Asia (BCA)", "sort_order": 2, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "009"), "code": "009", "name": "Bank Negara Indonesia (BNI)", "sort_order": 3, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "013"), "code": "013", "name": "Bank Rakyat Indonesia (BRI)", "sort_order": 4, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "022"), "code": "022", "name": "Bank CIMB Niaga", "sort_order": 5, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "028"), "code": "028", "name": "Bank OCBC NISP", "sort_order": 6, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "031"), "code": "031", "name": "Bank Danamon", "sort_order": 7, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "041"), "code": "041", "name": "Bank Pan Indonesia (Panin)", "sort_order": 8, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "011"), "code": "011", "name": "Bank Danamon", "sort_order": 9, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "016"), "code": "016", "name": "Bank Permata", "sort_order": 10, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "019"), "code": "019", "name": "Bank Panin Dubai Syariah", "sort_order": 11, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "023"), "code": "023", "name": "Bank UOB Indonesia", "sort_order": 12, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "042"), "code": "042", "name": "Bank Maybank Indonesia", "sort_order": 13, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "087"), "code": "087", "name": "Bank Syariah Indonesia (BSI)", "sort_order": 14, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "147"), "code": "147", "name": "Bank Muamalat Indonesia", "sort_order": 15, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "200"), "code": "200", "name": "Bank Tabungan Negara (BTN)", "sort_order": 16, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "213"), "code": "213", "name": "Bank BTPN", "sort_order": 17, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "451"), "code": "451", "name": "Bank Syariah Mandiri", "sort_order": 18, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "484"), "code": "484", "name": "Bank Mega", "sort_order": 19, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "494"), "code": "494", "name": "Bank Raya", "sort_order": 20, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "506"), "code": "506", "name": "Bank Mayora", "sort_order": 21, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "523"), "code": "523", "name": "Bank DKI", "sort_order": 22, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "526"), "code": "526", "name": "Bank BJB", "sort_order": 23, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "536"), "code": "536", "name": "Bank BPD Bali", "sort_order": 24, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("bank", "542"), "code": "542", "name": "Bank Jatim", "sort_order": 25, "created_at": time.Now(), "updated_at": time.Now()},
	}
	return batchInsert(db, "banks", data, 50)
}
// ── Job Families ──
// Schema: id(char(36) PK), code, name, description(text NULL), sort_order, timestamps
func seedJobFamilies(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		{"id": codeToUUID("job_family", "FIN"), "code": "FIN", "name": "Finance & Accounting", "description": "Keuangan, akuntansi, perpajakan", "sort_order": 1, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("job_family", "HR"), "code": "HR", "name": "Human Resources", "description": "SDM, pengembangan organisasi, rekrutmen", "sort_order": 2, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("job_family", "IT"), "code": "IT", "name": "Information Technology", "description": "Teknologi informasi, pengembangan software, infrastruktur", "sort_order": 3, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("job_family", "MKT"), "code": "MKT", "name": "Marketing", "description": "Pemasaran, branding, komunikasi", "sort_order": 4, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("job_family", "SLS"), "code": "SLS", "name": "Sales", "description": "Penjualan, bisnis development", "sort_order": 5, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("job_family", "OPS"), "code": "OPS", "name": "Operations", "description": "Operasional, logistik, supply chain", "sort_order": 6, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("job_family", "LEG"), "code": "LEG", "name": "Legal", "description": "Hukum, kepatuhan, regulasi", "sort_order": 7, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("job_family", "GA"), "code": "GA", "name": "General Affairs", "description": "Umum, administrasi, fasilitas", "sort_order": 8, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("job_family", "RND"), "code": "RND", "name": "Research & Development", "description": "Penelitian dan pengembangan produk", "sort_order": 9, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("job_family", "PROD"), "code": "PROD", "name": "Production", "description": "Produksi, manufaktur", "sort_order": 10, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("job_family", "QC"), "code": "QC", "name": "Quality Control", "description": "Kontrol kualitas, jaminan mutu", "sort_order": 11, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("job_family", "PUR"), "code": "PUR", "name": "Procurement", "description": "Pengadaan, pembelian", "sort_order": 12, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("job_family", "CS"), "code": "CS", "name": "Customer Service", "description": "Layanan pelanggan", "sort_order": 13, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("job_family", "INT"), "code": "INT", "name": "Internal Audit", "description": "Audit internal", "sort_order": 14, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("job_family", "EXC"), "code": "EXC", "name": "Executive", "description": "Direksi, komisaris, manajemen puncak", "sort_order": 15, "created_at": time.Now(), "updated_at": time.Now()},
	}
	return batchInsert(db, "job_families", data, 50)
}
// ── Provinces ──
// Schema: id(char(2) PK - kode Kemendagri), code(varchar(10)), name(varchar(100)), timestamps, deleted_at
func seedProvinces(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		{"id": "11", "code": "11", "name": "Aceh", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "12", "code": "12", "name": "Sumatera Utara", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "13", "code": "13", "name": "Sumatera Barat", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "14", "code": "14", "name": "Riau", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "15", "code": "15", "name": "Jambi", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "16", "code": "16", "name": "Sumatera Selatan", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "17", "code": "17", "name": "Bengkulu", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "18", "code": "18", "name": "Lampung", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "19", "code": "19", "name": "Kepulauan Bangka Belitung", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "21", "code": "21", "name": "Kepulauan Riau", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "31", "code": "31", "name": "DKI Jakarta", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "32", "code": "32", "name": "Jawa Barat", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "33", "code": "33", "name": "Jawa Tengah", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "34", "code": "34", "name": "DI Yogyakarta", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "35", "code": "35", "name": "Jawa Timur", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "36", "code": "36", "name": "Banten", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "51", "code": "51", "name": "Bali", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "52", "code": "52", "name": "Nusa Tenggara Barat", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "53", "code": "53", "name": "Nusa Tenggara Timur", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "61", "code": "61", "name": "Kalimantan Barat", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "62", "code": "62", "name": "Kalimantan Tengah", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "63", "code": "63", "name": "Kalimantan Selatan", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "64", "code": "64", "name": "Kalimantan Timur", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "65", "code": "65", "name": "Kalimantan Utara", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "71", "code": "71", "name": "Sulawesi Utara", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "72", "code": "72", "name": "Sulawesi Tengah", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "73", "code": "73", "name": "Sulawesi Selatan", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "74", "code": "74", "name": "Sulawesi Tenggara", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "75", "code": "75", "name": "Gorontalo", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "76", "code": "76", "name": "Sulawesi Barat", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "81", "code": "81", "name": "Maluku", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "82", "code": "82", "name": "Maluku Utara", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "84", "code": "84", "name": "Papua Barat", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "85", "code": "85", "name": "Papua Barat Daya", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "91", "code": "91", "name": "Papua", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "92", "code": "92", "name": "Papua Pegunungan", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "93", "code": "93", "name": "Papua Selatan", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "94", "code": "94", "name": "Papua Tengah", "created_at": time.Now(), "updated_at": time.Now()},
	}
	return batchInsert(db, "provinces", data, 50)
}

// ── Salary Grades ──
// Schema (migration 001_master_data): id, code, name, description, min_amount, max_amount, sort_order, timestamps
func seedSalaryGrades(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		{"id": codeToUUID("salary_grade", "IA"), "code": "IA", "name": "Golongan I-A", "description": "Operator/Staf Pemula", "min_amount": 1500000, "max_amount": 2500000, "sort_order": 1, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("salary_grade", "IB"), "code": "IB", "name": "Golongan I-B", "description": "Operator/Staf Dasar", "min_amount": 2500000, "max_amount": 3500000, "sort_order": 2, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("salary_grade", "IC"), "code": "IC", "name": "Golongan I-C", "description": "Operator/Staf Madya", "min_amount": 3500000, "max_amount": 4500000, "sort_order": 3, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("salary_grade", "IIA"), "code": "IIA", "name": "Golongan II-A", "description": "Staf Pelaksana Pemula", "min_amount": 4500000, "max_amount": 6000000, "sort_order": 4, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("salary_grade", "IIB"), "code": "IIB", "name": "Golongan II-B", "description": "Staf Pelaksana", "min_amount": 6000000, "max_amount": 8000000, "sort_order": 5, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("salary_grade", "IIC"), "code": "IIC", "name": "Golongan II-C", "description": "Staf Pelaksana Senior", "min_amount": 8000000, "max_amount": 10000000, "sort_order": 6, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("salary_grade", "IIIA"), "code": "IIIA", "name": "Golongan III-A", "description": "Supervisor/Pengawas Pemula", "min_amount": 10000000, "max_amount": 15000000, "sort_order": 7, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("salary_grade", "IIIB"), "code": "IIIB", "name": "Golongan III-B", "description": "Supervisor/Pengawas", "min_amount": 15000000, "max_amount": 20000000, "sort_order": 8, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("salary_grade", "IIIC"), "code": "IIIC", "name": "Golongan III-C", "description": "Supervisor/Pengawas Senior", "min_amount": 20000000, "max_amount": 30000000, "sort_order": 9, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("salary_grade", "IVA"), "code": "IVA", "name": "Golongan IV-A", "description": "Manager Pemula", "min_amount": 30000000, "max_amount": 50000000, "sort_order": 10, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("salary_grade", "IVB"), "code": "IVB", "name": "Golongan IV-B", "description": "Manager", "min_amount": 50000000, "max_amount": 100000000, "sort_order": 11, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("salary_grade", "IVC"), "code": "IVC", "name": "Golongan IV-C", "description": "Senior Manager/Direktur", "min_amount": 100000000, "max_amount": 200000000, "sort_order": 12, "created_at": time.Now(), "updated_at": time.Now()},
	}
	return batchInsert(db, "salary_grades", data, 50)
}

// ── PTKPs (Penghasilan Tidak Kena Pajak) ──
// Schema (migration 001_master_data): id, name, ptkp(BIGINT), `group`(CHAR(1)), created_by, updated_by, timestamps
func seedPTKPs(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		// Kategori TER sesuai aturan resmi PER-2/PJ/2024:
		// A = TK/0, TK/1, K/0 · B = TK/2, TK/3, K/1, K/2 · C = K/3 + K/I/*
		{"id": codeToUUID("ptkp", "TK0"), "name": "Tidak Kawin (TK/0)", "code": "TK0", "ptkp": 54000000, "group": "A", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ptkp", "TK1"), "name": "Tidak Kawin 1 Tanggungan (TK/1)", "code": "TK1", "ptkp": 58500000, "group": "A", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ptkp", "K0"), "name": "Kawin (K/0)", "code": "K0", "ptkp": 58500000, "group": "A", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ptkp", "TK2"), "name": "Tidak Kawin 2 Tanggungan (TK/2)", "code": "TK2", "ptkp": 63000000, "group": "B", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ptkp", "K1"), "name": "Kawin 1 Tanggungan (K/1)", "code": "K1", "ptkp": 63000000, "group": "B", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ptkp", "TK3"), "name": "Tidak Kawin 3 Tanggungan (TK/3)", "code": "TK3", "ptkp": 67500000, "group": "B", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ptkp", "K2"), "name": "Kawin 2 Tanggungan (K/2)", "code": "K2", "ptkp": 67500000, "group": "B", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ptkp", "K3"), "name": "Kawin 3 Tanggungan (K/3)", "code": "K3", "ptkp": 72000000, "group": "C", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ptkp", "KI0"), "name": "Kawin Penghasilan Istri Digabung (K/I/0)", "code": "KI0", "ptkp": 112500000, "group": "C", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ptkp", "KI1"), "name": "Kawin Penghasilan Istri Digabung 1 Tanggungan (K/I/1)", "code": "KI1", "ptkp": 117000000, "group": "C", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ptkp", "KI2"), "name": "Kawin Penghasilan Istri Digabung 2 Tanggungan (K/I/2)", "code": "KI2", "ptkp": 121500000, "group": "C", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ptkp", "KI3"), "name": "Kawin Penghasilan Istri Digabung 3 Tanggungan (K/I/3)", "code": "KI3", "ptkp": 126000000, "group": "C", "created_at": time.Now(), "updated_at": time.Now()},
	}
	return batchInsert(db, "ptkps", data, 50)
}

// ── PPh21 Tax Brackets ──
// Schema (migration 006_payroll_structure): id, bracket_order, lower_bound, upper_bound, rate_percent,
//   effective_start_date(DATE), effective_end_date(DATE), status(ENUM), created_by, updated_by, timestamps
func seedPPh21TaxBrackets(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		{"id": codeToUUID("pph21_bracket", "1"), "bracket_order": 1, "lower_bound": 0, "upper_bound": 60000000, "rate_percent": 5.0, "effective_start_date": "2024-01-01", "status": "ACTIVE", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("pph21_bracket", "2"), "bracket_order": 2, "lower_bound": 60000000, "upper_bound": 250000000, "rate_percent": 15.0, "effective_start_date": "2024-01-01", "status": "ACTIVE", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("pph21_bracket", "3"), "bracket_order": 3, "lower_bound": 250000000, "upper_bound": 500000000, "rate_percent": 25.0, "effective_start_date": "2024-01-01", "status": "ACTIVE", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("pph21_bracket", "4"), "bracket_order": 4, "lower_bound": 500000000, "upper_bound": 5000000000, "rate_percent": 30.0, "effective_start_date": "2024-01-01", "status": "ACTIVE", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("pph21_bracket", "5"), "bracket_order": 5, "lower_bound": 5000000000, "upper_bound": nil, "rate_percent": 35.0, "effective_start_date": "2024-01-01", "status": "ACTIVE", "created_at": time.Now(), "updated_at": time.Now()},
	}
	return batchInsert(db, "pph21_tax_brackets", data, 50)
}

// ── Tarif Efektif Rata-rata (TER) ──
// Schema (migration 001_master_data): id, `group`(CHAR(1)), bruto_min(BIGINT), bruto_max(BIGINT), rate(DECIMAL 10,2),
//   created_by, updated_by, timestamps, deleted_at
func seedTERs(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		// Group A: TK/0, TK/1, K/0 (bruto monthly) - 44 brackets
		{"id": codeToUUID("ter", "A1"), "group": "A", "bruto_min": 0, "bruto_max": 5400000, "rate": 0.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A2"), "group": "A", "bruto_min": 5400000, "bruto_max": 5650000, "rate": 0.25, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A3"), "group": "A", "bruto_min": 5650000, "bruto_max": 5950000, "rate": 0.50, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A4"), "group": "A", "bruto_min": 5950000, "bruto_max": 6300000, "rate": 0.75, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A5"), "group": "A", "bruto_min": 6300000, "bruto_max": 6750000, "rate": 1.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A6"), "group": "A", "bruto_min": 6750000, "bruto_max": 7500000, "rate": 1.25, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A7"), "group": "A", "bruto_min": 7500000, "bruto_max": 8550000, "rate": 1.50, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A8"), "group": "A", "bruto_min": 8550000, "bruto_max": 9650000, "rate": 1.75, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A9"), "group": "A", "bruto_min": 9650000, "bruto_max": 10050000, "rate": 2.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A10"), "group": "A", "bruto_min": 10050000, "bruto_max": 10350000, "rate": 2.25, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A11"), "group": "A", "bruto_min": 10350000, "bruto_max": 10700000, "rate": 2.50, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A12"), "group": "A", "bruto_min": 10700000, "bruto_max": 11050000, "rate": 3.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A13"), "group": "A", "bruto_min": 11050000, "bruto_max": 11600000, "rate": 3.50, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A14"), "group": "A", "bruto_min": 11600000, "bruto_max": 12500000, "rate": 4.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A15"), "group": "A", "bruto_min": 12500000, "bruto_max": 13750000, "rate": 5.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A16"), "group": "A", "bruto_min": 13750000, "bruto_max": 15100000, "rate": 6.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A17"), "group": "A", "bruto_min": 15100000, "bruto_max": 16950000, "rate": 7.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A18"), "group": "A", "bruto_min": 16950000, "bruto_max": 19750000, "rate": 8.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A19"), "group": "A", "bruto_min": 19750000, "bruto_max": 24150000, "rate": 9.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A20"), "group": "A", "bruto_min": 24150000, "bruto_max": 26450000, "rate": 10.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A21"), "group": "A", "bruto_min": 26450000, "bruto_max": 28000000, "rate": 11.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A22"), "group": "A", "bruto_min": 28000000, "bruto_max": 30050000, "rate": 12.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A23"), "group": "A", "bruto_min": 30050000, "bruto_max": 32400000, "rate": 13.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A24"), "group": "A", "bruto_min": 32400000, "bruto_max": 35400000, "rate": 14.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A25"), "group": "A", "bruto_min": 35400000, "bruto_max": 39100000, "rate": 15.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A26"), "group": "A", "bruto_min": 39100000, "bruto_max": 43850000, "rate": 16.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A27"), "group": "A", "bruto_min": 43850000, "bruto_max": 47800000, "rate": 17.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A28"), "group": "A", "bruto_min": 47800000, "bruto_max": 51400000, "rate": 18.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A29"), "group": "A", "bruto_min": 51400000, "bruto_max": 56300000, "rate": 19.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A30"), "group": "A", "bruto_min": 56300000, "bruto_max": 62200000, "rate": 20.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A31"), "group": "A", "bruto_min": 62200000, "bruto_max": 68600000, "rate": 21.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A32"), "group": "A", "bruto_min": 68600000, "bruto_max": 77500000, "rate": 22.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A33"), "group": "A", "bruto_min": 77500000, "bruto_max": 89000000, "rate": 23.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A34"), "group": "A", "bruto_min": 89000000, "bruto_max": 103000000, "rate": 24.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A35"), "group": "A", "bruto_min": 103000000, "bruto_max": 125000000, "rate": 25.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A36"), "group": "A", "bruto_min": 125000000, "bruto_max": 157000000, "rate": 26.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A37"), "group": "A", "bruto_min": 157000000, "bruto_max": 206000000, "rate": 27.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A38"), "group": "A", "bruto_min": 206000000, "bruto_max": 337000000, "rate": 28.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A39"), "group": "A", "bruto_min": 337000000, "bruto_max": 454000000, "rate": 29.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A40"), "group": "A", "bruto_min": 454000000, "bruto_max": 550000000, "rate": 30.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A41"), "group": "A", "bruto_min": 550000000, "bruto_max": 695000000, "rate": 31.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A42"), "group": "A", "bruto_min": 695000000, "bruto_max": 910000000, "rate": 32.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A43"), "group": "A", "bruto_min": 910000000, "bruto_max": 1400000000, "rate": 33.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "A44"), "group": "A", "bruto_min": 1400000000, "bruto_max": nil, "rate": 34.00, "created_at": time.Now(), "updated_at": time.Now()},

		// Group B: TK/2, TK/3, K/1, K/2 (bruto monthly) - 40 brackets
		{"id": codeToUUID("ter", "B1"), "group": "B", "bruto_min": 0, "bruto_max": 6200000, "rate": 0.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B2"), "group": "B", "bruto_min": 6200000, "bruto_max": 6500000, "rate": 0.25, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B3"), "group": "B", "bruto_min": 6500000, "bruto_max": 6850000, "rate": 0.50, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B4"), "group": "B", "bruto_min": 6850000, "bruto_max": 7300000, "rate": 0.75, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B5"), "group": "B", "bruto_min": 7300000, "bruto_max": 9200000, "rate": 1.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B6"), "group": "B", "bruto_min": 9200000, "bruto_max": 10750000, "rate": 1.50, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B7"), "group": "B", "bruto_min": 10750000, "bruto_max": 11250000, "rate": 2.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B8"), "group": "B", "bruto_min": 11250000, "bruto_max": 11600000, "rate": 2.50, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B9"), "group": "B", "bruto_min": 11600000, "bruto_max": 12600000, "rate": 3.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B10"), "group": "B", "bruto_min": 12600000, "bruto_max": 13600000, "rate": 4.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B11"), "group": "B", "bruto_min": 13600000, "bruto_max": 14950000, "rate": 5.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B12"), "group": "B", "bruto_min": 14950000, "bruto_max": 16400000, "rate": 6.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B13"), "group": "B", "bruto_min": 16400000, "bruto_max": 18450000, "rate": 7.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B14"), "group": "B", "bruto_min": 18450000, "bruto_max": 21850000, "rate": 8.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B15"), "group": "B", "bruto_min": 21850000, "bruto_max": 26000000, "rate": 9.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B16"), "group": "B", "bruto_min": 26000000, "bruto_max": 27700000, "rate": 10.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B17"), "group": "B", "bruto_min": 27700000, "bruto_max": 29350000, "rate": 11.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B18"), "group": "B", "bruto_min": 29350000, "bruto_max": 31450000, "rate": 12.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B19"), "group": "B", "bruto_min": 31450000, "bruto_max": 33950000, "rate": 13.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B20"), "group": "B", "bruto_min": 33950000, "bruto_max": 37100000, "rate": 14.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B21"), "group": "B", "bruto_min": 37100000, "bruto_max": 41100000, "rate": 15.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B22"), "group": "B", "bruto_min": 41100000, "bruto_max": 45800000, "rate": 16.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B23"), "group": "B", "bruto_min": 45800000, "bruto_max": 49500000, "rate": 17.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B24"), "group": "B", "bruto_min": 49500000, "bruto_max": 53800000, "rate": 18.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B25"), "group": "B", "bruto_min": 53800000, "bruto_max": 58500000, "rate": 19.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B26"), "group": "B", "bruto_min": 58500000, "bruto_max": 64000000, "rate": 20.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B27"), "group": "B", "bruto_min": 64000000, "bruto_max": 71000000, "rate": 21.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B28"), "group": "B", "bruto_min": 71000000, "bruto_max": 80000000, "rate": 22.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B29"), "group": "B", "bruto_min": 80000000, "bruto_max": 93000000, "rate": 23.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B30"), "group": "B", "bruto_min": 93000000, "bruto_max": 109000000, "rate": 24.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B31"), "group": "B", "bruto_min": 109000000, "bruto_max": 129000000, "rate": 25.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B32"), "group": "B", "bruto_min": 129000000, "bruto_max": 163000000, "rate": 26.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B33"), "group": "B", "bruto_min": 163000000, "bruto_max": 211000000, "rate": 27.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B34"), "group": "B", "bruto_min": 211000000, "bruto_max": 374000000, "rate": 28.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B35"), "group": "B", "bruto_min": 374000000, "bruto_max": 459000000, "rate": 29.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B36"), "group": "B", "bruto_min": 459000000, "bruto_max": 555000000, "rate": 30.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B37"), "group": "B", "bruto_min": 555000000, "bruto_max": 704000000, "rate": 31.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B38"), "group": "B", "bruto_min": 704000000, "bruto_max": 957000000, "rate": 32.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B39"), "group": "B", "bruto_min": 957000000, "bruto_max": 1405000000, "rate": 33.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "B40"), "group": "B", "bruto_min": 1405000000, "bruto_max": nil, "rate": 34.00, "created_at": time.Now(), "updated_at": time.Now()},

		// Group C: K/3, K/I/0, K/I/1, K/I/2, K/I/3 (bruto monthly) - 41 brackets
		{"id": codeToUUID("ter", "C0"), "group": "C", "bruto_min": 0, "bruto_max": 6600000, "rate": 0.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C1"), "group": "C", "bruto_min": 6600000, "bruto_max": 6950000, "rate": 0.25, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C2"), "group": "C", "bruto_min": 6950000, "bruto_max": 7350000, "rate": 0.50, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C3"), "group": "C", "bruto_min": 7350000, "bruto_max": 7800000, "rate": 0.75, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C4"), "group": "C", "bruto_min": 7800000, "bruto_max": 8850000, "rate": 1.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C5"), "group": "C", "bruto_min": 8850000, "bruto_max": 9800000, "rate": 1.25, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C6"), "group": "C", "bruto_min": 9800000, "bruto_max": 10950000, "rate": 1.50, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C7"), "group": "C", "bruto_min": 10950000, "bruto_max": 11200000, "rate": 1.75, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C8"), "group": "C", "bruto_min": 11200000, "bruto_max": 12050000, "rate": 2.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C9"), "group": "C", "bruto_min": 12050000, "bruto_max": 12950000, "rate": 3.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C10"), "group": "C", "bruto_min": 12950000, "bruto_max": 14150000, "rate": 4.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C11"), "group": "C", "bruto_min": 14150000, "bruto_max": 15550000, "rate": 5.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C12"), "group": "C", "bruto_min": 15550000, "bruto_max": 17050000, "rate": 6.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C13"), "group": "C", "bruto_min": 17050000, "bruto_max": 19500000, "rate": 7.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C14"), "group": "C", "bruto_min": 19500000, "bruto_max": 22700000, "rate": 8.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C15"), "group": "C", "bruto_min": 22700000, "bruto_max": 26600000, "rate": 9.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C16"), "group": "C", "bruto_min": 26600000, "bruto_max": 28100000, "rate": 10.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C17"), "group": "C", "bruto_min": 28100000, "bruto_max": 30100000, "rate": 11.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C18"), "group": "C", "bruto_min": 30100000, "bruto_max": 32600000, "rate": 12.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C19"), "group": "C", "bruto_min": 32600000, "bruto_max": 35400000, "rate": 13.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C20"), "group": "C", "bruto_min": 35400000, "bruto_max": 38900000, "rate": 14.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C21"), "group": "C", "bruto_min": 38900000, "bruto_max": 43000000, "rate": 15.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C22"), "group": "C", "bruto_min": 43000000, "bruto_max": 47400000, "rate": 16.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C23"), "group": "C", "bruto_min": 47400000, "bruto_max": 51200000, "rate": 17.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C24"), "group": "C", "bruto_min": 51200000, "bruto_max": 55800000, "rate": 18.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C25"), "group": "C", "bruto_min": 55800000, "bruto_max": 60400000, "rate": 19.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C26"), "group": "C", "bruto_min": 60400000, "bruto_max": 66700000, "rate": 20.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C27"), "group": "C", "bruto_min": 66700000, "bruto_max": 74500000, "rate": 21.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C28"), "group": "C", "bruto_min": 74500000, "bruto_max": 83200000, "rate": 22.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C29"), "group": "C", "bruto_min": 83200000, "bruto_max": 95600000, "rate": 23.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C30"), "group": "C", "bruto_min": 95600000, "bruto_max": 110000000, "rate": 24.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C31"), "group": "C", "bruto_min": 110000000, "bruto_max": 134000000, "rate": 25.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C32"), "group": "C", "bruto_min": 134000000, "bruto_max": 169000000, "rate": 26.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C33"), "group": "C", "bruto_min": 169000000, "bruto_max": 221000000, "rate": 27.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C34"), "group": "C", "bruto_min": 221000000, "bruto_max": 390000000, "rate": 28.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C35"), "group": "C", "bruto_min": 390000000, "bruto_max": 463000000, "rate": 29.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C36"), "group": "C", "bruto_min": 463000000, "bruto_max": 561000000, "rate": 30.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C37"), "group": "C", "bruto_min": 561000000, "bruto_max": 709000000, "rate": 31.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C38"), "group": "C", "bruto_min": 709000000, "bruto_max": 965000000, "rate": 32.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C39"), "group": "C", "bruto_min": 965000000, "bruto_max": 1419000000, "rate": 33.00, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("ter", "C40"), "group": "C", "bruto_min": 1419000000, "bruto_max": nil, "rate": 34.00, "created_at": time.Now(), "updated_at": time.Now()},
	}
	return batchInsert(db, "ters", data, 50)
}

// ── BPJS Settings ──
// Schema (migration 006_payroll_structure): id, setting_code, setting_name, base_source,
//   health_max_base_amount, pension_max_base_amount, default_jkk_risk_class, rounding_mode,
//   effective_start_date, effective_end_date, status, notes, timestamps
func seedBPJSSettings(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		{
			"id":                     codeToUUID("bpjs_setting", "2024"),
			"setting_code":           "BPJS-2024-DEFAULT",
			"setting_name":           "BPJS Default 2024",
			"base_source":            "BPJS_BASE_COMPONENTS",
			"health_max_base_amount": 12000000,
			"pension_max_base_amount": 10167400,
			"default_jkk_risk_class": "LOW",
			"rounding_mode":          "ROUND",
			"effective_start_date":   "2024-01-01",
			"status":                 "ACTIVE",
			"notes":                  "Default BPJS configuration untuk tahun 2024. Sesuai PP No. 82/2021 (Kesehatan), PP No. 44/2015, PP No. 45/2015, PP No. 46/2015 (Ketenagakerjaan), dan PP No. 37/2021 (JKP).",
			"created_at":             time.Now(),
			"updated_at":             time.Now(),
		},
	}
	return batchInsert(db, "bpjs_settings", data, 50)
}

// ── BPJS Rate Components ──
// Schema (migration 006_payroll_structure): id, bpjs_setting_id, rate_code, rate_name,
//   bpjs_program(ENUM: HEALTH/JHT/JP/JKK/JKM/JKP), paid_by(ENUM: EMPLOYEE/EMPLOYER),
//   salary_component_id(NULL), rate_percent, fixed_amount(NULL), min_base_amount(NULL),
//   max_base_amount(NULL), jkk_risk_class(NULL), is_employee_deduction,
//   is_employer_contribution, generate_to_payroll_item, print_on_payslip,
//   display_order, effective_start_date, effective_end_date, status, timestamps
func seedBPJSRateComponents(db *gorm.DB) (int, int, error) {
	settingID := codeToUUID("bpjs_setting", "2024")
	data := []map[string]interface{}{
		// ── BPJS Kesehatan (HEALTH) ──
		// Employer 4%, Employee 1% — max base Rp 12jt (PP 82/2021)
		{
			"id":                     codeToUUID("bpjs_rate", "HEALTH-EMP"),
			"bpjs_setting_id":        settingID,
			"rate_code":              "BPJS-KES-EMP",
			"rate_name":              "BPJS Kesehatan - Pemberi Kerja",
			"bpjs_program":           "HEALTH",
			"paid_by":                "EMPLOYER",
			"rate_percent":           4.0,
			"max_base_amount":        12000000,
			"is_employee_deduction":   0,
			"is_employer_contribution": 1,
			"generate_to_payroll_item": 1,
			"print_on_payslip":        1,
			"display_order":           1,
			"effective_start_date":    "2024-01-01",
			"status":                  "ACTIVE",
			"created_at":              time.Now(),
			"updated_at":              time.Now(),
		},
		{
			"id":                     codeToUUID("bpjs_rate", "HEALTH-EE"),
			"bpjs_setting_id":        settingID,
			"rate_code":              "BPJS-KES-EE",
			"rate_name":              "BPJS Kesehatan - Pekerja",
			"bpjs_program":           "HEALTH",
			"paid_by":                "EMPLOYEE",
			"rate_percent":           1.0,
			"max_base_amount":        12000000,
			"is_employee_deduction":   1,
			"is_employer_contribution": 0,
			"generate_to_payroll_item": 1,
			"print_on_payslip":        1,
			"display_order":           2,
			"effective_start_date":    "2024-01-01",
			"status":                  "ACTIVE",
			"created_at":              time.Now(),
			"updated_at":              time.Now(),
		},

		// ── JHT (Jaminan Hari Tua) ──
		// Employer 3.7%, Employee 2% (PP 46/2015)
		{
			"id":                     codeToUUID("bpjs_rate", "JHT-EMP"),
			"bpjs_setting_id":        settingID,
			"rate_code":              "BPJS-JHT-EMP",
			"rate_name":              "JHT - Pemberi Kerja",
			"bpjs_program":           "JHT",
			"paid_by":                "EMPLOYER",
			"rate_percent":           3.7,
			"is_employee_deduction":   0,
			"is_employer_contribution": 1,
			"generate_to_payroll_item": 1,
			"print_on_payslip":        1,
			"display_order":           3,
			"effective_start_date":    "2024-01-01",
			"status":                  "ACTIVE",
			"created_at":              time.Now(),
			"updated_at":              time.Now(),
		},
		{
			"id":                     codeToUUID("bpjs_rate", "JHT-EE"),
			"bpjs_setting_id":        settingID,
			"rate_code":              "BPJS-JHT-EE",
			"rate_name":              "JHT - Pekerja",
			"bpjs_program":           "JHT",
			"paid_by":                "EMPLOYEE",
			"rate_percent":           2.0,
			"is_employee_deduction":   1,
			"is_employer_contribution": 0,
			"generate_to_payroll_item": 1,
			"print_on_payslip":        1,
			"display_order":           4,
			"effective_start_date":    "2024-01-01",
			"status":                  "ACTIVE",
			"created_at":              time.Now(),
			"updated_at":              time.Now(),
		},

		// ── JP (Jaminan Pensiun) ──
		// Employer 2%, Employee 1% (PP 45/2015)
		// Max base amount: dikonfigurasi via pension_max_base_amount di setting
		{
			"id":                     codeToUUID("bpjs_rate", "JP-EMP"),
			"bpjs_setting_id":        settingID,
			"rate_code":              "BPJS-JP-EMP",
			"rate_name":              "JP - Pemberi Kerja",
			"bpjs_program":           "JP",
			"paid_by":                "EMPLOYER",
			"rate_percent":           2.0,
			"is_employee_deduction":   0,
			"is_employer_contribution": 1,
			"generate_to_payroll_item": 1,
			"print_on_payslip":        1,
			"display_order":           5,
			"effective_start_date":    "2024-01-01",
			"status":                  "ACTIVE",
			"created_at":              time.Now(),
			"updated_at":              time.Now(),
		},
		{
			"id":                     codeToUUID("bpjs_rate", "JP-EE"),
			"bpjs_setting_id":        settingID,
			"rate_code":              "BPJS-JP-EE",
			"rate_name":              "JP - Pekerja",
			"bpjs_program":           "JP",
			"paid_by":                "EMPLOYEE",
			"rate_percent":           1.0,
			"is_employee_deduction":   1,
			"is_employer_contribution": 0,
			"generate_to_payroll_item": 1,
			"print_on_payslip":        1,
			"display_order":           6,
			"effective_start_date":    "2024-01-01",
			"status":                  "ACTIVE",
			"created_at":              time.Now(),
			"updated_at":              time.Now(),
		},

		// ── JKK (Jaminan Kecelakaan Kerja) ──
		// Employer only: 0.24%–1.74% tergantung kelas risiko (PP 44/2015)
		{
			"id":                     codeToUUID("bpjs_rate", "JKK-VLOW"),
			"bpjs_setting_id":        settingID,
			"rate_code":              "BPJS-JKK-VLOW",
			"rate_name":              "JKK - Risiko Sangat Rendah (0.24%)",
			"bpjs_program":           "JKK",
			"paid_by":                "EMPLOYER",
			"rate_percent":           0.24,
			"jkk_risk_class":         "VERY_LOW",
			"is_employee_deduction":   0,
			"is_employer_contribution": 1,
			"generate_to_payroll_item": 1,
			"print_on_payslip":        1,
			"display_order":           7,
			"effective_start_date":    "2024-01-01",
			"status":                  "ACTIVE",
			"created_at":              time.Now(),
			"updated_at":              time.Now(),
		},
		{
			"id":                     codeToUUID("bpjs_rate", "JKK-LOW"),
			"bpjs_setting_id":        settingID,
			"rate_code":              "BPJS-JKK-LOW",
			"rate_name":              "JKK - Risiko Rendah (0.54%)",
			"bpjs_program":           "JKK",
			"paid_by":                "EMPLOYER",
			"rate_percent":           0.54,
			"jkk_risk_class":         "LOW",
			"is_employee_deduction":   0,
			"is_employer_contribution": 1,
			"generate_to_payroll_item": 1,
			"print_on_payslip":        1,
			"display_order":           8,
			"effective_start_date":    "2024-01-01",
			"status":                  "ACTIVE",
			"created_at":              time.Now(),
			"updated_at":              time.Now(),
		},
		{
			"id":                     codeToUUID("bpjs_rate", "JKK-MED"),
			"bpjs_setting_id":        settingID,
			"rate_code":              "BPJS-JKK-MED",
			"rate_name":              "JKK - Risiko Sedang (0.89%)",
			"bpjs_program":           "JKK",
			"paid_by":                "EMPLOYER",
			"rate_percent":           0.89,
			"jkk_risk_class":         "MEDIUM",
			"is_employee_deduction":   0,
			"is_employer_contribution": 1,
			"generate_to_payroll_item": 1,
			"print_on_payslip":        1,
			"display_order":           9,
			"effective_start_date":    "2024-01-01",
			"status":                  "ACTIVE",
			"created_at":              time.Now(),
			"updated_at":              time.Now(),
		},
		{
			"id":                     codeToUUID("bpjs_rate", "JKK-HIGH"),
			"bpjs_setting_id":        settingID,
			"rate_code":              "BPJS-JKK-HIGH",
			"rate_name":              "JKK - Risiko Tinggi (1.27%)",
			"bpjs_program":           "JKK",
			"paid_by":                "EMPLOYER",
			"rate_percent":           1.27,
			"jkk_risk_class":         "HIGH",
			"is_employee_deduction":   0,
			"is_employer_contribution": 1,
			"generate_to_payroll_item": 1,
			"print_on_payslip":        1,
			"display_order":           10,
			"effective_start_date":    "2024-01-01",
			"status":                  "ACTIVE",
			"created_at":              time.Now(),
			"updated_at":              time.Now(),
		},
		{
			"id":                     codeToUUID("bpjs_rate", "JKK-VHIGH"),
			"bpjs_setting_id":        settingID,
			"rate_code":              "BPJS-JKK-VHIGH",
			"rate_name":              "JKK - Risiko Sangat Tinggi (1.74%)",
			"bpjs_program":           "JKK",
			"paid_by":                "EMPLOYER",
			"rate_percent":           1.74,
			"jkk_risk_class":         "VERY_HIGH",
			"is_employee_deduction":   0,
			"is_employer_contribution": 1,
			"generate_to_payroll_item": 1,
			"print_on_payslip":        1,
			"display_order":           11,
			"effective_start_date":    "2024-01-01",
			"status":                  "ACTIVE",
			"created_at":              time.Now(),
			"updated_at":              time.Now(),
		},

		// ── JKM (Jaminan Kematian) ──
		// Employer only: 0.3% (PP 44/2015)
		{
			"id":                     codeToUUID("bpjs_rate", "JKM"),
			"bpjs_setting_id":        settingID,
			"rate_code":              "BPJS-JKM",
			"rate_name":              "JKM - Jaminan Kematian",
			"bpjs_program":           "JKM",
			"paid_by":                "EMPLOYER",
			"rate_percent":           0.3,
			"is_employee_deduction":   0,
			"is_employer_contribution": 1,
			"generate_to_payroll_item": 1,
			"print_on_payslip":        1,
			"display_order":           12,
			"effective_start_date":    "2024-01-01",
			"status":                  "ACTIVE",
			"created_at":              time.Now(),
			"updated_at":              time.Now(),
		},

		// ── JKP (Jaminan Kehilangan Pekerjaan) ──
		// Employer only: 0.46% (PP 37/2021)
		{
			"id":                     codeToUUID("bpjs_rate", "JKP"),
			"bpjs_setting_id":        settingID,
			"rate_code":              "BPJS-JKP",
			"rate_name":              "JKP - Jaminan Kehilangan Pekerjaan",
			"bpjs_program":           "JKP",
			"paid_by":                "EMPLOYER",
			"rate_percent":           0.46,
			"is_employee_deduction":   0,
			"is_employer_contribution": 1,
			"generate_to_payroll_item": 1,
			"print_on_payslip":        1,
			"display_order":           13,
			"effective_start_date":    "2024-01-01",
			"status":                  "ACTIVE",
			"created_at":              time.Now(),
			"updated_at":              time.Now(),
		},
	}
	return batchInsert(db, "bpjs_rate_components", data, 50)
}

// =============================================================================
// SQL-based Seeders (large datasets)
// =============================================================================

//go:embed seeddata
var seedDataFS embed.FS

// seedFromEmbeddedSQL membaca file SQL besar dari embedded FS (seedDataFS)
// dan mengeksekusinya ke database.
// File SQL districts/villages di-embed ke binary (folder seeddata/), sehingga
// tidak bergantung pada working directory (CWD) — provisioning aman dijalankan
// dari direktori mana pun.
// Disable FK checks sementara karena data master lintas-referensi.
func seedFromEmbeddedSQL(db *gorm.DB, tableName, fileName string) (int, int, error) {
	filePath := "seeddata/" + fileName
	data, err := seedDataFS.ReadFile(filePath)
	if err != nil {
		return 0, 0, fmt.Errorf("read embedded SQL for %s (%s) failed: %w", tableName, filePath, err)
	}

	// Disable FK checks untuk bulk insert
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	sql := string(data)
	if err := db.Exec(sql).Error; err != nil {
		db.Exec("SET FOREIGN_KEY_CHECKS = 1")
		return 0, 0, fmt.Errorf("execute SQL for %s (%s) failed: %w", tableName, filePath, err)
	}

	// Re-enable FK checks
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	return 1, 0, nil
}

// seedDistrictsFromSQL men-seed districts dari file SQL embedded.
func seedDistrictsFromSQL(db *gorm.DB) (int, int, error) {
	return seedFromEmbeddedSQL(db, "districts", "002_seed_districts.sql")
}

// seedVillagesFromSQL men-seed villages dari file SQL embedded.
func seedVillagesFromSQL(db *gorm.DB) (int, int, error) {
	return seedFromEmbeddedSQL(db, "villages", "003_seed_villages.sql")
}

// ── Regencies ──
// Schema: id(char(4) PK - kode Kemendagri), code(varchar(10)), province_id(char(2) FK -> provinces.id), name(varchar(100)), timestamps, deleted_at
func seedRegencies(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		// === ACEH (11) ===
		{"id": "1101", "code": "1101", "name": "KAB. ACEH SELATAN", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1102", "code": "1102", "name": "KAB. ACEH TENGGARA", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1103", "code": "1103", "name": "KAB. ACEH TIMUR", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1104", "code": "1104", "name": "KAB. ACEH TENGAH", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1105", "code": "1105", "name": "KAB. ACEH BARAT", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1106", "code": "1106", "name": "KAB. ACEH BESAR", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1107", "code": "1107", "name": "KAB. PIDIE", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1108", "code": "1108", "name": "KAB. ACEH UTARA", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1109", "code": "1109", "name": "KAB. SIMEULUE", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1110", "code": "1110", "name": "KAB. ACEH SINGKIL", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1111", "code": "1111", "name": "KAB. BIREUEN", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1112", "code": "1112", "name": "KAB. ACEH BARAT DAYA", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1113", "code": "1113", "name": "KAB. GAYO LUES", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1114", "code": "1114", "name": "KAB. ACEH JAYA", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1115", "code": "1115", "name": "KAB. NAGAN RAYA", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1116", "code": "1116", "name": "KAB. ACEH TAMIANG", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1117", "code": "1117", "name": "KAB. BENER MERIAH", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1118", "code": "1118", "name": "KAB. PIDIE JAYA", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1171", "code": "1171", "name": "KOTA BANDA ACEH", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1172", "code": "1172", "name": "KOTA SABANG", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1173", "code": "1173", "name": "KOTA LHOKSEUMAWE", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1174", "code": "1174", "name": "KOTA LANGSA", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1175", "code": "1175", "name": "KOTA SUBULUSSALAM", "province_id": "11", "created_at": time.Now(), "updated_at": time.Now()},

		// === SUMATERA UTARA (12) ===
		{"id": "1201", "code": "1201", "name": "KAB. TAPANULI TENGAH", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1202", "code": "1202", "name": "KAB. TAPANULI UTARA", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1203", "code": "1203", "name": "KAB. TAPANULI SELATAN", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1204", "code": "1204", "name": "KAB. NIAS", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1205", "code": "1205", "name": "KAB. LANGKAT", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1206", "code": "1206", "name": "KAB. KARO", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1207", "code": "1207", "name": "KAB. DELI SERDANG", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1208", "code": "1208", "name": "KAB. SIMALUNGUN", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1209", "code": "1209", "name": "KAB. ASAHAN", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1210", "code": "1210", "name": "KAB. LABUHANBATU", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1211", "code": "1211", "name": "KAB. DAIRI", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1212", "code": "1212", "name": "KAB. TOBA SAMOSIR", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1213", "code": "1213", "name": "KAB. MANDAILING NATAL", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1214", "code": "1214", "name": "KAB. NIAS SELATAN", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1215", "code": "1215", "name": "KAB. PAKPAK BHARAT", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1216", "code": "1216", "name": "KAB. HUMBANG HASUNDUTAN", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1217", "code": "1217", "name": "KAB. SAMOSIR", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1218", "code": "1218", "name": "KAB. SERDANG BEDAGAI", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1219", "code": "1219", "name": "KAB. BATU BARA", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1220", "code": "1220", "name": "KAB. PADANG LAWAS UTARA", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1221", "code": "1221", "name": "KAB. PADANG LAWAS", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1222", "code": "1222", "name": "KAB. LABUHANBATU SELATAN", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1223", "code": "1223", "name": "KAB. LABUHANBATU UTARA", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1224", "code": "1224", "name": "KAB. NIAS UTARA", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1225", "code": "1225", "name": "KAB. NIAS BARAT", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1271", "code": "1271", "name": "KOTA MEDAN", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1272", "code": "1272", "name": "KOTA PEMATANGSIANTAR", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1273", "code": "1273", "name": "KOTA SIBOLGA", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1274", "code": "1274", "name": "KOTA TANJUNG BALAI", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1275", "code": "1275", "name": "KOTA BINJAI", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1276", "code": "1276", "name": "KOTA TEBING TINGGI", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1277", "code": "1277", "name": "KOTA PADANG SIDEMPUAN", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1278", "code": "1278", "name": "KOTA GUNUNGSITOLI", "province_id": "12", "created_at": time.Now(), "updated_at": time.Now()},

		// === SUMATERA BARAT (13) ===
		{"id": "1301", "code": "1301", "name": "KAB. PESISIR SELATAN", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1302", "code": "1302", "name": "KAB. SOLOK", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1303", "code": "1303", "name": "KAB. SIJUNJUNG", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1304", "code": "1304", "name": "KAB. TANAH DATAR", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1305", "code": "1305", "name": "KAB. PADANG PARIAMAN", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1306", "code": "1306", "name": "KAB. AGAM", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1307", "code": "1307", "name": "KAB. LIMA PULUH KOTA", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1308", "code": "1308", "name": "KAB. PASAMAN", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1309", "code": "1309", "name": "KAB. KEPULAUAN MENTAWAI", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1310", "code": "1310", "name": "KAB. DHARMASRAYA", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1311", "code": "1311", "name": "KAB. SOLOK SELATAN", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1312", "code": "1312", "name": "KAB. PASAMAN BARAT", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1371", "code": "1371", "name": "KOTA PADANG", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1372", "code": "1372", "name": "KOTA SOLOK", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1373", "code": "1373", "name": "KOTA SAWAHLUNTO", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1374", "code": "1374", "name": "KOTA PADANG PANJANG", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1375", "code": "1375", "name": "KOTA BUKITTINGGI", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1376", "code": "1376", "name": "KOTA PAYAKUMBUH", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1377", "code": "1377", "name": "KOTA PARIAMAN", "province_id": "13", "created_at": time.Now(), "updated_at": time.Now()},

		// === RIAU (14) ===
		{"id": "1401", "code": "1401", "name": "KAB. KAMPAR", "province_id": "14", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1402", "code": "1402", "name": "KAB. INDRAGIRI HULU", "province_id": "14", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1403", "code": "1403", "name": "KAB. BENGKALIS", "province_id": "14", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1404", "code": "1404", "name": "KAB. INDRAGIRI HILIR", "province_id": "14", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1405", "code": "1405", "name": "KAB. PELALAWAN", "province_id": "14", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1406", "code": "1406", "name": "KAB. ROKAN HULU", "province_id": "14", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1407", "code": "1407", "name": "KAB. ROKAN HILIR", "province_id": "14", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1408", "code": "1408", "name": "KAB. SIAK", "province_id": "14", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1409", "code": "1409", "name": "KAB. KUANTAN SINGINGI", "province_id": "14", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1410", "code": "1410", "name": "KAB. KEPULAUAN MERANTI", "province_id": "14", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1471", "code": "1471", "name": "KOTA PEKANBARU", "province_id": "14", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1472", "code": "1472", "name": "KOTA DUMAI", "province_id": "14", "created_at": time.Now(), "updated_at": time.Now()},

		// === JAMBI (15) ===
		{"id": "1501", "code": "1501", "name": "KAB. KERINCI", "province_id": "15", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1502", "code": "1502", "name": "KAB. MERANGIN", "province_id": "15", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1503", "code": "1503", "name": "KAB. SAROLANGUN", "province_id": "15", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1504", "code": "1504", "name": "KAB. BATANG HARI", "province_id": "15", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1505", "code": "1505", "name": "KAB. MUARO JAMBI", "province_id": "15", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1506", "code": "1506", "name": "KAB. TANJUNG JABUNG TIMUR", "province_id": "15", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1507", "code": "1507", "name": "KAB. TANJUNG JABUNG BARAT", "province_id": "15", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1508", "code": "1508", "name": "KAB. TEBO", "province_id": "15", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1509", "code": "1509", "name": "KAB. BUNGO", "province_id": "15", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1571", "code": "1571", "name": "KOTA JAMBI", "province_id": "15", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1572", "code": "1572", "name": "KOTA SUNGAI PENUH", "province_id": "15", "created_at": time.Now(), "updated_at": time.Now()},

		// === SUMATERA SELATAN (16) ===
		{"id": "1601", "code": "1601", "name": "KAB. OGAN KOMERING ULU", "province_id": "16", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1602", "code": "1602", "name": "KAB. OGAN KOMERING ILIR", "province_id": "16", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1603", "code": "1603", "name": "KAB. MUARA ENIM", "province_id": "16", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1604", "code": "1604", "name": "KAB. LAHAT", "province_id": "16", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1605", "code": "1605", "name": "KAB. MUSI RAWAS", "province_id": "16", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1606", "code": "1606", "name": "KAB. MUSI BANYUASIN", "province_id": "16", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1607", "code": "1607", "name": "KAB. BANYUASIN", "province_id": "16", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1608", "code": "1608", "name": "KAB. OGAN KOMERING ULU SELATAN", "province_id": "16", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1609", "code": "1609", "name": "KAB. OGAN KOMERING ULU TIMUR", "province_id": "16", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1610", "code": "1610", "name": "KAB. OGAN ILIR", "province_id": "16", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1611", "code": "1611", "name": "KAB. EMPAT LAWANG", "province_id": "16", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1612", "code": "1612", "name": "KAB. PENUKAL ABAB LEMATANG ILIR", "province_id": "16", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1613", "code": "1613", "name": "KAB. MUSI RAWAS UTARA", "province_id": "16", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1671", "code": "1671", "name": "KOTA PALEMBANG", "province_id": "16", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1672", "code": "1672", "name": "KOTA PAGAR ALAM", "province_id": "16", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1673", "code": "1673", "name": "KOTA LUBUK LINGGAU", "province_id": "16", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1674", "code": "1674", "name": "KOTA PRABUMULIH", "province_id": "16", "created_at": time.Now(), "updated_at": time.Now()},

		// === BENGKULU (17) ===
		{"id": "1701", "code": "1701", "name": "KAB. BENGKULU SELATAN", "province_id": "17", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1702", "code": "1702", "name": "KAB. REJANG LEBONG", "province_id": "17", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1703", "code": "1703", "name": "KAB. BENGKULU UTARA", "province_id": "17", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1704", "code": "1704", "name": "KAB. KAUR", "province_id": "17", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1705", "code": "1705", "name": "KAB. SELUMA", "province_id": "17", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1706", "code": "1706", "name": "KAB. MUKOMUKO", "province_id": "17", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1707", "code": "1707", "name": "KAB. LEBONG", "province_id": "17", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1708", "code": "1708", "name": "KAB. KEPAHIANG", "province_id": "17", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1709", "code": "1709", "name": "KAB. BENGKULU TENGAH", "province_id": "17", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1771", "code": "1771", "name": "KOTA BENGKULU", "province_id": "17", "created_at": time.Now(), "updated_at": time.Now()},

		// === LAMPUNG (18) ===
		{"id": "1801", "code": "1801", "name": "KAB. LAMPUNG BARAT", "province_id": "18", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1802", "code": "1802", "name": "KAB. TANGGAMUS", "province_id": "18", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1803", "code": "1803", "name": "KAB. LAMPUNG SELATAN", "province_id": "18", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1804", "code": "1804", "name": "KAB. LAMPUNG TIMUR", "province_id": "18", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1805", "code": "1805", "name": "KAB. LAMPUNG TENGAH", "province_id": "18", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1806", "code": "1806", "name": "KAB. LAMPUNG UTARA", "province_id": "18", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1807", "code": "1807", "name": "KAB. WAY KANAN", "province_id": "18", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1808", "code": "1808", "name": "KAB. TULANGBAWANG", "province_id": "18", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1809", "code": "1809", "name": "KAB. PESAWARAN", "province_id": "18", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1810", "code": "1810", "name": "KAB. PRINGSEWU", "province_id": "18", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1811", "code": "1811", "name": "KAB. MESUJI", "province_id": "18", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1812", "code": "1812", "name": "KAB. TULANG BAWANG BARAT", "province_id": "18", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1813", "code": "1813", "name": "KAB. PESISIR BARAT", "province_id": "18", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1871", "code": "1871", "name": "KOTA BANDAR LAMPUNG", "province_id": "18", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1872", "code": "1872", "name": "KOTA METRO", "province_id": "18", "created_at": time.Now(), "updated_at": time.Now()},

		// === BANGKA BELITUNG (19) ===
		{"id": "1901", "code": "1901", "name": "KAB. BANGKA", "province_id": "19", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1902", "code": "1902", "name": "KAB. BELITUNG", "province_id": "19", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1903", "code": "1903", "name": "KAB. BANGKA BARAT", "province_id": "19", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1904", "code": "1904", "name": "KAB. BANGKA TENGAH", "province_id": "19", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1905", "code": "1905", "name": "KAB. BANGKA SELATAN", "province_id": "19", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1906", "code": "1906", "name": "KAB. BELITUNG TIMUR", "province_id": "19", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "1971", "code": "1971", "name": "KOTA PANGKAL PINANG", "province_id": "19", "created_at": time.Now(), "updated_at": time.Now()},

		// === KEPULAUAN RIAU (21) ===
		{"id": "2101", "code": "2101", "name": "KAB. KEPULAUAN RIAU", "province_id": "21", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "2102", "code": "2102", "name": "KAB. KARIMUN", "province_id": "21", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "2103", "code": "2103", "name": "KAB. NATUNA", "province_id": "21", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "2104", "code": "2104", "name": "KAB. LINGGA", "province_id": "21", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "2105", "code": "2105", "name": "KAB. KEPULAUAN ANAMBAS", "province_id": "21", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "2106", "code": "2106", "name": "KAB. BINTAN", "province_id": "21", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "2171", "code": "2171", "name": "KOTA BATAM", "province_id": "21", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "2172", "code": "2172", "name": "KOTA TANJUNG PINANG", "province_id": "21", "created_at": time.Now(), "updated_at": time.Now()},

		// === DKI JAKARTA (31) ===
		{"id": "3101", "code": "3101", "name": "KAB. KEPULAUAN SERIBU", "province_id": "31", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3171", "code": "3171", "name": "KOTA JAKARTA PUSAT", "province_id": "31", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3172", "code": "3172", "name": "KOTA JAKARTA UTARA", "province_id": "31", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3173", "code": "3173", "name": "KOTA JAKARTA BARAT", "province_id": "31", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3174", "code": "3174", "name": "KOTA JAKARTA SELATAN", "province_id": "31", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3175", "code": "3175", "name": "KOTA JAKARTA TIMUR", "province_id": "31", "created_at": time.Now(), "updated_at": time.Now()},

		// === JAWA BARAT (32) ===
		{"id": "3201", "code": "3201", "name": "KAB. BOGOR", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3202", "code": "3202", "name": "KAB. SUKABUMI", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3203", "code": "3203", "name": "KAB. CIANJUR", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3204", "code": "3204", "name": "KAB. BANDUNG", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3205", "code": "3205", "name": "KAB. GARUT", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3206", "code": "3206", "name": "KAB. TASIKMALAYA", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3207", "code": "3207", "name": "KAB. CIAMIS", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3208", "code": "3208", "name": "KAB. KUNINGAN", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3209", "code": "3209", "name": "KAB. CIREBON", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3210", "code": "3210", "name": "KAB. MAJALENGKA", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3211", "code": "3211", "name": "KAB. SUMEDANG", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3212", "code": "3212", "name": "KAB. INDRAMAYU", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3213", "code": "3213", "name": "KAB. SUBANG", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3214", "code": "3214", "name": "KAB. PURWAKARTA", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3215", "code": "3215", "name": "KAB. KARAWANG", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3216", "code": "3216", "name": "KAB. BEKASI", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3217", "code": "3217", "name": "KAB. BANDUNG BARAT", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3218", "code": "3218", "name": "KAB. PANGANDARAN", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3271", "code": "3271", "name": "KOTA BOGOR", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3272", "code": "3272", "name": "KOTA SUKABUMI", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3273", "code": "3273", "name": "KOTA BANDUNG", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3274", "code": "3274", "name": "KOTA CIREBON", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3275", "code": "3275", "name": "KOTA BEKASI", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3276", "code": "3276", "name": "KOTA DEPOK", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3277", "code": "3277", "name": "KOTA CIMAHI", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3278", "code": "3278", "name": "KOTA TASIKMALAYA", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
		{"id": "3279", "code": "3279", "name": "KOTA BANJAR", "province_id": "32", "created_at": time.Now(), "updated_at": time.Now()},
	}
	return batchInsert(db, "regencies", data, 50)
}

// =============================================================================
// Performance Management (KPI) Master Data Seeders
// =============================================================================

// ── Performance Perspectives (BSC Perspectives) ──
// Schema: id(char(36) PK), name, description, sort_order, timestamps
func seedPerformancePerspectives(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		{"id": codeToUUID("perf_perspective", "FIN"), "name": "Financial", "description": "Perspektif keuangan - mengukur kinerja finansial dan profitabilitas", "sort_order": 1, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("perf_perspective", "CUS"), "name": "Customer", "description": "Perspektif pelanggan - mengukur kepuasan dan loyalitas pelanggan", "sort_order": 2, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("perf_perspective", "INT"), "name": "Internal Process", "description": "Perspektif proses internal - mengukur efisiensi dan efektivitas proses bisnis", "sort_order": 3, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("perf_perspective", "LRN"), "name": "Learning & Growth", "description": "Perspektif pembelajaran dan pertumbuhan - mengukur pengembangan SDM dan inovasi", "sort_order": 4, "created_at": time.Now(), "updated_at": time.Now()},
	}
	return batchInsert(db, "performance_perspectives", data, 50)
}

// ── Performance Ratings ──
// Schema: id(char(36) PK), code, name, min_score, max_score, color, description, sort_order, timestamps
func seedPerformanceRatings(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		{"id": codeToUUID("perf_rating", "OUT"), "code": "OUT", "name": "Outstanding", "min_score": 95.00, "max_score": 100.00, "color": "success", "description": "Kinerja luar biasa, melampaui target dengan signifikan", "sort_order": 1, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("perf_rating", "EXC"), "code": "EXC", "name": "Excellent", "min_score": 85.00, "max_score": 94.99, "color": "primary", "description": "Kinerja sangat baik, melampaui target", "sort_order": 2, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("perf_rating", "GOO"), "code": "GOO", "name": "Good", "min_score": 75.00, "max_score": 84.99, "color": "info", "description": "Kinerja baik, memenuhi target", "sort_order": 3, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("perf_rating", "FAI"), "code": "FAI", "name": "Fair", "min_score": 60.00, "max_score": 74.99, "color": "warning", "description": "Kinerja cukup, mendekati target", "sort_order": 4, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("perf_rating", "POO"), "code": "POO", "name": "Poor", "min_score": 0.00, "max_score": 59.99, "color": "danger", "description": "Kinerja kurang, di bawah target", "sort_order": 5, "created_at": time.Now(), "updated_at": time.Now()},
	}
	return batchInsert(db, "performance_ratings", data, 50)
}

// ── Performance Indicator Formulas ──
// Schema: id(char(36) PK), code, name, formula_type, expression, description, sort_order, timestamps
func seedPerformanceIndicatorFormulas(db *gorm.DB) (int, int, error) {
	data := []map[string]interface{}{
		{"id": codeToUUID("perf_formula", "MANUAL"), "code": "MANUAL", "name": "Manual Score", "formula_type": "MANUAL", "expression": nil, "description": "Nilai diinput manual oleh reviewer berdasarkan penilaian kualitatif", "sort_order": 1, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("perf_formula", "HIGHER"), "code": "HIGHER", "name": "Higher Better", "formula_type": "HIGHER_BETTER", "expression": "(actual / target) * 100", "description": "Semakin tinggi nilai aktual semakin baik. Achievement = (Actual / Target) x 100", "sort_order": 2, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("perf_formula", "LOWER"), "code": "LOWER", "name": "Lower Better", "formula_type": "LOWER_BETTER", "expression": "(target / actual) * 100", "description": "Semakin rendah nilai aktual semakin baik. Achievement = (Target / Actual) x 100", "sort_order": 3, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("perf_formula", "RANGE"), "code": "RANGE", "name": "Range Score", "formula_type": "RANGE", "expression": nil, "description": "Nilai berdasarkan rentang tertentu yang telah ditentukan", "sort_order": 4, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("perf_formula", "BOOLEAN"), "code": "BOOLEAN", "name": "Boolean", "formula_type": "BOOLEAN", "expression": "actual == 1 ? 100 : 0", "description": "Nilai Ya/Tidak. Ya = 100%, Tidak = 0%", "sort_order": 5, "created_at": time.Now(), "updated_at": time.Now()},
		{"id": codeToUUID("perf_formula", "PERCENTAGE"), "code": "PERCENTAGE", "name": "Percentage", "formula_type": "PERCENTAGE", "expression": "actual", "description": "Nilai langsung berupa persentase (0-100)", "sort_order": 6, "created_at": time.Now(), "updated_at": time.Now()},
	}
	return batchInsert(db, "performance_indicator_formulas", data, 50)
}
