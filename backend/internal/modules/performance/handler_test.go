package performance

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func setupTestRouter() (*gin.Engine, *Handler, func()) {
	gin.SetMode(gin.TestMode)
	_, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger := zap.NewNop()
	svc := NewService(repo, logger)
	handler := NewHandler(svc)

	r := gin.New()
	tenant := r.Group("/api/v1/tenant")
	RegisterRoutes(tenant, handler)

	return r, handler, func() { cleanup() }
}

func performRequest(r http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// =========================================================================
// Performance Period Handler Tests
// =========================================================================

func TestHandler_CreatePerformancePeriod(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	body := map[string]interface{}{
		"period_code": "Q1",
		"period_type": "QUARTERLY",
		"year":        2026,
	}

	w := performRequest(r, "POST", "/api/v1/tenant/performance/periods", body)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Errorf("expected success true, got %v", resp["success"])
	}
}

func TestHandler_CreatePerformancePeriod_ValidationError(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	body := map[string]interface{}{
		"period_type": "INVALID",
	}

	w := performRequest(r, "POST", "/api/v1/tenant/performance/periods", body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandler_ListPerformancePeriods(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	// Create some periods first
	for i := 0; i < 3; i++ {
		performRequest(r, "POST", "/api/v1/tenant/performance/periods", map[string]interface{}{
			"period_code": "P",
			"period_type": "MONTHLY",
			"year":        2026,
		})
	}

	w := performRequest(r, "GET", "/api/v1/tenant/performance/periods?page=1&per_page=10", nil)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestHandler_GetPerformancePeriodByID(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	// Create first
	createW := performRequest(r, "POST", "/api/v1/tenant/performance/periods", map[string]interface{}{
		"period_code": "FY",
		"period_type": "ANNUAL",
		"year":        2026,
	})

	var createResp map[string]interface{}
	json.Unmarshal(createW.Body.Bytes(), &createResp)
	data := createResp["data"].(map[string]interface{})
	id := data["id"].(string)

	w := performRequest(r, "GET", "/api/v1/tenant/performance/periods/"+id, nil)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestHandler_GetPerformancePeriodByID_NotFound(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "GET", "/api/v1/tenant/performance/periods/non-existent-id", nil)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandler_UpdatePerformancePeriod(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	createW := performRequest(r, "POST", "/api/v1/tenant/performance/periods", map[string]interface{}{
		"period_code": "OLD",
		"period_type": "QUARTERLY",
		"year":        2025,
	})
	var createResp map[string]interface{}
	json.Unmarshal(createW.Body.Bytes(), &createResp)
	id := createResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "PUT", "/api/v1/tenant/performance/periods/"+id, map[string]interface{}{
		"period_code": "NEW",
	})
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DeletePerformancePeriod(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	createW := performRequest(r, "POST", "/api/v1/tenant/performance/periods", map[string]interface{}{
		"period_code": "DEL",
		"period_type": "MONTHLY",
		"year":        2026,
	})
	var createResp map[string]interface{}
	json.Unmarshal(createW.Body.Bytes(), &createResp)
	id := createResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "DELETE", "/api/v1/tenant/performance/periods/"+id, nil)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// =========================================================================
// Performance Perspective Handler Tests
// =========================================================================

func TestHandler_CreatePerformancePerspective(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "POST", "/api/v1/tenant/performance/kpi/perspectives", map[string]interface{}{
		"name":       "Financial",
		"sort_order": 1,
	})
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListPerformancePerspectives(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	performRequest(r, "POST", "/api/v1/tenant/performance/kpi/perspectives", map[string]interface{}{
		"name": "Financial",
	})
	performRequest(r, "POST", "/api/v1/tenant/performance/kpi/perspectives", map[string]interface{}{
		"name": "Customer",
	})

	w := performRequest(r, "GET", "/api/v1/tenant/performance/kpi/perspectives", nil)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// =========================================================================
// Performance Template Handler Tests
// =========================================================================

func TestHandler_CreatePerformanceTemplate(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "POST", "/api/v1/tenant/performance/kpi/templates", map[string]interface{}{
		"organization_id": createTestOrgID(),
		"name":            "Manager Template",
	})
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Performance Evaluation Handler Tests
// =========================================================================

func TestHandler_CreatePerformanceEvaluation(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	// Create dependencies
	periodW := performRequest(r, "POST", "/api/v1/tenant/performance/periods", map[string]interface{}{
		"period_code": "FY", "period_type": "ANNUAL", "year": 2026,
	})
	var periodResp map[string]interface{}
	json.Unmarshal(periodW.Body.Bytes(), &periodResp)
	periodID := periodResp["data"].(map[string]interface{})["id"].(string)

	tmplW := performRequest(r, "POST", "/api/v1/tenant/performance/kpi/templates", map[string]interface{}{
		"organization_id": createTestOrgID(), "name": "Template",
	})
	var tmplResp map[string]interface{}
	json.Unmarshal(tmplW.Body.Bytes(), &tmplResp)
	tmplID := tmplResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/performance/kpi/evaluations", map[string]interface{}{
		"employee_id":     createTestUUID(),
		"organization_id": createTestOrgID(),
		"period_id":       periodID,
		"template_id":     tmplID,
	})
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdateEvaluationStatus(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	// Create evaluation
	periodW := performRequest(r, "POST", "/api/v1/tenant/performance/periods", map[string]interface{}{
		"period_code": "FY", "period_type": "ANNUAL", "year": 2026,
	})
	var periodResp map[string]interface{}
	json.Unmarshal(periodW.Body.Bytes(), &periodResp)
	periodID := periodResp["data"].(map[string]interface{})["id"].(string)

	tmplW := performRequest(r, "POST", "/api/v1/tenant/performance/kpi/templates", map[string]interface{}{
		"organization_id": createTestOrgID(), "name": "Template",
	})
	var tmplResp map[string]interface{}
	json.Unmarshal(tmplW.Body.Bytes(), &tmplResp)
	tmplID := tmplResp["data"].(map[string]interface{})["id"].(string)

	evalW := performRequest(r, "POST", "/api/v1/tenant/performance/kpi/evaluations", map[string]interface{}{
		"employee_id":     createTestUUID(),
		"organization_id": createTestOrgID(),
		"period_id":       periodID,
		"template_id":     tmplID,
	})
	var evalResp map[string]interface{}
	json.Unmarshal(evalW.Body.Bytes(), &evalResp)
	evalID := evalResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "PUT", "/api/v1/tenant/performance/kpi/evaluations/"+evalID+"/status", map[string]interface{}{
		"status": "PLAN_SUBMITTED",
	})
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateEvaluationDetail(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	// Create dependencies
	periodW := performRequest(r, "POST", "/api/v1/tenant/performance/periods", map[string]interface{}{
		"period_code": "FY", "period_type": "ANNUAL", "year": 2026,
	})
	var periodResp map[string]interface{}
	json.Unmarshal(periodW.Body.Bytes(), &periodResp)
	periodID := periodResp["data"].(map[string]interface{})["id"].(string)

	tmplW := performRequest(r, "POST", "/api/v1/tenant/performance/kpi/templates", map[string]interface{}{
		"organization_id": createTestOrgID(), "name": "Template",
	})
	var tmplResp map[string]interface{}
	json.Unmarshal(tmplW.Body.Bytes(), &tmplResp)
	tmplID := tmplResp["data"].(map[string]interface{})["id"].(string)

	perspW := performRequest(r, "POST", "/api/v1/tenant/performance/kpi/perspectives", map[string]interface{}{
		"name": "Financial",
	})
	var perspResp map[string]interface{}
	json.Unmarshal(perspW.Body.Bytes(), &perspResp)
	perspID := perspResp["data"].(map[string]interface{})["id"].(string)

	evalW := performRequest(r, "POST", "/api/v1/tenant/performance/kpi/evaluations", map[string]interface{}{
		"employee_id":     createTestUUID(),
		"organization_id": createTestOrgID(),
		"period_id":       periodID,
		"template_id":     tmplID,
	})
	var evalResp map[string]interface{}
	json.Unmarshal(evalW.Body.Bytes(), &evalResp)
	evalID := evalResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/performance/kpi/evaluation-details", map[string]interface{}{
		"performance_evaluation_id": evalID,
		"perspective_id":           perspID,
		"achievement_percentage":    90.0,
		"weight":                    40.0,
		"score":                     36.0,
	})
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreatePerformanceTarget(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	// Dependencies
	periodW := performRequest(r, "POST", "/api/v1/tenant/performance/periods", map[string]interface{}{
		"period_code": "FY", "period_type": "ANNUAL", "year": 2026,
	})
	var periodResp map[string]interface{}
	json.Unmarshal(periodW.Body.Bytes(), &periodResp)
	periodID := periodResp["data"].(map[string]interface{})["id"].(string)

	tmplW := performRequest(r, "POST", "/api/v1/tenant/performance/kpi/templates", map[string]interface{}{
		"organization_id": createTestOrgID(), "name": "Template",
	})
	var tmplResp map[string]interface{}
	json.Unmarshal(tmplW.Body.Bytes(), &tmplResp)
	tmplID := tmplResp["data"].(map[string]interface{})["id"].(string)

	perspW := performRequest(r, "POST", "/api/v1/tenant/performance/kpi/perspectives", map[string]interface{}{
		"name": "Financial",
	})
	var perspResp map[string]interface{}
	json.Unmarshal(perspW.Body.Bytes(), &perspResp)
	perspID := perspResp["data"].(map[string]interface{})["id"].(string)

	evalW := performRequest(r, "POST", "/api/v1/tenant/performance/kpi/evaluations", map[string]interface{}{
		"employee_id":     createTestUUID(),
		"organization_id": createTestOrgID(),
		"period_id":       periodID,
		"template_id":     tmplID,
	})
	var evalResp map[string]interface{}
	json.Unmarshal(evalW.Body.Bytes(), &evalResp)
	evalID := evalResp["data"].(map[string]interface{})["id"].(string)

	indW := performRequest(r, "POST", "/api/v1/tenant/performance/kpi/indicators", map[string]interface{}{
		"performance_template_id": tmplID,
		"perspective_id":         perspID,
		"indicator_type":         "MAXIMIZATION",
		"title":                  "Revenue",
		"weight":                 100.0,
	})
	var indResp map[string]interface{}
	json.Unmarshal(indW.Body.Bytes(), &indResp)
	indID := indResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/performance/kpi/targets", map[string]interface{}{
		"performance_evaluation_id": evalID,
		"indicator_id":             indID,
		"target_value":             1000000,
		"weight":                   50.0,
	})
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}
