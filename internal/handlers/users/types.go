package users

import "data-api/internal/handlers"

type UserHandler struct {
	handlers.BaseHandler // Embedding the Handler struct for shared functionality.
	SubjectPrefix        string
}

// UserCreated represents the event structure for a registered user.
type UserCreateEvent struct {
	ID        string         `json:"id"`
	Data      UserCreateData `json:"data"`
	CreatedAt string         `json:"created_at"`
}

type UserCreateData struct {
	Email string `json:"email" binding:"required,email"`
}
