package enumtype

type TimeRegion string

const (
	TimeRegionUTC     TimeRegion = "UTC"
	TimeRegionSeoul   TimeRegion = "Asia/Seoul"
	TimeRegionLondon  TimeRegion = "Europe/London"
	TimeRegionNewYork TimeRegion = "America/New_York"
)
