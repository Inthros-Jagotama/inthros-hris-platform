package documenttemplate

import (
	"errors"
	"fmt"
)

var (
	ErrTemplateNotFound = errors.New("document template not found")
	ErrVersionNotFound  = errors.New("document template version not found")
)

type DuplicateActiveTemplateError struct {
	DocumentType string
}

func (e *DuplicateActiveTemplateError) Error() string {
	return fmt.Sprintf("an active template already exists for document type '%s'", e.DocumentType)
}

type DuplicateCodeError struct {
	Code string
}

func (e *DuplicateCodeError) Error() string {
	return fmt.Sprintf("template code '%s' already exists", e.Code)
}

type InvalidDocumentTypeError struct {
	DocumentType string
}

func (e *InvalidDocumentTypeError) Error() string {
	return fmt.Sprintf("invalid document type '%s'", e.DocumentType)
}

type InvalidMovementTypeError struct {
	MovementType string
}

func (e *InvalidMovementTypeError) Error() string {
	return fmt.Sprintf("invalid movement type '%s'", e.MovementType)
}

// MovementTypeNotAllowedError: movement_type hanya berlaku untuk document type
// MOVEMENT_SK (atau sebaliknya).
type MovementTypeNotAllowedError struct {
	DocumentType string
	MovementType string
}

func (e *MovementTypeNotAllowedError) Error() string {
	return fmt.Sprintf("movement_type '%s' is not allowed for document type '%s'", e.MovementType, e.DocumentType)
}
