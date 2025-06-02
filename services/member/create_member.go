package member

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"kithli-api/models"
	"kithli-api/repositories"
	"log"
	"net/http"
)

// Learning note: this is basically our class with dependency injection
type MemberService struct {
	Repo repositories.IMemberRepository
}

// Learning note: this is a constructor
func NewMemberService(repo repositories.IMemberRepository) *MemberService {
	return &MemberService{Repo: repo}
}

func (memberService *MemberService) CreateMember(ctx context.Context, req models.MemberRequest) (int, error) {
	if req.UID == "" || req.MyHeadline == "" {
		return 0, errors.New("UID and myHeadline are required")
	}

	// Step 1: Insert address
	address := &models.Address{
		Street:     req.StreetAddress,
		City:       req.City,
		State:      req.State,
		PostalCode: req.PostalCode,
		AptNumber:  req.AptNumber,
	}

	addressID, err := memberService.Repo.InsertAddress(ctx, address)
	if err != nil {
		return 0, err
	}

	// Step 2: Insert member
	member := &models.Member{
		MyHeadline:            req.MyHeadline,
		AboutMe:               req.AboutMe,
		AdditionalInformation: req.AdditionalInformation,
	}

	memberID, err := memberService.Repo.InsertMember(ctx, member)
	if err != nil {
		return 0, err
	}

	// Step 3: Link member to address
	err = memberService.Repo.LinkAddress(ctx, memberID, addressID)
	if err != nil {
		return 0, err
	}

	// Step 4: Update user with member ID
	err = memberService.Repo.UpdateUserMember(ctx, req.UID, memberID)
	if err != nil {
		return 0, err
	}

	return memberID, nil
}


func CreateMemberHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.MemberRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		// Start transaction
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			http.Error(w, `{"error": "Database transaction error"}`, http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		// Create the repository + service
		repo := repositories.NewMemberRepository(tx)            // your real DB implementation
		service := NewMemberService(repo)                 // inject the repo
		memberID, err := service.CreateMember(ctx, req)   // run the logic

		if err != nil {
			log.Println("[ERROR] Creating member:", err)
			http.Error(w, `{"error": "Failed to create member"}`, http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			log.Println("[ERROR] Commit failed:", err)
			http.Error(w, `{"error": "Failed to commit transaction"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"member_id": memberID,
			"message":   "Member created successfully",
		})
	}
}