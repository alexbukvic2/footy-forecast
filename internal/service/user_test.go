package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

// ---------- deriveDisplayName tests ----------

func TestDeriveDisplayName(t *testing.T) {
	tests := []struct {
		name           string
		rawDisplayName string
		email          string
		want           string
	}{
		{
			name:           "non-blank raw name wins",
			rawDisplayName: "Alice",
			email:          "alice@example.com",
			want:           "Alice",
		},
		{
			name:           "whitespace-only raw name falls back to email prefix",
			rawDisplayName: "   ",
			email:          "alice@example.com",
			want:           "alice",
		},
		{
			name:           "blank raw name falls back to email prefix",
			rawDisplayName: "",
			email:          "bob@example.com",
			want:           "bob",
		},
		{
			name:           "email with no @ uses full email",
			rawDisplayName: "",
			email:          "noemail",
			want:           "noemail",
		},
		{
			name:           "long name truncated to 50 runes",
			rawDisplayName: strings.Repeat("a", 60),
			email:          "x@y.com",
			want:           strings.Repeat("a", 50),
		},
		{
			name:           "unicode truncated at code-point boundary",
			rawDisplayName: strings.Repeat("á", 60),
			email:          "x@y.com",
			want:           strings.Repeat("á", 50),
		},
		{
			name:           "raw name trimmed before use",
			rawDisplayName: "  Carol  ",
			email:          "carol@example.com",
			want:           "Carol",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := deriveDisplayName(tc.rawDisplayName, tc.email)
			if got != tc.want {
				t.Errorf("deriveDisplayName(%q, %q) = %q, want %q", tc.rawDisplayName, tc.email, got, tc.want)
			}
		})
	}
}

// ---------- fakeUserRepo ----------

type fakeUserRepo struct {
	upsertFn            func(ctx context.Context, id uuid.UUID, cognitoSub, email, displayName string) (domain.User, error)
	updateDisplayNameFn func(ctx context.Context, id uuid.UUID, displayName string) (domain.User, error)
	updateTimezoneFn    func(ctx context.Context, id uuid.UUID, timezone string, silentFrom, silentUntil *domain.TimeOfDay) (domain.User, error)
}

func (f *fakeUserRepo) Upsert(ctx context.Context, id uuid.UUID, cognitoSub, email, displayName string) (domain.User, error) {
	return f.upsertFn(ctx, id, cognitoSub, email, displayName)
}

func (f *fakeUserRepo) UpdateDisplayName(ctx context.Context, id uuid.UUID, displayName string) (domain.User, error) {
	if f.updateDisplayNameFn != nil {
		return f.updateDisplayNameFn(ctx, id, displayName)
	}
	return domain.User{}, nil
}

func (f *fakeUserRepo) UpdateTimezone(ctx context.Context, id uuid.UUID, timezone string, silentFrom, silentUntil *domain.TimeOfDay) (domain.User, error) {
	if f.updateTimezoneFn != nil {
		return f.updateTimezoneFn(ctx, id, timezone, silentFrom, silentUntil)
	}
	return domain.User{Timezone: timezone, SilentFrom: silentFrom, SilentUntil: silentUntil}, nil
}

// countingRepo returns a repo whose Upsert always succeeds and counts calls.
func countingRepo(returnUser domain.User) (*fakeUserRepo, *atomic.Int32) {
	var calls atomic.Int32
	repo := &fakeUserRepo{
		upsertFn: func(_ context.Context, _ uuid.UUID, _, _, _ string) (domain.User, error) {
			calls.Add(1)
			return returnUser, nil
		},
	}
	return repo, &calls
}

// ---------- ProvisionFromCognito tests ----------

func TestUserService_ProvisionFromCognito_FirstRequest_HitsDB(t *testing.T) {
	want := domain.User{CognitoSub: "sub1", Email: "user@example.com", DisplayName: "Alice"}
	repo, calls := countingRepo(want)

	svc := NewUserService(repo)
	got, err := svc.ProvisionFromCognito(context.Background(), "sub1", "user@example.com", "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Email != want.Email {
		t.Errorf("email = %q, want %q", got.Email, want.Email)
	}
	if calls.Load() != 1 {
		t.Errorf("upsert called %d times, want 1", calls.Load())
	}
}

func TestUserService_ProvisionFromCognito_CacheHit_SkipsDB(t *testing.T) {
	want := domain.User{CognitoSub: "sub2", Email: "user@example.com"}
	repo, calls := countingRepo(want)

	svc := NewUserService(repo)
	// First request — populates cache.
	if _, err := svc.ProvisionFromCognito(context.Background(), "sub2", "user@example.com", ""); err != nil {
		t.Fatalf("first call: %v", err)
	}
	// Second request — same sub, same email — must not touch DB.
	if _, err := svc.ProvisionFromCognito(context.Background(), "sub2", "user@example.com", ""); err != nil {
		t.Fatalf("second call: %v", err)
	}
	if calls.Load() != 1 {
		t.Errorf("upsert called %d times, want 1 (cache should serve second request)", calls.Load())
	}
}

func TestUserService_ProvisionFromCognito_EmailMismatch_Upserts(t *testing.T) {
	repo, calls := countingRepo(domain.User{CognitoSub: "sub3", Email: "new@example.com"})

	svc := NewUserService(repo)
	// Seed cache with old email.
	svc.cache["sub3"] = domain.User{CognitoSub: "sub3", Email: "old@example.com"}

	got, err := svc.ProvisionFromCognito(context.Background(), "sub3", "new@example.com", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Email != "new@example.com" {
		t.Errorf("email = %q, want new@example.com", got.Email)
	}
	if calls.Load() != 1 {
		t.Errorf("upsert called %d times, want 1 on email mismatch", calls.Load())
	}
}

func TestUserService_ProvisionFromCognito_EmailMismatch_UpdatesCache(t *testing.T) {
	fresh := domain.User{CognitoSub: "sub4", Email: "new@example.com", DisplayName: "Updated"}
	repo, _ := countingRepo(fresh)

	svc := NewUserService(repo)
	svc.cache["sub4"] = domain.User{CognitoSub: "sub4", Email: "old@example.com"}

	if _, err := svc.ProvisionFromCognito(context.Background(), "sub4", "new@example.com", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Subsequent request with new email must hit cache, not DB.
	repo2, calls2 := countingRepo(fresh)
	svc.repo = repo2
	if _, err := svc.ProvisionFromCognito(context.Background(), "sub4", "new@example.com", ""); err != nil {
		t.Fatalf("unexpected error on third call: %v", err)
	}
	if calls2.Load() != 0 {
		t.Errorf("upsert called after email was already synced, want 0")
	}
}

func TestUserService_ProvisionFromCognito_RepoError_NotCached(t *testing.T) {
	repoErr := errors.New("db failure")
	repo := &fakeUserRepo{
		upsertFn: func(_ context.Context, _ uuid.UUID, _, _, _ string) (domain.User, error) {
			return domain.User{}, repoErr
		},
	}
	svc := NewUserService(repo)
	_, err := svc.ProvisionFromCognito(context.Background(), "sub5", "user@example.com", "Alice")
	if !errors.Is(err, repoErr) {
		t.Errorf("err = %v, want to wrap %v", err, repoErr)
	}
	// Cache must not hold a zero-value entry after a failure.
	svc.mu.RLock()
	_, cached := svc.cache["sub5"]
	svc.mu.RUnlock()
	if cached {
		t.Error("failed upsert must not populate cache")
	}
}

func TestUserService_InvalidateUser_ForcesDBOnNextRequest(t *testing.T) {
	want := domain.User{CognitoSub: "sub6", Email: "user@example.com"}
	repo, calls := countingRepo(want)

	svc := NewUserService(repo)
	// Populate cache.
	if _, err := svc.ProvisionFromCognito(context.Background(), "sub6", "user@example.com", ""); err != nil {
		t.Fatalf("first call: %v", err)
	}

	svc.InvalidateUser("sub6")

	// Next request must go to DB.
	if _, err := svc.ProvisionFromCognito(context.Background(), "sub6", "user@example.com", ""); err != nil {
		t.Fatalf("second call: %v", err)
	}
	if calls.Load() != 2 {
		t.Errorf("upsert called %d times, want 2 (once before and once after invalidation)", calls.Load())
	}
}

func TestUserService_ProvisionFromCognito_UpsertReceivesCorrectArgs(t *testing.T) {
	var capturedSub, capturedEmail, capturedDisplayName string
	var capturedID uuid.UUID
	repo := &fakeUserRepo{
		upsertFn: func(_ context.Context, id uuid.UUID, sub, email, displayName string) (domain.User, error) {
			capturedID = id
			capturedSub = sub
			capturedEmail = email
			capturedDisplayName = displayName
			return domain.User{CognitoSub: sub, Email: email, DisplayName: displayName}, nil
		},
	}
	svc := NewUserService(repo)
	if _, err := svc.ProvisionFromCognito(context.Background(), "sub7", "user@example.com", "Alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedID == (uuid.UUID{}) {
		t.Error("expected non-zero uuid")
	}
	if capturedSub != "sub7" {
		t.Errorf("sub = %q, want sub7", capturedSub)
	}
	if capturedEmail != "user@example.com" {
		t.Errorf("email = %q, want user@example.com", capturedEmail)
	}
	if capturedDisplayName != "Alice" {
		t.Errorf("displayName = %q, want Alice", capturedDisplayName)
	}
}

// ---------- UpdateDisplayName tests ----------

func TestUserService_UpdateDisplayName_HappyPath(t *testing.T) {
	id, _ := uuid.NewV7()
	want := domain.User{ID: id, CognitoSub: "sub-upd", Email: "u@example.com", DisplayName: "Alice", Status: domain.UserStatusActive}
	repo := &fakeUserRepo{
		upsertFn: func(_ context.Context, _ uuid.UUID, _, _, _ string) (domain.User, error) {
			return domain.User{}, nil
		},
		updateDisplayNameFn: func(_ context.Context, gotID uuid.UUID, gotName string) (domain.User, error) {
			if gotID != id {
				return domain.User{}, fmt.Errorf("unexpected id %v", gotID)
			}
			if gotName != "Alice" {
				return domain.User{}, fmt.Errorf("unexpected name %q", gotName)
			}
			return want, nil
		},
	}
	svc := NewUserService(repo)
	got, err := svc.UpdateDisplayName(context.Background(), id, "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.DisplayName != "Alice" {
		t.Errorf("display_name = %q, want Alice", got.DisplayName)
	}
}

func TestUserService_UpdateDisplayName_TrimsWhitespace(t *testing.T) {
	id, _ := uuid.NewV7()
	repo := &fakeUserRepo{
		upsertFn: func(_ context.Context, _ uuid.UUID, _, _, _ string) (domain.User, error) {
			return domain.User{}, nil
		},
		updateDisplayNameFn: func(_ context.Context, _ uuid.UUID, gotName string) (domain.User, error) {
			return domain.User{CognitoSub: "s", DisplayName: gotName}, nil
		},
	}
	svc := NewUserService(repo)
	got, err := svc.UpdateDisplayName(context.Background(), id, "  Bob  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.DisplayName != "Bob" {
		t.Errorf("display_name = %q, want Bob", got.DisplayName)
	}
}

func TestUserService_UpdateDisplayName_BlankName_ReturnsInvalid(t *testing.T) {
	id, _ := uuid.NewV7()
	svc := NewUserService(&fakeUserRepo{
		upsertFn: func(_ context.Context, _ uuid.UUID, _, _, _ string) (domain.User, error) {
			return domain.User{}, nil
		},
	})
	_, err := svc.UpdateDisplayName(context.Background(), id, "   ")
	if !errors.Is(err, domain.ErrInvalid) {
		t.Errorf("err = %v, want domain.ErrInvalid", err)
	}
}

func TestUserService_UpdateDisplayName_TooLong_ReturnsInvalid(t *testing.T) {
	id, _ := uuid.NewV7()
	svc := NewUserService(&fakeUserRepo{
		upsertFn: func(_ context.Context, _ uuid.UUID, _, _, _ string) (domain.User, error) {
			return domain.User{}, nil
		},
	})
	_, err := svc.UpdateDisplayName(context.Background(), id, strings.Repeat("a", 51))
	if !errors.Is(err, domain.ErrInvalid) {
		t.Errorf("err = %v, want domain.ErrInvalid", err)
	}
}

func TestUserService_UpdateDisplayName_InvalidatesCache(t *testing.T) {
	id, _ := uuid.NewV7()
	sub := "sub-invalidate"
	repo, calls := countingRepo(domain.User{ID: id, CognitoSub: sub, Email: "u@example.com"})
	repo.updateDisplayNameFn = func(_ context.Context, _ uuid.UUID, _ string) (domain.User, error) {
		return domain.User{ID: id, CognitoSub: sub, Email: "u@example.com", DisplayName: "New"}, nil
	}
	svc := NewUserService(repo)

	// Populate cache.
	if _, err := svc.ProvisionFromCognito(context.Background(), sub, "u@example.com", ""); err != nil {
		t.Fatalf("provision: %v", err)
	}
	if calls.Load() != 1 {
		t.Fatalf("expected 1 upsert call, got %d", calls.Load())
	}

	// Update display name — should invalidate cache.
	if _, err := svc.UpdateDisplayName(context.Background(), id, "New"); err != nil {
		t.Fatalf("UpdateDisplayName: %v", err)
	}

	// Next ProvisionFromCognito call must hit DB again (cache was invalidated).
	if _, err := svc.ProvisionFromCognito(context.Background(), sub, "u@example.com", ""); err != nil {
		t.Fatalf("second provision: %v", err)
	}
	if calls.Load() != 2 {
		t.Errorf("upsert called %d times after invalidation, want 2", calls.Load())
	}
}

// compile-time interface check
var _ UserRepo = (*fakeUserRepo)(nil)
