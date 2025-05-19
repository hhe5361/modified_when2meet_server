package router

import (
	"better-when2meet/internal/helper"
	"better-when2meet/internal/room"
	"better-when2meet/internal/user"
	"errors"

	"github.com/gin-gonic/gin"
)

// 방 생성
func CreateRoomHandler(strg *room.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newRoom room.ReqCreateRoom
		if err := c.ShouldBindJSON(&newRoom); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}

		if err := strg.InsertRoom(newRoom, helper.GenerateURL()); err != nil {
			c.JSON(500, gin.H{"error": "Failed to create room"})
			return
		}

		c.JSON(200, gin.H{"message": "Room created successfully"})
	}
}

// 방 정보 조회
// room 관련 정보, users , user의 availble gintime
func GetRoomInfoHandler(strRoom *room.Storage, strgUser *user.Storage) gin.HandlerFunc {
	type Data struct {
		roomInfo room.Room
		users    []user.UserDetail
	}

	return func(c *gin.Context) {
		url := c.Param("url")
		room, err := strRoom.GetRoomByUrl(url)
		if err != nil {
			c.JSON(404, gin.H{"error": "Room not found"})
			return
		}
		usersDetail, err := strgUser.UsersDetailByRoomId(int64(room.ID))
		if err != nil {
			c.JSON(404, gin.H{"error": "get user failed"})
			return
		}
		data := Data{room, usersDetail}

		c.JSON(200, gin.H{"message": "Success", "data": data})
	}
}

// register 후 jwt 반환 세션 연결용
func RegisterHandler(rstrg *room.Storage, ustrg *user.Storage) gin.HandlerFunc {

	return func(c *gin.Context) {
		var req user.ReqLogin
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.JSON(404, gin.H{"error": "Can not Bind Body to Json"}) //이 부분 status code 고쳐야 할 듯
		}
		//get roomId
		room, err := rstrg.GetRoomByUrl(c.Param("url"))
		if err != nil {
			c.JSON(404, gin.H{"error": "Room not found"})
			return
		}
		//check if user is existed
		userData, err := ustrg.Login(req.Name, req.Password, room.ID)
		if errors.Is(err, user.ErrUserNotFound) {
			//user가 없을 경우 create
		} else if errors.Is(err, user.ErrInvalidPassword) {
			c.JSON(404, gin.H{"error": "Invalid Password"})
			return
		} else {
			c.JSON(404, gin.H{"error": err})
			return
		}

		//login 성공했을 때
		userDetail, err := ustrg.UserDetailById(userData.ID)
		if err != nil {
			c.JSON(404, gin.H{"error": err})
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
			"data":    userDetail,
		})
	}
}

// user current info 확인
// func GetRoomInfoHandler(strg *room.Storage) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var req user.ReqLogin
// 		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
// 			c.JSON(400, gin.H{"error": err.Error()})
// 			return
// 		}

// 	}
// }

//available time vote 기능

//result 추첨 기능

//get result api ?

//jwt 인증?
