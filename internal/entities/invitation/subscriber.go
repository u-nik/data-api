package invitation

import (
	"context"
	"data-api/internal/const/subjects"
	"data-api/internal/entities/user"
	"data-api/internal/utils"
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
	// Subscriber for Invitation Create Events
	_, err := stream.Subscribe(subjects.Invitations.Create, func(msg *nats.Msg) {
		var evt struct {
			ID    string `json:"id"`
			Email string `json:"email"`
			// ... add more fields as needed
		}
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			s.Logger.Errorw("Error unmarshaling user event", "error", err)
			return
		}
		userService := user.NewUserService(user.NewUserRepository(s.DB))
		userEntity, err := userService.FindOrCreateInvitedUser(ctx, evt.Email)
		if err != nil {
			s.Logger.Errorw("Error creating invited user", "error", err)
			return
		}
		inviteService := NewInvitationService(s.DB, s.Logger, "https://app.localhost:3000")
		link, err := inviteService.CreateInvite(ctx, userEntity.ID)
		if err != nil {
			s.Logger.Errorw("Error creating invite link", "error", err)
			return
		}
		s.Logger.Infow("Invite link for user", "email", userEntity.Email, "link", link)
		msg.Ack()
	}, nats.Durable("invitations_create-invited-user"), nats.ManualAck())

	// Subscriber for Invitation Accept Events
	_, err2 := stream.Subscribe(subjects.Invitations.Accept, func(msg *nats.Msg) {
		var evt struct {
			Token      string `json:"token"`
			Name       string `json:"name"`
			Password   string `json:"password"`
			AcceptedAt string `json:"accepted_at"`
		}
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			s.Logger.Errorw("Error unmarshaling invitation accept event", "error", err)
			return
		}
		invitationRepo := NewInvitationRepository(s.DB)
		inv, err := invitationRepo.FindByToken(ctx, evt.Token)
		if err != nil {
			s.Logger.Errorw("Invitation not found for token", "token", evt.Token, "error", err)
			return
		}
		userRepo := user.NewUserRepository(s.DB)
		usr, err := userRepo.FindByID(ctx, inv.UserID)
		if err != nil {
			s.Logger.Errorw("User not found for invitation", "user_id", inv.UserID, "error", err)
			return
		}
		hash, err := utils.HashPassword(evt.Password)
		if err != nil {
			s.Logger.Errorw("Failed to hash password", "error", err)
			return
		}
		usr.PasswordHash = hash
		usr.CreatedAt = time.Now()
		usr.UpdatedAt = time.Now()
		if err := userRepo.Update(ctx, usr); err != nil {
			s.Logger.Errorw("Failed to update user after invitation accept", "user_id", usr.ID, "error", err)
			return
		}
		if err := invitationRepo.DeleteAllByUserID(ctx, usr.ID); err != nil {
			s.Logger.Warnw("Failed to delete invitations after registration", "user_id", usr.ID, "error", err)
		}
		s.Logger.Infow("Invitation accepted and user registered", "user_id", usr.ID)
		msg.Ack()
	}, nats.Durable("invitations_accept-invited-user"), nats.ManualAck())

	if err != nil {
		return err
	}
	if err2 != nil {
		return err2
	}
	return nil
}
