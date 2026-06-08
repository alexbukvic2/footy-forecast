package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository/dbgen"
)

// UserRepository handles persistence of User records.
type UserRepository struct {
	q *dbgen.Queries
}

// NewUserRepository constructs a UserRepository backed by pool.
func NewUserRepository(pool *db.Pool) *UserRepository {
	return &UserRepository{q: dbgen.New(pool)}
}

// Upsert inserts a new user or updates email/display_name on cognito_sub conflict.
// Returns domain.ErrConflict if the email uniqueness constraint is violated
// (which indicates a Cognito misconfiguration — two identities sharing one email).
func (r *UserRepository) Upsert(
	ctx context.Context,
	id uuid.UUID,
	cognitoSub, email, displayName string,
) (domain.User, error) {
	row, err := r.q.UpsertUser(ctx, dbgen.UpsertUserParams{
		ID:          id,
		CognitoSub:  cognitoSub,
		Email:       email,
		DisplayName: displayName,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return domain.User{}, fmt.Errorf("upsert user: %w", domain.ErrConflict)
		}
		return domain.User{}, fmt.Errorf("upsert user: %w", err)
	}
	return userFromRow(row), nil
}

// GetByID fetches a user by UUID.
// Returns domain.ErrNotFound if no row exists.
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user %s: %w", id, domain.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("get user by id: %w", err)
	}
	return userFromRow(row), nil
}

// UpdateDisplayName sets the user's display_name. If the user's current status is
// pending_profile, it is transitioned to active atomically in the same query.
// Returns domain.ErrNotFound if no row matches id.
func (r *UserRepository) UpdateDisplayName(ctx context.Context, id uuid.UUID, displayName string) (domain.User, error) {
	row, err := r.q.UpdateDisplayName(ctx, dbgen.UpdateDisplayNameParams{
		ID:          id,
		DisplayName: displayName,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user %s: %w", id, domain.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("update display name: %w", err)
	}
	return userFromRow(row), nil
}

// GetByCognitoSub fetches a user by their Cognito subject identifier.
// Returns domain.ErrNotFound if no row exists.
func (r *UserRepository) GetByCognitoSub(ctx context.Context, sub string) (domain.User, error) {
	row, err := r.q.GetUserByCognitoSub(ctx, sub)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user sub %q: %w", sub, domain.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("get user by cognito sub: %w", err)
	}
	return userFromRow(row), nil
}

// UpdateTimezone sets the user's timezone and optional silent window.
// Pass nil for silentFrom and silentUntil to clear the window.
// Returns domain.ErrNotFound if no row matches id.
func (r *UserRepository) UpdateTimezone(ctx context.Context, id uuid.UUID, timezone string, silentFrom, silentUntil *domain.TimeOfDay) (domain.User, error) {
	row, err := r.q.UpdateTimezone(ctx, dbgen.UpdateTimezoneParams{
		ID:          id,
		Timezone:    timezone,
		SilentFrom:  timeOfDayToPgTime(silentFrom),
		SilentUntil: timeOfDayToPgTime(silentUntil),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user %s: %w", id, domain.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("update timezone: %w", err)
	}
	return userFromRow(row), nil
}

// userFromRow maps a persistence row to the domain type.
func userFromRow(row dbgen.User) domain.User {
	return domain.User{
		ID:          row.ID,
		CognitoSub:  row.CognitoSub,
		Email:       row.Email,
		DisplayName: row.DisplayName,
		Status:      domain.UserStatus(row.Status),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		Timezone:    row.Timezone,
		SilentFrom:  pgTimeToTimeOfDay(row.SilentFrom),
		SilentUntil: pgTimeToTimeOfDay(row.SilentUntil),
	}
}

// timeOfDayToPgTime converts a *domain.TimeOfDay to a pgtype.Time.
// Nil input produces an invalid (NULL) pgtype.Time.
func timeOfDayToPgTime(tod *domain.TimeOfDay) pgtype.Time {
	if tod == nil {
		return pgtype.Time{Valid: false}
	}
	return pgtype.Time{
		Microseconds: int64(tod.Hour*3600+tod.Minute*60) * 1_000_000,
		Valid:        true,
	}
}

// pgTimeToTimeOfDay converts a pgtype.Time to a *domain.TimeOfDay.
// An invalid (NULL) pgtype.Time produces nil.
func pgTimeToTimeOfDay(t pgtype.Time) *domain.TimeOfDay {
	if !t.Valid {
		return nil
	}
	totalSeconds := t.Microseconds / 1_000_000
	return &domain.TimeOfDay{
		Hour:   int(totalSeconds / 3600),
		Minute: int((totalSeconds % 3600) / 60),
	}
}
