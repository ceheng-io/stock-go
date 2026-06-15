package types

// BlockTradeDateOptions configures block-trade date range queries.
type BlockTradeDateOptions struct {
	StartDate string
	EndDate   string
}

// BlockTradeMarketStatItem is a block-trade market statistics row.
type BlockTradeMarketStatItem struct {
	Date            string
	SHClose         *float64
	SHChangePercent *float64
	TotalAmount     *float64
	PremiumAmount   *float64
	PremiumRatio    *float64
	DiscountAmount  *float64
	DiscountRatio   *float64
}

// BlockTradeDetailItem is a block-trade deal detail row.
type BlockTradeDetailItem struct {
	Code          string
	Name          string
	Date          string
	Close         *float64
	ChangePercent *float64
	DealPrice     *float64
	DealVolume    *float64
	DealAmount    *float64
	PremiumRate   *float64
	BuyBranch     string
	SellBranch    string
}

// BlockTradeDailyStatItem is a block-trade daily stock statistics row.
type BlockTradeDailyStatItem struct {
	Code            string
	Name            string
	Date            string
	ChangePercent   *float64
	Close           *float64
	DealCount       *float64
	DealTotalAmount *float64
	DealTotalVolume *float64
	PremiumAmount   *float64
	DiscountAmount  *float64
}
