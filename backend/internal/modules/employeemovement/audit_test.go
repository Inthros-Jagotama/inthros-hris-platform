package employeemovement

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// createAuditedMovement creates a draft movement THROUGH the service (so a
// CREATED audit row is recorded, plan §12.6) and returns it.
func createAuditedMovement(t *testing.T, svc *Service, employeeID uuid.UUID) *EmployeeMovement {
	t.Helper()
	resp, err := svc.CreateMovement(ctx(), CreateMovementRequest{
		EmployeeID:           employeeID.String(),
		MovementType:         string(MovementTypeOther),
		DecisionLetterNumber: "SK-AUDIT-" + uuid.New().String()[:8],
		DecisionLetterDate:   "2026-07-01",
		EffectiveDate:        "2026-08-01",
	})
	if err != nil {
		t.Fatalf("CreateMovement failed: %v", err)
	}
	uid, _ := uuid.Parse(resp.ID)
	m, err := svc.repo.FindMovementByID(ctx(), uid)
	if err != nil {
		t.Fatalf("failed to reload created movement: %v", err)
	}
	return m
}

// TestService_CreateMovement_RecordsAudit verifies plan §12.6: creating a
// movement writes a CREATED audit row with new_status draft + snapshot.
func TestService_CreateMovement_RecordsAudit(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createAuditedMovement(t, svc, employeeID)

	resp, err := svc.ListMovementAudits(ctx(), created.ID.String(), 1, 100)
	if err != nil {
		t.Fatalf("ListMovementAudits failed: %v", err)
	}
	if resp.Total != 1 {
		t.Fatalf("expected 1 audit row, got %d", resp.Total)
	}
	items := resp.Data.([]MovementAuditResponse)
	entry := items[0]
	if entry.Action != string(MovementAuditActionCreated) {
		t.Errorf("expected action CREATED, got '%s'", entry.Action)
	}
	if entry.NewStatus == nil || *entry.NewStatus != string(MovementStatusDraft) {
		t.Errorf("expected new_status draft, got %v", entry.NewStatus)
	}
	if entry.NewData == nil || !strings.Contains(*entry.NewData, created.ID.String()) {
		t.Errorf("expected new_data snapshot to contain movement id, got %v", entry.NewData)
	}
}

// TestService_UpdateMovement_RecordsAudit verifies UPDATED audit carries
// before/after snapshots.
func TestService_UpdateMovement_RecordsAudit(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createAuditedMovement(t, svc, employeeID)

	newReason := "Alasan diperbarui"
	_, err := svc.UpdateMovement(ctx(), created.ID.String(), UpdateMovementRequest{
		Reason: &newReason,
	})
	if err != nil {
		t.Fatalf("UpdateMovement failed: %v", err)
	}

	resp, err := svc.ListMovementAudits(ctx(), created.ID.String(), 1, 100)
	if err != nil {
		t.Fatalf("ListMovementAudits failed: %v", err)
	}
	if resp.Total != 2 {
		t.Fatalf("expected 2 audit rows (CREATED + UPDATED), got %d", resp.Total)
	}
	items := resp.Data.([]MovementAuditResponse)
	updated := items[0]
	if updated.Action != string(MovementAuditActionUpdated) {
		t.Errorf("expected newest action UPDATED, got '%s'", updated.Action)
	}
	if updated.OldData == nil || updated.NewData == nil {
		t.Fatal("expected old_data and new_data snapshots for UPDATED")
	}
	if !strings.Contains(*updated.NewData, newReason) {
		t.Errorf("expected new_data to contain updated reason, got %s", *updated.NewData)
	}
	if strings.Contains(*updated.OldData, newReason) {
		t.Errorf("expected old_data to NOT contain updated reason, got %s", *updated.OldData)
	}
}

// TestService_SubmitMovement_RecordsAudit verifies SUBMITTED audit is recorded
// with draft → pending_approval transition.
func TestService_SubmitMovement_RecordsAudit(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	svc.SetApprovalEngine(&fakeApprovalEngine{resolvedFlowID: uuid.New().String()})

	employeeID := uuid.New()
	created := createAuditedMovement(t, svc, employeeID)

	if _, err := svc.SubmitMovement(ctx(), created.ID.String(), SubmitMovementRequest{}); err != nil {
		t.Fatalf("SubmitMovement failed: %v", err)
	}

	resp, err := svc.ListMovementAudits(ctx(), created.ID.String(), 1, 100)
	if err != nil {
		t.Fatalf("ListMovementAudits failed: %v", err)
	}
	if resp.Total != 2 {
		t.Fatalf("expected 2 audit rows, got %d", resp.Total)
	}
	items := resp.Data.([]MovementAuditResponse)
	submitted := items[0]
	if submitted.Action != string(MovementAuditActionSubmitted) {
		t.Errorf("expected action SUBMITTED, got '%s'", submitted.Action)
	}
	if submitted.OldStatus == nil || *submitted.OldStatus != string(MovementStatusDraft) {
		t.Errorf("expected old_status draft, got %v", submitted.OldStatus)
	}
	if submitted.NewStatus == nil || *submitted.NewStatus != string(MovementStatusPendingApproval) {
		t.Errorf("expected new_status pending_approval, got %v", submitted.NewStatus)
	}
}

// TestService_HandleApprovalStatusChange_RecordsAudit verifies APPROVED /
// REJECTED audits are recorded with the right transitions.
func TestService_HandleApprovalStatusChange_RecordsAudit(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createAuditedMovement(t, svc, employeeID)
	created.Status = MovementStatusPendingApproval
	if err := svc.repo.UpdateMovement(ctx(), created); err != nil {
		t.Fatalf("failed to set pending_approval: %v", err)
	}

	// APPROVED
	if err := svc.HandleApprovalStatusChange(ctx(), created.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange APPROVED failed: %v", err)
	}
	// REJECTED (from a fresh pending_approval movement)
	created2 := createAuditedMovement(t, svc, employeeID)
	created2.Status = MovementStatusPendingApproval
	if err := svc.repo.UpdateMovement(ctx(), created2); err != nil {
		t.Fatalf("failed to set pending_approval: %v", err)
	}
	rejectNote := "SK kurang lengkap"
	if err := svc.HandleApprovalStatusChange(ctx(), created2.ID, "REJECTED", rejectNote); err != nil {
		t.Fatalf("HandleApprovalStatusChange REJECTED failed: %v", err)
	}

	resp1, _ := svc.ListMovementAudits(ctx(), created.ID.String(), 1, 100)
	items1 := resp1.Data.([]MovementAuditResponse)
	if resp1.Total != 2 {
		t.Fatalf("expected 2 audit rows for movement 1 (CREATED + APPROVED), got %d", resp1.Total)
	}
	approved := items1[0]
	if approved.Action != string(MovementAuditActionApproved) {
		t.Errorf("expected action APPROVED, got '%s'", approved.Action)
	}
	if approved.OldStatus == nil || *approved.OldStatus != string(MovementStatusPendingApproval) {
		t.Errorf("expected old_status pending_approval, got %v", approved.OldStatus)
	}
	if approved.NewStatus == nil || *approved.NewStatus != string(MovementStatusApproved) {
		t.Errorf("expected new_status approved, got %v", approved.NewStatus)
	}

	resp2, _ := svc.ListMovementAudits(ctx(), created2.ID.String(), 1, 100)
	items2 := resp2.Data.([]MovementAuditResponse)
	rejected := items2[0]
	if rejected.Action != string(MovementAuditActionRejected) {
		t.Errorf("expected action REJECTED, got '%s'", rejected.Action)
	}
	if rejected.Reason == nil || *rejected.Reason != rejectNote {
		t.Errorf("expected reason '%s', got %v", rejectNote, rejected.Reason)
	}
	if rejected.NewStatus == nil || *rejected.NewStatus != string(MovementStatusRejected) {
		t.Errorf("expected new_status rejected, got %v", rejected.NewStatus)
	}
}

// TestService_CancelMovement_RecordsAudit verifies CANCELLED audit is recorded.
func TestService_CancelMovement_RecordsAudit(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createAuditedMovement(t, svc, employeeID)

	if err := svc.CancelMovement(ctx(), created.ID.String()); err != nil {
		t.Fatalf("CancelMovement failed: %v", err)
	}

	resp, err := svc.ListMovementAudits(ctx(), created.ID.String(), 1, 100)
	if err != nil {
		t.Fatalf("ListMovementAudits failed: %v", err)
	}
	if resp.Total != 2 {
		t.Fatalf("expected 2 audit rows, got %d", resp.Total)
	}
	items := resp.Data.([]MovementAuditResponse)
	cancelled := items[0]
	if cancelled.Action != string(MovementAuditActionCancelled) {
		t.Errorf("expected action CANCELLED, got '%s'", cancelled.Action)
	}
	if cancelled.NewStatus == nil || *cancelled.NewStatus != string(MovementStatusCancelled) {
		t.Errorf("expected new_status cancelled, got %v", cancelled.NewStatus)
	}
}

// TestService_ExecuteMovement_RecordsAudit verifies EXECUTED audit is recorded
// (non-tx path — contract_extension) with approved → executed transition.
func TestService_ExecuteMovement_RecordsAudit(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createAuditedMovement(t, svc, employeeID)
	created.MovementType = MovementTypeContractExtension
	created.Status = MovementStatusApproved
	if err := svc.repo.UpdateMovement(ctx(), created); err != nil {
		t.Fatalf("failed to set approved: %v", err)
	}

	executorID := uuid.New().String()
	if err := svc.ExecuteMovement(ctx(), created.ID.String(), executorID); err != nil {
		t.Fatalf("ExecuteMovement failed: %v", err)
	}

	resp, err := svc.ListMovementAudits(ctx(), created.ID.String(), 1, 100)
	if err != nil {
		t.Fatalf("ListMovementAudits failed: %v", err)
	}
	if resp.Total != 2 {
		t.Fatalf("expected 2 audit rows, got %d", resp.Total)
	}
	items := resp.Data.([]MovementAuditResponse)
	executed := items[0]
	if executed.Action != string(MovementAuditActionExecuted) {
		t.Errorf("expected action EXECUTED, got '%s'", executed.Action)
	}
	if executed.OldStatus == nil || *executed.OldStatus != string(MovementStatusApproved) {
		t.Errorf("expected old_status approved, got %v", executed.OldStatus)
	}
	if executed.NewStatus == nil || *executed.NewStatus != string(MovementStatusExecuted) {
		t.Errorf("expected new_status executed, got %v", executed.NewStatus)
	}
	if executed.ActedBy == nil || *executed.ActedBy != executorID {
		t.Errorf("expected acted_by executor %s, got %v", executorID, executed.ActedBy)
	}
}

// TestService_ListMovementAudits_Pagination verifies pagination works for the
// audit trail (list newest-first).
func TestService_ListMovementAudits_Pagination(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createAuditedMovement(t, svc, employeeID)

	svc.SetApprovalEngine(&fakeApprovalEngine{resolvedFlowID: uuid.New().String()})
	if _, err := svc.SubmitMovement(ctx(), created.ID.String(), SubmitMovementRequest{}); err != nil {
		t.Fatalf("SubmitMovement failed: %v", err)
	}

	page1, err := svc.ListMovementAudits(ctx(), created.ID.String(), 1, 1)
	if err != nil {
		t.Fatalf("ListMovementAudits page 1 failed: %v", err)
	}
	if page1.Total != 2 {
		t.Fatalf("expected total 2, got %d", page1.Total)
	}
	if page1.TotalPages != 2 {
		t.Fatalf("expected total_pages 2, got %d", page1.TotalPages)
	}
	items := page1.Data.([]MovementAuditResponse)
	if len(items) != 1 || items[0].Action != string(MovementAuditActionSubmitted) {
		t.Errorf("expected page 1 to contain SUBMITTED (newest first), got %+v", items)
	}
}

// TestHandler_ListMovementAudits verifies GET /movements/:id/audits endpoint.
func TestHandler_ListMovementAudits(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	// Create through the router so the service instance records the audit.
	employeeID := uuid.New()
	body := `{
		"employee_id": "` + employeeID.String() + `",
		"movement_type": "other",
		"decision_letter_number": "SK-AUDIT-H",
		"decision_letter_date": "2026-07-01",
		"effective_date": "2026-08-01"
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/employee-movements/movements", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 created, got %d: %s", w.Code, w.Body.String())
	}
	var createdResp struct {
		Data MovementResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &createdResp); err != nil {
		t.Fatalf("failed to unmarshal create response: %v", err)
	}

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/tenant/employee-movements/movements/"+createdResp.Data.ID+"/audits", nil)
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w2.Code, w2.Body.String())
	}

	var resp struct {
		Success bool `json:"success"`
		Data    []struct {
			ID        string  `json:"id"`
			Action    string  `json:"action"`
			NewStatus *string `json:"new_status"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w2.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 audit entry, got %d", len(resp.Data))
	}
	if resp.Data[0].Action != string(MovementAuditActionCreated) {
		t.Errorf("expected action CREATED, got '%s'", resp.Data[0].Action)
	}
}
