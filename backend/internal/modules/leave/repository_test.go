package leave

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// =========================================================================
// Leave Type Repository Tests
// =========================================================================

func TestRepo_CreateLeaveType_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	tp := &LeaveType{
		Name:        "Annual Leave",
		Description: "Annual paid leave",
		IsPaid:      true,
	}

	if err := repo.CreateLeaveType(ctx, tp); err != nil {
		t.Fatalf("CreateLeaveType failed: %v", err)
	}

	if tp.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindLeaveTypeByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestLeaveType(repo)

	found, err := repo.FindLeaveTypeByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("FindLeaveTypeByID failed: %v", err)
	}

	if found.Name == "" {
		t.Error("expected name to be non-empty")
	}
}

func TestRepo_FindLeaveTypeByID_NotFound(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	_, err := repo.FindLeaveTypeByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent leave type")
	}
}

func TestRepo_ListLeaveTypes_Pagination(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	for i := 0; i < 3; i++ {
		createTestLeaveType(repo)
	}

	types, total, err := repo.ListLeaveTypes(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("ListLeaveTypes failed: %v", err)
	}

	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(types) != 3 {
		t.Errorf("expected 3 types, got %d", len(types))
	}
}

func TestRepo_DeleteLeaveType_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	created := createTestLeaveType(repo)

	if err := repo.DeleteLeaveType(context.Background(), created.ID); err != nil {
		t.Fatalf("DeleteLeaveType failed: %v", err)
	}

	_, err := repo.FindLeaveTypeByID(context.Background(), created.ID)
	if err == nil {
		t.Fatal("expected error after deleting leave type")
	}
}

// =========================================================================
// Accrual Policy Repository Tests
// =========================================================================

func TestRepo_CreateAccrualPolicy_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	lType := createTestLeaveType(repo)

	p := &LeaveAccrualPolicy{
		LeaveTypeID:   lType.ID,
		BaseQuotaDays: 12,
		EffectiveFrom: "2026-01-01",
	}

	if err := repo.CreateAccrualPolicy(context.Background(), p); err != nil {
		t.Fatalf("CreateAccrualPolicy failed: %v", err)
	}

	if p.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindAccrualPolicyByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	lType := createTestLeaveType(repo)
	created := createTestAccrualPolicy(repo, lType.ID)

	found, err := repo.FindAccrualPolicyByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("FindAccrualPolicyByID failed: %v", err)
	}

	if found.BaseQuotaDays != 12 {
		t.Errorf("expected base_quota_days 12, got %.2f", found.BaseQuotaDays)
	}
}

func TestRepo_ListAccrualPolicies_ByType(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	lType1 := createTestLeaveType(repo)
	createTestAccrualPolicy(repo, lType1.ID)
	// Second policy needs different effective_from to avoid unique constraint
	ctx := context.Background()
	p2 := &LeaveAccrualPolicy{
		LeaveTypeID:   lType1.ID,
		BaseQuotaDays: 6,
		EffectiveFrom: "2026-06-01",
	}
	if err := repo.CreateAccrualPolicy(ctx, p2); err != nil {
		t.Fatalf("CreateAccrualPolicy p2 failed: %v", err)
	}

	policies, total, err := repo.ListAccrualPolicies(context.Background(), &lType1.ID, 1, 10)
	if err != nil {
		t.Fatalf("ListAccrualPolicies failed: %v", err)
	}

	if total != 2 {
		t.Errorf("expected 2 policies, got %d", total)
	}
	if len(policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(policies))
	}
}

// =========================================================================
// Leave Reason Repository Tests
// =========================================================================

func TestRepo_CreateLeaveReason_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	r := &LeaveReason{Name: "Sick", SortOrder: 1}
	if err := repo.CreateLeaveReason(context.Background(), r); err != nil {
		t.Fatalf("CreateLeaveReason failed: %v", err)
	}

	if r.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_ListLeaveReasons_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	createTestLeaveReason(repo)
	createTestLeaveReason(repo)

	reasons, err := repo.ListLeaveReasons(context.Background())
	if err != nil {
		t.Fatalf("ListLeaveReasons failed: %v", err)
	}

	if len(reasons) != 2 {
		t.Errorf("expected 2 reasons, got %d", len(reasons))
	}
}

// =========================================================================
// Leave Request Repository Tests
// =========================================================================

func TestRepo_CreateLeaveRequest_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	lType := createTestLeaveType(repo)
	empID := uuid.New()

	lr := &LeaveRequest{
		EmployeeID:      empID,
		LeaveTypeID:     lType.ID,
		RequestStartDate: "2026-01-15",
		RequestEndDate:  "2026-01-16",
		RequestedDays:   2,
		Status:          LeaveStatusSubmitted,
	}

	if err := repo.CreateLeaveRequest(context.Background(), lr); err != nil {
		t.Fatalf("CreateLeaveRequest failed: %v", err)
	}

	if lr.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_FindLeaveRequestByID_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	lType := createTestLeaveType(repo)
	empID := uuid.New()
	created := createTestLeaveRequest(repo, empID, lType.ID)

	found, err := repo.FindLeaveRequestByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("FindLeaveRequestByID failed: %v", err)
	}

	if found.Status != LeaveStatusSubmitted {
		t.Errorf("expected status SUBMITTED, got '%s'", found.Status)
	}
}

func TestRepo_ListLeaveRequests_ByEmployee(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	lType := createTestLeaveType(repo)
	empID := uuid.New()
	createTestLeaveRequest(repo, empID, lType.ID)
	createTestLeaveRequest(repo, empID, lType.ID)

	requests, total, err := repo.ListLeaveRequests(context.Background(), &empID, nil, 1, 10)
	if err != nil {
		t.Fatalf("ListLeaveRequests failed: %v", err)
	}

	if total != 2 {
		t.Errorf("expected 2 requests, got %d", total)
	}
	if len(requests) != 2 {
		t.Errorf("expected 2 requests, got %d", len(requests))
	}
}

// =========================================================================
// Leave Balance Repository Tests
// =========================================================================

func TestRepo_UpsertLeaveBalance_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	bal := &EmployeeLeaveBalance{
		EmployeeID:    uuid.New(),
		LeaveTypeID:   uuid.New(),
		PeriodYear:    2026,
		QuotaDays:     12,
		UsedDays:      2,
		RemainingDays: 10,
	}

	if err := repo.UpsertLeaveBalance(context.Background(), bal); err != nil {
		t.Fatalf("UpsertLeaveBalance failed: %v", err)
	}

	if bal.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_ListLeaveBalances_ByEmployee(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	empID := uuid.New()
	for i := 0; i < 3; i++ {
		bal := &EmployeeLeaveBalance{
			EmployeeID:    empID,
			LeaveTypeID:   uuid.New(),
			PeriodYear:    2026,
			QuotaDays:     12,
			UsedDays:      0,
			RemainingDays: 12,
		}
		if err := repo.UpsertLeaveBalance(context.Background(), bal); err != nil {
			t.Fatalf("UpsertLeaveBalance failed: %v", err)
		}
	}

	balances, total, err := repo.ListLeaveBalances(context.Background(), &empID, 1, 10)
	if err != nil {
		t.Fatalf("ListLeaveBalances failed: %v", err)
	}

	if total != 3 {
		t.Errorf("expected 3 balances, got %d", total)
	}
	if len(balances) != 3 {
		t.Errorf("expected 3 balances, got %d", len(balances))
	}
}

// =========================================================================
// Leave Balance Transaction (Ledger) Repository Tests
// =========================================================================

func TestRepo_CreateLeaveBalanceTransaction_Success(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	bal := &EmployeeLeaveBalance{
		EmployeeID:    uuid.New(),
		LeaveTypeID:   uuid.New(),
		PeriodYear:    2026,
		QuotaDays:     12,
		UsedDays:      0,
		RemainingDays: 12,
	}
	if err := repo.UpsertLeaveBalance(context.Background(), bal); err != nil {
		t.Fatalf("UpsertLeaveBalance failed: %v", err)
	}

	txn := &LeaveBalanceTransaction{
		EmployeeID:      bal.EmployeeID,
		LeaveTypeID:     bal.LeaveTypeID,
		BalanceID:       bal.ID,
		TransactionType: LeaveTransactionUsage,
		Amount:          -2,
		BalanceBefore:   12,
		BalanceAfter:    10,
	}
	if err := repo.CreateLeaveBalanceTransaction(context.Background(), txn); err != nil {
		t.Fatalf("CreateLeaveBalanceTransaction failed: %v", err)
	}
	if txn.ID == uuid.Nil {
		t.Error("expected ID to be auto-generated")
	}
}

func TestRepo_ListLeaveBalanceTransactions_ByEmployeeAndType(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	bal := &EmployeeLeaveBalance{
		EmployeeID:    uuid.New(),
		LeaveTypeID:   uuid.New(),
		PeriodYear:    2026,
		QuotaDays:     12,
		UsedDays:      0,
		RemainingDays: 12,
	}
	if err := repo.UpsertLeaveBalance(context.Background(), bal); err != nil {
		t.Fatalf("UpsertLeaveBalance failed: %v", err)
	}

	types := []LeaveTransactionType{LeaveTransactionAccrual, LeaveTransactionUsage, LeaveTransactionAdjustment}
	for _, tt := range types {
		txn := &LeaveBalanceTransaction{
			EmployeeID:      bal.EmployeeID,
			LeaveTypeID:     bal.LeaveTypeID,
			BalanceID:       bal.ID,
			TransactionType: tt,
			Amount:          1,
			BalanceBefore:   0,
			BalanceAfter:    1,
		}
		if err := repo.CreateLeaveBalanceTransaction(context.Background(), txn); err != nil {
			t.Fatalf("CreateLeaveBalanceTransaction failed: %v", err)
		}
	}

	// Unrelated employee's transaction must not leak into the result.
	other := &LeaveBalanceTransaction{
		EmployeeID:      uuid.New(),
		LeaveTypeID:     uuid.New(),
		BalanceID:       uuid.New(),
		TransactionType: LeaveTransactionAccrual,
		Amount:          1,
		BalanceBefore:   0,
		BalanceAfter:    1,
	}
	if err := repo.CreateLeaveBalanceTransaction(context.Background(), other); err != nil {
		t.Fatalf("CreateLeaveBalanceTransaction failed: %v", err)
	}

	txns, total, err := repo.ListLeaveBalanceTransactions(context.Background(), bal.EmployeeID, bal.LeaveTypeID, 1, 10)
	if err != nil {
		t.Fatalf("ListLeaveBalanceTransactions failed: %v", err)
	}
	if total != 3 {
		t.Errorf("expected 3 transactions, got %d", total)
	}
	if len(txns) != 3 {
		t.Errorf("expected 3 transactions, got %d", len(txns))
	}
}
