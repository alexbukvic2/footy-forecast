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

	mux := http.NewServeMux()

	// protected wraps a handler with auth middleware.
	// http.HandlerFunc conversion is hidden here so call sites stay clean.
	protected := func(fn http.HandlerFunc) http.Handler {
		return authMW(fn)
	}

	// Health.
	mux.HandleFunc("GET /health", healthH.Live)
	mux.HandleFunc("GET /health/ready", healthH.Ready)

	// Tournaments.
	mux.Handle("POST /tournaments", authMW(http.HandlerFunc(tournamentH.Create)))
	mux.HandleFunc("GET /tournaments", tournamentH.List)
	mux.HandleFunc("GET /tournaments/{id}", tournamentH.GetByID)
	mux.HandleFunc("GET /tournaments/slug/{slug}", tournamentH.GetBySlug)

	// Users (protected).
	mux.Handle("GET /users/me", protected(userH.Me))

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
