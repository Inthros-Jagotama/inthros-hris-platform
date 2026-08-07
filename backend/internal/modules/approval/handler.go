package approval

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

// Handler untuk HTTP endpoints Approval Engine.
type Handler struct {
	service *Service
}

// NewHandler membuat Handler baru.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// =========================================================================
// Approval Flow Handlers
// =========================================================================

// CreateFlow menangani POST /api/v1/tenant/approval/flows
func (h *Handler) CreateFlow(c *gin.Context) {
	var req CreateFlowRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.CreateFlow(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
		"message": "Approval flow created",
	})
}

// GetFlowByID menangani GET /api/v1/tenant/approval/flows/:flowId
// GetActiveFlowByModule resolves the active flow for a module (?module=xxx),
// used by consumers that auto-resolve their flow instead of picking a
// flow_id manually.
func (h *Handler) GetActiveFlowByModule(c *gin.Context) {
	module := c.Query("module")
	if module == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_PARAM",
				"message": "module query parameter is required",
			},
		})
		return
	}

	response, err := h.service.GetActiveFlowByModule(c.Request.Context(), module)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func (h *Handler) GetFlowByID(c *gin.Context) {
	id := c.Param("flowId")

	response, err := h.service.GetFlowByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// ListFlows menangani GET /api/v1/tenant/approval/flows
func (h *Handler) ListFlows(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	response, err := h.service.ListFlows(c.Request.Context(), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListAvailableModules menangani GET /api/v1/tenant/approval/available-modules
// Mengembalikan slug module yang aktif/disubscribe tenant — dipakai frontend
// agar flow builder hanya menampilkan module yang benar-benar tersedia.
func (h *Handler) ListAvailableModules(c *gin.Context) {
	modules, err := h.service.ListAvailableModules(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    modules,
	})
}

// UpdateFlow menangani PUT /api/v1/tenant/approval/flows/:flowId
func (h *Handler) UpdateFlow(c *gin.Context) {
	id := c.Param("flowId")

	var req UpdateFlowRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.UpdateFlow(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Approval flow updated",
	})
}

// DeleteFlow menangani DELETE /api/v1/tenant/approval/flows/:flowId
func (h *Handler) DeleteFlow(c *gin.Context) {
	id := c.Param("flowId")

	if err := h.service.DeleteFlow(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Approval flow deleted",
	})
}

// =========================================================================
// Approval Flow Steps Handlers
// =========================================================================

// CreateStep menangani POST /api/v1/tenant/approval/flows/:flowId/steps
func (h *Handler) CreateStep(c *gin.Context) {
	flowID := c.Param("flowId")

	var req CreateStepRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.CreateStep(c.Request.Context(), flowID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
		"message": "Approval step created",
	})
}

// ListSteps menangani GET /api/v1/tenant/approval/flows/:flowId/steps
func (h *Handler) ListSteps(c *gin.Context) {
	flowID := c.Param("flowId")

	responses, err := h.service.ListStepsByFlow(c.Request.Context(), flowID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    responses,
	})
}

// UpdateStep menangani PUT /api/v1/tenant/approval/flows/:flowId/steps/:stepId
func (h *Handler) UpdateStep(c *gin.Context) {
	flowID := c.Param("flowId")
	stepID := c.Param("stepId")

	var req UpdateStepRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.UpdateStep(c.Request.Context(), flowID, stepID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Approval step updated",
	})
}

// DeleteStep menangani DELETE /api/v1/tenant/approval/flows/:flowId/steps/:stepId
func (h *Handler) DeleteStep(c *gin.Context) {
	flowID := c.Param("flowId")
	stepID := c.Param("stepId")

	if err := h.service.DeleteStep(c.Request.Context(), flowID, stepID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Approval step deleted",
	})
}

// =========================================================================
// Approval Instance Handlers
// =========================================================================

// CreateInstance menangani POST /api/v1/tenant/approval/instances
func (h *Handler) CreateInstance(c *gin.Context) {
	var req CreateInstanceRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.CreateInstance(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
		"message": "Approval instance created",
	})
}

// GetInstanceByID menangani GET /api/v1/tenant/approval/instances/:id
func (h *Handler) GetInstanceByID(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.GetInstanceByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// ListInstances menangani GET /api/v1/tenant/approval/instances
func (h *Handler) ListInstances(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	module := c.Query("module")
	status := c.Query("status")

	response, err := h.service.ListInstances(c.Request.Context(), page, perPage, module, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CancelInstance menangani POST /api/v1/tenant/approval/instances/:id/cancel
func (h *Handler) CancelInstance(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.CancelInstance(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "CANCEL_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Approval instance cancelled",
	})
}

// =========================================================================
// Approval Action Handlers
// =========================================================================

// SubmitAction menangani POST /api/v1/tenant/approval/instances/:id/actions
func (h *Handler) SubmitAction(c *gin.Context) {
	id := c.Param("id")

	var req SubmitActionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	// Get user from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "user not authenticated",
			},
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "invalid user_id in context",
			},
		})
		return
	}

	response, err := h.service.SubmitAction(c.Request.Context(), id, userIDStr, req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "ACTION_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Approval action submitted",
	})
}

// =========================================================================
// Approval Task Handlers
// =========================================================================

// ListMyPendingTasks menangani GET /api/v1/tenant/approval/tasks/pending
func (h *Handler) ListMyPendingTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "user not authenticated",
			},
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "invalid user_id in context",
			},
		})
		return
	}

	response, err := h.service.ListMyPendingTasks(c.Request.Context(), userIDStr, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
