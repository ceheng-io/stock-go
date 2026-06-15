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
