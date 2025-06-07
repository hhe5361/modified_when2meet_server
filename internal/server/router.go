package server

import (
	"better-when2meet/internal/db"
	"better-when2meet/internal/domain/notice"
	"better-when2meet/internal/domain/room"
	"better-when2meet/internal/domain/user"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	database := db.InitDB()
	roomRepo := room.New(database)
	userRepo := user.New(database)
	noticeRepo := notice.New(database)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/rooms", CreateRoomHandler(roomRepo))
	r.POST("/rooms/:url/login", RegisterHandler(roomRepo, userRepo))

	r.GET("/rooms/:url", GetRoomInfoHandler(roomRepo, userRepo))
	r.GET("/rooms/:url/user", JWTRoomAuthMiddleware(), GetUserDetailHandler(userRepo))

	r.PUT("/rooms/:url/times", JWTRoomAuthMiddleware(), VoteTimeHandler(roomRepo, userRepo))

	r.GET("/rooms/:url/notices", GetNoticeHandler(noticeRepo, roomRepo))
	r.POST("/rooms/:url/notices", JWTRoomAuthMiddleware(), CreateNoticeHandler(noticeRepo))

	//delete and put will be added ... maybe

	return r
}

func InitServer() error {
	r := SetupRouter()
	return r.Run(":8080")
}
