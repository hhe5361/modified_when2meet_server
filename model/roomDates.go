package model

import "time"

type RoomDate struct {
	ID        int       `json:"id"`
	RoomID    int       `json:"room_id"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
