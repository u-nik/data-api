package invitations

import "data-api/internal/handlers"

type InvitationsHandler struct {
	handlers.BaseHandler // Embedding the Handler struct for shared functionality.
	SubjectPrefix        string
}

// UserCreated represents the event structure for a registered user.
type InvitationCreateEvent struct {
	ID        string                `json:"id"`
	Input     InvitationCreateInput `json:"user_data"`
	CreatedAt string                `json:"created_at"`
}

type InvitationCreateInput struct {
	Email string `json:"email" binding:"required,email"`
}
