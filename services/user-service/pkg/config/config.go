package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds the configuration for the user service
type Config struct {
	ServerPort string
	LogLevel   string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		ServerPort: getEnv("SERVER_PORT", "50052"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
	}

	return cfg, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt gets an environment variable as int or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsDuration gets an environment variable as duration or returns a default value
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		fmt.Printf("Warning: invalid duration for %s, using default\n", key)
		return defaultValue
	}
	return value
}
