package attendance

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestService_CreateCorrectionRequest_Success(t *testing.T) {
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
	session, err := repo.FindSessionByEmployeeAndDate(ctx(), empID, "2026-01-15")
	if err != nil {
		t.Fatalf("expected session: %v", err)
	}
	if session.Status != SessionStatusMissingCheckOut {
		t.Fatalf("expected MISSING_CHECKOUT, got %s", session.Status)
	}

	checkoutTime := "2026-01-15T17:00:00+07:00"
	req := CreateCorrectionRequest{
		EmployeeID:          empID.String(),
		AttendanceSessionID: session.ID.String(),
		CorrectionType:      string(CorrectionTypeMissingCheckout),
		RequestedCheckout:   &checkoutTime,
		Reason:              "Forgot to check out",
	}
	resp, err := svc.CreateCorrectionRequest(ctx(), req)
	if err != nil {
		t.Fatalf("CreateCorrectionRequest failed: %v", err)
	}
	if resp.Status != "SUBMITTED" {
		t.Errorf("expected status SUBMITTED, got '%s'", resp.Status)
	}
}

func TestService_HandleApprovalStatusChange_Correction_MissingCheckout_AppliesToSession(t *testing.T) {
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
	session, err := repo.FindSessionByEmployeeAndDate(ctx(), empID, "2026-01-15")
	if err != nil {
		t.Fatalf("expected session: %v", err)
	}

	checkoutTime := "2026-01-15T17:00:00+07:00"
	c := &AttendanceCorrectionRequest{
		EmployeeID:          empID,
		AttendanceSessionID: session.ID,
		CorrectionType:      CorrectionTypeMissingCheckout,
		RequestedCheckout:   timePtr(parseTime(checkoutTime)),
		Reason:              "Forgot to check out",
		Status:              CorrectionPendingApproval,
	}
	if err := repo.CreateCorrectionRequest(ctx(), c); err != nil {
		t.Fatalf("failed to seed correction request: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), c.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetCorrectionRequestByID(ctx(), c.ID.String())
	if err != nil {
		t.Fatalf("GetCorrectionRequestByID failed: %v", err)
	}
	if updated.Status != "APPROVED" {
		t.Errorf("expected status APPROVED, got '%s'", updated.Status)
	}

	refreshedSession, err := repo.FindSessionByEmployeeAndDate(ctx(), empID, "2026-01-15")
	if err != nil {
		t.Fatalf("expected session: %v", err)
	}
	if refreshedSession.Status != SessionStatusClosed {
		t.Errorf("expected session status CLOSED after correction, got '%s'", refreshedSession.Status)
	}
	if refreshedSession.WorkMinutes != 540 {
		t.Errorf("expected 540 work minutes (9h), got %d", refreshedSession.WorkMinutes)
	}
	if refreshedSession.CheckoutEventID == nil {
		t.Error("expected session to have a checkout event after correction")
	}
}

func TestService_HandleApprovalStatusChange_Correction_Rejected(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	shift := createTestShift(repo)
	empID := uuid.New()
	createTestEmployeeShift(repo, empID, shift.ID)
	session := &AttendanceSession{EmployeeID: empID, WorkDate: "2026-01-15", Status: SessionStatusMissingCheckOut}
	if err := repo.UpsertSession(ctx(), session); err != nil {
		t.Fatalf("failed to seed session: %v", err)
	}

	c := &AttendanceCorrectionRequest{
		EmployeeID:          empID,
		AttendanceSessionID: session.ID,
		CorrectionType:      CorrectionTypeMissingCheckout,
		Reason:              "Forgot to check out",
		Status:              CorrectionPendingApproval,
	}
	if err := repo.CreateCorrectionRequest(ctx(), c); err != nil {
		t.Fatalf("failed to seed correction request: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), c.ID, "REJECTED", "not enough evidence"); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetCorrectionRequestByID(ctx(), c.ID.String())
	if err != nil {
		t.Fatalf("GetCorrectionRequestByID failed: %v", err)
	}
	if updated.Status != "REJECTED" {
		t.Errorf("expected status REJECTED, got '%s'", updated.Status)
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
