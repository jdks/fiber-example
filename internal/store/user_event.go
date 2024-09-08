package store

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gofrs/uuid"

	table "github.com/jdks/fiber-example/internal/store/user_events"
)

type UserEvent struct {
	EventID           uuid.UUID
	UserID            uuid.UUID
	CreatedAt         time.Time
	AssociatedUserIDs []uuid.UUID  `db:"associated_user_ids"`
	Payload           EventPayload `json:"-" db:"payload"`
	Processed         bool         `db:"processed"`
}

func (s *Store) CreateUserEvent(event UserEvent) error {
	query := sq.Insert(table.Table).
		Columns(table.EventID, table.UserID, table.CreatedAt, table.AssociatedUserIDs, table.Payload, table.Processed).
		Values(event.EventID, event.UserID, event.CreatedAt, event.AssociatedUserIDs, event.Payload, event.Processed).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build user event query: %w", err)
	}

	_, err = s.db.Exec(context.Background(), sql, args...)
	if err != nil {
		return fmt.Errorf("failed to create user event: %w", err)
	}
	return nil
}

func (s Store) QueryEvents(ctx context.Context, params EventQueryParams, pageSize, pageNumber int) ([]*UserEvent, error) {
	userEvents := make([]*UserEvent, 0, pageSize)

	eventQuery := params.ToEventQuery()
	paginatedBuilder := withPagination(pageSize, pageNumber)

	builder := sq.Select(table.Columns...).
		From(table.Table).
		Where(sq.Eq{table.Processed: false}).
		OrderBy(fmt.Sprintf("%s ASC", table.CreatedAt)).
		PlaceholderFormat(sq.Dollar)

	if string(eventQuery.payloadQuery) != "{}" {
		builder = builder.Where(fmt.Sprintf("%s @> '%s'", table.Payload, eventQuery.payloadQuery))
	}
	if eventQuery.eventID != "" {
		builder = builder.Where(sq.Eq{table.EventID: eventQuery.eventID})
	}
	if eventQuery.userID != "" {
		builder = builder.Where(sq.Eq{table.UserID: eventQuery.userID})
	}

	builder = paginatedBuilder(builder)
	sql, args, err := builder.ToSql()
	if err != nil {
		s.log.Error().Err(err).Msg("")
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	err = pgxscan.ScanAll(&userEvents, rows)
	if err != nil {
		s.log.Error().Err(err).Msg("")
		return nil, fmt.Errorf("failed to scan user events: %w", err)
	}

	return userEvents, nil
}
