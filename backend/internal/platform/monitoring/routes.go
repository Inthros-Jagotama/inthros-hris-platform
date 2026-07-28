package monitoring

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan semua endpoint Monitoring ke router group.
// Semua endpoint memerlukan JWT authentication.
// RBAC permission middleware opsional (skip jika nil).
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW, rbacMW gin.HandlerFunc) {
	protected := rg.Group("")
	protected.Use(authMW)
	if rbacMW != nil {
		protected.Use(rbacMW)
	}
	{
		monitoring := protected.Group("/monitoring")
		{
			monitoring.GET("/health", handler.HealthCheck)
			monitoring.GET("/pool", handler.PoolStats)
			monitoring.GET("/tenants", handler.TenantHealth)
			monitoring.GET("/tenants/:id", handler.TenantDetail)
			monitoring.GET("/seed-status", handler.SeedDataStatus)
		}
	}
}
