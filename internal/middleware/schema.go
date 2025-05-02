package middleware

import (
	"data-api/internal/schema"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONSchemaValidator(schemaName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input = make(map[string]interface{})

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationResult, err := schema.Validate(schemaName, input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": validationResult,
			})
			c.Abort()
			return
		}

		c.Set("input", input)
		c.Next()
	}
}
