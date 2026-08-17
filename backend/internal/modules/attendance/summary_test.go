package attendance

import (
	"testing"

	"github.com/google/uuid"
)

func TestService_GetEmployeeCalendar_ReturnsSessionsInRange(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	shift := createTestShift(repo)
	empID := uuid.New()
	createTestEmployeeShift(repo, empID, shift.ID)

	for _, date := range []string{"2026-01-15", "2026-01-16", "2026-02-01"} {
		checkin := CreateEventRequest{
			EmployeeID:     empID.String(),
			EventType:      "CHECKIN",
			EventTimeUTC:   date + "T01:00:00Z",
			EventTimeLocal: date + "T08:00:00+07:00",
			Latitude:       -6.2088,
			Longitude:      106.8456,
		}
		if _, err := svc.CreateEvent(ctx(), checkin); err != nil {
			t.Fatalf("CreateEvent failed for %s: %v", date, err)
		}
		checkout := checkin
		checkout.EventType = "CHECKOUT"
		checkout.EventTimeUTC = date + "T10:00:00Z"
		checkout.EventTimeLocal = date + "T17:00:00+07:00"
		if _, err := svc.CreateEvent(ctx(), checkout); err != nil {
			t.Fatalf("checkout CreateEvent failed for %s: %v", date, err)
		}
	}

	calendar, err := svc.GetEmployeeCalendar(ctx(), empID.String(), "2026-01-01", "2026-01-31")
	if err != nil {
		t.Fatalf("GetEmployeeCalendar failed: %v", err)
	}
	if len(calendar) != 2 {
		t.Fatalf("expected 2 sessions in January, got %d", len(calendar))
	}
}

func TestService_GetAttendanceReport_TenantWide_AcrossEmployees(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	shift := createTestShift(repo)
	emp1 := uuid.New()
	emp2 := uuid.New()
	createTestEmployeeShift(repo, emp1, shift.ID)
	createTestEmployeeShift(repo, emp2, shift.ID)

	for _, empID := range []uuid.UUID{emp1, emp2} {
		checkin := CreateEventRequest{
			EmployeeID:     empID.String(),
			EventType:      "CHECKIN",
			EventTimeUTC:   "2026-01-15T01:00:00Z",
			EventTimeLocal: "2026-01-15T08:00:00+07:00",
			Latitude:       -6.2088,
			Longitude:      106.8456,
		}
		if _, err := svc.CreateEvent(ctx(), checkin); err != nil {
			t.Fatalf("CreateEvent failed: %v", err)
		}
	}

	report, err := svc.GetAttendanceReport(ctx(), "2026-01-15", "2026-01-15")
	if err != nil {
		t.Fatalf("GetAttendanceReport failed: %v", err)
	}
	if len(report) != 2 {
		t.Fatalf("expected 2 sessions (one per employee) on 2026-01-15, got %d", len(report))
	}
}

func TestService_GetEmployeeSummary_AggregatesCorrectly(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	shift := createTestShift(repo) // 08:00 - 17:00
	empID := uuid.New()
	createTestEmployeeShift(repo, empID, shift.ID)

	// Day 1: on-time full day (CLOSED, not late)
	checkin1 := CreateEventRequest{
		EmployeeID:     empID.String(),
		EventType:      "CHECKIN",
		EventTimeUTC:   "2026-01-15T01:00:00Z",
		EventTimeLocal: "2026-01-15T08:00:00+07:00",
		Latitude:       -6.2088,
		Longitude:      106.8456,
	}
	if _, err := svc.CreateEvent(ctx(), checkin1); err != nil {
		t.Fatalf("checkin1 failed: %v", err)
	}
	checkout1 := checkin1
	checkout1.EventType = "CHECKOUT"
	checkout1.EventTimeUTC = "2026-01-15T10:00:00Z"
	checkout1.EventTimeLocal = "2026-01-15T17:00:00+07:00"
	if _, err := svc.CreateEvent(ctx(), checkout1); err != nil {
		t.Fatalf("checkout1 failed: %v", err)
	}

	// Day 2: late check-in, missing checkout
	checkin2 := CreateEventRequest{
		EmployeeID:     empID.String(),
		EventType:      "CHECKIN",
		EventTimeUTC:   "2026-01-16T01:30:00Z",
		EventTimeLocal: "2026-01-16T08:30:00+07:00",
		Latitude:       -6.2088,
		Longitude:      106.8456,
	}
	if _, err := svc.CreateEvent(ctx(), checkin2); err != nil {
		t.Fatalf("checkin2 failed: %v", err)
	}

	// Day 3: approved leave, no events at all.
	leaveReqID := uuid.New()
	if err := svc.ApplyApprovedLeave(ctx(), empID, "2026-01-17", leaveReqID, 1.0); err != nil {
		t.Fatalf("ApplyApprovedLeave failed: %v", err)
	}

	summary, err := svc.GetEmployeeSummary(ctx(), empID.String(), "2026-01-01", "2026-01-31")
	if err != nil {
		t.Fatalf("GetEmployeeSummary failed: %v", err)
	}
	if summary.TotalSessions != 3 {
		t.Errorf("expected 3 sessions, got %d", summary.TotalSessions)
	}
	if summary.PresentDays != 1 {
		t.Errorf("expected 1 present day, got %d", summary.PresentDays)
	}
	if summary.MissingCheckoutDays != 1 {
		t.Errorf("expected 1 missing checkout day, got %d", summary.MissingCheckoutDays)
	}
	if summary.LeaveDays != 1.0 {
		t.Errorf("expected 1.0 leave days, got %.2f", summary.LeaveDays)
	}
	if summary.TotalWorkMinutes != 540 {
		t.Errorf("expected 540 total work minutes, got %d", summary.TotalWorkMinutes)
	}
}

func TestService_GetAttendanceStats_TenantWideAggregation(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	shift := createTestShift(repo) // 08:00 - 17:00
	emp1 := uuid.New()
	emp2 := uuid.New()
	createTestEmployeeShift(repo, emp1, shift.ID)
	createTestEmployeeShift(repo, emp2, shift.ID)

	// emp1: on-time full day → CLOSED, not late → present.
	checkin := CreateEventRequest{
		EmployeeID:     emp1.String(),
		EventType:      "CHECKIN",
		EventTimeUTC:   "2026-01-15T01:00:00Z",
		EventTimeLocal: "2026-01-15T08:00:00+07:00",
		Latitude:       -6.2088,
		Longitude:      106.8456,
	}
	if _, err := svc.CreateEvent(ctx(), checkin); err != nil {
		t.Fatalf("emp1 checkin failed: %v", err)
	}
	checkout := checkin
	checkout.EventType = "CHECKOUT"
	checkout.EventTimeUTC = "2026-01-15T10:00:00Z"
	checkout.EventTimeLocal = "2026-01-15T17:00:00+07:00"
	if _, err := svc.CreateEvent(ctx(), checkout); err != nil {
		t.Fatalf("emp1 checkout failed: %v", err)
	}

	// emp2 day 1: late check-in + checkout → CLOSED with lateness → late.
	lateIn := CreateEventRequest{
		EmployeeID:     emp2.String(),
		EventType:      "CHECKIN",
		EventTimeUTC:   "2026-01-15T01:30:00Z",
		EventTimeLocal: "2026-01-15T08:30:00+07:00",
		Latitude:       -6.2088,
		Longitude:      106.8456,
	}
	if _, err := svc.CreateEvent(ctx(), lateIn); err != nil {
		t.Fatalf("emp2 checkin failed: %v", err)
	}
	lateOut := lateIn
	lateOut.EventType = "CHECKOUT"
	lateOut.EventTimeUTC = "2026-01-15T10:30:00Z"
	lateOut.EventTimeLocal = "2026-01-15T17:30:00+07:00"
	if _, err := svc.CreateEvent(ctx(), lateOut); err != nil {
		t.Fatalf("emp2 checkout failed: %v", err)
	}

	// emp2 day 2: check-in only → MISSING_CHECKOUT.
	noOut := CreateEventRequest{
		EmployeeID:     emp2.String(),
		EventType:      "CHECKIN",
		EventTimeUTC:   "2026-01-16T01:00:00Z",
		EventTimeLocal: "2026-01-16T08:00:00+07:00",
		Latitude:       -6.2088,
		Longitude:      106.8456,
	}
	if _, err := svc.CreateEvent(ctx(), noOut); err != nil {
		t.Fatalf("emp2 day2 checkin failed: %v", err)
	}

	// emp2: approved leave the next day.
	leaveReqID := uuid.New()
	if err := svc.ApplyApprovedLeave(ctx(), emp2, "2026-01-17", leaveReqID, 1.0); err != nil {
		t.Fatalf("ApplyApprovedLeave failed: %v", err)
	}

	stats, err := svc.GetAttendanceStats(ctx(), "2026-01-01", "2026-01-31")
	if err != nil {
		t.Fatalf("GetAttendanceStats failed: %v", err)
	}
	if stats.TotalSessions != 4 {
		t.Errorf("expected 4 sessions tenant-wide, got %d", stats.TotalSessions)
	}
	if stats.Present != 1 {
		t.Errorf("expected 1 present, got %d", stats.Present)
	}
	if stats.Late != 1 {
		t.Errorf("expected 1 late, got %d", stats.Late)
	}
	if stats.MissingCheckout != 1 {
		t.Errorf("expected 1 missing checkout, got %d", stats.MissingCheckout)
	}
	if stats.LeaveDays != 1.0 {
		t.Errorf("expected 1.0 leave days, got %.2f", stats.LeaveDays)
	}
	if stats.TotalWorkMinutes != 1080 {
		t.Errorf("expected 1080 total work minutes (540 + 540), got %d", stats.TotalWorkMinutes)
	}
}

func TestService_GetAttendanceStats_IncludesOvertimeAndTravel(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	// Lembur: 1 pending + 1 approved (120 + 60 menit calculated).
	empID := uuid.New()
	pending := createTestOvertimeRequest(repo, empID) // work_date 2026-01-15, 120 menit
	pending.Status = OvertimePendingApproval
	if err := repo.UpdateOvertimeRequest(ctx(), pending); err != nil {
		t.Fatalf("update overtime status failed: %v", err)
	}
	approved := &AttendanceOvertimeRequest{
		EmployeeID:      empID,
		WorkDate:        "2026-01-16",
		StartTimeLocal:  parseTime("2026-01-16T18:00:00+07:00"),
		EndTimeLocal:    parseTime("2026-01-16T19:00:00+07:00"),
		RequestedMinutes: 60,
		Status:          OvertimeApproved,
	}
	calc := 60
	approved.CalculatedMinutes = &calc
	if err := repo.CreateOvertimeRequest(ctx(), approved); err != nil {
		t.Fatalf("create approved overtime failed: %v", err)
	}

	// Perjalanan dinas: 1 approved + 1 completed, start_date dalam rentang.
	purpose := "Client visit"
	travels := []*BusinessTravel{
		{RequesterID: empID, RequestNumber: "TRV-TEST-001", Title: "Visit A", Purpose: &purpose, StartDate: parseTime("2026-01-10T00:00:00+07:00"), EndDate: parseTime("2026-01-12T00:00:00+07:00"), Status: TravelStatusApproved},
		{RequesterID: empID, RequestNumber: "TRV-TEST-002", Title: "Visit B", Purpose: &purpose, StartDate: parseTime("2026-01-20T00:00:00+07:00"), EndDate: parseTime("2026-01-22T00:00:00+07:00"), Status: TravelStatusCompleted},
		// di luar rentang → tidak dihitung.
		{RequesterID: empID, RequestNumber: "TRV-TEST-003", Title: "Visit C", Purpose: &purpose, StartDate: parseTime("2026-02-05T00:00:00+07:00"), EndDate: parseTime("2026-02-06T00:00:00+07:00"), Status: TravelStatusApproved},
	}
	for _, tr := range travels {
		if err := repo.CreateBusinessTravel(ctx(), tr); err != nil {
			t.Fatalf("create business travel failed: %v", err)
		}
	}

	stats, err := svc.GetAttendanceStats(ctx(), "2026-01-01", "2026-01-31")
	if err != nil {
		t.Fatalf("GetAttendanceStats failed: %v", err)
	}
	if stats.OvertimeTotal != 2 {
		t.Errorf("expected 2 overtime requests, got %d", stats.OvertimeTotal)
	}
	if stats.OvertimePending != 1 {
		t.Errorf("expected 1 pending overtime, got %d", stats.OvertimePending)
	}
	if stats.OvertimeApproved != 1 {
		t.Errorf("expected 1 approved overtime, got %d", stats.OvertimeApproved)
	}
	if stats.OvertimeMinutes != 60 {
		t.Errorf("expected 60 approved overtime minutes, got %d", stats.OvertimeMinutes)
	}
	if stats.TravelTotal != 2 {
		t.Errorf("expected 2 business travels in range, got %d", stats.TravelTotal)
	}
	if stats.TravelApproved != 1 {
		t.Errorf("expected 1 approved travel, got %d", stats.TravelApproved)
	}
	if stats.TravelCompleted != 1 {
		t.Errorf("expected 1 completed travel, got %d", stats.TravelCompleted)
	}
}

func TestService_GetOvertimeTrend_GroupsByIsoWeek(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	empID := uuid.New()
	// 2026-01-05 adalah Senin (minggu 1) dan 2026-01-12 Senin berikutnya.
	// Senin: 1 approved (60 menit) + 1 pending.
	monday := &AttendanceOvertimeRequest{
		EmployeeID:      empID,
		WorkDate:        "2026-01-05",
		StartTimeLocal:  parseTime("2026-01-05T18:00:00+07:00"),
		EndTimeLocal:    parseTime("2026-01-05T19:00:00+07:00"),
		RequestedMinutes: 60,
		Status:          OvertimeApproved,
	}
	calc := 60
	monday.CalculatedMinutes = &calc
	if err := repo.CreateOvertimeRequest(ctx(), monday); err != nil {
		t.Fatalf("create monday overtime failed: %v", err)
	}
	pending := createTestOvertimeRequest(repo, empID) // work_date 2026-01-15, SUBMITTED
	pending.Status = OvertimePendingApproval
	if err := repo.UpdateOvertimeRequest(ctx(), pending); err != nil {
		t.Fatalf("update pending overtime failed: %v", err)
	}

	// Kamis 2026-01-08: approved 120 menit — masih minggu yang sama dengan Senin.
	thu := &AttendanceOvertimeRequest{
		EmployeeID:      empID,
		WorkDate:        "2026-01-08",
		StartTimeLocal:  parseTime("2026-01-08T18:00:00+07:00"),
		EndTimeLocal:    parseTime("2026-01-08T20:00:00+07:00"),
		RequestedMinutes: 120,
		Status:          OvertimeApproved,
	}
	calc2 := 120
	thu.CalculatedMinutes = &calc2
	if err := repo.CreateOvertimeRequest(ctx(), thu); err != nil {
		t.Fatalf("create thursday overtime failed: %v", err)
	}

	// Selasa 2026-01-13: approved 90 menit → minggu berikutnya (mulai 12 Jan).
	tue := &AttendanceOvertimeRequest{
		EmployeeID:      empID,
		WorkDate:        "2026-01-13",
		StartTimeLocal:  parseTime("2026-01-13T18:00:00+07:00"),
		EndTimeLocal:    parseTime("2026-01-13T19:30:00+07:00"),
		RequestedMinutes: 90,
		Status:          OvertimeApproved,
	}
	calc3 := 90
	tue.CalculatedMinutes = &calc3
	if err := repo.CreateOvertimeRequest(ctx(), tue); err != nil {
		t.Fatalf("create tuesday overtime failed: %v", err)
	}

	trend, err := svc.GetOvertimeTrend(ctx(), "2026-01-05", "2026-01-18")
	if err != nil {
		t.Fatalf("GetOvertimeTrend failed: %v", err)
	}
	if len(trend.Weeks) != 2 {
		t.Fatalf("expected 2 weeks (05 & 12 Jan), got %d: %+v", len(trend.Weeks), trend.Weeks)
	}
	w1 := trend.Weeks[0]
	if w1.WeekStart != "2026-01-05" {
		t.Errorf("expected first week start 2026-01-05, got %s", w1.WeekStart)
	}
	if w1.Count != 2 || w1.Approved != 2 || w1.Minutes != 180 {
		t.Errorf("week 1: expected count=2 approved=2 minutes=180, got count=%d approved=%d minutes=%d", w1.Count, w1.Approved, w1.Minutes)
	}
	w2 := trend.Weeks[1]
	if w2.WeekStart != "2026-01-12" {
		t.Errorf("expected second week start 2026-01-12, got %s", w2.WeekStart)
	}
	if w2.Count != 2 || w2.Approved != 1 || w2.Minutes != 90 {
		t.Errorf("week 2: expected count=2 approved=1 minutes=90, got count=%d approved=%d minutes=%d", w2.Count, w2.Approved, w2.Minutes)
	}
}
