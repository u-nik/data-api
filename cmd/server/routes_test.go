package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"data-api/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockHandler struct {
	SetupRoutesFunc func(rg *gin.RouterGroup)
}

func (m *mockHandler) SetupRoutes(rg *gin.RouterGroup) {
	if m.SetupRoutesFunc != nil {
		m.SetupRoutesFunc(rg)
	}
}

func (m *mockHandler) GetSubject() string {
	return "mock-subject"
}

func (m *mockHandler) Subscribe(ctx context.Context, js nats.JetStreamContext, logger *zap.SugaredLogger) bool {
	// Mock subscription logic
	return false
}

var middlewares = []func() gin.HandlerFunc{}

func TestSetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should register swagger route", func(t *testing.T) {
		r := gin.Default()
		handlers := make(map[string]handlers.HandlerInterface)
		logger := zap.NewNop()

		SetupRoutes(r, handlers, middlewares, logger)

		// Test if the swagger route is registered.
		// The swagger route is usually registered at "/swagger/*any" in the SetupRoutes function.
		for _, route := range r.Routes() {
			if route.Path == "/swagger/*any" && route.Method == http.MethodGet {
				return
			}
		}
		t.Error("Swagger route not registered")
	})

	t.Run("should register API routes with middleware", func(t *testing.T) {
		r := gin.Default()
		mockHandler := &mockHandler{
			SetupRoutesFunc: func(rg *gin.RouterGroup) {
				rg.GET("/test", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "test"})
				})
			},
		}
		handlers := map[string]handlers.HandlerInterface{
			"mock": mockHandler,
		}
		logger := zap.NewNop()

		SetupRoutes(r, handlers, middlewares, logger)

		req, _ := http.NewRequest(http.MethodGet, "/api/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "test")
	})

	t.Run("should apply middleware to API routes", func(t *testing.T) {
		r := gin.Default()
		middlewareCalled := false
		middlewares = append(middlewares, func() gin.HandlerFunc {
			return func(c *gin.Context) {
				middlewareCalled = true
				c.Next()
			}
		})

		mockHandler := &mockHandler{
			SetupRoutesFunc: func(rg *gin.RouterGroup) {
				rg.GET("/test", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "test"})
				})
			},
		}
		handlers := map[string]handlers.HandlerInterface{
			"mock": mockHandler,
		}
		logger := zap.NewNop()

		SetupRoutes(r, handlers, middlewares, logger)

		req, _ := http.NewRequest(http.MethodGet, "/api/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.True(t, middlewareCalled)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
