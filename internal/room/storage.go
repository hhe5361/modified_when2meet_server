package room

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

func (r *Storage) GetRoomById(id int) (Room, error) {
	query := "SELECT * FROM room WHERE id = ?"
	return db.QueryOnlyRow(r.db, query, db.ScanStruct[Room], id)
}

func (r *Storage) GetRoomByUrl(url string) (Room, error) {
	query := "SELECT * FROM room WHERE url = ?"
	return db.QueryOnlyRow(r.db, query, db.ScanStruct[Room], url)
}

func (r *Storage) InsertRoom(m ReqCreateRoom, url string) error {
	query := `INSERT INTO room (name, url, start_time, end_time, time_region, is_online) 
		VALUES (?, ?, ?, ?, ?, ?)`

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
	query := `SELECT * FROM roomdate WHERE room_id = ?`
	return db.QueryRows(r.db, query, db.ScanStruct[RoomDate], roomId)
}

func (r *Storage) GetRoomDetailByUrl(url string) (RoomDetail, error) {
	roomData, err := r.GetRoomByUrl(url)
	if err != nil {
		return RoomDetail{}, err
	}
	dates, err := r.GetRoomDatesByRoomID(roomData.ID)
	if err != nil {
		return RoomDetail{}, err
	}

	return RoomDetail{
		Room:  roomData,
		Dates: dates,
	}, nil
}
