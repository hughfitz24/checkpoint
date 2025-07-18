package healthcheck

import (
	"net/http"
	"time"
)

// Type definition
// HealthCheckResult represents the result of a health check
type HealthCheckResult struct {
	URL      string // The URL tested
	Status   string // UP or DOWN
	Latency  time.Duration // Latency (RTT)
	HTTPCode int // HTTP error code returned
	Error    string // Error (if error occured)
}

// HealthChecker performs health checks on URLs
type HealthChecker struct {
	client  *http.Client // HTTP client
	timeout time.Duration // Time to wait before dropping connection
}

// NewHealthChecker creates a new health checker with configurable timeout
func NewHealthChecker(timeout time.Duration) *HealthChecker {
	return &HealthChecker{ // Create a new HealthChecker struct
		client: &http.Client{ // client
			Timeout: timeout, // timeout of client is passed as input
		},
		timeout: timeout, // Make timeout value directly accessible from struct
	}
}
