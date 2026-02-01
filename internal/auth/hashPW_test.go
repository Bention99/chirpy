package auth

import "testing"

func TestHashPassword_ReturnsHash(t *testing.T) {
	pw := "super-secret-password$123"

	hash, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if hash == "" {
		t.Fatalf("expected non-empty hash")
	}
	if hash == pw {
		t.Fatalf("expected hash to differ from password")
	}
}
