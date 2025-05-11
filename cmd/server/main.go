package main

//go:generate swag init --dir ./,../../internal/handlers --output ../../api

import (
	"context"
	_ "data-api/api"
	"data-api/internal/handlers"
	_ "data-api/internal/handlers/user"
	"data-api/internal/middleware"
	"data-api/internal/schema"
	"data-api/internal/stream"
	"data-api/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	ctx = context.Background() // Global context for Redis operations.
	rdb *redis.Client          // Redis client instance.
)

// @title           Data API
// @version         1.0
// @description     Dies ist eine Data-API mit Gin und Swagger
// @host            localhost:8080
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in 				header
// @name 			Authorization
// @description 	Trage deinen Bearer Token ein: "Bearer &lt;token&gt;"
func main() {
	baseLogger := setupLogger()              // Initialize the logger.
	defer func() { _ = baseLogger.Sync() }() // Ensure all buffered log entries are flushed before the program exits.

	r := gin.Default() // Initialize the Gin router.
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger(baseLogger)) // Use the request logger middleware.

	// Redis Initialization
	// Get Redis DSN from environment variable, with fallback to default
	rdb = redis.NewClient(&redis.Options{
		Addr: utils.GetEnv("REDIS_URL", "localhost:6379"), // Address of the Redis server from environment variable
	})
	defer func() { _ = rdb.Close() }() // Close the Redis client when the function exits.

	// Initialize the schema manager.
	schema.Initialize(baseLogger)

	// Initialize the event streams.
	stream.Initialize(utils.GetEnv("NATS_URL", "localhost:4222"))

	handlerMap := handlers.SetupHandlers(baseLogger, ctx, rdb, *stream.Context) // Set up the handlers for the application.
	stream.RegisterSubscribers(ctx, rdb, handlerMap)                            // Set up the subscribers for the event streams.

	apiMiddlewares := []func(*zap.Logger) gin.HandlerFunc{
		middleware.Auth, // Add request logger middleware to the API routes.
	}

	SetupRoutes(r, handlerMap, apiMiddlewares, baseLogger) // Set up the routes for the Gin router.

	// Start the Gin server on port 8080.
	err := r.Run(utils.GetEnv("SERVER_URL", ":8080"))
	if err != nil {
		baseLogger.Fatal("Failed to start server", zap.Error(err))
	}
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
