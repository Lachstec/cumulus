package config

import (
	"github.com/joho/godotenv"
	"os"
)

// Config is the main configuration type for the application.
// All settings or credentials the backend needs is to be supplied here.
type Config struct {
	// Db configuration to access the primary database
	Db DbConfig
}

// LoadConfig loads the application configuration.
//
// It first checks if a .env file is available and loads it if so, overriding possibly set
// environment variables. If it is not present, it simply reads the set environment
// variables, supplying a default value if one should not be present.
// Expected variables:
//
// DB_HOST: hostname for the database (default: localhost)
// DB_PORT: port for the database (default: 5432)
// DB_USER: username for the database (default: postgres)
// DB_PASS: password for the database (default: postgres)
func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Db: DbConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASS", "postgres"),
		},
	}
	return cfg
}

// getEnv looks if an environment variable with the name key exists,
// returning it if so, else, the fallback is returned.
func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
