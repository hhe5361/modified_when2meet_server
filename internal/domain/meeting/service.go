package meeting //this packag

import (
	"better-when2meet/internal/domain/room"
	"better-when2meet/internal/domain/user"
	"fmt"
)

// userDetail 에 빈 슬라이스 주거나 그렇게 해야할 듯?
// 아직까지 storage 와 service 의 구분이 모호한 듯. 리팩토링할 때 개선 해야함. 두 개 따로 들고 오는게 그냥 좋을 것 같음.
func ToVoteTable(userDetail []user.UserDetail, roomDetail room.RoomDetail) (VoteTable, error) {
	//Make Votable data type\
	//room Deatil 에는 투표할 수 인는 date 가 저장되어 있으며
	result := make(VoteTable)

	//make maps
	//"%04d-%02d-%02d", year, month, day
	for _, date := range roomDetail.Dates {
		dateStr := fmt.Sprintf("%04d-%02d-%02d", date.Year, date.Month, date.Day)
		result[dateStr] = makeHourBlock(roomDetail.Room.StartTime, roomDetail.Room.EndTime)
	}

	for _, detail := range userDetail {
		for _, availableTime := range detail.AvailableTime {
			dateStr := fmt.Sprintf("%04d-%02d-%02d",
				availableTime.Date.Year(),
				int(availableTime.Date.Month()),
				availableTime.Date.Day())

			if blocks, exists := result[dateStr]; exists {
				for i := range blocks {
					if blocks[i].Hour >= availableTime.HourStartSlot && blocks[i].Hour <= availableTime.HourEndSlot {
						blocks[i].Users = append(blocks[i].Users, detail.User.Name)
					}
				}
				result[dateStr] = blocks
			}
		}
	}
	return result, nil
}

func makeHourBlock(start int, end int) []HourBlock {
	var result []HourBlock

	for i := start; i <= end; i++ {
		value := HourBlock{
			Hour:  i,
			Users: []string{},
		}
		result = append(result, value)
	}
	return result
}
