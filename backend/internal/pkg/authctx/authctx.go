// Package authctx menyediakan helper untuk mengekstrak informasi
// user dan authentication dari context.Context.
package authctx

import (
	"context"

	"github.com/google/uuid"
)

// GetUserID extracts the authenticated user ID from the request context.
// Returns nil if not found or invalid (e.g., system-generated/system actions).
// userId diset oleh middleware AuthJWT dan localize ke request context.
func GetUserID(ctx context.Context) *uuid.UUID {
	if userID, ok := ctx.Value("user_id").(string); ok && userID != "" {
		id, err := uuid.Parse(userID)
		if err == nil {
			return &id
		}
	}
	return nil
}

// GetCompanyID extracts the tenant company ID from the request context.
// Returns "" if not found. company_id diset oleh middleware.TenantRequired
// dan di-propagate ke request context.
func GetCompanyID(ctx context.Context) string {
	if companyID, ok := ctx.Value("company_id").(string); ok {
		return companyID
	}
	return ""
}

// GetPermissions extracts the caller's permission claims ("resource.action"
// strings) from the request context. Returns an empty slice if not found.
// permissions diset oleh middleware AuthJWT dan di-propagate oleh
// middleware.TenantRequired ke request context.
func GetPermissions(ctx context.Context) []string {
	if perms, ok := ctx.Value("permissions").([]string); ok {
		return perms
	}
	return []string{}
}

// HasPermission mengecek apakah caller punya permission "resource.action"
// tertentu, termasuk dukungan wildcard "*".
func HasPermission(ctx context.Context, resource, action string) bool {
	required := resource + "." + action
	for _, p := range GetPermissions(ctx) {
		if p == "*" || p == required {
			return true
		}
	}
	return false
}
