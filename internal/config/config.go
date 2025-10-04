package config

import (
	"os"
)

type Config struct {
	DatabaseURL          string
	GoogleClientID       string
	GoogleClientSecret   string
	GoogleRedirectURL    string
	SessionSecret        string
	FrontendURL          string
	Port                 string
}

func Load() *Config {
	return &Config{
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://clockworth:clockworth_dev_password@localhost:5432/clockworth?sslmode=disable"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
		SessionSecret:      getEnv("SESSION_SECRET", ""),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:5173"),
		Port:               getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
