package employee

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
	"github.com/google/uuid"

	"github.com/inthros/hris-platform/internal/modules/organization"
	"github.com/inthros/hris-platform/internal/modules/setting"
)

// setupTestDB creates an in-memory SQLite database and auto-migrates all employee models.
// Returns the GORM DB, a dbResolver function, and a cleanup function.
func setupTestDB() (*gorm.DB, func(ctx context.Context) (*gorm.DB, error), func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	// AutoMigrate all models + tabel referensi (settings/organization) yang
	// di-preload di FindEmployeeByID supaya preload tidak error saat FK terisi.
	if err := db.AutoMigrate(
		&Employee{},
		&EmployeeAddress{},
		&EmergencyContact{},
		&EmployeeFamily{},
		&EmployeeEducation{},
		&EmployeeExperience{},
		&EmployeeDocument{},
		&EmployeeInsurance{},
		&EmployeeBankAccount{},
		&Employment{},
		&SensitiveFieldSetting{},
		&setting.Religion{},
		&setting.MaritalStatus{},
		&setting.Nationality{},
		&setting.RelationshipType{},
		&setting.Education{},
		&setting.EducationMajor{},
		&setting.Insurance{},
		&setting.Bank{},
		&setting.EmploymentStatus{},
		&organization.Organization{},
	); err != nil {
		panic(fmt.Sprintf("failed to migrate test db: %v", err))
	}

	// Seed sensitive field settings (all disabled by default), mirroring the
	// real 151_sensitive_field_settings migration, so Create/Update's
	// encryptIfEnabled has real rows to check against.
	for _, d := range SensitiveFieldRegistry {
		db.Create(&SensitiveFieldSetting{ID: uuid.New().String(), FieldKey: d.Key, IsEncryptionEnabled: false})
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

// createTestEmployee inserts a test employee and returns it.
func createTestEmployee(ctx context.Context, repo *Repository) *Employee {
	emp := &Employee{
		EmployeeID: "EMP-TEST-001",
		Name:       "Test Employee",
		Gender:     strPtr("M"),
		Status:     "active",
	}
	if err := repo.CreateEmployee(ctx, emp); err != nil {
		panic(fmt.Sprintf("failed to create test employee: %v", err))
	}
	return emp
}

// strPtr returns a pointer to the given string.
func strPtr(s string) *string {
	return &s
}

// intPtr returns a pointer to the given int.
func intPtr(i int) *int {
	return &i
}


