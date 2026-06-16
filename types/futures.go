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
	Date          string   `json:"date"`
	Code          string   `json:"code"`
	Name          string   `json:"name"`
	Open          *float64 `json:"open"`
	Close         *float64 `json:"close"`
	High          *float64 `json:"high"`
	Low           *float64 `json:"low"`
	Volume        *float64 `json:"volume"`
	Amount        *float64 `json:"amount"`
	Amplitude     *float64 `json:"amplitude"`
	ChangePercent *float64 `json:"changePercent"`
	Change        *float64 `json:"change"`
	TurnoverRate  *float64 `json:"turnoverRate"`
	OpenInterest  *float64 `json:"openInterest"`
}

// GlobalFuturesQuote is a global futures spot quote row.
type GlobalFuturesQuote struct {
	Code          string   `json:"code"`
	Name          string   `json:"name"`
	Price         *float64 `json:"price"`
	Change        *float64 `json:"change"`
	ChangePercent *float64 `json:"changePercent"`
	Open          *float64 `json:"open"`
	High          *float64 `json:"high"`
	Low           *float64 `json:"low"`
	PrevSettle    *float64 `json:"prevSettle"`
	Volume        *float64 `json:"volume"`
	BuyVolume     *float64 `json:"buyVolume"`
	SellVolume    *float64 `json:"sellVolume"`
	OpenInterest  *float64 `json:"openInterest"`
}

// FuturesInventorySymbol 是期货库存品种信息。
type FuturesInventorySymbol struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	MarketCode string `json:"marketCode"`
}

// FuturesInventory 是国内期货库存数据行。
type FuturesInventory struct {
	Code      string   `json:"code"`
	Date      string   `json:"date"`
	Inventory *float64 `json:"inventory"`
	Change    *float64 `json:"change"`
}

// ComexInventory 是 COMEX 黄金或白银库存数据行。
type ComexInventory struct {
	Date         string   `json:"date"`
	Name         string   `json:"name"`
	StorageTon   *float64 `json:"storageTon"`
	StorageOunce *float64 `json:"storageOunce"`
}
