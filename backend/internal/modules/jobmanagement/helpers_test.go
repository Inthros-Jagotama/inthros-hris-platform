package jobmanagement

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

// setupTestDB creates an in-memory SQLite database and auto-migrates all job management models.
// Returns the GORM DB, a dbResolver function, and a cleanup function.
func setupTestDB() (*gorm.DB, func(ctx context.Context) (*gorm.DB, error), func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	// AutoMigrate all models
	if err := db.AutoMigrate(
		&JobTitle{},
		&JobTitleSub{},
		&JobValue{},
		&JobManagementValueCluster{},
		&JobObjective{},
		&JobIdentification{},
		&JobResponsibility{},
		&JobEducationExperience{},
		&JobManagementMajor{},
		&JobManagementJobFamily{},
		&JobHRAuthority{},
		&JobOperationalAuthority{},
		&JobWorkingActivity{},
		&JobWorkingRisk{},
		&JobRelationship{},
		&JobManagementRelationshipDetail{},
		&JobSubordinateControl{},
		&JobAsset{},
		&JobFinancial{},
		&JobPotencyCompetency{},
		&JobScore{},
		&JobCompetencyGroup{},
		// Tabel organizations minimal agar Preload("Organization") pada
		// JobManagementRelationshipDetail teruji (organizations asli di migrasi tenant).
		&OrganizationRef{},
	); err != nil {
		panic(fmt.Sprintf("failed to migrate test db: %v", err))
	}

	// Tabel milik module lain yang dipakai query raw dashboard
	// (GetOrganizationSummaryDashboard) — minimal schema, kolom yang disentuh
	// query saja.
	// Ganti tabel organizations hasil AutoMigrate OrganizationRef (hanya
	// id/nomenclature/full_code) dengan versi lengkap berisi kolom yang dipakai
	// query raw dashboard: organization_summary_id & code.
	rawTables := []string{
		`DROP TABLE IF EXISTS organizations`,
		`CREATE TABLE organizations (
			id CHAR(36) PRIMARY KEY,
			organization_summary_id CHAR(36) NULL,
			code VARCHAR(10) NOT NULL DEFAULT '',
			full_code VARCHAR(50) NOT NULL DEFAULT '',
			nomenclature VARCHAR(255) NOT NULL DEFAULT '',
			deleted_at DATETIME NULL
		)`,
		`CREATE TABLE IF NOT EXISTS organization_summaries (
			id CHAR(36) PRIMARY KEY,
			code VARCHAR(7) NOT NULL,
			decree_no VARCHAR(20) NOT NULL,
			decree_date DATE NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'inactive',
			deleted_at DATETIME NULL
		)`,
		`CREATE TABLE IF NOT EXISTS employments (
			id CHAR(36) PRIMARY KEY,
			employee_id CHAR(36) NULL,
			organization_id CHAR(36) NULL,
			effective_end_date DATE NULL,
			deleted_at DATETIME NULL
		)`,
	}
	for _, q := range rawTables {
		if err := db.Exec(q).Error; err != nil {
			panic(fmt.Sprintf("failed to create raw table: %v", err))
		}
	}

	dbResolver := func(ctx context.Context) (*gorm.DB, error) {
		return db, nil
	}

	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	return db, dbResolver, cleanup
}

// createTestJobTitle inserts a test job title and returns it.
func createTestJobTitle(ctx context.Context, repo *Repository) *JobTitle {
	name := "Test Title"
	t := &JobTitle{
		Name:   &name,
		Status: int8Ptr(1),
	}
	if err := repo.CreateJobTitle(ctx, t); err != nil {
		panic(fmt.Sprintf("failed to create test job title: %v", err))
	}
	return t
}

// createTestJobValue inserts a test job value with a specific type.
func createTestJobValue(ctx context.Context, repo *Repository, valueType string) *JobValue {
	v := &JobValue{
		Type: valueType,
		Sort: intPtr(1),
	}
	if err := repo.CreateJobValue(ctx, v); err != nil {
		panic(fmt.Sprintf("failed to create test job value: %v", err))
	}
	return v
}

// createTestOrgID returns a UUID string for use as organization_id in tests.
func createTestOrgID() string {
	return uuid.New().String()
}

// strPtr returns a pointer to the given string.
func strPtr(s string) *string {
	return &s
}

// intPtr returns a pointer to the given int.
func intPtr(i int) *int {
	return &i
}

// int8Ptr returns a pointer to the given int8.
func int8Ptr(i int8) *int8 {
	return &i
}

// float64Ptr returns a pointer to the given float64.
func float64Ptr(f float64) *float64 {
	return &f
}

// boolPtr returns a pointer to the given bool.
func boolPtr(b bool) *bool {
	return &b
}

// uint64Ptr returns a pointer to the given uint64.
func uint64Ptr(u uint64) *uint64 {
	return &u
}
