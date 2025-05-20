package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(
		func(ctx context.Context, db *bun.DB) error {
			_, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS invitations (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				user_id UUID NOT NULL,
				token TEXT UNIQUE NOT NULL,
				created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
			);
		`)
			return err
		},
		func(ctx context.Context, db *bun.DB) error {
			_, err := db.Exec(`DROP TABLE IF EXISTS invitations;`)
			return err
		},
	)
}
