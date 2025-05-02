package schema

import (
	"github.com/kaptinlin/jsonschema"
	"go.uber.org/zap"
)

type Manager struct {
	logger  *zap.Logger                   // Logger for logging events and errors.
	schemas map[string]*jsonschema.Schema // Map of JSON schemas for validation.
}
