package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxConns        = 25
	minConns        = 5
	maxConnLifetime = time.Hour
	maxConnIdleTime = 30 * time.Minute
)

func InitDatabase(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	config.MaxConns = maxConns
	config.MinConns = minConns
	config.MaxConnLifetime = maxConnLifetime
	config.MaxConnIdleTime = maxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database unreachable: %w", err)
	}

	return pool, nil
}
