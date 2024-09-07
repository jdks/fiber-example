package store

import (
	"github.com/gofrs/uuid"
)

type EventPayload struct {
	SessionID  uuid.UUID   `json:"session_id"`
	Action     EventAction `json:"action"`
	Hour       uint        `json:"hour"`
	DayInWeek  uint        `json:"day_in_week"`
	DayInMonth uint        `json:"day_in_month"`
	Month      uint        `json:"month"`
	Year       uint        `json:"year"`
}
