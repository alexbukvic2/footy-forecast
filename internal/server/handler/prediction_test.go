package handler_test

import (
	"bytes"
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

type fakePredictionService struct {
	upsertScoreFn func(context.Context, domain.UpsertScorePredictionInput) (*domain.ScorePrediction, error)
}

func (f *fakePredictionService) UpsertScore(
	ctx context.Context,
	in domain.UpsertScorePredictionInput,
) (*domain.ScorePrediction, error) {
	return f.upsertScoreFn(ctx, in)
}

// ---------- helper ----------

func putJSONAuthedWithPathValue(
	t *testing.T,
	h http.HandlerFunc,
	path, key, value string,
	body any,
) *httptest.ResponseRecorder {
	t.Helper()
	raw, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPut, path, bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(key, value)
	req = authedRequest(req)
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

// ---------- tests ----------

func TestPrediction_UpsertScore(t *testing.T) {
	t.Parallel()

	fixtureID := uuid.New()

	t.Run("returns 200 with prediction (points null)", func(t *testing.T) {
		t.Parallel()
		created := &domain.ScorePrediction{
			ID:        uuid.New(),
			FixtureID: fixtureID,
			GoalsHome: 2,
			GoalsAway: 1,
			Points:    nil,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		svc := &fakePredictionService{
			upsertScoreFn: func(_ context.Context, _ domain.UpsertScorePredictionInput) (*domain.ScorePrediction, error) {
				return created, nil
			},
		}
		h := handler.NewPrediction(silentLogger(), svc)

		rec := putJSONAuthedWithPathValue(
			t, h.UpsertScore,
			"/predictions/"+fixtureID.String(), "fixtureId", fixtureID.String(),
			map[string]int{"goals_home": 2, "goals_away": 1},
		)

		require.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]any
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.Equal(t, float64(2), resp["goals_home"])
		require.Nil(t, resp["points"])
	})

	t.Run("returns 400 for non-UUID fixtureId", func(t *testing.T) {
		t.Parallel()
		h := handler.NewPrediction(silentLogger(), &fakePredictionService{})
		rec := putJSONAuthedWithPathValue(
			t, h.UpsertScore,
			"/predictions/bad", "fixtureId", "bad",
			map[string]int{"goals_home": 1, "goals_away": 0},
		)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 400 when service returns ErrInvalid", func(t *testing.T) {
		t.Parallel()
		svc := &fakePredictionService{
			upsertScoreFn: func(_ context.Context, _ domain.UpsertScorePredictionInput) (*domain.ScorePrediction, error) {
				return nil, domain.ErrInvalid
			},
		}
		h := handler.NewPrediction(silentLogger(), svc)
		rec := putJSONAuthedWithPathValue(
			t, h.UpsertScore,
			"/predictions/"+fixtureID.String(), "fixtureId", fixtureID.String(),
			map[string]int{"goals_home": -1, "goals_away": 0},
		)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 404 when service returns ErrNotFound", func(t *testing.T) {
		t.Parallel()
		svc := &fakePredictionService{
			upsertScoreFn: func(_ context.Context, _ domain.UpsertScorePredictionInput) (*domain.ScorePrediction, error) {
				return nil, domain.ErrNotFound
			},
		}
		h := handler.NewPrediction(silentLogger(), svc)
		rec := putJSONAuthedWithPathValue(
			t, h.UpsertScore,
			"/predictions/"+fixtureID.String(), "fixtureId", fixtureID.String(),
			map[string]int{"goals_home": 1, "goals_away": 0},
		)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("returns 403 when service returns ErrForbidden", func(t *testing.T) {
		t.Parallel()
		svc := &fakePredictionService{
			upsertScoreFn: func(_ context.Context, _ domain.UpsertScorePredictionInput) (*domain.ScorePrediction, error) {
				return nil, domain.ErrForbidden
			},
		}
		h := handler.NewPrediction(silentLogger(), svc)
		rec := putJSONAuthedWithPathValue(
			t, h.UpsertScore,
			"/predictions/"+fixtureID.String(), "fixtureId", fixtureID.String(),
			map[string]int{"goals_home": 1, "goals_away": 0},
		)
		require.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("returns 500 when user not in context", func(t *testing.T) {
		t.Parallel()
		h := handler.NewPrediction(silentLogger(), &fakePredictionService{})
		raw, _ := json.Marshal(map[string]int{"goals_home": 1, "goals_away": 0})
		req := httptest.NewRequest(http.MethodPut, "/predictions/"+fixtureID.String(), bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("fixtureId", fixtureID.String())
		// no user in context
		rec := httptest.NewRecorder()
		h.UpsertScore(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
