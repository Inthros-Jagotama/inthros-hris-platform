package documenttemplate

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func setupTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	db, cleanup := setupTestDB()
	t.Cleanup(cleanup)
	svc := NewService(newTestRepo(db), zap.NewNop())
	handler := NewHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	group := r.Group("/api/v1/tenant")
	RegisterRoutes(group, handler)
	return r
}

func TestHandlerCreateAndList(t *testing.T) {
	r := setupTestRouter(t)

	payload := `{"name":"PKWT","code":"PKWT_H1","document_type":"CONTRACT_AGREEMENT"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenant/document-templates", bytes.NewReader([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/tenant/document-templates", nil)
	listW := httptest.NewRecorder()
	r.ServeHTTP(listW, listReq)
	if listW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", listW.Code, listW.Body.String())
	}
	var body struct {
		Data struct {
			Data  []DocumentTemplate `json:"data"`
			Total int64              `json:"total"`
		} `json:"data"`
	}
	if err := json.Unmarshal(listW.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.Data.Total != 1 {
		t.Fatalf("expected 1 template, got %d", body.Data.Total)
	}
}

func TestHandlerCreateRejectsInvalidDocumentType(t *testing.T) {
	r := setupTestRouter(t)
	payload := `{"name":"X","code":"XCODE","document_type":"NOT_REAL"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenant/document-templates", bytes.NewReader([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandlerFromDefaultRouteNotShadowedByIDParam(t *testing.T) {
	r := setupTestRouter(t)
	payload := `{"document_type":"NOT_REAL","name":"X","code":"XCODE2"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenant/document-templates/from-default", bytes.NewReader([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// Must reach the from-default handler (which validates document_type and
	// 400s) rather than being swallowed by the GET /:id route as id="from-default".
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 (from-default handler reached, invalid doc type), got %d: %s", w.Code, w.Body.String())
	}
}
