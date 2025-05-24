package events

// UserCreated represents the event structure for a registered user.
type InvitationCreateEvent struct {
	BaseEvent
	Data InvitationCreateData `json:"data"`
}

type InvitationAcceptEvent struct {
	BaseEvent
	Data InvitationAcceptData `json:"data"`
}

type InvitationCreateData struct {
	Email string `json:"email" binding:"required,email"`
}

type InvitationAcceptData struct {
	Token    string `json:"token" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}
