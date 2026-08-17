package approval

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestRequireApprovalSettings verifies requireApprovalSettings rejects
// callers that only have the module-level "approval.view" permission (which
// the default Employee role gets automatically) and only allows callers
// with the exact "approval.settings.<action>" permission or "*".
func TestRequireApprovalSettings(t *testing.T) {
	gin.SetMode(gin.TestMode)

	newRouter := func(perms []string) *gin.Engine {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("permissions", perms)
			c.Next()
		})
		r.PUT("/flows/:flowId", requireApprovalSettings("update"), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
		return r
	}

	cases := []struct {
		name       string
		perms      []string
		wantStatus int
	}{
		{"module-level approval.view saja ditolak", []string{"approval.view"}, http.StatusForbidden},
		{"tanpa permission sama sekali ditolak", []string{}, http.StatusForbidden},
		{"permission submenu approval.settings.update diizinkan", []string{"approval.settings.update"}, http.StatusOK},
		{"wildcard * diizinkan", []string{"*"}, http.StatusOK},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			router := newRouter(tc.perms)
			req := httptest.NewRequest(http.MethodPut, "/flows/some-id", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tc.wantStatus, rec.Body.String())
			}
		})
	}
}
