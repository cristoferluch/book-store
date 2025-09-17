package database

import (
	"book-store/pkg/config"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func NewPostgresDB() (*pgxpool.Pool, error) {
	ctx := context.Background()

	parseConfig, err := pgxpool.ParseConfig(config.Cfg.Database.URI)
	if err != nil {
		return nil, err
	}

	parseConfig.MaxConns = config.Cfg.Database.MaxConns
	parseConfig.MinConns = config.Cfg.Database.MinConns
	parseConfig.MaxConnLifetime = config.Cfg.Database.MaxConnLifetime * time.Minute
	parseConfig.MaxConnIdleTime = config.Cfg.Database.MaxConnIdleTime * time.Minute

	db, err := pgxpool.NewWithConfig(ctx, parseConfig)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
