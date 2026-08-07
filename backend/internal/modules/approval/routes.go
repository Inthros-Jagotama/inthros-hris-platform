package approval

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan semua endpoint Approval Engine ke router group tenant.
// Semua endpoint di bawah /api/v1/tenant/approval
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	approval := rg.Group("/approval")
	{
		// =================================================================
		// Approval Flows (Master alur persetujuan)
		// =================================================================
		approval.POST("/flows", handler.CreateFlow)
		approval.GET("/flows", handler.ListFlows)
		approval.GET("/flows/:flowId", handler.GetFlowByID)
		approval.PUT("/flows/:flowId", handler.UpdateFlow)
		approval.DELETE("/flows/:flowId", handler.DeleteFlow)
		approval.GET("/available-modules", handler.ListAvailableModules)

		// =================================================================
		// Approval Flow Steps (Langkah-langkah dalam alur)
		// =================================================================
		approval.POST("/flows/:flowId/steps", handler.CreateStep)
		approval.GET("/flows/:flowId/steps", handler.ListSteps)
		approval.PUT("/flows/:flowId/steps/:stepId", handler.UpdateStep)
		approval.DELETE("/flows/:flowId/steps/:stepId", handler.DeleteStep)

		// =================================================================
		// Approval Instances (Instance persetujuan)
		// =================================================================
		approval.POST("/instances", handler.CreateInstance)
		approval.GET("/instances", handler.ListInstances)
		approval.GET("/instances/:id", handler.GetInstanceByID)
		approval.POST("/instances/:id/cancel", handler.CancelInstance)

		// =================================================================
		// Approval Actions (Aksi approve/reject)
		// =================================================================
		approval.POST("/instances/:id/actions", handler.SubmitAction)

		// =================================================================
		// Approval Tasks (Task approval per approver)
		// =================================================================
		approval.GET("/tasks/pending", handler.ListMyPendingTasks)
	}
}
