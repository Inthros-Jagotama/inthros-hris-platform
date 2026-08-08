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
