package schema

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const ValidInputKey = "validInput"

func JSONSchemaValidator(schemaName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input = make(map[string]interface{})

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationResult, err := Validate(schemaName, input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": validationResult,
			})
			c.Abort()
			return
		}

		c.Set(ValidInputKey, input)
		c.Next()
	}
}

func ShouldBindValidInput(c *gin.Context, obj any) error {
	raw, exists := c.Get(ValidInputKey)
	if !exists {
		return fmt.Errorf("missing valid input in context")
	}

	bytes, err := json.Marshal(raw)
	if err != nil {
		return err
	}

	if obj == nil {
		return fmt.Errorf("obj is nil")
	}

	if err := json.Unmarshal(bytes, obj); err != nil {
		return err
	}

	return nil
}
