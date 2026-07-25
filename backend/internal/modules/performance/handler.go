package performance

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.CreatePerformancePeriod(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
}

func (h *Handler) ListPerformancePeriods(c *gin.Context) {
	page, perPage := parsePagination(c)
	resp, err := h.svc.ListPerformancePeriods(c.Request.Context(), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetPerformancePeriodByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetPerformancePeriodByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": "Performance period not found"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdatePerformancePeriod(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePerformancePeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.UpdatePerformancePeriod(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeletePerformancePeriod(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePerformancePeriod(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Performance period deleted"})
}

// =========================================================================
// Performance Perspectives
// =========================================================================

func (h *Handler) CreatePerformancePerspective(c *gin.Context) {
	var req CreatePerformancePerspectiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.CreatePerformancePerspective(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
}

func (h *Handler) ListPerformancePerspectives(c *gin.Context) {
	page, perPage := parsePagination(c)
	resp, err := h.svc.ListPerformancePerspectives(c.Request.Context(), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetPerformancePerspectiveByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetPerformancePerspectiveByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": "Performance perspective not found"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdatePerformancePerspective(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePerformancePerspectiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.UpdatePerformancePerspective(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeletePerformancePerspective(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePerformancePerspective(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Performance perspective deleted"})
}

// =========================================================================
// Performance Templates
// =========================================================================

func (h *Handler) CreatePerformanceTemplate(c *gin.Context) {
	var req CreatePerformanceTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.CreatePerformanceTemplate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
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
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetPerformanceTemplateByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetPerformanceTemplateByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": "Performance template not found"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdatePerformanceTemplate(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePerformanceTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.UpdatePerformanceTemplate(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeletePerformanceTemplate(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePerformanceTemplate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Performance template deleted"})
}

// =========================================================================
// Performance Indicators
// =========================================================================

func (h *Handler) CreatePerformanceIndicator(c *gin.Context) {
	var req CreatePerformanceIndicatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.CreatePerformanceIndicator(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
}

func (h *Handler) ListPerformanceIndicators(c *gin.Context) {
	page, perPage := parsePagination(c)
	templateID := c.Query("template_id")
	if templateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": "template_id query parameter is required"}})
		return
	}
	resp, err := h.svc.ListPerformanceIndicators(c.Request.Context(), templateID, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetPerformanceIndicatorByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetPerformanceIndicatorByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": "Performance indicator not found"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdatePerformanceIndicator(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePerformanceIndicatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.UpdatePerformanceIndicator(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeletePerformanceIndicator(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePerformanceIndicator(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Performance indicator deleted"})
}

// =========================================================================
// Performance Evaluations
// =========================================================================

func (h *Handler) CreatePerformanceEvaluation(c *gin.Context) {
	var req CreatePerformanceEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.CreatePerformanceEvaluation(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
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
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetPerformanceEvaluationByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetPerformanceEvaluationByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": "Performance evaluation not found"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdatePerformanceEvaluation(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePerformanceEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.UpdatePerformanceEvaluation(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateEvaluationStatus(c *gin.Context) {
	id := c.Param("id")
	var req UpdateEvaluationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.UpdateEvaluationStatus(c.Request.Context(), id, req.Status, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeletePerformanceEvaluation(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePerformanceEvaluation(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Performance evaluation deleted"})
}

// =========================================================================
// Evaluation Details
// =========================================================================

func (h *Handler) CreateEvaluationDetail(c *gin.Context) {
	var req CreateEvaluationDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.CreateEvaluationDetail(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
}

func (h *Handler) ListEvaluationDetails(c *gin.Context) {
	evalID := c.Param("id")
	items, err := h.svc.ListEvaluationDetails(c.Request.Context(), evalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) UpdateEvaluationDetail(c *gin.Context) {
	id := c.Param("id")
	var req UpdateEvaluationDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.UpdateEvaluationDetail(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteEvaluationDetail(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteEvaluationDetail(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Evaluation detail deleted"})
}

// =========================================================================
// Performance Targets
// =========================================================================

func (h *Handler) CreatePerformanceTarget(c *gin.Context) {
	var req CreatePerformanceTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.CreatePerformanceTarget(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
}

func (h *Handler) ListPerformanceTargets(c *gin.Context) {
	evalID := c.Param("id")
	items, err := h.svc.ListPerformanceTargets(c.Request.Context(), evalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) UpdatePerformanceTarget(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePerformanceTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.UpdatePerformanceTarget(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeletePerformanceTarget(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePerformanceTarget(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Performance target deleted"})
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
