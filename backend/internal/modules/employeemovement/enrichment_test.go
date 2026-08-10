package employeemovement

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// setupEnrichmentRefs creates the reference tables (employees, organizations,
// positions, employment_statuses) the enrichment queries JOIN against — the
// test DB only auto-migrates the two module models, so these must be created
// by hand (pola sama approval_integration_test.go).
func setupEnrichmentRefs(t *testing.T, dbResolver func(ctx context.Context) (*gorm.DB, error)) {
	t.Helper()
	db, err := dbResolver(context.Background())
	if err != nil {
		t.Fatalf("failed to get test db: %v", err)
	}
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS employees (
			id CHAR(36) PRIMARY KEY,
			employee_id VARCHAR(50) NOT NULL,
			name VARCHAR(255) NOT NULL,
			status VARCHAR(20) DEFAULT 'active'
		)`,
		`CREATE TABLE IF NOT EXISTS organizations (
			id CHAR(36) PRIMARY KEY,
			code VARCHAR(2) NOT NULL,
			full_code VARCHAR(50) NOT NULL,
			nomenclature VARCHAR(255) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS positions (
			id CHAR(36) PRIMARY KEY,
			organization_id CHAR(36) NOT NULL,
			code VARCHAR(50) NOT NULL,
			title VARCHAR(200) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS employment_statuses (
			id CHAR(36) PRIMARY KEY,
			code VARCHAR(20) NOT NULL,
			name VARCHAR(100) NOT NULL
		)`,
	}
	for _, s := range stmts {
		if err := db.Exec(s).Error; err != nil {
			t.Fatalf("failed to create reference table: %v\nSQL: %s", err, s)
		}
	}
}

// seedEnrichmentRefs inserts one employee, one organization, one position and
// one employment status; returns their ids.
func seedEnrichmentRefs(t *testing.T, dbResolver func(ctx context.Context) (*gorm.DB, error)) (empID, orgID, posID, statusID uuid.UUID) {
	t.Helper()
	db, err := dbResolver(context.Background())
	if err != nil {
		t.Fatalf("failed to get test db: %v", err)
	}
	empID, orgID, posID, statusID = uuid.New(), uuid.New(), uuid.New(), uuid.New()
	if err := db.Exec(`INSERT INTO employees (id, employee_id, name) VALUES (?, ?, ?)`,
		empID.String(), "EMP-001", "Budi Santoso").Error; err != nil {
		t.Fatalf("failed to seed employee: %v", err)
	}
	if err := db.Exec(`INSERT INTO organizations (id, code, full_code, nomenclature) VALUES (?, ?, ?, ?)`,
		orgID.String(), "IT", "IT-DIV", "Divisi Teknologi Informasi").Error; err != nil {
		t.Fatalf("failed to seed organization: %v", err)
	}
	if err := db.Exec(`INSERT INTO positions (id, organization_id, code, title) VALUES (?, ?, ?, ?)`,
		posID.String(), orgID.String(), "P-IT-01", "Software Engineer").Error; err != nil {
		t.Fatalf("failed to seed position: %v", err)
	}
	if err := db.Exec(`INSERT INTO employment_statuses (id, code, name) VALUES (?, ?, ?)`,
		statusID.String(), "PERM", "Permanent").Error; err != nil {
		t.Fatalf("failed to seed employment status: %v", err)
	}
	return empID, orgID, posID, statusID
}

// TestService_ListMovements_Enriched verifies G-4: list responses carry
// employee/organization/position/status display names via batch JOIN.
func TestService_ListMovements_Enriched(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	setupEnrichmentRefs(t, dbResolver)

	empID, orgID, posID, statusID := seedEnrichmentRefs(t, dbResolver)

	// Movement with to_* references populated.
	m := &EmployeeMovement{
		EmployeeID:           empID,
		MovementType:         MovementTypePromotion,
		DecisionLetterNumber: "SK-ENRICH-1",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-01",
		Status:               MovementStatusDraft,
		ToOrganizationID:     &orgID,
		ToPositionID:         &posID,
		ToEmploymentStatusID: &statusID,
	}
	if err := db.Create(m).Error; err != nil {
		t.Fatalf("failed to create movement: %v", err)
	}

	repo := NewRepository(dbResolver)
	logger := testLogger()
	svc := NewService(repo, logger)

	resp, err := svc.ListMovements(ctx(), 1, 20, "", "", "")
	if err != nil {
		t.Fatalf("ListMovements failed: %v", err)
	}
	if resp.Total != 1 {
		t.Fatalf("expected 1 movement, got %d", resp.Total)
	}
	items := resp.Data.([]MovementResponse)
	item := items[0]

	if item.EmployeeName != "Budi Santoso" {
		t.Errorf("expected employee_name 'Budi Santoso', got '%s'", item.EmployeeName)
	}
	if item.EmployeeCode != "EMP-001" {
		t.Errorf("expected employee_code 'EMP-001', got '%s'", item.EmployeeCode)
	}
	if item.ToOrganizationName != "Divisi Teknologi Informasi" {
		t.Errorf("expected to_organization_name 'Divisi Teknologi Informasi', got '%s'", item.ToOrganizationName)
	}
	if item.ToPositionName != "Software Engineer" {
		t.Errorf("expected to_position_name 'Software Engineer', got '%s'", item.ToPositionName)
	}
	if item.ToEmploymentStatusName != "Permanent" {
		t.Errorf("expected to_employment_status_name 'Permanent', got '%s'", item.ToEmploymentStatusName)
	}
}

// TestService_GetMovementByID_Enriched verifies G-4 single get enrichment.
func TestService_GetMovementByID_Enriched(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	setupEnrichmentRefs(t, dbResolver)

	empID, orgID, posID, statusID := seedEnrichmentRefs(t, dbResolver)

	m := &EmployeeMovement{
		EmployeeID:             empID,
		MovementType:           MovementTypePromotion,
		DecisionLetterNumber:   "SK-ENRICH-2",
		DecisionLetterDate:     "2026-07-01",
		EffectiveDate:          "2026-08-01",
		Status:                 MovementStatusDraft,
		FromOrganizationID:     &orgID,
		FromPositionID:         &posID,
		FromEmploymentStatusID: &statusID,
	}
	if err := db.Create(m).Error; err != nil {
		t.Fatalf("failed to create movement: %v", err)
	}

	repo := NewRepository(dbResolver)
	svc := NewService(repo, testLogger())

	item, err := svc.GetMovementByID(ctx(), m.ID.String())
	if err != nil {
		t.Fatalf("GetMovementByID failed: %v", err)
	}
	if item.EmployeeName != "Budi Santoso" {
		t.Errorf("expected employee_name 'Budi Santoso', got '%s'", item.EmployeeName)
	}
	if item.FromOrganizationName != "Divisi Teknologi Informasi" {
		t.Errorf("expected from_organization_name 'Divisi Teknologi Informasi', got '%s'", item.FromOrganizationName)
	}
	if item.FromPositionName != "Software Engineer" {
		t.Errorf("expected from_position_name 'Software Engineer', got '%s'", item.FromPositionName)
	}
	if item.FromEmploymentStatusName != "Permanent" {
		t.Errorf("expected from_employment_status_name 'Permanent', got '%s'", item.FromEmploymentStatusName)
	}
}

// TestService_ListContracts_Enriched verifies G-4 contract enrichment:
// employee name/code + previous contract number.
func TestService_ListContracts_Enriched(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	setupEnrichmentRefs(t, dbResolver)

	empID, _, _, _ := seedEnrichmentRefs(t, dbResolver)

	prev := &EmployeeContract{
		EmployeeID:     empID,
		ContractNumber: "CTR-001",
		ContractType:   ContractTypePKWT,
		StartDate:      "2026-01-01",
		Status:         ContractStatusActive,
	}
	if err := db.Create(prev).Error; err != nil {
		t.Fatalf("failed to create previous contract: %v", err)
	}

	cur := &EmployeeContract{
		EmployeeID:         empID,
		ContractNumber:     "CTR-002",
		ContractType:       ContractTypePKWT,
		StartDate:          "2026-06-01",
		Status:             ContractStatusActive,
		PreviousContractID: &prev.ID,
	}
	if err := db.Create(cur).Error; err != nil {
		t.Fatalf("failed to create current contract: %v", err)
	}

	repo := NewRepository(dbResolver)
	svc := NewService(repo, testLogger())

	resp, err := svc.ListContracts(ctx(), 1, 20, "", "")
	if err != nil {
		t.Fatalf("ListContracts failed: %v", err)
	}
	if resp.Total != 2 {
		t.Fatalf("expected 2 contracts, got %d", resp.Total)
	}
	items := resp.Data.([]ContractResponse)

	// Find the one with previous_contract_id.
	var current *ContractResponse
	for i := range items {
		if items[i].PreviousContractID != nil {
			current = &items[i]
			break
		}
	}
	if current == nil {
		t.Fatal("expected a contract with previous_contract_id")
	}
	if current.EmployeeName != "Budi Santoso" {
		t.Errorf("expected employee_name 'Budi Santoso', got '%s'", current.EmployeeName)
	}
	if current.EmployeeCode != "EMP-001" {
		t.Errorf("expected employee_code 'EMP-001', got '%s'", current.EmployeeCode)
	}
	if current.PreviousContractNumber != prev.ContractNumber {
		t.Errorf("expected previous_contract_number '%s', got '%s'", prev.ContractNumber, current.PreviousContractNumber)
	}
}

// TestService_CreateMovement_SnapshotPersisted verifies plan §12.5: the
// display names resolved at creation time are persisted onto the movement
// row (from_*/to_*_name columns), not just filled into the response.
func TestService_CreateMovement_SnapshotPersisted(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	setupEnrichmentRefs(t, dbResolver)

	empID, orgID, posID, statusID := seedEnrichmentRefs(t, dbResolver)

	repo := NewRepository(dbResolver)
	svc := NewService(repo, testLogger())

	resp, err := svc.CreateMovement(ctx(), CreateMovementRequest{
		EmployeeID:           empID.String(),
		MovementType:         string(MovementTypePromotion),
		ToOrganizationID:     strPtr(orgID.String()),
		ToPositionID:         strPtr(posID.String()),
		ToEmploymentStatusID: strPtr(statusID.String()),
		DecisionLetterNumber: "SK-SNAP-1",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-01",
	})
	if err != nil {
		t.Fatalf("CreateMovement failed: %v", err)
	}

	// Response carries the snapshot names.
	if resp.ToOrganizationName != "Divisi Teknologi Informasi" {
		t.Errorf("expected to_organization_name snapshot in response, got '%s'", resp.ToOrganizationName)
	}
	if resp.ToPositionName != "Software Engineer" {
		t.Errorf("expected to_position_name snapshot in response, got '%s'", resp.ToPositionName)
	}
	if resp.ToEmploymentStatusName != "Permanent" {
		t.Errorf("expected to_employment_status_name snapshot in response, got '%s'", resp.ToEmploymentStatusName)
	}

	// The row itself must have the snapshot persisted.
	var stored EmployeeMovement
	if err := db.Where("id = ?", resp.ID).First(&stored).Error; err != nil {
		t.Fatalf("failed to reload movement: %v", err)
	}
	if stored.ToOrganizationName != "Divisi Teknologi Informasi" {
		t.Errorf("expected persisted to_organization_name, got '%s'", stored.ToOrganizationName)
	}
	if stored.ToPositionName != "Software Engineer" {
		t.Errorf("expected persisted to_position_name, got '%s'", stored.ToPositionName)
	}
	if stored.ToEmploymentStatusName != "Permanent" {
		t.Errorf("expected persisted to_employment_status_name, got '%s'", stored.ToEmploymentStatusName)
	}
}

// TestService_Snapshot_ImmutableOnMasterRename verifies plan §12.5: renaming
// the master position after the movement was created does NOT rewrite the
// movement's history — the persisted snapshot name is returned.
func TestService_Snapshot_ImmutableOnMasterRename(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	setupEnrichmentRefs(t, dbResolver)

	empID, orgID, posID, statusID := seedEnrichmentRefs(t, dbResolver)

	repo := NewRepository(dbResolver)
	svc := NewService(repo, testLogger())

	created, err := svc.CreateMovement(ctx(), CreateMovementRequest{
		EmployeeID:           empID.String(),
		MovementType:         string(MovementTypePromotion),
		ToOrganizationID:     strPtr(orgID.String()),
		ToPositionID:         strPtr(posID.String()),
		ToEmploymentStatusID: strPtr(statusID.String()),
		DecisionLetterNumber: "SK-SNAP-2",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-01",
	})
	if err != nil {
		t.Fatalf("CreateMovement failed: %v", err)
	}

	// Rename the master position afterwards (the scenario §12.5 protects).
	if err := db.Exec(`UPDATE positions SET title = ? WHERE id = ?`, "Principal Engineer", posID.String()).Error; err != nil {
		t.Fatalf("failed to rename position: %v", err)
	}

	got, err := svc.GetMovementByID(ctx(), created.ID)
	if err != nil {
		t.Fatalf("GetMovementByID failed: %v", err)
	}
	if got.ToPositionName != "Software Engineer" {
		t.Errorf("expected snapshot name 'Software Engineer' after master rename, got '%s'", got.ToPositionName)
	}
}

// TestService_ListMovements_NoRefs_NoError verifies enrichment is a no-op
// (not an error) when reference rows don't exist in the DB.
func TestService_ListMovements_NoRefs_NoError(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	empID := uuid.New()
	m := &EmployeeMovement{
		EmployeeID:           empID,
		MovementType:         MovementTypePromotion,
		DecisionLetterNumber: "SK-NOREFS",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-01",
		Status:               MovementStatusDraft,
	}
	if err := svc.repo.CreateMovement(ctx(), m); err != nil {
		t.Fatalf("failed to create movement: %v", err)
	}

	resp, err := svc.ListMovements(ctx(), 1, 20, "", "", "")
	if err != nil {
		t.Fatalf("ListMovements failed: %v", err)
	}
	items := resp.Data.([]MovementResponse)
	if len(items) != 1 {
		t.Fatalf("expected 1 movement, got %d", len(items))
	}
	if items[0].EmployeeName != "" {
		t.Errorf("expected empty employee_name when no refs, got '%s'", items[0].EmployeeName)
	}
}
