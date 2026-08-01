package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/cache"
)

// errFakeUnavailable mensimulasikan sumber data tidak tersedia.
var errFakeUnavailable = &fakeUnavailableError{}

type fakeUnavailableError struct{}

func (e *fakeUnavailableError) Error() string { return "fake source unavailable" }

// fakeLister adalah implementasi CompanyModuleLister untuk test.
type fakeLister struct {
	slugs []string
	err   error
	calls int
}

func (f *fakeLister) EnabledModuleSlugs(companyID string) ([]string, error) {
	f.calls++
	return f.slugs, f.err
}

// newTestCache membuat cache dengan miniredis (dipakai sesuai konvensi cache_test.go).
func newTestCache(t *testing.T) *cache.Cache {
	t.Helper()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis.Run: %v", err)
	}
	t.Cleanup(mr.Close)

	c, err := cache.New(cache.Config{RedisAddr: mr.Addr(), DefaultTTL: time.Minute}, zap.NewNop())
	if err != nil {
		t.Fatalf("cache.New: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })
	return c
}

func TestModuleSlugFromPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"organizations list", "/api/v1/tenant/organizations", "organization"},
		{"organizations detail", "/api/v1/tenant/organizations/123", "organization"},
		{"organization summaries", "/api/v1/tenant/organization-summaries", "organization"},
		{"employees", "/api/v1/tenant/employees", "employee"},
		{"job management", "/api/v1/tenant/job-management/titles", "jobmanagement"},
		{"competency", "/api/v1/tenant/competency", "competency"},
		{"employee movements", "/api/v1/tenant/employee-movements", "employeemovement"},
		{"attendance", "/api/v1/tenant/attendance", "attendance"},
		{"approval", "/api/v1/tenant/approval", "approval"},
		{"payroll", "/api/v1/tenant/payroll/runs", "payroll"},
		{"leave", "/api/v1/tenant/leave", "leave"},
		{"performance", "/api/v1/tenant/performance", "performance"},
		{"recruitment", "/api/v1/tenant/recruitment", "recruitment"},
		{"reimbursements", "/api/v1/tenant/reimbursements", "reimbursement"},
		{"trainings", "/api/v1/tenant/trainings", "training"},
		{"workforce intelligence", "/api/v1/tenant/workforce-intelligence", "workforce-intelligence"},
		{"career intelligence", "/api/v1/tenant/career-intelligence", "career-intelligence"},
		{"settings", "/api/v1/tenant/settings/banks", "setting"},
		{"company modules (exempt)", "/api/v1/tenant/company-modules", ""},
		{"packages (exempt)", "/api/v1/tenant/packages", ""},
		{"platform route", "/api/v1/platform/companies", ""},
		{"non-tenant route", "/docs", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := moduleSlugFromPath(tt.path); got != tt.want {
				t.Errorf("moduleSlugFromPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

// newTestContext membuat gin context test dengan company_id dan path.
func newTestContext(method, path, companyID string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, nil)
	if companyID != "" {
		c.Set("company_id", companyID)
	}
	return c, w
}

func TestLicenseMiddleware_BlockWhenModuleNotLicensed(t *testing.T) {
	cm := newTestCache(t)
	lister := &fakeLister{slugs: []string{"organization", "employee"}} // tanpa "payroll"
	mw := LicenseMiddleware(cm, lister, zap.NewNop())

	c, w := newTestContext(http.MethodGet, "/api/v1/tenant/payroll/runs", "company-1")
	mw(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.Error.Code != "MODULE_NOT_LICENSED" {
		t.Errorf("error code = %q, want MODULE_NOT_LICENSED", body.Error.Code)
	}
}

func TestLicenseMiddleware_AllowWhenModuleLicensed(t *testing.T) {
	cm := newTestCache(t)
	lister := &fakeLister{slugs: []string{"organization", "employee"}}
	mw := LicenseMiddleware(cm, lister, zap.NewNop())

	c, w := newTestContext(http.MethodGet, "/api/v1/tenant/employees", "company-1")
	mw(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.Len() != 0 {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}

func TestLicenseMiddleware_AllowExemptPaths(t *testing.T) {
	cm := newTestCache(t)
	lister := &fakeLister{slugs: []string{}}
	mw := LicenseMiddleware(cm, lister, zap.NewNop())

	// Path yang tidak memetakan ke modul berlisensi tidak pernah diblokir,
	// meskipun daftar modul company kosong.
	c, w := newTestContext(http.MethodGet, "/api/v1/tenant/company-modules", "company-1")
	mw(c)
	if w.Code != http.StatusOK {
		t.Errorf("company-modules status = %d, want %d", w.Code, http.StatusOK)
	}

	c2, w2 := newTestContext(http.MethodGet, "/api/v1/tenant/packages", "company-1")
	mw(c2)
	if w2.Code != http.StatusOK {
		t.Errorf("packages status = %d, want %d", w2.Code, http.StatusOK)
	}
}

func TestLicenseMiddleware_MissingCompanyID_PassesThrough(t *testing.T) {
	cm := newTestCache(t)
	lister := &fakeLister{slugs: []string{}}
	mw := LicenseMiddleware(cm, lister, zap.NewNop())

	c, w := newTestContext(http.MethodGet, "/api/v1/tenant/employees", "")
	mw(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if lister.calls != 0 {
		t.Errorf("lister called %d times, want 0", lister.calls)
	}
}

func TestLicenseMiddleware_CacheMissThenHit(t *testing.T) {
	cm := newTestCache(t)
	lister := &fakeLister{slugs: []string{"employee"}}
	mw := LicenseMiddleware(cm, lister, zap.NewNop())

	// Request 1: cache miss → lister dipanggil → hasil di-cache.
	c1, _ := newTestContext(http.MethodGet, "/api/v1/tenant/employees", "company-cache")
	mw(c1)
	if lister.calls != 1 {
		t.Fatalf("after request 1, lister calls = %d, want 1", lister.calls)
	}

	// Request 2: cache hit → lister TIDAK dipanggil lagi.
	c2, _ := newTestContext(http.MethodGet, "/api/v1/tenant/employees", "company-cache")
	mw(c2)
	if lister.calls != 1 {
		t.Errorf("after request 2, lister calls = %d, want 1 (cache hit)", lister.calls)
	}

	// Company lain → cache miss baru.
	c3, _ := newTestContext(http.MethodGet, "/api/v1/tenant/employees", "company-lain")
	mw(c3)
	if lister.calls != 2 {
		t.Errorf("after request 3, lister calls = %d, want 2", lister.calls)
	}
}

func TestLicenseMiddleware_EmptyModuleList_FailOpen(t *testing.T) {
	cm := newTestCache(t)
	lister := &fakeLister{slugs: []string{}} // company tanpa data lisensi sama sekali
	mw := LicenseMiddleware(cm, lister, zap.NewNop())

	// Modul path tetap diizinkan (fail-open) karena company belum punya data lisensi,
	// supaya tenant legacy/tanpa package tidak terkunci total.
	c, w := newTestContext(http.MethodGet, "/api/v1/tenant/employees", "company-no-license")
	mw(c)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d (fail-open on empty module list)", w.Code, http.StatusOK)
	}
}

func TestLicenseMiddleware_ListerError_FailOpen(t *testing.T) {
	cm := newTestCache(t)
	lister := &fakeLister{err: errFakeUnavailable}
	mw := LicenseMiddleware(cm, lister, zap.NewNop())

	// Error dari sumber data → fail-open (allow) dengan log warning.
	c, w := newTestContext(http.MethodGet, "/api/v1/tenant/employees", "company-1")
	mw(c)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d (fail-open)", w.Code, http.StatusOK)
	}
}

func TestLicenseMiddleware_CacheKey(t *testing.T) {
	if got := LicenseCacheKey("abc-123"); got != "hris:platform:license:modules:abc-123" {
		t.Errorf("LicenseCacheKey = %q, want hris:platform:license:modules:abc-123", got)
	}
}
