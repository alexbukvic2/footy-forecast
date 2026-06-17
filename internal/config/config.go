// Package config loads runtime configuration from environment variables.
package config

import (
	"fmt"
	"time"

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

	CognitoRegion           string   `env:"COGNITO_REGION,required"`
	CognitoUserPoolID       string   `env:"COGNITO_USER_POOL_ID,required"`
	CognitoAllowedClientIDs []string `env:"COGNITO_ALLOWED_CLIENT_IDS,required" envSeparator:","`

	FootballAPIKey     string        `env:"API_FOOTBALL_KEY,required"`
	FootballAPIBaseURL string        `env:"FOOTBALL_API_BASE_URL"    envDefault:"https://v3.football.api-sports.io"`
	WorkerPollInterval time.Duration `env:"WORKER_POLL_INTERVAL"     envDefault:"60s"`
	// PredictionLockLeadMinutes is how many minutes before kickoff predictions are locked.
	PredictionLockLeadMinutes int `env:"PREDICTION_LOCK_LEAD_MINUTES" envDefault:"60"`

	BedrockRegion  string `env:"BEDROCK_REGION"   envDefault:"eu-central-1"`
	BedrockModelID string `env:"BEDROCK_MODEL_ID" envDefault:"eu.anthropic.claude-sonnet-4-6"`
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

	if len(cfg.CognitoAllowedClientIDs) == 0 {
		return nil, fmt.Errorf("COGNITO_ALLOWED_CLIENT_IDS must contain at least one value")
	}

	return cfg, nil
}
