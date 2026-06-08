package config

import (
	"github.com/spf13/viper"
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

// Load loads the configuration.
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../../")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// Config file not found, use defaults.
		return &Config{}, nil
	}

	var cfg *Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
