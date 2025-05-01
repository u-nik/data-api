package schema

import "log"

var (
	manager *Manager // Global schema manager instance.
)

func init() {
	log.Println("Initializing schema manager...")
	// Initialize the schema manager.
	manager = NewManager() // Initialize the schema manager.
}

func GetManager() *Manager {
	return manager // Return the global schema manager instance.
}
