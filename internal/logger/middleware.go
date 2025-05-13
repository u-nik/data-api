package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const ContextKey = "logger"
const RequestIDHeader = "X-Request-ID"

// GinMiddleware returns a Gin middleware that logs requests using zap.SugaredLogger.
func RequestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		L().Infow("HTTP request",
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", path,
			"ip", c.ClientIP(),
			"duration", time.Since(start),
			"request_id", c.GetHeader(RequestIDHeader),
		)
	}
}

func RequestLogAdditionsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use Request-ID from header or generate a new one
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Extend the logger with the request ID
		reqLogger := zap.L().Sugar().With(zap.String("request_id", requestID))

		// Set the extended logger in the context
		c.Set(ContextKey, reqLogger)

		// Write the request ID to the response header
		c.Writer.Header().Set(RequestIDHeader, requestID)

		c.Next()
	}
}
