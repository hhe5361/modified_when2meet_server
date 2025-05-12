package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"example.com/modified_when2meet_server/model"
	"github.com/gin-gonic/gin"
)

var db *sql.DB

func main() {
	var err error

	//connection to mysql db
	dsn := "root:password@tcp(localhost:3306)/when2meet?parseTime=true"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("DB not responding:", err)
	}

	r := gin.Default()
	r.POST("/rooms", createRoom)
	r.GET("/rooms/:url", getRoom)
	r.POST("/rooms/:url/login", getRoom) //login to room
	r.POST("/rooms/:url/vote", voteAvailableTime)

	//추가로 해야 하는 것 : 채팅 기능 등.
	r.Run(":8080")
}

func createRoom(c *gin.Context) {
	var room model.Room

	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 유저 존재 확인
	var userID int
	err := db.QueryRow("SELECT id FROM user WHERE id = ?", room.UserCreatorID).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	now := time.Now()
	result, err := db.Exec(`INSERT INTO room (name, url, user_creator_id, vote_open_at, vote_close_at, time_region, isonline, createdAt, updatedAt) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		room.Name, room.URL, room.UserCreatorID, room.VoteOpenAt, room.VoteCloseAt, room.TimeRegion, room.IsOnline, now, now)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
		return
	}

	roomID, _ := result.LastInsertId()

	// 유저에 room_id 업데이트
	_, err = db.Exec("UPDATE user SET room_id = ? WHERE id = ?", roomID, room.UserCreatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign user to room"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Room created", "room_id": roomID})
}

// GET /rooms/:id
func getRoom(c *gin.Context) {
	roomID := c.Param("id")

	var room model.Room
	err := db.QueryRow("SELECT id, name, url, user_creator_id, vote_open_at, vote_close_at, time_region, isonline, createdAt, updatedAt FROM room WHERE id = ?", roomID).
		Scan(&room.ID, &room.Name, &room.URL, &room.UserCreatorID, &room.VoteOpenAt, &room.VoteCloseAt, &room.TimeRegion, &room.IsOnline, &room.CreatedAt, &room.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	users := []model.User{}
	rows, err := db.Query("SELECT id, name, role FROM user WHERE room_id = ?", roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Role); err == nil {
			users = append(users, u)
		}
	}

	times := []model.AvailableTime{}
	timeRows, err := db.Query("SELECT user_id, date, hour FROM available_time WHERE room_id = ?", roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve available times"})
		return
	}
	defer timeRows.Close()

	for timeRows.Next() {
		var at model.AvailableTime
		if err := timeRows.Scan(&at.UserID, &at.Date, &at.Hour); err == nil {
			times = append(times, at)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"room":            room,
		"users":           users,
		"available_times": times,
	})
}

// POST /rooms/:id/vote
func voteAvailableTime(c *gin.Context) {
	roomID := c.Param("id")
	var vote model.AvailableTime

	if err := c.ShouldBindJSON(&vote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 유저 방 참가 확인
	var check int
	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE id = ? AND room_id = ?", vote.UserID, roomID).Scan(&check)
	if err != nil || check == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found in room"})
		return
	}

	now := time.Now()
	_, err = db.Exec(`INSERT INTO available_time (user_id, room_id, date, hour, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?)`,
		vote.UserID, roomID, vote.Date, vote.Hour, now, now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record vote"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Vote recorded successfully"})
}
