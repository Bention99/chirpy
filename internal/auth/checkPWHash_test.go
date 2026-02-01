package auth

import "testing"

func TestCheckPasswordHash_CorrectPassword(t *testing.T) {
	pw := "super-secret-password$123____321"

	hash, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}

	ok, err := CheckPasswordHash(pw, hash)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true, got false")
	}
}

func TestCheckPasswordHash_WrongPassword(t *testing.T) {
	pw := "super-secret-password$123"
	wrong := "wrong-password"

	hash, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}

	ok, err := CheckPasswordHash(wrong, hash)
	if err != nil {
		t.Fatalf("expected no error for wrong password, got %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false for wrong password, got true")
	}
}

func TestCheckPasswordHash_InvalidHash(t *testing.T) {
	ok, err := CheckPasswordHash("whatever", "not-a-valid-argon2-hash")
	if err == nil {
		t.Fatalf("expected error for invalid hash, got nil")
	}
	if ok {
		t.Fatalf("expected ok=false for invalid hash, got true")
	}
}
