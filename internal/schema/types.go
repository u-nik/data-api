package schema

import "github.com/kaptinlin/jsonschema"

type Manager struct {
	schemas map[string]*jsonschema.Schema // Map of JSON schemas for validation.
}
