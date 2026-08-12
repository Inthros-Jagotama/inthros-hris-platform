package employeemovement

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/modules/approval"
)

// setupTestRouter creates a Gin engine with the employee movement routes registered.
func setupTestRouter() (*gin.Engine, *Repository, func()) {
	gin.SetMode(gin.TestMode)

	_, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	handler := NewHandler(svc)

	r := gin.New()
	rg := r.Group("/api/v1/tenant")
	RegisterRoutes(rg, handler)

	return r, repo, func() {
		cleanup()
		_ = logger.Sync()
	}
}

// =========================================================================
// Movement Handler Tests
// =========================================================================

// TestHandler_SubmitMovement_ApprovalRoutingErrorBilingual verifies that an
// approval routing failure (e.g. "no supervisor found ... is vacant") is
// emitted as a bilingual 400 APPROVAL_ROUTING_FAILED following the request
// language instead of a raw internal error.
func TestHandler_SubmitMovement_ApprovalRoutingErrorBilingual(t *testing.T) {
	gin.SetMode(gin.TestMode)

	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()

	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	svc := NewService(repo, logger)

	svc.SetApprovalEngine(&fakeApprovalEngine{
		resolvedFlowID: uuid.New().String(),
		createErr: &approval.RoutingError{
			Key:    "approval.no_supervisor_vacant",
			Params: []string{"Persetujuan Supervisor"},
		},
	})
	handler := NewHandler(svc)

	r := gin.New()
	r.POST("/movements/:id/submit", handler.SubmitMovement)

	empID := uuid.New()
	movement := createTestMovement(repo, empID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/movements/"+movement.ID.String()+"/submit", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "id")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	errObj, ok := body["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing error object: %+v", body)
	}
	if errObj["code"] != "APPROVAL_ROUTING_FAILED" {
		t.Errorf("expected code APPROVAL_ROUTING_FAILED, got %v", errObj["code"])
	}
	msg, _ := errObj["message"].(string)
	if !strings.Contains(msg, "Supervisor tidak ditemukan") {
		t.Errorf("expected Indonesian supervisor message, got: %s", msg)
	}
}

func TestHandler_CreateMovement_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{
		"employee_id": "` + uuidStr() + `",
		"movement_type": "promotion",
		"to_position_id": "` + uuidStr() + `",
		"decision_letter_number": "SK-001",
		"decision_letter_date": "2026-07-01",
		"effective_date": "2026-08-01",
		"reason": "Kinerja baik"
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/employee-movements/movements", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool             `json:"success"`
		Data    MovementResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Data.MovementType != "promotion" {
		t.Errorf("expected movement_type 'promotion', got '%s'", resp.Data.MovementType)
	}
}

func TestHandler_CreateMovement_ValidationError(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	// Missing required fields
	body := `{"employee_id": "invalid"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/employee-movements/movements", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetMovementByID_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestMovement(repo, uuid.New())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/employee-movements/movements/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetMovementByID_NotFound(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/employee-movements/movements/"+uuidStr(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListMovements(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	createTestMovement(repo, uuid.New())
	createTestMovement(repo, uuid.New())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/employee-movements/movements", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListMovementsByEmployee(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	empID := uuid.New()
	createTestMovement(repo, empID)
	createTestMovement(repo, empID)
	createTestMovement(repo, uuid.New()) // another employee

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/employee-movements/employees/"+empID.String()+"/movements", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdateMovement_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestMovement(repo, uuid.New())

	body := `{"reason": "Updated reason"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/tenant/employee-movements/movements/"+created.ID.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DeleteMovement_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestMovement(repo, uuid.New())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/tenant/employee-movements/movements/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Contract Handler Tests
// =========================================================================

func TestHandler_CreateContract_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{
		"employee_id": "` + uuidStr() + `",
		"contract_number": "CTR-001",
		"contract_type": "pkwt",
		"start_date": "2026-01-01",
		"end_date": "2026-12-31"
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/employee-movements/contracts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateContract_ValidationError(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	// Missing required fields
	body := `{"employee_id": "invalid"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/employee-movements/contracts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetContractByID_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestContract(repo, uuid.New())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/employee-movements/contracts/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetContractByID_NotFound(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/employee-movements/contracts/"+uuidStr(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListContracts(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	createTestContract(repo, uuid.New())
	createTestContract(repo, uuid.New())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/employee-movements/contracts", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdateContract_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestContract(repo, uuid.New())

	body := `{"status": "expired"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/tenant/employee-movements/contracts/"+created.ID.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DeleteContract_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestContract(repo, uuid.New())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/tenant/employee-movements/contracts/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}
