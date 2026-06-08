package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

// PushTokenRepository handles persistence of push_tokens rows.
type PushTokenRepository struct {
	pool *db.Pool
}

// NewPushTokenRepository constructs a PushTokenRepository.
func NewPushTokenRepository(pool *db.Pool) *PushTokenRepository {
	return &PushTokenRepository{pool: pool}
}

// UpsertToken inserts or reassigns a push token to the given user.
// If the token already exists for a different user (device re-used), it is
// reassigned to userID.
func (r *PushTokenRepository) UpsertToken(ctx context.Context, userID uuid.UUID, token string) (domain.PushToken, error) {
	const q = `
INSERT INTO push_tokens (user_id, token)
VALUES ($1, $2)
ON CONFLICT (token) DO UPDATE SET user_id = EXCLUDED.user_id
RETURNING id, user_id, token, created_at`

	var id uuid.UUID
	var uid uuid.UUID
	var tok string
	var createdAt time.Time

	err := r.pool.QueryRow(ctx, q, userID, token).Scan(&id, &uid, &tok, &createdAt)
	if err != nil {
		return domain.PushToken{}, fmt.Errorf("upsert push token: %w", err)
	}

	return domain.PushToken{
		ID:        id.String(),
		UserID:    uid.String(),
		Token:     tok,
		CreatedAt: createdAt,
	}, nil
}

// DeleteToken deletes a push token (idempotent — no error if the token is missing).
func (r *PushTokenRepository) DeleteToken(ctx context.Context, token string) error {
	const q = `DELETE FROM push_tokens WHERE token = $1`
	_, err := r.pool.Exec(ctx, q, token)
	if err != nil {
		return fmt.Errorf("delete push token: %w", err)
	}
	return nil
}
