package employeemovement

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

// Handler untuk HTTP endpoints Employee Movement & Career Management.
type Handler struct {
	service *Service
}

// NewHandler membuat Handler baru.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// =========================================================================
// Employee Movement Handlers
// =========================================================================

// CreateMovement menangani POST /api/v1/tenant/employee-movements/movements
func (h *Handler) CreateMovement(c *gin.Context) {
	var req CreateMovementRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.CreateMovement(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, response, "success.created")
}

// GetMovementByID menangani GET /api/v1/tenant/employee-movements/movements/:id
func (h *Handler) GetMovementByID(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.GetMovementByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, response)
}

// ListMovements menangani GET /api/v1/tenant/employee-movements/movements
func (h *Handler) ListMovements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	response, err := h.service.ListMovements(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListMovementsByEmployee menangani GET /api/v1/tenant/employee-movements/employees/:employeeId/movements
func (h *Handler) ListMovementsByEmployee(c *gin.Context) {
	employeeID := c.Param("employeeId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	response, err := h.service.ListMovementsByEmployee(c.Request.Context(), employeeID, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateMovement menangani PUT /api/v1/tenant/employee-movements/movements/:id
func (h *Handler) UpdateMovement(c *gin.Context) {
	id := c.Param("id")

	var req UpdateMovementRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.UpdateMovement(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.UpdatedJSON(c, response, "success.updated")
}

// DeleteMovement menangani DELETE /api/v1/tenant/employee-movements/movements/:id
func (h *Handler) DeleteMovement(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteMovement(c.Request.Context(), id); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ApproveMovement menangani POST /api/v1/tenant/employee-movements/movements/:id/approve
func (h *Handler) ApproveMovement(c *gin.Context) {
	id := c.Param("id")

	// Get approver from JWT context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		httputil.ErrorRaw(c, http.StatusUnauthorized, "UNAUTHORIZED", "user not authenticated")
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

	if err := h.service.ApproveMovement(c.Request.Context(), id, userIDStr); err != nil {
		httputil.ErrorRaw(c, http.StatusConflict, "APPROVE_FAILED", err.Error())
		return
	}
	httputil.MessageJSON(c, "success.approved")
}

// ExecuteMovement menangani POST /api/v1/tenant/employee-movements/movements/:id/execute
func (h *Handler) ExecuteMovement(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		httputil.ErrorRaw(c, http.StatusUnauthorized, "UNAUTHORIZED", "user not authenticated")
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

	if err := h.service.ExecuteMovement(c.Request.Context(), id, userIDStr); err != nil {
		httputil.ErrorRaw(c, http.StatusConflict, "EXECUTE_FAILED", err.Error())
		return
	}
	httputil.MessageJSON(c, "success.executed")
}

// CancelMovement menangani POST /api/v1/tenant/employee-movements/movements/:id/cancel
func (h *Handler) CancelMovement(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.CancelMovement(c.Request.Context(), id); err != nil {
		httputil.ErrorRaw(c, http.StatusConflict, "CANCEL_FAILED", err.Error())
		return
	}
	httputil.MessageJSON(c, "success.cancelled")
}

// =========================================================================
// Employee Contract Handlers
// =========================================================================

// CreateContract menangani POST /api/v1/tenant/employee-movements/contracts
func (h *Handler) CreateContract(c *gin.Context) {
	var req CreateContractRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.CreateContract(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, response, "success.created")
}

// GetContractByID menangani GET /api/v1/tenant/employee-movements/contracts/:id
func (h *Handler) GetContractByID(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.GetContractByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, response)
}

// ListContracts menangani GET /api/v1/tenant/employee-movements/contracts
func (h *Handler) ListContracts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	response, err := h.service.ListContracts(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListContractsByEmployee menangani GET /api/v1/tenant/employee-movements/employees/:employeeId/contracts
func (h *Handler) ListContractsByEmployee(c *gin.Context) {
	employeeID := c.Param("employeeId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	response, err := h.service.ListContractsByEmployee(c.Request.Context(), employeeID, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateContract menangani PUT /api/v1/tenant/employee-movements/contracts/:id
func (h *Handler) UpdateContract(c *gin.Context) {
	id := c.Param("id")

	var req UpdateContractRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.UpdateContract(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.UpdatedJSON(c, response, "success.updated")
}

// DeleteContract menangani DELETE /api/v1/tenant/employee-movements/contracts/:id
func (h *Handler) DeleteContract(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteContract(c.Request.Context(), id); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}
