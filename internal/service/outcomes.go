package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

// OutcomesRepo is the data access interface for tournament outcomes.
type OutcomesRepo interface {
	ListPlayerOutcomes(ctx context.Context, tournamentID uuid.UUID) ([]*domain.PlayerOutcome, error)
	ListTeamOutcomes(ctx context.Context, tournamentID uuid.UUID) ([]*domain.TeamOutcome, error)
}

// OutcomesService handles tournament-outcomes use cases.
type OutcomesService struct {
	repo        OutcomesRepo
	tournaments TournamentGetter
}

// NewOutcomesService constructs an OutcomesService.
func NewOutcomesService(repo OutcomesRepo, tournaments TournamentGetter) *OutcomesService {
	return &OutcomesService{repo: repo, tournaments: tournaments}
}

// ListByTournament returns all recorded outcomes for the tournament.
// Returns domain.ErrNotFound if the tournament does not exist.
func (s *OutcomesService) ListByTournament(
	ctx context.Context,
	tournamentID uuid.UUID,
) (*domain.TournamentOutcomes, error) {
	if _, err := s.tournaments.GetByID(ctx, tournamentID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("tournament %s not found: %w", tournamentID, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get tournament: %w", err)
	}

	players, err := s.repo.ListPlayerOutcomes(ctx, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("list player outcomes: %w", err)
	}

	teams, err := s.repo.ListTeamOutcomes(ctx, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("list team outcomes: %w", err)
	}

	return &domain.TournamentOutcomes{
		PlayerOutcomes: players,
		TeamOutcomes:   teams,
	}, nil
}

// compile-time interface check
var _ OutcomesRepo = (*repository.OutcomesRepository)(nil)
