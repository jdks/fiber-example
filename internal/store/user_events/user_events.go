// Package user_events contains the table definitions for the user_events table.
package user_events

const (
	Table             = "user_events"
	EventID           = "event_id"
	UserID            = "user_id"
	CreatedAt         = "created_at"
	AssociatedUserIDs = "associated_user_ids"
	Payload           = "payload"
	Processed         = "processed"

	Identifier = EventID
)

var Columns = []string{
	EventID,
	UserID,
	CreatedAt,
	AssociatedUserIDs,
	Payload,
}
