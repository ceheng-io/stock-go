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
	Date              string
	Time              string
	ShanghaiNetInflow *float64
	ShenzhenNetInflow *float64
	TotalNetInflow    *float64
}

// NorthboundFlowSummary is a Shanghai/Shenzhen/HK connect flow summary row.
type NorthboundFlowSummary struct {
	Date               string
	Type               string
	BoardName          string
	Direction          string
	Status             string
	NetBuyAmount       *float64
	NetInflow          *float64
	RemainAmount       *float64
	UpCount            *float64
	FlatCount          *float64
	DownCount          *float64
	IndexCode          string
	IndexName          string
	IndexChangePercent *float64
}

// NorthboundHoldingRankItem is a northbound holding ranking row.
type NorthboundHoldingRankItem struct {
	Date                  string
	Code                  string
	Name                  string
	Close                 *float64
	ChangePercent         *float64
	HoldShares            *float64
	HoldMarketValue       *float64
	HoldRatioFloat        *float64
	HoldRatioTotal        *float64
	AddShares             *float64
	AddMarketValue        *float64
	AddMarketValuePercent *float64
	Sector                string
}

// NorthboundHistoryItem is a northbound or southbound daily fund-flow row.
type NorthboundHistoryItem struct {
	Date                  string
	NetBuyAmount          *float64
	BuyAmount             *float64
	SellAmount            *float64
	AccNetBuyAmount       *float64
	NetInflow             *float64
	RemainAmount          *float64
	TopStockCode          *string
	TopStockName          *string
	TopStockChangePercent *float64
}

// NorthboundIndividualItem is a stock's northbound holding history row.
type NorthboundIndividualItem struct {
	Date            string
	HoldShares      *float64
	HoldMarketValue *float64
	HoldRatioFloat  *float64
	HoldRatioTotal  *float64
	Close           *float64
	ChangePercent   *float64
}
