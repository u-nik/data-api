package invitations

import "data-api/internal/handlers"

type InvitationsHandler struct {
	handlers.BaseHandler // Embedding the Handler struct for shared functionality.
}

// UserCreated represents the event structure for a registered user.
type InvitationCreateEvent struct {
	ID        string               `json:"id"`
	Data      InvitationCreateData `json:"data"`
	CreatedAt string               `json:"created_at"`
}

type InvitationCreateData struct {
	Email string `json:"email" binding:"required,email"`
}
