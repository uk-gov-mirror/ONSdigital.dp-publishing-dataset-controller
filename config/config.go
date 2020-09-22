package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

var cfg *Config

// Config represents the configuration required for florence
type Config struct {
	BindAddr                  string        `envconfig:"BIND_ADDR"`
	APIRouterURL              string        `envconfig:"API_ROUTER_URL"`
	GracefulShutdownTimeout   time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval       time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCritialTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
}

// Get retrieves the config from the environment for florence
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                  ":24000",
		APIRouterURL:              "http://localhost:23200/v1",
		GracefulShutdownTimeout:   5 * time.Second,
		HealthCheckInterval:       30 * time.Second,
		HealthCheckCritialTimeout: 90 * time.Second,
	}

	return cfg, envconfig.Process("", cfg)
}
