package employeemovement

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/inthros/hris-platform/internal/modules/employee"
)

// seedCareerReferenceTables membuat + mengisi tabel master minimal yang
// diperlukan enrichment career history (organizations, positions,
// employment_statuses, employees). setupTestDB hanya meng-automigrate model
// movement/employment — tabel referensi sengaja dibuat di sini dengan skema
// sekecil mungkin agar jalur resolve-names (G-4) teruji penuh.
//
// employeeID dipakai juga untuk seed tabel employees sehingga enrichment nama
// karyawan ketemu (id yang sama dengan id karyawan movement).
func seedCareerReferenceTables(t *testing.T, repo *Repository, employeeID uuid.UUID) (orgID, posID, statusID uuid.UUID) {
	t.Helper()
	db, err := repo.getDB(ctx())
	if err != nil {
		t.Fatalf("failed to get test db: %v", err)
	}
	for _, stmt := range []string{
		"CREATE TABLE IF NOT EXISTS organizations (id CHAR(36) PRIMARY KEY, nomenclature VARCHAR(255))",
		"CREATE TABLE IF NOT EXISTS positions (id CHAR(36) PRIMARY KEY, title VARCHAR(255))",
		"CREATE TABLE IF NOT EXISTS employment_statuses (id CHAR(36) PRIMARY KEY, name VARCHAR(255))",
		// religions & marital_statuses + kolom employees yang dibutuhkan
		// GetEmployeeProfile (Generate Document, variable {{employee.*}}).
		"CREATE TABLE IF NOT EXISTS religions (id CHAR(36) PRIMARY KEY, name VARCHAR(200))",
		"CREATE TABLE IF NOT EXISTS marital_statuses (id CHAR(36) PRIMARY KEY, name VARCHAR(100))",
		"CREATE TABLE IF NOT EXISTS employees (id CHAR(36) PRIMARY KEY, name VARCHAR(255), employee_id VARCHAR(50), nik VARCHAR(16), family_id VARCHAR(16), mother_name VARCHAR(255), gender VARCHAR(10), nationality_type VARCHAR(10), nationality_id CHAR(2), passport VARCHAR(50), pob VARCHAR(255), dob DATE, phone_number VARCHAR(255), email VARCHAR(255), linkedin VARCHAR(255), ig VARCHAR(255), religion_id CHAR(36), marital_status_id CHAR(36), status VARCHAR(20) DEFAULT 'active')",
	} {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("failed to create reference table: %v", err)
		}
	}
	orgID = uuid.New()
	posID = uuid.New()
	statusID = uuid.New()
	religionID := uuid.New()
	maritalStatusID := uuid.New()
	if err := db.Exec("INSERT INTO organizations (id, nomenclature) VALUES (?, ?)", orgID.String(), "PT Maju Bersama").Error; err != nil {
		t.Fatalf("failed to seed organization: %v", err)
	}
	// position_id sebenarnya organizations.id (Organization = Position) --
	// GetPositionNamesByIDs sekarang resolve dari organizations, bukan tabel
	// positions yang mati; seed posID sebagai baris organizations terpisah.
	if err := db.Exec("INSERT INTO organizations (id, nomenclature) VALUES (?, ?)", posID.String(), "Software Engineer").Error; err != nil {
		t.Fatalf("failed to seed position: %v", err)
	}
	if err := db.Exec("INSERT INTO employment_statuses (id, name) VALUES (?, ?)", statusID.String(), "Permanent").Error; err != nil {
		t.Fatalf("failed to seed employment status: %v", err)
	}
	if err := db.Exec("INSERT INTO religions (id, name) VALUES (?, ?)", religionID.String(), "Islam").Error; err != nil {
		t.Fatalf("failed to seed religion: %v", err)
	}
	if err := db.Exec("INSERT INTO marital_statuses (id, name) VALUES (?, ?)", maritalStatusID.String(), "Menikah").Error; err != nil {
		t.Fatalf("failed to seed marital status: %v", err)
	}
	if err := db.Exec("INSERT INTO employees (id, name, employee_id, religion_id, marital_status_id, dob) VALUES (?, ?, ?, ?, ?, ?)",
		employeeID.String(), "Test Karyawan", "EMP-TL-001", religionID.String(), maritalStatusID.String(), "1990-01-01").Error; err != nil {
		t.Fatalf("failed to seed employee: %v", err)
	}
	return orgID, posID, statusID
}

// TestService_GetCareerHistory verifies plan §12.8: timeline dibangun dari
// employment (JOINED), movement executed (MOVEMENT), dan contract (CONTRACT),
// terurut kronologis, lengkap dengan current position dari employment aktif.
func TestService_GetCareerHistory(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	orgID, posID, statusID := seedCareerReferenceTables(t, repo, employeeID)

	// Employment pertama (JOINED) — terbuka (aktif).
	first := &employee.Employment{
		EmployeeID:           &employeeID,
		OrganizationID:       &orgID,
		PositionID:           &posID,
		EmploymentStatusID:   &statusID,
		DecisionLetterNumber: "SK-TL-001",
		DecisionLetterDate:   "2024-01-10",
		EffectiveDate:        "2024-01-01",
	}
	seedEmployment(t, repo, first)

	// Movement executed: mutation 2025 dengan snapshot names.
	mov := createTestMovement(repo, employeeID)
	mov.MovementType = MovementTypeMutation
	mov.Status = MovementStatusExecuted
	mov.EffectiveDate = "2025-03-01"
	mov.FromPositionName = "Software Engineer"
	mov.ToPositionName = "Senior Software Engineer"
	if err := repo.UpdateMovement(ctx(), mov); err != nil {
		t.Fatalf("failed to set executed movement: %v", err)
	}

	// Contract mulai 2024 (periode 2024-2025).
	contract := createTestContract(repo, employeeID)
	contract.StartDate = "2024-01-01"
	contract.EndDate = strPtr("2025-12-31")
	if err := repo.UpdateContract(ctx(), contract); err != nil {
		t.Fatalf("failed to update contract: %v", err)
	}

	resp, err := svc.GetCareerHistory(ctx(), employeeID.String())
	if err != nil {
		t.Fatalf("GetCareerHistory failed: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Data.EmployeeID != employeeID.String() {
		t.Errorf("expected employee_id %s, got %s", employeeID.String(), resp.Data.EmployeeID)
	}
	if resp.Data.EmployeeName != "Test Karyawan" {
		t.Errorf("expected employee_name 'Test Karyawan', got '%s'", resp.Data.EmployeeName)
	}

	// Timeline: JOINED (2024-01-01) → CONTRACT (2024-01-01, setelah JOINED pd
	// tanggal sama) → MOVEMENT (2025-03-01).
	if len(resp.Data.Timeline) != 3 {
		t.Fatalf("expected 3 timeline entries, got %d: %+v", len(resp.Data.Timeline), resp.Data.Timeline)
	}
	joined := resp.Data.Timeline[0]
	if joined.EventType != "JOINED" || joined.Date != "2024-01-01" {
		t.Errorf("expected JOINED at 2024-01-01, got %s at %s", joined.EventType, joined.Date)
	}
	if joined.Title != "Software Engineer" {
		t.Errorf("expected JOINED title 'Software Engineer', got '%s'", joined.Title)
	}
	if joined.EmploymentID == nil || *joined.EmploymentID != first.ID.String() {
		t.Errorf("expected JOINED employment_id %s, got %v", first.ID.String(), joined.EmploymentID)
	}
	contractEntry := resp.Data.Timeline[1]
	if contractEntry.EventType != "CONTRACT" || contractEntry.Date != "2024-01-01" {
		t.Errorf("expected CONTRACT at 2024-01-01, got %s at %s", contractEntry.EventType, contractEntry.Date)
	}
	if contractEntry.ContractID == nil || *contractEntry.ContractID != contract.ID.String() {
		t.Errorf("expected CONTRACT contract_id %s, got %v", contract.ID.String(), contractEntry.ContractID)
	}
	movement := resp.Data.Timeline[2]
	if movement.EventType != "MOVEMENT" || movement.Date != "2025-03-01" {
		t.Errorf("expected MOVEMENT at 2025-03-01, got %s at %s", movement.EventType, movement.Date)
	}
	if movement.MovementType == nil || *movement.MovementType != string(MovementTypeMutation) {
		t.Errorf("expected movement_type mutation, got %v", movement.MovementType)
	}
	if movement.Description == nil || *movement.Description != "Software Engineer → Senior Software Engineer" {
		t.Errorf("expected movement description from→to, got %v", movement.Description)
	}
	if movement.MovementID == nil || *movement.MovementID != mov.ID.String() {
		t.Errorf("expected MOVEMENT movement_id %s, got %v", mov.ID.String(), movement.MovementID)
	}

	// Current position: employment pertama masih terbuka → posisi diisi.
	if resp.Data.CurrentPosition == nil {
		t.Fatal("expected current_position to be set")
	}
	if resp.Data.CurrentPosition.PositionName != "Software Engineer" {
		t.Errorf("expected current_position position 'Software Engineer', got '%s'", resp.Data.CurrentPosition.PositionName)
	}
	if resp.Data.CurrentPosition.OrganizationName != "PT Maju Bersama" {
		t.Errorf("expected current_position organization 'PT Maju Bersama', got '%s'", resp.Data.CurrentPosition.OrganizationName)
	}
	if resp.Data.CurrentPosition.EmploymentStatusName != "Permanent" {
		t.Errorf("expected current_position status 'Permanent', got '%s'", resp.Data.CurrentPosition.EmploymentStatusName)
	}
}

// TestService_GetCareerHistory_OnlyExecutedMovements verifies draft/pending
// movements tidak masuk timeline (hanya status executed).
func TestService_GetCareerHistory_OnlyExecutedMovements(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	seedCareerReferenceTables(t, repo, employeeID)

	seedEmployment(t, repo, &employee.Employment{
		EmployeeID:           &employeeID,
		DecisionLetterNumber: "SK-TL-ONLY",
		DecisionLetterDate:   "2024-01-10",
		EffectiveDate:        "2024-01-01",
	})

	// Draft movement — TIDAK boleh masuk timeline.
	draft := createTestMovement(repo, employeeID)
	draft.MovementType = MovementTypePromotion
	draft.EffectiveDate = "2025-06-01"
	if err := repo.UpdateMovement(ctx(), draft); err != nil {
		t.Fatalf("failed to update draft movement: %v", err)
	}

	resp, err := svc.GetCareerHistory(ctx(), employeeID.String())
	if err != nil {
		t.Fatalf("GetCareerHistory failed: %v", err)
	}
	if len(resp.Data.Timeline) != 1 {
		t.Fatalf("expected 1 timeline entry (JOINED only), got %d: %+v", len(resp.Data.Timeline), resp.Data.Timeline)
	}
	if resp.Data.Timeline[0].EventType != "JOINED" {
		t.Errorf("expected only JOINED event, got '%s'", resp.Data.Timeline[0].EventType)
	}
}

// TestService_GetCareerHistory_CurrentPositionPrefersOpenEmployment verifies
// current position memilih employment terbuka (effective_end_date NULL) dengan
// tanggal efektif terbesar, bukan sekadar employment terakhir.
func TestService_GetCareerHistory_CurrentPositionPrefersOpenEmployment(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	orgID, posID, statusID := seedCareerReferenceTables(t, repo, employeeID)

	// Employment 1: 2024-01-01 — ditutup 2024-12-31 (bukan current).
	seedEmployment(t, repo, &employee.Employment{
		EmployeeID:           &employeeID,
		OrganizationID:       &orgID,
		PositionID:           &posID,
		EmploymentStatusID:   &statusID,
		DecisionLetterNumber: "SK-TL-OLD",
		DecisionLetterDate:   "2024-01-10",
		EffectiveDate:        "2024-01-01",
		EffectiveEndDate:     strPtr("2024-12-31"),
	})
	// Employment 2: 2025-01-01 — masih terbuka (current).
	seedEmployment(t, repo, &employee.Employment{
		EmployeeID:           &employeeID,
		OrganizationID:       &orgID,
		PositionID:           &posID,
		EmploymentStatusID:   &statusID,
		DecisionLetterNumber: "SK-TL-NEW",
		DecisionLetterDate:   "2025-01-10",
		EffectiveDate:        "2025-01-01",
	})

	resp, err := svc.GetCareerHistory(ctx(), employeeID.String())
	if err != nil {
		t.Fatalf("GetCareerHistory failed: %v", err)
	}
	if resp.Data.CurrentPosition == nil {
		t.Fatal("expected current_position to be set")
	}
	if resp.Data.CurrentPosition.EffectiveDate != "2025-01-01" {
		t.Errorf("expected current_position effective_date 2025-01-01 (open), got '%s'", resp.Data.CurrentPosition.EffectiveDate)
	}
}

// TestHandler_GetCareerHistory verifies GET /employees/:id/career-history
// endpoint returns 200 with the read model.
func TestHandler_GetCareerHistory(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	employeeID := uuid.New()
	seedCareerReferenceTables(t, repo, employeeID)

	seedEmployment(t, repo, &employee.Employment{
		EmployeeID:           &employeeID,
		DecisionLetterNumber: "SK-TL-H",
		DecisionLetterDate:   "2024-01-10",
		EffectiveDate:        "2024-01-01",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/employee-movements/employees/"+employeeID.String()+"/career-history", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Success bool `json:"success"`
		Data    struct {
			EmployeeID string `json:"employee_id"`
			Timeline   []struct {
				EventType string `json:"event_type"`
			} `json:"timeline"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Data.EmployeeID != employeeID.String() {
		t.Errorf("expected employee_id %s, got %s", employeeID.String(), resp.Data.EmployeeID)
	}
	if len(resp.Data.Timeline) != 1 || resp.Data.Timeline[0].EventType != "JOINED" {
		t.Errorf("expected 1 JOINED entry, got %+v", resp.Data.Timeline)
	}
}
