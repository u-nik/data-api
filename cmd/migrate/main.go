package main

import (
	"context"
	"data-api/internal/db"
	"data-api/internal/db/migrations"
	"data-api/internal/logger"
	"os"

	"github.com/uptrace/bun/migrate"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	defer func() { _ = zap.L().Sync() }() // Ensure all buffered log entries are flushed before the program exits.

	zap.L().Sugar().Info("Starting Data API Migrations...")

	// Initialize Bun DB
	if err := db.Init(); err != nil {
		zap.L().Fatal("Failed to initialize Bun/PostgreSQL DB", zap.Error(err))
	}
	zap.L().Sugar().Info("Connected to PostgreSQL via Bun ORM")

	migrator := migrate.NewMigrator(db.DB, migrations.Migrations)
	ctx := context.Background()
	if err := migrator.Init(ctx); err != nil {
		zap.L().Fatal("Migration init failed", zap.Error(err))
	}
	if _, err := migrator.Migrate(ctx); err != nil {
		zap.L().Fatal("Migration failed", zap.Error(err))
	}
	zap.L().Sugar().Info("Migration completed successfully")
	os.Exit(0)
}
