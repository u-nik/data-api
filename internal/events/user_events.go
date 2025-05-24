package events

type UserCreateEvent struct {
	BaseEvent
	Data UserCreateData `json:"data"`
}

type UserCreateData struct {
	Email string `json:"email" binding:"required,email"`
}
