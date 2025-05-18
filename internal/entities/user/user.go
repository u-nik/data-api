package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:admin_users,alias:au"`
	ID            string    `bun:",pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Email         string    `bun:",unique,notnull" json:"email"`
	PasswordHash  string    `bun:"password_hash,nullzero" json:"password_hash,omitempty"`
	WebAuthnData  string    `bun:"webauthn_data,nullzero" json:"webauthn_data,omitempty"`
	InvitedAt     time.Time `bun:",notnull" json:"invited_at"`
	CreatedAt     time.Time `bun:"created_at,nullzero" json:"created_at,omitempty"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero" json:"updated_at,omitempty"`
	LastSeenAt    time.Time `bun:"last_seen_at,nullzero" json:"last_seen_at,omitempty"`
}

// CreateInvitedUser creates a new invited user in the DB (without password/webauthn)
func CreateInvitedUser(ctx context.Context, db *bun.DB, email string, invitedAt time.Time) (*User, error) {
	user := &User{
		ID:        uuid.NewString(),
		Email:     email,
		InvitedAt: invitedAt,
	}
	_, err := db.NewInsert().Model(user).Exec(ctx)
	return user, err
}
