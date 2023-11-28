package entity

import "time"

type User struct {
	ID         int64     `json:"id"`
	UsernameTG string    `json:"tg_username"`
	CreatedAt  time.Time `json:"create_at"`
	Role       string    `json:"role"`
}
