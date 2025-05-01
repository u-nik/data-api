package logger

import (
	"log"

	"go.uber.org/zap"
)

func init() {
	// Initialize the logger.
	logger = setupLogger()
}

// setupLogger creates and configures a Zap logger with development settings.
// It returns a sugared logger which provides a more ergonomic API.
// The function will terminate the program with a fatal error if logger initialization fails.
// Note: The logger's Sync method is deferred within this function, which may flush
// any buffered log entries before the logger is returned.
func setupLogger() *zap.Logger {
	// Initialize the logger with development configuration.
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	return logger
}

func GetLogger() *zap.Logger {
	return logger // Return the global logger instance.
}
