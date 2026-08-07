package authz

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// MiddlewareConfig untuk middleware RBAC.
type MiddlewareConfig struct {
	Enforcer *Enforcer
}

// hasPermissionClaim memeriksa apakah JWT permissions claim memuat permission
// yang diminta (format Spatie "resource.action", mis. "employee.view").
// "*" berarti wildcard (super_admin / akses penuh).
func hasPermissionClaim(c *gin.Context, required string) bool {
	for _, p := range c.GetStringSlice("permissions") {
		if p == "*" || p == required {
			return true
		}
	}
	return false
}

// isTenantPath menentukan apakah path termasuk route tenant (/api/v1/tenant/...).
// Hanya route tenant yang boleh diizinkan lewat JWT permissions claim — route
// platform tetap murni platform enforcer agar tidak ada bypass dari claim role
// tenant (mis. tenant Admin punya rbac.* di claim, tidak boleh akses
// /api/v1/platform/rbac/*).
func isTenantPath(path string) bool {
	return strings.Contains(path, "/api/v1/tenant/")
}

// NewMiddleware membuat middleware Gin untuk permission checking.
//
// Middleware ini HARUS ditempatkan SETELAH AuthJWT middleware
// agar claims (role + permissions) sudah tersedia di context.
//
// Prioritas pengecekan:
//  1. JWT permissions claim (format Spatie "resource.action", mis. "employee.view").
//     User tenant (employee/admin) membawa permissions dari tenant DB di dalam
//     JWT — ini sumber kebenaran RBAC tenant, termasuk role yang tidak dikenal
//     platform enforcer (mis. "Employee"/"Admin").
//     HANYA berlaku untuk route tenant (/api/v1/tenant/...) — lihat isTenantPath.
//  2. Platform enforcer (role hierarchy) sebagai fallback untuk platform roles
//     (super_admin/company_admin/manager) yang permissions claim-nya kosong.
func NewMiddleware(cfg MiddlewareConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Role sudah di-set oleh AuthJWT middleware
		role := c.GetString("role")
		if role == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "Missing role in context. Auth middleware required.",
				},
			})
			return
		}

		// Approve/reject pada approval task adalah aksi berbasis assignment
		// (siapapun yang menjadi assignee task pending berhak bertindak), bukan
		// permission berbasis role generik seperti "approval.create" — role
		// Employee (default, view-only) TIDAK pernah punya "approval.create",
		// sehingga blanket check di bawah akan mem-block approver yang sah.
		// Otorisasi sebenarnya sudah divalidasi di service.SubmitAction
		// (actor harus punya ApprovalTask PENDING miliknya untuk step berjalan).
		if c.FullPath() == "/api/v1/tenant/approval/instances/:id/actions" {
			c.Next()
			return
		}

		// Extract resource dari path
		resource := ResourceFromPath(c.FullPath())
		if resource == "" {
			// Path tidak dikenal, allow by default
			c.Next()
			return
		}

		// Extract action dari HTTP method
		action := ActionFromMethod(c.Request.Method)

		// 1. Cek permissions claim dari JWT (RBAC tenant, format "resource.action").
		//    Hanya untuk route tenant — route platform murni enforcer (anti-bypass).
		if isTenantPath(c.FullPath()) && hasPermissionClaim(c, resource+"."+action) {
			c.Next()
			return
		}

		// 2. Fallback: platform enforcer (role hierarchy) untuk platform roles.
		if cfg.Enforcer.Check(role, resource, action) == DecisionDeny {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "You don't have permission to perform this action",
					"details": gin.H{
						"role":     role,
						"resource": resource,
						"action":   action,
					},
				},
			})
			return
		}

		c.Next()
	}
}

// RequirePermission adalah middleware factory untuk validasi
// permission spesifik pada suatu route. Digunakan untuk route
// yang memerlukan permission berbeda dari action default.
//
// Contoh:
//
//	router.POST("/companies/:id/suspend", RequirePermission("company", "suspend"), handler.Suspend)
//	router.POST("/companies/:id/activate", RequirePermission("company", "activate"), handler.Activate)
func RequirePermission(enforcer *Enforcer, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "Missing role in context.",
				},
			})
			return
		}

		// 1. Cek permissions claim dari JWT (RBAC tenant) — sama seperti NewMiddleware.
		//    Hanya untuk route tenant — route platform murni enforcer (anti-bypass).
		if isTenantPath(c.FullPath()) && hasPermissionClaim(c, resource+"."+action) {
			c.Next()
			return
		}

		// 2. Fallback: platform enforcer (role hierarchy) untuk platform roles.
		if enforcer.Check(role, resource, action) == DecisionDeny {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "You don't have permission to perform this action",
					"details": gin.H{
						"role":     role,
						"resource": resource,
						"action":   action,
					},
				},
			})
			return
		}

		c.Next()
	}
}
