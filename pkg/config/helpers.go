package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// buildEnvKey builds an environment variable key with optional prefix
// Example: buildEnvKey("AUTH", "DATABASE_HOST") => "AUTH_DATABASE_HOST"
func buildEnvKey(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "_" + key
}

// getEnvWithDefault gets an environment variable or returns default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as int or returns default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool gets an environment variable as bool or returns default value
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvAsDuration gets an environment variable as duration or returns default value
// Supports formats like "15m", "1h30m", "168h"
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// getEnvAsSlice gets an environment variable as string slice (comma-separated)
func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		parts := strings.Split(value, ",")
		// Trim spaces from each part
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}
	return defaultValue
}
