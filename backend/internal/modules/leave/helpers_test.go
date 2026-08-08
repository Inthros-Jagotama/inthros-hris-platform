package leave

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

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

// setupTestDB creates an in-memory SQLite database and auto-migrates all models.
func setupTestDB() (*gorm.DB, func(ctx context.Context) (*gorm.DB, error), func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	if err := db.AutoMigrate(
		&LeaveType{},
		&LeaveAccrualPolicy{},
		&LeaveReason{},
		&LeaveRequest{},
		&LeaveRequestDetail{},
		&EmployeeLeaveBalance{},
		&LeaveBalanceTransaction{},
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

var _leaveTypeCounter int

func createTestLeaveType(repo *Repository) *LeaveType {
	ctx := context.Background()
	_leaveTypeCounter++
	t := &LeaveType{
		Name:         fmt.Sprintf("Annual Leave %d", _leaveTypeCounter),
		Description:  "Annual paid leave",
		IsPaid:       true,
		AllowHalfDay: true,
	}
	if err := repo.CreateLeaveType(ctx, t); err != nil {
		panic(fmt.Sprintf("failed to create test leave type: %v", err))
	}
	return t
}

func createTestLeaveReason(repo *Repository) *LeaveReason {
	ctx := context.Background()
	r := &LeaveReason{
		Name:      "Sick",
		SortOrder: 1,
	}
	if err := repo.CreateLeaveReason(ctx, r); err != nil {
		panic(fmt.Sprintf("failed to create test leave reason: %v", err))
	}
	return r
}

func createTestLeaveRequest(repo *Repository, empID uuid.UUID, lTypeID uuid.UUID) *LeaveRequest {
	ctx := context.Background()
	lr := &LeaveRequest{
		EmployeeID:       empID,
		LeaveTypeID:      lTypeID,
		RequestStartDate: "2026-01-15",
		RequestEndDate:   "2026-01-16",
		DurationMode:     DurationFullDay,
		RequestedDays:    2,
		Status:           LeaveStatusSubmitted,
	}
	if err := repo.CreateLeaveRequest(ctx, lr); err != nil {
		panic(fmt.Sprintf("failed to create test leave request: %v", err))
	}
	return lr
}

func createTestAccrualPolicy(repo *Repository, lTypeID uuid.UUID) *LeaveAccrualPolicy {
	ctx := context.Background()
	p := &LeaveAccrualPolicy{
		LeaveTypeID:   lTypeID,
		BaseQuotaDays: 12,
		EffectiveFrom: "2026-01-01",
	}
	if err := repo.CreateAccrualPolicy(ctx, p); err != nil {
		panic(fmt.Sprintf("failed to create test accrual policy: %v", err))
	}
	return p
}

func uuidStr() string {
	return uuid.New().String()
}
