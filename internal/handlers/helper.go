package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *BaseHandler) GetInputFromContext(c *gin.Context) (map[string]interface{}, error) {

	value, exists := c.Get("input")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "input not found"})
		return nil, fmt.Errorf("input not found in context")
	}
	input, ok := value.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input format"})
		return nil, fmt.Errorf("invalid input format")
	}

	// This method should return the schema name based on the handler's context or configuration.
	// For example, it could be a field in the BaseHandler struct or derived from the request.
	return input, nil
}
