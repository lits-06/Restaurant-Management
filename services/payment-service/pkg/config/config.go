package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Environment     string
	GRPCPort        int
	LogLevel        string
	OrderServiceURL string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Environment:     getEnv("ENVIRONMENT", "development"),
		GRPCPort:        getEnvAsInt("GRPC_PORT", 50056),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		OrderServiceURL: getEnv("ORDER_SERVICE_URL", "localhost:50055"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Printf("Warning: Invalid value for %s, using default %d\n", key, defaultValue)
		return defaultValue
	}
	return value
}
