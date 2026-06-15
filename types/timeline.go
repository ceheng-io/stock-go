package types

// TodayTimeline 是腾讯当日分时项。
type TodayTimeline struct {
	Time      string
	Timestamp *int64
	TZ        string
	Price     float64
	AvgPrice  float64
	Volume    float64
	Amount    float64
}

// TodayTimelineResponse 是腾讯当日分时响应。
type TodayTimelineResponse struct {
	Code      string
	Date      string
	Timestamp *int64
	TZ        string
	PreClose  *float64
	Data      []TodayTimeline
}
