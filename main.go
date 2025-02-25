package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"kithli-api/config"
	"kithli-api/firebase"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

type MyHandler struct {
	db *sql.DB
}

func main() {
	// Only load .env in local development (Render will provide env variables directly)
	if os.Getenv("RENDER") == "" { // Render services automatically provide env variables
		if err := godotenv.Load(".env"); err != nil {
			log.Println("No .env file found, using system environment variables.")
		}
	}

	fmt.Println("Loading...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	db := config.NewDb()

	if err := goose.Up(db, "./migrations"); err != nil {
		log.Fatal("Migration failed:", err)
	}

	firebaseClient, err := firebase.InitFirebase("./config/firebase.json")
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	log.Printf("Running on port %s", port)

	router := chi.NewRouter()
	NewRoutes(router, db, firebaseClient)
	fmt.Println("Ready!")

	http.ListenAndServe(":"+port, router)
}
