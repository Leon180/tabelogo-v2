package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents application error codes
type ErrorCode string

const (
	// General errors
	ErrCodeInternal       ErrorCode = "INTERNAL_ERROR"
	ErrCodeInvalidRequest ErrorCode = "INVALID_REQUEST"
	ErrCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrCodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden      ErrorCode = "FORBIDDEN"
	ErrCodeConflict       ErrorCode = "CONFLICT"

	// Auth errors
	ErrCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrCodeTokenExpired       ErrorCode = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid       ErrorCode = "TOKEN_INVALID"

	// Business logic errors
	ErrCodeUserNotFound       ErrorCode = "USER_NOT_FOUND"
	ErrCodeRestaurantNotFound ErrorCode = "RESTAURANT_NOT_FOUND"
	ErrCodeBookingNotFound    ErrorCode = "BOOKING_NOT_FOUND"
	ErrCodeBookingConflict    ErrorCode = "BOOKING_CONFLICT"
)

// AppError represents an application error
type AppError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	HTTPStatus int                    `json:"-"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Err        error                  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap implements the errors.Unwrap interface
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: getHTTPStatus(code),
	}
}

// Wrap wraps an existing error with an AppError
func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: getHTTPStatus(code),
		Err:        err,
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// WithHTTPStatus sets a custom HTTP status code
func (e *AppError) WithHTTPStatus(status int) *AppError {
	e.HTTPStatus = status
	return e
}

// getHTTPStatus returns the HTTP status code for an error code
func getHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrCodeInvalidRequest:
		return http.StatusBadRequest
	case ErrCodeUnauthorized, ErrCodeInvalidCredentials, ErrCodeTokenExpired, ErrCodeTokenInvalid:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound, ErrCodeUserNotFound, ErrCodeRestaurantNotFound, ErrCodeBookingNotFound:
		return http.StatusNotFound
	case ErrCodeConflict, ErrCodeBookingConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// Common error constructors

// NewInternalError creates an internal server error
func NewInternalError(message string) *AppError {
	return New(ErrCodeInternal, message)
}

// NewInvalidRequestError creates an invalid request error
func NewInvalidRequestError(message string) *AppError {
	return New(ErrCodeInvalidRequest, message)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(message string) *AppError {
	return New(ErrCodeNotFound, message)
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return New(ErrCodeUnauthorized, message)
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string) *AppError {
	return New(ErrCodeForbidden, message)
}

// NewConflictError creates a conflict error
func NewConflictError(message string) *AppError {
	return New(ErrCodeConflict, message)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// AsAppError converts an error to AppError if possible
func AsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}
