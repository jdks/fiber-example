package store

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gofrs/uuid"

	table "github.com/jdks/fiber-example/internal/store/users"
)

type User struct {
	ID        uuid.UUID `db:"user_id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
}

func (s Store) CreateUser(ctx context.Context, user User) error {
	query := sq.Insert(table.Table).
		Columns(table.UserID, table.FirstName, table.LastName).
		Values(user.ID, user.FirstName, user.LastName).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	_, err = s.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (s Store) GetUser(ctx context.Context, id string) (User, error) {
	var user User

	sql, args, err := sq.Select(table.Columns...).
		From(table.Table).
		Where(sq.Eq{table.Identifier: id}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return user, fmt.Errorf("failed to build SQL query: %w", err)
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return user, fmt.Errorf("failed to execute query: %w", err)
	}

	err = pgxscan.ScanOne(&user, rows)
	if err != nil {
		return user, fmt.Errorf("failed to scan user: %w", err)
	}

	return user, nil
}

func (s Store) GetAllUsers(ctx context.Context, pageSize, pageNumber int) ([]*User, error) {
	users := make([]*User, 0, pageSize)

	builder := sq.Select(table.Columns...).From(table.Table)
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
