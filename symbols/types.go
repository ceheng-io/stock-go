package symbols

// Market is a trading region or market system.
type Market string

const (
	MarketCN     Market = "CN"
	MarketHK     Market = "HK"
	MarketUS     Market = "US"
	MarketGlobal Market = "GLOBAL"
)

// AssetType describes the normalized asset category.
type AssetType string

const (
	AssetStock   AssetType = "stock"
	AssetIndex   AssetType = "index"
	AssetFund    AssetType = "fund"
	AssetBond    AssetType = "bond"
	AssetFutures AssetType = "futures"
	AssetOption  AssetType = "option"
	AssetBoard   AssetType = "board"
)

// Exchange is a trading exchange code.
type Exchange string

const (
	ExchangeSSE    Exchange = "SSE"
	ExchangeSZSE   Exchange = "SZSE"
	ExchangeBSE    Exchange = "BSE"
	ExchangeHKEX   Exchange = "HKEX"
	ExchangeNASDAQ Exchange = "NASDAQ"
	ExchangeNYSE   Exchange = "NYSE"
	ExchangeAMEX   Exchange = "AMEX"
	ExchangeUS     Exchange = "US"
	ExchangeSHFE   Exchange = "SHFE"
	ExchangeDCE    Exchange = "DCE"
	ExchangeCZCE   Exchange = "CZCE"
	ExchangeINE    Exchange = "INE"
	ExchangeCFFEX  Exchange = "CFFEX"
	ExchangeGFEX   Exchange = "GFEX"
	ExchangeCOMEX  Exchange = "COMEX"
	ExchangeNYMEX  Exchange = "NYMEX"
	ExchangeCBOT   Exchange = "CBOT"
	ExchangeLME    Exchange = "LME"
)

// SymbolRef is a user input with optional market hints.
type SymbolRef struct {
	Code      string
	Market    Market
	AssetType AssetType
	Exchange  Exchange
}

// SymbolInput preserves the TypeScript SDK input type name.
//
// Normalize accepts string values and SymbolRef values; keeping this as an
// alias of any avoids narrowing the existing Go API while retaining the
// broad public input name.
type SymbolInput = any

// Hint provides optional parse hints for string inputs.
type Hint struct {
	Market    Market
	AssetType AssetType
	Exchange  Exchange
}

// Normalized is the SDK's internal symbol representation.
type Normalized struct {
	Market    Market
	AssetType AssetType
	Exchange  Exchange
	Code      string
	Variety   string
	Input     string
}
