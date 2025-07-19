package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadYamlConfig(filename string) (YamlConfig, error) {
	f, err := os.ReadFile(filename)

	if err != nil {
		return YamlConfig{}, fmt.Errorf("error reading config file: %w", err)
	}

	var config YamlConfig

	if err := yaml.Unmarshal(f, &config); err != nil {
		return YamlConfig{}, fmt.Errorf("error unmarshalling config file: %w", err)
	}
	return config, nil
}
