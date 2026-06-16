package types

// BlockTradeDateOptions configures block-trade date range queries.
type BlockTradeDateOptions struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// BlockTradeMarketStatItem is a block-trade market statistics row.
type BlockTradeMarketStatItem struct {
	Date            string   `json:"date"`
	SHClose         *float64 `json:"shClose"`
	SHChangePercent *float64 `json:"shChangePercent"`
	TotalAmount     *float64 `json:"totalAmount"`
	PremiumAmount   *float64 `json:"premiumAmount"`
	PremiumRatio    *float64 `json:"premiumRatio"`
	DiscountAmount  *float64 `json:"discountAmount"`
	DiscountRatio   *float64 `json:"discountRatio"`
}

// BlockTradeDetailItem is a block-trade deal detail row.
type BlockTradeDetailItem struct {
	Code          string   `json:"code"`
	Name          string   `json:"name"`
	Date          string   `json:"date"`
	Close         *float64 `json:"close"`
	ChangePercent *float64 `json:"changePercent"`
	DealPrice     *float64 `json:"dealPrice"`
	DealVolume    *float64 `json:"dealVolume"`
	DealAmount    *float64 `json:"dealAmount"`
	PremiumRate   *float64 `json:"premiumRate"`
	BuyBranch     string   `json:"buyBranch"`
	SellBranch    string   `json:"sellBranch"`
}

// BlockTradeDailyStatItem is a block-trade daily stock statistics row.
type BlockTradeDailyStatItem struct {
	Code            string   `json:"code"`
	Name            string   `json:"name"`
	Date            string   `json:"date"`
	ChangePercent   *float64 `json:"changePercent"`
	Close           *float64 `json:"close"`
	DealCount       *float64 `json:"dealCount"`
	DealTotalAmount *float64 `json:"dealTotalAmount"`
	DealTotalVolume *float64 `json:"dealTotalVolume"`
	PremiumAmount   *float64 `json:"premiumAmount"`
	DiscountAmount  *float64 `json:"discountAmount"`
}
