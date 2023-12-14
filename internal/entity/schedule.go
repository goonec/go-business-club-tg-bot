package entity

import "time"

type Schedule struct {
	ID          int       `json:"id"`
	PhotoFileID string    `json:"photo_file_id"`
	CreatedAt   time.Time `json:"created_at"`
}
