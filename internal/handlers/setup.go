package handlers

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func SetupHandlers(
	logger *zap.Logger,
	ctx context.Context,
	rdb *redis.Client,
	stream nats.JetStreamContext,
) map[string]HandlerInterface {
	baseHandler := BaseHandler{
		Logger: logger.Sugar(),
		Ctx:    ctx,    // Global context for Redis operations.
		Rdb:    rdb,    // Redis client instance.
		Stream: stream, // NATS JetStream context for event streaming.
	}

	// Register the handlers.
	return CreateHandlers(baseHandler)
}
