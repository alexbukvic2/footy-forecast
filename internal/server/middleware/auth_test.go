package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/cognito"
	"github.com/alexbukvic2/footy-forecast/internal/server/ctxutil"
)

// fakeValidator is a configurable JWTValidator for tests.
type fakeValidator struct {
	claims cognito.Claims
	err    error
}

func (f *fakeValidator) Validate(_ context.Context, _ string) (cognito.Claims, error) {
	return f.claims, f.err
}

// fakeProvisioner is a configurable UserProvisioner for tests.
type fakeProvisioner struct {
	user domain.User
	err  error
}

func (f *fakeProvisioner) ProvisionFromCognito(_ context.Context, _, _, _ string) (domain.User, error) {
	return f.user, f.err
}

func applyAuth(validator cognito.JWTValidator, provisioner UserProvisioner, r *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	Auth(validator, provisioner)(next).ServeHTTP(rec, r)
	return rec
}

func TestAuth_MissingAuthorizationHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	rec := applyAuth(&fakeValidator{}, &fakeProvisioner{}, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", rec.Code)
	}
}

func TestAuth_MalformedAuthorizationHeader(t *testing.T) {
	tests := []string{
		"",
		"Basic dXNlcjpwYXNz",
		"Bearer",
		"Bearer ",
	}
	for _, hdr := range tests {
		req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
		if hdr != "" {
			req.Header.Set("Authorization", hdr)
		}
		rec := applyAuth(&fakeValidator{}, &fakeProvisioner{}, req)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("header %q: status = %d, want 401", hdr, rec.Code)
		}
	}
}

func TestAuth_ValidatorError_Returns401(t *testing.T) {
	validator := &fakeValidator{err: errors.New("token expired")}
	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req.Header.Set("Authorization", "Bearer sometoken")
	rec := applyAuth(validator, &fakeProvisioner{}, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", rec.Code)
	}
}

func TestAuth_ProvisionError_Returns500(t *testing.T) {
	validator := &fakeValidator{claims: cognito.Claims{Sub: "s", Email: "e@e.com", Name: "N"}}
	provisioner := &fakeProvisioner{err: errors.New("db down")}
	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	rec := applyAuth(validator, provisioner, req)
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

func TestAuth_ValidTokenAndProvision_CallsNext(t *testing.T) {
	wantUser := domain.User{Email: "user@example.com", Status: domain.UserStatusActive}
	validator := &fakeValidator{claims: cognito.Claims{Sub: "s", Email: "user@example.com", Name: "User"}}
	provisioner := &fakeProvisioner{user: wantUser}

	var gotUser domain.User
	var gotOK bool
	rec := httptest.NewRecorder()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUser, gotOK = ctxutil.UserFromCtx(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	Auth(validator, provisioner)(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	if !gotOK {
		t.Fatal("UserFromCtx returned ok=false; middleware did not store user in context")
	}
	if gotUser.Email != wantUser.Email {
		t.Errorf("user.Email = %q, want %q", gotUser.Email, wantUser.Email)
	}
}

func TestAuth_SuspendedUser_Returns403(t *testing.T) {
	validator := &fakeValidator{claims: cognito.Claims{Sub: "s", Email: "e@e.com"}}
	provisioner := &fakeProvisioner{user: domain.User{Status: domain.UserStatusSuspended}}
	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	rec := applyAuth(validator, provisioner, req)
	if rec.Code != http.StatusForbidden {
		t.Errorf("status = %d, want 403 for suspended user", rec.Code)
	}
}

func TestAuth_DisplayNamePreference(t *testing.T) {
	// When Name is present, it should be passed to provisioner.
	var capturedDisplayName string
	provisioner := &fakeProvisioner{user: domain.User{}}
	provisioner.user = domain.User{}

	captureProvisioner := &captureProvisionerFn{
		fn: func(_ context.Context, _, _, displayName string) (domain.User, error) {
			capturedDisplayName = displayName
			return domain.User{Status: domain.UserStatusActive}, nil
		},
	}

	validator := &fakeValidator{claims: cognito.Claims{Sub: "s", Email: "e@e.com", Name: "Alice", GivenName: "Al"}}
	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req.Header.Set("Authorization", "Bearer tok")
	_ = provisioner
	applyAuth(validator, captureProvisioner, req)
	if capturedDisplayName != "Alice" {
		t.Errorf("displayName = %q, want Alice (Name takes priority over GivenName)", capturedDisplayName)
	}

	// When Name is empty, GivenName should be used.
	validator2 := &fakeValidator{claims: cognito.Claims{Sub: "s", Email: "e@e.com", Name: "", GivenName: "Bob"}}
	capturedDisplayName = ""
	req2 := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req2.Header.Set("Authorization", "Bearer tok")
	applyAuth(validator2, captureProvisioner, req2)
	if capturedDisplayName != "Bob" {
		t.Errorf("displayName = %q, want Bob (GivenName used when Name is empty)", capturedDisplayName)
	}
}

type captureProvisionerFn struct {
	fn func(ctx context.Context, sub, email, displayName string) (domain.User, error)
}

func (c *captureProvisionerFn) ProvisionFromCognito(ctx context.Context, sub, email, displayName string) (domain.User, error) {
	return c.fn(ctx, sub, email, displayName)
}
