package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPGXPool(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse pgxpool: %w", err)
	}
	// optional tuning:
	// cfg.MaxConns = 20
	// cfg.MinConns = 2
	// cfg.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, cfg) // ✅ context tidak nil
	if err != nil {
		return nil, fmt.Errorf("connect pgxpool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping pgxpool: %w", err)
	}
	return pool, nil
}
