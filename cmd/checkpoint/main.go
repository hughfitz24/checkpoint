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
		return
	}

	cfg, err := config.ConvertConfig(yamlConfig)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	fmt.Println("Running health checks...")

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop() // Clean up
	const maxIterations = 100
	iteration := 0
	var allResults []healthcheck.HealthCheckResult
	fmt.Printf("%-60s %-8s %-12s %-8s %s\n", "Request", "Status", "Latency", "HTTP", "Error")

	for range ticker.C {

		results := healthcheck.RunHealthChecks(cfg)

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
