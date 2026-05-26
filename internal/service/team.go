package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

// TeamRepo is the subset of the repository that TeamService needs.
type TeamRepo interface {
	ListWithHandicapsByTournament(ctx context.Context, tournamentID uuid.UUID) ([]domain.TeamWithHandicaps, error)
}

// TeamService orchestrates team use cases.
type TeamService struct {
	repo TeamRepo
}

// NewTeamService constructs a TeamService.
func NewTeamService(repo TeamRepo) *TeamService {
	return &TeamService{repo: repo}
}

// ListWithHandicaps returns all teams for the given tournament including their handicaps.
func (s *TeamService) ListWithHandicaps(ctx context.Context, tournamentID uuid.UUID) ([]domain.TeamWithHandicaps, error) {
	teams, err := s.repo.ListWithHandicapsByTournament(ctx, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("list teams with handicaps: %w", err)
	}
	return teams, nil
}

// compile-time interface check
var _ TeamRepo = (*repository.TeamRepository)(nil)
