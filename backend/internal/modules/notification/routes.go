package notification

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan semua endpoint Notification ke router group tenant.
// Semua endpoint di bawah /api/v1/tenant/notifications
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	n := rg.Group("/notifications")
	{
		n.GET("", handler.ListNotifications)
		n.GET("/unread-count", handler.GetUnreadCount)
		n.PATCH("/:id/read", handler.MarkAsRead)
		n.POST("/read-all", handler.MarkAllAsRead)
	}
}
