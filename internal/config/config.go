package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig   `validate:"required"`
	Database DatabaseConfig `validate:"required"`
	Auth     AuthConfig     `validate:"required"`
	GCS      GCSConfig      `validate:"required"`
}

type ServerConfig struct {
	Port string `validate:"required,numeric"`
}

type DatabaseConfig struct {
	DatabaseURL string `validate:"required"`
}

type AuthConfig struct {
	GoogleClientID string `validate:"required"`
	JwtSecret      string `validate:"required"`
}

type GCSConfig struct {
	GoogleCredentialPath string `validate:"required"`
	GCSBucketName        string `validate:"required"`
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
	if err := godotenv.Load(); err != nil {
		slog.Error("could not load initial env variables")
		os.Exit(1)
	}
	conf := &Config{
		Server: ServerConfig{
			Port: GetENV("PORT", "8000"),
		},
		Database: DatabaseConfig{
			DatabaseURL: GetENV("DATABASE_URL", ""),
		},
		Auth: AuthConfig{
			GoogleClientID: GetENV("GOOGLE_CLIENT_ID", ""),
			JwtSecret:      GetENV("JWT_SECRET", ""),
		},
		GCS: GCSConfig{
			GoogleCredentialPath: GetENV("GOOGLE_APPLICATION_CREDENTIALS", ""),
			GCSBucketName:        GetENV("GCS_BUCKET_NAME", ""),
		},
	}

	if err := validateConfig(conf); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return conf, nil
}
