package leave

import (
	"testing"

	"github.com/google/uuid"
)

func TestService_HandleApprovalStatusChange_Approved_DeductsBalance(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	lType := createTestLeaveType(repo)
	empID := uuid.New()
	lr := createTestLeaveRequest(repo, empID, lType.ID) // RequestedDays: 2
	lr.Status = LeaveStatusPendingApproval
	if err := repo.UpdateLeaveRequest(ctx(), lr); err != nil {
		t.Fatalf("failed to seed leave request: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), lr.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	bal, err := repo.FindLeaveBalance(ctx(), empID, lType.ID, 2026)
	if err != nil {
		t.Fatalf("expected leave balance to be created, got error: %v", err)
	}
	if bal.UsedDays != 2 {
		t.Errorf("expected used_days 2, got %.2f", bal.UsedDays)
	}
	if bal.RemainingDays != -2 {
		t.Errorf("expected remaining_days -2 (quota defaults to 0), got %.2f", bal.RemainingDays)
	}

	txns, total, err := repo.ListLeaveBalanceTransactions(ctx(), empID, lType.ID, 1, 10)
	if err != nil {
		t.Fatalf("ListLeaveBalanceTransactions failed: %v", err)
	}
	if total != 1 || len(txns) != 1 {
		t.Fatalf("expected 1 ledger transaction, got %d", total)
	}
	if txns[0].TransactionType != LeaveTransactionUsage || txns[0].Amount != -2 {
		t.Errorf("unexpected ledger transaction: %+v", txns[0])
	}
}

func TestService_HandleApprovalStatusChange_Approved_SkipsWhenNotCountingAgainstQuota(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	// A raw struct insert with CountsAgainstQuota: false gets re-defaulted to
	// true by GORM's `default:1` tag on INSERT (it treats the Go zero value
	// as "unset"), so create normally first, then flip it via an UPDATE
	// (unaffected by the default tag) to actually persist false.
	lType := createTestLeaveType(repo)
	lType.CountsAgainstQuota = false
	if err := repo.UpdateLeaveType(ctx(), lType); err != nil {
		t.Fatalf("failed to update leave type: %v", err)
	}
	empID := uuid.New()
	lr := createTestLeaveRequest(repo, empID, lType.ID)
	lr.Status = LeaveStatusPendingApproval
	if err := repo.UpdateLeaveRequest(ctx(), lr); err != nil {
		t.Fatalf("failed to seed leave request: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx(), lr.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	if _, err := repo.FindLeaveBalance(ctx(), empID, lType.ID, 2026); err == nil {
		t.Error("expected no leave balance row to be created for a leave type that does not count against quota")
	}
}

func TestService_UpdateLeaveRequestStatus_CancelApproved_ReversesBalance(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	lType := createTestLeaveType(repo)
	empID := uuid.New()
	// Seed directly as APPROVED_FINAL (skipping the PENDING_APPROVAL->APPROVED
	// transition) so ApprovedAt stays nil: the glebarez/sqlite test driver
	// can't round-trip a populated `type:timestamp(6)` column back through a
	// full-row SELECT (see the comment on leaveRequestStatusNote in
	// approval_integration_test.go), which would otherwise break the
	// FindLeaveRequestByID inside UpdateLeaveRequestStatus below.
	lr := createTestLeaveRequest(repo, empID, lType.ID)
	lr.Status = LeaveStatusApprovedFinal
	if err := repo.UpdateLeaveRequest(ctx(), lr); err != nil {
		t.Fatalf("failed to seed leave request: %v", err)
	}
	if err := svc.applyLeaveUsage(ctx(), lr, lType); err != nil {
		t.Fatalf("failed to seed balance deduction: %v", err)
	}

	if _, err := svc.UpdateLeaveRequestStatus(ctx(), lr.ID.String(), string(LeaveStatusCancelled), "cancelled by HR"); err != nil {
		t.Fatalf("UpdateLeaveRequestStatus failed: %v", err)
	}

	bal, err := repo.FindLeaveBalance(ctx(), empID, lType.ID, 2026)
	if err != nil {
		t.Fatalf("expected leave balance to exist: %v", err)
	}
	if bal.UsedDays != 0 {
		t.Errorf("expected used_days back to 0 after reversal, got %.2f", bal.UsedDays)
	}
	if bal.RemainingDays != 0 {
		t.Errorf("expected remaining_days back to 0 after reversal, got %.2f", bal.RemainingDays)
	}

	_, total, err := repo.ListLeaveBalanceTransactions(ctx(), empID, lType.ID, 1, 10)
	if err != nil {
		t.Fatalf("ListLeaveBalanceTransactions failed: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected 2 ledger transactions (USAGE + REVERSAL), got %d", total)
	}
}
