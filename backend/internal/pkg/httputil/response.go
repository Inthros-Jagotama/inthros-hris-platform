package httputil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- Helpers for extracting language ---

// LangFromCtx returns the detected language from the Gin context.
// It first checks the context (set by Localize middleware), then falls back
// to parsing the Accept-Language header.
func LangFromCtx(c *gin.Context) Lang {
	return GetLang(c)
}

// tCtx is a shorthand that detects the language from the context and translates
// the given key. Equivalent to t(GetLang(c), key, params...).
func tCtx(c *gin.Context, key string, params ...string) string {
	return t(GetLang(c), key, params...)
}

// --- Standard response helpers ---

// SuccessJSON sends a 200 response with success: true and the given data.
//
//	{ "success": true, "data": { ... } }
func SuccessJSON(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// CreatedJSON sends a 201 response with success: true, data, and a translated
// message identified by msgKey.
//
//	{ "success": true, "data": { ... }, "message": "Created successfully" }
func CreatedJSON(c *gin.Context, data interface{}, msgKey string) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
		"message": tCtx(c, msgKey),
	})
}

// UpdatedJSON sends a 200 response with success: true, data, and a translated
// message identified by msgKey.
//
//	{ "success": true, "data": { ... }, "message": "Updated successfully" }
func UpdatedJSON(c *gin.Context, data interface{}, msgKey string) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"message": tCtx(c, msgKey),
	})
}

// DeletedJSON sends a 200 response with success: true and a translated message.
//
//	{ "success": true, "message": "Deleted successfully" }
func DeletedJSON(c *gin.Context, msgKey string) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": tCtx(c, msgKey),
	})
}

// MessageJSON sends a 200 response with success: true and a translated message
// (data-less). Accepts optional params for the message template.
//
//	{ "success": true, "message": "..." }
func MessageJSON(c *gin.Context, msgKey string, params ...string) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": tCtx(c, msgKey, params...),
	})
}

// SuccessJSONMessage sends a 200 response with success: true, data, and a
// translated message identified by msgKey.
//
//	{ "success": true, "data": { ... }, "message": "..." }
func SuccessJSONMessage(c *gin.Context, data interface{}, msgKey string) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"message": tCtx(c, msgKey),
	})
}

// OKJSON is an alias for SuccessJSON.
func OKJSON(c *gin.Context, data interface{}) {
	SuccessJSON(c, data)
}

// --- Error response helpers ---

// ErrorJSON sends a JSON error response with the given HTTP status, error code,
// and a translated message key.
//
//	{ "success": false, "error": { "code": "NOT_FOUND", "message": "..." } }
func ErrorJSON(c *gin.Context, status int, code string, msgKey string, params ...string) {
	c.JSON(status, gin.H{
		"success": false,
		"error": gin.H{
			"code":    code,
			"message": tCtx(c, msgKey, params...),
		},
	})
}

// ErrorRaw sends a JSON error response with a raw (dynamic) message instead of
// a translated key. Use this when the message comes from the service layer.
//
//	{ "success": false, "error": { "code": "...", "message": "raw text" } }
func ErrorRaw(c *gin.Context, status int, code string, rawMsg string) {
	c.JSON(status, gin.H{
		"success": false,
		"error": gin.H{
			"code":    code,
			"message": rawMsg,
		},
	})
}

// ErrorRawExtra sends a JSON error response with a raw message plus additional
// fields merged into the error object (e.g. validation details).
//
//	{ "success": false, "error": { "code": "...", "message": "raw", "extra": ... } }
func ErrorRawExtra(c *gin.Context, status int, code string, rawMsg string, extra gin.H) {
	body := gin.H{"code": code, "message": rawMsg}
	for k, v := range extra {
		body[k] = v
	}
	c.JSON(status, gin.H{
		"success": false,
		"error":   body,
	})
}

// ErrorSimple sends a JSON error response with only a raw message (no code).
//
//	{ "success": false, "error": "raw message" }
func ErrorSimple(c *gin.Context, status int, rawMsg string) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   rawMsg,
	})
}

// --- Convenience shortcuts for common HTTP status codes ---

// NotFound sends a 404 error with a translated "not found" message.
// If rawMsg is not empty, it is used instead of the translated key.
func NotFound(c *gin.Context, rawMsg string) {
	if rawMsg != "" {
		ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", rawMsg)
		return
	}
	ErrorJSON(c, http.StatusNotFound, "NOT_FOUND", "error.not_found")
}

// InternalError sends a 500 error with the given raw message.
func InternalError(c *gin.Context, rawMsg string) {
	ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", rawMsg)
}

// Unauthorized sends a 401 error with the given raw message.
func Unauthorized(c *gin.Context, rawMsg string) {
	ErrorRaw(c, http.StatusUnauthorized, "UNAUTHORIZED", rawMsg)
}

// Conflict sends a 409 error with the given raw message.
func Conflict(c *gin.Context, code string, rawMsg string) {
	ErrorRaw(c, http.StatusConflict, code, rawMsg)
}

// BadRequest sends a 400 error with the given raw message.
func BadRequest(c *gin.Context, rawMsg string) {
	ErrorRaw(c, http.StatusBadRequest, "BAD_REQUEST", rawMsg)
}
