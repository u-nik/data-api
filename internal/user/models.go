package user

import (
	"time"

	"github.com/uptrace/bun"
)

// User represents a user in the database.
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            string    `bun:",pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Name          string    `bun:",notnull" json:"name"`
	Email         string    `bun:",unique,notnull" json:"email"`
	PasswordHash  string    `bun:",notnull" json:"password_hash"`                        // Store hashed password
	PasskeyID     string    `bun:"passkey_id,nullzero" json:"passkey_id,omitempty"`      // For WebAuthn credential ID
	PasskeyPub    string    `bun:"passkey_pub,nullzero" json:"passkey_pubkey,omitempty"` // For WebAuthn public key
	CreatedAt     time.Time `bun:",notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time `bun:",notnull,default:current_timestamp" json:"updated_at"`
}
