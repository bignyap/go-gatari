package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {

	dbURL := os.Getenv("DB_DSN")
	if dbURL == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	schemaDir := os.Getenv("DB_SCHEMA_DIR")
	if schemaDir == "" {
		log.Fatal("DB_SCHEMA_DIR environment variable is required")
	}

	db, err := goose.OpenDBWithDriver("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	if err := goose.Up(db, schemaDir); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
}
