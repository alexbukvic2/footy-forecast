package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/ctxutil"
)

// NotificationService is the subset of service.NotificationService the handler needs.
type NotificationService interface {
	RegisterToken(ctx context.Context, userID uuid.UUID, token string) (domain.PushToken, error)
	DeleteToken(ctx context.Context, token string) error
	GetPreferences(ctx context.Context, userID uuid.UUID) ([]domain.NotificationPreference, error)
	UpdatePreference(ctx context.Context, userID uuid.UUID, typ domain.NotificationType, enabled bool) (domain.NotificationPreference, error)
}

// Notification is the HTTP handler for push token and preference endpoints.
type Notification struct {
	logger *slog.Logger
	svc    NotificationService
}

// NewNotification constructs a Notification handler.
func NewNotification(logger *slog.Logger, svc NotificationService) *Notification {
	return &Notification{logger: logger, svc: svc}
}

// RegisterToken handles POST /users/me/push-tokens.
func (h *Notification) RegisterToken(w http.ResponseWriter, r *http.Request) {
	u, ok := ctxutil.UserFromCtx(r.Context())
	if !ok {
		h.logger.ErrorContext(r.Context(), "user not in context")
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	pt, err := h.svc.RegisterToken(r.Context(), u.ID, req.Token)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": pt.ID})
}

// DeleteToken handles DELETE /users/me/push-tokens/{token}.
func (h *Notification) DeleteToken(w http.ResponseWriter, r *http.Request) {
	_, ok := ctxutil.UserFromCtx(r.Context())
	if !ok {
		h.logger.ErrorContext(r.Context(), "user not in context")
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	token := r.PathValue("token")
	if err := h.svc.DeleteToken(r.Context(), token); err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type notificationPreferenceResponse struct {
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}

// GetPreferences handles GET /users/me/notification-preferences.
func (h *Notification) GetPreferences(w http.ResponseWriter, r *http.Request) {
	u, ok := ctxutil.UserFromCtx(r.Context())
	if !ok {
		h.logger.ErrorContext(r.Context(), "user not in context")
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	prefs, err := h.svc.GetPreferences(r.Context(), u.ID)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	resp := make([]notificationPreferenceResponse, 0, len(prefs))
	for _, p := range prefs {
		resp = append(resp, notificationPreferenceResponse{Type: string(p.Type), Enabled: p.Enabled})
	}
	writeJSON(w, http.StatusOK, resp)
}

// UpdatePreference handles PUT /users/me/notification-preferences/{type}.
func (h *Notification) UpdatePreference(w http.ResponseWriter, r *http.Request) {
	u, ok := ctxutil.UserFromCtx(r.Context())
	if !ok {
		h.logger.ErrorContext(r.Context(), "user not in context")
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	typ := domain.NotificationType(r.PathValue("type"))

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	pref, err := h.svc.UpdatePreference(r.Context(), u.ID, typ, req.Enabled)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	writeJSON(w, http.StatusOK, notificationPreferenceResponse{Type: string(pref.Type), Enabled: pref.Enabled})
}
