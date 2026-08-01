package rbac

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Repository struct {
	dbResolver func(ctx context.Context) (*gorm.DB, error)
}

func NewRepository(dbResolver func(ctx context.Context) (*gorm.DB, error)) *Repository {
	return &Repository{dbResolver: dbResolver}
}

func (r *Repository) getDB(ctx context.Context) (*gorm.DB, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required for tenant database resolution")
	}
	return r.dbResolver(ctx)
}

// ── Roles ────────────────────────────────────────────────────────────────────

func (r *Repository) CreateRole(ctx context.Context, role *Role) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(role).Error
}

func (r *Repository) FindRoleByID(ctx context.Context, id string) (*Role, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var role Role
	if err := db.First(&role, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}
	return &role, nil
}

func (r *Repository) FindRoleByName(ctx context.Context, name string) (*Role, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var role Role
	if err := db.First(&role, "name = ?", name).Error; err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}
	return &role, nil
}

func (r *Repository) ListRoles(ctx context.Context) ([]Role, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var roles []Role
	if err := db.Order("name ASC").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *Repository) UpdateRole(ctx context.Context, role *Role) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(role).Error
}

func (r *Repository) DeleteRole(ctx context.Context, id string) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	// Soft delete — roles punya kolom deleted_at
	return db.Delete(&Role{}, "id = ?", id).Error
}

// ── Permissions ──────────────────────────────────────────────────────────────

func (r *Repository) ListPermissions(ctx context.Context) ([]Permission, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var perms []Permission
	if err := db.Order("name ASC").Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

// ── Role ↔ Permission (role_has_permissions) ─────────────────────────────────

func (r *Repository) ListRolePermissionIDs(ctx context.Context, roleID string) ([]string, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var rows []RoleHasPermission
	if err := db.Where("role_id = ?", roleID).Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.PermissionID)
	}
	return ids, nil
}

// ReplaceRolePermissions menghapus semua permission role lalu assign ulang.
// Dijalankan dalam transaksi agar konsisten.
func (r *Repository) ReplaceRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&RoleHasPermission{}).Error; err != nil {
			return err
		}
		for _, pid := range permissionIDs {
			row := RoleHasPermission{PermissionID: pid, RoleID: roleID}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ── Users & User ↔ Role (model_has_roles) ────────────────────────────────────

func (r *Repository) ListUsers(ctx context.Context) ([]TenantUser, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var users []TenantUser
	if err := db.Order("name ASC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) FindUserByID(ctx context.Context, id string) (*TenantUser, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var user TenantUser
	if err := db.First(&user, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (r *Repository) ListUserRoleIDs(ctx context.Context, userID string) ([]string, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var rows []ModelHasRole
	if err := db.Where("model_id = ? AND model_type = ?", userID, "user").Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.RoleID)
	}
	return ids, nil
}

// ReplaceUserRoles menghapus semua role user lalu assign ulang.
// Dijalankan dalam transaksi agar konsisten.
func (r *Repository) ReplaceUserRoles(ctx context.Context, userID string, roleIDs []string) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("model_id = ? AND model_type = ?", userID, "user").Delete(&ModelHasRole{}).Error; err != nil {
			return err
		}
		for _, roleID := range roleIDs {
			row := ModelHasRole{RoleID: roleID, ModelType: "user", ModelID: userID}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// CountRoleUsers menghitung jumlah user yang punya role tertentu.
func (r *Repository) CountRoleUsers(ctx context.Context, roleID string) (int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return 0, err
	}
	var count int64
	if err := db.Model(&ModelHasRole{}).Where("role_id = ? AND model_type = ?", roleID, "user").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
