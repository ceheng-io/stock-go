package symbols

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ceheng-io/stock-go/internal/core"
)

var (
	numberPattern       = regexp.MustCompile(`^\d+$`)
	alnumPattern        = regexp.MustCompile(`^[0-9A-Za-z]+$`)
	letterTickerPattern = regexp.MustCompile(`^[A-Za-z][A-Za-z.\-]*$`)
	containsLetter      = regexp.MustCompile(`[A-Za-z]`)
)

var prefixMap = map[string]struct {
	market   Market
	exchange Exchange
}{
	"sh": {market: MarketCN, exchange: ExchangeSSE},
	"sz": {market: MarketCN, exchange: ExchangeSZSE},
	"bj": {market: MarketCN, exchange: ExchangeBSE},
	"hk": {market: MarketHK, exchange: ExchangeHKEX},
	"us": {market: MarketUS, exchange: ExchangeUS},
}

var suffixMap = map[string]struct {
	market   Market
	exchange Exchange
}{
	"SH": {market: MarketCN, exchange: ExchangeSSE},
	"SZ": {market: MarketCN, exchange: ExchangeSZSE},
	"BJ": {market: MarketCN, exchange: ExchangeBSE},
	"HK": {market: MarketHK, exchange: ExchangeHKEX},
	"US": {market: MarketUS, exchange: ExchangeUS},
}

var secidMap = map[string]struct {
	market   Market
	exchange Exchange
}{
	"0":   {market: MarketCN, exchange: ExchangeSZSE},
	"1":   {market: MarketCN, exchange: ExchangeSSE},
	"116": {market: MarketHK, exchange: ExchangeHKEX},
	"105": {market: MarketUS, exchange: ExchangeNASDAQ},
	"106": {market: MarketUS, exchange: ExchangeNYSE},
	"107": {market: MarketUS, exchange: ExchangeAMEX},
}

var orderedPrefixes = []string{"sh", "sz", "bj", "hk", "us"}

// Normalize parses a string or SymbolRef into a Normalized symbol.
func Normalize(input any, hint *Hint) (Normalized, error) {
	ref, rawInput, err := normalizeInput(input)
	if err != nil {
		return Normalized{}, err
	}
	code0 := strings.TrimSpace(ref.Code)
	if code0 == "" {
		return Normalized{}, invalidSymbolError(rawInput)
	}

	hintMarket := firstMarket(ref.Market, valueOrZero(hint, func(h *Hint) Market { return h.Market }))
	hintAsset := firstAsset(ref.AssetType, valueOrZero(hint, func(h *Hint) AssetType { return h.AssetType }))
	hintExchange := firstExchange(ref.Exchange, valueOrZero(hint, func(h *Hint) Exchange { return h.Exchange }))

	finish := func(market Market, exchange Exchange, code string, assetType AssetType, variety string) Normalized {
		if hintExchange != "" {
			exchange = hintExchange
		}
		if hintAsset != "" {
			assetType = hintAsset
		}
		return Normalized{
			Market:    market,
			Exchange:  exchange,
			AssetType: assetType,
			Code:      code,
			Variety:   variety,
			Input:     rawInput,
		}
	}

	if strings.Contains(code0, ".") {
		dot := strings.Index(code0, ".")
		left := code0[:dot]
		right := code0[dot+1:]
		upperLeft := strings.ToUpper(left)
		upperRight := strings.ToUpper(right)

		if numberPattern.MatchString(left) {
			if secid, ok := secidMap[left]; ok {
				return finish(secid.market, secid.exchange, right, AssetStock, ""), nil
			}
		}
		if suffix, ok := suffixMap[upperRight]; ok {
			cleanLeft := left
			lowerLeft := strings.ToLower(left)
			for _, prefix := range orderedPrefixes {
				if strings.HasPrefix(lowerLeft, prefix) && len(left) > len(prefix) {
					cleanLeft = left[len(prefix):]
					break
				}
			}
			return finish(suffix.market, suffix.exchange, cleanLeft, AssetStock, ""), nil
		}
		if futures, ok := futuresExchanges[upperLeft]; ok {
			return finish(futures.Market, futures.Exchange, upperRight, AssetFutures, ExtractVariety(right)), nil
		}
		if left == "90" {
			return finish(MarketCN, ExchangeSSE, right, AssetBoard, ""), nil
		}
	}

	lower := strings.ToLower(code0)
	for _, prefix := range orderedPrefixes {
		if strings.HasPrefix(lower, prefix) && len(code0) > len(prefix) {
			rest := code0[len(prefix):]
			mapped := prefixMap[prefix]
			restOK := false
			if mapped.market == MarketCN {
				restOK = numberPattern.MatchString(rest)
			} else {
				restOK = alnumPattern.MatchString(rest)
			}
			if restOK {
				code := rest
				if mapped.market == MarketHK {
					code = leftPad(rest, 5, "0")
				}
				if mapped.market == MarketUS {
					code = strings.ToUpper(rest)
				}
				return finish(mapped.market, mapped.exchange, code, AssetStock, ""), nil
			}
		}
	}

	if numberPattern.MatchString(code0) {
		if hintMarket == MarketUS {
			return finish(MarketUS, ExchangeUS, code0, AssetStock, ""), nil
		}
		if hintMarket == MarketHK || len(code0) == 5 || len(code0) == 4 {
			return finish(MarketHK, ExchangeHKEX, leftPad(code0, 5, "0"), AssetStock, ""), nil
		}
		return finish(MarketCN, InferAShareExchange(code0), code0, AssetStock, ""), nil
	}

	if (hintAsset == AssetFutures || hintMarket == MarketGlobal) && containsLetter.MatchString(code0) {
		futExchange := hintExchange
		if futExchange == "" && hintMarket == MarketGlobal {
			return Normalized{}, invalidSymbolError(fmt.Sprintf("%s (GLOBAL futures require an explicit exchange, e.g. { exchange: 'COMEX' })", rawInput))
		}
		if futExchange == "" {
			futExchange = ExchangeSHFE
		}
		market := hintMarket
		if market == "" {
			market = MarketCN
		}
		return finish(market, futExchange, strings.ToUpper(code0), AssetFutures, ExtractVariety(code0)), nil
	}

	if letterTickerPattern.MatchString(code0) {
		return finish(MarketUS, ExchangeUS, strings.ToUpper(code0), AssetStock, ""), nil
	}

	return Normalized{}, invalidSymbolError(rawInput)
}

// NormalizeSymbol preserves the TypeScript SDK function name.
func NormalizeSymbol(input SymbolInput, hint *Hint) (Normalized, error) {
	return Normalize(input, hint)
}

func normalizeInput(input any) (SymbolRef, string, error) {
	switch value := input.(type) {
	case string:
		return SymbolRef{Code: value}, value, nil
	case SymbolRef:
		return value, value.Code, nil
	case *SymbolRef:
		if value == nil {
			return SymbolRef{}, "", invalidSymbolError("<nil>")
		}
		return *value, value.Code, nil
	default:
		return SymbolRef{}, "", invalidSymbolError(fmt.Sprintf("%v", input))
	}
}

func invalidSymbolError(symbol string) error {
	return core.NewCodedError("INVALID_SYMBOL", fmt.Sprintf("Invalid symbol: %s", symbol), nil)
}

func valueOrZero[T any](hint *Hint, fn func(*Hint) T) T {
	var zero T
	if hint == nil {
		return zero
	}
	return fn(hint)
}

func firstMarket(values ...Market) Market {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func firstAsset(values ...AssetType) AssetType {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func firstExchange(values ...Exchange) Exchange {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func leftPad(value string, width int, pad string) string {
	for len(value) < width {
		value = pad + value
	}
	return value
}
