package user

type ReqLogin struct {
	Name     string `json:"name"`
	Password int    `json:"password"`
}
