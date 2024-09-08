// Package store provides functionality for interacting with the database.
package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/jdks/fiber-example/internal/config"
	"github.com/rs/zerolog"
)

type Store struct {
	db  *pgxpool.Pool
	log zerolog.Logger
}

func New(cfg config.StoreConfig) (Store, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DBConnectionURL)
	if err != nil {
		log.Error().Err(err).Msg("failed to create store")
		return Store{}, fmt.Errorf("failed to create store: %w", err)
	}

	return Store{
		db: pool,
	}, nil
}

func (s *Store) WithLogger(log zerolog.Logger) *Store {
	s.log = log.With().Str("component", "store").Logger()
	return s
}
