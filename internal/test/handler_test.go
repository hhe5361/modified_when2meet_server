package test

import (
	"better-when2meet/internal/db"
	"better-when2meet/internal/room"
	"better-when2meet/internal/server"
	"better-when2meet/internal/user"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type Response = server.Response

func setupTestRouter() (*gin.Engine, *room.Storage, *user.Storage) {
	database := db.InitDB()
	roomRepo := room.New(database)
	userRepo := user.New(database)

	gin.SetMode(gin.TestMode)
	r := server.SetupRouter()

	return r, roomRepo, userRepo
}

func TestCreateRoomHandler(t *testing.T) {
	router, _, _ := setupTestRouter()

	tests := []struct {
		name       string
		payload    room.ReqCreateRoom
		wantStatus int
		wantError  bool
	}{
		{
			name: "Valid Room Creation",
			payload: room.ReqCreateRoom{
				RoomName:   "Test Room",
				TimeRegion: "Asia/Seoul",
				StartTime:  9,
				EndTime:    18,
				IsOnline:   false,
				VoteableRooms: []room.ReqRoomDate{
					{
						Year:  2024,
						Month: 3,
						Day:   1,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name: "Invalid Time Range",
			payload: room.ReqCreateRoom{
				RoomName:   "Test Room",
				TimeRegion: "Asia/Seoul",
				StartTime:  18,
				EndTime:    9,
				IsOnline:   false,
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/rooms", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantError {
				var response Response
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.Error)
			} else {
				var response Response
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Room created successfully", response.Message)
			}
		})
	}
}

func TestRegisterHandler(t *testing.T) {
	router, _, _ := setupTestRouter()

	roomPayload := room.ReqCreateRoom{
		RoomName:   "Test Room",
		TimeRegion: "Asia/Seoul",
		StartTime:  9,
		EndTime:    18,
		IsOnline:   false,
		VoteableRooms: []room.ReqRoomDate{
			{
				Year:  2024,
				Month: 3,
				Day:   1,
			},
		},
	}
	jsonData, _ := json.Marshal(roomPayload)
	req := httptest.NewRequest("POST", "/rooms", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var roomResponse Response
	json.Unmarshal(w.Body.Bytes(), &roomResponse)
	roomURL := roomResponse.Data.(map[string]interface{})["url"].(string)

	tests := []struct {
		name       string
		url        string
		payload    user.ReqLogin
		wantStatus int
		wantError  bool
	}{
		{
			name: "Valid Registration",
			url:  roomURL,
			payload: user.ReqLogin{
				Name:       "testuser",
				Password:   "testpass",
				TimeRegion: "Asia/Seoul",
			},
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name: "Invalid Room URL",
			url:  "nonexistent-room",
			payload: user.ReqLogin{
				Name:       "testuser",
				Password:   "testpass",
				TimeRegion: "Asia/Seoul",
			},
			wantStatus: http.StatusNotFound,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/rooms/"+tt.url+"/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			t.Logf("Response Status: %d", w.Code)
			t.Logf("Response Body: %s", w.Body.String())

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantError {
				var response Response
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.Error)
				t.Logf("Error Message: %s", response.Error)
			} else {
				var response Response
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Success", response.Message)
				assert.NotEmpty(t, response.Data.(map[string]interface{})["jwt_token"])
			}
		})
	}
}

func TestVoteTimeHandler(t *testing.T) {
	router, _, _ := setupTestRouter()

	roomPayload := room.ReqCreateRoom{
		RoomName:   "Test Room",
		TimeRegion: "Asia/Seoul",
		StartTime:  9,
		EndTime:    18,
		IsOnline:   false,
		VoteableRooms: []room.ReqRoomDate{
			{
				Year:  2024,
				Month: 3,
				Day:   1,
			},
		},
	}
	jsonData, _ := json.Marshal(roomPayload)
	req := httptest.NewRequest("POST", "/rooms", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var roomResponse Response
	json.Unmarshal(w.Body.Bytes(), &roomResponse)
	roomURL := roomResponse.Data.(map[string]interface{})["url"].(string)

	// Register user
	loginPayload := user.ReqLogin{
		Name:       "testuser",
		Password:   "testpass",
		TimeRegion: "Asia/Seoul",
	}
	jsonData, _ = json.Marshal(loginPayload)
	req = httptest.NewRequest("POST", "/rooms/"+roomURL+"/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var loginResponse Response
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse.Data.(map[string]interface{})["jwt_token"].(string)

	tests := []struct {
		name       string
		token      string
		payload    user.ReqAvailableTime
		wantStatus int
		wantError  bool
	}{
		{
			name:  "Valid Vote Time",
			token: token,
			payload: user.ReqAvailableTime{
				Date:          time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
				HourStartSlot: 9,
				HourEndSlot:   12,
			},
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name:  "Invalid Time Slot",
			token: token,
			payload: user.ReqAvailableTime{
				Date:          time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
				HourStartSlot: 12,
				HourEndSlot:   9,
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.payload)

			req := httptest.NewRequest("PUT", "/rooms/"+roomURL+"/times", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tt.token)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantError {
				var response Response
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.Error)
			} else {
				var response Response
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Vote time recorded successfully", response.Message)
			}
		})
	}
}
