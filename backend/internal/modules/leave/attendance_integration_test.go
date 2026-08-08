package leave

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

type fakeAttendanceSessionUpdater struct {
	calls []struct {
		employeeID     uuid.UUID
		workDate       string
		leaveRequestID uuid.UUID
		dayFraction    float64
	}
}

func (f *fakeAttendanceSessionUpdater) ApplyApprovedLeave(_ context.Context, employeeID uuid.UUID, workDate string, leaveRequestID uuid.UUID, dayFraction float64) error {
	f.calls = append(f.calls, struct {
		employeeID     uuid.UUID
		workDate       string
		leaveRequestID uuid.UUID
		dayFraction    float64
	}{employeeID, workDate, leaveRequestID, dayFraction})
	return nil
}

func TestService_HandleApprovalStatusChange_Approved_PushesAttendanceIntegration(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	fakeApproval := &fakeApprovalEngine{}
	svc.SetApprovalEngine(fakeApproval)
	fakeAttendance := &fakeAttendanceSessionUpdater{}
	svc.SetAttendanceSessionUpdater(fakeAttendance)

	lType := createTestLeaveType(repo)
	empID := uuid.New()
	// Seeded directly (not via svc.CreateLeaveRequest) to avoid SubmittedAt
	// being set - the glebarez/sqlite test driver can't round-trip a
	// populated `type:timestamp(6)` column back through a full-row SELECT
	// (see the comment on leaveRequestStatusNote in
	// approval_integration_test.go), which would otherwise break the
	// FindLeaveRequestByID inside HandleApprovalStatusChange below. Details
	// are seeded directly too, mirroring what svc.CreateLeaveRequest does.
	lr := createTestLeaveRequest(repo, empID, lType.ID) // 2026-01-15..2026-01-16
	lr.Status = LeaveStatusPendingApproval
	if err := repo.UpdateLeaveRequest(ctx(), lr); err != nil {
		t.Fatalf("failed to seed leave request: %v", err)
	}
	for _, date := range []string{"2026-01-15", "2026-01-16"} {
		detail := &LeaveRequestDetail{
			LeaveRequestID: lr.ID,
			EmployeeID:     empID,
			LeaveDate:      date,
			DayFraction:    1.0,
		}
		if err := repo.CreateLeaveRequestDetail(ctx(), detail); err != nil {
			t.Fatalf("failed to seed leave request detail: %v", err)
		}
	}
	lrID := lr.ID

	if err := svc.HandleApprovalStatusChange(ctx(), lrID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	if len(fakeAttendance.calls) != 2 {
		t.Fatalf("expected 2 ApplyApprovedLeave calls (one per day), got %d", len(fakeAttendance.calls))
	}
	for _, c := range fakeAttendance.calls {
		if c.employeeID != empID {
			t.Errorf("expected employee_id %s, got %s", empID, c.employeeID)
		}
		if c.leaveRequestID != lrID {
			t.Errorf("expected leave_request_id %s, got %s", lrID, c.leaveRequestID)
		}
		if c.dayFraction != 1.0 {
			t.Errorf("expected day_fraction 1.0, got %.2f", c.dayFraction)
		}
	}
}

func TestService_HandleApprovalStatusChange_Approved_NoAttendanceUpdater_DoesNotFail(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	// No SetAttendanceSessionUpdater call - should degrade gracefully.

	lType := createTestLeaveType(repo)
	lr := createTestLeaveRequest(repo, uuid.New(), lType.ID)
	lr.Status = LeaveStatusPendingApproval
	if err := repo.UpdateLeaveRequest(ctx(), lr); err != nil {
		t.Fatalf("failed to seed leave request: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), lr.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}
}
