package room

import (
	"time"
)

type Room struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	StartTime  int       `json:"start_time"` //0
	EndTime    int       `json:"end_time"`   //23
	TimeRegion string    `json:"time_region"`
	IsOnline   bool      `json:"is_online"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type RoomDate struct {
	ID     int64 `json:"id"`
	RoomID int64 `json:"room_id"`
	Year   int   `json:"year"`
	Month  int   `json:"month"`
	Day    int   `json:"day"`
}

type RoomDetail struct {
	Room  Room       `json:"room"`
	Dates []RoomDate `json:"dates"`
}

var allowedTimeRegions = map[string]bool{
	"Asia/Seoul":       true,
	"UTC":              true,
	"America/New_York": true,
	"Europe/London":    true,
}
