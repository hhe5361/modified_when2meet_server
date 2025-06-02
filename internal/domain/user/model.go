package user

import "time"

type User struct {
	ID         int64     `json:"id"`
	RoomID     int64     `json:"room_id"`
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	TimeRegion string    `json:"time_region"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AvailableTime struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Date          time.Time `json:"date"`
	HourEndSlot   int       `json:"hour_end_slot"`
	HourStartSlot int       `json:"hour_start_slot"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (u *User) ToResUser() ResUser {
	return ResUser{
		ID:         u.ID,
		Name:       u.Name,
		TimeRegion: u.TimeRegion,
	}
}
