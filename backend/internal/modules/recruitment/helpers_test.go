package recruitment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

func setupTestDB() (*gorm.DB, func(ctx context.Context) (*gorm.DB, error), func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	if err := db.AutoMigrate(
		&JobRequisition{},
		&Candidate{},
		&JobApplication{},
		&Interview{},
		&OnboardingTaskTemplate{},
		&EmployeeOnboarding{},
		&OnboardingTaskItem{},
	); err != nil {
		panic(fmt.Sprintf("failed to migrate test db: %v", err))
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

func createTestUUID() string {
	return uuid.New().String()
}

func createTestOrgID() string {
	return uuid.New().String()
}

func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
