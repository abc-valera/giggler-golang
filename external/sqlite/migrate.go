package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	_ "embed"

	_ "modernc.org/sqlite"
)

var (
	sqlitePathEnv = strings.TrimSpace(os.Getenv("SQLITE_PATH"))

	//go:embed migrations.sql
	migrations string
)

func main() {
	if sqlitePathEnv == "" {
		log.Fatal("SQLITE_PATH env is required")
	}

	// Create the directory if it doesn't exist
	if _, err := os.Stat(path.Dir(sqlitePathEnv)); os.IsNotExist(err) {
		if err := os.Mkdir(path.Dir(sqlitePathEnv), 0o755); err != nil {
			log.Fatal(err)
		}
	}

	// Create a database connection string
	//
	// Consider the following:
	// - Set WAL mode, so readers and writers can access the database concurrently
	// - Set busy timeout, so concurrent writers wait on each other instead of erroring immediately
	// - Enable foreign key checks
	dsn := sqlitePathEnv +
		"?_pragma=journal_mode(WAL)" +
		"&_pragma=busy_timeout(5000)" +
		"&_pragma=foreign_keys(1)"

	database, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.Exec(migrations)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database created and migration applied successfully!")
}
