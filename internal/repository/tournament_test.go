//go:build integration

package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

func TestTournamentRepository_CRUD(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	repo := repository.NewTournamentRepository(pool)
	ctx := context.Background()

	params := repository.CreateTournamentParams{
		Slug:     "world-cup-2026",
		Name:     "FIFA World Cup 2026",
		StartsAt: time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
		EndsAt:   time.Date(2026, 7, 19, 23, 59, 59, 0, time.UTC),
	}

	t.Run(
		"create returns persisted tournament", func(t *testing.T) {
			created, err := repo.Create(ctx, params)
			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, created.ID)
			require.Equal(t, params.Slug, created.Slug)
			require.Equal(t, params.Name, created.Name)
			require.Equal(t, domain.TournamentStatusUpcoming, created.Status)
			require.WithinDuration(t, time.Now(), created.CreatedAt, 5*time.Second)
			require.WithinDuration(t, time.Now(), created.UpdatedAt, 5*time.Second)
		},
	)

	t.Run(
		"create with duplicate slug returns ErrConflict", func(t *testing.T) {
			_, err := repo.Create(ctx, params)
			require.Error(t, err)
			require.True(t, errors.Is(err, domain.ErrConflict), "got %T %v", err, err)
		},
	)

	t.Run(
		"get by id with unknown id returns ErrNotFound", func(t *testing.T) {
			_, err := repo.GetByID(ctx, uuid.New())
			require.True(t, errors.Is(err, domain.ErrNotFound), "got %T %v", err, err)
		},
	)

	t.Run(
		"list returns all tournaments", func(t *testing.T) {
			// Add a second tournament.
			_, err := repo.Create(
				ctx, repository.CreateTournamentParams{
					Slug:     "euro-2028",
					Name:     "UEFA Euro 2028",
					StartsAt: time.Date(2028, 6, 9, 0, 0, 0, 0, time.UTC),
					EndsAt:   time.Date(2028, 7, 9, 23, 59, 59, 0, time.UTC),
				},
			)
			require.NoError(t, err)

			list, err := repo.List(ctx)
			require.NoError(t, err)
			require.Len(t, list, 2)

			// Most recent first: euro-2028 starts in 2028, world-cup in 2026.
			require.Equal(t, "euro-2028", list[0].Slug)
			require.Equal(t, "world-cup-2026", list[1].Slug)
		},
	)
}

func TestTournamentRepository_UpdateAtChangesOnUpdate(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()

	// Skip via repository — we want to test the DB trigger directly,
	// and we don't have an Update method on the repo yet. Use raw SQL.
	id := uuid.Must(uuid.NewV7())
	_, err := pool.Exec(
		ctx, `
		INSERT INTO tournaments (id, slug, name, starts_at, ends_at)
		VALUES ($1, $2, $3, $4, $5)
	`, id, "test-cup", "Test Cup",
		time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 7, 19, 0, 0, 0, 0, time.UTC),
	)
	require.NoError(t, err)

	var firstUpdate, secondUpdate time.Time
	require.NoError(
		t, pool.QueryRow(
			ctx,
			"SELECT updated_at FROM tournaments WHERE id = $1", id,
		).Scan(&firstUpdate),
	)

	// Sleep just enough that NOW() advances.
	time.Sleep(10 * time.Millisecond)

	_, err = pool.Exec(ctx, "UPDATE tournaments SET name = 'Updated' WHERE id = $1", id)
	require.NoError(t, err)

	require.NoError(
		t, pool.QueryRow(
			ctx,
			"SELECT updated_at FROM tournaments WHERE id = $1", id,
		).Scan(&secondUpdate),
	)

	require.True(
		t, secondUpdate.After(firstUpdate),
		"updated_at should advance after UPDATE: first=%v second=%v", firstUpdate, secondUpdate,
	)
}
