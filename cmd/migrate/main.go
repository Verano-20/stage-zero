package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Verano-20/stage-zero/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

const (
	dialect = "pgx"
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "migrations", "directory with migration files")
)

func main() {
	// Check flags and get command
	flags.Usage = usage
	flags.Parse(os.Args[1:])

	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	command := args[0]

	config.InitConfig()
	dsn := config.Get().GetDBConnectionString()

	db, err := goose.OpenDBWithDriver(dialect, dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("failed to close database connection: %v", err)
		}
	}()

	// Set the embedded filesystem for migrations
	goose.SetBaseFS(embedMigrations)

	if err := goose.RunContext(context.Background(), command, db, *dir, args[1:]...); err != nil {
		log.Fatalf("failed to run migration command '%s': %v", command, err)
	}

	if command != "status" && command != "version" {
		log.Printf("Migration command '%s' completed successfully", command)
	}
}

func usage() {
	usageText := `Usage: migrate [OPTIONS] COMMAND

Example:
	migrate status

Commands:
    up                   Migrate the DB to the most recent version available
    down                 Roll back the version by 1
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database

Options:`
	fmt.Print(usageText)
	flags.PrintDefaults()
}
