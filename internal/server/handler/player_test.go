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

type fakePlayerService struct {
	searchFn func(context.Context, domain.SearchPlayersInput) ([]*domain.PlayerSearchResult, error)
}

func (f *fakePlayerService) Search(ctx context.Context, in domain.SearchPlayersInput) ([]*domain.PlayerSearchResult, error) {
	return f.searchFn(ctx, in)
}

// ---------- helpers ----------

func getPlayerSearch(t *testing.T, h http.HandlerFunc, tournamentID, q string) *httptest.ResponseRecorder {
	t.Helper()
	path := "/tournaments/" + tournamentID + "/players/search"
	if q != "" {
		path += "?q=" + q
	}
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.SetPathValue("tournament_id", tournamentID)
	req = authedRequest(req)
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

// ---------- tests ----------

func TestPlayer_Search_HappyPath(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	players := []*domain.PlayerSearchResult{
		{ID: uuid.New(), Name: "Kylian Mbappé", TeamName: "France", TeamLogo: "<svg/>"},
	}
	var gotInput domain.SearchPlayersInput
	svc := &fakePlayerService{
		searchFn: func(_ context.Context, in domain.SearchPlayersInput) ([]*domain.PlayerSearchResult, error) {
			gotInput = in
			return players, nil
		},
	}
	h := handler.NewPlayer(silentLogger(), svc)

	rec := getPlayerSearch(t, h.Search, tournamentID.String(), "mbappe")

	require.Equal(t, http.StatusOK, rec.Code)
	var resp struct {
		Players []map[string]any `json:"players"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Len(t, resp.Players, 1)
	require.Equal(t, "Kylian Mbappé", resp.Players[0]["name"])
	require.Equal(t, players[0].ID.String(), resp.Players[0]["id"])
	require.Equal(t, "France", resp.Players[0]["team_name"])
	require.Equal(t, "<svg/>", resp.Players[0]["team_logo"])
	require.NotContains(t, resp.Players[0], "team_flag", "team_flag must not appear in response")
	// Verify the handler correctly forwarded both the parsed tournament ID and the query.
	require.Equal(t, tournamentID, gotInput.TournamentID)
	require.Equal(t, "mbappe", gotInput.Query)
}

func TestPlayer_Search_EmptyResults(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	svc := &fakePlayerService{
		searchFn: func(context.Context, domain.SearchPlayersInput) ([]*domain.PlayerSearchResult, error) {
			return []*domain.PlayerSearchResult{}, nil
		},
	}
	h := handler.NewPlayer(silentLogger(), svc)

	rec := getPlayerSearch(t, h.Search, tournamentID.String(), "nobody")

	require.Equal(t, http.StatusOK, rec.Code)
	var resp struct {
		Players []map[string]any `json:"players"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.NotNil(t, resp.Players)
	require.Empty(t, resp.Players)
}

func TestPlayer_Search_InvalidTournamentID(t *testing.T) {
	t.Parallel()

	svc := &fakePlayerService{}
	h := handler.NewPlayer(silentLogger(), svc)

	req := httptest.NewRequest(http.MethodGet, "/tournaments/not-a-uuid/players/search?q=messi", nil)
	req.SetPathValue("tournament_id", "not-a-uuid")
	req = authedRequest(req)
	rec := httptest.NewRecorder()
	h.Search(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPlayer_Search_ServiceReturnsErrInvalid(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	svc := &fakePlayerService{
		searchFn: func(context.Context, domain.SearchPlayersInput) ([]*domain.PlayerSearchResult, error) {
			return nil, domain.ErrInvalid
		},
	}
	h := handler.NewPlayer(silentLogger(), svc)

	rec := getPlayerSearch(t, h.Search, tournamentID.String(), "")

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPlayer_Search_ServiceReturnsErrNotFound(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	svc := &fakePlayerService{
		searchFn: func(context.Context, domain.SearchPlayersInput) ([]*domain.PlayerSearchResult, error) {
			return nil, domain.ErrNotFound
		},
	}
	h := handler.NewPlayer(silentLogger(), svc)

	rec := getPlayerSearch(t, h.Search, tournamentID.String(), "messi")

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestPlayer_Search_ServiceReturnsUnexpectedError(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	svc := &fakePlayerService{
		searchFn: func(context.Context, domain.SearchPlayersInput) ([]*domain.PlayerSearchResult, error) {
			return nil, context.DeadlineExceeded
		},
	}
	h := handler.NewPlayer(silentLogger(), svc)

	rec := getPlayerSearch(t, h.Search, tournamentID.String(), "messi")

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
