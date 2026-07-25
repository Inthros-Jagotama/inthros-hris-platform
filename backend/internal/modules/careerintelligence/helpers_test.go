package careerintelligence

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

// setupTestDB creates an in-memory SQLite database and auto-migrates all models.
func setupTestDB() (*gorm.DB, func(ctx context.Context) (*gorm.DB, error), func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	if err := db.AutoMigrate(
		&CareerTalentMap{},
		&CareerInterest{},
		&CareerPath{},
		&CareerSuccessionPlan{},
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

// newTestService creates a Service with in-memory SQLite repository.
func newTestService() (*Service, *Repository, *gorm.DB, func()) {
	db, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	return svc, repo, db, func() {
		cleanup()
		_ = logger.Sync()
	}
}

// =========================================================================
// Fixture helpers
// =========================================================================

var _fixtureCounter int

func createTestTalentMap(repo *Repository) *CareerTalentMap {
	_fixtureCounter++
	ctx := context.Background()
	perf := "HIGH"
	pot := "MEDIUM"
	if _fixtureCounter%3 == 0 {
		perf = "MEDIUM"
		pot = "HIGH"
	} else if _fixtureCounter%3 == 1 {
		perf = "LOW"
		pot = "LOW"
	}
	tm := &CareerTalentMap{
		ID:           uuid.New(),
		EmployeeID:   uuid.New(),
		Period:       fmt.Sprintf("2026-%02d", (_fixtureCounter%12)+1),
		Performance:  perf,
		Potential:    pot,
		GridPosition: computeGridPosition(perf, pot),
		Notes:        fmt.Sprintf("Test talent map %d", _fixtureCounter),
	}
	if err := repo.CreateTalentMap(ctx, tm); err != nil {
		panic(fmt.Sprintf("failed to create test talent map: %v", err))
	}
	return tm
}

func createTestCareerInterest(repo *Repository) *CareerInterest {
	_fixtureCounter++
	ctx := context.Background()
	types := []string{"LEADERSHIP", "SPECIALIST", "INTERNATIONAL", "ENTREPRENEUR"}
	readiness := []string{"NOW", "1_YEAR", "2_3_YEARS", "3_PLUS"}
	ci := &CareerInterest{
		ID:              uuid.New(),
		EmployeeID:      uuid.New(),
		InterestType:    types[_fixtureCounter%len(types)],
		TargetPosition:  fmt.Sprintf("Position %d", _fixtureCounter),
		ReadinessLevel:  readiness[_fixtureCounter%len(readiness)],
		IsActive:        true,
	}
	if err := repo.CreateCareerInterest(ctx, ci); err != nil {
		panic(fmt.Sprintf("failed to create test career interest: %v", err))
	}
	return ci
}

func createTestCareerPath(repo *Repository) *CareerPath {
	_fixtureCounter++
	ctx := context.Background()
	pathTypes := []string{"PROMOTION", "LATERAL", "DEMOTION", "CROSSFUNCTIONAL"}
	cp := &CareerPath{
		ID:            uuid.New(),
		SourceTitleID:  uuid.New(),
		TargetTitleID:  uuid.New(),
		PathType:       pathTypes[_fixtureCounter%len(pathTypes)],
		TypicalTenure:  24,
		Requirements:   "Bachelor degree, 5 years experience",
		IsActive:       true,
	}
	if err := repo.CreateCareerPath(ctx, cp); err != nil {
		panic(fmt.Sprintf("failed to create test career path: %v", err))
	}
	return cp
}

func createTestSuccessionPlan(repo *Repository) *CareerSuccessionPlan {
	_fixtureCounter++
	ctx := context.Background()
	readiness := []string{"READY_NOW", "READY_1YR", "READY_2YR", "POTENTIAL"}
	sp := &CareerSuccessionPlan{
		ID:             uuid.New(),
		PositionID:     uuid.New(),
		SuccessorID:    uuid.New(),
		ReadinessLevel: readiness[_fixtureCounter%len(readiness)],
		PriorityOrder:  (_fixtureCounter % 3) + 1,
		Status:         "ACTIVE",
	}
	if err := repo.CreateSuccessionPlan(ctx, sp); err != nil {
		panic(fmt.Sprintf("failed to create test succession plan: %v", err))
	}
	return sp
}

func ctx() context.Context {
	return context.Background()
}
