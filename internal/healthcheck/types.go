package healthcheck

import (
	"net/http"
	"time"
)

// Type definition
// HealthCheckResult represents the result of a health check
type HealthCheckResult struct {
	Request  string        // The request made. Method and URL
	Status   string        // UP or DOWN
	Latency  time.Duration // Latency (RTT)
	HTTPCode int           // HTTP error code returned
	Error    string        // Error (if error occurred)
}

// HealthChecker performs health checks on URLs
type HealthChecker struct {
	client  *http.Client  // HTTP client
	timeout time.Duration // Time to wait before dropping connection
}
