package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

// NotificationPreferenceRepository handles persistence of notification_preferences rows.
type NotificationPreferenceRepository struct {
	pool *db.Pool
}

// NewNotificationPreferenceRepository constructs a NotificationPreferenceRepository.
func NewNotificationPreferenceRepository(pool *db.Pool) *NotificationPreferenceRepository {
	return &NotificationPreferenceRepository{pool: pool}
}

// GetForUser returns all preference rows stored for the user. Missing types are not
// included — the caller is responsible for materialising defaults.
func (r *NotificationPreferenceRepository) GetForUser(ctx context.Context, userID uuid.UUID) ([]domain.NotificationPreference, error) {
	const q = `SELECT user_id, type, enabled FROM notification_preferences WHERE user_id = $1`

	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("get notification preferences: %w", err)
	}
	defer rows.Close()

	var prefs []domain.NotificationPreference
	for rows.Next() {
		var uid uuid.UUID
		var typ string
		var enabled bool
		if err := rows.Scan(&uid, &typ, &enabled); err != nil {
			return nil, fmt.Errorf("scan notification preference: %w", err)
		}
		prefs = append(prefs, domain.NotificationPreference{
			UserID:  uid.String(),
			Type:    domain.NotificationType(typ),
			Enabled: enabled,
		})
	}
	return prefs, rows.Err()
}

// Upsert inserts or updates a single preference row.
func (r *NotificationPreferenceRepository) Upsert(ctx context.Context, userID uuid.UUID, typ domain.NotificationType, enabled bool) (domain.NotificationPreference, error) {
	const q = `
INSERT INTO notification_preferences (user_id, type, enabled)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, type) DO UPDATE SET enabled = EXCLUDED.enabled
RETURNING user_id, type, enabled`

	var uid uuid.UUID
	var typStr string
	var en bool
	err := r.pool.QueryRow(ctx, q, userID, string(typ), enabled).Scan(&uid, &typStr, &en)
	if err != nil {
		return domain.NotificationPreference{}, fmt.Errorf("upsert notification preference: %w", err)
	}
	return domain.NotificationPreference{
		UserID:  uid.String(),
		Type:    domain.NotificationType(typStr),
		Enabled: en,
	}, nil
}
