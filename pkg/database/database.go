package database

import (
	"book-store/pkg/config"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func NewPostgresDB(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx := context.Background()
	parseConfig, err := pgxpool.ParseConfig(cfg.Database.URI)
	if err != nil {
		return nil, err
	}
	parseConfig.MaxConns = cfg.Database.MaxConns
	parseConfig.MinConns = cfg.Database.MinConns
	parseConfig.MaxConnLifetime = cfg.Database.MaxConnLifetime * time.Minute
	parseConfig.MaxConnIdleTime = cfg.Database.MaxConnIdleTime * time.Minute

	db, err := pgxpool.NewWithConfig(ctx, parseConfig)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
