package performance

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
// Performance Periods
// =========================================================================

func (h *Handler) CreatePerformancePeriod(c *gin.Context) {
	var req CreatePerformancePeriodRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreatePerformancePeriod(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPerformancePeriods(c *gin.Context) {
	page, perPage := parsePagination(c)
	resp, err := h.svc.ListPerformancePeriods(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetPerformancePeriodByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetPerformancePeriodByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdatePerformancePeriod(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePerformancePeriodRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdatePerformancePeriod(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeletePerformancePeriod(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePerformancePeriod(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Performance Perspectives
// =========================================================================

func (h *Handler) CreatePerformancePerspective(c *gin.Context) {
	var req CreatePerformancePerspectiveRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreatePerformancePerspective(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPerformancePerspectives(c *gin.Context) {
	page, perPage := parsePagination(c)
	resp, err := h.svc.ListPerformancePerspectives(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetPerformancePerspectiveByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetPerformancePerspectiveByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdatePerformancePerspective(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePerformancePerspectiveRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdatePerformancePerspective(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeletePerformancePerspective(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePerformancePerspective(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Performance Templates
// =========================================================================

func (h *Handler) CreatePerformanceTemplate(c *gin.Context) {
	var req CreatePerformanceTemplateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreatePerformanceTemplate(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPerformanceTemplates(c *gin.Context) {
	page, perPage := parsePagination(c)
	orgID := c.Query("organization_id")
	var orgPtr *string
	if orgID != "" {
		orgPtr = &orgID
	}
	resp, err := h.svc.ListPerformanceTemplates(c.Request.Context(), orgPtr, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetPerformanceTemplateByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetPerformanceTemplateByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) GetMyKPIContext(c *gin.Context) {
	resp, err := h.svc.GetMyKPIContext(c.Request.Context())
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListTemplateOrganizationScope(c *gin.Context) {
	resp, err := h.svc.ListTemplateOrganizationScope(c.Request.Context())
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) RecalculatePeriodScoring(c *gin.Context) {
	periodID := c.Param("period_id")
	resp, err := h.svc.RecalculatePeriodScoring(c.Request.Context(), periodID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdatePerformanceTemplate(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePerformanceTemplateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdatePerformanceTemplate(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeletePerformanceTemplate(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePerformanceTemplate(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) DuplicatePerformanceTemplate(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.DuplicatePerformanceTemplate(c.Request.Context(), id)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

// =========================================================================
// Performance Indicators
// =========================================================================

func (h *Handler) CreatePerformanceIndicator(c *gin.Context) {
	var req CreatePerformanceIndicatorRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreatePerformanceIndicator(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPerformanceIndicators(c *gin.Context) {
	page, perPage := parsePagination(c)
	templateID := c.Query("template_id")
	if templateID == "" {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "performance.template_param_required")
		return
	}
	resp, err := h.svc.ListPerformanceIndicators(c.Request.Context(), templateID, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetPerformanceIndicatorByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetPerformanceIndicatorByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdatePerformanceIndicator(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePerformanceIndicatorRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdatePerformanceIndicator(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeletePerformanceIndicator(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePerformanceIndicator(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Performance Evaluations
// =========================================================================

func (h *Handler) CreatePerformanceEvaluation(c *gin.Context) {
	var req CreatePerformanceEvaluationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreatePerformanceEvaluation(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPerformanceEvaluations(c *gin.Context) {
	page, perPage := parsePagination(c)
	empID := c.Query("employee_id")
	periodID := c.Query("period_id")
	status := c.Query("status")
	var empPtr, perPtr *string
	if empID != "" {
		empPtr = &empID
	}
	if periodID != "" {
		perPtr = &periodID
	}
	resp, err := h.svc.ListPerformanceEvaluations(c.Request.Context(), empPtr, perPtr, &status, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetPerformanceEvaluationByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetPerformanceEvaluationByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdatePerformanceEvaluation(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePerformanceEvaluationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdatePerformanceEvaluation(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateEvaluationStatus(c *gin.Context) {
	id := c.Param("id")
	var req UpdateEvaluationStatusRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateEvaluationStatus(c.Request.Context(), id, req.Status, req.Notes)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeletePerformanceEvaluation(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePerformanceEvaluation(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Evaluation Details
// =========================================================================

func (h *Handler) CreateEvaluationDetail(c *gin.Context) {
	var req CreateEvaluationDetailRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateEvaluationDetail(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListEvaluationDetails(c *gin.Context) {
	evalID := c.Param("id")
	items, err := h.svc.ListEvaluationDetails(c.Request.Context(), evalID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) UpdateEvaluationDetail(c *gin.Context) {
	id := c.Param("id")
	var req UpdateEvaluationDetailRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateEvaluationDetail(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteEvaluationDetail(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteEvaluationDetail(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Performance Targets
// =========================================================================

func (h *Handler) CreatePerformanceTarget(c *gin.Context) {
	var req CreatePerformanceTargetRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreatePerformanceTarget(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPerformanceTargets(c *gin.Context) {
	evalID := c.Param("id")
	items, err := h.svc.ListPerformanceTargets(c.Request.Context(), evalID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) UpdatePerformanceTarget(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePerformanceTargetRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdatePerformanceTarget(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeletePerformanceTarget(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePerformanceTarget(c.Request.Context(), id); err != nil {
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

// =========================================================================
// Performance Progress
// =========================================================================

func (h *Handler) CreatePerformanceProgress(c *gin.Context) {
	var req CreatePerformanceProgressRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreatePerformanceProgress(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPerformanceProgressByDetailID(c *gin.Context) {
	detailID := c.Param("id")
	items, err := h.svc.ListPerformanceProgressByDetailID(c.Request.Context(), detailID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) GetPerformanceProgressByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetPerformanceProgressByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdatePerformanceProgress(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePerformanceProgressRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdatePerformanceProgress(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeletePerformanceProgress(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePerformanceProgress(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Performance Comments
// =========================================================================

func (h *Handler) CreatePerformanceComment(c *gin.Context) {
	var req CreatePerformanceCommentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreatePerformanceComment(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPerformanceCommentsByEvaluationID(c *gin.Context) {
	evalID := c.Param("id")
	items, err := h.svc.ListPerformanceCommentsByEvaluationID(c.Request.Context(), evalID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) GetPerformanceCommentByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetPerformanceCommentByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdatePerformanceComment(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePerformanceCommentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdatePerformanceComment(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeletePerformanceComment(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePerformanceComment(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Performance Attachments
// =========================================================================

func (h *Handler) CreatePerformanceAttachment(c *gin.Context) {
	var req CreatePerformanceAttachmentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreatePerformanceAttachment(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPerformanceAttachmentsByDetailID(c *gin.Context) {
	detailID := c.Param("id")
	items, err := h.svc.ListPerformanceAttachmentsByDetailID(c.Request.Context(), detailID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) GetPerformanceAttachmentByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetPerformanceAttachmentByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdatePerformanceAttachment(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePerformanceAttachmentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdatePerformanceAttachment(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeletePerformanceAttachment(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePerformanceAttachment(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Performance Ratings
// =========================================================================

func (h *Handler) CreatePerformanceRating(c *gin.Context) {
	var req CreatePerformanceRatingRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreatePerformanceRating(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPerformanceRatings(c *gin.Context) {
	page, perPage := parsePagination(c)
	resp, err := h.svc.ListPerformanceRatings(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetPerformanceRatingByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetPerformanceRatingByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdatePerformanceRating(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePerformanceRatingRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdatePerformanceRating(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeletePerformanceRating(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePerformanceRating(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Performance Indicator Formulas
// =========================================================================

func (h *Handler) CreatePerformanceIndicatorFormula(c *gin.Context) {
	var req CreatePerformanceIndicatorFormulaRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreatePerformanceIndicatorFormula(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPerformanceIndicatorFormulas(c *gin.Context) {
	page, perPage := parsePagination(c)
	resp, err := h.svc.ListPerformanceIndicatorFormulas(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetPerformanceIndicatorFormulaByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetPerformanceIndicatorFormulaByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdatePerformanceIndicatorFormula(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePerformanceIndicatorFormulaRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdatePerformanceIndicatorFormula(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeletePerformanceIndicatorFormula(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePerformanceIndicatorFormula(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Performance Logs (Read-only)
// =========================================================================

func (h *Handler) ListPerformanceLogs(c *gin.Context) {
	page, perPage := parsePagination(c)
	resp, err := h.svc.ListPerformanceLogs(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ListPerformanceLogsByEvaluationID(c *gin.Context) {
	evalID := c.Param("id")
	page, perPage := parsePagination(c)
	resp, err := h.svc.ListPerformanceLogsByEvaluationID(c.Request.Context(), evalID, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetPerformanceLogByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetPerformanceLogByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

// =========================================================================
// Phase 3 - Business Process Handlers
// =========================================================================

// CreateEvaluationWithSnapshot creates evaluation and snapshots KPIs from template
func (h *Handler) CreateEvaluationWithSnapshot(c *gin.Context) {
	var req CreateEvaluationWithSnapshotRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateEvaluationWithSnapshot(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

// GetEvaluationWithDetails returns evaluation with all details
func (h *Handler) GetEvaluationWithDetails(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetEvaluationWithDetails(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

// UpdateEvaluationActual updates actual value with auto calculation
func (h *Handler) UpdateEvaluationTarget(c *gin.Context) {
	id := c.Param("id")
	var req UpdateEvaluationTargetRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateEvaluationTarget(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateEvaluationActual(c *gin.Context) {
	id := c.Param("id")
	var req UpdateEvaluationActualRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateEvaluationActual(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// =========================================================================
// Performance Evaluation Program Items
// =========================================================================

func (h *Handler) CreateProgramItem(c *gin.Context) {
	var req CreateProgramItemRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateProgramItem(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListProgramItems(c *gin.Context) {
	evalID := c.Param("id")
	items, err := h.svc.ListProgramItems(c.Request.Context(), evalID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) UpdateProgramItemTarget(c *gin.Context) {
	id := c.Param("id")
	var req UpdateProgramItemTargetRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateProgramItemTarget(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateProgramItemActual(c *gin.Context) {
	id := c.Param("id")
	var req UpdateProgramItemActualRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateProgramItemActual(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteProgramItem(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteProgramItem(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// BulkUpdateEvaluationActuals updates multiple details with auto calculation
func (h *Handler) BulkUpdateEvaluationActuals(c *gin.Context) {
	var req BulkUpdateEvaluationActualRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.BulkUpdateEvaluationActuals(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// RecalculateEvaluationScore recalculates final score and rating
func (h *Handler) RecalculateEvaluationScore(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.RecalculateEvaluationScore(c.Request.Context(), id); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	// Return updated evaluation
	resp, err := h.svc.GetEvaluationWithDetails(c.Request.Context(), id)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetEvaluationProgressSummary returns progress summary
func (h *Handler) GetEvaluationProgressSummary(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetEvaluationProgressSummary(c.Request.Context(), id)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// SubmitEvaluation changes status to SUBMITTED
func (h *Handler) SubmitTarget(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.SubmitTarget(c.Request.Context(), id)
	if err != nil {
		emitKPIStatusError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ApproveTarget(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.ApproveTarget(c.Request.Context(), id)
	if err != nil {
		emitKPIStatusError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) RejectTarget(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Notes *string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Notes is optional, so ignore binding error
	}
	resp, err := h.svc.RejectTarget(c.Request.Context(), id, req.Notes)
	if err != nil {
		emitKPIStatusError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) SubmitEvaluation(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.SubmitEvaluation(c.Request.Context(), id)
	if err != nil {
		emitKPIStatusError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// ApproveEvaluation changes status to APPROVED
func (h *Handler) ApproveEvaluation(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.ApproveEvaluation(c.Request.Context(), id)
	if err != nil {
		emitKPIStatusError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// RejectEvaluation changes status back to DRAFT
func (h *Handler) RejectEvaluation(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Notes *string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Notes is optional, so ignore binding error
	}
	resp, err := h.svc.RejectEvaluation(c.Request.Context(), id, req.Notes)
	if err != nil {
		emitKPIStatusError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// CompleteEvaluation changes status to COMPLETED
func (h *Handler) CompleteEvaluation(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.CompleteEvaluation(c.Request.Context(), id)
	if err != nil {
		emitKPIStatusError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// emitKPIStatusError responds with a bilingual error for KPI status/approval
// actions. Approval routing/assignee-resolution failures surface as
// approval.RoutingError and are emitted via the shared approval.EmitRoutingError
// helper (400 APPROVAL_ROUTING_FAILED, message follows the request language);
// anything else falls back to the raw error with the legacy INVALID_STATUS
// code, keeping existing behavior unchanged.
func emitKPIStatusError(c *gin.Context, err error) {
	if approval.EmitRoutingError(c, err) {
		return
	}
	httputil.ErrorRaw(c, http.StatusBadRequest, "INVALID_STATUS", err.Error())
}

// =========================================================================
// Phase 4 - Dashboard Handlers
// =========================================================================

// GetEmployeeDashboard returns dashboard data for an employee
func (h *Handler) GetEmployeeDashboard(c *gin.Context) {
	employeeID := c.Param("employee_id")
	if employeeID == "" {
		httputil.ErrorJSON(c, http.StatusBadRequest, "INVALID_PARAM", "performance.employee_id_required")
		return
	}

	periodID := c.Query("period_id")
	var pID *string
	if periodID != "" {
		pID = &periodID
	}

	resp, err := h.svc.GetEmployeeDashboard(c.Request.Context(), employeeID, pID)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusInternalServerError, "DASHBOARD_ERROR", err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetManagerDashboard returns dashboard data for a manager
func (h *Handler) GetManagerDashboard(c *gin.Context) {
	managerID := c.Param("manager_id")
	if managerID == "" {
		httputil.ErrorJSON(c, http.StatusBadRequest, "INVALID_PARAM", "performance.manager_id_required")
		return
	}

	periodID := c.Query("period_id")
	var pID *string
	if periodID != "" {
		pID = &periodID
	}

	resp, err := h.svc.GetManagerDashboard(c.Request.Context(), managerID, pID)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusInternalServerError, "DASHBOARD_ERROR", err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetHRDashboard returns dashboard data for HR
func (h *Handler) GetHRDashboard(c *gin.Context) {
	periodID := c.Query("period_id")
	var pID *string
	if periodID != "" {
		pID = &periodID
	}

	resp, err := h.svc.GetHRDashboard(c.Request.Context(), pID)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusInternalServerError, "DASHBOARD_ERROR", err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// =========================================================================
// Performance Components (Phase 5 - Scoring Configuration)
// =========================================================================

func (h *Handler) CreatePerformanceComponent(c *gin.Context) {
	var req CreatePerformanceComponentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreatePerformanceComponent(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPerformanceComponents(c *gin.Context) {
	page, perPage := parsePagination(c)
	resp, err := h.svc.ListPerformanceComponents(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetPerformanceComponentByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetPerformanceComponentByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdatePerformanceComponent(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePerformanceComponentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdatePerformanceComponent(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeletePerformanceComponent(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePerformanceComponent(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Performance Organization Components
// =========================================================================

func (h *Handler) UpsertOrganizationComponent(c *gin.Context) {
	var req UpsertOrganizationComponentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpsertOrganizationComponent(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListOrganizationComponents(c *gin.Context) {
	orgID := c.Param("organization_id")
	resp, err := h.svc.ListOrganizationComponents(c.Request.Context(), orgID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteOrganizationComponent(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteOrganizationComponent(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Performance Evaluation Components (Scoring Engine)
// =========================================================================

func (h *Handler) ListEvaluationComponents(c *gin.Context) {
	evalID := c.Param("id")
	resp, err := h.svc.ListEvaluationComponents(c.Request.Context(), evalID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) CalculateEvaluationComponentScoring(c *gin.Context) {
	evalID := c.Param("id")
	resp, err := h.svc.CalculateEvaluationComponentScoring(c.Request.Context(), evalID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateEvaluationComponentScore(c *gin.Context) {
	evalID := c.Param("id")
	componentID := c.Param("component_id")
	var req UpdateEvaluationComponentScoreRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateEvaluationComponentScore(c.Request.Context(), evalID, componentID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
