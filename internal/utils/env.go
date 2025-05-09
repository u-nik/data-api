package utils

import "os"

// Helper function to get environment variable with fallback
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value != "" {
		return value
	}

	return defaultValue
}

func GetRequiredEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		panic("Environment variable " + key + " is required")
	}

	return value
}
