package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	Server   ServerConfig   `validate:"required"`
	Database DatabaseConfig `validate:"required"`
	Auth     AuthConfig     `validate:"required"`
}

type ServerConfig struct {
	Port string `validate:"required,numeric"`
}

type DatabaseConfig struct {
	DatabaseURL string `validate:"required"`
}

type AuthConfig struct {
	GoogleOAuth *oauth2.Config
	JwtSecret   string `validate:"required"`
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
			GoogleOAuth: &oauth2.Config{
				ClientID:     GetENV("GOOGLE_CLIENT_ID", ""),
				ClientSecret: GetENV("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  GetENV("GOOGLE_REDIRECT_URL", ""),
				Scopes: []string{
					"https://www.googleapis.com/auth/userinfo.email",
					"https://www.googleapis.com/auth/userinfo.profile",
				},
				Endpoint: google.Endpoint,
			},
		},
	}

	if err := validateConfig(conf); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return conf, nil
}
