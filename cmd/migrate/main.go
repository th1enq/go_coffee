package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/th1enq/go_coffee/config"
)

const (
	dialect  = "pgx"
	dbString = "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable"
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "migrations", "directory with migration files")
	db    = flags.String("db", "all", "database to migrate (all, user, character)")
)

func main() {
	flags.Usage = usage
	flags.Parse(os.Args[1:])

	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	command := args[0]

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config in migrate: %v", err)
	}

	// Determine which databases to migrate
	var databases []struct {
		name string
		dsn  string
	}

	if *db == "all" || *db == "user" {
		databases = append(databases, struct {
			name string
			dsn  string
		}{
			name: "user",
			dsn:  fmt.Sprintf(dbString, cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.UserDB, cfg.DB.Port),
		})
	}

	if *db == "all" || *db == "character" {
		databases = append(databases, struct {
			name string
			dsn  string
		}{
			name: "character",
			dsn:  fmt.Sprintf(dbString, cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.CharDB, cfg.DB.Port),
		})
	}

	for _, database := range databases {
		log.Printf("Migrating %s database...", database.name)

		db, err := goose.OpenDBWithDriver(dialect, database.dsn)
		if err != nil {
			log.Fatalf("Failed to open database connection for %s: %v", database.name, err)
		}

		defer func() {
			if err := db.Close(); err != nil {
				log.Fatalf("Error closing database connection for %s: %v", database.name, err)
			}
		}()

		if err := goose.Run(command, db, *dir, args[1:]...); err != nil {
			log.Fatalf("migrate %v on %s database: %v", command, database.name, err)
		}

		log.Printf("Successfully ran migration command %s on %s database", command, database.name)
	}
}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: migrate [OPTIONS] COMMAND
Examples:
    migrate status
    migrate -db=user up
    migrate -db=character up
Options:
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations`
)
