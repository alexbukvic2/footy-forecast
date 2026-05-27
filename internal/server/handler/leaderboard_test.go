package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/ctxutil"
	"github.com/alexbukvic2/footy-forecast/internal/server/handler"
)

// ---------- fake service ----------

type fakeLeaderboardService struct {
	leagueEntries     []*domain.LeaderboardEntry
	leagueErr         error
	tournamentEntries []*domain.LeaderboardEntry
	tournamentErr     error
}

func (f *fakeLeaderboardService) GetLeagueLeaderboard(_ context.Context, _, _ uuid.UUID) ([]*domain.LeaderboardEntry, error) {
	return f.leagueEntries, f.leagueErr
}
func (f *fakeLeaderboardService) GetTournamentLeaderboard(_ context.Context, _ uuid.UUID) ([]*domain.LeaderboardEntry, error) {
	return f.tournamentEntries, f.tournamentErr
}

// ---------- helpers ----------

func authedGetWithPathValue(t *testing.T, h http.HandlerFunc, path, key, value string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.SetPathValue(key, value)
	u := domain.User{ID: uuid.New(), Email: "caller@example.com", Status: domain.UserStatusActive}
	req = req.WithContext(ctxutil.WithUser(req.Context(), u))
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

func unauthGetWithPathValue(t *testing.T, h http.HandlerFunc, path, key, value string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.SetPathValue(key, value)
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

// ---------- GetForLeague ----------

func TestLeaderboard_GetForLeague(t *testing.T) {
	t.Parallel()

	t.Run("happy path returns 200 with leaderboard", func(t *testing.T) {
		t.Parallel()
		uid := uuid.New()
		entries := []*domain.LeaderboardEntry{
			{
				Position: 1, UserID: uid, DisplayName: "Alice",
				ScorePts: 5, GroupTopScorerPts: 3, TotalTopScorerPts: 2,
				GroupWinnerPts: 1, PlayoffPts: 2, SemifinalistPts: 1, WinnerPts: 4,
				TotalPoints: 18,
			},
		}
		svc := &fakeLeaderboardService{leagueEntries: entries}
		h := handler.NewLeaderboard(silentLogger(), svc)

		id := uuid.New()
		rec := authedGetWithPathValue(t, h.GetForLeague, "/leagues/"+id.String()+"/leaderboard", "id", id.String())

		require.Equal(t, http.StatusOK, rec.Code)
		var resp struct {
			Leaderboard []struct {
				Position    float64        `json:"position"`
				UserID      string         `json:"user_id"`
				DisplayName string         `json:"display_name"`
				TotalPoints float64        `json:"total_points"`
				Breakdown   map[string]any `json:"points_breakdown"`
			} `json:"leaderboard"`
		}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.Len(t, resp.Leaderboard, 1)
		e := resp.Leaderboard[0]
		require.Equal(t, float64(1), e.Position)
		require.Equal(t, uid.String(), e.UserID)
		require.Equal(t, "Alice", e.DisplayName)
		require.Equal(t, float64(18), e.TotalPoints)
		require.Equal(t, float64(5), e.Breakdown["score_pts"])
		require.Equal(t, float64(3), e.Breakdown["group_top_scorer_pts"])
		require.Equal(t, float64(2), e.Breakdown["total_top_scorer_pts"])
		require.Equal(t, float64(1), e.Breakdown["group_winner_pts"])
		require.Equal(t, float64(2), e.Breakdown["playoff_pts"])
		require.Equal(t, float64(1), e.Breakdown["semifinalist_pts"])
		require.Equal(t, float64(4), e.Breakdown["winner_pts"])
	})

	t.Run("non-member returns 403", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeaderboardService{leagueErr: fmt.Errorf("not a member: %w", domain.ErrForbidden)}
		h := handler.NewLeaderboard(silentLogger(), svc)

		id := uuid.New()
		rec := authedGetWithPathValue(t, h.GetForLeague, "/leagues/"+id.String()+"/leaderboard", "id", id.String())
		require.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("no auth returns 401", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeaderboardService{}
		h := handler.NewLeaderboard(silentLogger(), svc)

		id := uuid.New()
		rec := unauthGetWithPathValue(t, h.GetForLeague, "/leagues/"+id.String()+"/leaderboard", "id", id.String())
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid UUID returns 400", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeaderboardService{}
		h := handler.NewLeaderboard(silentLogger(), svc)

		rec := authedGetWithPathValue(t, h.GetForLeague, "/leagues/not-a-uuid/leaderboard", "id", "not-a-uuid")
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("league not found returns 404", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeaderboardService{leagueErr: fmt.Errorf("league not found: %w", domain.ErrNotFound)}
		h := handler.NewLeaderboard(silentLogger(), svc)

		id := uuid.New()
		rec := authedGetWithPathValue(t, h.GetForLeague, "/leagues/"+id.String()+"/leaderboard", "id", id.String())
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("service error returns 500", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeaderboardService{leagueErr: errors.New("unexpected db error")}
		h := handler.NewLeaderboard(silentLogger(), svc)

		id := uuid.New()
		rec := authedGetWithPathValue(t, h.GetForLeague, "/leagues/"+id.String()+"/leaderboard", "id", id.String())
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

// ---------- GetForTournament ----------

func TestLeaderboard_GetForTournament(t *testing.T) {
	t.Parallel()

	t.Run("happy path returns 200 with leaderboard", func(t *testing.T) {
		t.Parallel()
		uid := uuid.New()
		entries := []*domain.LeaderboardEntry{
			{Position: 1, UserID: uid, DisplayName: "Bob", TotalTopScorerPts: 4, WinnerPts: 6, TotalPoints: 10},
		}
		svc := &fakeLeaderboardService{tournamentEntries: entries}
		h := handler.NewLeaderboard(silentLogger(), svc)

		id := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/tournaments/"+id.String()+"/leaderboard", nil)
		req.SetPathValue("id", id.String())
		rec := httptest.NewRecorder()
		h.GetForTournament(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var resp struct {
			Leaderboard []map[string]any `json:"leaderboard"`
		}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.Len(t, resp.Leaderboard, 1)
		require.Equal(t, float64(1), resp.Leaderboard[0]["position"])
		require.Equal(t, "Bob", resp.Leaderboard[0]["display_name"])
	})

	t.Run("invalid UUID returns 400", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeaderboardService{}
		h := handler.NewLeaderboard(silentLogger(), svc)

		rec := getWithPathValue(t, h.GetForTournament, "/tournaments/bad/leaderboard", "id", "bad")
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("tournament not found returns 404", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeaderboardService{tournamentErr: fmt.Errorf("not found: %w", domain.ErrNotFound)}
		h := handler.NewLeaderboard(silentLogger(), svc)

		id := uuid.New()
		rec := getWithPathValue(t, h.GetForTournament, "/tournaments/"+id.String()+"/leaderboard", "id", id.String())
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}
