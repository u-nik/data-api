package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type BaseHandler struct {
	Logger *zap.SugaredLogger    // Logger for logging events and errors.
	Rdb    *redis.Client         // Redis client for caching and data storage.
	Ctx    context.Context       // Context for managing request lifecycle.
	Stream nats.JetStreamContext // NATS JetStream context for event streaming.
}

type HandlerInterface interface {
	GetSubject() string             // Method to get the subject for the handler.
	SetupRoutes(g *gin.RouterGroup) // Method to create routes for the handler.
}

type HandlerFactory func(baseHandler BaseHandler) HandlerInterface
