package main

//go:generate swag init --dir ./,../../internal/handlers --output ../../api

import (
	"context"
	_ "data-api/api"
	"data-api/internal/auth"
	"data-api/internal/handlers"
	_ "data-api/internal/handlers/user"
	"data-api/internal/logger"
	"data-api/internal/schema"
	"data-api/internal/stream"
	"data-api/internal/utils"

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
	// Initialize Logger
	logger.Init()
	defer func() { _ = zap.L().Sync() }() // Ensure all buffered log entries are flushed before the program exits.

	zap.L().Sugar().Info("Starting Data API...")

	auth.Initialize() // Set up OIDC authentication middleware.

	r := gin.New() // Initialize the Gin router.
	r.Use(gin.Recovery())
	r.Use(logger.RequestLoggingMiddleware())      // Use the request logger middleware.
	r.Use(logger.RequestLogAdditionsMiddleware()) // Use the request logger middleware.

	// Redis Initialization
	// Get Redis DSN from environment variable, with fallback to default
	rdb = redis.NewClient(&redis.Options{
		Addr: utils.GetEnv("REDIS_URL", "localhost:6379"), // Address of the Redis server from environment variable
	})
	defer func() { _ = rdb.Close() }() // Close the Redis client when the function exits.

	// Initialize the schema manager.
	schema.Initialize(zap.L())

	// Initialize the event streams.
	stream.Initialize(utils.GetEnv("NATS_URL", "localhost:4222"))

	handlerMap := handlers.SetupHandlers(zap.L(), ctx, rdb, *stream.Context) // Set up the handlers for the application.
	stream.RegisterSubscribers(ctx, rdb, handlerMap)                         // Set up the subscribers for the event streams.

	apiMiddlewares := []func() gin.HandlerFunc{
		auth.Auth, // Add request logger middleware to the API routes.
	}

	SetupRoutes(r, handlerMap, apiMiddlewares, zap.L()) // Set up the routes for the Gin router.

	// Start the Gin server on port 8080.
	err := r.Run(utils.GetEnv("SERVER_URL", ":8080"))
	if err != nil {
		zap.L().Fatal("Failed to start server", zap.Error(err))
	}
}
