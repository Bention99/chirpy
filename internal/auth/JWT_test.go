package auth

import (
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeJWT_SetsExpectedClaimsAndSigns(t *testing.T) {
	userID := uuid.New()
	secret := "super-secret"
	expiresIn := 10 * time.Minute

	before := time.Now().UTC()
	tokenString, err := MakeJWT(userID, secret, expiresIn)
	after := time.Now().UTC()

	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}
	if tokenString == "" {
		t.Fatalf("expected non-empty token string")
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("ParseWithClaims returned error: %v", err)
	}
	if !token.Valid {
		t.Fatalf("expected token to be valid")
	}

	if claims.Issuer != "chirpy" {
		t.Fatalf("expected issuer chirpy, got %q", claims.Issuer)
	}
	if claims.Subject != userID.String() {
		t.Fatalf("expected subject %q, got %q", userID.String(), claims.Subject)
	}
	if claims.IssuedAt == nil {
		t.Fatalf("expected IssuedAt to be set")
	}
	if claims.ExpiresAt == nil {
		t.Fatalf("expected ExpiresAt to be set")
	}

	iat := claims.IssuedAt.Time
	exp := claims.ExpiresAt.Time

	if iat.Before(before.Add(-2*time.Second)) || iat.After(after.Add(2*time.Second)) {
		t.Fatalf("IssuedAt out of expected range: iat=%s before=%s after=%s", iat, before, after)
	}

	wantExpMin := iat.Add(expiresIn - 2*time.Second)
	wantExpMax := iat.Add(expiresIn + 2*time.Second)
	if exp.Before(wantExpMin) || exp.After(wantExpMax) {
		t.Fatalf("ExpiresAt out of expected range: exp=%s want between %s and %s", exp, wantExpMin, wantExpMax)
	}
}

func TestValidateJWT_Success(t *testing.T) {
	userID := uuid.New()
	secret := "super-secret"

	tokenString, err := MakeJWT(userID, secret, 5*time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	got, err := ValidateJWT(tokenString, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}
	if got != userID {
		t.Fatalf("expected %s, got %s", userID, got)
	}
}

func TestValidateJWT_Expired(t *testing.T) {
	userID := uuid.New()
	secret := "super-secret"

	tokenString, err := MakeJWT(userID, secret, -1*time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = ValidateJWT(tokenString, secret)
	if err == nil {
		t.Fatalf("expected error for expired token, got nil")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()

	tokenString, err := MakeJWT(userID, "right-secret", 5*time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = ValidateJWT(tokenString, "wrong-secret")
	if err == nil {
		t.Fatalf("expected error for wrong secret, got nil")
	}
}

func TestValidateJWT_MalformedToken(t *testing.T) {
	_, err := ValidateJWT("not-a-jwt", "secret")
	if err == nil {
		t.Fatalf("expected error for malformed token, got nil")
	}
}

func TestValidateJWT_UnexpectedSigningMethod(t *testing.T) {
	userID := uuid.New()
	secret := "super-secret"

	header := map[string]any{
		"alg": "RS256",
		"typ": "JWT",
	}
	now := time.Now().UTC()
	payload := map[string]any{
		"iss": "chirpy",
		"sub": userID.String(),
		"iat": now.Unix(),
		"exp": now.Add(5 * time.Minute).Unix(),
	}

	tokenString, err := makeUnsignedJWT(header, payload)
	if err != nil {
		t.Fatalf("failed to build token: %v", err)
	}

	_, err = ValidateJWT(tokenString, secret)
	if err == nil {
		t.Fatalf("expected error for unexpected signing method, got nil")
	}
}

func makeUnsignedJWT(header, payload map[string]any) (string, error) {
	hb, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	pb, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	enc := base64.RawURLEncoding
	h := enc.EncodeToString(hb)
	p := enc.EncodeToString(pb)

	sig := "invalidsig"

	return h + "." + p + "." + sig, nil
}
