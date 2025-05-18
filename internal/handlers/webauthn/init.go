package webauthn

import (
	"context"
	"data-api/internal/handlers"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func init() {
	handlers.RegisterHandler("webauthn", func(baseHandler handlers.BaseHandler) handlers.HandlerInterface {
		return NewHandler(baseHandler)
	})
}

func NewHandler(baseHandler handlers.BaseHandler) *WebAuthnHandler {
	return &WebAuthnHandler{}
}

func (h WebAuthnHandler) GetSubject() string {
	return "webauthn.events"
}

func (h WebAuthnHandler) Subscribe(ctx context.Context, js nats.JetStreamContext, logger *zap.SugaredLogger) bool {
	// Custom subscriber registration logic for WebAuthn
	// This is where you would implement the subscription logic specific to WebAuthn events.
	// For example, you might want to subscribe to a specific subject and handle events accordingly.
	return false // Return true if a custom subscriber is registered, otherwise false.
}
