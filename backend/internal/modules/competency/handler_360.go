package competency

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

// =========================================================================
// Rating Scale Handlers
// =========================================================================

func (h *Handler) CreateRatingScale(c *gin.Context) {
	var req CreateRatingScaleRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateRatingScale(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetRatingScaleByID(c *gin.Context) {
	resp, err := h.service.GetRatingScaleByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListRatingScales(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListRatingScales(c.Request.Context(), page, perPage, c.Query("status"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateRatingScale(c *gin.Context) {
	var req UpdateRatingScaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	resp, err := h.service.UpdateRatingScale(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteRatingScale(c *gin.Context) {
	if err := h.service.DeleteRatingScale(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Assessment Template Handlers
// =========================================================================

func (h *Handler) CreateAssessmentTemplate(c *gin.Context) {
	var req CreateAssessmentTemplateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateAssessmentTemplate(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetAssessmentTemplateByID(c *gin.Context) {
	resp, err := h.service.GetAssessmentTemplateByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListAssessmentTemplates(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListAssessmentTemplates(c.Request.Context(), page, perPage, c.Query("status"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateAssessmentTemplate(c *gin.Context) {
	var req UpdateAssessmentTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	resp, err := h.service.UpdateAssessmentTemplate(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteAssessmentTemplate(c *gin.Context) {
	if err := h.service.DeleteAssessmentTemplate(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Indicator Handlers
// =========================================================================

func (h *Handler) CreateIndicator(c *gin.Context) {
	var req CreateIndicatorRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateIndicator(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetIndicatorByID(c *gin.Context) {
	resp, err := h.service.GetIndicatorByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListIndicators(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListIndicators(c.Request.Context(), page, perPage, c.Query("competency_id"), c.Query("status"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateIndicator(c *gin.Context) {
	var req UpdateIndicatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	resp, err := h.service.UpdateIndicator(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteIndicator(c *gin.Context) {
	if err := h.service.DeleteIndicator(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Template Indicator Handlers
// =========================================================================

func (h *Handler) SetTemplateIndicators(c *gin.Context) {
	var req []TemplateIndicatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	resp, err := h.service.SetTemplateIndicators(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListTemplateIndicators(c *gin.Context) {
	resp, err := h.service.ListTemplateIndicators(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// =========================================================================
// Rater Assignment Handlers (§9)
// =========================================================================

func (h *Handler) AssignRaters(c *gin.Context) {
	var req AssignRatersRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.AssignRaters(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListRatersByTarget(c *gin.Context) {
	resp, err := h.service.ListRatersByTarget(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteRater(c *gin.Context) {
	if err := h.service.DeleteRater(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// My Assessment Handlers (§24 Employee)
// =========================================================================

func (h *Handler) MyAssessments(c *gin.Context) {
	resp, err := h.service.MyAssessments(c.Request.Context())
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) GetAssessmentDetail(c *gin.Context) {
	resp, err := h.service.GetAssessmentDetail(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) SaveResponses(c *gin.Context) {
	var req SaveResponsesRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.SaveResponses(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) SubmitAssessment(c *gin.Context) {
	resp, err := h.service.SubmitAssessment(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// =========================================================================
// Approval Integration Handlers (§13)
// =========================================================================

func (h *Handler) SubmitAssessmentForApproval(c *gin.Context) {
	resp, err := h.service.SubmitAssessmentForApproval(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// =========================================================================
// Result & Gap Handlers (§22 Result)
// =========================================================================

func (h *Handler) GetEmployeeResult(c *gin.Context) {
	resp, err := h.service.GetEmployeeResult(c.Request.Context(), c.Param("employee"), c.Query("event_id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) GetEmployeeGap(c *gin.Context) {
	resp, err := h.service.GetEmployeeGap(c.Request.Context(), c.Param("employee"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
