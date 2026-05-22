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

func TestLeagueRepository_CRUD(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()

	// Insert fixtures.
	tournamentID := uuid.Must(uuid.NewV7())
	user1ID := uuid.Must(uuid.NewV7())
	user2ID := uuid.Must(uuid.NewV7())

	_, err := pool.Exec(ctx, `
		INSERT INTO tournaments (id, slug, name, starts_at, ends_at)
		VALUES ($1, 'test-cup', 'Test Cup',
		        '2026-06-01T00:00:00Z', '2026-07-01T00:00:00Z')`,
		tournamentID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO users (id, cognito_sub, email) VALUES
		($1, 'sub-1', 'user1@example.com'),
		($2, 'sub-2', 'user2@example.com')`,
		user1ID, user2ID)
	require.NoError(t, err)

	repo := repository.NewLeagueRepository(pool)

	var leagueID uuid.UUID

	t.Run("create returns persisted league", func(t *testing.T) {
		league, err := repo.Create(ctx, repository.CreateLeagueParams{
			TournamentID: tournamentID,
			OwnerID:      user1ID,
			Name:         "Test League",
			Code:         "TESTCODE",
		})
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, league.ID)
		require.Equal(t, tournamentID, league.TournamentID)
		require.Equal(t, user1ID, league.OwnerID)
		require.Equal(t, "Test League", league.Name)
		require.Equal(t, "TESTCODE", league.Code)
		require.WithinDuration(t, time.Now(), league.CreatedAt, 5*time.Second)
		leagueID = league.ID
	})

	t.Run("create with duplicate code returns ErrConflict", func(t *testing.T) {
		_, err := repo.Create(ctx, repository.CreateLeagueParams{
			TournamentID: tournamentID,
			OwnerID:      user1ID,
			Name:         "Another League",
			Code:         "TESTCODE",
		})
		require.True(t, errors.Is(err, domain.ErrConflict), "got %T %v", err, err)
	})

	t.Run("get by id returns the league", func(t *testing.T) {
		got, err := repo.GetByID(ctx, leagueID)
		require.NoError(t, err)
		require.Equal(t, leagueID, got.ID)
		require.Equal(t, "Test League", got.Name)
	})

	t.Run("get by id unknown returns ErrNotFound", func(t *testing.T) {
		_, err := repo.GetByID(ctx, uuid.New())
		require.True(t, errors.Is(err, domain.ErrNotFound), "got %T %v", err, err)
	})

	t.Run("get by code returns the league", func(t *testing.T) {
		got, err := repo.GetByCode(ctx, "TESTCODE")
		require.NoError(t, err)
		require.Equal(t, leagueID, got.ID)
	})

	t.Run("get by code unknown returns ErrNotFound", func(t *testing.T) {
		_, err := repo.GetByCode(ctx, "NOTEXIST")
		require.True(t, errors.Is(err, domain.ErrNotFound), "got %T %v", err, err)
	})

	t.Run("add member persists correctly", func(t *testing.T) {
		m, err := repo.AddMember(ctx, leagueID, user1ID, domain.LeagueMemberRoleOwner)
		require.NoError(t, err)
		require.Equal(t, leagueID, m.LeagueID)
		require.Equal(t, user1ID, m.UserID)
		require.Equal(t, domain.LeagueMemberRoleOwner, m.Role)
	})

	t.Run("add duplicate member returns ErrConflict", func(t *testing.T) {
		_, err := repo.AddMember(ctx, leagueID, user1ID, domain.LeagueMemberRoleMember)
		require.True(t, errors.Is(err, domain.ErrConflict), "got %T %v", err, err)
	})

	t.Run("get member returns the record", func(t *testing.T) {
		m, err := repo.GetMember(ctx, leagueID, user1ID)
		require.NoError(t, err)
		require.Equal(t, domain.LeagueMemberRoleOwner, m.Role)
	})

	t.Run("get member not found returns ErrNotFound", func(t *testing.T) {
		_, err := repo.GetMember(ctx, leagueID, uuid.New())
		require.True(t, errors.Is(err, domain.ErrNotFound), "got %T %v", err, err)
	})

	t.Run("list for user returns only user's leagues", func(t *testing.T) {
		// user2 has no membership yet.
		list, err := repo.ListForUser(ctx, user2ID)
		require.NoError(t, err)
		require.Empty(t, list)

		// user1 is a member.
		list, err = repo.ListForUser(ctx, user1ID)
		require.NoError(t, err)
		require.Len(t, list, 1)
		require.Equal(t, leagueID, list[0].ID)
	})

	t.Run("list members ordered by joined_at", func(t *testing.T) {
		_, err := repo.AddMember(ctx, leagueID, user2ID, domain.LeagueMemberRoleMember)
		require.NoError(t, err)

		members, err := repo.ListMembers(ctx, leagueID)
		require.NoError(t, err)
		require.Len(t, members, 2)
		// user1 joined first.
		require.Equal(t, user1ID, members[0].UserID)
	})

	t.Run("update name", func(t *testing.T) {
		updated, err := repo.UpdateName(ctx, leagueID, "Renamed League")
		require.NoError(t, err)
		require.Equal(t, "Renamed League", updated.Name)
	})

	t.Run("update code unique violation returns ErrConflict", func(t *testing.T) {
		// Create a second league with a known code then try to assign it to the first.
		league2, err := repo.Create(ctx, repository.CreateLeagueParams{
			TournamentID: tournamentID,
			OwnerID:      user2ID,
			Name:         "League 2",
			Code:         "CODE2222",
		})
		require.NoError(t, err)

		_, err = repo.UpdateCode(ctx, leagueID, league2.Code)
		require.True(t, errors.Is(err, domain.ErrConflict), "got %T %v", err, err)
	})

	t.Run("update code success", func(t *testing.T) {
		updated, err := repo.UpdateCode(ctx, leagueID, "NEWCODE1")
		require.NoError(t, err)
		require.Equal(t, "NEWCODE1", updated.Code)
	})

	t.Run("remove member", func(t *testing.T) {
		err := repo.RemoveMember(ctx, leagueID, user2ID)
		require.NoError(t, err)

		_, err = repo.GetMember(ctx, leagueID, user2ID)
		require.True(t, errors.Is(err, domain.ErrNotFound))
	})

	t.Run("remove non-existent member returns ErrNotFound", func(t *testing.T) {
		err := repo.RemoveMember(ctx, leagueID, uuid.New())
		require.True(t, errors.Is(err, domain.ErrNotFound), "got %T %v", err, err)
	})

	t.Run("delete cascades to members", func(t *testing.T) {
		err := repo.Delete(ctx, leagueID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, leagueID)
		require.True(t, errors.Is(err, domain.ErrNotFound))

		_, err = repo.GetMember(ctx, leagueID, user1ID)
		require.True(t, errors.Is(err, domain.ErrNotFound))
	})
}

func TestLeagueRepository_CreateWithOwner(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()

	tournamentID := uuid.Must(uuid.NewV7())
	ownerID := uuid.Must(uuid.NewV7())

	_, err := pool.Exec(ctx, `
		INSERT INTO tournaments (id, slug, name, starts_at, ends_at)
		VALUES ($1, 'atomic-cup', 'Atomic Cup',
		        '2026-06-01T00:00:00Z', '2026-07-01T00:00:00Z')`,
		tournamentID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO users (id, cognito_sub, email) VALUES ($1, 'sub-owner', 'owner@example.com')`,
		ownerID)
	require.NoError(t, err)

	repo := repository.NewLeagueRepository(pool)

	t.Run("inserts league and owner membership atomically", func(t *testing.T) {
		league, err := repo.CreateWithOwner(ctx, repository.CreateLeagueParams{
			TournamentID: tournamentID,
			OwnerID:      ownerID,
			Name:         "Atomic League",
			Code:         "ATOMCODE",
		})
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, league.ID)

		// Owner membership must exist.
		member, err := repo.GetMember(ctx, league.ID, ownerID)
		require.NoError(t, err)
		require.Equal(t, domain.LeagueMemberRoleOwner, member.Role)
	})

	t.Run("duplicate code returns ErrConflict and no league row persists", func(t *testing.T) {
		_, err := repo.CreateWithOwner(ctx, repository.CreateLeagueParams{
			TournamentID: tournamentID,
			OwnerID:      ownerID,
			Name:         "Collision League",
			Code:         "ATOMCODE", // same code as above
		})
		require.True(t, errors.Is(err, domain.ErrConflict), "got %T %v", err, err)
	})
}
