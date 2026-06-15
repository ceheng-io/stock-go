package types

// FuturesExchange is a domestic futures exchange code.
type FuturesExchange string

const (
	FuturesExchangeSHFE  FuturesExchange = "SHFE"
	FuturesExchangeDCE   FuturesExchange = "DCE"
	FuturesExchangeCZCE  FuturesExchange = "CZCE"
	FuturesExchangeINE   FuturesExchange = "INE"
	FuturesExchangeCFFEX FuturesExchange = "CFFEX"
	FuturesExchangeGFEX  FuturesExchange = "GFEX"
)

// FuturesKline is a domestic or global futures historical K-line row.
type FuturesKline struct {
	Date          string
	Code          string
	Name          string
	Open          *float64
	Close         *float64
	High          *float64
	Low           *float64
	Volume        *float64
	Amount        *float64
	Amplitude     *float64
	ChangePercent *float64
	Change        *float64
	TurnoverRate  *float64
	OpenInterest  *float64
}

// GlobalFuturesQuote is a global futures spot quote row.
type GlobalFuturesQuote struct {
	Code          string
	Name          string
	Price         *float64
	Change        *float64
	ChangePercent *float64
	Open          *float64
	High          *float64
	Low           *float64
	PrevSettle    *float64
	Volume        *float64
	BuyVolume     *float64
	SellVolume    *float64
	OpenInterest  *float64
}

// FuturesInventorySymbol 是期货库存品种信息。
type FuturesInventorySymbol struct {
	Code       string
	Name       string
	MarketCode string
}

// FuturesInventory 是国内期货库存数据行。
type FuturesInventory struct {
	Code      string
	Date      string
	Inventory *float64
	Change    *float64
}

// ComexInventory 是 COMEX 黄金或白银库存数据行。
type ComexInventory struct {
	Date         string
	Name         string
	StorageTon   *float64
	StorageOunce *float64
}
