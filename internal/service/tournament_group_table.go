package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

// TournamentGroupTableRepo is the data access interface for group-table entries.
type TournamentGroupTableRepo interface {
	ListByTournament(
		ctx context.Context,
		tournamentID uuid.UUID,
	) ([]*domain.TournamentGroupEntry, error)
}

// TournamentGroupTableService handles group-stage standings use cases.
type TournamentGroupTableService struct {
	repo        TournamentGroupTableRepo
	tournaments TournamentGetter
}

// NewTournamentGroupTableService constructs a TournamentGroupTableService.
func NewTournamentGroupTableService(
	repo TournamentGroupTableRepo,
	tournaments TournamentGetter,
) *TournamentGroupTableService {
	return &TournamentGroupTableService{repo: repo, tournaments: tournaments}
}

// ListByTournament returns all group-table entries for the tournament.
// Returns domain.ErrNotFound if the tournament does not exist.
func (s *TournamentGroupTableService) ListByTournament(
	ctx context.Context,
	tournamentID uuid.UUID,
) ([]*domain.TournamentGroupEntry, error) {
	if _, err := s.tournaments.GetByID(ctx, tournamentID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("tournament %s not found: %w", tournamentID, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get tournament: %w", err)
	}

	entries, err := s.repo.ListByTournament(ctx, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("list group table: %w", err)
	}
	return entries, nil
}

// compile-time interface check
var _ TournamentGroupTableRepo = (*repository.TournamentGroupTableRepository)(nil)
