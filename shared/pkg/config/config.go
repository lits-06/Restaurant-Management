package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	GRPC     GRPCConfig
}

type ServerConfig struct {
	Port        int
	Environment string
	ServiceName string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

type JWTConfig struct {
	SecretKey            string
	AccessTokenDuration  int // in minutes
	RefreshTokenDuration int // in hours
	Issuer               string
}

type GRPCConfig struct {
	Host string
	Port int
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.environment", "development")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("jwt.access_token_duration", 15)
	viper.SetDefault("jwt.refresh_token_duration", 168)

	// Enable reading from environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	config := &Config{
		Server: ServerConfig{
			Port:        viper.GetInt("server.port"),
			Environment: viper.GetString("server.environment"),
			ServiceName: viper.GetString("server.service_name"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetInt("database.port"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			Database: viper.GetString("database.database"),
			SSLMode:  viper.GetString("database.sslmode"),
		},
		JWT: JWTConfig{
			SecretKey:            viper.GetString("jwt.secret_key"),
			AccessTokenDuration:  viper.GetInt("jwt.access_token_duration"),
			RefreshTokenDuration: viper.GetInt("jwt.refresh_token_duration"),
			Issuer:               viper.GetString("jwt.issuer"),
		},
		GRPC: GRPCConfig{
			Host: viper.GetString("grpc.host"),
			Port: viper.GetInt("grpc.port"),
		},
	}

	return config, nil
}

// GetDatabaseDSN returns the database connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Database,
		c.Database.SSLMode,
	)
}
