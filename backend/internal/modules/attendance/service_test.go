package attendance

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func ctx() context.Context {
	return context.Background()
}

// =========================================================================
// Company Settings Service Tests
// =========================================================================

func TestService_UpsertCompanySetting_Create(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateCompanySettingRequest{
		IsLocationRequired: boolPtr(true),
		MaxDistanceMeter:   intPtr(100),
	}

	resp, err := svc.UpsertCompanySetting(ctx(), req)
	if err != nil {
		t.Fatalf("UpsertCompanySetting failed: %v", err)
	}

	if !resp.IsLocationRequired {
		t.Error("expected IsLocationRequired = true")
	}
	if resp.MaxDistanceMeter == nil || *resp.MaxDistanceMeter != 100 {
		t.Errorf("expected MaxDistanceMeter 100, got %v", resp.MaxDistanceMeter)
	}
	if resp.ID == "" {
		t.Error("expected ID to be set")
	}
}

func TestService_GetCompanySetting_NotFound(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetCompanySetting(ctx())
	if err == nil {
		t.Fatal("expected error when no setting exists")
	}
}

// =========================================================================
// Shift Service Tests
// =========================================================================

func TestService_CreateShift_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateCompanyShiftRequest{
		ShiftName:    "Morning Shift",
		CheckInTime:  "08:00:00",
		CheckOutTime: "17:00:00",
	}

	resp, err := svc.CreateShift(ctx(), req)
	if err != nil {
		t.Fatalf("CreateShift failed: %v", err)
	}

	if resp.ShiftName != "Morning Shift" {
		t.Errorf("expected shift_name 'Morning Shift', got '%s'", resp.ShiftName)
	}
	if resp.CheckInTime != "08:00:00" {
		t.Errorf("expected check_in_time '08:00:00', got '%s'", resp.CheckInTime)
	}
}

func TestService_GetShiftByID_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestShift(repo)

	found, err := svc.GetShiftByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetShiftByID failed: %v", err)
	}

	if found.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), found.ID)
	}
}

func TestService_GetShiftByID_InvalidUUID(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetShiftByID(ctx(), "not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_GetShiftByID_NotFound(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetShiftByID(ctx(), uuidStr())
	if err == nil {
		t.Fatal("expected error for non-existent shift")
	}
}

func TestService_ListShifts_DefaultPagination(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	for i := 0; i < 3; i++ {
		createTestShift(repo)
	}

	resp, err := svc.ListShifts(ctx(), 0, 0)
	if err != nil {
		t.Fatalf("ListShifts failed: %v", err)
	}

	if resp.Page != 1 {
		t.Errorf("expected page 1, got %d", resp.Page)
	}
	if resp.PerPage != 20 {
		t.Errorf("expected per_page 20 (default), got %d", resp.PerPage)
	}
	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}
}

func TestService_UpdateShift_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestShift(repo)

	updated, err := svc.UpdateShift(ctx(), created.ID.String(), UpdateCompanyShiftRequest{
		ShiftName: strPtr("Night Shift"),
	})
	if err != nil {
		t.Fatalf("UpdateShift failed: %v", err)
	}

	if updated.ShiftName != "Night Shift" {
		t.Errorf("expected shift_name 'Night Shift', got '%s'", updated.ShiftName)
	}
}

func TestService_DeleteShift_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestShift(repo)

	if err := svc.DeleteShift(ctx(), created.ID.String()); err != nil {
		t.Fatalf("DeleteShift failed: %v", err)
	}

	_, err := svc.GetShiftByID(ctx(), created.ID.String())
	if err == nil {
		t.Fatal("expected error after deleting shift")
	}
}

// =========================================================================
// Employee Shift Service Tests
// =========================================================================

func TestService_CreateEmployeeShift_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	shift := createTestShift(repo)

	req := CreateEmployeeShiftRequest{
		EmployeeID:        uuidStr(),
		AttendanceShiftID: shift.ID.String(),
		EffectiveDateFrom: "2026-01-01",
	}

	resp, err := svc.CreateEmployeeShift(ctx(), req)
	if err != nil {
		t.Fatalf("CreateEmployeeShift failed: %v", err)
	}

	if resp.EffectiveDateFrom != "2026-01-01" {
		t.Errorf("expected effective_date_from '2026-01-01', got '%s'", resp.EffectiveDateFrom)
	}
}

func TestService_CreateEmployeeShift_InvalidEmployeeID(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	shift := createTestShift(repo)

	req := CreateEmployeeShiftRequest{
		EmployeeID:        "not-a-uuid",
		AttendanceShiftID: shift.ID.String(),
		EffectiveDateFrom: "2026-01-01",
	}

	_, err := svc.CreateEmployeeShift(ctx(), req)
	if err == nil {
		t.Fatal("expected error for invalid employee_id")
	}
}

func TestService_ListEmployeeShifts_ByEmployee(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	shift := createTestShift(repo)
	empID := uuidStr()

	// Non-overlapping date ranges - overlap validation would otherwise
	// reject the 2nd/3rd assignment for the same employee.
	dateRanges := [][2]string{
		{"2026-01-01", "2026-01-31"},
		{"2026-02-01", "2026-02-28"},
		{"2026-03-01", "2026-03-31"},
	}
	for i, dr := range dateRanges {
		to := dr[1]
		req := CreateEmployeeShiftRequest{
			EmployeeID:        empID,
			AttendanceShiftID: shift.ID.String(),
			EffectiveDateFrom: dr[0],
			EffectiveDateTo:   &to,
		}
		if _, err := svc.CreateEmployeeShift(ctx(), req); err != nil {
			t.Fatalf("CreateEmployeeShift iteration %d failed: %v", i, err)
		}
	}

	resp, err := svc.ListEmployeeShifts(ctx(), &empID, 1, 10)
	if err != nil {
		t.Fatalf("ListEmployeeShifts failed: %v", err)
	}

	if resp.Total != 3 {
		t.Errorf("expected 3 employee shifts, got %d", resp.Total)
	}
}

func TestService_UpdateEmployeeShift_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	shift := createTestShift(repo)
	empID := uuid.New()
	created := createTestEmployeeShift(repo, empID, shift.ID)

	isDayOff := true
	updated, err := svc.UpdateEmployeeShift(ctx(), created.ID.String(), UpdateEmployeeShiftRequest{
		IsDayOff: &isDayOff,
	})
	if err != nil {
		t.Fatalf("UpdateEmployeeShift failed: %v", err)
	}

	if updated.IsDayOff == nil || !*updated.IsDayOff {
		t.Error("expected IsDayOff = true")
	}
}

func TestService_CreateEmployeeShift_OverlappingDates_ReturnsError(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	shift := createTestShift(repo)
	empID := uuidStr()

	firstTo := "2026-01-31"
	first := CreateEmployeeShiftRequest{
		EmployeeID:        empID,
		AttendanceShiftID: shift.ID.String(),
		EffectiveDateFrom: "2026-01-01",
		EffectiveDateTo:   &firstTo,
	}
	if _, err := svc.CreateEmployeeShift(ctx(), first); err != nil {
		t.Fatalf("first CreateEmployeeShift failed: %v", err)
	}

	overlapping := CreateEmployeeShiftRequest{
		EmployeeID:        empID,
		AttendanceShiftID: shift.ID.String(),
		EffectiveDateFrom: "2026-01-15",
	}
	if _, err := svc.CreateEmployeeShift(ctx(), overlapping); err == nil {
		t.Fatal("expected error for overlapping shift assignment dates")
	}
}

func TestService_CreateEmployeeShift_EffectiveDateToBeforeFrom_ReturnsError(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	shift := createTestShift(repo)
	to := "2026-01-01"
	req := CreateEmployeeShiftRequest{
		EmployeeID:        uuidStr(),
		AttendanceShiftID: shift.ID.String(),
		EffectiveDateFrom: "2026-01-31",
		EffectiveDateTo:   &to,
	}
	if _, err := svc.CreateEmployeeShift(ctx(), req); err == nil {
		t.Fatal("expected error when effective_date_to is before effective_date_from")
	}
}

func TestService_UpdateEmployeeShift_OverlappingDates_ReturnsError(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	shift := createTestShift(repo)
	empID := uuid.New()

	firstTo := "2026-01-31"
	first := createTestEmployeeShift(repo, empID, shift.ID)
	first.EffectiveDateFrom = "2026-01-01"
	first.EffectiveDateTo = &firstTo
	if err := repo.UpdateEmployeeShift(ctx(), first); err != nil {
		t.Fatalf("failed to seed first shift: %v", err)
	}

	secondTo := "2026-03-31"
	second := createTestEmployeeShift(repo, empID, shift.ID)
	second.EffectiveDateFrom = "2026-03-01"
	second.EffectiveDateTo = &secondTo
	if err := repo.UpdateEmployeeShift(ctx(), second); err != nil {
		t.Fatalf("failed to seed second shift: %v", err)
	}

	newFrom := "2026-01-20"
	if _, err := svc.UpdateEmployeeShift(ctx(), second.ID.String(), UpdateEmployeeShiftRequest{
		EffectiveDateFrom: &newFrom,
	}); err == nil {
		t.Fatal("expected error when updated date range overlaps another assignment")
	}
}

// =========================================================================
// Location Service Tests
// =========================================================================

func TestService_CreateLocation_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateLocationRequest{
		Name:      "Main Office",
		Latitude:  -6.2088,
		Longitude: 106.8456,
	}

	resp, err := svc.CreateLocation(ctx(), req)
	if err != nil {
		t.Fatalf("CreateLocation failed: %v", err)
	}

	if resp.Name != "Main Office" {
		t.Errorf("expected name 'Main Office', got '%s'", resp.Name)
	}
	if resp.RadiusM != 50 {
		t.Errorf("expected default radius 50, got %d", resp.RadiusM)
	}
}

func TestService_CreateLocation_CustomRadius(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	radius := 100
	req := CreateLocationRequest{
		Name:      "Branch Office",
		Latitude:  -6.2,
		Longitude: 106.8,
		RadiusM:   &radius,
	}

	resp, err := svc.CreateLocation(ctx(), req)
	if err != nil {
		t.Fatalf("CreateLocation failed: %v", err)
	}

	if resp.RadiusM != 100 {
		t.Errorf("expected radius 100, got %d", resp.RadiusM)
	}
}

func TestService_GetLocationByID_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestLocation(repo)

	found, err := svc.GetLocationByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetLocationByID failed: %v", err)
	}

	if found.Name != "Main Office" {
		t.Errorf("expected name 'Main Office', got '%s'", found.Name)
	}
}

func TestService_ListLocations_DefaultPagination(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	for i := 0; i < 5; i++ {
		createTestLocation(repo)
	}

	resp, err := svc.ListLocations(ctx(), 0, 0)
	if err != nil {
		t.Fatalf("ListLocations failed: %v", err)
	}

	if resp.Total != 5 {
		t.Errorf("expected total 5, got %d", resp.Total)
	}
}

// =========================================================================
// Event Service Tests
// =========================================================================

func TestService_CreateEvent_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateEventRequest{
		EmployeeID:     uuidStr(),
		EventType:      "CHECKIN",
		EventTimeUTC:   "2026-01-15T00:00:00Z",
		EventTimeLocal: "2026-01-15T07:00:00+07:00",
		Latitude:       -6.2088,
		Longitude:      106.8456,
	}

	resp, err := svc.CreateEvent(ctx(), req)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	if resp.EventType != "CHECKIN" {
		t.Errorf("expected event_type CHECKIN, got '%s'", resp.EventType)
	}
	if resp.ValidationStatus != "PENDING" {
		t.Errorf("expected validation_status PENDING, got '%s'", resp.ValidationStatus)
	}
}

func TestService_CreateEvent_LocationRequired_InsideGeofence_MarksValid(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	required := true
	if _, err := svc.UpsertCompanySetting(ctx(), CreateCompanySettingRequest{IsLocationRequired: &required}); err != nil {
		t.Fatalf("UpsertCompanySetting failed: %v", err)
	}
	loc := &AttendanceLocation{Name: "Office", Latitude: -6.2088, Longitude: 106.8456, RadiusM: 100}
	if err := repo.CreateLocation(ctx(), loc); err != nil {
		t.Fatalf("failed to create location: %v", err)
	}

	req := CreateEventRequest{
		EmployeeID:     uuidStr(),
		EventType:      "CHECKIN",
		EventTimeUTC:   "2026-01-15T00:00:00Z",
		EventTimeLocal: "2026-01-15T07:00:00+07:00",
		Latitude:       -6.2088,
		Longitude:      106.8456,
	}
	resp, err := svc.CreateEvent(ctx(), req)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}
	if resp.ValidationStatus != "VALID" {
		t.Errorf("expected validation_status VALID, got '%s'", resp.ValidationStatus)
	}
	if !resp.IsInGeofence {
		t.Errorf("expected is_in_geofence true, got %v", resp.IsInGeofence)
	}
}

func TestService_CreateEvent_LocationRequired_OutsideGeofence_MarksInvalid(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	required := true
	if _, err := svc.UpsertCompanySetting(ctx(), CreateCompanySettingRequest{IsLocationRequired: &required}); err != nil {
		t.Fatalf("UpsertCompanySetting failed: %v", err)
	}
	loc := &AttendanceLocation{Name: "Office", Latitude: -6.2088, Longitude: 106.8456, RadiusM: 100}
	if err := repo.CreateLocation(ctx(), loc); err != nil {
		t.Fatalf("failed to create location: %v", err)
	}

	req := CreateEventRequest{
		EmployeeID:     uuidStr(),
		EventType:      "CHECKIN",
		EventTimeUTC:   "2026-01-15T00:00:00Z",
		EventTimeLocal: "2026-01-15T07:00:00+07:00",
		Latitude:       -6.3000, // far from the office
		Longitude:      106.9000,
	}
	resp, err := svc.CreateEvent(ctx(), req)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}
	if resp.ValidationStatus != "INVALID" {
		t.Errorf("expected validation_status INVALID, got '%s'", resp.ValidationStatus)
	}
}

func TestService_CreateEvent_FaceRequired_LeavesStatusPending(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	faceRequired := true
	if _, err := svc.UpsertCompanySetting(ctx(), CreateCompanySettingRequest{IsFaceRequired: &faceRequired}); err != nil {
		t.Fatalf("UpsertCompanySetting failed: %v", err)
	}

	req := CreateEventRequest{
		EmployeeID:     uuidStr(),
		EventType:      "CHECKIN",
		EventTimeUTC:   "2026-01-15T00:00:00Z",
		EventTimeLocal: "2026-01-15T07:00:00+07:00",
		Latitude:       -6.2088,
		Longitude:      106.8456,
	}
	resp, err := svc.CreateEvent(ctx(), req)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}
	if resp.ValidationStatus != "PENDING" {
		t.Errorf("expected validation_status PENDING (face check has no implementation to run), got '%s'", resp.ValidationStatus)
	}
}

func TestService_CreateEvent_InvalidTime(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateEventRequest{
		EmployeeID:     uuidStr(),
		EventType:      "CHECKIN",
		EventTimeUTC:   "invalid-time",
		EventTimeLocal: "2026-01-15T07:00:00+07:00",
		Latitude:       -6.2088,
		Longitude:      106.8456,
	}

	_, err := svc.CreateEvent(ctx(), req)
	if err == nil {
		t.Fatal("expected error for invalid time format")
	}
}

func TestService_ListEvents_All(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	empID := uuid.New()
	createTestEvent(repo, empID)
	createTestEvent(repo, empID)

	resp, err := svc.ListEvents(ctx(), nil, 1, 10)
	if err != nil {
		t.Fatalf("ListEvents failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
}

// =========================================================================
// Session Service Tests
// =========================================================================

func TestService_ListSessions_ByEmployee(t *testing.T) {
	svc, _, db, cleanup := newTestService()
	defer cleanup()

	empID := uuid.New()
	empIDStr := empID.String()

	createTestSession(db, empID)       // WorkDate: 2026-01-15
	// Second session needs different work date - create directly
	s2 := &AttendanceSession{
		EmployeeID: empID,
		WorkDate:   "2026-01-16",
		Status:     SessionStatusClosed,
	}
	if err := db.WithContext(context.Background()).Create(s2).Error; err != nil {
		t.Fatalf("failed to create second session: %v", err)
	}

	resp, err := svc.ListSessions(ctx(), &empIDStr, 1, 10)
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2 sessions, got %d", resp.Total)
	}
}

// =========================================================================
// Overtime Request Service Tests
// =========================================================================

func TestService_CreateOvertimeRequest_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateOvertimeRequest{
		EmployeeID:       uuidStr(),
		WorkDate:         "2026-01-15",
		StartTimeLocal:   "2026-01-15T18:00:00+07:00",
		EndTimeLocal:     "2026-01-15T20:00:00+07:00",
		RequestedMinutes: 120,
		Reason:           "Project deadline",
	}

	resp, err := svc.CreateOvertimeRequest(ctx(), req)
	if err != nil {
		t.Fatalf("CreateOvertimeRequest failed: %v", err)
	}

	if resp.Status != "SUBMITTED" {
		t.Errorf("expected status SUBMITTED, got '%s'", resp.Status)
	}
	if resp.RequestedMinutes != 120 {
		t.Errorf("expected 120 minutes, got %d", resp.RequestedMinutes)
	}
	if resp.Reason == nil || *resp.Reason != "Project deadline" {
		t.Errorf("expected reason 'Project deadline', got %v", resp.Reason)
	}
}

func TestService_CreateOvertimeRequest_DefaultStatus(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateOvertimeRequest{
		EmployeeID:       uuidStr(),
		WorkDate:         "2026-01-15",
		StartTimeLocal:   "2026-01-15T18:00:00+07:00",
		EndTimeLocal:     "2026-01-15T20:00:00+07:00",
		RequestedMinutes: 60,
	}

	resp, err := svc.CreateOvertimeRequest(ctx(), req)
	if err != nil {
		t.Fatalf("CreateOvertimeRequest failed: %v", err)
	}

	if resp.Status != "SUBMITTED" {
		t.Errorf("expected default status SUBMITTED, got '%s'", resp.Status)
	}
	if resp.Reason != nil {
		t.Error("expected reason to be nil when not provided")
	}
}

func TestService_GetOvertimeRequestByID_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	empID := uuid.New()
	created := createTestOvertimeRequest(repo, empID)

	found, err := svc.GetOvertimeRequestByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetOvertimeRequestByID failed: %v", err)
	}

	if found.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), found.ID)
	}
}

// =========================================================================
// Exempt Position Service Tests
// =========================================================================

func TestService_CreateExemptPosition_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateExemptPositionRequest{
		OrganizationID: uuidStr(),
	}

	resp, err := svc.CreateExemptPosition(ctx(), req)
	if err != nil {
		t.Fatalf("CreateExemptPosition failed: %v", err)
	}

	if !resp.IsExempt {
		t.Error("expected default IsExempt = true")
	}
}

func TestService_CreateExemptPosition_NotExempt(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	notExempt := false
	req := CreateExemptPositionRequest{
		OrganizationID: uuidStr(),
		IsExempt:       &notExempt,
	}

	resp, err := svc.CreateExemptPosition(ctx(), req)
	if err != nil {
		t.Fatalf("CreateExemptPosition failed: %v", err)
	}

	if resp.IsExempt {
		t.Error("expected IsExempt = false")
	}
}

func TestService_ListExemptPositions_Pagination(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	for i := 0; i < 3; i++ {
		createTestExemptPosition(repo, uuid.New())
	}

	resp, err := svc.ListExemptPositions(ctx(), 1, 10)
	if err != nil {
		t.Fatalf("ListExemptPositions failed: %v", err)
	}

	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}
}

func TestService_UpdateExemptPosition_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	orgID := uuid.New()
	created := createTestExemptPosition(repo, orgID)

	isExempt := false
	updated, err := svc.UpdateExemptPosition(ctx(), created.ID.String(), UpdateExemptPositionRequest{
		IsExempt: &isExempt,
	})
	if err != nil {
		t.Fatalf("UpdateExemptPosition failed: %v", err)
	}

	if updated.IsExempt {
		t.Error("expected IsExempt = false after update")
	}
}

func TestService_DeleteExemptPosition_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	orgID := uuid.New()
	created := createTestExemptPosition(repo, orgID)

	if err := svc.DeleteExemptPosition(ctx(), created.ID.String()); err != nil {
		t.Fatalf("DeleteExemptPosition failed: %v", err)
	}
}

// boolPtr is defined in helpers_test.go
// strPtr is defined in service.go
