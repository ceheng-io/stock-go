package types

// TodayTimeline 是腾讯当日分时项。
type TodayTimeline struct {
	Time      string  `json:"time"`
	Timestamp *int64  `json:"timestamp"`
	TZ        string  `json:"tz"`
	Price     float64 `json:"price"`
	AvgPrice  float64 `json:"avgPrice"`
	Volume    float64 `json:"volume"`
	Amount    float64 `json:"amount"`
}

// TodayTimelineResponse 是腾讯当日分时响应。
type TodayTimelineResponse struct {
	Code      string          `json:"code"`
	Date      string          `json:"date"`
	Timestamp *int64          `json:"timestamp"`
	TZ        string          `json:"tz"`
	PreClose  *float64        `json:"preClose"`
	Data      []TodayTimeline `json:"data"`
}
