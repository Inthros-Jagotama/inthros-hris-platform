package employee

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// setupSensitiveFieldHandlerRouter wires a Handler backed by a fresh
// in-memory sensitive-field-settings DB (see setupSensitiveFieldTestDB in
// sensitive_field_settings_test.go) onto a real gin router with the
// package's route registration, mirroring the pattern used in
// setting/numbering_handler_test.go.
// perms: permission claim yang dibawa caller (seperti yang di-set middleware
// AuthJWT). Kosong = pakai wildcard "*" supaya test perilaku handler tidak
// bercampur dengan test gating (lihat sensitive_field_settings_authz_test.go).
func setupSensitiveFieldHandlerRouter(t *testing.T, perms ...string) *gin.Engine {
	t.Helper()
	if len(perms) == 0 {
		perms = []string{"*"}
	}
	db := setupSensitiveFieldTestDB(t)
	resolver := func(ctx context.Context) (*gorm.DB, error) { return db, nil }
	repo := NewRepository(resolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	handler := NewHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("permissions", perms)
		c.Next()
	})
	group := r.Group("/api/v1/tenant")
	RegisterRoutes(group, handler)
	return r
}

func TestHandler_ListSensitiveFieldSettings(t *testing.T) {
	r := setupSensitiveFieldHandlerRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tenant/employees/settings/sensitive-fields", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var body struct {
		Data []SensitiveFieldSetting `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(body.Data) != len(SensitiveFieldRegistry) {
		t.Fatalf("expected %d settings, got %d", len(SensitiveFieldRegistry), len(body.Data))
	}
	for _, s := range body.Data {
		if s.FieldKey == "" {
			t.Errorf("expected non-empty field_key in %+v", s)
		}
	}
}

func TestHandler_SetSensitiveFieldEnabled(t *testing.T) {
	r := setupSensitiveFieldHandlerRouter(t)

	putReq := httptest.NewRequest(http.MethodPut, "/api/v1/tenant/employees/settings/sensitive-fields/employee.nik", strings.NewReader(`{"is_encryption_enabled":true}`))
	putReq.Header.Set("Content-Type", "application/json")
	putW := httptest.NewRecorder()
	r.ServeHTTP(putW, putReq)

	if putW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", putW.Code, putW.Body.String())
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/tenant/employees/settings/sensitive-fields", nil)
	getW := httptest.NewRecorder()
	r.ServeHTTP(getW, getReq)

	var body struct {
		Data []SensitiveFieldSetting `json:"data"`
	}
	if err := json.Unmarshal(getW.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	found := false
	for _, s := range body.Data {
		if s.FieldKey == "employee.nik" {
			found = true
			if !s.IsEncryptionEnabled {
				t.Errorf("expected employee.nik to be enabled after PUT")
			}
		}
	}
	if !found {
		t.Fatalf("employee.nik not found in settings list")
	}
}

func TestHandler_SetSensitiveFieldEnabled_UnknownKey(t *testing.T) {
	r := setupSensitiveFieldHandlerRouter(t)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/tenant/employees/settings/sensitive-fields/not.a.field", strings.NewReader(`{"is_encryption_enabled":true}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK || w.Code == http.StatusInternalServerError {
		t.Fatalf("expected a 4xx client error, got %d: %s", w.Code, w.Body.String())
	}
}
