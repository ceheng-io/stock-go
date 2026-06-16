package types

// ZTPoolType identifies an Eastmoney limit-up pool category.
type ZTPoolType string

const (
	ZTPoolZT        ZTPoolType = "zt"
	ZTPoolYesterday ZTPoolType = "yesterday"
	ZTPoolStrong    ZTPoolType = "strong"
	ZTPoolSubNew    ZTPoolType = "sub_new"
	ZTPoolBroken    ZTPoolType = "broken"
	ZTPoolDT        ZTPoolType = "dt"
)

// ZTPoolItem is an Eastmoney limit-up pool row.
type ZTPoolItem struct {
	Code                 string   `json:"code"`
	Name                 string   `json:"name"`
	Price                *float64 `json:"price"`
	ChangePercent        *float64 `json:"changePercent"`
	LimitPrice           *float64 `json:"limitPrice"`
	Amount               *float64 `json:"amount"`
	FloatMarketValue     *float64 `json:"floatMarketValue"`
	TotalMarketValue     *float64 `json:"totalMarketValue"`
	TurnoverRate         *float64 `json:"turnoverRate"`
	ContinuousBoardCount *float64 `json:"continuousBoardCount"`
	FirstBoardTime       *string  `json:"firstBoardTime"`
	LastBoardTime        *string  `json:"lastBoardTime"`
	BoardAmount          *float64 `json:"boardAmount"`
	SealAmount           *float64 `json:"sealAmount"`
	FailedCount          *float64 `json:"failedCount"`
	Industry             string   `json:"industry"`
	ZTStatistics         string   `json:"ztStatistics"`
	Amplitude            *float64 `json:"amplitude"`
	Speed                *float64 `json:"speed"`
}

// StockChangeType identifies an Eastmoney intraday stock change category.
type StockChangeType string

const (
	StockChangeRocketLaunch   StockChangeType = "rocket_launch"
	StockChangeQuickRebound   StockChangeType = "quick_rebound"
	StockChangeLargeBuy       StockChangeType = "large_buy"
	StockChangeLimitUpSeal    StockChangeType = "limit_up_seal"
	StockChangeLimitDownOpen  StockChangeType = "limit_down_open"
	StockChangeBigBuyOrder    StockChangeType = "big_buy_order"
	StockChangeAuctionUp      StockChangeType = "auction_up"
	StockChangeHighOpen5D     StockChangeType = "high_open_5d"
	StockChangeGapUp          StockChangeType = "gap_up"
	StockChangeHigh60D        StockChangeType = "high_60d"
	StockChangeSurge60D       StockChangeType = "surge_60d"
	StockChangeAccelerateDown StockChangeType = "accelerate_down"
	StockChangeHighDive       StockChangeType = "high_dive"
	StockChangeLargeSell      StockChangeType = "large_sell"
	StockChangeLimitDownSeal  StockChangeType = "limit_down_seal"
	StockChangeLimitUpOpen    StockChangeType = "limit_up_open"
	StockChangeBigSellOrder   StockChangeType = "big_sell_order"
	StockChangeAuctionDown    StockChangeType = "auction_down"
	StockChangeLowOpen5D      StockChangeType = "low_open_5d"
	StockChangeGapDown        StockChangeType = "gap_down"
	StockChangeLow60D         StockChangeType = "low_60d"
	StockChangeDrop60D        StockChangeType = "drop_60d"
)

// StockChangeItem is an intraday stock change row.
type StockChangeItem struct {
	Time            string          `json:"time"`
	Code            string          `json:"code"`
	Name            string          `json:"name"`
	ChangeType      StockChangeType `json:"changeType"`
	ChangeTypeLabel string          `json:"changeTypeLabel"`
	Info            string          `json:"info"`
}

// BoardChangeItem is a board change summary row.
type BoardChangeItem struct {
	Name                   string             `json:"name"`
	ChangePercent          *float64           `json:"changePercent"`
	MainNetInflow          *float64           `json:"mainNetInflow"`
	TotalChangeCount       *float64           `json:"totalChangeCount"`
	TopStockCode           string             `json:"topStockCode"`
	TopStockName           string             `json:"topStockName"`
	TopStockDirection      string             `json:"topStockDirection"`
	ChangeTypeDistribution map[string]float64 `json:"changeTypeDistribution"`
}

// THSLimitUpOrderField identifies a Tonghuashun limit-up pool sort field.
type THSLimitUpOrderField string

const (
	THSLimitUpOrderFirstLimitUpTime THSLimitUpOrderField = "330323"
	THSLimitUpOrderLastLimitUpTime  THSLimitUpOrderField = "330324"
	THSLimitUpOrderOpenNum          THSLimitUpOrderField = "330325"
)

// THSLimitUpOrderType identifies Tonghuashun sort direction.
type THSLimitUpOrderType string

const (
	THSLimitUpOrderDesc THSLimitUpOrderType = "0"
	THSLimitUpOrderAsc  THSLimitUpOrderType = "1"
)

// THSLimitUpPoolOptions controls Tonghuashun limit-up pool requests.
type THSLimitUpPoolOptions struct {
	Date       string               `json:"date"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	Filter     string               `json:"filter"`
	OrderField THSLimitUpOrderField `json:"orderField"`
	OrderType  THSLimitUpOrderType  `json:"orderType"`
}

// THSLimitUpPoolResult is a Tonghuashun limit-up pool page.
type THSLimitUpPoolResult struct {
	Page           THSLimitUpPage    `json:"page"`
	Items          []THSLimitUpItem  `json:"items"`
	LimitUpCount   THSLimitStatGroup `json:"limitUpCount"`
	LimitDownCount THSLimitStatGroup `json:"limitDownCount"`
	Date           string            `json:"date"`
	Message        string            `json:"message"`
	TradeStatus    THSTradeStatus    `json:"tradeStatus"`
}

// THSLimitUpPage describes Tonghuashun pagination metadata.
type THSLimitUpPage struct {
	Limit int `json:"limit"`
	Total int `json:"total"`
	Count int `json:"count"`
	Page  int `json:"page"`
}

// THSLimitStatGroup contains today's and yesterday's up/down limit statistics.
type THSLimitStatGroup struct {
	Today     THSLimitStat `json:"today"`
	Yesterday THSLimitStat `json:"yesterday"`
}

// THSLimitStat describes Tonghuashun limit-up/down statistics.
type THSLimitStat struct {
	Num        int      `json:"num"`
	HistoryNum int      `json:"historyNum"`
	Rate       *float64 `json:"rate"`
	OpenNum    int      `json:"openNum"`
}

// THSTradeStatus describes Tonghuashun market status metadata.
type THSTradeStatus struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// THSLimitUpItem is a Tonghuashun limit-up pool row.
type THSLimitUpItem struct {
	Code                 string    `json:"code"`
	Name                 string    `json:"name"`
	Latest               *float64  `json:"latest"`
	ChangeRate           *float64  `json:"changeRate"`
	FirstLimitUpTime     *int64    `json:"firstLimitUpTime"`
	FirstLimitUpTimeText string    `json:"firstLimitUpTimeText"`
	LastLimitUpTime      *int64    `json:"lastLimitUpTime"`
	LastLimitUpTimeText  string    `json:"lastLimitUpTimeText"`
	OpenNum              *int      `json:"openNum"`
	LimitUpType          string    `json:"limitUpType"`
	OrderVolume          *float64  `json:"orderVolume"`
	OrderAmount          *float64  `json:"orderAmount"`
	TurnoverRate         *float64  `json:"turnoverRate"`
	CurrencyValue        *float64  `json:"currencyValue"`
	ReasonType           string    `json:"reasonType"`
	HighDays             string    `json:"highDays"`
	HighDaysValue        *int      `json:"highDaysValue"`
	ChangeTag            string    `json:"changeTag"`
	MarketType           string    `json:"marketType"`
	MarketID             *int      `json:"marketId"`
	IsNew                *int      `json:"isNew"`
	IsAgainLimit         *int      `json:"isAgainLimit"`
	LimitUpSuccessRate   *float64  `json:"limitUpSuccessRate"`
	TimePreview          []float64 `json:"timePreview"`
}
