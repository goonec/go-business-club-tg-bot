package entity

type UserResident struct {
	ID         int    `json:"id"`
	UserID     int64  `json:"user_id"`
	UsernameTG string `json:"tg_username"`
}
