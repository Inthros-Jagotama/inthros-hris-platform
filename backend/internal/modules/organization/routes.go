package organization

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	orgs := rg.Group("/organizations")
	{
		// Basic CRUD
		orgs.POST("", handler.Create)
		orgs.GET("", handler.List)
		orgs.GET("/:id", handler.GetByID)
		orgs.PUT("/:id", handler.Update)
		orgs.DELETE("/:id", handler.Delete)

		// History & Audit Trail
		orgs.GET("/history", handler.ListHistory)

		// Versioning & Restore
		orgs.POST("/versions", handler.CreateVersion)
		orgs.GET("/versions", handler.ListVersions)
		orgs.GET("/versions/:id", handler.GetVersion)
		orgs.GET("/versions/:id/diff/:targetId", handler.DiffVersions)
		orgs.POST("/versions/:id/restore", handler.RestoreVersion)

		// Cloning
		orgs.POST("/clone", handler.CloneTree)
	}

	// Organization Summaries CRUD
	summaries := rg.Group("/organization-summaries")
	{
		summaries.POST("", handler.CreateSummary)
		summaries.GET("", handler.ListSummaries)
		summaries.GET("/stats", handler.GetSummaryStats)
		summaries.GET("/:id", handler.GetSummaryByID)
		summaries.PUT("/:id", handler.UpdateSummary)
		summaries.DELETE("/:id", handler.DeleteSummary)
	}
}
