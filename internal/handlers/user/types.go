package user

import "data-api/internal/handlers"

type UserHandler struct {
	handlers.BaseHandler // Embedding the Handler struct for shared functionality.
	Subject              string
}

// UserRegistered represents the event structure for a registered user.
type UserRegistered struct {
	ID    string `json:"id" jsonschema:"uuid"`     // Unique identifier for the user.
	Email string `json:"email" jsonschema:"email"` // Email address of the user.
}
