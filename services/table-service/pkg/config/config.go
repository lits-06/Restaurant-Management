package config

import (
	"os"
	"strconv"
)

// Config holds the configuration for Table Service.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port        int
	Environment string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

// Load loads configuration from environment variables.
func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port:        getEnvAsInt("SERVER_PORT", 50053),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DATABASE_HOST", "localhost"),
			Port:     getEnvAsInt("DATABASE_PORT", 5432),
			User:     getEnv("DATABASE_USER", "restaurant_user"),
			Password: getEnv("DATABASE_PASSWORD", "restaurant_pass"),
			Database: getEnv("DATABASE_NAME", "restaurant_db"),
			SSLMode:  getEnv("DATABASE_SSLMODE", "disable"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return defaultValue
}
