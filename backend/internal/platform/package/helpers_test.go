package pkgmgr

import (
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

// seedModuleInfo menyimpan data modul yang akan di-seed ke tabel modules untuk test.
type seedModuleInfo struct {
	ID        string
	Name      string
	Slug      string
	DependsOn string
}

// testModules adalah daftar modul test yang digunakan untuk validasi dependensi.
var testModules = []seedModuleInfo{
	{Name: "Organization Management", Slug: "organization", DependsOn: ""},
	{Name: "Employee Management", Slug: "employee", DependsOn: "organization"},
	{Name: "Job Management", Slug: "job-management", DependsOn: "organization"},
	{Name: "Competency Management", Slug: "competency", DependsOn: "organization,employee"},
	{Name: "Employee Movement", Slug: "employee-movement", DependsOn: "employee,organization"},
	{Name: "Attendance Management", Slug: "attendance", DependsOn: "employee,organization"},
	{Name: "Payroll Management", Slug: "payroll", DependsOn: "employee,organization"},
	{Name: "Leave Management", Slug: "leave", DependsOn: "employee,organization"},
	{Name: "Performance Management", Slug: "performance", DependsOn: "organization,employee,job-management,competency"},
	{Name: "Reimbursement", Slug: "reimbursement", DependsOn: "employee"},
}

// setupTestDB creates an in-memory SQLite database and auto-migrates Package + PackageModule.
// Also creates the "modules" table and seeds test module data for dependency validation.
// Returns the GORM DB and a cleanup function.
func setupTestDB() (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	// Auto-migrate Package + PackageModule models
	if err := db.AutoMigrate(&Package{}, &PackageModule{}); err != nil {
		panic(fmt.Sprintf("failed to migrate test db: %v", err))
	}

	// Create "modules" table manually (simulates platform module registry)
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS modules (
			id CHAR(36) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(100) NOT NULL UNIQUE,
			version VARCHAR(20) NOT NULL DEFAULT '1.0.0',
			description TEXT,
			is_core TINYINT(1) NOT NULL DEFAULT 0,
			depends_on TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`).Error; err != nil {
		panic(fmt.Sprintf("failed to create modules table: %v", err))
	}

	// Seed test modules — gunakan index loop karena seedModuleInfo adalah value type
	for i := range testModules {
		mID := uuid.New()
		testModules[i].ID = mID.String()
		if err := db.Table("modules").Create(map[string]interface{}{
			"id":         mID.String(),
			"name":       testModules[i].Name,
			"slug":       testModules[i].Slug,
			"version":    "1.0.0",
			"is_core":    0,
			"depends_on": testModules[i].DependsOn,
		}).Error; err != nil {
			panic(fmt.Sprintf("failed to seed module '%s': %v", testModules[i].Slug, err))
		}
	}

	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	return db, cleanup
}

// createTestPackage inserts a test package directly into the database.
// Returns the created Package with its module associations.
func createTestPackage(db *gorm.DB, name, slug string, moduleIDs []uuid.UUID) *Package {
	pkg := &Package{
		Name:   name,
		Slug:   slug,
		Status: string(PackageDraft),
	}

	// Fetch module info from the test modules table
	for i, mid := range moduleIDs {
		var modName, modSlug string
		row := db.Table("modules").Select("name, slug").Where("id = ?", mid.String()).Row()
		if err := row.Scan(&modName, &modSlug); err != nil {
			panic(fmt.Sprintf("failed to find module %s: %v", mid.String(), err))
		}

		pkg.Modules = append(pkg.Modules, PackageModule{
			ModuleID:    mid,
			IsMandatory: i == 1, // second module is mandatory
			SortOrder:   i,
			ModuleName:  modName,
			ModuleSlug:  modSlug,
		})
	}

	if err := db.Create(pkg).Error; err != nil {
		panic(fmt.Sprintf("failed to create test package: %v", err))
	}

	return pkg
}

// findModuleUUIDBySlug mencari UUID modul test berdasarkan slug.
func findModuleUUIDBySlug(slug string) uuid.UUID {
	for _, m := range testModules {
		if m.Slug == slug {
			uid, err := uuid.Parse(m.ID)
			if err != nil {
				panic(fmt.Sprintf("invalid module uuid for slug '%s': %v", slug, err))
			}
			return uid
		}
	}
	panic(fmt.Sprintf("module slug '%s' not found in test data", slug))
}

// newTestService creates a Service with SQLite repository and test module data.
// Returns the service and a cleanup function.
func newTestService() (*Service, func()) {
	db, cleanup := setupTestDB()
	repo := NewRepository(db)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	return svc, func() {
		cleanup()
		_ = logger.Sync()
	}
}

// uuidStr returns a UUID string for test use.
func uuidStr() string {
	return uuid.New().String()
}
