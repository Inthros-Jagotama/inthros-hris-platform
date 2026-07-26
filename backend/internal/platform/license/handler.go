package license

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

// Handler untuk HTTP endpoints License Management.
type Handler struct {
	service *Service
}

// NewHandler membuat Handler baru.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateLicense menangani POST /api/v1/platform/licenses
func (h *Handler) CreateLicense(c *gin.Context) {
	var req CreateLicenseRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.CreateLicense(req)
	if err != nil {
		httputil.ErrorRaw(c, 409, "CONFLICT", err.Error())
		return
	}

	httputil.CreatedJSON(c, response, "license.created")
}

// GetLicense menangani GET /api/v1/platform/licenses/:id
func (h *Handler) GetLicense(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.GetLicense(id)
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, response)
}

// ListLicenses menangani GET /api/v1/platform/licenses
func (h *Handler) ListLicenses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	response, err := h.service.ListLicenses(page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(200, response)
}

// UpdateLicense menangani PUT /api/v1/platform/licenses/:id
func (h *Handler) UpdateLicense(c *gin.Context) {
	id := c.Param("id")

	var req UpdateLicenseRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.UpdateLicense(id, req)
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}

	httputil.UpdatedJSON(c, response, "license.updated")
}
