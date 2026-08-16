package reimbursement

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func setupTestRouter() (*gin.Engine, *Repository, func()) {
	gin.SetMode(gin.TestMode)

	db, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	handler := NewHandler(svc)

	// Simulate auth middleware: sets user_id (user-account UUID) in request
	// context, with a linked employee_accounts row so request creation
	// resolves to the employee the same way production does.
	userID := uuid.New()
	empID := uuid.New()
	createTestEmployeeAccount(db, empID, userID)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "user_id", userID.String())
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})
	rg := r.Group("/api/v1/tenant")
	RegisterRoutes(rg, handler)

	return r, repo, func() {
		cleanup()
		_ = logger.Sync()
	}
}

// =========================================================================
// Reimbursement Type Handler Tests
// =========================================================================

func TestHandler_CreateReimbursementType_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"name": "Medical Reimbursement", "description": "Biaya berobat", "code": "MED"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/reimbursements/types", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool                     `json:"success"`
		Data    ReimbursementTypeResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Data.Name != "Medical Reimbursement" {
		t.Errorf("expected name 'Medical Reimbursement', got '%s'", resp.Data.Name)
	}
}

func TestHandler_ListReimbursementTypes_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	createTestReimbursementType(repo)
	createTestReimbursementType(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/reimbursements/types", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetReimbursementTypeByID_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestReimbursementType(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/reimbursements/types/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetReimbursementTypeByID_NotFound(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/reimbursements/types/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdateReimbursementType_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestReimbursementType(repo)

	body := `{"name": "Updated Type"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/tenant/reimbursements/types/"+created.ID.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DeleteReimbursementType_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestReimbursementType(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/tenant/reimbursements/types/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Reimbursement Request Handler Tests
// =========================================================================

func TestHandler_CreateReimbursementRequest_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	rType := createTestReimbursementType(repo)

	body := `{
		"request_type_id": "` + rType.ID.String() + `",
		"title": "Medical Claim July 2026",
		"description": "Biaya berobat di klinik"
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/reimbursements/requests", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateReimbursementRequest_ValidationError(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"title": ""}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/reimbursements/requests", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListReimbursementRequests_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	createTestReimbursementRequest(repo, empID, rType.ID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/reimbursements/requests", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetReimbursementRequestByID_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	created := createTestReimbursementRequest(repo, empID, rType.ID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/reimbursements/requests/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetReimbursementRequestByID_NotFound(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/reimbursements/requests/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdateReimbursementRequestStatus_Submit(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	rr := createTestReimbursementRequest(repo, empID, rType.ID)
	createTestReimbursementItem(repo, rr.ID)

	body := `{"status": "SUBMITTED"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/tenant/reimbursements/requests/"+rr.ID.String()+"/status", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DeleteReimbursementRequest_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	rr := createTestReimbursementRequest(repo, empID, rType.ID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/tenant/reimbursements/requests/"+rr.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Reimbursement Item Handler Tests
// =========================================================================

func TestHandler_CreateReimbursementItem_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	rr := createTestReimbursementRequest(repo, empID, rType.ID)

	body := `{
		"expense_date": "2026-07-15",
		"expense_type": "MEDICAL",
		"description": "Doctor visit",
		"amount": 250000
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/reimbursements/requests/"+rr.ID.String()+"/items", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListReimbursementItems_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	rr := createTestReimbursementRequest(repo, empID, rType.ID)

	for i := 0; i < 2; i++ {
		createTestReimbursementItem(repo, rr.ID)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/reimbursements/requests/"+rr.ID.String()+"/items", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdateReimbursementItem_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	rr := createTestReimbursementRequest(repo, empID, rType.ID)
	item := createTestReimbursementItem(repo, rr.ID)

	body := `{"amount": 500000}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/tenant/reimbursements/requests/"+rr.ID.String()+"/items/"+item.ID.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DeleteReimbursementItem_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	rType := createTestReimbursementType(repo)
	empID := uuid.New()
	rr := createTestReimbursementRequest(repo, empID, rType.ID)
	item := createTestReimbursementItem(repo, rr.ID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/tenant/reimbursements/requests/"+rr.ID.String()+"/items/"+item.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}
