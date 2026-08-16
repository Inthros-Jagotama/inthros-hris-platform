package reimbursement

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
		&ReimbursementType{},
		&ReimbursementRequest{},
		&ReimbursementItem{},
		&testEmployeeAccount{},
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

var _reimbTypeCounter int

func createTestReimbursementType(repo *Repository) *ReimbursementType {
	ctx := context.Background()
	_reimbTypeCounter++
	t := &ReimbursementType{
		Code:        fmt.Sprintf("TYPE%d", _reimbTypeCounter),
		Name:        fmt.Sprintf("Reimbursement Type %d", _reimbTypeCounter),
		Description: "Test reimbursement type",
		IsActive:    true,
	}
	if err := repo.CreateReimbursementType(ctx, t); err != nil {
		panic(fmt.Sprintf("failed to create test reimbursement type: %v", err))
	}
	return t
}

func createTestReimbursementRequest(repo *Repository, empID uuid.UUID, rTypeID uuid.UUID) *ReimbursementRequest {
	ctx := context.Background()
	rr := &ReimbursementRequest{
		EmployeeID:    empID,
		RequestTypeID: rTypeID,
		Title:         "Test Reimbursement",
		Description:   "Test reimbursement request",
		TotalAmount:   0,
		Currency:      "IDR",
		Status:        ReimbStatusDraft,
	}
	if err := repo.CreateReimbursementRequest(ctx, rr); err != nil {
		panic(fmt.Sprintf("failed to create test reimbursement request: %v", err))
	}
	return rr
}

func createTestReimbursementItem(repo *Repository, requestID uuid.UUID) *ReimbursementItem {
	ctx := context.Background()
	item := &ReimbursementItem{
		ReimbursementRequestID: requestID,
		ExpenseDate:            "2026-07-15",
		ExpenseType:            "MEDICAL",
		Description:            "Doctor visit",
		Amount:                 250000,
	}
	if err := repo.CreateReimbursementItem(ctx, item); err != nil {
		panic(fmt.Sprintf("failed to create test reimbursement item: %v", err))
	}
	return item
}

// testEmployeeAccount mirrors useraccount.EmployeeAccount's employee_id <->
// user_id mapping (employee_accounts table) for FindUserIDByEmployeeID tests,
// without importing the useraccount package.
type testEmployeeAccount struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey"`
	EmployeeID uuid.UUID `gorm:"type:char(36);not null;uniqueIndex"`
	UserID     uuid.UUID `gorm:"type:char(36);not null"`
}

func (testEmployeeAccount) TableName() string { return "employee_accounts" }

func createTestEmployeeAccount(db *gorm.DB, employeeID, userID uuid.UUID) {
	acc := &testEmployeeAccount{ID: uuid.New(), EmployeeID: employeeID, UserID: userID}
	if err := db.Create(acc).Error; err != nil {
		panic(fmt.Sprintf("failed to create test employee account: %v", err))
	}
}

func uuidStr() string {
	return uuid.New().String()
}

func boolPtr(b bool) *bool {
	return &b
}

func float64Ptr(f float64) *float64 {
	return &f
}
