package stock

import (
	"fmt"
	"regexp"
	"strings"
)

var futuresVarietyPattern = regexp.MustCompile(`^([A-Za-z]+)`)

// ExtractFuturesVariety extracts the futures variety prefix from a contract code.
func ExtractFuturesVariety(symbol string) (string, error) {
	match := futuresVarietyPattern.FindStringSubmatch(strings.TrimSpace(symbol))
	if len(match) < 2 {
		return "", NewInvalidArgumentError(
			fmt.Sprintf("invalid futures symbol %q: expected variety + contract, e.g. rb2605, RBM, IF2604", symbol),
			map[string]any{"symbol": symbol},
		)
	}
	return match[1], nil
}

// FuturesMarketCode returns Eastmoney market code for a domestic futures variety.
func FuturesMarketCode(variety string) (int, error) {
	normalized := strings.TrimSpace(variety)
	exchange, ok := lookupFuturesExchange(normalized)
	if !ok && len(normalized) > 1 && strings.HasSuffix(normalized, "M") {
		exchange, ok = lookupFuturesExchange(normalized[:len(normalized)-1])
	}
	if !ok {
		return 0, NewInvalidArgumentError(
			fmt.Sprintf("unknown futures variety %q", variety),
			map[string]any{"variety": variety},
		)
	}
	marketCode, ok := FuturesExchangeMap()[exchange]
	if !ok {
		return 0, NewInvalidArgumentError(
			fmt.Sprintf("no market code found for futures exchange %q", exchange),
			map[string]any{"exchange": exchange, "variety": variety},
		)
	}
	return marketCode, nil
}

func lookupFuturesExchange(variety string) (string, bool) {
	varieties := FuturesVarietyExchange()
	if exchange, ok := varieties[variety]; ok {
		return exchange, true
	}
	if exchange, ok := varieties[strings.ToLower(variety)]; ok {
		return exchange, true
	}
	if exchange, ok := varieties[strings.ToUpper(variety)]; ok {
		return exchange, true
	}
	return "", false
}
