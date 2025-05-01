package main

import (
	"context"
	"data-api/internal/handlers"
	_ "data-api/internal/handlers/user"
	_ "data-api/internal/schema"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	ctx    = context.Background() // Global context for Redis operations.
	rdb    *redis.Client          // Redis client instance.
	logger *zap.Logger            // Global logger instance for logging.
)

func init() {
	logger = setupLogger() // Initialize the logger.
}

func main() {
	defer logger.Sync() // Flushes buffer, if any.

	r := gin.Default() // Initialize the Gin router.

	// Redis Initialization
	// Get Redis DSN from environment variable, with fallback to default
	rdb = redis.NewClient(&redis.Options{
		Addr: getEnv("REDIS_URL", "localhost:6379"), // Address of the Redis server from environment variable
	})
	defer rdb.Close() // Close the Redis client when the function exits.

	stream := SetupStream(getEnv("NATS_URL", "localhost:4222"))    // Set up the event streams.
	handlerMap := handlers.SetupHandlers(logger, ctx, rdb, stream) // Set up the handlers for the application.
	RegisterSubscribers(stream, ctx, rdb, handlerMap)              // Set up the subscribers for the event streams.
	SetupRoutes(r, handlerMap)                                     // Set up the routes for the Gin router.

	// Start the Gin server on port 8080.
	err := r.Run(getEnv("SERVER_URL", ":8080"))
	if err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

// Helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value != "" {
		return value
	}

	return fallback
}
