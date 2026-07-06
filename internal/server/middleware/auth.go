package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/cognito"
	"github.com/alexbukvic2/footy-forecast/internal/server/ctxutil"
)

// UserProvisioner is the narrow interface the auth middleware needs from UserService.
type UserProvisioner interface {
	ProvisionFromCognito(
		ctx context.Context,
		sub, email, displayName string,
	) (domain.User, error)
}

// Auth returns middleware that validates a Cognito Bearer token on every request.
// Only users with status "active" are allowed through; others receive 403.
//
// Flow:
//  1. Extract Authorization: Bearer <token> — 401 if absent or malformed.
//  2. Validate the token — 401 on error.
//  3. Provision the user via JIT upsert — 500 on error.
//  4. Enforce active status — 403 if suspended or pending_profile.
//  5. Store the user in context and call next.
func Auth(
	validator cognito.JWTValidator,
	users UserProvisioner,
) func(http.Handler) http.Handler {
	return authMiddleware(validator, users, false)
}

// AuthAllowPendingProfile is like Auth but also lets pending_profile users through.
// Use only for endpoints that are part of the profile-completion flow (PATCH /users/me).
func AuthAllowPendingProfile(
	validator cognito.JWTValidator,
	users UserProvisioner,
) func(http.Handler) http.Handler {
	return authMiddleware(validator, users, true)
}

func authMiddleware(
	validator cognito.JWTValidator,
	users UserProvisioner,
	allowPendingProfile bool,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				token, ok := bearerToken(r)
				if !ok {
					writeAuthError(w, http.StatusUnauthorized, "unauthorized")
					return
				}

				claims, err := validator.Validate(r.Context(), token)
				if err != nil {
					slog.WarnContext(
						r.Context(), "token validation failed",
						"request_id", RequestIDFrom(r.Context()),
						"err", err,
					)
					writeAuthError(w, http.StatusUnauthorized, "unauthorized")
					return
				}

				displayName := claims.GivenName
				if displayName == "" {
					displayName = claims.Name
				}

				user, err := users.ProvisionFromCognito(r.Context(), claims.Sub, claims.Email, displayName)
				if err != nil {
					slog.ErrorContext(
						r.Context(), "user provision failed",
						"request_id", RequestIDFrom(r.Context()),
						"err", fmt.Errorf("provision: %w", err),
						"token", token,
						"claims", claims,
					)
					writeAuthError(w, http.StatusInternalServerError, "internal server error")
					return
				}

				switch user.Status {
				case domain.UserStatusActive:
					// always allowed
				case domain.UserStatusPendingProfile:
					if !allowPendingProfile {
						writeAuthError(w, http.StatusForbidden, "forbidden")
						return
					}
				default:
					writeAuthError(w, http.StatusForbidden, "forbidden")
					return
				}

				next.ServeHTTP(w, r.WithContext(ctxutil.WithUser(r.Context(), user)))
			},
		)
	}
}

// bearerToken extracts the token from an "Authorization: Bearer <token>" header.
func bearerToken(r *http.Request) (string, bool) {
	hdr := r.Header.Get("Authorization")
	token, ok := strings.CutPrefix(hdr, "Bearer ")
	if !ok || token == "" {
		return "", false
	}
	return token, true
}

// writeAuthError writes a plain JSON error response.
// Defined inline to avoid importing the handler package (which would be circular).
func writeAuthError(
	w http.ResponseWriter,
	status int,
	msg string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
