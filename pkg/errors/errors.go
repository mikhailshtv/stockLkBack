package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN"
	ErrorTypeInternal     ErrorType = "INTERNAL_ERROR"
	ErrorTypeDatabase     ErrorType = "DATABASE_ERROR"
)

type AppError struct {
	Type       ErrorType `json:"type"`
	Message    string    `json:"message"`
	Code       int       `json:"code"`
	Internal   error     `json:"-"`
	StackTrace string    `json:"-"`
}

func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %s (internal: %v)", e.Type, e.Message, e.Internal)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Internal
}

func NewValidationError(message string, internal error) *AppError {
	return &AppError{
		Type:     ErrorTypeValidation,
		Message:  message,
		Code:     http.StatusBadRequest,
		Internal: internal,
	}
}

func NewNotFoundError(resource string, internal error) *AppError {
	return &AppError{
		Type:     ErrorTypeNotFound,
		Message:  fmt.Sprintf("%s не найден", resource),
		Code:     http.StatusNotFound,
		Internal: internal,
	}
}

func NewUnauthorizedError(message string, internal error) *AppError {
	return &AppError{
		Type:     ErrorTypeUnauthorized,
		Message:  message,
		Code:     http.StatusUnauthorized,
		Internal: internal,
	}
}

func NewForbiddenError(message string, internal error) *AppError {
	return &AppError{
		Type:     ErrorTypeForbidden,
		Message:  message,
		Code:     http.StatusForbidden,
		Internal: internal,
	}
}

func NewDatabaseError(operation string, internal error) *AppError {
	return &AppError{
		Type:     ErrorTypeDatabase,
		Message:  fmt.Sprintf("Ошибка базы данных: %s", operation),
		Code:     http.StatusInternalServerError,
		Internal: internal,
	}
}

func NewInternalError(message string, internal error) *AppError {
	if message == "" {
		message = "Внутренняя ошибка сервера"
	}
	return &AppError{
		Type:     ErrorTypeInternal,
		Message:  message,
		Code:     http.StatusInternalServerError,
		Internal: internal,
	}
}

func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if err != nil {
		var appErr *AppError
		if errors.As(err, &appErr) {
			return appErr, true
		}
	}
	return appErr, false
}
