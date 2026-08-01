package company

import "time"

// CreateCompanyRequest DTO untuk create company.
type CreateCompanyRequest struct {
	Name          string  `json:"name" binding:"required,min=3,max=255"`
	Subdomain     *string `json:"subdomain,omitempty" binding:"omitempty,min=2,max=100"`
	Domain        *string `json:"domain,omitempty" binding:"omitempty,max=255"`
	NPWP          *string `json:"npwp" binding:"omitempty,len=16"`
	NIB           *string `json:"nib" binding:"omitempty,max=25"`
	Address       *string `json:"address"`
	Email         *string `json:"email" binding:"omitempty,email"`
	Phone         *string `json:"phone" binding:"omitempty,max=20"`
	AdminName     string  `json:"admin_name" binding:"required,min=1"`
	AdminEmail    string  `json:"admin_email" binding:"required,email"`
	AdminPassword string  `json:"admin_password" binding:"required,min=6"`
	PackageID     string  `json:"package_id,omitempty"`
}

// RotateCredentialsRequest DTO untuk rotasi kredensial DB tenant.
// NewPassword opsional — jika kosong, backend meng-generate password acak kuat.
type RotateCredentialsRequest struct {
	NewPassword string `json:"new_password,omitempty" binding:"omitempty,min=8,max=128"`
}

// RotateCredentialsResponse DTO untuk response rotasi kredensial.
// NewPassword dikembalikan hanya jika auto-generated oleh backend.
type RotateCredentialsResponse struct {
	CompanyID   string `json:"company_id"`
	Rotated     bool   `json:"rotated"`
	NewPassword string `json:"new_password,omitempty"`
}

// UpdateCompanyRequest DTO untuk update company.
type UpdateCompanyRequest struct {
	Name      *string `json:"name" binding:"omitempty,min=3,max=255"`
	Subdomain *string `json:"subdomain,omitempty" binding:"omitempty,min=2,max=100"`
	Domain    *string `json:"domain,omitempty" binding:"omitempty,max=255"`
	NPWP      *string `json:"npwp" binding:"omitempty,len=16"`
	NIB       *string `json:"nib" binding:"omitempty,max=25"`
	Address   *string `json:"address"`
	Email     *string `json:"email" binding:"omitempty,email"`
	Phone     *string `json:"phone" binding:"omitempty,max=20"`
}

// CompanyResponse DTO untuk response company.
type CompanyResponse struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Slug             string            `json:"slug"`
	Subdomain        *string           `json:"subdomain,omitempty"`
	Domain           *string           `json:"domain,omitempty"`
	NPWP             *string           `json:"npwp,omitempty"`
	NIB              *string           `json:"nib,omitempty"`
	Address          *string           `json:"address,omitempty"`
	Email            *string           `json:"email,omitempty"`
	Phone            *string           `json:"phone,omitempty"`
	Status           string            `json:"status"`
	AdminUser        *AdminUserInfo    `json:"admin_user,omitempty"`
	LicenseInfo      *LicenseInfo      `json:"license_info,omitempty"`
	ProvisioningInfo *ProvisioningInfo `json:"provisioning_info,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

// ResolveCompanyResponse — subset info company untuk endpoint publik
// resolve by hostname/subdomain (dipakai tenant FE sebelum login).
type ResolveCompanyResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	Subdomain *string `json:"subdomain,omitempty"`
	Domain    *string `json:"domain,omitempty"`
	Status    string  `json:"status"`
}

// LicenseInfo menampilkan informasi lisensi yang dibuat saat signup.
type LicenseInfo struct {
	ID         string `json:"id"`
	LicenseKey string `json:"license_key"`
	PlanType   string `json:"plan_type"`
	PackageID  string `json:"package_id,omitempty"`
}

// ProvisioningInfo menampilkan status provisioning tenant database.
type ProvisioningInfo struct {
	Provisioned bool   `json:"provisioned"`
	IsActive    bool   `json:"is_active"`
	Driver      string `json:"driver,omitempty"`
	DBName      string `json:"db_name,omitempty"`
}

// AdminUserInfo menampilkan informasi admin yang dibuat saat create company.
type AdminUserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// PaginatedResponse DTO untuk response pagination.
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// ToResponse mengonversi Company model ke CompanyResponse.
func (c *Company) ToResponse() CompanyResponse {
	return CompanyResponse{
		ID:        c.ID.String(),
		Name:      c.Name,
		Slug:      c.Slug,
		Subdomain: c.Subdomain,
		Domain:    c.Domain,
		NPWP:      c.NPWP,
		NIB:       c.NIB,
		Address:   c.Address,
		Email:     c.Email,
		Phone:     c.Phone,
		Status:    string(c.Status),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// ToResolveResponse mengonversi Company ke ResolveCompanyResponse (publik).
func (c *Company) ToResolveResponse() ResolveCompanyResponse {
	return ResolveCompanyResponse{
		ID:        c.ID.String(),
		Name:      c.Name,
		Slug:      c.Slug,
		Subdomain: c.Subdomain,
		Domain:    c.Domain,
		Status:    string(c.Status),
	}
}
