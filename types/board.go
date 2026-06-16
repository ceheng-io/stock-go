package types

// Board describes an Eastmoney industry or concept board list item.
type Board struct {
	Rank                      int      `json:"rank"`
	Name                      string   `json:"name"`
	Code                      string   `json:"code"`
	Price                     *float64 `json:"price"`
	Change                    *float64 `json:"change"`
	ChangePercent             *float64 `json:"changePercent"`
	TotalMarketCap            *float64 `json:"totalMarketCap"`
	TurnoverRate              *float64 `json:"turnoverRate"`
	RiseCount                 *float64 `json:"riseCount"`
	FallCount                 *float64 `json:"fallCount"`
	LeadingStock              *string  `json:"leadingStock"`
	LeadingStockChangePercent *float64 `json:"leadingStockChangePercent"`
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
	Item  string   `json:"item"`
	Value *float64 `json:"value"`
}

// BoardConstituent is an Eastmoney board constituent stock.
type BoardConstituent struct {
	Rank          int      `json:"rank"`
	Code          string   `json:"code"`
	Name          string   `json:"name"`
	Price         *float64 `json:"price"`
	ChangePercent *float64 `json:"changePercent"`
	Change        *float64 `json:"change"`
	Volume        *float64 `json:"volume"`
	Amount        *float64 `json:"amount"`
	Amplitude     *float64 `json:"amplitude"`
	High          *float64 `json:"high"`
	Low           *float64 `json:"low"`
	Open          *float64 `json:"open"`
	PrevClose     *float64 `json:"prevClose"`
	TurnoverRate  *float64 `json:"turnoverRate"`
	PE            *float64 `json:"pe"`
	PB            *float64 `json:"pb"`
}

// BoardKline is an Eastmoney board historical daily/weekly/monthly K-line row.
type BoardKline struct {
	Date          string   `json:"date"`
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

// BoardMinuteTimeline is an Eastmoney board 1-minute intraday timeline row.
type BoardMinuteTimeline struct {
	Time   string   `json:"time"`
	Open   *float64 `json:"open"`
	Close  *float64 `json:"close"`
	High   *float64 `json:"high"`
	Low    *float64 `json:"low"`
	Volume *float64 `json:"volume"`
	Amount *float64 `json:"amount"`
	Price  *float64 `json:"price"`
}

// BoardMinuteKline is an Eastmoney board 5/15/30/60-minute K-line row.
type BoardMinuteKline struct {
	Time          string   `json:"time"`
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

// BoardMinuteKlineResult contains either board 1-minute timeline rows or minute K-line rows.
type BoardMinuteKlineResult struct {
	Timeline []BoardMinuteTimeline `json:"timeline"`
	Klines   []BoardMinuteKline    `json:"klines"`
}
