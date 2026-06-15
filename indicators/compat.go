package indicators

// IndicatorOptions keeps the TypeScript SDK indicator option name available.
type IndicatorOptions = Options

// IndicatorContext keeps the TypeScript SDK indicator context name available.
type IndicatorContext = Context

// IndicatorLookback keeps the TypeScript SDK lookback result name available.
type IndicatorLookback = Lookback

// IndicatorDescriptor describes one registered indicator.
type IndicatorDescriptor struct {
	Key              IndicatorKey
	EstimateLookback func(IndicatorOptions) IndicatorLookback
}

// BuildIndicatorContext extracts reusable indicator inputs from K-line rows.
func BuildIndicatorContext(rows []KlineInput) IndicatorContext {
	return BuildContext(rows)
}

// GetEnabledIndicatorKeys returns enabled indicator keys in registry order.
func GetEnabledIndicatorKeys(options IndicatorOptions) []IndicatorKey {
	return GetEnabledKeys(options)
}

// EstimateIndicatorLookback returns the amount of historical context needed.
func EstimateIndicatorLookback(options IndicatorOptions) IndicatorLookback {
	return EstimateLookback(options)
}

// IndicatorRegistry returns a copy of the indicator descriptor registry.
func IndicatorRegistry() map[IndicatorKey]IndicatorDescriptor {
	registry := make(map[IndicatorKey]IndicatorDescriptor, len(indicatorKeyOrder))
	for _, key := range indicatorKeyOrder {
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
