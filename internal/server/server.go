package server

import (
	"github.com/jdks/fiber-example/internal/config"
	"github.com/jdks/fiber-example/internal/store"
)

type Server struct {
	store  *store.Store
	config *config.Config
}

func New(cfg *config.Config) (Server, error) {
	store, err := store.New(cfg.Store)
	if err != nil {
		return Server{}, err
	}
	return Server{
		store: &store,
	}, nil
}
