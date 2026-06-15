package symbols

import (
	"fmt"

	"github.com/ceheng.io/stock-go/internal/core"
)

var exchangeToSecidPrefix = map[Exchange]string{
	ExchangeSSE:    "1",
	ExchangeSZSE:   "0",
	ExchangeBSE:    "0",
	ExchangeHKEX:   "116",
	ExchangeNASDAQ: "105",
	ExchangeNYSE:   "106",
	ExchangeAMEX:   "107",
	ExchangeUS:     "105",
}

var exchangeToTencentPrefix = map[Exchange]string{
	ExchangeSSE:  "sh",
	ExchangeSZSE: "sz",
	ExchangeBSE:  "bj",
}

// ToTencent maps a normalized symbol to Tencent quote format.
func ToTencent(symbol Normalized) string {
	value, _ := ToTencentE(symbol)
	return value
}

// ToTencentSymbol preserves the TypeScript SDK function name.
func ToTencentSymbol(symbol Normalized) string {
	return ToTencent(symbol)
}

// ToTencentE maps a normalized symbol to Tencent quote format with validation.
func ToTencentE(symbol Normalized) (string, error) {
	switch symbol.AssetType {
	case AssetBoard, AssetFutures, AssetOption:
		return "", invalidArgumentError(fmt.Sprintf("tencent quote symbol does not support asset type %q", symbol.AssetType))
	}

	switch symbol.Market {
	case MarketCN:
		prefix, ok := exchangeToTencentPrefix[symbol.Exchange]
		if !ok {
			return "", invalidArgumentError(fmt.Sprintf("unsupported CN exchange %q", symbol.Exchange))
		}
		return prefix + symbol.Code, nil
	case MarketHK:
		return "hk" + leftPad(symbol.Code, 5, "0"), nil
	case MarketUS:
		return "us" + symbol.Code, nil
	default:
		return "", invalidArgumentError(fmt.Sprintf("unsupported market %q", symbol.Market))
	}
}

// ToEastmoneySecID maps a normalized symbol to Eastmoney secid format.
func ToEastmoneySecID(symbol Normalized) string {
	value, _ := ToEastmoneySecIDE(symbol)
	return value
}

// ToEastmoneySecid preserves the TypeScript SDK function name.
func ToEastmoneySecid(symbol Normalized) string {
	return ToEastmoneySecID(symbol)
}

// ToEastmoneySecIDE maps a normalized symbol to Eastmoney secid format with validation.
func ToEastmoneySecIDE(symbol Normalized) (string, error) {
	if symbol.AssetType == AssetBoard {
		return "90." + symbol.Code, nil
	}
	if symbol.AssetType == AssetFutures || symbol.AssetType == AssetOption {
		return "", invalidArgumentError(fmt.Sprintf("eastmoney secid for asset type %q uses a separate scheme", symbol.AssetType))
	}
	if symbol.Market == MarketHK {
		return "116." + leftPad(symbol.Code, 5, "0"), nil
	}
	prefix, ok := exchangeToSecidPrefix[symbol.Exchange]
	if !ok {
		return "", invalidArgumentError(fmt.Sprintf("unsupported exchange %q", symbol.Exchange))
	}
	return prefix + "." + symbol.Code, nil
}

func invalidArgumentError(message string) error {
	return core.NewCodedError("INVALID_ARGUMENT", message, nil)
}

// ToPlainCode returns the plain normalized code.
func ToPlainCode(symbol Normalized) string {
	return symbol.Code
}
