package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

// Localize is a Gin middleware that detects the preferred language from the
// Accept-Language request header and stores it in the Gin context under the
// key "lang".
//
// Subsequent response helpers (httputil.SuccessJSON, httputil.ErrorJSON, etc.)
// automatically read this value from context, avoiding repeated header parsing
// on every response call.
//
// The middleware is idempotent — if "lang" is already set in the context
// (e.g., by a test helper), it will not override it.
func Localize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip if already set (e.g., by test helper or upstream middleware).
		if c.GetString("lang") != "" {
			c.Next()
			return
		}

		// Use httputil's canonical detection logic — no duplicate parsing.
		lang := httputil.GetLang(c)
		c.Set("lang", string(lang))
		c.Next()
	}
}
