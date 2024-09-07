// Package server provides the main server and HTTP handlers for the application.
package server

import (
	"os"

	"github.com/rs/zerolog"

	"github.com/jdks/fiber-example/internal/config"
	"github.com/jdks/fiber-example/internal/store"
)

type Server struct {
	store  *store.Store
	config *config.Config
	log    zerolog.Logger
}

func New(cfg *config.Config) (Server, error) {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	store, err := store.New(cfg.Store)
	if err != nil {
		return Server{}, err
	}

	return Server{
		store:  store.WithLogger(log),
		config: cfg,
		log:    log.With().Str("component", "server").Logger(),
	}, nil
}
