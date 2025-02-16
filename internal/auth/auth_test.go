package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenerateJWT(t *testing.T) {
	_, err := MakeJWT(uuid.New(), "TestTokenSecret", 10*time.Minute)
	if err != nil {
		t.Errorf("Expected no error but got error: %v", err)
	}
}

func TestValidateJWT(t *testing.T) {
	ID := uuid.New()
	token, _ := MakeJWT(ID, "TestTokenSecret", 10*time.Minute)
	userID, err := ValidateJWT(token, "TestTokenSecret")
	if err != nil {
		t.Errorf("Expected no error but got error: %v", err)
	}

	if userID.String() != ID.String() {
		t.Errorf("ID: %s and userID: %s not equal", ID.String(), userID.String())
	}
}
