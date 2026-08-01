package useraccount

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EmployeeAccount menghubungkan employee ↔ user tenant (tabel users, 022_users.sql)
// dan menyimpan setup token untuk alur "set password via link email".
type EmployeeAccount struct {
	ID               uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	CompanyID        uuid.UUID  `gorm:"type:char(36);not null" json:"company_id"`
	EmployeeID       uuid.UUID  `gorm:"type:char(36);not null;uniqueIndex:uk_employee_accounts_employee" json:"employee_id"`
	UserID           uuid.UUID  `gorm:"type:char(36);not null" json:"user_id"`
	Email            string     `gorm:"type:varchar(255);not null;uniqueIndex:uk_employee_accounts_email" json:"email"`
	SetupToken       *string    `gorm:"type:varchar(255)" json:"setup_token,omitempty"`
	SetupTokenExpiry *time.Time `gorm:"column:setup_token_expires" json:"setup_token_expires,omitempty"`
	CreatedBy        *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy        *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func (EmployeeAccount) TableName() string { return "employee_accounts" }

func (a *EmployeeAccount) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// TenantUser merepresentasikan baris di tabel `users` tenant (022_users.sql).
// Model terpisah agar module useraccount tidak bergantung pada module rbac.
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

// ModelHasRole merepresentasikan pivot `model_has_roles` (user ↔ role).
type ModelHasRole struct {
	RoleID    string `gorm:"type:char(36);primaryKey" json:"role_id"`
	ModelType string `gorm:"type:varchar(255);primaryKey" json:"model_type"`
	ModelID   string `gorm:"type:char(36);primaryKey" json:"model_id"`
}

func (ModelHasRole) TableName() string { return "model_has_roles" }
