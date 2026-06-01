package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/handler"
)

// ---------- fake service ----------

type fakeOutcomesSvc struct {
	listFn func(context.Context, uuid.UUID) (*domain.TournamentOutcomes, error)
}

func (f *fakeOutcomesSvc) ListByTournament(ctx context.Context, id uuid.UUID) (*domain.TournamentOutcomes, error) {
	return f.listFn(ctx, id)
}

// ---------- tests ----------

func TestOutcomes_ListOutcomes(t *testing.T) {
	tournamentID := uuid.Must(uuid.NewV7())
	playerID := uuid.Must(uuid.NewV7())
	teamID := uuid.Must(uuid.NewV7())
	now := time.Now().Truncate(time.Second).UTC()

	t.Run("happy path returns 200 with outcomes", func(t *testing.T) {
		svc := &fakeOutcomesSvc{
			listFn: func(_ context.Context, _ uuid.UUID) (*domain.TournamentOutcomes, error) {
				return &domain.TournamentOutcomes{
					PlayerOutcomes: []*domain.PlayerOutcome{
						{
							ID:         uuid.Must(uuid.NewV7()),
							Category:   domain.PlayerHandicapCategoryTotalTopScorer,
							PlayerID:   playerID,
							PlayerName: "Messi",
							TeamID:     teamID,
							TeamName:   "Argentina",
							RecordedAt: now,
						},
					},
					TeamOutcomes: []*domain.TeamOutcome{
						{
							ID:         uuid.Must(uuid.NewV7()),
							Category:   domain.TeamHandicapCategoryWinner,
							TeamID:     teamID,
							TeamName:   "Argentina",
							RecordedAt: now,
						},
					},
				}, nil
			},
		}
		h := handler.NewOutcomes(silentLogger(), svc)
		rec := getWithPathValue(t, h.ListOutcomes, "/tournaments/"+tournamentID.String()+"/outcomes", "tournamentId", tournamentID.String())

		require.Equal(t, http.StatusOK, rec.Code)
		var body map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

		players := body["player_outcomes"].([]any)
		require.Len(t, players, 1)
		p := players[0].(map[string]any)
		require.Equal(t, "total_top_scorer", p["category"])
		require.Equal(t, "Messi", p["player_name"])
		require.Equal(t, "Argentina", p["team_name"])

		teams := body["team_outcomes"].([]any)
		require.Len(t, teams, 1)
		tm := teams[0].(map[string]any)
		require.Equal(t, "winner", tm["category"])
		require.Equal(t, "Argentina", tm["team_name"])
	})

	t.Run("empty outcomes returns 200 with empty arrays", func(t *testing.T) {
		svc := &fakeOutcomesSvc{
			listFn: func(_ context.Context, _ uuid.UUID) (*domain.TournamentOutcomes, error) {
				return &domain.TournamentOutcomes{
					PlayerOutcomes: []*domain.PlayerOutcome{},
					TeamOutcomes:   []*domain.TeamOutcome{},
				}, nil
			},
		}
		h := handler.NewOutcomes(silentLogger(), svc)
		rec := getWithPathValue(t, h.ListOutcomes, "/tournaments/"+tournamentID.String()+"/outcomes", "tournamentId", tournamentID.String())

		require.Equal(t, http.StatusOK, rec.Code)
		var body map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
		require.Empty(t, body["player_outcomes"].([]any))
		require.Empty(t, body["team_outcomes"].([]any))
	})

	t.Run("bad tournamentId returns 400", func(t *testing.T) {
		svc := &fakeOutcomesSvc{
			listFn: func(_ context.Context, _ uuid.UUID) (*domain.TournamentOutcomes, error) {
				return nil, nil
			},
		}
		h := handler.NewOutcomes(silentLogger(), svc)
		req := httptest.NewRequest(http.MethodGet, "/tournaments/not-a-uuid/outcomes", nil)
		req.SetPathValue("tournamentId", "not-a-uuid")
		rec := httptest.NewRecorder()
		h.ListOutcomes(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("tournament not found returns 404", func(t *testing.T) {
		svc := &fakeOutcomesSvc{
			listFn: func(_ context.Context, _ uuid.UUID) (*domain.TournamentOutcomes, error) {
				return nil, domain.ErrNotFound
			},
		}
		h := handler.NewOutcomes(silentLogger(), svc)
		rec := getWithPathValue(t, h.ListOutcomes, "/tournaments/"+tournamentID.String()+"/outcomes", "tournamentId", tournamentID.String())

		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("service error returns 500", func(t *testing.T) {
		svc := &fakeOutcomesSvc{
			listFn: func(_ context.Context, _ uuid.UUID) (*domain.TournamentOutcomes, error) {
				return nil, errors.New("db down")
			},
		}
		h := handler.NewOutcomes(silentLogger(), svc)
		rec := getWithPathValue(t, h.ListOutcomes, "/tournaments/"+tournamentID.String()+"/outcomes", "tournamentId", tournamentID.String())

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
