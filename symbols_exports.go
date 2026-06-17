package stock

import "github.com/ceheng-io/stock-go/symbols"

type Market = symbols.Market

const (
	SymbolMarketCN     Market = symbols.MarketCN
	SymbolMarketHK     Market = symbols.MarketHK
	SymbolMarketUS     Market = symbols.MarketUS
	SymbolMarketGlobal Market = symbols.MarketGlobal
)

type AssetType = symbols.AssetType

const (
	AssetStock   AssetType = symbols.AssetStock
	AssetIndex   AssetType = symbols.AssetIndex
	AssetFund    AssetType = symbols.AssetFund
	AssetBond    AssetType = symbols.AssetBond
	AssetFutures AssetType = symbols.AssetFutures
	AssetOption  AssetType = symbols.AssetOption
	AssetBoard   AssetType = symbols.AssetBoard
)

type Exchange = symbols.Exchange

const (
	ExchangeSSE    Exchange = symbols.ExchangeSSE
	ExchangeSZSE   Exchange = symbols.ExchangeSZSE
	ExchangeBSE    Exchange = symbols.ExchangeBSE
	ExchangeHKEX   Exchange = symbols.ExchangeHKEX
	ExchangeNASDAQ Exchange = symbols.ExchangeNASDAQ
	ExchangeNYSE   Exchange = symbols.ExchangeNYSE
	ExchangeAMEX   Exchange = symbols.ExchangeAMEX
	ExchangeUS     Exchange = symbols.ExchangeUS
	ExchangeSHFE   Exchange = symbols.ExchangeSHFE
	ExchangeDCE    Exchange = symbols.ExchangeDCE
	ExchangeCZCE   Exchange = symbols.ExchangeCZCE
	ExchangeINE    Exchange = symbols.ExchangeINE
	ExchangeCFFEX  Exchange = symbols.ExchangeCFFEX
	ExchangeGFEX   Exchange = symbols.ExchangeGFEX
	ExchangeCOMEX  Exchange = symbols.ExchangeCOMEX
	ExchangeNYMEX  Exchange = symbols.ExchangeNYMEX
	ExchangeCBOT   Exchange = symbols.ExchangeCBOT
	ExchangeLME    Exchange = symbols.ExchangeLME
)

type SymbolRef = symbols.SymbolRef
type SymbolInput = symbols.SymbolInput
type SymbolHint = symbols.Hint
type NormalizedSymbol = symbols.Normalized
type FuturesExchange = symbols.FuturesExchange

// NormalizeSymbol parses a string or SymbolRef into a normalized symbol.
func NormalizeSymbol(input any, hint *SymbolHint) (NormalizedSymbol, error) {
	return symbols.Normalize(input, hint)
}

// ToTencentSymbol maps a normalized symbol to Tencent quote format.
func ToTencentSymbol(symbol NormalizedSymbol) string {
	return symbols.ToTencent(symbol)
}

// ToEastmoneySecID maps a normalized symbol to Eastmoney secid format.
func ToEastmoneySecID(symbol NormalizedSymbol) string {
	return symbols.ToEastmoneySecID(symbol)
}

// ToEastmoneySecid maps a normalized symbol to Eastmoney secid format.
func ToEastmoneySecid(symbol NormalizedSymbol) string {
	return ToEastmoneySecID(symbol)
}

// ToPlainCode returns the plain normalized code.
func ToPlainCode(symbol NormalizedSymbol) string {
	return symbols.ToPlainCode(symbol)
}

// InferAShareExchange infers an A-share exchange from a plain numeric code.
func InferAShareExchange(code string) Exchange {
	return symbols.InferAShareExchange(code)
}

// ExtractVariety extracts a futures variety from a contract code.
func ExtractVariety(contract string) string {
	return symbols.ExtractVariety(contract)
}

// FuturesExchanges returns a copy of known futures exchange metadata.
func FuturesExchanges() map[string]FuturesExchange {
	return symbols.FuturesExchanges()
}
