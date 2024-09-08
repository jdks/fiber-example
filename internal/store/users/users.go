// Package users contains the table definitions for the users table.
package users

const (
	Table     = "users"
	UserID    = "user_id"
	FirstName = "first_name"
	LastName  = "last_name"

	Identifier = UserID
)

var Columns = []string{
	UserID,
	FirstName,
	LastName,
}
