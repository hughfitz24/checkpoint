package config

import (
	"fmt"
	"os"
	"strings"
	"time"
	"net/url"

	"gopkg.in/yaml.v3"
	"github.com/go-playground/validator/v10"
)

func ReadYamlConfig(filename string) (*YamlConfig, error) {
	f, err := os.ReadFile(filename)

	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config YamlConfig

	if err := yaml.Unmarshal(f, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config file: %w", err)
	}

	validator := validator.New()

	if err := validator.Struct(&config); err != nil {
		return nil, fmt.Errorf("validation failed. %w", err)
	}
	return &config, nil
}

// ConvertConfig converts a YamlConfig struct to a HealthCheckConfig struct.

func ConvertConfig(yamlConfig *YamlConfig) (*HealthCheckConfig, error) {
	URLs := make([]string, 0, len(yamlConfig.Endpoints))
	baseURL := strings.TrimSuffix(yamlConfig.URL, "/")
	cfg := HealthCheckConfig{}

	for _, endpoint := range yamlConfig.Endpoints {
		URLs = append(URLs, baseURL+endpoint)
	}

	cfg.URLs = URLs
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("error in URLs created from config: %w", err)
	}
	cfg.Timeout = time.Millisecond * time.Duration(yamlConfig.Timeout)
	return &cfg, nil
}

func (hc *HealthCheckConfig) Validate() error {
	for _, urlStr := range hc.URLs {
		if _, err := url.Parse(urlStr); err != nil {
			return fmt.Errorf("invalid URL generated: %w", err)
		}
	}
	return nil
}