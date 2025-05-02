package user

import "data-api/internal/handlers"

type UserHandler struct {
	handlers.BaseHandler // Embedding the Handler struct for shared functionality.
	Subject              string
}

// UserCreated represents the event structure for a registered user.
type UserCreated struct {
	ID        string `json:"id"`         // Unique identifier for the user.
	Name      string `json:"name"`       // Name of the user.
	Email     string `json:"email"`      // Email address of the user.
	CreatedAt string `json:"created_at"` // Timestamp of when the user was created.
}
