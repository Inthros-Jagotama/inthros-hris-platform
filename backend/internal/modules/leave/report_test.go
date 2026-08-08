package leave

import (
	"testing"

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
