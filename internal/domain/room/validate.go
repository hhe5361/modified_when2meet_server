package room

import "errors"

func CheckFieldType(r ReqCreateRoom) (res bool, err error) {
	if !checkFieldTime(r.StartTime, r.EndTime) {
		return false, errors.New("start and end time format is not accepted")
	}
	if !checkFieldTimeRegions(r.TimeRegion) {
		return false, errors.New("time region is not existed")
	}
	return true, nil
}

func checkFieldTime(start int, end int) bool {
	if start > end {
		return false
	}
	if start < 0 || start > 23 || end < 0 || end > 23 {
		return false
	}
	return true
}

func checkFieldTimeRegions(r string) bool {
	_, ok := allowedTimeRegions[r]
	return ok
}
