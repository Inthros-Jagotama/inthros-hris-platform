package attendance

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// =========================================================================
// Company Settings Repository Tests
// =========================================================================

func TestRepo_UpsertCompanySetting_Create(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	s := &AttendanceCompanySetting{
		IsLocationRequired: true,
		MaxDistanceMeter:   intPtr(100),
	}

	if err := repo.UpsertCompanySetting(ctx, s); err != nil {
		t.Fatalf("UpsertCompanySetting failed: %v", err)
	}

	if s.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_UpsertCompanySetting_UpdatePreservesCreatedAt(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	first := &AttendanceCompanySetting{IsLocationRequired: true}
	if err := repo.UpsertCompanySetting(ctx, first); err != nil {
		t.Fatalf("initial UpsertCompanySetting failed: %v", err)
	}

	before, err := repo.FindCompanySetting(ctx)
	if err != nil {
		t.Fatalf("FindCompanySetting failed: %v", err)
	}
	if before.CreatedAt.IsZero() {
		t.Fatal("expected created_at to be set after create")
	}

	// Simulasikan request update: objek baru tanpa timestamps (CreatedAt zero).
	// Tanpa perbaikan, Save menimpa created_at dengan zero time (MySQL menolak
	// '0000-00-00' pada strict mode).
	updated := &AttendanceCompanySetting{IsOvertimeEnabled: true, AllowCheckinOnDayOff: true}
	if err := repo.UpsertCompanySetting(ctx, updated); err != nil {
		t.Fatalf("update UpsertCompanySetting failed: %v", err)
	}

	after, err := repo.FindCompanySetting(ctx)
	if err != nil {
		t.Fatalf("FindCompanySetting after update failed: %v", err)
	}
	if !after.CreatedAt.Equal(before.CreatedAt) {
		t.Errorf("expected created_at preserved (%v), got %v", before.CreatedAt, after.CreatedAt)
	}
	if !after.IsOvertimeEnabled {
		t.Error("expected IsOvertimeEnabled = true after update")
	}
}

func TestRepo_FindCompanySetting_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	s := &AttendanceCompanySetting{IsLocationRequired: true}
	if err := repo.UpsertCompanySetting(ctx, s); err != nil {
		t.Fatalf("UpsertCompanySetting failed: %v", err)
	}

	found, err := repo.FindCompanySetting(ctx)
	if err != nil {
		t.Fatalf("FindCompanySetting failed: %v", err)
	}

	if !found.IsLocationRequired {
		t.Error("expected IsLocationRequired = true")
	}
}

func TestRepo_FindCompanySetting_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	_, err := repo.FindCompanySetting(ctx)
	if err == nil {
		t.Fatal("expected error when no setting exists")
	}
}

// =========================================================================
// Company Shift Repository Tests
// =========================================================================

func TestRepo_CreateShift_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	s := &AttendanceCompanyShift{
		ShiftName:    "Morning Shift",
		CheckInTime:  "08:00:00",
		CheckOutTime: "17:00:00",
	}

	if err := repo.CreateShift(ctx, s); err != nil {
		t.Fatalf("CreateShift failed: %v", err)
	}

	if s.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindShiftByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestShift(repo)

	found, err := repo.FindShiftByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindShiftByID failed: %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("expected ID '%s', got '%s'", created.ID, found.ID)
	}
	if found.ShiftName != "Morning Shift" {
		t.Errorf("expected shift_name 'Morning Shift', got '%s'", found.ShiftName)
	}
}

func TestRepo_FindShiftByID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	_, err := repo.FindShiftByID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent shift")
	}
}

func TestRepo_ListShifts_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		createTestShift(repo)
	}

	shifts, total, err := repo.ListShifts(ctx, 1, 3)
	if err != nil {
		t.Fatalf("ListShifts failed: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(shifts) != 3 {
		t.Errorf("expected 3 shifts (page 1 of 2), got %d", len(shifts))
	}
}

func TestRepo_UpdateShift_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestShift(repo)

	created.ShiftName = "Night Shift"
	if err := repo.UpdateShift(ctx, created); err != nil {
		t.Fatalf("UpdateShift failed: %v", err)
	}

	found, _ := repo.FindShiftByID(ctx, created.ID)
	if found.ShiftName != "Night Shift" {
		t.Errorf("expected shift_name 'Night Shift', got '%s'", found.ShiftName)
	}
}

func TestRepo_DeleteShift_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestShift(repo)

	if err := repo.DeleteShift(ctx, created.ID); err != nil {
		t.Fatalf("DeleteShift failed: %v", err)
	}

	_, err := repo.FindShiftByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting shift")
	}
}

func TestRepo_DeleteShift_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	err := repo.DeleteShift(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error when deleting non-existent shift")
	}
}

// =========================================================================
// Employee Shift Repository Tests
// =========================================================================

func TestRepo_CreateEmployeeShift_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	shift := createTestShift(repo)
	empID := uuid.New()

	es := &AttendanceEmployeeShift{
		EmployeeID:        empID,
		AttendanceShiftID: shift.ID,
		EffectiveDateFrom: "2026-01-01",
	}

	if err := repo.CreateEmployeeShift(ctx, es); err != nil {
		t.Fatalf("CreateEmployeeShift failed: %v", err)
	}

	if es.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindEmployeeShiftByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	shift := createTestShift(repo)
	empID := uuid.New()
	created := createTestEmployeeShift(repo, empID, shift.ID)

	found, err := repo.FindEmployeeShiftByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindEmployeeShiftByID failed: %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("expected ID '%s', got '%s'", created.ID, found.ID)
	}
}

func TestRepo_ListEmployeeShifts_ByEmployee(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	shift := createTestShift(repo)
	empID := uuid.New()
	otherEmp := uuid.New()

	createTestEmployeeShift(repo, empID, shift.ID)
	createTestEmployeeShift(repo, empID, shift.ID)
	createTestEmployeeShift(repo, otherEmp, shift.ID)

	// List by specific employee
	shifts, total, err := repo.ListEmployeeShifts(ctx, &empID, 1, 10)
	if err != nil {
		t.Fatalf("ListEmployeeShifts failed: %v", err)
	}

	if total != 2 {
		t.Errorf("expected 2 shifts for employee, got %d", total)
	}
	if len(shifts) != 2 {
		t.Errorf("expected 2 shifts, got %d", len(shifts))
	}
}

func TestRepo_UpdateEmployeeShift_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	shift := createTestShift(repo)
	empID := uuid.New()
	created := createTestEmployeeShift(repo, empID, shift.ID)

	isDayOff := true
	created.IsDayOff = &isDayOff
	if err := repo.UpdateEmployeeShift(ctx, created); err != nil {
		t.Fatalf("UpdateEmployeeShift failed: %v", err)
	}

	found, _ := repo.FindEmployeeShiftByID(ctx, created.ID)
	if found.IsDayOff == nil || !*found.IsDayOff {
		t.Error("expected IsDayOff = true")
	}
}

func TestRepo_DeleteEmployeeShift_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	shift := createTestShift(repo)
	empID := uuid.New()
	created := createTestEmployeeShift(repo, empID, shift.ID)

	if err := repo.DeleteEmployeeShift(ctx, created.ID); err != nil {
		t.Fatalf("DeleteEmployeeShift failed: %v", err)
	}
}

// =========================================================================
// Location Repository Tests
// =========================================================================

func TestRepo_CreateLocation_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	l := &AttendanceLocation{
		Name:      "Head Office",
		Latitude:  -6.2088,
		Longitude: 106.8456,
		RadiusM:   50,
	}

	if err := repo.CreateLocation(ctx, l); err != nil {
		t.Fatalf("CreateLocation failed: %v", err)
	}

	if l.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindLocationByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestLocation(repo)

	found, err := repo.FindLocationByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindLocationByID failed: %v", err)
	}

	if found.Name != "Main Office" {
		t.Errorf("expected name 'Main Office', got '%s'", found.Name)
	}
}

func TestRepo_ListLocations_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		createTestLocation(repo)
	}

	locs, total, err := repo.ListLocations(ctx, 1, 2)
	if err != nil {
		t.Fatalf("ListLocations failed: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(locs) != 2 {
		t.Errorf("expected 2 locations (page 1 of 3), got %d", len(locs))
	}
}

func TestRepo_DeleteLocation_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	created := createTestLocation(repo)

	if err := repo.DeleteLocation(ctx, created.ID); err != nil {
		t.Fatalf("DeleteLocation failed: %v", err)
	}
}

// =========================================================================
// Event Repository Tests
// =========================================================================

func TestRepo_CreateEvent_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	empID := uuid.New()
	e := &AttendanceEvent{
		EmployeeID:       empID,
		EventType:        EventTypeCheckIn,
		EventTimeUTC:     parseTime("2026-01-15T00:00:00Z"),
		EventTimeLocal:   parseTime("2026-01-15T07:00:00+07:00"),
		Latitude:         -6.2088,
		Longitude:        106.8456,
		ValidationStatus: ValidationPending,
	}

	if err := repo.CreateEvent(ctx, e); err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	if e.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindEventByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	empID := uuid.New()
	created := createTestEvent(repo, empID)

	found, err := repo.FindEventByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindEventByID failed: %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("expected ID '%s', got '%s'", created.ID, found.ID)
	}
	if found.EventType != EventTypeCheckIn {
		t.Errorf("expected event_type CHECKIN, got '%s'", found.EventType)
	}
}

func TestRepo_ListEvents_ByEmployee(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	empID := uuid.New()
	otherEmp := uuid.New()

	createTestEvent(repo, empID)
	createTestEvent(repo, empID)
	createTestEvent(repo, otherEmp)

	events, total, err := repo.ListEvents(ctx, &empID, 1, 10)
	if err != nil {
		t.Fatalf("ListEvents failed: %v", err)
	}

	if total != 2 {
		t.Errorf("expected 2 events for employee, got %d", total)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
}

// =========================================================================
// Session Repository Tests
// =========================================================================

func TestRepo_FindSessionByEmployeeAndDate_Success(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	empID := uuid.New()
	createTestSession(db, empID)

	// Use FindSessions to check
	sessions, total, err := repo.ListSessions(ctx, &empID, 1, 10)
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}

	if total != 1 {
		t.Errorf("expected 1 session, got %d", total)
	}
	if len(sessions) != 1 {
		t.Errorf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].WorkDate != "2026-01-15" {
		t.Errorf("expected work_date '2026-01-15', got '%s'", sessions[0].WorkDate)
	}
}

func TestRepo_ListSessions_Empty(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	sessions, total, err := repo.ListSessions(ctx, nil, 1, 10)
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}

	if total != 0 {
		t.Errorf("expected 0 sessions, got %d", total)
	}
	if len(sessions) != 0 {
		t.Errorf("expected 0 sessions, got %d", len(sessions))
	}
}

// =========================================================================
// Overtime Request Repository Tests
// =========================================================================

func TestRepo_CreateOvertimeRequest_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	empID := uuid.New()
	created := createTestOvertimeRequest(repo, empID)

	if created.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindOvertimeRequestByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	empID := uuid.New()
	created := createTestOvertimeRequest(repo, empID)

	found, err := repo.FindOvertimeRequestByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindOvertimeRequestByID failed: %v", err)
	}

	if found.RequestedMinutes != 120 {
		t.Errorf("expected 120 minutes, got %d", found.RequestedMinutes)
	}
}

func TestRepo_ListOvertimeRequests_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	empID := uuid.New()
	createTestOvertimeRequest(repo, empID)
	createTestOvertimeRequest(repo, empID)

	requests, total, err := repo.ListOvertimeRequests(ctx, &empID, 1, 10)
	if err != nil {
		t.Fatalf("ListOvertimeRequests failed: %v", err)
	}

	if total != 2 {
		t.Errorf("expected 2 requests, got %d", total)
	}
	if len(requests) != 2 {
		t.Errorf("expected 2 requests, got %d", len(requests))
	}
}

// =========================================================================
// Exempt Position Repository Tests
// =========================================================================

func TestRepo_CreateExemptPosition_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	orgID := uuid.New()
	p := &AttendanceExemptPosition{
		OrganizationID: orgID,
		IsExempt:       true,
	}

	if err := repo.CreateExemptPosition(ctx, p); err != nil {
		t.Fatalf("CreateExemptPosition failed: %v", err)
	}

	if p.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindExemptPositionByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	orgID := uuid.New()
	created := createTestExemptPosition(repo, orgID)

	found, err := repo.FindExemptPositionByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindExemptPositionByID failed: %v", err)
	}

	if !found.IsExempt {
		t.Error("expected IsExempt = true")
	}
}

func TestRepo_FindExemptPositionByOrgID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	orgID := uuid.New()
	createTestExemptPosition(repo, orgID)

	found, err := repo.FindExemptPositionByOrgID(ctx, orgID)
	if err != nil {
		t.Fatalf("FindExemptPositionByOrgID failed: %v", err)
	}

	if found.OrganizationID != orgID {
		t.Errorf("expected org_id '%s', got '%s'", orgID, found.OrganizationID)
	}
}

func TestRepo_ListExemptPositions_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		createTestExemptPosition(repo, uuid.New())
	}

	positions, total, err := repo.ListExemptPositions(ctx, 1, 10)
	if err != nil {
		t.Fatalf("ListExemptPositions failed: %v", err)
	}

	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(positions) != 3 {
		t.Errorf("expected 3 positions, got %d", len(positions))
	}
}

func TestRepo_UpdateExemptPosition_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	orgID := uuid.New()
	created := createTestExemptPosition(repo, orgID)

	created.IsExempt = false
	if err := repo.UpdateExemptPosition(ctx, created); err != nil {
		t.Fatalf("UpdateExemptPosition failed: %v", err)
	}

	found, _ := repo.FindExemptPositionByID(ctx, created.ID)
	if found.IsExempt {
		t.Error("expected IsExempt = false after update")
	}
}

func TestRepo_DeleteExemptPosition_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	orgID := uuid.New()
	created := createTestExemptPosition(repo, orgID)

	if err := repo.DeleteExemptPosition(ctx, created.ID); err != nil {
		t.Fatalf("DeleteExemptPosition failed: %v", err)
	}
}

func TestRepo_FindOrganizationIDByEmployeeID_ReturnsCurrentOrg(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	employeeID := uuid.New()
	orgID := uuid.New()
	seedOrgWithZone(db, orgID, nil, "Test Org")
	seedEmployment(db, employeeID, orgID)

	got, err := repo.FindOrganizationIDByEmployeeID(ctx, employeeID)
	if err != nil {
		t.Fatalf("FindOrganizationIDByEmployeeID failed: %v", err)
	}
	if got == nil {
		t.Fatal("expected organization id, got nil")
	}
	if *got != orgID {
		t.Errorf("got %s, want %s", got, orgID)
	}
}

func TestRepo_FindOrganizationIDByEmployeeID_NoEmployment_ReturnsNil(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	got, err := repo.FindOrganizationIDByEmployeeID(ctx, uuid.New())
	if err != nil {
		t.Fatalf("FindOrganizationIDByEmployeeID failed: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil organization id, got %s", got)
	}
}

// =========================================================================
// (intPtr is defined in helpers_test.go)
// =========================================================================
