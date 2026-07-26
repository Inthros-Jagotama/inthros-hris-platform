package authz

import (
	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

// Handler untuk endpoint manajemen RBAC.
type Handler struct {
	service *Service
}

// NewHandler membuat Handler RBAC baru.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// =========================================================================
// Role Handlers
// =========================================================================

func (h *Handler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	resp, err := h.service.CreateRole(req)
	if err != nil {
		httputil.ErrorSimple(c, 500, err.Error())
		return
	}

	if err := h.service.Sync(); err != nil {
		httputil.ErrorSimple(c, 500, err.Error())
		return
	}

	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListRoles(c *gin.Context) {
	roles, err := h.service.ListRoles()
	if err != nil {
		httputil.ErrorSimple(c, 500, err.Error())
		return
	}

	httputil.SuccessJSON(c, roles)
}

func (h *Handler) GetRole(c *gin.Context) {
	resp, err := h.service.GetRole(c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateRole(c *gin.Context) {
	var req UpdateRoleRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	resp, err := h.service.UpdateRole(c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, 500, err.Error())
		return
	}

	if err := h.service.Sync(); err != nil {
		httputil.ErrorSimple(c, 500, err.Error())
		return
	}

	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteRole(c *gin.Context) {
	if err := h.service.DeleteRole(c.Param("id")); err != nil {
		httputil.ErrorSimple(c, 400, err.Error())
		return
	}

	if err := h.service.Sync(); err != nil {
		httputil.ErrorSimple(c, 500, err.Error())
		return
	}

	httputil.DeletedJSON(c, "rbac.role.deleted")
}

// =========================================================================
// Permission Handlers
// =========================================================================

func (h *Handler) CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	resp, err := h.service.CreatePermission(req)
	if err != nil {
		httputil.ErrorSimple(c, 500, err.Error())
		return
	}

	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPermissions(c *gin.Context) {
	perms, err := h.service.ListPermissions()
	if err != nil {
		httputil.ErrorSimple(c, 500, err.Error())
		return
	}

	httputil.SuccessJSON(c, perms)
}

func (h *Handler) DeletePermission(c *gin.Context) {
	if err := h.service.DeletePermission(c.Param("id")); err != nil {
		httputil.ErrorSimple(c, 400, err.Error())
		return
	}

	httputil.DeletedJSON(c, "rbac.permission.deleted")
}

// =========================================================================
// Role-Permission Assignment Handlers
// =========================================================================

func (h *Handler) AssignPermission(c *gin.Context) {
	var req AssignPermissionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	if err := h.service.AssignPermission(c.Param("id"), req.PermissionID); err != nil {
		httputil.ErrorSimple(c, 400, err.Error())
		return
	}

	if err := h.service.Sync(); err != nil {
		httputil.ErrorSimple(c, 500, err.Error())
		return
	}

	httputil.MessageJSON(c, "rbac.permission.assigned")
}

func (h *Handler) RevokePermission(c *gin.Context) {
	if err := h.service.RevokePermission(c.Param("id"), c.Param("permissionId")); err != nil {
		httputil.ErrorSimple(c, 400, err.Error())
		return
	}

	if err := h.service.Sync(); err != nil {
		httputil.ErrorSimple(c, 500, err.Error())
		return
	}

	httputil.MessageJSON(c, "rbac.permission.revoked")
}
