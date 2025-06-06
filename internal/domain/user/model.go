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

func (a *AvailableTime) ToRes() ResAvailableTime {
	return ResAvailableTime{
		ID:            a.ID,
		UserID:        a.UserID,
		Date:          a.Date.Format("2006-01-02"),
		HourEndSlot:   a.HourEndSlot,
		HourStartSlot: a.HourStartSlot,
	}
}
