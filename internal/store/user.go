package store

import (
	"github.com/gofrs/uuid"
)

type User struct {
	ID        uuid.UUID `db:"user_id"`
	FirstName string
	LastName  string
}
