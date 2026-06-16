package types

// NorthboundDirection is northbound or southbound flow direction.
type NorthboundDirection string

const (
	NorthboundNorth NorthboundDirection = "north"
	NorthboundSouth NorthboundDirection = "south"
)

// NorthboundMarket is a northbound holding rank market scope.
type NorthboundMarket string

const (
	NorthboundMarketAll      NorthboundMarket = "all"
	NorthboundMarketShanghai NorthboundMarket = "shanghai"
	NorthboundMarketShenzhen NorthboundMarket = "shenzhen"
)

// NorthboundRankPeriod is a northbound holding rank period.
type NorthboundRankPeriod string

const (
	NorthboundRankToday    NorthboundRankPeriod = "today"
	NorthboundRankThreeDay NorthboundRankPeriod = "3day"
	NorthboundRankFiveDay  NorthboundRankPeriod = "5day"
	NorthboundRankTenDay   NorthboundRankPeriod = "10day"
	NorthboundRankMonth    NorthboundRankPeriod = "month"
	NorthboundRankQuarter  NorthboundRankPeriod = "quarter"
	NorthboundRankYear     NorthboundRankPeriod = "year"
)

// NorthboundMinuteItem is a northbound or southbound intraday fund-flow row.
type NorthboundMinuteItem struct {
	Date              string   `json:"date"`
	Time              string   `json:"time"`
	ShanghaiNetInflow *float64 `json:"shanghaiNetInflow"`
	ShenzhenNetInflow *float64 `json:"shenzhenNetInflow"`
	TotalNetInflow    *float64 `json:"totalNetInflow"`
}

// NorthboundFlowSummary is a Shanghai/Shenzhen/HK connect flow summary row.
type NorthboundFlowSummary struct {
	Date               string   `json:"date"`
	Type               string   `json:"type"`
	BoardName          string   `json:"boardName"`
	Direction          string   `json:"direction"`
	Status             string   `json:"status"`
	NetBuyAmount       *float64 `json:"netBuyAmount"`
	NetInflow          *float64 `json:"netInflow"`
	RemainAmount       *float64 `json:"remainAmount"`
	UpCount            *float64 `json:"upCount"`
	FlatCount          *float64 `json:"flatCount"`
	DownCount          *float64 `json:"downCount"`
	IndexCode          string   `json:"indexCode"`
	IndexName          string   `json:"indexName"`
	IndexChangePercent *float64 `json:"indexChangePercent"`
}

// NorthboundHoldingRankItem is a northbound holding ranking row.
type NorthboundHoldingRankItem struct {
	Date                  string   `json:"date"`
	Code                  string   `json:"code"`
	Name                  string   `json:"name"`
	Close                 *float64 `json:"close"`
	ChangePercent         *float64 `json:"changePercent"`
	HoldShares            *float64 `json:"holdShares"`
	HoldMarketValue       *float64 `json:"holdMarketValue"`
	HoldRatioFloat        *float64 `json:"holdRatioFloat"`
	HoldRatioTotal        *float64 `json:"holdRatioTotal"`
	AddShares             *float64 `json:"addShares"`
	AddMarketValue        *float64 `json:"addMarketValue"`
	AddMarketValuePercent *float64 `json:"addMarketValuePercent"`
	Sector                string   `json:"sector"`
}

// NorthboundHistoryItem is a northbound or southbound daily fund-flow row.
type NorthboundHistoryItem struct {
	Date                  string   `json:"date"`
	NetBuyAmount          *float64 `json:"netBuyAmount"`
	BuyAmount             *float64 `json:"buyAmount"`
	SellAmount            *float64 `json:"sellAmount"`
	AccNetBuyAmount       *float64 `json:"accNetBuyAmount"`
	NetInflow             *float64 `json:"netInflow"`
	RemainAmount          *float64 `json:"remainAmount"`
	TopStockCode          *string  `json:"topStockCode"`
	TopStockName          *string  `json:"topStockName"`
	TopStockChangePercent *float64 `json:"topStockChangePercent"`
}

// NorthboundIndividualItem is a stock's northbound holding history row.
type NorthboundIndividualItem struct {
	Date            string   `json:"date"`
	HoldShares      *float64 `json:"holdShares"`
	HoldMarketValue *float64 `json:"holdMarketValue"`
	HoldRatioFloat  *float64 `json:"holdRatioFloat"`
	HoldRatioTotal  *float64 `json:"holdRatioTotal"`
	Close           *float64 `json:"close"`
	ChangePercent   *float64 `json:"changePercent"`
}
