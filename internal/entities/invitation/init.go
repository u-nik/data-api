package invitation

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

// RegisterSubscriber starts the user subscriber for admin_users DB writes.
func RegisterSubscriber(ctx context.Context, db *bun.DB, js nats.JetStreamContext, logger *zap.SugaredLogger) error {
	sub := InvitationSubscriber{DB: db, Logger: logger}
	return sub.Start(ctx, js, "invitations.*")
}
