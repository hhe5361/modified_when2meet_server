package test

import (
	"better-when2meet/internal/db"
	"better-when2meet/internal/room"
	"better-when2meet/internal/server"
	"better-when2meet/internal/user"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateRoom(t *testing.T) {
	database := db.InitDB()
	if database == nil {
		t.Fatal("Failed to initialize test database")
	}
	defer database.Close()

	router := server.SetupRouter()
	if router == nil {
		t.Error("Router should not be nil")
	}

	w := httptest.NewRecorder()

	var voteableDates []room.ReqRoomDate
	now := time.Now()
	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, i)
		voteableDates = append(voteableDates, room.ReqRoomDate{
			Year:  date.Year(),
			Month: int(date.Month()),
			Day:   date.Day(),
		})
	}

	exampleRoom := room.ReqCreateRoom{
		RoomName:      "Test Meeting Room",
		TimeRegion:    "Seoul/Asia",
		StartTime:     10,
		EndTime:       18,
		IsOnline:      false,
		VoteableRooms: voteableDates,
	}

	jsonData, err := json.Marshal(exampleRoom)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/rooms", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	// Read and log the response body
	body, _ := io.ReadAll(w.Body)
	t.Logf("Response body: %s", string(body))

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetRoomInfo(t *testing.T) {
	database := db.InitDB()
	if database == nil {
		t.Fatal("Failed to initialize test database")
	}
	defer database.Close()

	roomStorage := room.New(database)
	userStorage := user.New(database)

	var voteableDates []room.ReqRoomDate
	now := time.Now()
	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, i)
		voteableDates = append(voteableDates, room.ReqRoomDate{
			Year:  date.Year(),
			Month: int(date.Month()),
			Day:   date.Day(),
		})
	}

	testRoom := room.ReqCreateRoom{
		RoomName:      "Test Meeting Room",
		TimeRegion:    "Asia/Seoul",
		StartTime:     10,
		EndTime:       18,
		IsOnline:      false,
		VoteableRooms: voteableDates,
	}

	// Log the room creation attempt
	t.Log("Attempting to create test room...")
	if err := roomStorage.InsertRoom(testRoom, "test-room-123"); err != nil {
		t.Fatal("Failed to create test room:", err)
	}
	t.Log("Test room created successfully")

	// Verify room was created
	createdRoom, err := roomStorage.GetRoomByUrl("test-room-123")
	if err != nil {
		t.Fatal("Failed to verify room creation:", err)
	}
	t.Logf("Verified room creation: %+v", createdRoom)

	// Create test user
	testUser := user.ReqLogin{
		Name:       "Test User",
		Password:   "testpass",
		TimeRegion: "Seoul/Asia",
	}

	// Log user creation attempt
	t.Log("Attempting to create test user...")
	userId, err := userStorage.InsertUser(testUser, createdRoom.ID)
	if err != nil {
		t.Fatal("Failed to create test user:", err)
	}
	t.Logf("Test user created successfully with ID: %d", userId)

	// Set up router
	router := server.SetupRouter()
	if router == nil {
		t.Error("Router should not be nil")
	}

	t.Run("Get Room Info Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rooms/test-room-123", nil)
		router.ServeHTTP(w, req)

		// Log response for debugging
		body, _ := io.ReadAll(w.Body)
		t.Logf("Response status: %d", w.Code)
		t.Logf("Response body: %s", string(body))

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			t.Fatal("Failed to parse response:", err)
		}

		// Verify response structure
		if response["message"] != "Success" {
			t.Errorf("Expected message 'Success', got '%v'", response["message"])
		}

		data, ok := response["data"].(map[string]interface{})
		if !ok {
			t.Fatal("Response data is not in expected format")
		}

		// Verify room info
		roomInfo, ok := data["roomInfo"].(map[string]interface{})
		if !ok {
			t.Fatal("Room info is not in expected format")
		}
		if roomInfo["name"] != "Test Meeting Room" {
			t.Errorf("Expected room name 'Test Meeting Room', got '%v'", roomInfo["name"])
		}
	})

	t.Run("Get Room Info Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rooms/non-existent-room", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}

		body, _ := io.ReadAll(w.Body)
		var response map[string]interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			t.Fatal("Failed to parse response:", err)
		}

		if response["error"] != "Room not found" {
			t.Errorf("Expected error 'Room not found', got '%v'", response["error"])
		}
	})
}
