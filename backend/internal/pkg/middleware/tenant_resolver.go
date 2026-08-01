package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// HostCompanyResolver menentukan company_id dari hostname/subdomain URL.
// Diimplementasikan oleh company module (companyModule.ResolveByHost) via
// adapter di main.go — menghindari import cycle middleware ↔ company.
type HostCompanyResolver interface {
	ResolveByHost(host string) (companyID string, err error)
}

// TenantResolver mengembalikan middleware yang otomatis menentukan tenant
// (company) dari request — mode SaaS di mana tiap company punya URL/subdomain
// sendiri, sehingga frontend tidak perlu mengirim company_slug/company_id manual.
//
// Prioritas resolusi:
//  1. company_id dari JWT claims (di-set oleh AuthJWT) — IDENTITAS MENANG.
//     Host header tidak bisa memindahkan user terautentikasi ke tenant lain.
//  2. Header `X-Tenant-ID` eksplisit (dipakai ops/dev).
//  3. `X-Forwarded-Host` → `Host` → resolve via HostCompanyResolver.
//
// Setelah berhasil resolve, middleware:
//   - meng-set gin context `company_id` (dibaca oleh TenantRequired & GetCompanyID)
//   - meng-set response header `X-Tenant-ID` agar FE tahu company aktif
//
// Jika tidak ada company yang bisa di-resolve, request diteruskan — middleware
// berikutnya (TenantRequired) yang akan menolak dengan 403 bila company_id
// memang tidak tersedia.
func TenantResolver(resolver HostCompanyResolver, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. JWT company_id (dari AuthJWT) menang — jangan pernah override.
		if cid := c.GetString("company_id"); cid != "" {
			c.Header("X-Tenant-ID", cid)
			c.Next()
			return
		}

		// 2. Header X-Tenant-ID eksplisit (harus UUID valid — hardening
		//    terhadap nilai sembarangan di jalur publik/login tanpa JWT).
		if cid := c.GetHeader("X-Tenant-ID"); cid != "" {
			if parsed, err := uuid.Parse(cid); err == nil {
				// Normalisasi ke format standar (uuid.Parse juga menerima
				// bentuk {..} / urn:uuid:.. — simpan bentuk kanonik saja).
				canonical := parsed.String()
				c.Set("company_id", canonical)
				c.Header("X-Tenant-ID", canonical)
				c.Next()
				return
			}
			logger.Warn("Invalid X-Tenant-ID header ignored",
				zap.String("x-tenant-id", cid),
				zap.String("path", c.Request.URL.Path),
			)
		}

		// 3. Resolve dari hostname (X-Forwarded-Host untuk reverse proxy,
		// fallback Host). Tanpa auth/company, tenant ditentukan dari URL.
		host := c.GetHeader("X-Forwarded-Host")
		if host == "" {
			host = c.Request.Host
		}
		if host != "" && resolver != nil {
			if cid, err := resolver.ResolveByHost(host); err == nil && cid != "" {
				c.Set("company_id", cid)
				c.Header("X-Tenant-ID", cid)
			} else if err != nil {
				logger.Debug("Tenant resolve failed",
					zap.String("host", host),
					zap.Error(err),
				)
			}
		}

		c.Next()
	}
}
