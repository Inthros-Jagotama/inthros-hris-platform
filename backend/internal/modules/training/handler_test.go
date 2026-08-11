package training

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// setupRouter membuat router test dengan handler training.
func setupRouter(t *testing.T) (*gin.Engine, *Handler) {
	t.Helper()
	svc := testSvc(t)
	handler := NewHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	trn := r.Group("/api/v1/tenant/trainings")
	{
		// Categories
		trn.POST("/categories", handler.CreateCategory)
		trn.GET("/categories", handler.ListCategories)
		trn.GET("/categories/:id", handler.GetCategoryByID)
		trn.PUT("/categories/:id", handler.UpdateCategory)
		trn.DELETE("/categories/:id", handler.DeleteCategory)

		// Courses
		trn.POST("/courses", handler.CreateCourse)
		trn.GET("/courses", handler.ListCourses)
		trn.GET("/courses/:id", handler.GetCourseByID)
		trn.PUT("/courses/:id", handler.UpdateCourse)
		trn.DELETE("/courses/:id", handler.DeleteCourse)

		// Sessions
		trn.POST("/sessions", handler.CreateSession)
		trn.GET("/sessions", handler.ListSessions)
		trn.GET("/sessions/:id", handler.GetSessionByID)
		trn.PUT("/sessions/:id", handler.UpdateSession)
		trn.PUT("/sessions/:id/status", handler.UpdateSessionStatus)
		trn.DELETE("/sessions/:id", handler.DeleteSession)

		// Participants
		trn.POST("/participants", handler.CreateParticipant)
		trn.GET("/participants", handler.ListParticipants)
		trn.GET("/participants/:id", handler.GetParticipantByID)
		trn.PUT("/participants/:id", handler.UpdateParticipant)
		trn.DELETE("/participants/:id", handler.DeleteParticipant)
	}
	return r, handler
}

// TestRegisterRoutes_NoWildcardConflict memastikan registrasi SELURUH route
// training (via RegisterRoutes) tidak memicu panic konflik wildcard Gin.
// Regression guard: pernah panic "':form_id' ... conflicts with existing wildcard
// ':id' ... /evaluation-forms/:id" saat route questions & CRUD form hidup
// bersamaan — semua wildcard pada level path yang sama kini bernama konsisten
// (evaluation-forms memakai :form_id, sisanya :id).
func TestRegisterRoutes_NoWildcardConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	rg := r.Group("/api/v1/tenant")
	// Registrasi seluruh route — panic konflik wildcard Gin terjadi di sini
	// bila ada wildcard berbeda pada level path yang sama.
	RegisterRoutes(rg, NewHandler(testSvc(t)))

	// Request probe: route list (selalu 200 walaupun kosong) membuktikan
	// tree Gin ter-resolve tanpa konflik wildcard. Route detail yang
	// mengembalikan 404 legit (record tidak ada) tidak di-probe.
	for _, path := range []string{
		"/api/v1/tenant/trainings/evaluation-forms/00000000-0000-0000-0000-000000000000/questions",
		"/api/v1/tenant/trainings/courses/00000000-0000-0000-0000-000000000000/objectives",
		"/api/v1/tenant/trainings/sessions/00000000-0000-0000-0000-000000000000/costs",
		"/api/v1/tenant/trainings/reports/dashboard",
	} {
		req := httptest.NewRequest("GET", path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code == http.StatusNotFound {
			t.Errorf("route %s not reachable (wildcard conflict?)", path)
		}
	}
}

func TestHandler_CreateCategory(t *testing.T) {
	r, _ := setupRouter(t)

	body := `{"code":"HR","name":"Human Resources","description":"HR training"}`
	req := httptest.NewRequest("POST", "/api/v1/tenant/trainings/categories", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Error("expected success=true")
	}
}

func TestHandler_ListCategories(t *testing.T) {
	r, h := setupRouter(t)

	// Seed some categories
	svc := h.svc
	desc := "Test"
	svc.CreateCategory(testCtx(), CreateTrainingCategoryRequest{Code: "A", Name: "Alpha", Description: &desc})
	svc.CreateCategory(testCtx(), CreateTrainingCategoryRequest{Code: "B", Name: "Beta", Description: &desc})

	req := httptest.NewRequest("GET", "/api/v1/tenant/trainings/categories?page=1&per_page=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestHandler_GetCategoryByID_NotFound(t *testing.T) {
	r, _ := setupRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/tenant/trainings/categories/00000000-0000-0000-0000-000000000000", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandler_CreateCategory_ValidationError(t *testing.T) {
	r, _ := setupRouter(t)

	// Missing required 'name' field
	body := `{"code":"TEST"}`
	req := httptest.NewRequest("POST", "/api/v1/tenant/trainings/categories", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for validation error, got %d", w.Code)
	}
}

func TestHandler_UpdateCategory(t *testing.T) {
	r, h := setupRouter(t)

	// Create first
	svc := h.svc
	desc := "Original"
	created, _ := svc.CreateCategory(testCtx(), CreateTrainingCategoryRequest{
		Code: "OLD", Name: "Old", Description: &desc,
	})

	body := `{"name":"Updated Name"}`
	req := httptest.NewRequest("PUT", "/api/v1/tenant/trainings/categories/"+created.ID, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DeleteCategory(t *testing.T) {
	r, h := setupRouter(t)
	svc := h.svc
	desc := "To delete"
	created, _ := svc.CreateCategory(testCtx(), CreateTrainingCategoryRequest{
		Code: "DEL", Name: "Delete", Description: &desc,
	})

	req := httptest.NewRequest("DELETE", "/api/v1/tenant/trainings/categories/"+created.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestHandler_CreateSession(t *testing.T) {
	r, h := setupRouter(t)

	// Seed required data
	svc := h.svc
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	body := `{
		"course_id":"` + courseID + `",
		"session_code":"CLS-TEST",
		"trainer_name":"Test Trainer",
		"start_date":"2026-08-01",
		"end_date":"2026-08-05",
		"max_quota":25
	}`
	req := httptest.NewRequest("POST", "/api/v1/tenant/trainings/sessions", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListParticipants(t *testing.T) {
	r, h := setupRouter(t)

	// Seed data
	svc := h.svc
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessID := seedSession(t, svc, courseID)
	svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID: sessID, EmployeeID: "00000000-0000-0000-0000-000000000010",
	})

	req := httptest.NewRequest("GET", "/api/v1/tenant/trainings/participants?session_id="+sessID+"&page=1&per_page=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
