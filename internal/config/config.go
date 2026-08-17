package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	Server   ServerConfig   `validate:"required"`
	Database DatabaseConfig `validate:"required"`
}

type ServerConfig struct {
	Port string `validate:"required, numeric"`
}

type DatabaseConfig struct {
	DatabaseURL string `validate:"required"`
}

func validateConfig(cnf *Config) error {
	validate := validator.New()
	err := validate.Struct(cnf)

	if err != nil {
		return err
	}
	return nil
}

func GetENV(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func LoadConfig() (*Config, error) {
	conf := &Config{
		Server: ServerConfig{
			Port: GetENV("PORT", "8000"),
		},
		Database: DatabaseConfig{
			DatabaseURL: GetENV("DATABASE_URL", ""),
		},
	}

	if err := validateConfig(conf); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return conf, nil
}
