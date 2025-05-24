package schema

import (
	"fmt"

	"github.com/kaptinlin/jsonschema"
)

func (m *Manager) ValidateJSON(dataType string, data any) (*jsonschema.List, error) {
	schema, ok := m.schemas[dataType]

	if !ok {
		return nil, fmt.Errorf("schema not found for type: %s", dataType)
	}

	// Validate the document against the schema
	result := schema.Validate(data)
	if !result.IsValid() {
		// Collect all validation errors
		return result.ToList(), fmt.Errorf("validation failed for type: %s", dataType)
	}

	return nil, nil
}
