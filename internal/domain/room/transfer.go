package room

type ReqRoomDate struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type ReqCreateRoom struct {
	RoomName      string        `json:"room_name"`
	TimeRegion    string        `json:"time_region"`
	StartTime     int           `json:"start_time"`
	EndTime       int           `json:"end_time"`
	IsOnline      bool          `json:"is_online"`
	VoteableRooms []ReqRoomDate `json:"voteable_rooms"`
}

