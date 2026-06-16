package types

// StockFundFlow is an Eastmoney stock or sector historical fund-flow row.
type StockFundFlow struct {
	Date                       string   `json:"date"`
	Close                      *float64 `json:"close"`
	ChangePercent              *float64 `json:"changePercent"`
	MainNetInflow              *float64 `json:"mainNetInflow"`
	MainNetInflowPercent       *float64 `json:"mainNetInflowPercent"`
	SuperLargeNetInflow        *float64 `json:"superLargeNetInflow"`
	SuperLargeNetInflowPercent *float64 `json:"superLargeNetInflowPercent"`
	LargeNetInflow             *float64 `json:"largeNetInflow"`
	LargeNetInflowPercent      *float64 `json:"largeNetInflowPercent"`
	MediumNetInflow            *float64 `json:"mediumNetInflow"`
	MediumNetInflowPercent     *float64 `json:"mediumNetInflowPercent"`
	SmallNetInflow             *float64 `json:"smallNetInflow"`
	SmallNetInflowPercent      *float64 `json:"smallNetInflowPercent"`
}

type StockFundFlowDaily = StockFundFlow

// MarketFundFlow is an Eastmoney market fund-flow row.
type MarketFundFlow struct {
	Date                       string   `json:"date"`
	SHClose                    *float64 `json:"shClose"`
	SHChangePercent            *float64 `json:"shChangePercent"`
	SZClose                    *float64 `json:"szClose"`
	SZChangePercent            *float64 `json:"szChangePercent"`
	MainNetInflow              *float64 `json:"mainNetInflow"`
	MainNetInflowPercent       *float64 `json:"mainNetInflowPercent"`
	SuperLargeNetInflow        *float64 `json:"superLargeNetInflow"`
	SuperLargeNetInflowPercent *float64 `json:"superLargeNetInflowPercent"`
	LargeNetInflow             *float64 `json:"largeNetInflow"`
	LargeNetInflowPercent      *float64 `json:"largeNetInflowPercent"`
	MediumNetInflow            *float64 `json:"mediumNetInflow"`
	MediumNetInflowPercent     *float64 `json:"mediumNetInflowPercent"`
	SmallNetInflow             *float64 `json:"smallNetInflow"`
	SmallNetInflowPercent      *float64 `json:"smallNetInflowPercent"`
}

// FundFlowRankItem is an Eastmoney stock fund-flow ranking row.
type FundFlowRankItem struct {
	Code                       string   `json:"code"`
	Name                       string   `json:"name"`
	Price                      *float64 `json:"price"`
	ChangePercent              *float64 `json:"changePercent"`
	MainNetInflow              *float64 `json:"mainNetInflow"`
	MainNetInflowPercent       *float64 `json:"mainNetInflowPercent"`
	SuperLargeNetInflow        *float64 `json:"superLargeNetInflow"`
	SuperLargeNetInflowPercent *float64 `json:"superLargeNetInflowPercent"`
	LargeNetInflow             *float64 `json:"largeNetInflow"`
	LargeNetInflowPercent      *float64 `json:"largeNetInflowPercent"`
	MediumNetInflow            *float64 `json:"mediumNetInflow"`
	MediumNetInflowPercent     *float64 `json:"mediumNetInflowPercent"`
	SmallNetInflow             *float64 `json:"smallNetInflow"`
	SmallNetInflowPercent      *float64 `json:"smallNetInflowPercent"`
}

// SectorFundFlowItem is an Eastmoney sector fund-flow ranking row.
type SectorFundFlowItem struct {
	Code                 string   `json:"code"`
	Name                 string   `json:"name"`
	ChangePercent        *float64 `json:"changePercent"`
	MainNetInflow        *float64 `json:"mainNetInflow"`
	MainNetInflowPercent *float64 `json:"mainNetInflowPercent"`
	SuperLargeNetInflow  *float64 `json:"superLargeNetInflow"`
	LargeNetInflow       *float64 `json:"largeNetInflow"`
	MediumNetInflow      *float64 `json:"mediumNetInflow"`
	SmallNetInflow       *float64 `json:"smallNetInflow"`
	TopStockName         *string  `json:"topStockName"`
	TopStockCode         *string  `json:"topStockCode"`
}
