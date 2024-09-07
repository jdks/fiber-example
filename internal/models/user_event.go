package models

import "time"

type UserEvent struct {
	ID                string    `json:"id"`
	SessionID         string    `json:"session_id"`
	UserID            string    `json:"user_id"`
	Action            string    `json:"action"`
	CreatedAt         time.Time `json:"created_at"`
	AssociatedUserIDs []string  `json:"associated_user_ids"`
}
