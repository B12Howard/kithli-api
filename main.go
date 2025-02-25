package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
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
	// Only load .env in local development (Render provides env variables directly)
	if os.Getenv("RENDER") == "" {
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

	// Determine Firebase config file path
	firebaseConfigPath := "./config/firebase.json"

	// If on Render, decode FIREBASE_CONFIG_BASE64 and write to a temp file
	if os.Getenv("RENDER") != "" {
		firebaseBase64 := os.Getenv("FIREBASE_CONFIG")
		if firebaseBase64 == "" {
			log.Fatal("Missing FIREBASE_CONFIG environment variable")
		}

		decoded, err := base64.StdEncoding.DecodeString(firebaseBase64)
		if err != nil {
			log.Fatalf("Failed to decode Firebase config: %v", err)
		}

		firebaseConfigPath = "/tmp/firebase.json" // Use Render's writable temp directory
		err = ioutil.WriteFile(firebaseConfigPath, decoded, 0644)
		if err != nil {
			log.Fatalf("Failed to write Firebase config file: %v", err)
		}
		fmt.Println("Firebase config written to:", firebaseConfigPath)
	}

	// Initialize Firebase with the correct file path
	firebaseClient, err := firebase.InitFirebase(firebaseConfigPath)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	log.Printf("Running on port %s", port)

	router := chi.NewRouter()
	NewRoutes(router, db, firebaseClient)
	fmt.Println("Ready!")

	http.ListenAndServe(":"+port, router)
}
