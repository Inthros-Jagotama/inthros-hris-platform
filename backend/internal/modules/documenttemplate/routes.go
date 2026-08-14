package documenttemplate

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	templates := rg.Group("/document-templates")
	{
		templates.GET("", handler.List)
		templates.POST("", handler.Create)
		templates.POST("/from-default", handler.CreateFromDefault)
		templates.GET("/:id", handler.GetByID)
		templates.PUT("/:id", handler.Update)
		templates.DELETE("/:id", handler.Delete)
		templates.POST("/:id/activate", handler.Activate)
		templates.POST("/:id/deactivate", handler.Deactivate)
		templates.GET("/:id/versions", handler.ListVersions)
		templates.POST("/:id/versions", handler.CreateVersion)
	}
}
