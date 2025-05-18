package invitation

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type InvitationService struct {
	DB      *bun.DB
	Logger  *zap.SugaredLogger
	BaseURL string // e.g. https://app.localhost:3000
}

// CreateInvite creates an invite for a user and returns the invite link.
func (s *InvitationService) CreateInvite(ctx context.Context, userID string) (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	invite := &Invitation{
		ID:        uuid.NewString(),
		UserID:    userID,
		Token:     token,
		CreatedAt: time.Now(),
	}
	if _, err := s.DB.NewInsert().Model(invite).Exec(ctx); err != nil {
		return "", err
	}
	link := fmt.Sprintf("%s/invite/%s", s.BaseURL, token)
	s.Logger.Infow("User invite link generated", "user_id", userID, "link", link)
	return link, nil
}

// DeleteInvite removes an invite after it was used.
func (s *InvitationService) DeleteInvite(ctx context.Context, token string) error {
	_, err := s.DB.NewDelete().Model((*Invitation)(nil)).Where("token = ?", token).Exec(ctx)
	return err
}
