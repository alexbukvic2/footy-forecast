// Package service contains business logic.
// Services orchestrate repositories and enforce domain rules.
// They do not depend on HTTP or specific database drivers.
package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

// TournamentRepo is the subset of the repository that TournamentService needs.
type TournamentRepo interface {
	Create(
		ctx context.Context,
		p repository.CreateTournamentParams,
	) (*domain.Tournament, error)
	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.Tournament, error)
	List(ctx context.Context) ([]*domain.Tournament, error)
}

// TournamentService orchestrates tournament use cases.
type TournamentService struct {
	repo TournamentRepo
}

// NewTournamentService constructs a TournamentService.
func NewTournamentService(repo TournamentRepo) *TournamentService {
	return &TournamentService{repo: repo}
}

// Create validates input and persists a new tournament.
//
// Returns:
//   - domain.ErrInvalid if input fails validation
//   - domain.ErrConflict if a tournament with the same slug exists
func (s *TournamentService) Create(
	ctx context.Context,
	in domain.CreateTournamentInput,
) (*domain.Tournament, error) {
	in.Slug = strings.TrimSpace(strings.ToLower(in.Slug))
	in.Name = strings.TrimSpace(in.Name)

	if err := validateSlug(in.Slug); err != nil {
		return nil, err
	}
	if err := validateName(in.Name); err != nil {
		return nil, err
	}
	if err := validateDates(in.StartsAt, in.EndsAt); err != nil {
		return nil, err
	}

	return s.repo.Create(
		ctx, repository.CreateTournamentParams{
			Slug:     in.Slug,
			Name:     in.Name,
			StartsAt: in.StartsAt.UTC(),
			EndsAt:   in.EndsAt.UTC(),
		},
	)
}

// GetByID returns the tournament with the given ID, or domain.ErrNotFound.
func (s *TournamentService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Tournament, error) {
	return s.repo.GetByID(ctx, id)
}

// List returns all tournaments.
func (s *TournamentService) List(ctx context.Context) ([]*domain.Tournament, error) {
	return s.repo.List(ctx)
}

// ---------- validation ----------

// slugPattern allows lowercase letters, digits, and single hyphens between segments.
// Min 3, max 64 chars. No leading/trailing hyphens, no consecutive hyphens.
var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

const (
	minSlugLen = 3
	maxSlugLen = 64
	minNameLen = 1
	maxNameLen = 200
)

func validateSlug(slug string) error {
	switch {
	case len(slug) < minSlugLen:
		return fmt.Errorf("slug must be at least %d chars: %w", minSlugLen, domain.ErrInvalid)
	case len(slug) > maxSlugLen:
		return fmt.Errorf("slug must be at most %d chars: %w", maxSlugLen, domain.ErrInvalid)
	case !slugPattern.MatchString(slug):
		return fmt.Errorf("slug must be lowercase alphanumeric with single hyphens: %w", domain.ErrInvalid)
	}
	return nil
}

func validateName(name string) error {
	switch {
	case len(name) < minNameLen:
		return fmt.Errorf("name must not be empty: %w", domain.ErrInvalid)
	case len(name) > maxNameLen:
		return fmt.Errorf("name must be at most %d chars: %w", maxNameLen, domain.ErrInvalid)
	}
	return nil
}

func validateDates(startsAt, endsAt time.Time) error {
	if startsAt.IsZero() {
		return fmt.Errorf("starts_at is required: %w", domain.ErrInvalid)
	}
	if endsAt.IsZero() {
		return fmt.Errorf("ends_at is required: %w", domain.ErrInvalid)
	}
	if !endsAt.After(startsAt) {
		return fmt.Errorf("ends_at must be after starts_at: %w", domain.ErrInvalid)
	}
	return nil
}

// ensure repository.TournamentRepository satisfies our interface at compile time
var _ TournamentRepo = (*repository.TournamentRepository)(nil)
