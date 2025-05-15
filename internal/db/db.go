package db

import (
	"context"
	"data-api/internal/utils"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var DB *bun.DB

// Init initializes the Bun database connection using environment variables.
func Init() error {
	dsn := utils.GetEnv("POSTGRESQL_DSN", "postgres://postgres@localhost:5432/postgres")

	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	if err := db.PingContext(context.Background()); err != nil {
		return err
	}

	DB = bun.NewDB(db, pgdialect.New())

	return nil
}
