package handler

import (
	"context"
	"encoding/json"
	"fmt"
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
	UpdateTimezone(ctx context.Context, userID uuid.UUID, timezone string, silentFrom, silentUntil *domain.TimeOfDay) (domain.User, error)
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
	Timezone    string    `json:"timezone"`
	SilentFrom  *string   `json:"silent_from"`
	SilentUntil *string   `json:"silent_until"`
}

func timeOfDayToString(tod *domain.TimeOfDay) *string {
	if tod == nil {
		return nil
	}
	s := fmt.Sprintf("%02d:%02d", tod.Hour, tod.Minute)
	return &s
}

func toUserResponse(u domain.User) userResponse {
	return userResponse{
		ID:          u.ID.String(),
		Email:       u.Email,
		DisplayName: u.DisplayName,
		Status:      string(u.Status),
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Timezone:    u.Timezone,
		SilentFrom:  timeOfDayToString(u.SilentFrom),
		SilentUntil: timeOfDayToString(u.SilentUntil),
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

// patchMeRequest is the request body for PATCH /users/me.
// All fields are optional (pointer); absent fields are not updated.
type patchMeRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	Timezone    *string `json:"timezone,omitempty"`
	SilentFrom  *string `json:"silent_from"`  // null clears the window
	SilentUntil *string `json:"silent_until"` // null clears the window
}

// parseTimeOfDay parses an "HH:MM" string into a TimeOfDay.
func parseTimeOfDay(s string) (domain.TimeOfDay, error) {
	var h, m int
	n, err := fmt.Sscanf(s, "%d:%02d", &h, &m)
	if err != nil || n != 2 || h < 0 || h > 23 || m < 0 || m > 59 {
		return domain.TimeOfDay{}, fmt.Errorf("invalid time %q: must be HH:MM", s)
	}
	return domain.TimeOfDay{Hour: h, Minute: m}, nil
}

// CompleteProfile handles PATCH /users/me.
// Accepts optional fields: display_name, timezone, silent_from, silent_until.
func (h *User) CompleteProfile(w http.ResponseWriter, r *http.Request) {
	u, ok := ctxutil.UserFromCtx(r.Context())
	if !ok {
		h.logger.ErrorContext(r.Context(), "user not in context — auth middleware not applied")
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	var req patchMeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	current := u

	// Update display name if provided.
	if req.DisplayName != nil {
		updated, err := h.svc.UpdateDisplayName(r.Context(), u.ID, *req.DisplayName)
		if err != nil {
			writeError(w, r, h.logger, err)
			return
		}
		current = updated
	}

	// Update timezone/silent window if timezone is provided.
	if req.Timezone != nil {
		// Validate silent window: both must be present, or both null.
		silentFromPresent := req.SilentFrom != nil
		silentUntilPresent := req.SilentUntil != nil
		if silentFromPresent != silentUntilPresent {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "silent_from and silent_until must both be set or both be null"})
			return
		}

		var silentFrom, silentUntil *domain.TimeOfDay
		if silentFromPresent {
			from, err := parseTimeOfDay(*req.SilentFrom)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
				return
			}
			until, err := parseTimeOfDay(*req.SilentUntil)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
				return
			}
			if from.Hour == until.Hour && from.Minute == until.Minute {
				writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "silent_from and silent_until cannot be equal"})
				return
			}
			silentFrom = &from
			silentUntil = &until
		}

		updated, err := h.svc.UpdateTimezone(r.Context(), u.ID, *req.Timezone, silentFrom, silentUntil)
		if err != nil {
			writeError(w, r, h.logger, err)
			return
		}
		current = updated
	}

	writeJSON(w, http.StatusOK, toUserResponse(current))
}
