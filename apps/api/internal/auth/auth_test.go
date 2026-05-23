package auth

import (
	"testing"
	"time"
)

func TestHashAndCheckPassword(t *testing.T) {
	t.Parallel()

	hash, err := HashPassword("s3cret-password")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if err := CheckPassword(hash, "s3cret-password"); err != nil {
		t.Fatalf("CheckPassword() error = %v", err)
	}

	if err := CheckPassword(hash, "wrong-password"); err == nil {
		t.Fatal("CheckPassword() expected error for wrong password")
	}
}

func TestCreateAndParseToken(t *testing.T) {
	t.Parallel()

	secret := []byte("test-secret")
	token, err := CreateToken(secret, "user-1", "admin", time.Hour)
	if err != nil {
		t.Fatalf("CreateToken() error = %v", err)
	}

	claims, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}

	if claims.UserID != "user-1" {
		t.Fatalf("ParseToken() user id = %q, want %q", claims.UserID, "user-1")
	}

	if claims.Role != "admin" {
		t.Fatalf("ParseToken() role = %q, want %q", claims.Role, "admin")
	}
}
