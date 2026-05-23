package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

// PlayerHandicapRepo is the subset of the repository that PlayerHandicapService needs.
type PlayerHandicapRepo interface {
	Get(ctx context.Context, playerID uuid.UUID, category domain.PlayerHandicapCategory) (*domain.PlayerHandicap, error)
}

// PlayerHandicapService orchestrates player handicap use cases.
type PlayerHandicapService struct {
	repo PlayerHandicapRepo
}

// NewPlayerHandicapService constructs a PlayerHandicapService.
func NewPlayerHandicapService(repo PlayerHandicapRepo) *PlayerHandicapService {
	return &PlayerHandicapService{repo: repo}
}

// Get returns the handicap for the given player and category.
func (s *PlayerHandicapService) Get(
	ctx context.Context,
	playerID uuid.UUID,
	category domain.PlayerHandicapCategory,
) (*domain.PlayerHandicap, error) {
	h, err := s.repo.Get(ctx, playerID, category)
	if err != nil {
		return nil, fmt.Errorf("get player handicap: %w", err)
	}
	return h, nil
}

// compile-time interface check
var _ PlayerHandicapRepo = (*repository.PlayerHandicapRepository)(nil)
