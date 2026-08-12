package approval

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestEmitRoutingError_Bilingual verifies that a RoutingError is emitted as a
// 400 with the bilingual message following the request language and the
// APPROVAL_ROUTING_FAILED error code.
func TestEmitRoutingError_Bilingual(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/x", nil)
	c.Request.Header.Set("Accept-Language", "id")

	handled := EmitRoutingError(c, &RoutingError{
		Key:    "approval.no_supervisor_vacant",
		Params: []string{"Persetujuan Supervisor"},
	})
	if !handled {
		t.Fatal("expected EmitRoutingError to handle a *RoutingError")
	}
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	errObj, ok := body["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing error object: %+v", body)
	}
	if errObj["code"] != "APPROVAL_ROUTING_FAILED" {
		t.Errorf("expected code APPROVAL_ROUTING_FAILED, got %v", errObj["code"])
	}
	msg, _ := errObj["message"].(string)
	if !strings.Contains(msg, "Supervisor tidak ditemukan untuk langkah \"Persetujuan Supervisor\"") {
		t.Errorf("expected Indonesian supervisor message, got: %s", msg)
	}
}

// TestEmitRoutingError_EnglishDefault verifies the English catalog text is
// used when no Accept-Language header is present.
func TestEmitRoutingError_EnglishDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/x", nil)

	if !EmitRoutingError(c, &RoutingError{Key: "approval.flow_inactive"}) {
		t.Fatal("expected EmitRoutingError to handle a *RoutingError")
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	msg, _ := body["error"].(map[string]interface{})["message"].(string)
	if msg != "Flow is not active" {
		t.Errorf("expected English flow message %q, got: %s", "Flow is not active", msg)
	}
}

// TestEmitRoutingError_NonRoutingReturnsFalse ensures non-routing errors are
// left untouched (return false, no response written) so callers fall back to
// their own handling.
func TestEmitRoutingError_NonRoutingReturnsFalse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	if EmitRoutingError(c, fmt.Errorf("flow not found")) {
		t.Fatal("expected false for a non-RoutingError")
	}
	if w.Body.Len() != 0 {
		t.Fatalf("expected no response written, got: %s", w.Body.String())
	}
}

// TestEmitRoutingError_WrappedRoutingError verifies errors.As still classifies
// a RoutingError wrapped with %w (as consumer services do).
func TestEmitRoutingError_WrappedRoutingError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/x", nil)

	wrapped := fmt.Errorf("failed to create approval instance: %w",
		&RoutingError{Key: "approval.flow_no_steps"})

	if !EmitRoutingError(c, wrapped) {
		t.Fatal("expected wrapped RoutingError to be handled")
	}
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
