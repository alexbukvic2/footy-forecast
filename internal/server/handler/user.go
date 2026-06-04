package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/ctxutil"
)

// UserService is the subset of service.UserService that the User handler needs.
type UserService interface {
	UpdateDisplayName(ctx context.Context, userID uuid.UUID, displayName string) (domain.User, error)
}

// User is the HTTP handler for the /users resource.
type User struct {
	logger *slog.Logger
	svc    UserService
}

// NewUser constructs a User handler.
func NewUser(logger *slog.Logger, svc UserService) *User {
	return &User{logger: logger, svc: svc}
}

// userResponse is the JSON shape returned to clients for the caller's own profile.
// cognito_sub is intentionally omitted — it is an internal implementation detail.
type userResponse struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toUserResponse(u domain.User) userResponse {
	return userResponse{
		ID:          u.ID.String(),
		Email:       u.Email,
		DisplayName: u.DisplayName,
		Status:      string(u.Status),
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

// Me handles GET /users/me.
// The auth middleware guarantees a User is in context before this handler runs.
func (h *User) Me(w http.ResponseWriter, r *http.Request) {
	u, ok := ctxutil.UserFromCtx(r.Context())
	if !ok {
		h.logger.ErrorContext(r.Context(), "user not in context — auth middleware not applied")
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}
	writeJSON(w, http.StatusOK, toUserResponse(u))
}

type completeProfileRequest struct {
	DisplayName string `json:"display_name"`
}

// CompleteProfile handles PATCH /users/me.
// Accepts { "display_name": "..." }, validates it, persists it, and transitions
// the user from pending_profile to active if that is the current status.
func (h *User) CompleteProfile(w http.ResponseWriter, r *http.Request) {
	u, ok := ctxutil.UserFromCtx(r.Context())
	if !ok {
		h.logger.ErrorContext(r.Context(), "user not in context — auth middleware not applied")
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	var req completeProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	updated, err := h.svc.UpdateDisplayName(r.Context(), u.ID, req.DisplayName)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	writeJSON(w, http.StatusOK, toUserResponse(updated))
}
