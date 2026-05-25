// Package repository contains database access code.
// Each repository wraps the generated sqlc code, translating between
// persistence rows (dbgen.*) and domain types (domain.*).
package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository/dbgen"
)

// TournamentRepository handles persistence of Tournament aggregates.
type TournamentRepository struct {
	q *dbgen.Queries
}

// NewTournamentRepository constructs a TournamentRepository backed by pool.
func NewTournamentRepository(pool *db.Pool) *TournamentRepository {
	return &TournamentRepository{q: dbgen.New(pool)}
}

// CreateTournamentParams holds the fields needed to insert a tournament.
type CreateTournamentParams struct {
	Slug     string
	Name     string
	StartsAt time.Time
	EndsAt   time.Time
}

// Create inserts a new tournament and returns the persisted record.
// Returns domain.ErrConflict if a tournament with the same slug already exists.
func (r *TournamentRepository) Create(
	ctx context.Context,
	p CreateTournamentParams,
) (*domain.Tournament, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate uuid: %w", err)
	}

	row, err := r.q.CreateTournament(
		ctx, dbgen.CreateTournamentParams{
			ID:       id,
			Slug:     p.Slug,
			Name:     p.Name,
			StartsAt: p.StartsAt,
			EndsAt:   p.EndsAt,
		},
	)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("tournament with slug %q: %w", p.Slug, domain.ErrConflict)
		}
		return nil, fmt.Errorf("insert tournament: %w", err)
	}
	return tournamentFromRow(row), nil
}

// GetByID fetches a tournament by its UUID.
// Returns domain.ErrNotFound if no row exists.
func (r *TournamentRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Tournament, error) {
	row, err := r.q.GetTournamentByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("tournament %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get tournament: %w", err)
	}
	return tournamentFromRow(row), nil
}

// List returns all tournaments ordered by start date (most recent first).
func (r *TournamentRepository) List(ctx context.Context) ([]*domain.Tournament, error) {
	rows, err := r.q.ListTournaments(ctx)
	if err != nil {
		return nil, fmt.Errorf("list tournaments: %w", err)
	}
	out := make([]*domain.Tournament, 0, len(rows))
	for _, row := range rows {
		out = append(out, tournamentFromRow(row))
	}
	return out, nil
}

// tournamentFromRow maps a persistence row to the domain type.
// This is the boundary between DB shape and domain shape.
func tournamentFromRow(row dbgen.Tournament) *domain.Tournament {
	return &domain.Tournament{
		ID:        row.ID,
		Slug:      row.Slug,
		Name:      row.Name,
		Status:    domain.TournamentStatus(row.Status),
		StartsAt:  row.StartsAt,
		EndsAt:    row.EndsAt,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

// isUniqueViolation reports whether err is a Postgres unique constraint violation.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == pgerrcode.UniqueViolation
	}
	return false
}
