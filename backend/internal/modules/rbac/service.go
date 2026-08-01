package rbac

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// DuplicateRoleError dikembalikan saat nama role sudah dipakai tenant.
type DuplicateRoleError struct {
	Name string
}

func (e *DuplicateRoleError) Error() string {
	return fmt.Sprintf("role %q already exists", e.Name)
}

// SystemRoleError dikembalikan saat mencoba mengubah/menghapus role sistem.
type SystemRoleError struct {
	Name string
}

func (e *SystemRoleError) Error() string {
	return fmt.Sprintf("role %q is a system role and cannot be modified", e.Name)
}

// systemRoles adalah role bawaan hasil seeding yang dilindungi
// (tidak bisa dihapus, dan namanya tidak bisa diganti).
var systemRoles = map[string]bool{
	"Admin":    true,
	"Employee": true,
}

type Service struct {
	repo   *Repository
	logger *zap.Logger
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// ── Roles ────────────────────────────────────────────────────────────────────

func (s *Service) ListRoles(ctx context.Context) ([]RoleResponse, error) {
	roles, err := s.repo.ListRoles(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]RoleResponse, 0, len(roles))
	for _, role := range roles {
		permIDs, err := s.repo.ListRolePermissionIDs(ctx, role.ID)
		if err != nil {
			return nil, err
		}
		userCount, err := s.repo.CountRoleUsers(ctx, role.ID)
		if err != nil {
			return nil, err
		}
		resp = append(resp, RoleResponse{
			ID:            role.ID,
			Name:          role.Name,
			GuardName:     role.GuardName,
			Description:   role.Description,
			IsDefault:     boolVal(role.IsDefault),
			IsSystem:      systemRoles[role.Name],
			UserCount:     userCount,
			PermissionIDs: permIDs,
			CreatedAt:     fmtTime(role.CreatedAt),
			UpdatedAt:     fmtTime(role.UpdatedAt),
		})
	}
	return resp, nil
}

func (s *Service) CreateRole(ctx context.Context, req CreateRoleRequest) (*RoleResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("role name is required")
	}
	if _, err := s.repo.FindRoleByName(ctx, name); err == nil {
		return nil, &DuplicateRoleError{Name: name}
	}

	role := &Role{
		ID:          uuid.NewString(),
		Name:        name,
		GuardName:   "web",
		Description: req.Description,
		IsDefault:   req.IsDefault,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.repo.CreateRole(ctx, role); err != nil {
		return nil, err
	}
	return &RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		GuardName:   role.GuardName,
		Description: role.Description,
		IsDefault:   boolVal(role.IsDefault),
		IsSystem:    false,
		UserCount:   0,
		CreatedAt:   fmtTime(role.CreatedAt),
		UpdatedAt:   fmtTime(role.UpdatedAt),
	}, nil
}

func (s *Service) UpdateRole(ctx context.Context, id string, req UpdateRoleRequest) (*RoleResponse, error) {
	role, err := s.repo.FindRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Role sistem: nama tidak boleh diganti (Admin/Employee).
	if systemRoles[role.Name] && req.Name != nil && strings.TrimSpace(*req.Name) != role.Name {
		return nil, &SystemRoleError{Name: role.Name}
	}

	if req.Name != nil {
		newName := strings.TrimSpace(*req.Name)
		if newName == "" {
			return nil, errors.New("role name is required")
		}
		if newName != role.Name {
			if _, err := s.repo.FindRoleByName(ctx, newName); err == nil {
				return nil, &DuplicateRoleError{Name: newName}
			}
		}
		role.Name = newName
	}
	if req.Description != nil {
		role.Description = req.Description
	}
	if req.IsDefault != nil {
		role.IsDefault = req.IsDefault
	}
	role.UpdatedAt = time.Now()

	if err := s.repo.UpdateRole(ctx, role); err != nil {
		return nil, err
	}

	permIDs, err := s.repo.ListRolePermissionIDs(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	userCount, err := s.repo.CountRoleUsers(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	return &RoleResponse{
		ID:            role.ID,
		Name:          role.Name,
		GuardName:     role.GuardName,
		Description:   role.Description,
		IsDefault:     boolVal(role.IsDefault),
		IsSystem:      systemRoles[role.Name],
		UserCount:     userCount,
		PermissionIDs: permIDs,
		CreatedAt:     fmtTime(role.CreatedAt),
		UpdatedAt:     fmtTime(role.UpdatedAt),
	}, nil
}

func (s *Service) DeleteRole(ctx context.Context, id string) error {
	role, err := s.repo.FindRoleByID(ctx, id)
	if err != nil {
		return err
	}
	if systemRoles[role.Name] {
		return &SystemRoleError{Name: role.Name}
	}
	count, err := s.repo.CountRoleUsers(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("role %q is assigned to %d user(s) and cannot be deleted", role.Name, count)
	}
	return s.repo.DeleteRole(ctx, id)
}

// ── Permissions ──────────────────────────────────────────────────────────────

func (s *Service) ListPermissions(ctx context.Context) ([]PermissionResponse, error) {
	perms, err := s.repo.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]PermissionResponse, 0, len(perms))
	for _, p := range perms {
		resource, action := splitPermissionName(p.Name)
		resp = append(resp, PermissionResponse{
			ID:        p.ID,
			Name:      p.Name,
			GuardName: p.GuardName,
			Resource:  resource,
			Action:    action,
		})
	}
	return resp, nil
}

// AssignRolePermissions mengganti seluruh permission milik role.
func (s *Service) AssignRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	if _, err := s.repo.FindRoleByID(ctx, roleID); err != nil {
		return err
	}
	return s.repo.ReplaceRolePermissions(ctx, roleID, permissionIDs)
}

// ── Users & User-Role ────────────────────────────────────────────────────────

func (s *Service) ListUsers(ctx context.Context) ([]UserResponse, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	roles, err := s.repo.ListRoles(ctx)
	if err != nil {
		return nil, err
	}
	roleName := make(map[string]string, len(roles))
	for _, r := range roles {
		roleName[r.ID] = r.Name
	}

	resp := make([]UserResponse, 0, len(users))
	for _, u := range users {
		roleIDs, err := s.repo.ListUserRoleIDs(ctx, u.ID)
		if err != nil {
			return nil, err
		}
		roleNames := make([]string, 0, len(roleIDs))
		for _, rid := range roleIDs {
			if n, ok := roleName[rid]; ok {
				roleNames = append(roleNames, n)
			}
		}
		resp = append(resp, UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			IsActive:  u.IsActive != 0,
			RoleIDs:   roleIDs,
			RoleNames: roleNames,
			CreatedAt: fmtTime(u.CreatedAt),
		})
	}
	return resp, nil
}

func (s *Service) AssignUserRoles(ctx context.Context, userID string, roleIDs []string) error {
	if _, err := s.repo.FindUserByID(ctx, userID); err != nil {
		return err
	}
	return s.repo.ReplaceUserRoles(ctx, userID, roleIDs)
}

// splitPermissionName memecah "resource.action" menjadi resource & action.
func splitPermissionName(name string) (string, string) {
	idx := strings.LastIndex(name, ".")
	if idx < 0 {
		return name, ""
	}
	return name[:idx], name[idx+1:]
}
