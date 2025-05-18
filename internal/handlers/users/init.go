package users

import (
	"context"
	"data-api/internal/handlers"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func init() {
	// Register the UserHandler with the factory.
	handlers.RegisterHandler("users", func(baseHandler handlers.BaseHandler) handlers.HandlerInterface {
		return NewHandler(baseHandler)
	})
}

func NewHandler(baseHandler handlers.BaseHandler) *UserHandler {
	return &UserHandler{
		BaseHandler:   baseHandler,
		SubjectPrefix: SubjectPrefix,
	}
}

func (h UserHandler) GetSubject() string {
	return h.SubjectPrefix + "*" // Return the subject for the user handler.
}

func (h UserHandler) Subscribe(ctx context.Context, js nats.JetStreamContext, logger *zap.SugaredLogger) bool {
	return false
}
