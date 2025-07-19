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
	Error    string // Error (if error occurred)
}

// HealthChecker performs health checks on URLs
type HealthChecker struct {
		client  *http.Client // HTTP client
	timeout time.Duration // Time to wait before dropping connection
}


// HealthCheckConfig holds configuration for batch health checks
type HealthCheckConfig struct {
	URLs    []string // String of URLs to process
	Timeout time.Duration // Timeout (standard for each endpoint)
}

