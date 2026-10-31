package auth

import (
	"testing"
)

func TestPasswordHashing(t *testing.T) {
	password := "my-secure-password-123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error hashing password, got %v", err)
	}

	if hash == password {
		t.Fatalf("hash should not equal raw password")
	}

	if !CheckPasswordHash(password, hash) {
		t.Fatalf("expected password to match hash")
	}

	if CheckPasswordHash("wrong-password", hash) {
		t.Fatalf("expected wrong password not to match hash")
	}
}
