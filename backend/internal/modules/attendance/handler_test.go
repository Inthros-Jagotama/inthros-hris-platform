package attendance

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
	"gorm.io/gorm"
)

// setupTestRouter creates a Gin engine with attendance routes registered.
func setupTestRouter() (*gin.Engine, *Repository, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)

	db, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	handler := NewHandler(svc)

	r := gin.New()
	rg := r.Group("/api/v1/tenant")
	RegisterRoutes(rg, handler)

	return r, repo, db, func() {
		cleanup()
		_ = logger.Sync()
	}
}

// =========================================================================
// Company Settings Handler Tests
// =========================================================================

func TestHandler_UpsertCompanySetting_Success(t *testing.T) {
	router, _, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"is_location_required": true, "max_distance_meter": 100}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/tenant/attendance/settings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool                   `json:"success"`
		Data    CompanySettingResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success=true")
	}
	if !resp.Data.IsLocationRequired {
		t.Error("expected IsLocationRequired = true")
	}
}

func TestHandler_GetCompanySetting_Success(t *testing.T) {
	router, repo, _, cleanup := setupTestRouter()
	defer cleanup()

	// First upsert a setting
	s := &AttendanceCompanySetting{IsLocationRequired: true}
	_ = repo.UpsertCompanySetting(context.Background(), s)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/attendance/settings", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetCompanySetting_NoRow_ReturnsDefaults(t *testing.T) {
	router, _, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/attendance/settings", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK with defaults, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Shift Handler Tests
// =========================================================================

func TestHandler_CreateShift_Success(t *testing.T) {
	router, _, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"shift_name": "Morning Shift", "check_in_time": "08:00:00", "check_out_time": "17:00:00"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/attendance/shifts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool                 `json:"success"`
		Data    CompanyShiftResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Data.ShiftName != "Morning Shift" {
		t.Errorf("expected shift_name 'Morning Shift', got '%s'", resp.Data.ShiftName)
	}
}

func TestHandler_CreateShift_ValidationError(t *testing.T) {
	router, _, _, cleanup := setupTestRouter()
	defer cleanup()

	// Missing required fields
	body := `{"shift_name": "Morning Shift"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/attendance/shifts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetShiftByID_Success(t *testing.T) {
	router, repo, _, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestShift(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/attendance/shifts/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetShiftByID_NotFound(t *testing.T) {
	router, _, _, cleanup := setupTestRouter()
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/attendance/shifts/"+uuidStr(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListShifts_Success(t *testing.T) {
	router, repo, _, cleanup := setupTestRouter()
	defer cleanup()

	createTestShift(repo)
	createTestShift(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/attendance/shifts", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdateShift_Success(t *testing.T) {
	router, repo, _, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestShift(repo)

	body := `{"shift_name": "Night Shift"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/tenant/attendance/shifts/"+created.ID.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DeleteShift_Success(t *testing.T) {
	router, repo, _, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestShift(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/tenant/attendance/shifts/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Location Handler Tests
// =========================================================================

func TestHandler_CreateLocation_Success(t *testing.T) {
	router, _, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"name": "Main Office", "latitude": -6.2088, "longitude": 106.8456}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/attendance/locations", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetLocationByID_Success(t *testing.T) {
	router, repo, _, cleanup := setupTestRouter()
	defer cleanup()

	created := createTestLocation(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/attendance/locations/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListLocations_Success(t *testing.T) {
	router, repo, _, cleanup := setupTestRouter()
	defer cleanup()

	createTestLocation(repo)
	createTestLocation(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/attendance/locations", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Event Handler Tests
// =========================================================================

func TestHandler_CreateEvent_Success(t *testing.T) {
	router, _, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{
		"employee_id": "` + uuidStr() + `",
		"event_type": "CHECKIN",
		"event_time_utc": "2026-01-15T00:00:00Z",
		"event_time_local": "2026-01-15T07:00:00+07:00",
		"latitude": -6.2088,
		"longitude": 106.8456
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/attendance/events", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateEvent_ValidationError(t *testing.T) {
	router, _, _, cleanup := setupTestRouter()
	defer cleanup()

	// Missing required fields
	body := `{"event_type": "CHECKIN"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/attendance/events", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetEventByID_Success(t *testing.T) {
	router, repo, _, cleanup := setupTestRouter()
	defer cleanup()

	empID := uuid.New()
	created := createTestEvent(repo, empID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/attendance/events/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListEvents_Success(t *testing.T) {
	router, repo, _, cleanup := setupTestRouter()
	defer cleanup()

	empID := uuid.New()
	createTestEvent(repo, empID)
	createTestEvent(repo, empID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/attendance/events", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Session Handler Tests
// =========================================================================

func TestHandler_ListSessions_Success(t *testing.T) {
	router, _, db, cleanup := setupTestRouter()
	defer cleanup()

	empID := uuid.New()
	createTestSession(db, empID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/attendance/sessions", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Overtime Request Handler Tests
// =========================================================================

func TestHandler_CreateOvertimeRequest_Success(t *testing.T) {
	router, _, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{
		"employee_id": "` + uuidStr() + `",
		"work_date": "2026-01-15",
		"start_time_local": "2026-01-15T18:00:00+07:00",
		"end_time_local": "2026-01-15T20:00:00+07:00",
		"requested_minutes": 120
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/attendance/overtime-requests", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetOvertimeRequestByID_Success(t *testing.T) {
	router, repo, _, cleanup := setupTestRouter()
	defer cleanup()

	empID := uuid.New()
	created := createTestOvertimeRequest(repo, empID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/attendance/overtime-requests/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListOvertimeRequests_Success(t *testing.T) {
	router, repo, _, cleanup := setupTestRouter()
	defer cleanup()

	empID := uuid.New()
	createTestOvertimeRequest(repo, empID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/attendance/overtime-requests", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Exempt Position Handler Tests
// =========================================================================

func TestHandler_CreateExemptPosition_Success(t *testing.T) {
	router, _, _, cleanup := setupTestRouter()
	defer cleanup()

	body := `{"organization_id": "` + uuidStr() + `"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/attendance/exempt-positions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetExemptPositionByID_Success(t *testing.T) {
	router, repo, _, cleanup := setupTestRouter()
	defer cleanup()

	orgID := uuid.New()
	created := createTestExemptPosition(repo, orgID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/attendance/exempt-positions/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListExemptPositions_Success(t *testing.T) {
	router, repo, _, cleanup := setupTestRouter()
	defer cleanup()

	createTestExemptPosition(repo, uuid.New())
	createTestExemptPosition(repo, uuid.New())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tenant/attendance/exempt-positions", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdateExemptPosition_Success(t *testing.T) {
	router, repo, _, cleanup := setupTestRouter()
	defer cleanup()

	orgID := uuid.New()
	created := createTestExemptPosition(repo, orgID)

	body := `{"is_exempt": false}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/tenant/attendance/exempt-positions/"+created.ID.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DeleteExemptPosition_Success(t *testing.T) {
	router, repo, _, cleanup := setupTestRouter()
	defer cleanup()

	orgID := uuid.New()
	created := createTestExemptPosition(repo, orgID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/tenant/attendance/exempt-positions/"+created.ID.String(), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}
}
