package recruitment

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/modules/competency"
)

func setupTestRouter() (*gin.Engine, *Handler, func()) {
	gin.SetMode(gin.TestMode)
	db, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger := zap.NewNop()
	svc := NewService(repo, logger)
	seedDefaultRecruitmentStages(db)
	handler := NewHandler(svc)

	r := gin.New()
	rg := r.Group("/api/v1/tenant")
	RegisterRoutes(rg, handler)

	return r, handler, cleanup
}

// setupTestRouterWithUserID mirrors setupTestRouter but installs a
// middleware that sets Gin's "user_id" context key, the same key the real
// auth middleware (internal/pkg/middleware/auth.go) populates from JWT
// claims. Handlers read it via c.GetString("user_id") (see Fix 2, G-5 final
// review — changed_by plumbing).
func setupTestRouterWithUserID(userID string) (*gin.Engine, *Handler, func()) {
	gin.SetMode(gin.TestMode)
	db, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger := zap.NewNop()
	svc := NewService(repo, logger)
	seedDefaultRecruitmentStages(db)
	handler := NewHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	rg := r.Group("/api/v1/tenant")
	RegisterRoutes(rg, handler)

	return r, handler, cleanup
}

func performRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// =========================================================================
// Requisition Handler Tests
// =========================================================================

func TestHandler_CreateRequisition(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	req := CreateRequisitionRequest{
		OrganizationID: createTestOrgID(),
		Title:          "Backend Engineer",
		EmploymentType: "FULL_TIME",
		SlotsAvailable: intPtr(2),
	}
	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Fatal("expected success = true")
	}
}

func TestHandler_ListRequisitions(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	orgID := createTestOrgID()
	performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: orgID, Title: "Req 1",
	})
	performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: orgID, Title: "Req 2",
	})

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/requisitions?organization_id="+orgID, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Fatal("expected success = true")
	}
	if resp["total"] != float64(2) {
		t.Fatalf("expected total = 2, got %v", resp["total"])
	}
}

func TestHandler_GetRequisitionByID(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Get Test",
	})
	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	data := created["data"].(map[string]interface{})
	id := data["id"].(string)

	w2 := performRequest(r, "GET", "/api/v1/tenant/recruitment/requisitions/"+id, nil)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w2.Code, w2.Body.String())
	}
}

func TestHandler_GetRequisitionByID_NotFound(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/requisitions/"+uuid.New().String(), nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestHandler_UpdateRequisition(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Original",
	})
	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	id := created["data"].(map[string]interface{})["id"].(string)

	newTitle := "Updated Title"
	w2 := performRequest(r, "PUT", "/api/v1/tenant/recruitment/requisitions/"+id, UpdateRequisitionRequest{
		Title: &newTitle,
	})
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w2.Code, w2.Body.String())
	}
}

func TestHandler_DeleteRequisition(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Delete Me",
	})
	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	id := created["data"].(map[string]interface{})["id"].(string)

	w2 := performRequest(r, "DELETE", "/api/v1/tenant/recruitment/requisitions/"+id, nil)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w2.Code, w2.Body.String())
	}
}

// =========================================================================
// Candidate Handler Tests
// =========================================================================

func TestHandler_CreateCandidate(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Jane", LastName: "Smith", Email: "jane@example.com",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListCandidates(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Alice", LastName: "A", Email: "alice@test.com",
	})
	performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Bob", LastName: "B", Email: "bob@test.com",
	})

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/candidates", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["total"] != float64(2) {
		t.Fatalf("expected total = 2, got %v", resp["total"])
	}
}

func TestHandler_GetCandidateByID(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Get", LastName: "Test", Email: "get@test.com",
	})
	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	id := created["data"].(map[string]interface{})["id"].(string)

	w2 := performRequest(r, "GET", "/api/v1/tenant/recruitment/candidates/"+id, nil)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}
}

func TestHandler_UpdateCandidate(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Update", LastName: "Test", Email: "update@test.com",
	})
	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	id := created["data"].(map[string]interface{})["id"].(string)

	newPhone := "08111111111"
	w2 := performRequest(r, "PUT", "/api/v1/tenant/recruitment/candidates/"+id, UpdateCandidateRequest{
		Phone: &newPhone,
	})
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w2.Code, w2.Body.String())
	}
}

func TestHandler_DeleteCandidate(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Del", LastName: "Cand", Email: "delcand@test.com",
	})
	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	id := created["data"].(map[string]interface{})["id"].(string)

	w2 := performRequest(r, "DELETE", "/api/v1/tenant/recruitment/candidates/"+id, nil)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}
}

// =========================================================================
// Application Handler Tests
// =========================================================================

func TestHandler_CreateApplication(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	// Create prerequisite via handler
	reqW := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Engineer",
	})
	var reqResp map[string]interface{}
	json.Unmarshal(reqW.Body.Bytes(), &reqResp)
	reqID := reqResp["data"].(map[string]interface{})["id"].(string)

	candW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "App", LastName: "Cand", Email: "appcand@test.com",
	})
	var candResp map[string]interface{}
	json.Unmarshal(candW.Body.Bytes(), &candResp)
	candID := candResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/applications", CreateApplicationRequest{
		RequisitionID: reqID,
		CandidateID:   candID,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetApplicationByID(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	// Create requisition & candidate
	reqW := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "My Req",
	})
	var reqResp map[string]interface{}
	json.Unmarshal(reqW.Body.Bytes(), &reqResp)
	rid := reqResp["data"].(map[string]interface{})["id"].(string)

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "GetApp", LastName: "Test", Email: "getapp@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	appW := performRequest(r, "POST", "/api/v1/tenant/recruitment/applications", CreateApplicationRequest{
		RequisitionID: rid, CandidateID: cid,
	})
	var appResp map[string]interface{}
	json.Unmarshal(appW.Body.Bytes(), &appResp)
	appID := appResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/applications/"+appID, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHandler_UpdateApplicationStatus(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	reqW := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "My Req",
	})
	var reqResp map[string]interface{}
	json.Unmarshal(reqW.Body.Bytes(), &reqResp)
	rid := reqResp["data"].(map[string]interface{})["id"].(string)

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Status", LastName: "Test", Email: "status@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	appW := performRequest(r, "POST", "/api/v1/tenant/recruitment/applications", CreateApplicationRequest{
		RequisitionID: rid, CandidateID: cid,
	})
	var appResp map[string]interface{}
	json.Unmarshal(appW.Body.Bytes(), &appResp)
	appID := appResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "PUT", "/api/v1/tenant/recruitment/applications/"+appID+"/status", UpdateApplicationStatusRequest{
		Status: "SHORTLISTED",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// TestHandler_UpdateApplicationStatus_ChangedByFromContext verifies Fix 2
// (G-5 final review): a manual PUT .../status with "user_id" set in the
// Gin context (as the real auth middleware does from JWT claims) produces
// a history row whose changed_by matches that user id.
func TestHandler_UpdateApplicationStatus_ChangedByFromContext(t *testing.T) {
	actorID := uuid.New().String()
	r, _, cleanup := setupTestRouterWithUserID(actorID)
	defer cleanup()

	reqW := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "My Req",
	})
	var reqResp map[string]interface{}
	json.Unmarshal(reqW.Body.Bytes(), &reqResp)
	rid := reqResp["data"].(map[string]interface{})["id"].(string)

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "ChangedBy", LastName: "Test", Email: "changedby@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	appW := performRequest(r, "POST", "/api/v1/tenant/recruitment/applications", CreateApplicationRequest{
		RequisitionID: rid, CandidateID: cid,
	})
	var appResp map[string]interface{}
	json.Unmarshal(appW.Body.Bytes(), &appResp)
	appID := appResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "PUT", "/api/v1/tenant/recruitment/applications/"+appID+"/status", UpdateApplicationStatusRequest{
		Status: "SHORTLISTED",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	histW := performRequest(r, "GET", "/api/v1/tenant/recruitment/applications/"+appID+"/history", nil)
	if histW.Code != http.StatusOK {
		t.Fatalf("expected 200 on history, got %d: %s", histW.Code, histW.Body.String())
	}
	var histResp map[string]interface{}
	json.Unmarshal(histW.Body.Bytes(), &histResp)
	rows, ok := histResp["data"].([]interface{})
	if !ok || len(rows) == 0 {
		t.Fatalf("expected non-empty history data, got: %s", histW.Body.String())
	}

	var found bool
	for _, row := range rows {
		m := row.(map[string]interface{})
		if cb, ok := m["changed_by"].(string); ok && cb == actorID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected a history row with changed_by=%s, got: %s", actorID, histW.Body.String())
	}
}

func TestHandler_DeleteApplication(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	reqW := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Req",
	})
	var reqResp map[string]interface{}
	json.Unmarshal(reqW.Body.Bytes(), &reqResp)
	rid := reqResp["data"].(map[string]interface{})["id"].(string)

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "DelApp", LastName: "Test", Email: "delapp@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	appW := performRequest(r, "POST", "/api/v1/tenant/recruitment/applications", CreateApplicationRequest{
		RequisitionID: rid, CandidateID: cid,
	})
	var appResp map[string]interface{}
	json.Unmarshal(appW.Body.Bytes(), &appResp)
	appID := appResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "DELETE", "/api/v1/tenant/recruitment/applications/"+appID, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

// =========================================================================
// Interview Handler Tests
// =========================================================================

func TestHandler_CreateInterview(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	reqW := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "My Req",
	})
	var reqResp map[string]interface{}
	json.Unmarshal(reqW.Body.Bytes(), &reqResp)
	rid := reqResp["data"].(map[string]interface{})["id"].(string)

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "IntCand", LastName: "Test", Email: "intcand@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	appW := performRequest(r, "POST", "/api/v1/tenant/recruitment/applications", CreateApplicationRequest{
		RequisitionID: rid, CandidateID: cid,
	})
	var appResp map[string]interface{}
	json.Unmarshal(appW.Body.Bytes(), &appResp)
	appID := appResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/interviews", CreateInterviewRequest{
		ApplicationID: appID,
		InterviewerID: createTestUUID(),
		Stage:         "TECHNICAL",
		ScheduledAt:   1760000000,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListInterviews(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/interviews", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHandler_GetInterviewByID_NotFound(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/interviews/"+uuid.New().String(), nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// =========================================================================
// Onboarding Task Template Handler Tests
// =========================================================================

func TestHandler_CreateOnboardingTaskTemplate(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/onboarding-task-templates", CreateOnboardingTaskTemplateRequest{
		Name:      "IT Account Setup",
		Category:  "IT",
		DayOffset: intPtr(-7),
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListOnboardingTaskTemplates(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	performRequest(r, "POST", "/api/v1/tenant/recruitment/onboarding-task-templates", CreateOnboardingTaskTemplateRequest{
		Name: "IT Setup", Category: "IT", DayOffset: intPtr(0),
	})
	performRequest(r, "POST", "/api/v1/tenant/recruitment/onboarding-task-templates", CreateOnboardingTaskTemplateRequest{
		Name: "HR Setup", Category: "HR", DayOffset: intPtr(0),
	})

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/onboarding-task-templates?category=IT", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["total"] != float64(1) {
		t.Fatalf("expected total = 1, got %v", resp["total"])
	}
}

func TestHandler_UpdateOnboardingTaskTemplate(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/onboarding-task-templates", CreateOnboardingTaskTemplateRequest{
		Name: "Original Template", Category: "LEGAL", DayOffset: intPtr(0),
	})
	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	id := created["data"].(map[string]interface{})["id"].(string)

	newName := "Updated Template"
	w2 := performRequest(r, "PUT", "/api/v1/tenant/recruitment/onboarding-task-templates/"+id, UpdateOnboardingTaskTemplateRequest{
		Name: &newName,
	})
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w2.Code, w2.Body.String())
	}
}

func TestHandler_DeleteOnboardingTaskTemplate(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/onboarding-task-templates", CreateOnboardingTaskTemplateRequest{
		Name: "To Delete", Category: "TEST", DayOffset: intPtr(0),
	})
	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	id := created["data"].(map[string]interface{})["id"].(string)

	w2 := performRequest(r, "DELETE", "/api/v1/tenant/recruitment/onboarding-task-templates/"+id, nil)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}
}

// =========================================================================
// Employee Onboarding Handler Tests
// =========================================================================

func TestHandler_CreateEmployeeOnboarding(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	// Create a task template first (auto-created on onboarding)
	performRequest(r, "POST", "/api/v1/tenant/recruitment/onboarding-task-templates", CreateOnboardingTaskTemplateRequest{
		Name: "Contract Signing", Category: "LEGAL", DayOffset: intPtr(0),
	})

	reqW := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "My Req",
	})
	var reqResp map[string]interface{}
	json.Unmarshal(reqW.Body.Bytes(), &reqResp)
	rid := reqResp["data"].(map[string]interface{})["id"].(string)

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "OnbCand", LastName: "Test", Email: "onbcand@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	appW := performRequest(r, "POST", "/api/v1/tenant/recruitment/applications", CreateApplicationRequest{
		RequisitionID: rid, CandidateID: cid,
	})
	var appResp map[string]interface{}
	json.Unmarshal(appW.Body.Bytes(), &appResp)
	appID := appResp["data"].(map[string]interface{})["id"].(string)

	// Accept candidate first
	performRequest(r, "PUT", "/api/v1/tenant/recruitment/applications/"+appID+"/status", UpdateApplicationStatusRequest{
		Status: "ACCEPTED",
	})

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/employee-onboardings", CreateEmployeeOnboardingRequest{
		EmployeeID:    createTestUUID(),
		ApplicationID: appID,
		StartDate:     "2026-08-01",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListEmployeeOnboardings(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/employee-onboardings", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

// =========================================================================
// Onboarding Task Item Handler Tests
// =========================================================================

func TestHandler_CreateOnboardingTaskItem(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	// Create a template first
	performRequest(r, "POST", "/api/v1/tenant/recruitment/onboarding-task-templates", CreateOnboardingTaskTemplateRequest{
		Name: "Contract", Category: "LEGAL", DayOffset: intPtr(0),
	})

	// Create onboarding
	reqW := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "My Req",
	})
	var reqResp map[string]interface{}
	json.Unmarshal(reqW.Body.Bytes(), &reqResp)
	rid := reqResp["data"].(map[string]interface{})["id"].(string)

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "ItemCand", LastName: "Test", Email: "itemcand@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	appW := performRequest(r, "POST", "/api/v1/tenant/recruitment/applications", CreateApplicationRequest{
		RequisitionID: rid, CandidateID: cid,
	})
	var appResp map[string]interface{}
	json.Unmarshal(appW.Body.Bytes(), &appResp)
	appID := appResp["data"].(map[string]interface{})["id"].(string)

	performRequest(r, "PUT", "/api/v1/tenant/recruitment/applications/"+appID+"/status", UpdateApplicationStatusRequest{
		Status: "ACCEPTED",
	})

	onbW := performRequest(r, "POST", "/api/v1/tenant/recruitment/employee-onboardings", CreateEmployeeOnboardingRequest{
		EmployeeID: createTestUUID(), ApplicationID: appID, StartDate: "2026-08-01",
	})
	var onbResp map[string]interface{}
	json.Unmarshal(onbW.Body.Bytes(), &onbResp)
	onbID := onbResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/onboarding-task-items", CreateOnboardingTaskItemRequest{
		EmployeeOnboardingID: onbID,
		Name:                "Custom Task",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListOnboardingTaskItems(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	performRequest(r, "POST", "/api/v1/tenant/recruitment/onboarding-task-templates", CreateOnboardingTaskTemplateRequest{
		Name: "Contract", Category: "LEGAL", DayOffset: intPtr(0),
	})

	reqW := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Req",
	})
	var reqResp map[string]interface{}
	json.Unmarshal(reqW.Body.Bytes(), &reqResp)
	rid := reqResp["data"].(map[string]interface{})["id"].(string)

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "ListItems", LastName: "Test", Email: "listitems@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	appW := performRequest(r, "POST", "/api/v1/tenant/recruitment/applications", CreateApplicationRequest{
		RequisitionID: rid, CandidateID: cid,
	})
	var appResp map[string]interface{}
	json.Unmarshal(appW.Body.Bytes(), &appResp)
	appID := appResp["data"].(map[string]interface{})["id"].(string)

	performRequest(r, "PUT", "/api/v1/tenant/recruitment/applications/"+appID+"/status", UpdateApplicationStatusRequest{
		Status: "ACCEPTED",
	})

	onbW := performRequest(r, "POST", "/api/v1/tenant/recruitment/employee-onboardings", CreateEmployeeOnboardingRequest{
		EmployeeID: createTestUUID(), ApplicationID: appID, StartDate: "2026-08-01",
	})
	var onbResp map[string]interface{}
	json.Unmarshal(onbW.Body.Bytes(), &onbResp)
	onbID := onbResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/employee-onboardings/"+onbID+"/task-items", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHandler_UpdateOnboardingTaskItem(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	performRequest(r, "POST", "/api/v1/tenant/recruitment/onboarding-task-templates", CreateOnboardingTaskTemplateRequest{
		Name: "Contract", Category: "LEGAL", DayOffset: intPtr(0),
	})

	reqW := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Req",
	})
	var reqResp map[string]interface{}
	json.Unmarshal(reqW.Body.Bytes(), &reqResp)
	rid := reqResp["data"].(map[string]interface{})["id"].(string)

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "UpdItem", LastName: "Test", Email: "upditem@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	appW := performRequest(r, "POST", "/api/v1/tenant/recruitment/applications", CreateApplicationRequest{
		RequisitionID: rid, CandidateID: cid,
	})
	var appResp map[string]interface{}
	json.Unmarshal(appW.Body.Bytes(), &appResp)
	appID := appResp["data"].(map[string]interface{})["id"].(string)

	performRequest(r, "PUT", "/api/v1/tenant/recruitment/applications/"+appID+"/status", UpdateApplicationStatusRequest{
		Status: "ACCEPTED",
	})

	onbW := performRequest(r, "POST", "/api/v1/tenant/recruitment/employee-onboardings", CreateEmployeeOnboardingRequest{
		EmployeeID: createTestUUID(), ApplicationID: appID, StartDate: "2026-08-01",
	})
	var onbResp map[string]interface{}
	json.Unmarshal(onbW.Body.Bytes(), &onbResp)
	onbID := onbResp["data"].(map[string]interface{})["id"].(string)

	// Create task item via onboarding (auto-created from template)
	itemW := performRequest(r, "POST", "/api/v1/tenant/recruitment/onboarding-task-items", CreateOnboardingTaskItemRequest{
		EmployeeOnboardingID: onbID,
		Name:                "Initial Task",
	})
	var itemResp map[string]interface{}
	json.Unmarshal(itemW.Body.Bytes(), &itemResp)
	itemID := itemResp["data"].(map[string]interface{})["id"].(string)

	completed := true
	w := performRequest(r, "PUT", "/api/v1/tenant/recruitment/onboarding-task-items/"+itemID, UpdateOnboardingTaskItemRequest{
		IsCompleted: &completed,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DeleteOnboardingTaskItem(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "DELETE", "/api/v1/tenant/recruitment/onboarding-task-items/"+uuid.New().String(), nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestHandler_GetApplicationHistory(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	reqW := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "My Req",
	})
	var reqResp map[string]interface{}
	json.Unmarshal(reqW.Body.Bytes(), &reqResp)
	rid := reqResp["data"].(map[string]interface{})["id"].(string)

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Hist", LastName: "Test", Email: "hist@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	appW := performRequest(r, "POST", "/api/v1/tenant/recruitment/applications", CreateApplicationRequest{
		RequisitionID: rid, CandidateID: cid,
	})
	var appResp map[string]interface{}
	json.Unmarshal(appW.Body.Bytes(), &appResp)
	appID := appResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/applications/"+appID+"/history", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetApplicationHistory_NotFound(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/applications/"+uuid.New().String()+"/history", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestHandler_UpdateApplicationStatus_InvalidTransitionReturns400(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	reqW := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "My Req",
	})
	var reqResp map[string]interface{}
	json.Unmarshal(reqW.Body.Bytes(), &reqResp)
	rid := reqResp["data"].(map[string]interface{})["id"].(string)

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Bad", LastName: "Transition", Email: "badtrans@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	appW := performRequest(r, "POST", "/api/v1/tenant/recruitment/applications", CreateApplicationRequest{
		RequisitionID: rid, CandidateID: cid,
	})
	var appResp map[string]interface{}
	json.Unmarshal(appW.Body.Bytes(), &appResp)
	appID := appResp["data"].(map[string]interface{})["id"].(string)

	performRequest(r, "PUT", "/api/v1/tenant/recruitment/applications/"+appID+"/status", UpdateApplicationStatusRequest{Status: "REJECTED"})

	w := performRequest(r, "PUT", "/api/v1/tenant/recruitment/applications/"+appID+"/status", UpdateApplicationStatusRequest{Status: "SCREENED"})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for transition out of terminal status, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateCandidateEducation(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Edu", LastName: "Handler", Email: "eduhandler@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/educations", CreateCandidateEducationRequest{
		InstitutionName: "Universitas Test",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListCandidateEducations(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "List", LastName: "EduH", Email: "listeduh@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/educations", CreateCandidateEducationRequest{InstitutionName: "SMA 1"})

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/candidates/"+cid+"/educations", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdateAndDeleteCandidateEducation(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Upd", LastName: "EduH", Email: "updeduh@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	eW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/educations", CreateCandidateEducationRequest{InstitutionName: "Original"})
	var eResp map[string]interface{}
	json.Unmarshal(eW.Body.Bytes(), &eResp)
	eid := eResp["data"].(map[string]interface{})["id"].(string)

	newName := "Updated"
	w := performRequest(r, "PUT", "/api/v1/tenant/recruitment/educations/"+eid, UpdateCandidateEducationRequest{InstitutionName: &newName})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	w2 := performRequest(r, "DELETE", "/api/v1/tenant/recruitment/educations/"+eid, nil)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}
}

func TestHandler_CreateCandidateWorkExperience(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Exp", LastName: "Handler", Email: "exphandler@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/work-experiences", CreateCandidateWorkExperienceRequest{
		CompanyName: "Acme", JobTitle: "Engineer", StartDate: "2020-01-01",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateCandidateSkill(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	svc := NewService(repo, zap.NewNop())
	seedDefaultRecruitmentStages(db)
	handler := NewHandler(svc)
	r := gin.New()
	rg := r.Group("/api/v1/tenant")
	RegisterRoutes(rg, handler)

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Skill", LastName: "Handler", Email: "skillhandler@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	comp := &competency.Competency{Name: "Go"}
	if err := db.Create(comp).Error; err != nil {
		t.Fatalf("failed to seed competency: %v", err)
	}

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/skills", CreateCandidateSkillRequest{
		CompetencyID: comp.ID.String(),
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateCandidateCertification(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Cert", LastName: "Handler", Email: "certhandler@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/certifications", CreateCandidateCertificationRequest{
		Name: "AWS SAA",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateCandidateDocument(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Doc", LastName: "Handler", Email: "dochandler@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/documents", CreateCandidateDocumentRequest{
		DocumentType: "RESUME", Name: "resume.pdf", FileURL: "/uploads/attachments/x.pdf",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateCandidateDocument_InvalidDocumentType(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Bad", LastName: "Type", Email: "badtype@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/documents", CreateCandidateDocumentRequest{
		DocumentType: "NOT_A_REAL_TYPE", Name: "x.pdf", FileURL: "/u/x.pdf",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid document_type, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListCandidateDocuments(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "List", LastName: "DocH", Email: "listdoch@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/documents", CreateCandidateDocumentRequest{Name: "a.pdf", FileURL: "/u/a.pdf"})

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/candidates/"+cid+"/documents", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateCandidateConsent(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Consent", LastName: "Handler", Email: "consenthandler@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/consents", CreateCandidateConsentRequest{
		Action: "GRANTED",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateCandidateConsent_InvalidAction(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Bad", LastName: "Action", Email: "badaction@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/consents", CreateCandidateConsentRequest{
		Action: "MAYBE",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid action, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateCandidateConsent_ChangedByFromContext(t *testing.T) {
	userID := uuid.New().String()
	r, _, cleanup := setupTestRouterWithUserID(userID)
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Actor", LastName: "Handler", Email: "actorhandler@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/consents", CreateCandidateConsentRequest{
		Action: "GRANTED",
	})
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["changed_by"] != userID {
		t.Errorf("expected changed_by %s, got %v", userID, data["changed_by"])
	}
}

func TestHandler_ListCandidateConsents(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "List", LastName: "ConsentH", Email: "listconsenth@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/consents", CreateCandidateConsentRequest{Action: "GRANTED"})

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/candidates/"+cid+"/consents", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func createTestApplicationForScreening(r *gin.Engine, email string) string {
	reqW := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Engineer",
	})
	var reqResp map[string]interface{}
	json.Unmarshal(reqW.Body.Bytes(), &reqResp)
	reqID := reqResp["data"].(map[string]interface{})["id"].(string)

	candW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Scr", LastName: "H", Email: email,
	})
	var candResp map[string]interface{}
	json.Unmarshal(candW.Body.Bytes(), &candResp)
	candID := candResp["data"].(map[string]interface{})["id"].(string)

	appW := performRequest(r, "POST", "/api/v1/tenant/recruitment/applications", CreateApplicationRequest{
		RequisitionID: reqID, CandidateID: candID,
	})
	var appResp map[string]interface{}
	json.Unmarshal(appW.Body.Bytes(), &appResp)
	return appResp["data"].(map[string]interface{})["id"].(string)
}

func TestHandler_CreateApplicationScreening(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	appID := createTestApplicationForScreening(r, "screenh1@test.com")

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/applications/"+appID+"/screenings", CreateApplicationScreeningRequest{
		Result: "PASS",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListApplicationScreenings(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	appID := createTestApplicationForScreening(r, "screenh2@test.com")
	performRequest(r, "POST", "/api/v1/tenant/recruitment/applications/"+appID+"/screenings", CreateApplicationScreeningRequest{Result: "HOLD"})

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/applications/"+appID+"/screenings", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateAssessment(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/assessments", CreateAssessmentRequest{
		Name: "Technical Test Batch March", Type: "TECHNICAL",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateAssessment_InvalidType(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/assessments", CreateAssessmentRequest{
		Name: "Bad Type", Type: "NOT_A_REAL_TYPE",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_AddAndListAssessmentParticipants(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	assessW := performRequest(r, "POST", "/api/v1/tenant/recruitment/assessments", CreateAssessmentRequest{Name: "Batch"})
	var assessResp map[string]interface{}
	json.Unmarshal(assessW.Body.Bytes(), &assessResp)
	assessID := assessResp["data"].(map[string]interface{})["id"].(string)

	appID := createTestApplicationForScreening(r, "assesspart-h1@test.com")

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/assessments/"+assessID+"/participants", AddAssessmentParticipantRequest{
		ApplicationID: appID,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	listW := performRequest(r, "GET", "/api/v1/tenant/recruitment/assessments/"+assessID+"/participants", nil)
	if listW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", listW.Code, listW.Body.String())
	}
}

func createTestInterviewForScorecard(r *gin.Engine, email string) string {
	appID := createTestApplicationForScreening(r, email)
	ivW := performRequest(r, "POST", "/api/v1/tenant/recruitment/interviews", CreateInterviewRequest{
		ApplicationID: appID, InterviewerID: createTestUUID(), Stage: "FIRST_INTERVIEW", ScheduledAt: 1760000000,
	})
	var ivResp map[string]interface{}
	json.Unmarshal(ivW.Body.Bytes(), &ivResp)
	return ivResp["data"].(map[string]interface{})["id"].(string)
}

func TestHandler_AddAndListInterviewers(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	ivID := createTestInterviewForScorecard(r, "ivh1@test.com")

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/interviews/"+ivID+"/interviewers", AddInterviewerRequest{
		EmployeeID: createTestUUID(), Role: "HR",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	listW := performRequest(r, "GET", "/api/v1/tenant/recruitment/interviews/"+ivID+"/interviewers", nil)
	if listW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", listW.Code, listW.Body.String())
	}
}

func TestHandler_AddScorecardItemAndComplete(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	ivID := createTestInterviewForScorecard(r, "ivh2@test.com")

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/interviews/"+ivID+"/scorecard-items", AddScorecardItemRequest{
		Criterion: "Technical Skill", Weight: 100,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	completeW := performRequest(r, "POST", "/api/v1/tenant/recruitment/interviews/"+ivID+"/complete", nil)
	if completeW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", completeW.Code, completeW.Body.String())
	}
}

func TestHandler_CreateAndListRequisitionRequirement(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	reqW := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Engineer",
	})
	var reqResp map[string]interface{}
	json.Unmarshal(reqW.Body.Bytes(), &reqResp)
	rid := reqResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions/"+rid+"/requirements", CreateRequisitionRequirementRequest{
		RequirementType: "EXPERIENCE_YEARS", Name: "Min Experience",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	listW := performRequest(r, "GET", "/api/v1/tenant/recruitment/requisitions/"+rid+"/requirements", nil)
	if listW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", listW.Code, listW.Body.String())
	}
}

func TestHandler_CreateAndListRequisitionCompetency(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	svc := NewService(repo, zap.NewNop())
	seedDefaultRecruitmentStages(db)
	handler := NewHandler(svc)
	r := gin.New()
	rg := r.Group("/api/v1/tenant")
	RegisterRoutes(rg, handler)

	reqW := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Engineer",
	})
	var reqResp map[string]interface{}
	json.Unmarshal(reqW.Body.Bytes(), &reqResp)
	rid := reqResp["data"].(map[string]interface{})["id"].(string)

	comp := &competency.Competency{Name: "Go"}
	if err := db.Create(comp).Error; err != nil {
		t.Fatalf("failed to seed competency: %v", err)
	}

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions/"+rid+"/competencies", CreateRequisitionCompetencyRequest{
		CompetencyID: comp.ID.String(),
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	listW := performRequest(r, "GET", "/api/v1/tenant/recruitment/requisitions/"+rid+"/competencies", nil)
	if listW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", listW.Code, listW.Body.String())
	}
}

func TestHandler_GetCandidateMatchScore(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	svc := NewService(repo, zap.NewNop())
	seedDefaultRecruitmentStages(db)
	handler := NewHandler(svc)
	r := gin.New()
	rg := r.Group("/api/v1/tenant")
	RegisterRoutes(rg, handler)

	reqW := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Engineer",
	})
	var reqResp map[string]interface{}
	json.Unmarshal(reqW.Body.Bytes(), &reqResp)
	rid := reqResp["data"].(map[string]interface{})["id"].(string)

	candW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Match", LastName: "H", Email: "matchh@test.com",
	})
	var candResp map[string]interface{}
	json.Unmarshal(candW.Body.Bytes(), &candResp)
	candID := candResp["data"].(map[string]interface{})["id"].(string)

	appW := performRequest(r, "POST", "/api/v1/tenant/recruitment/applications", CreateApplicationRequest{
		RequisitionID: rid, CandidateID: candID,
	})
	var appResp map[string]interface{}
	json.Unmarshal(appW.Body.Bytes(), &appResp)
	appID := appResp["data"].(map[string]interface{})["id"].(string)

	comp := &competency.Competency{Name: "Go"}
	if err := db.Create(comp).Error; err != nil {
		t.Fatalf("failed to seed competency: %v", err)
	}
	performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions/"+rid+"/competencies", CreateRequisitionCompetencyRequest{
		CompetencyID: comp.ID.String(),
	})

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/applications/"+appID+"/match-score", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetRecruitmentAnalyticsSummary(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/analytics/summary", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetRecruitmentAnalyticsSummary_WithDateRange(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/analytics/summary?from=1000&to=9999999999999999", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
