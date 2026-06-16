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
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// DragonTigerDetailItem is a daily dragon-tiger billboard stock detail row.
type DragonTigerDetailItem struct {
	Code             string   `json:"code"`
	Name             string   `json:"name"`
	Date             string   `json:"date"`
	Close            *float64 `json:"close"`
	ChangePercent    *float64 `json:"changePercent"`
	NetBuyAmount     *float64 `json:"netBuyAmount"`
	BuyAmount        *float64 `json:"buyAmount"`
	SellAmount       *float64 `json:"sellAmount"`
	DealAmount       *float64 `json:"dealAmount"`
	TotalAmount      *float64 `json:"totalAmount"`
	NetBuyRatio      *float64 `json:"netBuyRatio"`
	DealAmountRatio  *float64 `json:"dealAmountRatio"`
	TurnoverRate     *float64 `json:"turnoverRate"`
	FloatMarketValue *float64 `json:"floatMarketValue"`
	Reason           string   `json:"reason"`
	AfterChange1D    *float64 `json:"afterChange1D"`
	AfterChange2D    *float64 `json:"afterChange2D"`
	AfterChange5D    *float64 `json:"afterChange5D"`
	AfterChange10D   *float64 `json:"afterChange10D"`
}

// DragonTigerStockStatItem is a stock dragon-tiger billboard statistics row.
type DragonTigerStockStatItem struct {
	Code            string   `json:"code"`
	Name            string   `json:"name"`
	LatestDate      string   `json:"latestDate"`
	Close           *float64 `json:"close"`
	ChangePercent   *float64 `json:"changePercent"`
	Count           *float64 `json:"count"`
	TotalBuyAmount  *float64 `json:"totalBuyAmount"`
	TotalSellAmount *float64 `json:"totalSellAmount"`
	TotalNetAmount  *float64 `json:"totalNetAmount"`
	TotalDealAmount *float64 `json:"totalDealAmount"`
	BuyOrgCount     *float64 `json:"buyOrgCount"`
	SellOrgCount    *float64 `json:"sellOrgCount"`
}

// DragonTigerInstitutionItem is an institution trade statistics row.
type DragonTigerInstitutionItem struct {
	Code          string   `json:"code"`
	Name          string   `json:"name"`
	Date          string   `json:"date"`
	Close         *float64 `json:"close"`
	ChangePercent *float64 `json:"changePercent"`
	BuyOrgCount   *float64 `json:"buyOrgCount"`
	SellOrgCount  *float64 `json:"sellOrgCount"`
	OrgBuyAmount  *float64 `json:"orgBuyAmount"`
	OrgSellAmount *float64 `json:"orgSellAmount"`
	OrgNetAmount  *float64 `json:"orgNetAmount"`
}

// DragonTigerBranchItem is a brokerage branch ranking row.
type DragonTigerBranchItem struct {
	Code            string   `json:"code"`
	Name            string   `json:"name"`
	TotalBuyAmount  *float64 `json:"totalBuyAmount"`
	TotalSellAmount *float64 `json:"totalSellAmount"`
	BuyCount        *float64 `json:"buyCount"`
	SellCount       *float64 `json:"sellCount"`
	TotalCount      *float64 `json:"totalCount"`
}

// DragonTigerSeatItem is a stock seat detail row.
type DragonTigerSeatItem struct {
	Rank            *float64 `json:"rank"`
	BranchName      string   `json:"branchName"`
	BuyAmount       *float64 `json:"buyAmount"`
	BuyAmountRatio  *float64 `json:"buyAmountRatio"`
	SellAmount      *float64 `json:"sellAmount"`
	SellAmountRatio *float64 `json:"sellAmountRatio"`
	NetAmount       *float64 `json:"netAmount"`
	Side            string   `json:"side"`
}
