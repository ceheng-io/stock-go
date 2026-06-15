package types

// Board describes an Eastmoney industry or concept board list item.
type Board struct {
	Rank                      int
	Name                      string
	Code                      string
	Price                     *float64
	Change                    *float64
	ChangePercent             *float64
	TotalMarketCap            *float64
	TurnoverRate              *float64
	RiseCount                 *float64
	FallCount                 *float64
	LeadingStock              *string
	LeadingStockChangePercent *float64
}

type IndustryBoard = Board
type ConceptBoard = Board
type IndustryBoardSpot = BoardSpot
type ConceptBoardSpot = BoardSpot
type IndustryBoardConstituent = BoardConstituent
type ConceptBoardConstituent = BoardConstituent
type IndustryBoardKline = BoardKline
type ConceptBoardKline = BoardKline
type IndustryBoardMinuteTimeline = BoardMinuteTimeline
type ConceptBoardMinuteTimeline = BoardMinuteTimeline
type IndustryBoardMinuteKline = BoardMinuteKline
type ConceptBoardMinuteKline = BoardMinuteKline

// BoardSpot is an Eastmoney board spot metric.
type BoardSpot struct {
	Item  string
	Value *float64
}

// BoardConstituent is an Eastmoney board constituent stock.
type BoardConstituent struct {
	Rank          int
	Code          string
	Name          string
	Price         *float64
	ChangePercent *float64
	Change        *float64
	Volume        *float64
	Amount        *float64
	Amplitude     *float64
	High          *float64
	Low           *float64
	Open          *float64
	PrevClose     *float64
	TurnoverRate  *float64
	PE            *float64
	PB            *float64
}

// BoardKline is an Eastmoney board historical daily/weekly/monthly K-line row.
type BoardKline struct {
	Date          string
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

// BoardMinuteTimeline is an Eastmoney board 1-minute intraday timeline row.
type BoardMinuteTimeline struct {
	Time   string
	Open   *float64
	Close  *float64
	High   *float64
	Low    *float64
	Volume *float64
	Amount *float64
	Price  *float64
}

// BoardMinuteKline is an Eastmoney board 5/15/30/60-minute K-line row.
type BoardMinuteKline struct {
	Time          string
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

// BoardMinuteKlineResult contains either board 1-minute timeline rows or minute K-line rows.
type BoardMinuteKlineResult struct {
	Timeline []BoardMinuteTimeline
	Klines   []BoardMinuteKline
}
