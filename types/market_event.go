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
	Code                 string
	Name                 string
	Price                *float64
	ChangePercent        *float64
	LimitPrice           *float64
	Amount               *float64
	FloatMarketValue     *float64
	TotalMarketValue     *float64
	TurnoverRate         *float64
	ContinuousBoardCount *float64
	FirstBoardTime       *string
	LastBoardTime        *string
	BoardAmount          *float64
	SealAmount           *float64
	FailedCount          *float64
	Industry             string
	ZTStatistics         string
	Amplitude            *float64
	Speed                *float64
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
	Time            string
	Code            string
	Name            string
	ChangeType      StockChangeType
	ChangeTypeLabel string
	Info            string
}

// BoardChangeItem is a board change summary row.
type BoardChangeItem struct {
	Name                   string
	ChangePercent          *float64
	MainNetInflow          *float64
	TotalChangeCount       *float64
	TopStockCode           string
	TopStockName           string
	TopStockDirection      string
	ChangeTypeDistribution map[string]float64
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
	Date       string
	Page       int
	Limit      int
	Filter     string
	OrderField THSLimitUpOrderField
	OrderType  THSLimitUpOrderType
}

// THSLimitUpPoolResult is a Tonghuashun limit-up pool page.
type THSLimitUpPoolResult struct {
	Page           THSLimitUpPage
	Items          []THSLimitUpItem
	LimitUpCount   THSLimitStatGroup
	LimitDownCount THSLimitStatGroup
	Date           string
	Message        string
	TradeStatus    THSTradeStatus
}

// THSLimitUpPage describes Tonghuashun pagination metadata.
type THSLimitUpPage struct {
	Limit int
	Total int
	Count int
	Page  int
}

// THSLimitStatGroup contains today's and yesterday's up/down limit statistics.
type THSLimitStatGroup struct {
	Today     THSLimitStat
	Yesterday THSLimitStat
}

// THSLimitStat describes Tonghuashun limit-up/down statistics.
type THSLimitStat struct {
	Num        int
	HistoryNum int
	Rate       *float64
	OpenNum    int
}

// THSTradeStatus describes Tonghuashun market status metadata.
type THSTradeStatus struct {
	ID        string
	Name      string
	StartTime string
	EndTime   string
}

// THSLimitUpItem is a Tonghuashun limit-up pool row.
type THSLimitUpItem struct {
	Code                 string
	Name                 string
	Latest               *float64
	ChangeRate           *float64
	FirstLimitUpTime     *int64
	FirstLimitUpTimeText string
	LastLimitUpTime      *int64
	LastLimitUpTimeText  string
	OpenNum              *int
	LimitUpType          string
	OrderVolume          *float64
	OrderAmount          *float64
	TurnoverRate         *float64
	CurrencyValue        *float64
	ReasonType           string
	HighDays             string
	HighDaysValue        *int
	ChangeTag            string
	MarketType           string
	MarketID             *int
	IsNew                *int
	IsAgainLimit         *int
	LimitUpSuccessRate   *float64
	TimePreview          []float64
}
