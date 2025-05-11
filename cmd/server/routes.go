package main

import (
	"data-api/internal/handlers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func SetupRoutes(
	r *gin.Engine,
	handlers map[string]handlers.HandlerInterface,
	apiMiddlewares []func(*zap.Logger) gin.HandlerFunc,
	baseLogger *zap.Logger,
) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api") // Create a new route group for API endpoints.

	for _, middleware := range apiMiddlewares {
		api.Use(middleware(baseLogger))
	}

	for _, handler := range handlers {
		handler.SetupRoutes(api) // Setup routes for each handler.
	}
}
