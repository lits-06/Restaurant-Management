package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort      int
	Environment     string
	UserServiceAddr string
	Redis           RedisConfig
	JWT             JWTConfig
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	SecretKey            string
	AccessTokenMinutes   int
	RefreshTokenHours    int
	Issuer               string
}

func Load() *Config {
	return &Config{
		ServerPort:      getEnvAsInt("SERVER_PORT", 50051),
		Environment:     getEnv("ENVIRONMENT", "development"),
		UserServiceAddr: getEnv("USER_SERVICE_ADDR", "localhost:50056"),
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			SecretKey:          getEnv("JWT_SECRET", "auth-service-secret-key-change-in-production"),
			AccessTokenMinutes: getEnvAsInt("JWT_ACCESS_MINUTES", 15),
			RefreshTokenHours:  getEnvAsInt("JWT_REFRESH_HOURS", 168),
			Issuer:             getEnv("JWT_ISSUER", "auth-service"),
		},
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvAsInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
