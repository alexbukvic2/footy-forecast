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

// TeamRepository handles persistence of Team aggregates.
type TeamRepository struct {
	q *dbgen.Queries
}

// NewTeamRepository constructs a TeamRepository backed by pool.
func NewTeamRepository(pool *db.Pool) *TeamRepository {
	return &TeamRepository{q: dbgen.New(pool)}
}

// GetByID fetches a team by its UUID.
// Returns domain.ErrNotFound if no row exists.
func (r *TeamRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Team, error) {
	row, err := r.q.GetTeamByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("team %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get team: %w", err)
	}
	return &domain.Team{
		ID:           row.ID,
		Name:         row.Name,
		ShortName:    row.ShortName,
		Logo:         row.Logo,
		TournamentID: row.TournamentID,
		GroupLetter:  row.GroupLetter,
	}, nil
}

// ListGroupLettersByTournament returns sorted distinct group letters for teams in a tournament.
// Returns empty slice when no teams have a group_letter assigned.
func (r *TeamRepository) ListGroupLettersByTournament(ctx context.Context, tournamentID uuid.UUID) ([]string, error) {
	rows, err := r.q.ListGroupLettersByTournament(ctx, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("list group letters: %w", err)
	}
	out := make([]string, 0, len(rows))
	for _, ptr := range rows {
		if ptr != nil {
			out = append(out, *ptr)
		}
	}
	return out, nil
}

// ListWithHandicapsByTournament returns all teams for the tournament with their handicaps.
// Teams with no handicap rows will have an empty Handicaps slice.
func (r *TeamRepository) ListWithHandicapsByTournament(ctx context.Context, tournamentID uuid.UUID) ([]domain.TeamWithHandicaps, error) {
	rows, err := r.q.ListTeamsWithHandicapsByTournament(ctx, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("list teams with handicaps: %w", err)
	}

	// Collapse LEFT JOIN rows (one per team×handicap) into per-team entries.
	index := make(map[uuid.UUID]int, len(rows))
	out := make([]domain.TeamWithHandicaps, 0, len(rows))

	for _, row := range rows {
		idx, seen := index[row.ID]
		if !seen {
			out = append(out, domain.TeamWithHandicaps{
				Team: domain.Team{
					ID:           row.ID,
					Name:         row.Name,
					ShortName:    row.ShortName,
					Logo:         row.Logo,
					TournamentID: row.TournamentID,
					GroupLetter:  row.GroupLetter,
				},
				Handicaps: []domain.TeamHandicapItem{},
			})
			idx = len(out) - 1
			index[row.ID] = idx
		}
		if row.HandicapCategory != nil && row.HandicapPoints != nil {
			out[idx].Handicaps = append(out[idx].Handicaps, domain.TeamHandicapItem{
				Category: domain.TeamHandicapCategory(*row.HandicapCategory),
				Points:   int(*row.HandicapPoints),
			})
		}
	}
	return out, nil
}
