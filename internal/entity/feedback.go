package entity

import "time"

type Feedback struct {
	ID         int       `json:"id"`
	Message    string    `json:"message"`
	UsernameTG string    `json:"tg_username"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"created_at"`
}
