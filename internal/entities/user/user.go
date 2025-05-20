package user

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            string    `bun:",pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Email         string    `bun:",unique,notnull" json:"email"`
	PasswordHash  string    `bun:"password_hash,nullzero" json:"password_hash,omitempty"`
	WebAuthnData  string    `bun:"webauthn_data,nullzero" json:"webauthn_data,omitempty"`
	InvitedAt     time.Time `bun:",notnull" json:"invited_at"`
	CreatedAt     time.Time `bun:"created_at,nullzero" json:"created_at,omitempty"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero" json:"updated_at,omitempty"`
	LastSeenAt    time.Time `bun:"last_seen_at,nullzero" json:"last_seen_at,omitempty"`
}
