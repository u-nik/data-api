package handlers

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func SetupHandlers(
	logger *zap.Logger,
	ctx context.Context,
	rdb *redis.Client,
) map[string]HandlerInterface {
	baseHandler := BaseHandler{
		Logger: logger.Sugar(),
		Ctx:    ctx, // Global context for Redis operations.
		Rdb:    rdb, // Redis client instance.
	}

	// Register the handlers.
	return CreateHandlers(baseHandler)
}
