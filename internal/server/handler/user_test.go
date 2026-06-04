package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/ctxutil"
	"github.com/alexbukvic2/footy-forecast/internal/server/handler"
)

// ---------- fake user service ----------

type fakeUserService struct {
	updateDisplayNameFn func(ctx context.Context, userID uuid.UUID, displayName string) (domain.User, error)
}

func (f *fakeUserService) UpdateDisplayName(ctx context.Context, userID uuid.UUID, displayName string) (domain.User, error) {
	if f.updateDisplayNameFn != nil {
		return f.updateDisplayNameFn(ctx, userID, displayName)
	}
	return domain.User{}, nil
}

// ---------- GET /users/me ----------

func TestUserHandler_Me_UserInContext(t *testing.T) {
	id, _ := uuid.NewV7()
	u := domain.User{
		ID:          id,
		CognitoSub:  "cognito-sub-secret",
		Email:       "user@example.com",
		DisplayName: "Alice",
		Status:      domain.UserStatusActive,
		CreatedAt:   time.Now().UTC().Truncate(time.Second),
		UpdatedAt:   time.Now().UTC().Truncate(time.Second),
	}

	h := handler.NewUser(silentLogger(), &fakeUserService{})
	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req = req.WithContext(ctxutil.WithUser(req.Context(), u))
	rec := httptest.NewRecorder()
	h.Me(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	for _, field := range []string{"id", "email", "display_name", "status", "created_at", "updated_at"} {
		if _, ok := body[field]; !ok {
			t.Errorf("response missing field %q", field)
		}
	}
	if _, ok := body["cognito_sub"]; ok {
		t.Error("response must not contain cognito_sub")
	}
	if got := body["email"]; got != "user@example.com" {
		t.Errorf("email = %v, want user@example.com", got)
	}
	if got := body["display_name"]; got != "Alice" {
		t.Errorf("display_name = %v, want Alice", got)
	}
	if got := body["status"]; got != "active" {
		t.Errorf("status = %v, want active", got)
	}
	if got := body["id"]; got != id.String() {
		t.Errorf("id = %v, want %s", got, id.String())
	}
}

func TestUserHandler_Me_NoUserInContext(t *testing.T) {
	h := handler.NewUser(silentLogger(), &fakeUserService{})
	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	rec := httptest.NewRecorder()
	h.Me(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

// ---------- PATCH /users/me ----------

func TestUserHandler_CompleteProfile_HappyPath(t *testing.T) {
	id, _ := uuid.NewV7()
	u := domain.User{ID: id, Email: "u@example.com", Status: domain.UserStatusPendingProfile}
	updated := domain.User{ID: id, Email: "u@example.com", DisplayName: "Alice", Status: domain.UserStatusActive}

	svc := &fakeUserService{
		updateDisplayNameFn: func(_ context.Context, gotID uuid.UUID, gotName string) (domain.User, error) {
			if gotID != id {
				t.Errorf("userID = %v, want %v", gotID, id)
			}
			if gotName != "Alice" {
				t.Errorf("displayName = %q, want Alice", gotName)
			}
			return updated, nil
		},
	}
	h := handler.NewUser(silentLogger(), svc)

	body, _ := json.Marshal(map[string]string{"display_name": "Alice"})
	req := httptest.NewRequest(http.MethodPatch, "/users/me", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctxutil.WithUser(req.Context(), u))
	rec := httptest.NewRecorder()
	h.CompleteProfile(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", rec.Code, rec.Body)
	}
	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got := resp["display_name"]; got != "Alice" {
		t.Errorf("display_name = %v, want Alice", got)
	}
	if got := resp["status"]; got != "active" {
		t.Errorf("status = %v, want active", got)
	}
}

func TestUserHandler_CompleteProfile_NoUserInContext(t *testing.T) {
	h := handler.NewUser(silentLogger(), &fakeUserService{})
	body, _ := json.Marshal(map[string]string{"display_name": "Alice"})
	req := httptest.NewRequest(http.MethodPatch, "/users/me", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.CompleteProfile(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

func TestUserHandler_CompleteProfile_MalformedJSON(t *testing.T) {
	id, _ := uuid.NewV7()
	u := domain.User{ID: id, Status: domain.UserStatusPendingProfile}
	h := handler.NewUser(silentLogger(), &fakeUserService{})

	req := httptest.NewRequest(http.MethodPatch, "/users/me", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctxutil.WithUser(req.Context(), u))
	rec := httptest.NewRecorder()
	h.CompleteProfile(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestUserHandler_CompleteProfile_ServiceValidationError(t *testing.T) {
	id, _ := uuid.NewV7()
	u := domain.User{ID: id, Status: domain.UserStatusPendingProfile}
	svc := &fakeUserService{
		updateDisplayNameFn: func(_ context.Context, _ uuid.UUID, _ string) (domain.User, error) {
			return domain.User{}, errors.New("display_name cannot be blank: " + domain.ErrInvalid.Error())
		},
	}
	h := handler.NewUser(silentLogger(), svc)

	body, _ := json.Marshal(map[string]string{"display_name": ""})
	req := httptest.NewRequest(http.MethodPatch, "/users/me", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctxutil.WithUser(req.Context(), u))
	rec := httptest.NewRecorder()
	h.CompleteProfile(rec, req)

	if rec.Code != http.StatusInternalServerError {
		// The fake returns a plain error (not wrapping domain.ErrInvalid), so 500 is expected.
		// Real service would return domain.ErrInvalid → 400; tested in service layer tests.
		t.Errorf("status = %d, want 500", rec.Code)
	}
}
