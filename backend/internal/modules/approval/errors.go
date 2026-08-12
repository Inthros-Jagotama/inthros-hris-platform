package approval

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

// RoutingError is a user-facing approval routing / assignee-resolution
// failure. It carries a translation key (plus optional params) so consumer
// handlers (e.g. the KPI self-assessment flow) can emit a bilingual error
// message via httputil.ErrorJSON instead of leaking the raw English service
// text. The English text in Error() is built from the same bilingual catalog,
// keeping logs and any existing string-based checks consistent.
type RoutingError struct {
	Key    string
	Params []string
	msg    string
}

// Error implements the error interface. Consumers usually classify via
// errors.As instead; Error() is used for logs and raw fallbacks.
func (e *RoutingError) Error() string { return e.msg }

// newRoutingError builds a RoutingError for the given translation key, using
// the catalog's English text as the Error() representation.
//
// params must match the number of %s placeholders in the catalog entry for
// key — a mismatch would make fmt.Sprintf append %!(EXTRA ...) noise to the
// message. Constructing a RoutingError directly (e.g. in tests) leaves msg
// empty, which is fine as long as only Key/Params are consumed.
func newRoutingError(key string, params ...string) *RoutingError {
	return &RoutingError{
		Key:    key,
		Params: params,
		msg:    httputil.Translate(httputil.LangEN, key, params...),
	}
}

// EmitRoutingError writes a bilingual error response (HTTP 400, error code
// APPROVAL_ROUTING_FAILED) when err classifies as a *RoutingError, and
// reports whether it did. Consumer handlers (leave, reimbursement, employee
// movement, payroll, KPI) call this first on errors coming back from the
// approval engine: when it returns true they return immediately, otherwise
// they fall back to their own error handling — keeping non-routing errors
// (DB failures, parse errors, ...) on their existing path.
func EmitRoutingError(c *gin.Context, err error) bool {
	var re *RoutingError
	if !errors.As(err, &re) {
		return false
	}
	httputil.ErrorJSON(c, http.StatusBadRequest, "APPROVAL_ROUTING_FAILED", re.Key, re.Params...)
	return true
}
