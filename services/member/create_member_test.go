package member_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"fmt"
	"kithli-api/models"
	"kithli-api/repositories"
	"kithli-api/services/member"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	root, fpErr := filepath.Abs("../../")
	if fpErr != nil {
		fmt.Printf("Could not resolve project root: %v", fpErr)
	}

	err := godotenv.Load(filepath.Join(root,".env"))
	if err != nil {
		fmt.Println("No .env file found, using system environment variables.")
	}
	
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	fmt.Print("DB String")
	fmt.Print(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}

	return db
}

func TestCreateMemberHandler_Rollback(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Start a transaction to be rolled back
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback() // Ensures test doesn't persist data

	// Wrap the service handler with the mocked transaction
	repo := repositories.NewMemberRepository(tx)
	service := member.NewMemberService(repo)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req models.MemberRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}

		memberID, err := service.CreateMember(r.Context(), req)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)

			return
		}

		json.NewEncoder(w).Encode(map[string]any{
			"member_id": memberID,
			"message":   "success",
		})
	})

	// Create a request payload
	reqBody := models.MemberRequest{
		UID:                   "test-uid-123",
		MyHeadline:            "Test Headline",
		AboutMe:               strPtr("This is about me"),
		PostalCode:            strPtr("94016"),
		StreetAddress:         strPtr("123 Test St"),
		AptNumber:             strPtr("1A"),
		City:                  strPtr("Testville"),
		State:                 strPtr("CA"),
		AdditionalInformation: strPtr("No allergies"),
	}
	
	jsonBody, _ := json.Marshal(reqBody)

	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req := httptest.NewRequest("POST", "/create-member", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var res map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatal("Failed to parse JSON response:", err)
	}

	if res["message"] != "success" {
		t.Errorf("unexpected message: %v", res["message"])
	}
}

func strPtr(s string) *string {
	return &s
}
