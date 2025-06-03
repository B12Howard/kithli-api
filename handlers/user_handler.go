package handlers

import (
	"encoding/json"
	"kithli-api/firebase"
	"kithli-api/services/user"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetUserDataHandler(service *user.UserService, authClient *firebase.FirebaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		externalID := chi.URLParam(r, "external_id")
		if externalID == "" {
			http.Error(w, `{"error": "external_id is required"}`, http.StatusBadRequest)
			return
		}

		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, `{"error": "Missing authentication token"}`, http.StatusUnauthorized)
			return
		}

		uid, err := user.ExtractUserIDFromToken(authClient, token)
		if err != nil || uid != externalID {
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusForbidden)
			return
		}

		userData, err := service.GetFullUserData(r.Context(), externalID)
		if err != nil {
			http.Error(w, `{"error": "Failed to retrieve user data"}`, http.StatusInternalServerError)
			log.Println("[ERROR]", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userData)
	}
}
