package cognito

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	testKid      = "test-kid-1"
	testClientID = "test-client-id"
	testRegion   = "eu-central-1"
)

// testKey generates a fresh RSA key for each test.
func testKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate rsa key: %v", err)
	}
	return key
}

// buildJWKSServer starts an httptest server that returns a JWKS containing key.
// The server records how many times it has been called.
func buildJWKSServer(t *testing.T, kid string, key *rsa.PublicKey) (serverURL string, callCount *int) {
	t.Helper()
	count := 0
	callCount = &count

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		count++
		n := base64.RawURLEncoding.EncodeToString(key.N.Bytes())
		eBig := new(big.Int).SetInt64(int64(key.E))
		e := base64.RawURLEncoding.EncodeToString(eBig.Bytes())
		payload := map[string]any{
			"keys": []map[string]any{
				{"kid": kid, "kty": "RSA", "n": n, "e": e},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
	t.Cleanup(srv.Close)
	return srv.URL, callCount
}

// buildToken creates a signed JWT for testing.
func buildToken(t *testing.T, kid, issuer, clientID string, key *rsa.PrivateKey, opts ...func(*jwt.RegisteredClaims)) string {
	t.Helper()
	claims := &cognitoClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "sub-test",
			Issuer:    issuer,
			Audience:  jwt.ClaimStrings{clientID},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Email:     "user@example.com",
		Name:      "Test User",
		GivenName: "Test",
	}
	for _, opt := range opts {
		opt(&claims.RegisteredClaims)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid
	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return signed
}

// buildValidator wires a Validator pointing at the given JWKS server URL.
func buildValidator(jwksURL, issuer string, allowedClientIDs []string) *Validator {
	return &Validator{
		jwks:             NewJWKSCache(jwksURL, time.Hour),
		issuer:           issuer,
		allowedClientIDs: allowedClientIDs,
	}
}

func TestValidator_Valid(t *testing.T) {
	key := testKey(t)
	jwksURL, _ := buildJWKSServer(t, testKid, &key.PublicKey)
	issuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/pool1", testRegion)
	v := buildValidator(jwksURL, issuer, []string{testClientID})

	tok := buildToken(t, testKid, issuer, testClientID, key)
	claims, err := v.Validate(context.Background(), tok)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.Sub != "sub-test" {
		t.Errorf("Sub = %q, want sub-test", claims.Sub)
	}
	if claims.Email != "user@example.com" {
		t.Errorf("Email = %q, want user@example.com", claims.Email)
	}
	if claims.Name != "Test User" {
		t.Errorf("Name = %q, want Test User", claims.Name)
	}
}

func TestValidator_Expired(t *testing.T) {
	key := testKey(t)
	jwksURL, _ := buildJWKSServer(t, testKid, &key.PublicKey)
	issuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/pool1", testRegion)
	v := buildValidator(jwksURL, issuer, []string{testClientID})

	tok := buildToken(t, testKid, issuer, testClientID, key, func(c *jwt.RegisteredClaims) {
		c.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-2 * time.Minute))
	})
	_, err := v.Validate(context.Background(), tok)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestValidator_WrongSigningKey(t *testing.T) {
	key := testKey(t)
	wrongKey := testKey(t)
	// JWKS has wrongKey, token signed with key
	jwksURL, _ := buildJWKSServer(t, testKid, &wrongKey.PublicKey)
	issuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/pool1", testRegion)
	v := buildValidator(jwksURL, issuer, []string{testClientID})

	tok := buildToken(t, testKid, issuer, testClientID, key)
	_, err := v.Validate(context.Background(), tok)
	if err == nil {
		t.Fatal("expected error for wrong signing key, got nil")
	}
}

func TestValidator_WrongIssuer(t *testing.T) {
	key := testKey(t)
	jwksURL, _ := buildJWKSServer(t, testKid, &key.PublicKey)
	issuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/pool1", testRegion)
	v := buildValidator(jwksURL, issuer, []string{testClientID})

	tok := buildToken(t, testKid, "https://wrong-issuer.example.com", testClientID, key)
	_, err := v.Validate(context.Background(), tok)
	if err == nil {
		t.Fatal("expected error for wrong issuer, got nil")
	}
}

func TestValidator_AudienceNotAllowed(t *testing.T) {
	key := testKey(t)
	jwksURL, _ := buildJWKSServer(t, testKid, &key.PublicKey)
	issuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/pool1", testRegion)
	v := buildValidator(jwksURL, issuer, []string{testClientID})

	tok := buildToken(t, testKid, issuer, "different-client-id", key)
	_, err := v.Validate(context.Background(), tok)
	if err == nil {
		t.Fatal("expected error for disallowed audience, got nil")
	}
}

func TestValidator_UnknownKidTriggersRefetch(t *testing.T) {
	key := testKey(t)
	jwksURL, callCount := buildJWKSServer(t, testKid, &key.PublicKey)
	issuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/pool1", testRegion)

	// Pre-populate the cache with a different kid so the first lookup misses.
	v := buildValidator(jwksURL, issuer, []string{testClientID})
	v.jwks.keys["other-kid"] = &key.PublicKey
	v.jwks.fetchedAt = time.Now() // mark cache as fresh (but wrong kid)

	tok := buildToken(t, testKid, issuer, testClientID, key)
	_, err := v.Validate(context.Background(), tok)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *callCount != 1 {
		t.Errorf("JWKS fetched %d times, want 1", *callCount)
	}
}

func TestValidator_UnknownKidStillAbsentAfterRefetch(t *testing.T) {
	key := testKey(t)
	// Server returns a different kid than what's in the token.
	jwksURL, callCount := buildJWKSServer(t, "completely-different-kid", &key.PublicKey)
	issuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/pool1", testRegion)
	v := buildValidator(jwksURL, issuer, []string{testClientID})

	tok := buildToken(t, testKid, issuer, testClientID, key)
	_, err := v.Validate(context.Background(), tok)
	if err == nil {
		t.Fatal("expected error when kid absent after re-fetch, got nil")
	}
	if *callCount != 1 {
		t.Errorf("JWKS fetched %d times after failed lookup, want exactly 1", *callCount)
	}
}

func TestValidator_MultipleAllowedClientIDs(t *testing.T) {
	key := testKey(t)
	jwksURL, _ := buildJWKSServer(t, testKid, &key.PublicKey)
	issuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/pool1", testRegion)
	v := buildValidator(jwksURL, issuer, []string{"client-a", "client-b"})

	tok := buildToken(t, testKid, issuer, "client-b", key)
	_, err := v.Validate(context.Background(), tok)
	if err != nil {
		t.Fatalf("unexpected error for second allowed client: %v", err)
	}
}

// Verify the interface is satisfied.
var _ JWTValidator = (*Validator)(nil)

// Verify errors package is used.
var _ = errors.New
