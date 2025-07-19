package main

import (
	"fmt"
	"time"

	"github.com/hughfitz24/checkpoint/internal/config"
	"github.com/hughfitz24/checkpoint/internal/healthcheck"
)

func main() {

	configFile := "configs/test.yml"

	yamlConfig, err := config.ReadYamlConfig(configFile)
	if err != nil {
		fmt.Println("err: ", err)
	}

	config := healthcheck.ConvertConfig(yamlConfig)

	// 	for _, endpoint := range endpoints {
	// 		urls = append(urls, api + endpoint)
	// 	}
	// 	config := healthcheck.HealthCheckConfig{
	// 		URLs:    urls,
	// 		Timeout: 5 * time.Second,
	// 	}

	fmt.Println("Running health checks...")

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop() // Clean up
	const maxIterations = 100
	iteration := 0
	var allResults []healthcheck.HealthCheckResult

	for range ticker.C {

		results := healthcheck.RunHealthChecks(config)

		allResults = append(allResults, results...)

		if iteration%10 == 0 {
			healthcheck.PrintResults(results)
		}

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
