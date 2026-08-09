package notification

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) requireUserID(c *gin.Context) (uuid.UUID, bool) {
	userID := authctx.GetUserID(c.Request.Context())
	if userID == nil {
		httputil.ErrorSimple(c, http.StatusUnauthorized, "unauthorized")
		return uuid.Nil, false
	}
	return *userID, true
}

// ListNotifications returns the paginated notification feed for the logged-in user.
func (h *Handler) ListNotifications(c *gin.Context) {
	userID, ok := h.requireUserID(c)
	if !ok {
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	var isRead *bool
	if v := c.Query("is_read"); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			httputil.BadRequest(c, "is_read must be a boolean")
			return
		}
		isRead = &parsed
	}

	notifications, total, err := h.service.ListNotifications(c.Request.Context(), userID, isRead, page, perPage, httputil.GetLang(c))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, gin.H{
		"data":     notifications,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

// GetUnreadCount returns the unread notification count badge for the logged-in user.
func (h *Handler) GetUnreadCount(c *gin.Context) {
	userID, ok := h.requireUserID(c)
	if !ok {
		return
	}
	count, err := h.service.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, gin.H{"unread_count": count})
}

// MarkAsRead marks a single notification belonging to the logged-in user as read.
func (h *Handler) MarkAsRead(c *gin.Context) {
	userID, ok := h.requireUserID(c)
	if !ok {
		return
	}
	notificationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "invalid notification id")
		return
	}
	if err := h.service.MarkAsRead(c.Request.Context(), notificationID, userID); err != nil {
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.SuccessJSON(c, gin.H{"success": true})
}

// MarkAllAsRead marks every unread notification of the logged-in user as read.
func (h *Handler) MarkAllAsRead(c *gin.Context) {
	userID, ok := h.requireUserID(c)
	if !ok {
		return
	}
	if err := h.service.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, gin.H{"success": true})
}
