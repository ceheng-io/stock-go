package types

// DragonTigerPeriod is a dragon-tiger billboard statistics period.
type DragonTigerPeriod string

const (
	DragonTigerPeriodOneMonth   DragonTigerPeriod = "1month"
	DragonTigerPeriodThreeMonth DragonTigerPeriod = "3month"
	DragonTigerPeriodSixMonth   DragonTigerPeriod = "6month"
	DragonTigerPeriodOneYear    DragonTigerPeriod = "1year"
)

// DragonTigerDateOptions configures dragon-tiger date range queries.
type DragonTigerDateOptions struct {
	StartDate string
	EndDate   string
}

// DragonTigerDetailItem is a daily dragon-tiger billboard stock detail row.
type DragonTigerDetailItem struct {
	Code             string
	Name             string
	Date             string
	Close            *float64
	ChangePercent    *float64
	NetBuyAmount     *float64
	BuyAmount        *float64
	SellAmount       *float64
	DealAmount       *float64
	TotalAmount      *float64
	NetBuyRatio      *float64
	DealAmountRatio  *float64
	TurnoverRate     *float64
	FloatMarketValue *float64
	Reason           string
	AfterChange1D    *float64
	AfterChange2D    *float64
	AfterChange5D    *float64
	AfterChange10D   *float64
}

// DragonTigerStockStatItem is a stock dragon-tiger billboard statistics row.
type DragonTigerStockStatItem struct {
	Code            string
	Name            string
	LatestDate      string
	Close           *float64
	ChangePercent   *float64
	Count           *float64
	TotalBuyAmount  *float64
	TotalSellAmount *float64
	TotalNetAmount  *float64
	TotalDealAmount *float64
	BuyOrgCount     *float64
	SellOrgCount    *float64
}

// DragonTigerInstitutionItem is an institution trade statistics row.
type DragonTigerInstitutionItem struct {
	Code          string
	Name          string
	Date          string
	Close         *float64
	ChangePercent *float64
	BuyOrgCount   *float64
	SellOrgCount  *float64
	OrgBuyAmount  *float64
	OrgSellAmount *float64
	OrgNetAmount  *float64
}

// DragonTigerBranchItem is a brokerage branch ranking row.
type DragonTigerBranchItem struct {
	Code            string
	Name            string
	TotalBuyAmount  *float64
	TotalSellAmount *float64
	BuyCount        *float64
	SellCount       *float64
	TotalCount      *float64
}

// DragonTigerSeatItem is a stock seat detail row.
type DragonTigerSeatItem struct {
	Rank            *float64
	BranchName      string
	BuyAmount       *float64
	BuyAmountRatio  *float64
	SellAmount      *float64
	SellAmountRatio *float64
	NetAmount       *float64
	Side            string
}
