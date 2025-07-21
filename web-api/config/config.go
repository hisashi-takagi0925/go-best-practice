package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort        int
	JSONPlaceholderURL string
}

func Load() *Config {
	return &Config{
		ServerPort:        getEnvAsInt("SERVER_PORT", 8080),
		JSONPlaceholderURL: getEnv("JSONPLACEHOLDER_URL", "https://jsonplaceholder.typicode.com"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	strValue := getEnv(key, "")
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return defaultValue
}