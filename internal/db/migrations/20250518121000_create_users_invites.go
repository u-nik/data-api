package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(
		func(ctx context.Context, db *bun.DB) error {
			_, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS admin_invit (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				user_id UUID NOT NULL,
				token TEXT UNIQUE NOT NULL,
				created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
			);
		`)
			return err
		},
		func(ctx context.Context, db *bun.DB) error {
			_, err := db.Exec(`DROP TABLE IF EXISTS users_invites;`)
			return err
		},
	)
}
