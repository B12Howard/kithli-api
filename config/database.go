package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" // For loading .env file in local development
	_ "github.com/lib/pq"
)

// NewDb initializes a PostgreSQL database connection
func NewDb() *sql.DB {
	// Load .env only in local development
	if os.Getenv("RENDER") == "" { // Render does not set "RENDER" env, so this ensures .env is only loaded locally
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("No .env file found, using system environment variables.")
		} else {
			log.Println(".env file loaded successfully.")
		}
	}

	// Retrieve environment variables
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	if host == "" || port == "" || user == "" || dbname == "" {
		log.Fatal("Missing required environment variables for database connection.")
	}

	// Create the PostgreSQL connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	// Connect to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	log.Println("Database connection established successfully!")
	return db
}
