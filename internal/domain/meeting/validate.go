package meeting

import (
	"better-when2meet/internal/domain/room"
	"better-when2meet/internal/domain/user"
	"errors"
)

func CheckValidDate(dates []room.RoomDate, req user.ReqAvailableTime) error {

	if req.HourEndSlot < req.HourStartSlot {
		return errors.New("hour end slot must be greater than hour start slot")
	}

	for _, v := range dates {
		if req.Date.Year() == v.Year && int(req.Date.Month()) == v.Month && req.Date.Day() == v.Day {
			return nil
		}
	}
	return errors.New("date is not valid")
}
