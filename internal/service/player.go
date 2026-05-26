package service

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

// PlayerRepo is the subset of the repository that PlayerService needs.
type PlayerRepo interface {
	Search(ctx context.Context, tournamentID uuid.UUID, escapedQuery, rawQuery string, groupLetter *string) ([]*domain.PlayerSearchResult, error)
}

// defaultHandicapPoints is returned for any player+category pair that has no handicap row.
// Placing this here keeps the business rule colocated with the code that applies it.
const defaultHandicapPoints = 20

// PlayerService orchestrates player search use cases.
type PlayerService struct {
	players     PlayerRepo
	tournaments TournamentGetter
}

// NewPlayerService constructs a PlayerService.
func NewPlayerService(players PlayerRepo, tournaments TournamentGetter) *PlayerService {
	return &PlayerService{players: players, tournaments: tournaments}
}

const maxPlayerQueryRunes = 100

// Search validates input, confirms the tournament exists, then delegates to the repo.
// The query is trimmed of whitespace and SQL wildcards are escaped before passing to the repo.
// Every result carries handicap points for all player handicap categories; players without
// a handicap row for a given category receive defaultHandicapPoints.
func (s *PlayerService) Search(
	ctx context.Context,
	in domain.SearchPlayersInput,
) ([]*domain.PlayerSearchResult, error) {
	in.Query = strings.TrimSpace(in.Query)

	switch n := utf8.RuneCountInString(in.Query); {
	case n == 0:
		return nil, fmt.Errorf("q must not be empty: %w", domain.ErrInvalid)
	case n > maxPlayerQueryRunes:
		return nil, fmt.Errorf("q must be at most %d characters: %w", maxPlayerQueryRunes, domain.ErrInvalid)
	}

	if _, err := s.tournaments.GetByID(ctx, in.TournamentID); err != nil {
		return nil, fmt.Errorf("get tournament: %w", err)
	}

	escaped := escapeWildcards(in.Query)
	players, err := s.players.Search(ctx, in.TournamentID, escaped, in.Query, in.GroupLetter)
	if err != nil {
		return nil, fmt.Errorf("search players: %w", err)
	}

	for _, p := range players {
		for _, cat := range domain.AllPlayerHandicapCategories {
			if _, ok := p.Handicaps[cat]; !ok {
				p.Handicaps[cat] = defaultHandicapPoints
			}
		}
	}

	return players, nil
}

// escapeWildcards escapes SQL LIKE wildcard characters so user input is treated
// as a literal substring. The repo query uses ESCAPE '\'.
func escapeWildcards(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

// compile-time interface check
var _ PlayerRepo = (*repository.PlayerRepository)(nil)
