package config

import "os"

type Config struct {
	ServerPort string
	Database   DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func Load() (*Config, error) {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "50052"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "restaurant_user"),
			Password: getEnv("DB_PASSWORD", "restaurant_pass"),
			Name:     getEnv("DB_NAME", "restaurant_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}, nil
}

func (c *DatabaseConfig) DSN() string {
	return "host=" + c.Host +
		" port=" + c.Port +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.Name +
		" sslmode=" + c.SSLMode
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
