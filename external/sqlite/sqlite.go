package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

var (
	SqlitePathEnv     = strings.TrimSpace(os.Getenv("SQLITE_PATH"))
	MigrationsPathEnv = strings.TrimSpace(os.Getenv("MIGRATIONS_PATH"))
)

func main() {
	if SqlitePathEnv == "" {
		log.Fatal("SQLITE_PATH env is required")
	}

	if MigrationsPathEnv == "" {
		log.Fatal("MIGRATIONS_PATH env is required")
	}

	// Create a database connection string
	//
	// Consider the following:
	// - Set WAL mode, so readers and writers can access the database concurrently
	// - Set busy timeout, so concurrent writers wait on each other instead of erroring immediately
	// - Enable foreign key checks
	dsn := SqlitePathEnv +
		"?_pragma=journal_mode(WAL)" +
		"&_pragma=busy_timeout(5000)" +
		"&_pragma=foreign_keys(1)"

	database, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}

	migrations, err := os.ReadFile(MigrationsPathEnv)
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.Exec(string(migrations))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database created and migration applied successfully!")
}
