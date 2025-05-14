package user

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

func scanUser(rows *sql.Rows) (User, error) {
	var u User
	err := rows.Scan(u.ID, u.RoomID, u.Name, u.Password, u.TimeRegion, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return User{}, errors.New("server Error: room scan is failed")
	}
	return u, nil
}

func scanAvailableTime(rows *sql.Rows) (AvailableTime, error) {
	var a AvailableTime
	if err := rows.Scan(&a.ID, &a.RoomID, &a.UserID, &a.Date, &a.HourEndSlot, &a.HourEndSlot, &a.CreatedAt, &a.UpdatedAt); err != nil {
		return AvailableTime{}, errors.New("fail to scan")
	}
	return a, nil
}

func (u *Storage) GetUserById(id int64) (User, error) {
	query := `select * from user where id = ?`
	return db.QueryOnlyRow(u.db, query, scanUser, id)
}

func (u *Storage) GetUsersByroomId(id int64) ([]User, error) {
	query := `select * from user where room_id = ?`
	return db.QueryRows(u.db, query, scanUser, id)
}

func (u *Storage) InsertUser(r ReqLogin, roomdId int64) (int64, error) {
	query := `insert into user (room_id, name, password, time_region,created_at,updated_at) VALUES (?, ?, ?, ?, ?, ?) `
	return db.QueryExec(u.db, query, roomdId, r.Name, r.Password, r.TimeRegion)
}

func (u *Storage) GetTimesByUserId(id int64) ([]AvailableTime, error) {
	query := `select * from available_time where user_id = ?`

	return db.QueryRows(u.db, query, scanAvailableTime, id)
}
