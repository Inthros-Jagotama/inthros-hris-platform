package approval

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// requireApprovalSettings menggating endpoint yang MENGUBAH (create/update/
// delete) alur persetujuan beserta langkah-langkahnya, dan GET langkah-
// langkahnya (mengungkap hierarki approval internal) — supaya HANYA
// permission submenu "approval.settings.<action>" (bukan module-level
// "approval.view") yang berlaku. Middleware RBAC global (authz.NewMiddleware)
// menganggap module-level otomatis mencakup semua submenu-nya — cocok untuk
// kebanyakan resource, tapi tidak untuk aksi admin-config ini, yang harus
// benar-benar terpisah dari sekadar melihat tugas approval sehari-hari
// ("approval.view", yang dimiliki hampir semua role termasuk Employee
// default). Pola sama seperti requireSensitiveFieldSettings di modul
// employee.
//
// GET /flows, /flows/:flowId, /active-flow, /available-modules SENGAJA
// TIDAK dibatasi middleware ini — daftar nama flow (tanpa detail step)
// juga dipakai halaman Approvals utama (filter dropdown per flow) yang
// harus tetap bisa diakses siapa pun dengan "approval.view".
func requireApprovalSettings(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		required := "approval.settings." + action
		for _, p := range c.GetStringSlice("permissions") {
			if p == "*" || p == required {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "You don't have permission to manage approval flows",
				"details": gin.H{"required": required},
			},
		})
	}
}

// RegisterRoutes mendaftarkan semua endpoint Approval Engine ke router group tenant.
// Semua endpoint di bawah /api/v1/tenant/approval
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	approval := rg.Group("/approval")
	{
		// =================================================================
		// Approval Flows (Master alur persetujuan)
		// =================================================================
		approval.POST("/flows", requireApprovalSettings("create"), handler.CreateFlow)
		approval.GET("/flows", handler.ListFlows)
		approval.GET("/active-flow", handler.GetActiveFlowByModule)
		approval.GET("/flows/:flowId", handler.GetFlowByID)
		approval.PUT("/flows/:flowId", requireApprovalSettings("update"), handler.UpdateFlow)
		approval.DELETE("/flows/:flowId", requireApprovalSettings("delete"), handler.DeleteFlow)
		approval.GET("/available-modules", handler.ListAvailableModules)

		// =================================================================
		// Approval Flow Steps (Langkah-langkah dalam alur)
		// =================================================================
		approval.POST("/flows/:flowId/steps", requireApprovalSettings("create"), handler.CreateStep)
		approval.GET("/flows/:flowId/steps", requireApprovalSettings("view"), handler.ListSteps)
		approval.PUT("/flows/:flowId/steps/:stepId", requireApprovalSettings("update"), handler.UpdateStep)
		approval.DELETE("/flows/:flowId/steps/:stepId", requireApprovalSettings("delete"), handler.DeleteStep)

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
		approval.GET("/tasks/done", handler.ListMyDoneTasks)
	}
}
