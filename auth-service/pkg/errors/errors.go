package errors

import (
	"fmt"
	"runtime"
	"strings"
)

type ErrorType string

const (
	ErrorTypeValidation    ErrorType = "VALIDATION"
	ErrorTypeAuthorization ErrorType = "AUTHORIZATION"
	ErrorTypeNotFound      ErrorType = "NOT_FOUND"
	ErrorTypeInternal      ErrorType = "INTERNAL"
	ErrorTypeConflict      ErrorType = "CONFLICT"
	ErrorTypeBadRequest    ErrorType = "BAD_REQUEST"
)

type AppError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Err     error     `json:"error,omitempty"`
	Stack   string    `json:"stack,omitempty"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AppError) WithStack() *AppError {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var builder strings.Builder
	for {
		frame, more := frames.Next()
		fmt.Fprintf(&builder, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	e.Stack = builder.String()
	return e
}

func NewValidationError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Err:     err,
	}
}

func NewAuthorizationError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeAuthorization,
		Message: message,
		Err:     err,
	}
}

func NewNotFoundError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: message,
		Err:     err,
	}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Err:     err,
	}
}

func NewConflictError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeConflict,
		Message: message,
		Err:     err,
	}
}

func NewBadRequestError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeBadRequest,
		Message: message,
		Err:     err,
	}
}
