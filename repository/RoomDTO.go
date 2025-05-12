package dto

import {
}

type createRoom struct{
    UserName    string     `json:"user_name"`    // 사용자 이름
    UserPassword string    `json:"user_password"`// 사용자 비밀번호
    TimeRegion  string     `json:"time_region"`  // 시간대
    RoomDates   []RoomDate `json:"room_dates"`   // 투표 가능한 날짜들 (슬라이스)
} 