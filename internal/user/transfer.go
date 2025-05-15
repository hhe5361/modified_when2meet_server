package user

type ReqLogin struct {
	Name       string `json:"name"`
	Password   string `json:"password"`
	TimeRegion string `json:"time_region"`
}

type UserDetail struct {
	User          User
	AvailableTime []AvailableTime
}
