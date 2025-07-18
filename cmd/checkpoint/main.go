package main

import (
	"fmt"
	"time"
	"github.com/hughfitz24/checkpoint/internal/healthcheck"
)

func main() {
	// Example usage
	urls := []string{
		"https://google.com",
		"https://github.com",
		// "https://httpbin.org/status/200",
		// "https://httpbin.org/status/500",
		// "https://httpbin.org/delay/2",
		// "https://nonexistent-url-12345.com",
	}

	config := healthcheck.HealthCheckConfig{
		URLs:    urls,
		Timeout: 5 * time.Second,
	}

	fmt.Println("Running health checks...")

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop() // Clean up
	const maxIterations = 20
	iteration := 0
	var allResults []healthcheck.HealthCheckResult

	for range ticker.C {
		fmt.Println("Tick at", time.Now())

		results := healthcheck.RunHealthChecks(config)
		fmt.Println("\nResults:")
		healthcheck.PrintResults(results)

		allResults = append(allResults, results...)

		iteration++
		if iteration >= maxIterations {
			break
		}

	}

	// Calculate summary statistics
	var totalLatency time.Duration

	for _, result := range allResults {
		totalLatency += result.Latency
		if result.Status != "UP" {
			// Placeholder for handling non-UP statuses in the future.
		}
	}

	if len(allResults) > 0 {
		avgLatency := totalLatency / time.Duration(len(allResults))
		fmt.Printf("\nSummary: Average latency: %.2fms\n",
			float64(avgLatency.Nanoseconds())/1000000)
	} else {
		fmt.Println("\nSummary: No results to calculate average latency.")
	}
}
