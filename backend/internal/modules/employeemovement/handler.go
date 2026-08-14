package employeemovement

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/modules/approval"
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

func actorID(c *gin.Context) string {
	if v, ok := c.Get("user_id"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// =========================================================================
// Employee Movement Handlers
// =========================================================================

// movementErrStatus maps service errors to HTTP status codes: business
// validation failures (MovementValidationError, plan G-7) become 400, and
// employment-state conflicts (MovementConflictError, enhancement plan §12.3 /
// §12.4 — target position occupied or effective date overlap) become 409.
// Other errors fall through to the caller's default handling.
func movementErrStatus(err error) (int, bool) {
	var ve *MovementValidationError
	if errors.As(err, &ve) {
		return http.StatusBadRequest, true
	}
	var ce *MovementConflictError
	if errors.As(err, &ce) {
		return http.StatusConflict, true
	}
	return 0, false
}

// movementErrCode returns the error code that matches the mapped HTTP status
// (VALIDATION_ERROR for 400, CONFLICT_ERROR for 409).
func movementErrCode(status int) string {
	if status == http.StatusConflict {
		return "CONFLICT_ERROR"
	}
	return "VALIDATION_ERROR"
}

// CreateMovement menangani POST /api/v1/tenant/employee-movements/movements
func (h *Handler) CreateMovement(c *gin.Context) {
	var req CreateMovementRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.CreateMovement(c.Request.Context(), req)
	if err != nil {
		if status, ok := movementErrStatus(err); ok {
			httputil.ErrorRaw(c, status, movementErrCode(status), err.Error())
			return
		}
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
	movementType := c.Query("movement_type")
	status := c.Query("status")
	search := c.Query("search")

	response, err := h.service.ListMovements(c.Request.Context(), page, perPage, movementType, status, search)
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
		if status, ok := movementErrStatus(err); ok {
			httputil.ErrorRaw(c, status, movementErrCode(status), err.Error())
			return
		}
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

// SubmitMovement menangani POST /api/v1/tenant/employee-movements/movements/:id/submit
// Routes the movement through the central approval module.
func (h *Handler) SubmitMovement(c *gin.Context) {
	id := c.Param("id")
	var req SubmitMovementRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.SubmitMovement(c.Request.Context(), id, req)
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
	httputil.SuccessJSON(c, response)
}

// ExecuteMovement menangani POST /api/v1/tenant/employee-movements/movements/:id/execute
func (h *Handler) ExecuteMovement(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		httputil.ErrorJSON(c, http.StatusUnauthorized, "UNAUTHORIZED", "error.user_not_authenticated")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		httputil.ErrorJSON(c, http.StatusInternalServerError, "INTERNAL_ERROR", "error.invalid_user_context")
		return
	}

	if err := h.service.ExecuteMovement(c.Request.Context(), id, userIDStr); err != nil {
		// Employment-state conflicts (plan §12.3/§12.4) get their own error
		// code so the FE can distinguish "state conflict" from a generic
		// execution failure.
		if status, ok := movementErrStatus(err); ok {
			httputil.ErrorRaw(c, status, movementErrCode(status), err.Error())
			return
		}
		httputil.ErrorRaw(c, http.StatusConflict, "EXECUTE_FAILED", err.Error())
		return
	}
	httputil.MessageJSON(c, "success.executed")
}

// CancelMovement menangani POST /api/v1/tenant/employee-movements/movements/:id/cancel
// (plan §12.16 — Movement Cancellation Approval). Body opsional (flow_id /
// reason). Draft dibatalkan langsung; approved menjadi Cancellation Request
// yang diproses Central Approval (respond dengan success.cancellation_requested).
func (h *Handler) CancelMovement(c *gin.Context) {
	id := c.Param("id")

	var req CancelMovementRequest
	// Body opsional — jangan gagal bila kosong / bukan JSON (FE lama mengirim
	// body `{}`).
	_ = c.ShouldBindJSON(&req)

	response, err := h.service.CancelMovement(c.Request.Context(), id, req)
	if err != nil {
		// Cancellation requests route through the central approval module
		// too — routing failures surface bilingually, everything else keeps
		// the existing CANCEL_FAILED contract.
		if approval.EmitRoutingError(c, err) {
			return
		}
		httputil.ErrorRaw(c, http.StatusConflict, "CANCEL_FAILED", err.Error())
		return
	}
	if response.Status == string(MovementStatusCancellationPending) {
		httputil.MessageJSON(c, "success.cancellation_requested")
		return
	}
	httputil.MessageJSON(c, "success.cancelled")
}

// ListMovementAudits menangani GET /api/v1/tenant/employee-movements/movements/:id/audits
// (enhancement plan §12.6 — Movement Audit Trail).
func (h *Handler) ListMovementAudits(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	response, err := h.service.ListMovementAudits(c.Request.Context(), id, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetCareerHistory menangani GET /api/v1/tenant/employee-movements/employees/:employeeId/career-history
// (enhancement plan §12.8 — Career Timeline, read model dari movements +
// employments + contracts).
func (h *Handler) GetCareerHistory(c *gin.Context) {
	employeeID := c.Param("employeeId")

	response, err := h.service.GetCareerHistory(c.Request.Context(), employeeID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// =========================================================================
// Movement Document Handlers (plan §12.15)
// =========================================================================

// ListMovementDocuments menangani GET /api/v1/tenant/employee-movements/movements/:id/documents
func (h *Handler) ListMovementDocuments(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	response, err := h.service.ListMovementDocuments(c.Request.Context(), id, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateMovementDocument menangani POST /api/v1/tenant/employee-movements/movements/:id/documents
func (h *Handler) CreateMovementDocument(c *gin.Context) {
	id := c.Param("id")

	var req CreateMovementDocumentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.CreateMovementDocument(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, response, "success.created")
}

// DeleteMovementDocument menangani DELETE /api/v1/tenant/employee-movements/movements/:id/documents/:documentId
func (h *Handler) DeleteMovementDocument(c *gin.Context) {
	documentID := c.Param("documentId")

	if err := h.service.DeleteMovementDocument(c.Request.Context(), documentID); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Generate Document (plan §16/§17)
// =========================================================================

// GenerateMovementDocument menangani POST /movements/:id/generate-document
func (h *Handler) GenerateMovementDocument(c *gin.Context) {
	id := c.Param("id")
	doc, err := h.service.GenerateMovementDocument(c.Request.Context(), id, actorID(c))
	if err != nil {
		h.generateError(c, err)
		return
	}
	httputil.CreatedJSON(c, doc, "success.document_generated")
}

// ListGeneratedMovementDocuments menangani GET /movements/:id/generated-documents
func (h *Handler) ListGeneratedMovementDocuments(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	docs, total, err := h.service.ListGeneratedDocuments(c.Request.Context(), "movement", id, page, perPage)
	if err != nil {
		h.generateError(c, err)
		return
	}
	httputil.SuccessJSON(c, map[string]interface{}{
		"data": docs,
		"total": total,
	})
}

// GenerateContractDocument menangani POST /contracts/:id/generate-document
func (h *Handler) GenerateContractDocument(c *gin.Context) {
	id := c.Param("id")
	doc, err := h.service.GenerateContractDocument(c.Request.Context(), id, actorID(c))
	if err != nil {
		h.generateError(c, err)
		return
	}
	httputil.CreatedJSON(c, doc, "success.document_generated")
}

// ListGeneratedContractDocuments menangani GET /contracts/:id/generated-documents
func (h *Handler) ListGeneratedContractDocuments(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	docs, total, err := h.service.ListGeneratedDocuments(c.Request.Context(), "contract", id, page, perPage)
	if err != nil {
		h.generateError(c, err)
		return
	}
	httputil.SuccessJSON(c, map[string]interface{}{
		"data": docs,
		"total": total,
	})
}

// generateError memetakan error Generate Document ke response HTTP yang tepat.
func (h *Handler) generateError(c *gin.Context, err error) {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "document generator is not configured"):
		httputil.ErrorJSON(c, http.StatusServiceUnavailable, "PDF_ENGINE_NOT_CONFIGURED", "documenttemplate.pdf_engine_unavailable")
	case strings.Contains(msg, "no active template"):
		httputil.ErrorSimple(c, http.StatusBadRequest, msg)
	case strings.Contains(msg, "document can only be generated"):
		httputil.ErrorSimple(c, http.StatusBadRequest, msg)
	case strings.Contains(msg, "invalid movement id") || strings.Contains(msg, "invalid contract id") ||
		strings.Contains(msg, "movement not found") || strings.Contains(msg, "contract not found"):
		httputil.NotFound(c, msg)
	default:
		httputil.InternalError(c, msg)
	}
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
	status := c.Query("status")
	search := c.Query("search")

	response, err := h.service.ListContracts(c.Request.Context(), page, perPage, status, search)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetMovementReport menangani GET /api/v1/tenant/employee-movements/reports/movements
// (plan §12.17 — Movement Report dengan filter periode/org/posisi/employee/tipe/status).
func (h *Handler) GetMovementReport(c *gin.Context) {
	response, err := h.service.GetMovementReport(
		c.Request.Context(),
		c.Query("date_from"),
		c.Query("date_to"),
		c.Query("organization_id"),
		c.Query("position_id"),
		c.Query("employee_id"),
		c.Query("movement_type"),
		c.Query("status"),
	)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetContractReport menangani GET /api/v1/tenant/employee-movements/reports/contracts
// (plan §12.17 — Contract Report: by status + expiring).
func (h *Handler) GetContractReport(c *gin.Context) {
	response, err := h.service.GetContractReport(c.Request.Context())
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetHRDashboard menangani GET /api/v1/tenant/employee-movements/dashboard
// (plan §12.18 — kartu HR Dashboard: movement by type, pending approval,
// effective this month, ringkasan kontrak).
func (h *Handler) GetHRDashboard(c *gin.Context) {
	response, err := h.service.GetHRDashboard(c.Request.Context())
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

// =========================================================================
// Career Path Handlers (plan §12.9) — configuration jenjang karier
// =========================================================================

// CreateCareerPath menangani POST /api/v1/tenant/career-paths
// =========================================================================
// Promotion Eligibility Handlers (plan §12.10/§12.11)
// =========================================================================

// GetMovementEligibility menangani GET /api/v1/tenant/employee-movements/employees/:employeeId/movement-eligibility
func (h *Handler) GetMovementEligibility(c *gin.Context) {
	employeeID := c.Param("employeeId")

	response, err := h.service.GetMovementEligibility(c.Request.Context(), employeeID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetPromotionEligibility menangani GET /api/v1/tenant/employee-movements/employees/:employeeId/promotion-eligibility
func (h *Handler) GetPromotionEligibility(c *gin.Context) {
	employeeID := c.Param("employeeId")

	response, err := h.service.GetPromotionEligibility(c.Request.Context(), employeeID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
