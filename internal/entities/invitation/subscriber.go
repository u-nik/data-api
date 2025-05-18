package invitation

import (
	"context"
	"data-api/internal/entities/user"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type InvitationSubscriber struct {
	DB     *bun.DB
	Logger *zap.SugaredLogger
}

// Start subscribes to the user subject and processes user creation events.
func (s *InvitationSubscriber) Start(ctx context.Context, stream nats.JetStreamContext, subject string) error {
	_, err := stream.Subscribe(subject+"create", func(msg *nats.Msg) {
		var evt struct {
			ID    string `json:"id"`
			Email string `json:"email"`
			// ... add more fields as needed
		}
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			s.Logger.Errorw("Error unmarshaling user event", "error", err)
			return
		}
		// Create invited user in DB
		user, err := user.CreateInvitedUser(ctx, s.DB, evt.Email, time.Now())
		if err != nil {
			s.Logger.Errorw("Error creating invited user", "error", err)
			return
		}
		// Create invite link
		inviteService := InvitationService{DB: s.DB, Logger: s.Logger, BaseURL: "https://app.localhost:3000"}
		link, err := inviteService.CreateInvite(ctx, user.ID)
		if err != nil {
			s.Logger.Errorw("Error creating invite link", "error", err)
			return
		}
		// TODO: Send email with invite link (here only log)
		s.Logger.Infow("Invite link for user", "email", user.Email, "link", link)
	}, nats.Durable("invitations-subscriber"))
	return err
}
