package useraccount

import (
	"time"

	"github.com/inthros/hris-platform/internal/modules/employee"
)

// CreateAccountRequest — body untuk membuat akun login employee.
type CreateAccountRequest struct {
	Email string `json:"email" binding:"required,email,max=255"`
}

// AccountResponse — status akun employee.
type AccountResponse struct {
	ID               string     `json:"id"`
	CompanyID        string     `json:"company_id,omitempty"`
	EmployeeID       string     `json:"employee_id"`
	UserID           string     `json:"user_id"`
	Email            string     `json:"email"`
	RoleName         string     `json:"role_name"`
	SetupToken       string     `json:"setup_token,omitempty"`
	SetupTokenExpiry *time.Time `json:"setup_token_expires,omitempty"`
	PasswordSet      bool       `json:"password_set"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	// Employee — data karyawan lengkap TANPA masking, hanya diisi oleh
	// GetMyAccount (endpoint self-service /user-accounts/me).
	Employee *employee.EmployeeResponse `json:"employee,omitempty"`
}

// SetPasswordRequest — body untuk set password via link email (publik).
type SetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=72"`
}

// ── Tenant Auth (login/refresh publik) ───────────────────────────────────────

// TenantLoginRequest — body login user tenant (employee) atau fallback
// platform user (company_admin) yang terikat company.
// Company bisa diidentifikasi via company_slug ATAU company_id:
//   - company_slug: dipakai saat FE hanya tahu slug (mis. dari subdomain URL)
//   - company_id:   dipakai development via env (VITE_COMPANY_ID) atau
//     setelah resolve company by hostname (endpoint /public/companies/resolve)
// Minimal salah satu wajib diisi.
type TenantLoginRequest struct {
	CompanySlug string `json:"company_slug,omitempty"`
	CompanyID   string `json:"company_id,omitempty"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8,max=72"`
}

// TenantUserResponse — data user dalam response login tenant.
type TenantUserResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	CompanyID   string   `json:"company_id"`
	// CompanyName — nama company dari platform DB. Dipakai FE (mis. sidebar)
	// agar nama perusahaan tampil tanpa harus fetch endpoint platform user.
	CompanyName string `json:"company_name,omitempty"`
}

// TenantLoginResponse — response login tenant (sama bentuk dengan platform).
type TenantLoginResponse struct {
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
	TokenType    string             `json:"token_type"`
	ExpiresIn    int                `json:"expires_in"`
	User         TenantUserResponse `json:"user"`
}

// TenantRefreshRequest — body refresh access token.
type TenantRefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// TenantRefreshResponse — response refresh access token.
type TenantRefreshResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
