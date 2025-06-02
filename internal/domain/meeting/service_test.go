package meeting

import (
	"better-when2meet/internal/domain/room"
	"better-when2meet/internal/domain/user"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToVoteTable(t *testing.T) {
	// Mock room data
	mockRoom := room.RoomDetail{
		Room: room.Room{
			ID:         1,
			Name:       "Test Room",
			StartTime:  9,
			EndTime:    17,
			TimeRegion: "Asia/Seoul",
		},
		Dates: []room.RoomDate{
			{Year: 2024, Month: 3, Day: 1},
			{Year: 2024, Month: 3, Day: 2},
		},
	}

	// Mock user data
	mockUsers := []user.UserDetail{
		{
			User: user.ResUser{
				ID:         1,
				Name:       "Alice",
				TimeRegion: "Asia/Seoul",
			},
			AvailableTime: []user.AvailableTime{
				{
					Date:          time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					HourStartSlot: 9,
					HourEndSlot:   12,
				},
			},
		},
		{
			User: user.ResUser{
				ID:         2,
				Name:       "Bob",
				TimeRegion: "Asia/Seoul",
			},
			AvailableTime: []user.AvailableTime{
				{
					Date:          time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					HourStartSlot: 10,
					HourEndSlot:   15,
				},
				{
					Date:          time.Date(2024, 3, 2, 0, 0, 0, 0, time.UTC),
					HourStartSlot: 13,
					HourEndSlot:   17,
				},
			},
		},
	}

	result, err := ToVoteTable(mockUsers, mockRoom)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	t.Log("\nVote Table Result:")
	for date, blocks := range result {
		t.Logf("\nDate: %s", date)
		t.Log("Hour | Available Users")
		t.Log("-------------------")
		for _, block := range blocks {
			t.Logf("%2d:00 | %v", block.Hour, block.Users)
		}
	}

	date1 := "2024-03-01"
	assert.Contains(t, result, date1)
	blocks1 := result[date1]
	assert.Len(t, blocks1, 9)

	for i := 0; i < 4; i++ {
		assert.Contains(t, blocks1[i].Users, "Alice")
	}
	for i := 4; i < len(blocks1); i++ { // 13-17
		assert.NotContains(t, blocks1[i].Users, "Alice")
	}

	// Check Bob's availability (10-15 on day 1)
	for i := 1; i < 7; i++ { // 10-15
		assert.Contains(t, blocks1[i].Users, "Bob")
	}
	for i := 7; i < len(blocks1); i++ { // 16-17
		assert.NotContains(t, blocks1[i].Users, "Bob")
	}

	date2 := "2024-03-02"
	assert.Contains(t, result, date2)
	blocks2 := result[date2]
	assert.Len(t, blocks2, 9)

	for i := 0; i < 4; i++ {
		assert.NotContains(t, blocks2[i].Users, "Bob")
	}
	for i := 4; i < len(blocks2); i++ {
		assert.Contains(t, blocks2[i].Users, "Bob")
	}
}
