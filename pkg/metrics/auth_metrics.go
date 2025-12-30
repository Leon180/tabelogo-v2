package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// NewTimer creates a new Timer for measuring duration
// Usage: timer := NewTimer(AuthLoginDuration); defer timer.ObserveDuration()
func NewTimer(h prometheus.Histogram) *prometheus.Timer {
	return prometheus.NewTimer(h)
}

var (
	// ==================== Auth - Login Metrics ====================

	// AuthLoginTotal counts total login attempts by status
	AuthLoginTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_login_total",
			Help: "Total number of login attempts",
		},
		[]string{"status"}, // success, failed, blocked
	)

	// AuthLoginDuration measures login request duration
	AuthLoginDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "auth_login_duration_seconds",
			Help:    "Login request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	// ==================== Auth - Registration Metrics ====================

	// AuthRegisterTotal counts total registration attempts by status
	AuthRegisterTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_register_total",
			Help: "Total number of registration attempts",
		},
		[]string{"status"}, // success, failed, email_exists
	)

	// ==================== Auth - Token Metrics ====================

	// AuthTokenRefreshTotal counts token refresh operations
	AuthTokenRefreshTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_token_refresh_total",
			Help: "Total number of token refresh operations",
		},
		[]string{"status"}, // success, failed
	)

	// AuthTokenValidationTotal counts token validation attempts
	AuthTokenValidationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_token_validation_total",
			Help: "Total number of token validation attempts",
		},
		[]string{"status"}, // valid, invalid, expired
	)

	// ==================== Auth - Security Metrics ====================

	// AuthFailedLoginAttempts counts failed login attempts by reason
	AuthFailedLoginAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_failed_login_attempts_total",
			Help: "Total failed login attempts (for security monitoring)",
		},
		[]string{"reason"}, // wrong_password, user_not_found, account_locked
	)

	// ==================== HTTP Request Metrics ====================

	// AuthHTTPRequestsTotal counts total HTTP requests
	AuthHTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_http_requests_total",
			Help: "Total number of HTTP requests to Auth Service",
		},
		[]string{"method", "endpoint", "status"},
	)

	// AuthHTTPRequestDuration measures HTTP request duration
	AuthHTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "auth_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)
