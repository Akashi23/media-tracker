package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection string - adjust as needed
	dsn := "host=localhost port=5432 user=postgres password=postgres123 dbname=media_tracker sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Create default user
	defaultUserID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	email := "admin@example.com"
	name := "Admin User"

	// Check if user already exists
	var existingID uuid.UUID
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&existingID)
	if err == nil {
		log.Printf("User %s already exists with ID: %s", email, existingID)
		return
	}

	if err != sql.ErrNoRows {
		log.Fatal("Error checking for existing user:", err)
	}

	// Insert default user
	_, err = db.ExecContext(context.Background(),
		"INSERT INTO users (id, email, name, created_at) VALUES ($1, $2, $3, $4)",
		defaultUserID, email, name, time.Now())

	if err != nil {
		log.Fatal("Failed to insert default user:", err)
	}

	log.Printf("Default user created successfully: %s (%s)", name, email)
	log.Printf("User ID: %s", defaultUserID)
	log.Printf("You can now login with email: %s", email)
}
