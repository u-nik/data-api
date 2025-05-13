package utils

import "os"

// Alias for GetEnvOrDefault function
func GetEnv(key string, defaultValue string) string {
	return GetEnvOrDefault(key, defaultValue)
}

func GetEnvOrDefault(key string, defaultValue string) string {
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

func Ptr[T any](v T) *T {
	return &v
}
