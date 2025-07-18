package healthcheck

import (
	"fmt"
	"time"
	"sync"
)

// CheckURL performs a health check on a single URL
// Method on HealthChecker struct
func (hc *HealthChecker) CheckURL(url string) HealthCheckResult {
	// Start of healthcheck
	start := time.Now()
    // Perform request on URL
	resp, err := hc.client.Get(url)
	// Get RTT
	latency := time.Since(start)
    // Define result as HealthCheckResult struct
	result := HealthCheckResult{
		URL:     url,
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
		result.Error = "N/A"
	} else {
		result.Status = "DOWN"
		result.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}
	// Return HealthCheckResult struct
	return result
}

// CheckURLs performs health checks on multiple URLs concurrently
// Method on HealthChecker struct
func (hc *HealthChecker) CheckURLs(urls []string) []HealthCheckResult {
	// Create waitGroup, co-ordinating goroutines
	var wg sync.WaitGroup
	// Create slice with same length as input URL list
	results := make([]HealthCheckResult, len(urls))
    // Iterate over each url, start a goroutine for each URL
	for i, url := range urls {
		wg.Add(1)
		// Start goroutine
		go func(index int, u string) {
			defer wg.Done()
			results[index] = hc.CheckURL(u)
		}(i, url)
	}

	wg.Wait()
	return results
}

// PrintResults prints the formatted health check results
func PrintResults(results []HealthCheckResult) {
	fmt.Printf("%-40s %-8s %-12s %-8s %s\n", "URL", "Status", "Latency", "HTTP", "Error")
	fmt.Println(string(make([]rune, 80)))

	for _, result := range results {
		latencyStr := fmt.Sprintf("%.2fms", float64(result.Latency.Nanoseconds())/1000000)
		httpStr := ""
		if result.HTTPCode > 0 {
			httpStr = fmt.Sprintf("%d", result.HTTPCode)
		}

		fmt.Printf("%-40s %-8s %-12s %-8s %s\n",
			result.URL,
			result.Status,
			latencyStr,
			httpStr,
			result.Error)
	}
}

// HealthCheckConfig holds configuration for batch health checks
type HealthCheckConfig struct {
	URLs    []string // String of URLs to process
	Timeout time.Duration // Timeout (standard for each endpoint)
}

// RunHealthChecks runs health checks with the given configuration
func RunHealthChecks(config HealthCheckConfig) []HealthCheckResult {
	checker := NewHealthChecker(config.Timeout) // Define checker
	return checker.CheckURLs(config.URLs)
}
