package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/handler"
)

// ---------- fake ----------

type fakeFixtureService struct {
	listForUserFn   func(context.Context, uuid.UUID, uuid.UUID) ([]*domain.UserFixtureView, error)
	listForLeagueFn func(context.Context, uuid.UUID, uuid.UUID) ([]*domain.LeagueFixtureView, error)
}

func (f *fakeFixtureService) ListForUser(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
) ([]*domain.UserFixtureView, error) {
	return f.listForUserFn(ctx, tournamentID, userID)
}

func (f *fakeFixtureService) ListForLeague(
	ctx context.Context,
	leagueID, userID uuid.UUID,
) ([]*domain.LeagueFixtureView, error) {
	return f.listForLeagueFn(ctx, leagueID, userID)
}

// ---------- tests ----------

func TestFixture_ListForUser(t *testing.T) {
	t.Parallel()

	fixture := domain.Fixture{
		ID:           uuid.New(),
		TournamentID: uuid.New(),
		HomeTeamID:   uuid.New(),
		AwayTeamID:   uuid.New(),
		KickoffAt:    time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC),
		Status:       domain.FixtureStatusUpcoming,
	}

	t.Run("returns 200 with fixtures and nil prediction", func(t *testing.T) {
		t.Parallel()
		svc := &fakeFixtureService{
			listForUserFn: func(_ context.Context, _, _ uuid.UUID) ([]*domain.UserFixtureView, error) {
				return []*domain.UserFixtureView{{Fixture: fixture}}, nil
			},
		}
		h := handler.NewFixture(silentLogger(), svc)

		tournamentID := fixture.TournamentID.String()
		req := httptest.NewRequest(http.MethodGet, "/tournaments/"+tournamentID+"/fixtures", nil)
		req.SetPathValue("tournamentId", tournamentID)
		req = authedRequest(req)
		rec := httptest.NewRecorder()
		h.ListForUser(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var resp struct {
			Fixtures []map[string]any `json:"fixtures"`
		}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.Len(t, resp.Fixtures, 1)
		require.Nil(t, resp.Fixtures[0]["prediction"])
	})

	t.Run("returns 400 on invalid tournamentId", func(t *testing.T) {
		t.Parallel()
		h := handler.NewFixture(silentLogger(), &fakeFixtureService{})
		req := httptest.NewRequest(http.MethodGet, "/tournaments/bad/fixtures", nil)
		req.SetPathValue("tournamentId", "bad")
		req = authedRequest(req)
		rec := httptest.NewRecorder()
		h.ListForUser(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 500 when user not in context", func(t *testing.T) {
		t.Parallel()
		h := handler.NewFixture(silentLogger(), &fakeFixtureService{})
		req := httptest.NewRequest(http.MethodGet, "/tournaments/"+fixture.TournamentID.String()+"/fixtures", nil)
		req.SetPathValue("tournamentId", fixture.TournamentID.String())
		rec := httptest.NewRecorder()
		h.ListForUser(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestFixture_ListForLeague(t *testing.T) {
	t.Parallel()

	leagueID := uuid.New()

	t.Run("returns 200 with predictions array", func(t *testing.T) {
		t.Parallel()
		memberID := uuid.New()
		ghome := 2
		gaway := 1
		view := &domain.LeagueFixtureView{
			Fixture: domain.Fixture{
				ID:        uuid.New(),
				KickoffAt: time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC),
				Status:    domain.FixtureStatusFinished,
			},
			Predictions: []domain.LeagueMemberPrediction{
				{UserID: memberID, DisplayName: "Alice", GoalsHome: &ghome, GoalsAway: &gaway},
			},
		}
		svc := &fakeFixtureService{
			listForLeagueFn: func(_ context.Context, _, _ uuid.UUID) ([]*domain.LeagueFixtureView, error) {
				return []*domain.LeagueFixtureView{view}, nil
			},
		}
		h := handler.NewFixture(silentLogger(), svc)

		req := httptest.NewRequest(http.MethodGet, "/leagues/"+leagueID.String()+"/predictions", nil)
		req.SetPathValue("leagueId", leagueID.String())
		req = authedRequest(req)
		rec := httptest.NewRecorder()
		h.ListForLeague(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var resp struct {
			Fixtures []map[string]any `json:"fixtures"`
		}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.Len(t, resp.Fixtures, 1)
		preds := resp.Fixtures[0]["predictions"].([]any)
		require.Len(t, preds, 1)
		p := preds[0].(map[string]any)
		require.Equal(t, "Alice", p["display_name"])
	})

	t.Run("returns 403 when service returns ErrForbidden", func(t *testing.T) {
		t.Parallel()
		svc := &fakeFixtureService{
			listForLeagueFn: func(_ context.Context, _, _ uuid.UUID) ([]*domain.LeagueFixtureView, error) {
				return nil, domain.ErrForbidden
			},
		}
		h := handler.NewFixture(silentLogger(), svc)
		req := httptest.NewRequest(http.MethodGet, "/leagues/"+leagueID.String()+"/predictions", nil)
		req.SetPathValue("leagueId", leagueID.String())
		req = authedRequest(req)
		rec := httptest.NewRecorder()
		h.ListForLeague(rec, req)
		require.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("returns 400 on invalid leagueId", func(t *testing.T) {
		t.Parallel()
		h := handler.NewFixture(silentLogger(), &fakeFixtureService{})
		req := httptest.NewRequest(http.MethodGet, "/leagues/bad/predictions", nil)
		req.SetPathValue("leagueId", "bad")
		req = authedRequest(req)
		rec := httptest.NewRecorder()
		h.ListForLeague(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
