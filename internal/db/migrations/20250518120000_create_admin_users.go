package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(
		func(ctx context.Context, db *bun.DB) error {
			_, err := db.Exec(`
				CREATE TABLE IF NOT EXISTS users (
					id UUID PRIMARY KEY,
					email TEXT NOT NULL UNIQUE,
					password_hash TEXT,
					webauthn_data TEXT,
					invited_at TIMESTAMPTZ NOT NULL,
					created_at TIMESTAMPTZ,
					updated_at TIMESTAMPTZ,
					last_seen_at TIMESTAMPTZ
				);
			`)
			return err
		},
		func(ctx context.Context, db *bun.DB) error {
			_, err := db.Exec(`DROP TABLE IF EXISTS users;`)
			return err
		},
	)
}
