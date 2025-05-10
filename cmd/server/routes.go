package main

import (
	"data-api/internal/handlers"
	"data-api/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func SetupRoutes(
	r *gin.Engine,
	handlers map[string]handlers.HandlerInterface,
	baseLogger *zap.Logger,
) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api") // Create a new route group for API endpoints.
	api.Use(middleware.Auth(baseLogger))
	{
		for _, handler := range handlers {
			handler.SetupRoutes(api) // Setup routes for each handler.
		}
	}
}
