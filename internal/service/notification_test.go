package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

// fakePushTokenRepo is a hand-written fake for PushTokenRepo.
type fakePushTokenRepo struct {
	tokens map[string]domain.PushToken
	err    error
}

func newFakePushTokenRepo() *fakePushTokenRepo {
	return &fakePushTokenRepo{tokens: make(map[string]domain.PushToken)}
}

func (f *fakePushTokenRepo) UpsertToken(_ context.Context, userID uuid.UUID, token string) (domain.PushToken, error) {
	if f.err != nil {
		return domain.PushToken{}, f.err
	}
	pt := domain.PushToken{
		ID:        uuid.New().String(),
		UserID:    userID.String(),
		Token:     token,
		CreatedAt: time.Now(),
	}
	f.tokens[token] = pt
	return pt, nil
}

func (f *fakePushTokenRepo) DeleteToken(_ context.Context, token string) error {
	if f.err != nil {
		return f.err
	}
	delete(f.tokens, token)
	return nil
}

// fakeNotificationPrefRepo is a hand-written fake for NotificationPrefRepo.
type fakeNotificationPrefRepo struct {
	prefs map[string]domain.NotificationPreference // key: userID+type
	err   error
}

func newFakeNotificationPrefRepo() *fakeNotificationPrefRepo {
	return &fakeNotificationPrefRepo{prefs: make(map[string]domain.NotificationPreference)}
}

func prefKey(userID uuid.UUID, typ domain.NotificationType) string {
	return userID.String() + ":" + string(typ)
}

func (f *fakeNotificationPrefRepo) GetForUser(_ context.Context, userID uuid.UUID) ([]domain.NotificationPreference, error) {
	if f.err != nil {
		return nil, f.err
	}
	var out []domain.NotificationPreference
	for _, p := range f.prefs {
		if p.UserID == userID.String() {
			out = append(out, p)
		}
	}
	return out, nil
}

func (f *fakeNotificationPrefRepo) Upsert(_ context.Context, userID uuid.UUID, typ domain.NotificationType, enabled bool) (domain.NotificationPreference, error) {
	if f.err != nil {
		return domain.NotificationPreference{}, f.err
	}
	p := domain.NotificationPreference{
		UserID:  userID.String(),
		Type:    typ,
		Enabled: enabled,
	}
	f.prefs[prefKey(userID, typ)] = p
	return p, nil
}

func TestNotificationService_RegisterToken(t *testing.T) {
	t.Parallel()

	tokenRepo := newFakePushTokenRepo()
	prefRepo := newFakeNotificationPrefRepo()
	svc := NewNotificationService(tokenRepo, prefRepo)

	userID := uuid.New()
	token := "ExponentPushToken[abc]"

	pt, err := svc.RegisterToken(context.Background(), userID, token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pt.Token != token {
		t.Errorf("token = %q, want %q", pt.Token, token)
	}
	if pt.UserID != userID.String() {
		t.Errorf("user_id = %q, want %q", pt.UserID, userID.String())
	}
}

func TestNotificationService_RegisterToken_Empty(t *testing.T) {
	t.Parallel()

	svc := NewNotificationService(newFakePushTokenRepo(), newFakeNotificationPrefRepo())

	_, err := svc.RegisterToken(context.Background(), uuid.New(), "")
	if !errors.Is(err, domain.ErrInvalid) {
		t.Errorf("expected ErrInvalid, got %v", err)
	}
}

func TestNotificationService_DeleteToken(t *testing.T) {
	t.Parallel()

	tokenRepo := newFakePushTokenRepo()
	svc := NewNotificationService(tokenRepo, newFakeNotificationPrefRepo())

	userID := uuid.New()
	token := "ExponentPushToken[xyz]"
	_, _ = tokenRepo.UpsertToken(context.Background(), userID, token)

	if err := svc.DeleteToken(context.Background(), token); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := tokenRepo.tokens[token]; ok {
		t.Error("token should have been deleted")
	}
}

func TestNotificationService_GetPreferences_Defaults(t *testing.T) {
	t.Parallel()

	svc := NewNotificationService(newFakePushTokenRepo(), newFakeNotificationPrefRepo())
	userID := uuid.New()

	prefs, err := svc.GetPreferences(context.Background(), userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should materialise defaults for AllNotificationTypes.
	if len(prefs) != len(domain.AllNotificationTypes) {
		t.Errorf("expected %d defaults, got %d", len(domain.AllNotificationTypes), len(prefs))
	}
	for _, p := range prefs {
		if !p.Enabled {
			t.Errorf("default preference for %q should be enabled", p.Type)
		}
	}
}

func TestNotificationService_GetPreferences_WithStoredRow(t *testing.T) {
	t.Parallel()

	prefRepo := newFakeNotificationPrefRepo()
	svc := NewNotificationService(newFakePushTokenRepo(), prefRepo)
	userID := uuid.New()

	// Pre-store a disabled matchday preference.
	prefRepo.prefs[prefKey(userID, domain.NotificationTypeMatchday)] = domain.NotificationPreference{
		UserID:  userID.String(),
		Type:    domain.NotificationTypeMatchday,
		Enabled: false,
	}

	prefs, err := svc.GetPreferences(context.Background(), userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	prefsMap := make(map[domain.NotificationType]bool)
	for _, p := range prefs {
		prefsMap[p.Type] = p.Enabled
	}

	if prefsMap[domain.NotificationTypeMatchday] != false {
		t.Error("matchday should be disabled (stored row)")
	}
	if prefsMap[domain.NotificationTypePreMatch] != true {
		t.Error("pre_match should be enabled (default)")
	}
}

func TestNotificationService_UpdatePreference_Valid(t *testing.T) {
	t.Parallel()

	svc := NewNotificationService(newFakePushTokenRepo(), newFakeNotificationPrefRepo())
	userID := uuid.New()

	pref, err := svc.UpdatePreference(context.Background(), userID, domain.NotificationTypeMatchday, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pref.Type != domain.NotificationTypeMatchday {
		t.Errorf("type = %q, want %q", pref.Type, domain.NotificationTypeMatchday)
	}
	if pref.Enabled {
		t.Error("expected enabled=false")
	}
}

func TestNotificationService_UpdatePreference_InvalidType(t *testing.T) {
	t.Parallel()

	svc := NewNotificationService(newFakePushTokenRepo(), newFakeNotificationPrefRepo())

	_, err := svc.UpdatePreference(context.Background(), uuid.New(), domain.NotificationType("unknown"), true)
	if !errors.Is(err, domain.ErrInvalid) {
		t.Errorf("expected ErrInvalid, got %v", err)
	}
}

func TestNotificationService_UpdatePreference_TournamentReminderValid(t *testing.T) {
	t.Parallel()

	svc := NewNotificationService(newFakePushTokenRepo(), newFakeNotificationPrefRepo())

	pref, err := svc.UpdatePreference(context.Background(), uuid.New(), domain.NotificationTypeTournamentReminder, false)
	if err != nil {
		t.Fatalf("tournament_reminder should be a valid type: %v", err)
	}
	if pref.Type != domain.NotificationTypeTournamentReminder {
		t.Errorf("type = %q, want tournament_reminder", pref.Type)
	}
}
