package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riazahmedshah/vfs/internal/config"
)

type Database struct {
	Pool *pgxpool.Pool
}

const DatabasePingTimeout = 10

func New(cfg *config.Config) (*Database, error) {
	dsn := cfg.Database.DatabaseURL

	pgxPoolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("falied to parse pgx pool config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	database := &Database{
		Pool: pool,
	}
	ctx, cancel := context.WithTimeout(context.Background(), DatabasePingTimeout*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	slog.Info("connected to database")

	return database, nil
}
