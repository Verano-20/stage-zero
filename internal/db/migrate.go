package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

// RunMigrations runs database migrations
func RunMigrations(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Set the migrations directory
	goose.SetBaseFS(migrations)

	// Run migrations
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %v", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}
