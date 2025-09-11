package database

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// ConnectDB initializes a connection to the database using environment variables.

func ConnectDB(ctx context.Context , connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}

	// Ping the database to verify connection
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// Log successful connection
	if os.Getenv("ENVIRONMENT") == "dev" {
		log.Println("Database connected successfully")
	}
	return db, nil
}
