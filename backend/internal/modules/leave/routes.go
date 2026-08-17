package leave

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// requireLeaveSettings menggating endpoint yang MENGUBAH (create/update/
// delete) master cuti (types, accrual-policies, reasons) beserta GET detail
// admin-nya — supaya HANYA permission submenu "leave.settings.<action>" yang
// berlaku (bukan module-level "leave.view"). Middleware RBAC global
// (authz.NewMiddleware) menganggap module-level otomatis mencakup semua
// submenu-nya — cocok untuk kebanyakan resource, tapi tidak untuk aksi
// admin-config ini, yang harus benar-benar terpisah dari sekadar melihat/
// mengajukan cuti sehari-hari ("leave.view", yang dimiliki hampir semua role
// termasuk Employee default). Pola sama seperti requireApprovalSettings di
// modul approval.
//
// GET /types dan /reasons (daftar) SENGAJA TIDAK dibatasi middleware ini —
// keduanya juga dipakai halaman Leave utama (dropdown tipe cuti & alasan
// saat mengajukan permohonan) yang harus tetap bisa diakses siapa pun dengan
// "leave.view". GET /accrual-policies tidak dipakai halaman Leave utama,
// jadi ikut di-gate.
func requireLeaveSettings(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		required := "leave.settings." + action
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
				"message": "You don't have permission to manage leave settings",
				"details": gin.H{"required": required},
			},
		})
	}
}

// RegisterRoutes mendaftarkan semua endpoint Leave ke router group tenant.
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	leave := rg.Group("/leave")
	{
		// Leave Types (admin-config)
		leave.POST("/types", requireLeaveSettings("create"), handler.CreateLeaveType)
		leave.GET("/types", handler.ListLeaveTypes)
		leave.GET("/types/:id", requireLeaveSettings("view"), handler.GetLeaveTypeByID)
		leave.PUT("/types/:id", requireLeaveSettings("update"), handler.UpdateLeaveType)
		leave.DELETE("/types/:id", requireLeaveSettings("delete"), handler.DeleteLeaveType)

		// Accrual Policies (admin-config)
		leave.POST("/accrual-policies", requireLeaveSettings("create"), handler.CreateAccrualPolicy)
		leave.GET("/accrual-policies", requireLeaveSettings("view"), handler.ListAccrualPolicies)
		leave.GET("/accrual-policies/:id", requireLeaveSettings("view"), handler.GetAccrualPolicyByID)
		leave.PUT("/accrual-policies/:id", requireLeaveSettings("update"), handler.UpdateAccrualPolicy)
		leave.DELETE("/accrual-policies/:id", requireLeaveSettings("delete"), handler.DeleteAccrualPolicy)

		// Leave Reasons (admin-config)
		leave.POST("/reasons", requireLeaveSettings("create"), handler.CreateLeaveReason)
		leave.GET("/reasons", handler.ListLeaveReasons)
		leave.GET("/reasons/:id", requireLeaveSettings("view"), handler.GetLeaveReasonByID)
		leave.PUT("/reasons/:id", requireLeaveSettings("update"), handler.UpdateLeaveReason)
		leave.DELETE("/reasons/:id", requireLeaveSettings("delete"), handler.DeleteLeaveReason)

		// Leave Requests
		leave.POST("/requests", handler.CreateLeaveRequest)
		leave.GET("/requests", handler.ListLeaveRequests)
		leave.GET("/requests/:id", handler.GetLeaveRequestByID)
		leave.PUT("/requests/:id/status", handler.UpdateLeaveRequestStatus)
		leave.DELETE("/requests/:id", handler.DeleteLeaveRequest)

		// Leave Request Details
		leave.GET("/requests/:id/details", handler.ListLeaveRequestDetails)

		// Leave Balances
		leave.GET("/balances", handler.ListLeaveBalances)
		leave.GET("/balances/employees/:employeeId/types/:leaveTypeId", handler.GetLeaveBalance)

		// Calendar
		leave.GET("/calendar", handler.GetEmployeeCalendar)

		// Reports
		leave.GET("/reports/usage", handler.GetLeaveUsageReport)
	}
}
