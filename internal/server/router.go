// Package server wires together the application's HTTP layer:
// routes, middleware, and handlers.
package server

import (
	"log/slog"
	"net/http"

	"github.com/alexbukvic2/footy-forecast/internal/config"
	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
	"github.com/alexbukvic2/footy-forecast/internal/server/cognito"
	"github.com/alexbukvic2/footy-forecast/internal/server/handler"
	"github.com/alexbukvic2/footy-forecast/internal/server/middleware"
	"github.com/alexbukvic2/footy-forecast/internal/service"
)

// NewRouter wires up all HTTP routes and returns the root handler.
func NewRouter(
	logger *slog.Logger,
	pool *db.Pool,
	cfg *config.Config,
) http.Handler {
	// Wire dependencies bottom-up.
	tournamentRepo := repository.NewTournamentRepository(pool)
	tournamentSvc := service.NewTournamentService(tournamentRepo)

	userRepo := repository.NewUserRepository(pool)
	userSvc := service.NewUserService(userRepo)

	leagueRepo := repository.NewLeagueRepository(pool)
	leagueSvc := service.NewLeagueService(leagueRepo, tournamentRepo)

	playerHandicapRepo := repository.NewPlayerHandicapRepository(pool)
	playerHandicapSvc := service.NewPlayerHandicapService(playerHandicapRepo)

	playerRepo := repository.NewPlayerRepository(pool)
	playerSvc := service.NewPlayerService(playerRepo, tournamentRepo)

	teamRepo := repository.NewTeamRepository(pool)
	teamSvc := service.NewTeamService(teamRepo)

	fixtureRepo := repository.NewFixtureRepository(pool)
	predictionRepo := repository.NewPredictionRepository(pool)
	predictionSvc := service.NewPredictionService(predictionRepo, fixtureRepo, service.RealClock{})
	fixtureSvc := service.NewFixtureService(fixtureRepo, leagueRepo)
	playerPredictionRepo := repository.NewPlayerPredictionRepository(pool)
	teamPredictionRepo := repository.NewTeamPredictionRepository(pool)
	tournamentPredictionSvc := service.NewTournamentPredictionService(
		playerPredictionRepo,
		teamPredictionRepo,
		playerRepo,
		teamRepo,
		teamRepo,
		fixtureRepo,
		leagueRepo,
		service.RealClock{},
	)

	validator := cognito.NewValidator(
		cfg.CognitoRegion,
		cfg.CognitoUserPoolID,
		cfg.CognitoAllowedClientIDs,
	)
	authMW := middleware.Auth(validator, userSvc)

	// Handlers.
	healthH := handler.NewHealth(logger, pool)
	tournamentH := handler.NewTournament(logger, tournamentSvc)
	userH := handler.NewUser(logger)
	leagueH := handler.NewLeague(logger, leagueSvc)
	playerH := handler.NewPlayer(logger, playerSvc)
	playerHandicapH := handler.NewPlayerHandicap(logger, playerHandicapSvc)
	teamH := handler.NewTeam(logger, teamSvc)
	tournamentPredictionH := handler.NewTournamentPrediction(logger, tournamentPredictionSvc)
	scorePredictionH := handler.NewScorePrediction(logger, predictionSvc)
	fixtureH := handler.NewFixture(logger, fixtureSvc)
	specH := handler.NewSpec()

	mux := http.NewServeMux()

	// protected wraps a handler with auth middleware.
	// http.HandlerFunc conversion is hidden here so call sites stay clean.
	protected := func(fn http.HandlerFunc) http.Handler {
		return authMW(fn)
	}

	// Spec (no auth).
	mux.HandleFunc("GET /openapi.json", specH.ServeJSON)
	mux.HandleFunc("GET /docs", specH.ServeUI)

	// Health.
	mux.HandleFunc("GET /health", healthH.Live)
	mux.HandleFunc("GET /health/ready", healthH.Ready)

	// Tournaments.
	mux.Handle("POST /tournaments", authMW(http.HandlerFunc(tournamentH.Create)))
	mux.HandleFunc("GET /tournaments", tournamentH.List)
	mux.HandleFunc("GET /tournaments/{id}", tournamentH.GetByID)

	// Teams (protected).
	mux.Handle("GET /tournaments/{tournament_id}/teams", protected(teamH.List))

	// Users (protected).
	mux.Handle("GET /users/me", protected(userH.Me))

	// Players (protected).
	mux.Handle("GET /tournaments/{tournament_id}/players/search", protected(playerH.Search))

	// Player handicaps (protected).
	mux.Handle("GET /player-handicaps/{player_id}/{category}", protected(playerHandicapH.Get))

	// Tournament predictions (protected).
	mux.Handle(
		"PUT /tournaments/{tournamentId}/predictions/players",
		protected(tournamentPredictionH.BulkUpsertPlayerPredictions),
	)
	mux.Handle(
		"GET /tournaments/{tournamentId}/predictions/players",
		protected(tournamentPredictionH.ListMyPlayerPredictions),
	)
	mux.Handle(
		"PUT /tournaments/{tournamentId}/predictions/teams",
		protected(tournamentPredictionH.BulkUpsertTeamPredictions),
	)
	mux.Handle(
		"GET /tournaments/{tournamentId}/predictions/teams",
		protected(tournamentPredictionH.ListMyTeamPredictions),
	)
	mux.Handle(
		"GET /leagues/{leagueId}/predictions/players",
		protected(tournamentPredictionH.ListLeaguePlayerPredictions),
	)
	mux.Handle("GET /leagues/{leagueId}/predictions/teams", protected(tournamentPredictionH.ListLeagueTeamPredictions))

	// Score predictions (protected).
	mux.Handle("PUT /predictions/{fixtureId}", protected(scorePredictionH.UpsertScore))
	mux.Handle("GET /tournaments/{tournamentId}/fixtures", protected(fixtureH.ListForUser))
	mux.Handle("GET /leagues/{leagueId}/predictions", protected(fixtureH.ListForLeague))

	// Leagues (protected).
	mux.Handle("POST /leagues", protected(leagueH.Create))
	mux.Handle("GET /leagues", protected(leagueH.List))
	mux.Handle("POST /leagues/join", protected(leagueH.Join))
	mux.Handle("GET /leagues/{id}", protected(leagueH.Get))
	mux.Handle("PATCH /leagues/{id}", protected(leagueH.UpdateName))
	mux.Handle("DELETE /leagues/{id}", protected(leagueH.Delete))
	mux.Handle("POST /leagues/{id}/code", protected(leagueH.RegenerateCode))
	mux.Handle("DELETE /leagues/{id}/members/{userId}", protected(leagueH.RemoveMember))

	// Apply middleware in outer-to-inner order.
	return middleware.RequestID(
		middleware.Logger(
			logger,
			middleware.Recoverer(logger, mux),
		),
	)
}
