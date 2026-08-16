package authz

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// =========================================================================
// NewMiddleware — JWT permissions claim-first behavior
// =========================================================================
//
// Middleware memeriksa permissions claim JWT (RBAC tenant, format
// "resource.action") terlebih dahulu, lalu fallback ke platform enforcer.
// Ini penting karena user tenant (employee) membawa role "Employee"/"Admin"
// yang tidak dikenal platform enforcer (hanya slug lowercase).

func setupMiddlewareTest(claims []string, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	enforcer := NewEnforcer()
	mw := NewMiddleware(MiddlewareConfig{Enforcer: enforcer})

	r.Use(func(c *gin.Context) {
		if role != "" {
			c.Set("role", role)
		}
		c.Set("permissions", claims)
		c.Next()
	})
	r.GET("/api/v1/tenant/employees", mw, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	r.GET("/api/v1/tenant/settings", mw, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	// Route submenu (level-submenu RBAC) — claim "resource.submenu.action"
	// memenuhi, begitu juga module-level "resource.action".
	r.GET("/api/v1/tenant/competency/events", mw, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	r.GET("/api/v1/tenant/competency/templates/:id", mw, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	// Route platform — claim tenant TIDAK boleh mem-bypass enforcer.
	r.GET("/api/v1/platform/rbac/roles", mw, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	r.PATCH("/api/v1/tenant/notifications/:id/read", mw, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	r.POST("/api/v1/tenant/notifications/read-all", mw, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func doMiddlewareRequest(r *gin.Engine, path string) *httptest.ResponseRecorder {
	return doMiddlewareRequestMethod(r, http.MethodGet, path)
}

func doMiddlewareRequestMethod(r *gin.Engine, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	r.ServeHTTP(w, req)
	return w
}

// Employee dengan permission employee.view → GET /tenant/employees diizinkan
// (meski role "Employee" tidak dikenal platform enforcer).
func TestNewMiddleware_ClaimMatch_AllowsTenantRole(t *testing.T) {
	r := setupMiddlewareTest([]string{"employee.view", "setting.view"}, "Employee")

	if w := doMiddlewareRequest(r, "/api/v1/tenant/employees"); w.Code != http.StatusOK {
		t.Errorf("GET employees: status = %d, want 200 (claim employee.view)", w.Code)
	}
	if w := doMiddlewareRequest(r, "/api/v1/tenant/settings"); w.Code != http.StatusOK {
		t.Errorf("GET settings: status = %d, want 200 (claim setting.view)", w.Code)
	}
}

// Employee tanpa permission yang cocok → 403 (deny by default).
func TestNewMiddleware_ClaimMismatch_Denies(t *testing.T) {
	// Hanya punya employee.view, minta resource lain
	r := setupMiddlewareTest([]string{"employee.view"}, "Employee")

	if w := doMiddlewareRequest(r, "/api/v1/tenant/settings"); w.Code != http.StatusForbidden {
		t.Errorf("GET settings: status = %d, want 403 (no setting.view claim)", w.Code)
	}
}

// Role dengan claim level-submenu "competency.events.view" (tanpa
// "competency.view") → GET /tenant/competency/events diizinkan.
func TestNewMiddleware_SubmenuClaim_AllowsSubmenuRoute(t *testing.T) {
	r := setupMiddlewareTest([]string{"competency.events.view"}, "Employee")

	if w := doMiddlewareRequest(r, "/api/v1/tenant/competency/events"); w.Code != http.StatusOK {
		t.Errorf("GET competency/events: status = %d, want 200 (claim competency.events.view)", w.Code)
	}
	// Route submenu lain yang TIDAK punya claim → 403.
	if w := doMiddlewareRequest(r, "/api/v1/tenant/competency/templates/abc"); w.Code != http.StatusForbidden {
		t.Errorf("GET competency/templates/abc: status = %d, want 403 (no templates claim)", w.Code)
	}
}

// Role dengan module-level "competency.view" → semua submenu competency diizinkan
// (backward compatible: claim module-level memenuhi semua submenu).
func TestNewMiddleware_ModuleClaim_GrantsSubmenus(t *testing.T) {
	r := setupMiddlewareTest([]string{"competency.view"}, "Employee")

	if w := doMiddlewareRequest(r, "/api/v1/tenant/competency/events"); w.Code != http.StatusOK {
		t.Errorf("GET competency/events: status = %d, want 200 (module claim competency.view)", w.Code)
	}
	if w := doMiddlewareRequest(r, "/api/v1/tenant/competency/templates/abc"); w.Code != http.StatusOK {
		t.Errorf("GET competency/templates/abc: status = %d, want 200 (module claim competency.view)", w.Code)
	}
}

// Wildcard "*" di claim → semua akses diizinkan.
func TestNewMiddleware_WildcardClaim_AllowsAll(t *testing.T) {
	r := setupMiddlewareTest([]string{"*"}, "super_admin")

	if w := doMiddlewareRequest(r, "/api/v1/tenant/employees"); w.Code != http.StatusOK {
		t.Errorf("GET employees: status = %d, want 200 (wildcard claim)", w.Code)
	}
}

// Claim kosong + role platform dikenal enforcer → fallback enforcer mengizinkan.
func TestNewMiddleware_EmptyClaim_FallsBackToEnforcer(t *testing.T) {
	// company_admin punya employee:* di enforcer non-DB default
	r := setupMiddlewareTest(nil, "company_admin")

	if w := doMiddlewareRequest(r, "/api/v1/tenant/employees"); w.Code != http.StatusOK {
		t.Errorf("GET employees: status = %d, want 200 (enforcer fallback)", w.Code)
	}
}

// Claim kosong + role tenant tak dikenal enforcer → deny (tidak ada bypass).
func TestNewMiddleware_EmptyClaim_UnknownRole_Denies(t *testing.T) {
	r := setupMiddlewareTest(nil, "Employee")

	if w := doMiddlewareRequest(r, "/api/v1/tenant/employees"); w.Code != http.StatusForbidden {
		t.Errorf("GET employees: status = %d, want 403 (unknown role, no claim)", w.Code)
	}
}

// Tanpa role di context → 403 (auth middleware required).
func TestNewMiddleware_MissingRole_Denies(t *testing.T) {
	r := setupMiddlewareTest([]string{"employee.view"}, "")

	if w := doMiddlewareRequest(r, "/api/v1/tenant/employees"); w.Code != http.StatusForbidden {
		t.Errorf("GET employees: status = %d, want 403 (missing role)", w.Code)
	}
}

// Anti-bypass: tenant role (Admin) dengan claim rbac.view TIDAK boleh mengakses
// route platform /api/v1/platform/rbac/* — claim hanya berlaku utk route tenant.
func TestNewMiddleware_ClaimDoesNotBypassPlatformRoutes(t *testing.T) {
	// Tenant Admin punya rbac.view di claim, tapi role "Admin" tak dikenal enforcer.
	r := setupMiddlewareTest([]string{"rbac.view", "employee.view"}, "Admin")

	if w := doMiddlewareRequest(r, "/api/v1/platform/rbac/roles"); w.Code != http.StatusForbidden {
		t.Errorf("GET platform rbac: status = %d, want 403 (claim must not bypass platform routes)", w.Code)
	}
}

// =========================================================================
// Notification mark-as-read routes — permission override
// =========================================================================
//
// PATCH .../notifications/:id/read and POST .../notifications/read-all used
// to be checked against the blanket path/method-derived permission
// ("notification.patch" / "notification.create"), neither of which was ever
// a declared, grantable permission — so these two self-service routes were
// permanently 403 regardless of what the user had. Fixed by checking
// "notification.view" instead (ownership itself is enforced in the service
// layer via authctx.GetUserID, same as the approval-actions bypass above).

func TestNewMiddleware_NotificationMarkAsRead_AllowsWithViewClaim(t *testing.T) {
	r := setupMiddlewareTest([]string{"notification.view"}, "Employee")

	if w := doMiddlewareRequestMethod(r, http.MethodPatch, "/api/v1/tenant/notifications/8c3c72a8-f17c-43f1-9396-e6dc5c5a23da/read"); w.Code != http.StatusOK {
		t.Errorf("PATCH notifications/:id/read: status = %d, want 200 (claim notification.view)", w.Code)
	}
}

func TestNewMiddleware_NotificationMarkAllAsRead_AllowsWithViewClaim(t *testing.T) {
	r := setupMiddlewareTest([]string{"notification.view"}, "Employee")

	if w := doMiddlewareRequestMethod(r, http.MethodPost, "/api/v1/tenant/notifications/read-all"); w.Code != http.StatusOK {
		t.Errorf("POST notifications/read-all: status = %d, want 200 (claim notification.view)", w.Code)
	}
}

// Regression guard for the exact bug reported: a "notification.manage" claim
// alone (no notification.view) does NOT satisfy these routes — they check
// specifically against "notification.view", the same permission gating the
// list/unread-count endpoints, not "notification.manage".
func TestNewMiddleware_NotificationMarkAsRead_ManageClaimAloneInsufficient(t *testing.T) {
	r := setupMiddlewareTest([]string{"notification.manage"}, "Employee")

	if w := doMiddlewareRequestMethod(r, http.MethodPatch, "/api/v1/tenant/notifications/8c3c72a8-f17c-43f1-9396-e6dc5c5a23da/read"); w.Code != http.StatusForbidden {
		t.Errorf("PATCH notifications/:id/read: status = %d, want 403 (manage claim alone, no view claim, no enforcer fallback for Employee)", w.Code)
	}
}

func TestNewMiddleware_NotificationMarkAsRead_NoClaim_Denies(t *testing.T) {
	r := setupMiddlewareTest(nil, "Employee")

	if w := doMiddlewareRequestMethod(r, http.MethodPatch, "/api/v1/tenant/notifications/8c3c72a8-f17c-43f1-9396-e6dc5c5a23da/read"); w.Code != http.StatusForbidden {
		t.Errorf("PATCH notifications/:id/read: status = %d, want 403 (no claim, unknown role)", w.Code)
	}
}

// =========================================================================
// isTenantPath helper
// =========================================================================

func TestIsTenantPath(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/api/v1/tenant/employees", true},
		{"/api/v1/tenant/settings/zones", true},
		{"/api/v1/platform/rbac/roles", false},
		{"/api/v1/platform/companies", false},
		{"/healthz", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isTenantPath(tt.path); got != tt.expected {
				t.Errorf("isTenantPath(%q) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}
}

// =========================================================================
// hasPermissionClaim helper
// =========================================================================

func TestHasPermissionClaim(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		perms     []string
		required  string
		expectYes bool
	}{
		{"exact match", []string{"employee.view", "setting.view"}, "employee.view", true},
		{"wildcard", []string{"*"}, "anything.anything", true},
		{"no match", []string{"employee.view"}, "setting.view", false},
		{"empty claim", nil, "employee.view", false},
		{"partial not enough", []string{"employee"}, "employee.view", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Set("permissions", tt.perms)
			got := hasPermissionClaim(c, tt.required)
			if got != tt.expectYes {
				t.Errorf("hasPermissionClaim(%v, %q) = %v, want %v", tt.perms, tt.required, got, tt.expectYes)
			}
		})
	}
}
