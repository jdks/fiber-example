package store

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jdks/fiber-example/internal/config"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(cfg config.StoreConfig) (Store, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DBConnectionURL)
	if err != nil {
		return Store{}, err
	}

	return Store{
		pool: pool,
	}, nil
}

func (s Store) GetUser(ctx context.Context, id string) (User, error) {
	rows, err := s.pool.Query(ctx, "select * from users where user_id=$1", id)
	if err != nil {
		return User{}, err
	}

	var user User
	if err = pgxscan.ScanOne(&user, rows); err != nil {
		return User{}, err
	}

	return user, nil
}

func (s Store) GetAllUsers(ctx context.Context) ([]*User, error) {
	users := make([]*User, 0, 10)
	rows, err := s.pool.Query(ctx, "select * from users limit 10")
	if err != nil {
		return users, err
	}

	if err = pgxscan.ScanAll(&users, rows); err != nil {
		return users, err
	}

	return users, nil
}
