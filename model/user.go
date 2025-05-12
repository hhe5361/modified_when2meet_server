package model

import "time"

type User struct {
	ID         int       `json:"id"`
	RoomID     int       `json:"room_id"`
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	Role       string    `json:"role"`
	TimeRegion string    `json:"time_region"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
