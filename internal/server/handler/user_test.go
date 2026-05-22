package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/ctxutil"
	"github.com/alexbukvic2/footy-forecast/internal/server/handler"
)

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

	h := handler.NewUser(silentLogger())
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

	// Required fields.
	for _, field := range []string{"id", "email", "display_name", "status", "created_at", "updated_at"} {
		if _, ok := body[field]; !ok {
			t.Errorf("response missing field %q", field)
		}
	}

	// cognito_sub must never be present.
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
	h := handler.NewUser(silentLogger())
	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	rec := httptest.NewRecorder()
	h.Me(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}
