package rbac

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan endpoint manajemen RBAC tenant.
// Semua route berada di bawah grup /api/v1/tenant/rbac (di-set oleh caller).
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	rbac := rg.Group("/rbac")
	{
		// Roles
		rbac.GET("/roles", handler.ListRoles)
		rbac.POST("/roles", handler.CreateRole)
		rbac.PUT("/roles/:id", handler.UpdateRole)
		rbac.DELETE("/roles/:id", handler.DeleteRole)

		// Permissions
		rbac.GET("/permissions", handler.ListPermissions)

		// Role ↔ Permission (assign/replace)
		rbac.PUT("/roles/:id/permissions", handler.AssignRolePermissions)

		// Users & User ↔ Role (assign/replace)
		rbac.GET("/users", handler.ListUsers)
		rbac.PUT("/users/:id/roles", handler.AssignUserRoles)
	}
}
