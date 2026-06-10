package config

import (
	"os"
	"strconv"
)

type Config struct {
	GRPCPort  int
	RedisHost string
	RedisPort int
	RedisPass string
	RedisDB   int
}

func Load() *Config {
	return &Config{
		GRPCPort:  getEnvInt("SERVER_PORT", 50058),
		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: getEnvInt("REDIS_PORT", 6379),
		RedisPass: getEnv("REDIS_PASSWORD", ""),
		RedisDB:   getEnvInt("REDIS_DB", 0),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
