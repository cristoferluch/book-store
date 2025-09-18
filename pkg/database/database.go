package database

import (
	"book-store/configs"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func NewPostgresDB() (*pgxpool.Pool, error) {
	ctx := context.Background()

	parseConfig, err := pgxpool.ParseConfig(configs.Cfg.Database.URI)
	if err != nil {
		return nil, err
	}

	parseConfig.MaxConns = configs.Cfg.Database.MaxConns
	parseConfig.MinConns = configs.Cfg.Database.MinConns
	parseConfig.MaxConnLifetime = configs.Cfg.Database.MaxConnLifetime * time.Minute
	parseConfig.MaxConnIdleTime = configs.Cfg.Database.MaxConnIdleTime * time.Minute

	db, err := pgxpool.NewWithConfig(ctx, parseConfig)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
