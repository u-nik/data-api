package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	Repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{Repo: repo}
}

// FindOrCreateInvitedUser finds an invited user by email or creates a new one if not found.
func (s *UserService) FindOrCreateInvitedUser(ctx context.Context, email string) (*User, error) {
	user, err := s.Repo.FindByEmail(ctx, email)
	if err == nil {
		return user, nil // User already exists
	}
	if err != sql.ErrNoRows {
		// Only create if not found
		if err.Error() != "sql: no rows in result set" && err.Error() != "bun: no rows in result set" {
			return nil, err
		}
	}
	user = &User{
		ID:        uuid.NewString(),
		Email:     email,
		InvitedAt: time.Now(),
	}
	err = s.Repo.Create(ctx, user)
	return user, err
}
