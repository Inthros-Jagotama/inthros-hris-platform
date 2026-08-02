package setting

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

// setupTestDB creates an in-memory SQLite database and auto-migrates
// all setting module models. Returns the GORM DB and a cleanup function.
func setupTestDB() (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	if err := db.AutoMigrate(
		&Zone{},
		&Province{},
		&Regency{},
		&District{},
		&Village{},
		&Education{},
		&Religion{},
		&MaritalStatus{},
		&RelationshipType{},
		&EmploymentStatus{},
		&Bank{},
		&Nationality{},
		&JobFamily{},
		&Grading{},
		&SalaryGrade{},
		&TER{},
		&PTKP{},
		&Insurance{},
		&CompanyHoliday{},
		&Competency{},
	); err != nil {
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

// testDBResolver creates a TenantDBFunc that always returns the given
// in-memory SQLite DB regardless of context.
func testDBResolver(db *gorm.DB) func(ctx context.Context) (*gorm.DB, error) {
	return func(ctx context.Context) (*gorm.DB, error) {
		return db, nil
	}
}

// newTestRepo creates a Repository backed by in-memory SQLite.
// Returns the repository and a cleanup function.
func newTestRepo() (*Repository, *gorm.DB, func()) {
	db, cleanup := setupTestDB()
	repo := NewRepository(testDBResolver(db))
	return repo, db, cleanup
}

// newTestService creates a Service backed by in-memory SQLite.
// Returns the service and a cleanup function.
func newTestService() (*Service, *gorm.DB, func()) {
	db, cleanup := setupTestDB()
	repo := NewRepository(testDBResolver(db))
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	return svc, db, func() {
		cleanup()
		_ = logger.Sync()
	}
}

// createTestZone is a test helper to insert a Zone directly via GORM.
func createTestZone(db *gorm.DB, code, name string) *Zone {
	z := &Zone{Code: code, Name: name, Zone: name}
	if err := db.Create(z).Error; err != nil {
		panic(fmt.Sprintf("failed to create test zone: %v", err))
	}
	return z
}

// createTestEducation is a test helper to insert an Education directly via GORM.
func createTestEducation(db *gorm.DB, code, name string) *Education {
	e := &Education{Code: code, Name: name}
	if err := db.Create(e).Error; err != nil {
		panic(fmt.Sprintf("failed to create test education: %v", err))
	}
	return e
}

// createTestReligion is a test helper to insert a Religion directly via GORM.
func createTestReligion(db *gorm.DB, code, name string) *Religion {
	r := &Religion{Code: code, Name: name}
	if err := db.Create(r).Error; err != nil {
		panic(fmt.Sprintf("failed to create test religion: %v", err))
	}
	return r
}

// createTestMaritalStatus is a test helper to insert a MaritalStatus directly via GORM.
func createTestMaritalStatus(db *gorm.DB, code, name string) *MaritalStatus {
	m := &MaritalStatus{Code: code, Name: name}
	if err := db.Create(m).Error; err != nil {
		panic(fmt.Sprintf("failed to create test marital status: %v", err))
	}
	return m
}

// createTestBank is a test helper to insert a Bank directly via GORM.
func createTestBank(db *gorm.DB, code, name string) *Bank {
	b := &Bank{Code: code, Name: name}
	if err := db.Create(b).Error; err != nil {
		panic(fmt.Sprintf("failed to create test bank: %v", err))
	}
	return b
}

// createTestJobFamily is a test helper to insert a JobFamily directly via GORM.
func createTestJobFamily(db *gorm.DB, code, name string) *JobFamily {
	jf := &JobFamily{Code: code, Name: name}
	if err := db.Create(jf).Error; err != nil {
		panic(fmt.Sprintf("failed to create test job family: %v", err))
	}
	return jf
}

// uuidStr returns a UUID string for test use.
func uuidStr() string {
	return uuid.New().String()
}
