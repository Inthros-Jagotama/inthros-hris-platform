package leave

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	dbResolver func(ctx context.Context) (*gorm.DB, error)
}

func NewRepository(resolver func(ctx context.Context) (*gorm.DB, error)) *Repository {
	return &Repository{dbResolver: resolver}
}

func (r *Repository) db(ctx context.Context) (*gorm.DB, error) {
	return r.dbResolver(ctx)
}

// =========================================================================
// Leave Types
// =========================================================================

func (r *Repository) CreateLeaveType(ctx context.Context, t *LeaveType) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(t).Error
}

func (r *Repository) FindLeaveTypeByID(ctx context.Context, id uuid.UUID) (*LeaveType, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var t LeaveType
	if err := db.WithContext(ctx).Where("id = ?", id).First(&t).Error; err != nil {
		return nil, fmt.Errorf("leave type not found: %w", err)
	}
	return &t, nil
}

func (r *Repository) ListLeaveTypes(ctx context.Context, page, perPage int) ([]LeaveType, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var total int64
	db.WithContext(ctx).Model(&LeaveType{}).Count(&total)

	var types []LeaveType
	offset := (page - 1) * perPage
	if err := db.WithContext(ctx).Offset(offset).Limit(perPage).Order("name ASC").Find(&types).Error; err != nil {
		return nil, 0, err
	}
	return types, total, nil
}

func (r *Repository) UpdateLeaveType(ctx context.Context, t *LeaveType) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(t).Error
}

func (r *Repository) DeleteLeaveType(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Where("id = ?", id).Delete(&LeaveType{})
	if result.RowsAffected == 0 {
		return fmt.Errorf("leave type not found")
	}
	return result.Error
}

// =========================================================================
// Leave Accrual Policies
// =========================================================================

func (r *Repository) CreateAccrualPolicy(ctx context.Context, p *LeaveAccrualPolicy) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(p).Error
}

func (r *Repository) FindAccrualPolicyByID(ctx context.Context, id uuid.UUID) (*LeaveAccrualPolicy, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var p LeaveAccrualPolicy
	if err := db.WithContext(ctx).Where("id = ?", id).First(&p).Error; err != nil {
		return nil, fmt.Errorf("accrual policy not found: %w", err)
	}
	return &p, nil
}

func (r *Repository) ListAccrualPolicies(ctx context.Context, leaveTypeID *uuid.UUID, page, perPage int) ([]LeaveAccrualPolicy, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	query := db.WithContext(ctx).Model(&LeaveAccrualPolicy{})
	if leaveTypeID != nil {
		query = query.Where("leave_type_id = ?", *leaveTypeID)
	}
	var total int64
	query.Count(&total)

	var policies []LeaveAccrualPolicy
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("effective_from DESC").Find(&policies).Error; err != nil {
		return nil, 0, err
	}
	return policies, total, nil
}

func (r *Repository) UpdateAccrualPolicy(ctx context.Context, p *LeaveAccrualPolicy) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(p).Error
}

func (r *Repository) DeleteAccrualPolicy(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Where("id = ?", id).Delete(&LeaveAccrualPolicy{})
	if result.RowsAffected == 0 {
		return fmt.Errorf("accrual policy not found")
	}
	return result.Error
}

// =========================================================================
// Leave Reasons
// =========================================================================

func (r *Repository) CreateLeaveReason(ctx context.Context, reason *LeaveReason) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(reason).Error
}

func (r *Repository) FindLeaveReasonByID(ctx context.Context, id uuid.UUID) (*LeaveReason, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var reason LeaveReason
	if err := db.WithContext(ctx).Where("id = ?", id).First(&reason).Error; err != nil {
		return nil, fmt.Errorf("leave reason not found: %w", err)
	}
	return &reason, nil
}

func (r *Repository) ListLeaveReasons(ctx context.Context) ([]LeaveReason, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var reasons []LeaveReason
	if err := db.WithContext(ctx).Order("sort_order ASC, name ASC").Find(&reasons).Error; err != nil {
		return nil, err
	}
	return reasons, nil
}

func (r *Repository) UpdateLeaveReason(ctx context.Context, reason *LeaveReason) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(reason).Error
}

func (r *Repository) DeleteLeaveReason(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Where("id = ?", id).Delete(&LeaveReason{})
	if result.RowsAffected == 0 {
		return fmt.Errorf("leave reason not found")
	}
	return result.Error
}

// =========================================================================
// Leave Requests
// =========================================================================

func (r *Repository) CreateLeaveRequest(ctx context.Context, req *LeaveRequest) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(req).Error
}

func (r *Repository) FindLeaveRequestByID(ctx context.Context, id uuid.UUID) (*LeaveRequest, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var req LeaveRequest
	if err := db.WithContext(ctx).Where("id = ?", id).First(&req).Error; err != nil {
		return nil, fmt.Errorf("leave request not found: %w", err)
	}
	return &req, nil
}

func (r *Repository) ListLeaveRequests(ctx context.Context, employeeID *uuid.UUID, status *string, page, perPage int) ([]LeaveRequest, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	query := db.WithContext(ctx).Model(&LeaveRequest{})
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}
	var total int64
	query.Count(&total)

	var requests []LeaveRequest
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&requests).Error; err != nil {
		return nil, 0, err
	}
	return requests, total, nil
}

// CountOverlappingLeaveRequests counts the employee's non-final leave
// requests whose date range overlaps [startDate, endDate], excluding
// REJECTED_FINAL/CANCELLED requests which no longer hold the dates.
func (r *Repository) CountOverlappingLeaveRequests(ctx context.Context, employeeID uuid.UUID, startDate, endDate string) (int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return 0, err
	}
	var count int64
	err = db.WithContext(ctx).Model(&LeaveRequest{}).
		Where("employee_id = ?", employeeID).
		Where("status NOT IN ?", []string{string(LeaveStatusRejectedFinal), string(LeaveStatusCancelled)}).
		Where("request_start_date <= ? AND request_end_date >= ?", endDate, startDate).
		Count(&count).Error
	return count, err
}

func (r *Repository) UpdateLeaveRequest(ctx context.Context, req *LeaveRequest) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(req).Error
}

func (r *Repository) DeleteLeaveRequest(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Where("id = ?", id).Delete(&LeaveRequest{})
	if result.RowsAffected == 0 {
		return fmt.Errorf("leave request not found")
	}
	return result.Error
}

// =========================================================================
// Leave Request Details
// =========================================================================

func (r *Repository) CreateLeaveRequestDetail(ctx context.Context, detail *LeaveRequestDetail) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(detail).Error
}

func (r *Repository) ListLeaveRequestDetails(ctx context.Context, leaveRequestID uuid.UUID) ([]LeaveRequestDetail, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var details []LeaveRequestDetail
	if err := db.WithContext(ctx).Where("leave_request_id = ?", leaveRequestID).Order("leave_date ASC").Find(&details).Error; err != nil {
		return nil, err
	}
	return details, nil
}

func (r *Repository) DeleteLeaveRequestDetails(ctx context.Context, leaveRequestID uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Where("leave_request_id = ?", leaveRequestID).Delete(&LeaveRequestDetail{}).Error
}

// =========================================================================
// Employee Leave Balances
// =========================================================================

func (r *Repository) FindLeaveBalance(ctx context.Context, employeeID, leaveTypeID uuid.UUID, year int) (*EmployeeLeaveBalance, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var bal EmployeeLeaveBalance
	if err := db.WithContext(ctx).
		Where("employee_id = ? AND leave_type_id = ? AND period_year = ?", employeeID, leaveTypeID, year).
		First(&bal).Error; err != nil {
		return nil, fmt.Errorf("leave balance not found: %w", err)
	}
	return &bal, nil
}

func (r *Repository) UpsertLeaveBalance(ctx context.Context, bal *EmployeeLeaveBalance) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(bal).Error
}

func (r *Repository) ListLeaveBalances(ctx context.Context, employeeID *uuid.UUID, page, perPage int) ([]EmployeeLeaveBalance, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	query := db.WithContext(ctx).Model(&EmployeeLeaveBalance{})
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	var total int64
	query.Count(&total)

	var balances []EmployeeLeaveBalance
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("period_year DESC").Find(&balances).Error; err != nil {
		return nil, 0, err
	}
	return balances, total, nil
}

// CreateLeaveBalanceTransaction inserts one ledger entry. Callers are
// responsible for computing BalanceBefore/BalanceAfter and updating the
// corresponding EmployeeLeaveBalance row — this is a pure append, no
// recalculation happens here (business logic lands in Phase 6).
func (r *Repository) CreateLeaveBalanceTransaction(ctx context.Context, txn *LeaveBalanceTransaction) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(txn).Error
}

// ListLeaveBalanceTransactions returns the ledger history for one employee's
// leave type, newest first — the audit trail behind a current-balance row.
func (r *Repository) ListLeaveBalanceTransactions(ctx context.Context, employeeID, leaveTypeID uuid.UUID, page, perPage int) ([]LeaveBalanceTransaction, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	query := db.WithContext(ctx).Model(&LeaveBalanceTransaction{}).
		Where("employee_id = ? AND leave_type_id = ?", employeeID, leaveTypeID)

	var total int64
	query.Count(&total)

	var txns []LeaveBalanceTransaction
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&txns).Error; err != nil {
		return nil, 0, err
	}
	return txns, total, nil
}

// FindUserIDByEmployeeID resolves an employee's platform user_id via
// employee_accounts, mirroring the employee_id -> user_id resolution already
// used by the approval module (GetUserIDsByOrganization). Returns nil if the
// employee has no linked user account.
func (r *Repository) FindUserIDByEmployeeID(ctx context.Context, employeeID uuid.UUID) (*uuid.UUID, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var userIDStrs []string
	err = db.WithContext(ctx).Table("employee_accounts").
		Where("employee_id = ?", employeeID).
		Limit(1).
		Pluck("user_id", &userIDStrs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to resolve employee user id: %w", err)
	}
	if len(userIDStrs) == 0 || userIDStrs[0] == "" {
		return nil, nil
	}
	userID, err := uuid.Parse(userIDStrs[0])
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	return &userID, nil
}

// CalendarEntry is one date-level row of an employee's leave calendar —
// a leave_request_details row joined back to its parent leave_requests for
// status/leave_type_id, since details themselves don't carry those fields.
type CalendarEntry struct {
	LeaveRequestID uuid.UUID `json:"leave_request_id"`
	LeaveDate      string    `json:"leave_date"`
	DayFraction    float64   `json:"day_fraction"`
	LeaveTypeID    uuid.UUID `json:"leave_type_id"`
	Status         string    `json:"status"`
}

// FindCalendarEntriesForEmployeeInRange returns the employee's leave dates in
// [fromDate, toDate], excluding REJECTED_FINAL requests (dates that never
// actually happened as leave, per §26's "Pending Leave -> tidak boleh
// dianggap sebagai approved leave" principle extended to the calendar: a
// rejected date shouldn't show up as if it were a leave day at all).
// CANCELLED is kept visible so the employee can still see why a date is free
// again, mirroring how the request list itself never hides cancelled rows.
func (r *Repository) FindCalendarEntriesForEmployeeInRange(ctx context.Context, employeeID uuid.UUID, fromDate, toDate string) ([]CalendarEntry, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var entries []CalendarEntry
	err = db.WithContext(ctx).
		Table("leave_request_details AS d").
		Select("d.leave_request_id, d.leave_date, d.day_fraction, r.leave_type_id, r.status").
		Joins("JOIN leave_requests AS r ON r.id = d.leave_request_id").
		Where("d.employee_id = ? AND d.leave_date BETWEEN ? AND ?", employeeID, fromDate, toDate).
		Where("r.status <> ?", string(LeaveStatusRejectedFinal)).
		Where("r.deleted_at IS NULL AND d.deleted_at IS NULL").
		Order("d.leave_date ASC").
		Find(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to load leave calendar entries: %w", err)
	}
	return entries, nil
}
