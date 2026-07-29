package setting

import "fmt"

// DuplicateCodeError represents a duplicate code validation error.
// It carries the entity table name and the duplicate code value so that
// the HTTP handler can construct a bilingual error response.
type DuplicateCodeError struct {
	Table string
	Code  string
}

func (e *DuplicateCodeError) Error() string {
	return fmt.Sprintf("code '%s' already exists in %s", e.Code, e.Table)
}
