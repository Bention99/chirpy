package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken_Success(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Bearer abc.def.ghi")

	got, err := GetBearerToken(h)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != "abc.def.ghi" {
		t.Fatalf("expected %q, got %q", "abc.def.ghi", got)
	}
}

func TestGetBearerToken_StripsWhitespace(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Bearer   abc.def.ghi   ")

	got, err := GetBearerToken(h)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != "abc.def.ghi" {
		t.Fatalf("expected %q, got %q", "abc.def.ghi", got)
	}
}

func TestGetBearerToken_MissingHeader(t *testing.T) {
	h := http.Header{}

	_, err := GetBearerToken(h)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetBearerToken_WrongPrefix(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Basic abcdef")

	_, err := GetBearerToken(h)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetBearerToken_EmptyToken(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Bearer    ")

	_, err := GetBearerToken(h)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
