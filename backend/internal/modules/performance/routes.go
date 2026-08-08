package performance

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan semua endpoint Performance Management ke router group tenant.
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	perf := rg.Group("/performance")
	{
		// Shared Master Data (used by both KPI and OKR)
		// Performance Periods
		perf.POST("/periods", handler.CreatePerformancePeriod)
		perf.GET("/periods", handler.ListPerformancePeriods)
		perf.GET("/periods/:id", handler.GetPerformancePeriodByID)
		perf.PUT("/periods/:id", handler.UpdatePerformancePeriod)
		perf.DELETE("/periods/:id", handler.DeletePerformancePeriod)

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
		perf.GET("/logs/:id", handler.GetPerformanceLogByID)

		// =====================================================================
		// KPI Sub-module: /performance/kpi/*
		// =====================================================================
		kpi := perf.Group("/kpi")
		{
			// Performance Perspectives (BSC Perspectives)
			kpi.POST("/perspectives", handler.CreatePerformancePerspective)
			kpi.GET("/perspectives", handler.ListPerformancePerspectives)
			kpi.GET("/perspectives/:id", handler.GetPerformancePerspectiveByID)
			kpi.PUT("/perspectives/:id", handler.UpdatePerformancePerspective)
			kpi.DELETE("/perspectives/:id", handler.DeletePerformancePerspective)

			// Performance Templates (BSC Templates)
			kpi.GET("/my-context", handler.GetMyKPIContext)
			kpi.GET("/templates/organization-scope", handler.ListTemplateOrganizationScope)
			kpi.POST("/templates", handler.CreatePerformanceTemplate)
			kpi.GET("/templates", handler.ListPerformanceTemplates)
			kpi.GET("/templates/:id", handler.GetPerformanceTemplateByID)
			kpi.PUT("/templates/:id", handler.UpdatePerformanceTemplate)
			kpi.DELETE("/templates/:id", handler.DeletePerformanceTemplate)

			// Performance Indicators (KPI Indicators)
			kpi.POST("/indicators", handler.CreatePerformanceIndicator)
			kpi.GET("/indicators", handler.ListPerformanceIndicators)
			kpi.GET("/indicators/:id", handler.GetPerformanceIndicatorByID)
			kpi.PUT("/indicators/:id", handler.UpdatePerformanceIndicator)
			kpi.DELETE("/indicators/:id", handler.DeletePerformanceIndicator)

			// Performance Evaluations
			kpi.POST("/evaluations", handler.CreatePerformanceEvaluation)
			kpi.GET("/evaluations", handler.ListPerformanceEvaluations)
			kpi.GET("/evaluations/:id", handler.GetPerformanceEvaluationByID)
			kpi.PUT("/evaluations/:id", handler.UpdatePerformanceEvaluation)
			kpi.PUT("/evaluations/:id/status", handler.UpdateEvaluationStatus)
			kpi.DELETE("/evaluations/:id", handler.DeletePerformanceEvaluation)

			// Evaluation with KPI Snapshot
			kpi.POST("/evaluations/snapshot", handler.CreateEvaluationWithSnapshot)
			kpi.GET("/evaluations/:id/full", handler.GetEvaluationWithDetails)
			kpi.GET("/evaluations/:id/logs", handler.ListPerformanceLogsByEvaluationID)

			// Evaluation Details
			kpi.POST("/evaluation-details", handler.CreateEvaluationDetail)
			kpi.GET("/evaluations/:id/details", handler.ListEvaluationDetails)
			kpi.PUT("/evaluation-details/:id", handler.UpdateEvaluationDetail)
			kpi.DELETE("/evaluation-details/:id", handler.DeleteEvaluationDetail)

			// Evaluation Target/Actual Input (two-phase: target while DRAFT,
			// actual once TARGET_APPROVED)
			kpi.PUT("/evaluation-details/:id/target", handler.UpdateEvaluationTarget)
			kpi.PUT("/evaluation-details/:id/actual", handler.UpdateEvaluationActual)
			kpi.PUT("/evaluations/:id/actuals", handler.BulkUpdateEvaluationActuals)

			// Program Items (employee-authored, no HR template)
			kpi.POST("/program-items", handler.CreateProgramItem)
			kpi.GET("/evaluations/:id/program-items", handler.ListProgramItems)
			kpi.PUT("/program-items/:id/target", handler.UpdateProgramItemTarget)
			kpi.PUT("/program-items/:id/actual", handler.UpdateProgramItemActual)
			kpi.DELETE("/program-items/:id", handler.DeleteProgramItem)

			// Performance Targets
			kpi.POST("/targets", handler.CreatePerformanceTarget)
			kpi.GET("/evaluations/:id/targets", handler.ListPerformanceTargets)
			kpi.PUT("/targets/:id", handler.UpdatePerformanceTarget)
			kpi.DELETE("/targets/:id", handler.DeletePerformanceTarget)

			// Performance Progress
			kpi.POST("/progress", handler.CreatePerformanceProgress)
			kpi.GET("/evaluation-details/:id/progress", handler.ListPerformanceProgressByDetailID)
			kpi.GET("/progress/:id", handler.GetPerformanceProgressByID)
			kpi.PUT("/progress/:id", handler.UpdatePerformanceProgress)
			kpi.DELETE("/progress/:id", handler.DeletePerformanceProgress)

			// Performance Comments
			kpi.POST("/comments", handler.CreatePerformanceComment)
			kpi.GET("/evaluations/:id/comments", handler.ListPerformanceCommentsByEvaluationID)
			kpi.GET("/comments/:id", handler.GetPerformanceCommentByID)
			kpi.PUT("/comments/:id", handler.UpdatePerformanceComment)
			kpi.DELETE("/comments/:id", handler.DeletePerformanceComment)

			// Performance Attachments
			kpi.POST("/attachments", handler.CreatePerformanceAttachment)
			kpi.GET("/evaluation-details/:id/attachments", handler.ListPerformanceAttachmentsByDetailID)
			kpi.GET("/attachments/:id", handler.GetPerformanceAttachmentByID)
			kpi.PUT("/attachments/:id", handler.UpdatePerformanceAttachment)
			kpi.DELETE("/attachments/:id", handler.DeletePerformanceAttachment)

			// Score Calculation & Progress
			kpi.POST("/evaluations/:id/recalculate", handler.RecalculateEvaluationScore)
			kpi.GET("/evaluations/:id/progress-summary", handler.GetEvaluationProgressSummary)

			// Workflow Status Transitions
			// Phase 1: DRAFT -> TARGET_SUBMITTED -> TARGET_APPROVED ("Ajukan Target")
			kpi.POST("/evaluations/:id/submit-target", handler.SubmitTarget)
			kpi.POST("/evaluations/:id/approve-target", handler.ApproveTarget)
			kpi.POST("/evaluations/:id/reject-target", handler.RejectTarget)
			// Phase 2: TARGET_APPROVED -> SUBMITTED -> APPROVED -> COMPLETED ("Ajukan Realisasi")
			kpi.POST("/evaluations/:id/submit", handler.SubmitEvaluation)
			kpi.POST("/evaluations/:id/approve", handler.ApproveEvaluation)
			kpi.POST("/evaluations/:id/reject", handler.RejectEvaluation)
			kpi.POST("/evaluations/:id/complete", handler.CompleteEvaluation)

			// KPI Dashboard
			kpi.GET("/dashboard/employee/:employee_id", handler.GetEmployeeDashboard)
			kpi.GET("/dashboard/manager/:manager_id", handler.GetManagerDashboard)
			kpi.GET("/dashboard/hr", handler.GetHRDashboard)

			// =================================================================
			// Phase 5: Performance Scoring Configuration
			// =================================================================

			// Performance Components (Master data) — locked to exactly 3 seeded
			// rows (KPI/PROGRAM/SUBORDINATE, migration 064); create/delete are
			// intentionally not exposed so the set can't grow or shrink. Only
			// GET/PUT (rename, reorder, activate/deactivate) remain.
			kpi.GET("/components", handler.ListPerformanceComponents)
			kpi.GET("/components/:id", handler.GetPerformanceComponentByID)
			kpi.PUT("/components/:id", handler.UpdatePerformanceComponent)

			// Organization Component Weight Configuration
			kpi.POST("/organization-components", handler.UpsertOrganizationComponent)
			kpi.GET("/organizations/:organization_id/components", handler.ListOrganizationComponents)
			kpi.DELETE("/organization-components/:id", handler.DeleteOrganizationComponent)

			// Evaluation Scoring Engine
			kpi.GET("/evaluations/:id/components", handler.ListEvaluationComponents)
			kpi.POST("/evaluations/:id/calculate-scoring", handler.CalculateEvaluationComponentScoring)
			kpi.POST("/periods/:period_id/recalculate-scoring", handler.RecalculatePeriodScoring)
			kpi.PUT("/evaluations/:id/components/:component_id", handler.UpdateEvaluationComponentScore)
		}
	}
}

// RegisterOKRRoutes mendaftarkan semua endpoint OKR ke router group tenant.
func RegisterOKRRoutes(rg *gin.RouterGroup, handler *OKRHandler) {
	okr := rg.Group("/performance/okr")
	{
		// OKR Templates
		templates := okr.Group("/templates")
		{
			templates.POST("", handler.CreateTemplate)
			templates.GET("/objective-scope", handler.GetObjectiveCreationScope)
			templates.GET("", handler.ListTemplates)
			templates.GET("/:id", handler.GetTemplateByID)
			templates.PUT("/:id", handler.UpdateTemplate)
			templates.DELETE("/:id", handler.DeleteTemplate)
			templates.POST("/:id/duplicate", handler.DuplicateTemplate)
			templates.GET("/:id/objectives", handler.ListObjectivesByTemplateID)
		}

		// OKR Objectives
		objectives := okr.Group("/objectives")
		{
			objectives.POST("", handler.CreateObjective)
			objectives.GET("/:id", handler.GetObjectiveByID)
			objectives.PUT("/:id", handler.UpdateObjective)
			objectives.DELETE("/:id", handler.DeleteObjective)
			objectives.GET("/:id/key-results", handler.ListKeyResultsByObjectiveID)
		}

		// OKR Key Results
		keyResults := okr.Group("/key-results")
		{
			keyResults.POST("", handler.CreateKeyResult)
			keyResults.GET("/:id", handler.GetKeyResultByID)
			keyResults.PUT("/:id", handler.UpdateKeyResult)
			keyResults.DELETE("/:id", handler.DeleteKeyResult)
		}

		// OKR Evaluations
		evaluations := okr.Group("/evaluations")
		{
			evaluations.POST("", handler.CreateEvaluationWithSnapshot)
			evaluations.GET("", handler.ListEvaluations)
			evaluations.GET("/:id", handler.GetEvaluationByID)
			evaluations.GET("/:id/details", handler.GetEvaluationWithDetails)
			evaluations.PUT("/:id", handler.UpdateEvaluation)
			evaluations.DELETE("/:id", handler.DeleteEvaluation)

			// Bulk update actuals
			evaluations.PUT("/:id/actuals", handler.BulkUpdateEvaluationActuals)

			// Recalculate score
			evaluations.POST("/:id/recalculate", handler.RecalculateEvaluationScore)

			// Employee-proposed Key Results (DRAFT phase)
			evaluations.POST("/:id/key-results", handler.CreateEvaluationKeyResult)

			// Workflow actions — Key Result proposal phase
			evaluations.POST("/:id/submit-key-results", handler.SubmitKeyResults)
			evaluations.POST("/:id/approve-key-results", handler.ApproveKeyResults)
			evaluations.POST("/:id/reject-key-results", handler.RejectKeyResults)

			// Workflow actions — assessment phase
			evaluations.POST("/:id/submit", handler.SubmitEvaluation)
			evaluations.POST("/:id/approve", handler.ApproveEvaluation)
			evaluations.POST("/:id/reject", handler.RejectEvaluation)
			evaluations.POST("/:id/complete", handler.CompleteEvaluation)

			// Evaluation comments
			evaluations.GET("/:id/comments", handler.ListCommentsByEvaluationID)
		}

		// OKR Evaluation Key Results (employee-proposed, DRAFT phase) — distinct
		// prefix from /key-results above, which is Template-level master data.
		evaluationKeyResults := okr.Group("/evaluation-key-results")
		{
			evaluationKeyResults.PUT("/:id/target", handler.UpdateEvaluationKeyResultTarget)
			evaluationKeyResults.DELETE("/:id", handler.DeleteEvaluationKeyResult)
		}

		// OKR Evaluation Details
		details := okr.Group("/evaluation-details")
		{
			details.PUT("/:id", handler.UpdateEvaluationDetailActual)
			details.GET("/:id/progress", handler.ListProgressByDetailID)
			details.GET("/:id/attachments", handler.ListAttachmentsByDetailID)
		}

		// OKR Progress
		progress := okr.Group("/progress")
		{
			progress.POST("", handler.CreateProgress)
			progress.GET("/:id", handler.GetProgressByID)
			progress.PUT("/:id", handler.UpdateProgress)
			progress.DELETE("/:id", handler.DeleteProgress)
		}

		// OKR Comments
		comments := okr.Group("/comments")
		{
			comments.POST("", handler.CreateComment)
			comments.PUT("/:id", handler.UpdateComment)
			comments.DELETE("/:id", handler.DeleteComment)
		}

		// OKR Attachments
		attachments := okr.Group("/attachments")
		{
			attachments.POST("", handler.CreateAttachment)
			attachments.DELETE("/:id", handler.DeleteAttachment)
		}

		// OKR Dashboard
		okr.GET("/dashboard/hr", handler.GetHRDashboard)

		// OKR Self-assessment context
		okr.GET("/my-context", handler.GetMyOKRContext)
	}
}
