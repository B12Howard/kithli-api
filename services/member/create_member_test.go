package member_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"fmt"
	"kithli-api/handlers"
	"kithli-api/models"
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

func strPtr(s string) *string {
	return &s
}

func TestCreateMemberHandler_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := member.NewMemberService(db)

	// Set up the handler using the actual HTTP handler function
	handler := handlers.CreateMemberHandler(service)

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

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req := httptest.NewRequest("POST", "/create-member", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Expect success
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var res map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatal("Failed to parse JSON response:", err)
	}

	if res["message"] != "Member created successfully" {
		t.Errorf("unexpected message: %v", res["message"])
	}

	if res["member_id"] == nil {
		t.Errorf("expected member_id to be returned, got nil")
	}
}
