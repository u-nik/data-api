package invitation

import (
	"time"

	"github.com/uptrace/bun"
)

// Invitation represents a user invitation link in the database.
type Invitation struct {
	bun.BaseModel `bun:"table:admin_invitations,alias:ai"`
	ID            string    `bun:",pk,type:uuid,default:gen_random_uuid()" json:"id"`
	UserID        string    `bun:",notnull" json:"user_id"`
	Token         string    `bun:",unique,notnull" json:"token"`
	CreatedAt     time.Time `bun:",notnull,default:current_timestamp" json:"created_at"`
}
