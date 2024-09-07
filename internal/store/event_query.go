package store

import (
	"encoding/json"
)

type EventQueryParams struct {
	EventID string `json:"-"`
	UserID  string `json:"-"`

	SessionID  string      `json:"session_id,omitempty"`
	Action     EventAction `json:"action,omitempty"`
	DayInWeek  int         `json:"day_in_week,omitempty"`
	DayInMonth int         `json:"day_in_month,omitempty"`
	Month      int         `json:"month,omitempty"`
	Year       int         `json:"year,omitempty"`
	Hour       int         `json:"hour,omitempty"`
}

type eventQuery struct {
	userID       string
	eventID      string
	payloadQuery []byte
}

func (p EventQueryParams) ToEventQuery() eventQuery {
	eq := eventQuery{}
	query, err := json.Marshal(&p)
	if err != nil {
		return eq
	}
	eq.payloadQuery = query

	if p.EventID != "" {
		eq.eventID = p.EventID
	}
	if p.UserID != "" {
		eq.userID = p.UserID
	}

	return eq
}
