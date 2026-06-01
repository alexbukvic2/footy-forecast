package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository/dbgen"
)

// OutcomesRepository handles persistence for player and team tournament outcomes.
type OutcomesRepository struct {
	q *dbgen.Queries
}

// NewOutcomesRepository constructs an OutcomesRepository backed by pool.
func NewOutcomesRepository(pool *db.Pool) *OutcomesRepository {
	return &OutcomesRepository{q: dbgen.New(pool)}
}

// ListPlayerOutcomes returns all recorded player outcomes for a tournament.
func (r *OutcomesRepository) ListPlayerOutcomes(
	ctx context.Context,
	tournamentID uuid.UUID,
) ([]*domain.PlayerOutcome, error) {
	rows, err := r.q.ListPlayerOutcomesByTournament(ctx, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("list player outcomes: %w", err)
	}
	out := make([]*domain.PlayerOutcome, 0, len(rows))
	for _, row := range rows {
		out = append(out, &domain.PlayerOutcome{
			ID:           row.ID,
			TournamentID: row.TournamentID,
			Category:     domain.PlayerHandicapCategory(row.Category),
			PlayerID:     row.PlayerID,
			PlayerName:   row.PlayerName,
			TeamID:       row.TeamID,
			TeamName:     row.TeamName,
			RecordedAt:   row.RecordedAt,
		})
	}
	return out, nil
}

// ListTeamOutcomes returns all recorded team outcomes for a tournament.
func (r *OutcomesRepository) ListTeamOutcomes(
	ctx context.Context,
	tournamentID uuid.UUID,
) ([]*domain.TeamOutcome, error) {
	rows, err := r.q.ListTeamOutcomesByTournament(ctx, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("list team outcomes: %w", err)
	}
	out := make([]*domain.TeamOutcome, 0, len(rows))
	for _, row := range rows {
		out = append(out, &domain.TeamOutcome{
			ID:           row.ID,
			TournamentID: row.TournamentID,
			Category:     domain.TeamHandicapCategory(row.Category),
			TeamID:       row.TeamID,
			TeamName:     row.TeamName,
			RecordedAt:   row.RecordedAt,
		})
	}
	return out, nil
}
