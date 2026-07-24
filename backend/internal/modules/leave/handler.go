package leave

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// =========================================================================
// Leave Types
// =========================================================================

func (h *Handler) CreateLeaveType(c *gin.Context) {
	var req CreateLeaveTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.CreateLeaveType(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
}

func (h *Handler) ListLeaveTypes(c *gin.Context) {
	page, perPage := parsePagination(c)
	resp, err := h.svc.ListLeaveTypes(c.Request.Context(), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetLeaveTypeByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetLeaveTypeByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "leave type not found" {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": "Leave type not found"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateLeaveType(c *gin.Context) {
	id := c.Param("id")
	var req UpdateLeaveTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.UpdateLeaveType(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "leave type not found" {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": "Leave type not found"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteLeaveType(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteLeaveType(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Leave type deleted"})
}

// =========================================================================
// Accrual Policies
// =========================================================================

func (h *Handler) CreateAccrualPolicy(c *gin.Context) {
	var req CreateAccrualPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.CreateAccrualPolicy(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
}

func (h *Handler) ListAccrualPolicies(c *gin.Context) {
	page, perPage := parsePagination(c)
	leaveTypeID := c.Query("leave_type_id")
	var lTypePtr *string
	if leaveTypeID != "" {
		lTypePtr = &leaveTypeID
	}
	resp, err := h.svc.ListAccrualPolicies(c.Request.Context(), lTypePtr, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetAccrualPolicyByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetAccrualPolicyByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "accrual policy not found" {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": "Accrual policy not found"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateAccrualPolicy(c *gin.Context) {
	id := c.Param("id")
	var req UpdateAccrualPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.UpdateAccrualPolicy(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteAccrualPolicy(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteAccrualPolicy(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Accrual policy deleted"})
}

// =========================================================================
// Leave Reasons
// =========================================================================

func (h *Handler) CreateLeaveReason(c *gin.Context) {
	var req CreateLeaveReasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.CreateLeaveReason(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
}

func (h *Handler) ListLeaveReasons(c *gin.Context) {
	reasons, err := h.svc.ListLeaveReasons(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": reasons})
}

func (h *Handler) GetLeaveReasonByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetLeaveReasonByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "leave reason not found" {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": "Leave reason not found"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateLeaveReason(c *gin.Context) {
	id := c.Param("id")
	var req UpdateLeaveReasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.UpdateLeaveReason(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteLeaveReason(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteLeaveReason(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Leave reason deleted"})
}

// =========================================================================
// Leave Requests
// =========================================================================

func (h *Handler) CreateLeaveRequest(c *gin.Context) {
	var req CreateLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.CreateLeaveRequest(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
}

func (h *Handler) ListLeaveRequests(c *gin.Context) {
	page, perPage := parsePagination(c)
	empID := c.Query("employee_id")
	status := c.Query("status")
	var empPtr *string
	if empID != "" {
		empPtr = &empID
	}
	resp, err := h.svc.ListLeaveRequests(c.Request.Context(), empPtr, &status, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetLeaveRequestByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetLeaveRequestByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "leave request not found" {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": "Leave request not found"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateLeaveRequestStatus(c *gin.Context) {
	id := c.Param("id")
	var req UpdateLeaveRequestStatus
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.UpdateLeaveRequestStatus(c.Request.Context(), id, req.Status, req.Note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteLeaveRequest(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteLeaveRequest(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Leave request deleted"})
}

// =========================================================================
// Leave Request Details
// =========================================================================

func (h *Handler) ListLeaveRequestDetails(c *gin.Context) {
	leaveRequestID := c.Param("id")
	details, err := h.svc.ListLeaveRequestDetails(c.Request.Context(), leaveRequestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": details})
}

// =========================================================================
// Leave Balances
// =========================================================================

func (h *Handler) ListLeaveBalances(c *gin.Context) {
	page, perPage := parsePagination(c)
	empID := c.Query("employee_id")
	var empPtr *string
	if empID != "" {
		empPtr = &empID
	}
	resp, err := h.svc.ListLeaveBalances(c.Request.Context(), empPtr, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetLeaveBalance(c *gin.Context) {
	empID := c.Param("employeeId")
	lTypeID := c.Param("leaveTypeId")
	year := time.Now().Year()
	if y := c.Query("year"); y != "" {
		if v, err := strconv.Atoi(y); err == nil && v > 0 {
			year = v
		}
	}
	resp, err := h.svc.GetLeaveBalance(c.Request.Context(), empID, lTypeID, year)
	if err != nil {
		if err.Error() == "leave balance not found" {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": "Leave balance not found"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// =========================================================================
// Helpers
// =========================================================================

func parsePagination(c *gin.Context) (int, int) {
	page := 1
	perPage := 20

	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if pp := c.Query("per_page"); pp != "" {
		if v, err := strconv.Atoi(pp); err == nil && v > 0 && v <= 100 {
			perPage = v
		}
	}
	return page, perPage
}
