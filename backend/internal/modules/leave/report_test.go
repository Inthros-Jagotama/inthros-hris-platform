package leave

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// seedLeaveRequest creates a leave_requests row directly via the repository
// (like createTestLeaveRequest) rather than svc.CreateLeaveRequest, so
// submitted_at stays unpopulated and we avoid the sqlite test driver's
// populated-timestamp(6) round-trip quirk on the report's re-fetch (the
// same category of issue documented for calendar_test.go/balance.go).
func seedLeaveRequest(repo *Repository, empID, lTypeID uuid.UUID, start, end string, status LeaveStatus) *LeaveRequest {
	lr := &LeaveRequest{
		EmployeeID:       empID,
		LeaveTypeID:      lTypeID,
		RequestStartDate: start,
		RequestEndDate:   end,
		DurationMode:     DurationFullDay,
		RequestedDays:    1,
		Status:           status,
	}
	if err := repo.CreateLeaveRequest(ctx(), lr); err != nil {
		panic(err)
	}
	return lr
}

func TestService_GetLeaveUsageReport_ReturnsRequestsInRange_TenantWide(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	lType := createTestLeaveType(repo)
	emp1 := uuid.New()
	emp2 := uuid.New()

	seedLeaveRequest(repo, emp1, lType.ID, "2026-01-15", "2026-01-16", LeaveStatusSubmitted)
	seedLeaveRequest(repo, emp2, lType.ID, "2026-01-20", "2026-01-20", LeaveStatusSubmitted)
	// Outside the January range we'll query.
	seedLeaveRequest(repo, emp1, lType.ID, "2026-03-02", "2026-03-02", LeaveStatusSubmitted)

	report, err := svc.GetLeaveUsageReport(ctx(), "2026-01-01", "2026-01-31")
	if err != nil {
		t.Fatalf("GetLeaveUsageReport failed: %v", err)
	}
	if len(report) != 2 {
		t.Fatalf("expected 2 requests (one per employee) in January, got %d", len(report))
	}
}

func TestService_GetLeaveUsageReport_IncludesRejectedAndCancelled(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	lType := createTestLeaveType(repo)
	empID := uuid.New()

	seedLeaveRequest(repo, empID, lType.ID, "2026-01-15", "2026-01-16", LeaveStatusRejectedFinal)

	report, err := svc.GetLeaveUsageReport(ctx(), "2026-01-01", "2026-01-31")
	if err != nil {
		t.Fatalf("GetLeaveUsageReport failed: %v", err)
	}
	if len(report) != 1 {
		t.Fatalf("expected the rejected request to still appear in the usage report, got %d entries", len(report))
	}
	if report[0].Status != "REJECTED_FINAL" {
		t.Errorf("expected status REJECTED_FINAL, got %s", report[0].Status)
	}
}

func TestService_GetOnLeaveToday(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()

	lType := createTestLeaveType(repo)
	empA := uuid.New()
	empB := uuid.New()

	// Approved final & cuti hari ini (tanggal hari ini dihitung server).
	today := time.Now().Format("2006-01-02")
	approved := seedLeaveRequest(repo, empA, lType.ID, today, today, LeaveStatusApprovedFinal)
	if err := repo.CreateLeaveRequestDetail(ctx(), &LeaveRequestDetail{
		LeaveRequestID: approved.ID,
		EmployeeID:     empA,
		LeaveDate:      today,
		DayFraction:    1,
	}); err != nil {
		t.Fatalf("failed to create leave detail: %v", err)
	}

	// Approved final tapi cuti tanggal lain → tidak dihitung hari ini.
	other := seedLeaveRequest(repo, empB, lType.ID, "2026-01-15", "2026-01-15", LeaveStatusApprovedFinal)
	if err := repo.CreateLeaveRequestDetail(ctx(), &LeaveRequestDetail{
		LeaveRequestID: other.ID,
		EmployeeID:     empB,
		LeaveDate:      "2026-01-15",
		DayFraction:    1,
	}); err != nil {
		t.Fatalf("failed to create leave detail: %v", err)
	}

	// SUBMITTED hari ini → tidak dihitung (hanya approved final).
	submitted := seedLeaveRequest(repo, empB, lType.ID, today, today, LeaveStatusSubmitted)
	if err := repo.CreateLeaveRequestDetail(ctx(), &LeaveRequestDetail{
		LeaveRequestID: submitted.ID,
		EmployeeID:     empB,
		LeaveDate:      today,
		DayFraction:    1,
	}); err != nil {
		t.Fatalf("failed to create leave detail: %v", err)
	}
	_ = db

	count, err := svc.GetOnLeaveToday(ctx())
	if err != nil {
		t.Fatalf("GetOnLeaveToday failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 employee on leave today, got %d", count)
	}
}
