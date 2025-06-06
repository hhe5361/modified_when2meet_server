package user

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
	User             ResUser            `json:"user"`
	ResAvailableTime []ResAvailableTime `json:"available_time"`
}

type ResAvailableTime struct {
	ID            int64  `json:"id"`
	UserID        int64  `json:"user_id"`
	Date          string `json:"date"`
	HourEndSlot   int    `json:"hour_end_slot"`
	HourStartSlot int    `json:"hour_start_slot"`
}

type ReqAvailableTime struct {
	Date          string `json:"date"`
	HourStartSlot int    `json:"hour_start_slot"`
	HourEndSlot   int    `json:"hour_end_slot"`
}

type ReqAvailableTimeList struct {
	Times []ReqAvailableTime `json:"times"`
}
