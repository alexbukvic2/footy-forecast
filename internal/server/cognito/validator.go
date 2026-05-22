package cognito

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims holds the identity fields extracted from a validated Cognito ID token.
type Claims struct {
	Sub       string
	Email     string
	Name      string
	GivenName string
}

// JWTValidator validates a raw token string and returns its claims.
type JWTValidator interface {
	Validate(ctx context.Context, tokenString string) (Claims, error)
}

// cognitoClaims extends the standard JWT registered claims with Cognito ID token fields.
type cognitoClaims struct {
	jwt.RegisteredClaims
	TokenUse  string `json:"token_use"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	GivenName string `json:"given_name"`
}

// Validator validates RS256 ID tokens issued by an AWS Cognito User Pool.
type Validator struct {
	jwks             *JWKSCache
	issuer           string
	allowedClientIDs []string
}

// NewValidator creates a validator for the given Cognito User Pool.
// region and userPoolID are used to derive the issuer URL and JWKS endpoint.
func NewValidator(region, userPoolID string, allowedClientIDs []string) *Validator {
	poolURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", region, userPoolID)
	jwksURL := poolURL + "/.well-known/jwks.json"
	return &Validator{
		jwks:             NewJWKSCache(jwksURL, time.Hour),
		issuer:           poolURL,
		allowedClientIDs: allowedClientIDs,
	}
}

// Validate parses and validates tokenString. On success it returns the identity Claims.
//
// Validation checks: RS256 signature, expiry (with 30s leeway), issuer, and audience.
// The audience must contain at least one of the configured allowedClientIDs.
func (v *Validator) Validate(ctx context.Context, tokenString string) (Claims, error) {
	parser := jwt.NewParser(
		jwt.WithIssuer(v.issuer),
		jwt.WithExpirationRequired(),
		jwt.WithLeeway(30*time.Second),
		jwt.WithValidMethods([]string{"RS256"}),
	)

	token, err := parser.ParseWithClaims(tokenString, &cognitoClaims{}, v.keyFunc(ctx))
	if err != nil {
		return Claims{}, fmt.Errorf("parse token: %w", err)
	}

	c, ok := token.Claims.(*cognitoClaims)
	if !ok || !token.Valid {
		return Claims{}, errors.New("invalid token claims")
	}

	if !v.audienceAllowed(c.Audience) {
		return Claims{}, errors.New("token audience not in allowed client IDs")
	}

	if c.TokenUse != "id" {
		return Claims{}, fmt.Errorf("token_use %q: want \"id\"", c.TokenUse)
	}

	return Claims{
		Sub:       c.Subject,
		Email:     c.Email,
		Name:      c.Name,
		GivenName: c.GivenName,
	}, nil
}

// keyFunc returns a jwt.Keyfunc that resolves the signing key from the JWKS cache.
func (v *Validator) keyFunc(ctx context.Context) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		kid, ok := token.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, errors.New("missing kid in token header")
		}
		return v.jwks.Get(ctx, kid)
	}
}

// audienceAllowed reports whether any value in audience matches an allowed client ID.
func (v *Validator) audienceAllowed(audience jwt.ClaimStrings) bool {
	for _, aud := range audience {
		for _, allowed := range v.allowedClientIDs {
			if aud == allowed {
				return true
			}
		}
	}
	return false
}

// ensure Validator satisfies JWTValidator at compile time.
var _ JWTValidator = (*Validator)(nil)
