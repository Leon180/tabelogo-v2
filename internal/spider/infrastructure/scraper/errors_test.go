package scraper

import (
	"errors"
	"testing"
)

func TestClassifyError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected ErrorType
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: ErrorTypePermanent,
		},
		{
			name:     "429 rate limit",
			err:      errors.New("HTTP 429: Too Many Requests"),
			expected: ErrorTypeRateLimit,
		},
		{
			name:     "rate limit text",
			err:      errors.New("rate limit exceeded"),
			expected: ErrorTypeRateLimit,
		},
		{
			name:     "timeout error",
			err:      errors.New("connection timeout"),
			expected: ErrorTypeTransient,
		},
		{
			name:     "503 service unavailable",
			err:      errors.New("HTTP 503 Service Unavailable"),
			expected: ErrorTypeTransient,
		},
		{
			name:     "connection refused",
			err:      errors.New("connection refused"),
			expected: ErrorTypeTransient,
		},
		{
			name:     "404 not found",
			err:      errors.New("HTTP 404 Not Found"),
			expected: ErrorTypePermanent,
		},
		{
			name:     "invalid input",
			err:      errors.New("invalid request"),
			expected: ErrorTypePermanent,
		},
		{
			name:     "parse error",
			err:      errors.New("parse error: invalid JSON"),
			expected: ErrorTypePermanent,
		},
		{
			name:     "unknown error defaults to transient",
			err:      errors.New("some unknown error"),
			expected: ErrorTypeTransient,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClassifyError(tt.err)
			if result != tt.expected {
				t.Errorf("ClassifyError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShouldRetry(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "transient error should retry",
			err:      errors.New("timeout"),
			expected: true,
		},
		{
			name:     "rate limit should retry",
			err:      errors.New("429 too many requests"),
			expected: true,
		},
		{
			name:     "permanent error should not retry",
			err:      errors.New("404 not found"),
			expected: false,
		},
		{
			name:     "nil error should not retry",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShouldRetry(tt.err)
			if result != tt.expected {
				t.Errorf("ShouldRetry() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsRetryableHTTPStatus(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expected   bool
	}{
		{"429 Too Many Requests", 429, true},
		{"500 Internal Server Error", 500, true},
		{"502 Bad Gateway", 502, true},
		{"503 Service Unavailable", 503, true},
		{"504 Gateway Timeout", 504, true},
		{"200 OK", 200, false},
		{"404 Not Found", 404, false},
		{"400 Bad Request", 400, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRetryableHTTPStatus(tt.statusCode)
			if result != tt.expected {
				t.Errorf("IsRetryableHTTPStatus(%d) = %v, want %v", tt.statusCode, result, tt.expected)
			}
		})
	}
}
