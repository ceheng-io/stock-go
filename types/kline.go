package types

// Kline is a normalized OHLCV bar.
type Kline struct {
	Code      string  `json:"code"`
	Market    Market  `json:"market"`
	Time      int64   `json:"time"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
	Amount    float64 `json:"amount"`
	Turnover  float64 `json:"turnover"`
	Amplitude float64 `json:"amplitude"`
}

// AnyHistoryKline is the common interface implemented by CN/HK/US history rows.
//
// It mirrors the TypeScript union:
// HistoryKline | HKHistoryKline | USHistoryKline.
type AnyHistoryKline interface {
	isAnyHistoryKline()
}

// HistoryKline is a CN historical daily/weekly/monthly K-line row.
type HistoryKline struct {
	Date          string   `json:"date"`
	Timestamp     *int64   `json:"timestamp"`
	TZ            string   `json:"tz"`
	Code          string   `json:"code"`
	Open          *float64 `json:"open"`
	Close         *float64 `json:"close"`
	High          *float64 `json:"high"`
	Low           *float64 `json:"low"`
	Volume        *float64 `json:"volume"`
	Amount        *float64 `json:"amount"`
	Amplitude     *float64 `json:"amplitude"`
	ChangePercent *float64 `json:"changePercent"`
	Change        *float64 `json:"change"`
	TurnoverRate  *float64 `json:"turnoverRate"`
}

func (HistoryKline) isAnyHistoryKline() {}

// ForeignHistoryKline is the common HK/US historical daily/weekly/monthly K-line row.
type ForeignHistoryKline struct {
	Date          string   `json:"date"`
	Timestamp     *int64   `json:"timestamp"`
	TZ            string   `json:"tz"`
	Code          string   `json:"code"`
	Name          string   `json:"name"`
	Open          *float64 `json:"open"`
	Close         *float64 `json:"close"`
	High          *float64 `json:"high"`
	Low           *float64 `json:"low"`
	Volume        *float64 `json:"volume"`
	Amount        *float64 `json:"amount"`
	Amplitude     *float64 `json:"amplitude"`
	ChangePercent *float64 `json:"changePercent"`
	Change        *float64 `json:"change"`
	TurnoverRate  *float64 `json:"turnoverRate"`
}

// HKHistoryKline is a Hong Kong historical daily/weekly/monthly K-line row.
type HKHistoryKline struct {
	ForeignHistoryKline
	Currency string   `json:"currency"`
	LotSize  *float64 `json:"lotSize"`
}

func (HKHistoryKline) isAnyHistoryKline() {}

// USHistoryKline is a US historical daily/weekly/monthly K-line row.
type USHistoryKline struct {
	ForeignHistoryKline
	Currency string `json:"currency"`
}

func (USHistoryKline) isAnyHistoryKline() {}

// ForeignMinuteTimeline is the common HK/US 1-minute intraday timeline row.
type ForeignMinuteTimeline struct {
	Time      string   `json:"time"`
	Timestamp *int64   `json:"timestamp"`
	TZ        string   `json:"tz"`
	Code      string   `json:"code"`
	Open      *float64 `json:"open"`
	Close     *float64 `json:"close"`
	High      *float64 `json:"high"`
	Low       *float64 `json:"low"`
	Volume    *float64 `json:"volume"`
	Amount    *float64 `json:"amount"`
	AvgPrice  *float64 `json:"avgPrice"`
}

// ForeignMinuteKline is the common HK/US 5/15/30/60-minute K-line row.
type ForeignMinuteKline struct {
	Time          string   `json:"time"`
	Timestamp     *int64   `json:"timestamp"`
	TZ            string   `json:"tz"`
	Code          string   `json:"code"`
	Open          *float64 `json:"open"`
	Close         *float64 `json:"close"`
	High          *float64 `json:"high"`
	Low           *float64 `json:"low"`
	Volume        *float64 `json:"volume"`
	Amount        *float64 `json:"amount"`
	Amplitude     *float64 `json:"amplitude"`
	ChangePercent *float64 `json:"changePercent"`
	Change        *float64 `json:"change"`
	TurnoverRate  *float64 `json:"turnoverRate"`
}

// HKMinuteTimeline is a Hong Kong 1-minute intraday timeline row.
type HKMinuteTimeline struct {
	ForeignMinuteTimeline
	Currency string `json:"currency"`
}

// HKMinuteKline is a Hong Kong 5/15/30/60-minute K-line row.
type HKMinuteKline struct {
	ForeignMinuteKline
	Currency string `json:"currency"`
}

// HKMinuteKlineResult contains either HK 1-minute timeline rows or minute K-line rows.
type HKMinuteKlineResult struct {
	Timeline []HKMinuteTimeline `json:"timeline"`
	Klines   []HKMinuteKline    `json:"klines"`
}

// USMinuteTimeline is a US 1-minute intraday timeline row.
type USMinuteTimeline struct {
	ForeignMinuteTimeline
	Currency string `json:"currency"`
}

// USMinuteKline is a US 5/15/30/60-minute K-line row.
type USMinuteKline struct {
	ForeignMinuteKline
	Currency string `json:"currency"`
}

// USMinuteKlineResult contains either US 1-minute timeline rows or minute K-line rows.
type USMinuteKlineResult struct {
	Timeline []USMinuteTimeline `json:"timeline"`
	Klines   []USMinuteKline    `json:"klines"`
}

// MinuteTimeline is a CN 1-minute intraday timeline row.
type MinuteTimeline struct {
	Time      string   `json:"time"`
	Timestamp *int64   `json:"timestamp"`
	TZ        string   `json:"tz"`
	Code      string   `json:"code"`
	Open      *float64 `json:"open"`
	Close     *float64 `json:"close"`
	High      *float64 `json:"high"`
	Low       *float64 `json:"low"`
	Volume    *float64 `json:"volume"`
	Amount    *float64 `json:"amount"`
	AvgPrice  *float64 `json:"avgPrice"`
}

// MinuteKline is a CN 5/15/30/60-minute K-line row.
type MinuteKline struct {
	Time          string   `json:"time"`
	Timestamp     *int64   `json:"timestamp"`
	TZ            string   `json:"tz"`
	Code          string   `json:"code"`
	Open          *float64 `json:"open"`
	Close         *float64 `json:"close"`
	High          *float64 `json:"high"`
	Low           *float64 `json:"low"`
	Volume        *float64 `json:"volume"`
	Amount        *float64 `json:"amount"`
	Amplitude     *float64 `json:"amplitude"`
	ChangePercent *float64 `json:"changePercent"`
	Change        *float64 `json:"change"`
	TurnoverRate  *float64 `json:"turnoverRate"`
}

// MinuteKlineResult contains either 1-minute timeline rows or minute K-line rows.
type MinuteKlineResult struct {
	Timeline []MinuteTimeline `json:"timeline"`
	Klines   []MinuteKline    `json:"klines"`
}
