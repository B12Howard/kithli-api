package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strings"

	"kithli-api/firebase"
	"kithli-api/handlers"
	"kithli-api/repositories"
	"kithli-api/services"
	"kithli-api/services/member"
	"kithli-api/services/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	// "github.com/gorilla/websocket"
)

type ContextKey string

const UserIDKey ContextKey = "user_id"

// AuthMiddleware verifies Firebase ID tokens and extracts the user ID
func AuthMiddleware(authClient *firebase.FirebaseClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, `{"error": "Missing authentication token"}`, http.StatusUnauthorized)
				return
			}

			// Remove 'Bearer ' prefix if present
			tokenString = strings.TrimSpace(strings.Replace(tokenString, "Bearer ", "", 1))

			// Verify Firebase token
			token, err := authClient.FirebaseAuth.VerifyIDToken(r.Context(), tokenString)
			if err != nil {
				http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			// Extract user ID
			uid, ok := token.Claims["user_id"].(string)
			if !ok {
				http.Error(w, `{"error": "User ID missing in token"}`, http.StatusUnauthorized)
				return
			}

			// Store user ID in request context
			ctx := context.WithValue(r.Context(), UserIDKey, uid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func NewRoutes(router *chi.Mux, db *sql.DB, firebaseClient *firebase.FirebaseClient) *chi.Mux {
	hub := services.NewHub()
	go hub.Run()

	router.Use(middleware.Logger)
	log.Println("Middleware Logger Enabled")

	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://localhost:3000"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowedOrigins: []string{"https://app.kithli.com"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	userRepo := repositories.NewUserRepository(db)
	userService := user.NewUserService(userRepo)

	memberService := member.NewMemberService(db)

	router.Group(func(router chi.Router) {
		router.Route("/create-user", func(router chi.Router) {
			router.Post("/", services.CreateUserHandler(db, firebaseClient))
		})
		router.Route("/create-member", func(router chi.Router) {
			router.Post("/", handlers.CreateMemberHandler(memberService))
		})
		router.Route("/update-member", func(router chi.Router) {
			router.Patch("/", handlers.UpdateMemberHandler(memberService))
		})
		router.Route("/getUser", func(router chi.Router) {
			router.Post("/", services.GetUser(db))
		})
		router.Route("/check-membership", func(router chi.Router) {
			router.Post("/", services.CheckUserMembershipHandler(db))
		})
		router.Route("/get-user-data", func(router chi.Router) {
			router.With(AuthMiddleware(firebaseClient)).Get("/{external_id}", handlers.GetUserDataHandler(userService, firebaseClient))
		})

	})
	return router
}
