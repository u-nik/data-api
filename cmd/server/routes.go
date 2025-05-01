package main

import (
	"data-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, handlers map[string]handlers.HandlerInterface) {
	api := r.Group("/api") // Create a new route group for API endpoints.
	{
		for _, handler := range handlers {
			handler.SetupRoutes(api) // Setup routes for each handler.
		}
	}
}
