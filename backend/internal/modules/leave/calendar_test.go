package leave

import (
	"testing"

	"github.com/google/uuid"
)

func TestService_GetEmployeeCalendar_ReturnsEntriesInRange(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	lType := createTestLeaveType(repo)
	empID := uuid.New()

	// 2026-01-15 (Thu) - 2026-01-16 (Fri): within range
	if _, err := svc.CreateLeaveRequest(ctx(), CreateLeaveRequest{
		EmployeeID:       empID.String(),
		LeaveTypeID:      lType.ID.String(),
		RequestStartDate: "2026-01-15",
		RequestEndDate:   "2026-01-16",
	}); err != nil {
		t.Fatalf("CreateLeaveRequest failed: %v", err)
	}

	// 2026-02-02 (Mon): outside the January range we'll query
	if _, err := svc.CreateLeaveRequest(ctx(), CreateLeaveRequest{
		EmployeeID:       empID.String(),
		LeaveTypeID:      lType.ID.String(),
		RequestStartDate: "2026-02-02",
		RequestEndDate:   "2026-02-02",
	}); err != nil {
		t.Fatalf("CreateLeaveRequest failed: %v", err)
	}

	entries, err := svc.GetEmployeeCalendar(ctx(), empID.String(), "2026-01-01", "2026-01-31")
	if err != nil {
		t.Fatalf("GetEmployeeCalendar failed: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 calendar entries in January, got %d", len(entries))
	}
	if entries[0].LeaveDate != "2026-01-15" || entries[1].LeaveDate != "2026-01-16" {
		t.Errorf("expected entries for 2026-01-15 and 2026-01-16, got %s and %s", entries[0].LeaveDate, entries[1].LeaveDate)
	}
	for _, e := range entries {
		if e.Status != "SUBMITTED" {
			t.Errorf("expected status SUBMITTED, got %s", e.Status)
		}
		if e.LeaveTypeID != lType.ID.String() {
			t.Errorf("expected leave_type_id %s, got %s", lType.ID.String(), e.LeaveTypeID)
		}
	}
}

func TestService_GetEmployeeCalendar_ExcludesRejectedRequests(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	lType := createTestLeaveType(repo)
	empID := uuid.New()

	// Seeded directly via the repository (like createTestLeaveRequest) rather
	// than svc.CreateLeaveRequest, so submitted_at stays unpopulated and we
	// avoid the sqlite driver's populated-timestamp(6) round-trip quirk on
	// the FindLeaveRequestByID re-fetch inside HandleApprovalStatusChange.
	lr := createTestLeaveRequest(repo, empID, lType.ID)
	lr.Status = LeaveStatusPendingApproval
	if err := repo.UpdateLeaveRequest(ctx(), lr); err != nil {
		t.Fatalf("failed to seed pending status: %v", err)
	}
	if err := repo.CreateLeaveRequestDetail(ctx(), &LeaveRequestDetail{
		LeaveRequestID: lr.ID,
		EmployeeID:     empID,
		LeaveDate:      "2026-01-15",
		DayFraction:    1.0,
		IsPaid:         true,
	}); err != nil {
		t.Fatalf("failed to seed leave request detail: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), lr.ID, "REJECTED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	entries, err := svc.GetEmployeeCalendar(ctx(), empID.String(), "2026-01-01", "2026-01-31")
	if err != nil {
		t.Fatalf("GetEmployeeCalendar failed: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 calendar entries after rejection, got %d", len(entries))
	}
}
