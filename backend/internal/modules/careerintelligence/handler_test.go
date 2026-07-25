package careerintelligence

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
	rg := r.Group("/api/v1/tenant")
	RegisterRoutes(rg, handler)

	return r, repo, func() {
		cleanup()
		_ = logger.Sync()
	}
}

// =========================================================================
// Talent Map Handler Tests
// =========================================================================

func TestHandler_CreateTalentMap_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"employee_id": "` + uuid.New().String() + `", "period": "2026-07", "performance": "HIGH", "potential": "HIGH"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/career-intelligence/talent-maps", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool              `json:"success"`
		Data    TalentMapResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Data.GridPosition != "9-BOX-1" {
		t.Errorf("expected grid position '9-BOX-1', got '%s'", resp.Data.GridPosition)
	}
}

func TestHandler_CreateTalentMap_ValidationError(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"employee_id": ""}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/career-intelligence/talent-maps", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListTalentMaps_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	for i := 0; i < 2; i++ {
		createTestTalentMap(repo)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/career-intelligence/talent-maps", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetTalentMapByID_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestTalentMap(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/career-intelligence/talent-maps/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetTalentMapByID_NotFound(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/career-intelligence/talent-maps/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetTalentGrid_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	for i := 0; i < 3; i++ {
		createTestTalentMap(repo)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/career-intelligence/talent-maps/grid", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetEmployeeTalentProfile_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestTalentMap(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/career-intelligence/talent-maps/employee/"+created.EmployeeID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdateTalentMap_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestTalentMap(repo)
	body := `{"performance": "LOW"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/tenant/career-intelligence/talent-maps/"+created.ID.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DeleteTalentMap_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestTalentMap(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/tenant/career-intelligence/talent-maps/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Career Interest Handler Tests
// =========================================================================

func TestHandler_CreateCareerInterest_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"employee_id": "` + uuid.New().String() + `", "interest_type": "LEADERSHIP", "target_position": "VP"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/career-intelligence/interests", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListCareerInterests_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	for i := 0; i < 2; i++ {
		createTestCareerInterest(repo)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/career-intelligence/interests", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetEmployeeCareerInterests_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestCareerInterest(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/career-intelligence/interests/employee/"+created.EmployeeID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Career Path Handler Tests
// =========================================================================

func TestHandler_CreateCareerPath_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"source_title_id": "` + uuid.New().String() + `", "target_title_id": "` + uuid.New().String() + `", "path_type": "PROMOTION"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/career-intelligence/paths", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListCareerPaths_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	for i := 0; i < 2; i++ {
		createTestCareerPath(repo)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/career-intelligence/paths", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetGapAnalysis_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/career-intelligence/paths/gap-analysis?employee_id="+uuid.New().String()+"&target_title_id="+uuid.New().String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetGapAnalysis_MissingParams(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/career-intelligence/paths/gap-analysis", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DeleteCareerPath_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestCareerPath(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/tenant/career-intelligence/paths/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Succession Plan Handler Tests
// =========================================================================

func TestHandler_CreateSuccessionPlan_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"position_id": "` + uuid.New().String() + `", "successor_id": "` + uuid.New().String() + `", "readiness_level": "READY_NOW"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/career-intelligence/successions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListSuccessionPlans_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	for i := 0; i < 2; i++ {
		createTestSuccessionPlan(repo)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/career-intelligence/successions", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetSuccessionPlanByID_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestSuccessionPlan(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/career-intelligence/successions/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdateSuccessionPlan_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestSuccessionPlan(repo)
	body := `{"readiness_level": "READY_2YR"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/tenant/career-intelligence/successions/"+created.ID.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DeleteSuccessionPlan_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestSuccessionPlan(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/tenant/career-intelligence/successions/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}
