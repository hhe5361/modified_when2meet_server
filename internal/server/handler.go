package server

import (
	"better-when2meet/internal/room"
	"better-when2meet/internal/user"

	// "better-when2meet/internal/user"

	"github.com/gin-gonic/gin"
)

// 방 생성
func CreateRoomHandler(strg *room.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newRoom room.ReqCreateRoom
		if err := c.ShouldBindJSON(&newRoom); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := strg.InsertRoom(newRoom, generateURL()); err != nil {
			c.JSON(500, gin.H{"error": "Failed to create room"})
			return
		}

		c.JSON(200, gin.H{"message": "Room created successfully"})
	}
}

// 방 정보 조회
// room 관련 정보, users , user의 availble time
func GetRoomInfoHandler(strRoom *room.Storage, strgUser *user.Storage) gin.HandlerFunc {
	type Res struct{
		roomInfo room.Room
		users [] map[] 
	} 
	return func(c *gin.Context) {
		url := c.Param("url")
		room, err := strRoom.GetRoomByUrl(url)
		if err != nil {
			c.JSON(404, gin.H{"error": "Room not found"})
			return
		}
		users, err := strgUser.GetUsersByroomId(int64(room.ID))
		if err != nil {
			c.JSON(404, gin.H{"error": "get user failed"})
			return
		}
		for 

	}
}

// // user 로그인
// func GetRoomInfoHandler(strg *room.Storage) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var req user.ReqLogin
// 		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
// 			c.JSON(400, gin.H{"error": err.Error()})
// 			return
// 		}

// 	}
// }

// r.DELETE("/rooms/:id", func(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(400, gin.H{"error": "Invalid ID format"})
// 		return
// 	}

// 	if err := roomRepo.DeleteRoom(id); err != nil {
// 		c.JSON(500, gin.H{"error": "Failed to delete room"})
// 		return
// 	}
// 	c.JSON(200, gin.H{"message": "Room deleted successfully"})
// })
