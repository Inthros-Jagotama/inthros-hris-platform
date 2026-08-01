package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// fakeHostResolver adalah stub HostCompanyResolver untuk test.
type fakeHostResolver struct {
	fn func(host string) (string, error)
}

func (f fakeHostResolver) ResolveByHost(host string) (string, error) {
	if f.fn == nil {
		return "", errors.New("no resolver stub")
	}
	return f.fn(host)
}

func newTestRouterWithResolver(resolver HostCompanyResolver) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mw := TenantResolver(resolver, zap.NewNop())
	r.Use(mw)
	r.GET("/api/v1/tenant/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"company_id": c.GetString("company_id"),
			"tenant_id":  c.GetHeader("X-Tenant-ID"),
		})
	})
	return r
}

func TestTenantResolver_NoHostNoResolver_Continues(t *testing.T) {
	r := newTestRouterWithResolver(nil)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/tenant/ping", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	// Tanpa JWT company_id & tanpa host yang bisa di-resolve → company_id kosong.
	if got := w.Header().Get("X-Tenant-ID"); got != "" {
		t.Errorf("expected no X-Tenant-ID header, got %q", got)
	}
}

func TestTenantResolver_ResolvesFromHost(t *testing.T) {
	r := newTestRouterWithResolver(fakeHostResolver{
		fn: func(host string) (string, error) {
			if host == "pt-inthros-jago-utama.localhost" {
				return "11111111-1111-1111-1111-111111111111", nil
			}
			return "", errors.New("not found")
		},
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/tenant/ping", nil)
	req.Host = "pt-inthros-jago-utama.localhost"
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if got := w.Header().Get("X-Tenant-ID"); got != "11111111-1111-1111-1111-111111111111" {
		t.Errorf("expected X-Tenant-ID from host resolve, got %q", got)
	}
}

func TestTenantResolver_XTenantIDHeader_TakesPrecedenceOverHost(t *testing.T) {
	r := newTestRouterWithResolver(fakeHostResolver{
		fn: func(host string) (string, error) {
			return "99999999-9999-9999-9999-999999999999", nil // harus TIDAK terpakai
		},
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/tenant/ping", nil)
	req.Host = "whatever.localhost"
	req.Header.Set("X-Tenant-ID", "55555555-5555-5555-5555-555555555555")
	r.ServeHTTP(w, req)

	if got := w.Header().Get("X-Tenant-ID"); got != "55555555-5555-5555-5555-555555555555" {
		t.Errorf("expected X-Tenant-ID from header, got %q", got)
	}
}

func TestTenantResolver_UnresolvableHost_Continues(t *testing.T) {
	r := newTestRouterWithResolver(fakeHostResolver{
		fn: func(host string) (string, error) {
			return "", errors.New("company not found for host")
		},
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/tenant/ping", nil)
	req.Host = "unknown.localhost"
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 (continue), got %d", w.Code)
	}
	if got := w.Header().Get("X-Tenant-ID"); got != "" {
		t.Errorf("expected no X-Tenant-ID for unresolvable host, got %q", got)
	}
}

func TestTenantResolver_XForwardedHost_TakesPrecedenceOverHost(t *testing.T) {
	resolved := []string{}
	r := newTestRouterWithResolver(fakeHostResolver{
		fn: func(host string) (string, error) {
			resolved = append(resolved, host)
			return "11111111-1111-1111-1111-111111111111", nil
		},
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/tenant/ping", nil)
	req.Host = "backend.localhost"
	req.Header.Set("X-Forwarded-Host", "pt-inthros-jago-utama.localhost")
	r.ServeHTTP(w, req)

	if got := w.Header().Get("X-Tenant-ID"); got != "11111111-1111-1111-1111-111111111111" {
		t.Errorf("expected X-Tenant-ID from X-Forwarded-Host resolve, got %q", got)
	}
	if len(resolved) != 1 || resolved[0] != "pt-inthros-jago-utama.localhost" {
		t.Errorf("expected X-Forwarded-Host to be used (not Host), got %v", resolved)
	}
}

func TestTenantResolver_InvalidXTenantIDHeader_Ignored(t *testing.T) {
	// Header X-Tenant-ID bukan UUID → harus diabaikan, fallback ke Host.
	r := newTestRouterWithResolver(fakeHostResolver{
		fn: func(host string) (string, error) {
			return "22222222-2222-2222-2222-222222222222", nil
		},
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/tenant/ping", nil)
	req.Host = "pt-inthros-jago-utama.localhost"
	req.Header.Set("X-Tenant-ID", "not-a-uuid")
	r.ServeHTTP(w, req)

	if got := w.Header().Get("X-Tenant-ID"); got != "22222222-2222-2222-2222-222222222222" {
		t.Errorf("expected fallback to host resolve after invalid header, got %q", got)
	}
}

func TestTenantResolver_JWTCompanyID_WinsOverHost(t *testing.T) {
	// Simulasi: AuthJWT sudah set company_id dari JWT claims.
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("company_id", "aaaa0000-aaaa-0000-aaaa-000000000000")
		c.Next()
	})
	r.Use(TenantResolver(fakeHostResolver{
		fn: func(host string) (string, error) {
			// Host berbeda — harus diabaikan karena JWT menang.
			return "99999999-9999-9999-9999-999999999999", nil
		},
	}, zap.NewNop()))
	r.GET("/api/v1/tenant/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"company_id": c.GetString("company_id"),
			"tenant_id":  c.GetHeader("X-Tenant-ID"),
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/tenant/ping", nil)
	req.Host = "attacker.localhost"
	r.ServeHTTP(w, req)

	if got := w.Header().Get("X-Tenant-ID"); got != "aaaa0000-aaaa-0000-aaaa-000000000000" {
		t.Errorf("expected JWT company_id to win over host, got %q", got)
	}
}
