package employeemovement

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan semua endpoint Employee Movement ke router group tenant.
// Semua endpoint di bawah /api/v1/tenant/employee-movements
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	em := rg.Group("/employee-movements")
	{
		// Employee Movements
		em.POST("/movements", handler.CreateMovement)
		em.GET("/movements", handler.ListMovements)
		em.GET("/movements/:id", handler.GetMovementByID)
		em.PUT("/movements/:id", handler.UpdateMovement)
		em.DELETE("/movements/:id", handler.DeleteMovement)
		// Approval hanya lewat submit (Central Approval) — endpoint approve
		// manual dihapus (keputusan plan §11.5 / G-5).
		em.POST("/movements/:id/submit", handler.SubmitMovement)
		em.POST("/movements/:id/execute", handler.ExecuteMovement)
		em.POST("/movements/:id/cancel", handler.CancelMovement)
		// Audit trail movement (enhancement plan §12.6).
		em.GET("/movements/:id/audits", handler.ListMovementAudits)
		// Movement documents (enhancement plan §12.15) — metadata dokumen;
		// file fisik di-upload via endpoint upload generik.
		em.GET("/movements/:id/documents", handler.ListMovementDocuments)
		em.POST("/movements/:id/documents", handler.CreateMovementDocument)
		em.DELETE("/movements/:id/documents/:documentId", handler.DeleteMovementDocument)

		// Movements by Employee
		em.GET("/employees/:employeeId/movements", handler.ListMovementsByEmployee)

		// Career timeline read model (enhancement plan §12.8).
		em.GET("/employees/:employeeId/career-history", handler.GetCareerHistory)

		// Promotion eligibility (enhancement plan §12.10/§12.11).
		em.GET("/employees/:employeeId/movement-eligibility", handler.GetMovementEligibility)
		em.GET("/employees/:employeeId/promotion-eligibility", handler.GetPromotionEligibility)

		// Movement & Contract Reports (plan §12.17) — agregasi untuk HR reporting.
		em.GET("/reports/movements", handler.GetMovementReport)
		em.GET("/reports/contracts", handler.GetContractReport)
		// HR Dashboard summary (plan §12.18) — kartu dashboard module ini.
		em.GET("/dashboard", handler.GetHRDashboard)

		// Employee Contracts
		em.POST("/contracts", handler.CreateContract)
		em.GET("/contracts", handler.ListContracts)
		em.GET("/contracts/:id", handler.GetContractByID)
		em.PUT("/contracts/:id", handler.UpdateContract)
		em.DELETE("/contracts/:id", handler.DeleteContract)

		// Contracts by Employee
		em.GET("/employees/:employeeId/contracts", handler.ListContractsByEmployee)
	}

	// Career paths (enhancement plan §12.9) — planning/configuration jenjang
	// karier; endpoint per plan §15: /api/v1/tenant/career-paths.
	cp := rg.Group("/career-paths")
	{
		cp.GET("", handler.ListCareerPaths)
		cp.POST("", handler.CreateCareerPath)
		cp.GET("/:id", handler.GetCareerPathByID)
		cp.PUT("/:id", handler.UpdateCareerPath)
		cp.DELETE("/:id", handler.DeleteCareerPath)
	}
}
