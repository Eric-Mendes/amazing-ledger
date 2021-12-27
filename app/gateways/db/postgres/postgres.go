package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

func ConnectPool(ctx context.Context, connString string, logger zerolog.Logger) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parsing postgres config: %w", err)
	}

	config.ConnConfig.Logger = zerologadapter.NewLogger(logger)

	db, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("connecting to poll: %w", err)
	}

	return db, nil
}

func Connect(ctx context.Context, connString string, logger zerolog.Logger) (*pgx.Conn, error) {
	config, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parsing postgres config: %w", err)
	}

	config.Logger = zerologadapter.NewLogger(logger)

	db, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("connecting to poll: %w", err)
	}

	return db, nil
}
