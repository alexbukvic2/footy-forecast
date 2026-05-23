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

type fakeTeamHandicapService struct {
	getFn func(context.Context, uuid.UUID, domain.TeamHandicapCategory) (*domain.TeamHandicap, error)
}

func (f *fakeTeamHandicapService) Get(ctx context.Context, teamID uuid.UUID, category domain.TeamHandicapCategory) (*domain.TeamHandicap, error) {
	return f.getFn(ctx, teamID, category)
}

// ---------- helpers ----------

func getTeamHandicap(t *testing.T, h http.HandlerFunc, teamID, category string) *httptest.ResponseRecorder {
	t.Helper()
	path := "/team-handicaps/" + teamID + "/" + category
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.SetPathValue("team_id", teamID)
	req.SetPathValue("category", category)
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

// ---------- tests ----------

func TestTeamHandicap_Get_HappyPath(t *testing.T) {
	t.Parallel()

	teamID := uuid.New()
	svc := &fakeTeamHandicapService{
		getFn: func(_ context.Context, id uuid.UUID, cat domain.TeamHandicapCategory) (*domain.TeamHandicap, error) {
			require.Equal(t, teamID, id)
			require.Equal(t, domain.TeamHandicapCategoryWinner, cat)
			return &domain.TeamHandicap{ID: uuid.New(), TeamID: id, Category: cat, Points: 10}, nil
		},
	}
	h := handler.NewTeamHandicap(silentLogger(), svc)

	rec := getTeamHandicap(t, h.Get, teamID.String(), "winner")

	require.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]any
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Equal(t, float64(10), resp["points"])
}

func TestTeamHandicap_Get_InvalidTeamID(t *testing.T) {
	t.Parallel()

	svc := &fakeTeamHandicapService{}
	h := handler.NewTeamHandicap(silentLogger(), svc)

	rec := getTeamHandicap(t, h.Get, "not-a-uuid", "winner")
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTeamHandicap_Get_UnknownCategory(t *testing.T) {
	t.Parallel()

	svc := &fakeTeamHandicapService{}
	h := handler.NewTeamHandicap(silentLogger(), svc)

	rec := getTeamHandicap(t, h.Get, uuid.New().String(), "invalid_category")
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTeamHandicap_Get_ServiceErrNotFound(t *testing.T) {
	t.Parallel()

	svc := &fakeTeamHandicapService{
		getFn: func(context.Context, uuid.UUID, domain.TeamHandicapCategory) (*domain.TeamHandicap, error) {
			return nil, domain.ErrNotFound
		},
	}
	h := handler.NewTeamHandicap(silentLogger(), svc)

	rec := getTeamHandicap(t, h.Get, uuid.New().String(), "group_winner")
	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTeamHandicap_Get_ServiceUnexpectedError(t *testing.T) {
	t.Parallel()

	svc := &fakeTeamHandicapService{
		getFn: func(context.Context, uuid.UUID, domain.TeamHandicapCategory) (*domain.TeamHandicap, error) {
			return nil, context.DeadlineExceeded
		},
	}
	h := handler.NewTeamHandicap(silentLogger(), svc)

	rec := getTeamHandicap(t, h.Get, uuid.New().String(), "semifinalist")
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
