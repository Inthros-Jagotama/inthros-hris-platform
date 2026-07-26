package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

// Handler untuk HTTP endpoints User & Auth.
type Handler struct {
	service *Service
}

// NewHandler membuat Handler baru.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Login menangani POST /api/v1/platform/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.Login(req)
	if err != nil {
		httputil.ErrorRaw(c, 401, "UNAUTHORIZED", err.Error())
		return
	}

	httputil.SuccessJSON(c, response)
}

// RefreshToken menangani POST /api/v1/platform/refresh
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		httputil.ErrorRaw(c, 401, "UNAUTHORIZED", err.Error())
		return
	}

	httputil.SuccessJSON(c, response)
}

// ListUsers menangani GET /api/v1/platform/users
func (h *Handler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	response, err := h.service.ListUsers(page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(200, response)
}

// GetUser menangani GET /api/v1/platform/users/:id
func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.GetUser(id)
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, response)
}

// CreateUser menangani POST /api/v1/platform/users
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.CreateUser(req)
	if err != nil {
		httputil.ErrorRaw(c, 409, "CONFLICT", err.Error())
		return
	}

	httputil.CreatedJSON(c, response, "user.created")
}

// DeleteUser menangani DELETE /api/v1/platform/users/:id
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteUser(id); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.DeletedJSON(c, "user.deleted")
}

// UpdateUser menangani PUT /api/v1/platform/users/:id
func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req UpdateUserRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	response, err := h.service.UpdateUser(id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.UpdatedJSON(c, response, "user.updated")
}

// ChangePassword menangani PUT /api/v1/platform/users/:id/password
func (h *Handler) ChangePassword(c *gin.Context) {
	id := c.Param("id")

	var req ChangePasswordRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	if err := h.service.ChangePassword(id, req); err != nil {
		httputil.ErrorRaw(c, 400, "BAD_REQUEST", err.Error())
		return
	}

	httputil.UpdatedJSON(c, nil, "password.updated")
}
