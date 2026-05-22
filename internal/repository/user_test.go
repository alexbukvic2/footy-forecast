//go:build integration

package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

func TestUserRepository(t *testing.T) {
	pool := startPostgres(t)
	repo := repository.NewUserRepository(pool)
	ctx := context.Background()

	newID := func() uuid.UUID {
		id, err := uuid.NewV7()
		if err != nil {
			t.Fatalf("generate uuid: %v", err)
		}
		return id
	}

	t.Run("Upsert creates a new row", func(t *testing.T) {
		id := newID()
		u, err := repo.Upsert(ctx, id, "sub-create", "create@example.com", "Creator")
		if err != nil {
			t.Fatalf("Upsert: %v", err)
		}
		if u.ID != id {
			t.Errorf("id = %v, want %v", u.ID, id)
		}
		if u.Email != "create@example.com" {
			t.Errorf("email = %q, want create@example.com", u.Email)
		}
		if u.DisplayName != "Creator" {
			t.Errorf("display_name = %q, want Creator", u.DisplayName)
		}
		if u.Status != domain.UserStatusActive {
			t.Errorf("status = %q, want active", u.Status)
		}
	})

	t.Run("second Upsert same cognito_sub returns same id", func(t *testing.T) {
		id := newID()
		u1, err := repo.Upsert(ctx, id, "sub-idempotent", "idempotent@example.com", "Name")
		if err != nil {
			t.Fatalf("first Upsert: %v", err)
		}
		u2, err := repo.Upsert(ctx, newID(), "sub-idempotent", "idempotent2@example.com", "Name2")
		if err != nil {
			t.Fatalf("second Upsert: %v", err)
		}
		if u2.ID != u1.ID {
			t.Errorf("id changed: got %v, want %v", u2.ID, u1.ID)
		}
	})

	t.Run("second Upsert refreshes email", func(t *testing.T) {
		id := newID()
		if _, err := repo.Upsert(ctx, id, "sub-email-refresh", "old@example.com", "N"); err != nil {
			t.Fatalf("first Upsert: %v", err)
		}
		u, err := repo.Upsert(ctx, newID(), "sub-email-refresh", "new@example.com", "N")
		if err != nil {
			t.Fatalf("second Upsert: %v", err)
		}
		if u.Email != "new@example.com" {
			t.Errorf("email = %q, want new@example.com", u.Email)
		}
	})

	t.Run("second Upsert does NOT overwrite non-empty display_name", func(t *testing.T) {
		id := newID()
		if _, err := repo.Upsert(ctx, id, "sub-dn-preserve", "dn@example.com", "Original"); err != nil {
			t.Fatalf("first Upsert: %v", err)
		}
		u, err := repo.Upsert(ctx, newID(), "sub-dn-preserve", "dn@example.com", "ShouldBeIgnored")
		if err != nil {
			t.Fatalf("second Upsert: %v", err)
		}
		if u.DisplayName != "Original" {
			t.Errorf("display_name = %q, want Original", u.DisplayName)
		}
	})

	t.Run("second Upsert DOES set display_name when stored value is empty", func(t *testing.T) {
		id := newID()
		if _, err := repo.Upsert(ctx, id, "sub-dn-empty", "dnempty@example.com", ""); err != nil {
			t.Fatalf("first Upsert: %v", err)
		}
		u, err := repo.Upsert(ctx, newID(), "sub-dn-empty", "dnempty@example.com", "NewName")
		if err != nil {
			t.Fatalf("second Upsert: %v", err)
		}
		if u.DisplayName != "NewName" {
			t.Errorf("display_name = %q, want NewName", u.DisplayName)
		}
	})

	t.Run("updated_at advances after second Upsert", func(t *testing.T) {
		id := newID()
		u1, err := repo.Upsert(ctx, id, "sub-updated-at", "updtime@example.com", "N")
		if err != nil {
			t.Fatalf("first Upsert: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
		u2, err := repo.Upsert(ctx, newID(), "sub-updated-at", "updtime-new@example.com", "N")
		if err != nil {
			t.Fatalf("second Upsert: %v", err)
		}
		if !u2.UpdatedAt.After(u1.UpdatedAt) {
			t.Errorf("updated_at did not advance: first=%v second=%v", u1.UpdatedAt, u2.UpdatedAt)
		}
	})

	t.Run("GetByID returns user for known id", func(t *testing.T) {
		id := newID()
		if _, err := repo.Upsert(ctx, id, "sub-getbyid", "getbyid@example.com", "G"); err != nil {
			t.Fatalf("Upsert: %v", err)
		}
		u, err := repo.GetByID(ctx, id)
		if err != nil {
			t.Fatalf("GetByID: %v", err)
		}
		if u.ID != id {
			t.Errorf("id = %v, want %v", u.ID, id)
		}
	})

	t.Run("GetByID returns ErrNotFound for unknown id", func(t *testing.T) {
		_, err := repo.GetByID(ctx, newID())
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("err = %v, want domain.ErrNotFound", err)
		}
	})

	t.Run("GetByCognitoSub returns user for known sub", func(t *testing.T) {
		id := newID()
		if _, err := repo.Upsert(ctx, id, "sub-getbysub", "getbysub@example.com", "G"); err != nil {
			t.Fatalf("Upsert: %v", err)
		}
		u, err := repo.GetByCognitoSub(ctx, "sub-getbysub")
		if err != nil {
			t.Fatalf("GetByCognitoSub: %v", err)
		}
		if u.CognitoSub != "sub-getbysub" {
			t.Errorf("cognito_sub = %q, want sub-getbysub", u.CognitoSub)
		}
	})

	t.Run("GetByCognitoSub returns ErrNotFound for unknown sub", func(t *testing.T) {
		_, err := repo.GetByCognitoSub(ctx, "nonexistent-sub")
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("err = %v, want domain.ErrNotFound", err)
		}
	})
}
