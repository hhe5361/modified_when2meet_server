package server

import (
	"better-when2meet/internal/db"
	"better-when2meet/internal/domain/room"
	"better-when2meet/internal/domain/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	database := db.InitDB()
	roomRepo := room.New(database)
	userRepo := user.New(database)

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/rooms", CreateRoomHandler(roomRepo))
	r.POST("/rooms/:url/login", RegisterHandler(roomRepo, userRepo))

	r.GET("/rooms/:url", GetRoomInfoHandler(roomRepo, userRepo))
	r.GET("/rooms/:url/user", JWTRoomAuthMiddleware(), GetUserDetailHandler(userRepo))

	r.PUT("/rooms/:url/times", JWTRoomAuthMiddleware(), VoteTimeHandler(roomRepo, userRepo))

	return r
}

func InitServer() error {
	r := SetupRouter()
	return r.Run(":8080")
}
