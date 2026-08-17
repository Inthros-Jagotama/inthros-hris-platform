package training

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestRequireTrainingSettings verifies requireTrainingSettings rejects callers
// that only have the module-level "training.view" permission (which the default
// Employee role gets automatically) and only allows callers with the exact
// "training.settings.<action>" permission or "*".
func TestRequireTrainingSettings(t *testing.T) {
	gin.SetMode(gin.TestMode)

	newRouter := func(perms []string) *gin.Engine {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("permissions", perms)
			c.Next()
		})
		r.POST("/courses", requireTrainingSettings("create"), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
		return r
	}

	cases := []struct {
		name       string
		perms      []string
		wantStatus int
	}{
		{"module-level training.view saja ditolak", []string{"training.view"}, http.StatusForbidden},
		{"tanpa permission sama sekali ditolak", []string{}, http.StatusForbidden},
		{"permission submenu training.settings.create diizinkan", []string{"training.settings.create"}, http.StatusOK},
		{"wildcard * diizinkan", []string{"*"}, http.StatusOK},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			router := newRouter(tc.perms)
			req := httptest.NewRequest(http.MethodPost, "/courses", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tc.wantStatus, rec.Body.String())
			}
		})
	}
}
