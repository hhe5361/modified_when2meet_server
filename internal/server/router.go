package server

import (
	"better-when2meet/internal/db"
	"better-when2meet/internal/room"

	"github.com/gin-gonic/gin"
)

func InitServer() error {
	// Initialize database connection
	database := db.InitDB()
	roomRepo := room.New(database)

	// Initialize Gin router
	r := gin.Default()

	// Room endpoints
	r.POST("/rooms", CreateRoomHandler(roomRepo))
	// r.GET("/rooms/:url", GetRoomInfoHandler(roomRepo))
	r.POST("/rooms/:url")
	// Start server
	return r.Run(":8080")
}
