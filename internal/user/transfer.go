package user

import "time"

type ReqLogin struct {
	Name       string `json:"name"`
	Password   string `json:"password"`
	TimeRegion string `json:"time_region"`
}

type UserDetail struct {
	User          User
	AvailableTime []AvailableTime
}

type ReqAvailableTime struct {
	Date          time.Time `json:"date"`
	HourEndSlot   int       `json:"hour_end_slot"`
	HourStartSlot int       `json:"hour_start_slot"`
}
