package scraper

import (
	"fmt"
	"net/http"
	"strings"
)

// ErrorType represents the classification of an error
type ErrorType int

const (
	// ErrorTypeTransient represents temporary errors that should be retried
	ErrorTypeTransient ErrorType = iota
	// ErrorTypePermanent represents permanent errors that should not be retried
	ErrorTypePermanent
	// ErrorTypeRateLimit represents rate limiting errors
	ErrorTypeRateLimit
)

// ClassifyError determines the type of error for retry logic
func ClassifyError(err error) ErrorType {
	if err == nil {
		return ErrorTypePermanent
	}

	errMsg := strings.ToLower(err.Error())

	// Rate limit errors
	if strings.Contains(errMsg, "429") ||
		strings.Contains(errMsg, "too many requests") ||
		strings.Contains(errMsg, "rate limit") {
		return ErrorTypeRateLimit
	}

	// Transient errors (should retry)
	transientPatterns := []string{
		"timeout",
		"connection refused",
		"connection reset",
		"temporary failure",
		"503",
		"502",
		"500",
		"network",
		"dns",
		"i/o timeout",
		"context deadline exceeded",
	}

	for _, pattern := range transientPatterns {
		if strings.Contains(errMsg, pattern) {
			return ErrorTypeTransient
		}
	}

	// Permanent errors (should not retry)
	permanentPatterns := []string{
		"404",
		"not found",
		"invalid",
		"parse error",
		"unmarshal",
		"400",
		"401",
		"403",
	}

	for _, pattern := range permanentPatterns {
		if strings.Contains(errMsg, pattern) {
			return ErrorTypePermanent
		}
	}

	// Default to transient for unknown errors
	return ErrorTypeTransient
}

// IsRetryableHTTPStatus checks if an HTTP status code is retryable
func IsRetryableHTTPStatus(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests, // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout:      // 504
		return true
	default:
		return false
	}
}

// ShouldRetry determines if an error should be retried
func ShouldRetry(err error) bool {
	errType := ClassifyError(err)
	return errType == ErrorTypeTransient || errType == ErrorTypeRateLimit
}

// FormatErrorWithContext adds context to an error message
func FormatErrorWithContext(err error, context map[string]string) error {
	if err == nil {
		return nil
	}

	msg := err.Error()
	for key, value := range context {
		msg = fmt.Sprintf("%s [%s=%s]", msg, key, value)
	}

	return fmt.Errorf("%s", msg)
}
