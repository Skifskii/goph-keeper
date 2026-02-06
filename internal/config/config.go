package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTP            HTTPConfig
	DatabaseDSN     string        `env:"DATABASE_DSN"`
	MasterKey       string        `env:"MASTER_KEY"`
	SecretKey       string        `env:"SECRET_KEY"`
	JWTTokenTTL     time.Duration `env:"JWT_TOKEN_TTL"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT"`
}

type HTTPConfig struct {
	Address string `env:"HTTP_ADDRESS"`
}

func New() (*Config, error) {
	cfg := Config{}

	// loading env vars from an .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, proceeding without it")
	}

	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error loading from env vars: %v", err)
	}

	return &cfg, nil
}
