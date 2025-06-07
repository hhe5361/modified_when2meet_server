package notice

import "time"

type Notice struct {
	ID        int64     `json:"id"`
	RoomID    int64     `json:"room_id"`
	USERID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
