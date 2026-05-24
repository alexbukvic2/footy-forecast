package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository/dbgen"
)

// LeagueRepository handles persistence of League aggregates.
type LeagueRepository struct {
	pool *db.Pool
	q    *dbgen.Queries
}

// NewLeagueRepository constructs a LeagueRepository backed by pool.
func NewLeagueRepository(pool *db.Pool) *LeagueRepository {
	return &LeagueRepository{pool: pool, q: dbgen.New(pool)}
}

// CreateLeagueParams holds the fields needed to insert a league.
type CreateLeagueParams struct {
	TournamentID uuid.UUID
	OwnerID      uuid.UUID
	Name         string
	Code         string
}

// Create inserts a new league and returns the persisted record.
// Returns domain.ErrConflict if a league with the same code already exists.
// Use CreateWithOwner when you need league + owner membership inserted atomically.
func (r *LeagueRepository) Create(
	ctx context.Context,
	p CreateLeagueParams,
) (*domain.League, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate uuid: %w", err)
	}

	row, err := r.q.CreateLeague(ctx, dbgen.CreateLeagueParams{
		ID:           id,
		TournamentID: p.TournamentID,
		OwnerID:      p.OwnerID,
		Name:         p.Name,
		Code:         p.Code,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("league code %q: %w", p.Code, domain.ErrConflict)
		}
		return nil, fmt.Errorf("insert league: %w", err)
	}
	return leagueFromRow(row), nil
}

// CreateWithOwner inserts the league row and its owner membership atomically.
// Returns domain.ErrConflict if the invite code is already in use.
func (r *LeagueRepository) CreateWithOwner(
	ctx context.Context,
	p CreateLeagueParams,
) (*domain.League, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate uuid: %w", err)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := r.q.WithTx(tx)

	row, err := q.CreateLeague(ctx, dbgen.CreateLeagueParams{
		ID:           id,
		TournamentID: p.TournamentID,
		OwnerID:      p.OwnerID,
		Name:         p.Name,
		Code:         p.Code,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("league code %q: %w", p.Code, domain.ErrConflict)
		}
		return nil, fmt.Errorf("insert league: %w", err)
	}

	if _, err := q.AddLeagueMember(ctx, dbgen.AddLeagueMemberParams{
		LeagueID: row.ID,
		UserID:   p.OwnerID,
		Role:     dbgen.LeagueMemberRoleOwner,
	}); err != nil {
		return nil, fmt.Errorf("add owner as member: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return leagueFromRow(row), nil
}

// GetByID fetches a league by its UUID.
// Returns domain.ErrNotFound if no row exists.
func (r *LeagueRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.League, error) {
	row, err := r.q.GetLeagueByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("league %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get league: %w", err)
	}
	return leagueFromRow(row), nil
}

// GetByCode fetches a league by its invite code.
// Returns domain.ErrNotFound if no row exists.
func (r *LeagueRepository) GetByCode(
	ctx context.Context,
	code string,
) (*domain.League, error) {
	row, err := r.q.GetLeagueByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("league code %q: %w", code, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get league by code: %w", err)
	}
	return leagueFromRow(row), nil
}

// ListForUser returns all leagues the user is a member of, ordered by created_at DESC.
func (r *LeagueRepository) ListForUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]*domain.League, error) {
	rows, err := r.q.ListLeaguesForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list leagues for user: %w", err)
	}
	out := make([]*domain.League, 0, len(rows))
	for _, row := range rows {
		out = append(out, leagueFromRow(row))
	}
	return out, nil
}

// UpdateName updates the league name and returns the updated record.
// Returns domain.ErrNotFound if the league no longer exists.
func (r *LeagueRepository) UpdateName(
	ctx context.Context,
	id uuid.UUID,
	name string,
) (*domain.League, error) {
	row, err := r.q.UpdateLeagueName(ctx, dbgen.UpdateLeagueNameParams{ID: id, Name: name})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("league %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("update league name: %w", err)
	}
	return leagueFromRow(row), nil
}

// UpdateCode updates the league invite code and returns the updated record.
// Returns domain.ErrConflict if the new code is already in use.
// Returns domain.ErrNotFound if the league no longer exists.
func (r *LeagueRepository) UpdateCode(
	ctx context.Context,
	id uuid.UUID,
	code string,
) (*domain.League, error) {
	row, err := r.q.UpdateLeagueCode(ctx, dbgen.UpdateLeagueCodeParams{ID: id, Code: code})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("league %s: %w", id, domain.ErrNotFound)
		}
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("league code %q: %w", code, domain.ErrConflict)
		}
		return nil, fmt.Errorf("update league code: %w", err)
	}
	return leagueFromRow(row), nil
}

// Delete removes a league by ID. Cascades to league_members.
func (r *LeagueRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.q.DeleteLeague(ctx, id); err != nil {
		return fmt.Errorf("delete league: %w", err)
	}
	return nil
}

// AddMember inserts a league membership record.
// Returns domain.ErrConflict if the user is already a member.
func (r *LeagueRepository) AddMember(
	ctx context.Context,
	leagueID, userID uuid.UUID,
	role domain.LeagueMemberRole,
) (*domain.LeagueMember, error) {
	row, err := r.q.AddLeagueMember(ctx, dbgen.AddLeagueMemberParams{
		LeagueID: leagueID,
		UserID:   userID,
		Role:     dbgen.LeagueMemberRole(role),
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("user %s already member of league %s: %w", userID, leagueID, domain.ErrConflict)
		}
		return nil, fmt.Errorf("add league member: %w", err)
	}
	return leagueMemberFromRow(row), nil
}

// RemoveMember deletes a membership record.
// Returns domain.ErrNotFound if the membership does not exist.
func (r *LeagueRepository) RemoveMember(
	ctx context.Context,
	leagueID, userID uuid.UUID,
) error {
	tag, err := r.q.RemoveLeagueMember(ctx, dbgen.RemoveLeagueMemberParams{
		LeagueID: leagueID,
		UserID:   userID,
	})
	if err != nil {
		return fmt.Errorf("remove league member: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("member %s in league %s: %w", userID, leagueID, domain.ErrNotFound)
	}
	return nil
}

// GetMember fetches a single membership record.
// Returns domain.ErrNotFound if the user is not a member.
func (r *LeagueRepository) GetMember(
	ctx context.Context,
	leagueID, userID uuid.UUID,
) (*domain.LeagueMember, error) {
	row, err := r.q.GetLeagueMember(ctx, dbgen.GetLeagueMemberParams{
		LeagueID: leagueID,
		UserID:   userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("member %s in league %s: %w", userID, leagueID, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get league member: %w", err)
	}
	return leagueMemberFromRow(row), nil
}

// IsMember reports whether userID is a member of leagueID.
func (r *LeagueRepository) IsMember(
	ctx context.Context,
	leagueID, userID uuid.UUID,
) (bool, error) {
	_, err := r.GetMember(ctx, leagueID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ListMembers returns all members of a league ordered by joined_at ASC.
func (r *LeagueRepository) ListMembers(
	ctx context.Context,
	leagueID uuid.UUID,
) ([]*domain.LeagueMember, error) {
	rows, err := r.q.ListLeagueMembersForLeague(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("list league members: %w", err)
	}
	out := make([]*domain.LeagueMember, 0, len(rows))
	for _, row := range rows {
		out = append(out, leagueMemberFromRow(row))
	}
	return out, nil
}

func leagueFromRow(row dbgen.League) *domain.League {
	return &domain.League{
		ID:           row.ID,
		TournamentID: row.TournamentID,
		OwnerID:      row.OwnerID,
		Name:         row.Name,
		Code:         row.Code,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}

func leagueMemberFromRow(row dbgen.LeagueMember) *domain.LeagueMember {
	return &domain.LeagueMember{
		LeagueID: row.LeagueID,
		UserID:   row.UserID,
		Role:     domain.LeagueMemberRole(row.Role),
		JoinedAt: row.JoinedAt,
	}
}
