package user

import (
	"better-when2meet/internal/db"
	"better-when2meet/internal/helper"
	"database/sql"
	"errors"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (u *Storage) UserById(id int64) (User, error) {
	query := `select * from user where id = ?`
	return db.QueryOnlyRow(u.db, query, db.ScanStruct[User], id)
}

func (u *Storage) UsersByroomId(id int64) ([]User, error) {
	query := `select * from user where room_id = ?`
	return db.QueryRows(u.db, query, db.ScanStruct[User], id)
}

func (u *Storage) InsertUser(r ReqLogin, roomdId int64) (int64, error) {
	query := `insert into user (room_id, name, password, time_region) VALUES (?, ?, ?, ?) `
	return db.QueryExec(u.db, query, roomdId, r.Name, r.Password, r.TimeRegion)
}

func (u *Storage) TimesByUserId(id int64) ([]AvailableTime, error) {
	query := `select * from available_time where user_id = ?`

	return db.QueryRows(u.db, query, db.ScanStruct[AvailableTime], id)
}

func (u *Storage) UserDetailById(id int64) (UserDetail, error) {
	user, e := u.UserById(id)
	if e != nil {
		return UserDetail{}, e
	}

	times, e := u.TimesByUserId(id)
	if e != nil {
		return UserDetail{}, e
	}

	return UserDetail{user, times}, nil
}

func (u *Storage) UsersDetailByRoomId(id int64) ([]UserDetail, error) {
	var userDetails []UserDetail

	users, e := u.UsersByroomId(id)
	if e != nil {
		return nil, e
	}

	for _, user := range users {
		var userDetail UserDetail
		userDetail, e = u.UserDetailById(user.ID)
		if e != nil {
			return nil, e
		}

		userDetails = append(userDetails, userDetail)
	}
	return userDetails, nil
}

func (u *Storage) Login(name string, pwd string, roomId int64) (User, error) {
	query := `SELECT * FROM user WHERE name = ? AND room_id = ?`
	user, err := db.QueryOnlyRow(u.db, query, db.ScanStruct[User], name, roomId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	if !helper.CheckPasswordHash(pwd, user.Password) {
		return User{}, ErrInvalidPassword
	}

	return user, nil
}

func (u *Storage) InsertVoteTime(userId int64, times ReqAvailableTime) error {
	query := `INSERT into available_time (user_id, date, hour_start_slot, hour_end_slot) values (?,?,?,?)`
	_, err := db.QueryExec(u.db, query, userId, times.Date.Format("2006-01-02"), times.HourStartSlot, times.HourEndSlot)
	if err != nil {
		return err
	}
	return nil
}

func (u *Storage) DeleteVoteTime(userId int64, date time.Time) error {
	query := `DELETE FROM available_time WHERE user_id = ? AND date = ?`
	_, err := db.QueryExec(u.db, query, userId, date.Format("2006-01-02"))
	return err
}
