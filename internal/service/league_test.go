package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
	"github.com/alexbukvic2/footy-forecast/internal/service"
)

// ---------- fakes ----------

type fakeLeagueRepo struct {
	createWithOwnerFn func(context.Context, repository.CreateLeagueParams) (*domain.League, error)
	getByIDFn         func(context.Context, uuid.UUID) (*domain.League, error)
	getByCodeFn       func(context.Context, string) (*domain.League, error)
	listForUserFn     func(context.Context, uuid.UUID) ([]*domain.League, error)
	updateNameFn      func(context.Context, uuid.UUID, string) (*domain.League, error)
	updateCodeFn      func(context.Context, uuid.UUID, string) (*domain.League, error)
	deleteFn          func(context.Context, uuid.UUID) error
	addMemberFn       func(context.Context, uuid.UUID, uuid.UUID, domain.LeagueMemberRole) (*domain.LeagueMember, error)
	removeMemberFn    func(context.Context, uuid.UUID, uuid.UUID) error
	getMemberFn       func(context.Context, uuid.UUID, uuid.UUID) (*domain.LeagueMember, error)
	listMembersFn     func(context.Context, uuid.UUID) ([]*domain.LeagueMember, error)
}

func (f *fakeLeagueRepo) CreateWithOwner(ctx context.Context, p repository.CreateLeagueParams) (*domain.League, error) {
	return f.createWithOwnerFn(ctx, p)
}
func (f *fakeLeagueRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.League, error) {
	return f.getByIDFn(ctx, id)
}
func (f *fakeLeagueRepo) GetByCode(ctx context.Context, code string) (*domain.League, error) {
	return f.getByCodeFn(ctx, code)
}
func (f *fakeLeagueRepo) ListForUser(ctx context.Context, userID uuid.UUID) ([]*domain.League, error) {
	return f.listForUserFn(ctx, userID)
}
func (f *fakeLeagueRepo) UpdateName(ctx context.Context, id uuid.UUID, name string) (*domain.League, error) {
	return f.updateNameFn(ctx, id, name)
}
func (f *fakeLeagueRepo) UpdateCode(ctx context.Context, id uuid.UUID, code string) (*domain.League, error) {
	return f.updateCodeFn(ctx, id, code)
}
func (f *fakeLeagueRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return f.deleteFn(ctx, id)
}
func (f *fakeLeagueRepo) AddMember(ctx context.Context, leagueID, userID uuid.UUID, role domain.LeagueMemberRole) (*domain.LeagueMember, error) {
	return f.addMemberFn(ctx, leagueID, userID, role)
}
func (f *fakeLeagueRepo) RemoveMember(ctx context.Context, leagueID, userID uuid.UUID) error {
	return f.removeMemberFn(ctx, leagueID, userID)
}
func (f *fakeLeagueRepo) GetMember(ctx context.Context, leagueID, userID uuid.UUID) (*domain.LeagueMember, error) {
	return f.getMemberFn(ctx, leagueID, userID)
}
func (f *fakeLeagueRepo) ListMembers(ctx context.Context, leagueID uuid.UUID) ([]*domain.LeagueMember, error) {
	return f.listMembersFn(ctx, leagueID)
}

type fakeTournamentGetter struct {
	getByIDFn func(context.Context, uuid.UUID) (*domain.Tournament, error)
}

func (f *fakeTournamentGetter) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tournament, error) {
	return f.getByIDFn(ctx, id)
}

// ---------- helpers ----------

func upcomingTournament(id uuid.UUID) *domain.Tournament {
	return &domain.Tournament{
		ID:       id,
		Status:   domain.TournamentStatusUpcoming,
		StartsAt: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		EndsAt:   time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
	}
}

// ---------- CreateLeague ----------

func TestLeagueService_CreateLeague_NameValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input string
	}{
		{"empty name", ""},
		{"whitespace only", "   "},
		{"name too long", string(make([]rune, 101))},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tournamentID := uuid.New()
			createCalled := false
			repo := &fakeLeagueRepo{
				createWithOwnerFn: func(context.Context, repository.CreateLeagueParams) (*domain.League, error) {
					createCalled = true
					return nil, nil
				},
			}
			tournaments := &fakeTournamentGetter{
				getByIDFn: func(context.Context, uuid.UUID) (*domain.Tournament, error) {
					return upcomingTournament(tournamentID), nil
				},
			}
			svc := service.NewLeagueService(repo, tournaments)

			_, err := svc.CreateLeague(context.Background(), uuid.New(), domain.CreateLeagueInput{
				TournamentID: tournamentID,
				Name:         tc.input,
			})

			require.ErrorIs(t, err, domain.ErrInvalid)
			require.False(t, createCalled, "repo.CreateWithOwner should not be called on invalid input")
		})
	}
}

func TestLeagueService_CreateLeague_TournamentNotFound(t *testing.T) {
	t.Parallel()

	repo := &fakeLeagueRepo{}
	tournaments := &fakeTournamentGetter{
		getByIDFn: func(context.Context, uuid.UUID) (*domain.Tournament, error) {
			return nil, domain.ErrNotFound
		},
	}
	svc := service.NewLeagueService(repo, tournaments)

	_, err := svc.CreateLeague(context.Background(), uuid.New(), domain.CreateLeagueInput{
		TournamentID: uuid.New(),
		Name:         "My League",
	})

	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestLeagueService_CreateLeague_TournamentConcluded(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	repo := &fakeLeagueRepo{}
	tournaments := &fakeTournamentGetter{
		getByIDFn: func(context.Context, uuid.UUID) (*domain.Tournament, error) {
			return &domain.Tournament{ID: tournamentID, Status: domain.TournamentStatusConcluded}, nil
		},
	}
	svc := service.NewLeagueService(repo, tournaments)

	_, err := svc.CreateLeague(context.Background(), uuid.New(), domain.CreateLeagueInput{
		TournamentID: tournamentID,
		Name:         "My League",
	})

	require.ErrorIs(t, err, domain.ErrInvalid)
}

func TestLeagueService_CreateLeague_Success(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	ownerID := uuid.New()
	leagueID := uuid.New()

	repo := &fakeLeagueRepo{
		createWithOwnerFn: func(_ context.Context, p repository.CreateLeagueParams) (*domain.League, error) {
			return &domain.League{ID: leagueID, TournamentID: p.TournamentID, OwnerID: p.OwnerID, Name: p.Name, Code: p.Code}, nil
		},
	}
	tournaments := &fakeTournamentGetter{
		getByIDFn: func(context.Context, uuid.UUID) (*domain.Tournament, error) {
			return upcomingTournament(tournamentID), nil
		},
	}
	svc := service.NewLeagueService(repo, tournaments)

	league, err := svc.CreateLeague(context.Background(), ownerID, domain.CreateLeagueInput{
		TournamentID: tournamentID,
		Name:         "My League",
	})

	require.NoError(t, err)
	require.Equal(t, leagueID, league.ID)
}

func TestLeagueService_CreateLeague_CodeCollisionRetry(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	leagueID := uuid.New()
	callCount := 0

	repo := &fakeLeagueRepo{
		createWithOwnerFn: func(_ context.Context, p repository.CreateLeagueParams) (*domain.League, error) {
			callCount++
			if callCount == 1 {
				return nil, domain.ErrConflict
			}
			return &domain.League{ID: leagueID, Code: p.Code}, nil
		},
	}
	tournaments := &fakeTournamentGetter{
		getByIDFn: func(context.Context, uuid.UUID) (*domain.Tournament, error) {
			return upcomingTournament(tournamentID), nil
		},
	}
	svc := service.NewLeagueService(repo, tournaments)

	league, err := svc.CreateLeague(context.Background(), uuid.New(), domain.CreateLeagueInput{
		TournamentID: tournamentID,
		Name:         "My League",
	})

	require.NoError(t, err)
	require.Equal(t, leagueID, league.ID)
	require.Equal(t, 2, callCount, "should retry once on code collision")
}

// ---------- GetLeague ----------

func TestLeagueService_GetLeague_NotMember(t *testing.T) {
	t.Parallel()

	leagueID := uuid.New()
	repo := &fakeLeagueRepo{
		getByIDFn: func(context.Context, uuid.UUID) (*domain.League, error) {
			return &domain.League{ID: leagueID}, nil
		},
		getMemberFn: func(context.Context, uuid.UUID, uuid.UUID) (*domain.LeagueMember, error) {
			return nil, domain.ErrNotFound
		},
	}
	svc := service.NewLeagueService(repo, &fakeTournamentGetter{})

	_, _, err := svc.GetLeague(context.Background(), leagueID, uuid.New())
	require.ErrorIs(t, err, domain.ErrNotFound)
}

// ---------- UpdateLeagueName ----------

func TestLeagueService_UpdateLeagueName_NotOwner(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	leagueID := uuid.New()
	repo := &fakeLeagueRepo{
		getByIDFn: func(context.Context, uuid.UUID) (*domain.League, error) {
			return &domain.League{ID: leagueID, OwnerID: ownerID}, nil
		},
	}
	svc := service.NewLeagueService(repo, &fakeTournamentGetter{})

	_, err := svc.UpdateLeagueName(context.Background(), leagueID, uuid.New(), "New Name")
	require.ErrorIs(t, err, domain.ErrForbidden)
}

// ---------- DeleteLeague ----------

func TestLeagueService_DeleteLeague_NotOwner(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	leagueID := uuid.New()
	repo := &fakeLeagueRepo{
		getByIDFn: func(context.Context, uuid.UUID) (*domain.League, error) {
			return &domain.League{ID: leagueID, OwnerID: ownerID}, nil
		},
	}
	svc := service.NewLeagueService(repo, &fakeTournamentGetter{})

	err := svc.DeleteLeague(context.Background(), leagueID, uuid.New())
	require.ErrorIs(t, err, domain.ErrForbidden)
}

// ---------- RegenerateCode ----------

func TestLeagueService_RegenerateCode_NotOwner(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	leagueID := uuid.New()
	repo := &fakeLeagueRepo{
		getByIDFn: func(context.Context, uuid.UUID) (*domain.League, error) {
			return &domain.League{ID: leagueID, OwnerID: ownerID}, nil
		},
	}
	svc := service.NewLeagueService(repo, &fakeTournamentGetter{})

	_, err := svc.RegenerateCode(context.Background(), leagueID, uuid.New())
	require.ErrorIs(t, err, domain.ErrForbidden)
}

// ---------- JoinLeague ----------

func TestLeagueService_JoinLeague_BadCode(t *testing.T) {
	t.Parallel()

	repo := &fakeLeagueRepo{
		getByCodeFn: func(context.Context, string) (*domain.League, error) {
			return nil, domain.ErrNotFound
		},
	}
	svc := service.NewLeagueService(repo, &fakeTournamentGetter{})

	_, err := svc.JoinLeague(context.Background(), "BADCODE1", uuid.New())
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestLeagueService_JoinLeague_AlreadyMember(t *testing.T) {
	t.Parallel()

	leagueID := uuid.New()
	userID := uuid.New()
	repo := &fakeLeagueRepo{
		getByCodeFn: func(context.Context, string) (*domain.League, error) {
			return &domain.League{ID: leagueID}, nil
		},
		getMemberFn: func(context.Context, uuid.UUID, uuid.UUID) (*domain.LeagueMember, error) {
			return &domain.LeagueMember{}, nil // already a member
		},
	}
	svc := service.NewLeagueService(repo, &fakeTournamentGetter{})

	_, err := svc.JoinLeague(context.Background(), "SOMECODE", userID)
	require.ErrorIs(t, err, domain.ErrConflict)
}

func TestLeagueService_JoinLeague_GetMemberErrorPropagated(t *testing.T) {
	t.Parallel()

	leagueID := uuid.New()
	repo := &fakeLeagueRepo{
		getByCodeFn: func(context.Context, string) (*domain.League, error) {
			return &domain.League{ID: leagueID}, nil
		},
		getMemberFn: func(context.Context, uuid.UUID, uuid.UUID) (*domain.LeagueMember, error) {
			return nil, context.DeadlineExceeded // non-ErrNotFound error
		},
	}
	svc := service.NewLeagueService(repo, &fakeTournamentGetter{})

	_, err := svc.JoinLeague(context.Background(), "SOMECODE", uuid.New())
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

// ---------- RemoveMember ----------

func TestLeagueService_RemoveMember_ForbiddenNonOwner(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	leagueID := uuid.New()
	targetID := uuid.New()
	requesterID := uuid.New() // neither owner nor target

	repo := &fakeLeagueRepo{
		getByIDFn: func(context.Context, uuid.UUID) (*domain.League, error) {
			return &domain.League{ID: leagueID, OwnerID: ownerID}, nil
		},
	}
	svc := service.NewLeagueService(repo, &fakeTournamentGetter{})

	err := svc.RemoveMember(context.Background(), leagueID, targetID, requesterID)
	require.ErrorIs(t, err, domain.ErrForbidden)
}

func TestLeagueService_RemoveMember_OwnerCannotLeave(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	leagueID := uuid.New()

	repo := &fakeLeagueRepo{
		getByIDFn: func(context.Context, uuid.UUID) (*domain.League, error) {
			return &domain.League{ID: leagueID, OwnerID: ownerID}, nil
		},
	}
	svc := service.NewLeagueService(repo, &fakeTournamentGetter{})

	// Owner tries to remove themselves.
	err := svc.RemoveMember(context.Background(), leagueID, ownerID, ownerID)
	require.ErrorIs(t, err, domain.ErrInvalid)
}

func TestLeagueService_RemoveMember_SelfLeave(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	leagueID := uuid.New()
	memberID := uuid.New()
	removed := false

	repo := &fakeLeagueRepo{
		getByIDFn: func(context.Context, uuid.UUID) (*domain.League, error) {
			return &domain.League{ID: leagueID, OwnerID: ownerID}, nil
		},
		removeMemberFn: func(_ context.Context, lgID, uID uuid.UUID) error {
			require.Equal(t, leagueID, lgID)
			require.Equal(t, memberID, uID)
			removed = true
			return nil
		},
	}
	svc := service.NewLeagueService(repo, &fakeTournamentGetter{})

	err := svc.RemoveMember(context.Background(), leagueID, memberID, memberID)
	require.NoError(t, err)
	require.True(t, removed)
}

func TestLeagueService_RemoveMember_OwnerRemovesMember(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	leagueID := uuid.New()
	memberID := uuid.New()
	removed := false

	repo := &fakeLeagueRepo{
		getByIDFn: func(context.Context, uuid.UUID) (*domain.League, error) {
			return &domain.League{ID: leagueID, OwnerID: ownerID}, nil
		},
		removeMemberFn: func(_ context.Context, _, _ uuid.UUID) error {
			removed = true
			return nil
		},
	}
	svc := service.NewLeagueService(repo, &fakeTournamentGetter{})

	err := svc.RemoveMember(context.Background(), leagueID, memberID, ownerID)
	require.NoError(t, err)
	require.True(t, removed)
}
