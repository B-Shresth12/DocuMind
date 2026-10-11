package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	AppPort     string
	DatabaseURL string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return Config{}, fmt.Errorf("Load .env: %w", err)
	}

	cfg := Config{
		AppEnv:      strings.TrimSpace(os.Getenv("APP_ENV")),
		AppPort:     strings.TrimSpace(os.Getenv("APP_PORT")),
		DatabaseURL: strings.TrimSpace(os.Getenv("DATABASE_URL")),
	}

	if cfg.AppEnv == "" {
		cfg.AppEnv = "development"
	}

	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}

	port, err := strconv.Atoi(cfg.AppPort)
	if err != nil || port < 1 || port > 64435 {
		return Config{}, fmt.Errorf("APP_PORT must be a valid port between 1 and 65535")
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}
