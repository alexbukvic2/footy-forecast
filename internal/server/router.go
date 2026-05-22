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

	mux := http.NewServeMux()

	// Health.
	mux.HandleFunc("GET /health", healthH.Live)
	mux.HandleFunc("GET /health/ready", healthH.Ready)

	// Tournaments (public).
	mux.HandleFunc("POST /tournaments", tournamentH.Create)
	mux.HandleFunc("GET /tournaments", tournamentH.List)
	mux.HandleFunc("GET /tournaments/{id}", tournamentH.GetByID)
	mux.HandleFunc("GET /tournaments/slug/{slug}", tournamentH.GetBySlug)

	// Users (protected).
	mux.Handle("GET /users/me", authMW(http.HandlerFunc(userH.Me)))

	// Apply middleware in outer-to-inner order.
	return middleware.RequestID(
		middleware.Logger(
			logger,
			middleware.Recoverer(logger, mux),
		),
	)
}
