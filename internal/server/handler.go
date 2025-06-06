package server

import (
	"better-when2meet/internal/auth"
	"better-when2meet/internal/domain/meeting"
	"better-when2meet/internal/domain/notice"
	"better-when2meet/internal/domain/room"
	"better-when2meet/internal/domain/user"
	"better-when2meet/internal/helper"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// 방 생성
func CreateRoomHandler(strg *room.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newRoom room.ReqCreateRoom
		if err := c.ShouldBindJSON(&newRoom); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Error: "Invalid request format: " + err.Error(),
			})
			return
		}

		//check field type
		if valid, err := room.CheckFieldType(newRoom); !valid {
			c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
			return
		}

		url := helper.GenerateURL()
		if err := strg.InsertRoom(newRoom, url); err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Error: "Failed to create room: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, Response{
			Message: "Room created successfully",
			Data: gin.H{
				"url": url,
			},
		})
	}
}

// 방 정보 조회
// room 관련 정보, users , user의 availble gintime
func GetRoomInfoHandler(strRoom *room.Storage, strgUser *user.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Param("url")
		roomDetail, err := strRoom.GetRoomDetailByUrl(url)
		if err != nil {
			c.JSON(http.StatusNotFound, Response{
				Error: "Room not found",
			})
			return
		}
		usersDetail, err := strgUser.UsersDetailByRoomId(int64(roomDetail.Room.ID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Error: "Failed to get users: " + err.Error(),
			})
			return
		}

		data, err := meeting.ToVoteTable(usersDetail, roomDetail)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Error: "Failed to convert voteTable: " + err.Error(),
			})
		}

		c.JSON(http.StatusOK, Response{
			Message: "Success",
			Data: gin.H{
				"roomInfo":   roomDetail.Room,
				"vote_table": data, //users 에 해당 내용 들어감.
			},
		})
	}
}

// register 후 jwt 반환 세션 연결용 아 함수 너무 긴 것 같음. 리팩토링 ㄱㄱ 책임 분리 중복 최소화
func RegisterHandler(rstrg *room.Storage, ustrg *user.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req user.ReqLogin
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Error: "Invalid request format: " + err.Error(),
			})
			return
		}
		//get roomId
		room, err := rstrg.GetRoomByUrl(c.Param("url"))
		if err != nil {
			c.JSON(http.StatusNotFound, Response{
				Error: "Room not found",
			})
			return
		}
		//check if user is existed
		userData, err := ustrg.Login(req.Name, req.Password, room.ID)

		if errors.Is(err, user.ErrUserNotFound) {
			//user가 없을 경우 create
			encrypt, err := helper.EncryptPassword(req.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, Response{
					Error: "Failed to encrypt user: " + err.Error(),
				})
			}
			req.Password = encrypt

			userId, err := ustrg.InsertUser(req, room.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, Response{
					Error: "Failed to create user: " + err.Error(),
				})
				return
			}
			userData.ID = userId

		} else if errors.Is(err, user.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, Response{
				Error: "Invalid password",
			})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Error: "Login failed: " + err.Error(),
			})
			return
		}

		//login 성공했을 때
		userDetail, err := ustrg.UserDetailById(userData.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Error: "Failed to get user details: " + err.Error(),
			})
			return
		}

		token, err := auth.GenerateJWT(userData.ID, room.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Error: "Failed to generate token: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, Response{
			Message: "Success",
			Data: gin.H{
				"user":      userDetail,
				"jwt_token": token,
			},
		})
	}
}

// jwt 토큰 가진 유저가 vote (변경될때마다 전송하는 식으로 작동하면 될 것 같은데 )
func VoteTimeHandler(rstrg *room.Storage, ustrg *user.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, Response{
				Error: "Unauthorized: missing user ID",
			})
			return
		}
		roomID, exists := c.Get("roomId")
		if !exists {
			c.JSON(http.StatusUnauthorized, Response{
				Error: "Unauthorized: missing room ID",
			})
			return
		}

		roomIDInt64 := int64(roomID.(float64)) // JWT 는 기본적으로 숫자 타입을 float에 매칭한다네..?
		userIdInt64 := int64(userID.(float64)) // JWT 는 기본적으로 숫자 타입을 float에 매칭한다네..?

		roomDates, err := rstrg.GetRoomDatesByRoomID(roomIDInt64)
		if err != nil {
			c.JSON(http.StatusNotFound, Response{
				Error: "Room not found",
			})
			return
		}

		var req user.ReqAvailableTimeList
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Error: "Invalid request format: " + err.Error(),
			})
			return
		}

		if err := ustrg.DeleteVoteTime(userIdInt64); err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Error: "Failed to update vote time: " + err.Error(),
			})
			return
		}

		// Insert new vote time
		for _, timeSlot := range req.Times {
			//check if date is valid
			if err := meeting.CheckValidDate(roomDates, timeSlot); err != nil {
				c.JSON(http.StatusBadRequest, Response{
					Error: err.Error(),
				})
				return
			}

			if err := ustrg.InsertVoteTime(userIdInt64, timeSlot); err != nil {
				c.JSON(http.StatusInternalServerError, Response{
					Error: "Failed to insert vote time: " + err.Error(),
				})
				return
			}
		}

		c.JSON(http.StatusOK, Response{
			Message: "Vote time updated successfully",
		})
	}
}

// User Detail
func GetUserDetailHandler(ustrg *user.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, Response{
				Error: "Unauthorized: missing user ID",
			})
			return
		}

		userIdInt64 := int64(userID.(float64))

		userDetail, err := ustrg.UserDetailById(userIdInt64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Error: "Failed to get user details: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, Response{
			Message: "Success",
			Data: gin.H{
				"user": userDetail,
			},
		})
	}
}

func GetResultHandler(rstrg *room.Storage, ustrg *user.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Param("url")
		room, err := rstrg.GetRoomByUrl(url)
		if err != nil {
			c.JSON(http.StatusNotFound, Response{
				Error: "Room not found",
			})
			return
		}

		roomDates, err := rstrg.GetRoomDatesByRoomID(room.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, Response{
				Error: "Failed to get room dates",
			})
			return
		}

		// TODO: Implement result calculation logic
		c.JSON(http.StatusOK, Response{
			Message: "Success",
			Data: gin.H{
				"room":  room,
				"dates": roomDates,
			},
		})
	}
}

func GetNoticeHandler(nstrg *notice.Storage, rstrg *room.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Param("url")
		room, err := rstrg.GetRoomByUrl(url)
		if err != nil {
			c.JSON(http.StatusNotFound, Response{
				Error: "Room not found",
			})
			return
		}

		res, err := nstrg.GetByRoomID(int(room.ID))

		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Error: "Failed to create notice: " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, Response{
			Message: "Success",
			Data:    res,
		})

	}
}

func CreateNoticeHandler(nstrg *notice.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, Response{
				Error: "Unauthorized: missing user ID",
			})
			return
		}
		roomID, exists := c.Get("roomId")
		if !exists {
			c.JSON(http.StatusUnauthorized, Response{
				Error: "Unauthorized: missing room ID",
			})
			return
		}

		var req notice.CreateNoticeReq
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Error: "Invalid request format: " + err.Error(),
			})
			return
		}
		if err := nstrg.Insert(req, int(roomID.(float64)), int(userID.(float64))); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Error: "Failed to create notice: " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, Response{
			Message: "Success",
		})

	}
}

// 리팩토링할 때 기본 인증 로직은 그냥 하나로 합치는게 나을 듯 너무 맣ㄴ다. delete 기능은 일단 추후에 구현 예정.
// func DeleteNoticeHandler(nstrg *notice.Storage) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userID, exists := c.Get("userId")
// 		if !exists {
// 			c.JSON(http.StatusUnauthorized, Response{
// 				Error: "Unauthorized: missing user ID",
// 			})
// 			return
// 		}

// 		noticeID := c.Param("id")
// 		if noticeID == "" {
// 			c.JSON(http.StatusBadRequest, Response{
// 				Error: "Notice ID is required",
// 			})
// 			return
// 		}

// 		noticeIDInt, err := strconv.ParseInt(noticeID, 10, 64)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, Response{
// 				Error: "Invalid notice ID format",
// 			})
// 			return
// 		}

// 		notices, err := nstrg.GetByRoomID(int(userID.(float64)))
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, Response{
// 				Error: "Failed to get notice: " + err.Error(),
// 			})
// 			return
// 		}

// 		found := false
// 		for _, n := range notices {
// 			if n.ID == noticeIDInt && n.USERID == int64(userID.(float64)) {
// 				found = true
// 				break
// 			}
// 		}

// 		if !found {
// 			c.JSON(http.StatusForbidden, Response{
// 				Error: "Notice not found or unauthorized",
// 			})
// 			return
// 		}

// 		if err := nstrg.DeleteById(int(noticeIDInt)); err != nil {
// 			c.JSON(http.StatusInternalServerError, Response{
// 				Error: "Failed to delete notice: " + err.Error(),
// 			})
// 			return
// 		}

// 		c.JSON(http.StatusOK, Response{
// 			Message: "Notice deleted successfully",
// 		})
// 	}
// }
