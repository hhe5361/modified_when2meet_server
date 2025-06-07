package notice

type CreateNoticeReq struct {
	Content string `json:"content"`
}

type DetailNotice struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	//Join
	UserNickname string `json:"user_name"`
}
