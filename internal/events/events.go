package events

import (
	"time"

	"github.com/google/uuid"
)

// NewBaseEvent erstellt ein neues BaseEvent mit den Standardwerten
func NewBaseEvent(eventType string) BaseEvent {
	return BaseEvent{
		ID:        eventType,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
}

// EventFactory ist eine generische Factory-Funktion, die jede Struct, die ein BaseEvent
// eingebettet hat, instantiiert und die Standardfelder mit automatisch generierten Werten befüllt
func EventFactory[T any](constructor func(BaseEvent) T) T {
	// Generate a unique ID for the user.
	uuidObj, _ := uuid.NewV7()

	// Erstelle das BaseEvent mit Standardwerten
	baseEvent := BaseEvent{
		ID:        uuidObj.String(),
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	// Rufe den Constructor mit dem vorbereiteten BaseEvent auf
	return constructor(baseEvent)
}
