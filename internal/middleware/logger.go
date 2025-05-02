package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const LoggerKey = "logger"

func RequestLogger(baseLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request-ID aus Header übernehmen oder neu erzeugen
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Logger mit request_id erweitern
		reqLogger := baseLogger.With(zap.String("request_id", requestID))

		// Logger im Context speichern
		c.Set(LoggerKey, reqLogger)

		// Request-ID auch im Header zurückgeben (optional)
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}
