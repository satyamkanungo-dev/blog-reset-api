package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/Error"
)

var (
	envFileNotFound     = errors.New("Critical Warning: .env file not found")
	envMissingVariables = errors.New("critical: missing variables")
)

type Config struct {
	DatabaseURL string
	Port        string
	Secret      string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("%w: %w", Error.ErrFailedConfiguration, envFileNotFound)
	}

	// env config
	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
		Secret:      os.Getenv("JWT_SECRET"),
	}, nil

}

func (c *Config) Validate() error {
	if c.Secret == "" || c.DatabaseURL == "" || c.Port == "" {
		return fmt.Errorf("%w: %w", Error.ErrFailedConfiguration, envMissingVariables)

	}
	return nil
}
