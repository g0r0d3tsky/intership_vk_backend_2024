package repository

import (
	"cinema_service/config"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var (
	ErrURLNotFound    = errors.New("url not found")
	ErrDuplicateLogin = errors.New("duplicate login")
)

func Connect(c *config.Config) (*pgxpool.Pool, error) {
	connectionString := c.PostgresDSN()

	poolConfig, err := pgxpool.ParseConfig(connectionString)
	//fmt.Println(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx pool config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return pool, nil
}