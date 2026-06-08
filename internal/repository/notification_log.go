package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/db"
)

// NotificationLogRepository handles persistence of notification_log rows.
type NotificationLogRepository struct {
	pool *db.Pool
}

// NewNotificationLogRepository constructs a NotificationLogRepository.
func NewNotificationLogRepository(pool *db.Pool) *NotificationLogRepository {
	return &NotificationLogRepository{pool: pool}
}

// Insert records a sent notification. If the (user_id, type, reference_id) triple
// already exists, it does nothing and returns nil (idempotent).
func (r *NotificationLogRepository) Insert(ctx context.Context, userID uuid.UUID, typ, referenceID string) error {
	const q = `
INSERT INTO notification_log (user_id, type, reference_id)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, type, reference_id) DO NOTHING`

	_, err := r.pool.Exec(ctx, q, userID, typ, referenceID)
	if err != nil {
		return fmt.Errorf("insert notification log: %w", err)
	}
	return nil
}
