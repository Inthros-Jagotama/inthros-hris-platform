package employeemovement

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository untuk database operations Employee Movement & Career Management.
type Repository struct {
	dbResolver func(ctx context.Context) (*gorm.DB, error)
}

// NewRepository membuat Repository baru.
func NewRepository(dbResolver func(ctx context.Context) (*gorm.DB, error)) *Repository {
	return &Repository{dbResolver: dbResolver}
}

func (r *Repository) getDB(ctx context.Context) (*gorm.DB, error) {
	return r.dbResolver(ctx)
}

// =========================================================================
// Employee Movement
// =========================================================================

func (r *Repository) CreateMovement(ctx context.Context, m *EmployeeMovement) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Create(m).Error; err != nil {
		return fmt.Errorf("failed to create employee movement: %w", err)
	}
	return nil
}

// FindUserIDByEmployeeID resolves the employee's linked user account id for
// notifications (same pattern attendance/leave repository). Returns nil if
// the employee has no linked user account.
func (r *Repository) FindUserIDByEmployeeID(ctx context.Context, employeeID uuid.UUID) (*uuid.UUID, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var userIDStrs []string
	err = db.Table("employee_accounts").
		Where("employee_id = ?", employeeID).
		Limit(1).
		Pluck("user_id", &userIDStrs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to resolve employee user id: %w", err)
	}
	if len(userIDStrs) == 0 || userIDStrs[0] == "" {
		return nil, nil
	}
	uid, parseErr := uuid.Parse(userIDStrs[0])
	if parseErr != nil {
		return nil, fmt.Errorf("invalid user id for employee: %w", parseErr)
	}
	return &uid, nil
}

// employeeRefInfo holds display info (name + employee code) for an employee.
type employeeRefInfo struct {
	EmployeeID string
	Name       string
	Code       string
}

// GetEmployeeInfoByIDs resolves display info (name, employee code) for a batch
// of employees at once — used to enrich movement/contract responses without an
// N+1 query per row (pola sama attendance.GetEmployeeInfoByIDs).
func (r *Repository) GetEmployeeInfoByIDs(ctx context.Context, employeeIDs []uuid.UUID) (map[string]employeeRefInfo, error) {
	result := make(map[string]employeeRefInfo, len(employeeIDs))
	if len(employeeIDs) == 0 {
		return result, nil
	}
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var rows []employeeRefInfo
	err = db.Table("employees").
		Where("id IN ?", employeeIDs).
		Select("id AS employee_id, name AS name, employee_id AS code").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to resolve employee info: %w", err)
	}
	for _, rrow := range rows {
		result[rrow.EmployeeID] = rrow
	}
	return result, nil
}

// GetOrganizationNamesByIDs resolves organization nomenclatures for a batch of
// ids (used for from/to_organization_name enrichment).
func (r *Repository) GetOrganizationNamesByIDs(ctx context.Context, ids []uuid.UUID) (map[string]string, error) {
	return r.resolveNamesByIDs(ctx, "organizations", "nomenclature", ids)
}

// GetPositionNamesByIDs resolves position titles for a batch of ids
// (used for from/to_position_name enrichment).
func (r *Repository) GetPositionNamesByIDs(ctx context.Context, ids []uuid.UUID) (map[string]string, error) {
	return r.resolveNamesByIDs(ctx, "positions", "title", ids)
}

// GetEmploymentStatusNamesByIDs resolves employment status names for a batch
// of ids (used for from/to_employment_status_name enrichment).
func (r *Repository) GetEmploymentStatusNamesByIDs(ctx context.Context, ids []uuid.UUID) (map[string]string, error) {
	return r.resolveNamesByIDs(ctx, "employment_statuses", "name", ids)
}

// resolveNamesByIDs is a shared helper that maps reference-table ids to a
// display column (name/nomenclature/title) using a single batch query.
func (r *Repository) resolveNamesByIDs(ctx context.Context, table, nameColumn string, ids []uuid.UUID) (map[string]string, error) {
	result := make(map[string]string, len(ids))
	if len(ids) == 0 {
		return result, nil
	}
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	type row struct {
		ID   string
		Name string
	}
	var rows []row
	err = db.Table(table).
		Where("id IN ?", ids).
		Select("id AS id, " + nameColumn + " AS name").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to resolve %s names: %w", table, err)
	}
	for _, rrow := range rows {
		result[rrow.ID] = rrow.Name
	}
	return result, nil
}

// GetContractNumbersByIDs resolves contract numbers for a batch of contract
// ids (used for previous_contract_number enrichment).
func (r *Repository) GetContractNumbersByIDs(ctx context.Context, ids []uuid.UUID) (map[string]string, error) {
	result := make(map[string]string, len(ids))
	if len(ids) == 0 {
		return result, nil
	}
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	type row struct {
		ID     string
		Number string
	}
	var rows []row
	err = db.Table("employee_contracts").
		Where("id IN ?", ids).
		Select("id AS id, contract_number AS number").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to resolve contract numbers: %w", err)
	}
	for _, rrow := range rows {
		result[rrow.ID] = rrow.Number
	}
	return result, nil
}

func (r *Repository) FindMovementByID(ctx context.Context, id uuid.UUID) (*EmployeeMovement, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var m EmployeeMovement
	if err := db.Where("id = ?", id).First(&m).Error; err != nil {
		return nil, fmt.Errorf("employee movement not found: %w", err)
	}
	return &m, nil
}

func (r *Repository) FindMovementsByEmployeeID(ctx context.Context, employeeID uuid.UUID, page, perPage int) ([]EmployeeMovement, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var movements []EmployeeMovement
	var total int64

	query := db.Model(&EmployeeMovement{}).Where("employee_id = ?", employeeID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count employee movements: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&movements).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list employee movements: %w", err)
	}
	return movements, total, nil
}

func (r *Repository) ListMovements(ctx context.Context, page, perPage int) ([]EmployeeMovement, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var movements []EmployeeMovement
	var total int64

	query := db.Model(&EmployeeMovement{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count movements: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&movements).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list movements: %w", err)
	}
	return movements, total, nil
}

func (r *Repository) UpdateMovement(ctx context.Context, m *EmployeeMovement) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Save(m).Error; err != nil {
		return fmt.Errorf("failed to update employee movement: %w", err)
	}
	return nil
}

func (r *Repository) DeleteMovement(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	result := db.Where("id = ?", id).Delete(&EmployeeMovement{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete employee movement: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("employee movement not found")
	}
	return nil
}

// =========================================================================
// Employee Contract
// =========================================================================

func (r *Repository) CreateContract(ctx context.Context, c *EmployeeContract) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Create(c).Error; err != nil {
		return fmt.Errorf("failed to create employee contract: %w", err)
	}
	return nil
}

func (r *Repository) FindContractByID(ctx context.Context, id uuid.UUID) (*EmployeeContract, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var c EmployeeContract
	if err := db.Where("id = ?", id).First(&c).Error; err != nil {
		return nil, fmt.Errorf("employee contract not found: %w", err)
	}
	return &c, nil
}

func (r *Repository) FindContractsByEmployeeID(ctx context.Context, employeeID uuid.UUID, page, perPage int) ([]EmployeeContract, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var contracts []EmployeeContract
	var total int64

	query := db.Model(&EmployeeContract{}).Where("employee_id = ?", employeeID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count employee contracts: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&contracts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list employee contracts: %w", err)
	}
	return contracts, total, nil
}

func (r *Repository) ListContracts(ctx context.Context, page, perPage int) ([]EmployeeContract, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var contracts []EmployeeContract
	var total int64

	query := db.Model(&EmployeeContract{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count contracts: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&contracts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list contracts: %w", err)
	}
	return contracts, total, nil
}

func (r *Repository) UpdateContract(ctx context.Context, c *EmployeeContract) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Save(c).Error; err != nil {
		return fmt.Errorf("failed to update employee contract: %w", err)
	}
	return nil
}

func (r *Repository) DeleteContract(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	result := db.Where("id = ?", id).Delete(&EmployeeContract{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete employee contract: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("employee contract not found")
	}
	return nil
}

// =========================================================================
// Approval flows
// =========================================================================

// ExecuteMovement menandai movement sebagai executed. Bila toEmploymentID
// tidak nil, to_employment_id ikut dipersist (hasil eksekusi G-1: employment
// baru yang dibuat dari to_* fields movement).
func (r *Repository) ExecuteMovement(ctx context.Context, id uuid.UUID, executedBy uuid.UUID, toEmploymentID *uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	now := time.Now()
	updates := map[string]interface{}{
		"status":      MovementStatusExecuted,
		"executed_by": executedBy.String(),
		"executed_at": now,
	}
	if toEmploymentID != nil {
		updates["to_employment_id"] = toEmploymentID.String()
	}
	result := db.Model(&EmployeeMovement{}).
		Where("id = ? AND status = ?", id, MovementStatusApproved).
		Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to execute movement: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("movement not found or not in approved status")
	}
	return nil
}

func (r *Repository) CancelMovement(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	result := db.Model(&EmployeeMovement{}).
		Where("id = ? AND status IN ?", id, []MovementStatus{MovementStatusDraft, MovementStatusApproved}).
		Update("status", MovementStatusCancelled)
	if result.Error != nil {
		return fmt.Errorf("failed to cancel movement: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("movement not found or already executed/cancelled")
	}
	return nil
}

// ExtendContract membuat kontrak baru sebagai perpanjangan dari kontrak sebelumnya.
func (r *Repository) ExtendContract(ctx context.Context, newContract *EmployeeContract, previousID uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}

	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Set previous contract as extended
	if err := tx.Model(&EmployeeContract{}).
		Where("id = ?", previousID).
		Update("status", ContractStatusExtended).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update previous contract status: %w", err)
	}

	// Create new contract
	if err := tx.Create(newContract).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create new contract: %w", err)
	}

	return tx.Commit().Error
}
