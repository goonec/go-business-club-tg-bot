package entity

import "time"

type User struct {
	ID         int64     `json:"id"`
	UsernameTG string    `json:"tg_username"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"create_at"`
}
