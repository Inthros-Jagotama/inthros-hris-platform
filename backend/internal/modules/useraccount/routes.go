package useraccount

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan endpoint authenticated user-account
// di bawah grup /api/v1/tenant (dengan auth + RBAC oleh caller).
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	acc := rg.Group("/user-accounts")
	{
		// GET /api/v1/tenant/user-accounts/me — employee_id milik user yang login
		acc.GET("/me", handler.GetMyAccount)
		// POST /api/v1/tenant/user-accounts/employees/:employeeId — buat akun
		acc.POST("/employees/:employeeId", handler.CreateAccount)
		// Status & resend per employee
		acc.GET("/employees/:employeeId", handler.GetAccountStatus)
		acc.POST("/employees/:employeeId/resend", handler.ResendSetupEmail)
	}
}
