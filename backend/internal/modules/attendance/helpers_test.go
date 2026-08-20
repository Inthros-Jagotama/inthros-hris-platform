package attendance

import (
	"context"
	"fmt"
	"testing"
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
		&BusinessTravel{},
		&testEmployeeAccount{},
	); err != nil {
		panic(fmt.Sprintf("failed to migrate test db: %v", err))
	}

	// Raw tables owned by other modules (employee/organization) that
	// GetEmployeeInfoByIDs queries directly via db.Table(...) — attendance
	// must not import those packages (circular dependency risk). Minimal
	// schema covering only the columns the raw query touches, same pattern
	// as approval/helpers_test.go's setupTestDB.
	rawTables := []string{
		`CREATE TABLE IF NOT EXISTS employees (
			id CHAR(36) PRIMARY KEY,
			employee_id VARCHAR(50) NOT NULL DEFAULT '',
			name VARCHAR(255) NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS organizations (
			id CHAR(36) PRIMARY KEY,
			nomenclature VARCHAR(255) NOT NULL DEFAULT '',
			parent_id CHAR(36) NULL,
			zone_id CHAR(36) NULL,
			deleted_at DATETIME NULL
		)`,
		`CREATE TABLE IF NOT EXISTS employments (
			id CHAR(36) PRIMARY KEY,
			employee_id CHAR(36) NOT NULL,
			organization_id CHAR(36) NOT NULL,
			effective_date DATE NOT NULL,
			effective_end_date DATE NULL
		)`,
		// Raw zones table (owned by module `setting`) — used only by
		// ResolveOrganizationTimezone's LEFT JOIN. Minimal schema covering
		// the column touched (timezone).
		`CREATE TABLE IF NOT EXISTS zones (
			id CHAR(36) PRIMARY KEY,
			timezone VARCHAR(64) NULL,
			deleted_at DATETIME NULL
		)`,
	}
	for _, stmt := range rawTables {
		if err := db.Exec(stmt).Error; err != nil {
			panic(fmt.Sprintf("failed to create raw test table: %v", err))
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

// seedEmployeeOrg inserts a minimal employees + organizations + employments
// row set so GetEmployeeInfoByIDs can resolve employeeID -> (name, org name).
func seedEmployeeOrg(db *gorm.DB, employeeID uuid.UUID, name string, orgID uuid.UUID, orgName string) {
	if err := db.Exec("INSERT INTO employees (id, name) VALUES (?, ?)", employeeID.String(), name).Error; err != nil {
		panic(fmt.Sprintf("failed to seed employee: %v", err))
	}
	if err := db.Exec("INSERT INTO organizations (id, nomenclature) VALUES (?, ?)", orgID.String(), orgName).Error; err != nil {
		panic(fmt.Sprintf("failed to seed organization: %v", err))
	}
	if err := db.Exec(
		"INSERT INTO employments (id, employee_id, organization_id, effective_date, effective_end_date) VALUES (?, ?, ?, ?, NULL)",
		uuid.New().String(), employeeID.String(), orgID.String(), "2026-01-01",
	).Error; err != nil {
		panic(fmt.Sprintf("failed to seed employment: %v", err))
	}
}

// seedUserEmployee inserts a user_id -> employee_id mapping in
// employee_accounts so FindEmployeeIDByUserID can resolve the logged-in
// user's employee (used by the §32b ownership guard tests).
func seedUserEmployee(db *gorm.DB, userID, employeeID uuid.UUID) {
	createTestEmployeeAccount(db, employeeID, userID)
}

// seedOrg inserts an organization row (parentID may be nil for a root org)
// into the raw organizations table — used by the assignable-employee
// hierarchy tests.
func seedOrg(db *gorm.DB, id uuid.UUID, parentID *uuid.UUID, name string) {
	sql := "INSERT INTO organizations (id, nomenclature, parent_id) VALUES (?, ?, NULL)"
	args := []interface{}{id.String(), name}
	if parentID != nil {
		sql = "INSERT INTO organizations (id, nomenclature, parent_id) VALUES (?, ?, ?)"
		args = append(args, parentID.String())
	}
	if err := db.Exec(sql, args...).Error; err != nil {
		panic(fmt.Sprintf("failed to seed organization: %v", err))
	}
}

// setupTestPlatformDB creates a separate in-memory SQLite database
// containing just a minimal companies table (id, timezone), mirroring the
// platform DB that Repository.getCompanyTimezone reads companies.timezone
// from. Mirrors setting/helpers_test.go's setupTestPlatformDB, but avoids
// importing the company package (raw table, same convention as this file's
// other cross-module fixtures) to sidestep any import-cycle risk.
func setupTestPlatformDB() (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&platform=1"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to open test platform db: %v", err))
	}
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS companies (
		id CHAR(36) PRIMARY KEY,
		timezone VARCHAR(64) NOT NULL DEFAULT ''
	)`).Error; err != nil {
		panic(fmt.Sprintf("failed to create test companies table: %v", err))
	}
	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
	return db, cleanup
}

// newTestRepository creates a Repository wired with both an in-memory
// tenant DB and an in-memory platform DB (with a seeded company row keyed
// by companyID, timezone companyTz), and a ctx carrying that company ID via
// authctx (mirrors how middleware.TenantRequired populates request
// context). Used by ResolveOrganizationTimezone tests.
func newTestRepository(t *testing.T, companyID uuid.UUID, companyTz string) (*Repository, *gorm.DB, context.Context) {
	t.Helper()
	db, dbResolver, cleanup := setupTestDB()
	platformDB, platformCleanup := setupTestPlatformDB()
	t.Cleanup(func() {
		cleanup()
		platformCleanup()
	})
	if err := platformDB.Exec("INSERT INTO companies (id, timezone) VALUES (?, ?)", companyID.String(), companyTz).Error; err != nil {
		t.Fatalf("failed to seed test company: %v", err)
	}
	repo := NewRepositoryWithPlatformDB(dbResolver, platformDB)
	ctx := context.WithValue(context.Background(), "company_id", companyID.String())
	return repo, db, ctx
}

// seedZone inserts a zone row (raw table) with the given timezone override.
func seedZone(db *gorm.DB, id uuid.UUID, tz string) {
	if err := db.Exec("INSERT INTO zones (id, timezone) VALUES (?, ?)", id.String(), tz).Error; err != nil {
		panic(fmt.Sprintf("failed to seed zone: %v", err))
	}
}

// seedOrgWithZone inserts an organization row with a zone_id reference.
func seedOrgWithZone(db *gorm.DB, id uuid.UUID, zoneID *uuid.UUID, name string) {
	var zid interface{}
	if zoneID != nil {
		zid = zoneID.String()
	}
	if err := db.Exec("INSERT INTO organizations (id, nomenclature, parent_id, zone_id) VALUES (?, ?, NULL, ?)", id.String(), name, zid).Error; err != nil {
		panic(fmt.Sprintf("failed to seed organization with zone: %v", err))
	}
}

// seedEmployment inserts an active employment (effective_end_date NULL) for
// employeeID in orgID — used by the assignable-employee hierarchy tests.
func seedEmployment(db *gorm.DB, employeeID, orgID uuid.UUID) {
	if err := db.Exec(
		"INSERT INTO employments (id, employee_id, organization_id, effective_date, effective_end_date) VALUES (?, ?, ?, ?, NULL)",
		uuid.New().String(), employeeID.String(), orgID.String(), "2026-01-01",
	).Error; err != nil {
		panic(fmt.Sprintf("failed to seed employment: %v", err))
	}
}

// seedEmployeeRow inserts a minimal employees row so raw JOINs against the
// employees table (e.g. FindEmployeesByOrganizationIDs) can resolve names.
func seedEmployeeRow(db *gorm.DB, employeeID uuid.UUID, name string) {
	code := "EMP-" + employeeID.String()[:8]
	if err := db.Exec("INSERT INTO employees (id, employee_id, name) VALUES (?, ?, ?)", employeeID.String(), code, name).Error; err != nil {
		panic(fmt.Sprintf("failed to seed employee row: %v", err))
	}
}

// ctxWithUser returns a context carrying the given user_id, mirroring what
// the AuthJWT middleware injects so authctx.GetUserID can resolve it.
func ctxWithUser(userID uuid.UUID) context.Context {
	return context.WithValue(context.Background(), "user_id", userID.String())
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
