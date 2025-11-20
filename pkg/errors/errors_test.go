package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(ErrCodeNotFound, "resource not found")

	if err.Code != ErrCodeNotFound {
		t.Errorf("Expected code %s, got %s", ErrCodeNotFound, err.Code)
	}
	if err.Message != "resource not found" {
		t.Errorf("Expected message 'resource not found', got %s", err.Message)
	}
	if err.HTTPStatus != http.StatusNotFound {
		t.Errorf("Expected HTTP status %d, got %d", http.StatusNotFound, err.HTTPStatus)
	}
}

func TestWrap(t *testing.T) {
	originalErr := errors.New("database connection failed")
	err := Wrap(originalErr, ErrCodeInternal, "failed to connect to database")

	if err.Code != ErrCodeInternal {
		t.Errorf("Expected code %s, got %s", ErrCodeInternal, err.Code)
	}
	if err.Err != originalErr {
		t.Errorf("Expected wrapped error to be %v, got %v", originalErr, err.Err)
	}

	// Test Unwrap
	if unwrapped := err.Unwrap(); unwrapped != originalErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, originalErr)
	}
}

func TestWithDetails(t *testing.T) {
	details := map[string]interface{}{
		"field": "email",
		"value": "invalid",
	}

	err := NewInvalidRequestError("validation failed").WithDetails(details)

	if err.Details == nil {
		t.Error("Expected details to be set")
	}
	if err.Details["field"] != "email" {
		t.Errorf("Expected field 'email', got %v", err.Details["field"])
	}
}

func TestWithHTTPStatus(t *testing.T) {
	err := New(ErrCodeInternal, "internal error").WithHTTPStatus(http.StatusServiceUnavailable)

	if err.HTTPStatus != http.StatusServiceUnavailable {
		t.Errorf("Expected HTTP status %d, got %d", http.StatusServiceUnavailable, err.HTTPStatus)
	}
}

func TestErrorString(t *testing.T) {
	tests := []struct {
		name     string
		err      *AppError
		expected string
	}{
		{
			name:     "without wrapped error",
			err:      New(ErrCodeNotFound, "user not found"),
			expected: "NOT_FOUND: user not found",
		},
		{
			name:     "with wrapped error",
			err:      Wrap(errors.New("db error"), ErrCodeInternal, "database query failed"),
			expected: "INTERNAL_ERROR: database query failed (caused by: db error)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetHTTPStatus(t *testing.T) {
	tests := []struct {
		code       ErrorCode
		wantStatus int
	}{
		{ErrCodeInvalidRequest, http.StatusBadRequest},
		{ErrCodeUnauthorized, http.StatusUnauthorized},
		{ErrCodeForbidden, http.StatusForbidden},
		{ErrCodeNotFound, http.StatusNotFound},
		{ErrCodeConflict, http.StatusConflict},
		{ErrCodeInternal, http.StatusInternalServerError},
		{ErrCodeUserNotFound, http.StatusNotFound},
		{ErrCodeBookingConflict, http.StatusConflict},
	}

	for _, tt := range tests {
		t.Run(string(tt.code), func(t *testing.T) {
			got := getHTTPStatus(tt.code)
			if got != tt.wantStatus {
				t.Errorf("getHTTPStatus(%s) = %v, want %v", tt.code, got, tt.wantStatus)
			}
		})
	}
}

func TestCommonErrorConstructors(t *testing.T) {
	tests := []struct {
		name       string
		err        *AppError
		wantCode   ErrorCode
		wantStatus int
	}{
		{
			name:       "NewInternalError",
			err:        NewInternalError("internal error"),
			wantCode:   ErrCodeInternal,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "NewInvalidRequestError",
			err:        NewInvalidRequestError("invalid input"),
			wantCode:   ErrCodeInvalidRequest,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NewNotFoundError",
			err:        NewNotFoundError("not found"),
			wantCode:   ErrCodeNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "NewUnauthorizedError",
			err:        NewUnauthorizedError("unauthorized"),
			wantCode:   ErrCodeUnauthorized,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "NewForbiddenError",
			err:        NewForbiddenError("forbidden"),
			wantCode:   ErrCodeForbidden,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "NewConflictError",
			err:        NewConflictError("conflict"),
			wantCode:   ErrCodeConflict,
			wantStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.wantCode {
				t.Errorf("Code = %v, want %v", tt.err.Code, tt.wantCode)
			}
			if tt.err.HTTPStatus != tt.wantStatus {
				t.Errorf("HTTPStatus = %v, want %v", tt.err.HTTPStatus, tt.wantStatus)
			}
		})
	}
}

func TestIsAppError(t *testing.T) {
	appErr := New(ErrCodeNotFound, "not found")
	stdErr := errors.New("standard error")

	if !IsAppError(appErr) {
		t.Error("IsAppError() should return true for AppError")
	}
	if IsAppError(stdErr) {
		t.Error("IsAppError() should return false for standard error")
	}
}

func TestAsAppError(t *testing.T) {
	appErr := New(ErrCodeNotFound, "not found")
	stdErr := errors.New("standard error")

	// Test with AppError
	converted, ok := AsAppError(appErr)
	if !ok {
		t.Error("AsAppError() should succeed for AppError")
	}
	if converted != appErr {
		t.Error("AsAppError() should return the same AppError")
	}

	// Test with standard error
	_, ok = AsAppError(stdErr)
	if ok {
		t.Error("AsAppError() should fail for standard error")
	}
}
