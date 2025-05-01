package schema

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/kaptinlin/jsonschema"
)

// NewManager initializes and returns a new instance of Manager.
// It loads the JSON schemas during initialization.
func NewManager() *Manager {
	m := &Manager{
		schemas: loadJsonSchemas(),
	}

	return m
}

// loadJsonSchemas loads JSON schema files from the "./schema/*.json" directory pattern.
// It reads each JSON file, compiles the schemas using a JSON compiler that's configured
// to use the sonic JSON library for marshaling and unmarshaling, and returns a map where
// the keys are the filenames and the values are the compiled schema objects.
//
// The function will log fatal errors and terminate the program if it encounters issues
// with finding, reading, or compiling the schema files.
//
// Returns:
//   - map[string]*jsonschema.Schema: A map of compiled JSON schemas indexed by filename.
func loadJsonSchemas() map[string]*jsonschema.Schema {
	log.Println("Loading JSON schemas...")

	// Load JSON schema from a file or define it as a string.
	return loadJsonSchemasFromGlobPattern("./schemas/*.json")
}

// loadJsonSchemasFromGlobPattern loads and compiles JSON schemas from files matching the given glob pattern.
// It uses the jsonschema compiler with Sonic for JSON marshaling and unmarshaling.
//
// Parameters:
//   - pattern: A glob pattern string used to match JSON schema files (e.g., "schemas/*.json")
//
// Returns:
//   - A map where keys are the base filenames and values are the compiled JSON schemas
//
// The function will log fatal errors if it encounters issues with finding files,
// reading schema content, or compiling the schemas.
func loadJsonSchemasFromGlobPattern(pattern string) map[string]*jsonschema.Schema {
	log.Println("Compiling JSON schemas from glob pattern: ", pattern)

	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalf("Error finding schema files: %v", err)
	}

	// JSON schema compiler initialization
	compiler := jsonschema.NewCompiler()
	compiler.WithEncoderJSON(sonic.Marshal)
	compiler.WithDecoderJSON(sonic.Unmarshal)

	schemas := make(map[string]*jsonschema.Schema)

	for _, file := range files {
		schemaFile, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer schemaFile.Close()

		schemaData, err := io.ReadAll(schemaFile)
		if err != nil {
			log.Fatalf("Error reading schema file: %v", err)
		}

		schema, err := compiler.Compile([]byte(schemaData))
		if err != nil {
			log.Fatalf("Error compiling schema: %v", err)
		}

		filename := strings.TrimSuffix(filepath.Base(file), ".json")
		schemas[filename] = schema

		GetLogger.Infof("Loaded schema: %s", filename)
	}

	return schemas
}
