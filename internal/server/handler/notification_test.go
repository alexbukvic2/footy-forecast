package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/ctxutil"
	"github.com/alexbukvic2/footy-forecast/internal/server/handler"
)

// ---------- fake notification service ----------

type fakeNotificationService struct {
	registerTokenFn    func(ctx context.Context, userID uuid.UUID, token string) (domain.PushToken, error)
	deleteTokenFn      func(ctx context.Context, token string) error
	getPreferencesFn   func(ctx context.Context, userID uuid.UUID) ([]domain.NotificationPreference, error)
	updatePreferenceFn func(ctx context.Context, userID uuid.UUID, typ domain.NotificationType, enabled bool) (domain.NotificationPreference, error)
}

func (f *fakeNotificationService) RegisterToken(ctx context.Context, userID uuid.UUID, token string) (domain.PushToken, error) {
	if f.registerTokenFn != nil {
		return f.registerTokenFn(ctx, userID, token)
	}
	return domain.PushToken{ID: uuid.New().String(), UserID: userID.String(), Token: token, CreatedAt: time.Now()}, nil
}

func (f *fakeNotificationService) DeleteToken(ctx context.Context, token string) error {
	if f.deleteTokenFn != nil {
		return f.deleteTokenFn(ctx, token)
	}
	return nil
}

func (f *fakeNotificationService) GetPreferences(ctx context.Context, userID uuid.UUID) ([]domain.NotificationPreference, error) {
	if f.getPreferencesFn != nil {
		return f.getPreferencesFn(ctx, userID)
	}
	return []domain.NotificationPreference{}, nil
}

func (f *fakeNotificationService) UpdatePreference(ctx context.Context, userID uuid.UUID, typ domain.NotificationType, enabled bool) (domain.NotificationPreference, error) {
	if f.updatePreferenceFn != nil {
		return f.updatePreferenceFn(ctx, userID, typ, enabled)
	}
	return domain.NotificationPreference{UserID: userID.String(), Type: typ, Enabled: enabled}, nil
}

func testUser() domain.User {
	id, _ := uuid.NewV7()
	return domain.User{
		ID:     id,
		Email:  "test@example.com",
		Status: domain.UserStatusActive,
	}
}

// ---------- POST /users/me/push-tokens ----------

func TestNotificationHandler_RegisterToken_HappyPath(t *testing.T) {
	t.Parallel()

	u := testUser()
	h := handler.NewNotification(silentLogger(), &fakeNotificationService{})

	body, _ := json.Marshal(map[string]string{"token": "ExponentPushToken[abc]"})
	req := httptest.NewRequest(http.MethodPost, "/users/me/push-tokens", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctxutil.WithUser(req.Context(), u))
	rec := httptest.NewRecorder()
	h.RegisterToken(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body: %s", rec.Code, rec.Body)
	}
	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["id"] == "" {
		t.Error("expected non-empty id in response")
	}
}

func TestNotificationHandler_RegisterToken_NoUser(t *testing.T) {
	t.Parallel()

	h := handler.NewNotification(silentLogger(), &fakeNotificationService{})
	body, _ := json.Marshal(map[string]string{"token": "ExponentPushToken[abc]"})
	req := httptest.NewRequest(http.MethodPost, "/users/me/push-tokens", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.RegisterToken(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

func TestNotificationHandler_RegisterToken_InvalidBody(t *testing.T) {
	t.Parallel()

	u := testUser()
	h := handler.NewNotification(silentLogger(), &fakeNotificationService{})

	req := httptest.NewRequest(http.MethodPost, "/users/me/push-tokens", bytes.NewReader([]byte("bad")))
	req = req.WithContext(ctxutil.WithUser(req.Context(), u))
	rec := httptest.NewRecorder()
	h.RegisterToken(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestNotificationHandler_RegisterToken_ServiceError(t *testing.T) {
	t.Parallel()

	u := testUser()
	svc := &fakeNotificationService{
		registerTokenFn: func(_ context.Context, _ uuid.UUID, _ string) (domain.PushToken, error) {
			return domain.PushToken{}, fmt.Errorf("token cannot be empty: %w", domain.ErrInvalid)
		},
	}
	h := handler.NewNotification(silentLogger(), svc)

	body, _ := json.Marshal(map[string]string{"token": ""})
	req := httptest.NewRequest(http.MethodPost, "/users/me/push-tokens", bytes.NewReader(body))
	req = req.WithContext(ctxutil.WithUser(req.Context(), u))
	rec := httptest.NewRecorder()
	h.RegisterToken(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

// ---------- DELETE /users/me/push-tokens/{token} ----------

func TestNotificationHandler_DeleteToken_HappyPath(t *testing.T) {
	t.Parallel()

	u := testUser()
	h := handler.NewNotification(silentLogger(), &fakeNotificationService{})

	req := httptest.NewRequest(http.MethodDelete, "/users/me/push-tokens/ExponentPushToken%5Babc%5D", nil)
	req.SetPathValue("token", "ExponentPushToken[abc]")
	req = req.WithContext(ctxutil.WithUser(req.Context(), u))
	rec := httptest.NewRecorder()
	h.DeleteToken(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want 204", rec.Code)
	}
}

func TestNotificationHandler_DeleteToken_NoUser(t *testing.T) {
	t.Parallel()

	h := handler.NewNotification(silentLogger(), &fakeNotificationService{})
	req := httptest.NewRequest(http.MethodDelete, "/users/me/push-tokens/tok", nil)
	req.SetPathValue("token", "tok")
	rec := httptest.NewRecorder()
	h.DeleteToken(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

func TestNotificationHandler_DeleteToken_ServiceError(t *testing.T) {
	t.Parallel()

	u := testUser()
	svc := &fakeNotificationService{
		deleteTokenFn: func(_ context.Context, _ string) error {
			return errors.New("storage failure")
		},
	}
	h := handler.NewNotification(silentLogger(), svc)

	req := httptest.NewRequest(http.MethodDelete, "/users/me/push-tokens/tok", nil)
	req.SetPathValue("token", "tok")
	req = req.WithContext(ctxutil.WithUser(req.Context(), u))
	rec := httptest.NewRecorder()
	h.DeleteToken(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

// ---------- GET /users/me/notification-preferences ----------

func TestNotificationHandler_GetPreferences_HappyPath(t *testing.T) {
	t.Parallel()

	u := testUser()
	svc := &fakeNotificationService{
		getPreferencesFn: func(_ context.Context, _ uuid.UUID) ([]domain.NotificationPreference, error) {
			return []domain.NotificationPreference{
				{Type: domain.NotificationTypeMatchday, Enabled: true},
				{Type: domain.NotificationTypePreMatch, Enabled: false},
			}, nil
		},
	}
	h := handler.NewNotification(silentLogger(), svc)

	req := httptest.NewRequest(http.MethodGet, "/users/me/notification-preferences", nil)
	req = req.WithContext(ctxutil.WithUser(req.Context(), u))
	rec := httptest.NewRecorder()
	h.GetPreferences(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var resp []map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp) != 2 {
		t.Fatalf("expected 2 preferences, got %d", len(resp))
	}
}

func TestNotificationHandler_GetPreferences_NoUser(t *testing.T) {
	t.Parallel()

	h := handler.NewNotification(silentLogger(), &fakeNotificationService{})
	req := httptest.NewRequest(http.MethodGet, "/users/me/notification-preferences", nil)
	rec := httptest.NewRecorder()
	h.GetPreferences(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

func TestNotificationHandler_GetPreferences_ServiceError(t *testing.T) {
	t.Parallel()

	u := testUser()
	svc := &fakeNotificationService{
		getPreferencesFn: func(_ context.Context, _ uuid.UUID) ([]domain.NotificationPreference, error) {
			return nil, errors.New("db error")
		},
	}
	h := handler.NewNotification(silentLogger(), svc)

	req := httptest.NewRequest(http.MethodGet, "/users/me/notification-preferences", nil)
	req = req.WithContext(ctxutil.WithUser(req.Context(), u))
	rec := httptest.NewRecorder()
	h.GetPreferences(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

// ---------- PUT /users/me/notification-preferences/{type} ----------

func TestNotificationHandler_UpdatePreference_HappyPath(t *testing.T) {
	t.Parallel()

	u := testUser()
	h := handler.NewNotification(silentLogger(), &fakeNotificationService{})

	body, _ := json.Marshal(map[string]bool{"enabled": false})
	req := httptest.NewRequest(http.MethodPut, "/users/me/notification-preferences/matchday", bytes.NewReader(body))
	req.SetPathValue("type", "matchday")
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctxutil.WithUser(req.Context(), u))
	rec := httptest.NewRecorder()
	h.UpdatePreference(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", rec.Code, rec.Body)
	}
	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["type"] != "matchday" {
		t.Errorf("type = %v, want matchday", resp["type"])
	}
	if resp["enabled"] != false {
		t.Errorf("enabled = %v, want false", resp["enabled"])
	}
}

func TestNotificationHandler_UpdatePreference_NoUser(t *testing.T) {
	t.Parallel()

	h := handler.NewNotification(silentLogger(), &fakeNotificationService{})
	body, _ := json.Marshal(map[string]bool{"enabled": true})
	req := httptest.NewRequest(http.MethodPut, "/users/me/notification-preferences/matchday", bytes.NewReader(body))
	req.SetPathValue("type", "matchday")
	rec := httptest.NewRecorder()
	h.UpdatePreference(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

func TestNotificationHandler_UpdatePreference_InvalidBody(t *testing.T) {
	t.Parallel()

	u := testUser()
	h := handler.NewNotification(silentLogger(), &fakeNotificationService{})

	req := httptest.NewRequest(http.MethodPut, "/users/me/notification-preferences/matchday", bytes.NewReader([]byte("bad")))
	req.SetPathValue("type", "matchday")
	req = req.WithContext(ctxutil.WithUser(req.Context(), u))
	rec := httptest.NewRecorder()
	h.UpdatePreference(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestNotificationHandler_UpdatePreference_InvalidType(t *testing.T) {
	t.Parallel()

	u := testUser()
	svc := &fakeNotificationService{
		updatePreferenceFn: func(_ context.Context, _ uuid.UUID, typ domain.NotificationType, _ bool) (domain.NotificationPreference, error) {
			return domain.NotificationPreference{}, fmt.Errorf("unknown notification type %q: %w", typ, domain.ErrInvalid)
		},
	}
	h := handler.NewNotification(silentLogger(), svc)

	body, _ := json.Marshal(map[string]bool{"enabled": true})
	req := httptest.NewRequest(http.MethodPut, "/users/me/notification-preferences/unknown_type", bytes.NewReader(body))
	req.SetPathValue("type", "unknown_type")
	req = req.WithContext(ctxutil.WithUser(req.Context(), u))
	rec := httptest.NewRecorder()
	h.UpdatePreference(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}
