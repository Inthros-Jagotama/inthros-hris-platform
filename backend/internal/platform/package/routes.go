package pkgmgr

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mendaftarkan semua endpoint Package Management ke router group.
// Endpoint admin memerlukan JWT authentication + RBAC permission.
// Endpoint public di /api/v1/public/packages tidak memerlukan auth.
func RegisterRoutes(rg *gin.RouterGroup, publicRG *gin.RouterGroup, handler *Handler, authMW, rbacMW gin.HandlerFunc) {
	// Public endpoints (no auth) — didaftarkan langsung di main.go pada root router.
	// publicRG bisa nil jika dipanggil dari module.go (public endpoint sudah di main.go).
	if publicRG != nil {
		public := publicRG.Group("/packages")
		{
			public.GET("", handler.ListPublishedPackages)
		}
	}

	// Protected admin endpoints
	protected := rg.Group("")
	protected.Use(authMW)
	if rbacMW != nil {
		protected.Use(rbacMW)
	}
	{
		packages := protected.Group("/packages")
		{
			packages.GET("", handler.ListPackages)
			packages.POST("", handler.CreatePackage)
			packages.GET("/:id", handler.GetPackage)
			packages.PUT("/:id", handler.UpdatePackage)
			packages.DELETE("/:id", handler.DeletePackage)
			packages.POST("/:id/publish", handler.PublishPackage)
			packages.POST("/:id/unpublish", handler.UnpublishPackage)
			packages.GET("/:id/validate", handler.ValidateDependencies)
		}
	}
}
