package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config represents the configuration required for florence
type Config struct {
	BindAddr      string `envconfig:"BIND_ADDR"`
	ZebedeeURL    string `envconfig:"ZEBEDEE_URL"`
	DatasetAPIURL string `envconfig:"DATASET_API_URL"`
}

var cfg *Config

// Get retrieves the config from the environment for florence
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:      ":24000",
		ZebedeeURL:    "http://localhost:8082",
		DatasetAPIURL: "http://localhost:22000",
	}

	return cfg, envconfig.Process("", cfg)
}
