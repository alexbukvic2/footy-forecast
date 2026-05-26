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
	Search(ctx context.Context, tournamentID uuid.UUID, escapedQuery, rawQuery string, groupLetter *string, hasHandicap bool) ([]*domain.PlayerSearchResult, error)
}

// defaultPlayerHandicapPoints maps each player handicap category to the points
// returned when no handicap row exists for a player in that category.
// Business rule: group scorer slots are cheaper (20 pts) than the overall top-scorer (100 pts).
var defaultPlayerHandicapPoints = map[domain.PlayerHandicapCategory]int{
	domain.PlayerHandicapCategoryGroupTopScorer: 20,
	domain.PlayerHandicapCategoryTotalTopScorer: 100,
}

// FillDefaultPlayerHandicaps ensures every known category is present in p.Handicaps,
// inserting the per-category default for any that are missing.
// It is exported so other service packages can reuse the same rule.
func FillDefaultPlayerHandicaps(p *domain.PlayerSearchResult) {
	for cat, def := range defaultPlayerHandicapPoints {
		if _, ok := p.Handicaps[cat]; !ok {
			p.Handicaps[cat] = def
		}
	}
}

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

	if utf8.RuneCountInString(in.Query) > maxPlayerQueryRunes {
		return nil, fmt.Errorf("q must be at most %d characters: %w", maxPlayerQueryRunes, domain.ErrInvalid)
	}

	if _, err := s.tournaments.GetByID(ctx, in.TournamentID); err != nil {
		return nil, fmt.Errorf("get tournament: %w", err)
	}

	escaped := escapeWildcards(in.Query)
	players, err := s.players.Search(ctx, in.TournamentID, escaped, in.Query, in.GroupLetter, in.HasHandicap)
	if err != nil {
		return nil, fmt.Errorf("search players: %w", err)
	}

	for _, p := range players {
		FillDefaultPlayerHandicaps(p)
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
