package employee

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// =============================================================================
// Regression: endpoint setelan field sensitif tidak boleh terbuka hanya karena
// caller punya employee.view / employee.update.
//
// Middleware RBAC global menurunkan resource dari path, sehingga
// /employees/settings/sensitive-fields tercek sebagai employee.view /
// employee.update — artinya siapa pun yang boleh melihat atau mengedit data
// karyawan bisa mengubah setelan enkripsi at-rest seluruh tenant.
// Route sekarang butuh setting.sensitive-fields.view / .manage
// (fallback module-level setting.view / setting.update).
// =============================================================================

const sensitiveFieldsPath = "/api/v1/tenant/employees/settings/sensitive-fields"

func doSensitiveFieldsRequest(t *testing.T, method, path string, body string, perms ...string) int {
	t.Helper()
	r := setupSensitiveFieldHandlerRouter(t, perms...)
	var req *http.Request
	if body == "" {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Code
}

func TestSensitiveFieldSettings_ListRequiresDedicatedPermission(t *testing.T) {
	cases := []struct {
		name     string
		perms    []string
		wantCode int
	}{
		{"employee.view saja ditolak", []string{"employee.view", "employee.update"}, http.StatusForbidden},
		{"tanpa permission ditolak", []string{}, http.StatusForbidden},
		{"permission khusus diizinkan", []string{"setting.sensitive-fields.view"}, http.StatusOK},
		{"fallback module-level setting.view diizinkan", []string{"setting.view"}, http.StatusOK},
		{"wildcard super admin diizinkan", []string{"*"}, http.StatusOK},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// perms kosong harus benar-benar kosong, bukan default wildcard.
			perms := tc.perms
			if len(perms) == 0 {
				perms = []string{"none.none"}
			}
			got := doSensitiveFieldsRequest(t, http.MethodGet, sensitiveFieldsPath, "", perms...)
			if got != tc.wantCode {
				t.Errorf("GET %s with %v = %d, want %d", sensitiveFieldsPath, perms, got, tc.wantCode)
			}
		})
	}
}

func TestSensitiveFieldSettings_ToggleRequiresManagePermission(t *testing.T) {
	path := sensitiveFieldsPath + "/employee.nik"
	const body = `{"is_encryption_enabled":true}`

	cases := []struct {
		name     string
		perms    []string
		wantCode int
	}{
		{"employee.update saja ditolak", []string{"employee.view", "employee.update"}, http.StatusForbidden},
		{"permission view saja ditolak", []string{"setting.sensitive-fields.view"}, http.StatusForbidden},
		{"permission manage diizinkan", []string{"setting.sensitive-fields.manage"}, http.StatusOK},
		{"fallback module-level setting.update diizinkan", []string{"setting.update"}, http.StatusOK},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := doSensitiveFieldsRequest(t, http.MethodPut, path, body, tc.perms...)
			if got != tc.wantCode {
				t.Errorf("PUT %s with %v = %d, want %d", path, tc.perms, got, tc.wantCode)
			}
		})
	}
}
