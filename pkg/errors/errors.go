package errors

import (
	"errors"
	"fmt"
)

type ResourceNotFoundError struct {
	Kind string
	ID   string
}

func (e *ResourceNotFoundError) Error() string {
	return fmt.Sprintf("%s %s not found", e.Kind, e.ID)
}

func NewResourceNotFoundError(kind string, id string) *ResourceNotFoundError {
	return &ResourceNotFoundError{Kind: kind, ID: id}
}

func IsResourceNotFoundError(err error) bool {
	var e *ResourceNotFoundError
	return errors.As(err, &e)
}

type DuplicateResourceError struct {
	Kind  string
	Field string
	Value string
}

func (e *DuplicateResourceError) Error() string {
	return fmt.Sprintf("%s with %s %q already exists", e.Kind, e.Field, e.Value)
}

func NewDuplicateResourceError(kind, field, value string) *DuplicateResourceError {
	return &DuplicateResourceError{Kind: kind, Field: field, Value: value}
}

func IsDuplicateResourceError(err error) bool {
	var e *DuplicateResourceError
	return errors.As(err, &e)
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(msg string) *ValidationError {
	return &ValidationError{Message: msg}
}

func IsValidationError(err error) bool {
	var e *ValidationError
	return errors.As(err, &e)
}
