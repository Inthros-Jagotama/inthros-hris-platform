package company

import (
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

// testEnv menyimpan lingkungan test untuk handler tests.
type testEnv struct {
	router  *gin.Engine
	db      *gorm.DB
	fakeTM  *FakeTenantManager
	cleanup func()
}

// setupTestEnv creates a complete test environment for handler tests.
func setupTestEnv() *testEnv {
	gin.SetMode(gin.TestMode)

	db, cleanup := setupTestDB()
	repo := NewRepository(db)
	logger, _ := zap.NewDevelopment()
	fakeTM := &FakeTenantManager{}
	svc := NewService(repo, fakeTM, nil, nil, logger)
	handler := NewHandler(svc)

	r := gin.New()
	rg := r.Group("/api/v1/platform")
	{
		companies := rg.Group("/companies")
			{
				companies.POST("/:id/terminate", handler.Terminate)
				companies.POST("/:id/rotate-credentials", handler.RotateCredentials)
				companies.POST("/", handler.Create)
				companies.GET("/:id", handler.GetByID)
			}
		}

	return &testEnv{
		router:  r,
		db:      db,
		fakeTM:  fakeTM,
		cleanup: func() { cleanup(); logger.Sync() },
	}
}

// =========================================================================
// Unit Tests — isCompanyEditor (authz gate untuk PUT /tenant/companies/me)
// =========================================================================

func TestIsCompanyEditor(t *testing.T) {
	tests := []struct {
		name string
		role string
		want bool
	}{
		{"platform super_admin", "super_admin", true},
		{"platform company_admin", "company_admin", true},
		{"tenant admin (capitalized)", "Admin", true},
		{"tenant admin lowercase", "admin", true},
		{"whitespace + mixed case", "  Company_Admin ", true},
		{"tenant employee", "Employee", false},
		{"platform employee", "employee", false},
		{"manager", "manager", false},
		{"unknown role", "ceo", false},
		{"empty role (fail-closed)", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCompanyEditor(tt.role); got != tt.want {
				t.Errorf("isCompanyEditor(%q) = %v, want %v", tt.role, got, tt.want)
			}
		})
	}
}

// =========================================================================
// Handler Tests — Rotate Credentials Endpoint
// =========================================================================

func TestHandler_RotateCredentials_Success(t *testing.T) {
	env := setupTestEnv()
	defer env.cleanup()

	company := createTestCompany(env.db, "Rotate Handler")

	var rotatedPassword string
	env.fakeTM.RotateTenantCredFunc = func(companyID, newPassword string) error {
		rotatedPassword = newPassword
		return nil
	}

	body := `{"new_password":"handler-pass-123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/platform/companies/"+company.ID.String()+"/rotate-credentials",
		strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	if rotatedPassword != "handler-pass-123" {
		t.Errorf("expected rotated password 'handler-pass-123', got %q", rotatedPassword)
	}

	var resp struct {
		Success bool `json:"success"`
		Data    struct {
			CompanyID string `json:"company_id"`
			Rotated   bool   `json:"rotated"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if !resp.Success || !resp.Data.Rotated {
		t.Errorf("expected success+rotated, got %+v", resp)
	}
}

func TestHandler_RotateCredentials_InvalidBody(t *testing.T) {
	env := setupTestEnv()
	defer env.cleanup()

	company := createTestCompany(env.db, "Rotate Bad Body")

	// new_password terlalu pendek (< 8) → validasi gagal
	body := `{"new_password":"short"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/platform/companies/"+company.ID.String()+"/rotate-credentials",
		strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for short password, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_RotateCredentials_Terminated(t *testing.T) {
	env := setupTestEnv()
	defer env.cleanup()

	company := createTestCompany(env.db, "Rotate Term Handler")
	company.Status = CompanyStatusTerminated
	if err := env.db.Save(company).Error; err != nil {
		t.Fatalf("failed to save company: %v", err)
	}

	body := `{"new_password":"handler-pass-123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/platform/companies/"+company.ID.String()+"/rotate-credentials",
		strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409 for terminated company, got %d: %s", w.Code, w.Body.String())
	}
}

// =========================================================================
// Handler Tests — Terminate Endpoint
// =========================================================================

func TestHandler_Terminate_Success(t *testing.T) {
	env := setupTestEnv()
	defer env.cleanup()

	// Create a company directly in DB
	company := createTestCompany(env.db, "To Be Terminated")

	env.fakeTM.DropTenantDBFunc = func(companyID string) error {
		return nil
	}
	env.fakeTM.RemoveTenantConnFunc = func(companyID string) error {
		return nil
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/platform/companies/"+company.ID.String()+"/terminate", nil)
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool            `json:"success"`
		Data    CompanyResponse `json:"data"`
		Message string          `json:"message"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Data.Status != "terminated" {
		t.Errorf("expected status 'terminated', got '%s'", resp.Data.Status)
	}
	if resp.Message == "" {
		t.Error("expected a message")
	}
}

func TestHandler_Terminate_NotFound(t *testing.T) {
	env := setupTestEnv()
	defer env.cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/platform/companies/"+uuid.New().String()+"/terminate", nil)
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409 Conflict, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Terminate_AlreadyTerminated(t *testing.T) {
	env := setupTestEnv()
	defer env.cleanup()

	company := createTestCompany(env.db, "Already Terminated")
	company.Status = CompanyStatusTerminated
	if err := env.db.Save(company).Error; err != nil {
		t.Fatalf("failed to update company status: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/platform/companies/"+company.ID.String()+"/terminate", nil)
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409 Conflict, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Terminate_InvalidUUID(t *testing.T) {
	env := setupTestEnv()
	defer env.cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/platform/companies/not-a-uuid/terminate", nil)
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409 Conflict, got %d: %s", w.Code, w.Body.String())
	}
}
