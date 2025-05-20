package main

import (
	"bufio"
	"context"
	"data-api/internal/db"
	"data-api/internal/db/migrations"
	"data-api/internal/logger"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/uptrace/bun/migrate"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	defer func() { _ = zap.L().Sync() }() // Ensure all buffered log entries are flushed before the program exits.

	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate [migrate|generate] [name]")
		os.Exit(1)
	}
	action := os.Args[1]

	switch action {
	case "migrate":
		zlog := zap.L().Sugar()
		zlog.Info("Starting Data API Migrations...")
		// Initialize Bun DB
		if err := db.Init(); err != nil {
			zap.L().Fatal("Failed to initialize Bun/PostgreSQL DB", zap.Error(err))
		}
		zlog.Info("Connected to PostgreSQL via Bun ORM")
		migrator := migrate.NewMigrator(db.DB, migrations.Migrations)
		ctx := context.Background()
		if err := migrator.Init(ctx); err != nil {
			zap.L().Fatal("Migration init failed", zap.Error(err))
		}
		if _, err := migrator.Migrate(ctx); err != nil {
			zap.L().Fatal("Migration failed", zap.Error(err))
		}
		zlog.Info("Migration completed successfully")
		os.Exit(0)
	case "generate":
		var name string
		if len(os.Args) < 3 {
			fmt.Print("Enter migration name: ")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("No name provided. Aborting.")
				os.Exit(1)
			}
			name = strings.TrimSpace(input)
			if name == "" {
				fmt.Println("No name provided. Aborting.")
				os.Exit(1)
			}
		} else {
			name = os.Args[2]
		}
		name = strings.ReplaceAll(strings.ToLower(name), " ", "_")
		timestamp := time.Now().Format("20060102150405")
		filename := fmt.Sprintf("internal/db/migrations/%s_%s.go", timestamp, name)
		template := `package migrations

import (
	"context"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(
		func(ctx context.Context, db *bun.DB) error {
			// TODO: Write migration up logic here
			return nil
		},
		func(ctx context.Context, db *bun.DB) error {
			// TODO: Write migration down logic here
			return nil
		},
	)
}
`
		err := os.WriteFile(filename, []byte(template), 0644)
		if err != nil {
			fmt.Printf("Failed to create migration file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created migration: %s\n", filename)
		os.Exit(0)
	default:
		fmt.Println("Unknown action. Usage: migrate [migrate|generate] [name]")
		os.Exit(1)
	}
}
