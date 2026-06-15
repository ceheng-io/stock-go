package symbols

import (
	"regexp"
	"strings"
)

// FuturesExchange maps futures exchange codes to market metadata.
type FuturesExchange struct {
	Market   Market
	Exchange Exchange
}

var futuresExchanges = map[string]FuturesExchange{
	"SHFE":  {Market: MarketCN, Exchange: ExchangeSHFE},
	"DCE":   {Market: MarketCN, Exchange: ExchangeDCE},
	"CZCE":  {Market: MarketCN, Exchange: ExchangeCZCE},
	"INE":   {Market: MarketCN, Exchange: ExchangeINE},
	"CFFEX": {Market: MarketCN, Exchange: ExchangeCFFEX},
	"GFEX":  {Market: MarketCN, Exchange: ExchangeGFEX},
	"COMEX": {Market: MarketGlobal, Exchange: ExchangeCOMEX},
	"NYMEX": {Market: MarketGlobal, Exchange: ExchangeNYMEX},
	"CBOT":  {Market: MarketGlobal, Exchange: ExchangeCBOT},
	"LME":   {Market: MarketGlobal, Exchange: ExchangeLME},
}

// FUTURES_EXCHANGES keeps the TypeScript SDK futures exchange map name available.
var FUTURES_EXCHANGES = FuturesExchanges()

var varietyPattern = regexp.MustCompile(`^[A-Za-z]+`)

// FuturesExchanges returns a copy of known futures exchange metadata.
func FuturesExchanges() map[string]FuturesExchange {
	exchanges := make(map[string]FuturesExchange, len(futuresExchanges))
	for code, exchange := range futuresExchanges {
		exchanges[code] = exchange
	}
	return exchanges
}

// ExtractVariety extracts a futures variety from a contract code.
func ExtractVariety(contract string) string {
	match := varietyPattern.FindString(contract)
	if match == "" {
		return strings.ToUpper(contract)
	}
	return strings.ToUpper(match)
}
