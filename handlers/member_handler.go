package handlers

import (
	"encoding/json"
	"kithli-api/models"
	"kithli-api/services/member"
	"log"
	"net/http"
)

func UpdateMemberHandler(service *member.MemberService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.MemberRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		memberID, err := service.UpdateMember(ctx, req)

		if err != nil {
			log.Println("[ERROR] Updating member:", err)
			http.Error(w, `{"error": "Failed to update member"}`, http.StatusInternalServerError)
			return
		}

		response := models.UpdateMemberResponse{
			MemberID: memberID,
			Message:  "Member updated successfully",
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		
	}
}

func CreateMemberHandler(service *member.MemberService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.MemberRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		memberID, err := service.CreateMember(ctx, req)

		if err != nil {
			log.Println("[ERROR] Creating member:", err)
			http.Error(w, `{"error": "Failed to create member"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"member_id": memberID,
			"message":   "Member created successfully",
		})
	}
}
