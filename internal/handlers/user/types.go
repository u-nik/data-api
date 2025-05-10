package user

import "data-api/internal/handlers"

type UserHandler struct {
	handlers.BaseHandler // Embedding the Handler struct for shared functionality.
	Subject              string
}

// UserCreated represents the event structure for a registered user.
type UserCreatedEvent struct {
	ID        string    `json:"id"`
	UserInput UserInput `json:"user_data"`
	CreatedAt string    `json:"created_at"`
}

type UserInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}
