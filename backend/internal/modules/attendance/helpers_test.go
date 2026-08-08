package attendance

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

// testEmployeeAccount mirrors useraccount.EmployeeAccount's employee_id <->
// user_id mapping (employee_accounts table) for FindUserIDByEmployeeID tests,
// without importing the useraccount package (same pattern as
// leave/helpers_test.go's testEmployeeAccount).
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
		&AttendanceCompanySetting{},
		&AttendanceCompanyShift{},
		&AttendanceEmployeeShift{},
		&AttendanceLocation{},
		&AttendanceDeviceCapture{},
		&AttendanceFaceCapture{},
		&AttendanceEvent{},
		&AttendanceSession{},
		&AttendanceOvertimeRequest{},
		&AttendanceExemptPosition{},
		&AttendanceCorrectionRequest{},
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

// newTestHandler creates a Gin router with attendance routes registered.
func newTestHandler() (*Service, *Repository, *gorm.DB, func()) {
	return newTestService()
}

// =========================================================================
// Fixture helpers
// =========================================================================

// createTestEvent inserts a test attendance event.
func createTestEvent(repo *Repository, employeeID uuid.UUID) *AttendanceEvent {
	ctx := context.Background()
	e := &AttendanceEvent{
		EmployeeID:       employeeID,
		EventType:        EventTypeCheckIn,
		EventTimeUTC:     parseTime("2026-01-15T00:00:00Z"),
		EventTimeLocal:   parseTime("2026-01-15T07:00:00+07:00"),
		Latitude:         -6.2088,
		Longitude:        106.8456,
		ValidationStatus: ValidationPending,
	}
	if err := repo.CreateEvent(ctx, e); err != nil {
		panic(fmt.Sprintf("failed to create test event: %v", err))
	}
	return e
}

// createTestOvertimeRequest inserts a test overtime request.
func createTestOvertimeRequest(repo *Repository, employeeID uuid.UUID) *AttendanceOvertimeRequest {
	ctx := context.Background()
	reason := "Test overtime"
	r := &AttendanceOvertimeRequest{
		EmployeeID:       employeeID,
		WorkDate:         "2026-01-15",
		StartTimeLocal:   parseTime("2026-01-15T18:00:00+07:00"),
		EndTimeLocal:     parseTime("2026-01-15T20:00:00+07:00"),
		RequestedMinutes: 120,
		Reason:           &reason,
		Status:           OvertimeSubmitted,
	}
	if err := repo.CreateOvertimeRequest(ctx, r); err != nil {
		panic(fmt.Sprintf("failed to create test overtime request: %v", err))
	}
	return r
}

// createTestShift inserts a test company shift.
func createTestShift(repo *Repository) *AttendanceCompanyShift {
	ctx := context.Background()
	s := &AttendanceCompanyShift{
		ShiftName:    "Morning Shift",
		CheckInTime:  "08:00:00",
		CheckOutTime: "17:00:00",
	}
	if err := repo.CreateShift(ctx, s); err != nil {
		panic(fmt.Sprintf("failed to create test shift: %v", err))
	}
	return s
}

// createTestEmployeeShift inserts a test employee shift assignment.
func createTestEmployeeShift(repo *Repository, employeeID, shiftID uuid.UUID) *AttendanceEmployeeShift {
	ctx := context.Background()
	es := &AttendanceEmployeeShift{
		EmployeeID:        employeeID,
		AttendanceShiftID: shiftID,
		EffectiveDateFrom: "2026-01-01",
	}
	if err := repo.CreateEmployeeShift(ctx, es); err != nil {
		panic(fmt.Sprintf("failed to create test employee shift: %v", err))
	}
	return es
}

// createTestLocation inserts a test geofence location.
func createTestLocation(repo *Repository) *AttendanceLocation {
	ctx := context.Background()
	l := &AttendanceLocation{
		Name:      "Main Office",
		Latitude:  -6.2088,
		Longitude: 106.8456,
		RadiusM:   50,
	}
	if err := repo.CreateLocation(ctx, l); err != nil {
		panic(fmt.Sprintf("failed to create test location: %v", err))
	}
	return l
}

// createTestSession inserts a test attendance session directly via DB.
// Note: uses *gorm.DB directly because Repository doesn't have CreateSession (sessions are computed).
func createTestSession(db *gorm.DB, employeeID uuid.UUID) *AttendanceSession {
	ctx := context.Background()
	s := &AttendanceSession{
		EmployeeID: employeeID,
		WorkDate:   "2026-01-15",
		Status:     SessionStatusOpen,
	}
	if err := db.WithContext(ctx).Create(s).Error; err != nil {
		panic(fmt.Sprintf("failed to create test session: %v", err))
	}
	return s
}

// createTestExemptPosition inserts a test exempt position.
func createTestExemptPosition(repo *Repository, orgID uuid.UUID) *AttendanceExemptPosition {
	ctx := context.Background()
	p := &AttendanceExemptPosition{
		OrganizationID: orgID,
		IsExempt:       true,
	}
	if err := repo.CreateExemptPosition(ctx, p); err != nil {
		panic(fmt.Sprintf("failed to create test exempt position: %v", err))
	}
	return p
}

// =========================================================================
// Helpers
// =========================================================================

// uuidStr returns a UUID string for test use.
func uuidStr() string {
	return uuid.New().String()
}

// boolPtr returns a pointer to the given bool.
func boolPtr(b bool) *bool {
	return &b
}

// parseTime parses an RFC3339 time string and returns a time.Time value.
func parseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		// Try without timezone
		t, err = time.Parse("2006-01-02T15:04:05", s)
		if err != nil {
			panic(fmt.Sprintf("failed to parse time '%s': %v", s, err))
		}
	}
	return t
}

// intPtr returns a pointer to the given int.
func intPtr(i int) *int {
	return &i
}
