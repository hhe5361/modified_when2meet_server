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
	UserID        int64     `json:"user_id"` //user id 에만 종속으로 할 것인지 ?
	Date          time.Time `json:"date"`    //date 를 이따구로 받네 .. 아오 ..
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
