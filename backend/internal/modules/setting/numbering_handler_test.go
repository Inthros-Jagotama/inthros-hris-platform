package setting

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

	sqlite "github.com/glebarez/sqlite"

	"github.com/inthros/hris-platform/internal/pkg/numbering"
)

func setupNumberingTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	// Use a private (non-shared-cache) in-memory database so this test's
	// tables don't leak into the shared-cache DSN used by other tests in
	// this package (helpers_test.go's setupTestDB uses
	// "file::memory:?cache=shared", which persists for the whole test
	// binary as long as any connection to it stays open).
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() {
		if sqlDB, dbErr := db.DB(); dbErr == nil {
			_ = sqlDB.Close()
		}
	})
	if err := db.AutoMigrate(&numbering.DocumentNumberingSetting{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	seed := numbering.DocumentNumberingSetting{
		ID: "11111111-1111-1111-1111-111111111111", DocumentType: numbering.DocumentTypeEmployeeMovement,
		FormatTemplate: "SK/{sequence:3}/{year}", ResetPeriod: numbering.ResetPeriodYearly,
	}
	if err := db.Create(&seed).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	resolver := func(ctx context.Context) (*gorm.DB, error) { return db, nil }
	numberingSvc := numbering.NewService(resolver, zap.NewNop())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewNumberingHandler(numberingSvc)
	group := r.Group("/api/v1/tenant/settings")
	RegisterNumberingRoutes(group, handler)
	return r
}

func TestListNumberingSettings(t *testing.T) {
	r := setupNumberingTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tenant/settings/document-numbering", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var body struct {
		Data []numbering.DocumentNumberingSetting `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(body.Data) != 1 {
		t.Fatalf("expected 1 setting, got %d", len(body.Data))
	}
}

func TestUpdateNumberingSettingRejectsBadDocumentType(t *testing.T) {
	r := setupNumberingTestRouter(t)
	payload := `{"format_template":"X/{sequence}","reset_period":"yearly"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/tenant/settings/document-numbering/not_a_type", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}
