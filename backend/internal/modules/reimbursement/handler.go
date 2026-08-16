package reimbursement

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/modules/approval"
	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// =========================================================================
// Reimbursement Types
// =========================================================================

func (h *Handler) CreateReimbursementType(c *gin.Context) {
	var req CreateReimbursementTypeRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateReimbursementType(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListReimbursementTypes(c *gin.Context) {
	page, perPage := parsePagination(c)
	resp, err := h.svc.ListReimbursementTypes(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetReimbursementTypeByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetReimbursementTypeByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "reimbursement type not found" {
			httputil.NotFound(c, "")
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateReimbursementType(c *gin.Context) {
	id := c.Param("id")
	var req UpdateReimbursementTypeRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateReimbursementType(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "reimbursement type not found" {
			httputil.NotFound(c, "")
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteReimbursementType(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteReimbursementType(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Reimbursement Requests
// =========================================================================

func (h *Handler) CreateReimbursementRequest(c *gin.Context) {
	var req CreateReimbursementRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateReimbursementRequest(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListReimbursementRequests(c *gin.Context) {
	page, perPage := parsePagination(c)
	empID := c.Query("employee_id")
	status := c.Query("status")
	var empPtr *string
	if empID != "" {
		empPtr = &empID
	}
	resp, err := h.svc.ListReimbursementRequests(c.Request.Context(), empPtr, &status, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetReimbursementRequestByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetReimbursementRequestByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "reimbursement request not found" {
			httputil.NotFound(c, "")
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateReimbursementRequest(c *gin.Context) {
	id := c.Param("id")
	var req UpdateReimbursementRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateReimbursementRequest(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateReimbursementRequestStatus(c *gin.Context) {
	id := c.Param("id")
	var req UpdateReimbursementStatusRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateReimbursementRequestStatus(c.Request.Context(), id, req.Status, req.Note, req.Amount, req.FlowID, req.PaymentMethod, req.PaymentReference, req.PaymentNote)
	if err != nil {
		// Approval routing/assignee failures (approval.RoutingError) get a
		// bilingual 400 so the user sees why their submission didn't reach
		// an approver instead of a raw internal error.
		if approval.EmitRoutingError(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteReimbursementRequest(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteReimbursementRequest(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Reimbursement Items
// =========================================================================

func (h *Handler) CreateReimbursementItem(c *gin.Context) {
	requestID := c.Param("id")
	var req CreateReimbursementItemRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateReimbursementItem(c.Request.Context(), requestID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListReimbursementItems(c *gin.Context) {
	requestID := c.Param("id")
	items, err := h.svc.ListReimbursementItems(c.Request.Context(), requestID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) UpdateReimbursementItem(c *gin.Context) {
	requestID := c.Param("id")
	itemID := c.Param("itemId")
	var req UpdateReimbursementItemRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateReimbursementItem(c.Request.Context(), requestID, itemID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteReimbursementItem(c *gin.Context) {
	requestID := c.Param("id")
	itemID := c.Param("itemId")
	if err := h.svc.DeleteReimbursementItem(c.Request.Context(), requestID, itemID); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
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
