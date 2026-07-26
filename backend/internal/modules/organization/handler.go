package organization

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// POST /api/v1/tenant/organizations
func (h *Handler) Create(c *gin.Context) {
	var req CreateOrganizationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	resp, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
			httputil.InternalError(c, err.Error())
		return
	}

	httputil.CreatedJSON(c, resp, "success.created")
}

// GET /api/v1/tenant/organizations/:id
func (h *Handler) GetByID(c *gin.Context) {
	resp, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   gin.H{"code": "NOT_FOUND", "message": err.Error()},
		})
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GET /api/v1/tenant/organizations
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	// If tree=true, return tree structure
	if c.Query("tree") == "true" {
		tree, err := h.service.GetTree(c.Request.Context())
		if err != nil {
			httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
			return
		}
		httputil.SuccessJSON(c, tree)
		return
	}

	// Default: paginated flat list
	resp, err := h.service.List(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PUT /api/v1/tenant/organizations/:id
func (h *Handler) Update(c *gin.Context) {
	var req UpdateOrganizationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	resp, err := h.service.Update(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.SuccessJSON(c, resp)
}

// DELETE /api/v1/tenant/organizations/:id
func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// History Handlers
// =========================================================================

// GET /api/v1/tenant/organizations/history
// GET /api/v1/tenant/organizations/history?organization_id=xxx
func (h *Handler) ListHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	orgID := c.Query("organization_id")

	resp, err := h.service.GetHistory(c.Request.Context(), orgID, page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// =========================================================================
// Version Handlers
// =========================================================================

// POST /api/v1/tenant/organizations/versions
func (h *Handler) CreateVersion(c *gin.Context) {
	var req CreateVersionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	resp, err := h.service.CreateVersion(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.CreatedJSON(c, resp, "success.created")
}

// GET /api/v1/tenant/organizations/versions
func (h *Handler) ListVersions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	resp, err := h.service.ListVersions(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GET /api/v1/tenant/organizations/versions/:id
func (h *Handler) GetVersion(c *gin.Context) {
	includeSnapshot := c.Query("snapshot") == "true"
	resp, err := h.service.GetVersion(c.Request.Context(), c.Param("id"), includeSnapshot)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   gin.H{"code": "NOT_FOUND", "message": err.Error()},
		})
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GET /api/v1/tenant/organizations/versions/:sourceId/diff/:targetId
func (h *Handler) DiffVersions(c *gin.Context) {
	resp, err := h.service.DiffVersions(c.Request.Context(), c.Param("id"), c.Param("targetId"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// POST /api/v1/tenant/organizations/versions/:id/restore
func (h *Handler) RestoreVersion(c *gin.Context) {
	resp, err := h.service.RestoreVersion(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// =========================================================================
// Clone Handlers
// =========================================================================

// POST /api/v1/tenant/organizations/clone
func (h *Handler) CloneTree(c *gin.Context) {
	var req CloneRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	resp, err := h.service.CloneVersion(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.CreatedJSON(c, resp, "success.created")
}
