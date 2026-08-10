package employeemovement

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

// =========================================================================
// Repository Report Tests (plan §12.17)
// =========================================================================

// seedReportMovements creates movements across types/statuses/orgs/positions
// so aggregation queries have deterministic input.
func seedReportMovements(t *testing.T, repo *Repository, ctx context.Context) (orgID, posID, empID uuid.UUID) {
	t.Helper()
	orgID = uuid.New()
	posID = uuid.New()
	empID = uuid.New()

	makeMovement := func(mt MovementType, status MovementStatus, date string, toOrg, toPos *uuid.UUID) *EmployeeMovement {
		m := &EmployeeMovement{
			EmployeeID:           empID,
			MovementType:         mt,
			DecisionLetterNumber: uuid.New().String(),
			DecisionLetterDate:   "2026-07-01",
			EffectiveDate:        date,
			Status:               status,
			ToOrganizationID:     toOrg,
			ToPositionID:         toPos,
		}
		if err := repo.CreateMovement(ctx, m); err != nil {
			t.Fatalf("failed to create report movement: %v", err)
		}
		return m
	}

	// 2 promotion (executed + approved), 1 mutation, 1 offboarding.
	makeMovement(MovementTypePromotion, MovementStatusExecuted, "2026-01-15", &orgID, &posID)
	makeMovement(MovementTypePromotion, MovementStatusApproved, "2026-02-20", &orgID, &posID)
	makeMovement(MovementTypeMutation, MovementStatusExecuted, "2026-03-10", &orgID, &posID)
	makeMovement(MovementTypeOffboarding, MovementStatusExecuted, "2026-05-05", nil, nil)
	return orgID, posID, empID
}

func TestRepo_CountMovementsByType(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	orgID, _, _ := seedReportMovements(t, repo, ctx)

	// Without filters — all types present.
	byType, err := repo.CountMovementsByType(ctx, MovementReportFilter{})
	if err != nil {
		t.Fatalf("CountMovementsByType failed: %v", err)
	}
	if byType["promotion"] != 2 {
		t.Errorf("expected 2 promotions, got %d", byType["promotion"])
	}
	if byType["mutation"] != 1 {
		t.Errorf("expected 1 mutation, got %d", byType["mutation"])
	}
	if byType["offboarding"] != 1 {
		t.Errorf("expected 1 offboarding, got %d", byType["offboarding"])
	}

	// With period filter (effective_date range) — only movements in range.
	byType, err = repo.CountMovementsByType(ctx, MovementReportFilter{DateFrom: "2026-02-01", DateTo: "2026-04-30"})
	if err != nil {
		t.Fatalf("CountMovementsByType(period) failed: %v", err)
	}
	if byType["promotion"] != 1 {
		t.Errorf("expected 1 promotion in Feb-Apr, got %d", byType["promotion"])
	}
	if byType["offboarding"] != 0 {
		t.Errorf("expected 0 offboarding in Feb-Apr, got %d", byType["offboarding"])
	}

	// With organization filter — matches to OR from.
	byType, err = repo.CountMovementsByType(ctx, MovementReportFilter{OrganizationID: &orgID})
	if err != nil {
		t.Fatalf("CountMovementsByType(org) failed: %v", err)
	}
	if byType["offboarding"] != 0 {
		t.Errorf("expected 0 offboarding for org filter (no to_org), got %d", byType["offboarding"])
	}
	if byType["promotion"] != 2 {
		t.Errorf("expected 2 promotions for org filter, got %d", byType["promotion"])
	}
}

func TestRepo_CountMovementsByStatus(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	seedReportMovements(t, repo, ctx)

	byStatus, err := repo.CountMovementsByStatus(ctx, MovementReportFilter{})
	if err != nil {
		t.Fatalf("CountMovementsByStatus failed: %v", err)
	}
	if byStatus["executed"] != 3 {
		t.Errorf("expected 3 executed, got %d", byStatus["executed"])
	}
	if byStatus["approved"] != 1 {
		t.Errorf("expected 1 approved, got %d", byStatus["approved"])
	}

	// Status filter narrows to that status only.
	byStatus, err = repo.CountMovementsByStatus(ctx, MovementReportFilter{Status: "approved"})
	if err != nil {
		t.Fatalf("CountMovementsByStatus(status) failed: %v", err)
	}
	if byStatus["approved"] != 1 {
		t.Errorf("expected 1 approved with filter, got %d", byStatus["approved"])
	}
}

func TestRepo_CountContractsByStatus(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	empID := uuid.New()
	createTestContract(repo, empID) // active
	c2 := createTestContract(repo, empID)
	c3 := createTestContract(repo, empID)
	// Mark c2 expired, c3 terminated via direct update.
	db, err := dbResolver(ctx)
	if err != nil {
		t.Fatalf("dbResolver failed: %v", err)
	}
	if err := db.Model(&EmployeeContract{}).Where("id = ?", c2.ID).Update("status", ContractStatusExpired).Error; err != nil {
		t.Fatalf("update expired failed: %v", err)
	}
	if err := db.Model(&EmployeeContract{}).Where("id = ?", c3.ID).Update("status", ContractStatusTerminated).Error; err != nil {
		t.Fatalf("update terminated failed: %v", err)
	}

	byStatus, err := repo.CountContractsByStatus(ctx)
	if err != nil {
		t.Fatalf("CountContractsByStatus failed: %v", err)
	}
	if byStatus[string(ContractStatusActive)] != 1 {
		t.Errorf("expected 1 active, got %d", byStatus[string(ContractStatusActive)])
	}
	if byStatus[string(ContractStatusExpired)] != 1 {
		t.Errorf("expected 1 expired, got %d", byStatus[string(ContractStatusExpired)])
	}
	if byStatus[string(ContractStatusTerminated)] != 1 {
		t.Errorf("expected 1 terminated, got %d", byStatus[string(ContractStatusTerminated)])
	}
}

func TestRepo_CountExpiringContracts(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	empID := uuid.New()
	// Active contract expiring within window.
	expiring := &EmployeeContract{
		EmployeeID:     empID,
		ContractNumber: "CTR-EXP-1",
		ContractType:   ContractTypePKWT,
		StartDate:      "2026-01-01",
		EndDate:        strPtr("2026-08-20"),
		Status:         ContractStatusActive,
	}
	if err := repo.CreateContract(ctx, expiring); err != nil {
		t.Fatalf("CreateContract failed: %v", err)
	}
	// Active contract expiring outside window (ignored).
	late := &EmployeeContract{
		EmployeeID:     empID,
		ContractNumber: "CTR-LATE-1",
		ContractType:   ContractTypePKWT,
		StartDate:      "2026-01-01",
		EndDate:        strPtr("2027-01-01"),
		Status:         ContractStatusActive,
	}
	if err := repo.CreateContract(ctx, late); err != nil {
		t.Fatalf("CreateContract failed: %v", err)
	}
	// Expired contract with end_date in window — status active only counts.
	expiredInWindow := &EmployeeContract{
		EmployeeID:     empID,
		ContractNumber: "CTR-EXPD-1",
		ContractType:   ContractTypePKWT,
		StartDate:      "2026-01-01",
		EndDate:        strPtr("2026-08-10"),
		Status:         ContractStatusExpired,
	}
	if err := repo.CreateContract(ctx, expiredInWindow); err != nil {
		t.Fatalf("CreateContract failed: %v", err)
	}

	count, err := repo.CountExpiringContracts(ctx, "2026-08-01", "2026-08-31")
	if err != nil {
		t.Fatalf("CountExpiringContracts failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 expiring contract, got %d", count)
	}
}

// =========================================================================
// Service Report Tests (plan §12.17)
// =========================================================================

func TestService_GetMovementReport(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()
	svc := NewService(repo, testLogger())

	orgID, posID, _ := seedReportMovements(t, repo, ctx)

	// Full report (no filters).
	resp, err := svc.GetMovementReport(ctx, "", "", "", "", "", "", "")
	if err != nil {
		t.Fatalf("GetMovementReport failed: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Data.Total != 4 {
		t.Errorf("expected total 4, got %d", resp.Data.Total)
	}
	if resp.Data.ByType["promotion"] != 2 {
		t.Errorf("expected 2 promotions, got %d", resp.Data.ByType["promotion"])
	}
	if resp.Data.ByStatus["executed"] != 3 {
		t.Errorf("expected 3 executed, got %d", resp.Data.ByStatus["executed"])
	}

	// With period + org + position filters.
	resp, err = svc.GetMovementReport(ctx, "2026-01-01", "2026-12-31", orgID.String(), posID.String(), "", "", "")
	if err != nil {
		t.Fatalf("GetMovementReport(period+org+pos) failed: %v", err)
	}
	if resp.Data.Total != 3 {
		t.Errorf("expected 3 movements with org/pos filter (offboarding excluded), got %d", resp.Data.Total)
	}

	// Invalid organization id → MovementValidationError (400).
	if _, err := svc.GetMovementReport(ctx, "", "", "not-a-uuid", "", "", "", ""); err == nil {
		t.Error("expected error for invalid organization_id")
	} else if _, ok := err.(*MovementValidationError); !ok {
		t.Errorf("expected MovementValidationError, got %T", err)
	}

	// Terbalik periode (date_from > date_to) → MovementValidationError (400).
	if _, err := svc.GetMovementReport(ctx, "2026-12-31", "2026-01-01", "", "", "", "", ""); err == nil {
		t.Error("expected error when date_from > date_to")
	} else if _, ok := err.(*MovementValidationError); !ok {
		t.Errorf("expected MovementValidationError for inverted period, got %T", err)
	}

	// Periode valid (date_from <= date_to) → tidak error.
	if _, err := svc.GetMovementReport(ctx, "2026-01-01", "2026-12-31", "", "", "", "", ""); err != nil {
		t.Errorf("expected no error for valid period, got %v", err)
	}
}

func TestService_GetContractReport(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()
	svc := NewService(repo, testLogger())

	empID := uuid.New()
	createTestContract(repo, empID) // active, end 2026-12-31 (helper default)
	// One active contract expiring soon (deterministic relative to "today" is
	// impossible in tests, so assert structural correctness instead).
	c2 := createTestContract(repo, empID)
	db, err := dbResolver(ctx)
	if err != nil {
		t.Fatalf("dbResolver failed: %v", err)
	}
	if err := db.Model(&EmployeeContract{}).Where("id = ?", c2.ID).Update("status", ContractStatusExtended).Error; err != nil {
		t.Fatalf("update extended failed: %v", err)
	}

	resp, err := svc.GetContractReport(ctx)
	if err != nil {
		t.Fatalf("GetContractReport failed: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Data.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Data.Total)
	}
	if resp.Data.ByStatus[string(ContractStatusActive)] != 1 {
		t.Errorf("expected 1 active, got %d", resp.Data.ByStatus[string(ContractStatusActive)])
	}
	if resp.Data.ByStatus[string(ContractStatusExtended)] != 1 {
		t.Errorf("expected 1 extended, got %d", resp.Data.ByStatus[string(ContractStatusExtended)])
	}
}

// =========================================================================
// Handler Report Tests (plan §12.17)
// =========================================================================

func TestHandler_GetMovementReport(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	seedReportMovements(t, repo, context.Background())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/employee-movements/reports/movements", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
	var body struct {
		Success bool `json:"success"`
		Data    struct {
			Total    int64            `json:"total"`
			ByType   map[string]int64 `json:"by_type"`
			ByStatus map[string]int64 `json:"by_status"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if !body.Success {
		t.Error("expected success=true")
	}
	if body.Data.Total != 4 {
		t.Errorf("expected total 4, got %d", body.Data.Total)
	}
}

func TestHandler_GetMovementReport_WithFilters(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	orgID, _, _ := seedReportMovements(t, repo, context.Background())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/employee-movements/reports/movements?date_from=2026-01-01&date_to=2026-12-31&organization_id="+orgID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
	var body struct {
		Data struct {
			Total int64 `json:"total"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if body.Data.Total != 3 {
		t.Errorf("expected total 3 with org filter, got %d", body.Data.Total)
	}
}

func TestHandler_GetContractReport(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	empID := uuid.New()
	createTestContract(repo, empID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/employee-movements/reports/contracts", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
	var body struct {
		Success bool `json:"success"`
		Data    struct {
			Total    int64            `json:"total"`
			ByStatus map[string]int64 `json:"by_status"`
			Expiring int64            `json:"expiring"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if !body.Success {
		t.Error("expected success=true")
	}
	if body.Data.Total != 1 {
		t.Errorf("expected total 1, got %d", body.Data.Total)
	}
	if body.Data.ByStatus["active"] != 1 {
		t.Errorf("expected 1 active contract, got %d", body.Data.ByStatus["active"])
	}
}

// =========================================================================
// HR Dashboard Tests (plan §12.18)
// =========================================================================

func TestRepo_CountMovementsEffectiveBetween(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	empID := uuid.New()
	m1 := &EmployeeMovement{
		EmployeeID:           empID,
		MovementType:         MovementTypePromotion,
		DecisionLetterNumber: "SK-DASH-1",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-05",
		Status:               MovementStatusExecuted,
	}
	if err := repo.CreateMovement(ctx, m1); err != nil {
		t.Fatalf("CreateMovement failed: %v", err)
	}
	m2 := &EmployeeMovement{
		EmployeeID:           empID,
		MovementType:         MovementTypeMutation,
		DecisionLetterNumber: "SK-DASH-2",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-20",
		Status:               MovementStatusApproved,
	}
	if err := repo.CreateMovement(ctx, m2); err != nil {
		t.Fatalf("CreateMovement failed: %v", err)
	}
	m3 := &EmployeeMovement{
		EmployeeID:           empID,
		MovementType:         MovementTypeOffboarding,
		DecisionLetterNumber: "SK-DASH-3",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-09-01",
		Status:               MovementStatusDraft,
	}
	if err := repo.CreateMovement(ctx, m3); err != nil {
		t.Fatalf("CreateMovement failed: %v", err)
	}

	// August window — 2 movements effective in range.
	count, err := repo.CountMovementsEffectiveBetween(ctx, "2026-08-01", "2026-08-31")
	if err != nil {
		t.Fatalf("CountMovementsEffectiveBetween failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 movements effective in Aug, got %d", count)
	}

	// Empty window.
	count, err = repo.CountMovementsEffectiveBetween(ctx, "2027-01-01", "2027-01-31")
	if err != nil {
		t.Fatalf("CountMovementsEffectiveBetween(empty) failed: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 movements in empty window, got %d", count)
	}
}

func TestService_GetHRDashboard(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()
	svc := NewService(repo, testLogger())

	empID := uuid.New()
	seedReportMovements(t, repo, ctx) // 2 promotion, 1 mutation, 1 offboarding

	// One movement pending approval + one executed effective this month.
	pending := &EmployeeMovement{
		EmployeeID:           empID,
		MovementType:         MovementTypeStatusChange,
		DecisionLetterNumber: "SK-DASH-P",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        time.Now().Format("2006-01-02"),
		Status:               MovementStatusPendingApproval,
	}
	if err := repo.CreateMovement(ctx, pending); err != nil {
		t.Fatalf("CreateMovement failed: %v", err)
	}
	executed := &EmployeeMovement{
		EmployeeID:           empID,
		MovementType:         MovementTypeRetirement,
		DecisionLetterNumber: "SK-DASH-E",
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        time.Now().Format("2006-01-02"),
		Status:               MovementStatusExecuted,
	}
	if err := repo.CreateMovement(ctx, executed); err != nil {
		t.Fatalf("CreateMovement failed: %v", err)
	}
	// Contract untuk ringkasan.
	createTestContract(repo, empID) // active

	resp, err := svc.GetHRDashboard(ctx)
	if err != nil {
		t.Fatalf("GetHRDashboard failed: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Data.MovementByType["promotion"] != 2 {
		t.Errorf("expected 2 promotions, got %d", resp.Data.MovementByType["promotion"])
	}
	if resp.Data.PendingApproval != 1 {
		t.Errorf("expected 1 pending approval, got %d", resp.Data.PendingApproval)
	}
	if resp.Data.EffectiveThisMonth < 2 {
		t.Errorf("expected at least 2 movements effective this month, got %d", resp.Data.EffectiveThisMonth)
	}
	if resp.Data.Contracts.Active != 1 {
		t.Errorf("expected 1 active contract, got %d", resp.Data.Contracts.Active)
	}
}

func TestHandler_GetHRDashboard(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	seedReportMovements(t, repo, context.Background())
	createTestContract(repo, uuid.New())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/employee-movements/dashboard", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
	var body struct {
		Success bool `json:"success"`
		Data    struct {
			MovementByType     map[string]int64 `json:"movement_by_type"`
			PendingApproval    int64            `json:"pending_approval"`
			EffectiveThisMonth int64            `json:"effective_this_month"`
			Contracts          struct {
				Active   int64 `json:"active"`
				Expiring int64 `json:"expiring"`
				Expired  int64 `json:"expired"`
			} `json:"contracts"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if !body.Success {
		t.Error("expected success=true")
	}
	if body.Data.MovementByType["mutation"] != 1 {
		t.Errorf("expected 1 mutation, got %d", body.Data.MovementByType["mutation"])
	}
	if body.Data.Contracts.Active != 1 {
		t.Errorf("expected 1 active contract, got %d", body.Data.Contracts.Active)
	}
}
