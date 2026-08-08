package leave

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan semua endpoint Leave ke router group tenant.
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	leave := rg.Group("/leave")
	{
		// Leave Types
		leave.POST("/types", handler.CreateLeaveType)
		leave.GET("/types", handler.ListLeaveTypes)
		leave.GET("/types/:id", handler.GetLeaveTypeByID)
		leave.PUT("/types/:id", handler.UpdateLeaveType)
		leave.DELETE("/types/:id", handler.DeleteLeaveType)

		// Accrual Policies
		leave.POST("/accrual-policies", handler.CreateAccrualPolicy)
		leave.GET("/accrual-policies", handler.ListAccrualPolicies)
		leave.GET("/accrual-policies/:id", handler.GetAccrualPolicyByID)
		leave.PUT("/accrual-policies/:id", handler.UpdateAccrualPolicy)
		leave.DELETE("/accrual-policies/:id", handler.DeleteAccrualPolicy)

		// Leave Reasons
		leave.POST("/reasons", handler.CreateLeaveReason)
		leave.GET("/reasons", handler.ListLeaveReasons)
		leave.GET("/reasons/:id", handler.GetLeaveReasonByID)
		leave.PUT("/reasons/:id", handler.UpdateLeaveReason)
		leave.DELETE("/reasons/:id", handler.DeleteLeaveReason)

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
	}
}
