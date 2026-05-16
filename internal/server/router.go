// Package server wires together the application's HTTP layer:
// routes, middleware, and handlers.
package server

import (
	"log/slog"
	"net/http"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
	"github.com/alexbukvic2/footy-forecast/internal/server/handler"
	"github.com/alexbukvic2/footy-forecast/internal/server/middleware"
	"github.com/alexbukvic2/footy-forecast/internal/service"
)

// NewRouter wires up all HTTP routes and returns the root handler.
func NewRouter(
	logger *slog.Logger,
	pool *db.Pool,
) http.Handler {
	// Wire dependencies bottom-up.
	tournamentRepo := repository.NewTournamentRepository(pool)
	tournamentSvc := service.NewTournamentService(tournamentRepo)

	// Handlers.
	healthH := handler.NewHealth(logger, pool)
	tournamentH := handler.NewTournament(logger, tournamentSvc)

	mux := http.NewServeMux()

	// Health.
	mux.HandleFunc("GET /health", healthH.Live)
	mux.HandleFunc("GET /health/ready", healthH.Ready)

	// Tournaments.
	mux.HandleFunc("POST /tournaments", tournamentH.Create)
	mux.HandleFunc("GET /tournaments", tournamentH.List)
	mux.HandleFunc("GET /tournaments/{id}", tournamentH.GetByID)
	mux.HandleFunc("GET /tournaments/slug/{slug}", tournamentH.GetBySlug)

	// Apply middleware in outer-to-inner order.
	return middleware.RequestID(
		middleware.Logger(
			logger,
			middleware.Recoverer(logger, mux),
		),
	)
}
