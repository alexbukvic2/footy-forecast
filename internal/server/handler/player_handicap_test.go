package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/handler"
)

// ---------- fake service ----------

type fakePlayerHandicapService struct {
	getFn func(context.Context, uuid.UUID, domain.PlayerHandicapCategory) (*domain.PlayerHandicap, error)
}

func (f *fakePlayerHandicapService) Get(ctx context.Context, playerID uuid.UUID, category domain.PlayerHandicapCategory) (*domain.PlayerHandicap, error) {
	return f.getFn(ctx, playerID, category)
}

// ---------- helpers ----------

func getPlayerHandicap(t *testing.T, h http.HandlerFunc, playerID, category string) *httptest.ResponseRecorder {
	t.Helper()
	path := "/player-handicaps/" + playerID + "/" + category
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.SetPathValue("player_id", playerID)
	req.SetPathValue("category", category)
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

// ---------- tests ----------

func TestPlayerHandicap_Get_HappyPath(t *testing.T) {
	t.Parallel()

	playerID := uuid.New()
	svc := &fakePlayerHandicapService{
		getFn: func(_ context.Context, id uuid.UUID, cat domain.PlayerHandicapCategory) (*domain.PlayerHandicap, error) {
			require.Equal(t, playerID, id)
			require.Equal(t, domain.PlayerHandicapCategoryGroupTopScorer, cat)
			return &domain.PlayerHandicap{ID: uuid.New(), PlayerID: id, Category: cat, Points: 5}, nil
		},
	}
	h := handler.NewPlayerHandicap(silentLogger(), svc)

	rec := getPlayerHandicap(t, h.Get, playerID.String(), "group_top_scorer")

	require.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]any
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Equal(t, float64(5), resp["points"])
}

func TestPlayerHandicap_Get_InvalidPlayerID(t *testing.T) {
	t.Parallel()

	svc := &fakePlayerHandicapService{}
	h := handler.NewPlayerHandicap(silentLogger(), svc)

	rec := getPlayerHandicap(t, h.Get, "not-a-uuid", "group_top_scorer")
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPlayerHandicap_Get_UnknownCategory(t *testing.T) {
	t.Parallel()

	svc := &fakePlayerHandicapService{}
	h := handler.NewPlayerHandicap(silentLogger(), svc)

	rec := getPlayerHandicap(t, h.Get, uuid.New().String(), "invalid_category")
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPlayerHandicap_Get_ServiceErrNotFound(t *testing.T) {
	t.Parallel()

	svc := &fakePlayerHandicapService{
		getFn: func(context.Context, uuid.UUID, domain.PlayerHandicapCategory) (*domain.PlayerHandicap, error) {
			return nil, domain.ErrNotFound
		},
	}
	h := handler.NewPlayerHandicap(silentLogger(), svc)

	rec := getPlayerHandicap(t, h.Get, uuid.New().String(), "total_top_scorer")
	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestPlayerHandicap_Get_ServiceUnexpectedError(t *testing.T) {
	t.Parallel()

	svc := &fakePlayerHandicapService{
		getFn: func(context.Context, uuid.UUID, domain.PlayerHandicapCategory) (*domain.PlayerHandicap, error) {
			return nil, context.DeadlineExceeded
		},
	}
	h := handler.NewPlayerHandicap(silentLogger(), svc)

	rec := getPlayerHandicap(t, h.Get, uuid.New().String(), "group_top_scorer")
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
