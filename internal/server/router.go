package router

import (
	"better-when2meet/internal/db"
	"better-when2meet/internal/room"
	"better-when2meet/internal/user"

	"github.com/gin-gonic/gin"
)

func InitServer() error {
	// Initialize database connection
	database := db.InitDB()
	roomRepo := room.New(database)
	userRepo := user.New(database)
	// Initialize Gin router
	r := gin.Default()

	// Room endpoints
	r.POST("/rooms", CreateRoomHandler(roomRepo))
	// r.GET("/rooms/:url", GetRoomInfoHandler(roomRepo))
	r.POST("/rooms/:url")

	//user login or register
	r.POST("/rooms/:url/login", RegisterHandler(roomRepo, userRepo))

	// Start server
	return r.Run(":8080")
}
