package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/wiradana/backend/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Use direct connection (migration_dsn) — transaction pooler breaks advisory locks
	m, err := migrate.New("file://db/migrations", cfg.Database.MigrationDSN)
	if err != nil {
		log.Fatalf("failed to init migrate: %v", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("migrations applied successfully")
}
