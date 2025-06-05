package user

import "time"

type ReqLogin struct {
	Name       string `json:"name"`
	Password   string `json:"password"`
	TimeRegion string `json:"time_region"`
}

type ResUser struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	TimeRegion string `json:"time_region"`
}

type UserDetail struct {
	User          ResUser         `json:"user"`
	AvailableTime []AvailableTime `json:"available_time"`
}

type ReqAvailableTime struct {
	Date          time.Time `json:"date"`
	HourStartSlot int       `json:"hour_start_slot"`
	HourEndSlot   int       `json:"hour_end_slot"`
}

type ReqAvailableTimeList struct {
	Times []ReqAvailableTime `json:"times"`
}
