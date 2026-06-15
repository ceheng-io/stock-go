package types

import "time"

// Kline is a normalized OHLCV bar.
type Kline struct {
	Code      string
	Market    Market
	Time      time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	Amount    float64
	Turnover  float64
	Amplitude float64
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
	Date          string
	Timestamp     *int64
	TZ            string
	Code          string
	Open          *float64
	Close         *float64
	High          *float64
	Low           *float64
	Volume        *float64
	Amount        *float64
	Amplitude     *float64
	ChangePercent *float64
	Change        *float64
	TurnoverRate  *float64
}

func (HistoryKline) isAnyHistoryKline() {}

// ForeignHistoryKline is the common HK/US historical daily/weekly/monthly K-line row.
type ForeignHistoryKline struct {
	Date          string
	Timestamp     *int64
	TZ            string
	Code          string
	Name          string
	Open          *float64
	Close         *float64
	High          *float64
	Low           *float64
	Volume        *float64
	Amount        *float64
	Amplitude     *float64
	ChangePercent *float64
	Change        *float64
	TurnoverRate  *float64
}

// HKHistoryKline is a Hong Kong historical daily/weekly/monthly K-line row.
type HKHistoryKline struct {
	ForeignHistoryKline
	Currency string
	LotSize  *float64
}

func (HKHistoryKline) isAnyHistoryKline() {}

// USHistoryKline is a US historical daily/weekly/monthly K-line row.
type USHistoryKline struct {
	ForeignHistoryKline
	Currency string
}

func (USHistoryKline) isAnyHistoryKline() {}

// ForeignMinuteTimeline is the common HK/US 1-minute intraday timeline row.
type ForeignMinuteTimeline struct {
	Time      string
	Timestamp *int64
	TZ        string
	Code      string
	Open      *float64
	Close     *float64
	High      *float64
	Low       *float64
	Volume    *float64
	Amount    *float64
	AvgPrice  *float64
}

// ForeignMinuteKline is the common HK/US 5/15/30/60-minute K-line row.
type ForeignMinuteKline struct {
	Time          string
	Timestamp     *int64
	TZ            string
	Code          string
	Open          *float64
	Close         *float64
	High          *float64
	Low           *float64
	Volume        *float64
	Amount        *float64
	Amplitude     *float64
	ChangePercent *float64
	Change        *float64
	TurnoverRate  *float64
}

// HKMinuteTimeline is a Hong Kong 1-minute intraday timeline row.
type HKMinuteTimeline struct {
	ForeignMinuteTimeline
	Currency string
}

// HKMinuteKline is a Hong Kong 5/15/30/60-minute K-line row.
type HKMinuteKline struct {
	ForeignMinuteKline
	Currency string
}

// HKMinuteKlineResult contains either HK 1-minute timeline rows or minute K-line rows.
type HKMinuteKlineResult struct {
	Timeline []HKMinuteTimeline
	Klines   []HKMinuteKline
}

// USMinuteTimeline is a US 1-minute intraday timeline row.
type USMinuteTimeline struct {
	ForeignMinuteTimeline
	Currency string
}

// USMinuteKline is a US 5/15/30/60-minute K-line row.
type USMinuteKline struct {
	ForeignMinuteKline
	Currency string
}

// USMinuteKlineResult contains either US 1-minute timeline rows or minute K-line rows.
type USMinuteKlineResult struct {
	Timeline []USMinuteTimeline
	Klines   []USMinuteKline
}

// MinuteTimeline is a CN 1-minute intraday timeline row.
type MinuteTimeline struct {
	Time      string
	Timestamp *int64
	TZ        string
	Code      string
	Open      *float64
	Close     *float64
	High      *float64
	Low       *float64
	Volume    *float64
	Amount    *float64
	AvgPrice  *float64
}

// MinuteKline is a CN 5/15/30/60-minute K-line row.
type MinuteKline struct {
	Time          string
	Timestamp     *int64
	TZ            string
	Code          string
	Open          *float64
	Close         *float64
	High          *float64
	Low           *float64
	Volume        *float64
	Amount        *float64
	Amplitude     *float64
	ChangePercent *float64
	Change        *float64
	TurnoverRate  *float64
}

// MinuteKlineResult contains either 1-minute timeline rows or minute K-line rows.
type MinuteKlineResult struct {
	Timeline []MinuteTimeline
	Klines   []MinuteKline
}
