package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	GRPCPort         int
	LogLevel         string
	Database         DatabaseConfig
	MenuServiceAddr  string
	TableServiceAddr string // empty = auto-assign disabled
}

// DatabaseConfig holds PostgreSQL connection settings.
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	grpcPort := getEnvAsInt("GRPC_PORT", 50055)
	logLevel := getEnv("LOG_LEVEL", "info")
	dbHost := getEnv("DATABASE_HOST", "localhost")
	dbPort := getEnvAsInt("DATABASE_PORT", 5432)
	dbUser := getEnv("DATABASE_USER", "restaurant_user")
	dbPassword := getEnv("DATABASE_PASSWORD", "restaurant_pass")
	dbName := getEnv("DATABASE_NAME", "restaurant_db")
	sslMode := getEnv("DATABASE_SSLMODE", "disable")
	menuServiceAddr  := getEnv("MENU_SERVICE_ADDR", "localhost:50054")
	tableServiceAddr := getEnv("TABLE_SERVICE_ADDR", "")

	return &Config{
		GRPCPort: grpcPort,
		LogLevel: logLevel,
		Database: DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Name:     dbName,
			SSLMode:  sslMode,
		},
		MenuServiceAddr:  menuServiceAddr,
		TableServiceAddr: tableServiceAddr,
	}, nil
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
		return defaultValue
	}
	return value
}
