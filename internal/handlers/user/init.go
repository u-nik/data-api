package user

import (
	"data-api/internal/handlers"
)

func init() {
	// Register the UserHandler with the factory.
	handlers.RegisterHandler("user", func(baseHandler handlers.BaseHandler) handlers.HandlerInterface {
		return NewHandler(baseHandler)
	})
}

func NewHandler(baseHandler handlers.BaseHandler) *UserHandler {
	return &UserHandler{
		BaseHandler: baseHandler,
		Subject:     "user.events",
	}
}

func (h UserHandler) GetSubject() string {
	return h.Subject // Return the subject for the user handler.
}
