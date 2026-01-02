package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port               string
	DatabaseURL        string
	MaxConcurrentJobs  int
	BrowserTimeout     time.Duration
	UserAgentRotation  bool
	ProxyEnabled       bool
	APIKey             string
	Environment        string
	LogLevel           string
}

func Load() *Config {
	return &Config{
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/scraper_db?sslmode=disable"),
		MaxConcurrentJobs:  getEnvAsInt("MAX_CONCURRENT_JOBS", 3),
		BrowserTimeout:     time.Duration(getEnvAsInt("BROWSER_TIMEOUT_SECONDS", 30)) * time.Second,
		UserAgentRotation:  getEnvAsBool("USER_AGENT_ROTATION", true),
		ProxyEnabled:       getEnvAsBool("PROXY_ENABLED", false),
		APIKey:             getEnv("API_KEY", ""),
		Environment:        getEnv("ENVIRONMENT", "development"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}