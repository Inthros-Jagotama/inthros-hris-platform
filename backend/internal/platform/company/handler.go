package company

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

// Handler untuk HTTP endpoints Company.
type Handler struct {
	service *Service
}

// NewHandler membuat Handler baru.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ResolveByHost menangani GET /api/v1/public/companies/resolve?host=...
// Endpoint publik (tanpa auth) untuk menentukan company dari hostname/subdomain
// URL aplikasi. Dipakai tenant FE sebelum login untuk prefill company.
func (h *Handler) ResolveByHost(c *gin.Context) {
	host := c.Query("host")
	resp, err := h.service.ResolveByHost(host)
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// Create menangani POST /api/v1/platform/companies
func (h *Handler) Create(c *gin.Context) {
	var req CreateCompanyRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.Create(req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.CreatedJSON(c, response, "company.created")
}

// GetByID menangani GET /api/v1/platform/companies/:id
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.GetByID(id)
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, response)
}

// List menangani GET /api/v1/platform/companies
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	response, err := h.service.List(page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(200, response)
}

// Update menangani PUT /api/v1/platform/companies/:id
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateCompanyRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.Update(id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.UpdatedJSON(c, response, "company.updated")
}

// Delete menangani DELETE /api/v1/platform/companies/:id
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.DeletedJSON(c, "company.deleted")
}

// Suspend menangani POST /api/v1/platform/companies/:id/suspend
func (h *Handler) Suspend(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.Suspend(id)
	if err != nil {
		httputil.ErrorRaw(c, 409, "SUSPEND_FAILED", err.Error())
		return
	}

	httputil.UpdatedJSON(c, response, "company.suspended")
}

// Activate menangani POST /api/v1/platform/companies/:id/activate
func (h *Handler) Activate(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.Activate(id)
	if err != nil {
		httputil.ErrorRaw(c, 409, "ACTIVATE_FAILED", err.Error())
		return
	}

	httputil.UpdatedJSON(c, response, "company.activated")
}

// Terminate menangani POST /api/v1/platform/companies/:id/terminate
func (h *Handler) Terminate(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.Terminate(id)
	if err != nil {
		httputil.ErrorRaw(c, 409, "TERMINATE_FAILED", err.Error())
		return
	}

	httputil.UpdatedJSON(c, response, "company.terminated")
}

// RotateCredentials menangani POST /api/v1/platform/companies/:id/rotate-credentials
func (h *Handler) RotateCredentials(c *gin.Context) {
	id := c.Param("id")

	var req RotateCredentialsRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.RotateCredentials(id, req.NewPassword)
	if err != nil {
		httputil.ErrorRaw(c, 409, "ROTATE_FAILED", err.Error())
		return
	}

	httputil.SuccessJSON(c, response)
}

// Backup menangani POST /api/v1/platform/companies/:id/backup
func (h *Handler) Backup(c *gin.Context) {
	id := c.Param("id")

	// TODO: Implement actual backup logic (Phase 2)
	_ = id

	httputil.MessageJSON(c, "success.backed_up")
}

// Restore menangani POST /api/v1/platform/companies/:id/restore
func (h *Handler) Restore(c *gin.Context) {
	id := c.Param("id")

	// TODO: Implement actual restore logic (Phase 2)
	_ = id

	httputil.MessageJSON(c, "success.restored_data")
}
