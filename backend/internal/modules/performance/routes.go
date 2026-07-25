package performance

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan semua endpoint Performance Management ke router group tenant.
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	perf := rg.Group("/performance")
	{
		// Performance Periods
		perf.POST("/periods", handler.CreatePerformancePeriod)
		perf.GET("/periods", handler.ListPerformancePeriods)
		perf.GET("/periods/:id", handler.GetPerformancePeriodByID)
		perf.PUT("/periods/:id", handler.UpdatePerformancePeriod)
		perf.DELETE("/periods/:id", handler.DeletePerformancePeriod)

		// Performance Perspectives (BSC Perspectives)
		perf.POST("/perspectives", handler.CreatePerformancePerspective)
		perf.GET("/perspectives", handler.ListPerformancePerspectives)
		perf.GET("/perspectives/:id", handler.GetPerformancePerspectiveByID)
		perf.PUT("/perspectives/:id", handler.UpdatePerformancePerspective)
		perf.DELETE("/perspectives/:id", handler.DeletePerformancePerspective)

		// Performance Templates (BSC Templates)
		perf.POST("/templates", handler.CreatePerformanceTemplate)
		perf.GET("/templates", handler.ListPerformanceTemplates)
		perf.GET("/templates/:id", handler.GetPerformanceTemplateByID)
		perf.PUT("/templates/:id", handler.UpdatePerformanceTemplate)
		perf.DELETE("/templates/:id", handler.DeletePerformanceTemplate)

		// Performance Indicators (KPI Indicators)
		perf.POST("/indicators", handler.CreatePerformanceIndicator)
		perf.GET("/indicators", handler.ListPerformanceIndicators)
		perf.GET("/indicators/:id", handler.GetPerformanceIndicatorByID)
		perf.PUT("/indicators/:id", handler.UpdatePerformanceIndicator)
		perf.DELETE("/indicators/:id", handler.DeletePerformanceIndicator)

		// Performance Evaluations
		perf.POST("/evaluations", handler.CreatePerformanceEvaluation)
		perf.GET("/evaluations", handler.ListPerformanceEvaluations)
		perf.GET("/evaluations/:id", handler.GetPerformanceEvaluationByID)
		perf.PUT("/evaluations/:id", handler.UpdatePerformanceEvaluation)
		perf.PUT("/evaluations/:id/status", handler.UpdateEvaluationStatus)
		perf.DELETE("/evaluations/:id", handler.DeletePerformanceEvaluation)

		// Evaluation Details (nested under evaluations)
		perf.POST("/evaluation-details", handler.CreateEvaluationDetail)
		perf.GET("/evaluations/:id/details", handler.ListEvaluationDetails)
		perf.PUT("/evaluation-details/:id", handler.UpdateEvaluationDetail)
		perf.DELETE("/evaluation-details/:id", handler.DeleteEvaluationDetail)

		// Performance Targets (nested under evaluations)
		perf.POST("/targets", handler.CreatePerformanceTarget)
		perf.GET("/evaluations/:id/targets", handler.ListPerformanceTargets)
		perf.PUT("/targets/:id", handler.UpdatePerformanceTarget)
		perf.DELETE("/targets/:id", handler.DeletePerformanceTarget)
	}
}
