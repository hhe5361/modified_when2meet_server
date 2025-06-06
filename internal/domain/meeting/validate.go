package meeting

import (
	"better-when2meet/internal/domain/room"
	"better-when2meet/internal/domain/user"
	"errors"
	"time"
)

func CheckValidDate(dates []room.RoomDate, req user.ReqAvailableTime) error {
	if req.HourEndSlot < req.HourStartSlot {
		return errors.New("hour end slot must be greater than hour start slot")
	}

	date, _ := time.Parse("2006-01-02", req.Date)
	utcDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	for _, v := range dates {
		roomDate := time.Date(v.Year, time.Month(v.Month), v.Day, 0, 0, 0, 0, time.UTC)

		if utcDate.Equal(roomDate) {
			return nil
		}
	}
	return errors.New("date is not valid")
}
