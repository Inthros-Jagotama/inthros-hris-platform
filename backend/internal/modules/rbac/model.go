package rbac

import (
	"time"

	"gorm.io/gorm"
)

// ─────────────────────────────────────────────────────────────────────────────
// Tenant RBAC models (pola Spatie) — tabel sudah dibuat oleh migration SQL
// tenant (011_settings.sql: permissions, roles, role_has_permissions,
// model_has_roles, model_has_permissions; 022_users.sql: users).
// Module ini hanya melakukan query/manajemen — TIDAK AutoMigrate ulang.
// ─────────────────────────────────────────────────────────────────────────────

// Role merepresentasikan tabel `roles` tenant.
type Role struct {
	ID          string         `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	GuardName   string         `gorm:"type:varchar(255);not null;default:web" json:"guard_name"`
	Description *string        `gorm:"type:varchar(255)" json:"description"`
	IsDefault   *bool          `gorm:"type:tinyint(1);default:0" json:"is_default"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (Role) TableName() string { return "roles" }

// Permission merepresentasikan tabel `permissions` tenant.
type Permission struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	GuardName string    `gorm:"type:varchar(255);not null;default:web" json:"guard_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Permission) TableName() string { return "permissions" }

// RoleHasPermission merepresentasikan tabel pivot `role_has_permissions`.
// Composite PK (permission_id, role_id) — tanpa kolom id.
type RoleHasPermission struct {
	PermissionID string `gorm:"type:char(36);primaryKey" json:"permission_id"`
	RoleID       string `gorm:"type:char(36);primaryKey" json:"role_id"`
}

func (RoleHasPermission) TableName() string { return "role_has_permissions" }

// ModelHasRole merepresentasikan tabel pivot `model_has_roles` (user ↔ role).
// Composite PK (role_id, model_type, model_id). model_type diisi "user".
type ModelHasRole struct {
	RoleID    string `gorm:"type:char(36);primaryKey" json:"role_id"`
	ModelType string `gorm:"type:varchar(255);primaryKey" json:"model_type"`
	ModelID   string `gorm:"type:char(36);primaryKey" json:"model_id"`
}

func (ModelHasRole) TableName() string { return "model_has_roles" }

// TenantUser merepresentasikan tabel `users` tenant (022_users.sql).
type TenantUser struct {
	ID         string         `gorm:"type:char(36);primaryKey" json:"id"`
	Name       string         `gorm:"type:varchar(255);not null" json:"name"`
	Email      string         `gorm:"type:varchar(255);not null" json:"email"`
	Password   string         `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	IsActive   int8           `gorm:"type:smallint;default:1" json:"is_active"`
	LastLogin  *time.Time     `gorm:"column:last_login_at" json:"last_login_at,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (TenantUser) TableName() string { return "users" }
