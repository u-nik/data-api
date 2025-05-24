package schema

import (
	"github.com/kaptinlin/jsonschema"
	"go.uber.org/zap"
)

var (
	manager *Manager // Global schema manager instance.
)

func Initialize(logger *zap.Logger) {
	logger.Sugar().Infoln("Initializing schema manager...")

	// Initialize the schema manager.
	manager = NewManager(logger)
}

// Validate checks if the provided data is valid according to the schema for the given dataType
func Validate(dataType string, data any) (*jsonschema.List, error) {
	return manager.ValidateJSON(dataType, data)
}
