package server

import (
	"better-when2meet/internal/db"
	"better-when2meet/internal/room"
	"better-when2meet/internal/user"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	database := db.InitDB()
	roomRepo := room.New(database)
	userRepo := user.New(database)

	r := gin.Default()

	r.POST("/rooms", CreateRoomHandler(roomRepo))
	r.POST("/rooms/:url/login", RegisterHandler(roomRepo, userRepo))
	r.GET("/rooms/:url", GetRoomInfoHandler(roomRepo, userRepo))
	return r
}

func InitServer() error {
	r := SetupRouter()
	return r.Run(":8080")
}
