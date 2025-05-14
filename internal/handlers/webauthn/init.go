package webauthn

import "data-api/internal/handlers"

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
