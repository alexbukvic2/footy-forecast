// Package config loads runtime configuration from environment variables.
package config

import (
	"errors"
	"fmt"
	"os"
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
	Env  Env
	Port string
}

// Load reads configuration from the environment and validates it.
// Returns an error if any value is invalid.
func Load() (*Config, error) {
	env := Env(getenv("APP_ENV", string(EnvDev)))
	if env != EnvDev && env != EnvProd {
		return nil, fmt.Errorf("invalid APP_ENV: %q", env)
	}

	port := getenv("PORT", "8080")
	if port == "" {
		return nil, errors.New("PORT must not be empty")
	}

	return &Config{
		Env:  env,
		Port: port,
	}, nil
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
