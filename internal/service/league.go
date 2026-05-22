package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

// LeagueRepo is the subset of the repository that LeagueService needs.
type LeagueRepo interface {
	// CreateWithOwner inserts the league and owner membership in one transaction.
	CreateWithOwner(
		ctx context.Context,
		p repository.CreateLeagueParams,
	) (*domain.League, error)
	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.League, error)
	GetByCode(
		ctx context.Context,
		code string,
	) (*domain.League, error)
	ListForUser(
		ctx context.Context,
		userID uuid.UUID,
	) ([]*domain.League, error)
	UpdateName(
		ctx context.Context,
		id uuid.UUID,
		name string,
	) (*domain.League, error)
	UpdateCode(
		ctx context.Context,
		id uuid.UUID,
		code string,
	) (*domain.League, error)
	Delete(
		ctx context.Context,
		id uuid.UUID,
	) error
	AddMember(
		ctx context.Context,
		leagueID, userID uuid.UUID,
		role domain.LeagueMemberRole,
	) (*domain.LeagueMember, error)
	RemoveMember(
		ctx context.Context,
		leagueID, userID uuid.UUID,
	) error
	GetMember(
		ctx context.Context,
		leagueID, userID uuid.UUID,
	) (*domain.LeagueMember, error)
	ListMembers(
		ctx context.Context,
		leagueID uuid.UUID,
	) ([]*domain.LeagueMember, error)
}

// TournamentGetter is the narrow interface used for tournament validation.
type TournamentGetter interface {
	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.Tournament, error)
}

// LeagueService orchestrates league use cases.
type LeagueService struct {
	leagues     LeagueRepo
	tournaments TournamentGetter
}

// NewLeagueService constructs a LeagueService.
func NewLeagueService(
	leagues LeagueRepo,
	tournaments TournamentGetter,
) *LeagueService {
	return &LeagueService{leagues: leagues, tournaments: tournaments}
}

const (
	minLeagueNameRunes = 1
	maxLeagueNameRunes = 100
)

const codeChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generateCode generates an 8-char uppercase alphanumeric invite code.
// Rejection sampling eliminates modulo bias (256 % 36 == 4 would otherwise
// make the first four characters ~14% more likely without it).
func generateCode() (string, error) {
	const n = 8
	// Largest multiple of len(codeChars) that fits in a byte; bytes >= threshold
	// are discarded so every accepted byte maps uniformly onto the alphabet.
	const threshold = 256 - (256 % len(codeChars)) // 252

	out := make([]byte, 0, n)
	buf := make([]byte, n+4) // slightly oversized; rejects are rare (~1.6%)
	for len(out) < n {
		if _, err := rand.Read(buf); err != nil {
			return "", fmt.Errorf("read random bytes: %w", err)
		}
		for _, v := range buf {
			if int(v) < threshold {
				out = append(out, codeChars[int(v)%len(codeChars)])
				if len(out) == n {
					break
				}
			}
		}
	}
	return string(out), nil
}

func validateLeagueName(name string) error {
	n := utf8.RuneCountInString(name)
	switch {
	case n < minLeagueNameRunes:
		return fmt.Errorf("name must not be empty: %w", domain.ErrInvalid)
	case n > maxLeagueNameRunes:
		return fmt.Errorf("name must be at most %d characters: %w", maxLeagueNameRunes, domain.ErrInvalid)
	}
	return nil
}

// CreateLeague validates input, verifies the tournament exists and is not concluded,
// generates an invite code, and creates the league with the caller as owner.
func (s *LeagueService) CreateLeague(
	ctx context.Context,
	userID uuid.UUID,
	in domain.CreateLeagueInput,
) (*domain.League, error) {
	in.Name = strings.TrimSpace(in.Name)

	if err := validateLeagueName(in.Name); err != nil {
		return nil, err
	}

	tournament, err := s.tournaments.GetByID(ctx, in.TournamentID)
	if err != nil {
		return nil, fmt.Errorf("get tournament: %w", err)
	}
	if tournament.Status == domain.TournamentStatusConcluded {
		return nil, fmt.Errorf("tournament %s is concluded: %w", in.TournamentID, domain.ErrInvalid)
	}

	code, err := generateCode()
	if err != nil {
		return nil, fmt.Errorf("generate invite code: %w", err)
	}

	p := repository.CreateLeagueParams{
		TournamentID: in.TournamentID,
		OwnerID:      userID,
		Name:         in.Name,
		Code:         code,
	}

	league, err := s.leagues.CreateWithOwner(ctx, p)
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			// Retry once on code collision (astronomically rare).
			p.Code, err = generateCode()
			if err != nil {
				return nil, fmt.Errorf("generate invite code (retry): %w", err)
			}
			league, err = s.leagues.CreateWithOwner(ctx, p)
			if err != nil {
				return nil, fmt.Errorf("create league (retry): %w", err)
			}
		} else {
			return nil, fmt.Errorf("create league: %w", err)
		}
	}

	return league, nil
}

// GetLeague returns a league with its member list, visible only to members.
// Returns ErrNotFound when the league does not exist or the requester is not a member
// (existence is not leaked to non-members).
func (s *LeagueService) GetLeague(
	ctx context.Context,
	leagueID, requesterID uuid.UUID,
) (*domain.League, []*domain.LeagueMember, error) {
	league, err := s.leagues.GetByID(ctx, leagueID)
	if err != nil {
		return nil, nil, err
	}

	if _, err := s.leagues.GetMember(ctx, leagueID, requesterID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			// Return ErrNotFound, not ErrForbidden — do not leak existence.
			return nil, nil, fmt.Errorf("league %s: %w", leagueID, domain.ErrNotFound)
		}
		return nil, nil, fmt.Errorf("check membership: %w", err)
	}

	members, err := s.leagues.ListMembers(ctx, leagueID)
	if err != nil {
		return nil, nil, fmt.Errorf("list members: %w", err)
	}

	return league, members, nil
}

// ListLeaguesForUser returns all leagues the user is a member of.
func (s *LeagueService) ListLeaguesForUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]*domain.League, error) {
	return s.leagues.ListForUser(ctx, userID)
}

// UpdateLeagueName renames a league. Only the owner may do this.
func (s *LeagueService) UpdateLeagueName(
	ctx context.Context,
	leagueID, requesterID uuid.UUID,
	name string,
) (*domain.League, error) {
	name = strings.TrimSpace(name)

	if err := validateLeagueName(name); err != nil {
		return nil, err
	}

	league, err := s.leagues.GetByID(ctx, leagueID)
	if err != nil {
		return nil, err
	}
	if league.OwnerID != requesterID {
		return nil, fmt.Errorf("update league name: %w", domain.ErrForbidden)
	}

	return s.leagues.UpdateName(ctx, leagueID, name)
}

// DeleteLeague deletes a league and all its members. Only the owner may do this.
func (s *LeagueService) DeleteLeague(
	ctx context.Context,
	leagueID, requesterID uuid.UUID,
) error {
	league, err := s.leagues.GetByID(ctx, leagueID)
	if err != nil {
		return err
	}
	if league.OwnerID != requesterID {
		return fmt.Errorf("delete league: %w", domain.ErrForbidden)
	}

	return s.leagues.Delete(ctx, leagueID)
}

// RegenerateCode generates and stores a new invite code. Only the owner may do this.
func (s *LeagueService) RegenerateCode(
	ctx context.Context,
	leagueID, requesterID uuid.UUID,
) (*domain.League, error) {
	league, err := s.leagues.GetByID(ctx, leagueID)
	if err != nil {
		return nil, err
	}
	if league.OwnerID != requesterID {
		return nil, fmt.Errorf("regenerate code: %w", domain.ErrForbidden)
	}

	code, err := generateCode()
	if err != nil {
		return nil, fmt.Errorf("generate invite code: %w", err)
	}

	updated, err := s.leagues.UpdateCode(ctx, leagueID, code)
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			code, err = generateCode()
			if err != nil {
				return nil, fmt.Errorf("generate invite code (retry): %w", err)
			}
			updated, err = s.leagues.UpdateCode(ctx, leagueID, code)
			if err != nil {
				return nil, fmt.Errorf("update code (retry): %w", err)
			}
		} else {
			return nil, fmt.Errorf("update code: %w", err)
		}
	}

	return updated, nil
}

// JoinLeague adds the user as a member of the league identified by code.
// Returns ErrNotFound for an unknown code, ErrConflict if already a member.
func (s *LeagueService) JoinLeague(
	ctx context.Context,
	code string,
	userID uuid.UUID,
) (*domain.League, error) {
	league, err := s.leagues.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	_, err = s.leagues.GetMember(ctx, league.ID, userID)
	if err == nil {
		return nil, fmt.Errorf("user already a member of league %s: %w", league.ID, domain.ErrConflict)
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("check membership: %w", err)
	}

	if _, err := s.leagues.AddMember(ctx, league.ID, userID, domain.LeagueMemberRoleMember); err != nil {
		return nil, fmt.Errorf("add member: %w", err)
	}

	return league, nil
}

// RemoveMember removes a user from a league.
// A member may remove themselves. An owner may remove any member.
// The owner cannot remove themselves (must delete the league instead).
func (s *LeagueService) RemoveMember(
	ctx context.Context,
	leagueID, targetUserID, requesterID uuid.UUID,
) error {
	league, err := s.leagues.GetByID(ctx, leagueID)
	if err != nil {
		return err
	}

	requesterIsOwner := league.OwnerID == requesterID
	selfAction := requesterID == targetUserID

	if !selfAction && !requesterIsOwner {
		return fmt.Errorf("remove member: %w", domain.ErrForbidden)
	}

	if targetUserID == league.OwnerID {
		return fmt.Errorf("owner cannot leave; delete the league instead: %w", domain.ErrInvalid)
	}

	return s.leagues.RemoveMember(ctx, leagueID, targetUserID)
}

// compile-time interface checks
var _ LeagueRepo = (*repository.LeagueRepository)(nil)
var _ TournamentGetter = (*repository.TournamentRepository)(nil)
