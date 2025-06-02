package meeting

// to use time tables
type HourBlock struct {
	Hour  int      `json:"hour"`
	Users []string `json:"user_name"`
}

// this saved like this Date  : [~~~] front should draw based on this struct data
type VoteTable map[string][]HourBlock //string
