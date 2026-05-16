// Package db provides database connection pool setup and lifecycle management.
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool wraps a pgx connection pool with our own helpers.
type Pool struct {
	*pgxpool.Pool
}

// New creates a connection pool to the database at databaseURL,
// verifies connectivity, and returns it ready to use.
//
// The caller is responsible for calling Close() when shutting down.
func New(
	ctx context.Context,
	databaseURL string,
) (*Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}

	// Pool tuning for a small EC2 box. RDS db.t4g.micro has a default
	// max_connections of ~80; staying well below that leaves headroom.
	cfg.MaxConns = 20
	cfg.MinConns = 2
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	// Verify the DB is actually reachable, not just that the pool was created.
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Pool{Pool: pool}, nil
}

// HealthCheck verifies the DB is reachable. Used by the readiness probe.
func (p *Pool) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return p.Ping(ctx)
}
