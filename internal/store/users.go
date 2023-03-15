package store

type User struct {
	ID        string `db:"user_id"`
	FirstName string
	LastName  string
}
