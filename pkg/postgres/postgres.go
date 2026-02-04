package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(url string, poolSize int) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres:New:ParseConfig - %w", err)
	}

	poolConfig.MaxConns = int32(poolSize)

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("postgres:New:NewWithConfig - %w", err)
	}

	return pool, nil
}
