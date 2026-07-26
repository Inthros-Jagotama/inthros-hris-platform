package httputil

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidationError sends a structured 400 response with field-level validation errors.
//
// Language is auto-detected from the Accept-Language request header.
// Default language is English.
//
// Response format (English):
//
//	{
//	  "success": false,
//	  "error": {
//	    "code": "VALIDATION_ERROR",
//	    "message": "Validation failed",
//	    "fields": {
//	      "name": ["This field is required"],
//	      "email": ["Must be a valid email address"]
//	    }
//	  }
//	}
//
// Response format (Indonesian):
//
//	{
//	  "success": false,
//	  "error": {
//	    "code": "VALIDATION_ERROR",
//	    "message": "Validasi gagal",
//	    "fields": {
//	      "name": ["Field ini wajib diisi"],
//	      "email": ["Format email tidak valid"]
//	    }
//	  }
//	}
//
// req must be the same struct pointer passed to ShouldBindJSON so that the
// function can read the json tags to map struct fields to request body keys.
func ValidationError(c *gin.Context, req interface{}, err error) {
	lang := GetLang(c)
	fields := make(map[string][]string)

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		t := reflect.TypeOf(req)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		for _, fe := range validationErrors {
			fieldName := jsonFieldName(t, fe.StructField())
			if fieldName == "" {
				fieldName = toSnakeCase(fe.StructField())
			}
			msg := validationMessage(lang, fe)
			fields[fieldName] = append(fields[fieldName], msg)
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "VALIDATION_ERROR",
			"message": t(lang, "_title"),
			"fields":  fields,
		},
	})
}

// jsonFieldName returns the JSON key name for a struct field by reading its
// `json` tag. If the struct cannot be inspected or the field is not found an
// empty string is returned.
func jsonFieldName(t reflect.Type, structField string) string {
	if t.Kind() != reflect.Struct {
		return ""
	}
	field, ok := t.FieldByName(structField)
	if !ok {
		// Search embedded / anonymous fields recursively.
		for i := range t.NumField() {
			f := t.Field(i)
			if f.Anonymous && f.Type.Kind() == reflect.Struct {
				if name := jsonFieldName(f.Type, structField); name != "" {
					return name
				}
			}
		}
		return ""
	}
	jsonTag := field.Tag.Get("json")
	if jsonTag != "" {
		if idx := strings.Index(jsonTag, ","); idx != -1 {
			jsonTag = jsonTag[:idx]
		}
		if jsonTag != "" && jsonTag != "-" {
			return jsonTag
		}
	}
	// Fallback to form tag.
	formTag := field.Tag.Get("form")
	if formTag != "" {
		if idx := strings.Index(formTag, ","); idx != -1 {
			formTag = formTag[:idx]
		}
		return formTag
	}
	return ""
}

// validationMessage returns a localised error message for the given field error.
func validationMessage(lang Lang, fe validator.FieldError) string {
	tag := fe.Tag()
	param := fe.Param()

	switch tag {
	case "oneof":
		return t(lang, tag, strings.ReplaceAll(param, " ", ", "))
	case "min", "max", "len", "min_len", "max_len", "datetime":
		return t(lang, tag, param)
	default:
		// Check if the tag is registered in the locale catalog.
		if _, ok := localeMessages[tag]; ok {
			return t(lang, tag)
		}
		// Unknown tag — use the generic fallback.
		return t(lang, "_unknown", tag)
	}
}

// toSnakeCase converts PascalCase or camelCase to snake_case.
func toSnakeCase(s string) string {
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(r + 32) // to lower
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// BindAndValidate is a convenience wrapper that calls ShouldBindJSON and sends
// a structured validation error response immediately when binding or validation
// fails. It returns true when the request was successfully bound.
//
// Usage:
//
//	var req CreateEmployeeRequest
//	if !httputil.BindAndValidate(c, &req) {
//	    return
//	}
func BindAndValidate(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		ValidationError(c, req, err)
		return false
	}
	return true
}
