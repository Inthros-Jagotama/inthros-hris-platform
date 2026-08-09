package attendance

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

// fakeApprovalEngine is a test double for ApprovalEngine.
type fakeApprovalEngine struct {
	createCalls []struct {
		module     string
		documentID string
		flowID     string
	}
	instanceID     string
	createErr      error
	activeFlowID   string
	flowResolveErr error
}

func (f *fakeApprovalEngine) CreateApprovalInstance(ctx context.Context, module, documentID, flowID string) (string, error) {
	f.createCalls = append(f.createCalls, struct {
		module     string
		documentID string
		flowID     string
	}{module, documentID, flowID})
	if f.createErr != nil {
		return "", f.createErr
	}
	if f.instanceID == "" {
		f.instanceID = uuid.New().String()
	}
	return f.instanceID, nil
}

func (f *fakeApprovalEngine) GetApprovalInstanceStatus(ctx context.Context, instanceID string) (string, error) {
	return "PENDING", nil
}

func (f *fakeApprovalEngine) GetActiveFlowIDForModule(ctx context.Context, module string) (string, error) {
	if f.flowResolveErr != nil {
		return "", f.flowResolveErr
	}
	return f.activeFlowID, nil
}

func TestService_CreateOvertimeRequest_WithApprovalEngine_CreatesInstance(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{}
	svc.SetApprovalEngine(fake)

	flowID := uuid.New().String()
	req := CreateOvertimeRequest{
		EmployeeID:       uuid.New().String(),
		WorkDate:         "2026-01-15",
		StartTimeLocal:   "2026-01-15T18:00:00+07:00",
		EndTimeLocal:     "2026-01-15T20:00:00+07:00",
		RequestedMinutes: 120,
		Reason:           "Deadline crunch",
		FlowID:           &flowID,
	}

	resp, err := svc.CreateOvertimeRequest(ctx(), req)
	if err != nil {
		t.Fatalf("CreateOvertimeRequest failed: %v", err)
	}

	if resp.Status != "PENDING_APPROVAL" {
		t.Errorf("expected status PENDING_APPROVAL, got '%s'", resp.Status)
	}
	if resp.ApprovalInstanceID == nil || *resp.ApprovalInstanceID != fake.instanceID {
		t.Errorf("expected approval_instance_id %s, got %v", fake.instanceID, resp.ApprovalInstanceID)
	}
	if len(fake.createCalls) != 1 {
		t.Fatalf("expected 1 CreateApprovalInstance call, got %d", len(fake.createCalls))
	}
	if fake.createCalls[0].module != "attendance" || fake.createCalls[0].flowID != flowID {
		t.Errorf("unexpected call params: %+v", fake.createCalls[0])
	}
}

func TestService_CreateOvertimeRequest_NoFlowID_ResolvesActiveFlow(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	activeFlow := uuid.New().String()
	fake := &fakeApprovalEngine{activeFlowID: activeFlow}
	svc.SetApprovalEngine(fake)

	req := CreateOvertimeRequest{
		EmployeeID:       uuid.New().String(),
		WorkDate:         "2026-01-15",
		StartTimeLocal:   "2026-01-15T18:00:00+07:00",
		EndTimeLocal:     "2026-01-15T20:00:00+07:00",
		RequestedMinutes: 120,
		Reason:           "Deadline crunch",
	}

	resp, err := svc.CreateOvertimeRequest(ctx(), req)
	if err != nil {
		t.Fatalf("CreateOvertimeRequest failed: %v", err)
	}
	if resp.Status != "PENDING_APPROVAL" {
		t.Errorf("expected status PENDING_APPROVAL, got '%s'", resp.Status)
	}
	if resp.ApprovalInstanceID == nil || *resp.ApprovalInstanceID != fake.instanceID {
		t.Errorf("expected approval_instance_id %s, got %v", fake.instanceID, resp.ApprovalInstanceID)
	}
	if len(fake.createCalls) != 1 {
		t.Fatalf("expected 1 CreateApprovalInstance call, got %d", len(fake.createCalls))
	}
	if fake.createCalls[0].module != "attendance" || fake.createCalls[0].flowID != activeFlow {
		t.Errorf("unexpected call params: %+v", fake.createCalls[0])
	}
}

func TestService_CreateOvertimeRequest_NoFlowID_NoActiveFlow_SkipsApproval(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{flowResolveErr: errors.New("no active flow for module")}
	svc.SetApprovalEngine(fake)

	req := CreateOvertimeRequest{
		EmployeeID:       uuid.New().String(),
		WorkDate:         "2026-01-15",
		StartTimeLocal:   "2026-01-15T18:00:00+07:00",
		EndTimeLocal:     "2026-01-15T20:00:00+07:00",
		RequestedMinutes: 120,
		Reason:           "Deadline crunch",
	}

	resp, err := svc.CreateOvertimeRequest(ctx(), req)
	if err != nil {
		t.Fatalf("CreateOvertimeRequest failed: %v", err)
	}
	if resp.Status != "SUBMITTED" {
		t.Errorf("expected status SUBMITTED, got '%s'", resp.Status)
	}
	if len(fake.createCalls) != 0 {
		t.Errorf("expected no CreateApprovalInstance calls, got %d", len(fake.createCalls))
	}
}

func TestService_CreateCorrectionRequest_WithApprovalEngine_CreatesInstance(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{}
	svc.SetApprovalEngine(fake)

	flowID := uuid.New().String()
	req := CreateCorrectionRequest{
		EmployeeID:          uuid.New().String(),
		AttendanceSessionID: uuid.New().String(),
		CorrectionType:      string(CorrectionTypeMissingCheckout),
		RequestedCheckout:   strPtr("2026-01-15T17:00:00+07:00"),
		Reason:              "Forgot to check out",
		FlowID:              &flowID,
	}

	resp, err := svc.CreateCorrectionRequest(ctx(), req)
	if err != nil {
		t.Fatalf("CreateCorrectionRequest failed: %v", err)
	}

	if resp.Status != "PENDING_APPROVAL" {
		t.Errorf("expected status PENDING_APPROVAL, got '%s'", resp.Status)
	}
	if resp.ApprovalInstanceID == nil || *resp.ApprovalInstanceID != fake.instanceID {
		t.Errorf("expected approval_instance_id %s, got %v", fake.instanceID, resp.ApprovalInstanceID)
	}
	if len(fake.createCalls) != 1 {
		t.Fatalf("expected 1 CreateApprovalInstance call, got %d", len(fake.createCalls))
	}
	if fake.createCalls[0].module != "attendance" || fake.createCalls[0].flowID != flowID {
		t.Errorf("unexpected call params: %+v", fake.createCalls[0])
	}
}

func TestService_CreateCorrectionRequest_NoFlowID_ResolvesActiveFlow(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	activeFlow := uuid.New().String()
	fake := &fakeApprovalEngine{activeFlowID: activeFlow}
	svc.SetApprovalEngine(fake)

	req := CreateCorrectionRequest{
		EmployeeID:          uuid.New().String(),
		AttendanceSessionID: uuid.New().String(),
		CorrectionType:      string(CorrectionTypeMissingCheckout),
		RequestedCheckout:   strPtr("2026-01-15T17:00:00+07:00"),
		Reason:              "Forgot to check out",
	}

	resp, err := svc.CreateCorrectionRequest(ctx(), req)
	if err != nil {
		t.Fatalf("CreateCorrectionRequest failed: %v", err)
	}
	if resp.Status != "PENDING_APPROVAL" {
		t.Errorf("expected status PENDING_APPROVAL, got '%s'", resp.Status)
	}
	if resp.ApprovalInstanceID == nil || *resp.ApprovalInstanceID != fake.instanceID {
		t.Errorf("expected approval_instance_id %s, got %v", fake.instanceID, resp.ApprovalInstanceID)
	}
	if len(fake.createCalls) != 1 {
		t.Fatalf("expected 1 CreateApprovalInstance call, got %d", len(fake.createCalls))
	}
	if fake.createCalls[0].module != "attendance" || fake.createCalls[0].flowID != activeFlow {
		t.Errorf("unexpected call params: %+v", fake.createCalls[0])
	}
}

func TestService_CreateCorrectionRequest_NoFlowID_NoActiveFlow_SkipsApproval(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	fake := &fakeApprovalEngine{flowResolveErr: errors.New("no active flow for module")}
	svc.SetApprovalEngine(fake)

	req := CreateCorrectionRequest{
		EmployeeID:          uuid.New().String(),
		AttendanceSessionID: uuid.New().String(),
		CorrectionType:      string(CorrectionTypeMissingCheckout),
		RequestedCheckout:   strPtr("2026-01-15T17:00:00+07:00"),
		Reason:              "Forgot to check out",
	}

	resp, err := svc.CreateCorrectionRequest(ctx(), req)
	if err != nil {
		t.Fatalf("CreateCorrectionRequest failed: %v", err)
	}
	if resp.Status != "SUBMITTED" {
		t.Errorf("expected status SUBMITTED, got '%s'", resp.Status)
	}
	if len(fake.createCalls) != 0 {
		t.Errorf("expected no CreateApprovalInstance calls, got %d", len(fake.createCalls))
	}
}

func TestService_HandleApprovalStatusChange_Approved(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	o := createTestOvertimeRequest(repo, uuid.New())
	o.Status = OvertimePendingApproval
	if err := repo.UpdateOvertimeRequest(ctx(), o); err != nil {
		t.Fatalf("failed to seed overtime request: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), o.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetOvertimeRequestByID(ctx(), o.ID.String())
	if err != nil {
		t.Fatalf("GetOvertimeRequestByID failed: %v", err)
	}
	if updated.Status != "APPROVED" {
		t.Errorf("expected status APPROVED, got '%s'", updated.Status)
	}
	if updated.ApprovedAt == nil {
		t.Error("expected ApprovedAt to be set")
	}
}

func TestService_HandleApprovalStatusChange_Approved_CalculatesOvertimeFromSession(t *testing.T) {
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
	checkout.EventTimeUTC = "2026-01-15T12:00:00Z"
	checkout.EventTimeLocal = "2026-01-15T19:00:00+07:00" // 2h (120m) past planned 17:00
	if _, err := svc.CreateEvent(ctx(), checkout); err != nil {
		t.Fatalf("checkout CreateEvent failed: %v", err)
	}

	o := createTestOvertimeRequest(repo, empID) // WorkDate 2026-01-15, RequestedMinutes 120
	o.Status = OvertimePendingApproval
	if err := repo.UpdateOvertimeRequest(ctx(), o); err != nil {
		t.Fatalf("failed to seed overtime request: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), o.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetOvertimeRequestByID(ctx(), o.ID.String())
	if err != nil {
		t.Fatalf("GetOvertimeRequestByID failed: %v", err)
	}
	if updated.ActualMinutes == nil || *updated.ActualMinutes != 120 {
		t.Errorf("expected actual_minutes 120, got %v", updated.ActualMinutes)
	}
	if updated.CalculatedMinutes == nil || *updated.CalculatedMinutes != 120 {
		t.Errorf("expected calculated_minutes 120, got %v", updated.CalculatedMinutes)
	}

	session, err := repo.FindSessionByEmployeeAndDate(ctx(), empID, "2026-01-15")
	if err != nil {
		t.Fatalf("expected session: %v", err)
	}
	if !session.IsOvertimeDay {
		t.Error("expected session.IsOvertimeDay = true")
	}
	if session.OvertimeMinutes != 120 {
		t.Errorf("expected session.OvertimeMinutes 120, got %d", session.OvertimeMinutes)
	}
}

func TestService_HandleApprovalStatusChange_Approved_CalculatedMinutesCappedByRequested(t *testing.T) {
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
	checkout.EventTimeUTC = "2026-01-15T13:00:00Z"
	checkout.EventTimeLocal = "2026-01-15T20:00:00+07:00" // 3h (180m) past planned 17:00
	if _, err := svc.CreateEvent(ctx(), checkout); err != nil {
		t.Fatalf("checkout CreateEvent failed: %v", err)
	}

	o := createTestOvertimeRequest(repo, empID) // RequestedMinutes 120
	o.Status = OvertimePendingApproval
	if err := repo.UpdateOvertimeRequest(ctx(), o); err != nil {
		t.Fatalf("failed to seed overtime request: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), o.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetOvertimeRequestByID(ctx(), o.ID.String())
	if err != nil {
		t.Fatalf("GetOvertimeRequestByID failed: %v", err)
	}
	if updated.ActualMinutes == nil || *updated.ActualMinutes != 180 {
		t.Errorf("expected actual_minutes 180, got %v", updated.ActualMinutes)
	}
	if updated.CalculatedMinutes == nil || *updated.CalculatedMinutes != 120 {
		t.Errorf("expected calculated_minutes capped at requested 120, got %v", updated.CalculatedMinutes)
	}
}

func TestService_HandleApprovalStatusChange_Rejected(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	o := createTestOvertimeRequest(repo, uuid.New())
	o.Status = OvertimePendingApproval
	if err := repo.UpdateOvertimeRequest(ctx(), o); err != nil {
		t.Fatalf("failed to seed overtime request: %v", err)
	}

	note := "insufficient justification"
	if err := svc.HandleApprovalStatusChange(ctx(), o.ID, "REJECTED", note); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetOvertimeRequestByID(ctx(), o.ID.String())
	if err != nil {
		t.Fatalf("GetOvertimeRequestByID failed: %v", err)
	}
	if updated.Status != "REJECTED" {
		t.Errorf("expected status REJECTED, got '%s'", updated.Status)
	}
	if updated.ApprovalNote == nil || *updated.ApprovalNote != note {
		t.Errorf("expected approval note %q, got %v", note, updated.ApprovalNote)
	}
}

func TestService_HandleApprovalStatusChange_NotPendingApproval_NoOp(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	o := createTestOvertimeRequest(repo, uuid.New())
	// Status is left as SUBMITTED (the default from createTestOvertimeRequest).

	if err := svc.HandleApprovalStatusChange(ctx(), o.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetOvertimeRequestByID(ctx(), o.ID.String())
	if err != nil {
		t.Fatalf("GetOvertimeRequestByID failed: %v", err)
	}
	if updated.Status != "SUBMITTED" {
		t.Errorf("expected status to remain SUBMITTED, got '%s'", updated.Status)
	}
}
