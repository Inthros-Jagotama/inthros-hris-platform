package rbac

import "time"

// ── Role DTOs ────────────────────────────────────────────────────────────────

type CreateRoleRequest struct {
	Name        string  `json:"name" binding:"required,max=255"`
	Description *string `json:"description"`
	IsDefault   *bool   `json:"is_default"`
}

type UpdateRoleRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=255"`
	Description *string `json:"description"`
	IsDefault   *bool   `json:"is_default"`
}

type RoleResponse struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	GuardName     string   `json:"guard_name"`
	Description   *string  `json:"description"`
	IsDefault     bool     `json:"is_default"`
	IsSystem      bool     `json:"is_system"`
	UserCount     int64    `json:"user_count"`
	PermissionIDs []string `json:"permission_ids"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

// ── Permission DTOs ──────────────────────────────────────────────────────────

type PermissionResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	GuardName string `json:"guard_name"`
	Resource  string `json:"resource"`
	Submenu   string `json:"submenu,omitempty"`
	Action    string `json:"action"`
}

// AssignPermissionsRequest untuk assign (replace) daftar permission ke role.
type AssignPermissionsRequest struct {
	PermissionIDs []string `json:"permission_ids" binding:"required"`
}

// ── User DTOs ────────────────────────────────────────────────────────────────

type UserResponse struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	IsActive  bool     `json:"is_active"`
	RoleIDs   []string `json:"role_ids"`
	RoleNames []string `json:"role_names"`
	CreatedAt string   `json:"created_at"`
}

// AssignUserRolesRequest untuk assign (replace) daftar role ke user tenant.
type AssignUserRolesRequest struct {
	RoleIDs []string `json:"role_ids" binding:"required"`
}

// helpers
func fmtTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func boolVal(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
