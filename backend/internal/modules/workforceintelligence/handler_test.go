package workforceintelligence

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

	_, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	handler := NewHandler(svc)

	r := gin.New()
	// Simulate auth middleware that sets user_id in request context
	r.Use(func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "user_id", uuid.New().String())
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
// Headcount Plan Handler Tests
// =========================================================================

func TestHandler_CreateHeadcountPlan_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"period": "2026-Q3", "organization_id": "` + uuid.New().String() + `", "planned_hc": 100}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/workforce-intelligence/planning/headcounts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool                    `json:"success"`
		Data    HeadcountPlanResponse   `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Data.PlannedHC != 100 {
		t.Errorf("expected PlannedHC 100, got %d", resp.Data.PlannedHC)
	}
}

func TestHandler_CreateHeadcountPlan_ValidationError(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"period": ""}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/workforce-intelligence/planning/headcounts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListHeadcountPlans_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	createTestHeadcountPlan(repo)
	createTestHeadcountPlan(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/planning/headcounts", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetHeadcountPlanByID_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestHeadcountPlan(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/planning/headcounts/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetHeadcountPlanByID_NotFound(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/planning/headcounts/"+uuid.New().String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdateHeadcountPlan_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestHeadcountPlan(repo)
	body := `{"planned_hc": 200}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/tenant/workforce-intelligence/planning/headcounts/"+created.ID.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DeleteHeadcountPlan_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestHeadcountPlan(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/tenant/workforce-intelligence/planning/headcounts/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Forecast Handler Tests
// =========================================================================

func TestHandler_CreateForecast_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{
		"period": "2026-Q3",
		"organization_id": "` + uuid.New().String() + `",
		"forecast_type": "DEMAND",
		"headcount": 150
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/workforce-intelligence/planning/forecasts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListForecasts_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	createTestForecast(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/planning/forecasts", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// KPI Handler Tests
// =========================================================================

func TestHandler_ListKPIs_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	createTestKPI(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/kpi", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetKPISummary_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/kpi/summary", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetKPIByCode_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestKPI(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/kpi/"+created.KpiCode, nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetKPIByCode_NotFound(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/kpi/NONEXISTENT", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Analytics Handler Tests
// =========================================================================

func TestHandler_GetHeadcountAnalytics_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/analytics/headcount", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Capacity Handler Tests
// =========================================================================

func TestHandler_GetCapacityForecast_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/capacity/forecast", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Cost Handler Tests
// =========================================================================

func TestHandler_GetCostPerEmployee_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/cost/per-employee", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Scenario Handler Tests
// =========================================================================

func TestHandler_CreateScenario_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{
		"name": "Growth Scenario",
		"scenario_type": "GROWTH",
		"parameters": {"growth_rate": 10}
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/workforce-intelligence/scenarios", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListScenarios_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	createTestScenario(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/scenarios", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetScenarioByID_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestScenario(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/scenarios/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdateScenario_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestScenario(repo)
	body := `{"name": "Updated Scenario"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/tenant/workforce-intelligence/scenarios/"+created.ID.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DeleteScenario_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestScenario(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/tenant/workforce-intelligence/scenarios/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_RunScenario_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestScenario(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/workforce-intelligence/scenarios/"+created.ID.String()+"/run", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CloneScenario_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestScenario(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/workforce-intelligence/scenarios/"+created.ID.String()+"/clone", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Risk Handler Tests
// =========================================================================

func TestHandler_ListRiskIndicators_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	createTestRiskIndicator(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/risk/indicators", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetRiskHighTurnover_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/risk/high-turnover", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetRiskRetirement_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/risk/retirement", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Health Score Handler Tests
// =========================================================================

func TestHandler_ListHealthScores_Success(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	createTestHealthScore(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/health/scores", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetHealthDashboard_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/health/dashboard", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetSpanOfControl_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/health/span-of-control", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetSuccessionReadiness_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/health/succession", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// People Analytics Handler Tests
// =========================================================================

func TestHandler_GetTrainingVsPerformance_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/people-analytics/training-vs-performance", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetOvertimeVsProductivity_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/people-analytics/overtime-vs-productivity", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Executive Dashboard Handler Tests
// =========================================================================

func TestHandler_GetExecutiveSummary_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/executive/summary", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetExecutiveGrowth_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/executive/growth", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetExecutiveHealthScore_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/executive/health-score", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Gap Analysis Handler Tests
// =========================================================================

func TestHandler_GetGapAnalysis_NoEmployeesTable(t *testing.T) {
	router, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/workforce-intelligence/planning/gap-analysis", nil)
	router.ServeHTTP(w, req)

	// Expected 500: GetGapAnalysis calls GetActiveEmployeeCount which queries
	// the employees table (source module) that doesn't exist in test DB.
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 Internal Server Error (missing employees table), got %d: %s", w.Code, w.Body.String())
	}
}
