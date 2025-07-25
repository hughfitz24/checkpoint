package config

import (
	"time"
)

type YamlConfig struct {
	URL       string   `yaml:"url" validate:"required,url"`
	Endpoints []map[string]string `yaml:"endpoints" validate:"required,min=1"`
	Timeout   int      `yaml:"timeout" validate:"required,min=100,max=30000"`
}

// HealthCheckConfig holds configuration for batch health checks
type HealthCheckConfig struct {
	URLs    []map[string]string // Map of URLs to process. {url: "full_url", method: "GET/POST"}
	Timeout time.Duration // Timeout (standard for each endpoint)
}
