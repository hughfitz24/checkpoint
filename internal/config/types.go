package config

import (
	"time"
)

type YamlConfig struct {
	URL       string   `yaml:"url" validate:"required,url"`
	Endpoints []string `yaml:"endpoints" validate:"required,min=1,dive,startswith=/"`
	Timeout   int      `yaml:"timeout" validate:"required,min=100,max=30000"`
}

// HealthCheckConfig holds configuration for batch health checks
type HealthCheckConfig struct {
	URLs    []string // String of URLs to process
	Timeout time.Duration // Timeout (standard for each endpoint)
}