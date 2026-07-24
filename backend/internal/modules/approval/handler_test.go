package approval

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// setupTestRouter creates a Gin engine with the approval routes registered.
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
// Flow Handler Tests
// =========================================================================

func TestHandler_CreateFlow_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"module": "leave", "name": "Leave Approval Flow", "version": 1}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/approval/flows", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool         `json:"success"`
		Data    FlowResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Data.Module != "leave" {
		t.Errorf("expected module 'leave', got '%s'", resp.Data.Module)
	}
}

func TestHandler_CreateFlow_ValidationError(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	// Missing required module and name
	body := `{"version": 1}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/approval/flows", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetFlowByID_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestFlow(repo, "leave")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/approval/flows/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool         `json:"success"`
		Data    FlowResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Data.ID != created.ID.String() {
		t.Errorf("expected ID '%s', got '%s'", created.ID.String(), resp.Data.ID)
	}
}

func TestHandler_GetFlowByID_NotFound(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/approval/flows/"+uuidStr(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdateFlow_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestFlow(repo, "leave")

	body := `{"name": "Updated Flow Name"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/tenant/approval/flows/"+created.ID.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DeleteFlow_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestFlow(repo, "leave")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/tenant/approval/flows/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListFlows(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	createTestFlow(repo, "leave")
	createTestFlow(repo, "leave")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/approval/flows", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Step Handler Tests
// =========================================================================

func TestHandler_CreateStep_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	approverID := uuidStr()

	body := `{
		"step_name": "Manager Approval",
		"approver_type": "USER",
		"approver_user_id": "` + approverID + `"
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/approval/flows/"+flow.ID.String()+"/steps", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListSteps_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)
	createTestStep(repo, flow.ID, 2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/approval/flows/"+flow.ID.String()+"/steps", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Instance Handler Tests
// =========================================================================

func TestHandler_CreateInstance_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)

	body := `{
		"module": "leave",
		"document_id": "` + uuidStr() + `",
		"flow_id": "` + flow.ID.String() + `"
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/approval/instances", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool             `json:"success"`
		Data    InstanceResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Data.Module != "leave" {
		t.Errorf("expected module 'leave', got '%s'", resp.Data.Module)
	}
}

func TestHandler_CreateInstance_ValidationError(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	// Missing required fields
	body := `{"module": "leave"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/approval/instances", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetInstanceByID_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	createTestStep(repo, flow.ID, 1)
	inst := createTestInstance(repo, flow, uuid.New())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/approval/instances/"+inst.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListInstances(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	flow := createTestFlow(repo, "leave")
	createTestInstance(repo, flow, uuid.New())
	createTestInstance(repo, flow, uuid.New())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/approval/instances", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Task Handler Tests
// =========================================================================

func TestHandler_ListMyPendingTasks_Unauthenticated(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	// Without user_id in context (no auth middleware) → expected error
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/approval/tasks/pending", nil)
	router.ServeHTTP(w, req)

	// Should return 500 because user_id is not in context
	// (No auth middleware so c.Get("user_id") returns nil)
	if w.Code != http.StatusInternalServerError && w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 or 500 without auth, got %d: %s", w.Code, w.Body.String())
	}
}
