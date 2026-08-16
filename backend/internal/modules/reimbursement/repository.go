package reimbursement

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	dbFunc TenantDBFunc
}

func NewRepository(dbFunc TenantDBFunc) *Repository {
	return &Repository{dbFunc: dbFunc}
}

func (r *Repository) db(ctx context.Context) (*gorm.DB, error) {
	return r.dbFunc(ctx)
}

// =========================================================================
// Reimbursement Types
// =========================================================================

func (r *Repository) CreateReimbursementType(ctx context.Context, t *ReimbursementType) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(t).Error
}

func (r *Repository) FindReimbursementTypeByID(ctx context.Context, id uuid.UUID) (*ReimbursementType, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var t ReimbursementType
	if err := db.WithContext(ctx).First(&t, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("reimbursement type not found")
		}
		return nil, err
	}
	return &t, nil
}

func (r *Repository) ListReimbursementTypes(ctx context.Context, page, perPage int) ([]ReimbursementType, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var types []ReimbursementType
	var total int64

	query := db.WithContext(ctx).Model(&ReimbursementType{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&types).Error; err != nil {
		return nil, 0, err
	}
	return types, total, nil
}

// ListReimbursementTypeCodes returns every existing type code (non-deleted),
// used to auto-generate unique codes for new types without asking the user.
func (r *Repository) ListReimbursementTypeCodes(ctx context.Context) ([]string, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var codes []string
	if err := db.WithContext(ctx).Model(&ReimbursementType{}).Pluck("code", &codes).Error; err != nil {
		return nil, err
	}
	return codes, nil
}

func (r *Repository) UpdateReimbursementType(ctx context.Context, t *ReimbursementType) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(t).Error
}

func (r *Repository) DeleteReimbursementType(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&ReimbursementType{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("reimbursement type not found")
	}
	return result.Error
}

// =========================================================================
// Reimbursement Requests
// =========================================================================

func (r *Repository) CreateReimbursementRequest(ctx context.Context, req *ReimbursementRequest) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(req).Error
}

func (r *Repository) FindReimbursementRequestByID(ctx context.Context, id uuid.UUID) (*ReimbursementRequest, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var req ReimbursementRequest
	if err := db.WithContext(ctx).Preload("Items").First(&req, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("reimbursement request not found")
		}
		return nil, err
	}
	return &req, nil
}

func (r *Repository) ListReimbursementRequests(ctx context.Context, employeeID *uuid.UUID, status *string, page, perPage int) ([]ReimbursementRequest, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var requests []ReimbursementRequest
	var total int64

	query := db.WithContext(ctx).Model(&ReimbursementRequest{})
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&requests).Error; err != nil {
		return nil, 0, err
	}
	return requests, total, nil
}

func (r *Repository) UpdateReimbursementRequest(ctx context.Context, req *ReimbursementRequest) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(req).Error
}

func (r *Repository) DeleteReimbursementRequest(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&ReimbursementRequest{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("reimbursement request not found")
	}
	return result.Error
}

// =========================================================================
// Reimbursement Items
// =========================================================================

func (r *Repository) CreateReimbursementItem(ctx context.Context, item *ReimbursementItem) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(item).Error
}

func (r *Repository) FindReimbursementItemByID(ctx context.Context, id uuid.UUID) (*ReimbursementItem, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var item ReimbursementItem
	if err := db.WithContext(ctx).First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("reimbursement item not found")
		}
		return nil, err
	}
	return &item, nil
}

func (r *Repository) ListReimbursementItems(ctx context.Context, requestID uuid.UUID) ([]ReimbursementItem, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var items []ReimbursementItem
	if err := db.WithContext(ctx).Where("reimbursement_request_id = ?", requestID).
		Order("expense_date ASC, created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) UpdateReimbursementItem(ctx context.Context, item *ReimbursementItem) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(item).Error
}

func (r *Repository) DeleteReimbursementItem(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&ReimbursementItem{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("reimbursement item not found")
	}
	return result.Error
}

// =========================================================================
// User resolution
// =========================================================================

// FindUserIDByEmployeeID resolves an employee's platform user_id via
// employee_accounts, mirroring the employee_id -> user_id resolution already
// used by the approval module (GetUserIDsByOrganization) and leave. Returns
// nil if the employee has no linked user account.
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

// FindEmployeeIDByUserID resolves the logged-in user's own employee record
// via employee_accounts (user_id -> employee_id), mirroring the reverse of
// FindUserIDByEmployeeID and the /user-accounts/me resolution — the JWT
// user_id is the user-account UUID, which is a different UUID from the
// employee record, so requests must be attributed to the employee id.
// Returns nil if the user has no linked employee account.
func (r *Repository) FindEmployeeIDByUserID(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var empIDStrs []string
	err = db.WithContext(ctx).Table("employee_accounts").
		Where("user_id = ?", userID).
		Limit(1).
		Pluck("employee_id", &empIDStrs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to resolve employee id: %w", err)
	}
	if len(empIDStrs) == 0 || empIDStrs[0] == "" {
		return nil, nil
	}
	empID, err := uuid.Parse(empIDStrs[0])
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}
	return &empID, nil
}

// =========================================================================
// Aggregation
// =========================================================================

func (r *Repository) SumReimbursementItems(ctx context.Context, requestID uuid.UUID) (float64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return 0, err
	}
	var sum float64
	if err := db.WithContext(ctx).Model(&ReimbursementItem{}).
		Where("reimbursement_request_id = ?", requestID).
		Select("COALESCE(SUM(amount), 0)").Scan(&sum).Error; err != nil {
		return 0, err
	}
	return sum, nil
}
