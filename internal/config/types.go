package config

type YamlConfig struct {
	URL       string   `yaml:"url"`
	Endpoints []string `yaml:"endpoints"`
	Timeout   int      `yaml:"timeout"`
}
