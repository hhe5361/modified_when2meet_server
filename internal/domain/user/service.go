package user

import (
	"better-when2meet/internal/db"
	"better-when2meet/internal/helper"
	"database/sql"
	"errors"
	"time"
)

// Service defines the interface for user-related operations
type Service interface {
	TryLogin(pwd string, name string, roomId int64) (ResUser, error)
	InsertUser(r ReqLogin, roomId int64) (int64, error)
	UserDetailById(id int64) (UserDetail, error)
	InsertVoteTime(userId int64, times ReqAvailableTime) error
	DeleteVoteTime(userId int64, date time.Time) error
}

// service implements the Service interface
type service struct {
	db *sql.DB
}

// NewService creates a new instance of the user service
func NewService() Service {
	return &service{
		db: db.InitDB(),
	}
}

// TryLogin attempts to authenticate a user
func (s *service) TryLogin(pwd string, name string, roomId int64) (ResUser, error) {
	query := `SELECT * FROM user WHERE name = ? AND room_id = ?`

	user, err := db.QueryOnlyRow(s.db, query, db.ScanStruct[User], name, roomId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ResUser{}, ErrUserNotFound
		}
		return ResUser{}, err
	}

	if !helper.CheckPasswordHash(pwd, user.Password) {
		return ResUser{}, ErrInvalidPassword
	}

	return user.ToResUser(), nil
}

// InsertUser creates a new user
func (s *service) InsertUser(r ReqLogin, roomId int64) (int64, error) {
	query := `INSERT INTO user (room_id, name, password, time_region) VALUES (?, ?, ?, ?)`
	return db.QueryExec(s.db, query, roomId, r.Name, r.Password, r.TimeRegion)
}

// UserDetailById retrieves user details including available times
func (s *service) UserDetailById(id int64) (UserDetail, error) {
	user, err := s.UserById(id)
	if err != nil {
		return UserDetail{}, err
	}

	times, err := s.TimesByUserId(id)
	if err != nil {
		return UserDetail{}, err
	}

	if times == nil {
		times = []AvailableTime{}
	}

	return UserDetail{user.ToResUser(), times}, nil
}

// UserById retrieves a user by ID
func (s *service) UserById(id int64) (User, error) {
	query := `SELECT * FROM user WHERE id = ?`
	return db.QueryOnlyRow(s.db, query, db.ScanStruct[User], id)
}

// TimesByUserId retrieves available times for a user
func (s *service) TimesByUserId(id int64) ([]AvailableTime, error) {
	query := `SELECT * FROM available_time WHERE user_id = ?`
	return db.QueryRows(s.db, query, db.ScanStruct[AvailableTime], id)
}

// InsertVoteTime records a user's available time
func (s *service) InsertVoteTime(userId int64, times ReqAvailableTime) error {
	query := `INSERT INTO available_time (user_id, date, hour_start_slot, hour_end_slot) VALUES (?, ?, ?, ?)`
	_, err := db.QueryExec(s.db, query, userId, times.Date.Format("2006-01-02"), times.HourStartSlot, times.HourEndSlot)
	return err
}

// DeleteVoteTime removes a user's available time for a specific date
func (s *service) DeleteVoteTime(userId int64, date time.Time) error {
	query := `DELETE FROM available_time WHERE user_id = ? AND date = ?`
	_, err := db.QueryExec(s.db, query, userId, date.Format("2006-01-02"))
	return err
}
