package service

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

// PushTokenRepo is the subset of the push token repository that NotificationService needs.
type PushTokenRepo interface {
	UpsertToken(ctx context.Context, userID uuid.UUID, token string) (domain.PushToken, error)
	DeleteToken(ctx context.Context, token string) error
}

// NotificationPrefRepo is the subset of the preference repository that NotificationService needs.
type NotificationPrefRepo interface {
	GetForUser(ctx context.Context, userID uuid.UUID) ([]domain.NotificationPreference, error)
	Upsert(ctx context.Context, userID uuid.UUID, typ domain.NotificationType, enabled bool) (domain.NotificationPreference, error)
}

// NotificationService manages push tokens and notification preferences.
type NotificationService struct {
	tokenRepo PushTokenRepo
	prefRepo  NotificationPrefRepo
}

// NewNotificationService constructs a NotificationService.
func NewNotificationService(tokenRepo PushTokenRepo, prefRepo NotificationPrefRepo) *NotificationService {
	return &NotificationService{tokenRepo: tokenRepo, prefRepo: prefRepo}
}

// RegisterToken registers a device push token for the user. The token is upserted:
// if it already exists for another user (device re-used after sign-out) it is reassigned.
func (s *NotificationService) RegisterToken(ctx context.Context, userID uuid.UUID, token string) (domain.PushToken, error) {
	if token == "" {
		return domain.PushToken{}, fmt.Errorf("token cannot be empty: %w", domain.ErrInvalid)
	}
	pt, err := s.tokenRepo.UpsertToken(ctx, userID, token)
	if err != nil {
		return domain.PushToken{}, fmt.Errorf("register token: %w", err)
	}
	return pt, nil
}

// DeleteToken removes a push token (idempotent — no error if missing).
func (s *NotificationService) DeleteToken(ctx context.Context, token string) error {
	if err := s.tokenRepo.DeleteToken(ctx, token); err != nil {
		return fmt.Errorf("delete token: %w", err)
	}
	return nil
}

// GetPreferences returns all preferences for the user. Types with no stored row are
// materialised with their default value (enabled=true) for every type in
// domain.AllNotificationTypes.
func (s *NotificationService) GetPreferences(ctx context.Context, userID uuid.UUID) ([]domain.NotificationPreference, error) {
	rows, err := s.prefRepo.GetForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get preferences: %w", err)
	}

	// Materialise defaults for types not yet in the DB.
	existing := make(map[domain.NotificationType]bool, len(rows))
	for _, r := range rows {
		existing[r.Type] = true
	}
	for _, t := range domain.AllNotificationTypes {
		if !existing[t] {
			rows = append(rows, domain.NotificationPreference{
				UserID:  userID.String(),
				Type:    t,
				Enabled: true,
			})
		}
	}
	return rows, nil
}

// UpdatePreference upserts the enabled state for a single notification type.
// Returns domain.ErrInvalid if the type is not one of the recognised values.
func (s *NotificationService) UpdatePreference(ctx context.Context, userID uuid.UUID, typ domain.NotificationType, enabled bool) (domain.NotificationPreference, error) {
	if !slices.Contains(domain.ValidNotificationTypes, typ) {
		return domain.NotificationPreference{}, fmt.Errorf("unknown notification type %q: %w", typ, domain.ErrInvalid)
	}

	pref, err := s.prefRepo.Upsert(ctx, userID, typ, enabled)
	if err != nil {
		return domain.NotificationPreference{}, fmt.Errorf("update preference: %w", err)
	}
	return pref, nil
}
