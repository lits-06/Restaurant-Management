// Package config contains configuration for the report service.
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds configuration for the report service.
type Config struct {
	// gRPC Server
	GrpcPort int
	GrpcHost string

	// Service Addresses
	OrderServiceAddr      string
	PaymentServiceAddr    string
	InventoryServiceAddr  string
	UserServiceAddr       string
	MenuServiceAddr       string

	// Environment
	Environment string
	LogLevel    string
}

// Load loads configuration from environment variables.
func Load() *Config {
	return &Config{
		GrpcPort:             getIntEnv("REPORT_GRPC_PORT", 50059),
		GrpcHost:             getEnv("REPORT_GRPC_HOST", "0.0.0.0"),
		OrderServiceAddr:     getEnv("ORDER_SERVICE_ADDR", "localhost:50055"),
		PaymentServiceAddr:   getEnv("PAYMENT_SERVICE_ADDR", "localhost:50056"),
		InventoryServiceAddr: getEnv("INVENTORY_SERVICE_ADDR", "localhost:50057"),
		UserServiceAddr:      getEnv("USER_SERVICE_ADDR", "localhost:50052"),
		MenuServiceAddr:      getEnv("MENU_SERVICE_ADDR", "localhost:50054"),
		Environment:          getEnv("ENVIRONMENT", "development"),
		LogLevel:             getEnv("LOG_LEVEL", "info"),
	}
}

// GetGrpcAddress returns the full gRPC address.
func (c *Config) GetGrpcAddress() string {
	return fmt.Sprintf("%s:%d", c.GrpcHost, c.GrpcPort)
}

// Helper functions
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}
