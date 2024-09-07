package store

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
)

type UserEvent struct {
	EventID           uuid.UUID
	UserID            uuid.UUID
	CreatedAt         time.Time
	AssociatedUserIDs []uuid.UUID  `db:"associated_user_ids"`
	Payload           EventPayload `json:"-" db:"payload"`
}

func (s *Store) CreateUserEvent(event UserEvent) error {
	query := sq.Insert("user_events").
		Columns("event_id", "user_id", "created_at", "associated_user_ids", "payload").
		Values(event.EventID, event.UserID, event.CreatedAt, event.AssociatedUserIDs, event.Payload).
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
