package room

import (
	"better-when2meet/internal/db"
	"database/sql"
	"errors"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

// Scan(*sql.Rows) (T, error)

func scanRoom(rows *sql.Rows) (Room, error) {
	var r Room
	err := rows.Scan(&r.ID, &r.Name, &r.URL, &r.StartTime, &r.EndTime, &r.TimeRegion, &r.IsOnline, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return Room{}, errors.New("server Error: room scan is failed")
	}
	return r, nil
}

func scanRoomDate(rows *sql.Rows) (RoomDate, error) {
	var r RoomDate
	err := rows.Scan(&r.ID, &r.RoomID, &r.Year, &r.Month, &r.Day)
	if err != nil {
		return RoomDate{}, errors.New("server Error : roomdate scan is failed")
	}
	return r, nil
}

func (r *Storage) GetRoomById(id int) (Room, error) {
	query := "SELECT * FROM room WHERE id = ?"
	return db.QueryOnlyRow(r.db, query, scanRoom, id)
}

func (r *Storage) GetRoomByUrl(url string) (Room, error) {
	query := "SELECT * FROM room WHERE url = ?"
	return db.QueryOnlyRow(r.db, query, scanRoom, url)
}

func (r *Storage) InsertRoom(m ReqCreateRoom, url string) error {
	query := `INSERT INTO room (name, url, start_time, end_time, time_region, is_online, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`

	id, err := db.QueryExec(r.db, query, m.RoomName, url, m.StartTime, m.EndTime, m.TimeRegion, m.IsOnline)

	if err != nil {
		return err
	}

	return r.InsertRoomDate(m.VoteableRooms, id)
}

// 연쇄 삭제 일어나게 하려면 db 에서 그냥 설정하는게 빠르려나 ;;
func (r *Storage) DeleteRoom(id int) error {
	query := "DELETE FROM room WHERE id = ?"
	_, e := db.QueryExec(r.db, query, id)
	return e
}

func (r *Storage) InsertRoomDate(dates []ReqRoomDate, roomId int64) error {
	query := `INSERT INTO roomdate (room_id, year, month, day) VALUES (?, ?, ?, ?)`

	for _, v := range dates {
		if _, err := db.QueryExec(r.db, query, roomId, v.Year, v.Month, v.Day); err != nil {
			return err
		}
	}
	return nil
}

func (r *Storage) GetRoomDatesByRoomID(roomId int64) ([]RoomDate, error) {
	query := `SELECT * FROM room_dates WHERE room_id = ?`
	return db.QueryRows(r.db, query, scanRoomDate, roomId)
}
