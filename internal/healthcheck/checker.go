package healthcheck

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/hughfitz24/checkpoint/internal/config"
)

// CheckURL performs a health check on a single URL. A requestConfig map is passed as input, which contains the URL, method, and body (if applicable).
// Method on HealthChecker struct
func (hc *HealthChecker) CheckURL(requestConfig config.RequestConfig) HealthCheckResult {
	url := requestConfig.URL       // Get URL from requestConfig
	method := requestConfig.Method // Get method from requestConfig

	// TODO: Create warming function to warm endpoints before iterative checks are made.
	// The below warms the HTTP client for each request, which preserves request accuracy but slows down the process.

	// // Warm http client, so connection is cached by client before timed requests are made.
	// switch strings.ToUpper(method) {
	// case "GET":
	// 	resp, err := hc.client.Get(url)
	// 	if err == nil {
	// 		io.Copy(io.Discard, resp.Body) // Discard response body to avoid memory leak
	// 		resp.Body.Close() // Close response body
	// 	}
	// case "POST":
	// 	// If method is POST, use http.Post
	// 	resp, err := hc.client.Post(url, requestConfig.ContentType, strings.NewReader(requestConfig.Body))
	// 	if err == nil {
	// 		io.Copy(io.Discard, resp.Body) // Discard response body to avoid memory leak
	// 		resp.Body.Close() // Close response body
	// 	}
	// default:
	// 	// If method is not recognized, return Error
	// 	return HealthCheckResult{
	// 		Request:  fmt.Sprintf("%s %s", method, url), // Format request strings
	// 		Status:  "DOWN",
	// 		Latency: 0,
	// 		Error:   fmt.Sprintf("Unsupported method: %s", method), // Return error for unsupported Method
	// 	}
	// }
	// Perform request on URL depending on method
	var resp *http.Response
	var err error
	// Start of healthcheck
	start := time.Time{}
	switch strings.ToUpper(method) {
	case "GET":
		// If method is GET, use http.GET
		start = time.Now() // Start time for latency calculation
		resp, err = hc.client.Get(url)
	case "POST":
		// If method is POST, use http.POST
		reader := strings.NewReader(requestConfig.Body) // Get body from requestConfig
		start = time.Now()                              // Start time for latency calculation
		resp, err = hc.client.Post(url, requestConfig.ContentType, reader)
	default:
		// If method is not recognized, return error
		return HealthCheckResult{
			Request: fmt.Sprintf("%s %s", method, url), // Format request string
			Status:  "DOWN",
			Latency: 0,
			Error:   fmt.Sprintf("Unsupported method: %s", method),
		}
	}
	// Get RTT
	latency := time.Since(start)
	// Define result as HealthCheckResult struct
	result := HealthCheckResult{
		Request: fmt.Sprintf("%s %s", method, url), // Format request string
		Latency: latency,
	}
	// If error occurs, raise DOWN alert
	if err != nil {
		result.Status = "DOWN"
		result.Error = err.Error()
		return result
	}
	// Defer closing the HTTP client
	defer resp.Body.Close()

	result.HTTPCode = resp.StatusCode

	// Consider 2xx and 3xx status codes as healthy
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		result.Status = "UP"
		result.Error = ""
	} else {
		result.Status = "DOWN"
		result.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}
	// Return HealthCheckResult struct
	return result
}

// CheckURLs performs health checks on multiple URLs concurrently
// Method on HealthChecker struct
func (hc *HealthChecker) CheckURLs(requests []config.RequestConfig) []HealthCheckResult {
	// Create waitGroup, co-ordinating goroutines
	var wg sync.WaitGroup
	// Create slice with same length as input URL list
	results := make([]HealthCheckResult, len(requests))
	// Iterate over each url, start a goroutine for each URL
	for i, request := range requests {
		wg.Add(1)
		// Start goroutine
		go func(index int, u config.RequestConfig) {
			defer wg.Done()
			results[index] = hc.CheckURL(u)
		}(i, request)
	}

	wg.Wait()
	return results
}

// PrintResults prints the formatted health check results
func PrintResults(results []HealthCheckResult) {
	fmt.Println(strings.Repeat("-", 100))

	for _, result := range results {
		latencyStr := fmt.Sprintf("%.2fms", float64(result.Latency.Nanoseconds())/1000000)
		httpStr := ""
		if result.HTTPCode > 0 {
			httpStr = fmt.Sprintf("%d", result.HTTPCode)
		}

		fmt.Printf("%-60s %-8s %-12s %-8s %s\n",
			result.Request,
			result.Status,
			latencyStr,
			httpStr,
			result.Error)
	}
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

// RunHealthChecks runs health checks with the given configuration
func RunHealthChecks(config *config.HealthCheckConfig) []HealthCheckResult {
	checker := NewHealthChecker(config.Timeout) // Define checker
	return checker.CheckURLs(config.URLs)
}
