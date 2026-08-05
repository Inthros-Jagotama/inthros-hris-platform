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

		// =====================================================================
		// Phase 2 KPI Enhancement Routes
		// =====================================================================

		// Performance Progress (nested under evaluation-details)
		perf.POST("/progress", handler.CreatePerformanceProgress)
		perf.GET("/evaluation-details/:id/progress", handler.ListPerformanceProgressByDetailID)
		perf.GET("/progress/:id", handler.GetPerformanceProgressByID)
		perf.PUT("/progress/:id", handler.UpdatePerformanceProgress)
		perf.DELETE("/progress/:id", handler.DeletePerformanceProgress)

		// Performance Comments (nested under evaluations)
		perf.POST("/comments", handler.CreatePerformanceComment)
		perf.GET("/evaluations/:id/comments", handler.ListPerformanceCommentsByEvaluationID)
		perf.GET("/comments/:id", handler.GetPerformanceCommentByID)
		perf.PUT("/comments/:id", handler.UpdatePerformanceComment)
		perf.DELETE("/comments/:id", handler.DeletePerformanceComment)

		// Performance Attachments (nested under evaluation-details)
		perf.POST("/attachments", handler.CreatePerformanceAttachment)
		perf.GET("/evaluation-details/:id/attachments", handler.ListPerformanceAttachmentsByDetailID)
		perf.GET("/attachments/:id", handler.GetPerformanceAttachmentByID)
		perf.PUT("/attachments/:id", handler.UpdatePerformanceAttachment)
		perf.DELETE("/attachments/:id", handler.DeletePerformanceAttachment)

		// Performance Ratings (Master data)
		perf.POST("/ratings", handler.CreatePerformanceRating)
		perf.GET("/ratings", handler.ListPerformanceRatings)
		perf.GET("/ratings/:id", handler.GetPerformanceRatingByID)
		perf.PUT("/ratings/:id", handler.UpdatePerformanceRating)
		perf.DELETE("/ratings/:id", handler.DeletePerformanceRating)

		// Performance Indicator Formulas (Master data)
		perf.POST("/indicator-formulas", handler.CreatePerformanceIndicatorFormula)
		perf.GET("/indicator-formulas", handler.ListPerformanceIndicatorFormulas)
		perf.GET("/indicator-formulas/:id", handler.GetPerformanceIndicatorFormulaByID)
		perf.PUT("/indicator-formulas/:id", handler.UpdatePerformanceIndicatorFormula)
		perf.DELETE("/indicator-formulas/:id", handler.DeletePerformanceIndicatorFormula)

		// Performance Logs (Audit trail - Read only)
		perf.GET("/logs", handler.ListPerformanceLogs)
		perf.GET("/evaluations/:id/logs", handler.ListPerformanceLogsByEvaluationID)
		perf.GET("/logs/:id", handler.GetPerformanceLogByID)

		// =====================================================================
		// Phase 3 Business Process Routes
		// =====================================================================

		// Evaluation with KPI Snapshot
		perf.POST("/evaluations/snapshot", handler.CreateEvaluationWithSnapshot)
		perf.GET("/evaluations/:id/full", handler.GetEvaluationWithDetails)

		// Evaluation Actuals Input
		perf.PUT("/evaluation-details/:id/actual", handler.UpdateEvaluationActual)
		perf.PUT("/evaluations/:id/actuals", handler.BulkUpdateEvaluationActuals)

		// Score Calculation & Progress
		perf.POST("/evaluations/:id/recalculate", handler.RecalculateEvaluationScore)
		perf.GET("/evaluations/:id/progress-summary", handler.GetEvaluationProgressSummary)

		// Workflow Status Transitions
		perf.POST("/evaluations/:id/submit", handler.SubmitEvaluation)
		perf.POST("/evaluations/:id/approve", handler.ApproveEvaluation)
		perf.POST("/evaluations/:id/reject", handler.RejectEvaluation)
		perf.POST("/evaluations/:id/complete", handler.CompleteEvaluation)

		// =====================================================================
		// Phase 4 Dashboard Routes
		// =====================================================================

		// Employee Dashboard
		perf.GET("/dashboard/employee/:employee_id", handler.GetEmployeeDashboard)

		// Manager Dashboard
		perf.GET("/dashboard/manager/:manager_id", handler.GetManagerDashboard)

		// HR Dashboard
		perf.GET("/dashboard/hr", handler.GetHRDashboard)
	}
}
