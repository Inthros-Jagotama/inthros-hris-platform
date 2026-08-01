// Package middleware menyediakan Gin middleware untuk
// authentication, authorization, tenant resolution, logging, CORS, dll.
//
// LicenseMiddleware (PlatformLicenseMiddleware) adalah lapisan guard kedua
// setelah UserAuthMiddleware: memvalidasi bahwa modul yang diakses oleh request
// termasuk dalam lisensi (company_modules) perusahaan tenant.
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/cache"
)

// licenseCacheTTL adalah durasi cache daftar module aktif per company.
const licenseCacheTTL = 5 * time.Minute

// licenseCachePrefix adalah prefix Redis key untuk lisensi modul tenant.
const licenseCachePrefix = cache.PrefixPlatform + "license:modules:"

// CompanyModuleLister menyediakan daftar slug modul yang aktif (licensed)
// untuk sebuah company. Diimplementasikan via adapter di main.go yang
// membungkus modulemgmt.Service (ListCompanyModules + filter Enabled).
type CompanyModuleLister interface {
	EnabledModuleSlugs(companyID string) ([]string, error)
}

// moduleByPath memetakan segmen path pertama (setelah /api/v1/tenant/)
// ke slug modul. Harus sinkron dengan prefix route di tiap modul tenant.
var moduleByPath = map[string]string{
	"organizations":          "organization",
	"organization-summaries": "organization",
	"employees":              "employee",
	"job-management":         "jobmanagement",
	"competency":             "competency",
	"employee-movements":     "employeemovement",
	"attendance":             "attendance",
	"approval":               "approval",
	"payroll":                "payroll",
	"leave":                  "leave",
	"performance":            "performance",
	"recruitment":            "recruitment",
	"reimbursements":         "reimbursement",
	"trainings":              "training",
	"workforce-intelligence": "workforce-intelligence",
	"career-intelligence":    "career-intelligence",
	"settings":               "setting",
}

// LicenseCacheKey mengembalikan Redis key untuk daftar modul aktif company.
// Diekspor agar pemanggil (mis. handler subscribe/unsubscribe di main.go)
// dapat meng-invalidate cache saat lisensi berubah.
func LicenseCacheKey(companyID string) string {
	return licenseCachePrefix + companyID
}

// LicenseMiddleware memvalidasi lisensi modul tenant untuk setiap request.
//
// Alur:
//  1. Ekstrak company_id dari gin context (di-set oleh AuthJWT + TenantRequired).
//  2. Tentukan slug modul dari path request (pemetaan prefix route → module slug).
//  3. Ambil daftar slug modul aktif company dari Redis (cache), fallback ke
//     CompanyModuleLister jika cache miss, lalu cache hasilnya.
//  4. Jika slug modul tidak ada dalam daftar → 403 MODULE_NOT_LICENSED.
//
// Path yang tidak memetakan ke modul berlisensi (mis. /company-modules,
// /packages) selalu diizinkan. Fail-open berlaku untuk:
//   - Sumber data (cache/lister) error — agar outage cache/DB tidak mengunci tenant.
//   - Company tanpa data lisensi sama sekali (daftar modul kosong) — artinya lisensi
//     belum dikonfigurasi (mis. tenant legacy / dibuat tanpa package); memblokir
//     semua route akan membrukan tenant. Setelah lisensi ada, enforcement ketat.
//
// Lapisan RBAC tetap menjadi pengaman akhir.
func LicenseMiddleware(cacheManager *cache.Cache, lister CompanyModuleLister, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		companyID := c.GetString("company_id")
		if companyID == "" {
			// Tanpa konteks tenant — biarkan middleware lain (TenantRequired) yang menangani.
			c.Next()
			return
		}

		slug := moduleSlugFromPath(c.Request.URL.Path)
		if slug == "" {
			// Path tidak terkait modul berlisensi.
			c.Next()
			return
		}

		enabled, err := enabledModuleSlugs(c.Request.Context(), cacheManager, lister, companyID, logger)
		if err != nil {
			logger.Warn("License check skipped (module list unavailable)",
				zap.String("company_id", companyID),
				zap.String("module", slug),
				zap.Error(err),
			)
			c.Next()
			return
		}

		// Fail-open saat company belum punya data lisensi sama sekali
		// (daftar modul kosong) — hindari mengunci tenant legacy/baru.
		if len(enabled) == 0 {
			logger.Warn("License check skipped (company has no module license data)",
				zap.String("company_id", companyID),
				zap.String("module", slug),
			)
			c.Next()
			return
		}

		if !enabled[slug] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "MODULE_NOT_LICENSED",
					"message": "Modul '" + slug + "' tidak termasuk dalam lisensi perusahaan ini",
				},
			})
			return
		}

		c.Next()
	}
}

// moduleSlugFromPath mengekstrak slug modul dari path request tenant.
// Mengembalikan "" jika path tidak memetakan ke modul berlisensi.
func moduleSlugFromPath(path string) string {
	const tenantPrefix = "/api/v1/tenant/"
	if !strings.HasPrefix(path, tenantPrefix) {
		return ""
	}

	rest := strings.TrimPrefix(path, tenantPrefix)
	first := rest
	if idx := strings.IndexByte(first, '/'); idx >= 0 {
		first = first[:idx]
	}
	return moduleByPath[first]
}

// enabledModuleSlugs mengambil daftar slug modul aktif company dari cache
// (Redis + local), fallback ke lister saat cache miss, lalu meng-cache hasilnya.
func enabledModuleSlugs(ctx context.Context, cm *cache.Cache, lister CompanyModuleLister, companyID string, logger *zap.Logger) (map[string]bool, error) {
	key := LicenseCacheKey(companyID)

	if cm != nil {
		if data, ok := cm.Get(ctx, key); ok {
			var slugs []string
			if err := json.Unmarshal(data, &slugs); err == nil {
				return slugSet(slugs), nil
			}
		}
	}

	slugs, err := lister.EnabledModuleSlugs(companyID)
	if err != nil {
		return nil, err
	}

	if cm != nil {
		if err := cm.SetJSON(ctx, key, slugs, licenseCacheTTL); err != nil {
			logger.Warn("Failed to cache enabled module slugs",
				zap.String("company_id", companyID),
				zap.Error(err),
			)
		}
	}

	return slugSet(slugs), nil
}

func slugSet(slugs []string) map[string]bool {
	set := make(map[string]bool, len(slugs))
	for _, s := range slugs {
		set[s] = true
	}
	return set
}
