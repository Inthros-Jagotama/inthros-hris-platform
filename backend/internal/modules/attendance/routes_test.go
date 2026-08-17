package attendance

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestRequireAttendanceSettings verifies requireAttendanceSettings rejects
// callers that only have the module-level "attendance.view" permission (which
// the default Employee role gets automatically) and only allows callers with
// the exact "attendance.settings.<action>" permission or "*".
func TestRequireAttendanceSettings(t *testing.T) {
	gin.SetMode(gin.TestMode)

	newRouter := func(perms []string) *gin.Engine {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("permissions", perms)
			c.Next()
		})
		r.PUT("/settings/:id", requireAttendanceSettings("update"), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
		return r
	}

	cases := []struct {
		name       string
		perms      []string
		wantStatus int
	}{
		{"module-level attendance.view saja ditolak", []string{"attendance.view"}, http.StatusForbidden},
		{"tanpa permission sama sekali ditolak", []string{}, http.StatusForbidden},
		{"permission submenu attendance.settings.update diizinkan", []string{"attendance.settings.update"}, http.StatusOK},
		{"wildcard * diizinkan", []string{"*"}, http.StatusOK},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			router := newRouter(tc.perms)
			req := httptest.NewRequest(http.MethodPut, "/settings/some-id", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tc.wantStatus, rec.Body.String())
			}
		})
	}
}

// TestRequireAttendanceReport verifies requireAttendanceReport rejects callers
// that only have the module-level "attendance.view" permission and only allows
// the exact "attendance.report.view" permission or "*".
func TestRequireAttendanceReport(t *testing.T) {
	gin.SetMode(gin.TestMode)

	newRouter := func(perms []string) *gin.Engine {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("permissions", perms)
			c.Next()
		})
		r.GET("/reports/sessions", requireAttendanceReport("view"), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
		return r
	}

	cases := []struct {
		name       string
		perms      []string
		wantStatus int
	}{
		{"module-level attendance.view saja ditolak", []string{"attendance.view"}, http.StatusForbidden},
		{"tanpa permission sama sekali ditolak", []string{}, http.StatusForbidden},
		{"permission submenu attendance.report.view diizinkan", []string{"attendance.report.view"}, http.StatusOK},
		{"wildcard * diizinkan", []string{"*"}, http.StatusOK},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			router := newRouter(tc.perms)
			req := httptest.NewRequest(http.MethodGet, "/reports/sessions", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tc.wantStatus, rec.Body.String())
			}
		})
	}
}
