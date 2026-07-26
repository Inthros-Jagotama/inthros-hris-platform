package modulemgmt

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

// Handler untuk HTTP endpoints Module Management.
type Handler struct {
	service *Service
}

// NewHandler membuat Handler baru.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateModule menangani POST /api/v1/platform/modules
func (h *Handler) CreateModule(c *gin.Context) {
	var req CreateModuleRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.CreateModule(req)
	if err != nil {
		httputil.ErrorRaw(c, 409, "CONFLICT", err.Error())
		return
	}

	httputil.CreatedJSON(c, response, "module.created")
}

// GetModule menangani GET /api/v1/platform/modules/:id
func (h *Handler) GetModule(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.GetModule(id)
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, response)
}

// ListModules menangani GET /api/v1/platform/modules
// Query params:
//   - page: halaman (default 1)
//   - per_page: item per halaman (default 20, max 100)
//   - module_type: filter berdasarkan tipe ("platform" atau "tenant"), kosongkan untuk semua
func (h *Handler) ListModules(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	moduleType := c.Query("module_type")

	response, err := h.service.ListModules(page, perPage, moduleType)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(200, response)
}

// UpdateModule menangani PUT /api/v1/platform/modules/:id
func (h *Handler) UpdateModule(c *gin.Context) {
	id := c.Param("id")

	var req UpdateModuleRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.UpdateModule(id, req)
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}

	httputil.UpdatedJSON(c, response, "module.updated")
}

// ListCompanyModules menangani GET /api/v1/platform/modules/:id/companies
func (h *Handler) ListCompanyModules(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		httputil.ErrorRaw(c, 400, "VALIDATION_ERROR", "company_id query parameter is required")
		return
	}

	modules, err := h.service.ListCompanyModules(companyID)
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, modules)
}

// ActivateModule menangani POST /api/v1/platform/modules/:id/activate
func (h *Handler) ActivateModule(c *gin.Context) {
	id := c.Param("id")

	var req ToggleModuleRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.ActivateModule(id, req.CompanyID)
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}

	httputil.UpdatedJSON(c, response, "success.activated")
}

// DeactivateModule menangani POST /api/v1/platform/modules/:id/deactivate
func (h *Handler) DeactivateModule(c *gin.Context) {
	id := c.Param("id")

	var req ToggleModuleRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.DeactivateModule(id, req.CompanyID)
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}

	httputil.UpdatedJSON(c, response, "success.suspended")
}
