package types

// StockFundFlow is an Eastmoney stock or sector historical fund-flow row.
type StockFundFlow struct {
	Date                       string
	Close                      *float64
	ChangePercent              *float64
	MainNetInflow              *float64
	MainNetInflowPercent       *float64
	SuperLargeNetInflow        *float64
	SuperLargeNetInflowPercent *float64
	LargeNetInflow             *float64
	LargeNetInflowPercent      *float64
	MediumNetInflow            *float64
	MediumNetInflowPercent     *float64
	SmallNetInflow             *float64
	SmallNetInflowPercent      *float64
}

type StockFundFlowDaily = StockFundFlow

// MarketFundFlow is an Eastmoney market fund-flow row.
type MarketFundFlow struct {
	Date                       string
	SHClose                    *float64
	SHChangePercent            *float64
	SZClose                    *float64
	SZChangePercent            *float64
	MainNetInflow              *float64
	MainNetInflowPercent       *float64
	SuperLargeNetInflow        *float64
	SuperLargeNetInflowPercent *float64
	LargeNetInflow             *float64
	LargeNetInflowPercent      *float64
	MediumNetInflow            *float64
	MediumNetInflowPercent     *float64
	SmallNetInflow             *float64
	SmallNetInflowPercent      *float64
}

// FundFlowRankItem is an Eastmoney stock fund-flow ranking row.
type FundFlowRankItem struct {
	Code                       string
	Name                       string
	Price                      *float64
	ChangePercent              *float64
	MainNetInflow              *float64
	MainNetInflowPercent       *float64
	SuperLargeNetInflow        *float64
	SuperLargeNetInflowPercent *float64
	LargeNetInflow             *float64
	LargeNetInflowPercent      *float64
	MediumNetInflow            *float64
	MediumNetInflowPercent     *float64
	SmallNetInflow             *float64
	SmallNetInflowPercent      *float64
}

// SectorFundFlowItem is an Eastmoney sector fund-flow ranking row.
type SectorFundFlowItem struct {
	Code                 string
	Name                 string
	ChangePercent        *float64
	MainNetInflow        *float64
	MainNetInflowPercent *float64
	SuperLargeNetInflow  *float64
	LargeNetInflow       *float64
	MediumNetInflow      *float64
	SmallNetInflow       *float64
	TopStockName         *string
	TopStockCode         *string
}
