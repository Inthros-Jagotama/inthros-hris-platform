package leave

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func ctx() context.Context {
	return context.Background()
}

// =========================================================================
// Leave Type Service Tests
// =========================================================================

func TestService_CreateLeaveType_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateLeaveTypeRequest{
		Name:        "Annual Leave",
		Description: "Annual paid leave",
		IsPaid:      boolPtr(true),
	}

	resp, err := svc.CreateLeaveType(ctx(), req)
	if err != nil {
		t.Fatalf("CreateLeaveType failed: %v", err)
	}

	if resp.Name != "Annual Leave" {
		t.Errorf("expected name 'Annual Leave', got '%s'", resp.Name)
	}
	if !resp.IsPaid {
		t.Error("expected IsPaid = true")
	}
	if resp.ID == "" {
		t.Error("expected ID to be set")
	}
}

func TestService_GetLeaveTypeByID_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestLeaveType(repo)

	found, err := svc.GetLeaveTypeByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetLeaveTypeByID failed: %v", err)
	}

	if found.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), found.ID)
	}
}

func TestService_ListLeaveTypes_DefaultPagination(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	for i := 0; i < 3; i++ {
		createTestLeaveType(repo)
	}

	resp, err := svc.ListLeaveTypes(ctx(), 0, 0)
	if err != nil {
		t.Fatalf("ListLeaveTypes failed: %v", err)
	}

	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}
}

func TestService_UpdateLeaveType_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestLeaveType(repo)
	newName := "Updated Leave"
	req := UpdateLeaveTypeRequest{Name: &newName}

	updated, err := svc.UpdateLeaveType(ctx(), created.ID.String(), req)
	if err != nil {
		t.Fatalf("UpdateLeaveType failed: %v", err)
	}

	if updated.Name != "Updated Leave" {
		t.Errorf("expected name 'Updated Leave', got '%s'", updated.Name)
	}
}

func TestService_DeleteLeaveType_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	created := createTestLeaveType(repo)

	if err := svc.DeleteLeaveType(ctx(), created.ID.String()); err != nil {
		t.Fatalf("DeleteLeaveType failed: %v", err)
	}

	_, err := svc.GetLeaveTypeByID(ctx(), created.ID.String())
	if err == nil {
		t.Fatal("expected error after deleting leave type")
	}
}

// =========================================================================
// Accrual Policy Service Tests
// =========================================================================

func TestService_CreateAccrualPolicy_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	lType := createTestLeaveType(repo)

	req := CreateAccrualPolicyRequest{
		LeaveTypeID:   lType.ID.String(),
		BaseQuotaDays: 12,
		EffectiveFrom: "2026-01-01",
	}

	resp, err := svc.CreateAccrualPolicy(ctx(), req)
	if err != nil {
		t.Fatalf("CreateAccrualPolicy failed: %v", err)
	}

	if resp.BaseQuotaDays != 12 {
		t.Errorf("expected base_quota_days 12, got %.2f", resp.BaseQuotaDays)
	}
}

func TestService_ListAccrualPolicies_DefaultPagination(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	lType := createTestLeaveType(repo)
	createTestAccrualPolicy(repo, lType.ID)
	// Second policy needs different effective_from to avoid unique constraint
	p2 := &LeaveAccrualPolicy{
		LeaveTypeID:   lType.ID,
		BaseQuotaDays: 6,
		EffectiveFrom: "2026-06-01",
	}
	if err := repo.CreateAccrualPolicy(context.Background(), p2); err != nil {
		t.Fatalf("CreateAccrualPolicy p2 failed: %v", err)
	}

	resp, err := svc.ListAccrualPolicies(ctx(), nil, 0, 0)
	if err != nil {
		t.Fatalf("ListAccrualPolicies failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
}

// =========================================================================
// Leave Reason Service Tests
// =========================================================================

func TestService_CreateLeaveReason_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	req := CreateLeaveReasonRequest{
		Name:      "Personal",
		SortOrder: intPtr(2),
	}

	resp, err := svc.CreateLeaveReason(ctx(), req)
	if err != nil {
		t.Fatalf("CreateLeaveReason failed: %v", err)
	}

	if resp.Name != "Personal" {
		t.Errorf("expected name 'Personal', got '%s'", resp.Name)
	}
}

func TestService_ListLeaveReasons_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	createTestLeaveReason(repo)
	createTestLeaveReason(repo)

	reasons, err := svc.ListLeaveReasons(ctx())
	if err != nil {
		t.Fatalf("ListLeaveReasons failed: %v", err)
	}

	if len(reasons) != 2 {
		t.Errorf("expected 2 reasons, got %d", len(reasons))
	}
}

// =========================================================================
// Leave Request Service Tests
// =========================================================================

func TestService_CreateLeaveRequest_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	lType := createTestLeaveType(repo)

	req := CreateLeaveRequest{
		EmployeeID:      uuidStr(),
		LeaveTypeID:     lType.ID.String(),
		RequestStartDate: "2026-01-15",
		RequestEndDate:  "2026-01-16",
		RequestedDays:   2,
	}

	resp, err := svc.CreateLeaveRequest(ctx(), req)
	if err != nil {
		t.Fatalf("CreateLeaveRequest failed: %v", err)
	}

	if resp.Status != "SUBMITTED" {
		t.Errorf("expected status SUBMITTED, got '%s'", resp.Status)
	}
	if resp.RequestedDays != 2 {
		t.Errorf("expected 2 days, got %.2f", resp.RequestedDays)
	}
}

func TestService_GetLeaveRequestByID_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	lType := createTestLeaveType(repo)
	empID := uuid.New()
	created := createTestLeaveRequest(repo, empID, lType.ID)

	found, err := svc.GetLeaveRequestByID(ctx(), created.ID.String())
	if err != nil {
		t.Fatalf("GetLeaveRequestByID failed: %v", err)
	}

	if found.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), found.ID)
	}
}

func TestService_UpdateLeaveRequestStatus_Approve(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	lType := createTestLeaveType(repo)
	empID := uuid.New()
	created := createTestLeaveRequest(repo, empID, lType.ID)

	updated, err := svc.UpdateLeaveRequestStatus(ctx(), created.ID.String(), "APPROVED_FINAL", "Approved by HR")
	if err != nil {
		t.Fatalf("UpdateLeaveRequestStatus failed: %v", err)
	}

	if updated.Status != "APPROVED_FINAL" {
		t.Errorf("expected status APPROVED_FINAL, got '%s'", updated.Status)
	}
	if updated.ApprovedAt == nil {
		t.Error("expected ApprovedAt to be set")
	}
}

func TestService_ListLeaveRequests_DefaultPagination(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	lType := createTestLeaveType(repo)
	empID := uuid.New()
	createTestLeaveRequest(repo, empID, lType.ID)
	createTestLeaveRequest(repo, empID, lType.ID)

	resp, err := svc.ListLeaveRequests(ctx(), nil, nil, 0, 0)
	if err != nil {
		t.Fatalf("ListLeaveRequests failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
}

// =========================================================================
// Helpers
// =========================================================================

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}
