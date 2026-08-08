package attendance

import (
	"testing"

	"github.com/google/uuid"
)

func TestService_CreateEvent_SessionGeneration_OnTimeCheckinAndCheckout(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	shift := createTestShift(repo) // 08:00 - 17:00
	empID := uuid.New()
	createTestEmployeeShift(repo, empID, shift.ID)

	checkin := CreateEventRequest{
		EmployeeID:     empID.String(),
		EventType:      "CHECKIN",
		EventTimeUTC:   "2026-01-15T01:00:00Z",
		EventTimeLocal: "2026-01-15T08:00:00+07:00",
		Latitude:       -6.2088,
		Longitude:      106.8456,
	}
	if _, err := svc.CreateEvent(ctx(), checkin); err != nil {
		t.Fatalf("checkin CreateEvent failed: %v", err)
	}
	checkout := checkin
	checkout.EventType = "CHECKOUT"
	checkout.EventTimeUTC = "2026-01-15T10:00:00Z"
	checkout.EventTimeLocal = "2026-01-15T17:00:00+07:00"
	if _, err := svc.CreateEvent(ctx(), checkout); err != nil {
		t.Fatalf("checkout CreateEvent failed: %v", err)
	}

	session, err := repo.FindSessionByEmployeeAndDate(ctx(), empID, "2026-01-15")
	if err != nil {
		t.Fatalf("expected session to be generated: %v", err)
	}
	if session.Status != SessionStatusClosed {
		t.Errorf("expected status CLOSED, got '%s'", session.Status)
	}
	if session.LatenessMinutes != 0 {
		t.Errorf("expected 0 lateness minutes, got %d", session.LatenessMinutes)
	}
	if session.EarlyLeaveMinutes != 0 {
		t.Errorf("expected 0 early leave minutes, got %d", session.EarlyLeaveMinutes)
	}
	if session.WorkMinutes != 540 {
		t.Errorf("expected 540 work minutes (9h), got %d", session.WorkMinutes)
	}
}

func TestService_CreateEvent_SessionGeneration_LateWithTolerance(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	shift := createTestShift(repo) // 08:00 - 17:00
	empID := uuid.New()
	createTestEmployeeShift(repo, empID, shift.ID)

	tolerance := 10
	if _, err := svc.UpsertCompanySetting(ctx(), CreateCompanySettingRequest{LateToleranceMinutes: &tolerance}); err != nil {
		t.Fatalf("UpsertCompanySetting failed: %v", err)
	}

	checkin := CreateEventRequest{
		EmployeeID:     empID.String(),
		EventType:      "CHECKIN",
		EventTimeUTC:   "2026-01-15T01:15:00Z",
		EventTimeLocal: "2026-01-15T08:15:00+07:00", // 15 minutes late
		Latitude:       -6.2088,
		Longitude:      106.8456,
	}
	if _, err := svc.CreateEvent(ctx(), checkin); err != nil {
		t.Fatalf("checkin CreateEvent failed: %v", err)
	}

	session, err := repo.FindSessionByEmployeeAndDate(ctx(), empID, "2026-01-15")
	if err != nil {
		t.Fatalf("expected session to be generated: %v", err)
	}
	if session.Status != SessionStatusMissingCheckOut {
		t.Errorf("expected status MISSING_CHECKOUT, got '%s'", session.Status)
	}
	if session.LatenessMinutes != 5 {
		t.Errorf("expected 5 lateness minutes (15 - 10 tolerance), got %d", session.LatenessMinutes)
	}
}

func TestService_CreateEvent_SessionGeneration_CrossMidnight(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	nightShift := &AttendanceCompanyShift{
		ShiftName:       "Night Shift",
		CheckInTime:     "22:00:00",
		CheckOutTime:    "06:00:00",
		IsCrossMidnight: true,
	}
	if err := repo.CreateShift(ctx(), nightShift); err != nil {
		t.Fatalf("failed to create night shift: %v", err)
	}
	empID := uuid.New()
	createTestEmployeeShift(repo, empID, nightShift.ID)

	checkin := CreateEventRequest{
		EmployeeID:     empID.String(),
		EventType:      "CHECKIN",
		EventTimeUTC:   "2026-01-15T15:00:00Z",
		EventTimeLocal: "2026-01-15T22:00:00+07:00",
		Latitude:       -6.2088,
		Longitude:      106.8456,
	}
	if _, err := svc.CreateEvent(ctx(), checkin); err != nil {
		t.Fatalf("checkin CreateEvent failed: %v", err)
	}
	checkout := checkin
	checkout.EventType = "CHECKOUT"
	checkout.EventTimeUTC = "2026-01-15T23:00:00Z"
	checkout.EventTimeLocal = "2026-01-16T06:00:00+07:00"
	if _, err := svc.CreateEvent(ctx(), checkout); err != nil {
		t.Fatalf("checkout CreateEvent failed: %v", err)
	}

	// Session should be attributed to the CHECKIN's local date (15th), not
	// the CHECKOUT's (16th), per §24.
	session, err := repo.FindSessionByEmployeeAndDate(ctx(), empID, "2026-01-15")
	if err != nil {
		t.Fatalf("expected session on the check-in's date: %v", err)
	}
	if session.Status != SessionStatusClosed {
		t.Errorf("expected status CLOSED, got '%s'", session.Status)
	}
	if session.WorkMinutes != 480 {
		t.Errorf("expected 480 work minutes (8h), got %d", session.WorkMinutes)
	}

	if _, err := repo.FindSessionByEmployeeAndDate(ctx(), empID, "2026-01-16"); err == nil {
		t.Error("did not expect a separate session on the check-out's date")
	}
}

func TestService_CreateEvent_SessionGeneration_DayOff(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	// Setting: check-in pada hari libur DIBLOKIR.
	if err := repo.UpsertCompanySetting(ctx(), &AttendanceCompanySetting{AllowCheckinOnDayOff: false}); err != nil {
		t.Fatalf("failed to create company setting: %v", err)
	}

	shift := createTestShift(repo)
	empID := uuid.New()
	dayOff := true
	es := &AttendanceEmployeeShift{
		EmployeeID:        empID,
		AttendanceShiftID: shift.ID,
		EffectiveDateFrom: "2026-01-01",
		IsDayOff:          &dayOff,
	}
	if err := repo.CreateEmployeeShift(ctx(), es); err != nil {
		t.Fatalf("failed to create day-off employee shift: %v", err)
	}

	req := CreateEventRequest{
		EmployeeID:     empID.String(),
		EventType:      "CHECKIN",
		EventTimeUTC:   "2026-01-15T01:00:00Z",
		EventTimeLocal: "2026-01-15T08:00:00+07:00",
		Latitude:       -6.2088,
		Longitude:      106.8456,
	}
	if _, err := svc.CreateEvent(ctx(), req); err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	session, err := repo.FindSessionByEmployeeAndDate(ctx(), empID, "2026-01-15")
	if err != nil {
		t.Fatalf("expected session to be generated: %v", err)
	}
	if session.Status != SessionStatusDayOff {
		t.Errorf("expected status DAY_OFF, got '%s'", session.Status)
	}
}

func TestService_CreateEvent_SessionGeneration_DayOff_AllowCheckin(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	// Setting: check-in pada hari libur DIIZINKAN (default migration 075).
	if err := repo.UpsertCompanySetting(ctx(), &AttendanceCompanySetting{AllowCheckinOnDayOff: true}); err != nil {
		t.Fatalf("failed to create company setting: %v", err)
	}

	shift := createTestShift(repo)
	empID := uuid.New()
	dayOff := true
	es := &AttendanceEmployeeShift{
		EmployeeID:        empID,
		AttendanceShiftID: shift.ID,
		EffectiveDateFrom: "2026-01-01",
		IsDayOff:          &dayOff,
	}
	if err := repo.CreateEmployeeShift(ctx(), es); err != nil {
		t.Fatalf("failed to create day-off employee shift: %v", err)
	}

	req := CreateEventRequest{
		EmployeeID:     empID.String(),
		EventType:      "CHECKIN",
		EventTimeUTC:   "2026-01-15T01:00:00Z",
		EventTimeLocal: "2026-01-15T08:00:00+07:00",
		Latitude:       -6.2088,
		Longitude:      106.8456,
	}
	if _, err := svc.CreateEvent(ctx(), req); err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	session, err := repo.FindSessionByEmployeeAndDate(ctx(), empID, "2026-01-15")
	if err != nil {
		t.Fatalf("expected session to be generated: %v", err)
	}
	// Hari libur tapi diizinkan check-in → status mencerminkan event,
	// sehingga tombol Check Out muncul di halaman absensi.
	if session.Status != SessionStatusMissingCheckOut {
		t.Errorf("expected status MISSING_CHECKOUT, got '%s'", session.Status)
	}
	if session.CheckinEventID == nil {
		t.Error("expected checkin_event_id to be linked to the session")
	}
}

func TestService_ApplyApprovedLeave_NoExistingSession_MarksLeave(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	empID := uuid.New()
	leaveReqID := uuid.New()
	if err := svc.ApplyApprovedLeave(ctx(), empID, "2026-01-20", leaveReqID, 1.0); err != nil {
		t.Fatalf("ApplyApprovedLeave failed: %v", err)
	}

	session, err := repo.FindSessionByEmployeeAndDate(ctx(), empID, "2026-01-20")
	if err != nil {
		t.Fatalf("expected session to be created: %v", err)
	}
	if session.Status != SessionStatusLeave {
		t.Errorf("expected status LEAVE, got '%s'", session.Status)
	}
	if session.LeaveRequestID == nil || *session.LeaveRequestID != leaveReqID {
		t.Errorf("expected leave_request_id %s, got %v", leaveReqID, session.LeaveRequestID)
	}
	if session.LeaveFraction == nil || *session.LeaveFraction != 1.0 {
		t.Errorf("expected leave_fraction 1.0, got %v", session.LeaveFraction)
	}
}

func TestService_ApplyApprovedLeave_ClosedSession_DoesNotOverwriteStatus(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	shift := createTestShift(repo)
	empID := uuid.New()
	createTestEmployeeShift(repo, empID, shift.ID)

	checkin := CreateEventRequest{
		EmployeeID:     empID.String(),
		EventType:      "CHECKIN",
		EventTimeUTC:   "2026-01-15T01:00:00Z",
		EventTimeLocal: "2026-01-15T08:00:00+07:00",
		Latitude:       -6.2088,
		Longitude:      106.8456,
	}
	if _, err := svc.CreateEvent(ctx(), checkin); err != nil {
		t.Fatalf("checkin CreateEvent failed: %v", err)
	}
	checkout := checkin
	checkout.EventType = "CHECKOUT"
	checkout.EventTimeUTC = "2026-01-15T05:00:00Z"
	checkout.EventTimeLocal = "2026-01-15T12:00:00+07:00" // half day
	if _, err := svc.CreateEvent(ctx(), checkout); err != nil {
		t.Fatalf("checkout CreateEvent failed: %v", err)
	}

	leaveReqID := uuid.New()
	if err := svc.ApplyApprovedLeave(ctx(), empID, "2026-01-15", leaveReqID, 0.5); err != nil {
		t.Fatalf("ApplyApprovedLeave failed: %v", err)
	}

	session, err := repo.FindSessionByEmployeeAndDate(ctx(), empID, "2026-01-15")
	if err != nil {
		t.Fatalf("expected session: %v", err)
	}
	if session.Status != SessionStatusClosed {
		t.Errorf("expected status to remain CLOSED (real attendance happened), got '%s'", session.Status)
	}
	if session.LeaveFraction == nil || *session.LeaveFraction != 0.5 {
		t.Errorf("expected leave_fraction 0.5 to still be recorded, got %v", session.LeaveFraction)
	}
}

func TestService_CreateEvent_SessionGeneration_NoShiftAssignment(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	empID := uuid.New()
	req := CreateEventRequest{
		EmployeeID:     empID.String(),
		EventType:      "CHECKIN",
		EventTimeUTC:   "2026-01-15T01:00:00Z",
		EventTimeLocal: "2026-01-15T08:00:00+07:00",
		Latitude:       -6.2088,
		Longitude:      106.8456,
	}
	resp, err := svc.CreateEvent(ctx(), req)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}
	if resp.EventType != "CHECKIN" {
		t.Errorf("expected event to still be created, got %+v", resp)
	}
	// No assertion on the session itself here beyond "doesn't crash" -
	// recalculateSession degrades to DAY_OFF per §25 when there's no shift.
}
