package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgreClient struct {
	Pool *pgxpool.Pool
}

func NewPostgreClient(ctx context.Context, dsn string) (*PostgreClient, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errors.New("failed to parse database configuration")
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.New("failed to create database pool")
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, errors.New("failed to ping database")
	}

	return &PostgreClient{Pool: pool}, nil
}

func (p *PostgreClient) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
