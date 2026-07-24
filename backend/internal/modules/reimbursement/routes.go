package reimbursement

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan semua endpoint Reimbursement ke router group tenant.
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	reimb := rg.Group("/reimbursements")
	{
		// Reimbursement Types
		reimb.POST("/types", handler.CreateReimbursementType)
		reimb.GET("/types", handler.ListReimbursementTypes)
		reimb.GET("/types/:id", handler.GetReimbursementTypeByID)
		reimb.PUT("/types/:id", handler.UpdateReimbursementType)
		reimb.DELETE("/types/:id", handler.DeleteReimbursementType)

		// Reimbursement Requests
		reimb.POST("/requests", handler.CreateReimbursementRequest)
		reimb.GET("/requests", handler.ListReimbursementRequests)
		reimb.GET("/requests/:id", handler.GetReimbursementRequestByID)
		reimb.PUT("/requests/:id", handler.UpdateReimbursementRequest)
		reimb.PUT("/requests/:id/status", handler.UpdateReimbursementRequestStatus)
		reimb.DELETE("/requests/:id", handler.DeleteReimbursementRequest)

		// Reimbursement Items (nested under requests)
		reimb.GET("/requests/:id/items", handler.ListReimbursementItems)
		reimb.POST("/requests/:id/items", handler.CreateReimbursementItem)
		reimb.PUT("/requests/:id/items/:itemId", handler.UpdateReimbursementItem)
		reimb.DELETE("/requests/:id/items/:itemId", handler.DeleteReimbursementItem)
	}
}
