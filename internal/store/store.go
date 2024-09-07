// Package store provides functionality for interacting with the database.
package store

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/jdks/fiber-example/internal/config"
	"github.com/rs/zerolog"
)

const (
	usersTable      = "users"
	userEventsTable = "user_events"
	maxPageSize     = 500
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

func (s Store) GetUser(ctx context.Context, id string) (User, error) {
	var user User

	sql, _, err := sq.Select("*").From("users").Where(sq.Eq{"user_id": id}).ToSql()
	if err != nil {
		return user, fmt.Errorf("failed to build SQL query: %w", err)
	}

	rows, err := s.db.Query(ctx, sql)
	if err != nil {
		return user, fmt.Errorf("failed to execute query: %w", err)
	}

	err = pgxscan.ScanOne(&user, rows)
	if err != nil {
		return user, fmt.Errorf("failed to scan user: %w", err)
	}

	return user, nil
}

func withPagination(pageSize, pageNumber int) func(sq.SelectBuilder) sq.SelectBuilder {
	return func(builder sq.SelectBuilder) sq.SelectBuilder {
		if pageSize > maxPageSize {
			pageSize = maxPageSize
		}
		offset := (pageNumber - 1) * pageSize
		return builder.Limit(uint64(pageSize)).Offset(uint64(offset))
	}
}

func (s Store) QueryEvents(ctx context.Context, params EventQueryParams, pageSize, pageNumber int) ([]*UserEvent, error) {
	userEvents := make([]*UserEvent, 0, pageSize)

	eventQuery := params.ToEventQuery()
	paginatedBuilder := withPagination(pageSize, pageNumber)
	columns := []string{
		"event_id",
		"user_id",
		"created_at",
		"payload",
		"associated_user_ids",
	}

	builder := sq.Select(columns...).
		From(userEventsTable).
		Where(sq.Eq{"processed": false}).
		OrderBy("created_at ASC").
		PlaceholderFormat(sq.Dollar)

	if string(eventQuery.payloadQuery) != "{}" {
		builder = builder.Where(fmt.Sprintf("payload @> '%s'", eventQuery.payloadQuery))
	}
	if eventQuery.eventID != "" {
		builder = builder.Where(sq.Eq{"event_id": eventQuery.eventID})
	}
	if eventQuery.userID != "" {
		builder = builder.Where(sq.Eq{"user_id": eventQuery.userID})
	}

	builder = paginatedBuilder(builder)
	sql, args, err := builder.ToSql()
	s.log.Debug().Str("sql", sql).Msg("")
	if err != nil {
		s.log.Error().Err(err).Msg("")
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	s.log.Debug().Msgf("rows: %v", rows)
	err = pgxscan.ScanAll(&userEvents, rows)
	if err != nil {
		s.log.Error().Err(err).Msg("")
		return nil, fmt.Errorf("failed to scan user events: %w", err)
	}

	return userEvents, nil
}

func (s Store) GetAllUsers(ctx context.Context, pageSize, pageNumber int) ([]*User, error) {
	users := make([]*User, 0, pageSize)

	builder := sq.Select("*").From(usersTable)
	paginatedBuilder := withPagination(pageSize, pageNumber)
	builder = paginatedBuilder(builder)
	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	err = pgxscan.ScanAll(&users, rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan users: %w", err)
	}

	return users, nil
}

func (s *Store) WithLogger(log zerolog.Logger) *Store {
	s.log = log.With().Str("component", "store").Logger()
	return s
}
