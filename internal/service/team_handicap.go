package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

// TeamHandicapRepo is the subset of the repository that TeamHandicapService needs.
type TeamHandicapRepo interface {
	Get(ctx context.Context, teamID uuid.UUID, category domain.TeamHandicapCategory) (*domain.TeamHandicap, error)
}

// TeamHandicapService orchestrates team handicap use cases.
type TeamHandicapService struct {
	repo TeamHandicapRepo
}

// NewTeamHandicapService constructs a TeamHandicapService.
func NewTeamHandicapService(repo TeamHandicapRepo) *TeamHandicapService {
	return &TeamHandicapService{repo: repo}
}

// Get returns the handicap for the given team and category.
func (s *TeamHandicapService) Get(
	ctx context.Context,
	teamID uuid.UUID,
	category domain.TeamHandicapCategory,
) (*domain.TeamHandicap, error) {
	h, err := s.repo.Get(ctx, teamID, category)
	if err != nil {
		return nil, fmt.Errorf("get team handicap: %w", err)
	}
	return h, nil
}

// compile-time interface check
var _ TeamHandicapRepo = (*repository.TeamHandicapRepository)(nil)
