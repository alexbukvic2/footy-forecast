// Package config loads runtime configuration from environment variables.
package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Env represents the runtime environment the app is running in.
type Env string

// Environment values understood by the app.
const (
	EnvDev  Env = "dev"
	EnvProd Env = "prod"
)

// Config holds all runtime configuration.
type Config struct {
	Env         Env    `env:"APP_ENV"       envDefault:"dev"`
	Port        string `env:"PORT"          envDefault:"8080"`
	DatabaseURL string `env:"DATABASE_URL,required"`
}

// Load reads configuration from the environment and validates it.
// Returns an error if any required value is missing or invalid.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	if cfg.Env != EnvDev && cfg.Env != EnvProd {
		return nil, fmt.Errorf("invalid APP_ENV: %q", cfg.Env)
	}

	return cfg, nil
}
