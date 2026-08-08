package leave

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// getOrCreateLeaveBalance returns the employee's current-balance cache row
// for (leaveTypeID, year), creating one seeded from LeaveType.DefaultQuotaDays
// if it doesn't exist yet.
func (s *Service) getOrCreateLeaveBalance(ctx context.Context, employeeID, leaveTypeID uuid.UUID, year int, leaveType *LeaveType) (*EmployeeLeaveBalance, error) {
	bal, err := s.repo.FindLeaveBalance(ctx, employeeID, leaveTypeID, year)
	if err == nil {
		return bal, nil
	}
	quota := 0.0
	if leaveType.DefaultQuotaDays != nil {
		quota = float64(*leaveType.DefaultQuotaDays)
	}
	bal = &EmployeeLeaveBalance{
		EmployeeID:    employeeID,
		LeaveTypeID:   leaveTypeID,
		PeriodYear:    year,
		QuotaDays:     quota,
		UsedDays:      0,
		RemainingDays: quota,
	}
	if err := s.repo.UpsertLeaveBalance(ctx, bal); err != nil {
		return nil, fmt.Errorf("failed to initialize leave balance: %w", err)
	}
	return bal, nil
}

// applyLeaveUsage deducts an approved leave request's days from the
// employee's balance and records a USAGE ledger entry (§9/§10/§34). Leave
// types with CountsAgainstQuota = false never touch the balance (§12).
func (s *Service) applyLeaveUsage(ctx context.Context, lr *LeaveRequest, leaveType *LeaveType) error {
	if !leaveType.CountsAgainstQuota || lr.RequestedDays <= 0 {
		return nil
	}
	year, err := requestYear(lr.RequestStartDate)
	if err != nil {
		return err
	}
	bal, err := s.getOrCreateLeaveBalance(ctx, lr.EmployeeID, lr.LeaveTypeID, year, leaveType)
	if err != nil {
		return err
	}
	before := bal.RemainingDays
	bal.UsedDays += lr.RequestedDays
	bal.RemainingDays -= lr.RequestedDays
	if err := s.repo.UpsertLeaveBalance(ctx, bal); err != nil {
		return fmt.Errorf("failed to update leave balance: %w", err)
	}
	return s.writeLeaveBalanceTransaction(ctx, bal, LeaveTransactionUsage, -lr.RequestedDays, before, "leave_request", lr.ID)
}

// reverseLeaveUsage reverses a previously-applied USAGE deduction — used
// when an APPROVED_FINAL leave request is cancelled (§18).
func (s *Service) reverseLeaveUsage(ctx context.Context, lr *LeaveRequest, leaveType *LeaveType) error {
	if !leaveType.CountsAgainstQuota || lr.RequestedDays <= 0 {
		return nil
	}
	year, err := requestYear(lr.RequestStartDate)
	if err != nil {
		return err
	}
	bal, err := s.getOrCreateLeaveBalance(ctx, lr.EmployeeID, lr.LeaveTypeID, year, leaveType)
	if err != nil {
		return err
	}
	before := bal.RemainingDays
	bal.UsedDays -= lr.RequestedDays
	bal.RemainingDays += lr.RequestedDays
	if err := s.repo.UpsertLeaveBalance(ctx, bal); err != nil {
		return fmt.Errorf("failed to update leave balance: %w", err)
	}
	return s.writeLeaveBalanceTransaction(ctx, bal, LeaveTransactionReversal, lr.RequestedDays, before, "leave_request", lr.ID)
}

func (s *Service) writeLeaveBalanceTransaction(ctx context.Context, bal *EmployeeLeaveBalance, txnType LeaveTransactionType, amount, before float64, referenceType string, referenceID uuid.UUID) error {
	txn := &LeaveBalanceTransaction{
		EmployeeID:      bal.EmployeeID,
		LeaveTypeID:     bal.LeaveTypeID,
		BalanceID:       bal.ID,
		TransactionType: txnType,
		ReferenceType:   &referenceType,
		ReferenceID:     &referenceID,
		Amount:          amount,
		BalanceBefore:   before,
		BalanceAfter:    bal.RemainingDays,
	}
	return s.repo.CreateLeaveBalanceTransaction(ctx, txn)
}

// requestYear extracts the year from a request_start_date. Some DB drivers
// (e.g. the sqlite driver used in tests) round-trip a "date"-typed column as
// a full RFC3339 timestamp string instead of plain "2006-01-02", so both
// layouts are accepted here.
func requestYear(requestStartDate string) (int, error) {
	if t, err := time.Parse(dateLayout, requestStartDate); err == nil {
		return t.Year(), nil
	}
	if t, err := time.Parse(time.RFC3339, requestStartDate); err == nil {
		return t.Year(), nil
	}
	return 0, fmt.Errorf("invalid request_start_date: %q", requestStartDate)
}
