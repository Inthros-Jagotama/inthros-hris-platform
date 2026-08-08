package useraccount

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// POST /api/v1/tenant/auth/login — publik, login user tenant (employee).
// company diidentifikasi via company_slug/company_id (FE) ATAU otomatis dari
// Host header (company_id di-set oleh middleware.TenantResolver).
func (h *Handler) Login(c *gin.Context) {
	var req TenantLoginRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	// Jika FE tidak mengirim company (mode SaaS — company di-resolve dari URL),
	// gunakan company_id yang sudah di-set TenantResolver middleware (Host header).
	if req.CompanyID == "" && req.CompanySlug == "" {
		if cid := c.GetString("company_id"); cid != "" {
			req.CompanyID = cid
		}
	}

	resp, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// POST /api/v1/tenant/auth/refresh — publik, refresh access token.
func (h *Handler) Refresh(c *gin.Context) {
	var req TenantRefreshRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.Refresh(req.RefreshToken)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// POST /api/v1/tenant/user-accounts/employees/:employeeId — buat akun login employee.
func (h *Handler) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateAccount(c.Request.Context(), c.Param("employeeId"), req.Email)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

// GET /api/v1/tenant/user-accounts/employees/:employeeId — status akun.
func (h *Handler) GetAccountStatus(c *gin.Context) {
	resp, err := h.service.GetAccountStatus(c.Request.Context(), c.Param("employeeId"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GET /api/v1/tenant/user-accounts/me — resolusi employee_id milik user yang
// sedang login, dipakai FE untuk fitur self-service (mis. Attendance My
// Dashboard) yang butuh employee_id tanpa harus di-supply manual oleh client.
func (h *Handler) GetMyAccount(c *gin.Context) {
	resp, err := h.service.GetMyAccount(c.Request.Context())
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// POST /api/v1/tenant/user-accounts/employees/:employeeId/resend — kirim ulang email setup.
func (h *Handler) ResendSetupEmail(c *gin.Context) {
	resp, err := h.service.ResendSetupEmail(c.Request.Context(), c.Param("employeeId"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// POST /api/v1/public/account/setup-password — publik, dipanggil dari halaman
// SetPassword (tujuan link email). Tidak butuh auth.
//
// Route ini TIDAK punya tenant middleware, jadi company_id harus di-resolve
// dari employee_accounts (disimpan saat CreateAccount) dan di-inject ke
// context agar service bisa mengakses tenant database.
func (h *Handler) SetPassword(c *gin.Context) {
	var req SetPasswordRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	ctx := c.Request.Context()
	// company_id dari query (dipakai route publik tanpa JWT).
	if cid := c.Query("company_id"); cid != "" {
		ctx = context.WithValue(ctx, "company_id", cid)
	}
	resp, err := h.service.SetPassword(ctx, req.Token, req.NewPassword)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
