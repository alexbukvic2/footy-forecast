package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/ctxutil"
)

// User is the HTTP handler for the /users resource.
type User struct {
	logger *slog.Logger
}

// NewUser constructs a User handler.
func NewUser(logger *slog.Logger) *User {
	return &User{logger: logger}
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
