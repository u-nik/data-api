package users

import "data-api/internal/handlers"

type UserHandler struct {
	handlers.BaseHandler // Embedding the Handler struct for shared functionality.
	SubjectPrefix        string
}

// UserCreated represents the event structure for a registered user.
type UserCreateEvent struct {
	ID        string          `json:"id"`
	UserInput UserCreateInput `json:"user_data"`
	CreatedAt string          `json:"created_at"`
}

type UserCreateInput struct {
	Email string `json:"email" binding:"required,email"`
}
