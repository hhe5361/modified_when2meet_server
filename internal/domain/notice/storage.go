package notice

import (
	"better-when2meet/internal/db"
	"database/sql"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (n *Storage) GetByRoomID(id int) ([]DetailNotice, error) {
	query := `
	SELECT 
		n.id AS id,
		n.content AS content,
		n.created_at AS created_at,
		n.updated_at AS updated_at,
		u.name AS user_name
	FROM notices n
	JOIN user u ON n.user_id = u.id
	WHERE n.room_id = ?
	`
	return db.QueryRows(n.db, query, db.ScanStruct[DetailNotice], id)
}

func (n *Storage) Insert(data CreateNoticeReq, roomID int, userID int) error {
	query := `INSERT INTO notices (room_id, user_id, content) 
		VALUES (?, ?, ?)`

	_, err := db.QueryExec(n.db, query, roomID, userID, data.Content)
	return err
}

func (n *Storage) DeleteById(id int) error {
	query := "DELETE FROM notices WHERE id = ?"
	_, e := db.QueryExec(n.db, query, id)
	return e
}

//PUT 은 생략 ..
