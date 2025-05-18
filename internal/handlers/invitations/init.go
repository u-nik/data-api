package invitations

import (
	"context"
	"data-api/internal/db"
	"data-api/internal/entities/invitation"
	"data-api/internal/handlers"
	"data-api/internal/stream"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func init() {
	// Register the UserHandler with the factory.
	handlers.RegisterHandler("invitations", func(baseHandler handlers.BaseHandler) handlers.HandlerInterface {
		return NewHandler(baseHandler)
	})
}

func NewHandler(baseHandler handlers.BaseHandler) *InvitationsHandler {
	return &InvitationsHandler{
		BaseHandler:   baseHandler,
		SubjectPrefix: SubjectPrefix,
	}
}

func (h InvitationsHandler) GetSubject() string {
	return h.SubjectPrefix + "*"
}

func (h InvitationsHandler) Subscribe(ctx context.Context, js nats.JetStreamContext, logger *zap.SugaredLogger) bool {
	// --- User Subscriber for admin_users DB ---
	if err := invitation.RegisterSubscriber(ctx, db.DB, (*stream.Context), zap.L().Sugar()); err != nil {
		zap.L().Fatal("Failed to start user subscriber", zap.Error(err))
	}
	return true
}
