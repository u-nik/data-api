package schema

import (
	"encoding/json"
	"fmt"
)

func (m *Manager) ValidateJSON(dataType string, data []byte) error {
	schema, ok := m.schemas[dataType]

	if !ok {
		return fmt.Errorf("schema not found for type: %s", dataType)
	}

	var dataInterface interface{}

	if err := json.Unmarshal(data, &dataInterface); err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}

	// Validate the document against the schema
	if err := schema.Validate(dataInterface); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	return nil
}
