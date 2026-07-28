package organization

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

// POST /api/v1/tenant/summaries
func (h *Handler) CreateSummary(c *gin.Context) {
	var req CreateOrganizationSummaryRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	resp, err := h.service.CreateSummary(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.CreatedJSON(c, resp, "success.created")
}

// GET /api/v1/tenant/summaries/:id
func (h *Handler) GetSummaryByID(c *gin.Context) {
	resp, err := h.service.GetSummaryByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   gin.H{"code": "NOT_FOUND", "message": err.Error()},
		})
		return
	}

	httputil.SuccessJSON(c, resp)
}

// GET /api/v1/tenant/summaries
func (h *Handler) ListSummaries(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	resp, err := h.service.ListSummaries(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PUT /api/v1/tenant/summaries/:id
func (h *Handler) UpdateSummary(c *gin.Context) {
	var req UpdateOrganizationSummaryRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	resp, err := h.service.UpdateSummary(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.SuccessJSON(c, resp)
}

// DELETE /api/v1/tenant/summaries/:id
func (h *Handler) DeleteSummary(c *gin.Context) {
	if err := h.service.DeleteSummary(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.DeletedJSON(c, "success.deleted")
}

// GET /api/v1/tenant/summaries/stats
func (h *Handler) GetSummaryStats(c *gin.Context) {
	stats, err := h.service.GetSummaryStats(c.Request.Context())
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}
