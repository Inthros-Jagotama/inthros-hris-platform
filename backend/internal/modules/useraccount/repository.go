package useraccount

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	dbResolver func(ctx context.Context) (*gorm.DB, error)
	platformDB *gorm.DB
}

func NewRepository(dbResolver func(ctx context.Context) (*gorm.DB, error), platformDB *gorm.DB) *Repository {
	return &Repository{dbResolver: dbResolver, platformDB: platformDB}
}

func (r *Repository) getDB(ctx context.Context) (*gorm.DB, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required for tenant database resolution")
	}
	return r.dbResolver(ctx)
}

// ── EmployeeAccount ──────────────────────────────────────────────────────────

func (r *Repository) CreateAccount(ctx context.Context, acc *EmployeeAccount) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(acc).Error
}

func (r *Repository) FindAccountByEmployeeID(ctx context.Context, employeeID uuid.UUID) (*EmployeeAccount, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var acc EmployeeAccount
	if err := db.First(&acc, "employee_id = ?", employeeID).Error; err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}
	return &acc, nil
}

func (r *Repository) FindAccountBySetupToken(ctx context.Context, token string) (*EmployeeAccount, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var acc EmployeeAccount
	if err := db.First(&acc, "setup_token = ?", token).Error; err != nil {
		return nil, fmt.Errorf("invalid or expired setup link: %w", err)
	}
	return &acc, nil
}

func (r *Repository) UpdateAccount(ctx context.Context, acc *EmployeeAccount) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(acc).Error
}

// ── TenantUser (tabel users) ─────────────────────────────────────────────────

func (r *Repository) CreateUser(ctx context.Context, u *TenantUser) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(u).Error
}

func (r *Repository) FindUserByID(ctx context.Context, id string) (*TenantUser, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var u TenantUser
	if err := db.First(&u, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &u, nil
}

func (r *Repository) UpdateUser(ctx context.Context, u *TenantUser) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(u).Error
}

// FindUserByEmail — dipakai untuk deteksi duplikat email saat buat akun.
func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*TenantUser, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var u TenantUser
	if err := db.First(&u, "email = ?", email).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &u, nil
}

// ── Tenant Auth (login) ──────────────────────────────────────────────────────

// FindCompanyBySlug — cari company di platform DB berdasarkan slug.
// Mengembalikan companyID dan status. Dipakai tenant login untuk resolve
// tenant database (karena tidak ada subdomain).
func (r *Repository) FindCompanyBySlug(ctx context.Context, slug string) (companyID, status string, err error) {
	if r.platformDB == nil {
		return "", "", fmt.Errorf("platform database not configured")
	}
	var c struct {
		ID     string
		Status string
	}
	if err := r.platformDB.Table("companies").Select("id, status").Where("slug = ?", slug).First(&c).Error; err != nil {
		return "", "", fmt.Errorf("company not found: %w", err)
	}
	return c.ID, c.Status, nil
}

// FindCompanyByID — cari company di platform DB berdasarkan ID.
// Dipakai tenant login saat FE mengirim company_id (mis. dari env VITE_COMPANY_ID
// di development, atau hasil resolve endpoint /public/companies/resolve).
func (r *Repository) FindCompanyByID(ctx context.Context, companyID string) (status string, err error) {
	if r.platformDB == nil {
		return "", fmt.Errorf("platform database not configured")
	}
	var c struct {
		Status string
	}
	if err := r.platformDB.Table("companies").Select("status").Where("id = ?", companyID).First(&c).Error; err != nil {
		return "", fmt.Errorf("company not found: %w", err)
	}
	return c.Status, nil
}

// FindCompanyNameByID — ambil nama company dari platform DB berdasarkan ID.
// Dipakai login tenant agar response berisi company_name (FE tampilkan di
// sidebar tanpa perlu fetch endpoint platform user).
func (r *Repository) FindCompanyNameByID(ctx context.Context, companyID string) (string, error) {
	if r.platformDB == nil {
		return "", fmt.Errorf("platform database not configured")
	}
	var c struct {
		Name string
	}
	if err := r.platformDB.Table("companies").Select("name").Where("id = ?", companyID).First(&c).Error; err != nil {
		return "", fmt.Errorf("company not found: %w", err)
	}
	return c.Name, nil
}

// FindRolesForUser — daftar nama role user tenant via pivot model_has_roles.
func (r *Repository) FindRolesForUser(ctx context.Context, userID string) ([]string, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var names []string
	if err := db.Table("roles").
		Select("roles.name").
		Joins("JOIN model_has_roles mhr ON mhr.role_id = roles.id").
		Where("mhr.model_type = 'user' AND mhr.model_id = ?", userID).
		Pluck("name", &names).Error; err != nil {
		return nil, fmt.Errorf("failed to load roles: %w", err)
	}
	return names, nil
}

// FindPermissionsForUser — daftar nama permission user tenant
// (via model_has_roles → role_has_permissions → permissions).
// Gunakan .Distinct() (bukan Select("DISTINCT ...")) agar kompatibel
// dengan Pluck() di semua versi GORM v2.
func (r *Repository) FindPermissionsForUser(ctx context.Context, userID string) ([]string, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var names []string
	if err := db.Table("permissions").
		Distinct().
		Joins("JOIN role_has_permissions rhp ON rhp.permission_id = permissions.id").
		Joins("JOIN model_has_roles mhr ON mhr.role_id = rhp.role_id").
		Where("mhr.model_type = 'user' AND mhr.model_id = ?", userID).
		Pluck("name", &names).Error; err != nil {
		return nil, fmt.Errorf("failed to load permissions: %w", err)
	}
	return names, nil
}

// ── Platform user fallback (company_admin) ───────────────────────────────────

// PlatformUserRef adalah subset tabel platform_users untuk fallback login
// company_admin (platform user yang terikat company).
type PlatformUserRef struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	Role         string
	IsActive     bool
	CompanyID    *string
}

// FindPlatformUserByEmail — cari platform user berdasarkan email.
func (r *Repository) FindPlatformUserByEmail(ctx context.Context, email string) (*PlatformUserRef, error) {
	if r.platformDB == nil {
		return nil, fmt.Errorf("platform database not configured")
	}
	var u PlatformUserRef
	if err := r.platformDB.Table("platform_users").
		Select("id, name, email, password_hash, role, is_active, company_id").
		Where("email = ?", email).
		First(&u).Error; err != nil {
		return nil, fmt.Errorf("platform user not found: %w", err)
	}
	return &u, nil
}

// FindPlatformPermissionsForRole — daftar permission platform (resource.action)
// untuk role platform user (mis. company_admin) via rbac_* tables.
func (r *Repository) FindPlatformPermissionsForRole(ctx context.Context, role string) ([]string, error) {
	if r.platformDB == nil {
		return nil, fmt.Errorf("platform database not configured")
	}
	// Super admin → wildcard.
	if role == "super_admin" {
		return []string{"*"}, nil
	}
	var rows []struct {
		Resource string
		Action   string
	}
	if err := r.platformDB.Table("rbac_permissions p").
		Select("p.resource, p.action").
		Joins("JOIN rbac_role_permissions rp ON rp.permission_id = p.id").
		Joins("JOIN rbac_roles r ON r.id = rp.role_id").
		Where("r.name = ? OR r.slug = ?", role, role).
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("failed to load platform permissions: %w", err)
	}
	perms := make([]string, 0, len(rows))
	seen := make(map[string]bool)
	for _, row := range rows {
		name := row.Resource + "." + row.Action
		if !seen[name] {
			seen[name] = true
			perms = append(perms, name)
		}
	}
	return perms, nil
}

// UpdateLastLogin — update kolom last_login_at user tenant.
func (r *Repository) UpdateLastLogin(ctx context.Context, userID string) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Model(&TenantUser{}).Where("id = ?", userID).Update("last_login_at", time.Now()).Error
}

// ── Roles & assignment ───────────────────────────────────────────────────────

// FindRoleByName — cari role default tenant (mis. "Employee").
func (r *Repository) FindRoleByName(ctx context.Context, name string) (string, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return "", err
	}
	var role struct {
		ID string
	}
	if err := db.Table("roles").Select("id").Where("name = ?", name).First(&role).Error; err != nil {
		return "", fmt.Errorf("default role %q not found: %w", name, err)
	}
	return role.ID, nil
}

// AssignRoleToUser — insert pivot model_has_roles (user ↔ role).
func (r *Repository) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	pivot := ModelHasRole{
		RoleID:    roleID,
		ModelType: "user",
		ModelID:   userID,
	}
	return db.Create(&pivot).Error
}
