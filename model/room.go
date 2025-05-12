package model

import "time"

type Room struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	VoteOpenAt  time.Time `json:"vote_open_at"`
	VoteCloseAt time.Time `json:"vote_close_at"`
	TimeRegion  string    `json:"time_region"`
	IsOnline    bool      `json:"isonline"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
