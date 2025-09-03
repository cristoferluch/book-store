package database

import (
	"book-store/pkg/config"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.Config) (*pgxpool.Pool, error) {

	ctx := context.Background()

	parseConfig, err := pgxpool.ParseConfig(cfg.Database.URI)
	if err != nil {
		log.Fatal(err)
	}

	parseConfig.MaxConns = cfg.Database.MaxConns

	db, err := pgxpool.NewWithConfig(ctx, parseConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
