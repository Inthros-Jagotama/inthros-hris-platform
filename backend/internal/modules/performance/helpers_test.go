package performance

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

// setupTestDB creates an in-memory SQLite database and auto-migrates all performance models.
func setupTestDB() (*gorm.DB, func(ctx context.Context) (*gorm.DB, error), func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	if err := db.AutoMigrate(
		&PerformancePeriod{},
		&PerformancePerspective{},
		&PerformanceTemplate{},
		&PerformanceIndicator{},
		&PerformanceEvaluation{},
		&PerformanceEvaluationDetail{},
		&PerformanceTarget{},
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

// createTestUUID returns a UUID string for general test use.
func createTestUUID() string {
	return uuid.New().String()
}

// createTestOrgID returns a UUID string for organization_id in tests.
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

// float64Ptr returns a pointer to the given float64.
func float64Ptr(f float64) *float64 {
	return &f
}
