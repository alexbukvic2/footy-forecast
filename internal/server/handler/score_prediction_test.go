package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/handler"
)

// ---------- fake services ----------

type fakeScorePredictionSvc struct {
	upsertFn func(
		context.Context,
		domain.UpsertScorePredictionInput,
	) (*domain.ScorePrediction, error)
}

func (f *fakeScorePredictionSvc) UpsertScore(
	ctx context.Context,
	in domain.UpsertScorePredictionInput,
) (*domain.ScorePrediction, error) {
	return f.upsertFn(ctx, in)
}

type fakeFixtureSvc struct {
	listForUserFn func(
		context.Context,
		uuid.UUID,
		uuid.UUID,
	) ([]*domain.UserFixtureView, error)
	listForLeagueFn func(
		context.Context,
		uuid.UUID,
		uuid.UUID,
	) ([]*domain.LeagueFixtureView, error)
	listForLeaguePagedFn func(
		context.Context,
		uuid.UUID,
		uuid.UUID,
		int,
		int,
	) ([]*domain.LeagueFixtureView, error)
}

func (f *fakeFixtureSvc) ListForUser(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
) ([]*domain.UserFixtureView, error) {
	return f.listForUserFn(ctx, tournamentID, userID)
}

func (f *fakeFixtureSvc) ListForLeague(
	ctx context.Context,
	leagueID, userID uuid.UUID,
) ([]*domain.LeagueFixtureView, error) {
	return f.listForLeagueFn(ctx, leagueID, userID)
}

func (f *fakeFixtureSvc) ListForLeaguePaged(
	ctx context.Context,
	leagueID, userID uuid.UUID,
	n, skip int,
) ([]*domain.LeagueFixtureView, error) {
	return f.listForLeaguePagedFn(ctx, leagueID, userID, n, skip)
}

// ---------- helpers ----------

func putJSON(
	method, path string,
	body []byte,
) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return authedRequest(req)
}

func getAuthed(path string) *http.Request {
	return authedRequest(httptest.NewRequest(http.MethodGet, path, nil))
}

// ---------- ScorePrediction handler tests ----------

func TestScorePrediction_UpsertScore_HappyPath(t *testing.T) {
	t.Parallel()

	fixtureID := uuid.New()
	want := &domain.ScorePrediction{
		ID:        uuid.New(),
		FixtureID: fixtureID,
		GoalsHome: 2,
		GoalsAway: 1,
	}
	svc := &fakeScorePredictionSvc{
		upsertFn: func(
			_ context.Context,
			in domain.UpsertScorePredictionInput,
		) (*domain.ScorePrediction, error) {
			require.Equal(t, fixtureID, in.FixtureID)
			require.Equal(t, 2, in.GoalsHome)
			require.Equal(t, 1, in.GoalsAway)
			return want, nil
		},
	}
	h := handler.NewScorePrediction(silentLogger(), svc)

	body, _ := json.Marshal(map[string]int{"goals_home": 2, "goals_away": 1})
	req := putJSON(http.MethodPut, "/predictions/"+fixtureID.String(), body)
	req.SetPathValue("fixtureId", fixtureID.String())
	rec := httptest.NewRecorder()
	h.UpsertScore(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]any
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Equal(t, fixtureID.String(), resp["fixture_id"])
	require.Equal(t, float64(2), resp["goals_home"])
	require.Equal(t, float64(1), resp["goals_away"])
	require.Nil(t, resp["points"])
}

func TestScorePrediction_UpsertScore_NonUUIDFixtureID(t *testing.T) {
	t.Parallel()

	svc := &fakeScorePredictionSvc{}
	h := handler.NewScorePrediction(silentLogger(), svc)

	body, _ := json.Marshal(map[string]int{"goals_home": 1, "goals_away": 0})
	req := putJSON(http.MethodPut, "/predictions/not-a-uuid", body)
	req.SetPathValue("fixtureId", "not-a-uuid")
	rec := httptest.NewRecorder()
	h.UpsertScore(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestScorePrediction_UpsertScore_NegativeGoals(t *testing.T) {
	t.Parallel()

	fixtureID := uuid.New()
	svc := &fakeScorePredictionSvc{
		upsertFn: func(
			_ context.Context,
			_ domain.UpsertScorePredictionInput,
		) (*domain.ScorePrediction, error) {
			return nil, fmt.Errorf("goals must be non-negative: %w", domain.ErrInvalid)
		},
	}
	h := handler.NewScorePrediction(silentLogger(), svc)

	body, _ := json.Marshal(map[string]int{"goals_home": -1, "goals_away": 0})
	req := putJSON(http.MethodPut, "/predictions/"+fixtureID.String(), body)
	req.SetPathValue("fixtureId", fixtureID.String())
	rec := httptest.NewRecorder()
	h.UpsertScore(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestScorePrediction_UpsertScore_FixtureNotFound(t *testing.T) {
	t.Parallel()

	fixtureID := uuid.New()
	svc := &fakeScorePredictionSvc{
		upsertFn: func(
			context.Context,
			domain.UpsertScorePredictionInput,
		) (*domain.ScorePrediction, error) {
			return nil, domain.ErrNotFound
		},
	}
	h := handler.NewScorePrediction(silentLogger(), svc)

	body, _ := json.Marshal(map[string]int{"goals_home": 1, "goals_away": 0})
	req := putJSON(http.MethodPut, "/predictions/"+fixtureID.String(), body)
	req.SetPathValue("fixtureId", fixtureID.String())
	rec := httptest.NewRecorder()
	h.UpsertScore(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestScorePrediction_UpsertScore_Locked(t *testing.T) {
	t.Parallel()

	fixtureID := uuid.New()
	svc := &fakeScorePredictionSvc{
		upsertFn: func(
			context.Context,
			domain.UpsertScorePredictionInput,
		) (*domain.ScorePrediction, error) {
			return nil, fmt.Errorf("locked: %w", domain.ErrForbidden)
		},
	}
	h := handler.NewScorePrediction(silentLogger(), svc)

	body, _ := json.Marshal(map[string]int{"goals_home": 1, "goals_away": 0})
	req := putJSON(http.MethodPut, "/predictions/"+fixtureID.String(), body)
	req.SetPathValue("fixtureId", fixtureID.String())
	rec := httptest.NewRecorder()
	h.UpsertScore(rec, req)

	require.Equal(t, http.StatusForbidden, rec.Code)
}

func TestScorePrediction_UpsertScore_Unauthenticated(t *testing.T) {
	t.Parallel()

	svc := &fakeScorePredictionSvc{}
	h := handler.NewScorePrediction(silentLogger(), svc)

	body, _ := json.Marshal(map[string]int{"goals_home": 1, "goals_away": 0})
	// no user in context
	req := httptest.NewRequest(http.MethodPut, "/predictions/"+uuid.New().String(), bytes.NewReader(body))
	req.SetPathValue("fixtureId", uuid.New().String())
	rec := httptest.NewRecorder()
	h.UpsertScore(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

// ---------- Fixture handler tests ----------

func TestFixture_ListForUser_HappyPath(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	kickoff := time.Date(2026, 6, 1, 15, 0, 0, 0, time.UTC)
	want := []*domain.UserFixtureView{
		{
			Fixture:    domain.Fixture{ID: uuid.New(), HomeTeamName: "Argentina", AwayTeamName: "France", KickoffAt: kickoff},
			Prediction: nil,
		},
	}
	svc := &fakeFixtureSvc{
		listForUserFn: func(
			_ context.Context,
			tid, _ uuid.UUID,
		) ([]*domain.UserFixtureView, error) {
			require.Equal(t, tournamentID, tid)
			return want, nil
		},
	}
	h := handler.NewFixture(silentLogger(), svc)

	req := getAuthed("/tournaments/" + tournamentID.String() + "/fixtures")
	req.SetPathValue("tournamentId", tournamentID.String())
	rec := httptest.NewRecorder()
	h.ListForUser(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var resp []map[string]any
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Len(t, resp, 1)
	require.Equal(t, "Argentina", resp[0]["home_team_name"])
	require.Equal(t, "France", resp[0]["away_team_name"])
	require.Nil(t, resp[0]["prediction"])
}

func TestFixture_ListForUser_InvalidTournamentID(t *testing.T) {
	t.Parallel()

	svc := &fakeFixtureSvc{}
	h := handler.NewFixture(silentLogger(), svc)

	req := getAuthed("/tournaments/bad-id/fixtures")
	req.SetPathValue("tournamentId", "bad-id")
	rec := httptest.NewRecorder()
	h.ListForUser(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestFixture_ListForUser_Unauthenticated(t *testing.T) {
	t.Parallel()

	svc := &fakeFixtureSvc{}
	h := handler.NewFixture(silentLogger(), svc)

	req := httptest.NewRequest(http.MethodGet, "/tournaments/"+uuid.New().String()+"/fixtures", nil)
	req.SetPathValue("tournamentId", uuid.New().String())
	rec := httptest.NewRecorder()
	h.ListForUser(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestFixture_ListForLeague_HappyPath(t *testing.T) {
	t.Parallel()

	leagueID := uuid.New()
	kickoff := time.Date(2026, 6, 1, 15, 0, 0, 0, time.UTC)
	want := []*domain.LeagueFixtureView{
		{
			Fixture: domain.Fixture{ID: uuid.New(), HomeTeamName: "Brazil", AwayTeamName: "Germany", KickoffAt: kickoff, Status: domain.FixtureStatusFinished},
			Predictions: []domain.LeagueMemberPrediction{
				{UserID: uuid.New(), DisplayName: "alice", GoalsHome: ptr(2), GoalsAway: ptr(1)},
			},
		},
	}
	svc := &fakeFixtureSvc{
		listForLeagueFn: func(
			_ context.Context,
			lid, _ uuid.UUID,
		) ([]*domain.LeagueFixtureView, error) {
			require.Equal(t, leagueID, lid)
			return want, nil
		},
	}
	h := handler.NewFixture(silentLogger(), svc)

	req := getAuthed("/leagues/" + leagueID.String() + "/predictions")
	req.SetPathValue("leagueId", leagueID.String())
	rec := httptest.NewRecorder()
	h.ListForLeague(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var resp []map[string]any
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Len(t, resp, 1)
	require.Equal(t, "Brazil", resp[0]["home_team_name"])
	preds := resp[0]["predictions"].([]any)
	require.Len(t, preds, 1)
	require.Equal(t, "alice", preds[0].(map[string]any)["display_name"])
}

func TestFixture_ListForLeague_NotMember(t *testing.T) {
	t.Parallel()

	leagueID := uuid.New()
	svc := &fakeFixtureSvc{
		listForLeagueFn: func(
			context.Context,
			uuid.UUID,
			uuid.UUID,
		) ([]*domain.LeagueFixtureView, error) {
			return nil, fmt.Errorf("not a member: %w", domain.ErrForbidden)
		},
	}
	h := handler.NewFixture(silentLogger(), svc)

	req := getAuthed("/leagues/" + leagueID.String() + "/predictions")
	req.SetPathValue("leagueId", leagueID.String())
	rec := httptest.NewRecorder()
	h.ListForLeague(rec, req)

	require.Equal(t, http.StatusForbidden, rec.Code)
}

func TestFixture_ListForLeague_InvalidLeagueID(t *testing.T) {
	t.Parallel()

	svc := &fakeFixtureSvc{}
	h := handler.NewFixture(silentLogger(), svc)

	req := getAuthed("/leagues/bad-id/predictions")
	req.SetPathValue("leagueId", "bad-id")
	rec := httptest.NewRecorder()
	h.ListForLeague(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestFixture_ListForLeague_Unauthenticated(t *testing.T) {
	t.Parallel()

	svc := &fakeFixtureSvc{}
	h := handler.NewFixture(silentLogger(), svc)

	req := httptest.NewRequest(http.MethodGet, "/leagues/"+uuid.New().String()+"/predictions", nil)
	req.SetPathValue("leagueId", uuid.New().String())
	rec := httptest.NewRecorder()
	h.ListForLeague(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

// ---------- Fixture paged handler tests ----------

func TestFixture_ListForLeague_Paged_HappyPath(t *testing.T) {
	t.Parallel()

	leagueID := uuid.New()
	kickoff := time.Date(2026, 6, 1, 15, 0, 0, 0, time.UTC)
	want := []*domain.LeagueFixtureView{
		{
			Fixture:     domain.Fixture{ID: uuid.New(), HomeTeamName: "Spain", AwayTeamName: "Italy", KickoffAt: kickoff, Status: domain.FixtureStatusFinished},
			Predictions: []domain.LeagueMemberPrediction{},
		},
	}
	svc := &fakeFixtureSvc{
		listForLeaguePagedFn: func(_ context.Context, lid, _ uuid.UUID, n, skip int) ([]*domain.LeagueFixtureView, error) {
			require.Equal(t, leagueID, lid)
			require.Equal(t, 3, n)
			require.Equal(t, 5, skip)
			return want, nil
		},
	}
	h := handler.NewFixture(silentLogger(), svc)

	req := getAuthed("/leagues/" + leagueID.String() + "/predictions?n=3&skip=5")
	req.SetPathValue("leagueId", leagueID.String())
	rec := httptest.NewRecorder()
	h.ListForLeague(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var resp []map[string]any
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Len(t, resp, 1)
	require.Equal(t, "Spain", resp[0]["home_team_name"])
}

func TestFixture_ListForLeague_Paged_NOnlyDefaultsSkipZero(t *testing.T) {
	t.Parallel()

	leagueID := uuid.New()
	called := false
	svc := &fakeFixtureSvc{
		listForLeaguePagedFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID, n, skip int) ([]*domain.LeagueFixtureView, error) {
			called = true
			require.Equal(t, 2, n)
			require.Equal(t, 0, skip)
			return []*domain.LeagueFixtureView{}, nil
		},
	}
	h := handler.NewFixture(silentLogger(), svc)

	req := getAuthed("/leagues/" + leagueID.String() + "/predictions?n=2")
	req.SetPathValue("leagueId", leagueID.String())
	rec := httptest.NewRecorder()
	h.ListForLeague(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.True(t, called)
}

func TestFixture_ListForLeague_Paged_SkipWithoutN(t *testing.T) {
	t.Parallel()

	leagueID := uuid.New()
	svc := &fakeFixtureSvc{}
	h := handler.NewFixture(silentLogger(), svc)

	req := getAuthed("/leagues/" + leagueID.String() + "/predictions?skip=2")
	req.SetPathValue("leagueId", leagueID.String())
	rec := httptest.NewRecorder()
	h.ListForLeague(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestFixture_ListForLeague_Paged_InvalidN(t *testing.T) {
	t.Parallel()

	leagueID := uuid.New()
	svc := &fakeFixtureSvc{}
	h := handler.NewFixture(silentLogger(), svc)

	req := getAuthed("/leagues/" + leagueID.String() + "/predictions?n=0")
	req.SetPathValue("leagueId", leagueID.String())
	rec := httptest.NewRecorder()
	h.ListForLeague(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestFixture_ListForLeague_Paged_NegativeSkip(t *testing.T) {
	t.Parallel()

	leagueID := uuid.New()
	svc := &fakeFixtureSvc{}
	h := handler.NewFixture(silentLogger(), svc)

	req := getAuthed("/leagues/" + leagueID.String() + "/predictions?n=3&skip=-1")
	req.SetPathValue("leagueId", leagueID.String())
	rec := httptest.NewRecorder()
	h.ListForLeague(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestFixture_ListForLeague_Paged_NonIntegerSkip(t *testing.T) {
	t.Parallel()

	leagueID := uuid.New()
	svc := &fakeFixtureSvc{}
	h := handler.NewFixture(silentLogger(), svc)

	req := getAuthed("/leagues/" + leagueID.String() + "/predictions?n=3&skip=abc")
	req.SetPathValue("leagueId", leagueID.String())
	rec := httptest.NewRecorder()
	h.ListForLeague(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------- test helpers ----------

func ptr(v int) *int { return &v }
