package leave

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

func setupTestRouter() (*gin.Engine, *Repository, func()) {
	gin.SetMode(gin.TestMode)

	_, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	handler := NewHandler(svc)

	r := gin.New()
	// requireLeaveSettings (routes.go) checks c.GetStringSlice("permissions")
	// — simulate a fully-privileged caller so existing handler-level tests
	// (which don't exercise the RBAC layer) keep passing; permission-scoped
	// behavior itself is covered by TestRequireLeaveSettings.
	r.Use(func(c *gin.Context) {
		c.Set("permissions", []string{"*"})
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
// Leave Type Handler Tests
// =========================================================================

func TestHandler_CreateLeaveType_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"name": "Annual Leave", "description": "Annual paid leave"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/leave/types", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool           `json:"success"`
		Data    LeaveTypeResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Data.Name != "Annual Leave" {
		t.Errorf("expected name 'Annual Leave', got '%s'", resp.Data.Name)
	}
}

func TestHandler_ListLeaveTypes_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	createTestLeaveType(repo)
	createTestLeaveType(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/leave/types", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetLeaveTypeByID_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestLeaveType(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/leave/types/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Leave Reason Handler Tests
// =========================================================================

func TestHandler_CreateLeaveReason_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"name": "Sick", "sort_order": 1}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/leave/reasons", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListLeaveReasons_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	createTestLeaveReason(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/leave/reasons", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Leave Request Handler Tests
// =========================================================================

func TestHandler_CreateLeaveRequest_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	lType := createTestLeaveType(repo)

	body := `{
		"employee_id": "` + uuid.New().String() + `",
		"leave_type_id": "` + lType.ID.String() + `",
		"request_start_date": "2026-01-15",
		"request_end_date": "2026-01-16",
		"requested_days": 2
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/leave/requests", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateLeaveRequest_ValidationError(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"employee_id": "invalid"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/leave/requests", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListLeaveRequests_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	lType := createTestLeaveType(repo)
	empID := uuid.New()
	createTestLeaveRequest(repo, empID, lType.ID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/leave/requests", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Accrual Policy Handler Tests
// =========================================================================

func TestHandler_CreateAccrualPolicy_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	lType := createTestLeaveType(repo)

	body := `{
		"leave_type_id": "` + lType.ID.String() + `",
		"base_quota_days": 12,
		"effective_from": "2026-01-01"
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/leave/accrual-policies", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListAccrualPolicies_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	lType := createTestLeaveType(repo)
	createTestAccrualPolicy(repo, lType.ID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/leave/accrual-policies", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Leave Balance Handler Tests
// =========================================================================

func TestHandler_ListLeaveBalances_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/leave/balances", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}
