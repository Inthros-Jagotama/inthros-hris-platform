package employeemovement

import (
	"context"
	"errors"
	"fmt"
	"strings"
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

func (r *Repository) ListMovements(ctx context.Context, page, perPage int, movementType, status, search string) ([]EmployeeMovement, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var movements []EmployeeMovement
	var total int64

	query := db.Model(&EmployeeMovement{})
	if movementType != "" {
		query = query.Where("movement_type = ?", movementType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		// Escape LIKE wildcards agar input user tidak diperlakukan sebagai pola.
		escaped := strings.NewReplacer("\\", "\\\\", "%", "\\%", "_", "\\_").Replace(search)
		like := "%" + escaped + "%"
		// Search by decision letter number or employee name (join employees).
		query = query.
			Joins("JOIN employees ON employees.id = employee_movements.employee_id").
			Where("(employee_movements.decision_letter_number LIKE ? OR employees.name LIKE ? OR employees.employee_id LIKE ?)", like, like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count movements: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("employee_movements.created_at DESC").Find(&movements).Error; err != nil {
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

// HasActiveContractByEmployeeID reports whether the employee currently has a
// contract with status = active. Used by the contract_extension movement
// validation (plan G-7: contract_extension wajib merujuk kontrak aktif).
func (r *Repository) HasActiveContractByEmployeeID(ctx context.Context, employeeID uuid.UUID) (bool, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return false, err
	}
	var count int64
	if err := db.Model(&EmployeeContract{}).
		Where("employee_id = ? AND status = ?", employeeID, ContractStatusActive).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to count active contracts: %w", err)
	}
	return count > 0, nil
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

func (r *Repository) ListContracts(ctx context.Context, page, perPage int, status, search string) ([]EmployeeContract, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var contracts []EmployeeContract
	var total int64

	query := db.Model(&EmployeeContract{})
	if status != "" {
		query = query.Where("employee_contracts.status = ?", status)
	}
	if search != "" {
		// Escape LIKE wildcards agar input user tidak diperlakukan sebagai pola.
		escaped := strings.NewReplacer("\\", "\\\\", "%", "\\%", "_", "\\_").Replace(search)
		like := "%" + escaped + "%"
		// Search by contract number or employee name (join employees).
		query = query.
			Joins("JOIN employees ON employees.id = employee_contracts.employee_id").
			Where("(employee_contracts.contract_number LIKE ? OR employees.name LIKE ? OR employees.employee_id LIKE ?)", like, like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count contracts: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("employee_contracts.created_at DESC").Find(&contracts).Error; err != nil {
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
// Movement Audit Trail (plan §12.6)
// =========================================================================

// CreateAudit menyimpan satu baris audit trail movement. Kegagalan di sini
// tidak boleh menggagalkan operasi movement utama (dipanggil best-effort oleh
// service).
func (r *Repository) CreateAudit(ctx context.Context, audit *EmployeeMovementAudit) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Create(audit).Error; err != nil {
		return fmt.Errorf("failed to create employee movement audit: %w", err)
	}
	return nil
}

// ListAuditsByMovementID mengembalikan audit trail satu movement, terurut
// acted_at DESC (baru dulu).
func (r *Repository) ListAuditsByMovementID(ctx context.Context, movementID uuid.UUID, page, perPage int) ([]EmployeeMovementAudit, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var audits []EmployeeMovementAudit
	var total int64

	query := db.Model(&EmployeeMovementAudit{}).Where("movement_id = ?", movementID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count movement audits: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("acted_at DESC, id DESC").Find(&audits).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list movement audits: %w", err)
	}
	return audits, total, nil
}

// =========================================================================
// Career History (plan §12.8) — read model
// =========================================================================

// careerEmploymentRow memuat kolom employments yang dibutuhkan career history.
// Struct lokal (bukan model module employee) karena employeemovement memakai
// narrow-interface pattern: hanya membaca tabel via query langsung, sama
// seperti PositionConflict / EmploymentEffectiveDateConflict.
type careerEmploymentRow struct {
	ID                   uuid.UUID
	OrganizationID       *uuid.UUID
	PositionID           *uuid.UUID
	EmploymentStatusID   *uuid.UUID
	DecisionLetterNumber string
	DecisionLetterDate   string
	EffectiveDate        string
	EffectiveEndDate     *string
}

// FindEmploymentsByEmployeeID mengembalikan seluruh riwayat employment seorang
// karyawan (terurut effective_date ASC) — sumber JOINED + current position
// pada career timeline (plan §12.8).
func (r *Repository) FindEmploymentsByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]careerEmploymentRow, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var rows []careerEmploymentRow
	err = db.Table("employments").
		Where("employee_id = ?", employeeID.String()).
		Order("effective_date ASC, created_at ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list employments for career history: %w", err)
	}
	return rows, nil
}

// FindExecutedMovementsByEmployeeID mengembalikan movement yang sudah
// dieksekusi seorang karyawan (terurut effective_date ASC) — sumber transaksi
// career timeline. Hanya status executed: movement draft/approved belum
// menjadi histori nyata.
func (r *Repository) FindExecutedMovementsByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]EmployeeMovement, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var movements []EmployeeMovement
	err = db.Where("employee_id = ? AND status = ?", employeeID, MovementStatusExecuted).
		Order("effective_date ASC, created_at ASC").
		Find(&movements).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list executed movements for career history: %w", err)
	}
	return movements, nil
}

// FindAllContractsByEmployeeID mengembalikan seluruh kontrak seorang karyawan
// (terurut start_date ASC) — sumber CONTRACT pada career timeline.
func (r *Repository) FindAllContractsByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]EmployeeContract, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var contracts []EmployeeContract
	err = db.Where("employee_id = ?", employeeID).
		Order("start_date ASC, created_at ASC").
		Find(&contracts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list contracts for career history: %w", err)
	}
	return contracts, nil
}

// =========================================================================
// Contract Expiry Management (plan §12.13)
// =========================================================================

// FindContractsExpiringOn mengembalikan kontrak status=active yang berakhir
// tepat pada tanggal yang diberikan (end_date = date) — dipakai scheduler
// reminder H-30/H-14/H-7/H-1 (plan §12.13). Tanggal dibandingkan sebagai
// string YYYY-MM-DD, konsisten dengan format kolom DATE di kedua driver.
func (r *Repository) FindContractsExpiringOn(ctx context.Context, date string) ([]EmployeeContract, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var contracts []EmployeeContract
	err = db.Where("status = ? AND end_date = ?", ContractStatusActive, date).
		Order("end_date ASC").
		Find(&contracts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list contracts expiring on %s: %w", date, err)
	}
	return contracts, nil
}

// FindContractsExpiredBefore mengembalikan kontrak status=active yang sudah
// melewati tanggal akhir (end_date < date) — kandidat dipindah ke status
// expired oleh scheduler (plan §12.13).
func (r *Repository) FindContractsExpiredBefore(ctx context.Context, date string) ([]EmployeeContract, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var contracts []EmployeeContract
	err = db.Where("status = ? AND end_date IS NOT NULL AND end_date < ?", ContractStatusActive, date).
		Order("end_date ASC").
		Find(&contracts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list expired contracts before %s: %w", date, err)
	}
	return contracts, nil
}

// MarkContractsExpired mengubah status kontrak menjadi expired (batch by ids).
// Dipanggil scheduler setelah FindContractsExpiredBefore: hanya kontrak yang
// masih berstatus active yang dipindah (guard ulang di query).
func (r *Repository) MarkContractsExpired(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	result := db.Model(&EmployeeContract{}).
		Where("id IN ? AND status = ?", ids, ContractStatusActive).
		Update("status", ContractStatusExpired)
	if result.Error != nil {
		return fmt.Errorf("failed to mark contracts expired: %w", result.Error)
	}
	return nil
}

// FindUserIDsWithPermission mengembalikan daftar user id yang memiliki sebuah
// permission — langsung (model_has_permissions) atau melalui role
// (model_has_roles + role_has_permissions). Dipakai scheduler untuk mengirim
// reminder kontrak kepada user HR (plan §12.13: "Notification dikirim kepada
// HR"). Query UNION berjalan di MySQL & PostgreSQL.
func (r *Repository) FindUserIDsWithPermission(ctx context.Context, permission string) ([]uuid.UUID, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var ids []string
	err = db.Raw(`
		SELECT DISTINCT mhr.model_id AS id FROM model_has_roles mhr
		JOIN role_has_permissions rhp ON rhp.role_id = mhr.role_id
		JOIN permissions p ON p.id = rhp.permission_id
		WHERE mhr.model_type = 'user' AND p.name = ?
		UNION
		SELECT DISTINCT mhp.model_id AS id FROM model_has_permissions mhp
		JOIN permissions p ON p.id = mhp.permission_id
		WHERE mhp.model_type = 'user' AND p.name = ?
	`, permission, permission).Scan(&ids).Error
	if err != nil {
		return nil, fmt.Errorf("failed to resolve users with permission %s: %w", permission, err)
	}
	result := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		uid, err := uuid.Parse(id)
		if err != nil {
			continue
		}
		result = append(result, uid)
	}
	return result, nil
}

// =========================================================================
// Movement Documents (plan §12.15)
// =========================================================================

// CreateMovementDocument menyimpan metadata satu dokumen movement.
func (r *Repository) CreateMovementDocument(ctx context.Context, d *EmployeeMovementDocument) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Create(d).Error; err != nil {
		return fmt.Errorf("failed to create employee movement document: %w", err)
	}
	return nil
}

// ListDocumentsByMovementID mengembalikan dokumen satu movement, terurut
// created_at DESC (terbaru dulu) dengan pagination.
func (r *Repository) ListDocumentsByMovementID(ctx context.Context, movementID uuid.UUID, page, perPage int) ([]EmployeeMovementDocument, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var documents []EmployeeMovementDocument
	var total int64

	query := db.Model(&EmployeeMovementDocument{}).Where("movement_id = ?", movementID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count movement documents: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC, id DESC").Find(&documents).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list movement documents: %w", err)
	}
	return documents, total, nil
}

// FindMovementDocumentByID mengambil satu dokumen movement berdasarkan id.
func (r *Repository) FindMovementDocumentByID(ctx context.Context, id uuid.UUID) (*EmployeeMovementDocument, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var d EmployeeMovementDocument
	if err := db.Where("id = ?", id).First(&d).Error; err != nil {
		return nil, fmt.Errorf("employee movement document not found: %w", err)
	}
	return &d, nil
}

// DeleteMovementDocument menghapus metadata dokumen movement. RowsAffected 0
// → dokumen tidak ditemukan (atau sudah terhapus via CASCADE).
func (r *Repository) DeleteMovementDocument(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	result := db.Where("id = ?", id).Delete(&EmployeeMovementDocument{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete employee movement document: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("employee movement document not found")
	}
	return nil
}

// =========================================================================
// Approval flows
// =========================================================================

// ExecuteMovementTx executes a movement atomically (enhancement plan §12.2):
// the HR data change supplied by hrChanges (conflict detection + employment /
// employee updates, run through the career executor's Tx variants) and the
// movement's own status update happen on ONE database transaction, so a
// failure anywhere rolls everything back — old employment stays intact and
// the movement stays approved (retryable by HR).
//
// hrChanges receives the transaction and returns the id of the newly created
// employment (nil when the movement creates none, e.g. offboarding). The
// callback runs on the same tx used for the movement status update.
func (r *Repository) ExecuteMovementTx(ctx context.Context, id uuid.UUID, executedBy uuid.UUID, hrChanges func(tx *gorm.DB) (*uuid.UUID, error)) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}

	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	// Reload the movement inside the transaction so the approved-status guard
	// is checked atomically with the HR data change (blocks double-execute).
	var m EmployeeMovement
	if err := tx.Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("employee movement not found")
		}
		return fmt.Errorf("failed to load movement for execution: %w", err)
	}
	if m.Status != MovementStatusApproved {
		return fmt.Errorf("movement not found or not in approved status")
	}

	var toEmploymentID *uuid.UUID
	if hrChanges != nil {
		toEmploymentID, err = hrChanges(tx)
		if err != nil {
			return err
		}
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
	if err := tx.Model(&EmployeeMovement{}).
		Where("id = ? AND status = ?", id, MovementStatusApproved).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to execute movement: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit movement execution: %w", err)
	}
	return nil
}

// PositionConflict reports whether another employee already occupies the
// target position at the given effective date (plan §12.3 — 1 position = 1
// employee). An existing employment conflicts when it is still open at the
// effective date: no effective_end_date (currently or future-planned) or it
// closes on/after the movement's effective date. `excludeEmployeeID` skips the
// movement's own employee (their current employment on the position is not a
// conflict).
//
// `tx` may be nil for the create/update soft-check (a fresh DB connection is
// resolved); ExecuteMovement passes the open transaction so the check runs
// atomically with the rest of the execution.
func (r *Repository) PositionConflict(ctx context.Context, tx *gorm.DB, positionID, excludeEmployeeID uuid.UUID, effectiveDate string) (bool, error) {
	db := tx
	if db == nil {
		var err error
		db, err = r.getDB(ctx)
		if err != nil {
			return false, err
		}
	}
	var count int64
	err := db.Table("employments").
		Where("position_id = ? AND employee_id <> ?", positionID.String(), excludeEmployeeID.String()).
		Where("(effective_end_date IS NULL OR effective_end_date >= ?)", effectiveDate).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check position conflict: %w", err)
	}
	return count > 0, nil
}

// EmploymentEffectiveDateConflict reports whether creating a new open-ended
// employment at `effectiveDate` would overlap an existing employment period
// of the same employee (plan §12.4). Overlap occurs when the employee already
// has an open employment (effective_end_date IS NULL) that starts on/after
// the movement's effective date — either a backdated movement into the
// current employment, or a collision with a future-dated employment created
// by an earlier execution.
//
// `tx` may be nil for reads outside a transaction; ExecuteMovement passes the
// open transaction so the check runs atomically.
func (r *Repository) EmploymentEffectiveDateConflict(ctx context.Context, tx *gorm.DB, employeeID uuid.UUID, effectiveDate string) (bool, error) {
	db := tx
	if db == nil {
		var err error
		db, err = r.getDB(ctx)
		if err != nil {
			return false, err
		}
	}
	var count int64
	err := db.Table("employments").
		Where("employee_id = ? AND effective_end_date IS NULL AND effective_date >= ?", employeeID.String(), effectiveDate).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check employment period conflict: %w", err)
	}
	return count > 0, nil
}

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
// ExtendContract creates a new contract as the extension of a previous one
// inside a single transaction. The new contract's extension_count is derived
// from the previous contract's count (previous + 1) so chained extensions
// accumulate correctly (plan G-6 — previously hardcoded to 1 by the caller).
// The previous contract is marked as status = extended.
func (r *Repository) ExtendContract(ctx context.Context, newContract *EmployeeContract, previousID uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}

	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer tx.Rollback()

	// Load previous contract to derive its extension_count (and confirm it
	// exists — avoid silently extending a deleted/missing contract).
	var previous EmployeeContract
	if err := tx.First(&previous, "id = ?", previousID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("previous contract not found")
		}
		return fmt.Errorf("failed to find previous contract: %w", err)
	}

	// Set previous contract as extended
	if err := tx.Model(&EmployeeContract{}).
		Where("id = ?", previousID).
		Update("status", ContractStatusExtended).Error; err != nil {
		return fmt.Errorf("failed to update previous contract status: %w", err)
	}

	// Derive the new contract's extension count from the chain (G-6).
	newContract.ExtensionCount = previous.ExtensionCount + 1

	// Create new contract
	if err := tx.Create(newContract).Error; err != nil {
		return fmt.Errorf("failed to create new contract: %w", err)
	}

	return tx.Commit().Error
}

// =========================================================================
// Career Path (plan §12.9) — planning/configuration jenjang karier
// =========================================================================

// CreateCareerPathTx menyimpan career path beserta steps-nya dalam satu
// transaksi: path header + seluruh steps harus berhasil atau semuanya batal.
func (r *Repository) CreateCareerPathTx(ctx context.Context, path *CareerPath, steps []CareerPathStep) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}

	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	if err := tx.Create(path).Error; err != nil {
		return fmt.Errorf("failed to create career path: %w", err)
	}
	if len(steps) > 0 {
		if err := tx.Create(&steps).Error; err != nil {
			return fmt.Errorf("failed to create career path steps: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit career path creation: %w", err)
	}
	return nil
}

// ListCareerPaths mengembalikan daftar career path terurut name ASC dengan
// pagination. Keyword (opsional) memfilter berdasarkan substring nama path.
func (r *Repository) ListCareerPaths(ctx context.Context, page, perPage int, keyword string) ([]CareerPath, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var paths []CareerPath
	var total int64

	query := db.Model(&CareerPath{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count career paths: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("name ASC").Find(&paths).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list career paths: %w", err)
	}
	return paths, total, nil
}

// FindCareerPathByID mengambil satu career path berdasarkan id.
func (r *Repository) FindCareerPathByID(ctx context.Context, id uuid.UUID) (*CareerPath, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var p CareerPath
	if err := db.Where("id = ?", id).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("career path not found")
		}
		return nil, fmt.Errorf("failed to find career path: %w", err)
	}
	return &p, nil
}

// ListCareerPathStepsByPathID mengembalikan steps satu career path terurut
// sequence ASC (urutan jenjang dari awal ke akhir).
func (r *Repository) ListCareerPathStepsByPathID(ctx context.Context, pathID uuid.UUID) ([]CareerPathStep, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var steps []CareerPathStep
	if err := db.Where("career_path_id = ?", pathID).Order("sequence ASC").Find(&steps).Error; err != nil {
		return nil, fmt.Errorf("failed to list career path steps: %w", err)
	}
	return steps, nil
}

// UpdateCareerPathTx memperbarui header path dan mengganti SELURUH steps-nya
// dalam satu transaksi (semantik full-replace: steps lama dihapus, steps baru
// dibuat ulang). Aman karena tidak ada referensi eksternal ke career_path_steps.
func (r *Repository) UpdateCareerPathTx(ctx context.Context, path *CareerPath, steps []CareerPathStep) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}

	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	// Header path — use a map so Description can be explicitly cleared (NULL)
	// and IsActive toggled regardless of zero-value semantics in struct updates.
	// Nil di-normalisasi ke nil interface (bukan typed nil pointer) agar GORM
	// konsisten menulis NULL ke kolom description.
	updates := map[string]interface{}{
		"name":      path.Name,
		"is_active": path.IsActive,
	}
	if path.Description != nil {
		updates["description"] = *path.Description
	} else {
		updates["description"] = nil
	}
	if path.UpdatedBy != nil {
		updates["updated_by"] = path.UpdatedBy.String()
	} else {
		updates["updated_by"] = nil
	}
	if err := tx.Model(&CareerPath{}).Where("id = ?", path.ID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update career path: %w", err)
	}

	// Full-replace steps.
	if err := tx.Where("career_path_id = ?", path.ID).Delete(&CareerPathStep{}).Error; err != nil {
		return fmt.Errorf("failed to delete career path steps: %w", err)
	}
	if len(steps) > 0 {
		if err := tx.Create(&steps).Error; err != nil {
			return fmt.Errorf("failed to recreate career path steps: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit career path update: %w", err)
	}
	return nil
}

// DeleteCareerPath menghapus career path beserta steps-nya. Steps dihapus
// eksplisit dulu agar konsisten lintas driver (MySQL/PG CASCADE + SQLite test
// yang tidak selalu mengaktifkan FK).
func (r *Repository) DeleteCareerPath(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}

	if err := db.Where("career_path_id = ?", id).Delete(&CareerPathStep{}).Error; err != nil {
		return fmt.Errorf("failed to delete career path steps: %w", err)
	}
	result := db.Where("id = ?", id).Delete(&CareerPath{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete career path: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("career path not found")
	}
	return nil
}
