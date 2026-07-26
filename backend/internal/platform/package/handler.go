package pkgmgr

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

// Handler untuk HTTP endpoints Package Management.
type Handler struct {
	service *Service
}

// NewHandler membuat Handler baru.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreatePackage menangani POST /api/v1/platform/packages
func (h *Handler) CreatePackage(c *gin.Context) {
	var req CreatePackageRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.CreatePackage(req)
	if err != nil {
		httputil.ErrorRaw(c, 409, "CONFLICT", err.Error())
		return
	}

	httputil.CreatedJSON(c, response, "package.created")
}

// GetPackage menangani GET /api/v1/platform/packages/:id
func (h *Handler) GetPackage(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.GetPackage(id)
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, response)
}

// ListPackages menangani GET /api/v1/platform/packages
// Query params:
//   - page: halaman (default 1)
//   - per_page: item per halaman (default 20, max 100)
//   - module_type: filter paket berdasarkan tipe modul ("platform" atau "tenant"), kosongkan untuk semua
func (h *Handler) ListPackages(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	moduleType := c.Query("module_type")

	response, err := h.service.ListPackages(page, perPage, moduleType)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(200, response)
}

// UpdatePackage menangani PUT /api/v1/platform/packages/:id
func (h *Handler) UpdatePackage(c *gin.Context) {
	id := c.Param("id")

	var req UpdatePackageRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.UpdatePackage(id, req)
	if err != nil {
		httputil.ErrorRaw(c, 409, "CONFLICT", err.Error())
		return
	}

	httputil.UpdatedJSON(c, response, "package.updated")
}

// DeletePackage menangani DELETE /api/v1/platform/packages/:id
func (h *Handler) DeletePackage(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeletePackage(id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}

	httputil.DeletedJSON(c, "package.deleted")
}

// PublishPackage menangani POST /api/v1/platform/packages/:id/publish
func (h *Handler) PublishPackage(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.PublishPackage(id)
	if err != nil {
		httputil.ErrorRaw(c, 400, "VALIDATION_ERROR", err.Error())
		return
	}

	httputil.UpdatedJSON(c, response, "package.published")
}

// UnpublishPackage menangani POST /api/v1/platform/packages/:id/unpublish
func (h *Handler) UnpublishPackage(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.UnpublishPackage(id)
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}

	httputil.UpdatedJSON(c, response, "package.unpublished")
}

// ValidateDependencies menangani GET /api/v1/platform/packages/:id/validate
func (h *Handler) ValidateDependencies(c *gin.Context) {
	id := c.Param("id")

	deps, err := h.service.ValidatePackageDependencies(id)
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, deps)
}

// ListPublishedPackages menangani GET /api/v1/public/packages
// Endpoint publik tanpa autentikasi.
// Query params:
//   - module_type: filter paket berdasarkan tipe modul ("platform" atau "tenant"), kosongkan untuk semua
func (h *Handler) ListPublishedPackages(c *gin.Context) {
	moduleType := c.Query("module_type")

	packages, err := h.service.ListPublishedPackages(moduleType)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, packages)
}
