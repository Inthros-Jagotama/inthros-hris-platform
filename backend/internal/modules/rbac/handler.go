package rbac

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// handleConflict checks duplicate/system role errors and sends the
// appropriate bilingual response. Returns true if handled.
func handleConflict(c *gin.Context, err error) bool {
	var dupErr *DuplicateRoleError
	if errors.As(err, &dupErr) {
		httputil.ErrorJSON(c, 409, "DUPLICATE_ROLE", "rbac.duplicate_role", dupErr.Name)
		return true
	}
	var sysErr *SystemRoleError
	if errors.As(err, &sysErr) {
		httputil.ErrorJSON(c, 409, "SYSTEM_ROLE", "rbac.system_role", sysErr.Name)
		return true
	}
	return false
}

// ── Roles ──
func (h *Handler) ListRoles(c *gin.Context) {
	resp, err := h.service.ListRoles(c.Request.Context())
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateRole(c.Request.Context(), req)
	if err != nil {
		if handleConflict(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) UpdateRole(c *gin.Context) {
	var req UpdateRoleRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateRole(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if handleConflict(c, err) {
			return
		}
		if isNotFound(err) {
			httputil.NotFound(c, "")
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteRole(c *gin.Context) {
	err := h.service.DeleteRole(c.Request.Context(), c.Param("id"))
	if err != nil {
		if handleConflict(c, err) {
			return
		}
		if isNotFound(err) {
			httputil.NotFound(c, "")
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── Permissions ──
func (h *Handler) ListPermissions(c *gin.Context) {
	resp, err := h.service.ListPermissions(c.Request.Context())
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) AssignRolePermissions(c *gin.Context) {
	var req AssignPermissionsRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	if err := h.service.AssignRolePermissions(c.Request.Context(), c.Param("id"), req.PermissionIDs); err != nil {
		if isNotFound(err) {
			httputil.NotFound(c, "")
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.MessageJSON(c, "rbac.permissions_updated")
}

// ── Users & User-Role ──
func (h *Handler) ListUsers(c *gin.Context) {
	resp, err := h.service.ListUsers(c.Request.Context())
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) AssignUserRoles(c *gin.Context) {
	var req AssignUserRolesRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	if err := h.service.AssignUserRoles(c.Request.Context(), c.Param("id"), req.RoleIDs); err != nil {
		if isNotFound(err) {
		httputil.NotFound(c, "")
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.MessageJSON(c, "rbac.user_roles_updated")
}

// isNotFound detects gorm ErrRecordNotFound wrapped errors.
func isNotFound(err error) bool {
	return err != nil && errors.Is(err, gorm.ErrRecordNotFound)
}
