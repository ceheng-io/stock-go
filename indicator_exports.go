package stock

import "github.com/ceheng.io/stock-go/indicators"

type IndicatorValue = indicators.Value
type OHLCV = indicators.OHLCV

type MAType = indicators.MAType

const (
	MATypeSMA = indicators.MATypeSMA
	MATypeEMA = indicators.MATypeEMA
	MATypeWMA = indicators.MATypeWMA
)

type MAOptions = indicators.MAOptions
type MAResult = indicators.MAResult
type MACDOptions = indicators.MACDOptions
type MACDResult = indicators.MACDResult
type BOLLOptions = indicators.BOLLOptions
type BOLLResult = indicators.BOLLResult
type KDJOptions = indicators.KDJOptions
type KDJResult = indicators.KDJResult
type RSIOptions = indicators.RSIOptions
type RSIResult = indicators.RSIResult
type WROptions = indicators.WROptions
type WRResult = indicators.WRResult
type BIASOptions = indicators.BIASOptions
type BIASResult = indicators.BIASResult
type CCIOptions = indicators.CCIOptions
type CCIResult = indicators.CCIResult
type ATROptions = indicators.ATROptions
type ATRResult = indicators.ATRResult
type OBVOptions = indicators.OBVOptions
type OBVResult = indicators.OBVResult
type ROCOptions = indicators.ROCOptions
type ROCResult = indicators.ROCResult
type DMIOptions = indicators.DMIOptions
type DMIResult = indicators.DMIResult
type SAROptions = indicators.SAROptions
type SARResult = indicators.SARResult
type KCOptions = indicators.KCOptions
type KCResult = indicators.KCResult
type IndicatorOptions = indicators.Options
type KlineInput = indicators.KlineInput
type KlineWithIndicators = indicators.KlineWithIndicators
type IndicatorLookback = indicators.Lookback
type IndicatorContext = indicators.Context
type IndicatorKey = indicators.IndicatorKey

// IndicatorDescriptor describes one registered indicator.
type IndicatorDescriptor struct {
	Key              IndicatorKey
	EstimateLookback func(IndicatorOptions) IndicatorLookback
}

const (
	IndicatorMA   = indicators.IndicatorMA
	IndicatorMACD = indicators.IndicatorMACD
	IndicatorBOLL = indicators.IndicatorBOLL
	IndicatorKDJ  = indicators.IndicatorKDJ
	IndicatorRSI  = indicators.IndicatorRSI
	IndicatorWR   = indicators.IndicatorWR
	IndicatorBIAS = indicators.IndicatorBIAS
	IndicatorCCI  = indicators.IndicatorCCI
	IndicatorATR  = indicators.IndicatorATR
	IndicatorOBV  = indicators.IndicatorOBV
	IndicatorROC  = indicators.IndicatorROC
	IndicatorDMI  = indicators.IndicatorDMI
	IndicatorSAR  = indicators.IndicatorSAR
	IndicatorKC   = indicators.IndicatorKC
)

func Float(value float64) IndicatorValue {
	return indicators.Float(value)
}

func Null() IndicatorValue {
	return indicators.Null()
}

func Values(values ...float64) []IndicatorValue {
	return indicators.Values(values...)
}

func CalcSMA(data []IndicatorValue, period int) []IndicatorValue {
	return indicators.CalcSMA(data, period)
}

func CalcEMA(data []IndicatorValue, period int) []IndicatorValue {
	return indicators.CalcEMA(data, period)
}

func CalcWMA(data []IndicatorValue, period int) []IndicatorValue {
	return indicators.CalcWMA(data, period)
}

func CalcMA(closes []IndicatorValue, options MAOptions) []MAResult {
	return indicators.CalcMA(closes, options)
}

func CalcMACD(closes []IndicatorValue, options MACDOptions) []MACDResult {
	return indicators.CalcMACD(closes, options)
}

func CalcBOLL(closes []IndicatorValue, options BOLLOptions) []BOLLResult {
	return indicators.CalcBOLL(closes, options)
}

func CalcKDJ(data []OHLCV, options KDJOptions) []KDJResult {
	return indicators.CalcKDJ(data, options)
}

func CalcRSI(closes []IndicatorValue, options RSIOptions) []RSIResult {
	return indicators.CalcRSI(closes, options)
}

func CalcWR(data []OHLCV, options WROptions) []WRResult {
	return indicators.CalcWR(data, options)
}

func CalcBIAS(closes []IndicatorValue, options BIASOptions) []BIASResult {
	return indicators.CalcBIAS(closes, options)
}

func CalcCCI(data []OHLCV, options CCIOptions) []CCIResult {
	return indicators.CalcCCI(data, options)
}

func CalcATR(data []OHLCV, options ATROptions) []ATRResult {
	return indicators.CalcATR(data, options)
}

func CalcOBV(data []OHLCV, options OBVOptions) []OBVResult {
	return indicators.CalcOBV(data, options)
}

func CalcROC(data []OHLCV, options ROCOptions) []ROCResult {
	return indicators.CalcROC(data, options)
}

func CalcDMI(data []OHLCV, options DMIOptions) []DMIResult {
	return indicators.CalcDMI(data, options)
}

func CalcSAR(data []OHLCV, options SAROptions) []SARResult {
	return indicators.CalcSAR(data, options)
}

func CalcKC(data []OHLCV, options KCOptions) []KCResult {
	return indicators.CalcKC(data, options)
}

func AddIndicators(rows []KlineInput, options IndicatorOptions) []KlineWithIndicators {
	return indicators.AddIndicators(rows, options)
}

func BuildIndicatorContext(rows []KlineInput) IndicatorContext {
	return indicators.BuildContext(rows)
}

func GetEnabledIndicatorKeys(options IndicatorOptions) []IndicatorKey {
	return indicators.GetEnabledKeys(options)
}

func EstimateIndicatorLookback(options IndicatorOptions) IndicatorLookback {
	return indicators.EstimateLookback(options)
}

// IndicatorRegistry returns a copy of the indicator descriptor registry.
func IndicatorRegistry() map[IndicatorKey]IndicatorDescriptor {
	registry := make(map[IndicatorKey]IndicatorDescriptor, len(indicatorRegistryKeys))
	for _, key := range indicatorRegistryKeys {
		key := key
		registry[key] = IndicatorDescriptor{
			Key: key,
			EstimateLookback: func(options IndicatorOptions) IndicatorLookback {
				return EstimateIndicatorLookback(optionForIndicator(key, options))
			},
		}
	}
	return registry
}

var INDICATOR_REGISTRY = IndicatorRegistry()

var indicatorRegistryKeys = []IndicatorKey{
	IndicatorMA,
	IndicatorMACD,
	IndicatorBOLL,
	IndicatorKDJ,
	IndicatorRSI,
	IndicatorWR,
	IndicatorBIAS,
	IndicatorCCI,
	IndicatorATR,
	IndicatorOBV,
	IndicatorROC,
	IndicatorDMI,
	IndicatorSAR,
	IndicatorKC,
}

func optionForIndicator(key IndicatorKey, options IndicatorOptions) IndicatorOptions {
	switch key {
	case IndicatorMA:
		return IndicatorOptions{MA: options.MA}
	case IndicatorMACD:
		return IndicatorOptions{MACD: options.MACD}
	case IndicatorBOLL:
		return IndicatorOptions{BOLL: options.BOLL}
	case IndicatorKDJ:
		return IndicatorOptions{KDJ: options.KDJ}
	case IndicatorRSI:
		return IndicatorOptions{RSI: options.RSI}
	case IndicatorWR:
		return IndicatorOptions{WR: options.WR}
	case IndicatorBIAS:
		return IndicatorOptions{BIAS: options.BIAS}
	case IndicatorCCI:
		return IndicatorOptions{CCI: options.CCI}
	case IndicatorATR:
		return IndicatorOptions{ATR: options.ATR}
	case IndicatorOBV:
		return IndicatorOptions{OBV: options.OBV}
	case IndicatorROC:
		return IndicatorOptions{ROC: options.ROC}
	case IndicatorDMI:
		return IndicatorOptions{DMI: options.DMI}
	case IndicatorSAR:
		return IndicatorOptions{SAR: options.SAR}
	case IndicatorKC:
		return IndicatorOptions{KC: options.KC}
	default:
		return IndicatorOptions{}
	}
}
