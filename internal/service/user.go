package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

// UserRepo is the subset of the repository that UserService needs.
type UserRepo interface {
	Upsert(ctx context.Context, id uuid.UUID, cognitoSub, email, displayName string) (domain.User, error)
	UpdateDisplayName(ctx context.Context, id uuid.UUID, displayName string) (domain.User, error)
	UpdateTimezone(ctx context.Context, id uuid.UUID, timezone string, silentFrom, silentUntil *domain.TimeOfDay) (domain.User, error)
}

// UserService orchestrates user provisioning with an in-process cache.
// The cache is keyed by cognito_sub and has no TTL — entries are invalidated
// explicitly when user data changes (e.g. display_name update, suspension).
// Email is kept in sync by detecting a mismatch between the cached value and
// the claim arriving in the token.
type UserService struct {
	repo  UserRepo
	mu    sync.RWMutex
	cache map[string]domain.User
}

// NewUserService constructs a UserService.
func NewUserService(repo UserRepo) *UserService {
	return &UserService{
		repo:  repo,
		cache: make(map[string]domain.User),
	}
}

// ProvisionFromCognito returns the User for the given Cognito identity,
// writing to the DB only when necessary:
//   - cache miss (first request from this user since process start)
//   - email mismatch (user changed their email in Cognito)
//
// All other requests are served from the in-process cache with no DB round-trip.
func (s *UserService) ProvisionFromCognito(
	ctx context.Context,
	cognitoSub, email, rawDisplayName string,
) (domain.User, error) {
	s.mu.RLock()
	cached, ok := s.cache[cognitoSub]
	s.mu.RUnlock()

	if ok && cached.Email == email {
		return cached, nil
	}

	// Cache miss or email changed — upsert to DB.
	displayName := deriveDisplayName(rawDisplayName, email)

	id, err := uuid.NewV7()
	if err != nil {
		return domain.User{}, fmt.Errorf("generate uuid: %w", err)
	}

	u, err := s.repo.Upsert(ctx, id, cognitoSub, email, displayName)
	if err != nil {
		return domain.User{}, fmt.Errorf("provision user: %w", err)
	}

	s.mu.Lock()
	s.cache[cognitoSub] = u
	s.mu.Unlock()

	return u, nil
}

// InvalidateUser removes the cache entry for cognitoSub.
//
// CONTRACT: every method on UserService that mutates a user record in the DB
// (e.g. UpdateDisplayName, Suspend, Unsuspend) MUST call InvalidateUser before
// returning, so the next request re-reads the authoritative state from the DB.
// Forgetting this will serve stale data — suspended users would remain active.
func (s *UserService) InvalidateUser(cognitoSub string) {
	s.mu.Lock()
	delete(s.cache, cognitoSub)
	s.mu.Unlock()
}

// UpdateDisplayName validates and persists a new display name for the user.
// Whitespace is trimmed; the result must be 1–50 Unicode code points.
// If the user's status is pending_profile it is atomically transitioned to active.
// The in-process cache entry is invalidated so the next request reads fresh state.
func (s *UserService) UpdateDisplayName(ctx context.Context, userID uuid.UUID, displayName string) (domain.User, error) {
	name := strings.TrimSpace(displayName)
	if name == "" {
		return domain.User{}, fmt.Errorf("display_name cannot be blank: %w", domain.ErrInvalid)
	}
	if len([]rune(name)) > 50 {
		return domain.User{}, fmt.Errorf("display_name must be at most 50 characters: %w", domain.ErrInvalid)
	}

	u, err := s.repo.UpdateDisplayName(ctx, userID, name)
	if err != nil {
		return domain.User{}, fmt.Errorf("update display name: %w", err)
	}

	s.InvalidateUser(u.CognitoSub)
	return u, nil
}

// UpdateTimezone validates and persists the user's IANA timezone and optional
// silent window. Both silentFrom and silentUntil must be non-nil together, or
// both nil (to clear the window). silentFrom == silentUntil is rejected.
// The in-process cache entry is invalidated after a successful update.
func (s *UserService) UpdateTimezone(ctx context.Context, userID uuid.UUID, timezone string, silentFrom, silentUntil *domain.TimeOfDay) (domain.User, error) {
	if _, err := time.LoadLocation(timezone); err != nil {
		return domain.User{}, fmt.Errorf("invalid timezone %q: %w", timezone, domain.ErrInvalid)
	}

	u, err := s.repo.UpdateTimezone(ctx, userID, timezone, silentFrom, silentUntil)
	if err != nil {
		return domain.User{}, fmt.Errorf("update timezone: %w", err)
	}
	s.InvalidateUser(u.CognitoSub)
	return u, nil
}

// deriveDisplayName produces a display name from the raw name claim and email.
// Priority:
//  1. TrimSpace(rawDisplayName) if non-empty
//  2. Part of email before @
//  3. Full email if no @ present
//
// Result is truncated to 50 Unicode code points.
func deriveDisplayName(rawDisplayName, email string) string {
	s := strings.TrimSpace(rawDisplayName)
	if s == "" {
		prefix, _, found := strings.Cut(email, "@")
		if found {
			s = prefix
		} else {
			s = email
		}
	}
	runes := []rune(s)
	if len(runes) > 50 {
		runes = runes[:50]
	}
	return string(runes)
}
